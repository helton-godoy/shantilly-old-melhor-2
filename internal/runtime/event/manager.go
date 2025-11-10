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
	handlers []OnHandler
}

// NewManager cria um EventManager a partir de uma lista de handlers declarativos.
// Chamadores devem garantir que esta lista foi derivada de AppConfig.On validado.
func NewManager(handlers []OnHandler) *Manager {
	return &Manager{
		handlers: handlers,
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
