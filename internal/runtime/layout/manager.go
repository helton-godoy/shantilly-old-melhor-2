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
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/runtime/event"
	"shantilly/internal/runtime/runner"
	"shantilly/pkg/tui"
)

// ShantillyComponent é um alias para o contrato definido em pkg/tui.
type ShantillyComponent = tui.ShantillyComponent

// LayoutNodeRef é uma view mínima sobre o modelo declarativo de layout.
// Em implementação completa, este tipo será abastecido a partir de pkg/declarative.
type LayoutNodeRef struct {
	ID          string
	Type        string // "column" | "row" | "box"
	Width       *int
	Height      *int
	Flex        *int
	Items       []LayoutNodeRef
	ComponentID string // válido apenas para type: box
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

// Init inicializa componentes conforme o layout.
// Implementação inicial mínima: prepara estrutura interna.
// Expansão detalhada (instanciação e cmds) é responsabilidade de bmad-dev/bmad-master.
func (m *Manager) Init() tea.Cmd {
	if m.initialized {
		return nil
	}

	// Placeholder: em implementação futura,
	// - percorrer árvore layout,
	// - resolver ComponentID -> ShantillyComponent,
	// - chamar Init() de cada componente.
	m.initialized = true
	return nil
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Em implementação futura:
		// - recalcular bounds para todos os nós.
		// - chamar SetDimensions em cada ShantillyComponent.
		return m, nil

	// Pedido de execução de script vindo do EventManager.
	case tui.RunScriptRequestMsg:
		if m.scriptRunner != nil {
			cmd = runner.HandleRunRequest(m.scriptRunner, msg)
			cmds = append(cmds, cmd)
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
	// Implementação futura:
	// - Percorrer árvore layout,
	// - Concatenar views dos componentes respeitando bounds calculados.
	return ""
}
