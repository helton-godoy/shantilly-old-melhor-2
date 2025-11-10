// Package modal fornece funções de integração entre a Modal Stack e o pipeline
// declarativo baseado em ShantillyEvent + EventManager.
//
// Este arquivo define helpers puros para:
//
// - Consumir ModalRequest (produzidos via `on:` pelo EventManager).
// - Empilhar/remover modais na ModalStack.
// - Gerar ShantillyEvent `modal.*` de resposta, que SEMPRE retornam ao EventManager.
//
// Invariantes normativos (E1.5-1/2/4, gates 1.x.modal-stack.yml, 1.x.security-jit-anti-trojan.yml):
//
//   - Nenhum modal é criado/gerido fora de internal/runtime/modal/**.
//   - Qualquer interação com modais flui por:
//     ShantillyEvent -> EventManager -> on: -> ModalRequest -> ModalStack
//     ModalStack -> ShantillyEvent (`modal.*`) -> EventManager -> on:
//   - Apenas o topo da pilha recebe eventos; fundo bloqueado enquanto HasActive() == true.
//   - Segredos nunca são logados/persistidos; estes helpers não fazem logging.
//   - Nenhum acesso direto ao ScriptRunner aqui; sempre via EventManager + on:.
//
// IMPORTANTE: Este pacote não conhece detalhes de UI (Bubble Tea, etc.) nem YAML;
// ele opera apenas sobre ModalRequest, Modal e ShantillyEvent.
package modal

// OpenModalFromRequest aplica um ModalRequest na stack,
// retornando o Modal criado e (opcionalmente) um evento `modal.opened`.
//
// Cabe ao chamador decidir se deseja emitir o evento `modal.opened`.
// Esta função NUNCA faz logging ou efeitos externos.
func OpenModalFromRequest(stack *ModalStack, req ModalRequest) (Modal, *ShantillyEvent) {
	m := stack.Push(req)

	// Evento opcional de rastreabilidade interna.
	ev := &ShantillyEvent{
		SourceComponentID: m.ID,
		Type:              ModalOpenedType,
		Payload: map[string]interface{}{
			ModalIDKey:   m.ID,
			ModalTypeKey: m.Request.Type,
		},
	}

	return m, ev
}

// ConfirmModal registra a confirmação de um modal identificado por modalID,
// gera o evento `modal.confirmed` e fecha o modal (pop).
//
// Esta função assume que o chamador já validou que o modalID corresponde
// ao topo ou é roteado corretamente via EventManager.
//
// Regras de segurança:
// - values pode conter segredos; não são logados aqui.
// - stack.Pop() garante remoção do modal e liberação de referência.
func ConfirmModal(stack *ModalStack, modalID string, values map[string]interface{}) *ShantillyEvent {
	// Removemos sempre o topo — o EventManager deve garantir roteamento correto.
	m, ok := stack.Pop()
	if !ok {
		return nil
	}
	// Se o ID não bater, o chamador deve tratar; aqui seguimos política minimalista.
	if m.ID != modalID && modalID != "" {
		// Em caso de mismatch, não reaplicamos o modal; devolvemos nil para
		// permitir que camada superior trate como erro lógico.
		return nil
	}

	payload := map[string]interface{}{
		ModalIDKey:        m.ID,
		ModalTypeKey:      m.Request.Type,
		ModalConfirmedKey: true,
	}
	if len(values) > 0 {
		payload[ModalValuesKey] = values
	}

	return &ShantillyEvent{
		SourceComponentID: m.ID,
		Type:              ModalConfirmedType,
		Payload:           payload,
	}
}

// CancelModal registra o cancelamento de um modal (ex.: usuário abortou),
// gera `modal.cancelled` e fecha o topo da pilha.
//
// Nenhum dado sensível é incluído.
func CancelModal(stack *ModalStack, modalID string) *ShantillyEvent {
	m, ok := stack.Pop()
	if !ok {
		return nil
	}
	if m.ID != modalID && modalID != "" {
		return nil
	}

	return &ShantillyEvent{
		SourceComponentID: m.ID,
		Type:              ModalCancelledType,
		Payload: map[string]interface{}{
			ModalIDKey:        m.ID,
			ModalTypeKey:      m.Request.Type,
			ModalConfirmedKey: false,
		},
	}
}

// InputModal registra entrada de dados para o modal de topo sem fechá-lo,
// emitindo `modal.input`.
//
// Uso típico: campos de texto, segredos, etc., enviados pelo TUI como ShantillyEvent
// específicos que o EventManager converte e encaminha.
//
// Regras de segurança:
// - values pode conter segredos; o chamador deve garantir manuseio somente em memória.
// - Esta função não persiste nem loga; apenas monta o evento.
func InputModal(stack *ModalStack, modalID string, values map[string]interface{}) *ShantillyEvent {
	top := stack.Top()
	if top == nil {
		return nil
	}
	if top.ID != modalID && modalID != "" {
		return nil
	}

	if len(values) == 0 {
		return nil
	}

	payload := map[string]interface{}{
		ModalIDKey:     top.ID,
		ModalTypeKey:   top.Request.Type,
		ModalValuesKey: values,
	}

	return &ShantillyEvent{
		SourceComponentID: top.ID,
		Type:              ModalInputType,
		Payload:           payload,
	}
}
