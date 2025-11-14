package runtime

import (
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/runtime/layout"
	"shantilly/internal/runtime/modal"
	"shantilly/pkg/tui"
)

// MainModel é o modelo Bubbletea raiz.
// Nesta implementação simplificada, ele apenas delega para o LayoutManager
// sem usar uma biblioteca externa de overlay. A pilha de modais poderá ser
// reintroduzida futuramente com uma dependência estável.
type MainModel struct {
	layout tea.Model // O LayoutManager
	modal  tea.Model // Modal atual (overlay), se existir
}

func NewMainModel(layoutManager tea.Model) *MainModel {
	return &MainModel{
		layout: layoutManager,
	}
}

func (m *MainModel) Init() tea.Cmd {
	return m.layout.Init()
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Tratamento de mensagens relacionadas a modal.
	switch msg := msg.(type) {
	case tui.ShowModalMsg:
		// Quando recebemos ShowModalMsg, ativamos o modal como overlay.
		// Se o conteúdo vier nil (caso atual do EventManager), criamos um modal
		// padrão de confirmação.
		if msg.Content != nil {
			m.modal = msg.Content
		} else {
			m.modal = modal.NewConfirm("Confirmar execução?")
		}
		if m.modal != nil {
			initCmd := m.modal.Init()
			if initCmd != nil {
				cmds = append(cmds, initCmd)
			}
		}
		return m, tea.Batch(cmds...)
	}

	// Se houver modal ativo, ele recebe prioridade de atualização.
	if m.modal != nil {
		var modalCmd tea.Cmd
		m.modal, modalCmd = m.modal.Update(msg)
		if modalCmd != nil {
			cmds = append(cmds, modalCmd)
		}

		// O modal de confirmação emite tui.ModalResultMsg quando é concluído.
		if _, ok := msg.(tui.ModalResultMsg); ok {
			// Ao receber o resultado, fechamos o modal (o LayoutManager/EventManager
			// já sabe tratar ModalResultMsg e continuará o fluxo JIT).
			m.modal = nil
		}

		return m, tea.Batch(cmds...)
	}

	// Sem modal ativo: delega normalmente para o LayoutManager.
	var cmd tea.Cmd
	m.layout, cmd = m.layout.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if len(cmds) == 0 {
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *MainModel) View() string {
	// Renderiza o layout base e, se houver, o modal como overlay simples.
	base := m.layout.View()
	if m.modal == nil {
		return base
	}
	return base + "\n" + m.modal.View()
}

// NewLayoutMainModel é um helper para criar o MainModel a partir de um LayoutManager.
func NewLayoutMainModel(lm *layout.Manager) *MainModel {
	return NewMainModel(lm)
}
