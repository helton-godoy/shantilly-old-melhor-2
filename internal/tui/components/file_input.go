package components

import (
	"shantilly/internal/config"

	"github.com/charmbracelet/huh"
)

// FileInput cria um campo de seleção de arquivo com validação
type FileInput struct {
	field *config.Field
	value *string
}

// NewFileInput cria um novo campo de seleção de arquivo
func NewFileInput(field *config.Field) *FileInput {
	var value string
	if field.Value != "" {
		value = field.Value
	}

	return &FileInput{
		field: field,
		value: &value,
	}
}

// Field retorna o campo huh configurado
func (fi *FileInput) Field() huh.Field {
	title := fi.field.Label
	if fi.field.Required {
		title += " *" // Indica campo obrigatório
	}

	placeholder := fi.field.Placeholder
	if placeholder == "" {
		if len(fi.field.FileTypes) > 0 {
			placeholder = "Caminho do arquivo (tipos aceitos: " + fi.formatFileTypes() + ")"
		} else {
			placeholder = "Caminho do arquivo"
		}
	}

	return huh.NewInput().
		Key(fi.field.Key).
		Title(title).
		Placeholder(placeholder).
		Value(fi.value)
}

// GetValue retorna o valor atual como string
func (fi *FileInput) GetValue() string {
	return *fi.value
}

// formatFileTypes formata os tipos de arquivo para exibição
func (fi *FileInput) formatFileTypes() string {
	if len(fi.field.FileTypes) == 0 {
		return "qualquer"
	}

	result := ""
	for i, fileType := range fi.field.FileTypes {
		if i > 0 {
			result += ", "
		}
		result += fileType
	}
	return result
}

// Validate valida o valor do campo usando o sistema de validação
func (fi *FileInput) Validate(validator *config.FieldValidator) *config.ValidationResult {
	result := validator.ValidateField(*fi.field, *fi.value)
	return &result
}
