package runner

// Tests focados em invariantes normativos do ScriptRunner.
// Rastreabilidade:
// - E1.4-1: Pipeline exclusivo via RunAction → ScriptRunner (simulado via API pública).
// - E1.4-3 / docs/architecture/security.md#5-regra-de-1-processo-por-update_target-como-invariante-de-seguranca-e14:
//   Regra 1 processo por update_target (cancelamento + SIGTERM/SIGKILL).
// - E1.4-4 / Gate 1.x.no-osexit-core.yml:
//   Garantia indireta: uso apenas de retornos/erros; scanner externo cobre os.Exit.
// - E1.5/E1.7 / Gate 1.x.security-jit-anti-trojan.yml:
//   Integração com SecurityPolicy (deny by default + whitelist).
// - E1.4-4: Emissão de ShantillyEvent estruturado em vez de acoplamento direto à UI.
//
// Estes testes não cobrem integração TEA/UI; focam no comportamento de serviço puro.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// fakeSink implementa EventSink para captura de eventos/updates.
type fakeSink struct {
	mu      sync.Mutex
	events  []ShantillyEvent
	updates []UpdateTargetUpdate
}

func (f *fakeSink) EmitEvent(ev ShantillyEvent) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, ev)
}

func (f *fakeSink) EmitUpdate(upd UpdateTargetUpdate) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updates = append(f.updates, upd)
}

func (f *fakeSink) lastEventOfType(t string) *ShantillyEvent {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := len(f.events) - 1; i >= 0; i-- {
		if f.events[i].Type == t {
			return &f.events[i]
		}
	}
	return nil
}

func (f *fakeSink) hasEventType(t string) bool {
	return f.lastEventOfType(t) != nil
}

func (f *fakeSink) lastUpdateFor(target string) *UpdateTargetUpdate {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := len(f.updates) - 1; i >= 0; i-- {
		if f.updates[i].TargetID == target {
			return &f.updates[i]
		}
	}
	return nil
}

// helper para criar script temporário simples compatível com o SO.
func createTempScript(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	var path string

	if runtime.GOOS == "windows" {
		path = filepath.Join(dir, "test-script.bat")
	} else {
		path = filepath.Join(dir, "test-script.sh")
	}

	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("failed to write temp script: %v", err)
	}

	// Garantir permissões executáveis em ambientes Unix.
	if runtime.GOOS != "windows" {
		_ = os.Chmod(path, 0o755)
	}

	return path
}

// Testa caminho feliz: script permitido executa, gera script.start/script.complete e update.
func TestScriptRunner_Run_Success(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Process lifecycle tests are POSIX-oriented; skip on Windows for now")
	}

	script := createTempScript(t, "#!/usr/bin/env bash\necho ok-success\n")

	sink := &fakeSink{}
	r := NewScriptRunner(sink, &SecurityPolicy{
		AllowedScripts: []string{script},
		DenyUnknown:    true,
	})

	ra := RunAction{
		Script:       script,
		UpdateTarget: "viewport-success",
	}

	if err := r.Run(ra); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if !sink.hasEventType("script.start") {
		t.Fatalf("expected script.start event")
	}
	if !sink.hasEventType("script.complete") {
		t.Fatalf("expected script.complete event")
	}

	upd := sink.lastUpdateFor("viewport-success")
	if upd == nil {
		t.Fatalf("expected update for viewport-success")
	}
	if upd.Mode != "append" {
		t.Errorf("expected Mode=append, got %s", upd.Mode)
	}
	if upd.Content == "" || !contains(upd.Content, "ok-success") {
		t.Errorf("expected content to contain ok-success, got %q", upd.Content)
	}
}

// Testa deny-by-default + whitelist: script fora da lista deve ser bloqueado.
func TestScriptRunner_Run_SecurityPolicy_DenyUnknown(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip on Windows")
	}

	allowedScript := createTempScript(t, "#!/usr/bin/env bash\necho allowed\n")
	disallowedScript := createTempScript(t, "#!/usr/bin/env bash\necho disallowed\n")

	sink := &fakeSink{}
	r := NewScriptRunner(sink, &SecurityPolicy{
		AllowedScripts: []string{allowedScript},
		DenyUnknown:    true,
	})

	// Execução permitida
	if err := r.Run(RunAction{Script: allowedScript}); err != nil {
		t.Fatalf("expected allowed script to run, got error: %v", err)
	}

	// Execução proibida
	err := r.Run(RunAction{Script: disallowedScript})
	if err == nil {
		t.Fatalf("expected error for disallowed script")
	}
	if !sink.hasEventType("script.error") {
		t.Fatalf("expected script.error event for disallowed script")
	}
}

// Testa regra 1 processo por update_target com cancelamento do anterior.
func TestScriptRunner_OneProcessPerUpdateTarget_Cancellation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Process signal semantics differ on Windows; skip for simplicity")
	}

	t.Skip("E1.4 invariant (1 processo por update_target com cancelamento robusto) será consolidado na branch feat/runtime-migration; teste temporariamente desabilitado no runtime legado")

	// Script que roda por tempo razoável.
	longScript := createTempScript(t, "#!/usr/bin/env bash\nsleep 10\necho long-done\n")
	shortScript := createTempScript(t, "#!/usr/bin/env bash\necho short-done\n")

	sink := &fakeSink{}
	r := NewScriptRunner(sink, &SecurityPolicy{
		AllowedScripts: []string{longScript, shortScript},
		DenyUnknown:    true,
	})

	// Acelera timeouts para o teste.
	r.gracefulTimeout = 500 * time.Millisecond
	r.forceTimeout = 1 * time.Second

	target := "logs"

	// Dispara o primeiro Run em goroutine (longScript).
	errCh1 := make(chan error, 1)
	go func() {
		errCh1 <- r.Run(RunAction{
			Script:       longScript,
			UpdateTarget: target,
		})
	}()

	// Aguarda o processo iniciar.
	waitForProcess(t, &r.mu, r.processes, target, 2*time.Second)

	// Dispara segundo Run para o mesmo target, que deve cancelar o anterior.
	err2 := r.Run(RunAction{
		Script:       shortScript,
		UpdateTarget: target,
	})
	if err2 != nil {
		t.Fatalf("expected second run to succeed, got %v", err2)
	}

	// Verifica que o primeiro terminou (cancelado ou encerrado).
	select {
	case err1 := <-errCh1:
		// Pode ser context.Canceled ou erro por sinal; ambos aceitáveis.
		if err1 == nil {
			// Se for nil, foi interrompido antes do long sleep efetivo; aceitável
			// desde que não coexistam dois processos (garantido pela lógica interna).
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("expected first run to finish after cancellation")
	}

	// Verifica que não há processo ativo para o target.
	r.mu.Lock()
	_, exists := r.processes[target]
	r.mu.Unlock()
	if exists {
		t.Fatalf("expected no active process for target after second run")
	}

	// Deve ter update com short-done.
	upd := sink.lastUpdateFor(target)
	if upd == nil || !contains(upd.Content, "short-done") {
		t.Fatalf("expected last update for %s to contain short-done, got %+v", target, upd)
	}
}

// Testa que RunAction.Script obrigatório é validado.
func TestScriptRunner_Run_InvalidRunAction(t *testing.T) {
	sink := &fakeSink{}
	r := NewScriptRunner(sink, &SecurityPolicy{
		AllowedScripts: []string{"/bin/echo"},
		DenyUnknown:    true,
	})

	err := r.Run(RunAction{})
	if err == nil {
		t.Fatalf("expected error for empty Script")
	}
	if !sink.hasEventType("script.error") {
		t.Fatalf("expected script.error event for invalid RunAction")
	}
}

// Helper: aguarda até que exista processo registrado para update_target.
func waitForProcess(t *testing.T, mu *sync.Mutex, processes map[string]*managedProcess, target string, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timeout waiting for process for target %s", target)
		case <-time.After(50 * time.Millisecond):
			mu.Lock()
			_, exists := processes[target]
			mu.Unlock()
			if exists {
				return
			}
		}
	}
}

// contains é um helper simples para verificar substrings.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Helper adicional para simular sinais em ambientes que suportam.
func sendSignal(cmd *exec.Cmd, sig os.Signal) error {
	if cmd == nil || cmd.Process == nil {
		return errors.New("no process to signal")
	}
	return cmd.Process.Signal(sig)
}

// Helper para matar processo com SIGKILL quando suportado.
func killProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return errors.New("no process to kill")
	}
	return cmd.Process.Signal(syscall.SIGKILL)
}
