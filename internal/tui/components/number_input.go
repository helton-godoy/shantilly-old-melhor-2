package components

import (
	"strconv"
	"strings"

	"shantilly/internal/config"

	"github.com/charmbracelet/huh"
)

// NumberInput cria um campo de entrada numérica com validação
type NumberInput struct {
	field *config.Field
	value *string
}

// NewNumberInput cria um novo campo de entrada numérica
func NewNumberInput(field *config.Field) *NumberInput {
	var value string
	if field.Value != "" {
		value = field.Value
	}

	return &NumberInput{
		field: field,
		value: &value,
	}
}

// Field retorna o campo huh configurado
func (ni *NumberInput) Field() huh.Field {
	title := ni.field.Label
	if ni.field.Required {
		title += " *" // Indica campo obrigatório
	}

	return huh.NewInput().
		Key(ni.field.Key).
		Title(title).
		Placeholder(ni.field.Placeholder).
		Value(ni.value)
}

// GetValue retorna o valor atual como string
func (ni *NumberInput) GetValue() string {
	return *ni.value
}

// GetFloatValue retorna o valor como float64
func (ni *NumberInput) GetFloatValue() (float64, error) {
	if *ni.value == "" {
		return 0, nil
	}
	return strconv.ParseFloat(strings.TrimSpace(*ni.value), 64)
}

// Validate valida o valor do campo usando o sistema de validação
func (ni *NumberInput) Validate(validator *config.FieldValidator) *config.ValidationResult {
	result := validator.ValidateField(*ni.field, *ni.value)
	return &result
}
