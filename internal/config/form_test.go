package config

import (
	"strings"
	"testing"
)

func TestFormConfig_StructValidation(t *testing.T) {
	// Test that FormConfig struct can be created
	form := FormConfig{
		Title: "Test Form",
		Fields: []Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	if form.Title != "Test Form" {
		t.Errorf("Expected title 'Test Form', got '%s'", form.Title)
	}

	if len(form.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(form.Fields))
	}

	if form.Fields[0].Key != "name" {
		t.Errorf("Expected field key 'name', got '%s'", form.Fields[0].Key)
	}
}

func TestParse_ValidYAML(t *testing.T) {
	yamlData := `
title: Test Form
fields:
  - key: name
    label: Name
    type: input
    placeholder: Enter your name
  - key: email
    label: Email
    type: input
    placeholder: Enter your email
`

	form, err := Parse([]byte(yamlData))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if form.Title != "Test Form" {
		t.Errorf("Expected title 'Test Form', got '%s'", form.Title)
	}

	if len(form.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(form.Fields))
	}

	// Check first field
	if form.Fields[0].Key != "name" {
		t.Errorf("Expected first field key 'name', got '%s'", form.Fields[0].Key)
	}
	if form.Fields[0].Type != "input" {
		t.Errorf("Expected first field type 'input', got '%s'", form.Fields[0].Type)
	}

	// Check second field
	if form.Fields[1].Key != "email" {
		t.Errorf("Expected second field key 'email', got '%s'", form.Fields[1].Key)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	invalidYAML := `
title: Test Form
fields:
  - key: name
    label: Name
    type: input
  - invalid yaml structure [
`

	_, err := Parse([]byte(invalidYAML))
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}
}

func TestParse_EmptyYAML(t *testing.T) {
	emptyYAML := ""

	_, err := Parse([]byte(emptyYAML))
	if err == nil {
		t.Fatal("Expected error for empty YAML, got nil")
	}
}

func TestParse_ValidationErrors(t *testing.T) {
	// Test missing key
	yamlWithoutKey := `
fields:
  - label: Name
    type: input
`
	_, err := Parse([]byte(yamlWithoutKey))
	if err == nil {
		t.Fatal("Expected error for missing key, got nil")
	}
	if !strings.Contains(err.Error(), "campo 1 deve ter uma chave") {
		t.Errorf("Expected error about missing key, got: %v", err)
	}

	// Test missing label
	yamlWithoutLabel := `
fields:
  - key: name
    type: input
`
	_, err = Parse([]byte(yamlWithoutLabel))
	if err == nil {
		t.Fatal("Expected error for missing label, got nil")
	}
	if !strings.Contains(err.Error(), "campo 1 deve ter um rótulo") {
		t.Errorf("Expected error about missing label, got: %v", err)
	}

	// Test missing type
	yamlWithoutType := `
fields:
  - key: name
    label: Name
`
	_, err = Parse([]byte(yamlWithoutType))
	if err == nil {
		t.Fatal("Expected error for missing type, got nil")
	}
	if !strings.Contains(err.Error(), "campo 1 deve ter um tipo") {
		t.Errorf("Expected error about missing type, got: %v", err)
	}

	// Test invalid type
	yamlWithInvalidType := `
fields:
  - key: name
    label: Name
    type: invalid_type
`
	_, err = Parse([]byte(yamlWithInvalidType))
	if err == nil {
		t.Fatal("Expected error for invalid type, got nil")
	}
	if !strings.Contains(err.Error(), "tipo de campo 'invalid_type' inválido") {
		t.Errorf("Expected error about invalid type, got: %v", err)
	}
}

func TestParse_MinimalValidYAML(t *testing.T) {
	minimalYAML := `
fields:
  - key: test
    label: Test
    type: input
`

	form, err := Parse([]byte(minimalYAML))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(form.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(form.Fields))
	}

	if form.Fields[0].Key != "test" {
		t.Errorf("Expected field key 'test', got '%s'", form.Fields[0].Key)
	}
}
