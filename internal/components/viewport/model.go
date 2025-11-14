package viewport

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
	tuiapi "shantilly/pkg/tui"
)

// Model é o componente responsável por exibir conteúdo de texto (incluindo
// saída de scripts) em uma área com scroll.
// Ele implementa o contrato tui.ShantillyComponent.

type Model struct {
	id     string
	theme  *tui.Theme
	source *declarative.Source

	viewport viewport.Model
	content  string
}

func New(id string, theme *tui.Theme, source *declarative.Source) *Model {
	vp := viewport.New(0, 0)
	m := &Model{
		id:       id,
		theme:    theme,
		source:   source,
		viewport: vp,
	}
	if source != nil && source.Type == "static" {
		m.SetContent(source.Content, source.ContentType)
	}
	return m
}

func (m *Model) Init() tea.Cmd { return nil }

// Update trata mensagens de streaming e de navegação do viewport.
func (m *Model) Update(msg tea.Msg) (tuiapi.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// Streaming de saída de scripts (linha a linha).
	case tuiapi.ScriptStdoutMsg:
		if msg.TargetID == m.id {
			// Para cada nova execução, substituímos o conteúdo anterior pelo
			// resultado mais recente, mantendo o viewport focado apenas na
			// última run (útil para debug e para evitar ruído acumulado).
			m.content = msg.Line
			m.viewport.SetContent(m.content)
			m.viewport.GotoBottom()
		}

	default:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.viewport.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h
}

func (m *Model) ID() string {
	return m.id
}

// SetContent define o conteúdo inicial do viewport, com suporte opcional a markdown.
func (m *Model) SetContent(content, contentType string) {
	if strings.ToLower(contentType) == "markdown" {
		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.viewport.Width),
		)
		if err == nil {
			if rendered, err2 := r.Render(content); err2 == nil {
				m.content = rendered
			} else {
				m.content = content
			}
		} else {
			m.content = content
		}
	} else {
		m.content = content
	}

	m.viewport.SetContent(m.content)
}
