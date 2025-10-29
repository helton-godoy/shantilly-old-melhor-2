# Test Strategy and Standards

## Testing Philosophy

* **Approach:** Test-After for MVP. Focus on unit tests first, with integration tests for TUI components using `teatest`.
* **Coverage Goals:** No strict % target for MVP, but high coverage (\>80%) for `internal/config` and critical TUI flows.
* **Test Pyramid (MVP):** Heavy on Unit Tests, complemented by TUI Integration Tests using `teatest`, and Manual TUI tests.

## Test Types and Organization

**Unit Tests**

* **Framework:** Go `testing` package (v1.24.2+).
* **File Convention:** `_test.go` in the same package.
* **Location:** Primarily `internal/config/` and `internal/tui/`.
* **Mocking:** No external dependencies to mock in MVP. Use interfaces for potential future manual fakes/stubs if needed between internal components.
* **AI Agent Requirements:** Generate comprehensive tests for `internal/config/parser.go`, covering valid/invalid YAML cases. Follow AAA pattern. Maintain pure unit tests for `Update` logic in TUI components.

**TUI Integration Tests (teatest)**

* **Framework:** `charmbracelet/bubbles/teatest` for testing Bubble Tea components.
* **Purpose:** Validate critical user flows and assert on textual output (string output), including basic layout and presence of styled elements via `lipgloss`.
* **Key Features:**
  * Use `teatest.WithInitialTermSize` to run tests with different terminal sizes for responsive behavior validation.
  * Focus on testing complete TUI workflows rather than individual component rendering.
  * Assert on final rendered output strings to verify layout and styling.
* **Location:** `internal/tui/` alongside unit tests.
* **AI Agent Requirements:** Create integration tests for critical TUI flows using `teatest`, ensuring proper handling of `tea.WindowSizeMsg` and responsive layout across different terminal dimensions.

**Integration Tests**

* **N/A:** Out of scope for MVP.

**E2E Tests**

* **N/A:** Out of scope for MVP. Manual testing covers this.

## Test Data Management

* **Strategy:** Example `.yaml` files in `examples/` directory act as fixtures for `ConfigParser` tests.

## Continuous Testing

* **CI Integration:** GitHub Actions runs `go test ./...` (via `lint.sh` or build workflow).
* **Local Testing:** Developers use `./lint.sh` (from template) which includes `go test -v -race ./...`.
* **Performance/Security Tests:** N/A for MVP.
