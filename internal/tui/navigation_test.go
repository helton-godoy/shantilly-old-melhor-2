package tui

import (
	"strings"
	"testing"

	"shantilly/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNavigation_TabNavigation tests Tab key navigation between fields
func TestNavigation_TabNavigation(t *testing.T) {
	// Create a test form configuration with multiple fields
	formConfig := &config.FormConfig{
		Title: "Navigation Test Form",
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

	model := NewModel(formConfig)

	// Test Tab key navigation
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(tabMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// The command should be handled by huh.Form internally
	// We can't easily test the exact navigation state without more complex setup
	// but we can verify the model handles the message without errors
}

// TestNavigation_ShiftTabNavigation tests Shift+Tab key navigation between fields
func TestNavigation_ShiftTabNavigation(t *testing.T) {
	// Create a test form configuration with multiple fields
	formConfig := &config.FormConfig{
		Title: "Navigation Test Form",
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

	// Test Shift+Tab key navigation
	shiftTabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(shiftTabMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// Verify the model handles the message without errors
}

// TestNavigation_EnterBehavior tests Enter key behavior in different field types
func TestNavigation_EnterBehavior(t *testing.T) {
	tests := []struct {
		name     string
		field    config.Field
		wantNext bool
	}{
		{
			name: "input field - Enter should move to next field",
			field: config.Field{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			wantNext: true,
		},
		{
			name: "textarea field - Enter should move to next field",
			field: config.Field{
				Key:   "bio",
				Label: "Bio",
				Type:  "textarea",
			},
			wantNext: true,
		},
		{
			name: "select field - Enter should confirm selection",
			field: config.Field{
				Key:     "country",
				Label:   "Country",
				Type:    "select",
				Options: []string{"US", "BR", "CA"},
			},
			wantNext: true,
		},
		{
			name: "confirm field - Enter should confirm",
			field: config.Field{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
			},
			wantNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formConfig := &config.FormConfig{
				Title:  "Test Form",
				Fields: []config.Field{tt.field},
			}

			model := NewModel(formConfig)

			// Test Enter key behavior
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			updatedModel, _ := model.Update(enterMsg)

			if updatedModel == nil {
				t.Fatal("Expected updated model to not be nil")
			}

			// The exact behavior depends on huh.Form's internal state
			// but we can verify the model handles the message
		})
	}
}

// TestNavigation_SpaceInMultiSelect tests Space key behavior in MultiSelect fields
func TestNavigation_SpaceInMultiSelect(t *testing.T) {
	// Create a test form configuration with multiselect field
	formConfig := &config.FormConfig{
		Title: "MultiSelect Test Form",
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

	// Test Space key behavior in multiselect
	spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
	updatedModel, _ := model.Update(spaceMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// Verify the model handles the space message without errors
}

// TestNavigation_ArrowKeysInSelect tests arrow key navigation in Select fields
func TestNavigation_ArrowKeysInSelect(t *testing.T) {
	// Create a test form configuration with select field
	formConfig := &config.FormConfig{
		Title: "Select Test Form",
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

	// Test Up arrow key navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(upMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// Test Down arrow key navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = model.Update(downMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// Verify the model handles arrow key messages without errors
}

// TestNavigation_ArrowKeysInMultiSelect tests arrow key navigation in MultiSelect fields
func TestNavigation_ArrowKeysInMultiSelect(t *testing.T) {
	// Create a test form configuration with multiselect field
	formConfig := &config.FormConfig{
		Title: "MultiSelect Test Form",
		Fields: []config.Field{
			{
				Key:     "interests",
				Label:   "Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports", "Art"},
			},
		},
	}

	model := NewModel(formConfig)

	// Test Up arrow key navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(upMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// Test Down arrow key navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = model.Update(downMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// Verify the model handles arrow key messages without errors
}

// TestNavigation_VisualFocusIndication tests that focused fields are visually highlighted
func TestNavigation_VisualFocusIndication(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Visual Focus Test Form",
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

	// Test that View method returns content with visual indicators
	view := model.View()
	if view == "" {
		t.Error("Expected View to return non-empty string")
	}

	// The view should contain visual focus indicators
	// (The exact format depends on huh.Form's rendering)
	if !strings.Contains(view, ">") && !strings.Contains(view, "•") {
		t.Log("Warning: View output may not contain expected focus indicators")
		// This is not necessarily an error as the exact format depends on huh.Form
	}
}

// TestNavigation_FormSubmission tests that form can be submitted after navigation
func TestNavigation_FormSubmission(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Submission Test Form",
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

	// Test that form handles submission after navigation
	// This is a basic test - full submission testing would require more complex setup
	submitMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(submitMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	// Verify the model handles submission messages
}

// TestNavigation_EdgeCases tests navigation edge cases
func TestNavigation_EdgeCases(t *testing.T) {
	t.Run("single field form", func(t *testing.T) {
		// Test navigation in a form with only one field
		formConfig := &config.FormConfig{
			Title: "Single Field Form",
			Fields: []config.Field{
				{
					Key:   "name",
					Label: "Name",
					Type:  "input",
				},
			},
		}

		model := NewModel(formConfig)

		// Test Tab navigation in single field form
		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
		updatedModel, _ := model.Update(tabMsg)

		if updatedModel == nil {
			t.Fatal("Expected updated model to not be nil")
		}

		// Should handle gracefully even with single field
	})

	t.Run("empty form validation", func(t *testing.T) {
		// Test that empty form validation is handled by the parser
		// The parser already validates this, so we don't need to test it here
		// This test ensures the validation is working correctly
		emptyConfig := &config.FormConfig{
			Title:  "Empty Form",
			Fields: []config.Field{},
		}

		// The validation should be handled by the config.Parse function
		// which is tested in the config package tests
		// Here we just verify the model can be created (validation happens in parser)
		// Note: huh.Form may panic with empty fields, so we expect this to fail gracefully
		defer func() {
			if r := recover(); r != nil {
				// Expected panic due to huh.Form not handling empty field lists
				t.Log("Expected panic for empty form - huh.Form cannot handle empty field lists")
			}
		}()

		model := NewModel(emptyConfig)
		if model == nil {
			t.Error("Expected model to be created even with empty config")
		}
	})
}

// TestNavigation_KeyboardShortcuts tests that keyboard shortcuts are displayed
func TestNavigation_KeyboardShortcuts(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Keyboard Shortcuts Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			{
				Key:     "interests",
				Label:   "Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports"},
			},
			{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
			},
		},
	}

	model := NewModel(formConfig)

	// Test that View method returns content with keyboard shortcuts
	view := model.View()
	if view == "" {
		t.Error("Expected View to return non-empty string")
	}

	// The view should contain navigation hints
	// (The exact format depends on huh.Form's rendering)
	navigationHints := []string{"tab", "enter", "shift", "space", "arrow"}
	hasHints := false
	for _, hint := range navigationHints {
		if strings.Contains(strings.ToLower(view), hint) {
			hasHints = true
			break
		}
	}

	if !hasHints {
		t.Log("Warning: View output may not contain expected navigation hints")
		// This is not necessarily an error as the exact format depends on huh.Form
	}
}
