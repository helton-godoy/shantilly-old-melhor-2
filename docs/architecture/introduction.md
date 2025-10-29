# Introduction

This document outlines the overall technical architecture for the `shantilly` project. Its primary goal is to serve as the guiding architectural blueprint for AI-driven development, ensuring consistency and adherence to chosen patterns and technologies.

`shantilly` is a monolithic Command Line Interface (CLI) application, written in Go. It operates as a terminal pipeline: receiving a TUI definition in YAML format via `stdin` or `--file`, rendering an interactive TUI (using `charmbracelet/bubbletea` and `charmbracelet/huh`), collecting user input, and upon submission, emitting the collected data as a structured JSON object to `stdout`.

This document focuses on the backend/CLI systems and non-web UI concerns.

## Starter Template or Existing Project

N/A. This is a greenfield project that will not be based on a starter template. The structure will be built from scratch, following standard Go conventions and utilizing the core libraries identified in the PRD (Go, `spf13/cobra`, `charmbracelet/bubbletea`, `charmbracelet/huh`, `charmbracelet/lipgloss`, `gopkg.in/yaml.v3`) and the user-provided project template files.

## Change Log

| Date       | Version | Description                                          | Author              |
|:-----------|:--------|:-----------------------------------------------------|:--------------------|
| 2025-10-23 | 0.1.0   | Initial draft based on PRD v0.2.0.                   | Winston (Architect) |
| 2025-10-23 | 0.2.0   | Incorporated user template files and decisions.      | Winston (Architect) |
| 2025-10-23 | 0.2.1   | Added `--file` input, cancel workflow, wizard logic. | Winston (Architect) |
| 2025-10-23 | 0.3.0   | Added Checklist Results and Next Steps sections.     | Winston (Architect) |
