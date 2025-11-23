# Tech Stack

## Build and Distribution Infrastructure

* **Build Platform:** GitHub Actions.
* **Key Services:**
  * `golangci-lint` (for code linting, configured via `.golangci.yml`).
  * `GoReleaser` (for build automation, cross-compilation, packaging, configured via `.goreleaser.yaml`).
* **Deployment Target (Distribution):** GitHub Releases.
* **Regions:** N/A (global binary distribution).

## Technology Stack Table

This table is the single source of truth for project dependencies. Versions are based on the user-provided template and current stable releases (as of late 2025).

| Category           | Technology                | Version | Purpose                                     | Rationale                                                                   |
|:-------------------|:--------------------------|:--------|:--------------------------------------------|:----------------------------------------------------------------------------|
| **Language**       | Go                        | 1.24.2+ | Primary development language                | Requirement (NFR3), User Template (`go.mod`). Modern, fast, static builds.  |
| **CLI Framework**  | `spf13/cobra`             | v1.8.x  | Core CLI structure (commands, flags, stdin) | Requirement (FR1). Go standard.                                             |
| **TUI Engine**     | `charmbracelet/bubbletea` | v0.26.x | TUI state management (TEA)                  | Requirement (FR4), User Template (`go.mod`). Powerful, flexible.            |
| **TUI Components** | `charmbracelet/huh`       | v0.4.x  | Declarative TUI form generation             | Requirement (FR3). Drastically simplifies form creation.                    |
| **TUI Styling**    | `charmbracelet/lipgloss`  | v0.11.x | Terminal styling (colors, layout)           | Requirement (FR7), User Template (`go.mod`). Needed for alignment.          |
| **YAML Parsing**   | `gopkg.in/yaml.v3`        | v3.0.x  | Decoding `stdin`/file YAML to Go structs    | Requirement (Story 1.2). Robust and widely used.                            |
| **Testing**        | `testing` (Go Stdlib)     | 1.24.2+ | Unit testing (MVP focus)                    | Go standard (PRD Tech Assumptions).                                         |
| **Linter**         | `golangci-lint`           | v1.59.x | Static analysis, code quality               | Best practice. User Template (`.golangci.yml`, `.pre-commit-config.yaml`).  |
| **Formatter**      | `gofumpt`                 | latest  | Strict code formatting                      | Best practice. User Template (`.pre-commit-config.yaml`, `lint.sh`).        |
| **Build/Release**  | `GoReleaser`              | v1.26.x | Build automation, cross-compilation (NFR2)  | Best practice. User Template (`.goreleaser.yaml`). Simplifies distribution. |
| **Pre-commit**     | `pre-commit` framework    | latest  | Local quality checks before commit          | Best practice. User Template (`.pre-commit-config.yaml`).                   |
