package layout

import (
	"fmt"

	"shantilly/internal/components/button"
	components_buttongroup "shantilly/internal/components/buttongroup"
	components_input "shantilly/internal/components/input"
	components_list "shantilly/internal/components/list"
	components_multiselect "shantilly/internal/components/multiselect"
	components_select "shantilly/internal/components/select"
	components_viewport "shantilly/internal/components/viewport"
	"shantilly/internal/tui"
	"shantilly/pkg/declarative"
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
		}, "replace", true)
	}

	switch comp.Type {
	case "viewport":
		// Espera-se que Props["source"] venha do YAML como tipos genéricos (map[string]any).
		// Nesta wave, convertemos para declarative.Source com fallback seguro.
		var src *declarative.Source
		if raw, ok := comp.Props["source"]; ok {
			if converted, err := convertSource(raw); err == nil {
				src = converted
			} else {
				// Fallback: viewport estático com mensagem de erro de conversão.
				return components_viewport.New(comp.ID, r.Theme, &declarative.Source{
					Type:    "static",
					Content: fmt.Sprintf("Erro ao interpretar source do componente %s: %v", comp.ID, err),
				}, "replace", true)
			}
		}
		mode, _ := comp.Props["mode"].(string)
		wrap := true
		if rawWrap, ok := comp.Props["wrap"]; ok {
			if b, okb := rawWrap.(bool); okb {
				wrap = b
			}
		}
		return components_viewport.New(comp.ID, r.Theme, src, mode, wrap)

	case "list":
		var items []declarative.Item
		if raw, ok := comp.Props["items"]; ok {
			if converted, err := convertItems(raw); err == nil {
				items = converted
			} else {
				// Fallback: lista vazia com mensagem de erro na viewport.
				return components_viewport.New(comp.ID, r.Theme, &declarative.Source{
					Type:    "static",
					Content: fmt.Sprintf("Erro ao interpretar items do componente %s: %v", comp.ID, err),
				}, "replace", true)
			}
		}

		// Fallback de diagnóstico: se após a conversão não houver itens,
		// criamos alguns itens artificiais para verificar se o problema está
		// no layout/list ou na conversão de props.
		if len(items) == 0 {
			items = []declarative.Item{
				{ID: "debug_1", Text: "Item de debug 1"},
				{ID: "debug_2", Text: "Item de debug 2"},
			}
		}
		return components_list.New(comp.ID, r.Theme, items)

	case "select":
		var items []declarative.Item
		if raw, ok := comp.Props["items"]; ok {
			if converted, err := convertItems(raw); err == nil {
				items = converted
			} else {
				return components_viewport.New(comp.ID, r.Theme, &declarative.Source{
					Type:    "static",
					Content: fmt.Sprintf("Erro ao interpretar items do componente %s: %v", comp.ID, err),
				}, "replace", true)
			}
		}
		return components_select.New(comp.ID, r.Theme, items)

	case "multiselect":
		var items []declarative.Item
		if raw, ok := comp.Props["items"]; ok {
			if converted, err := convertItems(raw); err == nil {
				items = converted
			} else {
				return components_viewport.New(comp.ID, r.Theme, &declarative.Source{
					Type:    "static",
					Content: fmt.Sprintf("Erro ao interpretar items do componente %s: %v", comp.ID, err),
				}, "replace", true)
			}
		}
		return components_multiselect.New(comp.ID, r.Theme, items)

	case "buttongroup":
		var items []declarative.Item
		if raw, ok := comp.Props["items"]; ok {
			if converted, err := convertItems(raw); err == nil {
				items = converted
			} else {
				// Fallback: viewport com mensagem de erro de conversão.
				return components_viewport.New(comp.ID, r.Theme, &declarative.Source{
					Type:    "static",
					Content: fmt.Sprintf("Erro ao interpretar items do componente %s: %v", comp.ID, err),
				}, "replace", true)
			}
		}
		return components_buttongroup.New(comp.ID, r.Theme, items)

	case "input":
		label, _ := comp.Props["label"].(string)
		placeholder, _ := comp.Props["placeholder"].(string)
		initial, _ := comp.Props["default"].(string)
		secret, _ := comp.Props["secret"].(bool)
		return components_input.New(comp.ID, r.Theme, label, placeholder, initial, secret)

	case "button":
		// Extrair props do button
		text, _ := comp.Props["text"].(string)
		if text == "" {
			text = "Button"
		}
		icon, _ := comp.Props["icon"].(string)
		style, _ := comp.Props["style"].(string)
		disabled, _ := comp.Props["disabled"].(bool)

		props := button.Props{
			Text:     text,
			Icon:     icon,
			Style:    style,
			Disabled: disabled,
		}
		return button.New(comp.ID, r.Theme, props)

	default:
		// Fallback: viewport com mensagem sobre tipo desconhecido.
		return components_viewport.New(comp.ID, r.Theme, &declarative.Source{
			Type:    "static",
			Content: fmt.Sprintf("Tipo de componente desconhecido: %s", comp.Type),
		}, "replace", true)
	}
}

// convertItems converte props genéricos vindos do YAML (ex.: []any ou []map[string]any)
// para []declarative.Item. Aceita também []declarative.Item diretamente para
// compatibilidade com chamadas internas.
func convertItems(raw any) ([]declarative.Item, error) {
	if raw == nil {
		return nil, nil
	}

	// Já tipado corretamente.
	if items, ok := raw.([]declarative.Item); ok {
		return items, nil
	}

	// Slice genérico vindo do yaml.Unmarshal.
	switch v := raw.(type) {
	case []any:
		items := make([]declarative.Item, 0, len(v))
		for _, elem := range v {
			m, ok := elem.(map[string]any)
			if !ok {
				continue
			}
			var it declarative.Item
			if id, ok := m["id"].(string); ok {
				it.ID = id
			}
			if text, ok := m["text"].(string); ok {
				it.Text = text
			}
			if label, ok := m["label"].(string); ok {
				it.Label = label
			}
			items = append(items, it)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("formato inesperado para items: %T", raw)
	}
}

// convertSource converte props genéricos vindos do YAML (ex.: map[string]any)
// para *declarative.Source. Aceita declarative.Source diretamente para
// compatibilidade com chamadas internas.
func convertSource(raw any) (*declarative.Source, error) {
	if raw == nil {
		return nil, nil
	}

	// Já tipado corretamente.
	if s, ok := raw.(declarative.Source); ok {
		return &s, nil
	}
	if sp, ok := raw.(*declarative.Source); ok {
		return sp, nil
	}

	switch v := raw.(type) {
	case map[string]any:
		var src declarative.Source
		if t, ok := v["type"].(string); ok {
			src.Type = t
		}
		if c, ok := v["content"].(string); ok {
			src.Content = c
		}
		if ct, ok := v["contentType"].(string); ok {
			src.ContentType = ct
		}
		return &src, nil
	default:
		return nil, fmt.Errorf("formato inesperado para source: %T", raw)
	}
}
