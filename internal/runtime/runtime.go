package runtime

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"shantilly/internal/runtime/event"
	"shantilly/internal/runtime/layout"
	"shantilly/internal/runtime/runner"
	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
)

// Start é o ponto de entrada do Runtime TUI declarativo v2.0.
//
// Ele recebe uma AppConfig já carregada/validada, constrói o LayoutManager
// a partir da árvore de layout declarativa e inicializa o MainModel com
// suporte a overlay de modais.
func Start(cfg *declarative.AppConfig) error {
	if cfg == nil {
		return fmt.Errorf("runtime.Start: config is nil")
	}

	// 1. Criar o tema padrão para os componentes.
	theme := tui.NewDefaultTheme()

	// 2. Converter LayoutNode (declarative) -> LayoutNodeRef (layout engine atual).
	rootRef := convertLayoutNode(cfg.Layout)

	// 3. Criar o registry concreto de componentes a partir de AppConfig.Components.
	registry := layout.NewDefaultRegistry(theme, cfg.Components)

	// 4. Criar LayoutManager com o registry de componentes.
	lm := layout.New(rootRef, registry)

	// 5. Instanciar motores de eventos e scripts.
	em := event.New(cfg.On)
	sr := runner.NewScriptRunner(nil, convertSecurityPolicy(cfg.Security))
	lm.SetEventManager(em)
	lm.SetScriptRunner(sr)

	// 6. Criar MainModel com overlay, embrulhando o LayoutManager.
	mainModel := NewLayoutMainModel(lm)

	// 7. Iniciar o programa Bubble Tea.
	// Nesta fase, utilizamos AltScreen para fornecer uma experiência mais
	// próxima da interface final para o usuário.
	p := tea.NewProgram(
		mainModel,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("erro durante a execução do runtime TUI: %w", err)
	}

	return nil
}

// convertLayoutNode mapeia o modelo declarativo normativo (LayoutNode)
// para o LayoutNodeRef usado internamente pelo LayoutManager nesta wave.
func convertLayoutNode(n declarative.LayoutNode) layout.LayoutNodeRef {
	ref := layout.LayoutNodeRef{
		ID:          n.ID,
		Type:        n.Type,
		Width:       n.Width,
		Height:      n.Height,
		Flex:        n.Flex,
		ComponentID: n.ComponentID,
	}

	if len(n.Items) > 0 {
		children := make([]layout.LayoutNodeRef, 0, len(n.Items))
		for _, child := range n.Items {
			children = append(children, convertLayoutNode(child))
		}
		ref.Items = children
	}

	return ref
}

// convertSecurityPolicy adapta declarative.SecurityPolicy para o contrato do ScriptRunner.
func convertSecurityPolicy(p *declarative.SecurityPolicy) *runner.SecurityPolicy {
	if p == nil {
		return nil
	}
	policy := &runner.SecurityPolicy{
		AllowedScripts: append([]string(nil), p.AllowedScripts...),
		DenyUnknown:    p.DenyUnknown,
		ExtraRules:     make(map[string]string, len(p.ExtraRules)),
	}
	for k, v := range p.ExtraRules {
		policy.ExtraRules[k] = v
	}
	return policy
}
