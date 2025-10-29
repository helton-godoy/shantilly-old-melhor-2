# Coding Standards

Defined primarily by the user-provided template files. Adherence is mandatory and checked automatically.

## Core Standards

* **Language & Runtime:** Go `1.24.2+` (per template `go.mod`).
* **Style & Linting:** Governed by `.golangci.yml` (per template). Checked via `golangci-lint run ./...` and pre-commit hooks.
* **Formatting:** Governed by `gofumpt`. Checked via `gofumpt -w .` and pre-commit hooks.
* **Test Organization:** `_test.go` files in the same package (Go standard).

## Naming Conventions

* Standard Go conventions (`camelCase`, `PascalCase`). Enforced by `stylecheck` linter in `.golangci.yml`.

## Critical Rules (from `.golangci.yml`)

1. **Error Handling (`errcheck`, `errorlint`, `wrapcheck`):** Check/wrap all errors. Use `errors.Is/As`.
2. **No Unused Code (`unused`, `ineffassign`):** Remove dead code.
3. **Simplicity (`gocyclo`, `funlen`, `nestif`):** Keep functions short, low complexity.
4. **Performance (`prealloc`, `gocritic`):** Pre-allocate slices, follow `gocritic` advice.
5. **No Magic Numbers (`mnd`):** Use named constants.
6. **Full Struct Init (`exhaustivestruct`):** Initialize all struct fields explicitly.

## TUI-Specific Guidelines (Bubble Tea + Charm)

**Preventing TUI Rendering Failures:**

1. **Window Size Handling (`tea.WindowSizeMsg`):**
   * **Mandatory:** Always handle `tea.WindowSizeMsg` in the `Update` function to ensure responsive behavior.
   * **Pattern:** Update model dimensions and trigger re-renders when terminal size changes.
   * **Testing:** Use `teatest.WithInitialTermSize` in integration tests to validate different terminal dimensions.

2. **Layout Primitives (`lipgloss`):**
   * **Strong Recommendation:** Use `lipgloss` primitives (`Width`, `Height`, `JoinHorizontal`, `JoinVertical`) for layout management.
   * **Purpose:** Ensures consistent spacing, alignment, and responsive behavior across different terminal sizes.
   * **Pattern:** Define layout constraints explicitly rather than relying on hardcoded spacing.

3. **Responsive Components:**
   * **Terminal Width Awareness:** Always consider terminal width constraints when designing TUI layouts.
   * **Truncation Strategy:** Implement text truncation or responsive components (`bubbles/list`, `bubbles/viewport`) for content that may exceed terminal width.
   * **Component Selection:** Prefer existing `bubbles` components over custom implementations for complex UI patterns.

4. **Centralized Styling:**
   * **Theme Definition:** Define a centralized `lipgloss` theme to ensure consistent styling across all TUI components.
   * **Consistency:** Apply theme styles uniformly to maintain visual coherence and prevent styling-related rendering issues.

## Documentation Standards

### Cross-Reference Format

To maintain consistency and ensure all references are valid and accessible:

* **Standard Format:** `docs/{pasta}/{arquivo}.md#{nome-da-seção-nível-2-ou-superior}`
* **Mandatory Prefix:** All cross-references MUST include the `docs/` prefix for clarity
* **Section Anchors:** Use level 2+ section headers (##, ###) as anchors, not arbitrary text
* **Validation:** All cross-references should be validated to ensure target sections exist
* **Examples:**
  * ✅ `docs/architecture/tech-stack.md#technology-stack-table`
  * ✅ `docs/architecture/data-models.md#formconfig-input-yaml`
  * ❌ `architecture/tech-stack.md#huh` (missing docs/ prefix)
  * ❌ `docs/architecture/components.md#bubbletea-model-structure` (section doesn't exist)

### Reference Maintenance

* **Validation Tools:** Consider implementing automated link checking in CI/CD pipeline
* **Review Process:** All cross-references must be verified during story validation
* **Update Process:** When sections are moved/renamed, update all references accordingly

## Language-Specific Guidelines (Go)

* Follow "Effective Go".
* Use pointers judiciously.
* Prefer small interfaces.
