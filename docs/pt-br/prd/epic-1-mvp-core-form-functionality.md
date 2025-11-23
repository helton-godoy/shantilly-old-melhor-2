# Epic 1: MVP - Core Form Functionality

**Expanded Goal:** Deliver the end-to-end functionality for the `form` subcommand, proving the viability of the declarative approach (`YAML -> TUI -> JSON`). This includes the initial Go project setup, the CLI structure with `cobra`, parsing the YAML input, rendering the TUI with `bubbletea` and `huh`, basic keyboard interaction, collecting submitted data, and formatting the JSON output for integration with shell scripts.

## Story 1.1 CLI Foundation and Project Structure

**As a** shantilly developer,
**I want** to set up the initial Go project structure (monorepo) and the CLI using `cobra`, with a basic `form` subcommand,
**So that** the application foundation is ready for subsequent features.

### Acceptance Criteria

1. AC1: The Go repository MUST be initialized with `cmd/shantilly` and `internal/` directory structures.
2. AC2: The CLI MUST be implemented using `spf13/cobra`.
3. AC3: A root `shantilly` command and a `form` subcommand MUST exist.
4. AC4: Running `shantilly form` MUST read data from `stdin` (no specific parsing yet).
5. AC5: Running `shantilly form` MUST print a placeholder message to `stdout` (e.g., "Form executed").
6. AC6: The project MUST include a basic `Makefile` for compilation (`go build`).

## Story 1.2 YAML Configuration Parsing

**As a** shantilly developer,
**I want** the `form` subcommand to parse the YAML received via `stdin` using `gopkg.in/yaml.v3` into an internal Go structure,
**So that** the TUI definition can be processed by the application.
**Prerequisite:** Story 1.1

### Acceptance Criteria

1. AC1: A Go struct (`internal/config` or similar) representing the expected YAML structure for a form (as defined in "Expected YAML Structure") MUST be defined.
2. AC2: The `form` subcommand MUST use `gopkg.in/yaml.v3` to unmarshal `stdin` into the defined Go struct.
3. AC3: If YAML parsing fails (invalid YAML), the application MUST print a clear error message to `stderr` and exit with an error status (non-zero) (NFR8).
4. AC4: If parsing succeeds, the application MUST (temporarily for testing) print a representation of the parsed Go struct to `stdout`.

## Story 1.3 Basic TUI Structure with Bubbletea

**As a** shantilly developer,
**I want** to integrate `bubbletea` and `lipgloss` to create a basic TUI model that can be launched by the `form` command, display a message, and be exited,
**So that** the TUI foundation is established.
**Prerequisite:** Story 1.1

### Acceptance Criteria

1. AC1: The project MUST include `charmbracelet/bubbletea` and `charmbracelet/lipgloss` dependencies.
2. AC2: A basic `bubbletea` model (`internal/tui` or similar) with `Init`, `Update`, `View` methods MUST be created.
3. AC3: The `form` subcommand MUST launch the `bubbletea` application with the basic model.
4. AC4: The `View` method MUST render a simple placeholder message (e.g., "Shantilly TUI") using `lipgloss` for basic styling.
5. AC5: The user MUST be able to exit the TUI application by pressing `Ctrl+C` or `Esc` (NFR8).

## Story 1.4 Form Rendering with `huh`

**As a** shantilly developer,
**I want** to map the Go structure (parsed from YAML) to `charmbracelet/huh` components and render the corresponding TUI form,
**So that** the user-defined interface is displayed.
**Prerequisites:** Story 1.2, Story 1.3

### Acceptance Criteria

1. AC1: The project MUST include the `charmbracelet/huh` dependency.
2. AC2: The TUI logic (`internal/tui`) MUST be able to receive the parsed Go struct (from Story 1.2).
3. AC3: The TUI logic MUST iterate over the field definitions and create instances of the corresponding `huh` components (Input, Textarea, Select, MultiSelect, Confirm, Note), using `key` and `label` from the YAML.
4. AC4: The `bubbletea` `View` method MUST use `huh.Form` to render the created components.
5. AC5: A simple YAML form (using the defined structure) with 2-3 fields passed via `stdin` MUST be rendered correctly in the TUI.
6. AC6: The `title` (if provided in the YAML) MUST be displayed above the form.

## Story 1.5 Keyboard Navigation and Interaction

**As a** shantilly user,
**I want** to be able to navigate between form fields and interact with them using only the keyboard,
**So that** I can fill out the form efficiently.
**Prerequisite:** Story 1.4

### Acceptance Criteria

1. AC1: The `Tab` key MUST move focus to the next form field.
2. AC2: `Shift+Tab` MUST move focus to the previous field.
3. AC3: `Enter` MUST confirm selection in fields like `Select` and `Confirm`, or move to the next field in `Input`/`Textarea` (`huh` default behavior).
4. AC4: `Space` MUST toggle selection in `MultiSelect` fields.
5. AC5: `Arrow` keys (Up/Down) MUST allow navigation within options of `Select` and `MultiSelect` fields.
6. AC6: The currently focused field MUST be visually highlighted.

## Story 1.6 Submission and JSON Output

**As a** script developer,
**I want** shantilly to collect the data entered in the TUI form and print it in JSON format to `stdout` upon submission,
**So that** my script can easily consume the results.
**Prerequisite:** Story 1.5

### Acceptance Criteria

1. AC1: The `huh` library MUST be configured to allow form submission (typically after the last field or via an implicit button).
2. AC2: After submission, the `bubbletea` application MUST terminate.
3. AC3: The values entered in each form field MUST be collected.
4. AC4: The collected data MUST be formatted as a JSON object, where keys are the `key` values defined in the YAML for each field.
5. AC5: The resulting JSON object MUST be printed to `stdout`.
6. AC6: The TUI MUST be cleared from the screen before printing the JSON.

## Story 1.7 Basic Styling and Alignment

**As a** shantilly user,
**I want** the TUI form to be presented centered with adequate spacing, and to adapt to basic terminal resizing,
**So that** the interface is visually pleasant and functional.
**Prerequisite:** Story 1.4

### Acceptance Criteria

1. AC1: The main container for the `huh` form MUST be styled using `lipgloss`.
2. AC2: The form MUST be rendered horizontally centered within the terminal window if there is sufficient space.
3. AC3: Consistent padding MUST be applied around the form.
4. AC4: The TUI MUST handle resize messages (`tea.WindowSizeMsg`) and re-render the basic layout (maintaining centering if possible).
5. AC5: If the terminal is resized to a very small width, the form MUST still be minimally usable (`huh` might truncate or wrap lines by default).

## Story 1.8 Static Build and Distribution

**As a** shantilly developer,
**I want** a build process that generates static, cross-compiled binaries for the target platforms,
**So that** the application can be easily distributed.
**Prerequisite:** Story 1.1

### Acceptance Criteria

1. AC1: The `Makefile` (or build script) MUST include targets to compile the `shantilly` binary.
2. AC2: The build process MUST use flags (`CGO_ENABLED=0`, `ldflags="-s -w"`) to ensure a static and optimized binary.
3. AC3: Targets (or a script) MUST exist to cross-compile the binary for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64).
4. AC4: The generated binaries MUST be executable on their respective platforms (basic verification).
