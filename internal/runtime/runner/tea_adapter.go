package runner

import (
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/pkg/declarative"
	"shantilly/pkg/tui"
)

// bufferingSink é um EventSink que opcionalmente delega para um sink interno
// e, ao mesmo tempo, acumula as atualizações de UpdateTargetUpdate em memória.
//
// Ele é usado apenas neste adaptador TEA para transformar o output agregado
// do ScriptRunner em uma única tui.ScriptStdoutMsg consumível pelo viewport.
type bufferingSink struct {
	inner   EventSink
	updates []UpdateTargetUpdate
}

func (s *bufferingSink) EmitEvent(ev ShantillyEvent) {
	if s.inner != nil {
		s.inner.EmitEvent(ev)
	}
}

func (s *bufferingSink) EmitUpdate(update UpdateTargetUpdate) {
	s.updates = append(s.updates, update)
	if s.inner != nil {
		s.inner.EmitUpdate(update)
	}
}

// HandleRunRequest converte um tui.RunScriptRequestMsg (baseado em declarative.OnHandler)
// em um RunAction local e dispara o ScriptRunner via tea.Cmd.
//
// Nesta integração inicial com viewport, usamos um bufferingSink para capturar
// as atualizações de UpdateTargetUpdate emitidas pelo ScriptRunner e, ao final,
// retornamos uma única tui.ScriptStdoutMsg com o conteúdo agregado.
func HandleRunRequest(r *ScriptRunner, msg tui.RunScriptRequestMsg) tea.Cmd {
	if r == nil {
		return nil
	}

	handler := msg.Handler
	if handler.Run == nil {
		return nil
	}

	// Converte declarative.RunAction -> runner.RunAction.
	converted := convertRunAction(*handler.Run)

	return func() tea.Msg {
		// Captura o sink atual (se houver) e envolve em um bufferingSink.
		origSink := r.sink
		bufSink := &bufferingSink{inner: origSink}
		r.sink = bufSink
		defer func() {
			// Garante restauração do sink original após a execução.
			r.sink = origSink
		}()

		// Execução síncrona dentro do Cmd: o Bubble Tea já roda o Cmd em goroutine.
		_ = r.Run(converted)

		// Se não houver UpdateTarget definido, não há viewport para atualizar.
		if converted.UpdateTarget == "" {
			return nil
		}

		// Agrega o conteúdo de todas as atualizações em uma única string.
		var content string
		for _, u := range bufSink.updates {
			if u.Content == "" {
				continue
			}
			if content != "" {
				content += "\n"
			}
			content += u.Content
		}

		if content == "" {
			return nil
		}

		// Envia uma única mensagem ScriptStdoutMsg para o target declarado.
		return tui.ScriptStdoutMsg{
			TargetID: converted.UpdateTarget,
			Line:     content,
		}
	}
}

// convertRunAction faz o mapeamento entre o modelo declarativo normativo e o
// modelo interno usado pelo ScriptRunner na wave atual.
func convertRunAction(src declarative.RunAction) RunAction {
	return RunAction{
		Script:        src.Script,
		Args:          src.Args,
		Stdin:         src.Stdin,
		Env:           src.Env,
		UpdateTarget:  src.UpdateTarget,
		Confirm:       src.Confirm,
		PromptSecrets: src.PromptSecrets,
	}
}
