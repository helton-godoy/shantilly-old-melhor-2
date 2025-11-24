Executado. Concluímos as Fases 1.1, 1.2 e 1.3.

Temos agora uma aplicação que renderiza um layout (FR1, FR2) e exibe componentes reais (FR4, FR5, FR6). No entanto, a interatividade é zero. Os componentes emitem eventos (`ShantillyEvent`), mas o `LayoutManager` (como vimos na Fase 1.3) apenas os exibe para depuração.

Iniciamos agora a **Fase 1.4: Motores de Lógica (Eventos e Ações)**.

Esta é a fase que dá "vida" ao TUI. Vamos implementar os motores `EventManager` e `ScriptRunner`.

  * ***Otimização de Risco (Visão do Futuro):*** Esta é a implementação da separação de responsabilidades (Layout vs. Lógica). O `LayoutManager` (que gere o foco) irá *delegar* eventos ao `EventManager`, que por sua vez irá *despachar* ações para o `ScriptRunner`.

Aqui estão os artefatos da Fase 1.4.

-----

### Tarefa 1.4.3: Implementar o Motor de Templates (FR10)

Primeiro, criamos o *helper* de *templating* (Estória 1.5 AC3). Vamos usar o `text/template` do Go para substituir valores como `{{ .Component.list_id.Payload }}` ou `{{ .Form.form_id.field_name }}`.

```go
// internal/runtime/runner/templating.go
package runner

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// TemplateState é o objeto de dados que expomos ao motor de templates.
// Fonte: docs/prd.md (Estória 1.5 AC3)
type TemplateState struct {
	// Armazena o último evento de componente, indexado por ID
	Component map[string]events.ShantillyEvent
	// Armazena o último payload de formulário, indexado por ID
	Form map[string]interface{}
	// (Na Fase 1.6, adicionaremos: Secret map[string]string)
}

// NewTemplateState cria um estado de template inicial.
func NewTemplateState() *TemplateState {
	return &TemplateState{
		Component: make(map[string]events.ShantillyEvent),
		Form:      make(map[string]interface{}),
	}
}

// UpdateState atualiza o estado com o último evento.
// Isto é crucial para que os templates tenham dados recentes.
func (ts *TemplateState) UpdateState(event events.ShantillyEvent) {
	ts.Component[event.ComponentID] = event

	// Se for um evento de formulário (Fase 1.5), armazena os dados
	if event.Type == event.ComponentID+":submit" {
		ts.Form[event.ComponentID] = event.Payload
	}
}

// ProcessString aplica o *templating* a uma string (para 'args').
func (ts *TemplateState) ProcessString(input string) (string, error) {
	tmpl, err := template.New("string").Parse(input)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ts); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ProcessStdin aplica o *templating* (FR10) e serializa para JSON (para 'stdin').
func (ts *TemplateState) ProcessStdin(input string) ([]byte, error) {
	// 1. Processa o template
	tmpl, err := template.New("stdin").Parse(input)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ts); err != nil {
		return nil, err
	}
	
	// 2. O resultado do template (ex: "{{ .Form.user_form }}")
	// precisa ser serializado como JSON.
	// (Esta é uma implementação simplificada; uma real precisaria
	// extrair o objeto do estado, não a string renderizada)
	
	// Implementação Simplificada (Estória 1.5 AC2):
	// Se o template for exatamente "{{ .Form.user_form }}",
	// serializa o 'user_form' diretamente.
	// (Vamos assumir que o 'input' é o nome do objeto por agora)
	
	// NOTA: A implementação robusta disto (FR10) é complexa.
	// Por agora, vamos serializar o *estado completo* se o input for '{{ . }}'
	// e o payload do formulário se for '{{ .Form.form_id }}'
	
	// Placeholder para FR10 / Estória 1.5 AC2:
	// A Tarefa 1.5 (Form) irá preencher 'ts.Form'.
	// Por agora, apenas retornamos a string processada.
	// Vamos refatorar isto na Fase 1.5 quando tivermos dados de formulário.
	
	// Serializa o *resultado* do template como JSON.
	return json.Marshal(buf.String())
}
```

-----

### Tarefa 1.4.1: Implementar o EventManager (FR8)

Este motor recebe eventos do `LayoutManager` e os mapeia para as regras `on:`, disparando Ações de Script (FR8).

```go
// internal/runtime/event/manager.go
package event

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Manager é o motor de lógica que trata do bloco 'on:'.
// Fonte: docs/architecture.md (Seção 5.4)
type Manager struct {
	rules []declarative.Logic
}

func New(rules []declarative.Logic) *Manager {
	return &Manager{
		rules: rules,
	}
}

// ProcessEvent é chamado pelo LayoutManager.
// Ele verifica se algum evento da UI corresponde a uma regra 'on:'.
// Fonte: docs/prd.md (FR8)
func (m *Manager) ProcessEvent(event events.ShantillyEvent) tea.Cmd {
	var cmds []tea.Cmd

	for _, rule := range m.rules {
		// Verifica se o 'event.Type' (ex: "menu:select")
		// corresponde à regra 'on.event'.
		if rule.Event == event.Type {
			
			// (Na Fase 1.6, verificaremos 'rule.Confirm' e 'rule.PromptSecrets' aqui)

			// FR9: Dispara a ação do script.
			// Criamos uma nova mensagem para o ScriptRunner (Tarefa 1.4.2)
			req := events.RunScriptRequestMsg{
				Rule: rule,
				// O estado do template será preenchido pelo LayoutManager,
				// que é o dono do estado.
			}
			cmds = append(cmds, func() tea.Msg { return req })
		}
	}

	return tea.Batch(cmds...)
}
```

-----

### Tarefa 1.4.2: Definir a Mensagem de Ação (Modificação)

Precisamos de uma nova mensagem em `pkg/tui/events.go` (da Fase 1.1) para que o `EventManager` possa *solicitar* que o `ScriptRunner` execute uma ação.

```go
// pkg/tui/events.go (MODIFICADO)
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/pkg/declarative" // <-- ADICIONADO
)

// ShantillyEvent é a mensagem padronizada que os componentes emitem
// para serem capturados pelo EventManager.
// Fonte: docs/architecture.md (Seção 5.3)
type ShantillyEvent struct {
	ComponentID string      // ID do componente que originou o evento
	Type        string      // Tipo do evento (ex: "list_select", "button_press", "form_submit")
	Payload     interface{} // Dados associados (ex: ID do item, dados do formulário)
}

// RuntimeErrorMsg é usado para notificar o TUI sobre um erro interno
// que precisa ser exibido ao usuário.
type RuntimeErrorMsg struct {
	Err error
}

// ScriptStdoutMsg é usado pelo ScriptRunner para enviar streaming de
// stdout/stderr para o ViewportComponent alvo.
type ScriptStdoutMsg struct {
	TargetID string // ID do viewport (update_target)
	Chunk    []byte // O pedaço de dados do stdout/stderr
}

// ShowModalMsg é emitida pelo EventManager (para Segurança JIT)
// e tratada pelo LayoutManager para exibir um modal.
type ShowModalMsg struct {
	Title    string
	Content  tea.Model // O modelo do modal a ser exibido (ex: confirmação, prompt de segredo)
	OnSubmit func(tea.Msg) tea.Cmd // O que fazer quando o modal for submetido
}

// RunScriptRequestMsg (ADICIONADO)
// Esta msg é enviada pelo EventManager (Tarefa 1.4.1)
// e tratada pelo ScriptRunner (Tarefa 1.4.2)
type RunScriptRequestMsg struct {
	Rule declarative.Logic // A regra 'on:' inteira
}
```

-----

### Tarefa 1.4.2: Implementar o ScriptRunner (FR9, FR11)

Este é o motor de execução. Ele ouve `RunScriptRequestMsg`, executa o script e, o mais importante, gere o ciclo de vida do processo (FR11).

  * ***Otimização de Risco (Visão do Futuro):*** Implementamos o `SIGTERM` (FR11) desde o início, conforme o fluxo de trabalho 7.2, para evitar processos zombies.

<!-- end list -->

```go
// internal/runtime/runner/runner.go
package runner

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Runner é o motor de execução de scripts (FR9).
// Fonte: docs/architecture.md (Seção 5.4, 7.2)
type Runner struct {
	// Gerencia processos ativos para implementar FR11 (Ciclo de Vida)
	activeProcesses map[string]*os.Process
}

func New() *Runner {
	return &Runner{
		activeProcesses: make(map[string]*os.Process),
	}
}

// HandleRunRequest é chamado pelo LayoutManager quando um EventManager
// solicita a execução de um script.
func (r *Runner) HandleRunRequest(req events.RunScriptRequestMsg, state *TemplateState) tea.Cmd {
	rule := req.Rule
	targetID := rule.Run.UpdateTarget

	// 1. Implementar FR11: Ciclo de Vida do Target
	// Fonte: docs/prd.md (FR11), docs/architecture.md (Seção 7.2)
	if proc, exists := r.activeProcesses[targetID]; exists && proc != nil {
		// Tenta terminar (SIGTERM) o processo anterior.
		// (A implementação robusta lidaria com SIGKILL após timeout)
		_ = proc.Signal(syscall.SIGTERM)
		delete(r.activeProcesses, targetID)
	}

	// 2. Preparar o comando (FR9, FR10)
	// (Usando o TemplateState da Fase 1.4.1)
	
	// Processar 'args'
	processedArgs := make([]string, len(rule.Run.Args))
	for i, arg := range rule.Run.Args {
		processedArgs[i], _ = state.ProcessString(arg)
		// (Ignorar erro de template por enquanto)
	}
	
	// (Shell 'sh -c' é necessário para interpretar o scriptPath)
	cmdArgs := append([]string{"-c", rule.Run.Script}, processedArgs...)
	cmd := exec.Command("sh", cmdArgs...)

	// Processar 'stdin' (FR10)
	// (Esta implementação será melhorada na Fase 1.5)
	stdinData, _ := state.ProcessStdin("{{ . }}") // Envia o estado inteiro
	cmd.Stdin = bytes.NewReader(stdinData)

	// 3. Preparar streaming de Stdout/Stderr (para FR11)
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	// 4. Iniciar o comando
	if err := cmd.Start(); err != nil {
		return func() tea.Msg { return events.RuntimeErrorMsg{Err: err} }
	}

	// 5. Armazenar o novo processo (FR11)
	if targetID != "" {
		r.activeProcesses[targetID] = cmd.Process
	}

	// 6. Criar o comando (Cmd) do Bubbletea para escutar o streaming
	return r.watchScript(cmd, stdoutPipe, stderrPipe, targetID)
}

// watchScript retorna um tea.Cmd que escuta o stdout/stderr do processo
// e emite ScriptStdoutMsg.
func (r *Runner) watchScript(cmd *exec.Cmd, stdout, stderr io.Reader, targetID string) tea.Cmd {
	return func() tea.Msg {
		// Canal para agregar stdout e stderr
		stream := make(chan []byte)

		// Goroutine para ler stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				stream <- scanner.Bytes()
			}
			close(stream) // Fecha o canal quando o stdout terminar
		}()
		
		// Goroutine para ler stderr (mesmo canal)
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				stream <- scanner.Bytes()
			}
		}()

		// Aguarda o término do comando em outra goroutine
		go func() {
			cmd.Wait()
			if targetID != "" {
				delete(r.activeProcesses, targetID)
			}
		}()
		
		// Escuta o canal e retorna a primeira linha de saída.
		// (Uma implementação robusta usaria tea.Tick para
		// enviar múltiplas mensagens)
		
		// Implementação de Streaming (corrigida):
		// Retorna a *primeira* mensagem, o LayoutManager
		// tratará as seguintes através do Tick (se necessário).
		// Por agora, vamos retornar a primeira linha.
		
		// (Implementação Simplificada Fase 1.4: envia a primeira linha)
		// (Implementação Robusta Fase 1.4: deve usar um Tick)
		
		// Vamos usar um tea.Cmd que retorna o *canal*
		// O Bubbletea (via LayoutManager) precisará de um
		// tea.Tick para ler este canal.
		
		// Implementação Correta (FR11 Streaming):
		// A função 'watchScript' em si não pode ser um tea.Cmd
		// se ela precisa emitir múltiplas mensagens.
		
		// REFAZENDO: O 'HandleRunRequest' deve retornar um tea.Cmd
		// que *inicia* o processo, e o LayoutManager deve
		// iniciar um 'Tick' para ler os pipes.
		
		// REFAZENDO (Simplificado para Fase 1.4):
		// Vamos ler o stdout/stderr em *uma* goroutine
		// e enviar *múltiplas* mensagens ScriptStdoutMsg.
		
		// O 'watchScript' não é o tea.Cmd. O tea.Cmd é 'func() tea.Msg'
		// que *inicia* o processo e as goroutines de escuta.
		
		// A implementação anterior (HandleRunRequest) está correta
		// em retornar `watchScript`. Bubbletea tratará a emissão
		// de Múltiplas mensagens se a função Cmd as enviar.
		
		// (Confusão: Um tea.Cmd só pode retornar *uma* tea.Msg)
		
		// SOLUÇÃO (Padrão Charm):
		// O `HandleRunRequest` inicia o Cmd.
		// O `HandleRunRequest` também retorna um `tea.Tick`
		// que chama uma nova `checkPipesMsg`
		
		// (Esta Fase 1.4 está ficando complexa. Vamos simplificar.)
		
		// ***Otimização de Risco (Simplificação):***
		// O tea.Cmd (watchScript) irá bloquear, ler *todo* o stdout/stderr,
		// e retornar *uma* ScriptStdoutMsg com o *conteúdo completo*.
		// O streaming real (FR11) será refatorado na Fase 1.5.
		
		// (Implementação da Otimização de Risco acima)
		
		var outBuf bytes.Buffer
		merged := io.MultiReader(stdout, stderr)
		io.Copy(&outBuf, merged) // Bloqueia até o Cmd terminar
		
		cmd.Wait()
		if targetID != "" {
			delete(r.activeProcesses, targetID)
		}
		
		// Retorna UMA mensagem com o conteúdo completo.
		return events.ScriptStdoutMsg{
			TargetID: targetID,
			Chunk:    outBuf.Bytes(),
		}
	}
}
```

-----

### Tarefa 1.4.4: Integrar Motores no LayoutManager (Modificação)

Finalmente, **modificamos** o `internal/runtime/layout/manager.go` (da Fase 1.3) para *usar* os novos motores de lógica.

  * ***Otimização de Risco:*** O `LayoutManager` agora implementa o "Padrão 3: Gestor Duplo". Ele gere o Foco, e delega Eventos (Lógica) ao `EventManager` e Ações (Execução) ao `ScriptRunner`.

<!-- end list -->

```go
// internal/runtime/layout/manager.go (MODIFICADO)
package layout

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/components/buttongroup"
	"github.com/helton-godoy/shantilly/internal/components/list"
	"github.com/helton-godoy/shantilly/internal/components/viewport"
	"github.com/helton-godoy/shantilly/internal/runtime/event" // <-- ADICIONADO
	"github.com/helton-godoy/shantilly/internal/runtime/runner" // <-- ADICIONADO
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Manager é o modelo Bubbletea raiz, implementando o "Gestor Duplo"
// (Layout e Foco) conforme Padrão 3 da arquitetura v2.0.
// Fonte: docs/architecture.md (Seção 5.4)
type Manager struct {
	theme *tui.Theme
	root  *declarative.LayoutNode

	// Estado do Layout
	width  int
	height int

	// Estado do Foco
	components map[string]tui.ShantillyComponent
	focusOrder []string
	focusedIdx int

	// Estado de Erro
	lastErr error
	
	// (Debug - removido 'lastEvent')
	
	// Motores de Lógica (ADICIONADO)
	// Fonte: docs/architecture/high-level-architecture.md
	eventManager *event.Manager
	scriptRunner *runner.Runner
	templateState *runner.TemplateState
}

func NewManager(root *declarative.LayoutNode, theme *tui.Theme, rules []declarative.Logic) *Manager { // <-- Assinatura mudou
	m := &Manager{
		theme:      theme,
		root:       root,
		components: make(map[string]tui.ShantillyComponent),
		focusOrder: []string{},
		focusedIdx: 0,
		
		// Motores de Lógica (ADICIONADO)
		eventManager: event.New(rules),
		scriptRunner: runner.New(),
		templateState: runner.NewTemplateState(),
	}

	m.buildComponentTree(root)
	return m
}

// (buildComponentTree não muda - permanece o da Fase 1.3)
func (m *Manager) buildComponentTree(node *declarative.LayoutNode) {
	if node == nil { return }
	if node.Type == "box" && node.Component != nil {
		if node.ID == "" {
			m.lastErr = fmt.Errorf("componente do tipo '%s' não possui 'id'", node.Component.Type)
			return
		}
		var cmp tui.ShantillyComponent
		switch node.Component.Type {
		case "list":
			cmp = list.New(node.ID, m.theme, node.Component.Items)
		case "viewport":
			cmp = viewport.New(node.ID, m.theme, node.Component.Source)
		case "buttongroup":
			cmp = buttongroup.New(node.ID, m.theme, node.Component.Items)
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


func (m *Manager) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range m.components {
		cmds = append(cmds, c.Init())
	}
	return tea.Batch(cmds...)
}

// Update (MODIFICADO) - agora delega lógica.
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// --- Delegação de Lógica (ADICIONADO) ---
	// 1. Componente emite evento (ex: list:select)
	case events.ShantillyEvent:
		// Atualiza o estado que o template usa (FR10)
		m.templateState.UpdateState(msg) 
		// Devolve ao EventManager para verificar regras 'on:' (FR8)
		cmd = m.eventManager.ProcessEvent(msg)
		cmds = append(cmds, cmd)

	// 2. EventManager pede para rodar script
	case events.RunScriptRequestMsg:
		// Devolve ao ScriptRunner para executar (FR9, FR11)
		cmd = m.scriptRunner.HandleRunRequest(msg, m.templateState)
		cmds = append(cmds, cmd)

	// 3. ScriptRunner envia saída (stdout/stderr)
	case events.ScriptStdoutMsg:
		// Envia para o componente 'viewport' correto (FR11)
		if targetCmp, ok := m.components[msg.TargetID]; ok {
			newCmp, cmd := targetCmp.Update(msg)
			m.components[msg.TargetID] = newCmp
			cmds = append(cmds, cmd)
		}

	// ------------------------------------------

	case events.RuntimeErrorMsg:
		m.lastErr = msg.Err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if len(m.focusOrder) > 0 {
				m.focusedIdx = (m.focusedIdx + 1) % len(m.focusOrder)
			}
		case "shift+tab":
			if len(m.focusOrder) > 0 {
				m.focusedIdx = (m.focusedIdx - 1 + len(m.focusOrder)) % len(m.focusOrder)
			}
		
		// Encaminha outras teclas para o componente focado
		default:
			if m.lastErr == nil && len(m.focusOrder) > 0 {
				focusedID := m.focusOrder[m.focusedIdx]
				focusedCmp := m.components[focusedID]
				newCmp, cmd := focusedCmp.Update(msg)
				m.components[focusedID] = newCmp
				cmds = append(cmds, cmd)
			}
		}
	
	// Encaminha mensagens não-tecla para o componente focado
	default:
		if m.lastErr == nil && len(m.focusOrder) > 0 {
			focusedID := m.focusOrder[m.focusedIdx]
			focusedCmp := m.components[focusedID]
			newCmp, cmd := focusedCmp.Update(msg)
			m.components[focusedID] = newCmp
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View (MODIFICADO) - removido o debug 'lastEvent'
func (m *Manager) View() string {
	if m.width == 0 || m.height == 0 {
		return "Inicializando..."
	}

	// 1. Renderizar a árvore de layout principal
	layoutView := m.renderNode(m.root, m.width, m.height)

	// 2. Renderizar Erros
	if m.lastErr != nil {
		errorView := m.theme.ErrorText.Width(m.width).Render("Erro: " + m.lastErr.Error())
		return errorView
	}

	return layoutView
}

// (render.go e getComponent não mudam - permanecem os da Fase 1.3)
// ... (funções renderNode e getComponent da Fase 1.3) ...
```

-----

### Tarefa 1.4.4: Modificar o Ponto de Entrada do Runtime (Integração)

Finalmente, **modificamos** o `internal/runtime/runtime.go` (da Fase 1.2) para passar as regras `on:` para o `LayoutManager`.

```go
// internal/runtime/runtime.go (MODIFICADO)
package runtime

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/runtime/layout"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
)

// Start inicializa e executa o TUI Runtime com a configuração fornecida.
func Start(config *declarative.Config) error {
	// 1. Criar o tema central
	theme := tui.NewDefaultTheme()

	// 2. Criar o LayoutManager (MODIFICADO)
	// Agora passamos a raiz do layout E as regras de lógica (FR8)
	layoutManager := layout.NewManager(&config.Root, theme, config.On)

	// 3. Iniciar o programa Bubbletea
	p := tea.NewProgram(
		layoutManager,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("erro durante a execução do runtime TUI: %w", err)
	}

	return nil
}
```