package config

// FormConfig representa a estrutura de nível superior de um arquivo de definição de formulário.
type FormConfig struct {
	Title  string  `yaml:"title,omitempty" json:"title,omitempty"`
	Fields []Field `yaml:"fields" json:"fields"`
}

// Field representa um único campo no formulário.
type Field struct {
	Key         string   `yaml:"key" json:"key"`
	Label       string   `yaml:"label" json:"label"`
	Type        string   `yaml:"type" json:"type"`
	Placeholder string   `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	Value       string   `yaml:"value,omitempty" json:"value,omitempty"`
	Options     []string `yaml:"options,omitempty" json:"options,omitempty"`
	Limit       int      `yaml:"limit,omitempty" json:"limit,omitempty"`
	Affirmative string   `yaml:"affirmative,omitempty" json:"affirmative,omitempty"`
	Negative    string   `yaml:"negative,omitempty" json:"negative,omitempty"`
	Title       string   `yaml:"title,omitempty" json:"title,omitempty"`
	Detail      string   `yaml:"detail,omitempty" json:"detail,omitempty"`
	// Novos campos de validação
	Required  bool     `yaml:"required,omitempty" json:"required,omitempty"`
	Min       *float64 `yaml:"min,omitempty" json:"min,omitempty"`
	Max       *float64 `yaml:"max,omitempty" json:"max,omitempty"`
	Pattern   string   `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	FileTypes []string `yaml:"fileTypes,omitempty" json:"fileTypes,omitempty"`
	MinLength int      `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength int      `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	// Novos campos para UX aprimorada
	Help string `yaml:"help,omitempty" json:"help,omitempty"`
}
