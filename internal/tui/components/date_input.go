package components

import (
	"fmt"
	"time"

	"shantilly/internal/config"

	"github.com/charmbracelet/huh"
)

// DateInput cria um campo de entrada de data com validação
type DateInput struct {
	field *config.Field
	value *string
}

// NewDateInput cria um novo campo de entrada de data
func NewDateInput(field *config.Field) *DateInput {
	var value string
	if field.Value != "" {
		value = field.Value
	}

	return &DateInput{
		field: field,
		value: &value,
	}
}

// Field retorna o campo huh configurado
func (di *DateInput) Field() huh.Field {
	title := di.field.Label
	if di.field.Required {
		title += " *" // Indica campo obrigatório
	}

	placeholder := di.field.Placeholder
	if placeholder == "" {
		placeholder = "YYYY-MM-DD, DD/MM/YYYY ou MM/DD/YYYY"
	}

	return huh.NewInput().
		Key(di.field.Key).
		Title(title).
		Placeholder(placeholder).
		Value(di.value)
}

// GetValue retorna o valor atual como string
func (di *DateInput) GetValue() string {
	return *di.value
}

// GetTimeValue retorna o valor como time.Time
func (di *DateInput) GetTimeValue() (time.Time, error) {
	if *di.value == "" {
		return time.Time{}, fmt.Errorf("valor vazio não é uma data válida")
	}

	// Formatos de data aceitos
	formats := []string{
		"2006-01-02", // YYYY-MM-DD
		"02/01/2006", // DD/MM/YYYY
		"01/02/2006", // MM/DD/YYYY
	}

	for _, format := range formats {
		if t, err := time.Parse(format, *di.value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("formato de data inválido: %s", *di.value)
}

// Validate valida o valor do campo usando o sistema de validação
func (di *DateInput) Validate(validator *config.FieldValidator) *config.ValidationResult {
	result := validator.ValidateField(*di.field, *di.value)
	return &result
}
