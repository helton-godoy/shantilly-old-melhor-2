package wizard

import (
	"fmt"

	"shantilly/internal/util"
	"shantilly/pkg/declarative"

	"github.com/charmbracelet/huh"
)

// Run executa um wizard TUI para coletar informações básicas
// e gerar um AppConfig declarativo v2.0 em memória.
func Run() (*declarative.AppConfig, error) {
	var (
		name        string
		description string
		layoutName  string
		withScripts bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Nome da aplicação").Value(&name),
			huh.NewText().Title("Descrição").Value(&description),
			huh.NewSelect[string]().
				Title("Layout base").
				Options(
					huh.NewOption("Menu + Viewport (demo)", "demo"),
					huh.NewOption("Menu simples", "simple"),
				).
				Value(&layoutName),
			// Pergunta opcional sobre inclusão de scripts de exemplo.
			huh.NewConfirm().
				Title("Incluir scripts de exemplo (on + security.allowed_scripts)?").
				Value(&withScripts),
		),
	)

	if err := form.Run(); err != nil {
		// Tratamos qualquer erro do wizard como cancelamento pelo usuário,
		// para que a casca CLI possa mapear para ExitCancelled.
		return nil, fmt.Errorf("wizard cancelado: %w", util.ErrAborted)
	}

	if layoutName == "" {
		layoutName = "demo"
	}

	cfg := &declarative.AppConfig{
		Version: "1.0",
		Metadata: map[string]string{
			"name":        name,
			"description": description,
		},
	}

	switch layoutName {
	case "simple":
		applySimpleTemplate(cfg, withScripts)
	default:
		applyDemoTemplate(cfg, withScripts)
	}

	return cfg, nil
}

func applySimpleTemplate(cfg *declarative.AppConfig, withScripts bool) {
	cfg.Layout = declarative.LayoutNode{
		ID:   "root",
		Type: "column",
		Items: []declarative.LayoutNode{
			{
				ID:   "main_row",
				Type: "row",
				Items: []declarative.LayoutNode{
					{
						ID:          "menu_box",
						Type:        "box",
						ComponentID: "menu",
					},
				},
			},
		},
	}

	cfg.Components = []declarative.Component{
		{
			ID:   "menu",
			Type: "list",
			Props: map[string]any{
				"items": []map[string]any{
					{"id": "item1", "text": "Item 1"},
					{"id": "item2", "text": "Item 2"},
				},
			},
		},
	}

	if withScripts {
		cfg.On = []declarative.OnHandler{
			{
				ID:    "on_menu_select_item1",
				Event: "menu:select_item1",
				Run: &declarative.RunAction{
					Script: "./scripts/hello.sh",
					Args:   []string{},
					Stdin:  map[string]interface{}{},
					Env:    map[string]string{},
					// Para o layout simple não há viewport; ainda assim deixamos update_target
					// vazio para futuras integrações.
					UpdateTarget:  "",
					Confirm:       false,
					PromptSecrets: []string{},
				},
			},
		}

		cfg.Security = &declarative.SecurityPolicy{
			AllowedScripts: []string{"./scripts/hello.sh"},
			DenyUnknown:    true,
			ExtraRules:     map[string]string{},
		}
		return
	}

	cfg.Security = &declarative.SecurityPolicy{}
}

func applyDemoTemplate(cfg *declarative.AppConfig, withScripts bool) {
	cfg.Layout = declarative.LayoutNode{
		ID:   "root",
		Type: "column",
		Items: []declarative.LayoutNode{
			{
				ID:   "main_row",
				Type: "row",
				Items: []declarative.LayoutNode{
					{
						ID:          "menu_box",
						Type:        "box",
						ComponentID: "menu",
					},
					{
						ID:          "output_box",
						Type:        "box",
						ComponentID: "output",
					},
				},
			},
		},
	}

	cfg.Components = []declarative.Component{
		{
			ID:   "menu",
			Type: "list",
			Props: map[string]any{
				"items": []map[string]any{
					{"id": "check_legacy", "text": "1) Verificar encapsulamento do legado (script real)"},
					{"id": "noop", "text": "2) Não fazer nada (placeholder)"},
					{"id": "hello_demo", "text": "3) Rodar script de Hello (demo rápida do viewport)"},
				},
			},
		},
		{
			ID:   "output",
			Type: "viewport",
			Props: map[string]any{
				"source": map[string]any{
					"type":         "static",
					"content_type": "text",
					"content":      "Bem-vindo ao Shantilly Runtime v2.0\n\nUse o menu à esquerda para disparar scripts de exemplo.",
				},
			},
		},
	}

	if withScripts {
		cfg.On = []declarative.OnHandler{
			{
				ID:    "on_menu_select_check_legacy",
				Event: "menu:select_check_legacy",
				Run: &declarative.RunAction{
					Script: "./scripts/check-legacy-encapsulation.sh",
					Args:   []string{},
					Stdin:  map[string]interface{}{},
					Env:    map[string]string{},
					// Segue o app.example.yaml usando "output" como update_target.
					UpdateTarget:  "output",
					Confirm:       true,
					PromptSecrets: []string{},
				},
			},
			{
				ID:    "on_menu_select_hello_demo",
				Event: "menu:select_hello_demo",
				Run: &declarative.RunAction{
					Script:        "./scripts/hello.sh",
					Args:          []string{},
					Stdin:         map[string]interface{}{},
					Env:           map[string]string{},
					UpdateTarget:  "output",
					Confirm:       false,
					PromptSecrets: []string{},
				},
			},
		}

		cfg.Security = &declarative.SecurityPolicy{
			AllowedScripts: []string{
				"./scripts/check-legacy-encapsulation.sh",
				"./scripts/hello.sh",
			},
			DenyUnknown: true,
			ExtraRules:  map[string]string{},
		}
		return
	}

	cfg.Security = &declarative.SecurityPolicy{}
}
