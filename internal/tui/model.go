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
	// Novos campos para UX aprimorada
	errors    map[string]string // field key -> error message
	helpTexts map[string]string // field key -> help text
	completed map[string]bool   // field key -> completion status
	progress  float64           // completion percentage
	showHelp  bool              // flag to show/hide help text
	// Componentes de UX aprimorada
	errorDisplay      *components.ErrorDisplay
	helpText          *components.HelpText
	progressIndicator *components.ProgressIndicator
}

const (
	defaultTerminalWidth  = 80 // Largura padrão do terminal
	defaultTerminalHeight = 24 // Altura padrão do terminal
)

// NewModel cria um novo modelo com a configuração do formulário.
func NewModel(formConfig *config.FormConfig) tea.Model {
	multiselectValues := &multiselectData{values: make(map[string]*[]string)}
	form := createForm(formConfig, multiselectValues)

	// Inicializar help texts dos campos
	helpTexts := make(map[string]string)
	for _, field := range formConfig.Fields {
		if field.Help != "" {
			helpTexts[field.Key] = field.Help
		}
	}

	// Criar componentes de UX aprimorada
	theme := DefaultTheme()
	errorDisplay := components.NewErrorDisplay(&components.Theme{
		FieldError: theme.FieldError,
	})
	helpText := components.NewHelpText(&components.Theme{
		FieldInput: theme.FieldInput,
	})
	// Inicializar help texts no componente
	for key, text := range helpTexts {
		helpText.SetHelpText(key, text)
	}
	progressIndicator := components.NewProgressIndicator(&components.Theme{
		FieldLabel: theme.FieldLabel,
	})

	return &Model{
		form:              form,
		formConfig:        formConfig,
		multiselectValues: multiselectValues,
		theme:             theme,
		terminalWidth:     defaultTerminalWidth,  // Largura padrão inicial
		terminalHeight:    defaultTerminalHeight, // Altura padrão inicial
		// Inicializar novos campos para UX aprimorada
		errors:    make(map[string]string),
		helpTexts: helpTexts,
		completed: make(map[string]bool),
		progress:  0.0,
		showHelp:  false,
		// Inicializar componentes
		errorDisplay:      errorDisplay,
		helpText:          helpText,
		progressIndicator: progressIndicator,
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
		case "?":
			// Toggle help display
			m.showHelp = !m.showHelp
			m.helpText.ToggleHelp()
		}
	}

	// Return any commands from form updates
	return m, cmd
}

// View renderiza a interface do usuário.
func (m *Model) View() string {
	// Atualizar o progresso antes de renderizar
	m.updateProgress()

	formView := m.form.View()

	// Adicionar indicadores de progresso
	progressView := m.renderProgressIndicator()
	if progressView != "" {
		formView = progressView + "\n\n" + formView
	}

	// Adicionar erros se houver
	errorView := m.errorDisplay.Render()
	if errorView != "" {
		formView += "\n\n" + errorView
	}

	// Adicionar texto de ajuda se ativado
	if m.showHelp {
		helpView := m.renderHelpText()
		if helpView != "" {
			formView += "\n\n" + helpView
		}
	}

	// Adicionar atalhos de teclado
	shortcutsView := m.renderKeyboardShortcuts()
	if shortcutsView != "" {
		formView += "\n\n" + shortcutsView
	}

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

// createForm cria um huh.Form a partir da configuração e retorna os valores multiselect.
func createForm(formConfig *config.FormConfig, multiselectValues *multiselectData) *huh.Form {
	var fields []huh.Field

	// Create fields based on configuration
	for i := range formConfig.Fields {
		field := &formConfig.Fields[i] // Use pointer to avoid copying

		// Inicializar help text se disponível
		if field.Help != "" {
			// helpTexts será inicializado no NewModel
		}

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
	// Limpar erros anteriores
	m.errorDisplay.ClearErrors()

	validator := config.NewFieldValidator()
	allErrors := validator.GetAllErrors(results)

	// Adicionar erros ao componente de exibição
	for _, err := range allErrors {
		m.errorDisplay.SetError(err.Field, err.Message)
		// Populate the model's own error map for other logic to use
		m.errors[err.Field] = err.Message
	}

	// Se houver erros, não submeter - manter formulário ativo
	if len(allErrors) > 0 {
		// Os erros serão exibidos na próxima renderização da View()
		return
	}
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

// setError define uma mensagem de erro para um campo específico
func (m *Model) setError(fieldKey, message string) {
	m.errors[fieldKey] = message
}

// clearError remove o erro de um campo específico
func (m *Model) clearError(fieldKey string) {
	delete(m.errors, fieldKey)
}

// getFieldHelp retorna o texto de ajuda para um campo específico
func (m *Model) getFieldHelp(fieldKey string) string {
	if help, exists := m.helpTexts[fieldKey]; exists {
		return help
	}
	return ""
}

// updateProgress calcula e atualiza a porcentagem de progresso do formulário
func (m *Model) updateProgress() {
	totalFields := 0
	completedFields := 0

	formData := m.collectFormData()

	for _, field := range m.formConfig.Fields {
		if field.Type == "note" {
			continue // Notes não contam para progresso
		}
		totalFields++
		// Check if the field has a non-empty value
		if val, ok := formData[field.Key]; ok && val != nil {
			isCompleted := false
			switch v := val.(type) {
			case string:
				if v != "" {
					isCompleted = true
				}
			case []string:
				if len(v) > 0 {
					isCompleted = true
				}
			case bool:
				if v { // Consider 'true' for confirm as completed
					isCompleted = true
				}
			default:
				// For other types like number, date, file, a non-nil value is enough
				isCompleted = true
			}
			m.completed[field.Key] = isCompleted
		}

		if m.completed[field.Key] {
			completedFields++
		}
	}

	if totalFields > 0 {
		m.progress = (float64(completedFields) / float64(totalFields)) * 100
	} else {
		m.progress = 0
	}
	m.progressIndicator.UpdateProgress(completedFields, totalFields)
}

// renderProgressIndicator renderiza o indicador de progresso
func (m *Model) renderProgressIndicator() string {
	completed := 0
	total := 0
	for _, field := range m.formConfig.Fields {
		if field.Type != "note" {
			total++
			if m.completed[field.Key] {
				completed++
			}
		}
	}
	m.progressIndicator.UpdateProgress(completed, total)
	return m.progressIndicator.Render()
}

// renderHelpText renderiza o texto de ajuda atual
func (m *Model) renderHelpText() string {
	// Solução simplificada: mostra a ajuda para o primeiro campo com erro.
	if len(m.errors) > 0 {
		for key := range m.errors {
			// Apenas pegue o primeiro erro que encontrarmos. A ordem não é garantida.
			m.helpText.SetCurrentField(key)
			return m.helpText.Render()
		}
	}

	// Se não houver erros, não mostre nenhuma ajuda contextual.
	m.helpText.SetCurrentField("")
	return m.helpText.Render()
}

// renderKeyboardShortcuts renderiza os atalhos de teclado disponíveis
func (m *Model) renderKeyboardShortcuts() string {
	shortcuts := []string{
		"? - Mostrar/ocultar ajuda",
		"Ctrl+C - Sair",
	}
	return m.theme.Base.Render("Atalhos: " + strings.Join(shortcuts, " | "))
}

// navigateToNextError (removido) - A API do huh não suporta foco programático
// de uma forma que torne esta funcionalidade trivial. A sinalização visual
// do erro é a principal forma de feedback por agora.
