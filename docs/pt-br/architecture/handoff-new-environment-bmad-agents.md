# Handoff Hands-on — Novo Ambiente Shantilly v2.0 (Agentes BMAD)

Este artefato é o prompt operativo e normativo para qualquer agente (humano ou IA) iniciando trabalho em um NOVO ambiente a partir do estado atual do projeto Shantilly v2.0.

Sempre carregue ESTE TEXTO no agente antes de pedir "qual a próxima estória" ou antes de implementar qualquer coisa.

---

## 1. Identidade do Projeto (NÃO ALTERAR)

Você está trabalhando no:

- `shantilly` v2.0 — Runtime TUI Declarativo orientado a eventos.
- Binário único.
- Um único YAML de aplicação contendo:
  - `layout` (column/row/box),
  - `components` (list, viewport, form, buttongroup),
  - `on:` orientado a `ShantillyEvent`,
  - `run:` com `args`, `stdin`, `update_target`,
  - Segurança JIT + Anti-Trojan YAML.
- Código legado v1.x (`shantilly form`, `FormConfig`, `internal/tui`) é APENAS backend interno do `FormComponent` legado, nunca produto ou arquitetura ativa.

Qualquer visão diferente disso está errada para este repositório.

Fontes normativas:

- [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1)
- [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)
- [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)
- [`docs/architecture/components.md`](docs/architecture/components.md:1)
- [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
- [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)

---

## 2. Regra-Mestra: Como escolher a PRÓXIMA ESTÓRIA

Sempre responda "qual a próxima estória?" usando SOMENTE este algoritmo:

1. Use o Story Map canônico:
   - [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
2. Use a Matriz de Cobertura:
   - [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)
3. Para Epic 1 (E1.1–E1.5), a próxima estória é SEMPRE:
   - O primeiro item ainda em `planned`/`draft` cujos pré-requisitos (blocos e stories anteriores) estejam `implemented` ou `passed` na matriz.
4. É proibido:
   - Escolher histórias por número do arquivo apenas.
   - Saltar direto para hardenings avançados (ex.: `7.1.hardening-seguranca-jit-anti-trojan-e15.story.md`) ignorando blocos centrais.

No snapshot atual deste repositório:

- A próxima estória OBRIGATÓRIA é:
  - E1.3 — Event Engine (Story [`docs/stories/3.3.event-engine-and-on-routing.story.md`](docs/stories/3.3.event-engine-and-on-routing.story.md:1))
- Qualquer agente que proponha outra coisa como próximo passo está violando este handoff.

---

## 3. Prompt Operacional Único para Agentes no Novo Ambiente

Sempre inicialize o agente com este prompt (pode ser colado literalmente):

"""
Você é um agente BMAD trabalhando no repositório Shantilly v2.0 (Runtime TUI Declarativo).

SUAS REGRAS OBRIGATÓRIAS:

1. Fonte de verdade:
   - Leia e respeite:
     - [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
     - [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)
     - [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1)
     - [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)
     - [`docs/architecture/components.md`](docs/architecture/components.md:1)
   - Trate esses arquivos como contratos normativos. Não crie uma ordem de implementação própria.

2. Como decidir a PRÓXIMA ESTÓRIA:
   - Use o seguinte algoritmo determinístico:
     a) Percorra os blocos Epic 1 na ordem: E1.1 → E1.2 → E1.3 → E1.4 → E1.5.
     b) Para cada bloco, consulte a Matriz de Cobertura para ver o status.
     c) A "próxima estória" é o primeiro item com status `planned` ou `draft` cujo conjunto de pré-requisitos (stories/linhas anteriores / blocos anteriores) esteja `implemented` ou `passed`.
   - Aplique isso antes de responder qualquer pergunta do tipo "qual a próxima estória?" ou antes de propor implementação.

3. Estado atual esperado (snapshot deste repo):
   - E1.1, E1.2: contratos e visão consolidados em docs; implementação full ainda em evolução.
   - E1.4 (ScriptRunner) e E1.5 (Modal Stack + Security JIT/Anti-Trojan): contratos e gates definidos; partes marcadas como implemented.
   - E1.3 — Event Engine (EventManager + ShantillyEvent + on:):
     - Story [`docs/stories/3.3.event-engine-and-on-routing.story.md`](docs/stories/3.3.event-engine-and-on-routing.story.md:1) está em `Status: draft`.
     - Matriz aponta E1.3 como planned.
     - Portanto, E1.3 é a próxima estória obrigatória.
   - NUNCA pule para stories de Wave 7 (ex.: `7.1.hardening-seguranca-jit-anti-trojan-e15.story.md`) sem E1.3 consolidado e gates anteriores satisfeitos.

4. Invariantes arquiteturais que você deve proteger:
   - Pipeline único obrigatório:
     - `ShantillyEvent` → `EventManager` → `on:` → (`RunAction` | `ModalRequest` | `UpdateState`) → ScriptRunner/Modal Stack/Layout.
   - Proibições:
     - Nenhum `os.Exit` fora de [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1).
     - Nenhuma execução de script fora de `internal/runtime/runner/**`.
     - Nenhum modal fora de `internal/runtime/modal/**`.
     - `LayoutManager` não executa scripts, não avalia `on:`, não abre modal.
     - Componentes não executam scripts nem hardcodam lógica de automação.
   - Legado:
     - `internal/tui/**` e `FormConfig` só podem existir encapsulados no `FormComponent`.
     - Qualquer uso direto novo é bug arquitetural.

5. Sua próxima ação concreta AO INICIAR NO NOVO AMBIENTE:
   - Reafirme explicitamente:
     - "Próxima estória: implementar e consolidar E1.3 — Event Engine (Story 3.3) como pipeline único de eventos."
   - Em seguida:
     - Detalhe plano de implementação em cima de:
       - [`docs/stories/3.3.event-engine-and-on-routing.story.md`](docs/stories/3.3.event-engine-and-on-routing.story.md:1)
       - [`docs/architecture/components.md#4-eventmanager--shantillyevent--e13`](docs/architecture/components.md:257)
       - [`docs/architecture/core-workflows.md#workflow-1--execucao-declarativa-completa-happy-path--e11e14`](docs/architecture/core-workflows.md:23)
   - Não avance para nenhuma história de Wave 5–7 sem:
     - E1.3 implementado,
     - QA Gate `1.x.event-engine.yml` definido e respeitado.

6. Em caso de dúvida:
   - Consulte sempre Story Map + Matriz.
   - Se uma escolha de próxima estória entrar em conflito com esses artefatos, considere-a inválida e corrija-a.
"""

Esse é o prompt canônico de handoff. Qualquer agente no novo ambiente deve ser inicializado com esse texto antes de sugerir roadmap ou próxima estória.
