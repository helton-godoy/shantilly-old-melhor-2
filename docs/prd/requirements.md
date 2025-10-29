# Requirements

## Functional

* **FR1:** The CLI MUST implement a `form` subcommand (based on `spf13/cobra`) [cite: Project Brief_ shantilly.md].
* **FR2:** The `form` subcommand MUST be able to read a form definition in YAML format from `stdin` (allowing *here-docs* or *pipes* in the shell) [cite: Project Brief_ shantilly.md].
* **FR3:** The YAML definition MUST support the following form components, mapped to the `charmbracelet/huh` library: Input, Textarea, Select, MultiSelect, Confirm, and Note [cite: Project Brief_ shantilly.md].
* **FR4:** shantilly MUST render an interactive TUI (based on `charmbracelet/bubbletea`) corresponding to the received YAML definition [cite: Project Brief_ shantilly.md].
* **FR5:** The user MUST be able to navigate and fill the TUI form using the keyboard [cite: Project Brief_ shantilly.md].
* **FR6:** Upon successful form submission, the application MUST print the collected data in JSON format to `stdout` [cite: Project Brief_ shantilly.md].
* **FR7:** The TUI MUST include basic styling and responsiveness to terminal size changes (using `charmbracelet/lipgloss`) [cite: Project Brief_ shantilly.md].

## Non-Functional

* **NFR1:** The application MUST be distributed as a single statically compiled executable binary [cite: Project Brief_ shantilly.md].
* **NFR2:** The binary MUST be compatible with the following architectures: Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64) [cite: Project Brief_ shantilly.md].
* **NFR3:** The application MUST be written in Go (latest stable version) [cite: Project Brief_ shantilly.md].
* **NFR4:** The TUI startup time (from command execution to rendering) MUST be fast (target < 500ms) [cite: Project Brief_ shantilly.md].
* **NFR5:** The application MUST NOT depend on external tools like `dialog` or `whiptail` [cite: Project Brief_ shantilly.md].
* **NFR6:** Mouse support is explicitly out of scope for the MVP (MUST NOT be implemented) [cite: Project Brief_ shantilly.md].
* **NFR7:** Support for complex layouts (multi-panel), advanced menus, modal dialogs, and SSH are out of scope for the MVP [cite: Project Brief_ shantilly.md].
* **NFR8:** The application MUST handle errors gracefully. Specifically:
  * YAML parsing errors (FR2) MUST result in a clear error message on `stderr` (indicating the problem, e.g., line/column if possible) and the application MUST exit with an error status (non-zero).
  * The TUI MUST allow the user to exit at any time (e.g., `Ctrl+C`, `Esc`) without corrupting the terminal.
  * Unexpected errors during TUI execution MUST be caught and ideally reported on `stderr` before exiting.
