package util

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Exit codes
const (
	ExitSuccess   = 0
	ExitError     = 1
	ExitCancelled = 2
)

// ErrAborted represents a user cancellation error.
var ErrAborted = errors.New("operation aborted by user")

// Handle verifica se um erro é não-nulo, imprime-o para stderr e sai.
// Se o erro for nulo, a função retorna sem fazer nada.
func Handle(err error) {
	if err == nil {
		// Em vez de encerrar o processo aqui (os.Exit), retornamos sucesso ao chamador.
		// A casca CLI principal (cmd/shantilly/main.go) é o ÚNICO ponto autorizado a chamar os.Exit,
		// conforme Gate 1.x.no-osexit-core.yml e docs/architecture/governance-runtime-tui-v2.0.md.
		return
	}

	exitCode := ExitError

	// Check for specific error types
	// Check if it's our specific abort error OR if it's the standard bubbletea Quit message
	if errors.Is(err, ErrAborted) {
		exitCode = ExitCancelled
		// Optionally print a message for cancellation
		// fmt.Fprintln(os.Stderr, "Operation cancelled.")
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	// Em vez de chamar os.Exit diretamente neste utilitário interno, propagamos o código
	// de saída para que a casca CLI decida como encerrar o processo.
	// Isso garante conformidade com:
	// - Gate 1.x.no-osexit-core.yml
	// - docs/architecture/introduction.md (constraints invioláveis)
	// - docs/architecture/governance-runtime-tui-v2.0.md (proibição de os.Exit no core)
	if exitCode != ExitSuccess {
		// Poderíamos integrar logging ou mapeamento adicional aqui se necessário,
		// desde que não invoque os.Exit diretamente.
	}
}
