# Introduction

Este documento descreve a arquitetura ativa do `shantilly` v2.0 e serve como blueprint vinculante para todos os agentes (Analyst, Architect, Dev/Master, QA, Orchestrator).

A partir do pivot definido no Epic 1: Runtime TUI Declarativo — Fundação do Runtime (v2.0), o projeto é, de forma definitiva, um:

- Runtime TUI Declarativo orientado a eventos.
- Implementado como binário único.
- Que consome um único YAML de aplicação como fonte de verdade para:
  - Layout hierárquico (`column`, `row`, `box`),
  - Componentes (`list`, `viewport`, `form`, `buttongroup`),
  - Bloco `on:` orientado a eventos (`ShantillyEvent`),
  - Ações `run:` com `args`, `stdin`, `update_target`,
  - Com segurança JIT e proteções anti-Trojan YAML como requisitos de primeira classe.

Qualquer visão anterior baseada em “shantilly form” como produto central ou arquitetura v1.x é, a partir deste documento, classificada como LEGADO restrito ao `FormComponent` dentro do runtime declarativo, apenas como referência histórica/técnica, não como roadmap ou baseline futura.

## Runtime TUI Declarativo v2.0 — Princípios Arquiteturais Centrais

1. Contrato YAML Único
   - Há um único contrato YAML de aplicação que define:
     - A árvore de layout (E1.1 Layout Engine).
     - O modelo de componentes (E1.2 Component Model).
     - O modelo de eventos e handlers `on:` (E1.3 Event Engine).
     - As ações `run:` (E1.4 ScriptRunner).
     - As regras de Modal Stack e segurança (E1.5 Modal Stack + Security).
   - Nenhum outro artefato pode redefinir o modelo conceitual fora do alinhamento entre:
     - [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1),
     - Este conjunto de documentos de arquitetura v2.0,
     - O story map canônico em [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1),
     - A matriz de QA v2.0.

2. Confinamento do Legado v1.0
   - O design original centrado em `form` é tratado como LEGADO.
   - Sua função presente:
     - Servir como base técnica para a implementação do `FormComponent`.
   - Restrições:
     - É proibido tratar “shantilly form” ou a arquitetura v1.x como produto principal, roadmap vigente ou baseline conceitual futura.
     - Qualquer documento, história, teste ou código que sugira isso deve ser corrigido ou marcado explicitamente como `legacy`.

3. Constraints Técnicas Invioláveis
   - Sem `os.Exit` no core do runtime:
     - Encerramentos são geridos por mensagens/erros internos tipados.
     - `os.Exit` só pode existir na casca de CLI, nunca dentro dos motores do runtime.
   - Um processo por `update_target`:
     - O `ScriptRunner` deve sempre cancelar o processo anterior associado ao mesmo `update_target` antes de iniciar um novo.
   - Modal Stack obrigatória:
     - Toda sobreposição de UI (ex.: confirmações, prompts sensíveis) é gerida via uma pilha de modais explícita, integrada ao `LayoutManager`.
   - Segurança JIT + Anti-Trojan YAML:
     - `deny by default` para chaves não reconhecidas.
     - Whitelists de ações e destinos permitidos.
     - Validações que previnem uso malicioso do YAML (Trojan YAML).
   - Essas constraints são requisitos arquiteturais e de QA, não recomendações opcionais.

4. Rastreabilidade Estruturada (E1.x)
   - Cada decisão arquitetural é vinculada aos blocos:
     - E1.1 — Layout Engine
     - E1.2 — Component Model
     - E1.3 — Event Engine (EventManager + ShantillyEvent + on:)
     - E1.4 — ScriptRunner (run:/args/stdin/update_target)
     - E1.5 — Modal Stack + Security (JIT + Anti-Trojan)
   - Todos os artefatos devem:
     - Referenciar explicitamente esses IDs quando relevante.
     - Manter links cruzados consistentes usando o padrão descrito em [`docs/architecture/coding-standards.md`](docs/architecture/coding-standards.md:50).

## Relationship to Other Architecture Docs

Este arquivo funciona como a “âncora normativa” dos demais documentos em `docs/architecture/`:

- [`docs/architecture/index.md`](docs/architecture/index.md:1) organiza a arquitetura v2.0 do Runtime TUI Declarativo.
- [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1) detalha o estilo arquitetural declarativo, fluxo de dados/eventos e binário único.
- [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1) deve evoluir para refletir o contrato YAML único (substituindo o modelo linear `FormConfig` como visão primária).
- [`docs/architecture/components.md`](docs/architecture/components.md:1) define `ShantillyComponent`, componentes centrais (list/viewport/form/buttongroup) e como o legado é encapsulado como `FormComponent`.
- [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1) mostra os fluxos de execução baseados no runtime declarativo (serão atualizados para o modelo v2.0).
- [`docs/architecture/security.md`](docs/architecture/security.md:1) consolida as garantias de Segurança JIT e anti-Trojan YAML.
- [`docs/architecture/error-handling-strategy.md`](docs/architecture/error-handling-strategy.md:1) será alinhado para remover qualquer dependência de `os.Exit` no core.
- Documentos como [`docs/architecture/checklist-results-report.md`](docs/architecture/checklist-results-report.md:1) e [`docs/architecture/next-steps.md`](docs/architecture/next-steps.md:1) passam a ser interpretados à luz desta arquitetura pivot.

Quaisquer seções ou arquivos que reflitam a arquitetura v1.x como estado presente ou recomendado:

- Devem ser:
  - Atualizados para o contexto v2.0, ou
  - Explicitamente marcados como `## Legacy (não vinculante para v2.0)`.

## Papel dos Agentes neste Contexto

- bmad-analyst:
  - Usa este documento como referência para mapear stories a E1.1–E1.5 e ao contrato YAML único.
- bmad-architect:
  - Garante que todos os docs de arquitetura reflitam estes princípios e definam contratos formais coerentes.
- bmad-dev / bmad-master:
  - Implementam o runtime v2.0 estritamente aderente a estes contratos e constraints.
- bmad-qa:
  - Define estratégia e matriz de testes que verifiquem cada constraint e cada bloco E1.x.
- bmad-orchestrator:
  - Audita continuamente todos os artefatos para garantir aderência a este pivot e confinamento do legado.

Este documento substitui, como referência ativa, qualquer interpretação anterior da arquitetura que contrarie o Runtime TUI Declarativo v2.0, o YAML único ou o confinamento do legado ao FormComponent.
