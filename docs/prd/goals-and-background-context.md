# Goals and Background Context

## Goals

* Simplify the creation of interactive and modern TUI interfaces for shell scripts [cite: shantilly Product Requirements Document (PRD).md].
* Offer a declarative (YAML) and more powerful alternative to `dialog` and `whiptail` [cite: shantilly Product Requirements Document (PRD).md].
* Ensure portability and consistency across Linux, macOS, and Windows via a single static binary [cite: shantilly Product Requirements Document (PRD).md].
* Facilitate integration with shell script pipelines (read `stdin` YAML, write `stdout` JSON) [cite: shantilly Product Requirements Document (PRD).md].
* Achieve broad adoption by the script developer community (Bash, Zsh, PowerShell, etc.) [cite: Project Brief_ shantilly.md].
* Deliver a robust and viable MVP (v1.0), focused on the `form` subcommand, within 3 months [cite: Project Brief_ shantilly.md].
* Enable the creation of significantly richer interfaces (layouts, rich components) than traditional tools [cite: Project Brief_ shantilly.md].

## Background Context

Shell script developers needing complex user interactions are currently underserved. Traditional tools like `dialog` and `whiptail` are functional but severely limited in components, layout control, and aesthetics, lacking mouse support [cite: Project Brief_shantilly.md, shantilly Product Requirements Document (PRD).md]. The alternative, using programmatic TUI libraries (like Bubbletea), requires knowledge of languages like Go or Python, adding disproportionate complexity for script-based automation [cite: Project Brief_ shantilly.md, shantilly Product Requirements Document (PRD).md].

shantilly fills this gap. It's a portable CLI tool (single static binary) enabling the *declarative creation* of rich, modern TUIs [cite: Project Brief_shantilly.md]. The user defines the interface (e.g., a form) in YAML, passes it to `shantilly` via `stdin`, and the tool renders the interactive TUI [cite: Project Brief_ shantilly.md]. Collected data is then returned as structured JSON on `stdout`, allowing easy integration into script pipelines [cite: Project Brief_shantilly.md]. This PRD focuses on the MVP (Minimum Viable Product) scope to validate the core functionality [cite: Project Brief_ shantilly.md].

## Change Log

| Date       | Version | Description                                                                                                    | Author     |
| :--------- | :------ | :------------------------------------------------------------------------------------------------------------- | :--------- |
| 2025-10-22 | 0.1.0   | Initial PRD draft based on Project Brief.                                                                      | John (PM)  |
| 2025-10-23 | 0.1.1   | Added UI/UX section and refined MVP layout.                                                                    | John (PM)  |
| 2025-10-23 | 0.1.2   | Added Technical Assumptions section.                                                                           | John (PM)  |
| 2025-10-23 | 0.1.3   | Added Epic List (MVP).                                                                                         | John (PM)  |
| 2025-10-23 | 0.1.4   | Added Epic 1 Details (MVP) with Stories.                                                                       | John (PM)  |
| 2025-10-23 | 0.1.5   | Completed PM Checklist and Next Steps section.                                                                 | John (PM)  |
| 2025-10-23 | 0.2.0   | Implemented PM Checklist recommendations (YAML Structure and Error Handling - NFR8). Updated Architect prompt. | John (PM)  |
| 2025-10-27 | 0.3.0   | Added preventive epics (3-7) for future roadmap planning and process improvement.                              | Sarah (PO) |
