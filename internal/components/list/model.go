package list

import (
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
	listItems := make([]list.Item, len(items))
	for i, it := range items {
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
		// Enter dispara evento declarativo para o EventManager.
		if msg.String() == "enter" {
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				// Para qualquer item, geramos um tipo de evento padronizado
				// "select_<id>", permitindo que o YAML use event: "menu:select_<id>".
				eventType := "select"
				if selectedItem.id != "" {
					eventType = "select_" + selectedItem.id
				}
				ev := tuiapi.ShantillyEvent{
					ComponentID: m.id,
					Type:        eventType,
					Payload:     selectedItem.id, // ex.: "users"
				}
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
	// Adiciona um cabeçalho simples para tornar a área de menu mais
	// evidente na renderização do layout.
	return "Menu principal:\n" + m.list.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model) ID() string { return m.id }
