# Action Items Checklist - Code Quality Assessment

Generated: 2025-11-12
Status: READY FOR IMPLEMENTATION

---

## PRIORITY 1: IMMEDIATE (Current Sprint)

### P1.1: Enable Coverage Threshold in CI/CD ⏱️ 30 min
- [ ] Edit `.github/workflows/build.yml`
- [ ] Add coverage gate step after test run:
  ```yaml
  - name: Check Coverage Threshold
    run: |
      COVERAGE=$(go tool cover -func=coverage.txt | tail -1 | awk '{print $3}' | sed 's/%//')
      echo "Coverage: $COVERAGE%"
      if (( $(echo "$COVERAGE < 70" | bc -l) )); then
        echo "FAIL: Coverage $COVERAGE% below 70% threshold"
        exit 1
      fi
  ```
- [ ] Create `.codecov.yml` with:
  ```yaml
  coverage:
    status:
      project:
        default:
          target: 70
          threshold: 5
  ```
- [ ] Test locally: `go test -coverprofile=coverage.txt -covermode=atomic ./...`
- [ ] Verify gate works in CI
- **Assignee**: DevOps / CI/CD Lead
- **Expected Outcome**: Coverage gating active, unknown coverage visible

---

### P1.2: Expand ScriptRunner Tests (4 hours) ⏱️
- [ ] Open `internal/runtime/runner/runner_test.go`
- [ ] Add test for process timeout:
  ```go
  func TestScriptRunner_ProcessTimeout(t *testing.T)
  // Test ctx.Done() handling, SIGTERM propagation
  ```
- [ ] Add test for rapid replacement on same target:
  ```go
  func TestScriptRunner_RapidReplacementSameTarget(t *testing.T)
  // Start script, immediately start another before first completes
  // Verify first is cancelled and second runs
  ```
- [ ] Add test for large stdout/stderr:
  ```go
  func TestScriptRunner_LargeOutput(t *testing.T)
  // Generate >1MB output, verify no buffer overflow
  ```
- [ ] Add test for signal escalation:
  ```go
  func TestScriptRunner_SignalEscalation(t *testing.T)
  // Verify SIGTERM -> SIGKILL escalation on timeout
  ```
- [ ] Add test for permission denied:
  ```go
  func TestScriptRunner_PermissionDenied(t *testing.T)
  // Script file not executable, verify proper error event
  ```
- [ ] Run `go test ./internal/runtime/runner/ -v` to verify
- **Assignee**: Backend Developer
- **Expected Outcome**: 5-10 new tests, improved runner coverage

---

### P1.3: Expand Form Model Tests (3 hours) ⏱️
- [ ] Open `internal/tui/model_test.go`
- [ ] Add test for rapid terminal resize:
  ```go
  func TestModel_RapidTerminalResize(t *testing.T)
  // Send multiple WindowSizeMsg rapidly, verify no panic/corruption
  ```
- [ ] Add test for error message ordering:
  ```go
  func TestModel_ErrorMessageOrdering(t *testing.T)
  // Verify error display order is deterministic (sorted)
  ```
- [ ] Add test for completion percentage:
  ```go
  func TestModel_CompletionPercentageEdgeCases(t *testing.T)
  // Test 0%, 50%, 100%, single field, many fields
  ```
- [ ] Add test for view rendering consistency:
  ```go
  func TestModel_ViewRenderingConsistency(t *testing.T)
  // Call View() multiple times with same state, verify output identical
  ```
- [ ] Run `go test ./internal/tui/ -v` to verify
- **Assignee**: Frontend Developer
- **Expected Outcome**: 4-5 new tests, improved model coverage

---

## PRIORITY 2: SHORT TERM (Next 2 Weeks)

### P2.1: Add Security Test Suite (3 hours) ⏱️
- [ ] Open `internal/runtime/runner/runner_test.go`
- [ ] Add test for path traversal:
  ```go
  func TestScriptRunner_SecurityPolicy_PathTraversal(t *testing.T)
  // Attempt: ../../../etc/passwd, /etc/passwd, etc.
  // Verify all denied by SecurityPolicy
  ```
- [ ] Add test for symlink traversal:
  ```go
  func TestScriptRunner_SecurityPolicy_SymlinkTraversal(t *testing.T)
  // Create symlink to /etc/passwd, attempt execution
  // Verify blocked
  ```
- [ ] Add test for command injection:
  ```go
  func TestScriptRunner_SecurityPolicy_CommandInjection(t *testing.T)
  // Args with shell metacharacters (;, |, &&, etc.)
  // Verify treated as literal args, not command injection
  ```
- [ ] Add test for whitelist bypass via unicode:
  ```go
  func TestScriptRunner_SecurityPolicy_UnicodeBypass(t *testing.T)
  // Try unicode normalization to bypass whitelist
  // Verify blocked
  ```
- [ ] Run `go test ./internal/runtime/runner/ -v` to verify
- **Assignee**: Security Engineer / Backend Developer
- **Expected Outcome**: 4-5 security tests, reduced vulnerability surface

---

### P2.2: Add Validation Boundary Case Tests (3 hours) ⏱️
- [ ] Open `internal/config/validation_test.go`
- [ ] Add test for number boundaries:
  ```go
  func TestFieldValidator_NumberBoundaries(t *testing.T)
  // Test: max float64, min float64, -0.0, infinity, NaN
  ```
- [ ] Add test for date edge cases:
  ```go
  func TestFieldValidator_DateEdgeCases(t *testing.T)
  // Test: leap year Feb 29, year boundaries, DST transitions
  ```
- [ ] Add test for file validation edge cases:
  ```go
  func TestFieldValidator_FileSymlinks(t *testing.T)
  // Test: symlinks, broken symlinks, permission denied
  ```
- [ ] Add test for text unicode edge cases:
  ```go
  func TestFieldValidator_TextUnicodeEdgeCases(t *testing.T)
  // Test: null bytes, RTL text, zero-width chars
  ```
- [ ] Run `go test ./internal/config/ -v` to verify
- **Assignee**: QA / Test Developer
- **Expected Outcome**: 4-5 new tests, improved validation coverage

---

### P2.3: Document Coverage Baseline (1 week) ⏱️
- [ ] Run coverage collection:
  ```bash
  go test -coverprofile=coverage.txt -covermode=atomic ./...
  go tool cover -func=coverage.txt > coverage-by-function.txt
  go tool cover -html=coverage.txt -o coverage.html  # For visual inspection
  ```
- [ ] Extract per-package percentages
- [ ] Create `docs/qa/coverage-baseline.md` with:
  - Overall coverage %
  - Per-package breakdown
  - Coverage targets by priority
  - Trending metrics
- [ ] Document current baseline (expected ~60-70% based on test distribution)
- [ ] Set targets: 75% overall, 80% core packages (runner, modal, event)
- [ ] Commit baseline to git for future tracking
- **Assignee**: QA Lead / Project Manager
- **Expected Outcome**: Documented coverage baseline, visibility, targets

---

### P2.4: Resolve Outstanding TODOs (1 hour) ⏱️
- [ ] Fix TODO in `internal/tui/components/error_display.go` line 38:
  ```go
  // Before: errors map iterated in arbitrary order
  // After: Sort keys before rendering
  var errorMessages []string
  var keys []string
  for k := range e.errors {
    keys = append(keys, k)
  }
  sort.Strings(keys)
  for _, k := range keys {
    errorMessages = append(errorMessages, e.theme.FieldError.Render("• "+e.errors[k]))
  }
  ```
- [ ] Add unit test for sorted error output
- [ ] Document reactive flow TODO in `internal/runtime/event/manager.go` line 16:
  - Link to `docs/architecture/core-workflows.md` 
  - Add reference to EventManager responsibilities
- [ ] Verify no new TODOs introduced
- **Assignee**: Any Developer
- **Expected Outcome**: Deterministic error display, documentation updated

---

## PRIORITY 3: MEDIUM TERM (Next 4 Weeks)

### P3.1: Refactor ScriptRunner Package (2-4 weeks) ⏱️
- [ ] Plan refactoring into sub-packages:
  - `internal/runtime/runner/executor.go` - Command execution, process lifecycle
  - `internal/runtime/runner/security.go` - SecurityPolicy validation, path checks
  - `internal/runtime/runner/events.go` - Event emission logic
  - Keep `internal/runtime/runner/runner.go` as public API, delegation
- [ ] Move methods to appropriate files
- [ ] Update test organization to match
- [ ] Verify all tests still pass
- [ ] Document sub-package responsibilities
- **Assignee**: Architecture / Senior Developer
- **Expected Outcome**: 3 focused packages instead of 596-line monolith

---

### P3.2: Refactor Form Model (2 weeks) ⏱️
- [ ] Plan refactoring:
  - `internal/tui/view/render.go` - All View() logic
  - `internal/tui/view/layout.go` - renderProgressIndicator, renderHelpText, etc.
  - Keep `internal/tui/model.go` for state and Update
- [ ] Move rendering methods to view sub-package
- [ ] Create view.Component interface for easier testing
- [ ] Update tests to test view separately from model
- **Assignee**: Frontend Developer
- **Expected Outcome**: Cleaner separation of concerns

---

### P3.3: Add Concurrency Stress Tests (1-2 weeks) ⏱️
- [ ] Add stress test for modal stack under concurrent load:
  ```go
  func TestModalStack_ConcurrentOperations(t *testing.T)
  // 100 goroutines pushing/popping simultaneously
  // Verify no data corruption, deadlock, or panic
  ```
- [ ] Add stress test for process replacement:
  ```go
  func TestScriptRunner_ConcurrentReplacementStress(t *testing.T)
  // Start 10 rapid sequences of script replacement
  // Verify proper cancellation and cleanup
  ```
- [ ] Run with `-race` flag extensively
- **Assignee**: QA / Test Developer
- **Expected Outcome**: Confidence in concurrency safety

---

## PRIORITY 4: LONG TERM (Next 8+ Weeks)

### P4.1: Coverage Trending Dashboard (2 weeks) ⏱️
- [ ] Set up Codecov or similar service
- [ ] Configure branch protection to require coverage uploads
- [ ] Create dashboard with trend lines
- [ ] Weekly reporting of coverage %

---

### P4.2: Benchmark Suite (1-2 weeks) ⏱️
- [ ] Add benchmarks for:
  - Form parsing speed
  - Modal stack operations
  - Validation performance
- [ ] Track performance trends

---

### P4.3: Chaos Testing (2+ weeks) ⏱️
- [ ] Simulate process crashes
- [ ] Simulate terminal resize during modal
- [ ] Simulate slow I/O

---

## TRACKING & METRICS

Use this table to track progress:

| Item | Status | Start Date | End Date | Assignee | PR Link |
|------|--------|-----------|----------|----------|---------|
| P1.1: Coverage Gate | 🔴 TODO | — | — | — | — |
| P1.2: Runner Tests | 🔴 TODO | — | — | — | — |
| P1.3: Model Tests | 🔴 TODO | — | — | — | — |
| P2.1: Security Tests | 🔴 TODO | — | — | — | — |
| P2.2: Validation Tests | 🔴 TODO | — | — | — | — |
| P2.3: Coverage Baseline | 🔴 TODO | — | — | — | — |
| P2.4: Resolve TODOs | 🔴 TODO | — | — | — | — |
| P3.1: Refactor Runner | 🔴 TODO | — | — | — | — |
| P3.2: Refactor Model | 🔴 TODO | — | — | — | — |
| P3.3: Stress Tests | 🔴 TODO | — | — | — | — |
| P4.*: Long-term | 🔴 TODO | — | — | — | — |

**Legend**: 🔴 TODO | 🟠 IN PROGRESS | 🟡 REVIEW | 🟢 DONE

---

## REFERENCES

- Assessment Date: 2025-11-12
- Full Assessment: `docs/qa/CODE-QUALITY-ASSESSMENT-20251112.md`
- Summary: `docs/qa/QUALITY-SUMMARY.md`
- Detailed Findings: `docs/qa/FINDINGS-DETAILED.md`
- Architecture: `docs/architecture/`
