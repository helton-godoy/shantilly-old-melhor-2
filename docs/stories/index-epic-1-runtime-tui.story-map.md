# Story Map — Epic 1: Runtime TUI Declarativo — Fundação do Runtime (v2.0)

Este documento é o mapa canônico de histórias do **Epic 1** e a âncora funcional para todos os agentes.
Ele consolida a rastreabilidade entre:

- PRD v2.0: [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura v2.0 (Runtime TUI Declarativo): [`docs/architecture.md`](docs/architecture.md:1) e [`docs/architecture/index.md`](docs/architecture/index.md:1)
- Roadmap v2.0: [`docs/prd/epic-list.md`](docs/prd/epic-list.md:1)
- Stories: `docs/stories/*.story.md`
- QA / Gates / Matriz de Cobertura: `docs/qa/`

Premissas invioláveis:

- Um único binário (`shantilly`) consumindo um **único YAML de aplicação**, contendo:
  - Layout hierárquico (`column` / `row` / `box`)
  - Componentes (`list`, `viewport`, `form`, `buttongroup`)
  - Bloco `on:` orientado a eventos (`ShantillyEvent`)
  - Ações `run:` com `args`, `stdin`, `update_target`
- Segurança JIT + anti-Trojan YAML como requisitos de primeira classe.
- Legado v1.0 estritamente confinado como implementação de `FormComponent` dentro do runtime declarativo.
- Proibições técnicas obrigatórias:
  - Sem `os.Exit` no core do runtime.
  - Um processo por `update_target` (cancelando o anterior).
  - Modal Stack obrigatória para sobreposições.
  - `deny by default` para execuções suspeitas + whitelist de ações + lockdown de chaves YAML não autorizadas.
- Nenhum artefato (docs, histórias, arquitetura, código, QA) pode tratar "shantilly form" ou a arquitetura v1.x como produto/roadmap principal ou baseline futura.

Estrutura de rastreabilidade:

Cada bloco abaixo é identificado como **E1.x** e vincula:

- Requisitos de PRD
- Contratos de arquitetura
- Stories (existentes/novas)
- Alvos de implementação
- Linhas na matriz de QA

---

## E1.1 — Layout Engine (Layout hierárquico declarativo)

Escopo:

- Interpretar o YAML único e construir um layout hierárquico usando `column`, `row`, `box`.
- Integrar com foco global e futura Modal Stack.
- Não contém regras de negócio; apenas estrutura visual/comportamental do layout.

Fontes:

- PRD: seção de layout no [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura:
  - Visão geral: [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)
  - Data Models (layout): [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)

Stories (exemplo de mapeamento inicial; bmad-analyst deve ajustar/expandir):

- `docs/stories/1.3.basic-tui-structure-with-bubbletea.story.md`
  - Epic: 1
  - Bloco: E1.1
  - Status: legacy_refactored (ajustar narrativa para runtime declarativo)
- `docs/stories/1.7.basic-styling-and-alignment.story.md`
  - Epic: 1
  - Bloco: E1.1
  - Status: legacy_refactored

Implementação alvo (conceitual):

- `internal/runtime/layout/...` (a ser definido no plano de implementação)

Notas:

- Ajustar todas as stories de layout para referenciar explicitamente o modelo declarativo YAML e não o fluxo linear v1.0.

---

## E1.2 — Component Model (list, viewport, form, buttongroup)

Escopo:

- Definir `ShantillyComponent` como padrão.
- Componentes oficiais: `list`, `viewport`, `form`, `buttongroup`.
- `FormComponent` encapsula o legado v1.0 dentro do runtime:
  - Pode reutilizar internals v1.0,
  - Não redefine visão de produto.

Fontes:

- PRD: componentes essenciais no [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura:
  - Componentes: [`docs/architecture/components.md`](docs/architecture/components.md:1)
  - High Level Architecture: [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)

Stories:

- `docs/stories/1.4.form-rendering-with-huh.story.md`
  - Epic: 1
  - Bloco: E1.2
  - Status: legacy_refactored (como comportamento do FormComponent dentro do runtime declarativo)
- `docs/stories/1.5.keyboard-navigation-and-interaction.story.md`
  - Epic: 1
  - Bloco: E1.2
  - Status: legacy_refactored
- Novas stories a definir para `list`, `viewport`, `buttongroup` (bmad-analyst deve criar com cabeçalho padrão).

Implementação alvo:

- `internal/components/...`
- `FormComponent` refatorado a partir de `internal/tui/...`, encapsulado.

Notas:

- Qualquer menção a "shantilly form" como produto deve ser corrigida para "FormComponent no runtime declarativo".

---

## E1.3 — Event Engine (EventManager + ShantillyEvent + on:)

Escopo:

- Definir o motor de eventos declarativos:
  - `ShantillyEvent` como tipo central.
  - `EventManager` roteando todos os eventos:
    - Interações de UI,
    - Respostas de scripts,
    - Ações internas (open modal, update component, etc).
  - Bloco `on:` no YAML como orquestrador único.

Fontes:

- PRD: blocos `on:` e eventos em [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura:
  - Workflows: [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)
  - Data Models (events): [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)

Stories:

- Ajustar/introduzir stories específicas para:
  - Emissão de eventos por componentes (`form:submit`, `list:select`, etc.).
  - Roteamento centralizado via `EventManager`.
  - Proibição de rotas alternativas fora de `on:`.

Implementação alvo:

- `internal/runtime/event/...`
- Definição de `ShantillyEvent` como estrutura única de eventos.

Notas:

- Toda automação deve ser vinculada a `on:`; nenhum comportamento "hardcoded" fora do modelo declarativo.

---

## E1.4 — ScriptRunner (run:, args, stdin, update_target)

Escopo:

- Executar ações `run:` definidas no YAML:
  - Suportar `args`, `stdin`, variáveis a partir do estado.
  - Implementar regra de 1 processo por `update_target`:
    - Sempre cancelar o anterior antes de iniciar o próximo.
- Integrar com EventManager:
  - Resultados/erros retornam como eventos.
- Respeitar constraints:
  - Sem `os.Exit` no core.

Fontes:

- PRD: seção de execução e `run:` em [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura:
  - Data Models (RunAction): [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)
  - Workflows de execução: [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)

Stories:

- Criar/ajustar stories para:
  - Execução segura de scripts.
  - update_target e cancelamento.
  - Integração com eventos e YAML único.

Implementação alvo:

- `internal/runtime/runner/...`

Notas:

- Qualquer referência a execução "ad hoc" fora do modelo `run:`/YAML deve ser removida ou marcada como inválida.

---

## E1.5 — Modal Stack + Security (JIT + Anti-Trojan YAML)

Escopo:

- Modal Stack:
  - Gerenciar telas/overlays empilhados (ex.: prompts, confirmações).
  - Integrar com LayoutEngine e foco global.
- Segurança:
  - Schema formal do YAML.
  - `deny by default` para chaves desconhecidas ou execuções suspeitas.
  - Whitelist de ações/paths permitidos.
  - Proteção contra uso malicioso de `run:` (anti-Trojan YAML).

Fontes:

- PRD: restrições de segurança em [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:210)
- Arquitetura:
  - Segurança: [`docs/architecture/security.md`](docs/architecture/security.md:1)
  - Workflows principais: [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)

Stories:

- Criar/ajustar stories para:
  - Modal Stack (push/pop, foco, integração com `on:`).
  - Validações de YAML e bloqueios.
  - Casos de uso de segurança (inputs maliciosos, comandos proibidos).

Implementação alvo:

- `internal/runtime/modal/...` (a definir)
- Módulo de validação YAML + políticas de segurança.

Notas:

- Testes e gates de segurança devem ser vinculados diretamente a este bloco.

---

## Convenções para TODAS as stories de Epic 1

bmad-analyst deve garantir que cada arquivo `*.story.md` relacionado ao Epic 1:

1. Contém cabeçalho obrigatório:
   - `Epic: 1 - Runtime TUI Declarativo — Fundação do Runtime`
   - `Bloco: E1.x - <nome>`
   - `Trace: PRD (<arquivo>#<seção>), Arquitetura (<arquivo>#<seção>), QA (matrix-epic-1-runtime-tui-coverage.md#<id>)`
   - `Status: {draft|approved|legacy_refactored|new}`

2. Está explicitamente alinhado ao:
   - Contrato do YAML único,
   - Confinamento do legado v1.0 ao FormComponent,
   - Constraints técnicas globais.

3. Não contém:
   - Qualquer formulação que eleve v1.x a roadmap principal,
   - Qualquer ambiguidade sobre múltiplos contratos ou múltiplos modos divergentes do modelo declarativo.

---

## Próximos passos após este mapa

Com este mapa criado:

- bmad-analyst:
  - Refatora imediatamente as stories existentes para aderir a este padrão.
  - Cria as stories faltantes para cobrir integralmente E1.1–E1.5.
- bmad-architect:
  - Usa este mapa como referência para normalizar os contratos em `docs/architecture/`.
- bmad-dev/bmad-master:
  - Usa E1.1–E1.5 e este story map para estruturar `implementation-plan-epic-1-runtime-tui.md`.
- bmad-qa:
  - Constrói `matrix-epic-1-runtime-tui-coverage.md` usando este índice como backbone.
- bmad-orchestrator:
  - Passa a auditar todos os artefatos com base neste story map + contratos + matriz.

Este arquivo é agora parte do trilho de rastreabilidade formal do Epic 1 e deve ser mantido sincronizado a cada mudança relevante.
