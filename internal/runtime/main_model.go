package runtime

import (
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/runtime/layout"
)

// MainModel é o modelo Bubbletea raiz.
// Nesta implementação simplificada, ele apenas delega para o LayoutManager
// sem usar uma biblioteca externa de overlay. A pilha de modais poderá ser
// reintroduzida futuramente com uma dependência estável.
type MainModel struct {
	layout tea.Model // O LayoutManager
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
	var cmd tea.Cmd
	m.layout, cmd = m.layout.Update(msg)
	return m, cmd
}

func (m *MainModel) View() string {
	return m.layout.View()
}

// NewLayoutMainModel é um helper para criar o MainModel a partir de um LayoutManager.
func NewLayoutMainModel(lm *layout.Manager) *MainModel {
	return NewMainModel(lm)
}
