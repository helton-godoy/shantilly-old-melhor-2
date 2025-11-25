package tui

import (
	"io"

	tea "github.com/charmbracelet/bubbletea"

	"shantilly/pkg/declarative"
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

// RunScriptRequestMsg é enviada pelo EventManager
// e tratada pelo ScriptRunner.
type RunScriptRequestMsg struct {
	Handler declarative.OnHandler // Handler completo do bloco on:
}

// ScriptStdoutMsg é usado pelo ScriptRunner para enviar streaming de
// stdout/stderr para o componente alvo (ex: ViewportComponent).
type ScriptStdoutMsg struct {
	TargetID string
	Line     string // Linha de saída (já com newline, se aplicável)
}

// ShowModalMsg é emitida pelo EventManager (para Segurança JIT)
// e tratada pelo MainModel/LayoutManager para exibir um modal.
// Fonte: docs/architecture/core-workflows.md (Workflow 7.3)
type ShowModalMsg struct {
	Content tea.Model // O modelo do modal (ex: modal.Confirm)
}

// ModalResultMsg é emitida pelo modal (ex: modal.Confirm) quando é fechado.
type ModalResultMsg struct {
	Confirmed bool        // true se 'Sim', false se 'Não'/Abortado
	Payload   interface{} // Payload opcional (ex: segredos)
}

// processPipeMsg é uma mensagem interna usada pelo ScriptRunner para
// comunicar linhas lidas de stdout/stderr para o LayoutManager.
type processPipeMsg struct {
	TargetID string
	Line     string
}

// processFinishedMsg é uma mensagem interna usada pelo ScriptRunner para
// indicar que o processo terminou (com sucesso ou erro).
type processFinishedMsg struct {
	TargetID string
	Err      error
}

// PipeReader é uma interface mínima compatível com io.Reader
// usada nas rotinas de streaming.
type PipeReader interface {
	io.Reader
}
