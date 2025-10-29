# Core Workflows

## Workflow 1: Successful Submission (Happy Path)

Illustrates FR2, FR4, FR6.

```mermaid
sequenceDiagram
    actor User
    participant Shell
    participant CLI(shantilly form)
    participant ConfigParser
    participant TUIEngine(bubbletea + huh)
    participant JSONEncoder

    User->>Shell: Executes script (e.g., shantilly form --file form.yaml)
    Shell->>CLI(shantilly form): 1. Reads YAML (file or stdin)

    activate CLI
    CLI->>ConfigParser: 2. Parse(yaml bytes)
    activate ConfigParser
    ConfigParser-->>CLI(shantilly form): 3. Returns *config.FormConfig
    deactivate ConfigParser

    CLI->>TUIEngine(bubbletea + huh): 4. Run(config)
    activate TUIEngine
    TUIEngine->>User: 5. Renders interactive TUI

    User->>TUIEngine(bubbletea + huh): 6. Fills form
    User->>TUIEngine(bubbletea + huh): 7. Submits form

    TUIEngine-->>CLI(shantilly form): 8. Returns data (map[string]interface{})
    deactivate TUIEngine

    CLI->>JSONEncoder: 9. Encode(data)
    activate JSONEncoder
    JSONEncoder-->>CLI(shantilly form): 10. Returns JSON string
    deactivate JSONEncoder

    CLI->>Shell: 11. Writes JSON string to stdout
    deactivate CLI
    Shell->>User: Displays JSON output
```

## Workflow 2: YAML Parse Error (Error Path)

Illustrates NFR8.

```mermaid
sequenceDiagram
    actor User
    participant Shell
    participant CLI(shantilly form)
    participant ConfigParser
    participant ErrorHandler

    User->>Shell: Executes script (e.g., shantilly form --file bad_form.yaml)
    Shell->>CLI(shantilly form): 1. Reads invalid YAML

    activate CLI
    CLI->>ConfigParser: 2. Parse(yaml bytes)
    activate ConfigParser
    ConfigParser-->>CLI(shantilly form): 3. Returns error (e.g., malformed YAML)
    deactivate ConfigParser

    CLI->>ErrorHandler: 4. Handle(error)
    activate ErrorHandler
    ErrorHandler->>Shell: 5. Writes error message to stderr
    ErrorHandler->>CLI(shantilly form): 6. Signals exit(1) via os.Exit()
    deactivate ErrorHandler

    deactivate CLI # Already exited
    Shell->>User: Displays error message
```

## Workflow 3: User Cancellation (Abort Path - NFR8)

Illustrates NFR8 requirement for clean exit on Esc/Ctrl+C.

```mermaid
sequenceDiagram
    participant User
    participant CLI(shantilly form)
    participant TUIEngine(bubbletea + huh)
    participant ErrorHandler

    User->>CLI(shantilly form): 1. Starts TUI
    activate CLI
    CLI->>TUIEngine(bubbletea + huh): 2. Run(config)
    activate TUIEngine
    TUIEngine->>User: 3. Renders TUI

    User->>TUIEngine(bubbletea + huh): 4. Presses 'Esc' or 'Ctrl+C'

    TUIEngine-->>CLI(shantilly form): 5. Returns error (e.g., bubbletea.QuitMsg or ErrAborted)
    deactivate TUIEngine

    CLI->>ErrorHandler: 6. Handle(ErrAborted or QuitMsg)
    activate ErrorHandler
    ErrorHandler->>Shell: 7. Optionally writes "Cancelled" to stderr
    ErrorHandler->>CLI(shantilly form): 8. Signals exit(2) via os.Exit()
    deactivate ErrorHandler

    deactivate CLI # Already exited
    User->>User: (Terminal is clean, NO JSON output)
```
