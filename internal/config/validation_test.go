package config

import (
	"testing"
)

func TestFieldValidator_ValidateField_Required(t *testing.T) {
	validator := NewFieldValidator()

	tests := []struct {
		name     string
		field    Field
		value    interface{}
		expected bool
	}{
		{
			name: "campo obrigatório vazio deve falhar",
			field: Field{
				Key:      "test",
				Required: true,
			},
			value:    "",
			expected: false,
		},
		{
			name: "campo obrigatório preenchido deve passar",
			field: Field{
				Key:      "test",
				Required: true,
			},
			value:    "valor",
			expected: true,
		},
		{
			name: "campo não obrigatório vazio deve passar",
			field: Field{
				Key:      "test",
				Required: false,
			},
			value:    "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateField(tt.field, tt.value)
			if result.IsValid != tt.expected {
				t.Errorf("esperado %v, mas obteve %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestFieldValidator_ValidateField_Number(t *testing.T) {
	validator := NewFieldValidator()

	min := 10.0
	max := 100.0

	tests := []struct {
		name     string
		field    Field
		value    interface{}
		expected bool
	}{
		{
			name: "número válido deve passar",
			field: Field{
				Key:  "test",
				Type: "number",
			},
			value:    "50",
			expected: true,
		},
		{
			name: "número abaixo do mínimo deve falhar",
			field: Field{
				Key:  "test",
				Type: "number",
				Min:  &min,
			},
			value:    "5",
			expected: false,
		},
		{
			name: "número acima do máximo deve falhar",
			field: Field{
				Key:  "test",
				Type: "number",
				Max:  &max,
			},
			value:    "150",
			expected: false,
		},
		{
			name: "valor não numérico deve falhar",
			field: Field{
				Key:  "test",
				Type: "number",
			},
			value:    "abc",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateField(tt.field, tt.value)
			if result.IsValid != tt.expected {
				t.Errorf("esperado %v, mas obteve %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestFieldValidator_ValidateField_Date(t *testing.T) {
	validator := NewFieldValidator()

	tests := []struct {
		name     string
		field    Field
		value    interface{}
		expected bool
	}{
		{
			name: "data YYYY-MM-DD válida deve passar",
			field: Field{
				Key:  "test",
				Type: "date",
			},
			value:    "2023-12-25",
			expected: true,
		},
		{
			name: "data DD/MM/YYYY válida deve passar",
			field: Field{
				Key:  "test",
				Type: "date",
			},
			value:    "25/12/2023",
			expected: true,
		},
		{
			name: "data MM/DD/YYYY válida deve passar",
			field: Field{
				Key:  "test",
				Type: "date",
			},
			value:    "12/25/2023",
			expected: true,
		},
		{
			name: "data inválida deve falhar",
			field: Field{
				Key:  "test",
				Type: "date",
			},
			value:    "25-12-2023",
			expected: false,
		},
		{
			name: "valor não string deve falhar",
			field: Field{
				Key:  "test",
				Type: "date",
			},
			value:    123,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateField(tt.field, tt.value)
			if result.IsValid != tt.expected {
				t.Errorf("esperado %v, mas obteve %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestFieldValidator_ValidateField_File(t *testing.T) {
	validator := NewFieldValidator()

	tests := []struct {
		name     string
		field    Field
		value    interface{}
		expected bool
	}{
		{
			name: "arquivo com extensão permitida deve passar",
			field: Field{
				Key:       "test",
				Type:      "file",
				FileTypes: []string{".pdf", ".docx"},
			},
			value:    "document.pdf",
			expected: true,
		},
		{
			name: "arquivo com extensão não permitida deve falhar",
			field: Field{
				Key:       "test",
				Type:      "file",
				FileTypes: []string{".pdf", ".docx"},
			},
			value:    "image.jpg",
			expected: false,
		},
		{
			name: "arquivo sem restrição de tipo deve passar",
			field: Field{
				Key:  "test",
				Type: "file",
			},
			value:    "anyfile.txt",
			expected: true,
		},
		{
			name: "valor não string deve falhar",
			field: Field{
				Key:       "test",
				Type:      "file",
				FileTypes: []string{".pdf"},
			},
			value:    123,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateField(tt.field, tt.value)
			if result.IsValid != tt.expected {
				t.Errorf("esperado %v, mas obteve %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestFieldValidator_ValidateField_Text(t *testing.T) {
	validator := NewFieldValidator()

	tests := []struct {
		name     string
		field    Field
		value    interface{}
		expected bool
	}{
		{
			name: "texto dentro dos limites deve passar",
			field: Field{
				Key:       "test",
				Type:      "input",
				MinLength: 2,
				MaxLength: 10,
			},
			value:    "hello",
			expected: true,
		},
		{
			name: "texto abaixo do mínimo deve falhar",
			field: Field{
				Key:       "test",
				Type:      "input",
				MinLength: 5,
			},
			value:    "hi",
			expected: false,
		},
		{
			name: "texto acima do máximo deve falhar",
			field: Field{
				Key:       "test",
				Type:      "input",
				MaxLength: 3,
			},
			value:    "hello world",
			expected: false,
		},
		{
			name: "email válido deve passar",
			field: Field{
				Key:     "test",
				Type:    "input",
				Pattern: "email",
			},
			value:    "user@example.com",
			expected: true,
		},
		{
			name: "email inválido deve falhar",
			field: Field{
				Key:     "test",
				Type:    "input",
				Pattern: "email",
			},
			value:    "invalid-email",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateField(tt.field, tt.value)
			if result.IsValid != tt.expected {
				t.Errorf("esperado %v, mas obteve %v", tt.expected, result.IsValid)
			}
		})
	}
}

func TestFieldValidator_ValidateForm(t *testing.T) {
	validator := NewFieldValidator()

	form := &FormConfig{
		Fields: []Field{
			{
				Key:      "name",
				Label:    "Nome",
				Type:     "input",
				Required: true,
			},
			{
				Key:   "age",
				Label: "Idade",
				Type:  "number",
				Min:   func() *float64 { v := 0.0; return &v }(),
				Max:   func() *float64 { v := 120.0; return &v }(),
			},
		},
	}

	values := map[string]interface{}{
		"name": "João",
		"age":  "25",
	}

	results := validator.ValidateForm(form, values)

	// Verificar que ambos os campos são válidos
	if !results["name"].IsValid {
		t.Errorf("campo 'name' deveria ser válido")
	}
	if !results["age"].IsValid {
		t.Errorf("campo 'age' deveria ser válido")
	}

	// Verificar que não há erros
	if validator.HasErrors(results) {
		t.Errorf("formulário não deveria ter erros")
	}
}

func TestFieldValidator_HasErrors(t *testing.T) {
	validator := NewFieldValidator()

	results := map[string]ValidationResult{
		"valid":   {IsValid: true, Errors: []ValidationError{}},
		"invalid": {IsValid: false, Errors: []ValidationError{{Code: "TEST"}}},
	}

	if !validator.HasErrors(results) {
		t.Errorf("deveria detectar erros")
	}

	// Testar com todos os campos válidos
	validResults := map[string]ValidationResult{
		"valid1": {IsValid: true, Errors: []ValidationError{}},
		"valid2": {IsValid: true, Errors: []ValidationError{}},
	}
	if validator.HasErrors(validResults) {
		t.Errorf("não deveria detectar erros quando todos são válidos")
	}
}
