Executado. Concluímos com sucesso o **Épico 1: Fundação do Runtime TUI (Genérico)**.

Temos uma base sólida: o `LayoutManager` (v2.0) está separado dos motores de lógica (`EventManager`, `ScriptRunner`), o formulário v1.0 está encapsulado e a segurança JIT (via `confirm:`) está funcional.

Seguindo o **Plano Mestre Sequencial**, iniciamos o **Épico 2: O Runner Especialista (Ansible Fase 2)**.

O objetivo é implementar o *runner* `run: { ansible_playbook: ... }` e o modal `ask_vault_pass: true`.

-----

### Fase 2.1: Extensão dos Contratos (Secrets & Ansible)

Primeiro, atualizamos os modelos de dados (Fase 1.1) para reconhecer a nova sintaxe do YAML e suportar o estado dos segredos.

**Tarefa 2.1.1: Modificar Modelos de Dados (pkg/declarative/models.go)**

Vamos **modificar** o `pkg/declarative/models.go` (da Fase 1.1) para adicionar a estrutura `AnsiblePlaybookAction` e ligá-la à `RunAction`.

```go
// pkg/declarative/models.go (MODIFICADO)
package declarative

// (Structs Config, LayoutNode, Component, Item, Source, Logic - permanecem os da Fase 1.1)
// ...
type Config struct {
	Root LayoutNode `yaml:",inline"`
	On   []Logic    `yaml:"on,omitempty"`
}
type LayoutNode struct {
	Type      string       `yaml:"type"`
	ID        string       `yaml:"id,omitempty"`
	Width     string       `yaml:"width,omitempty"`
	Height    int          `yaml:"height,omitempty"`
	Flex      int          `yaml:"flex,omitempty"`
	Items     []LayoutNode `yaml:"items,omitempty"`
	Component *Component `yaml:"component,omitempty"`
}
type Component struct {
	Type    string      `yaml:"type"`
	ID      string      `yaml:"id,omitempty"`
	Items   []Item      `yaml:"items,omitempty"`
	Source  *Source     `yaml:"source,omitempty"`
	Fields  interface{} `yaml:"fields,omitempty"`
	Actions interface{} `yaml:"actions,omitempty"`
	Content string      `yaml:"content,omitempty"`
}
type Item struct {
	ID    string `yaml:"id"`
	Text  string `yaml:"text,omitempty"`
	Label string `yaml:"label,omitempty"`
	Role  string `yaml:"role,omitempty"`
}
type Source struct {
	Type        string `yaml:"type"`
	Content     string `yaml:"content,omitempty"`
	Exec        string `yaml:"exec,omitempty"`
	ContentType string `yaml:"content_type,omitempty"`
}
type Logic struct {
	Event         string   `yaml:"event"`
	Run           RunAction  `yaml:"run"` // <-- Modificado abaixo
	Confirm       bool     `yaml:"confirm,omitempty"`
	PromptSecrets []string `yaml:"prompt_secrets,omitempty"` // (Suporte para Fase 2.x)
}

// RunAction (MODIFICADO)
// Agora suporta múltiplos tipos de runners (Script ou Ansible)
type RunAction struct {
	// Runner Genérico (FR9)
	Script string `yaml:"script,omitempty"`
	
	// Runner Especialista (Épico 2)
	AnsiblePlaybook *AnsiblePlaybookAction `yaml:"ansible_playbook,omitempty"` // <-- ADICIONADO

	// Argumentos (Comum a ambos os runners)
	Args         []string `yaml:"args,omitempty"`         // FR10
	Stdin        string   `yaml:"stdin,omitempty"`          // FR10
	UpdateTarget string   `yaml:"update_target,omitempty"`  // FR11
}

// AnsiblePlaybookAction (ADICIONADO)
// Define os parâmetros para o runner 'ansible_playbook'
// Fonte: docs/prd.md (Épico 2)
type AnsiblePlaybookAction struct {
	Playbook     string                 `yaml:"playbook"`      // Caminho para o playbook
	Inventory    string                 `yaml:"inventory,omitempty"` // Caminho para o inventário
	Vars         map[string]interface{} `yaml:"vars,omitempty"` // Variáveis (passadas como --extra-vars)
	AskVaultPass bool                   `yaml:"ask_vault_pass,omitempty"` // Ativa o modal JIT
}
```

-----

### Tarefa 2.1.2: Modificar Estado do Template

**Modificamos** o `internal/runtime/runner/templating.go` (da Fase 1.5) para adicionar o mapa `Secret` ao `TemplateState`. Isto é necessário para que o `AnsibleRunner` (Fase 2.4) possa aceder à *vault pass* recolhida pelo modal JIT (Fase 2.2).

```go
// internal/runtime/runner/templating.go (MODIFICADO)
package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// TemplateState (MODIFICADO)
// Fonte: docs/architecture.md (Seção 7.3 - "enriquecidos")
type TemplateState struct {
	Component map[string]events.ShantillyEvent
	Form      map[string]interface{}
	// ADICIONADO: Armazena segredos recolhidos JIT
	Secret map[string]string
}

// NewTemplateState (MODIFICADO)
func NewTemplateState() *TemplateState {
	return &TemplateState{
		Component: make(map[string]events.ShantillyEvent),
		Form:      make(map[string]interface{}),
		Secret:    make(map[string]string), // <-- ADICIONADO
	}
}

// (UpdateState, ProcessString, ProcessStdin - permanecem os da Fase 1.5)
// ...
func (ts *TemplateState) UpdateState(event events.ShantillyEvent) {
	ts.Component[event.ComponentID] = event
	if strings.HasSuffix(event.Type, ":submit") {
		ts.Form[event.ComponentID] = event.Payload
	}
}
func (ts *TemplateState) ProcessString(input string) (string, error) {
	tmpl, err := template.New("string").Parse(input)
	if err != nil { return "", err }
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ts); err != nil { return "", err }
	return buf.String(), nil
}
func (ts *TemplateState) ProcessStdin(input string) ([]byte, error) {
	key := strings.TrimSpace(input)
	key = strings.TrimPrefix(key, "{{")
	key = strings.TrimSuffix(key, "}}")
	key = strings.TrimSpace(key)
	key = strings.TrimPrefix(key, ".")
	parts := strings.Split(key, ".")
	var dataToSerialize interface{}
	if len(parts) == 2 && parts[0] == "Form" && parts[1] != "" {
		dataToSerialize = ts.Form[parts[1]]
	} else if len(parts) == 3 && parts[0] == "Component" && parts[2] == "Payload" {
		dataToSerialize = ts.Component[parts[1]].Payload
	} else if key == "" || key == "." {
		dataToSerialize = ts
	} else {
		str, err := ts.ProcessString(input)
		if err != nil { return nil, err }
		dataToSerialize = str
	}
	if dataToSerialize == nil {
		return nil, fmt.Errorf("template stdin '%s' resultou em nil", input)
	}
	return json.Marshal(dataToSerialize)
}
```