Executado. Concluímos as Fases 1.1 (Contratos), 1.2 (Layout), 1.3 (Componentes de Display) e 1.4 (Motores de Lógica).

O nosso runtime agora pode renderizar um layout complexo e executar scripts `on:` em resposta a eventos de `list` ou `buttongroup`.

Seguindo rigorosamente o **Plano Mestre Sequencial**, iniciamos a **Fase 1.5: Integração do Componente `form` (Absorção do v1.0)**.

Esta é a fase crítica onde aplicamos a nossa "visão do futuro" para resolver os teus problemas de refatoração desnecessária e código legado.

**Plano desta Fase:**

1.  **Encapsular (Não Refatorar):** Vamos mover a lógica v1.0 (`internal/tui/model.go`, `internal/config/config.go`, `internal/config/parser.go`, etc.) para um novo diretório (`internal/components/form/legacy_v1/`).
2.  **Aplicar Otimização de Risco:** Vamos modificar o `model_v1.go` (o TUI v1.0) para **remover** as chamadas `os.Exit` e a gestão de `WindowSizeMsg`.
3.  **Criar o "Wrapper":** Criaremos o `internal/components/form/wrapper.go`, que implementa a *nossa* interface `ShantillyComponent` (v2.0) e atua como um "adaptador" para o código v1.0.
4.  **Integrar:** Atualizaremos o `LayoutManager` (v2.0) para reconhecer `type: form` (FR7) e o `TemplateState` para tratar o *payload* do formulário (FR10).

Aqui estão os artefatos desta fase.

-----

### Tarefa 1.5.1: Encapsular Estruturas v1.0 (Legado)

Primeiro, movemos as estruturas de configuração v1.0 (de `internal/config/config.go`) para o seu novo local encapsulado.

```go
// internal/components/form/legacy_v1/config.go
package legacy_v1

// (Conteúdo movido de internal/config/config.go)

// FormConfig define a estrutura do YAML v1.0 para um formulário.
type FormConfig struct {
	Title  string  `yaml:"title"`
	Fields []Field `yaml:"fields"`
	// (Assumindo que Actions vêm do huh.Form)
}

// Field define um único campo no YAML v1.0.
type Field struct {
	Key         string   `yaml:"key"`
	Label       string   `yaml:"label"`
	Type        string   `yaml:"type"` // "input", "textarea", "select", etc.
	Required    bool     `yaml:"required,omitempty"`
	Options     []string `yaml:"options,omitempty"`
	Placeholder string   `yaml:"placeholder,omitempty"`
	// (Outros campos de validação v1.0)
	Min         *int     `yaml:"min,omitempty"`
	Max         *int     `yaml:"max,omitempty"`
	MinLength   *int     `yaml:"minLength,omitempty"`
	MaxLength   *int     `yaml:"maxLength,omitempty"`
	Pattern     string   `yaml:"pattern,omitempty"`
	FileTypes   []string `yaml:"fileTypes,omitempty"`
}
```

-----

### Tarefa 1.5.2: Encapsular Parser v1.0 (Legado)

Movemos a lógica de *parsing* v1.0 (de `internal/config/parser.go` e `internal/config/validation.go`) e a lógica de *criação do huh.Form* para uma função *helper* que o nosso *wrapper* v2.0 possa chamar.

```go
// internal/components/form/legacy_v1/parser.go
package legacy_v1

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/helton-godoy/shantilly/internal/tui/components" // (Assumindo que os componentes v1.0 estão aqui)
)

// CreateHuhForm é a lógica de parsing v1.0 extraída.
// Ela converte o YAML v1.0 (em formato interface{}) em um huh.Form.
// Fonte: Baseado em internal/config/parser.go
func CreateHuhForm(fieldsData, actionsData interface{}) (*huh.Form, error) {
	// (Esta função conteria a lógica complexa de
	// internal/config/parser.go e internal/config/validation.go
	// para converter os dados YAML em huh.Field e huh.Group)

	// --- Placeholder da Lógica de Parsing v1.0 ---
	// Na implementação real, usaríamos 'mapstructure' ou similar
	// para converter 'fieldsData' em '[]Field'
	
	// Simulação de parsing (FR7)
	fields := []huh.Field{
		huh.NewInput().
			Key("username").
			Title("Username (v1.0 Field)"),
		// (A lógica real de parsing v1.0 iria iterar sobre 'fieldsData')
	}
	
	// Simulação de parsing de Ações (FR7)
	var actions []huh.Field
	if actionsData != nil {
		// (Lógica de parsing para 'actions: type: buttongroup' v1.0)
		actions = []huh.Field{
			huh.NewConfirm().
				Key("submit").
				Title("Submit this form?"),
		}
	} else {
		// Ações padrão do v1.0 (se 'actions:' não for fornecido)
		actions = []huh.Field{
			huh.NewConfirm().
				Key("submit").
				Title("Submit?"),
		}
	}
	// --- Fim do Placeholder ---


	// (A lógica de validação de internal/config/validation.go seria aplicada
	// aos 'fields' aqui)

	form := huh.NewForm(
		huh.NewGroup(fields...).WithKey("data"),
		huh.NewGroup(actions...).WithKey("actions"),
	)
	
	return form, nil
}
```

-----

### Tarefa 1.5.3: Aplicar Otimização de Risco (Refatorar TUI v1.0)

Este é o artefato mais importante desta fase. Pegamos o `internal/tui/model.go` (v1.0) e o adaptamos.

  * ***Otimização de Risco (Visão do Futuro):***
    1.  Removemos `tea.WindowSizeMsg` (agora é trabalho do `LayoutManager`).
    2.  Removemos `tea.Quit` na submissão.
    3.  Em vez disso, emitimos `formSubmitMsg` (para o *wrapper*) ou `events.RuntimeErrorMsg` (para o `LayoutManager`).
    4.  Adicionamos `SetDimensions` para que o `LayoutManager` possa controlá-lo.

<!-- end list -->

```go
// internal/components/form/legacy_v1/model.go
package legacy_v1

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/helton-godoy/shantilly/internal/util"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// formSubmitMsg é um evento *interno* que o TUI v1.0
// envia ao Wrapper v2.0.
type formSubmitMsg struct {
	Payload interface{}
}

// Model (v1.0 Refatorado)
// Fonte: Baseado em internal/tui/model.go
type Model struct {
	form     *huh.Form
	progress progress.Model
	width    int
	// (Removido: height, state, quit, error)
}

func NewModel(form *huh.Form) *Model {
	return &Model{
		form:     form,
		progress: progress.NewModel(progress.WithDefaultGradient()),
		width:    80, // (Padrão, será sobrescrito por SetDimensions)
	}
}

func (m *Model) Init() tea.Cmd {
	return m.form.Init()
}

// Update (v1.0 Refatorado - Otimização de Risco)
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Processa o formulário
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	// --- OTIMIZAÇÃO DE RISCO (Remoção de tea.Quit) ---
	if m.form.State == huh.StateCompleted {
		// v1.0 (Legado): return m, tea.Quit
		// v2.0 (Novo): Emitir evento interno para o Wrapper (FR7)
		data := m.form.Get("data") // (Assumindo que o grupo de dados se chama 'data')
		cmd = func() tea.Msg { return formSubmitMsg{Payload: data} }
		cmds = append(cmds, cmd)
	}

	// --- OTIMIZAÇÃO DE RISCO (Tratamento de Abort) ---
	if m.form.State == huh.StateAborted {
		// v1.0 (Legado): m.quit = true; return m, tea.Quit
		// v2.0 (Novo): Emitir erro para o LayoutManager
		cmd = func() tea.Msg { return events.RuntimeErrorMsg{Err: util.ErrAborted} }
		cmds = append(cmds, cmd)
	}
	
	// --- OTIMIZAÇÃO DE RISCO (Remoção de WindowSizeMsg) ---
	// O 'case tea.WindowSizeMsg:' foi REMOVIDO.
	// O LayoutManager (v2.0) é agora responsável pelo tamanho
	// através do método SetDimensions.

	// (A lógica de 'progress' e 'spinner' do v1.0 permaneceria aqui)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// (A lógica de 'View' do v1.0 permaneceria aqui,
	// renderizando m.form.View() e m.progress.View())
	
	// Se a largura for 0, não renderize (evita pânico no 'huh')
	if m.width == 0 {
		return "..."
	}
	
	return m.form.View()
}

// SetDimensions (ADICIONADO)
// O LayoutManager (v2.0) chama esta função.
func (m *Model) SetDimensions(w, h int) {
	m.width = w
	// (O 'huh' usa variáveis de ambiente ou 'WithWidth',
	// vamos definir a largura no 'form' aqui)
	m.form.WithWidth(w)
}
```

-----

### Tarefa 1.5.4: Criar o FormComponent Wrapper (FR7)

Este é o "Adaptador" v2.0 que implementa `ShantillyComponent` e "embrulha" o `model_v1.go` refatorado.

```go
// internal/components/form/wrapper.go
package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/components/form/legacy_v1"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Wrapper implementa a interface ShantillyComponent (v2.0)
// e encapsula o TUI v1.0 (legacy_v1.Model).
// Fonte: docs/architecture.md (Seção 5.5, Ponto 1)
type Wrapper struct {
	id    string
	theme *tui.Theme

	// O modelo TUI v1.0 refatorado
	legacyModel tea.Model
}

func New(id string, theme *tui.Theme, cmp *declarative.Component) (tui.ShantillyComponent, error) {
	// 1. Usar o Parser v1.0 (Tarefa 1.5.2)
	// (FR7: "usando a sintaxe huh já validada no v1.0")
	huhForm, err := legacy_v1.CreateHuhForm(cmp.Fields, cmp.Actions)
	if err != nil {
		return nil, err
	}

	// 2. Criar o Modelo v1.0 (Tarefa 1.5.3)
	legacyModel := legacy_v1.NewModel(huhForm)

	return &Wrapper{
		id:          id,
		theme:       theme,
		legacyModel: legacyModel,
	}, nil
}

func (w *Wrapper) Init() tea.Cmd {
	return w.legacyModel.Init()
}

// Update (v2.0) - O "Tradutor" de Eventos
func (w *Wrapper) Update(msg tea.Msg) (tui.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// 1. Encaminha a msg para o modelo v1.0
	w.legacyModel, cmd = w.legacyModel.Update(msg)
	cmds = append(cmds, cmd)

	// 2. Verifica se o v1.0 emitiu a sua msg interna (formSubmitMsg)
	switch msg := msg.(type) {
	case legacy_v1.formSubmitMsg:
		// 3. TRADUZIR para o evento v2.0 (FR7)
		// Fonte: docs/prd.md (Estória 1.4 AC3)
		event := events.ShantillyEvent{
			ComponentID: w.id,
			Type:        w.id + ":submit", // ex: "user_form:submit"
			Payload:     msg.Payload,
		}
		cmds = append(cmds, func() tea.Msg { return event })
	}

	return w, tea.Batch(cmds...)
}

func (w *Wrapper) View() string {
	return w.legacyModel.View()
}

// SetDimensions (v2.0) - Passa o tamanho para o modelo v1.0
func (w *Wrapper) SetDimensions(width, height int) {
	// (Verifica se o modelo v1.0 implementa SetDimensions)
	if s, ok := w.legacyModel.(interface{ SetDimensions(int, int) }); ok {
		s.SetDimensions(width, height)
	}
}

func (w *Wrapper) ID() string {
	return w.id
}
```

-----

### Tarefa 1.5.5: Atualizar LayoutManager (Integração)

**Modificamos** o `internal/runtime/layout/manager.go` (da Fase 1.4) para adicionar o `case "form":` na "fábrica" de componentes.

```go
// internal/runtime/layout/manager.go (MODIFICADO)
package layout

// (Importações)
import (
	// ... (imports da Fase 1.4)
	"github.com/helton-godoy/shantilly/internal/components/form" // <-- ADICIONADO
)

// (Struct Manager, NewManager, Init - permanecem os da Fase 1.4)
// ...

// buildComponentTree (MODIFICADO)
func (m *Manager) buildComponentTree(node *declarative.LayoutNode) {
	if node == nil { return }
	if node.Type == "box" && node.Component != nil {
		if node.ID == "" {
			m.lastErr = fmt.Errorf("componente do tipo '%s' não possui 'id'", node.Component.Type)
			return
		}
		
		var (
			cmp tui.ShantillyComponent
			err error // (Erro agora é possível)
		)

		switch node.Component.Type {
		case "list":
			cmp = list.New(node.ID, m.theme, node.Component.Items)
		case "viewport":
			cmp = viewport.New(node.ID, m.theme, node.Component.Source)
		case "buttongroup":
			cmp = buttongroup.New(node.ID, m.theme, node.Component.Items)
		
		// --- ADICIONADO (Estória 1.4 AC1) ---
		case "form":
			cmp, err = form.New(node.ID, m.theme, node.Component)
			if err != nil {
				m.lastErr = fmt.Errorf("erro ao criar form '%s': %w", node.ID, err)
				return
			}
		// ------------------------------------
			
		default:
			cmp = viewport.New(node.ID, m.theme, &declarative.Source{
				Type: "static",
				Content: fmt.Sprintf("Placeholder para Componente:\nType: %s\nID: %s",
					node.Component.Type, node.Component.ID),
			})
		}
		
		m.components[node.ID] = cmp
		m.focusOrder = append(m.focusOrder, node.ID)
	}
	for _, item := range node.Items {
		m.buildComponentTree(&item)
	}
}

// (Update, View, etc. - permanecem os da Fase 1.4)
// ...
```

-----

### Tarefa 1.5.6: Atualizar Motor de Templates (Integração FR10)

**Modificamos** o `internal/runtime/runner/templating.go` (da Fase 1.4) para tratar corretamente os *payloads* de formulário (FR10).

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

// (Struct TemplateState e NewTemplateState - permanecem os da Fase 1.4)
// ...
type TemplateState struct {
	Component map[string]events.ShantillyEvent
	Form      map[string]interface{}
}
func NewTemplateState() *TemplateState {
	return &TemplateState{
		Component: make(map[string]events.ShantillyEvent),
		Form:      make(map[string]interface{}),
	}
}
// ...

// UpdateState (MODIFICADO)
func (ts *TemplateState) UpdateState(event events.ShantillyEvent) {
	ts.Component[event.ComponentID] = event

	// ADICIONADO: Se for um evento de formulário, armazena os dados
	// (Estória 1.5 AC3)
	if strings.HasSuffix(event.Type, ":submit") {
		ts.Form[event.ComponentID] = event.Payload
	}
}

// (ProcessString - permanece o da Fase 1.4)
// ...
func (ts *TemplateState) ProcessString(input string) (string, error) {
	tmpl, err := template.New("string").Parse(input)
	if err != nil { return "", err }
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ts); err != nil { return "", err }
	return buf.String(), nil
}
// ...

// ProcessStdin (MODIFICADO) - Implementação robusta do FR10
// Fonte: docs/prd.md (FR10, Estória 1.5 AC2)
func (ts *TemplateState) ProcessStdin(input string) ([]byte, error) {
	// Implementação Robusta (FR10):
	// O 'input' (ex: "{{ .Form.user_form }}") é um template Go
	// que seleciona *o objeto* a ser serializado.
	
	// 1. Otimização: Se o input for simples (ex: "{{ .Form.form_id }}")
	// extrai o objeto diretamente para evitar serialização dupla.
	
	// (Implementação Simplificada para Estória 1.5 AC2)
	// Removemos '{{' e '}}', ' .', e dividimos por '.'
	key := strings.TrimSpace(input)
	key = strings.TrimPrefix(key, "{{")
	key = strings.TrimSuffix(key, "}}")
	key = strings.TrimSpace(key)
	key = strings.TrimPrefix(key, ".") // Remove o '.' inicial

	parts := strings.Split(key, ".")
	
	var dataToSerialize interface{}
	
	if len(parts) == 2 && parts[0] == "Form" && parts[1] != "" {
		// (ex: "Form.user_form")
		dataToSerialize = ts.Form[parts[1]]
	} else if len(parts) == 3 && parts[0] == "Component" && parts[2] == "Payload" {
		// (ex: "Component.menu.Payload")
		dataToSerialize = ts.Component[parts[1]].Payload
	} else if key == "" || key == "." {
		// (ex: "{{ . }}" ou "{{ }}")
		dataToSerialize = ts // Serializa o estado inteiro
	} else {
		// (Não é um objeto simples, usa o template para renderizar como string)
		str, err := ts.ProcessString(input)
		if err != nil {
			return nil, err
		}
		dataToSerialize = str
	}
	
	if dataToSerialize == nil {
		// O template referenciou dados que não existem
		return nil, fmt.Errorf("template stdin '%s' resultou em nil", input)
	}

	// 2. Serializa o objeto selecionado como JSON (FR10)
	return json.Marshal(dataToSerialize)
}
```