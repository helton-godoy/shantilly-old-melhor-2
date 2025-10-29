# Source Tree

Based on standard Go project layout and the user-provided template.

```plaintext
shantilly/                     # Project Root (replace with actual name later)
├── .github/
│   └── workflows/
│       └── release.yml        # GitHub Actions pipeline for GoReleaser
├── .golangci.yml              # Linter configuration (from template)
├── .goreleaser.yaml           # GoReleaser configuration (from template, needs project_name update)
├── .pre-commit-config.yaml    # Pre-commit hooks (from template)
├── .gitignore                 # Git ignore rules (from template)
├── cmd/
│   └── shantilly/             # Main application package
│       └── main.go            # Entry point: Cobra setup, CLI logic, stdin/file reading
├── internal/
│   ├── config/                # Configuration parsing and validation
│   │   ├── models.go        # Go structs (FormConfig, FieldConfig)
│   │   ├── parser.go        # ConfigParser component (YAML parsing logic)
│   │   └── parser_test.go   # Unit tests for parser
│   ├── tui/                   # TUI rendering and interaction logic
│   │   ├── engine.go        # TUIEngine component (Bubbletea model)
│   │   └── mapper.go        # Logic to map config.FieldConfig -> huh components
│   └── util/                  # Shared utility functions
│       └── errorhandler.go    # ErrorHandler component (stderr output, exit codes)
├── examples/                  # Example YAML form definitions for testing/docs
│   └── basic_form.yaml
├── Makefile                   # Build, Lint, Test scripts (integrates lint.sh logic)
├── go.mod                     # Go module definition (from template, adjusted)
├── go.sum                     # Go module checksums
├── lint.sh                    # Local linting script (from template)
├── LICENSE                    # Project License (MIT from template)
└── README.md                  # Project README (from template, needs customization)

```
