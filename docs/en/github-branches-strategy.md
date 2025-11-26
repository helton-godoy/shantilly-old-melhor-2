# Branch Strategy - Shantilly

## Goals

- Keep the repository **clean and predictable**.
- Ensure only **stable** code reaches `main`.
- Support safe, continuous development on `develop`.
- Use short-lived, focused branches for features and fixes.

## Permanent Branches

### `main`

- Represents the **production / stable releases** state.
- Rules:
  - Only **tested and approved** code is allowed.
  - Changes land in `main` via **release PRs** (e.g. `integration-v2.0`) or very well justified hotfixes.
  - Version tags (e.g. `v0.2.0`) must always point to commits on `main`.

### `develop`

- Is the **continuous development trunk**.
- Always starts from a stable point of `main` (for example, right after a release).
- Rules:
  - New features and fixes branch off from `develop`.
  - `develop` receives merges from `feat/*` and `fix/*` branches via PRs.
  - Periodically, when stable, the state of `develop` is promoted to `main` as a new release.

## Temporary Branches

### `feat/*`

- For **new features**.
- Created from `develop`.
- Examples: `feat/runtime-logging`, `feat/new-component-select`.
- Flow:
  1. `git checkout develop`
  2. `git checkout -b feat/feature-name`
  3. Implement + run local tests.
  4. Open a PR `feat/feature-name` → `develop`.
  5. After merge, **delete** the branch (remote and local).

### `fix/*`

- For **bug fixes and small corrections** (runtime bugs, CI tweaks, etc.).
- Usually created from `develop`.
- Examples: `fix/runtime-panic-empty-config`, `fix/weekly-report-permissions`.
- Flow:
  1. `git checkout develop`
  2. `git checkout -b fix/describe-the-bug`
  3. Fix the issue and test (`go test ./...`, relevant smoke tests).
  4. Open a PR `fix/...` → `develop` (or → `main` in case of a critical production hotfix).
  5. After merge, **delete** the branch.

## Integration Policy

1. **Code review is mandatory** for merges into both `develop` and `main`.
2. **Automated tests** should run on PRs:
   - `go test ./...`
   - Runtime smoke tests (e.g. `./shantilly runtime --file app.example.yaml`).
3. **No direct commits to `main`**, except for extremely well-justified hotfixes.
4. **`develop` is the daily work base**, but should still be kept reasonably stable.

For a detailed review checklist, see:

- [Pull Request Checklist (EN)](github-pr-checklist.md)

## Branch Cleanup

- Legacy branches such as `docs-i18n-and-advanced-site`, `feat/runtime-migration`, `fix/package-lock`, `fix/weekly-report-permissions`, `fix/package-lock2` and `integration-v2.0` were **evaluated**, and their valuable content was **incorporated** into `main` and/or `develop`.
- After integration:
  - Old branches were removed from the remote to reduce noise.
  - History is preserved via commits and PRs (for example, PR #68 for v2.0 integration).

### Practical Rule

- Every `feat/*` or `fix/*` branch must:
  - Be created from `develop`.
  - Be integrated via PR.
  - Be deleted after the merge.

## Releases

- Formal releases are made from `main`, using tags:
  - Example: `v0.2.0` for the release that integrates runtime v2.0 + bilingual documentation.
- It is recommended to maintain a release notes file (e.g. `RELEASES.md`) summarizing key changes per version.

## Long-running features

Instead of keeping huge, long-lived branches (which accumulate conflicts and are hard to review), the preferred strategy is:

- **Break large features into smaller increments**, each delivered in a focused `feat/*` branch.
- Whenever possible, use **feature flags** or configuration toggles so that you can:
  - Merge partially implemented code into `develop`.
  - Keep the functionality disabled by default until it is ready.
- When several large features need to be coordinated:
  - If needed, create a **temporary integration branch** from `develop` (e.g. `integration/epic-1-runtime-v2`).
  - Merge related `feat/*` branches into it.
  - Once stable, promote the result back into `develop` (via PR) and **delete** the integration branch.

Important rules for long-running features:

- Avoid working for weeks on an isolated branch without pulling from `develop`.
- Prefer short cycles: small PRs, fast reviews, continuous feedback.
- Always use PRs to provide historical visibility and a place to document design decisions.

## Branch flow diagram

```mermaid
flowchart LR
    A[main\n(stable releases)] --> B[develop\n(continuous development)]

    B --> C[feat/new-feature]
    C --> B

    B --> D[fix/specific-bug]
    D --> B

    B --> E[integration/epic-X\n(temporary branch)]
    C --> E
    D --> E
    E --> B
```

## Summary

- **`main`**: production / stable releases.
- **`develop`**: continuous development, always based on a stable `main`.
- **`feat/*` and `fix/*`**: short-lived, focused branches, always integrated via PR and removed after merge.
- History is preserved in **commits, PRs and tags**, not in long-lived, unnecessary branches.
