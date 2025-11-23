# Error Handling Strategy

## General Approach

* **Model:** Standard Go `error` interface.
* **Propagation:** Errors propagated up the call stack to `CLI`.
* **Centralized Handling:** `CLI` catches errors, passes to `ErrorHandler`.
* **Output Separation:** `stdout` for success JSON (FR6), `stderr` for errors/status (NFR8).
* **Exit Codes:** Non-zero exit codes for errors and cancellations defined in `internal/util`.

## Error Handling Patterns

**YAML Parse Errors (`ConfigParser`)**

* **Detection:** `gopkg.in/yaml.v3` errors.
* **Output:** Formatted message (incl. line number if possible) to `stderr`.
* **Exit Code:** `util.ExitError` (e.g., 1).
* **Workflow:** See "Workflow 2: YAML Parse Error".

**TUI Cancellation (`TUIEngine`)**

* **Detection:** `bubbletea.QuitMsg` or specific `util.ErrAborted`.
* **Output:** No `stdout`. Optional "Cancelled" message to `stderr`.
* **Exit Code:** `util.ExitCancelled` (e.g., 2).
* **Workflow:** See "Workflow 3: User Cancellation".

**Unexpected TUI Errors (`TUIEngine`)**

* **Detection:** Internal `bubbletea`/`huh` errors.
* **Output:** Detailed error message (maybe stack trace) to `stderr`.
* **Exit Code:** `util.ExitError` (e.g., 1).

**Other Errors (e.g., JSON Encoding)**

* **Detection:** Standard library errors.
* **Output:** Error message to `stderr`.
* **Exit Code:** `util.ExitError` (e.g., 1).

## Implementation (`ErrorHandler`)

Located in `internal/util/errorhandler.go`. Provides `Handle(err error)`.

```go
package util

import (
 "errors" // Import errors package
 "fmt"
 "os"

 tea "[github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)" // Import bubbletea
)

// Exit codes for clarity
const (
 ExitOK        = 0
 ExitError     = 1 // General error
 ExitCancelled = 2 // User cancellation
)

// Define a specific error for cancellation if bubbletea.QuitMsg isn't enough context
var ErrAborted = errors.New("operation aborted by user")

// Handle prints the error (if not nil) to stderr and exits with the code.
// It assumes successful exit (code 0) if err is nil.
func Handle(err error) {
 if err == nil {
  os.Exit (apenas permitido em cmd/shantilly; é proibido em internal/runtime/** conforme governance-runtime-tui-v2.0)(ExitOK) // Success path
 }

 exitCode := ExitError // Default to general error

 // Check for specific error types or messages
 // Check if it's our specific abort error OR if it's the standard bubbletea Quit message
 if errors.Is(err, ErrAborted) || errors.As(err, &tea.QuitMsg{}) {
  exitCode = ExitCancelled
  // Optionally print a message to stderr for cancellation
  // fmt.Fprintln(os.Stderr, "Operation cancelled.") // Keep it silent for better script integration
 } else {
  // Print actual error for non-cancellation cases
  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
  // TODO: Add more sophisticated error type checking and formatting here
  // e.g., check for YAML parsing errors and provide line numbers if possible from yaml.v3.
 }

 os.Exit (apenas permitido em cmd/shantilly; é proibido em internal/runtime/** conforme governance-runtime-tui-v2.0)(exitCode)
}
```
