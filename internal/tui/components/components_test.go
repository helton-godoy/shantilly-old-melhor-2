package components

import (
	"testing"

	"shantilly/internal/config"
)

func TestNumberInput_NewNumberInput(t *testing.T) {
	field := &config.Field{
		Key:      "test_number",
		Label:    "Número de Teste",
		Type:     "number",
		Required: true,
		Min:      func() *float64 { v := 0.0; return &v }(),
		Max:      func() *float64 { v := 100.0; return &v }(),
	}

	ni := NewNumberInput(field)

	if ni.field.Key != "test_number" {
		t.Errorf("esperado chave 'test_number', mas obteve '%s'", ni.field.Key)
	}

	if ni.GetValue() != "" {
		t.Errorf("esperado valor vazio inicial, mas obteve '%s'", ni.GetValue())
	}
}

func TestNumberInput_GetFloatValue(t *testing.T) {
	field := &config.Field{
		Key:   "test",
		Label: "Teste",
		Type:  "number",
	}

	ni := NewNumberInput(field)

	// Testar valor vazio
	val, err := ni.GetFloatValue()
	if err != nil {
		t.Errorf("não esperava erro para valor vazio, mas obteve: %v", err)
	}
	if val != 0 {
		t.Errorf("esperado 0 para valor vazio, mas obteve %f", val)
	}

	// Simular entrada de valor
	*ni.value = "42.5"
	val, err = ni.GetFloatValue()
	if err != nil {
		t.Errorf("não esperava erro para '42.5', mas obteve: %v", err)
	}
	if val != 42.5 {
		t.Errorf("esperado 42.5, mas obteve %f", val)
	}
}

func TestDateInput_NewDateInput(t *testing.T) {
	field := &config.Field{
		Key:      "test_date",
		Label:    "Data de Teste",
		Type:     "date",
		Required: true,
	}

	di := NewDateInput(field)

	if di.field.Key != "test_date" {
		t.Errorf("esperado chave 'test_date', mas obteve '%s'", di.field.Key)
	}

	if di.GetValue() != "" {
		t.Errorf("esperado valor vazio inicial, mas obteve '%s'", di.GetValue())
	}
}

func TestDateInput_GetTimeValue(t *testing.T) {
	field := &config.Field{
		Key:   "test",
		Label: "Teste",
		Type:  "date",
	}

	di := NewDateInput(field)

	// Testar valor vazio
	_, err := di.GetTimeValue()
	if err == nil {
		t.Error("esperava erro para valor vazio")
	}

	// Testar data válida
	*di.value = "2023-12-25"
	timeVal, err := di.GetTimeValue()
	if err != nil {
		t.Errorf("não esperava erro para '2023-12-25', mas obteve: %v", err)
	}
	if timeVal.Year() != 2023 || timeVal.Month() != 12 || timeVal.Day() != 25 {
		t.Errorf("data incorreta: esperado 2023-12-25, mas obteve %v", timeVal)
	}
}

func TestFileInput_NewFileInput(t *testing.T) {
	field := &config.Field{
		Key:       "test_file",
		Label:     "Arquivo de Teste",
		Type:      "file",
		Required:  true,
		FileTypes: []string{".pdf", ".docx"},
	}

	fi := NewFileInput(field)

	if fi.field.Key != "test_file" {
		t.Errorf("esperado chave 'test_file', mas obteve '%s'", fi.field.Key)
	}

	if fi.GetValue() != "" {
		t.Errorf("esperado valor vazio inicial, mas obteve '%s'", fi.GetValue())
	}
}

func TestFileInput_FormatFileTypes(t *testing.T) {
	field := &config.Field{
		Key:  "test",
		Type: "file",
	}

	fi := NewFileInput(field)

	// Testar sem tipos de arquivo
	if fi.formatFileTypes() != "qualquer" {
		t.Errorf("esperado 'qualquer' para lista vazia, mas obteve '%s'", fi.formatFileTypes())
	}

	// Testar com tipos de arquivo
	fi.field.FileTypes = []string{".pdf", ".docx"}
	expected := ".pdf, .docx"
	if fi.formatFileTypes() != expected {
		t.Errorf("esperado '%s', mas obteve '%s'", expected, fi.formatFileTypes())
	}
}

func TestComponents_Validate(t *testing.T) {
	validator := config.NewFieldValidator()

	// Testar NumberInput com validação
	numberField := &config.Field{
		Key:      "number",
		Label:    "Número",
		Type:     "number",
		Required: true,
		Min:      func() *float64 { v := 0.0; return &v }(),
		Max:      func() *float64 { v := 100.0; return &v }(),
	}

	ni := NewNumberInput(numberField)
	*ni.value = "50"
	result := ni.Validate(validator)

	if !result.IsValid {
		t.Errorf("esperado número válido, mas obteve erros: %v", result.Errors)
	}

	// Testar com valor inválido
	*ni.value = "150" // Acima do máximo
	result = ni.Validate(validator)
	if result.IsValid {
		t.Error("esperado erro para número acima do máximo")
	}

	// Testar DateInput
	dateField := &config.Field{
		Key:      "date",
		Label:    "Data",
		Type:     "date",
		Required: true,
	}

	di := NewDateInput(dateField)
	*di.value = "2023-12-25"
	result = di.Validate(validator)

	if !result.IsValid {
		t.Errorf("esperado data válida, mas obteve erros: %v", result.Errors)
	}

	// Testar FileInput
	fileField := &config.Field{
		Key:       "file",
		Label:     "Arquivo",
		Type:      "file",
		Required:  true,
		FileTypes: []string{".pdf", ".docx"},
	}

	fi := NewFileInput(fileField)
	*fi.value = "document.pdf"
	result = fi.Validate(validator)

	if !result.IsValid {
		t.Errorf("esperado arquivo válido, mas obteve erros: %v", result.Errors)
	}

	// Testar com tipo inválido
	*fi.value = "image.jpg"
	result = fi.Validate(validator)
	if result.IsValid {
		t.Error("esperado erro para tipo de arquivo inválido")
	}
}
