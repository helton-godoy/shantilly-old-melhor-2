package tui

import (
	"strings"
	"testing"

	"shantilly/internal/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestNewModel(t *testing.T) {
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

	// Test model creation
	model := NewModel(formConfig)

	// Verify model is not nil
	if model == nil {
		t.Fatal("Expected model to be created, got nil")
	}

	// Verify model implements tea.Model interface
	var _ tea.Model = model
}

func TestModelInit(t *testing.T) {
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

	// Test Init method
	cmd := model.Init()
	if cmd != nil {
		t.Errorf("Expected Init to return nil command, got: %v", cmd)
	}
}

func TestModelUpdate(t *testing.T) {
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

	// Test Update with nil message
	updatedModel, cmd := model.Update(nil)
	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}
	if cmd != nil {
		t.Errorf("Expected no command for nil message, got: %v", cmd)
	}

	// Test Update with quit key
	quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd = model.Update(quitMsg)
	if updatedModel == nil {
		t.Fatal("Expected updated model to not be nil")
	}
	if cmd == nil {
		t.Fatal("Expected quit command for Ctrl+C")
	}

	// Verify it's a quit command by checking the command type
	// Note: This is a basic test, in a real scenario we'd check the command type
}

func TestModelView(t *testing.T) {
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

	// View should contain the form title or field label
	// (This is a basic check, the actual content depends on the huh.Form implementation)
}

func TestCreateForm(t *testing.T) {
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

	// Test form creation
	multiselectValues := &multiselectData{values: make(map[string]*[]string)}
	form := createForm(formConfig, multiselectValues)
	if form == nil {
		t.Fatal("Expected form to be created, got nil")
	}
}

func TestCreateFormWithDifferentFieldTypes(t *testing.T) {
	tests := []struct {
		name    string
		field   config.Field
		wantErr bool
	}{
		{
			name: "input field",
			field: config.Field{
				Key:   "name",
				Label: "Name",
				Type:  "input",
			},
			wantErr: false,
		},
		{
			name: "textarea field",
			field: config.Field{
				Key:   "description",
				Label: "Description",
				Type:  "textarea",
			},
			wantErr: false,
		},
		{
			name: "select field",
			field: config.Field{
				Key:     "country",
				Label:   "Country",
				Type:    "select",
				Options: []string{"US", "BR", "CA"},
			},
			wantErr: false,
		},
		{
			name: "multiselect field",
			field: config.Field{
				Key:     "interests",
				Label:   "Interests",
				Type:    "multiselect",
				Options: []string{"Tech", "Music", "Sports"},
				Value:   "Tech, Music",
				Limit:   2,
			},
			wantErr: false,
		},
		{
			name: "confirm field",
			field: config.Field{
				Key:   "agree",
				Label: "Agree",
				Type:  "confirm",
			},
			wantErr: false,
		},
		{
			name: "note field",
			field: config.Field{
				Key:   "info",
				Label: "Information",
				Type:  "note",
			},
			wantErr: false,
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
			if form == nil {
				t.Errorf("Expected form to be created for field type %s", tt.field.Type)
			}
		})
	}
}

// TestThemeCreation testa a criação do tema padrão
func TestThemeCreation(t *testing.T) {
	theme := DefaultTheme()
	if theme == nil {
		t.Fatal("Expected theme to be created, got nil")
	}

	// Verificar se os estilos foram inicializados
	if theme.Container.GetPaddingTop() != 1 {
		t.Errorf("Expected container padding top to be 1, got %d", theme.Container.GetPaddingTop())
	}

	if theme.Container.GetPaddingLeft() != 2 {
		t.Errorf("Expected container padding left to be 2, got %d", theme.Container.GetPaddingLeft())
	}
}

// TestWindowSizeHandling testa o tratamento de mensagens de redimensionamento
func TestWindowSizeHandling(t *testing.T) {
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

	// Verificar dimensões iniciais
	concreteModel := model.(*Model)
	if concreteModel.terminalWidth != 80 {
		t.Errorf("Expected initial terminal width to be 80, got %d", concreteModel.terminalWidth)
	}

	if concreteModel.terminalHeight != 24 {
		t.Errorf("Expected initial terminal height to be 24, got %d", concreteModel.terminalHeight)
	}

	// Testar atualização de dimensões
	windowSizeMsg := tea.WindowSizeMsg{Width: 120, Height: 30}
	updatedModel, _ := model.Update(windowSizeMsg)

	updatedConcreteModel := updatedModel.(*Model)
	if updatedConcreteModel.terminalWidth != 120 {
		t.Errorf("Expected terminal width to be updated to 120, got %d", updatedConcreteModel.terminalWidth)
	}

	if updatedConcreteModel.terminalHeight != 30 {
		t.Errorf("Expected terminal height to be updated to 30, got %d", updatedConcreteModel.terminalHeight)
	}
}

// TestViewWithCentralization testa a centralização no método View
func TestViewWithCentralization(t *testing.T) {
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
	concreteModel := model.(*Model)

	// Testar com largura suficiente para centralização
	concreteModel.terminalWidth = 100
	view := model.View()

	if view == "" {
		t.Error("Expected View to return non-empty string")
	}

	// Verificar se contém caracteres de borda (do container style)
	if !strings.Contains(view, "╭") && !strings.Contains(view, "┌") {
		t.Error("Expected view to contain border characters from container styling")
	}
}

// TestViewSmallTerminal testa o comportamento em terminais pequenos
func TestViewSmallTerminal(t *testing.T) {
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
	concreteModel := model.(*Model)

	// Testar com largura muito pequena
	concreteModel.terminalWidth = 30
	view := model.View()

	if view == "" {
		t.Error("Expected View to return non-empty string even in small terminal")
	}

	// Mesmo em terminal pequeno, deve haver algum conteúdo
	lines := strings.Split(view, "\n")
	if len(lines) == 0 {
		t.Error("Expected view to have at least one line")
	}
}

// TestLayoutResponsiveness testa a responsividade do layout
func TestLayoutResponsiveness(t *testing.T) {
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
	concreteModel := model.(*Model)

	testCases := []struct {
		name     string
		width    int
		expected bool // se deve haver centralização
	}{
		{"Wide terminal", 120, true},
		{"Medium terminal", 80, false}, // container é mais largo que 80
		{"Narrow terminal", 50, false}, // muito estreito para centralizar
		{"Very narrow", 30, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			concreteModel.terminalWidth = tc.width
			view := model.View()

			if view == "" {
				t.Errorf("Expected non-empty view for width %d", tc.width)
			}

			// Verificar largura do container
			containerStyle := concreteModel.theme.Container
			formView := concreteModel.form.View()
			containerView := containerStyle.Render(formView)
			containerWidth := lipgloss.Width(containerView)

			// Se o container for menor que a largura do terminal, deve centralizar
			shouldCenter := containerWidth < tc.width

			if shouldCenter != tc.expected {
				t.Errorf("For width %d, expected centering=%v but got centering=%v (containerWidth=%d)",
					tc.width, tc.expected, shouldCenter, containerWidth)
			}
		})
	}
}
