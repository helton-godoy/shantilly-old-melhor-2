package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
	tuiapi "shantilly/pkg/tui"
)

// Model implementa a interface ShantillyComponent para listas/menu.
// Fonte: docs/prd.md (FR5), docs/architecture.md (Seção 5.5).
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

// New cria um componente de lista a partir de itens declarativos.
func New(id string, theme *tui.Theme, items []declarative.Item) *Model {
	fmt.Printf("[list] New: id=%s items=%d\n", id, len(items))
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		fmt.Printf("[list]   item[%d]: id=%s text=%s\n", i, it.ID, it.Text)
		listItems[i] = item{id: it.ID, title: it.Text}
	}

	// Inicializa com dimensões mínimas para garantir que a lista
	// tenha área visível mesmo antes de receber SetDimensions
	// do LayoutManager. O layout irá ajustar depois via SetDimensions.
	l := list.New(listItems, list.NewDefaultDelegate(), 40, 10)
	l.Title = id // placeholder; estilos podem vir do theme

	return &Model{
		id:    id,
		theme: theme,
		list:  l,
	}
}

func (m *Model) Init() tea.Cmd { return nil }

// Update trata navegação e emite eventos selecionados via ShantillyEvent.
func (m *Model) Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		fmt.Printf("[list] KeyMsg: id=%s key=%q type=%v\n", m.id, msg.String(), msg.Type)
		// Enter dispara evento declarativo para o EventManager.
		if msg.String() == "enter" {
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				ev := tuiapi.ShantillyEvent{
					ComponentID: m.id,
					Type:        "select",        // o EventManager combina "id:select" via OnHandler.Event
					Payload:     selectedItem.id, // ex.: "users"
				}
				fmt.Printf("[list] Emitting ShantillyEvent: componentID=%s type=%s payload=%s\n", ev.ComponentID, ev.Type, selectedItem.id)
				cmds = append(cmds, func() tea.Msg { return ev })
			}
		}
	}

	// Navegação interna da lista.
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.list.View()
}

func (m *Model) SetDimensions(w, h int) {
	fmt.Printf("[list] SetDimensions: id=%s w=%d h=%d\n", m.id, w, h)
	m.list.SetSize(w, h)
}

func (m *Model) ID() string { return m.id }
