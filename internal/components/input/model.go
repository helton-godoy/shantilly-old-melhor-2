package input

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/tui"
	tuiapi "shantilly/pkg/tui"
)

// Model implementa a interface ShantillyComponent para campos de entrada de texto.
// Ele emite um evento "submit" quando o usuário pressiona Enter.
type Model struct {
	id    string
	theme *tui.Theme

	input textinput.Model
}

func New(id string, theme *tui.Theme, label, placeholder, initial string, secret bool) *Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(initial)
	// Estilos mínimos; o tema pode ser integrado em waves futuras.
	if strings.TrimSpace(label) != "" {
		ti.Prompt = label + " "
	} else {
		ti.Prompt = "> "
	}
	if secret {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	// Garante que o campo já comece focado, exibindo cursor e aceitando edição
	// imediatamente ao abrir o runtime.
	ti.Focus()

	return &Model{
		id:    id,
		theme: theme,
		input: ti,
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update trata edição de texto e emite evento de submit ao pressionar Enter.
func (m *Model) Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			value := m.input.Value()
			// Emite evento declarativo para o EventManager.
			ev := tuiapi.ShantillyEvent{
				ComponentID: m.id,
				Type:        "submit",
				Payload: map[string]interface{}{
					"value": value,
				},
			}
			return m, func() tea.Msg { return ev }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.input.View()
}

func (m *Model) SetDimensions(w, h int) {
	// textinput se ajusta bem só com largura; altura é irrelevante aqui.
	m.input.Width = w
	_ = h
}

func (m *Model) ID() string { return m.id }
