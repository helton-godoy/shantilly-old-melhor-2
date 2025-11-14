Executado. Concluímos a **Fase 2.1 (Extensão dos Contratos)**. Os nossos modelos de dados agora reconhecem o runner `ansible_playbook` e o nosso `TemplateState` está pronto para armazenar segredos.

Seguindo o **Plano Mestre Sequencial**, continuamos no **Épico 2**, iniciando a **Fase 2.2: Implementação do Modal JIT (PromptSecrets / Vault Pass)**.

Esta fase implementa a lógica JIT (Just-In-Time) para `prompt_secrets:` e `ask_vault_pass: true`, que dependem da infraestrutura de modais da Fase 1.6.

-----

### Tarefa 2.2.1: Implementar Componente Modal (Secret)

Primeiro, criamos o componente TUI para solicitar segredos (como o `vault_pass`).

  * ***Otimização de Risco:*** Para evitar a "baixa qualidade de código" que mencionaste, usamos `huh.Input` (da Fase 1.5) com `Password(true)`, garantindo consistência visual e funcionalidade correta de "password".

<!-- end list -->

```go
// internal/runtime/modal/secret.go
package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// SecretModel é um componente TUI (tea.Model) para
// recolha de segredos JIT (Just-In-Time).
// Fonte: docs/prd.md (Épico 2 - ask_vault_pass)
type Model struct {
	form      *huh.Form
	secretKeys []string // As chaves que estamos a pedir (ex: "vault_pass")
}

// NewSecret solicita uma lista de segredos ao usuário.
func NewSecret(title string, keys []string) *Model {
	var fields []huh.Field
	for _, key := range keys {
		fields = append(fields,
			huh.NewInput().
				Key(key).
				Title("Enter "+key).
				Password(true), // Garante que a entrada seja mascarada
		)
	}

	form := huh.NewForm(
		huh.NewGroup(fields...).WithTitle(title),
	).WithWidth(50) // Largura fixa para o modal

	return &Model{
		form:       form,
		secretKeys: keys,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	// Se o formulário de segredos for concluído
	if m.form.State == huh.StateCompleted {
		secrets := make(map[string]string)
		for _, key := range m.secretKeys {
			secrets[key] = m.form.Get(key).(string)
		}

		// Emite o evento de resultado (Tarefa 1.6.1) com o Payload
		cmd = func() tea.Msg {
			return events.ModalResultMsg{
				Confirmed: true,
				Payload:   secrets, // (Payload agora contém os segredos)
			}
		}
		cmds = append(cmds, cmd)
	}

	// Se for Abortado (Ctrl+C, Esc)
	if m.form.State == huh.StateAborted {
		// Emite "Não Confirmado"
		cmd = func() tea.Msg { return events.ModalResultMsg{Confirmed: false} }
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// Reutiliza o estilo do modal de confirmação (Fase 1.6.2)
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(m.form.View())
}
```

-----

### Tarefa 2.2.2: Modificar EventManager (Lógica JIT em Cadeia)

Esta é a "visão do futuro" em ação. **Modificamos** o `internal/runtime/event/manager.go` (da Fase 1.6.3) para suportar uma *cadeia* de modais JIT.

  * ***Otimização de Risco:*** A lógica JIT agora é uma *cadeia*:
    1.  `ProcessEvent` (Evento da UI) -\> Verifica `Confirm`? -\> Mostra Modal `Confirm`.
    2.  `HandleModalResult` (Resultado do `Confirm`) -\> Foi confirmado? -\> Verifica `PromptSecrets`? -\> Mostra Modal `Secret`.
    3.  `HandleModalResult` (Resultado do `Secret`) -\> Foi confirmado? -\> Executa a Ação (agora com segredos).

<!-- end list -->

```go
// internal/runtime/event/manager.go (MODIFICADO)
package event

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/runtime/modal"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)

// Manager (Struct permanece o da Fase 1.6.3)
type Manager struct {
	rules []declarative.Logic
	// Armazena a ação pendente de confirmação JIT
	pendingRequest *events.RunScriptRequestMsg
}

func New(rules []declarative.Logic) *Manager {
	return &Manager{
		rules: rules,
	}
}

// ProcessEvent (MODIFICADO)
// Agora é o *início* da cadeia JIT.
func (m *Manager) ProcessEvent(event events.ShantillyEvent) tea.Cmd {
	for _, rule := range m.rules {
		if rule.Event == event.Type {

			req := events.RunScriptRequestMsg{Rule: rule}

			// --- Início da Cadeia JIT (MODIFICADO) ---
			// Fonte: docs/architecture.md (Seção 7.3)

			// 1. Armazena a ação pendente
			m.pendingRequest = &req

			// 2. Inicia a cadeia: Verifica Confirmação (Fase 1.6)
			if rule.Confirm {
				title := fmt.Sprintf("Executar '%s'?", m.getRunnerName(rule))
				modalCmd := func() tea.Msg {
					return events.ShowModalMsg{
						Content: modal.NewConfirm(title),
					}
				}
				return modalCmd
			}

			// 3. (Se não houver 'Confirm') Verifica Segredos (Fase 2.2)
			// (O 'HandleModalResult' tratará o caso de 'Confirm' ser verdadeiro)
			return m.checkNextJITStep(true, nil)
		}
	}
	return nil
}

// checkNextJITStep (ADICIONADO)
// Avança na cadeia JIT. Chamado por ProcessEvent ou HandleModalResult.
func (m *Manager) checkNextJITStep(lastStepConfirmed bool, lastPayload interface{}) tea.Cmd {
	if m.pendingRequest == nil {
		return nil // Nenhuma ação pendente
	}

	if !lastStepConfirmed {
		m.pendingRequest = nil // Abortar cadeia
		return nil
	}
	
	rule := m.pendingRequest.Rule

	// --- Lógica de Segredos (ADICIONADO) ---
	// Se o último passo foi 'Confirm', ou se não havia 'Confirm'.
	
	// 1. Verifica se 'prompt_secrets' ou 'ask_vault_pass' estão definidos
	// (Esta é a verificação *antes* de executar a ação)
	secretsToAsk := []string{}
	secretsToAsk = append(secretsToAsk, rule.PromptSecrets...)
	if rule.Run.AnsiblePlaybook != nil && rule.Run.AnsiblePlaybook.AskVaultPass {
		// (Evita duplicados se o usuário também definir 'vault_pass' em 'prompt_secrets')
		found := false
		for _, s := range secretsToAsk {
			if s == "vault_pass" {
				found = true
				break
			}
		}
		if !found {
			secretsToAsk = append(secretsToAsk, "vault_pass")
		}
	}

	// 2. Se *ainda* não pedimos segredos (verificando o payload)
	// E se *precisamos* pedir segredos
	if _, isSecretPayload := lastPayload.(map[string]string); !isSecretPayload && len(secretsToAsk) > 0 {
		// 3. Emite o ShowModalMsg (Tarefa 2.2.1)
		title := fmt.Sprintf("Segredos para '%s'", m.getRunnerName(rule))
		modalCmd := func() tea.Msg {
			return events.ShowModalMsg{
				Content: modal.NewSecret(title, secretsToAsk),
			}
		}
		return modalCmd
	}
	
	// --- Fim da Lógica de Segredos ---
	
	// Se chegamos aqui, todas as etapas JIT (Confirm e Secrets) passaram.
	req := m.pendingRequest
	m.pendingRequest = nil // Limpa a ação pendente
	
	// 4. Executa a Ação de Script (retornando a msg da Fase 1.4)
	return func() tea.Msg { return *req }
}


// HandleModalResult (MODIFICADO)
// Agora avança na cadeia JIT.
func (m *Manager) HandleModalResult(msg events.ModalResultMsg) tea.Cmd {
	// (O payload do 'Confirm' é nil, o payload do 'Secret' é map[string]string)
	return m.checkNextJITStep(msg.Confirmed, msg.Payload)
}

// getRunnerName (ADICIONADO Helper)
func (m *Manager) getRunnerName(rule declarative.Logic) string {
	if rule.Run.Script != "" {
		return rule.Run.Script
	}
	if rule.Run.AnsiblePlaybook != nil {
		return rule.Run.AnsiblePlaybook.Playbook
	}
	return "Ação Desconhecida"
}
```

-----

### Tarefa 2.2.3: Modificar LayoutManager (Atualização de Estado Secreto)

**Modificamos** o `internal/runtime/layout/manager.go` (da Fase 1.6.6) para garantir que, quando um `ModalResultMsg` (com segredos) chegar, ele atualize o `templateState.Secret`.

```go
// internal/runtime/layout/manager.go (MODIFICADO)
package layout

// (Importações ... permanecem as da Fase 1.6)
import (
	// ...
	"github.com/helton-godoy/shantilly/pkg/tui/events"
)
// (Struct Manager, NewManager, buildComponentTree, Init - permanecem os da Fase 1.6)
// ...

// Update (MODIFICADO) - Atualiza TemplateState.Secret
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// ... (case tea.WindowSizeMsg, case events.ShantillyEvent,
	//      case events.RunScriptRequestMsg, case events.ScriptStdoutMsg
	//      permanecem os da Fase 1.6.6)

	// --- (Lógica da Fase 1.6.6) ---
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case events.ShantillyEvent:
		m.templateState.UpdateState(msg)
		cmd = m.eventManager.ProcessEvent(msg)
		cmds = append(cmds, cmd)
	case events.RunScriptRequestMsg:
		cmd = m.scriptRunner.HandleRunRequest(msg, m.templateState)
		cmds = append(cmds, cmd)
	case events.ScriptStdoutMsg:
		if targetCmp, ok := m.components[msg.TargetID]; ok {
			newCmp, cmd := targetCmp.Update(msg)
			m.components[msg.TargetID] = newCmp
			cmds = append(cmds, cmd)
		}
	// --- Fim da Lógica da Fase 1.6.6 ---

	// --- MODIFICADO (Tarefa 2.2.3) ---
	// 4. Modal fecha e envia resultado
	case events.ModalResultMsg:
		// ADICIONADO: Atualiza o TemplateState com os segredos
		// (Fonte: Fase 2.1, Tarefa 2.1.2)
		if secrets, ok := msg.Payload.(map[string]string); ok {
			for k, v := range secrets {
				m.templateState.Secret[k] = v
			}
		}

		// Devolve ao EventManager para tratar a lógica JIT
		cmd = m.eventManager.HandleModalResult(msg)
		cmds = append(cmds, cmd)
	// ------------------------------------

	case events.ShowModalMsg:
		return m, func() tea.Msg { return msg }

	// (Restante: events.RuntimeErrorMsg, tea.KeyMsg, default
	// permanecem os da Fase 1.6.6)
	// ...
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
		default:
			if m.lastErr == nil && len(m.focusOrder) > 0 {
				focusedID := m.focusOrder[m.focusedIdx]
				focusedCmp := m.components[focusedID]
				newCmp, cmd := focusedCmp.Update(msg)
				m.components[focusedID] = newCmp
				cmds = append(cmds, cmd)
			}
		}
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

// (View, renderNode, getComponent, isFocused - permanecem os da Fase 1.3)
// ...
```