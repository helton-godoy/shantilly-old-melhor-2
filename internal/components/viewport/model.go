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

	// mode controla como o conteúdo dinâmico é aplicado quando chegam
	// mensagens de saída de script:
	// - "replace" (padrão): cada execução sobrescreve o conteúdo anterior.
	// - "append": acumula histórico, adicionando novas execuções ao final.
	mode string

	// wrap controla se o conteúdo de texto "plain" deve ser reformatado para
	// caber na largura do viewport. Quando false, o conteúdo é tratado como
	// preformatado/raw (por exemplo, saída de ferramentas como gum) e é
	// repassado sem rewrap para o bubbles/viewport.
	wrap bool

	viewport       viewport.Model
	content        string
	rawContent     string
	rawContentType string
	baseContent    string // conteúdo estático inicial (ex.: mensagem de boas-vindas)
}

func New(id string, theme *tui.Theme, source *declarative.Source, mode string, wrap bool) *Model {
	vp := viewport.New(0, 0)
	if mode == "" {
		mode = "replace"
	}
	if !wrap {
		// Valor explícito false é respeitado; qualquer outro caso vira true.
	} else {
		wrap = true
	}
	m := &Model{
		id:       id,
		theme:    theme,
		source:   source,
		viewport: vp,
		mode:     strings.ToLower(mode),
		wrap:     wrap,
	}
	if source != nil && source.Type == "static" {
		m.SetContent(source.Content, source.ContentType)
		m.baseContent = m.rawContent
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
			// Aplica política de atualização baseada em mode.
			// - replace: cada nova execução sobrescreve o conteúdo anterior,
			//   preservando apenas o baseContent estático (se houver).
			// - append: acumula histórico, comportamento anterior.
			if m.mode == "replace" {
				if m.baseContent != "" {
					m.rawContent = m.baseContent
					if !strings.HasSuffix(m.rawContent, "\n\n") {
						// Garante um espaçamento mínimo entre conteúdo estático e saída dinâmica.
						m.rawContent += "\n\n"
					}
				} else {
					m.rawContent = ""
				}
			}

			m.rawContent += msg.Line
			m.rawContentType = "text"
			if m.wrap && m.viewport.Width > 0 {
				m.content = m.wrapPlainText(m.rawContent)
			} else {
				m.content = m.rawContent
			}
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
	if m.rawContentType != "markdown" && m.rawContent != "" {
		if m.wrap && m.viewport.Width > 0 {
			m.content = m.wrapPlainText(m.rawContent)
		} else {
			m.content = m.rawContent
		}
		m.viewport.SetContent(m.content)
	}
}

func (m *Model) ID() string {
	return m.id
}

// SetContent define o conteúdo inicial do viewport, com suporte opcional a markdown.
func (m *Model) SetContent(content, contentType string) {
	m.rawContent = content
	m.rawContentType = strings.ToLower(contentType)
	if m.rawContentType == "markdown" {
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
		if m.wrap && m.viewport.Width > 0 {
			m.content = m.wrapPlainText(m.rawContent)
		} else {
			m.content = m.rawContent
		}
	}

	m.viewport.SetContent(m.content)
}

func (m *Model) wrapPlainText(content string) string {
	width := m.viewport.Width
	if width <= 0 {
		return content
	}

	var wrappedLines []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		for len(line) > width {
			wrappedLines = append(wrappedLines, line[:width])
			line = line[width:]
		}
		wrappedLines = append(wrappedLines, line)
	}
	return strings.Join(wrappedLines, "\n")
}
