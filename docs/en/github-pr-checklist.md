# Pull Request Checklist - Shantilly

This checklist should be followed **before approving a PR** to `develop` or `main`.

## 1. Scope and purpose

- [ ] PR title is **clear and descriptive**.
- [ ] Description explains **the problem** and **the proposed solution**.
- [ ] The PR has a **reasonable scope** (avoid huge, mixed changes).

## 2. Source branch

- [ ] Branch name follows the convention:
  - `feat/feature-name` for new features.
  - `fix/bug-description` for fixes.
- [ ] Branch was created from `develop` (or from `main` for a critical hotfix).

## 3. Code and architecture

- [ ] Changes respect the **branch strategy** defined in `github-branches-strategy.md`.
- [ ] There are no obvious violations of architectural rules (e.g. avoid `os.Exit` outside `main`, etc.).
- [ ] Function/file/package names are consistent with the rest of the project.
- [ ] There is no leftover commented-out or junk code.

## 4. Tests

- [ ] Relevant local tests were run:
  - [ ] `go test ./...`
  - [ ] Runtime smoke test when applicable:
        `./shantilly runtime --file app.example.yaml`
- [ ] New critical behavior has **automated tests** where it makes sense.
- [ ] The PR does **not break** existing tests.

## 5. Documentation and communication

- [ ] Documentation was updated when needed:
  - [ ] `README.md`
  - [ ] Relevant docs under `docs/pt-br/` and `docs/en/`.
- [ ] Changelog / release notes were updated when the PR affects a release.
- [ ] Code comments only explain what is truly non-obvious.

## 6. CI/CD impact

- [ ] GitHub Actions workflows were reviewed if the PR touches `.github/workflows/`.
- [ ] Confirmed that the PR does not break build, lint or deploy in automated environments.

## 7. Quality review

- [ ] Diff was read by at least **one reviewer** (can be the author in a solo project, but with care).
- [ ] There are no hidden "surprises" (huge files, secrets, etc.).
- [ ] For long-running features, the strategy was considered:
  - Break into multiple smaller PRs when possible.
  - Use feature flags for partially implemented code.

## 8. Specific rules for `main`

Before approving a PR with base `main` (release or hotfix):

- [ ] The PR comes from a **stable** `develop` point (except for very small hotfixes).
- [ ] All tests pass in CI.
- [ ] The planned version tag (e.g. `v0.2.1`) is defined and documented.
- [ ] There is a clear description of the impact for end users.

Following this checklist helps keep the repository **healthy**, predictable and with a clear history of decisions.
