package event

import (
	"testing"

	"shantilly/internal/runtime/modal"
	"shantilly/pkg/declarative"
)

func TestCoordinator_RunActionConfirm_AbreModalNaoExecuta(t *testing.T) {
	stack := modal.NewStack()
	c := NewCoordinator(stack)

	events, toRun := c.ProcessActions(EngineActions{
		RunActions: []*RunAction{{Script: "/bin/echo", Confirm: true}},
	})
	if len(events) != 1 || events[0].Type != modal.ModalOpenedType {
		t.Fatalf("esperava 1 modal.opened, obtive: %#v", events)
	}
	if len(toRun) != 0 {
		t.Fatalf("não deveria executar antes de confirmar")
	}
	if !stack.HasActive() || !stack.Top().Active {
		t.Fatalf("topo deveria estar ativo e plano de fundo bloqueado")
	}
}

func TestCoordinator_ModalConfirmed_AplicaSegredosEmStdin_ExecutaPendencia(t *testing.T) {
	stack := modal.NewStack()
	c := NewCoordinator(stack)

	_, _ = c.ProcessActions(EngineActions{
		RunActions: []*RunAction{{Script: "/bin/echo", Confirm: true, UpdateTarget: "viewport:out"}},
	})
	top := stack.Top()
	if top == nil {
		t.Fatalf("esperava modal no topo")
	}

	events2, toRun := c.OnEvent(ShantillyEventFromModalConfirm(top.ID, map[string]interface{}{"token": "s3cr3t"}))
	if stack.HasActive() {
		t.Fatalf("após confirmar, não deveria haver modal ativo")
	}
	if len(toRun) != 1 {
		t.Fatalf("esperava 1 execução, obtive %d", len(toRun))
	}
	if toRun[0].Stdin["token"] != "s3cr3t" {
		t.Fatalf("segredos não aplicados corretamente ao RunAction.Stdin")
	}
	if len(events2) != 0 {
		t.Fatalf("coordenador não deveria emitir eventos adicionais aqui")
	}
}

func TestCoordinator_PromptSecrets_CompoeModalSecret_NaoExecuta(t *testing.T) {
	stack := modal.NewStack()
	c := NewCoordinator(stack)

	events, toRun := c.ProcessActions(EngineActions{
		RunActions: []*RunAction{{Script: "/bin/echo", PromptSecrets: []string{"token", "pass"}}},
	})
	if len(events) != 1 || events[0].Type != modal.ModalOpenedType {
		t.Fatalf("esperava 1 modal.opened")
	}
	m := stack.Top()
	if m == nil || m.Request.Type != "prompt" || !m.Sensitive || len(m.Request.Fields) != 2 || m.Request.Fields[0].Type != "secret" {
		t.Fatalf("modal prompt secreto mal formado: %#v", m)
	}
	if len(toRun) != 0 {
		t.Fatalf("não deveria executar antes da confirmação")
	}
}

func TestCoordinator_ModalCancelled_NaoExecutaLimpaPendencia(t *testing.T) {
	stack := modal.NewStack()
	c := NewCoordinator(stack)

	_, _ = c.ProcessActions(EngineActions{
		RunActions: []*RunAction{{Script: "/bin/echo", Confirm: true}},
	})
	top := stack.Top()
	if top == nil {
		t.Fatalf("esperava modal no topo")
	}

	events2, toRun := c.OnEvent(ShantillyEventFromModalCancel(top.ID))
	if len(toRun) != 0 {
		t.Fatalf("não deveria executar após cancelamento")
	}
	if stack.HasActive() {
		t.Fatalf("após cancelamento, não deveria haver modal ativo")
	}
	if len(events2) != 0 {
		t.Fatalf("coordenador não deveria emitir eventos adicionais aqui")
	}
}

func TestCoordinator_MultiplosPendentes_CorrelacaoPorModalID(t *testing.T) {
	stack := modal.NewStack()
	c := NewCoordinator(stack)

	_, _ = c.ProcessActions(EngineActions{RunActions: []*RunAction{{Script: "/bin/a", Confirm: true}}})
	idA := stack.Top().ID
	_, _ = c.ProcessActions(EngineActions{RunActions: []*RunAction{{Script: "/bin/b", Confirm: true}}})
	idB := stack.Top().ID

	_, toRun := c.OnEvent(ShantillyEventFromModalConfirm(idB, nil))
	if len(toRun) != 1 || toRun[0].Script != "/bin/b" {
		t.Fatalf("deveria executar pendência B primeiro")
	}

	_, toRun2 := c.OnEvent(ShantillyEventFromModalConfirm(idA, nil))
	if len(toRun2) != 1 || toRun2[0].Script != "/bin/a" {
		t.Fatalf("deveria executar pendência A depois")
	}
}

func ShantillyEventFromModalConfirm(id string, vals map[string]interface{}) declarative.ShantillyEvent {
	payload := map[string]interface{}{modal.ModalIDKey: id}
	if vals != nil {
		payload[modal.ModalValuesKey] = vals
	}
	return declarative.ShantillyEvent{
		Type:    modal.ModalConfirmedType,
		Payload: payload,
	}
}

func ShantillyEventFromModalCancel(id string) declarative.ShantillyEvent {
	return declarative.ShantillyEvent{
		Type:    modal.ModalCancelledType,
		Payload: map[string]interface{}{modal.ModalIDKey: id},
	}
}
