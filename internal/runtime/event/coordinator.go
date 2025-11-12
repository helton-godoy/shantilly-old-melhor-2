package event

import (
	"shantilly/internal/runtime/modal"
	"shantilly/pkg/declarative"
)

// Coordinator: sequencia EngineActions e aplica o Gate de Modais (Confirm/PromptSecrets).
// - Não executa scripts; apenas retorna []declarative.RunAction para o orquestrador chamar o ScriptRunner.
// - Não publica eventos diretamente; retorna []declarative.ShantillyEvent para o orquestrador emitir.
type Coordinator struct {
	stack   *modal.ModalStack
	pending map[string]declarative.RunAction // modalID -> RunAction aguardando confirmação
}

func NewCoordinator(stack *modal.ModalStack) *Coordinator {
	return &Coordinator{
		stack:   stack,
		pending: make(map[string]declarative.RunAction),
	}
}

// ProcessActions recebe EngineActions produzidas por EventManager.Handle(...)
// Retorna eventos de abertura de modal e ações de execução imediata (sem gate).
func (c *Coordinator) ProcessActions(actions EngineActions) (events []declarative.ShantillyEvent, toRun []declarative.RunAction) {
	// 1) Modais diretos oriundos de on:
	for _, m := range actions.Modals {
		req := toModalRequest(m)
		modalInst := c.stack.Push(req)
		events = append(events, declarative.ShantillyEvent{
			Type:    modal.ModalOpenedType,
			Payload: map[string]interface{}{modal.ModalIDKey: modalInst.ID},
		})
	}
	// 2) Gate para RunAction (Confirm/PromptSecrets)
	for _, rptr := range actions.RunActions {
		r := toDeclarativeRunAction(rptr)
		if needsModalGate(r) {
			req := modalForRun(r)
			modalInst := c.stack.Push(req)
			c.pending[modalInst.ID] = r
			events = append(events, declarative.ShantillyEvent{
				Type:    modal.ModalOpenedType,
				Payload: map[string]interface{}{modal.ModalIDKey: modalInst.ID},
			})
			continue
		}
		toRun = append(toRun, r)
	}
	return events, toRun
}

// OnEvent consome modal.confirmed/cancelled e libera execuções pendentes.
func (c *Coordinator) OnEvent(ev declarative.ShantillyEvent) (events []declarative.ShantillyEvent, toRun []declarative.RunAction) {
	switch ev.Type {
	case modal.ModalConfirmedType:
		mid, _ := ev.Payload[modal.ModalIDKey].(string)
		if mid == "" {
			return nil, nil
		}
		_, _ = c.stack.Pop()

		run, ok := c.pending[mid]
		if !ok {
			return nil, nil
		}
		delete(c.pending, mid)

		// Aplica segredos/valores exclusivamente em memória no Stdin, sem log/persistência.
		if vals, ok := ev.Payload[modal.ModalValuesKey].(map[string]interface{}); ok && len(vals) > 0 {
			run = applySecretsToRunAction(run, vals)
		}
		toRun = append(toRun, run)
		return events, toRun

	case modal.ModalCancelledType:
		mid, _ := ev.Payload[modal.ModalIDKey].(string)
		if mid == "" {
			return nil, nil
		}
		_, _ = c.stack.Pop()
		delete(c.pending, mid)
		return nil, nil

	default:
		return nil, nil
	}
}

func needsModalGate(r declarative.RunAction) bool {
	return r.Confirm || len(r.PromptSecrets) > 0
}

func modalForRun(r declarative.RunAction) modal.ModalRequest {
	req := modal.ModalRequest{
		Title:     "Confirmar execução",
		Message:   "Deseja prosseguir com a ação?",
		Sensitive: false,
		Metadata:  map[string]interface{}{"update_target": r.UpdateTarget},
	}
	switch {
	case r.Confirm && len(r.PromptSecrets) > 0:
		req.Type = "confirm+prompt"
	case r.Confirm:
		req.Type = "confirm"
	default:
		req.Type = "prompt"
	}
	for _, name := range r.PromptSecrets {
		req.Fields = append(req.Fields, modal.ModalField{Name: name, Type: "secret", Sensitive: true})
		req.Sensitive = true
	}
	return req
}

// Segredos em memória, aplicados ao Stdin (contrato v2.0). Nunca logar/persistir.
func applySecretsToRunAction(r declarative.RunAction, vals map[string]interface{}) declarative.RunAction {
	if r.Stdin == nil {
		r.Stdin = map[string]interface{}{}
	}
	for k, v := range vals {
		r.Stdin[k] = v
	}
	return r
}

// Conversões auxiliares entre tipos locais (package event) e canônicos (pkg/declarative).

func toDeclarativeRunAction(r *RunAction) declarative.RunAction {
	if r == nil {
		return declarative.RunAction{}
	}
	return declarative.RunAction{
		Script:        r.Script,
		Args:          r.Args,
		Stdin:         r.Stdin,
		Env:           r.Env,
		UpdateTarget:  r.UpdateTarget,
		Confirm:       r.Confirm,
		PromptSecrets: r.PromptSecrets,
	}
}

func toModalRequest(m *ModalRequest) modal.ModalRequest {
	if m == nil {
		return modal.ModalRequest{}
	}
	fields := make([]modal.ModalField, 0, len(m.Fields))
	for _, f := range m.Fields {
		fields = append(fields, modal.ModalField{
			Name:      f.Name,
			Type:      f.Type,
			Sensitive: f.Type == "secret",
		})
	}
	return modal.ModalRequest{
		Type:    m.Type,
		Title:   m.Title,
		Message: m.Message,
		Fields:  fields,
	}
}
