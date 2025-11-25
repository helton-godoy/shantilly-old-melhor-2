Executado. Concluímos com sucesso o **Épico 2: O Runner Especialista**. O nosso runtime agora não só executa scripts genéricos (Fase 1.4) com *streaming* (Fase 2.3), como também é um *dispatcher* que pode invocar o *runner* `ansible-playbook` (Fase 2.4), suportado pela nossa infraestrutura JIT de recolha de segredos (Fase 2.2).

Seguindo rigorosamente o **Plano Mestre Sequencial**, iniciamos o **Épico 3: O Runtime Preditivo (Ansible Fase 3)**.

**Objetivo:** Implementar os componentes de descoberta `playbook_explorer` e `inventory_explorer`.

  * **Fase 3.1 (Contratos):** Conforme a nossa análise anterior, o `struct Component` (da Fase 1.1) já suporta um `Type: string`. Não são necessárias modificações nos modelos de dados.
  * **Fase 3.2 (Implementação):** Vamos criar os novos componentes.
  * **Fase 3.3 (Integração):** Vamos atualizar a "fábrica" no `LayoutManager` (da Fase 1.5) para instanciar os novos tipos.

-----

### Tarefa 3.2.1: Implementar PlaybookExplorerComponent (NFR4)

Este componente implementa `ShantillyComponent`.

  * ***Otimização de Risco (Visão do Futuro):*** Implementamos rigorosamente o NFR4 (filtrar "ruído"). Ele usa `filepath.WalkDir` para encontrar ficheiros `.yml` ou `.yaml` e ignora explicitamente os diretórios `roles/` e `tasks/`.

<!-- end list -->

```go
// internal/components/playbook_explorer/model.go
package playbook_explorer

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Model implementa a interface ShantillyComponent.
// Fonte: docs/prd.md (Épico 3)
type Model struct {
	id    string
	theme *tui.Theme
	list  list.Model

	// (Propriedades específicas do componente, ex: 'baseDir')
	// (Assumindo que procura a partir do diretório atual ".")
}

// item (Reutiliza o padrão da Fase 1.3.2)
type item struct {
	path, title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.path }
func (i item) FilterValue() string { return i.title }

func New(id string, theme *tui.Theme, cmp *declarative.Component) *Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Playbook Explorer"

	return &Model{
		id:    id,
		theme: theme,
		list:  l,
	}
}

// scanPlaybooksCmd (ADICIONADO)
// Um tea.Cmd que executa a varredura do NFR4.
// Fonte: docs/prd.md (NFR4)
func (m *Model) scanPlaybooksCmd() tea.Msg {
	var playbooks []list.Item
	baseDir := "." // (Poderia ser configurável via 'cmp *declarative.Component')

	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// NFR4: Filtrar "ruído"
		if d.IsDir() {
			dirName := d.Name()
			if dirName == "roles" || dirName == "tasks" || dirName == ".git" {
				return filepath.SkipDir // Pular este diretório
			}
		}

		// NFR4: Encontrar playbooks
		if !d.IsDir() && (strings.HasSuffix(d.Name(), ".yml") || strings.HasSuffix(d.Name(), ".yaml")) {
			playbooks = append(playbooks, item{
				path:  path,
				title: d.Name(),
			})
		}
		return nil
	})

	if err != nil {
		return events.RuntimeErrorMsg{Err: err}
	}

	return m.list.SetItems(playbooks) // Retorna uma msg que o 'list' entende
}

func (m *Model) Init() tea.Cmd {
	// Inicia a varredura (NFR4) como um comando
	return m.scanPlaybooksCmd
}

func (m *Model) Update(msg tea.Msg) (tui.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				// Emite evento "playbook_explorer_id:select"
				event := events.ShantillyEvent{
					ComponentID: m.id,
					Type:        m.id + ":select",
					Payload:     selectedItem.path, // Envia o caminho do playbook
				}
				cmds = append(cmds, func() tea.Msg { return event })
			}
		}
	}

	// Trata a navegação interna e as msgs (ex: SetItems)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.list.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model) ID() string {
	return m.id
}
```

-----

### Tarefa 3.2.2: Implementar InventoryExplorerComponent (NFR4)

Este componente é mais complexo.

  * ***Otimização de Risco (Visão do Futuro):*** Implementamos o NFR4 (usar `ansible-inventory` como "Oráculo") usando o mesmo padrão `tea.Cmd` do `Init()` que acabámos de usar, mas este executa um `exec.Cmd`.

<!-- end list -->

```go
// internal/components/inventory_explorer/model.go
package inventory_explorer

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Model implementa a interface ShantillyComponent.
// Fonte: docs/prd.md (Épico 3)
type Model struct {
	id    string
	theme *tui.Theme
	list  list.Model
	// (Propriedades específicas, ex: 'inventoryFile')
}

// item (Reutiliza o padrão)
type item struct {
	id, title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.title }

// Estrutura simples para o JSON de 'ansible-inventory --list'
type ansibleInventory struct {
	All struct {
		Hosts []string `json:"hosts"`
	} `json:"all"`
	// (Poderia analisar 'Groups' também)
}

func New(id string, theme *tui.Theme, cmp *declarative.Component) *Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Inventory Explorer"

	return &Model{
		id:    id,
		theme: theme,
		list:  l,
	}
}

// runAnsibleInventoryCmd (ADICIONADO)
// Um tea.Cmd que executa o "Oráculo" (NFR4)
// Fonte: docs/prd.md (NFR4)
func (m *Model) runAnsibleInventoryCmd() tea.Msg {
	// (O 'inventoryFile' poderia ser configurável via 'cmp *declarative.Component')
	inventoryFile := "inventory.yml" // (Placeholder)
	
	// NFR4: Usar 'ansible-inventory' como "Oráculo"
	cmd := exec.Command("ansible-inventory", "-i", inventoryFile, "--list")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return events.RuntimeErrorMsg{Err: fmt.Errorf("falha ao executar 'ansible-inventory': %w. Output: %s", err, string(output))}
	}

	var inventory ansibleInventory
	if err := json.Unmarshal(output, &inventory); err != nil {
		return events.RuntimeErrorMsg{Err: fmt.Errorf("falha ao analisar JSON do 'ansible-inventory': %w", err)}
	}

	var items []list.Item
	for _, host := range inventory.All.Hosts {
		items = append(items, item{
			id:    host,
			title: host,
		})
	}

	return m.list.SetItems(items) // Retorna a msg que o 'list' entende
}

func (m *Model) Init() tea.Cmd {
	// Inicia a execução do "Oráculo" (NFR4)
	return m.runAnsibleInventoryCmd
}

func (m *Model) Update(msg tea.Msg) (tui.ShantillyComponent, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			selectedItem, ok := m.list.SelectedItem().(item)
			if ok {
				// Emite evento "inventory_explorer_id:select"
				event := events.ShantillyEvent{
					ComponentID: m.id,
					Type:        m.id + ":select",
					Payload:     selectedItem.id, // Envia o host/grupo
				}
				cmds = append(cmds, func() tea.Msg { return event })
			}
		}
	}

	// Trata a navegação interna e as msgs (ex: SetItems)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.list.View()
}

func (m *Model) SetDimensions(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model) ID() string {
	return m.id
}
```

-----

### Tarefa 3.3.1: Modificar LayoutManager (Integração)

Finalmente, **modificamos** o `internal/runtime/layout/manager.go` (da Fase 1.5) para adicionar os `case` "playbook\_explorer" e "inventory\_explorer" na "fábrica" de componentes (`buildComponentTree`).

```go
// internal/runtime/layout/manager.go (MODIFICADO)
package layout

// (Importações)
import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/components/buttongroup"
	"github.com/helton-godoy/shantilly/internal/components/form"
	"github.com/helton-godoy/shantilly/internal/components/inventory_explorer" // <-- ADICIONADO
	"github.com/helton-godoy/shantilly/internal/components/list"
	"github.com/helton-godoy/shantilly/internal/components/playbook_explorer" // <-- ADICIONADO
	"github.com/helton-godoy/shantilly/internal/components/viewport"
	"github.com/helton-godoy/shantilly/internal/runtime/event"
	"github.com/helton-godoy/shantilly/internal/runtime/runner"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// (Struct Manager, NewManager, Init - permanecem os da Fase 2.3.3)
// ...
type Manager struct {
	theme         *tui.Theme
	root          *declarative.LayoutNode
	width         int
	height        int
	components    map[string]tui.ShantillyComponent
	focusOrder    []string
	focusedIdx    int
	lastErr       error
	eventManager  *event.Manager
	scriptRunner  *runner.Runner
	templateState *runner.TemplateState
}
func NewManager(root *declarative.LayoutNode, theme *tui.Theme, rules []declarative.Logic) *Manager {
	m := &Manager{
		theme:         theme,
		root:          root,
		components:    make(map[string]tui.ShantillyComponent),
		focusOrder:    []string{},
		focusedIdx:    0,
		eventManager:  event.New(rules),
		scriptRunner:  runner.New(),
		templateState: runner.NewTemplateState(),
	}
	m.buildComponentTree(root)
	return m
}
func (m *Manager) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range m.components {
		cmds = append(cmds, c.Init())
	}
	return tea.Batch(cmds...)
}
// ...


// buildComponentTree (MODIFICADO)
// Fonte: docs/prd.md (Épico 3)
func (m *Manager) buildComponentTree(node *declarative.LayoutNode) {
	if node == nil { return }
	if node.Type == "box" && node.Component != nil {
		if node.ID == "" {
			m.lastErr = fmt.Errorf("componente do tipo '%s' não possui 'id'", node.Component.Type)
			return
		}
		
		var (
			cmp tui.ShantillyComponent
			err error
		)

		switch node.Component.Type {
		// (Fase 1.3)
		case "list":
			cmp = list.New(node.ID, m.theme, node.Component.Items)
		case "viewport":
			cmp = viewport.New(node.ID, m.theme, node.Component.Source)
		case "buttongroup":
			cmp = buttongroup.New(node.ID, m.theme, node.Component.Items)
		
		// (Fase 1.5)
		case "form":
			cmp, err = form.New(node.ID, m.theme, node.Component)
			if err != nil {
				m.lastErr = fmt.Errorf("erro ao criar form '%s': %w", node.ID, err)
				return
			}
			
		// --- ADICIONADO (Épico 3) ---
		case "playbook_explorer":
			cmp = playbook_explorer.New(node.ID, m.theme, node.Component)
			
		case "inventory_explorer":
			cmp = inventory_explorer.New(node.ID, m.theme, node.Component)
		// -----------------------------
			
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

// (Update, View, renderNode, getComponent, isFocused - 
//  permanecem os da Fase 2.3.3)
// ...
```