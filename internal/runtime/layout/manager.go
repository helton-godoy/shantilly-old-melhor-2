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
)

// ShantillyComponent é o contrato esperado dos componentes concretos do runtime v2.0.
// Definido conceitualmente em docs/architecture/components.md#1-shantillycomponent-contrato-base--e12.
// A interface é espelhada aqui para desacoplamento, podendo ser movida para um pacote compartilhado
// quando a implementação concreta for introduzida por bmad-dev/bmad-master.
type ShantillyComponent interface {
	ID() string
	Init() tea.Cmd
	Update(msg tea.Msg) (ShantillyComponent, tea.Cmd)
	View() string
	SetBounds(width, height int)
}

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

// Update implementa o loop TEA raiz para o LayoutManager.
// Regras nesta wave:
// - Tratar mensagens de resize para atualizar width/height.
// - Roteamento básico de msgs para componente focado (quando existir).
// - NÃO executar automação nem avaliar on:/run: (delegado ao EventManager).
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Em implementação futura:
		// - recalcular bounds para todos os nós.
		// - chamar SetBounds em cada ShantillyComponent.
		return m, nil
	default:
	}

	// Roteamento mínimo para componente focado (se houver).
	if m.focusedID != "" {
		if c, ok := m.components[m.focusedID]; ok && c != nil {
			next, cmd := c.Update(msg)
			m.components[m.focusedID] = next
			return m, cmd
		}
	}

	return m, nil
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
