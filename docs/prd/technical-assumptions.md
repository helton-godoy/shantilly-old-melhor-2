# Technical Assumptions

## Repository Structure: Monorepo

The project will be structured as a simple Go monorepo. This facilitates initial organization and potential future addition of subcommands or related libraries [cite: Project Brief_ shantilly.md, Relatório Técnico do Brainstorming Shantilly.md]. The structure will follow standard Go conventions, with CLI logic in `cmd/shantilly` and TUI/parsing logic in `internal/`.

## Service Architecture: Monolithic CLI Application

For the MVP, shantilly is a monolithic CLI application that executes, processes input, and terminates [cite: Project Brief_ shantilly.md]. The architecture will be modular internally (separating CLI and TUI), but there will be no microservices or network communication [cite: Relatório Técnico do Brainstorming Shantilly.md]. SSH server mode is a Post-MVP vision (NFR7) [cite: docs/prd.md].

## Testing Requirements: Focus on Unit Tests

Given the MVP time constraints and the challenges of testing TUIs [cite: Project Brief_shantilly.md], the primary focus will be on robust unit tests for the YAML parsing logic and internal business logic [cite: Relatório Técnico do Brainstorming Shantilly.md]. UI integration tests will be limited to manual smoke tests on the main platforms [cite: Project Brief_ shantilly.md].

## Expected YAML Structure (`stdin`)

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

## Additional Technical Assumptions and Requests

* **Language:** Go (NFR3) [cite: docs/prd.md].
* **Core Frameworks/Libraries:** `spf13/cobra` (FR1) [cite: docs/prd.md], `charmbracelet/bubbletea` (FR4) [cite: docs/prd.md], `charmbracelet/lipgloss` (FR7) [cite: docs/prd.md], `charmbracelet/huh` (FR3) [cite: docs/prd.md], `gopkg.in/yaml.v3` [cite: Project Brief\_ shantilly.md].
* **Build & Distribution:** Standard Go build process with flags for static compilation (`CGO_ENABLED=0`) and cross-compilation for NFR2 architectures [cite: docs/prd.md].
