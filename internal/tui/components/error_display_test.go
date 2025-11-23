package components

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestErrorDisplay(t *testing.T) {
	theme := &Theme{
		FieldError: lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
	}
	ed := NewErrorDisplay(theme)

	// Test initial state
	if ed.HasErrors() {
		t.Error("Expected no errors initially")
	}

	// Test setting error
	ed.SetError("name", "Name is required")
	if !ed.HasErrors() {
		t.Error("Expected errors after setting error")
	}

	// Test getting errors
	errors := ed.GetErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
	if errors["name"] != "Name is required" {
		t.Errorf("Expected error message 'Name is required', got '%s'", errors["name"])
	}

	// Test clearing error
	ed.ClearError("name")
	if ed.HasErrors() {
		t.Error("Expected no errors after clearing")
	}
}

func TestErrorDisplayRender(t *testing.T) {
	theme := &Theme{
		FieldError: lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
	}
	ed := NewErrorDisplay(theme)

	// Test render with no errors
	render := ed.Render()
	if render != "" {
		t.Errorf("Expected empty render, got '%s'", render)
	}

	// Test render with errors
	ed.SetError("email", "Invalid email format")
	render = ed.Render()
	if render == "" {
		t.Error("Expected non-empty render with errors")
	}
	if !contains(render, "email") {
		t.Error("Expected render to contain field name")
	}
	if !contains(render, "Invalid email format") {
		t.Error("Expected render to contain error message")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
