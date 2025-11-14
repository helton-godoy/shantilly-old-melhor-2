package runner

import (
	tea "github.com/charmbracelet/bubbletea"

	"shantilly/pkg/declarative"
	"shantilly/pkg/tui"
)

// HandleRunRequest converte um tui.RunScriptRequestMsg (baseado em declarative.OnHandler)
// em um RunAction local e dispara o ScriptRunner via tea.Cmd.
//
// Nesta primeira integração, aproveitamos a implementação síncrona de ScriptRunner.Run
// e simplesmente a executamos em uma goroutine, sem streaming fino de output.
// Futuras waves podem especializar este adaptador para FR11 (streaming linha a linha).
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
		// Execução em goroutine para não bloquear o loop TEA.
		go func() {
			_ = r.Run(converted)
		}()
		return nil
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
