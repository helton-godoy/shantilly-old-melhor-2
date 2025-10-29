package tui

import (
	"testing"

	"shantilly/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

// TestTUI_Initialization tests that the TUI starts without errors
func TestTUI_Initialization(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	// Test that model creation doesn't panic
	model := NewModel(formConfig)
	if model == nil {
		t.Fatal("Expected model to be created")
	}

	// Test that the model implements the required interface
	var _ tea.Model = model
}

// TestTUI_QuitFlow tests the quit functionality using official Bubble Tea patterns
func TestTUI_QuitFlow(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Test quit functionality directly through Update method (safer than full program)
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := model.Update(quitMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	if cmd == nil {
		t.Fatal("Expected quit command for 'q' key")
	}

	// Verify the model handles quit messages correctly
}

// TestTUI_CtrlCQuit tests Ctrl+C quit functionality
func TestTUI_CtrlCQuit(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Test Update method directly (pure unit test approach)
	quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := model.Update(quitMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	if cmd == nil {
		t.Fatal("Expected quit command for Ctrl+C")
	}
}

// TestTUI_EscQuit tests Esc quit functionality
func TestTUI_EscQuit(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Test Update method directly (pure unit test approach)
	quitMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd := model.Update(quitMsg)

	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}

	if cmd == nil {
		t.Fatal("Expected quit command for Esc")
	}
}

// TestTUI_ViewOutput tests that View method returns content
func TestTUI_ViewOutput(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Test View method
	view := model.View()
	if view == "" {
		t.Error("Expected View to return non-empty string")
	}

	// View should contain form-related content
	// (The exact content depends on the huh.Form implementation)
}

// TestTUI_FormCreation tests that forms are created correctly for different field types
func TestTUI_FormCreation(t *testing.T) {
	tests := []struct {
		name     string
		field    config.Field
		wantForm bool
	}{
		{
			name: "input field",
			field: config.Field{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			wantForm: true,
		},
		{
			name: "textarea field",
			field: config.Field{
				Key:   "description",
				Label: "Description",
				Type:  "textarea",
			},
			wantForm: true,
		},
		{
			name: "select field",
			field: config.Field{
				Key:     "country",
				Label:   "Country",
				Type:    "select",
				Options: []string{"US", "BR", "CA"},
			},
			wantForm: true,
		},
		{
			name: "confirm field",
			field: config.Field{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
			},
			wantForm: true,
		},
		{
			name: "note field",
			field: config.Field{
				Key:   "info",
				Label: "Information",
				Type:  "note",
			},
			wantForm: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formConfig := &config.FormConfig{
				Title:  "Test Form",
				Fields: []config.Field{tt.field},
			}

			multiselectValues := &multiselectData{values: make(map[string]*[]string)}
			form := createForm(formConfig, multiselectValues)
			if tt.wantForm && form == nil {
				t.Errorf("Expected form to be created for field type %s", tt.field.Type)
			}
		})
	}
}

// TestTUI_ModelInterface tests that the model properly implements the tea.Model interface
func TestTUI_ModelInterface(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Test that all required methods exist and work
	if model.Init() != nil {
		t.Error("Init should return nil for this simple model")
	}

	// Test Update with various messages
	_, cmd := model.Update(nil)
	if cmd != nil {
		t.Error("Update with nil message should return nil command")
	}

	// Test View
	view := model.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}
}
