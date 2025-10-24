# Technical Brainstorming Report: Shantilly

Session Date: 2025-10-22

Participants: User, Mary (BMad Business Analyst)

Topic: Initial planning and design for the shantilly CLI tool.

## Executive Summary

This session aimed to define the vision, architecture, requirements, and initial implementation plan for `shantilly`, a modern CLI tool for creating declarative TUIs for shell scripts, serving as an alternative to `dialog` and `whiptail`. We defined the core technology stack (Go + Charmbracelet Ecosystem), a modular architecture separating CLI and TUI, an incremental implementation plan focused on delivering value quickly (MVP with `form` subcommand), and the preferred configuration approach (YAML via `stdin`). We also explored advanced features like complex layouts, mouse support, floating windows (modals), and SSH access, positioning them in post-MVP phases. The main decisions were consolidated into the Project Brief (`docs/project-brief.md`).

## Session Details

### Topic and Goals

- **Topic:** Design the `shantilly` tool.

- **Goals:**
  
  - Define the problem and the proposed solution.
  
  - Outline the architecture and technology stack.
  
  - Identify key functional and non-functional requirements.
  
  - Create an incremental implementation plan (MVP and post-MVP).
  
  - Evaluate configuration formats for use in shell scripts.
  
  - Discuss advanced features (layouts, windows, mouse, SSH).

### Techniques Used (Implicitly)

- Requirements Analysis (based on the initial description).

- Technical Discussion and Library Comparison (Charm vs. tview, `winman` vs. `bubbletea-overlay`).

- Incremental Planning (defining MVP and subsequent phases).

- Pros and Cons Analysis (configuration formats, implementation approaches).

### Key Ideas and Concepts Generated

- **Declarative Tool:** `shantilly` will interpret a definition (YAML) to render the TUI.

- **TUI Core:** Based on `bubbletea` + `lipgloss`.

- **Components:** Use `huh` for forms, `bubbles` for lists/menus, `bubbleboxer` for layouts, `bubblezone` for mouse, `harmonica` for animations, `wish` for SSH.

- **Architecture:** Strict separation between `cmd/` (CLI/Cobra) and `internal/` (TUI/Bubbletea).

- **MVP:** Focus on the `form` subcommand reading YAML from `stdin` and returning JSON to `stdout`.

- **Post-MVP:** Layouts, menus, mouse, modal dialogs (`bubbletea-overlay` or custom), SSH.

- **Floating Windows:** Implement functionality inspired by `winman` within `bubbletea`, possibly using `bubbletea-overlay` as an initial base.

- **Configuration:** YAML via `stdin` as the primary method, supplemented by simple flags.

- **Distribution:** Single static binary for Linux, macOS, Windows.

### Decisions Made

- **TUI Framework:** Confirmation of the exclusive use of the Charmbracelet ecosystem (`bubbletea`, `lipgloss`, `huh`, etc.). `tview` and `winman` (based on `tview`) will not be used directly due to incompatibility.

- **Primary Configuration Format:** YAML read via `stdin` (here-docs). Flags will be used for simple overrides (e.g., `--title`).

- **Incremental Plan:** Adoption of the proposed phase plan (Phase 0: CLI Foundation -> Phase 1: TUI Engine -> Phase 2: MVP `form` -> Phase 3+: Advanced Components, Windows, Mouse -> Phase 4+: SSH, Tests).

- **Window Implementation:** Floating/modal window functionality will be implemented within `bubbletea`, inspired by `winman`, using `bubbletea-overlay` as a possible starting point.

### Immediate Next Steps (Completed)

- Finalization and approval of the Project Brief (`docs/project-brief.md`).

### Next Step (Pending)

- Handoff to the Product Manager (PM - John) to initiate PRD creation based on the Project Brief.

*Report generated using the BMAD-METHOD™ framework*
