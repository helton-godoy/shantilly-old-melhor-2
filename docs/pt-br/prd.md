# shantilly Product Requirements Document (PRD)

## Goals and Background Context

### Goals

* Simplify the creation of interactive and modern TUI interfaces for shell scripts [cite: shantilly Product Requirements Document (PRD).md].
* Offer a declarative (YAML) and more powerful alternative to `dialog` and `whiptail` [cite: shantilly Product Requirements Document (PRD).md].
* Ensure portability and consistency across Linux, macOS, and Windows via a single static binary [cite: shantilly Product Requirements Document (PRD).md].
* Facilitate integration with shell script pipelines (read `stdin` YAML, write `stdout` JSON) [cite: shantilly Product Requirements Document (PRD).md].
* Achieve broad adoption by the script developer community (Bash, Zsh, PowerShell, etc.) [cite: Project Brief_ shantilly.md].
* Deliver a robust and viable MVP (v1.0), focused on the `form` subcommand, within 3 months [cite: Project Brief_ shantilly.md].
* Enable the creation of significantly richer interfaces (layouts, rich components) than traditional tools [cite: Project Brief_ shantilly.md].

### Background Context

Shell script developers needing complex user interactions are currently underserved. Traditional tools like `dialog` and `whiptail` are functional but severely limited in components, layout control, and aesthetics, lacking mouse support [cite: Project Brief_shantilly.md, shantilly Product Requirements Document (PRD).md]. The alternative, using programmatic TUI libraries (like Bubbletea), requires knowledge of languages like Go or Python, adding disproportionate complexity for script-based automation [cite: Project Brief_ shantilly.md, shantilly Product Requirements Document (PRD).md].

shantilly fills this gap. It's a portable CLI tool (single static binary) enabling the *declarative creation* of rich, modern TUIs [cite: Project Brief_shantilly.md]. The user defines the interface (e.g., a form) in YAML, passes it to `shantilly` via `stdin`, and the tool renders the interactive TUI [cite: Project Brief_ shantilly.md]. Collected data is then returned as structured JSON on `stdout`, allowing easy integration into script pipelines [cite: Project Brief_shantilly.md]. This PRD focuses on the MVP (Minimum Viable Product) scope to validate the core functionality [cite: Project Brief_ shantilly.md].

### Change Log

| Date       | Version | Description                                                                                                    | Author     |
| :--------- | :------ | :------------------------------------------------------------------------------------------------------------- | :--------- |
| 2025-10-22 | 0.1.0   | Initial PRD draft based on Project Brief.                                                                      | John (PM)  |
| 2025-10-23 | 0.1.1   | Added UI/UX section and refined MVP layout.                                                                    | John (PM)  |
| 2025-10-23 | 0.1.2   | Added Technical Assumptions section.                                                                           | John (PM)  |
| 2025-10-23 | 0.1.3   | Added Epic List (MVP).                                                                                         | John (PM)  |
| 2025-10-23 | 0.1.4   | Added Epic 1 Details (MVP) with Stories.                                                                       | John (PM)  |
| 2025-10-23 | 0.1.5   | Completed PM Checklist and Next Steps section.                                                                 | John (PM)  |
| 2025-10-23 | 0.2.0   | Implemented PM Checklist recommendations (YAML Structure and Error Handling - NFR8). Updated Architect prompt. | John (PM)  |
| 2025-10-27 | 0.3.0   | Added preventive epics (3-7) for future roadmap planning and process improvement.                              | Sarah (PO) |
| 2025-11-08 | 2.0.0   | Pivô para Runtime TUI Declarativo v2.0, incluindo layouts, componentes múltiplos, lógica de eventos.           | John (PM)  |

## Requirements

### Functional (FRs) - O Runtime TUI

Estes requisitos definem o nosso novo MVP: a "Fundação Genérica".

* **FR1 (Layout):** O `shantilly` DEVE analisar e renderizar uma estrutura de layout hierárquica definida em YAML, usando os tipos `type: column`, `type: row`, e `type: box`.

* **FR2 (Estilo/Flex):** O layout DEVE suportar propriedades de dimensionamento como `height: <int>`, `width: 'N%'`, e `flex: <int>` para controlar o espaço.

* **FR3 (Componentes Embutidos):** O `shantilly` DEVE suportar a definição de componentes de UI diretamente dentro de um `box` usando a chave `component:` (o foco do nosso MVP).

* **FR4 (Componente: `viewport`):** DEVE suportar `component: { type: viewport }`, capaz de exibir `source: { type: static, content: "..." }` (incluindo markdown) e `source: { type: command, exec: "..." }` (para streaming de stdout).

* **FR5 (Componente: `list`):** DEVE suportar `component: { type: list }`, com um `id:` de grupo e `items:` (cada um com `id:` e `text`), e DEVE emitir um evento `list_id:select`.

* **FR6 (Componente: `buttongroup`):** DEVE suportar `component: { type: buttongroup }`, com um `id:` de grupo, `items:` (com `id`, `label`, `role`), e DEVE emitir um evento `buttongroup_id:press`.

* **FR7 (Componente: `form`):** DEVE suportar `component: { type: form }`, que contém `fields:` (usando a sintaxe `huh` já validada no v1.0) e `actions:`. DEVE emitir um evento `form_id:submit` contendo o *payload* de dados do formulário.

* **FR8 (Lógica: `on:`):** O `shantilly` DEVE analisar um bloco `on:` na raiz do YAML para definir a lógica de automação.

* **FR9 (Ação: `script`):** O bloco `on:` DEVE suportar o *runner* de fundação: `run: { script: "/path/to/script.sh" }`.

* **FR10 (Fluxo de Dados):** O *runner* `script:` DEVE suportar duas chaves para passagem de dados:
  1. **`args: []string`**: Uma lista de *strings* que serão passadas como argumentos de linha de comando para o script, com suporte para *templates* (ex: `{{ form.field_name }}`).
  2. **`stdin: any`**: Um objeto (ex: `{{ form }}`) que o `shantilly` irá serializar como JSON e passar para o `stdin` do *script*.

* **FR11 (Ciclo de Vida do Target):** O bloco `run:` DEVE suportar uma chave `update_target: "id_do_viewport"`. Se um novo evento `run:` for disparado para o *mesmo* `update_target`, o `shantilly` DEVE primeiro **terminar (enviar `SIGTERM`)** o processo anterior antes de iniciar o novo.

### Non-Functional (NFRs) - O Runtime TUI

* **NFR1 (Fundação v1.0):** Todos os NFRs do PRD v1.0 permanecem válidos: binário estático único, cross-platform (Linux, macOS, Windows), escrito em Go, arranque rápido (<500ms), e gestão de erros com `stderr` e códigos de saída não-zero.

* **NFR2 (Layout Fluido):** O motor de layout (`column`/`row`/`box`) DEVE responder a mensagens de redimensionamento do terminal (`tea.WindowSizeMsg`) e re-calcular o layout de forma fluida e "flicker-free" (sem piscar).

* **NFR3 (Precedência de Conteúdo):** O `shantilly` DEVE seguir a "Lógica de Precedência Unificada" para conteúdo de componentes (1º: CLI `--set`, 2º: YAML `model:`/`component:`, 3º: Vazio).

* **NFR4 (Descoberta Preditiva - Ansible):** Para a Fase 2/3, o `playbook_explorer` DEVE filtrar "ruído" (pastas `roles/`, `tasks/`) e o `inventory_explorer` DEVE usar `ansible-inventory` como "Oráculo".

* **NFR5 (Ficheiro-Sombra):** O ficheiro de catálogo (`.shantilly.yml`) DEVE ser opcional e usado apenas para *refinar* a descoberta automática, não sendo obrigatório.

## User Interface Design Goals

### Overall UX (User Experience) Vision

A UX deve ser a de um **"Runtime TUI Declarativo"**. A interface não é mais um formulário linear único, mas sim um *dashboard* composto, definido inteiramente pelo YAML. A experiência deve ser semelhante ao Appsmith: limpa, responsiva (ao terminal) e orientada a componentes.

### Key Interaction Paradigms

* **Orientada a Eventos (Nova):** A interação principal não é linear. O utilizador seleciona itens em listas (`list:select`) ou pressiona botões (`buttongroup:press`), que disparam ações no bloco `on:`.

* **Foco no Teclado (Mantido):** A navegação DEVE continuar a ser primariamente baseada no teclado (Tab, Setas, Enter).

* **Feedback Imediato (Mantido):** O componente focado DEVE ser claramente destacado. O `update_target` (FR11) DEVE exibir o *output* de comandos em tempo real.

* **Gestão de Foco Global (Nova):** A UI DEVE ter um mecanismo claro para indicar qual painel/componente (ex: `sidebar` vs `content`) está "em foco", e DEVE fornecer navegação intuitiva *entre* painéis (ex: Ctrl+Tab).

* **Ligação de Dados (Nova):** A UI DEVE ser reativa. Componentes (ex: um `viewport` estático) DEVEM ser capazes de exibir dados de outros componentes (ex: `Olá, {{ form.username }}`).

### Core Screens and Views

Não há ecrãs "pré-definidos". Os ecrãs são *definidos dinamicamente* pelo utilizador através do **`layout` YAML** (FR1). A UI é uma composição de `type: column`, `type: row`, e `type: box`.

### Alignment and Layout (Nova Visão)

O layout linear do MVP v1.0 está obsoleto. O novo requisito é:

* O `shantilly` DEVE renderizar com precisão o layout `column`/`row` definido pelo utilizador.

* O `shantilly` DEVE respeitar as propriedades de dimensionamento (`height`, `width`, `flex`) para distribuir o espaço.

* O `shantilly` DEVE responder a mensagens de redimensionamento do terminal (`tea.WindowSizeMsg`) e re-calcular o layout fluido e "flicker-free" (NFR2, Meta de UI Refinada).

### Accessibility, Branding, Plataformas Alvo

Estes requisitos permanecem os mesmos do PRD v1.0 (WCAG AA, estética Charmbracelet, terminais modernos em Linux/macOS/Windows).

## Technical Assumptions

### Repository Structure: Monorepo

The project will be structured as a simple Go monorepo. This facilitates initial organization and potential future addition of subcommands or related libraries [cite: Project Brief_ shantilly.md, Relatório Técnico do Brainstorming Shantilly.md]. The structure will follow standard Go conventions, with CLI logic in `cmd/shantilly` and TUI/parsing logic in `internal/`.

### Service Architecture: Monolithic CLI Application

For the MVP, shantilly is a monolithic CLI application that executes, processes input, and terminates [cite: Project Brief_ shantilly.md]. The architecture will be modular internally (separating CLI and TUI), but there will be no microservices or network communication [cite: Relatório Técnico do Brainstorming Shantilly.md]. SSH server mode is a Post-MVP vision (NFR7) [cite: docs/prd.md].

### Testing Requirements: Focus on Unit Tests

Given the MVP time constraints and the challenges of testing TUIs [cite: Project Brief_shantilly.md], the primary focus will be on robust unit tests for the YAML parsing logic and internal business logic [cite: Relatório Técnico do Brainstorming Shantilly.md]. UI integration tests will be limited to manual smoke tests on the main platforms [cite: Project Brief_ shantilly.md].

### Estrutura YAML Esperada (A Nova Fonte da Verdade)

* **Assunção (Nova):** A estrutura YAML do v1.0 (lista simples de `fields:`) está obsoleta. A nova assunção de arquitetura é o YAML "Appsmith-style" que definimos, composto por **Layout**, **Componentes** e **Lógica**:

```yaml
# 1. LAYOUT (Define o "onde")

type: column

items:

- type: row

flex: 1

items:

- type: box

id: "sidebar"

width: "30%"

# 2. COMPONENTE (Define o "o quê")

component:

type: list

id: "menu"

items:

- { id: "users", text: "Gerir Utilizadores" }

- type: box

id: "content"

width: "70%"

component:

type: form

id: "user_form"

fields:

- { name: "username", label: "Nome", type: "input" }

actions:

type: buttongroup

items:

- { id: "submit", label: "Criar", role: "primary" }



# 3. LÓGICA (Define o "como")

on:

- event: "user_form:submit"

run:

script: "/opt/scripts/create_user.sh"

# Nomenclatura refinada (Opção C)

args:

- "--mode=production"

stdin: "{{ form }}" # Passa o payload JSON para o stdin
```

### Assunções Técnicas Adicionais

* **Pilha de Tecnologias (Mantida e Validada):** Go (1.24.2+), `spf13/cobra`, `charmbracelet/bubbletea`, `charmbracelet/lipgloss`, `charmbracelet/huh`, `gopkg.in/yaml.v3`.

* **Bibliotecas Relevantes (Nova):** A arquitetura dependerá de `charmbracelet/bubbles` (para `list`, `viewport`), `charmbracelet/glamour` (para markdown), `rmhubbert/bubbletea-overlay` (para modais Fase 2), e inspiração de `76creates/stickers` (para layout).

* **Build e Distribuição (Mantido):** `GoReleaser` para binários estáticos cross-platform.

## Epic List (Roadmap v2.0)

Este *roadmap* substitui a lista de épicos do PRD v1.0.

* **Épico 1: Fundação do Runtime TUI (Genérico)**

* **Meta:** Construir o motor central: o layout (`column`/`row`/`box`), os componentes essenciais (`list`, `viewport`, `form`, `buttongroup`), e a lógica de eventos (`on:`, `run: { script: ... }`). (Este épico absorve todo o trabalho já concluído no v1.0).

* **Épico 2: O Runner Especialista (Ansible Fase 2)**

* **Meta:** Implementar o *runner* de conveniência `run: { ansible_playbook: ... }`, focando na gestão de `vars:` e no popup modal `ask_vault_pass: true`.

* **Épico 3: O Runtime Preditivo (Ansible Fase 3)**

* **Meta:** Implementar os componentes `playbook_explorer` (Magia 1: descobrir playbooks) e `inventory_explorer` (Magia 2: descobrir inventário).

* **Épico 4: Administração SSH (Visão de Longo Prazo)**

* **Meta:** Integrar o `charmbracelet/wish` para servir o Runtime TUI sobre SSH.

## Epic 1: Fundação do Runtime TUI (Genérico)

**Meta do Épico:** Construir o motor central do `shantilly`: o motor de layout (`column`/`row`/`box`), os componentes essenciais de dashboard (`list`, `viewport`, `form`, `buttongroup`), e a lógica de eventos (`on:`, `run: { script: ... }`). Este épico irá refatorar o trabalho concluído do v1.0 para que ele funcione como o componente `type: form` dentro deste novo runtime.

### Estória 1.1: O Motor de Layout (Renderização)

**Como um** SysAdmin, **Eu quero** definir um layout TUI usando `column`, `row`, e `box` no meu YAML, **Para que** eu possa criar dashboards complexos e organizados.

#### Critérios de Aceitação – Estória 1.1

1. O parser DEVE suportar as chaves `type: column`, `type: row`, e `type: box` (FR1).
2. O motor de renderização (`lipgloss`) DEVE respeitar as propriedades `height: <int>`, `width: 'N%'`, e `flex: <int>` (FR2).
3. O layout DEVE recalcular-se fluidamente (sem piscar) ao receber uma mensagem de redimensionamento (`tea.WindowSizeMsg`) (NFR2).
4. Um `box` DEVE renderizar o seu `component: { type: static, content: "..." }` (para testes de layout).

### Estória 1.2: O Motor de Lógica (Eventos e Ações)

**Como um** SysAdmin, **Eu quero** que a minha UI TUI possa "ouvir" eventos e executar ações (`scripts`) em resposta, **Para que** o meu dashboard seja interativo e possa orquestrar automações.

#### Critérios de Aceitação – Estória 1.2

1. O `shantilly` DEVE analisar um bloco `on:` na raiz do YAML (FR8).
2. O motor DEVE suportar o *runner* de fundação: `run: { script: "/path/to/script.sh" }` (FR9).
3. O `run:` DEVE suportar `update_target: "id_do_viewport"`, direcionando o `stdout` do script para o *viewport* alvo (FR11).
4. O motor DEVE garantir que, se um novo script for direcionado para um `update_target` já ocupado, o script anterior seja terminado (SIGTERM) antes de o novo começar (Refinamento FR11).

### Estória 1.3: Componentes Essenciais de Display (List, Viewport, Button)

**Como um** SysAdmin, **Eu quero** usar componentes de `list` (para menus), `viewport` (para saída de log) e `buttongroup` (para ações), **Para que** eu possa construir um dashboard funcional.

#### Critérios de Aceitação – Estória 1.3

1. DEVE implementar `component: { type: viewport }` (FR4), incluindo `source: { type: command, exec: "..." }` (ex: `tail -f`) e `content_type: markdown`.
2. DEVE implementar `component: { type: list }` (FR5), que emite um evento `list_id:select` quando um item é selecionado.
3. DEVE implementar `component: { type: buttongroup }` (FR6), que emite um evento `buttongroup_id:press` (com o `item.id`) quando um botão é pressionado.
4. A navegação por teclado DEVE permitir "saltar" entre estes novos painéis/componentes (Refinamento de Meta de UI).

### Estória 1.4: Integração do Componente `form` (Absorção do v1.0)

**Como um** SysAdmin, **Eu quero** usar o `type: form` (que já construímos no v1.0) como um componente *dentro* do meu novo layout, **Para que** eu possa coletar dados de forma organizada.

#### Critérios de Aceitação – Estória 1.4

1. Refatorar o código dos Épicos 1 e 2 (v1.0) para que funcione como um `component: { type: form }` (FR7).
2. O `form` DEVE renderizar e funcionar corretamente quando colocado dentro de um `box` do layout.
3. Quando a ação `id: "submit"` do formulário for pressionada, o componente `form` DEVE emitir um evento `form_id:submit`.
4. O *payload* do evento `form_id:submit` DEVE conter o JSON de dados do formulário (o output do v1.0).

### Estória 1.5: O Fluxo de Dados (Args & Stdin)

**Como um** SysAdmin, **Eu quero** passar os dados coletados no meu `form` (ou a seleção de uma `list`) para os meus `scripts` de forma robusta, **Para que** a minha automação possa usar a entrada do utilizador.

#### Critérios de Aceitação – Estória 1.5

1. O *runner* `script:` (FR9) DEVE suportar a chave `args: []string`, que passa argumentos "templatados" para a linha de comando do script (Refinamento FR10 / Opção C).
2. O *runner* `script:` (FR9) DEVE suportar a chave `stdin: any`, que serializa o valor (ex: `{{ form }}`) como JSON e o passa para o `stdin` do script (Refinamento FR10 / Opção C).
3. O motor de templates DEVE suportar "binding" de dados (ex: `{{ form.field_name }}`, `{{ component.menu.selected_id }}`) (Refinamento de Meta de UI).
4. Deve existir um exemplo de script (Bash ou PowerShell) que leia dados tanto de `args:` como de `stdin:` (via `jq` ou similar).

## Checklist Results Report

[[LLM: PM CHECKLIST EXECUTION INSTRUCTIONS (`pm-checklist.md`)

**Project:** shantilly (Greenfield, CLI/TUI - No Web/Mobile UI)
**Document Under Review:** `docs/prd.md` (current version in Canvas)

**Mode:** Comprehensive (YOLO) - Analyze all and present final report.

**Context:**

* PRD derived from detailed Project Brief [cite: Project Brief\_ shantilly.md].
* Scope strictly MVP: `form` subcommand, YAML stdin -\> TUI -\> JSON stdout [cite: docs/prd.md].
* Tech stack is Go + Charmbracelet [cite: docs/prd.md].
* Project has NO web/mobile UI, skip related checklist sections.

**Process:**

1. Read EACH item in `pm-checklist.md` [cite: team-fullstack.txt].
2. Verify if the item is covered in the current `docs/prd.md`.
3. Evaluate the *quality* and *completeness* of coverage.
4. Mark each item as ✅ PASS, ❌ FAIL, ⚠️ PARTIAL, or N/A.
5. For ❌ FAIL or ⚠️ PARTIAL, note *why*.
6. Calculate pass rates per section.
7. Synthesize results in the report template below. Be specific.
    ]]

### Executive Summary

* **Overall PRD Completeness:** 100%
* **MVP Scope Appropriateness:** Ideal
* **Architecture Readiness:** Ready
* **Critical Gaps or Concerns:** No critical gaps remaining.

### Category Analysis

| Category                         | Status | Critical Issues                     |
| :------------------------------- | :----- | :---------------------------------- |
| 1. Problem Definition & Context  | ✅ PASS |                                     |
| 2. MVP Scope Definition          | ✅ PASS |                                     |
| 3. User Experience Requirements  | ✅ PASS | (Simplified for TUI)                |
| 4. Functional Requirements       | ✅ PASS |                                     |
| 5. Non-Functional Requirements   | ✅ PASS | (NFR8 added for errors)             |
| 6. Epic & Story Structure        | ✅ PASS |                                     |
| 7. Technical Guidance            | ✅ PASS | (YAML added to Assumptions)         |
| 8. Cross-Functional Requirements | ✅ PASS | (YAML Structure & Errors addressed) |
| 9. Clarity & Communication       | ✅ PASS |                                     |

### Top Issues by Priority

* **BLOCKERS:** None.
* **HIGH:** None (Previous recommendations implemented).
* **MEDIUM:** None.
* **LOW:** None.

### MVP Scope Assessment

The MVP scope remains **Ideal** and aligned with the Project Brief [cite: Project Brief\_ shantilly.md].

### Technical Readiness

The **Technical Assumptions**, now including the expected YAML structure, provide a complete foundation for the Architect [cite: docs/prd.md].

### Recommendations

No further recommendations. Previous items have been addressed.

### Final Decision

**READY FOR ARCHITECT**: The PRD is complete, consistent, and addresses the points raised by the checklist. It is ready for handoff to the Architecture phase.

## Next Steps

### Prompt for UX Expert

N/A (The project is a TUI CLI and does not require a Web/Mobile UI/UX expert at this time).

### Prompt for Architect

**To:** Winston (Architect) 🏗️
**From:** John (PM) 📋

**Subject:** Start Architecture Phase for shantilly (MVP) - Final PRD Attached

Winston,

The PRD for the shantilly MVP is finalized (see `docs/prd.md`). It defines the requirements for building a CLI tool (Go + Charmbracelet) that renders TUI forms from a YAML definition on `stdin` (structure now defined in PRD) and returns results as JSON on `stdout`. Error handling has also been specified (NFR8).

Please review the finalized PRD, especially the **Technical Assumptions** (including YAML structure) and the **Functional/Non-Functional Requirements**.

Your next step is to create the **Architecture Document** (`architecture-tmpl.yaml`) [cite: team-fullstack.txt]. Key points to address:

1. **Detailed Project Structure:** Detail the Go package organization (`cmd/`, `internal/`, sub-packages like `tui`, `config`, `parser`).
2. **Data Structures:** Define the Go structs to represent the YAML configuration defined in the PRD.
3. **Internal Data Flow:** How data flows from YAML parsing (`gopkg.in/yaml.v3`) to the Go struct, into the `bubbletea`/`huh` model, and finally collected for JSON output.
4. **Bubbletea Model:** Specific structure of the `bubbletea` model (state, messages, `Init`/`Update`/`View`).
5. **`huh` Integration:** How `huh` components will be dynamically created from the YAML config and managed within `bubbletea`.
6. **Error Handling:** Architecture for capturing and reporting parsing errors (stderr, non-zero exit status) and TUI errors, per NFR8.
7. **Coding Standards:** Define specific Go patterns for this project.
8. **Testing Strategy:** Detail the approach for unit tests.

I am available to clarify any requirements.
