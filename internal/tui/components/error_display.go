package components

import (
	"strings"
)

// ErrorDisplay é um componente para exibir mensagens de erro.
type ErrorDisplay struct {
	errors map[string]string
	theme  *Theme
}

// NewErrorDisplay cria um novo ErrorDisplay.
func NewErrorDisplay(theme *Theme) *ErrorDisplay {
	return &ErrorDisplay{
		errors: make(map[string]string),
		theme:  theme,
	}
}

// SetError define uma mensagem de erro para um campo específico.
func (e *ErrorDisplay) SetError(fieldKey, message string) {
	e.errors[fieldKey] = message
}

// ClearErrors limpa todas as mensagens de erro.
func (e *ErrorDisplay) ClearErrors() {
	e.errors = make(map[string]string)
}

// Render renderiza as mensagens de erro.
func (e *ErrorDisplay) Render() string {
	if len(e.errors) == 0 {
		return ""
	}

	var errorMessages []string
	// TODO: Sort keys for consistent order in tests
	for _, msg := range e.errors {
		errorMessages = append(errorMessages, e.theme.FieldError.Render("• "+msg))
	}

	return strings.Join(errorMessages, "\n")
}

// HasErrors retorna verdadeiro se houver algum erro.
func (e *ErrorDisplay) HasErrors() bool {
	return len(e.errors) > 0
}

// GetErrors retorna o mapa de erros.
func (e *ErrorDisplay) GetErrors() map[string]string {
	return e.errors
}

// ClearError remove um erro de um campo específico.
func (e *ErrorDisplay) ClearError(fieldKey string) {
	delete(e.errors, fieldKey)
}
