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
		os.Exit(ExitSuccess)
	}

	exitCode := ExitError

	// Check for specific error types
	// Check if it's our specific abort error OR if it's the standard bubbletea Quit message
	if errors.Is(err, ErrAborted) || errors.As(err, &tea.QuitMsg{}) {
		exitCode = ExitCancelled
		// Optionally print a message for cancellation
		// fmt.Fprintln(os.Stderr, "Operation cancelled.")
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	os.Exit(exitCode)
}
