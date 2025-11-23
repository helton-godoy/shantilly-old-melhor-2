# Next Steps

## Handoff Prompt for Developer Agent

**To:** Dev Agent 👨‍💻
**From:** Winston (Architect) 🏗️
**Subject:** Implement Shantilly MVP - Story 1.1 (CLI Foundation)

**Context:**
The architecture for the `shantilly` CLI (MVP) is finalized and validated (see `docs/architecture.md`). This project uses Go (1.24.2+) and the Charmbracelet ecosystem (`bubbletea`, `huh`, `lipgloss`) along with `cobra` and `yaml.v3`. We are adopting the user-provided project template files (`.golangci.yml`, `.goreleaser.yaml`, `go.mod`, etc.) for tooling and standards.

**Goal:**
Implement Story 1.1 from the PRD (`prd.md`) - "CLI Foundation and Project Structure".

**Key Architectural Guidance (from `docs/architecture.md`):**

* **Source Tree:** Follow the defined structure, place main logic in `cmd/shantilly/main.go`.
* **Tech Stack:** Use `spf13/cobra` (v1.8.x).
* **Components:** Implement the basic `CLI (CobraCmd)` component.
* **Coding Standards:** Adhere strictly to `.golangci.yml` rules. Use `gofumpt` for formatting. Ensure `./lint.sh` passes before marking complete.

**Tasks (derived from Story 1.1 ACs):**

1. Initialize the Go module structure as defined in the `Source Tree` section (if not already present from the template). Ensure `go.mod` reflects the project name and Go version.
2. Implement the root `shantilly` command using `spf13/cobra` in `cmd/shantilly/main.go`.
3. Add the `form` subcommand to the root command (also in `main.go`). Define the `--file` string flag for it.
4. Implement the `RunE` (or `Run`) function for the `form` subcommand:
      * Check if the `--file` flag was provided.
          * If yes, read the content of the specified file. Handle file read errors.
          * If no, read all bytes from `os.Stdin`. Handle potential stdin read errors.
      * Propagate any read errors up to be handled by the `ErrorHandler`.
      * *(For this story only)*: Print a placeholder message like "Form executed. Read X bytes." to `os.Stdout`.
      * Return `nil` on success.
5. Ensure the `ErrorHandler` (`internal/util/errorhandler.go`) is implemented as defined in the architecture doc (handling `nil` error, printing non-nil errors to `stderr`, using `os.Exit` with defined codes). Update `main.go`'s `RunE` to call `util.Handle(err)` appropriately on error return.
6. Create/Update the `Makefile` with a basic `build` target: `go build -o shantilly ./cmd/shantilly`. Add `lint` (`./lint.sh`) and `test` (`go test -v -race ./...`) targets.
7. Ensure the code passes `make lint` and `make test` (though no specific tests for this story yet).

**Acceptance Criteria (from Story 1.1):**

* [AC1] Go repo initialized (`cmd/`, `internal/`).
* [AC2] `spf13/cobra` used.
* [AC3] `shantilly` root and `form` subcommand exist.
* [AC4] `shantilly form` reads `stdin` OR uses `--file` flag.
* [AC5] `shantilly form` prints placeholder to `stdout` on success.
* [AC6] `Makefile` includes `build`, `lint`, `test` targets.
* (Implicit) Errors during input reading lead to `stderr` message and non-zero exit code via `ErrorHandler`.

**Definition of Done:**

* All tasks completed.
* All ACs met.
* Code is formatted (`gofumpt`).
* Linter passes (`make lint`).
* Build succeeds (`make build`).
* Provide the complete content for `cmd/shantilly/main.go` and `internal/util/errorhandler.go`. List any other created/modified files (like `Makefile`).
