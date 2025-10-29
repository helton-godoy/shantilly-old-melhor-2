# Infrastructure and Deployment

## Infrastructure as Code (IaC)

* **Tool:** `GoReleaser` (`.goreleaser.yaml` from template).
* **Location:** `.goreleaser.yaml` (root directory).
* **Approach:** Defines declarative build, cross-compilation, packaging, and release process.

## Deployment Strategy (Release)

* **Strategy:** GitHub Releases.
* **CI/CD Platform:** GitHub Actions (Tech Stack).
* **Pipeline Configuration:** `.github/workflows/release.yml`.

## Environments

* **N/A:** Not applicable for a CLI. Target environments are user machines (Linux, macOS, Windows - NFR2).

## Promotion Flow (Release)

```mermaid
graph TD
    A[Dev: Push to 'main'] --> B(CI: Run 'golangci-lint' & 'go test')
    B -- Success --> C[Dev: Create Git Tag (e.g., 'v1.0.0')]
    C --> D[Dev: Push Tag to GitHub]
    D -- Trigger (on tag) --> E[GitHub Actions: Run 'goreleaser release']

    subgraph "GoReleaser (in CI)"
        direction TB
        F(1. Static Compile<br/>CGO_ENABLED=0)
        F --> G(2. Cross-Compile<br/>NFR2 Targets)
        G --> H(3. Create Checksums)
        H --> I(4. Package Archives<br/>.zip / .tar.gz)
    end

    E --> F
    I --> J[GitHub Actions: Publish GitHub Release 'v1.0.0' with binaries]

    style J fill:#bbf,stroke:#333,stroke-width:2px
```

## Rollback Strategy

* **Method:** Delete problematic GitHub Release, publish new patch release (e.g., `v1.0.1`) with fix.
* **Triggers:** Critical bug reports from users.
