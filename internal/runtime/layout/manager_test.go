package layout

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// fakeComponent é um stub mínimo de ShantillyComponent para validar o contrato
// estrutural do LayoutManager sem acoplar a componentes reais (governança E1.1).
type fakeComponent struct {
	id         string
	initCalled bool
	updateMsg  tea.Msg
	view       string
	width      int
	height     int
}

func (f *fakeComponent) ID() string    { return f.id }
func (f *fakeComponent) Init() tea.Cmd { f.initCalled = true; return nil }
func (f *fakeComponent) Update(msg tea.Msg) (ShantillyComponent, tea.Cmd) {
	f.updateMsg = msg
	return f, nil
}
func (f *fakeComponent) View() string       { return f.view }
func (f *fakeComponent) SetBounds(w, h int) { f.width, f.height = w, h }

// fakeRegistry garante que o Manager permaneça dependente apenas de ComponentRegistry.
type fakeRegistry struct {
	components map[string]ShantillyComponent
}

func (r *fakeRegistry) Resolve(id string) ShantillyComponent {
	if r.components == nil {
		return nil
	}
	return r.components[id]
}

// TestNew_RespectsGovernanceContracts valida que o construtor cria uma instância
// neutra, sem execução de automação ou efeitos colaterais externos.
func TestNew_RespectsGovernanceContracts(t *testing.T) {
	root := LayoutNodeRef{
		ID:   "root",
		Type: "column",
	}
	reg := &fakeRegistry{}
	m := New(root, reg)

	if m == nil {
		t.Fatalf("New must not return nil")
	}

	if m.layout.ID != "root" {
		t.Errorf("expected layout ID 'root', got %q", m.layout.ID)
	}

	if m.registry != reg {
		t.Errorf("registry was not stored correctly")
	}

	if len(m.components) != 0 {
		t.Errorf("expected no components to be initialized eagerly")
	}

	if m.initialized {
		t.Errorf("manager must start uninitialized")
	}
}

// TestInit_IdempotentAndNoAutomation garante apenas comportamento estrutural:
// - Init é idempotente;
// - Não há side effects visíveis além do flag initialized (sem execução de run/on).
func TestInit_IdempotentAndNoAutomation(t *testing.T) {
	root := LayoutNodeRef{
		ID:   "root",
		Type: "column",
	}
	reg := &fakeRegistry{}
	m := New(root, reg)

	if m.initialized {
		t.Fatalf("expected initialized=false before Init")
	}

	cmd1 := m.Init()
	if !m.initialized {
		t.Errorf("expected initialized=true after first Init")
	}
	if cmd1 != nil {
		t.Errorf("expected nil cmd from Init placeholder, got non-nil")
	}

	// Segunda chamada deve ser neutra.
	cmd2 := m.Init()
	if cmd2 != nil {
		t.Errorf("expected nil cmd from second Init (idempotent), got non-nil")
	}
}

// TestUpdate_HandlesWindowSizeMsgWithoutAutomation
// Foca apenas na responsabilidade de layout: armazenar width/height.
func TestUpdate_HandlesWindowSizeMsgWithoutAutomation(t *testing.T) {
	root := LayoutNodeRef{
		ID:   "root",
		Type: "column",
	}
	reg := &fakeRegistry{}
	m := New(root, reg)

	if m.width != 0 || m.height != 0 {
		t.Fatalf("expected initial width/height to be zero")
	}

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	next, cmd := m.Update(msg)

	if next != m {
		t.Errorf("expected Update to return same manager instance")
	}
	if cmd != nil {
		t.Errorf("expected no cmd from WindowSizeMsg handling")
	}
	if m.width != 120 || m.height != 40 {
		t.Errorf("expected width=120,height=40, got width=%d,height=%d", m.width, m.height)
	}
}

// TestUpdate_RoutesToFocusedComponentConfirmsIsolation
// Garante apenas roteamento mínimo ao componente focado,
// sem qualquer integração com EventManager, on: ou run:.
func TestUpdate_RoutesToFocusedComponentConfirmsIsolation(t *testing.T) {
	root := LayoutNodeRef{
		ID:   "root",
		Type: "column",
	}
	focusedComp := &fakeComponent{id: "comp1"}
	reg := &fakeRegistry{
		components: map[string]ShantillyComponent{
			"comp1": focusedComp,
		},
	}
	m := New(root, reg)

	// Registrar manualmente como se já tivesse sido instanciado pelo Init futuro.
	m.components["comp1"] = focusedComp
	m.focusedID = "comp1"

	customMsg := tea.KeyMsg{} // qualquer msg TEA neutra
	nextModel, cmd := m.Update(customMsg)

	if nextModel != m {
		t.Errorf("expected manager to remain root model")
	}
	if cmd != nil {
		// Placeholder atual sempre retorna nil; qualquer cmd aqui seria suspeito.
		t.Errorf("expected no cmd from focused component in this skeleton phase")
	}
	if focusedComp.updateMsg == nil {
		t.Errorf("expected focused component to receive routed message, got nil")
	}
}

// TestView_PlaceholderIsStructuralOnly
// Assegura que o esqueleto de View não introduz lógica de automação.
func TestView_PlaceholderIsStructuralOnly(t *testing.T) {
	root := LayoutNodeRef{
		ID:   "root",
		Type: "column",
	}
	reg := &fakeRegistry{}
	m := New(root, reg)

	view := m.View()
	if view != "" {
		// Na fase atual, qualquer conteúdo não-vazio pode indicar comportamento além do esqueleto.
		t.Errorf("expected empty placeholder view, got %q", view)
	}
}
