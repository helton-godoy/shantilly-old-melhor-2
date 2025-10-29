package tui

import (
	"encoding/json"
	"strings"
	"testing"

	"shantilly/internal/config"

	"github.com/charmbracelet/huh"
)

const testUserName = "John Doe"

// TestSubmission_BasicFormSubmission tests basic form submission and JSON output
func TestSubmission_BasicFormSubmission(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "Test Form",
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

	// Simulate form completion by setting field values
	m := model.(*Model)
	m.formConfig.Fields[0].Value = testUserName
	m.formConfig.Fields[1].Value = "john@example.com"

	// Test CollectFormData (public method)
	formData := m.CollectFormData()

	// Verify form data structure
	if len(formData) != 2 {
		t.Errorf("Expected 2 fields in form data, got %d", len(formData))
	}

	if formData["name"] != testUserName {
		t.Errorf("Expected name to be '%s', got %v", testUserName, formData["name"])
	}

	if formData["email"] != "john@example.com" {
		t.Errorf("Expected email to be 'john@example.com', got %v", formData["email"])
	}

	// Test JSON serialization
	jsonData, err := json.MarshalIndent(formData, "", "  ")
	if err != nil {
		t.Errorf("Expected JSON serialization to succeed, got error: %v", err)
	}

	// Verify JSON structure
	var jsonResult map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonResult); err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}

	if jsonResult["name"] != "John Doe" {
		t.Errorf("Expected JSON name to be 'John Doe', got %v", jsonResult["name"])
	}
}

// TestSubmission_MultiselectField tests multiselect field handling
func TestSubmission_MultiselectField(t *testing.T) {
	// Create a test form configuration with multiselect
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:     "interests",
				Label:   "Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports"},
				Value:   "Tech, Music",
			},
		},
	}

	model := NewModel(formConfig)

	// Test CollectFormData for multiselect
	m := model.(*Model)
	formData := m.CollectFormData()

	// Verify multiselect data
	interests, ok := formData["interests"]
	if !ok {
		t.Fatal("Expected interests field in form data")
	}

	interestsSlice, ok := interests.([]string)
	if !ok {
		t.Errorf("Expected interests to be []string, got %T", interests)
	}

	expected := []string{"Tech", "Music"}
	if len(interestsSlice) != len(expected) {
		t.Errorf("Expected %d interests, got %d", len(expected), len(interestsSlice))
	}

	for i, expectedInterest := range expected {
		if i >= len(interestsSlice) || interestsSlice[i] != expectedInterest {
			t.Errorf("Expected interest %d to be %s, got %s", i, expectedInterest, interestsSlice[i])
		}
	}
}

// TestSubmission_ConfirmField tests confirm field handling
func TestSubmission_ConfirmField(t *testing.T) {
	// Create a test form configuration with confirm field
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
				Value: "true",
			},
		},
	}

	model := NewModel(formConfig)

	// Test CollectFormData for confirm field
	m := model.(*Model)
	formData := m.CollectFormData()

	// Verify confirm data
	agree, ok := formData["agree"]
	if !ok {
		t.Fatal("Expected agree field in form data")
	}

	agreeBool, ok := agree.(bool)
	if !ok {
		t.Errorf("Expected agree to be bool, got %T", agree)
	}

	if !agreeBool {
		t.Error("Expected agree to be true")
	}
}

// TestSubmission_NoteField tests note field handling
func TestSubmission_NoteField(t *testing.T) {
	// Create a test form configuration with note field
	formConfig := &config.FormConfig{
		Title: "Test Form",
		Fields: []config.Field{
			{
				Key:   "info",
				Label: "Information",
				Type:  "note",
			},
		},
	}

	model := NewModel(formConfig)

	// Test CollectFormData for note field
	m := model.(*Model)
	formData := m.CollectFormData()

	// Verify note data
	info, ok := formData["info"]
	if !ok {
		t.Fatal("Expected info field in form data")
	}

	if info != nil {
		t.Errorf("Expected note field to be nil, got %v", info)
	}
}

// TestSubmission_CompleteFormData tests complete form data collection
func TestSubmission_CompleteFormData(t *testing.T) {
	// Create a comprehensive test form configuration
	formConfig := &config.FormConfig{
		Title: "Complete Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
				Value: "Jane Doe",
			},
			{
				Key:   "description",
				Label: "Description",
				Type:  "textarea",
				Value: "A test description",
			},
			{
				Key:     "country",
				Label:   "Country",
				Type:    "select",
				Options: []string{"US", "BR", "CA"},
				Value:   "BR",
			},
			{
				Key:     "interests",
				Label:   "Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports"},
				Value:   "Tech, Sports",
			},
			{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
				Value: "true",
			},
			{
				Key:   "info",
				Label: "Information",
				Type:  "note",
			},
		},
	}

	model := NewModel(formConfig)

	// Test CollectFormData for complete form
	m := model.(*Model)
	formData := m.CollectFormData()

	// Verify all fields are present
	expectedFields := []string{"name", "description", "country", "interests", "agree", "info"}
	if len(formData) != len(expectedFields) {
		t.Errorf("Expected %d fields, got %d", len(expectedFields), len(formData))
	}

	// Verify each field type and value
	if formData["name"] != "Jane Doe" {
		t.Errorf("Expected name to be 'Jane Doe', got %v", formData["name"])
	}

	if formData["description"] != "A test description" {
		t.Errorf("Expected description to be 'A test description', got %v", formData["description"])
	}

	if formData["country"] != "BR" {
		t.Errorf("Expected country to be 'BR', got %v", formData["country"])
	}

	interests, ok := formData["interests"].([]string)
	if !ok {
		t.Errorf("Expected interests to be []string, got %T", formData["interests"])
	} else {
		expectedInterests := []string{"Tech", "Sports"}
		if len(interests) != len(expectedInterests) {
			t.Errorf("Expected %d interests, got %d", len(expectedInterests), len(interests))
		}
		for i, expected := range expectedInterests {
			if i >= len(interests) || interests[i] != expected {
				t.Errorf("Expected interest %d to be %s, got %s", i, expected, interests[i])
			}
		}
	}

	if formData["agree"] != true {
		t.Errorf("Expected agree to be true, got %v", formData["agree"])
	}

	if formData["info"] != nil {
		t.Errorf("Expected info to be nil, got %v", formData["info"])
	}
}

// TestSubmission_JSONOutputFormat tests JSON output formatting
func TestSubmission_JSONOutputFormat(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "JSON Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
				Value: "Test User",
			},
			{
				Key:     "tags",
				Label:   "Tags",
				Type:    "multiselect",
				Options: []string{"tag1", "tag2", "tag3"},
				Value:   "tag1, tag3",
			},
		},
	}

	model := NewModel(formConfig)

	// Test JSON output generation
	m := model.(*Model)
	formData := m.CollectFormData()
	jsonData, err := json.MarshalIndent(formData, "", "  ")
	if err != nil {
		t.Errorf("Expected JSON serialization to succeed, got error: %v", err)
	}

	// Verify JSON is properly formatted
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, "Test User") {
		t.Error("Expected JSON to contain 'Test User'")
	}

	if !strings.Contains(jsonStr, "tag1") {
		t.Error("Expected JSON to contain 'tag1'")
	}

	if !strings.Contains(jsonStr, "tag3") {
		t.Error("Expected JSON to contain 'tag3'")
	}

	// Verify JSON is valid
	var jsonResult map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonResult); err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}
}

// TestSubmission_EmptyMultiselect tests empty multiselect handling
func TestSubmission_EmptyMultiselect(t *testing.T) {
	// Create a test form configuration with empty multiselect
	formConfig := &config.FormConfig{
		Title: "Empty Multiselect Test",
		Fields: []config.Field{
			{
				Key:     "empty_interests",
				Label:   "Empty Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports"},
				Value:   "", // Empty value
			},
		},
	}

	model := NewModel(formConfig)

	// Test CollectFormData for empty multiselect
	m := model.(*Model)
	formData := m.CollectFormData()

	// Verify empty multiselect data
	emptyInterests, ok := formData["empty_interests"]
	if !ok {
		t.Fatal("Expected empty_interests field in form data")
	}

	emptyInterestsSlice, ok := emptyInterests.([]string)
	if !ok {
		t.Errorf("Expected empty_interests to be []string, got %T", emptyInterests)
	}

	if len(emptyInterestsSlice) != 0 {
		t.Errorf("Expected empty multiselect to have 0 items, got %d", len(emptyInterestsSlice))
	}
}

// TestSubmission_FormStateHandling tests form state handling for submission
func TestSubmission_FormStateHandling(t *testing.T) {
	// Create a test form configuration
	formConfig := &config.FormConfig{
		Title: "State Test Form",
		Fields: []config.Field{
			{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
		},
	}

	model := NewModel(formConfig)

	// Initially, form should not be in completed state
	m := model.(*Model)
	if m.form.State == huh.StateCompleted {
		t.Error("Expected form to not be completed initially")
	}

	// Test that handleSubmission command is created when form is completed
	// Note: This is a basic test - in reality, the form state would be set by huh.Form
	// when the user completes all fields and submits
}
