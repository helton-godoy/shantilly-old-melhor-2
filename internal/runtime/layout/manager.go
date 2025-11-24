package layout

// Epic: 1 - Runtime TUI Declarativo — Fundação do Runtime (v2.0)
// Bloco: E1.1 - LayoutManager (Layout Engine hierárquico declarativo)
// Referências normativas:
// - Docs Arquitetura:
//   - docs/architecture/components.md#3-layoutmanager--e11
//   - docs/architecture/data-models.md#2-layoutnode-e11--layoutmanager
//   - docs/architecture/high-level-architecture.md#high-level-overview
// - Story Map:
//   - docs/stories/index-epic-1-runtime-tui.story-map.md#e11--layout-engine-layout-hierarquico-declarativo
// - QA Matrix:
//   - docs/qa/matrix-epic-1-runtime-tui-coverage.md#bloco-e11--layout-engine-layout-hierarquico-declarativo
//
// Invariantes (enforced pelo Orchestrator):
// - Consome exclusivamente modelos v2.0 (AppConfig/LayoutNode/Component) do pacote pkg/declarative.
// - Não executa run:, não avalia on:, não dispara automação.
// - Não conhece FormConfig nem internal/tui como contrato ativo (apenas componentes via ShantillyComponent).
// - Opera como root bubbletea.Model responsável por layout, foco global e integração futura com Modal Stack.

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	boxer "github.com/treilik/bubbleboxer"

	"shantilly/internal/runtime/event"
	"shantilly/internal/runtime/runner"
	"shantilly/pkg/tui"
)

// ShantillyComponent é um alias para o contrato definido em pkg/tui.
type ShantillyComponent = tui.ShantillyComponent

// LayoutNodeRef é uma view mínima sobre o modelo declarativo de layout.
// Em implementação completa, este tipo será abastecido a partir de pkg/declarative.
type LayoutNodeRef struct {
	ID             string
	Type           string // "column" | "row" | "box"
	Width          *int
	Height         *int
	Flex           *int
	Padding        *int
	Border         bool
	ComputedWidth  int
	ComputedHeight int
	Items          []LayoutNodeRef
	ComponentID    string // válido apenas para type: box
}

// ComponentRegistry expõe os componentes instanciáveis pelo LayoutManager.
// Implementação detalhada (factory) é responsabilidade de bmad-dev/bmad-master.
type ComponentRegistry interface {
	// Resolve retorna um componente concreto para um ID declarado em AppConfig.Components.
	// Retorno nil indica configuração inválida ou componente não encontrado (tratado em camada superior).
	Resolve(id string) ShantillyComponent
}

// Manager é o LayoutManager raiz (bubbletea.Model) responsável por:
// - Interpretar LayoutNodeRef (derivado de AppConfig.Layout).
// - Construir estrutura interna de layout e calcular bounds.
// - Instanciar e manter ShantillyComponent(s).
// - Gerenciar foco global e roteamento básico de mensagens.
// - Delegar integração com Modal Stack/Events a módulos dedicados (futuras Waves).
type Manager struct {
	layout      LayoutNodeRef
	registry    ComponentRegistry
	components  map[string]ShantillyComponent
	focusedID   string
	width       int
	height      int
	initialized bool

	// Integração com EventManager v2.0 (E1.3) e fluxo JIT (modais).
	eventManager *event.Manager

	// Integração com ScriptRunner (E1.4) para execução de RunAction.
	scriptRunner *runner.ScriptRunner
}

// New cria um LayoutManager a partir de um LayoutNodeRef e um registry de componentes.
// Notas de governança:
// - Chamadores devem garantir que layout/IDs já estejam validados conforme pkg/declarative.
// - Este construtor não executa scripts nem acessa on:/run:.
func New(root LayoutNodeRef, registry ComponentRegistry) *Manager {
	return &Manager{
		layout:     root,
		registry:   registry,
		components: make(map[string]ShantillyComponent),
	}
}

// init configura parâmetros globais do BubbleBoxer utilizados pelo LayoutManager.
func init() {
	// Mantemos os separadores padrão do BubbleBoxer (│ e ─). O LayoutManager
	// passa a tratar bordas dos componentes de forma mais neutra quando o
	// layout é controlado pelo Boxer.
	boxer.HorizontalSeparator = "│"
	boxer.VerticalSeparator = "─"
}

// boxerLeafModel é um tea.Model mínimo usado apenas pelo BubbleBoxer para
// compor layout. Ele delega View() para o ShantillyComponent correspondente
// no Manager, sem interferir no ciclo de Update/foco do runtime.
type boxerLeafModel struct {
	componentID string
	manager     *Manager
}

func (m boxerLeafModel) Init() tea.Cmd { return nil }

func (m boxerLeafModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m boxerLeafModel) View() string {
	if m.manager == nil || m.componentID == "" {
		return ""
	}
	if c, ok := m.manager.components[m.componentID]; ok && c != nil {
		return c.View()
	}
	return ""
}

// buildBubbleBoxerLayout converte o LayoutNodeRef raiz em uma árvore
// bubbleboxer.Boxer. Em vez de recalcular o layout, ele reutiliza os valores
// ComputedWidth/ComputedHeight já definidos por applyDimensions, gerando
// SizeFunc fixas que refletem a distribuição já calculada.
func (m *Manager) buildBubbleBoxerLayout() boxer.Boxer {
	box := boxer.Boxer{}
	rootNode := m.buildBoxerNode(&box, m.layout)
	box.LayoutTree = rootNode
	return box
}

// buildBoxerNode mapeia recursivamente um LayoutNodeRef para um boxer.Node.
//   - Para type: box, cria um Leaf associado ao ComponentID.
//   - Para type: row, monta um Node horizontal com SizeFunc fixa baseada em
//     ComputedWidth dos filhos.
//   - Para type: column, monta um Node vertical com SizeFunc fixa baseada em
//     ComputedHeight dos filhos.
func (m *Manager) buildBoxerNode(box *boxer.Boxer, node LayoutNodeRef) boxer.Node {
	switch node.Type {
	case "box":
		if node.ComponentID == "" {
			// Leaf vazio: apenas um placeholder sem conteúdo.
			leaf, _ := box.CreateLeaf(node.ID, boxerLeafModel{componentID: "", manager: m})
			return leaf
		}
		leafModel := boxerLeafModel{componentID: node.ComponentID, manager: m}
		leaf, _ := box.CreateLeaf(node.ComponentID, leafModel)
		return leaf
	case "row":
		children := make([]boxer.Node, len(node.Items))
		ratios := make([]int, len(node.Items))
		for i := range node.Items {
			children[i] = m.buildBoxerNode(box, node.Items[i])
			r := 1
			if node.Items[i].Flex != nil && *node.Items[i].Flex > 0 {
				r = *node.Items[i].Flex
			}
			ratios[i] = r
		}
		return boxer.Node{
			Children:        children,
			VerticalStacked: false,
			SizeFunc:        flexHorizontalSizeFunc(ratios),
		}
	case "column":
		children := make([]boxer.Node, len(node.Items))
		ratios := make([]int, len(node.Items))
		for i := range node.Items {
			children[i] = m.buildBoxerNode(box, node.Items[i])
			r := 1
			if node.Items[i].Flex != nil && *node.Items[i].Flex > 0 {
				r = *node.Items[i].Flex
			}
			ratios[i] = r
		}
		return boxer.Node{
			Children:        children,
			VerticalStacked: true,
			SizeFunc:        flexVerticalSizeFunc(ratios),
		}
	default:
		// Nó desconhecido: retorna um leaf vazio para evitar panics.
		leaf, _ := box.CreateLeaf(node.ID, boxerLeafModel{componentID: "", manager: m})
		return leaf
	}
}

// flexHorizontalSizeFunc cria uma SizeFunc que distribui a largura entre os
// filhos proporcionalmente aos ratios informados, seguindo o padrão usado no
// spike bubbleboxer-layout-spike.
func flexHorizontalSizeFunc(ratios []int) func(node boxer.Node, width int) []int {
	return func(node boxer.Node, width int) []int {
		count := len(node.Children)
		if count == 0 {
			return nil
		}
		// Se o tamanho da lista de ratios não bater, fazemos divisão igualitária.
		if len(ratios) != count {
			res := make([]int, count)
			base := width / count
			for i := range res {
				res[i] = base
			}
			rest := width - base*count
			for i := 0; i < rest && i < count; i++ {
				res[i]++
			}
			return res
		}

		// Normaliza ratios, tratando zeros como 1.
		norm := make([]int, count)
		total := 0
		for i, r := range ratios {
			if r <= 0 {
				r = 1
			}
			norm[i] = r
			total += r
		}
		if total <= 0 {
			res := make([]int, count)
			base := width / count
			for i := range res {
				res[i] = base
			}
			return res
		}

		res := make([]int, count)
		remaining := width
		remainingTotal := total
		for i, r := range norm {
			if i == count-1 {
				res[i] = remaining
				break
			}
			v := remaining * r / remainingTotal
			res[i] = v
			remaining -= v
			remainingTotal -= r
		}
		return res
	}
}

// flexVerticalSizeFunc cria uma SizeFunc que distribui a altura entre os
// filhos proporcionalmente aos ratios informados.
func flexVerticalSizeFunc(ratios []int) func(node boxer.Node, height int) []int {
	return func(node boxer.Node, height int) []int {
		count := len(node.Children)
		if count == 0 {
			return nil
		}
		if len(ratios) != count {
			res := make([]int, count)
			base := height / count
			for i := range res {
				res[i] = base
			}
			rest := height - base*count
			for i := 0; i < rest && i < count; i++ {
				res[i]++
			}
			return res
		}

		norm := make([]int, count)
		total := 0
		for i, r := range ratios {
			if r <= 0 {
				r = 1
			}
			norm[i] = r
			total += r
		}
		if total <= 0 {
			res := make([]int, count)
			base := height / count
			for i := range res {
				res[i] = base
			}
			return res
		}

		res := make([]int, count)
		remaining := height
		remainingTotal := total
		for i, r := range norm {
			if i == count-1 {
				res[i] = remaining
				break
			}
			v := remaining * r / remainingTotal
			res[i] = v
			remaining -= v
			remainingTotal -= r
		}
		return res
	}
}

// Init inicializa componentes conforme o layout.
// Implementação inicial mínima: prepara estrutura interna.
// Expansão detalhada (instanciação e cmds) é responsabilidade de bmad-dev/bmad-master.
func (m *Manager) Init() tea.Cmd {
	if m.initialized {
		return nil
	}

	// Se ainda não recebemos WindowSizeMsg, aplicamos dimensões padrão
	// para evitar que componentes sejam inicializados com largura/altura zero.
	if m.width == 0 || m.height == 0 {
		m.width = 80
		m.height = 24
	}
	cmd := m.buildComponents(&m.layout)
	m.applyDimensions()
	m.initialized = true
	return cmd
}

// SetEventManager injeta o EventManager responsável por resolver on:/RunAction
// a partir de tui.ShantillyEvent.
func (m *Manager) SetEventManager(em *event.Manager) {
	m.eventManager = em
}

// SetScriptRunner injeta o ScriptRunner responsável por executar RunAction
// declarativos disparados pelo EventManager.
func (m *Manager) SetScriptRunner(sr *runner.ScriptRunner) {
	m.scriptRunner = sr
}

// Update implementa o loop TEA raiz para o LayoutManager.
// Regras nesta wave:
// - Tratar mensagens de resize para atualizar width/height.
// - Roteamento básico de msgs para componente focado (quando existir).
// - NÃO executar automação nem avaliar on:/run: (delegado ao EventManager).
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.applyDimensions()
		return m, nil

	// Pedido de execução de script vindo do EventManager.
	case tui.RunScriptRequestMsg:
		if m.scriptRunner != nil {
			cmd = runner.HandleRunRequest(m.scriptRunner, msg)
			cmds = append(cmds, cmd)
		}

	// Saída de scripts destinada a um componente específico (ex.: viewport).
	case tui.ScriptStdoutMsg:
		if c, ok := m.components[msg.TargetID]; ok && c != nil {
			var ccmd tea.Cmd
			next, ccmd := c.Update(msg)
			m.components[msg.TargetID] = next
			if ccmd != nil {
				cmds = append(cmds, ccmd)
			}
		}

	// Eventos declarativos emitidos por componentes (via pkg/tui).
	case tui.ShantillyEvent:
		if m.eventManager != nil {
			cmd = m.eventManager.ProcessEvent(msg)
			cmds = append(cmds, cmd)
		}

	// Resultado de modal JIT vindo do MainModel/modal.
	case tui.ModalResultMsg:
		if m.eventManager != nil {
			cmd = m.eventManager.HandleModalResult(msg)
			cmds = append(cmds, cmd)
		}

	// Mensagem para exibir modal: devolvemos como tea.Cmd para o MainModel
	// interceptar e renderizar o modal (overlay).
	case tui.ShowModalMsg:
		return m, func() tea.Msg { return msg }
	}

	// Roteamento mínimo para componente focado (se houver).
	if m.focusedID != "" {
		if c, ok := m.components[m.focusedID]; ok && c != nil {
			next, ccmd := c.Update(msg)
			m.components[m.focusedID] = next
			cmds = append(cmds, ccmd)
		}
	}

	if len(cmds) == 0 {
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// View monta a string final com base no layout calculado.
// Nesta fase, fornecemos apenas um placeholder neutro alinhado à governança:
// nenhuma lógica de automação, apenas render estrutural.
func (m *Manager) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	return m.renderNode(m.layout)
}

func (m *Manager) buildComponents(node *LayoutNodeRef) tea.Cmd {
	if node == nil {
		return nil
	}

	var cmds []tea.Cmd

	if node.Type == "box" && node.ComponentID != "" {
		c := m.registry.Resolve(node.ComponentID)
		if c != nil {
			m.components[node.ComponentID] = c
			if m.focusedID == "" {
				m.focusedID = node.ComponentID
			}
			if initCmd := c.Init(); initCmd != nil {
				cmds = append(cmds, initCmd)
			}
		}
	}

	for i := range node.Items {
		if cmd := m.buildComponents(&node.Items[i]); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) == 0 {
		return nil
	}

	return tea.Batch(cmds...)
}

func (m *Manager) applyDimensions() {
	if m.width <= 0 || m.height <= 0 {
		return
	}
	m.applyDimensionsNode(&m.layout, 0, 0, m.width, m.height)
}

func (m *Manager) applyDimensionsNode(node *LayoutNodeRef, x, y, w, h int) {
	if node == nil {
		return
	}

	// Registra dimensões calculadas para uso posterior na renderização.
	node.ComputedWidth = w
	node.ComputedHeight = h

	switch node.Type {
	case "box":
		if node.ComponentID == "" {
			return
		}
		if c, ok := m.components[node.ComponentID]; ok && c != nil {
			c.SetDimensions(w, h)
		}
	case "row":
		count := len(node.Items)
		if count == 0 {
			return
		}

		// Primeiro, verificamos se algum filho define Flex explicitamente.
		// Se sim, usamos divisão proporcional por Flex; caso contrário,
		// mantemos o comportamento anterior (divisão igualitária).
		hasFlex := false
		totalFlex := 0
		childFlex := make([]int, count)
		for i := range node.Items {
			f := 0
			if node.Items[i].Flex != nil && *node.Items[i].Flex > 0 {
				f = *node.Items[i].Flex
				hasFlex = true
			}
			childFlex[i] = f
			totalFlex += f
		}

		if !hasFlex || totalFlex == 0 {
			// Fallback: divisão em colunas iguais (comportamento anterior).
			childWidth := w / count
			for i := range node.Items {
				cw := childWidth
				if i == count-1 {
					cw = w - childWidth*(count-1)
				}
				m.applyDimensionsNode(&node.Items[i], x+i*cw, y, cw, h)
			}
			return
		}

		// Quando há Flex, consideramos todos os filhos com Flex==0 como 1
		// para evitar colunas invisíveis.
		for i := range childFlex {
			if childFlex[i] == 0 {
				childFlex[i] = 1
				totalFlex++
			}
		}

		remainingWidth := w
		for i := range node.Items {
			// Distribui largura proporcionalmente ao Flex.
			cw := remainingWidth * childFlex[i] / totalFlex
			if i == count-1 {
				// Garante que o somatório não ultrapasse/escape por causa de divisão inteira.
				cw = remainingWidth
			}
			m.applyDimensionsNode(&node.Items[i], x, y, cw, h)
			x += cw
			remainingWidth -= cw
			totalFlex -= childFlex[i]
		}
	case "column":
		count := len(node.Items)
		if count == 0 {
			return
		}

		// Similar ao caso row, mas dividindo altura.
		hasFlex := false
		totalFlex := 0
		childFlex := make([]int, count)
		for i := range node.Items {
			f := 0
			if node.Items[i].Flex != nil && *node.Items[i].Flex > 0 {
				f = *node.Items[i].Flex
				hasFlex = true
			}
			childFlex[i] = f
			totalFlex += f
		}

		if !hasFlex || totalFlex == 0 {
			childHeight := h / count
			for i := range node.Items {
				ch := childHeight
				if i == count-1 {
					ch = h - childHeight*(count-1)
				}
				m.applyDimensionsNode(&node.Items[i], x, y+i*ch, w, ch)
			}
			return
		}

		for i := range childFlex {
			if childFlex[i] == 0 {
				childFlex[i] = 1
				totalFlex++
			}
		}

		remainingHeight := h
		for i := range node.Items {
			ch := remainingHeight * childFlex[i] / totalFlex
			if i == count-1 {
				ch = remainingHeight
			}
			m.applyDimensionsNode(&node.Items[i], x, y, w, ch)
			y += ch
			remainingHeight -= ch
			totalFlex -= childFlex[i]
		}
	}
}

func (m *Manager) renderNode(node LayoutNodeRef) string {
	switch node.Type {
	case "box":
		if node.ComponentID == "" {
			return ""
		}
		inner := ""
		if c, ok := m.components[node.ComponentID]; ok && c != nil {
			inner = c.View()
		}
		// Aplica padding/borda declarativos apenas na renderização textual.
		pad := 0
		if node.Padding != nil && *node.Padding > 0 {
			pad = *node.Padding
		}
		if !node.Border && pad == 0 {
			return inner
		}
		boxView := renderBox(inner, pad, node.Border)
		// Normaliza a largura final do box para a largura calculada do layout,
		// garantindo que bordas e colunas alinhem visualmente.
		if node.ComputedWidth > 0 {
			lines := strings.Split(boxView, "\n")
			for i, l := range lines {
				if len(l) < node.ComputedWidth {
					lines[i] = l + strings.Repeat(" ", node.ComputedWidth-len(l))
				} else if len(l) > node.ComputedWidth {
					lines[i] = l[:node.ComputedWidth]
				}
			}
			boxView = strings.Join(lines, "\n")
		}
		return boxView
	case "row":
		if len(node.Items) == 0 {
			return ""
		}
		// Renderização horizontal usando Lip Gloss, respeitando larguras/alturas
		// calculadas para cada filho.
		childViews := make([]string, 0, len(node.Items))
		for _, child := range node.Items {
			v := m.renderNode(child)
			style := lipgloss.NewStyle()
			if child.ComputedWidth > 0 {
				style = style.Width(child.ComputedWidth)
			}
			if child.ComputedHeight > 0 {
				style = style.Height(child.ComputedHeight)
			}
			childViews = append(childViews, style.Render(v))
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, childViews...)
	case "column":
		if len(node.Items) == 0 {
			return ""
		}
		childViews := make([]string, 0, len(node.Items))
		for _, child := range node.Items {
			v := m.renderNode(child)
			style := lipgloss.NewStyle()
			if child.ComputedWidth > 0 {
				style = style.Width(child.ComputedWidth)
			}
			if child.ComputedHeight > 0 {
				style = style.Height(child.ComputedHeight)
			}
			childViews = append(childViews, style.Render(v))
		}
		return lipgloss.JoinVertical(lipgloss.Left, childViews...)
	default:
		return ""
	}
}

// canRenderColumnWithFlexbox verifica se uma coluna é simples o suficiente
// para ser renderizada por um único FlexBox, onde cada row declarativa vira
// uma Row do FlexBox e cada box dentro dela vira uma Cell.
// Critérios atuais:
// - coluna tem largura/altura calculadas (>0);
// - todos os filhos imediatos são rows;
// - cada row contém apenas boxes com ComponentID definido.

// renderBox aplica padding e borda usando Lip Gloss.
// - padding controla a quantidade de espaços internos ao redor do conteúdo.
// - quando border==true, desenha uma borda normal.
func renderBox(content string, padding int, border bool) string {
	style := lipgloss.NewStyle()
	if padding > 0 {
		style = style.Padding(padding)
	}
	if border {
		style = style.Border(lipgloss.NormalBorder())
	}
	return style.Render(content)
}
