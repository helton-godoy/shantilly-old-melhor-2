package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Parse parses the YAML data into a FormConfig struct.
func Parse(data []byte) (*FormConfig, error) {
	var form FormConfig
	err := yaml.Unmarshal(data, &form)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear YAML: %w", err)
	}

	// Validate that we have at least one field
	if len(form.Fields) == 0 {
		return nil, fmt.Errorf("erro ao parsear YAML: formulário deve conter pelo menos um campo")
	}

	// Validate each field has required properties
	for i := range form.Fields {
		field := &form.Fields[i]
		if field.Key == "" {
			return nil, fmt.Errorf("erro ao parsear YAML: campo %d deve ter uma chave (key)", i+1)
		}
		if field.Label == "" {
			return nil, fmt.Errorf("erro ao parsear YAML: campo %d deve ter um rótulo (label)", i+1)
		}
		if field.Type == "" {
			return nil, fmt.Errorf("erro ao parsear YAML: campo %d deve ter um tipo (type)", i+1)
		}

		// Validate field type
		validTypes := []string{"input", "textarea", "select", "multiselect", "confirm", "note", "number", "date", "file"}
		isValidType := false
		for _, validType := range validTypes {
			if field.Type == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			return nil, fmt.Errorf("erro ao parsear YAML: tipo de campo '%s' inválido para o campo %d. Tipos válidos: input, textarea, select, multiselect, confirm, note, number, date, file", field.Type, i+1)
		}
	}

	return &form, nil
}
