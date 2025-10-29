# Data Models

## FormConfig (Input - YAML)

**Purpose:** Represents the YAML structure read from `stdin` or `--file`. Defines the form to be rendered. Used by `gopkg.in/yaml.v3` for unmarshalling.

**Key Attributes (Go Structs):**

```go
// Location: internal/config/config.go

package config

// FormConfig is the root structure parsed from YAML.
type FormConfig struct {
 Title  string  `yaml:"title,omitempty" json:"title,omitempty"` // Optional title (PRD).
 Fields []Field `yaml:"fields" json:"fields"`                   // Required list of fields (PRD).
}

// Field defines a single form field (component).
type Field struct {
 // --- Required Attributes (MVP) ---
 Key   string `yaml:"key" json:"key"`   // Unique ID for JSON output (FR6).
 Label string `yaml:"label" json:"label"` // Display text in TUI.
 Type  string `yaml:"type" json:"type"`  // Maps to 'huh' component (FR3).
                           // Expected: "input", "textarea", "select", "multiselect", "confirm", "note"

 // --- Optional Attributes (Based on PRD examples) ---
 Placeholder string   `yaml:"placeholder,omitempty" json:"placeholder,omitempty"` // Hint text for input/textarea.
 Value       string   `yaml:"value,omitempty" json:"value,omitempty"`       // Default value.
 Options     []string `yaml:"options,omitempty" json:"options,omitempty"`     // Choices for select/multiselect.
 Limit       int      `yaml:"limit,omitempty" json:"limit,omitempty"`       // Max selections for multiselect (0=unlimited).
 Affirmative string   `yaml:"affirmative,omitempty" json:"affirmative,omitempty"` // "Yes" text for confirm.
 Negative    string   `yaml:"negative,omitempty" json:"negative,omitempty"`    // "No" text for confirm.
 Title       string   `yaml:"title,omitempty" json:"title,omitempty"`        // Optional field title.
 Detail      string   `yaml:"detail,omitempty" json:"detail,omitempty"`      // Body text for note.
}
```

**Relationships:**

* A `FormConfig` HAS MANY `Field`.

## FormData (Output - JSON)

**Purpose:** Represents the data collected from the user upon successful form submission. Will be serialized to `stdout` (FR6).

**Key Attributes (Go Type):**

* As keys are user-defined in YAML (`key`), this is not a static Go struct.
* The Go type used for collection and serialization will be: `map[string]interface{}`
* The `string` key is the `key` from `Field`, and the `interface{}` value is the user input (`string`, `[]string`, or `bool`).

**Relationships:**

* This map is the final output of the `huh.Form` interaction, generated from `FormConfig`.
