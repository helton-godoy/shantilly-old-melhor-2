// Package modal implementa a Modal Stack normativa do Runtime TUI Declarativo v2.0.
//
// Mandatos (E1.5, Wave 5 — Modal Stack + Segurança JIT):
//   - Implementação exclusiva em internal/runtime/modal/** (docs/architecture/components.md#6-modal-stack--e15-wave-5-parte-1).
//   - Pilha explícita de modais; apenas o topo recebe foco/eventos; fundo bloqueado enquanto houver modal ativo.
//   - Toda abertura de modal via pipeline único:
//     ShantillyEvent -> EventManager -> on: -> ModalRequest -> Modal Stack
//     (docs/architecture/governance-runtime-tui-v2.0.md#4-wave-5--modal-stack--seguranca-jit-e15--parte-1).
//   - Toda resposta de modal retorna como ShantillyEvent (`modal.*`) para o EventManager, nunca via atalhos.
//
// Regras de segurança JIT (E1.5-1, E1.5-4, 1.x.modal-stack.yml, 1.x.security-jit-anti-trojan.yml):
// - `RunAction.Confirm` e `RunAction.PromptSecrets` devem acionar Modal Stack antes de qualquer execução.
// - Segredos coletados via Modal Stack:
//   - Nunca são logados.
//   - Nunca são persistidos em disco.
//   - Vivem apenas em memória, pelo tempo mínimo necessário para compor o RunAction final.
//
// - É proibido criar/gerir modais fora deste pacote.
package modal

// ShantillyEvent é definido aqui como contrato mínimo local para evitar dependência circular.
// Integração real com o EventManager e demais engines deve usar este shape.
// Referência: docs/architecture/data-models.md#5-shantillyevent-e13--tipo-de-evento-normalizado
type ShantillyEvent struct {
	SourceComponentID string                 `json:"source_component_id"`
	Type              string                 `json:"type"`
	Payload           map[string]interface{} `json:"payload,omitempty"`
}

// ModalRequest representa o pedido declarativo para abrir um modal,
// oriundo do pipeline `on:` -> `OpenModal`/`ModalRequest`.
//
// Fonte normativa: docs/architecture/data-models.md#7-modalrequest-e15--modal-stack
//
// Tipos previstos (não exaustivo, porém restrito a usos seguros):
// - "confirm"       — confirmação simples.
// - "prompt"        — coleta de um ou mais campos.
// - "confirm+prompt"— combinação de confirmação + segredos.
// Campos sensíveis DEVEM ser marcados adequadamente em ModalField.
type ModalRequest struct {
	// ID lógico opcional fornecido pelo chamador para correlacionar respostas.
	// Não é obrigatório; se vazio, a Modal Stack pode gerar um ID.
	ID string `json:"id,omitempty"`

	// Type indica o tipo de modal. Implementações devem validar/whitelistar.
	Type string `json:"type"`

	Title   string       `json:"title,omitempty"`
	Message string       `json:"message,omitempty"`
	Fields  []ModalField `json:"fields,omitempty"`

	// Sensitive indica se o modal lida com informações sensíveis (ex.: prompt_secrets).
	// Quando true:
	// - Valores não podem aparecer em logs.
	// - Implementações devem garantir limpeza em memória após uso.
	Sensitive bool `json:"sensitive,omitempty"`

	// Metadata opcional para integração futura (ex.: origem declarativa, hints de layout).
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ModalField define um campo solicitado por um ModalRequest.
//
// Exemplos típicos:
// - Name: "vault_pass", Type: "secret", Sensitive: true
// - Name: "confirm", Type: "bool"
type ModalField struct {
	Name string `json:"name"`
	Type string `json:"type"` // ex.: "text", "secret", "bool"

	// Sensitive indica se o valor deste campo é segredo.
	// Deve ser respeitado por toda a cadeia (Modal Stack, EventManager, ScriptRunner).
	Sensitive bool `json:"sensitive,omitempty"`
}

// Modal representa uma instância ativa na pilha de modais.
//
// É derivado de um ModalRequest e enriquecido com metadados de runtime.
type Modal struct {
	// ID único da instância do modal na pilha.
	// Preferencialmente derivado de ModalRequest.ID ou gerado pela stack.
	ID string

	// Request original, preservado para contexto de renderização/validação.
	Request ModalRequest

	// Active indica se o modal está atualmente em foco/topo.
	// Apenas o modal no topo da pilha pode receber eventos.
	Active bool

	// Sensitive indica se QUALQUER parte do modal é sensível.
	// É true se ModalRequest.Sensitive ou qualquer campo sensível.
	Sensitive bool
}

// ModalEventType contém tipos de eventos `modal.*` emitidos pela Modal Stack.
// Esses eventos SEMPRE retornam ao EventManager como ShantillyEvent.Type,
// respeitando o pipeline único.
//
// Eventos normativos mínimos:
//
// - "modal.opened"    — opcional, indica que um modal foi aberto/empilhado.
// - "modal.confirmed" — usuário confirmou uma ação.
// - "modal.cancelled" — usuário cancelou/fechou sem confirmar.
// - "modal.input"     — usuário forneceu entradas (incluindo segredos).
const (
	ModalOpenedType    = "modal.opened"
	ModalConfirmedType = "modal.confirmed"
	ModalCancelledType = "modal.cancelled"
	ModalInputType     = "modal.input"
)

// ModalEventPayloadKeys define chaves padrão para Payload de ShantillyEvent
// emitidos pela Modal Stack.
const (
	// ModalIDKey identifica o modal alvo na resposta/evento.
	ModalIDKey = "modal_id"

	// ModalTypeKey replica o tipo (confirm/prompt/etc.) para roteamento declarativo.
	ModalTypeKey = "modal_type"

	// ModalValuesKey carrega valores inseridos pelo usuário.
	// IMPORTANTE:
	// - Quando contiver segredos, deve ser tratado exclusivamente em memória.
	// - É proibido logar ou persistir esses valores.
	ModalValuesKey = "values"

	// ModalConfirmedKey indica confirmação explícita (bool).
	ModalConfirmedKey = "confirmed"
)

// NewModalFromRequest constrói um Modal a partir de um ModalRequest,
// calculando a flag Sensitive de forma conservadora.
//
// Esta função NÃO realiza logging nem qualquer I/O externo, em respeito
// às regras de segurança JIT.
func NewModalFromRequest(req ModalRequest, id string) Modal {
	sensitive := req.Sensitive
	for _, f := range req.Fields {
		if f.Sensitive || f.Type == "secret" {
			sensitive = true
			break
		}
	}

	return Modal{
		ID:        id,
		Request:   req,
		Active:    false, // Ativado pela Modal Stack quando estiver no topo.
		Sensitive: sensitive,
	}
}
