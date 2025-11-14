// Package declarative
//
// Wave 1 — Data Models v2.0 + YAML Parser Declarativo
// Epic 1: Runtime TUI Declarativo — Fundação do Runtime (v2.0)
// Blocos: E1.1 (LayoutManager), E1.2 (Component Model), E1.3 (Event Engine),
//
//	E1.4 (ScriptRunner), E1.5 (Modal Stack + Security)
//
// Rastreio:
//   - PRD: docs/prd/epic-1-runtime-tui-foundation.md
//   - Arquitetura: docs/architecture/data-models.md,
//     docs/architecture/high-level-architecture.md,
//     docs/architecture/components.md
//   - Story Map: docs/stories/index-epic-1-runtime-tui.story-map.md
//   - QA: docs/qa/matrix-epic-1-runtime-tui-coverage.md
//
// Este pacote define os modelos declarativos NORMATIVOS usados pelo runtime v2.0.
// Ele substitui o modelo linear centrado em FormConfig como visão ativa.
// FormConfig v1.x é considerado LEGADO e só pode ser utilizado internamente
// pelo FormComponent via contrato ShantillyComponent.
package declarative

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// AppConfig representa o documento YAML único consumido pelo runtime v2.0.
//
// Fonte normativa: docs/architecture/data-models.md#1-appconfig-e11e15--raiz-do-yaml-unico
//
// Invariantes principais (validados parcialmente em Validate()):
// - Um único documento raiz para layout, componentes, regras on: e política de segurança.
// - Nenhum outro arquivo/estrutura paralela redefine estes conceitos.
type AppConfig struct {
	Version    string            `yaml:"version,omitempty" json:"version,omitempty"`
	Metadata   map[string]string `yaml:"meta,omitempty" json:"meta,omitempty"`
	Layout     LayoutNode        `yaml:"layout" json:"layout"`
	Components []Component       `yaml:"components" json:"components"`
	On         []OnHandler       `yaml:"on,omitempty" json:"on,omitempty"`
	Security   *SecurityPolicy   `yaml:"security,omitempty" json:"security,omitempty"`
}

// LayoutNode define a árvore de layout hierárquico declarativo.
//
// Fonte normativa: docs/architecture/data-models.md#2-layoutnode-e11--layoutmanager
type LayoutNode struct {
	ID          string       `yaml:"id,omitempty" json:"id,omitempty"`
	Type        string       `yaml:"type" json:"type"` // "column" | "row" | "box"
	Width       *int         `yaml:"width,omitempty" json:"width,omitempty"`
	Height      *int         `yaml:"height,omitempty" json:"height,omitempty"`
	Flex        *int         `yaml:"flex,omitempty" json:"flex,omitempty"`
	Items       []LayoutNode `yaml:"items,omitempty" json:"items,omitempty"`
	ComponentID string       `yaml:"component,omitempty" json:"component,omitempty"` // apenas para type: box
}

// Component representa componentes TUI declarativos instanciados no layout.
//
// Fonte normativa: docs/architecture/data-models.md#3-component-e12--shantillycomponent--componentes-oficiais
type Component struct {
	ID       string            `yaml:"id" json:"id"`     // Referenciado por layout, on:, run:
	Type     string            `yaml:"type" json:"type"` // "list" | "viewport" | "form" | "buttongroup"
	Props    map[string]any    `yaml:"props,omitempty" json:"props,omitempty"`
	Bindings map[string]string `yaml:"bind,omitempty" json:"bind,omitempty"` // ex.: estados, seleção, etc.
}

// Item representa opções declarativas usadas por componentes como listas e grupos de botões.
type Item struct {
	ID    string `yaml:"id" json:"id"`
	Text  string `yaml:"text,omitempty" json:"text,omitempty"`
	Label string `yaml:"label,omitempty" json:"label,omitempty"`
}

// Source descreve conteúdos exibidos por componentes como Viewport.
type Source struct {
	Type        string `yaml:"type" json:"type"`
	Content     string `yaml:"content,omitempty" json:"content,omitempty"`
	ContentType string `yaml:"content_type,omitempty" json:"content_type,omitempty"`
}

// OnHandler modela uma regra declarativa on:.
//
// Fonte normativa: docs/architecture/data-models.md#4-onhandler-e13--bloco-on-como-unica-orquestracao
type OnHandler struct {
	ID          string                 `yaml:"id,omitempty" json:"id,omitempty"`
	Event       string                 `yaml:"event" json:"event"` // ex.: "menu:select", "user_form:submit"
	When        string                 `yaml:"when,omitempty" json:"when,omitempty"`
	Run         *RunAction             `yaml:"run,omitempty" json:"run,omitempty"`
	OpenModal   *ModalRequest          `yaml:"open_modal,omitempty" json:"open_modal,omitempty"`
	UpdateState map[string]interface{} `yaml:"update_state,omitempty" json:"update_state,omitempty"`
}

// ShantillyEvent é o tipo canônico de evento usado pelo EventManager.
//
// Fonte normativa: docs/architecture/data-models.md#5-shantillyevent-e13--tipo-de-evento-normalizado
type ShantillyEvent struct {
	SourceComponentID string                 `json:"source_component_id"`
	Type              string                 `json:"type"`              // ex.: "submit", "select", "click", "script.complete", "modal.confirmed"
	Payload           map[string]interface{} `json:"payload,omitempty"` // dados contextuais
}

// RunAction especifica ações run: disparadas por OnHandler.
//
// Fonte normativa: docs/architecture/data-models.md#6-runaction-e14--execucao-declarativa
type RunAction struct {
	Script        string                 `yaml:"script" json:"script"`
	Args          []string               `yaml:"args,omitempty" json:"args,omitempty"`
	Stdin         map[string]interface{} `yaml:"stdin,omitempty" json:"stdin,omitempty"`
	Env           map[string]string      `yaml:"env,omitempty" json:"env,omitempty"`
	UpdateTarget  string                 `yaml:"update_target,omitempty" json:"update_target,omitempty"`
	Confirm       bool                   `yaml:"confirm,omitempty" json:"confirm,omitempty"`
	PromptSecrets []string               `yaml:"prompt_secrets,omitempty" json:"prompt_secrets,omitempty"`
}

// ModalRequest representa pedido declarativo para abrir um modal.
//
// Fonte normativa: docs/architecture/data-models.md#7-modalrequest-e15--modal-stack
type ModalRequest struct {
	Type    string       `yaml:"type" json:"type"` // "confirm" | "prompt" | tipos seguros
	Title   string       `yaml:"title,omitempty" json:"title,omitempty"`
	Message string       `yaml:"message,omitempty" json:"message,omitempty"`
	Fields  []ModalField `yaml:"fields,omitempty" json:"fields,omitempty"`
}

// ModalField define campos de coleta em modais.
type ModalField struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"` // ex.: "secret"
}

// SecurityPolicy define políticas declarativas anti-Trojan YAML.
//
// Fonte normativa: docs/architecture/data-models.md#8-securitypolicy-e15--anti-trojan-yaml
// Wave 7 fará o hardening completo; aqui aplicamos apenas estrutura mínima.
type SecurityPolicy struct {
	AllowedScripts []string          `yaml:"allowed_scripts,omitempty" json:"allowed_scripts,omitempty"`
	DenyUnknown    bool              `yaml:"deny_unknown,omitempty" json:"deny_unknown,omitempty"`
	ExtraRules     map[string]string `yaml:"extra_rules,omitempty" json:"extra_rules,omitempty"`
}

// LoadAppConfig lê o YAML da AppConfig a partir de um io.Reader,
// aplica unmarshalling para AppConfig e executa validações estruturais mínimas.
//
// Wave 1 — Parser Declarativo (mínimo, sem hardening completo de segurança).
//
// Critérios mínimos (alinhados ao Implementation Plan Wave 1):
// - Campo layout obrigatório.
// - Pelo menos um componente quando layout fizer referência a component.
// - IDs de componentes únicos.
// - Referências de LayoutNode.ComponentID devem apontar para componentes existentes.
// - Não introduz dependências novas no legado v1.x (FormConfig).
func LoadAppConfig(r io.Reader) (*AppConfig, error) {
	if r == nil {
		return nil, fmt.Errorf("LoadAppConfig: reader is nil")
	}

	var cfg AppConfig
	dec := yaml.NewDecoder(r)
	dec.KnownFields(false) // Hardening completo (deny unknown) será feito na Wave 7.

	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("LoadAppConfig: failed to decode YAML: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate executa validações mínimas de estrutura da AppConfig.
// Não cobre ainda todas as regras de segurança JIT / anti-Trojan (Wave 7).
func (c *AppConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("app config is nil")
	}

	// Layout.type deve ser válido.
	if err := validateLayoutNode(&c.Layout); err != nil {
		return fmt.Errorf("layout validation failed: %w", err)
	}

	// Componentes: IDs únicos e tipos válidos.
	componentIndex := make(map[string]Component, len(c.Components))
	for _, comp := range c.Components {
		if comp.ID == "" {
			return fmt.Errorf("component id is required")
		}
		if _, exists := componentIndex[comp.ID]; exists {
			return fmt.Errorf("duplicate component id: %s", comp.ID)
		}
		switch comp.Type {
		case "list", "viewport", "form", "buttongroup":
			// ok
		default:
			return fmt.Errorf("invalid component type for %s: %s", comp.ID, comp.Type)
		}
		componentIndex[comp.ID] = comp
	}

	// Validar referencias componentID em nós box.
	if err := validateLayoutComponentRefs(&c.Layout, componentIndex); err != nil {
		return err
	}

	// OnHandlers mínimos: quando presentes, evento não vazio.
	for _, h := range c.On {
		if h.Event == "" {
			return fmt.Errorf("on: handler missing event")
		}
	}

	// SecurityPolicy mínima não exige validação complexa nesta wave.
	return nil
}

func validateLayoutNode(n *LayoutNode) error {
	if n == nil {
		return fmt.Errorf("layout node is nil")
	}

	switch n.Type {
	case "column", "row":
		// Pode ter items, não pode ter ComponentID obrigatório.
	case "box":
		// Box pode ter no máximo 0 ou 1 ComponentID; não deve ter children de layout.
		if len(n.Items) > 0 {
			return fmt.Errorf("layout node %q (box) must not have items", n.ID)
		}
	default:
		return fmt.Errorf("invalid layout type %q in node %q", n.Type, n.ID)
	}

	for i := range n.Items {
		if err := validateLayoutNode(&n.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

func validateLayoutComponentRefs(n *LayoutNode, components map[string]Component) error {
	if n == nil {
		return nil
	}

	if n.Type == "box" && n.ComponentID != "" {
		if _, ok := components[n.ComponentID]; !ok {
			return fmt.Errorf("layout node %q references unknown component %q", n.ID, n.ComponentID)
		}
	}

	for i := range n.Items {
		if err := validateLayoutComponentRefs(&n.Items[i], components); err != nil {
			return err
		}
	}
	return nil
}
