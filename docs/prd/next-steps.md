# Next Steps

## Prompt for UX Expert

N/A (The project is a TUI CLI and does not require a Web/Mobile UI/UX expert at this time).

## Prompt for Architect

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
