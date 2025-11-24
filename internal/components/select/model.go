package selectcmp

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
	tuiapi "shantilly/pkg/tui"
)

// Model implementa a interface ShantillyComponent para um seletor de opção única.
// É semelhante ao componente de lista, mas com semântica de "select" genérica
// e evento padronizado de mudança (change).
type Model struct {
	id    string
	theme *tui.Theme
	list  list.Model
}

// item é a implementação interna para bubbles/list.
type item struct {
	id, title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

// New cria um componente de select a partir de itens declarativos.
func New(id string, theme *tui.Theme, items []declarative.Item) *Model {
	listItems := make([]list.Item, len(items))
	for idx, it := range items {
		listItems[idx] = item{id: it.ID, title: it.Text}
	}

	// Dimensões mínimas; LayoutManager ajusta depois via SetDimensions.
	l := list.New(listItems, list.NewDefaultDelegate(), 40, 10)
	l.Title = id

	return &Model{
		id:    id,
		theme: theme,
		list:  l,
	}
}

func (m *Model) Init() tea.Cmd { return nil }

// Update trata navegação e emite evento de mudança ao confirmar a seleção.
func (m *Model) Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				payload := map[string]interface{}{
					"id":   selectedItem.id,
					"text": selectedItem.title,
				}
				ev := tuiapi.ShantillyEvent{
					ComponentID: m.id,
					Type:        "change",
					Payload:     payload,
				}
				cmds = append(cmds, func() tea.Msg { return ev })
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if len(cmds) == 0 {
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return "Selecione uma opção:\n" + m.list.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model) ID() string { return m.id }
