package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ValidationError representa um erro de validação com detalhes específicos
type ValidationError struct {
	Field   string
	Code    string
	Message string
	Value   interface{}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("campo '%s': %s", e.Field, e.Message)
}

// ValidationResult contém o resultado da validação de um campo
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// Validator define a interface para validadores
type Validator interface {
	Validate(value interface{}, field Field) ValidationResult
}

// FieldValidator implementa a validação de campos
type FieldValidator struct{}

// NewFieldValidator cria um novo validador de campos
func NewFieldValidator() *FieldValidator {
	return &FieldValidator{}
}

// ValidateField valida um campo individual
func (v *FieldValidator) ValidateField(field Field, value interface{}) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []ValidationError{}}

	// Validação de campo obrigatório
	if field.Required {
		if v.isEmpty(value) {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   field.Key,
				Code:    "REQUIRED",
				Message: "Este campo é obrigatório",
				Value:   value,
			})
			return result // Retorna imediatamente para campos obrigatórios vazios
		}
	}

	// Validações específicas por tipo
	switch field.Type {
	case "number":
		result = v.validateNumber(field, value)
	case "date":
		result = v.validateDate(field, value)
	case "file":
		result = v.validateFile(field, value)
	case "input", "textarea":
		result = v.validateText(field, value)
	}

	return result
}

// isEmpty verifica se um valor está vazio
func (v *FieldValidator) isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch val := value.(type) {
	case string:
		return strings.TrimSpace(val) == ""
	case []string:
		return len(val) == 0
	default:
		return false
	}
}

// validateNumber valida campos numéricos
func (v *FieldValidator) validateNumber(field Field, value interface{}) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []ValidationError{}}

	strValue, ok := value.(string)
	if !ok {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "INVALID_TYPE",
			Message: "Valor deve ser um número",
			Value:   value,
		})
		return result
	}

	// Tentar converter para float64
	numValue, err := strconv.ParseFloat(strings.TrimSpace(strValue), 64)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "INVALID_NUMBER",
			Message: "Valor deve ser um número válido",
			Value:   value,
		})
		return result
	}

	// Validar limites mínimo e máximo
	if field.Min != nil && numValue < *field.Min {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "MIN_VALUE",
			Message: fmt.Sprintf("Valor deve ser maior ou igual a %.2f", *field.Min),
			Value:   value,
		})
	}

	if field.Max != nil && numValue > *field.Max {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "MAX_VALUE",
			Message: fmt.Sprintf("Valor deve ser menor ou igual a %.2f", *field.Max),
			Value:   value,
		})
	}

	return result
}

// validateDate valida campos de data
func (v *FieldValidator) validateDate(field Field, value interface{}) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []ValidationError{}}

	strValue, ok := value.(string)
	if !ok {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "INVALID_TYPE",
			Message: "Valor deve ser uma data",
			Value:   value,
		})
		return result
	}

	// Formatos de data aceitos
	formats := []string{
		"2006-01-02", // YYYY-MM-DD
		"02/01/2006", // DD/MM/YYYY
		"01/02/2006", // MM/DD/YYYY
	}

	var parseErr error
	for _, format := range formats {
		_, parseErr = time.Parse(format, strings.TrimSpace(strValue))
		if parseErr == nil {
			return result // Data válida
		}
	}

	// Se nenhum formato funcionou
	result.IsValid = false
	result.Errors = append(result.Errors, ValidationError{
		Field:   field.Key,
		Code:    "INVALID_DATE",
		Message: "Data deve estar no formato YYYY-MM-DD, DD/MM/YYYY ou MM/DD/YYYY",
		Value:   value,
	})

	return result
}

// validateFile valida campos de arquivo
func (v *FieldValidator) validateFile(field Field, value interface{}) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []ValidationError{}}

	strValue, ok := value.(string)
	if !ok {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "INVALID_TYPE",
			Message: "Valor deve ser um caminho de arquivo",
			Value:   value,
		})
		return result
	}

	// Verificar se o arquivo tem extensão permitida
	if len(field.FileTypes) > 0 {
		validExtension := false
		fileExt := strings.ToLower(strings.TrimSpace(strValue))

		// Extrair extensão do arquivo
		parts := strings.Split(fileExt, ".")
		if len(parts) > 1 {
			ext := "." + parts[len(parts)-1]
			for _, allowedType := range field.FileTypes {
				if strings.ToLower(allowedType) == ext {
					validExtension = true
					break
				}
			}
		}

		if !validExtension {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   field.Key,
				Code:    "INVALID_FILE_TYPE",
				Message: fmt.Sprintf("Tipo de arquivo não permitido. Tipos aceitos: %s", strings.Join(field.FileTypes, ", ")),
				Value:   value,
			})
		}
	}

	return result
}

// validateText valida campos de texto (input, textarea)
func (v *FieldValidator) validateText(field Field, value interface{}) ValidationResult {
	result := ValidationResult{IsValid: true, Errors: []ValidationError{}}

	strValue, ok := value.(string)
	if !ok {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "INVALID_TYPE",
			Message: "Valor deve ser texto",
			Value:   value,
		})
		return result
	}

	// Validar comprimento mínimo
	if field.MinLength > 0 && len(strings.TrimSpace(strValue)) < field.MinLength {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "MIN_LENGTH",
			Message: fmt.Sprintf("Texto deve ter pelo menos %d caracteres", field.MinLength),
			Value:   value,
		})
	}

	// Validar comprimento máximo
	if field.MaxLength > 0 && len(strings.TrimSpace(strValue)) > field.MaxLength {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   field.Key,
			Code:    "MAX_LENGTH",
			Message: fmt.Sprintf("Texto deve ter no máximo %d caracteres", field.MaxLength),
			Value:   value,
		})
	}

	// Validar padrão regex se especificado
	if field.Pattern != "" {
		// Nota: Para simplificar, vamos implementar validação básica de padrão
		// Uma implementação completa usaria regexp.MatchString
		if field.Pattern == "email" && !strings.Contains(strValue, "@") {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   field.Key,
				Code:    "INVALID_PATTERN",
				Message: "Formato de email inválido",
				Value:   value,
			})
		}
	}

	return result
}

// ValidateForm valida todos os campos de um formulário
func (v *FieldValidator) ValidateForm(form *FormConfig, values map[string]interface{}) map[string]ValidationResult {
	results := make(map[string]ValidationResult)

	for _, field := range form.Fields {
		value, exists := values[field.Key]
		if !exists {
			value = ""
		}
		results[field.Key] = v.ValidateField(field, value)
	}

	return results
}

// HasErrors verifica se há erros de validação nos resultados
func (v *FieldValidator) HasErrors(results map[string]ValidationResult) bool {
	for _, result := range results {
		if !result.IsValid {
			return true
		}
	}
	return false
}

// GetAllErrors retorna todos os erros de validação
func (v *FieldValidator) GetAllErrors(results map[string]ValidationResult) []ValidationError {
	var allErrors []ValidationError
	for _, result := range results {
		allErrors = append(allErrors, result.Errors...)
	}
	return allErrors
}
