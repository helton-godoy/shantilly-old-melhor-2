# shantilly Product Requirements Document (PRD)

## Goals and Background Context

### Goals

* Simplify the creation of interactive and modern TUI interfaces for shell scripts [cite: shantilly Product Requirements Document (PRD).md].
* Offer a declarative (YAML) and more powerful alternative to `dialog` and `whiptail` [cite: shantilly Product Requirements Document (PRD).md].
* Ensure portability and consistency across Linux, macOS, and Windows via a single static binary [cite: shantilly Product Requirements Document (PRD).md].
* Facilitate integration with shell script pipelines (read `stdin` YAML, write `stdout` JSON) [cite: shantilly Product Requirements Document (PRD).md].
* Achieve broad adoption by the script developer community (Bash, Zsh, PowerShell, etc.) [cite: Project Brief_ shantilly.md].
* Deliver a robust and viable MVP (v1.0), focused on the `form` subcommand, within 3 months [cite: Project Brief_ shantilly.md].
* Enable the creation of significantly richer interfaces (layouts, rich components) than traditional tools [cite: Project Brief_ shantilly.md].

### Background Context

Shell script developers needing complex user interactions are currently underserved. Traditional tools like `dialog` and `whiptail` are functional but severely limited in components, layout control, and aesthetics, lacking mouse support [cite: Project Brief_shantilly.md, shantilly Product Requirements Document (PRD).md]. The alternative, using programmatic TUI libraries (like Bubbletea), requires knowledge of languages like Go or Python, adding disproportionate complexity for script-based automation [cite: Project Brief_ shantilly.md, shantilly Product Requirements Document (PRD).md].

shantilly fills this gap. It's a portable CLI tool (single static binary) enabling the *declarative creation* of rich, modern TUIs [cite: Project Brief_shantilly.md]. The user defines the interface (e.g., a form) in YAML, passes it to `shantilly` via `stdin`, and the tool renders the interactive TUI [cite: Project Brief_ shantilly.md]. Collected data is then returned as structured JSON on `stdout`, allowing easy integration into script pipelines [cite: Project Brief_shantilly.md]. This PRD focuses on the MVP (Minimum Viable Product) scope to validate the core functionality [cite: Project Brief_ shantilly.md].

### Change Log

| Date       | Version | Description                                                                                                    | Author    |
|:-----------|:--------|:---------------------------------------------------------------------------------------------------------------|:----------|
| 2025-10-22 | 0.1.0   | Initial PRD draft based on Project Brief.                                                                      | John (PM) |
| 2025-10-23 | 0.1.1   | Added UI/UX section and refined MVP layout.                                                                    | John (PM) |
| 2025-10-23 | 0.1.2   | Added Technical Assumptions section.                                                                           | John (PM) |
| 2025-10-23 | 0.1.3   | Added Epic List (MVP).                                                                                         | John (PM) |
| 2025-10-23 | 0.1.4   | Added Epic 1 Details (MVP) with Stories.                                                                       | John (PM) |
| 2025-10-23 | 0.1.5   | Completed PM Checklist and Next Steps section.                                                                 | John (PM) |
| 2025-10-23 | 0.2.0   | Implemented PM Checklist recommendations (YAML Structure and Error Handling - NFR8). Updated Architect prompt. | John (PM) |
| 2025-10-27 | 0.3.0   | Added preventive epics (3-7) for future roadmap planning and process improvement. | Sarah (PO) |

## Requirements

### Functional

* **FR1:** The CLI MUST implement a `form` subcommand (based on `spf13/cobra`) [cite: Project Brief_ shantilly.md].
* **FR2:** The `form` subcommand MUST be able to read a form definition in YAML format from `stdin` (allowing *here-docs* or *pipes* in the shell) [cite: Project Brief_ shantilly.md].
* **FR3:** The YAML definition MUST support the following form components, mapped to the `charmbracelet/huh` library: Input, Textarea, Select, MultiSelect, Confirm, and Note [cite: Project Brief_ shantilly.md].
* **FR4:** shantilly MUST render an interactive TUI (based on `charmbracelet/bubbletea`) corresponding to the received YAML definition [cite: Project Brief_ shantilly.md].
* **FR5:** The user MUST be able to navigate and fill the TUI form using the keyboard [cite: Project Brief_ shantilly.md].
* **FR6:** Upon successful form submission, the application MUST print the collected data in JSON format to `stdout` [cite: Project Brief_ shantilly.md].
* **FR7:** The TUI MUST include basic styling and responsiveness to terminal size changes (using `charmbracelet/lipgloss`) [cite: Project Brief_ shantilly.md].

### Non-Functional

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

## User Interface Design Goals

### Overall UX (User Experience) Vision

The user experience should be clean, intuitive, and focused on efficient keyboard data entry. The aesthetics should follow the modern, minimalist standard popularized by the Charmbracelet ecosystem [cite: Project Brief_shantilly.md, Relatório Técnico do Brainstorming Shantilly.md], serving as a direct upgrade from the dated appearance of `dialog` and `whiptail` [cite: Project Brief_ shantilly.md]. Responsiveness to terminal resizing is crucial [cite: docs/prd.md]. **Important: For the MVP, the form layout will be linear (one question below the other), following the `huh` library's standard.**

### Key Interaction Paradigms

* **Keyboard Focus:** Navigation MUST be primarily keyboard-based, following standard form patterns (e.g., `Tab` / `Shift+Tab` to navigate fields, `Enter` to submit or select, `Space` to toggle selections, `Arrows` for lists).
* **Immediate Feedback:** The currently focused component MUST be clearly highlighted.
* **Clean Exit:** The user MUST be able to exit the form at any time (e.g., `Ctrl+C` or `Esc`), and submission should clear the screen and return control to the script cleanly.

### Core Screens and Views

For the MVP scope, there is only one main view:

* **Form View (`form view`):** Renders the TUI form based on YAML [cite: docs/prd.md], processed by the `huh` library [cite: docs/prd.md].

### Alignment and Layout (MVP)

Although the MVP (based on `huh`) does not support complex layouts (multiple columns) [cite: docs/prd.md], the form view MUST, whenever possible, be rendered aesthetically (e.g., centered on screen, with adequate padding), using `lipgloss` [cite: docs/prd.md] to manage the overall alignment of the form container.

### Accessibility

* **Standard:** WCAG AA (As applicable to terminal text).
* **Requirements:** The interface MUST ensure sufficient color contrast between text, background, and focus elements, adhering to `lipgloss` standard themes [cite: docs/prd.md].

### Branding

The visual identity will be defined by the standard components of `charmbracelet/huh` and `charmbracelet/lipgloss` [cite: docs/prd.md, Relatório Técnico do Brainstorming Shantilly.md]. There will be no custom branding (logos, etc.) in the MVP.

### Target Device and Platforms

* **Platforms:** Modern terminals on Linux, macOS, and Windows [cite: docs/prd.md].
* **Requirements:** Requires a terminal with adequate support for colors (TrueColor recommended) and UTF-8 characters for correct component rendering [cite: Project Brief_ shantilly.md].

## Technical Assumptions

### Repository Structure: Monorepo

The project will be structured as a simple Go monorepo. This facilitates initial organization and potential future addition of subcommands or related libraries [cite: Project Brief_ shantilly.md, Relatório Técnico do Brainstorming Shantilly.md]. The structure will follow standard Go conventions, with CLI logic in `cmd/shantilly` and TUI/parsing logic in `internal/`.

### Service Architecture: Monolithic CLI Application

For the MVP, shantilly is a monolithic CLI application that executes, processes input, and terminates [cite: Project Brief_ shantilly.md]. The architecture will be modular internally (separating CLI and TUI), but there will be no microservices or network communication [cite: Relatório Técnico do Brainstorming Shantilly.md]. SSH server mode is a Post-MVP vision (NFR7) [cite: docs/prd.md].

### Testing Requirements: Focus on Unit Tests

Given the MVP time constraints and the challenges of testing TUIs [cite: Project Brief_shantilly.md], the primary focus will be on robust unit tests for the YAML parsing logic and internal business logic [cite: Relatório Técnico do Brainstorming Shantilly.md]. UI integration tests will be limited to manual smoke tests on the main platforms [cite: Project Brief_ shantilly.md].

### Expected YAML Structure (`stdin`)

The form definition passed via `stdin` (FR2) MUST follow a basic YAML structure. The Architect will detail the corresponding Go struct, but conceptually, the YAML must allow at least:

```yaml
# Optional: Title to be displayed above the form
title: "Example Shantilly Form"

# List of form fields
fields:
  - key: "username"           # Key used in the JSON output (FR6)
    label: "Username:"        # Text displayed in the TUI
    type: "input"             # Field type (mapped to huh - FR3)
    # placeholder: "Enter your name" # Optional for inputs
    # value: "Default Value"   # Optional

  - key: "description"
    label: "Description:"
    type: "textarea"

  - key: "server"
    label: "Select Server:"
    type: "select"
    options:                  # Required for select/multiselect
      - "Production"
      - "Staging"
      - "Development"

  - key: "modules"
    label: "Modules to Install:"
    type: "multiselect"
    options:
      - "Web Server"
      - "Database"
      - "Cache"
    # limit: 2                  # Optional for multiselect

  - key: "confirm_install"
    label: "Proceed with installation?"
    type: "confirm"
    # affirmative: "Yes"       # Optional
    # negative: "No"        # Optional

  - key: "warning_note"
    label: "Important Note:"    # 'label' here acts as the note title
    type: "note"
    # title: "Attention:"       # Alternative to 'label' for Note
    # detail: "Details..." # Body of the note (if not using 'label')

# Add more global configuration options here in the future (Post-MVP)
```

**Note:** This is a *minimum* structure for the MVP. The Architect may refine key names and add optional fields (like `placeholder`, `value`, `limit`) as needed during architectural design. Validation of the YAML structure after parsing is implied in Story 1.2/1.4.

### Additional Technical Assumptions and Requests

* **Language:** Go (NFR3) [cite: docs/prd.md].
* **Core Frameworks/Libraries:** `spf13/cobra` (FR1) [cite: docs/prd.md], `charmbracelet/bubbletea` (FR4) [cite: docs/prd.md], `charmbracelet/lipgloss` (FR7) [cite: docs/prd.md], `charmbracelet/huh` (FR3) [cite: docs/prd.md], `gopkg.in/yaml.v3` [cite: Project Brief\_ shantilly.md].
* **Build & Distribution:** Standard Go build process with flags for static compilation (`CGO_ENABLED=0`) and cross-compilation for NFR2 architectures [cite: docs/prd.md].

## Epic List

* **Epic 1: MVP - Core Form Functionality**
  * **Goal:** Establish the CLI structure, implement YAML parsing from stdin, render the linear TUI form using `huh`, allow keyboard navigation and submission, and return the collected data as JSON on stdout, delivering the core MVP functionality.

* **Epic 2: Advanced Form Features & User Experience**
  * **Goal:** Expand form capabilities with advanced field types (numeric, date, file), validation, and improved user experience through better error handling and feedback.

* **Epic 3: Advanced Layouts & Multi-Panel Forms**
  * **Goal:** Support complex layouts with multiple panels, form sections, and progressive disclosure for sophisticated form experiences.

* **Epic 4: SSH Server Mode & Remote Forms**
  * **Goal:** Implement SSH server mode for remote form execution, enabling administration and distributed system integration.

* **Epic 5: Advanced Components & Interactions**
  * **Goal:** Add advanced UI components like modals, enhanced selection controls, progress indicators, and sophisticated keyboard navigation.

* **Epic 6: Form Templates & Reusability**
  * **Goal:** Create template system and component library for rapid form development and consistency across applications.

* **Epic 7: Internationalization & Localization**
  * **Goal:** Implement full internationalization support for global adoption with multiple language support and cultural adaptation.

## Epic 1: MVP - Core Form Functionality

**Expanded Goal:** Deliver the end-to-end functionality for the `form` subcommand, proving the viability of the declarative approach (`YAML -> TUI -> JSON`). This includes the initial Go project setup, the CLI structure with `cobra`, parsing the YAML input, rendering the TUI with `bubbletea` and `huh`, basic keyboard interaction, collecting submitted data, and formatting the JSON output for integration with shell scripts.

### Story 1.1 CLI Foundation and Project Structure

**As a** shantilly developer,
**I want** to set up the initial Go project structure (monorepo) and the CLI using `cobra`, with a basic `form` subcommand,
**So that** the application foundation is ready for subsequent features.

#### Acceptance Criteria

1. AC1: The Go repository MUST be initialized with `cmd/shantilly` and `internal/` directory structures.
2. AC2: The CLI MUST be implemented using `spf13/cobra`.
3. AC3: A root `shantilly` command and a `form` subcommand MUST exist.
4. AC4: Running `shantilly form` MUST read data from `stdin` (no specific parsing yet).
5. AC5: Running `shantilly form` MUST print a placeholder message to `stdout` (e.g., "Form executed").
6. AC6: The project MUST include a basic `Makefile` for compilation (`go build`).

### Story 1.2 YAML Configuration Parsing

**As a** shantilly developer,
**I want** the `form` subcommand to parse the YAML received via `stdin` using `gopkg.in/yaml.v3` into an internal Go structure,
**So that** the TUI definition can be processed by the application.
**Prerequisite:** Story 1.1

#### Acceptance Criteria

1. AC1: A Go struct (`internal/config` or similar) representing the expected YAML structure for a form (as defined in "Expected YAML Structure") MUST be defined.
2. AC2: The `form` subcommand MUST use `gopkg.in/yaml.v3` to unmarshal `stdin` into the defined Go struct.
3. AC3: If YAML parsing fails (invalid YAML), the application MUST print a clear error message to `stderr` and exit with an error status (non-zero) (NFR8).
4. AC4: If parsing succeeds, the application MUST (temporarily for testing) print a representation of the parsed Go struct to `stdout`.

### Story 1.3 Basic TUI Structure with Bubbletea

**As a** shantilly developer,
**I want** to integrate `bubbletea` and `lipgloss` to create a basic TUI model that can be launched by the `form` command, display a message, and be exited,
**So that** the TUI foundation is established.
**Prerequisite:** Story 1.1

#### Acceptance Criteria

1. AC1: The project MUST include `charmbracelet/bubbletea` and `charmbracelet/lipgloss` dependencies.
2. AC2: A basic `bubbletea` model (`internal/tui` or similar) with `Init`, `Update`, `View` methods MUST be created.
3. AC3: The `form` subcommand MUST launch the `bubbletea` application with the basic model.
4. AC4: The `View` method MUST render a simple placeholder message (e.g., "Shantilly TUI") using `lipgloss` for basic styling.
5. AC5: The user MUST be able to exit the TUI application by pressing `Ctrl+C` or `Esc` (NFR8).

### Story 1.4 Form Rendering with `huh`

**As a** shantilly developer,
**I want** to map the Go structure (parsed from YAML) to `charmbracelet/huh` components and render the corresponding TUI form,
**So that** the user-defined interface is displayed.
**Prerequisites:** Story 1.2, Story 1.3

#### Acceptance Criteria

1. AC1: The project MUST include the `charmbracelet/huh` dependency.
2. AC2: The TUI logic (`internal/tui`) MUST be able to receive the parsed Go struct (from Story 1.2).
3. AC3: The TUI logic MUST iterate over the field definitions and create instances of the corresponding `huh` components (Input, Textarea, Select, MultiSelect, Confirm, Note), using `key` and `label` from the YAML.
4. AC4: The `bubbletea` `View` method MUST use `huh.Form` to render the created components.
5. AC5: A simple YAML form (using the defined structure) with 2-3 fields passed via `stdin` MUST be rendered correctly in the TUI.
6. AC6: The `title` (if provided in the YAML) MUST be displayed above the form.

### Story 1.5 Keyboard Navigation and Interaction

**As a** shantilly user,
**I want** to be able to navigate between form fields and interact with them using only the keyboard,
**So that** I can fill out the form efficiently.
**Prerequisite:** Story 1.4

#### Acceptance Criteria

1. AC1: The `Tab` key MUST move focus to the next form field.
2. AC2: `Shift+Tab` MUST move focus to the previous field.
3. AC3: `Enter` MUST confirm selection in fields like `Select` and `Confirm`, or move to the next field in `Input`/`Textarea` (`huh` default behavior).
4. AC4: `Space` MUST toggle selection in `MultiSelect` fields.
5. AC5: `Arrow` keys (Up/Down) MUST allow navigation within options of `Select` and `MultiSelect` fields.
6. AC6: The currently focused field MUST be visually highlighted.

### Story 1.6 Submission and JSON Output

**As a** script developer,
**I want** shantilly to collect the data entered in the TUI form and print it in JSON format to `stdout` upon submission,
**So that** my script can easily consume the results.
**Prerequisite:** Story 1.5

#### Acceptance Criteria

1. AC1: The `huh` library MUST be configured to allow form submission (typically after the last field or via an implicit button).
2. AC2: After submission, the `bubbletea` application MUST terminate.
3. AC3: The values entered in each form field MUST be collected.
4. AC4: The collected data MUST be formatted as a JSON object, where keys are the `key` values defined in the YAML for each field.
5. AC5: The resulting JSON object MUST be printed to `stdout`.
6. AC6: The TUI MUST be cleared from the screen before printing the JSON.

### Story 1.7 Basic Styling and Alignment

**As a** shantilly user,
**I want** the TUI form to be presented centered with adequate spacing, and to adapt to basic terminal resizing,
**So that** the interface is visually pleasant and functional.
**Prerequisite:** Story 1.4

#### Acceptance Criteria

1. AC1: The main container for the `huh` form MUST be styled using `lipgloss`.
2. AC2: The form MUST be rendered horizontally centered within the terminal window if there is sufficient space.
3. AC3: Consistent padding MUST be applied around the form.
4. AC4: The TUI MUST handle resize messages (`tea.WindowSizeMsg`) and re-render the basic layout (maintaining centering if possible).
5. AC5: If the terminal is resized to a very small width, the form MUST still be minimally usable (`huh` might truncate or wrap lines by default).

### Story 1.8 Static Build and Distribution

**As a** shantilly developer,
**I want** a build process that generates static, cross-compiled binaries for the target platforms,
**So that** the application can be easily distributed.
**Prerequisite:** Story 1.1

#### Acceptance Criteria

1. AC1: The `Makefile` (or build script) MUST include targets to compile the `shantilly` binary.
2. AC2: The build process MUST use flags (`CGO_ENABLED=0`, `ldflags="-s -w"`) to ensure a static and optimized binary.
3. AC3: Targets (or a script) MUST exist to cross-compile the binary for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64).
4. AC4: The generated binaries MUST be executable on their respective platforms (basic verification).

## Checklist Results Report

[[LLM: PM CHECKLIST EXECUTION INSTRUCTIONS (`pm-checklist.md`)

**Project:** shantilly (Greenfield, CLI/TUI - No Web/Mobile UI)
**Document Under Review:** `docs/prd.md` (current version in Canvas)

**Mode:** Comprehensive (YOLO) - Analyze all and present final report.

**Context:**

* PRD derived from detailed Project Brief [cite: Project Brief\_ shantilly.md].
* Scope strictly MVP: `form` subcommand, YAML stdin -\> TUI -\> JSON stdout [cite: docs/prd.md].
* Tech stack is Go + Charmbracelet [cite: docs/prd.md].
* Project has NO web/mobile UI, skip related checklist sections.

**Process:**

1. Read EACH item in `pm-checklist.md` [cite: team-fullstack.txt].
2. Verify if the item is covered in the current `docs/prd.md`.
3. Evaluate the *quality* and *completeness* of coverage.
4. Mark each item as ✅ PASS, ❌ FAIL, ⚠️ PARTIAL, or N/A.
5. For ❌ FAIL or ⚠️ PARTIAL, note *why*.
6. Calculate pass rates per section.
7. Synthesize results in the report template below. Be specific.
    ]]

### Executive Summary

* **Overall PRD Completeness:** 100%
* **MVP Scope Appropriateness:** Ideal
* **Architecture Readiness:** Ready
* **Critical Gaps or Concerns:** No critical gaps remaining.

### Category Analysis

| Category                         | Status | Critical Issues                     |
|:---------------------------------|:-------|:------------------------------------|
| 1. Problem Definition & Context  | ✅ PASS |                                     |
| 2. MVP Scope Definition          | ✅ PASS |                                     |
| 3. User Experience Requirements  | ✅ PASS | (Simplified for TUI)                |
| 4. Functional Requirements       | ✅ PASS |                                     |
| 5. Non-Functional Requirements   | ✅ PASS | (NFR8 added for errors)             |
| 6. Epic & Story Structure        | ✅ PASS |                                     |
| 7. Technical Guidance            | ✅ PASS | (YAML added to Assumptions)         |
| 8. Cross-Functional Requirements | ✅ PASS | (YAML Structure & Errors addressed) |
| 9. Clarity & Communication       | ✅ PASS |                                     |

### Top Issues by Priority

* **BLOCKERS:** None.
* **HIGH:** None (Previous recommendations implemented).
* **MEDIUM:** None.
* **LOW:** None.

### MVP Scope Assessment

The MVP scope remains **Ideal** and aligned with the Project Brief [cite: Project Brief\_ shantilly.md].

### Technical Readiness

The **Technical Assumptions**, now including the expected YAML structure, provide a complete foundation for the Architect [cite: docs/prd.md].

### Recommendations

No further recommendations. Previous items have been addressed.

### Final Decision

**READY FOR ARCHITECT**: The PRD is complete, consistent, and addresses the points raised by the checklist. It is ready for handoff to the Architecture phase.

## Next Steps

### Prompt for UX Expert

N/A (The project is a TUI CLI and does not require a Web/Mobile UI/UX expert at this time).

### Prompt for Architect

**To:** Winston (Architect) 🏗️
**From:** John (PM) 📋

**Subject:** Start Architecture Phase for shantilly (MVP) - Final PRD Attached

Winston,

The PRD for the shantilly MVP is finalized (see `docs/prd.md`). It defines the requirements for building a CLI tool (Go + Charmbracelet) that renders TUI forms from a YAML definition on `stdin` (structure now defined in PRD) and returns results as JSON on `stdout`. Error handling has also been specified (NFR8).

Please review the finalized PRD, especially the **Technical Assumptions** (including YAML structure) and the **Functional/Non-Functional Requirements**.

Your next step is to create the **Architecture Document** (`architecture-tmpl.yaml`) [cite: team-fullstack.txt]. Key points to address:

1. **Detailed Project Structure:** Detail the Go package organization (`cmd/`, `internal/`, sub-packages like `tui`, `config`, `parser`).
2. **Data Structures:** Define the Go structs to represent the YAML configuration defined in the PRD.
3. **Internal Data Flow:** How data flows from YAML parsing (`gopkg.in/yaml.v3`) to the Go struct, into the `bubbletea`/`huh` model, and finally collected for JSON output.
4. **Bubbletea Model:** Specific structure of the `bubbletea` model (state, messages, `Init`/`Update`/`View`).
5. **`huh` Integration:** How `huh` components will be dynamically created from the YAML config and managed within `bubbletea`.
6. **Error Handling:** Architecture for capturing and reporting parsing errors (stderr, non-zero exit status) and TUI errors, per NFR8.
7. **Coding Standards:** Define specific Go patterns for this project.
8. **Testing Strategy:** Detail the approach for unit tests.

I am available to clarify any requirements.
