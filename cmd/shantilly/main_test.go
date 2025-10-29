package main

import (
	"strings"
	"testing"

	"shantilly/internal/config"
)

func TestFormCommand_Stdin_ValidYAML(t *testing.T) {
	// Test parsing only - TUI integration is tested separately
	validYAML := `title: Test Form
fields:
  - key: name
    label: Name
    type: input
`

	_, err := config.Parse([]byte(validYAML))
	if err != nil {
		t.Fatalf("Expected no error when parsing YAML, got: %v", err)
	}
}

func TestFormCommand_Stdin_InvalidYAML(t *testing.T) {
	// Test parsing only - TUI integration is tested separately
	invalidYAML := `title: Test Form
fields:
  - key: name
    label: Name
    type: input
  - invalid yaml structure [
`

	_, err := config.Parse([]byte(invalidYAML))
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}

	if !strings.Contains(err.Error(), "erro ao parsear YAML") {
		t.Errorf("Expected error message to contain 'erro ao parsear YAML', got: %v", err)
	}
}

func TestFormCommand_Stdin_MissingRequiredFields(t *testing.T) {
	// Test parsing only - TUI integration is tested separately
	yamlWithoutKey := `fields:
  - label: Name
    type: input
`

	_, err := config.Parse([]byte(yamlWithoutKey))
	if err == nil {
		t.Fatal("Expected error for missing key field, got nil")
	}

	if !strings.Contains(err.Error(), "campo 1 deve ter uma chave") {
		t.Errorf("Expected error message to contain 'campo 1 deve ter uma chave', got: %v", err)
	}
}

func TestFormCommand_File(t *testing.T) {
	// Test parsing from file content directly - TUI integration is tested separately
	validYAML := `title: File Test Form
fields:
  - key: email
    label: Email
    type: input
`

	_, err := config.Parse([]byte(validYAML))
	if err != nil {
		t.Fatalf("Expected no error when parsing YAML from file content, got: %v", err)
	}
}

// TestParseOnly tests only the parsing logic without running the TUI
func TestParseOnly(t *testing.T) {
	// Test valid YAML parsing
	validYAML := `title: Test Form
fields:
  - key: name
    label: Name
    type: input
  - key: email
    label: Email
    type: input
`

	formConfig, err := config.Parse([]byte(validYAML))
	if err != nil {
		t.Fatalf("Expected no error when parsing valid YAML, got: %v", err)
	}

	if formConfig.Title != "Test Form" {
		t.Errorf("Expected title 'Test Form', got: %s", formConfig.Title)
	}

	if len(formConfig.Fields) != 2 {
		t.Errorf("Expected 2 fields, got: %d", len(formConfig.Fields))
	}

	// Test invalid YAML parsing
	invalidYAML := `title: Test Form
fields:
  - key: name
    label: Name
    type: input
  - invalid yaml structure [
`

	_, err = config.Parse([]byte(invalidYAML))
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}

	if !strings.Contains(err.Error(), "erro ao parsear YAML") {
		t.Errorf("Expected error message to contain 'erro ao parsear YAML', got: %v", err)
	}

	// Test missing required fields
	yamlWithoutKey := `fields:
  - label: Name
    type: input
`

	_, err = config.Parse([]byte(yamlWithoutKey))
	if err == nil {
		t.Fatal("Expected error for missing key field, got nil")
	}

	if !strings.Contains(err.Error(), "campo 1 deve ter uma chave") {
		t.Errorf("Expected error message to contain 'campo 1 deve ter uma chave', got: %v", err)
	}
}
