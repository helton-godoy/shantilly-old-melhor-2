package buttongroup

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
	tuiapi "shantilly/pkg/tui"
)

// Model implementa a interface ShantillyComponent para grupos de botões lógicos
// (ex.: ações de formulário).
// Fonte: docs/prd.md (FR6), docs/architecture.md (Seção 5.5).
type Model struct {
	id    string
	theme *tui.Theme
	items []declarative.Item

	focusedIdx int
	width      int
	height     int
}

func New(id string, theme *tui.Theme, items []declarative.Item) *Model {
	return &Model{
		id:         id,
		theme:      theme,
		items:      items,
		focusedIdx: 0,
	}
}

func (m *Model) Init() tea.Cmd { return nil }

// Update trata navegação horizontal e emite eventos de "press" ao selecionar.
func (m *Model) Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Navegação horizontal entre botões.
		case "left", "h":
			if m.focusedIdx > 0 {
				m.focusedIdx--
			}
		case "right", "l":
			if m.focusedIdx < len(m.items)-1 {
				m.focusedIdx++
			}
		// Seleção do botão focado.
		case "enter":
			if len(m.items) > 0 {
				item := m.items[m.focusedIdx]
				// Emite evento declarativo para o EventManager.
				ev := tuiapi.ShantillyEvent{
					ComponentID: m.id,
					Type:        "press", // OnHandler.Event deve usar "id:press" como chave
					Payload:     item.ID,
				}
				cmds = append(cmds, func() tea.Msg { return ev })
			}
		}
	}

	if len(cmds) == 0 {
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var buttonViews []string

	btnStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	btnFocusedStyle := btnStyle.Copy().
		BorderForeground(lipgloss.Color("63"))

	for i, it := range m.items {
		style := btnStyle
		if i == m.focusedIdx {
			style = btnFocusedStyle
		}
		label := it.Label
		if strings.TrimSpace(label) == "" {
			label = it.ID
		}
		buttonViews = append(buttonViews, style.Render(label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, buttonViews...)
}

func (m *Model) SetDimensions(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) ID() string { return m.id }
