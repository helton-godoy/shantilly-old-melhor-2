package button

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"shantilly/internal/tui"
	tuiapi "shantilly/pkg/tui"
)

// Props representa as propriedades configuráveis do button
type Props struct {
	Text     string `yaml:"text"`
	Icon     string `yaml:"icon"`
	Style    string `yaml:"style"`
	Disabled bool   `yaml:"disabled"`
}

// Model implementa a interface ShantillyComponent para buttons.
type Model struct {
	id            string
	theme         *tui.Theme
	props         Props
	baseStyle     lipgloss.Style
	focusedStyle  lipgloss.Style
	disabledStyle lipgloss.Style
	focused       bool
	renderCount   int
}

// New cria um componente button a partir de props declarativas
func New(id string, theme *tui.Theme, props Props) *Model {
	m := &Model{
		id:          id,
		theme:       theme,
		props:       props,
		focused:     false,
		renderCount: 0,
	}
	m.setupStyles()
	return m
}

func (m *Model) Init() tea.Cmd { return nil }

// Update trata eventos de teclado e emite eventos de click
func (m *Model) Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.props.Disabled {
			return m, nil
		}

		// Enter ou espaço para ativar o button
		switch msg.Type {
		case tea.KeyEnter, tea.KeySpace:
			// Emite evento de click
			eventType := "click_deploy_now"

			ev := tuiapi.ShantillyEvent{
				ComponentID: m.id,
				Type:        eventType,
				Payload:     m.props.Text,
			}
			cmds = append(cmds, func() tea.Msg { return ev })
		}
	}

	return m, tea.Batch(cmds...)
}

// View implementa a interface ShantillyComponent
func (m *Model) View() string {
	if m.props.Text == "" {
		return "[BUTTON: no text]"
	}

	result := m.renderNormal()
	if m.focused {
		return m.focusedStyle.Render(result)
	}

	return m.baseStyle.Render(result)
}

// renderNormal renderiza o button no estado normal
func (m *Model) renderNormal() string {
	text := m.props.Text
	if m.props.Icon != "" {
		text = fmt.Sprintf("%s %s", m.props.Icon, text)
	}

	// Adiciona indicador visual de foco
	if m.focused {
		text = fmt.Sprintf("[▶ %s ◀]", text)
	} else {
		text = fmt.Sprintf("[ %s ]", text)
	}

	return text
}

// Focus implementa a interface para controle de foco
func (m *Model) Focus(focused bool) {
	m.focused = focused
	m.renderCount++
}

func (m *Model) SetDimensions(w, h int) {
	// Button tem controle próprio sobre suas dimensões via estilos
}

func (m *Model) ID() string { return m.id }

// setupStyles configura os estilos do button baseado no theme
func (m *Model) setupStyles() {
	// Estilo base - cores super contrastantes
	m.baseStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("196")).
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Width(15)

	// Estilo quando focado
	m.focusedStyle = m.baseStyle.Copy().
		Background(lipgloss.Color("226")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	// Estilo quando desabilitado
	m.disabledStyle = m.baseStyle.Copy().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("245"))
}
