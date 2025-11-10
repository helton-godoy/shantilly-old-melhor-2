package runner

// Epic: 1 - Runtime TUI Declarativo — Fundação do Runtime (v2.0)
// Bloco: E1.4 — ScriptRunner (Wave 4) + integrações futuras E1.5/E1.7.
//
// Mandatos normativos (resumo, ver docs oficiais):
// - Implementação exclusiva em internal/runtime/runner/**.
// - Entrada exclusiva: RunAction produzido pelo pipeline:
//   ShantillyEvent → EventManager → on: (OnHandler) → RunAction.
// - Saída exclusiva:
//   - ShantillyEvent estruturado (ex.: script.start, script.output, script.error, script.complete, script.cancelled),
//   - e/ou atualização declarativa de update_target (ex.: conteúdo para Viewport).
// - Um processo por update_target (E1.4-3):
//   - Novo RunAction com mesmo UpdateTarget deve:
//     - Enviar SIGTERM ao processo anterior,
//     - Aguardar timeout configurado,
//     - Aplicar SIGKILL se necessário,
//     - Só então iniciar o novo processo.
// - Segurança (E1.4/E1.5/E1.7):
//   - Integração obrigatória com SecurityPolicy (deny by default, whitelist).
//   - Execuções fora da política DEVEM ser bloqueadas com erro estruturado.
// - Proibições:
//   - Nenhum os.Exit aqui (Gate: 1.x.no-osexit-core.yml).
//   - Nenhum uso de os/exec/syscall equivalente fora de internal/runtime/runner/**.
//   - Nenhuma execução de script fora de RunAction declarado.
// - Rastreabilidade QA/Governança:
//   - Gate 1.x.scriptrunner-and-update-target.yml,
//   - Gate 1.x.no-osexit-core.yml,
//   - Gate 1.x.security-jit-anti-trojan.yml.
//
// Este arquivo implementa o ScriptRunner como serviço puro, pronto para integração
// pelo EventManager e pelo loop TEA, sem dependências diretas de UI/legacy.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// ShantillyEvent espelha o modelo normativo
// (docs/architecture/data-models.md#5-shantillyevent-e13--tipo-de-evento-normalizado).
// Em implementação futura, deve ser importado de pkg/declarative.
type ShantillyEvent struct {
	SourceComponentID string
	Type              string
	Payload           map[string]interface{}
}

// RunAction espelha o contrato normativo
// (docs/architecture/data-models.md#6-runaction-e14--execucao-declarativa).
// Em implementação futura, deve ser importado de pkg/declarative.
type RunAction struct {
	Script        string
	Args          []string
	Stdin         map[string]interface{}
	Env           map[string]string
	UpdateTarget  string
	Confirm       bool
	PromptSecrets []string
}

// SecurityPolicy espelha o contrato normativo
// (docs/architecture/data-models.md#8-securitypolicy-e15--anti-trojan-yaml).
// Em implementação futura, deve ser importado de pkg/declarative.
type SecurityPolicy struct {
	AllowedScripts []string
	DenyUnknown    bool
	ExtraRules     map[string]string
}

// UpdateTargetUpdate representa uma atualização declarativa para um destino lógico
// (ex.: viewport). Não acopla diretamente à UI; o EventManager/Runtime decide como aplicar.
type UpdateTargetUpdate struct {
	TargetID string
	// Mode define como aplicar (ex.: "append", "replace"); por ora ilustrativo.
	Mode string
	// Content é o payload textual (stdout/stderr agregados) destinado ao target.
	Content string
}

// EventSink define a API mínima para o ScriptRunner emitir eventos de lifecycle
// de forma desacoplada (E1.4-4).
//
// Implementações típicas:
// - Adaptador para EventManager (internal/runtime/event) gerando ShantillyEvent interno/tea.Msg.
// - Adaptador de teste que captura eventos para asserções.
type EventSink interface {
	EmitEvent(ev ShantillyEvent)
	EmitUpdate(update UpdateTargetUpdate)
}

// ScriptRunner é o executor único de RunAction (E1.4).
//
// Responsabilidades centrais:
// - Validar RunAction conforme SecurityPolicy (deny by default + whitelist).
// - Garantir 1 processo ativo por UpdateTarget (regra de segurança).
// - Orquestrar ciclo SIGTERM→SIGKILL em substituições.
// - Executar scripts via os/exec (somente aqui, nunca fora de internal/runtime/runner/**).
// - Emitir eventos estruturados e atualizações de target via EventSink.
// - Nunca chamar os.Exit; erros fluem via retornos/eventos.
//
// NOTA: Esta implementação é síncrona no método Run para simplificar integração inicial;
// ela pode ser adaptada para modelo baseado em goroutine/tea.Cmd mantendo as invariantes.
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

// managedProcess guarda estado de um processo associado a um update_target.
type managedProcess struct {
	cmd          *exec.Cmd
	cancel       context.CancelFunc
	updateTarget string
	scriptPath   string
	startedAt    time.Time
}

// NewScriptRunner cria um ScriptRunner aderente às regras normativas.
//
// - sink: implementação de EventSink para integração com EventManager.
// - policy: SecurityPolicy opcional; se nil, aplica deny-by-default efetivo.
//
// Invariantes de segurança aplicadas aqui (E1.4/E1.5/E1.7).
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

// normalizePolicy aplica defaults normativos:
//
// - DenyUnknown: true efetivo mesmo se não informado.
// - Lista de scripts vazia + DenyUnknown true -> nenhuma execução permitida.
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

// deriveAllowedRoots extrai diretórios raiz permitidos a partir da SecurityPolicy.
//
// Estratégia simples:
// - Para cada caminho em AllowedScripts, considera-se o dir como root permitido.
// - Esta função é deliberadamente conservadora para gates Anti-Trojan.
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

// Run executa um RunAction em conformidade com os contratos.
//
// Fluxo:
// 1) Valida RunAction + SecurityPolicy (E1.4/E1.5/E1.7).
// 2) Aplica regra 1 processo por update_target (cancelamento anterior).
// 3) Monta comando seguro (args/stdin/env).
// 4) Executa script e agrega stdout/stderr.
// 5) Emite:
//   - script.start,
//   - script.output/update_target,
//   - script.complete ou script.error/script.cancelled.
//
// Este método NÃO chama os.Exit e NÃO toca UI diretamente.
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

// validateRunAction aplica invariantes básicos + política de segurança.
//
// Pontos principais:
//   - Script não vazio.
//   - Política deny-by-default + AllowedScripts/roots.
//   - Proíbe execuções fora dos caminhos permitidos.
//   - Confirm/PromptSecrets são tratados pelo pipeline Modal Stack (Wave 5);
//     aqui apenas assumimos que RunAction recebido já respeitou esses gates.
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

// isAllowedScript verifica se o script está permitido pela SecurityPolicy.
//
// Regras conservadoras:
// - Se AllowedScripts vazio e DenyUnknown true: nada é permitido.
// - Match exato por caminho OU prefixo (diretório raiz permitido).
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

// isSubpathOf verifica se path está dentro de root (prefixo de diretório seguro).
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

// ensureSingleProcess aplica a regra 1 processo por update_target.
//
// Se já existe processo associado ao target:
// - Solicita cancelamento via contexto,
// - Envia SIGTERM,
// - Aguarda gracefulTimeout,
// - Caso ainda ativo, envia SIGKILL até forceTimeout.
// - Emite evento script.cancelled apropriado (feito em Run).
func (r *ScriptRunner) ensureSingleProcess(updateTarget string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.processes[updateTarget]
	if !ok || existing == nil || existing.cmd == nil {
		return nil
	}

	// Solicita cancelamento de contexto se disponível.
	if existing.cancel != nil {
		existing.cancel()
	}

	// Envia SIGTERM.
	_ = existing.cmd.Process.Signal(syscall.SIGTERM)

	done := make(chan error, 1)
	go func(cmd *exec.Cmd) {
		done <- cmd.Wait()
	}(existing.cmd)

	select {
	case <-time.After(r.gracefulTimeout):
		// Tenta SIGKILL se ainda não terminou.
		_ = existing.cmd.Process.Signal(syscall.SIGKILL)
	case <-done:
		// Processo terminou dentro do grace period.
		return nil
	}

	// Aguarda final após SIGKILL (limitado por forceTimeout).
	select {
	case <-done:
	case <-time.After(r.forceTimeout):
		// Se ainda não saiu aqui, tratamos como erro grave mas não chamamos os.Exit.
		return fmt.Errorf("E1.4 violation: existing process for update_target '%s' did not terminate after SIGKILL", updateTarget)
	}

	return nil
}

// registerProcess registra um processo ativo para o update_target.
func (r *ScriptRunner) registerProcess(updateTarget string, mp *managedProcess) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processes[updateTarget] = mp
}

// clearProcess remove o processo ativo se corresponder ao mp fornecido.
func (r *ScriptRunner) clearProcess(updateTarget string, mp *managedProcess) {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.processes[updateTarget]
	if ok && current == mp {
		delete(r.processes, updateTarget)
	}
}

// buildCommand monta o comando exec.Cmd com contexto cancelável.
//
//   - Usa context.WithCancel para permitir cancelamento por ensureSingleProcess.
//   - Popula args/env a partir do RunAction.
//   - Serialização simples de Stdin: se presente, escreve JSON ou similar em stdin.
//     Aqui mantemos placeholder mínimo (não serializa para evitar dependências extras);
//     a integração real pode ser plugada posteriormente mantendo o contrato.
func (r *ScriptRunner) buildCommand(run RunAction) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Monta comando com script e args.
	args := append([]string{run.Script}, run.Args...)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Ambiente: aplicamos apenas variáveis declaradas em RunAction.Env.
	if len(run.Env) > 0 {
		var env []string
		for k, v := range run.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = append(cmd.Env, env...)
	}

	// STDIN: placeholder — para Epic 1, podemos não enviar nada
	// ou aplicar serialização simples em iterações futuras.
	// A ausência de manipulação aqui evita acoplamento prematuro.

	return cmd, ctx, cancel, nil
}

// emitStartEvent emite evento script.start.
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

// emitOutputUpdate envia atualização declarativa de output para update_target (se definido).
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

// emitCompleteEvent emite evento script.complete.
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

// emitCancelledEvent emite evento script.cancelled.
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

// emitErrorEvent emite evento de erro estruturado.
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
	r.sink.EmitEvent(ShantillyEvent{
		SourceComponentID: source,
		Type:              eventType,
		Payload:           payload,
	})
}

// isExitSignalTerminated detecta se um erro de execução foi causado por sinal de término.
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

// updateTargetOr retorna UpdateTarget se definido, senão fallback.
func updateTargetOr(run RunAction, fallback string) string {
	if run.UpdateTarget != "" {
		return run.UpdateTarget
	}
	return fallback
}
