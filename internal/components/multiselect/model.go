package multiselect

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
	tuiapi "shantilly/pkg/tui"
)

// Model implementa a interface ShantillyComponent para seleção múltipla.
// Baseado em bubbles/list, permitindo marcar/desmarcar itens e confirmar
// a seleção com Enter.
type Model struct {
	id    string
	theme *tui.Theme
	list  list.Model
}

// item é a implementação interna para bubbles/list, com flag de seleção.
type item struct {
	id       string
	title    string
	selected bool
}

func (i *item) Title() string {
	if i.selected {
		return "[x] " + i.title
	}
	return "[ ] " + i.title
}

func (i *item) Description() string { return "" }
func (i *item) FilterValue() string { return i.title }

// New cria um componente de multiselect a partir de itens declarativos.
func New(id string, theme *tui.Theme, items []declarative.Item) *Model {
	listItems := make([]list.Item, len(items))
	for idx, it := range items {
		listItems[idx] = &item{id: it.ID, title: it.Text, selected: false}
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

// Update trata navegação, toggle com espaço e confirma com Enter.
func (m *Model) Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		// Espaço: alterna seleção do item atual.
		if key == " " {
			if selectedItem, ok := m.list.SelectedItem().(*item); ok {
				selectedItem.selected = !selectedItem.selected
			}
		}
		// Enter: confirma seleção atual e emite evento change com todos selecionados.
		if key == "enter" {
			var (
				ids   []string
				items []map[string]string
			)
			for _, li := range m.list.Items() {
				if it, ok := li.(*item); ok && it.selected {
					ids = append(ids, it.id)
					items = append(items, map[string]string{
						"id":   it.id,
						"text": it.title,
					})
				}
			}
			payload := map[string]interface{}{
				"ids":   ids,
				"items": items,
			}
			ev := tuiapi.ShantillyEvent{
				ComponentID: m.id,
				Type:        "change",
				Payload:     payload,
			}
			cmds = append(cmds, func() tea.Msg { return ev })
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
	return "Selecione uma ou mais opções (Espaço = marcar, Enter = confirmar):\n" + m.list.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model) ID() string { return m.id }
