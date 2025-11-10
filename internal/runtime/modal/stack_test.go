package modal

import "testing"

// E1.5 — Modal Stack + Security
// Testes focados em:
// - Pilha explícita (push/pop/top).
// - Foco exclusivo no topo.
// - Bloqueio lógico do fundo via HasActive().
// - Ausência de efeitos colaterais externos (sem logging, sem integração direta).
// - Eventos `modal.*` gerados apenas via helpers deste pacote.

func TestModalStack_PushPopTop_HasActive(t *testing.T) {
	stack := NewStack()

	if stack.HasActive() {
		t.Fatalf("expected no active modal on new stack")
	}
	if stack.Top() != nil {
		t.Fatalf("expected nil Top on new stack")
	}
	if stack.Size() != 0 {
		t.Fatalf("expected size 0 on new stack, got %d", stack.Size())
	}

	// Push primeiro modal
	m1Req := ModalRequest{
		Type:  "confirm",
		Title: "Confirmar ação",
	}
	m1 := stack.Push(m1Req)
	if !stack.HasActive() {
		t.Fatalf("expected HasActive() after first push")
	}
	if stack.Size() != 1 {
		t.Fatalf("expected size 1 after first push, got %d", stack.Size())
	}
	top := stack.Top()
	if top == nil || top.ID != m1.ID || !top.Active {
		t.Fatalf("expected top to be first modal active, got %+v", top)
	}

	// Push segundo modal — deve assumir foco exclusivo.
	m2Req := ModalRequest{
		Type:  "prompt",
		Title: "Informe dado sensível",
		Fields: []ModalField{
			{Name: "secret", Type: "secret", Sensitive: true},
		},
		Sensitive: true,
	}
	m2 := stack.Push(m2Req)
	if stack.Size() != 2 {
		t.Fatalf("expected size 2 after second push, got %d", stack.Size())
	}
	top = stack.Top()
	if top == nil || top.ID != m2.ID || !top.Active {
		t.Fatalf("expected top to be second modal active, got %+v", top)
	}
	// Primeiro modal não deve estar ativo.
	if m1After := stack.modals[0]; m1After.Active {
		t.Fatalf("expected first modal to be inactive after second push")
	}

	// Pop topo — deve remover m2 e reativar m1.
	popped, ok := stack.Pop()
	if !ok {
		t.Fatalf("expected pop to succeed")
	}
	if popped.ID != m2.ID {
		t.Fatalf("expected popped to be second modal, got %+v", popped)
	}
	if stack.Size() != 1 {
		t.Fatalf("expected size 1 after pop, got %d", stack.Size())
	}
	top = stack.Top()
	if top == nil || top.ID != m1.ID || !top.Active {
		t.Fatalf("expected top to revert to first modal active, got %+v", top)
	}

	// Pop final — pilha vazia, sem ativo.
	_, ok = stack.Pop()
	if !ok {
		t.Fatalf("expected second pop to succeed")
	}
	if stack.HasActive() {
		t.Fatalf("expected no active modal after popping all")
	}
	if stack.Top() != nil {
		t.Fatalf("expected nil Top after popping all")
	}
	if stack.Size() != 0 {
		t.Fatalf("expected size 0 after popping all, got %d", stack.Size())
	}
}

func TestModalStack_SensitiveFlagPropagation(t *testing.T) {
	stack := NewStack()

	req := ModalRequest{
		Type:      "prompt",
		Title:     "Segredo",
		Sensitive: false,
		Fields: []ModalField{
			{Name: "vault_pass", Type: "secret", Sensitive: true},
		},
	}

	m := stack.Push(req)
	if !m.Sensitive {
		t.Fatalf("expected modal to be marked Sensitive due to secret field")
	}

	top := stack.Top()
	if top == nil || !top.Sensitive {
		t.Fatalf("expected top modal Sensitive flag true, got %+v", top)
	}
}

func TestHandlers_OpenModalFromRequest_EmitsOpenedEvent(t *testing.T) {
	stack := NewStack()

	req := ModalRequest{
		ID:    "m-open",
		Type:  "confirm",
		Title: "Confirmar",
	}
	m, ev := OpenModalFromRequest(stack, req)

	if m.ID != "m-open" {
		t.Fatalf("expected modal ID 'm-open', got %s", m.ID)
	}
	if !m.Active {
		t.Fatalf("expected opened modal to be active")
	}
	if stack.Top() == nil || stack.Top().ID != m.ID {
		t.Fatalf("expected opened modal to be on top")
	}

	if ev == nil {
		t.Fatalf("expected modal.opened event")
	}
	if ev.Type != ModalOpenedType {
		t.Fatalf("expected event type %s, got %s", ModalOpenedType, ev.Type)
	}
	if ev.Payload[ModalIDKey] != m.ID {
		t.Fatalf("expected payload[%s] = %s", ModalIDKey, m.ID)
	}
}

func TestHandlers_ConfirmModal_PopAndEmit(t *testing.T) {
	stack := NewStack()

	req := ModalRequest{
		ID:   "m-confirm",
		Type: "confirm",
	}
	stack.Push(req)

	values := map[string]interface{}{"ok": true}

	ev := ConfirmModal(stack, "m-confirm", values)
	if ev == nil {
		t.Fatalf("expected modal.confirmed event")
	}
	if ev.Type != ModalConfirmedType {
		t.Fatalf("expected event type %s, got %s", ModalConfirmedType, ev.Type)
	}
	if ev.Payload[ModalIDKey] != "m-confirm" {
		t.Fatalf("expected modal_id to be 'm-confirm'")
	}
	if !ev.Payload[ModalConfirmedKey].(bool) {
		t.Fatalf("expected confirmed=true")
	}
	if v, ok := ev.Payload[ModalValuesKey]; !ok || v == nil {
		t.Fatalf("expected values in payload")
	}

	// Após confirmação, pilha vazia.
	if stack.HasActive() || stack.Size() != 0 {
		t.Fatalf("expected empty stack after confirm")
	}
}

func TestHandlers_CancelModal_PopAndEmit(t *testing.T) {
	stack := NewStack()

	req := ModalRequest{
		ID:   "m-cancel",
		Type: "confirm",
	}
	stack.Push(req)

	ev := CancelModal(stack, "m-cancel")
	if ev == nil {
		t.Fatalf("expected modal.cancelled event")
	}
	if ev.Type != ModalCancelledType {
		t.Fatalf("expected event type %s, got %s", ModalCancelledType, ev.Type)
	}
	if ev.Payload[ModalIDKey] != "m-cancel" {
		t.Fatalf("expected modal_id 'm-cancel'")
	}
	if ev.Payload[ModalConfirmedKey].(bool) {
		t.Fatalf("expected confirmed=false on cancel")
	}

	if stack.HasActive() || stack.Size() != 0 {
		t.Fatalf("expected empty stack after cancel")
	}
}

func TestHandlers_InputModal_TopOnly_NoSideEffects(t *testing.T) {
	stack := NewStack()

	// Dois modais empilhados para garantir roteamento apenas ao topo.
	stack.Push(ModalRequest{ID: "m1", Type: "confirm"})
	stack.Push(ModalRequest{ID: "m2", Type: "prompt"})

	// Input errado (para ID que não é topo) deve retornar nil.
	if ev := InputModal(stack, "m1", map[string]interface{}{"k": "v"}); ev != nil {
		t.Fatalf("expected nil event when modalID is not top")
	}

	// Input correto (para topo) gera modal.input.
	values := map[string]interface{}{"field": "value"}
	ev := InputModal(stack, "m2", values)
	if ev == nil {
		t.Fatalf("expected modal.input event for top modal")
	}
	if ev.Type != ModalInputType {
		t.Fatalf("expected event type %s, got %s", ModalInputType, ev.Type)
	}
	if ev.Payload[ModalIDKey] != "m2" {
		t.Fatalf("expected modal_id 'm2'")
	}
	if got, ok := ev.Payload[ModalValuesKey].(map[string]interface{}); !ok || got["field"] != "value" {
		t.Fatalf("expected values payload preserved, got %#v", ev.Payload[ModalValuesKey])
	}

	// Input não deve alterar tamanho da pilha nem foco.
	if stack.Size() != 2 {
		t.Fatalf("expected stack size 2 after input, got %d", stack.Size())
	}
	if top := stack.Top(); top == nil || top.ID != "m2" {
		t.Fatalf("expected top modal to remain 'm2'")
	}
}
