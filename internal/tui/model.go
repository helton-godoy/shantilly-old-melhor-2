package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"shantilly/internal/config"
	"shantilly/internal/tui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// multiselectData armazena os valores dos campos multiselect
type multiselectData struct {
	values map[string]*[]string
}

// Model representa o estado da aplicação TUI.
type Model struct {
	form              *huh.Form
	formConfig        *config.FormConfig
	multiselectValues *multiselectData // Armazena valores de campos multiselect
	theme             *Theme           // Tema centralizado para styling
	terminalWidth     int              // Largura atual do terminal
	terminalHeight    int              // Altura atual do terminal
}

const (
	defaultTerminalWidth  = 80 // Largura padrão do terminal
	defaultTerminalHeight = 24 // Altura padrão do terminal
)

// NewModel cria um novo modelo com a configuração do formulário.
func NewModel(formConfig *config.FormConfig) tea.Model {
	multiselectValues := &multiselectData{values: make(map[string]*[]string)}
	form := createForm(formConfig, multiselectValues)
	return &Model{
		form:              form,
		formConfig:        formConfig,
		multiselectValues: multiselectValues,
		theme:             DefaultTheme(),
		terminalWidth:     defaultTerminalWidth,  // Largura padrão inicial
		terminalHeight:    defaultTerminalHeight, // Altura padrão inicial
	}
}

// Init é o primeiro comando a ser executado quando o programa inicia.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update é chamado quando mensagens (eventos) são recebidas.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size messages
	if windowSizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.terminalWidth = windowSizeMsg.Width
		m.terminalHeight = windowSizeMsg.Height
	}

	// Handle form messages
	var cmd tea.Cmd
	var updatedForm tea.Model
	updatedForm, cmd = m.form.Update(msg)
	if form, ok := updatedForm.(*huh.Form); ok {
		m.form = form

		// Check if form is completed (submitted)
		if m.form.State == huh.StateCompleted {
			return m, tea.Batch(cmd, m.handleSubmission())
		}
	}

	// Handle other messages
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	// Return any commands from form updates
	return m, cmd
}

// View renderiza a interface do usuário.
func (m *Model) View() string {
	formView := m.form.View()

	// Aplicar container styling com padding
	containerStyle := m.theme.Container
	containerView := containerStyle.Render(formView)

	// Centralizar horizontalmente se houver espaço suficiente
	if m.terminalWidth > 0 {
		containerWidth := lipgloss.Width(containerView)
		if containerWidth < m.terminalWidth {
			containerView = lipgloss.PlaceHorizontal(m.terminalWidth, lipgloss.Center, containerView)
		}
	}

	return containerView
}

// Styles define os estilos para a TUI.
type Styles struct {
	Title lipgloss.Style
	Body  lipgloss.Style
}

// DefaultStyles retorna os estilos padrão.
func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Align(lipgloss.Center),
		Body: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center),
	}
}

// createForm cria um huh.Form a partir da configuração e retorna os valores multiselect.
func createForm(formConfig *config.FormConfig, multiselectValues *multiselectData) *huh.Form {
	var fields []huh.Field

	// Create fields based on configuration
	for i := range formConfig.Fields {
		field := &formConfig.Fields[i] // Use pointer to avoid copying
		switch field.Type {
		case "input":
			fields = append(fields, huh.NewInput().
				Key(field.Key).
				Title(field.Label).
				Placeholder(field.Placeholder).
				Value(&field.Value))
		case "number":
			numberInput := components.NewNumberInput(field)
			fields = append(fields, numberInput.Field())
		case "date":
			dateInput := components.NewDateInput(field)
			fields = append(fields, dateInput.Field())
		case "file":
			fileInput := components.NewFileInput(field)
			fields = append(fields, fileInput.Field())
		case "textarea":
			fields = append(fields, huh.NewText().
				Key(field.Key).
				Title(field.Label).
				Placeholder(field.Placeholder).
				Value(&field.Value))
		case "select":
			fields = append(fields, huh.NewSelect[string]().
				Key(field.Key).
				Title(field.Label).
				Options(huh.NewOptions(field.Options...)...).
				Value(&field.Value))
		case "multiselect":
			var value []string
			if field.Value != "" {
				// Handle comma-separated values for multiselect
				// This allows YAML to specify default selected options
				value = strings.Split(field.Value, ",")
				// Trim spaces from each value
				for i, v := range value {
					value[i] = strings.TrimSpace(v)
				}
			}
			// Store the value in the multiselectValues map
			multiselectValues.values[field.Key] = &value
			fields = append(fields, huh.NewMultiSelect[string]().
				Key(field.Key).
				Title(field.Label).
				Options(huh.NewOptions(field.Options...)...).
				Value(multiselectValues.values[field.Key]).
				Limit(field.Limit))
		case "confirm":
			var value bool
			if field.Value == "true" {
				value = true
			}
			fields = append(fields, huh.NewConfirm().
				Key(field.Key).
				Title(field.Label).
				Affirmative(field.Affirmative).
				Negative(field.Negative).
				Value(&value))
		case "note":
			fields = append(fields, huh.NewNote().
				Title(field.Label).
				Description(field.Detail))
		}
	}

	// Create a group with all fields
	group := huh.NewGroup(fields...)

	// Create form with the group
	form := huh.NewForm(group)

	if formConfig.Title != "" {
		// Note: huh.Form doesn't have WithTitle method, title is handled by the group
		group = group.Title(formConfig.Title)
		form = huh.NewForm(group)
	}

	return form
}

// handleSubmission processa a submissão do formulário e gera a saída JSON.
func (m *Model) handleSubmission() tea.Cmd {
	return func() tea.Msg {
		// Validar todos os campos antes de submeter
		validator := config.NewFieldValidator()
		formData := m.collectFormData()

		// Executar validação
		validationResults := validator.ValidateForm(m.formConfig, formData)

		// Verificar se há erros de validação
		if validator.HasErrors(validationResults) {
			// Exibir erros de validação e impedir submissão
			m.displayValidationErrors(validationResults)
			return nil // Não submeter, manter formulário ativo
		}

		// Limpar a tela
		fmt.Print("\033[2J\033[H")

		// Converter para JSON e imprimir
		jsonData, err := json.MarshalIndent(formData, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao formatar JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonData))
		os.Exit(0)
		return nil
	}
}

// CollectFormData coleta os valores dos campos do formulário (método público para testes).
func (m *Model) CollectFormData() map[string]interface{} {
	return m.collectFormData()
}

// collectFormData coleta os valores dos campos do formulário.
func (m *Model) collectFormData() map[string]interface{} {
	formData := make(map[string]interface{})

	// Iterar sobre os campos da configuração para obter os valores atualizados
	// Os valores são atualizados através dos ponteiros passados para os campos huh
	for i := range m.formConfig.Fields {
		field := &m.formConfig.Fields[i]
		key := field.Key
		switch field.Type {
		case "input", "textarea", "select", "number", "date", "file":
			// Para string fields, o valor é armazenado em field.Value
			formData[key] = field.Value
		case "multiselect":
			// Para multiselect, os valores são armazenados no multiselectValues map
			if multiselectValue, exists := m.multiselectValues.values[key]; exists && multiselectValue != nil {
				formData[key] = *multiselectValue
			} else {
				formData[key] = []string{}
			}
		case "confirm":
			// Para confirm, converter string para bool
			formData[key] = field.Value == "true"
		case "note":
			// Notes não têm valores, apenas informações
			formData[key] = nil
		}
	}

	return formData
}

// displayValidationErrors exibe os erros de validação na tela
func (m *Model) displayValidationErrors(results map[string]config.ValidationResult) {
	fmt.Print("\033[2J\033[H") // Limpar tela

	validator := config.NewFieldValidator()
	allErrors := validator.GetAllErrors(results)

	fmt.Println(m.theme.FormTitle.Render("Erros de Validação"))
	fmt.Println()

	if len(allErrors) == 0 {
		fmt.Println("Nenhum erro encontrado.")
		return
	}

	for _, err := range allErrors {
		errorStyle := m.theme.FieldError
		fmt.Println(errorStyle.Render(err.Message))
	}

	fmt.Println()
	fmt.Println("Pressione Enter para continuar e corrigir os erros...")

	// Aguardar entrada do usuário
	fmt.Scanln()
}

// Start inicia a aplicação Bubble Tea.
func Start(cfg *config.FormConfig) error {
	p := tea.NewProgram(NewModel(cfg))
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("erro ao iniciar a TUI: %w", err)
	}
	return nil
}
