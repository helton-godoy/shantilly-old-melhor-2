package config

// FormConfig representa o formato legacy v1.x consumido pelo comando `shantilly form`.
// Este código existe apenas para manter compatibilidade com o fluxo antigo.
type FormConfig struct {
	Title  string  `yaml:"title"`
	Fields []Field `yaml:"fields"`
}

// Field representa um campo individual do formulário legacy.
type Field struct {
	Key     string   `yaml:"key"`
	Label   string   `yaml:"label"`
	Type    string   `yaml:"type"`
	Options []string `yaml:"options,omitempty"`
}
