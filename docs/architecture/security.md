# Security

## Input Validation

* **Focus:** YAML input via `stdin` or `--file`.
* **Location:** `ConfigParser` (`internal/config/parser.go`).
* **Required Rules:**
  * Validate `FieldConfig.Type` against known `huh` component types (e.g., "input", "textarea", "select", "multiselect", "confirm", "note"). Reject invalid types via `ErrorHandler` (NFR8).
  * Handle malformed YAML gracefully via `ErrorHandler` (NFR8).

## AuthN / AuthZ / Secrets / API Security / Data Protection

* **N/A:** Not applicable for MVP (local CLI, no network, no sensitive data persistence).

## Dependency Security

* **Scanning Tool:** `govulncheck`.
* **Update Policy:** Review/update dependencies regularly (e.g., via Dependabot).
* **Approval Process:** Evaluate new dependencies before adding.

## Security Testing

* **SAST:** `gosec` (integrated via `golangci-lint` in `.golangci.yml`).
* **DAST / Pentest:** N/A for MVP.
