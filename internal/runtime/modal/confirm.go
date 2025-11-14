package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"shantilly/pkg/tui"
)

// Model é um componente TUI (tea.Model) para
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
	).WithWidth(40)

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

	if m.form.State == huh.StateCompleted {
		confirmed, _ := m.form.Get("confirm").(bool)
		cmds = append(cmds, func() tea.Msg { return tui.ModalResultMsg{Confirmed: confirmed} })
	}

	if m.form.State == huh.StateAborted {
		cmds = append(cmds, func() tea.Msg { return tui.ModalResultMsg{Confirmed: false} })
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(m.form.View())
}
