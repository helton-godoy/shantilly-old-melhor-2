package components

import (
	"strings"
	"testing"

	"shantilly/internal/tui"

	"github.com/charmbracelet/lipgloss"
)

func TestProgressIndicator(t *testing.T) {
	theme := &tui.Theme{
		FieldLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
	}
	pi := NewProgressIndicator(theme)

	// Testar com um total > 0 para evitar a ambiguidade do estado 0/0
	pi.UpdateProgress(0, 4)

	if pi.GetProgress() != 0.0 {
		t.Errorf("Expected initial progress 0.0 for 0/4, got %.1f", pi.GetProgress())
	}
	if pi.IsComplete() {
		t.Error("Expected not complete initially for 0/4")
	}

	pi.UpdateProgress(2, 4)
	expectedProgress := 50.0
	if pi.GetProgress() != expectedProgress {
		t.Errorf("Expected progress %.1f, got %.1f", expectedProgress, pi.GetProgress())
	}

	pi.UpdateProgress(4, 4)
	if !pi.IsComplete() {
		t.Error("Expected complete with 100% progress")
	}
	if pi.GetProgress() != 100.0 {
		t.Errorf("Expected progress 100.0, got %.1f", pi.GetProgress())
	}

	render := pi.Render()
	if !strings.Contains(render, "(100%)") {
		t.Errorf("Expected render to contain progress percentage '(100%%)', got '%s'", render)
	}
	if !strings.Contains(render, "4/4") {
		t.Error("Expected render to contain completed/total count")
	}
}

// ... (resto do arquivo de teste inalterado)

func TestProgressIndicatorRenderCompact(t *testing.T) {
	theme := &tui.Theme{
		FieldLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
	}
	pi := NewProgressIndicator(theme)

	pi.UpdateProgress(3, 5)
	render := pi.RenderCompact()
	if render == "" {
		t.Error("Expected non-empty compact render")
	}
	if !strings.Contains(render, "3/5") {
		t.Error("Expected compact render to contain count")
	}
}

func TestProgressIndicatorZeroTotal(t *testing.T) {
	theme := &tui.Theme{
		FieldLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
	}
	pi := NewProgressIndicator(theme)

	pi.UpdateProgress(0, 0)
	if pi.GetProgress() != 100.0 {
		t.Errorf("Expected progress 100.0 with zero total, got %.1f", pi.GetProgress())
	}
	if !pi.IsComplete() {
		t.Error("Expected complete with zero total fields")
	}
}
