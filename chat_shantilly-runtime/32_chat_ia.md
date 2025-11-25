Executado. Concluímos as Fases 1.1-1.6 (Épico 1) e as Fases 2.1-2.3 (preparação do Épico 2).

Temos uma base TUI funcional, a nossa infraestrutura JIT (Just-In-Time) foi expandida para recolher segredos (como `vault_pass`), e o nosso `ScriptRunner` (Fase 2.3) agora suporta *streaming* real (FR11).

Seguindo o **Plano Mestre Sequencial**, iniciamos a **Fase 2.4: Implementação do Runner Especialista (Ansible)**.

Esta é a fase final do Épico 2, onde implementamos a lógica que realmente executa o `ansible-playbook`.

-----

### Tarefa 2.4.1: Implementar o AnsibleRunner

Criamos a lógica de execução específica do Ansible. Este ficheiro traduz o *struct* `AnsiblePlaybookAction` (da Fase 2.1) num `exec.Cmd`.

  * ***Otimização de Risco (Visão do Futuro):*** Esta lógica integra-se *imediatamente* com o `TemplateState` (Fase 2.2) para extrair o `vault_pass` dos segredos JIT. Também utiliza o `watchPipes` (da Fase 2.3) para garantir o streaming.

<!-- end list -->

```go
// internal/runtime/runner/ansible.go
package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// HandleAnsibleRequest (ADICIONADO)
// Lógica específica para o runner 'ansible_playbook' (Épico 2)
// Fonte: docs/prd.md (Épico 2), pkg/declarative/models.go
func (r *Runner) HandleAnsibleRequest(rule declarative.Logic, state *TemplateState) tea.Cmd {
	config := rule.Run.AnsiblePlaybook
	if config == nil {
		return func() tea.Msg { return events.RuntimeErrorMsg{Err: fmt.Errorf("ansible_playbook config está nil")} }
	}
	
	targetID := rule.Run.UpdateTarget

	// 1. Implementar FR11: Ciclo de Vida do Target (Permanece)
	if proc, exists := r.activeProcesses[targetID]; exists && proc != nil {
		_ = proc.Signal(syscall.SIGTERM)
		delete(r.activeProcesses, targetID)
	}

	// 2. Construir o comando 'ansible-playbook'
	var cmdArgs []string
	cmdArgs = append(cmdArgs, config.Playbook) // Caminho do Playbook

	if config.Inventory != "" {
		cmdArgs = append(cmdArgs, "-i", config.Inventory)
	}

	// 3. Processar 'vars' (FR10)
	if len(config.Vars) > 0 {
		// Serializa o mapa 'vars' como JSON
		varsJSON, err := json.Marshal(config.Vars)
		if err != nil {
			return func() tea.Msg { return events.RuntimeErrorMsg{Err: fmt.Errorf("falha ao serializar extra-vars: %w", err)} }
		}
		cmdArgs = append(cmdArgs, "--extra-vars", string(varsJSON))
	}
	
	// 4. Processar 'args' (templates) (FR10)
	for _, arg := range rule.Run.Args {
		processedArg, _ := state.ProcessString(arg)
		cmdArgs = append(cmdArgs, processedArg)
	}
	
	cmd := exec.Command("ansible-playbook", cmdArgs...)

	// 5. Processar 'stdin' (FR10)
	stdinData, _ := state.ProcessStdin(rule.Run.Stdin)
	cmd.Stdin = bytes.NewReader(stdinData)

	// 6. Integrar com JIT (ask_vault_pass) (Épico 2)
	// Fonte: Fase 2.2 (internal/runtime/runner/templating.go)
	if config.AskVaultPass {
		if vaultPass, ok := state.Secret["vault_pass"]; ok && vaultPass != "" {
			// (Otimização: Passa via variável de ambiente)
			// (Uma implementação mais robusta usaria --vault-password-file)
			cmd.Env = append(os.Environ(), "ANSIBLE_VAULT_PASSWORD="+vaultPass)
		} else {
			// (JIT foi acionado mas o segredo não foi encontrado;
			// EventManager (Fase 2.2) já deve ter tratado isso, mas
			// por segurança, retornamos um erro)
			return func() tea.Msg { return events.RuntimeErrorMsg{Err: fmt.Errorf("ask_vault_pass foi true, mas 'vault_pass' não foi encontrado no estado")} }
		}
	}


	// 7. Preparar streaming (FR11) - Reutiliza a lógica da Fase 2.3
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	// 8. Iniciar o comando
	if err := cmd.Start(); err != nil {
		return func() tea.Msg { return events.RuntimeErrorMsg{Err: err} }
	}

	// 9. Armazenar o novo processo (FR11)
	if targetID != "" {
		r.activeProcesses[targetID] = cmd.Process
	}

	// 10. Retorna o tea.Cmd que assiste os pipes (da Fase 2.3)
	return r.watchPipes(cmd, stdoutPipe, stderrPipe, targetID)
}
```

-----

### Tarefa 2.4.2: Modificar o ScriptRunner (Dispatch)

**Modificamos** o `internal/runtime/runner/runner.go` (da Fase 2.3.2) para se tornar um *dispatcher*. A função `HandleRunRequest` agora decide qual *runner* (Script ou Ansible) deve tratar a requisição.

```go
// internal/runtime/runner/runner.go (MODIFICADO)
package runner

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Runner (Struct permanece o da Fase 2.3.2)
type Runner struct {
	activeProcesses map[string]*os.Process
}

func New() *Runner {
	return &Runner{
		activeProcesses: make(map[string]*os.Process),
	}
}

// HandleRunRequest (MODIFICADO - Agora é um DISPATCHER)
// Fonte: Fase 2.1 (pkg/declarative/models.go)
func (r *Runner) HandleRunRequest(req events.RunScriptRequestMsg, state *TemplateState) tea.Cmd {
	rule := req.Rule

	// 1. Decide qual runner usar
	
	// Épico 2: Runner Ansible (Tarefa 2.4.1)
	if rule.Run.AnsiblePlaybook != nil {
		return r.HandleAnsibleRequest(rule, state)
	}

	// Épico 1: Runner Genérico (Script) (FR9)
	if rule.Run.Script != "" {
		return r.HandleScriptRequest(rule, state)
	}

	// 2. Nenhum runner válido encontrado
	return func() tea.Msg { return events.RuntimeErrorMsg{Err: fmt.Errorf("regra 'run:' não especificou um runner válido (script: ou ansible_playbook:)")} }
}

// HandleScriptRequest (ADICIONADO - Lógica da Fase 2.3.2 movida)
// Lógica específica para o runner 'script:'
func (r *Runner) HandleScriptRequest(rule declarative.Logic, state *TemplateState) tea.Cmd {
	targetID := rule.Run.UpdateTarget

	// 1. Implementar FR11: Ciclo de Vida do Target (Permanece)
	if proc, exists := r.activeProcesses[targetID]; exists && proc != nil {
		_ = proc.Signal(syscall.SIGTERM)
		delete(r.activeProcesses, targetID)
	}

	// 2. Preparar o comando (FR9, FR10)
	processedArgs := make([]string, len(rule.Run.Args))
	for i, arg := range rule.Run.Args {
		processedArgs[i], _ = state.ProcessString(arg)
	}

	cmdArgs := append([]string{"-c", rule.Run.Script}, processedArgs...)
	cmd := exec.Command("sh", cmdArgs...)

	stdinData, _ := state.ProcessStdin(rule.Run.Stdin)
	cmd.Stdin = bytes.NewReader(stdinData)

	// 3. Preparar streaming (FR11)
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	// 4. Iniciar o comando
	if err := cmd.Start(); err != nil {
		return func() tea.Msg { return events.RuntimeErrorMsg{Err: err} }
	}

	// 5. Armazenar o novo processo (FR11)
	if targetID != "" {
		r.activeProcesses[targetID] = cmd.Process
	}

	// 6. Retorna o tea.Cmd que assiste os pipes (da Fase 2.3)
	return r.watchPipes(cmd, stdoutPipe, stderrPipe, targetID)
}


// watchPipes (Permanece o da Fase 2.3.2)
// Retorna um tea.Cmd que escuta os pipes e emite múltiplas mensagens.
func (r *Runner) watchPipes(cmd *exec.Cmd, stdout, stderr io.Reader, targetID string) tea.Cmd {
	return func() tea.Msg {
		// Goroutine para Stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				tea.NewProgram(nil).Send(processPipeMsg{
					TargetID: targetID,
					Line:     scanner.Text(),
				})
			}
		}()

		// Goroutine para Stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				tea.NewProgram(nil).Send(processPipeMsg{
					TargetID: targetID,
					Line:     scanner.Text(),
				})
			}
		}()

		// Goroutine para esperar o fim
		go func() {
			err := cmd.Wait()
			tea.NewProgram(nil).Send(processFinishedMsg{
				TargetID: targetID,
				Err:      err,
			})
		}()
		
		return nil
	}
}

// HandleProcessFinished (Permanece o da Fase 2.3.2)
// Limpa o processo quando o 'processFinishedMsg' é recebido.
func (r *Runner) HandleProcessFinished(msg processFinishedMsg) {
	if msg.TargetID != "" {
		delete(r.activeProcesses, msg.TargetID)
	}
}
```