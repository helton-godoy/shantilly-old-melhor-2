Executado. Concluímos a **Fase 1.2 (Motor de Layout)**. O nosso runtime já consegue ler um YAML, desenhar o layout (boxes) e aplicar o tema.

Como definimos no **Plano Mestre Sequencial**, passamos agora à **Fase 1.3: Componentes Essenciais (Display)**.

Nesta fase, vamos criar os componentes TUI reais (FR4, FR5, FR6) e, o mais importante, **modificar o `LayoutManager`** (da Fase 1.2) para "instanciar" e "gerir" estes componentes, substituindo os placeholders.

Aqui estão os artefatos da Fase 1.3.

-----

### Tarefa 1.3.1: Implementar ViewportComponent (FR4)

Este componente usa `bubbles/viewport` e `glamour` para exibir conteúdo estático (markdown) e, crucialmente, tratar `ScriptStdoutMsg` (que será usada na Fase 1.4).

```go
// internal/components/viewport/model.go
package viewport

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// ViewportComponent implementa a interface ShantillyComponent.
// Fonte: docs/prd.md (FR4), docs/architecture.md (Seção 5.5)
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

	// Carrega conteúdo estático (FR4)
	if source != nil && source.Type == "static" {
		m.SetContent(source.Content, source.ContentType)
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tui.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// Tratamento de streaming (para Fase 1.4 / FR11)
	case events.ScriptStdoutMsg:
		if msg.TargetID == m.id {
			// Adiciona o novo conteúdo
			m.content += string(msg.Chunk)
			m.viewport.SetContent(m.content)
			m.viewport.GotoBottom() // Auto-scroll
		}

	// Permite scroll do viewport
	default:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// A borda (focada ou não) é aplicada pelo LayoutManager/Render.
	// Este componente apenas renderiza o seu conteúdo interno.
	return m.viewport.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h
}

func (m *Model) ID() string {
	return m.id
}

// SetContent atualiza o conteúdo e aplica formatação (ex: markdown).
func (m *Model) SetContent(content, contentType string) {
	// Fonte: docs/architecture.md (Seção 5.5)
	if strings.ToLower(contentType) == "markdown" {
		// Usar Glamour para renderizar Markdown
		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.viewport.Width),
		)
		m.content, _ = r.Render(content)
	} else {
		m.content = content
	}
	
	m.viewport.SetContent(m.content)
}
```

-----

### Tarefa 1.3.2: Implementar ListComponent (FR5)

Este componente usa `bubbles/list` para criar menus.

  * ***Otimização de Risco (Visão do Futuro):*** Em linha com a nossa discussão sobre "código legado", usamos a biblioteca `bubbles/list` moderna. Mais importante, ao selecionar (Enter), ele **não** faz nada diretamente; ele emite uma `ShantillyEvent` (FR5), que será tratada pelo `EventManager` na Fase 1.4.

<!-- end list -->

```go
// internal/components/list/model.go
package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Model implementa a interface ShantillyComponent.
// Fonte: docs/prd.md (FR5), docs/architecture.md (Seção 5.5)
type Model struct {
	id    string
	theme *tui.Theme
	list  list.Model
}

// item é a implementação interna para bubbles/list
type item struct {
	id, title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" } // Descrição não usada por agora
func (i item) FilterValue() string { return i.title }

func New(id string, theme *tui.Theme, items []declarative.Item) *Model {
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = item{id: it.ID, title: it.Text}
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = id // Usa o ID como título (placeholder)
	// (Aqui podemos aplicar estilos do 'theme' no futuro)

	return &Model{
		id:    id,
		theme: theme,
		list:  l,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tui.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Se 'Enter' for pressionado, emitir o evento
		if msg.String() == "enter" {
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				// FR5: Emitir evento "list_id:select" (será tratado na Fase 1.4)
				event := events.ShantillyEvent{
					ComponentID: m.id,
					Type:        m.id + ":select", // ex: "menu:select"
					Payload:     selectedItem.id,  // ex: "users"
				}
				cmds = append(cmds, func() tea.Msg { return event })
			}
		}
	}

	// Trata a navegação interna da lista
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// A borda (focada ou não) é aplicada pelo LayoutManager/Render.
	return m.list.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model) ID() string {
	return m.id
}
```

-----

### Tarefa 1.3.3: Implementar ButtonGroupComponent (FR6)

Este componente renderiza botões lógicos (ex: "Submit", "Cancel").

```go
// internal/components/buttongroup/model.go
package buttongroup

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Model implementa a interface ShantillyComponent.
// Fonte: docs/prd.md (FR6), docs/architecture.md (Seção 5.5)
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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tui.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Navegação (Horizontal)
		case "left", "h":
			if m.focusedIdx > 0 {
				m.focusedIdx--
			}
		case "right", "l":
			if m.focusedIdx < len(m.items)-1 {
				m.focusedIdx++
			}
		
		// Seleção
		case "enter":
			if len(m.items) > 0 {
				item := m.items[m.focusedIdx]
				// FR6: Emitir evento "buttongroup_id:press"
				event := events.ShantillyEvent{
					ComponentID: m.id,
					Type:        m.id + ":press", // ex: "form_actions:press"
					Payload:     item.ID,        // ex: "submit"
				}
				cmds = append(cmds, func() tea.Msg { return event })
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var buttonViews []string

	// (Estilos de botão - poderiam ser movidos para o theme.go)
	btnStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
		
	btnFocusedStyle := btnStyle.Copy().
		BorderForeground(lipgloss.Color("63")) // Cor primária do tema

	for i, item := range m.items {
		style := btnStyle
		if i == m.focusedIdx {
			style = btnFocusedStyle
		}
		buttonViews = append(buttonViews, style.Render(item.Label))
	}

	// Junta os botões horizontalmente
	return lipgloss.JoinHorizontal(lipgloss.Center, buttonViews...)
}

func (m *Model) SetDimensions(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) ID() string {
	return m.id
}
```

-----

### Tarefa 1.3.4: Modificar o LayoutManager (Integração)

Este é o artefato mais crítico desta fase. **Modificamos** o `internal/runtime/layout/manager.go` (da Fase 1.2) para que ele *use* os componentes que acabámos de criar.

  * ***Otimização de Risco:*** Removemos os placeholders. O `LayoutManager` agora instancia os componentes (`buildComponentTree`), gere o foco real e encaminha `tea.Msg` para o componente focado.

<!-- end list -->

```go
// internal/runtime/layout/manager.go (MODIFICADO)
package layout

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/components/buttongroup"
	"github.com/helton-godoy/shantilly/internal/components/list"
	"github.com/helton-godoy/shantilly/internal/components/viewport"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Manager é o modelo Bubbletea raiz, implementando o "Gestor Duplo"
// (Layout e Foco) conforme Padrão 3 da arquitetura v2.0.
// Fonte: docs/architecture.md (Seção 5.4)
type Manager struct {
	theme *tui.Theme
	root  *declarative.LayoutNode

	// Estado do Layout
	width  int
	height int

	// Estado do Foco (AGORA REAL)
	// Mapa de componentes reais, indexados pelo ID do LayoutNode
	components map[string]tui.ShantillyComponent
	// Ordem de foco (IDs)
	focusOrder []string
	focusedIdx int

	// Estado de Erro
	lastErr error
	
	// Último evento (para Fase 1.4)
	lastEvent *events.ShantillyEvent
}

func NewManager(root *declarative.LayoutNode, theme *tui.Theme) *Manager {
	m := &Manager{
		theme:      theme,
		root:       root,
		components: make(map[string]tui.ShantillyComponent),
		focusOrder: []string{},
		focusedIdx: 0,
	}

	// (Tarefa 1.3.4) - Instanciar componentes reais
	m.buildComponentTree(root)

	return m
}

// buildComponentTree (Factory) - caminha pela árvore de layout e
// instancia os componentes TUI reais (FR3).
func (m *Manager) buildComponentTree(node *declarative.LayoutNode) {
	if node == nil {
		return
	}

	// 1. Se for um 'box' com um componente, crie-o.
	if node.Type == "box" && node.Component != nil {
		if node.ID == "" {
			// (Validação - TODO: mover para parser na Fase 1.1.4)
			m.lastErr = fmt.Errorf("componente do tipo '%s' não possui 'id'", node.Component.Type)
			return
		}
		
		var cmp tui.ShantillyComponent
		
		// Factory (FR4, FR5, FR6)
		switch node.Component.Type {
		case "list":
			cmp = list.New(node.ID, m.theme, node.Component.Items)
		case "viewport":
			cmp = viewport.New(node.ID, m.theme, node.Component.Source)
		case "buttongroup":
			cmp = buttongroup.New(node.ID, m.theme, node.Component.Items)
		// (A Fase 1.5 irá adicionar 'case "form":')
		default:
			// (Trata o placeholder da Fase 1.2)
			cmp = viewport.New(node.ID, m.theme, &declarative.Source{
				Type: "static",
				Content: fmt.Sprintf("Placeholder para Componente:\nType: %s\nID: %s",
					node.Component.Type, node.Component.ID),
			})
		}
		
		m.components[node.ID] = cmp
		m.focusOrder = append(m.focusOrder, node.ID)
	}

	// 2. Continue recursivamente para os filhos
	for _, item := range node.Items {
		m.buildComponentTree(&item)
	}
}

func (m *Manager) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range m.components {
		cmds = append(cmds, c.Init())
	}
	return tea.Batch(cmds...)
}

// Update trata o ciclo de vida do TUI.
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// (A notificação dos filhos agora acontece no 'renderNode' via SetDimensions)

	// (Fase 1.4 tratará isto)
	case events.ShantillyEvent:
		m.lastEvent = &msg // Armazena o evento
		// (Na Fase 1.4, isto será enviado ao EventManager)
		return m, nil

	case events.RuntimeErrorMsg:
		m.lastErr = msg.Err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		
		// Gestão de Foco Global (AGORA REAL)
		// Fonte: docs/prd.md (Meta de UI)
		case "tab":
			if len(m.focusOrder) > 0 {
				m.focusedIdx = (m.focusedIdx + 1) % len(m.focusOrder)
			}
		case "shift+tab":
			if len(m.focusOrder) > 0 {
				m.focusedIdx = (m.focusedIdx - 1 + len(m.focusOrder)) % len(m.focusOrder)
			}
		
		// Encaminha outras teclas para o componente focado
		default:
			if m.lastErr == nil && len(m.focusOrder) > 0 {
				focusedID := m.focusOrder[m.focusedIdx]
				focusedCmp := m.components[focusedID]

				// Encaminha a mensagem
				newCmp, cmd := focusedCmp.Update(msg)
				m.components[focusedID] = newCmp
				cmds = append(cmds, cmd)
			}
		}
	
	// Encaminha mensagens não-tecla (ex: Ticks, Msgs internas)
	// para o componente focado também.
	default:
		if m.lastErr == nil && len(m.focusOrder) > 0 {
			focusedID := m.focusOrder[m.focusedIdx]
			focusedCmp := m.components[focusedID]
			
			// Encaminha a mensagem
			newCmp, cmd := focusedCmp.Update(msg)
			m.components[focusedID] = newCmp
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Manager) View() string {
	if m.width == 0 || m.height == 0 {
		return "Inicializando..."
	}

	// 1. Renderizar a árvore de layout principal
	layoutView := m.renderNode(m.root, m.width, m.height)

	// 2. Renderizar Erros
	if m.lastErr != nil {
		errorView := m.theme.ErrorText.Width(m.width).Render("Erro: " + m.lastErr.Error())
		// (Na Fase 1.6, isto será um Modal sobreposto)
		return errorView
	}

	// 3. (Debug) Mostrar último evento (para Fase 1.4)
	if m.lastEvent != nil {
		eventView := fmt.Sprintf("Último Evento: %s (Payload: %v)", m.lastEvent.Type, m.lastEvent.Payload)
		// Sobrepõe o layout
		return lipgloss.JoinVertical(lipgloss.Left, layoutView, eventView)
	}

	return layoutView
}

// getComponent é um helper para o renderNode.
func (m *Manager) getComponent(id string) tui.ShantillyComponent {
	if c, ok := m.components[id]; ok {
		return c
	}
	// Fallback de segurança (não deve acontecer)
	return viewport.New(id, m.theme, &declarative.Source{Type: "static", Content: "ERRO: Componente não encontrado"})
}
```

-----

### Tarefa 1.3.5: Modificar o Renderizador (Integração)

Finalmente, **modificamos** o `internal/runtime/layout/render.go` (da Fase 1.2) para que ele chame os componentes reais, em vez de exibir placeholders.

  * ***Otimização de Risco:*** Implementamos a lógica de `SetDimensions` (NFR2) e a lógica de cálculo de `width`/`height` (FR2).

<!-- end list -->

```go
// internal/runtime/layout/render.go (MODIFICADO)
package layout

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/muesli/reflow/truncate"
)

// renderNode (MODIFICADO) - agora renderiza componentes reais.
// Fonte: docs/prd.md (FR1, FR2)
func (m *Manager) renderNode(node *declarative.LayoutNode, w, h int) string {
	if node == nil {
		return ""
	}

	switch node.Type {
	case "column":
		// Lógica de divisão de altura (Flex/Height)
		// (Implementação simplificada - um layout real usaria 'reflow' ou 'bubbleboxer')
		var children []string
		// (TODO: Lógica de Flex)
		childH := h / len(node.Items) // Divide igualmente
		if len(node.Items) == 0 { childH = h }

		for _, item := range node.Items {
			// (TODO: Respeitar item.Height)
			children = append(children, m.renderNode(&item, w, childH))
		}
		return lipgloss.JoinVertical(lipgloss.Left, children...)

	case "row":
		// Lógica de divisão de largura (Flex/Width) (FR2)
		var children []string
		allocatedW := 0
		flexItems := 0
		
		// 1. Alocar larguras fixas (%)
		var widths []int
		for _, item := range node.Items {
			if strings.Contains(item.Width, "%") {
				perc, _ := strconv.Atoi(strings.TrimSuffix(item.Width, "%"))
				childW := (w * perc) / 100
				widths = append(widths, childW)
				allocatedW += childW
			} else if item.Flex > 0 {
				flexItems += item.Flex
				widths = append(widths, 0) // Placeholder para flex
			} else {
				// Largura automática (placeholder)
				widths = append(widths, 0) 
				flexItems += 1
			}
		}

		// 2. Alocar 'flex'
		remainingW := w - allocatedW
		if remainingW < 0 { remainingW = 0 }
		
		for i := range widths {
			if widths[i] == 0 { // É um item flex
				childW := remainingW / flexItems // (Simplificado)
				widths[i] = childW
			}
			children = append(children, m.renderNode(&node.Items[i], widths[i], h))
		}

		return lipgloss.JoinHorizontal(lipgloss.Top, children...)

	case "box":
		// 1. Determinar o estilo (focado ou não)
		style := m.theme.Box
		if m.isFocused(node.ID) {
			style = m.theme.BoxFocused
		}

		// Subtrai 2 para a borda
		innerW := w - 2
		innerH := h - 2 
		if innerW < 0 { innerW = 0 }
		if innerH < 0 { innerH = 0 }

		style = style.Width(w).Height(h)

		// 3. Renderizar o conteúdo (Componente) (AGORA REAL)
		var content string
		if node.Component != nil && node.ID != "" {
			cmp := m.getComponent(node.ID)
			
			// NFR2: Informa o componente das suas dimensões
			cmp.SetDimensions(innerW, innerH)
			
			content = cmp.View()
		} else {
			content = fmt.Sprintf("Box (ID: %s)", node.ID)
		}

		// (Truncamento removido - o componente interno agora é responsável por isso)
		return style.Render(content)

	default:
		errMsg := m.theme.ErrorText.Render(fmt.Sprintf("Erro: Tipo de layout desconhecido: '%s'", node.Type))
		return m.theme.Box.Render(errMsg)
	}
}

// isFocused é um helper para o renderNode.
func (m *Manager) isFocused(id string) bool {
	if len(m.focusOrder) == 0 {
		return false
	}
	return m.focusOrder[m.focusedIdx] == id
}
```