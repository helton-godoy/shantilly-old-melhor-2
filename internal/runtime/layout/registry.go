package layout

import (
	"fmt"

	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
	components_buttongroup "shantilly/internal/components/buttongroup"
	components_list "shantilly/internal/components/list"
	components_viewport "shantilly/internal/components/viewport"
)

// DefaultRegistry é uma implementação concreta de ComponentRegistry que
// instancia componentes oficiais a partir de declarative.Component.
//
// Este registry não executa automação nem avalia on:/run:; ele apenas
// cria instâncias de ShantillyComponent com base em AppConfig.Components.
type DefaultRegistry struct {
	Theme      *tui.Theme
	Components map[string]declarative.Component
}

// NewDefaultRegistry constrói um registry a partir da lista de componentes
// declarativos carregados em AppConfig.
func NewDefaultRegistry(theme *tui.Theme, comps []declarative.Component) *DefaultRegistry {
	index := make(map[string]declarative.Component, len(comps))
	for _, c := range comps {
		index[c.ID] = c
	}
	return &DefaultRegistry{
		Theme:      theme,
		Components: index,
	}
}

// Resolve implementa ComponentRegistry.Resolve.
// Recebe o ID lógico do componente (referenciado pelo LayoutNode.ComponentID)
// e devolve a instância concreta correspondente.
func (r *DefaultRegistry) Resolve(id string) ShantillyComponent {
	if r == nil {
		return nil
	}
	comp, ok := r.Components[id]
	if !ok {
		// Fallback seguro: viewport estático com mensagem de erro.
		return components_viewport.New(id, r.Theme, &declarative.Source{
			Type:    "static",
			Content: fmt.Sprintf("Componente desconhecido: %s", id),
		})
	}

	switch comp.Type {
	case "viewport":
		// Espera-se que Props contenha informações de Source, quando aplicável.
		// Nesta wave, usamos apenas um Source simples, se houver.
		var src *declarative.Source
		if sAny, ok := comp.Props["source"]; ok {
			if s, ok2 := sAny.(declarative.Source); ok2 {
				src = &s
			}
		}
		return components_viewport.New(comp.ID, r.Theme, src)

	case "list":
		var items []declarative.Item
		if raw, ok := comp.Props["items"]; ok {
			if cast, ok2 := raw.([]declarative.Item); ok2 {
				items = cast
			}
		}
		return components_list.New(comp.ID, r.Theme, items)

	case "buttongroup":
		var items []declarative.Item
		if raw, ok := comp.Props["items"]; ok {
			if cast, ok2 := raw.([]declarative.Item); ok2 {
				items = cast
			}
		}
		return components_buttongroup.New(comp.ID, r.Theme, items)

	default:
		// Fallback: viewport com mensagem sobre tipo desconhecido.
		return components_viewport.New(comp.ID, r.Theme, &declarative.Source{
			Type: "static",
			Content: fmt.Sprintf("Tipo de componente desconhecido: %s", comp.Type),
		})
	}
}
