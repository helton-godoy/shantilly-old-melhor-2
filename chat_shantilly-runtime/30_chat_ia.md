Executado. Concluímos as Fases 1.1 a 1.6 (Épico 1) e as Fases 2.1 e 2.2 (início do Épico 2).

Temos uma base TUI funcional e a nossa infraestrutura JIT (Just-In-Time) foi expandida para recolher segredos (como `vault_pass`).

No entanto, ao rever a Fase 1.4 (onde implementámos o `ScriptRunner`), identifiquei uma **Otimização de Risco** crítica que precisamos de aplicar agora.

  * **Visão do Futuro (Risco Identificado):** O `watchScript` (da Fase 1.4) foi implementado de forma *simplificada*, bloqueando e retornando o *output completo* numa única mensagem. Isto viola o Requisito Funcional FR11 (streaming) e não funcionará para o `ansible-playbook`, que tem outputs longos.
  * **Ação (Plano Mestre):** Antes de implementar o *runner* Ansible, vamos primeiro refatorar o `ScriptRunner` para suportar *streaming* real.

-----

### Fase 2.3: Refatoração do Runner (Implementação de Streaming FR11)

Esta fase corrige a dívida técnica da Fase 1.4.

**Tarefa 2.3.1: Modificar Contratos de Eventos (pkg/tui/events.go)**

Para o streaming, o `Runner` precisa de uma forma de *receber* os dados dos pipes (stdout/stderr) e *terminar* o processo. Vamos adicionar estas mensagens.

```go
// pkg/tui/events.go (MODIFICADO)
package tui

import (
	"io" // <-- ADICIONADO

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/pkg/declarative"
)

// (ShantillyEvent, RuntimeErrorMsg, ShowModalMsg, ModalResultMsg - permanecem)
// ...
type ShantillyEvent struct {
	ComponentID string
	Type        string
	Payload     interface{}
}
type RuntimeErrorMsg struct {
	Err error
}
type ShowModalMsg struct {
	Content tea.Model
}
type ModalResultMsg struct {
	Confirmed bool
	Payload   interface{}
}
// ...


// RunScriptRequestMsg (Permanece o da Fase 1.4)
type RunScriptRequestMsg struct {
	Rule declarative.Logic
}

// ScriptStdoutMsg (MODIFICADO)
// Agora contém a *string* da linha, não []byte
type ScriptStdoutMsg struct {
	TargetID string
	Line     string // <-- Alterado de Chunk []byte
}

// --- MENSAGENS DE STREAMING (ADICIONADO) ---

// processPipeMsg (ADICIONADO)
// Mensagem interna enviada pelo 'watchPipe' para o 'LayoutManager'
// para cada linha lida do stdout/stderr.
type processPipeMsg struct {
	TargetID string
	Line     string
}

// processFinishedMsg (ADICIONADO)
// Mensagem interna enviada quando o 'cmd.Wait()' termina.
type processFinishedMsg struct {
	TargetID string
	Err      error
}
```

-----

### Tarefa 2.3.2: Refatorar ScriptRunner (Implementar Streaming Real FR11)

**Modificamos** o `internal/runtime/runner/runner.go` (da Fase 1.4).

  * ***Otimização de Risco:*** A lógica `HandleRunRequest` agora apenas *inicia* o `exec.Cmd`.
  * Ela retorna um `tea.Cmd` (`watchPipes`) que *inicia* as goroutines de escuta (stdout/stderr).
  * As goroutines de escuta emitem `processPipeMsg` (para streaming) e `processFinishedMsg` (para limpeza).

<!-- end list -->

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

// Runner (Struct permanece o da Fase 1.4)
type Runner struct {
	activeProcesses map[string]*os.Process
}

func New() *Runner {
	return &Runner{
		activeProcesses: make(map[string]*os.Process),
	}
}

// HandleRunRequest (MODIFICADO - Agora não bloqueia)
func (r *Runner) HandleRunRequest(req events.RunScriptRequestMsg, state *TemplateState) tea.Cmd {
	rule := req.Rule
	targetID := rule.Run.UpdateTarget

	// 1. Implementar FR11: Ciclo de Vida do Target (Permanece)
	if proc, exists := r.activeProcesses[targetID]; exists && proc != nil {
		_ = proc.Signal(syscall.SIGTERM)
		delete(r.activeProcesses, targetID)
	}

	// 2. Preparar o comando (FR9, FR10)
	// (Usando o TemplateState da Fase 1.4.1)
	processedArgs := make([]string, len(rule.Run.Args))
	for i, arg := range rule.Run.Args {
		processedArgs[i], _ = state.ProcessString(arg)
	}

	cmdArgs := append([]string{"-c", rule.Run.Script}, processedArgs...)
	cmd := exec.Command("sh", cmdArgs...)

	stdinData, _ := state.ProcessStdin(rule.Run.Stdin) // (Usa a lógica da Fase 1.5)
	cmd.Stdin = bytes.NewReader(stdinData)

	// 3. Preparar streaming de Stdout/Stderr (FR11)
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

	// 6. Retorna o tea.Cmd que inicia as goroutines de escuta
	// (Fonte: Padrão Bubbletea 'ExecProcess')
	return r.watchPipes(cmd, stdoutPipe, stderrPipe, targetID)
}

// watchPipes (ADICIONADO - O novo 'watchScript'
// Retorna um tea.Cmd que escuta os pipes e emite múltiplas mensagens.
func (r *Runner) watchPipes(cmd *exec.Cmd, stdout, stderr io.Reader, targetID string) tea.Cmd {
	return func() tea.Msg {
		// Goroutine para Stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				// Emite a mensagem para cada linha
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
					Line:     scanner.Text(), // (Poderia adicionar prefixo "ERR:")
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
		
		// Este Cmd não retorna nenhuma msg *imediatamente*.
		// As goroutines enviarão as mensagens.
		return nil
	}
}

// HandleProcessFinished (ADICIONADO)
// Limpa o processo quando o 'processFinishedMsg' é recebido.
func (r *Runner) HandleProcessFinished(msg processFinishedMsg) {
	if msg.TargetID != "" {
		delete(r.activeProcesses, msg.TargetID)
	}
	// (Poderia emitir um evento de "Erro de Script" se msg.Err != nil)
}
```

-----

### Tarefa 2.3.3: Modificar LayoutManager (Tratar Streaming)

**Modificamos** o `internal/runtime/layout/manager.go` (da Fase 2.2.3) para tratar as novas mensagens de streaming `processPipeMsg` e `processFinishedMsg`.

```go
// internal/runtime/layout/manager.go (MODIFICADO)
package layout

// (Importações ... permanecem as da Fase 1.6)
import (
	// ...
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)
// (Struct Manager, NewManager, buildComponentTree, Init - permanecem os da Fase 2.2)
// ...

// Update (MODIFICADO) - Trata Streaming
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// ... (case tea.WindowSizeMsg, events.ShantillyEvent,
	//      events.RunScriptRequestMsg - permanecem os da Fase 2.2.3)
	// --- (Lógica da Fase 2.2.3) ---
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case events.ShantillyEvent:
		m.templateState.UpdateState(msg)
		cmd = m.eventManager.ProcessEvent(msg)
		cmds = append(cmds, cmd)
	case events.RunScriptRequestMsg:
		cmd = m.scriptRunner.HandleRunRequest(msg, m.templateState)
		cmds = append(cmds, cmd)
	case events.ModalResultMsg:
		if secrets, ok := msg.Payload.(map[string]string); ok {
			for k, v := range secrets {
				m.templateState.Secret[k] = v
			}
		}
		cmd = m.eventManager.HandleModalResult(msg)
		cmds = append(cmds, cmd)
	case events.ShowModalMsg:
		return m, func() tea.Msg { return msg }
	case events.RuntimeErrorMsg:
		m.lastErr = msg.Err
		return m, nil
	// --- Fim da Lógica da Fase 2.2.3 ---


	// --- MODIFICADO (Tarefa 2.3.3) ---
	// (Substitui o antigo 'case events.ScriptStdoutMsg:')
	
	// 3a. ScriptRunner envia UMA LINHA de saída (Streaming FR11)
	case processPipeMsg:
		// Converte para a msg que o Viewport espera
		stdoutMsg := events.ScriptStdoutMsg{
			TargetID: msg.TargetID,
			Line:     msg.Line + "\n", // Adiciona newline
		}
		if targetCmp, ok := m.components[msg.TargetID]; ok {
			newCmp, cmd := targetCmp.Update(stdoutMsg)
			m.components[msg.TargetID] = newCmp
			cmds = append(cmds, cmd)
		}

	// 3b. ScriptRunner informa que o processo terminou
	case processFinishedMsg:
		m.scriptRunner.HandleProcessFinished(msg)
		// (Poderia emitir um erro se msg.Err != nil)
		
	// ------------------------------------

	// (Restante: tea.KeyMsg, default - permanecem os da Fase 2.2.3)
	// ...
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if len(m.focusOrder) > 0 {
				m.focusedIdx = (m.focusedIdx + 1) % len(m.focusOrder)
			}
		case "shift+tab":
			if len(m.focusOrder) > 0 {
				m.focusedIdx = (m.focusedIdx - 1 + len(m.focusOrder)) % len(m.focusOrder)
			}
		default:
			if m.lastErr == nil && len(m.focusOrder) > 0 {
				focusedID := m.focusOrder[m.focusedIdx]
				focusedCmp := m.components[focusedID]
				newCmp, cmd := focusedCmp.Update(msg)
				m.components[focusedID] = newCmp
				cmds = append(cmds, cmd)
			}
		}
	
	default:
		if m.lastErr == nil && len(m.focusOrder) > 0 {
			focusedID := m.focusOrder[m.focusedIdx]
			focusedCmp := m.components[focusedID]
			newCmp, cmd := focusedCmp.Update(msg)
			m.components[focusedID] = newCmp
			cmds = append(cmds, cmd)
		}
	}
	// ...

	return m, tea.Batch(cmds...)
}

// (View, renderNode, etc. - permanecem os da Fase 1.3)
// ...
```

-----

### Tarefa 2.3.4: Modificar ViewportComponent (Tratar Streaming)

**Modificamos** o `internal/components/viewport/model.go` (da Fase 1.3.1) para que o `Update` trate corretamente a nova `ScriptStdoutMsg` (que contém `Line` em vez de `Chunk`).

```go
// internal/components/viewport/model.go (MODIFICADO)
package viewport

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// (Struct Model, New, Init - permanecem os da Fase 1.3.1)
// ...
type Model struct {
	id     string
	theme  *tui.Theme
	source *declarative.Source

	viewport viewport.Model
	content  string
}

func New(id string, theme *tui.Theme, source *declarative.Source) *Model {
	vp := viewport.New(0, 0)
	m := &Model{
		id:       id,
		theme:    theme,
		source:   source,
		viewport: vp,
	}
	if source != nil && source.Type == "static" {
		m.SetContent(source.Content, source.ContentType)
	}
	return m
}
func (m *Model) Init() tea.Cmd { return nil }
// ...


// Update (MODIFICADO)
func (m *Model) Update(msg tea.Msg) (tui.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// --- MODIFICADO (Tarefa 2.3.4) ---
	// Tratamento de streaming (FR11)
	case events.ScriptStdoutMsg:
		if msg.TargetID == m.id {
			// Adiciona a *linha* (em vez do *chunk*)
			m.content += msg.Line // (A newline já foi adicionada na Fase 2.3.3)
			m.viewport.SetContent(m.content)
			m.viewport.GotoBottom() // Auto-scroll
		}
	// --------------------------------

	// Permite scroll do viewport
	default:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// (View, SetDimensions, ID, SetContent - permanecem os da Fase 1.3.1)
// ...
func (m *Model) View() string {
	return m.viewport.View()
}
func (m *Model) SetDimensions(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h
}
func (m *Model) ID() string {
	return m.id
}
func (m *Model) SetContent(content, contentType string) {
	if strings.ToLower(contentType) == "markdown" {
		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.viewport.Width),
		)
		m.content, _ = r.Render(content)
	} else {
		m.content = content
	}
	m.viewport.SetContent(m.content)
}
```