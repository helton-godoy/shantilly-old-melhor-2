package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// ... (structs e outras funções)

func (r *ScriptRunner) buildCommand(run RunAction) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	args := append([]string{run.Script}, run.Args...)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}
	// ... (resto da função)
	return cmd, ctx, cancel, nil
}

func (r *ScriptRunner) ensureSingleProcess(updateTarget string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.processes[updateTarget]
	if !ok || existing == nil || existing.cmd == nil || existing.cmd.Process == nil {
		return nil
	}

	if existing.cancel != nil {
		existing.cancel()
	}

	pid := existing.cmd.Process.Pid
	if runtime.GOOS != "windows" {
		// Envia sinal para todo o grupo de processos
		if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
			// Ignora o erro se o processo não existir mais
		}
	} else {
		_ = existing.cmd.Process.Signal(syscall.SIGTERM)
	}

	done := make(chan error, 1)
	go func(cmd *exec.Cmd) {
		done <- cmd.Wait()
	}(existing.cmd)

	select {
	case <-time.After(r.gracefulTimeout):
		if runtime.GOOS != "windows" {
			_ = syscall.Kill(-pid, syscall.SIGKILL)
		} else {
			_ = existing.cmd.Process.Signal(syscall.SIGKILL)
		}
	case <-done:
		return nil
	}

	select {
	case <-done:
	case <-time.After(r.forceTimeout):
		return fmt.Errorf("E1.4 violation: existing process for update_target '%s' did not terminate after SIGKILL", updateTarget)
	}

	return nil
}

// ... (resto do arquivo inalterado)
type ScriptRunner struct {
	mu sync.Mutex
	// processos ativos por update_target
	processes map[string]*managedProcess

	security *SecurityPolicy
	sink     EventSink

	// Configuração de timeout para encerramento gracioso.
	gracefulTimeout time.Duration
	// Timeout total máximo antes de SIGKILL; se zero, usa gracefulTimeout * 2.
	forceTimeout time.Duration

	// Caminhos-base permitidos para scripts (derivado de SecurityPolicy / implementações).
	// Mantido explícito para facilitar gates anti-Trojan.
	allowedRoots []string
}
type managedProcess struct {
	cmd          *exec.Cmd
	cancel       context.CancelFunc
	updateTarget string
	scriptPath   string
	startedAt    time.Time
}
type ShantillyEvent struct {
	SourceComponentID string
	Type              string
	Payload           map[string]interface{}
}
type RunAction struct {
	Script        string
	Args          []string
	Stdin         map[string]interface{}
	Env           map[string]string
	UpdateTarget  string
	Confirm       bool
	PromptSecrets []string
}
type SecurityPolicy struct {
	AllowedScripts []string
	DenyUnknown    bool
	ExtraRules     map[string]string
}
type UpdateTargetUpdate struct {
	TargetID string
	// Mode define como aplicar (ex.: "append", "replace"); por ora ilustrativo.
	Mode string
	// Content é o payload textual (stdout/stderr agregados) destinado ao target.
	Content string
}
type EventSink interface {
	EmitEvent(ev ShantillyEvent)
	EmitUpdate(update UpdateTargetUpdate)
}

func NewScriptRunner(sink EventSink, policy *SecurityPolicy) *ScriptRunner {
	r := &ScriptRunner{
		processes:       make(map[string]*managedProcess),
		security:        normalizePolicy(policy),
		sink:            sink,
		gracefulTimeout: 3 * time.Second,
		forceTimeout:    6 * time.Second,
	}
	r.allowedRoots = deriveAllowedRoots(r.security)
	return r
}
func normalizePolicy(p *SecurityPolicy) *SecurityPolicy {
	if p == nil {
		return &SecurityPolicy{
			AllowedScripts: nil,
			DenyUnknown:    true,
			ExtraRules:     map[string]string{},
		}
	}
	if !p.DenyUnknown {
		// Arquitetura define deny-by-default como baseline;
		// se configurado false, normalizamos para true para evitar brechas.
		p.DenyUnknown = true
	}
	if p.ExtraRules == nil {
		p.ExtraRules = map[string]string{}
	}
	return p
}
func deriveAllowedRoots(p *SecurityPolicy) []string {
	if p == nil {
		return nil
	}
	var roots []string
	for _, s := range p.AllowedScripts {
		if s == "" {
			continue
		}
		dir := s
		if !isDirLike(s) {
			dir = filepath.Dir(s)
		}
		if dir != "" && dir != "." {
			roots = append(roots, dir)
		}
	}
	return roots
}
func isDirLike(path string) bool {
	// Heurística mínima: se termina com '/', tratamos como diretório.
	return len(path) > 0 && path[len(path)-1] == '/'
}
func (r *ScriptRunner) Run(run RunAction) error {
	if err := r.validateRunAction(run); err != nil {
		r.emitErrorEvent("scriptrunner", "script.error", run, err)
		return err
	}

	updateTarget := run.UpdateTarget
	if updateTarget == "" {
		// Mesmo sem update_target, continuamos execução,
		// mas regra 1 processo por target não se aplica.
	}

	if updateTarget != "" {
		if err := r.ensureSingleProcess(updateTarget); err != nil {
			r.emitErrorEvent(updateTarget, "script.error", run, err)
			return err
		}
	}

	cmd, ctx, cancel, err := r.buildCommand(run)
	if err != nil {
		r.emitErrorEvent(updateTargetOr(run, "scriptrunner"), "script.error", run, err)
		return err
	}

	mp := &managedProcess{
		cmd:          cmd,
		cancel:       cancel,
		updateTarget: updateTarget,
		scriptPath:   run.Script,
		startedAt:    time.Now(),
	}

	if updateTarget != "" {
		r.registerProcess(updateTarget, mp)
	}

	r.emitStartEvent(run)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Execução síncrona; integração TEA pode envolver goroutine.
	err = cmd.Run()
	exitErr := err

	// Limpa processo ativo se ainda apontando para este.
	if updateTarget != "" {
		r.clearProcess(updateTarget, mp)
	}

	// Context cancellation (via ensureSingleProcess) sinaliza cancelamento concorrente.
	select {
	case <-ctx.Done():
		// Cancelado por novo run ou timeout.
		// Se o erro foi por contexto, tratamos como script.cancelled.
		if errors.Is(exitErr, context.Canceled) || isExitSignalTerminated(exitErr) {
			r.emitCancelledEvent(updateTargetOr(run, "scriptrunner"), run, stdoutBuf.String(), stderrBuf.String())
			return exitErr
		}
	default:
	}

	combinedOut := stdoutBuf.String()
	combinedErr := stderrBuf.String()

	if exitErr != nil {
		// Erro de execução.
		r.emitOutputUpdate(updateTarget, combinedOut, combinedErr)
		r.emitErrorEvent(updateTargetOr(run, "scriptrunner"), "script.error", run, exitErr)
		return exitErr
	}

	// Sucesso.
	r.emitOutputUpdate(updateTarget, combinedOut, combinedErr)
	r.emitCompleteEvent(updateTargetOr(run, "scriptrunner"), run)
	return nil
}
func (r *ScriptRunner) validateRunAction(run RunAction) error {
	if run.Script == "" {
		return fmt.Errorf("E1.4 violation: RunAction.Script is required")
	}

	// Whitelist / deny-by-default.
	if r.security != nil && r.security.DenyUnknown {
		if !r.isAllowedScript(run.Script) {
			return fmt.Errorf("E1.5/E1.7 violation: script '%s' not allowed by SecurityPolicy", run.Script)
		}
	}

	// Campos Confirm/PromptSecrets não são tratados aqui:
	// espera-se que EventManager/Modal Stack tenham aplicado gate JIT
	// antes de invocar o ScriptRunner.

	return nil
}
func (r *ScriptRunner) isAllowedScript(script string) bool {
	if r.security == nil {
		return false
	}
	if len(r.security.AllowedScripts) == 0 && r.security.DenyUnknown {
		return false
	}

	abs, err := filepath.Abs(script)
	if err != nil {
		return false
	}

	for _, allowed := range r.security.AllowedScripts {
		if allowed == "" {
			continue
		}
		allowedAbs, err := filepath.Abs(allowed)
		if err != nil {
			continue
		}
		if abs == allowedAbs {
			return true
		}
	}

	for _, root := range r.allowedRoots {
		rootAbs, err := filepath.Abs(root)
		if err != nil {
			continue
		}
		if isSubpathOf(abs, rootAbs) {
			return true
		}
	}

	return false
}
func isSubpathOf(path, root string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	// Se rel não sobe diretórios, path está dentro de root.
	return rel != ".." && !startsWithDotDot(rel)
}
func startsWithDotDot(rel string) bool {
	return len(rel) >= 2 && rel[0:2] == ".."
}
func (r *ScriptRunner) registerProcess(updateTarget string, mp *managedProcess) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processes[updateTarget] = mp
}
func (r *ScriptRunner) clearProcess(updateTarget string, mp *managedProcess) {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.processes[updateTarget]
	if ok && current == mp {
		delete(r.processes, updateTarget)
	}
}
func (r *ScriptRunner) emitStartEvent(run RunAction) {
	if r.sink == nil {
		return
	}
	r.sink.EmitEvent(ShantillyEvent{
		SourceComponentID: run.UpdateTarget,
		Type:              "script.start",
		Payload: map[string]interface{}{
			"script": run.Script,
			"args":   run.Args,
		},
	})
}
func (r *ScriptRunner) emitOutputUpdate(updateTarget string, stdout, stderr string) {
	if r.sink == nil {
		return
	}
	if updateTarget == "" && stdout == "" && stderr == "" {
		return
	}

	content := stdout
	if stderr != "" {
		if content != "" {
			content += "\n"
		}
		content += stderr
	}

	if content == "" {
		return
	}

	r.sink.EmitUpdate(UpdateTargetUpdate{
		TargetID: updateTarget,
		Mode:     "append",
		Content:  content,
	})
}
func (r *ScriptRunner) emitCompleteEvent(source string, run RunAction) {
	if r.sink == nil {
		return
	}
	r.sink.EmitEvent(ShantillyEvent{
		SourceComponentID: source,
		Type:              "script.complete",
		Payload: map[string]interface{}{
			"script": run.Script,
		},
	})
}
func (r *ScriptRunner) emitCancelledEvent(source string, run RunAction, stdout, stderr string) {
	if r.sink == nil {
		return
	}
	r.emitOutputUpdate(run.UpdateTarget, stdout, stderr)
	r.sink.EmitEvent(ShantillyEvent{
		SourceComponentID: source,
		Type:              "script.cancelled",
		Payload: map[string]interface{}{
			"script": run.Script,
		},
	})
}
func (r *ScriptRunner) emitErrorEvent(source string, eventType string, run RunAction, err error) {
	if r.sink == nil {
		return
	}
	if eventType == "" {
		eventType = "script.error"
	}
	payload := map[string]interface{}{
		"error":  err.Error(),
		"script": run.Script,
	}
	if run.UpdateTarget != "" {
		payload["update_target"] = run.UpdateTarget
	}

	// Além do evento estruturado, quando houver update_target associado,
	// enviamos também uma atualização textual para o viewport correspondente,
	// garantindo feedback visível para o usuário (ex.: permissão negada,
	// script não permitido pela SecurityPolicy, script inexistente, etc.).
	if run.UpdateTarget != "" {
		msg := fmt.Sprintf("[erro] Falha ao executar script '%s'\nMotivo: %s", run.Script, err.Error())
		r.sink.EmitUpdate(UpdateTargetUpdate{
			TargetID: run.UpdateTarget,
			Mode:     "append",
			Content:  msg,
		})
	}

	r.sink.EmitEvent(ShantillyEvent{
		SourceComponentID: source,
		Type:              eventType,
		Payload:           payload,
	})
}
func isExitSignalTerminated(err error) bool {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		ws, ok := exitErr.Sys().(syscall.WaitStatus)
		if ok {
			sig := ws.Signal()
			return sig == syscall.SIGTERM || sig == syscall.SIGKILL
		}
	}
	return false
}
func updateTargetOr(run RunAction, fallback string) string {
	if run.UpdateTarget != "" {
		return run.UpdateTarget
	}
	return fallback
}
