package components

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestHelpText(t *testing.T) {
	theme := &Theme{
		FieldInput: lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
	}
	ht := NewHelpText(theme)

	// Test initial state
	if ht.IsHelpVisible() {
		t.Error("Expected help not visible initially")
	}

	// Test setting help text
	ht.SetHelpText("name", "Enter your full name")
	help, ok := ht.GetHelpText("name")
	if !ok || help != "Enter your full name" {
		t.Errorf("Expected help text 'Enter your full name', got '%s'", help)
	}

	// Test toggle help
	ht.ToggleHelp()
	if !ht.IsHelpVisible() {
		t.Error("Expected help visible after toggle")
	}

	// Test render with help visible
	ht.SetCurrentField("name")
	render := ht.Render()
	if render == "" {
		t.Error("Expected non-empty render with help visible")
	}
	if !contains(render, "Enter your full name") {
		t.Error("Expected render to contain help text")
	}

	// Test render with no current field
	ht.SetCurrentField("")
	render = ht.Render()
	if render == "" {
		t.Error("Expected non-empty render even with no current field")
	}
}

func TestHelpTextRenderFieldHelp(t *testing.T) {
	theme := &Theme{
		FieldInput: lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
	}
	ht := NewHelpText(theme)

	ht.SetHelpText("email", "Use format: user@domain.com")

	render := ht.RenderFieldHelp("email")
	if render == "" {
		t.Error("Expected non-empty field help render")
	}
	if !contains(render, "user@domain.com") {
		t.Error("Expected render to contain help text")
	}

	// Test render for field without help
	render = ht.RenderFieldHelp("phone")
	if render != "" {
		t.Errorf("Expected empty render for field without help, got '%s'", render)
	}
}
