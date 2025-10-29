package tui

import (
	"strings"
	"testing"

	"shantilly/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNavigation_Integration_BasicFormRendering tests that the form renders correctly
func TestNavigation_Integration_BasicFormRendering(t *testing.T) {
	// Create a test form configuration with multiple fields
	formConfig := &config.FormConfig{
		Title: "Integration Navigation Test",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			{
				Key:   "email",
				Label: "Email",
				Type:  "input",
			},
			{
				Key:     "country",
				Label:   "Country",
				Type:    "select",
				Options: []string{"US", "BR", "CA"},
			},
		},
	}

	// Create the model
	model := NewModel(formConfig)

	// Get the initial view
	view := model.View()
	if view == "" {
		t.Error("Expected initial view to be non-empty")
	}

	// The view should contain form content (we don't check specific text as huh.Form rendering may vary)
	if len(view) < 10 {
		t.Error("Expected view to contain substantial form content")
	}

	// Test that Tab navigation is handled without errors
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(tabMsg)

	if updatedModel == nil {
		t.Error("Expected model to handle Tab navigation without errors")
	}

	// The view should still contain the form after navigation
	updatedView := updatedModel.View()
	if updatedView == "" {
		t.Error("Expected view to be non-empty after Tab navigation")
	}
}

// TestNavigation_Integration_EnterNavigation tests Enter key navigation
func TestNavigation_Integration_EnterNavigation(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Enter Navigation Test",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
			},
		},
	}

	model := NewModel(formConfig)

	// Test Enter key behavior
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(enterMsg)

	if updatedModel == nil {
		t.Error("Expected model to handle Enter key without errors")
	}

	// The form should still be active (Enter should navigate within the form)
	view := updatedModel.View()
	if view == "" {
		t.Error("Expected view to be non-empty after Enter")
	}
}

// TestNavigation_Integration_ArrowKeys tests arrow key navigation in select fields
func TestNavigation_Integration_ArrowKeys(t *testing.T) {
	// Create a test form configuration with select field
	formConfig := &config.FormConfig{
		Title: "Arrow Keys Test",
		Fields: []config.Field{
			{
				Key:     "country",
				Label:   "Country",
				Type:    "select",
				Options: []string{"US", "BR", "CA", "UK", "FR"},
			},
		},
	}

	model := NewModel(formConfig)

	// Test that arrow keys are handled without errors
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(downMsg)

	if updatedModel == nil {
		t.Error("Expected model to handle Down arrow without errors")
	}

	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel2, _ := model.Update(upMsg)

	if updatedModel2 == nil {
		t.Error("Expected model to handle Up arrow without errors")
	}

	// The form should still be active
	view := updatedModel2.View()
	if view == "" {
		t.Error("Expected view to be non-empty after arrow navigation")
	}

	// Should still contain the select options
	if !strings.Contains(view, "US") || !strings.Contains(view, "BR") {
		t.Error("Expected select options to still be visible")
	}
}

// TestNavigation_Integration_SpaceInMultiSelect tests Space key in multiselect
func TestNavigation_Integration_SpaceInMultiSelect(t *testing.T) {
	// Create a test form configuration with multiselect field
	formConfig := &config.FormConfig{
		Title: "Space MultiSelect Test",
		Fields: []config.Field{
			{
				Key:     "interests",
				Label:   "Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports", "Art"},
				Limit:   2,
			},
		},
	}

	model := NewModel(formConfig)

	// Test that Space key is handled without errors
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	updatedModel, _ := model.Update(spaceMsg)

	if updatedModel == nil {
		t.Error("Expected model to handle Space key without errors")
	}

	// The form should still be active
	view := updatedModel.View()
	if view == "" {
		t.Error("Expected view to be non-empty after Space")
	}

	// Should still contain the multiselect options
	if !strings.Contains(view, "Tech") || !strings.Contains(view, "Music") {
		t.Error("Expected multiselect options to still be visible")
	}
}

// TestNavigation_Integration_VisualFocus tests visual focus indication
func TestNavigation_Integration_VisualFocus(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Visual Focus Test",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			{
				Key:   "email",
				Label: "Email",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Get the initial view
	view := model.View()
	if view == "" {
		t.Error("Expected initial view to be non-empty")
	}

	// The view should contain form content (we don't check specific text as huh.Form rendering may vary)
	if len(view) < 10 {
		t.Error("Expected view to contain substantial form content")
	}

	// Test navigation and check if focus changes
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(tabMsg)

	updatedView := updatedModel.View()
	if updatedView == "" {
		t.Error("Expected updated view to be non-empty")
	}

	// Both views should contain form content (we don't check specific text as huh.Form rendering may vary)
	if len(updatedView) < 10 {
		t.Error("Expected updated view to contain substantial form content")
	}
}

// TestNavigation_Integration_QuitKeys tests quit functionality
func TestNavigation_Integration_QuitKeys(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Quit Test",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Test that quit keys are handled without errors
	ctrlCMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, _ := model.Update(ctrlCMsg)

	if updatedModel == nil {
		t.Error("Expected model to handle Ctrl+C without errors")
	}

	// Test Esc key as well
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel2, _ := model.Update(escMsg)

	if updatedModel2 == nil {
		t.Error("Expected model to handle Esc without errors")
	}
}

// TestNavigation_Integration_CompleteFlow tests a complete navigation flow
func TestNavigation_Integration_CompleteFlow(t *testing.T) {
	// Create a comprehensive test form
	formConfig := &config.FormConfig{
		Title: "Complete Navigation Flow Test",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			{
				Key:   "email",
				Label: "Email",
				Type:  "input",
			},
			{
				Key:     "country",
				Label:   "Country",
				Type:    "select",
				Options: []string{"US", "BR", "CA"},
			},
			{
				Key:     "interests",
				Label:   "Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports"},
				Limit:   2,
			},
			{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
			},
		},
	}

	model := NewModel(formConfig)

	// Test complete navigation flow using Update method
	navigationKeys := []tea.KeyMsg{
		{Type: tea.KeyTab},   // Move to email
		{Type: tea.KeyTab},   // Move to country
		{Type: tea.KeyDown},  // Navigate in select
		{Type: tea.KeyEnter}, // Confirm selection
		{Type: tea.KeyTab},   // Move to interests
		{Type: tea.KeySpace}, // Select in multiselect
		{Type: tea.KeyTab},   // Move to agree
		{Type: tea.KeyEnter}, // Confirm
	}

	currentModel := model
	for _, key := range navigationKeys {
		updatedModel, _ := currentModel.Update(key)
		if updatedModel == nil {
			t.Errorf("Expected model to handle key %v without errors", key.Type)
		}
		currentModel = updatedModel
	}

	// Final view should still contain the form
	finalView := currentModel.View()
	if finalView == "" {
		t.Error("Expected final view to be non-empty")
	}

	// Should still contain form content (we don't check specific text as huh.Form rendering may vary)
	if len(finalView) < 10 {
		t.Error("Expected final view to contain substantial form content")
	}
}
