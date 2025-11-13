package util

import (
	"errors"
	"fmt"
	"os"
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
func Handle(err error) {
	if err == nil {
		return
	}

	exitCode := ExitError

	if errors.Is(err, ErrAborted) {
		exitCode = ExitCancelled
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	os.Exit(exitCode)
}
