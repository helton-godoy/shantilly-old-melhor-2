package event

// Epic: 1 - Runtime TUI Declarativo — Fundação do Runtime (v2.0)
// Bloco: E1.3 - Event Engine (EventManager + ShantillyEvent + on:)
//
// Referências normativas:
// - docs/stories/3.3.event-engine-and-on-routing.story.md
// - docs/architecture/data-models.md#4-onhandler-e13--bloco-on-como-unica-orquestracao
// - docs/architecture/data-models.md#5-shantillyevent-e13--tipo-de-evento-normalizado
// - docs/architecture/components.md#4-eventmanager--shantillyevent--e13
// - docs/architecture/core-workflows.md
// - docs/qa/matrix-epic-1-runtime-tui-coverage.md#bloco-e13--event-engine-eventmanager--shantillyevent--on
//
// Invariantes (enforced pelo Orchestrator):
// - O EventManager é o roteador ÚNICO de eventos declarativos (E1.3).
// - TODO fluxo reativo relevante segue:
//   ShantillyComponent/motores -> ShantillyEvent -> EventManager -> on: -> (RunAction|ModalRequest|UpdateState) -> motores dedicados.
// - É proibido executar scripts aqui: apenas delegação para ScriptRunner.
// - É proibido abrir modais diretamente: apenas via ModalRequest produzido a partir de on:.
// - Handlers soltos fora de on: são proibidos.
// - Opera apenas sobre modelos v2.0 (ShantillyEvent, OnHandler, RunAction, ModalRequest) definidos em pkg/declarative.

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"shantilly/pkg/declarative"
	"shantilly/pkg/tui"
)

// ShantillyEvent é o tipo canônico de evento interno do runtime v2.0.
// (Espelha docs/architecture/data-models.md#5-shantillyevent-e13--tipo-de-evento-normalizado.)
// Em implementação completa, deve ser importado de pkg/declarative para evitar duplicação.
type ShantillyEvent struct {
	SourceComponentID string
	Type              string
	Payload           map[string]interface{}
}

// RunAction minimal espelha o contrato declarativo (ver docs/architecture/data-models.md#6-runaction-e14--execucao-declarativa).
// Aqui mantido como placeholder para permitir o esqueleto do EventManager sem acoplamento direto ao ScriptRunner.
type RunAction struct {
	Script        string
	Args          []string
	Stdin         map[string]interface{}
	Env           map[string]string
	UpdateTarget  string
	Confirm       bool
	PromptSecrets []string
}

// ModalRequest minimal espelha o contrato declarativo de abertura de modais.
type ModalRequest struct {
	Type    string
	Title   string
	Message string
	Fields  []ModalField
}

type ModalField struct {
	Name string
	Type string
}

// OnHandler minimal espelha docs/architecture/data-models.md#4-onhandler-e13--bloco-on-como-unica-orquestracao.
type OnHandler struct {
	ID          string
	Event       string
	When        string
	Run         *RunAction
	OpenModal   *ModalRequest
	UpdateState map[string]interface{}
}

// EngineActions representa as saídas puramente declarativas do EventManager
// para outros motores do runtime.
//
// Nenhuma execução ou side-effect ocorre aqui: é apenas o resultado
// da resolução de on: para um determinado ShantillyEvent.
type EngineActions struct {
	RunActions   []*RunAction
	Modals       []*ModalRequest
	StateUpdates []map[string]interface{}
}

// Manager é o EventManager normativo.
//
// Responsabilidades:
// - Manter a lista de OnHandler (proveniente de AppConfig.On).
// - Receber ShantillyEvent.
// - Encontrar handlers compatíveis (Event match).
// - Gerar EngineActions para ScriptRunner/Modal Stack/estado declarativo.
//
// Ele NÃO:
// - Executa scripts.
// - Abre modais diretamente.
// - Atualiza estado de forma ad-hoc.
//
// Implementação concreta de avaliação de `When` e integração TEA
// será feita por bmad-dev/bmad-master conforme contratos.
type Manager struct {
	// handlers normativos baseados em OnHandler local (Wave atual)
	handlers []OnHandler

	// rules normativos baseados em declarative.OnHandler (v2.0)
	// usados pelo pipeline TEA + pkg/tui.
	rules []declarative.OnHandler

	// pendingRequest armazena a RunScriptRequestMsg pendente de confirmação JIT
	// (inspirado em 24_chat_ia.md, Tarefa 1.6.3).
	pendingRequest *tui.RunScriptRequestMsg
}

// NewManager cria um EventManager a partir de uma lista de handlers declarativos
// no formato local OnHandler (wave anterior). Mantido para compatibilidade.
// Chamadores devem garantir que esta lista foi derivada de AppConfig.On validado.
func NewManager(handlers []OnHandler) *Manager {
	return &Manager{
		handlers: handlers,
	}
}

// New cria um EventManager a partir de uma lista de declarative.OnHandler (v2.0).
// Este é o construtor preferencial para o runtime declarativo atual.
func New(rules []declarative.OnHandler) *Manager {
	return &Manager{
		rules: rules,
	}
}

// Handle recebe um ShantillyEvent normalizado e retorna ações declarativas
// derivadas do bloco on:, sem executar efeitos colaterais.
//
// Regras baseline:
// - Match exato por Event == fmt.Sprintf("%s:%s", SourceComponentID, Type).
// - Ordem determinística: mesma ordem dos handlers no YAML.
// - `When` é reservado para futura avaliação; por ora, ignorado se não vazio (não executa código arbitrário).
func (m *Manager) Handle(ev ShantillyEvent) EngineActions {
	key := fmt.Sprintf("%s:%s", ev.SourceComponentID, ev.Type)

	var actions EngineActions

	for _, h := range m.handlers {
		if h.Event != key {
			continue
		}

		// Placeholder para When: manter sem interpretar nesta wave.
		// Invariante: nenhum código condicional ad-hoc complexo aqui.

		if h.Run != nil {
			actions.RunActions = append(actions.RunActions, h.Run)
		}
		if h.OpenModal != nil {
			actions.Modals = append(actions.Modals, h.OpenModal)
		}
		if len(h.UpdateState) > 0 {
			actions.StateUpdates = append(actions.StateUpdates, h.UpdateState)
		}
	}

	return actions
}

// ProcessEvent é a API orientada a TEA que trabalha com tui.ShantillyEvent
// e produz comandos Bubble Tea. Ela implementa a lógica JIT descrita na
// arquitetura (Seção 4.6 e 7.3) e em 24_chat_ia.md.
//
// Fluxo (simplificado):
//  1. Normaliza chave de evento "source:type".
//  2. Encontra declarative.OnHandler compatíveis.
//  3. Se handler.Run for nil, não há execução de script.
//  4. Constrói tui.RunScriptRequestMsg com o handler completo.
//  5. Se handler.Run.Confirm == true, aplica JIT:
//     - Armazena em pendingRequest.
//     - Emite tui.ShowModalMsg com modal de confirmação (tratado em camada superior).
//  6. Caso contrário, emite diretamente a RunScriptRequestMsg.
func (m *Manager) ProcessEvent(ev tui.ShantillyEvent) tea.Cmd {
	if len(m.rules) == 0 {
		return nil
	}

	key := fmt.Sprintf("%s:%s", ev.ComponentID, ev.Type)

	for _, rule := range m.rules {
		if rule.Event != key {
			continue
		}

		if rule.Run == nil {
			// Nada para executar neste handler.
			continue
		}

		// Clonamos o handler para poder enriquecer Run.Env com dados do evento
		// sem mutar a configuração original carregada do YAML.
		updated := rule
		if updated.Run != nil {
			// Copia superficial de RunAction
			runCopy := *updated.Run
			// Copia defensiva do mapa Env (para evitar aliasing entre handlers).
			if runCopy.Env != nil {
				copiedEnv := make(map[string]string, len(runCopy.Env)+1)
				for k, v := range runCopy.Env {
					copiedEnv[k] = v
				}
				runCopy.Env = copiedEnv
			}

			// Binding automático input -> env
			if valueStr, ok := extractInputValue(ev); ok {
				if runCopy.Env == nil {
					runCopy.Env = make(map[string]string, 1)
				}
				key := buildInputEnvKey(ev.ComponentID)
				if _, exists := runCopy.Env[key]; !exists {
					runCopy.Env[key] = valueStr
				}
			}

			// Binding automático select -> env (ID único selecionado).
			if selectID, ok := extractSelectValue(ev); ok {
				if runCopy.Env == nil {
					runCopy.Env = make(map[string]string, 1)
				}
				key := buildSelectEnvKey(ev.ComponentID)
				if _, exists := runCopy.Env[key]; !exists {
					runCopy.Env[key] = selectID
				}
			}

			// Binding automático multiselect -> env (lista de IDs selecionados).
			if ids, ok := extractMultiSelectValues(ev); ok {
				if runCopy.Env == nil {
					runCopy.Env = make(map[string]string, 1)
				}
				key := buildMultiSelectEnvKey(ev.ComponentID)
				if _, exists := runCopy.Env[key]; !exists {
					runCopy.Env[key] = strings.Join(ids, ",")
				}
			}

			updated.Run = &runCopy
		}

		req := tui.RunScriptRequestMsg{Handler: updated}

		// Lógica JIT básica: se Confirm ligado, armazenamos a requisição
		// e deixamos o modal/overlay decidir.
		if rule.Run.Confirm {
			m.pendingRequest = &req
			// O modal em si será criado em nível superior (LayoutManager/MainModel),
			// aqui apenas sinalizamos que deve haver um modal de confirmação.
			return func() tea.Msg {
				return tui.ShowModalMsg{
					Content: nil, // o MainModel/layout é responsável por instanciar o modal concreto
				}
			}
		}

		// Sem JIT: executa imediatamente.
		return func() tea.Msg { return req }
	}

	return nil
}

// HandleModalResult trata o resultado vindo do modal de confirmação
// (tui.ModalResultMsg). Se houver uma pendingRequest e o usuário confirmar,
// devolve o tui.RunScriptRequestMsg correspondente; caso contrário, apenas
// limpa o estado pendente.
func (m *Manager) HandleModalResult(msg tui.ModalResultMsg) tea.Cmd {
	if m.pendingRequest == nil {
		return nil
	}

	req := m.pendingRequest
	m.pendingRequest = nil

	if !msg.Confirmed {
		return nil
	}

	return func() tea.Msg { return *req }
}

// extractInputValue tenta extrair um valor textual do payload do evento de input.
// Espera payload no formato map[string]interface{}{ "value": <string> } enviado pelo
// componente input, mas é tolerante a outros formatos simples.
func extractInputValue(ev tui.ShantillyEvent) (string, bool) {
	if ev.Payload == nil {
		return "", false
	}

	switch v := ev.Payload.(type) {
	case map[string]interface{}:
		if raw, ok := v["value"]; ok && raw != nil {
			return fmt.Sprint(raw), true
		}
	case string:
		return v, true
	default:
		return fmt.Sprint(v), true
	}

	return "", false
}

// buildInputEnvKey constrói a chave de variável de ambiente padronizada
// para o valor de um componente input, seguindo a convenção:
//
//	SHANTILLY_INPUT_<ID_EM_UPPER_SNAKE>
func buildInputEnvKey(componentID string) string {
	id := strings.TrimSpace(componentID)
	if id == "" {
		return "SHANTILLY_INPUT"
	}

	// Normaliza separadores comuns para underscore e aplica UPPERCASE.
	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "/", "_")
	id = strings.ToUpper(id)

	return "SHANTILLY_INPUT_" + id
}

func extractSelectValue(ev tui.ShantillyEvent) (string, bool) {
	if ev.Payload == nil {
		return "", false
	}

	switch v := ev.Payload.(type) {
	case map[string]interface{}:
		if raw, ok := v["id"]; ok && raw != nil {
			return fmt.Sprint(raw), true
		}
	case string:
		return v, true
	default:
		return fmt.Sprint(v), true
	}

	return "", false
}

func buildSelectEnvKey(componentID string) string {
	id := strings.TrimSpace(componentID)
	if id == "" {
		return "SHANTILLY_SELECT"
	}

	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "/", "_")
	id = strings.ToUpper(id)

	return "SHANTILLY_SELECT_" + id
}

func extractMultiSelectValues(ev tui.ShantillyEvent) ([]string, bool) {
	if ev.Payload == nil {
		return nil, false
	}

	var ids []string

	switch v := ev.Payload.(type) {
	case map[string]interface{}:
		raw, ok := v["ids"]
		if !ok || raw == nil {
			return nil, false
		}
		switch vv := raw.(type) {
		case []string:
			if len(vv) == 0 {
				return nil, false
			}
			ids = append(ids, vv...)
		case []interface{}:
			for _, e := range vv {
				ids = append(ids, fmt.Sprint(e))
			}
		default:
			ids = append(ids, fmt.Sprint(vv))
		}
	case []string:
		if len(v) == 0 {
			return nil, false
		}
		ids = append(ids, v...)
	case []interface{}:
		for _, e := range v {
			ids = append(ids, fmt.Sprint(e))
		}
	default:
		return nil, false
	}

	if len(ids) == 0 {
		return nil, false
	}

	return ids, true
}

func buildMultiSelectEnvKey(componentID string) string {
	id := strings.TrimSpace(componentID)
	if id == "" {
		return "SHANTILLY_MULTISELECT"
	}

	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "/", "_")
	id = strings.ToUpper(id)

	return "SHANTILLY_MULTISELECT_" + id
}
