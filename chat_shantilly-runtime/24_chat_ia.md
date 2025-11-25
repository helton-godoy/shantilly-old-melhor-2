Executado. Concluímos as Fases 1.1, 1.2, 1.3, 1.4 e 1.5. O nosso runtime TUI está funcional: renderiza layouts, exibe componentes (incluindo o formulário v1.0 encapsulado) e executa scripts (sem segurança).

Iniciamos agora a **Fase 1.6: Segurança JIT (Pilha Modal)**.

Esta é a fase final do Épico 1 e implementa a funcionalidade de segurança `confirm:` (definida na Seção 4.6 do `architecture.md`).

  * ***Otimização de Risco (Visão do Futuro):*** Para evitar os erros de "baixa qualidade" do Charm que mencionaste, vamos usar a biblioteca `rmhubbert/bubbletea-overlay` (mencionada como inspiração no PRD) para gerir os modais. Isto desacopla o `LayoutManager` da gestão da pilha modal, alinhando-se perfeitamente com a nossa arquitetura de separação de responsabilidades.

Aqui estão os artefatos da Fase 1.6.

-----

### Tarefa 1.6.1: Modificar Contratos de Eventos

**Modificamos** o `pkg/tui/events.go` (da Fase 1.4) para adicionar o `ModalResultMsg`. Este é o evento que o modal emite quando é fechado (Confirmado ou Abortado).

```go
// pkg/tui/events.go (MODIFICADO)
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/pkg/declarative"
)

// (ShantillyEvent, RuntimeErrorMsg, ScriptStdoutMsg - permanecem os da Fase 1.4)
// ...
type ShantillyEvent struct {
	ComponentID string
	Type        string
	Payload     interface{}
}
type RuntimeErrorMsg struct {
	Err error
}
type ScriptStdoutMsg struct {
	TargetID string
	Chunk    []byte
}
type RunScriptRequestMsg struct {
	Rule declarative.Logic
}

// ShowModalMsg (da Fase 1.1) - O 'Content' é o modelo do modal.
// Fonte: docs/architecture/core-workflows.md (Workflow 7.3)
type ShowModalMsg struct {
	Content tea.Model // O modelo do modal (ex: modal.Confirm)
}

// ModalResultMsg (ADICIONADO)
// Emitido pelo modal (ex: modal.Confirm) quando é fechado.
type ModalResultMsg struct {
	Confirmed bool        // true se 'Sim', false se 'Não'/Abortado
	Payload   interface{} // (Para Fase 1.6.x - PromptSecrets)
}
```

-----

### Tarefa 1.6.2: Implementar Componente Modal (Confirm)

Criamos o componente TUI real para o modal de confirmação.

  * ***Otimização de Risco:*** Usamos `huh.Confirm` (que sabemos estar disponível da Fase 1.5) para garantir consistência visual e evitar a invenção de um novo componente de UI.

<!-- end list -->

```go
// internal/runtime/modal/confirm.go
package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// ConfirmModel é um componente TUI (tea.Model) para
// confirmação JIT (Just-In-Time).
// Fonte: docs/architecture.md (Seção 7.3)
type Model struct {
	form *huh.Form
}

func NewConfirm(title string) *Model {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title(title).
				Affirmative("Sim").
				Negative("Não"),
		),
	).WithWidth(40) // Largura fixa para o modal

	return &Model{
		form: form,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	// Se o 'huh.Confirm' for concluído (Sim ou Não)
	if m.form.State == huh.StateCompleted {
		confirmed := m.form.Get("confirm").(bool)
		// Emite o evento de resultado (Tarefa 1.6.1)
		cmd = func() tea.Msg { return events.ModalResultMsg{Confirmed: confirmed} }
		cmds = append(cmds, cmd)
	}

	// Se for Abortado (Ctrl+C, Esc)
	if m.form.State == huh.StateAborted {
		// Emite "Não Confirmado"
		cmd = func() tea.Msg { return events.ModalResultMsg{Confirmed: false} }
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// Centraliza o formulário do modal
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true).
		BorderForeground(lipgloss.Color("63")). // Cor primária do tema
		Padding(1, 2).
		Render(m.form.View())
}
```

-----

### Tarefa 1.6.3: Modificar EventManager (Lógica JIT)

**Modificamos** o `internal/runtime/event/manager.go` (da Fase 1.4) para implementar a lógica JIT.

  * ***Otimização de Risco:*** O `EventManager` agora armazena a `RunScriptRequestMsg` pendente e, em vez de executá-la, emite `ShowModalMsg` (Workflow 7.3).

<!-- end list -->

```go
// internal/runtime/event/manager.go (MODIFICADO)
package event

import (
	"fmt"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/runtime/modal" // <-- ADICIONADO
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Manager é o motor de lógica que trata do bloco 'on:'.
type Manager struct {
	rules []declarative.Logic
	
	// ADICIONADO: Armazena a ação pendente de confirmação JIT
	// Fonte: docs/architecture.md (Seção 7.3)
	pendingRequest *events.RunScriptRequestMsg
}

func New(rules []declarative.Logic) *Manager {
	return &Manager{
		rules: rules,
	}
}

// ProcessEvent (MODIFICADO)
func (m *Manager) ProcessEvent(event events.ShantillyEvent) tea.Cmd {
	for _, rule := range m.rules {
		if rule.Event == event.Type {
			
			// Constrói a Ação de Script
			req := events.RunScriptRequestMsg{ Rule: rule }

			// --- LÓGICA JIT (ADICIONADO) ---
			// Fonte: docs/architecture.md (Seção 4.6, 7.3)
			if rule.Confirm {
				// 1. Armazena a ação pendente
				m.pendingRequest = &req
				
				// 2. Emite o ShowModalMsg (Tarefa 1.6.1)
				// em vez de RunScriptRequestMsg
				title := fmt.Sprintf("Executar '%s'?", rule.Run.Script)
				modalCmd := func() tea.Msg {
					return events.ShowModalMsg{
						Content: modal.NewConfirm(title),
					}
				}
				return modalCmd
			}
			// (Na Fase 2.x, verificaremos 'rule.PromptSecrets' aqui)

			// --- Fim da Lógica JIT ---

			// Se não houver JIT, executa imediatamente (como na Fase 1.4)
			return func() tea.Msg { return req }
		}
	}
	return nil
}

// HandleModalResult (ADICIONADO)
// Chamado pelo LayoutManager/MainModel quando o modal (Tarefa 1.6.2) fecha.
func (m *Manager) HandleModalResult(msg events.ModalResultMsg) tea.Cmd {
	if m.pendingRequest == nil {
		return nil // Nenhuma ação pendente
	}

	// 1. Limpa a ação pendente
	req := m.pendingRequest
	m.pendingRequest = nil

	// 2. Se confirmado, emite a ação pendente (RunScriptRequestMsg)
	if msg.Confirmed {
		return func() tea.Msg { return *req }
	}

	// 3. Se não confirmado, cancela a ação (não faz nada)
	return nil
}
```

-----

### Tarefa 1.6.4: Criar o Modelo Raiz (com Overlay)

Conforme a nossa "Otimização de Risco", criamos um novo modelo raiz (`MainModel`) que usa `bubbletea-overlay` para "hospedar" o `LayoutManager`.

```go
// internal/runtime/main_model.go
package runtime

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/runtime/layout"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
	"github.com/rmhubbert/bubbletea-overlay/overlay"
)

// MainModel é o novo modelo Bubbletea raiz.
// Ele gere a sobreposição (overlay) de modais sobre o layout principal.
// Fonte: docs/prd.md (Inspiração bubbletea-overlay)
type MainModel struct {
	overlay overlay.Model // O gestor de overlay
	layout  tea.Model     // O LayoutManager (Fase 1.4)
}

func NewMainModel(layoutManager *layout.Manager) *MainModel {
	return &MainModel{
		overlay: overlay.New(layoutManager),
		layout:  layoutManager,
	}
}

func (m *MainModel) Init() tea.Cmd {
	return m.layout.Init()
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Encaminha eventos de overlay PRIMEIRO
	if m.overlay.Focused() {
		m.overlay, cmd = m.overlay.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	// Trata mensagens de controle JIT
	switch msg := msg.(type) {
	// (Do EventManager, Tarefa 1.6.3)
	case events.ShowModalMsg:
		m.overlay.Set(msg.Content) // Mostra o modal
		return m, nil // Não encaminha para o layout

	// (Do Modal, Tarefa 1.6.2)
	case events.ModalResultMsg:
		m.overlay.Clear() // Fecha o modal
		// Encaminha o *resultado* para o LayoutManager,
		// que o encaminhará para o EventManager.
	}

	// Encaminha todas as outras mensagens para o LayoutManager
	m.layout, cmd = m.layout.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainModel) View() string {
	// Renderiza o layout (View do LayoutManager)
	// e o overlay (View do Modal) por cima.
	return m.overlay.View()
}
```

-----

### Tarefa 1.6.5: Modificar Ponto de Entrada do Runtime (Integração Final)

Finalmente, **modificamos** o `internal/runtime/runtime.go` (da Fase 1.4) para usar o novo `MainModel` (Tarefa 1.6.4) como o ponto de entrada do `bubbletea`, em vez do `LayoutManager`.

```go
// internal/runtime/runtime.go (MODIFICADO)
package runtime

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/runtime/layout"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	// (overlay.Model é usado em main_model.go)
)

// Start (MODIFICADO)
func Start(config *declarative.Config) error {
	// 1. Criar o tema central
	theme := tui.NewDefaultTheme()

	// 2. Criar o LayoutManager (como na Fase 1.4)
	layoutManager := layout.NewManager(&config.Root, theme, config.On)

	// 3. Criar o MainModel (Tarefa 1.6.4)
	// Este é o novo modelo raiz, que "embrulha" o layout.
	mainModel := NewMainModel(layoutManager)

	// 4. Iniciar o programa Bubbletea
	p := tea.NewProgram(
		mainModel, // <-- Inicia o MainModel, não o LayoutManager
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("erro durante a execução do runtime TUI: %w", err)
	}

	return nil
}
```

-----

### Tarefa 1.6.6: Modificar LayoutManager (Integração JIT)

Uma última modificação no `internal/runtime/layout/manager.go` (da Fase 1.4) para garantir que ele encaminhe os `ModalResultMsg` para o `EventManager`.

```go
// internal/runtime/layout/manager.go (MODIFICADO)
package layout

// (Importações ... permanecem as da Fase 1.4)
import (
	// ...
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)
// (Struct Manager, NewManager, buildComponentTree, Init - permanecem os da Fase 1.4)
// ...

// Update (MODIFICADO)
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// --- Delegação de Lógica (Modificado) ---
	// 1. Componente emite evento (ex: list:select)
	case events.ShantillyEvent:
		m.templateState.UpdateState(msg)
		cmd = m.eventManager.ProcessEvent(msg)
		cmds = append(cmds, cmd)

	// 2. EventManager pede para rodar script
	case events.RunScriptRequestMsg:
		cmd = m.scriptRunner.HandleRunRequest(msg, m.templateState)
		cmds = append(cmds, cmd)

	// 3. ScriptRunner envia saída (stdout/stderr)
	case events.ScriptStdoutMsg:
		if targetCmp, ok := m.components[msg.TargetID]; ok {
			newCmp, cmd := targetCmp.Update(msg)
			m.components[msg.TargetID] = newCmp
			cmds = append(cmds, cmd)
		}
		
	// --- ADICIONADO (Tarefa 1.6.6) ---
	// 4. Modal fecha e envia resultado
	case events.ModalResultMsg:
		// Devolve ao EventManager para tratar a lógica JIT
		cmd = m.eventManager.HandleModalResult(msg)
		cmds = append(cmds, cmd)
	
	// 5. EventManager (via JIT) pede para mostrar modal
	case events.ShowModalMsg:
		// Apenas retorna a mensagem. O MainModel (Tarefa 1.6.4)
		// irá intercetá-la e mostrá-la.
		return m, func() tea.Msg { return msg }
	// ------------------------------------

	case events.RuntimeErrorMsg:
		m.lastErr = msg.Err
		return m, nil

	// (Restante da lógica 'case tea.KeyMsg:' e 'default:'
	// permanecem os da Fase 1.4)
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

	return m, tea.Batch(cmds...)
}

// (View, renderNode, getComponent, isFocused - permanecem os da Fase 1.3)
// ...
```