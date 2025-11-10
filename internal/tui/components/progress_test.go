package components

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestProgressIndicator(t *testing.T) {
	theme := &Theme{
		FieldLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
	}
	pi := NewProgressIndicator(theme)

	// Test initial state
	if pi.GetProgress() != 0.0 {
		t.Errorf("Expected initial progress 0.0, got %.1f", pi.GetProgress())
	}
	if pi.IsComplete() {
		t.Error("Expected not complete initially")
	}

	// Test update progress
	pi.UpdateProgress(2, 4) // 2 completed out of 4 total
	expectedProgress := 50.0
	if pi.GetProgress() != expectedProgress {
		t.Errorf("Expected progress %.1f, got %.1f", expectedProgress, pi.GetProgress())
	}
	if pi.IsComplete() {
		t.Error("Expected not complete with 50% progress")
	}

	// Test complete progress
	pi.UpdateProgress(4, 4)
	if !pi.IsComplete() {
		t.Error("Expected complete with 100% progress")
	}
	if pi.GetProgress() != 100.0 {
		t.Errorf("Expected progress 100.0, got %.1f", pi.GetProgress())
	}

	// Test render
	render := pi.Render()
	if render == "" {
		t.Error("Expected non-empty render")
	}
	if !contains(render, "100.0%") {
		t.Error("Expected render to contain progress percentage")
	}
	if !contains(render, "4/4") {
		t.Error("Expected render to contain completed/total count")
	}
}

func TestProgressIndicatorRenderCompact(t *testing.T) {
	theme := &Theme{
		FieldLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
	}
	pi := NewProgressIndicator(theme)

	pi.UpdateProgress(3, 5)
	render := pi.RenderCompact()
	if render == "" {
		t.Error("Expected non-empty compact render")
	}
	if !contains(render, "60%") {
		t.Error("Expected compact render to contain percentage")
	}
}

func TestProgressIndicatorZeroTotal(t *testing.T) {
	theme := &Theme{
		FieldLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
	}
	pi := NewProgressIndicator(theme)

	// Test with zero total fields
	pi.UpdateProgress(0, 0)
	if pi.GetProgress() != 100.0 {
		t.Errorf("Expected progress 100.0 with zero total, got %.1f", pi.GetProgress())
	}
	if !pi.IsComplete() {
		t.Error("Expected complete with zero total fields")
	}
}
