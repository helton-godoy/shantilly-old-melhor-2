# Project Brief: shantilly

## Executive Summary

shantilly is a modern command-line tool designed to build complex and interactive TUI (Terminal User Interfaces) for shell scripts in a declarative and simplified manner. The project aims to solve the limitations of traditional tools like `dialog` and `whiptail`, offering script developers (Linux, macOS, Windows) a powerful and easy way to create rich UIs, including forms, layouts, mouse support, and themes, distributed as a single static binary. The key value proposition lies in simplifying the creation of modern and portable TUIs for automation and interaction in terminal environments.

## Problem Statement

Shell script developers (Bash, Zsh, PowerShell, etc.) often need user interactions that go beyond simple text prompts or confirmations. Existing standard tools, such as `dialog` and `whiptail`, while functional for basic dialog boxes (text input, simple menus, message boxes), have significant limitations:

- **Limited Components:** Lack native support for modern and complex UI elements (e.g., date pickers, advanced form fields, tabs, multi-panel layouts).

- **Rigid Layout:** Offer little control over the positioning and organization of elements on the screen.

- **Dated Aesthetics:** The generated interfaces generally look outdated and do not align with modern terminal aesthetics.

- **Lack of Interactivity:** Lack native mouse support, limiting interaction to the keyboard.

- **Cross-Platform Inconsistency:** Behavior and appearance can vary significantly between Linux, macOS, and Windows (especially in non-pure VT100 terminals).

- **Complexity for Rich UIs:** Building slightly more complex interfaces requires chaining multiple `dialog`/`whiptail` commands or resorting to complex and fragile script hacks.

The impact of these limitations is that developers often avoid creating rich interactive interfaces in their scripts, resulting in less intuitive and efficient CLI tools. The alternative, using TUI libraries directly (like `ncurses`, `tview`, or even `bubbletea` programmatically), requires knowledge of programming languages (like Go, Python, Rust) and a development effort disproportionate to creating simple interactive scripts. There is a growing need for a tool that fills this gap, allowing for the easy, declarative creation of modern and portable TUIs directly from shell scripts, improving developer productivity and the end-user experience in command-line environments.

## Proposed Solution

The proposed solution is **shantilly**, a CLI tool written in Go that acts as an interpreter for declarative TUI definitions (preferably in YAML). The core concept is to allow shell script developers to define complex interfaces through simple configuration, which `shantilly` then renders and manages interactively in the terminal.

**Key Differentiators:**

- **Declarative Approach:** Focus on YAML (or JSON) to define the UI, separating the *description* from the *implementation*.

- **Modern Charm Ecosystem:** Leverages libraries like `bubbletea`, `lipgloss`, `huh`, `bubbles`, `bubbletea-overlay`, and `bubblezone` to offer rich components, flexible layouts, advanced styling, and mouse support.

- **Single Portable Binary:** Statically compiled for Linux, macOS, and Windows, ensuring easy distribution and consistency across platforms.

- **Shell Integration:** Designed to read configuration from `stdin` (allowing `here-docs`) and return structured data (JSON) to `stdout`, seamlessly integrating into script pipelines.

**Why it will succeed:** `shantilly` will succeed by combining the power and aesthetics of modern TUI libraries with the ease of use of a declarative CLI tool. It abstracts the complexity of TUI programming in Go, making advanced features accessible directly from shell scripts, something that neither `dialog`/`whiptail` (too simple) nor the direct use of TUI libraries (too complex for scripts) currently offer.

**High-Level Vision:** To make `shantilly` the standard tool for creating any interactive TUI within shell scripts, from simple forms to complex dashboards, with an excellent user and developer experience. In the future, it may include **native SSH support (via `wish`) for secure and direct access to administrative TUIs**, as well as internationalization (i18n).

## Target Users

### Primary User Segment: Script Developers / CLI Power Users

Profile: Developers (Backend, DevOps, Frontend using CLI), System Administrators, QA Engineers, and advanced users who utilize shell scripts (Bash, Zsh, PowerShell, etc.) for automation, internal tools, system management, build/deploy processes, or creating interactive CLI applications. They primarily work in terminal environments, including SSH sessions, on Linux, macOS, or Windows systems.

Current Behaviors: Write scripts to automate repetitive or complex tasks. For user interaction, they use read for simple prompts, select for basic menus, or resort to dialog/whiptail for slightly richer interfaces but feel their limitations. They may have experience with other programming languages but prefer the convenience of shell scripting for certain tasks. Some may have tried using programmatic TUI libraries but found the effort excessive for the context of a script.

Needs and Pain Points:

- Need to create more intuitive and modern interfaces (forms with validation, selectors, organized layouts, visual feedback) within scripts.

- Difficulty building complex UIs with `dialog`/`whiptail` due to lack of components, limited layout control, and dated aesthetics.

- Frustration with the lack of mouse support and visual/behavioral inconsistencies across different terminals and operating systems.

- Desire for a declarative solution that doesn't require learning a new programming language or complex framework just to add a UI to a script.

- Need for easy integration with existing script logic (receiving input data, returning results in a structured format).

- **For System Administrators/DevOps:** Need to create secure and remotely accessible administrative interfaces without exposing a full shell or relying on complex web interfaces.

**Goals:**

- Increase the usability and clarity of their CLI tools and automation scripts.

- Quickly create interactive and visually appealing TUIs without leaving the scripting paradigm.

- Ensure their interfaces work consistently across different terminal environments and OS.

- Simplify the collection of complex user data within a script.

- Improve the overall experience for both themselves (when developing) and the end-users of their scripts.

- **Provide a secure and robust way for remote system administration through dedicated TUIs via SSH.**

## Goals & Success Metrics

### Business Objectives

- **Adoption:** Achieve X stars on GitHub and Y binary downloads within the first 6 months after the v1.0 release, indicating community acceptance.

- **Viability:** Release version 1.0 (MVP) with robust support for the `form` subcommand and basic shell integration (stdin/stdout) within 3 months.

- **Community:** Establish a public repository with clear documentation, examples, and a defined contribution process by the v1.0 release.

### User Success Metrics

- **Ease of Use:** Developers should be able to create a basic TUI form (3-5 fields) by reading the documentation and examples in less than 15 minutes. (Measured via feedback and examples)

- **Expressive Power:** Allow the creation of significantly richer interfaces (layouts, multiple components) than `dialog`/`whiptail`, as demonstrated by documentation examples.

- **Integration:** Scripts should be able to consume `shantilly`'s JSON output reliably and easily.

### Key Performance Indicators (KPIs)

- **GitHub Stars/Downloads:** Indicator of initial interest and adoption. (Target: See Business Objectives)

- **Issues Opened/Closed:** Track stability and responsiveness to bugs and feature requests. (Target: Maintain response time < 48h for new issues)

- **Community Contributions:** Number of PRs, issues reported by external users as an indicator of engagement. (Target: >5 external contributions in the first 6 months)

- **Number of Supported Components (MVP+):** Track the growth of functionality over time. (MVP Target: Robust support for forms via `huh`)

## MVP Scope

### Core Features (Must Have)

- **`form` subcommand:** Central implementation for rendering TUI forms.
  
  - **Rationale:** Core functionality, addressing the primary need for interactive data input.

- **Read YAML Configuration via `stdin`:** Allow defining forms via `here-docs` or pipes.
  
  - **Rationale:** Essential for declarative integration with shell scripts.

- **Support for Basic `huh` Field Types:** Include `Input`, `Textarea`, `Select`, `MultiSelect`, `Confirm`, `Note`. (Excluding `FilePicker` and complex custom fields in MVP).
  
  - **Rationale:** Covers most form use cases without adding excessive initial complexity.

- **JSON Data Output via `stdout`:** Return user-filled data in a structured format.
  
  - **Rationale:** Allows easy consumption by scripts calling `shantilly`.

- **Basic TUI Engine (`bubbletea` + `lipgloss`):** Fundamental structure for rendering, basic event handling (exit, resize), and initial styling.
  
  - **Rationale:** Necessary base for any TUI interface.

- **CLI Structure (`cobra`):** Setup of the root command and the `form` subcommand.
  
  - **Rationale:** Organization of the command-line tool.

- **Cross-Platform Static Compilation:** Builds for Linux, macOS, Windows.
  
  - **Rationale:** Meets the portability requirement.

### Out of Scope for MVP

- `layout`, `menu`, `dialog`, etc. subcommands.

- Mouse support via `bubblezone`.

- Floating/modal windows (`bubbletea-overlay` or similar).

- Animations (`harmonica`).

- Charts (`ntcharts`).

- Advanced themes (beyond basic styling via `lipgloss`).

- SSH support (`wish`).

- Internationalization (i18n).

- Reading configuration via complex flags (beyond `--title` and perhaps `-f file.yml`).

- Complex integration tests (focus on unit tests for parsers and core logic).

### MVP Success Criteria

The MVP will be considered a success if a script developer can:

1. Define a form with 3-5 supported field types in YAML.

2. Run `cat form.yml | shantilly form --title "My Form"`.

3. Interact with the rendered TUI form using the keyboard.

4. Submit the form.

5. Receive the correct JSON output on `stdout`.

6. The process works consistently on the three target operating systems (Linux, macOS, Windows).

## Post-MVP Vision

### Phase 2 Features

After the MVP launch (focusing on `form`), the next phase will concentrate on expanding UI building capabilities:

- **`layout` subcommand:** Introduce support for multi-panel layouts (sidebar/content, columns, rows) using `bubbleboxer` and `lipgloss`.

- **`menu`/`list` subcommand:** Implement interactive lists and selection menus based on `bubbles/list`.

- **Basic Mouse Support:** Integrate `bubblezone` to allow clicks on buttons and selection in lists/menus.

- **Simple `dialog` subcommand:** Use `bubbletea-overlay` to create basic modal dialog boxes (confirmation, alert, simple prompt).

### Long-term Vision

The long-term vision is to make `shantilly` a complete and robust declarative TUI tool:

- **Administration via SSH:** Implement direct and secure access to `shantilly` interfaces via SSH, utilizing `wish`, ideal for remote server management (like the NAS project).

- **Advanced Components:** Add support for tabs, tables, charts (`ntcharts`), date pickers (`datepicker`), etc.

- **Window Management:** Evolve the dialog functionality into a more robust system, inspired by `winman`, allowing multiple floating windows and focus management.

- **Theming:** Allow extensive appearance customization via theme files (`lipgloss`).

- **Animations:** Integrate `harmonica` to add fluidity to interactions and transitions.

### Expansion Opportunities

- **Internationalization (i18n):** Add support for multiple languages in the generated interfaces.

- **Extensibility:** Possibly allow user-defined custom components (though this increases complexity).

- **Tool Integration:** Explore integrations with other CLI tools or automation platforms.

## Technical Considerations

### Platform Requirements

- **Target Platforms:** Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64).

- **Browser/OS Support:** N/A (terminal application). Requires a modern terminal with good color support (TrueColor recommended) and ANSI/VT100 sequences.

- **Performance Requirements:** Must have fast startup (<500ms) and be responsive to user interaction, even in SSH sessions. The final binary should be lightweight (<20MB preferably).

### Technology Preferences

- **Language:** Go (latest stable version).

- **TUI Framework:** `charmbracelet/bubbletea` and its ecosystem (`lipgloss`, `bubbles`, `huh`, `bubblezone`, `harmonica`, `bubbletea-overlay`, `wish`, `log`).

- **CLI Framework:** `spf13/cobra`.

- **Configuration:** `spf13/viper` (for flags), `gopkg.in/yaml.v3` (for `stdin` parsing).

- **Filesystem (if needed):** `spf13/afero` for abstraction and testing.

- **Terminal Libraries:** `mattn/go-tty`, `mattn/go-isatty`, `mattn/go-colorable` (primarily for Windows support).

### Architecture Considerations

- **Repository Structure:** Simple monorepo with clear separation between `cmd/` (Cobra CLI logic) and `internal/` (Bubbletea TUI logic, parsing, etc.).

- **Service Architecture:** Monolithic CLI application. The SSH mode (`wish`) will function as a dedicated TUI server.

- **Integration Requirements:** Primary integration with the shell via `stdin` (YAML config) and `stdout` (JSON output). Future integration via SSH server.

- **Security/Compliance:** For SSH mode, leverage the inherent security of the SSH protocol and its authentication options. No specific compliance requirements beyond standard security best practices.

## Constraints & Assumptions

### Constraints

- **Budget:** Open-source project, no formal budget. Development depends on available time.

- **Timeline:** Target of 3 months for MVP (v1.0) release.

- **Resources:** Initial development primarily carried out by one contributor (implied).

- **Technical:** Limited by the capabilities of the Charmbracelet TUI libraries and the inherent limitations of terminal environments (e.g., rendering performance, font/character support). Compatibility with older or less capable terminals may be restricted.

### Key Assumptions

- The Charmbracelet ecosystem libraries (`bubbletea`, `lipgloss`, `huh`, etc.) are stable, performant, and suitable for building the desired functionality.

- YAML is an appropriate and user-friendly configuration format for shell script developers.

- Cross-platform static compilation in Go for Linux, macOS, and Windows is feasible and will produce portable, functional binaries.

- Target users have modern terminals with adequate support for colors and ANSI/VT100 sequences.

- Integration via `stdin`/`stdout` (YAML/JSON) is an effective and sufficient method for interaction with shell scripts in the MVP.

## Risks & Open Questions

### Key Risks

- **Charm Ecosystem Dependency:** The stability and evolution of Charmbracelet libraries directly impact `shantilly`. API changes or bugs in these dependencies may require adaptation effort. (Impact: Medium)

- **Declarative Interpretation Complexity:** Translating flexible YAML configurations into interactive TUI components and complex layouts can be challenging and prone to bugs, especially with nesting and interactions. (Impact: High)

- **Terminal Compatibility:** Ensuring consistent behavior and appearance across a wide range of terminals and emulators (especially on Windows), resolutions, window sizes, and fonts is notoriously difficult and a common point of failure in previous attempts. (Impact: Medium-High)

- **Performance:** Very complex TUIs rendered by `bubbletea` might exhibit performance issues on slower terminals or over high-latency SSH connections, especially with many components or animations. (Impact: Low to Medium)

- **Post-MVP Scope:** Implementing advanced features like custom window management and animations might prove more complex than anticipated. (Impact: Medium)

### Open Questions

- What will be the best strategy for managing focus between multiple TUI components in complex layouts (beyond what `bubbletea` offers natively)?

- How to handle YAML configuration parsing errors elegantly and informatively for the user?

- What level of theme customization will be supported initially (Post-MVP)?

- What will the exact API for defining multi-panel layouts in YAML be (Post-MVP)?

- Are there common `dialog`/`whiptail` use cases that current Charm components don't adequately cover?

### Areas Needing Further Research

- **Robust Cross-Terminal Rendering Strategies:** Investigate techniques and best practices used by other successful TUI tools (including Charm's own) to handle terminal inconsistencies, fonts, and resizing.

- **`bubbletea` Performance with Complex Layouts:** Look into benchmarks or existing examples of `bubbletea` applications with many sub-components and dynamic layouts.

- **Alternatives/Best Practices for Window Management in `bubbletea`:** Research if new libraries or patterns have emerged in the Charm community since `winman` (for `tview`) and `bubbletea-overlay`.

- **Special Character/Unicode Compatibility:** Validate how different terminals/OS handle borders, icons, and other characters used in styling.

## Appendices

### A. Research Summary

(Empty for now)

### B. Stakeholder Input

(Empty for now)

### C. References

- **Shantilly Repository:** [helton-godoy/shantilly](https://www.google.com/search?q=https://github.com/helton-godoy/shantilly "null")

- **Charmbracelet Ecosystem:**
  
  - [Charmbracelet on GitHub](https://github.com/charmbracelet "null")
  
  - [Bubble Tea](https://github.com/charmbracelet/bubbletea "null") (TUI Framework)
  
  - [Lipgloss](https://github.com/charmbracelet/lipgloss "null") (Styling)
  
  - [Huh](https://github.com/charmbracelet/huh "null") (Forms)
  
  - [Bubbles](https://github.com/charmbracelet/bubbles "null") (TUI Components)
  
  - [Bubblezone](https://github.com/lrstanley/bubblezone "null") (Mouse Support)
  
  - [Harmonica](https://github.com/charmbracelet/harmonica "null") (Animations)
  
  - [Wish](https://github.com/charmbracelet/wish "null") (SSH Server)

- **Other Mentioned Libraries:**
  
  - [Bubbletea Overlay](https://github.com/rmhubbert/bubbletea-overlay "null") (Simple modal alternative)
  
  - [Bubble Boxer](https://github.com/treilik/bubbleboxer "null") (Layouts)
  
  - [Winman (for tview)](https://github.com/epiclabs-io/winman "null") (Window management inspiration)
  
  - [tview](https://github.com/rivo/tview "null") (Alternative TUI framework)
  
  - [Cobra](https://github.com/spf13/cobra "null") (CLI Framework)
  
  - [Viper](https://github.com/spf13/viper "null") (Configuration)
  
  - [Afero](https://github.com/spf13/afero "null") (FS Abstraction)

- **Traditional Tools:**
  
  - `dialog`
  
  - `whiptail`

## Next Steps

### Immediate Actions

1. Formally review and approve this Project Brief.

2. Begin Phase 0 of the Implementation Plan: Go project setup, CLI structure with Cobra.

3. Prepare handoff to the Product Manager (PM) to start PRD creation.

### PM Handoff

This Project Brief provides the complete context for the **shantilly** project. The next step is to engage the Product Manager (PM), John, to begin creating the Product Requirements Document (PRD) using the `prd-tmpl.yaml` template. Please review this brief thoroughly and collaborate with the user to create the PRD section by section, as indicated by the template, requesting necessary clarifications or suggesting improvements based on this brief. Initial focus on the MVP scope.
