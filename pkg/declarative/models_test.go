// Wave 1 — Data Models v2.0 + YAML Parser Declarativo
// Epic 1: Runtime TUI Declarativo — Fundação do Runtime (v2.0)
// Blocos cobertos: E1.1 (Layout), E1.2 (Component Model), E1.3 (on:), E1.4 (RunAction), E1.5 (Security/Modal)
// Rastreio:
// - PRD: docs/prd/epic-1-runtime-tui-foundation.md
// - Arquitetura: docs/architecture/data-models.md, docs/architecture/high-level-architecture.md
// - QA Matrix: docs/qa/matrix-epic-1-runtime-tui-coverage.md (linhas Wave 1 - YAML config)
//
// Objetivo destes testes:
// - Validar LoadAppConfig(r io.Reader) como parser declarativo mínimo.
// - Garantir validações estruturais mínimas: layout obrigatório, tipos de layout válidos,
//   componentes com IDs únicos e tipos válidos, referências de layout -> componentes consistentes.
// - Cobertura de cenários positivos e negativos, alinhado à governança v2.0.
//
// Observação:
// - Hardening completo (deny by default, SecurityPolicy estrita) será implementado na Wave 7.
// - Nenhum acoplamento ao legado v1.x (FormConfig) é introduzido aqui.

package declarative

import (
	"strings"
	"testing"
)

func TestLoadAppConfig_ValidMinimalConfig(t *testing.T) {
	yamlSrc := `
version: "1.0"
meta:
  app: "test"
layout:
  id: "root"
  type: "column"
  items:
    - id: "row-1"
      type: "row"
      items:
        - id: "box-1"
          type: "box"
          component: "main_view"
components:
  - id: "main_view"
    type: "viewport"
`

	cfg, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Layout.ID != "root" || cfg.Layout.Type != "column" {
		t.Errorf("unexpected root layout: %+v", cfg.Layout)
	}

	if len(cfg.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(cfg.Components))
	}

	if cfg.Components[0].ID != "main_view" || cfg.Components[0].Type != "viewport" {
		t.Errorf("unexpected component: %+v", cfg.Components[0])
	}
}

func TestLoadAppConfig_MissingLayout(t *testing.T) {
	yamlSrc := `
version: "1.0"
components:
  - id: "c1"
    type: "viewport"
`

	_, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err == nil {
		t.Fatal("expected error due to missing/invalid layout, got nil")
	}
}

func TestLoadAppConfig_InvalidLayoutType(t *testing.T) {
	yamlSrc := `
layout:
  id: "root"
  type: "grid" # invalid
components:
  - id: "c1"
    type: "viewport"
`

	_, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err == nil {
		t.Fatal("expected error for invalid layout type, got nil")
	}
	if !strings.Contains(err.Error(), "invalid layout type") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoadAppConfig_BoxWithItemsShouldFail(t *testing.T) {
	yamlSrc := `
layout:
  id: "root"
  type: "box"
  items:
    - id: "child"
      type: "box"
components:
  - id: "c1"
    type: "viewport"
`

	_, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err == nil {
		t.Fatal("expected error for box node with items, got nil")
	}
}

func TestLoadAppConfig_UnknownComponentReferenceInLayout(t *testing.T) {
	yamlSrc := `
layout:
  id: "root"
  type: "column"
  items:
    - id: "box-1"
      type: "box"
      component: "missing"
components:
  - id: "c1"
    type: "viewport"
`

	_, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err == nil {
		t.Fatal("expected error for unknown component reference, got nil")
	}
	if !strings.Contains(err.Error(), "references unknown component") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoadAppConfig_DuplicateComponentID(t *testing.T) {
	yamlSrc := `
layout:
  id: "root"
  type: "column"
components:
  - id: "dup"
    type: "viewport"
  - id: "dup"
    type: "form"
`

	_, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err == nil {
		t.Fatal("expected error for duplicate component id, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate component id") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoadAppConfig_InvalidComponentType(t *testing.T) {
	yamlSrc := `
layout:
  id: "root"
  type: "column"
components:
  - id: "c1"
    type: "unknown"
`

	_, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err == nil {
		t.Fatal("expected error for invalid component type, got nil")
	}
	if !strings.Contains(err.Error(), "invalid component type") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoadAppConfig_OnHandlersRequireEvent(t *testing.T) {
	yamlSrc := `
layout:
  id: "root"
  type: "column"
components:
  - id: "main"
    type: "viewport"
on:
  - id: "no_event"
    run:
      script: "./do.sh"
`

	_, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err == nil {
		t.Fatal("expected error for on: handler missing event, got nil")
	}
	if !strings.Contains(err.Error(), "handler missing event") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoadAppConfig_WithSecurityPolicyStructurallyAccepted(t *testing.T) {
	yamlSrc := `
layout:
  id: "root"
  type: "column"
components:
  - id: "main"
    type: "viewport"
security:
  allowed_scripts:
    - "./safe.sh"
  deny_unknown: false
  extra_rules:
    foo: "bar"
`

	cfg, err := LoadAppConfig(strings.NewReader(yamlSrc))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Security == nil {
		t.Fatal("expected security policy to be parsed")
	}
	if len(cfg.Security.AllowedScripts) != 1 || cfg.Security.AllowedScripts[0] != "./safe.sh" {
		t.Errorf("unexpected security policy: %+v", cfg.Security)
	}
}
