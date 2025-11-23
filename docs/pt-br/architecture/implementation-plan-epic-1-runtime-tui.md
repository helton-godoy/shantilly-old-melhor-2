# Implementation Plan — Epic 1: Runtime TUI Declarativo — Fundação do Runtime (v2.0)

Este plano orquestra a execução da Fase 3 (bmad-dev/bmad-master), conectando diretamente:

- PRD v2.0: [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura v2.0 (pivot): [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1), [`docs/architecture/index.md`](docs/architecture/index.md:1), [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)
- Story Map canônico: [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
- QA v2.0 (a ser consolidado): `docs/qa/`
- Governança (a ser definida): `docs/architecture/governance-runtime-tui-v2.0.md`

Premissas vinculantes (herdadas do pivot):

- Binário único (`shantilly`).
- YAML único por aplicação como fonte de verdade:
  - Layout (`column` / `row` / `box`)
  - Componentes (`list`, `viewport`, `form`, `buttongroup`)
  - `on:` orientado a eventos
  - `run:` com `args`, `stdin`, `update_target`
- Legado v1.0:
  - Restrito ao `FormComponent` dentro do runtime.
  - Nunca descrito como produto/roadmap principal.
- Constraints técnicas obrigatórias:
  - Sem `os.Exit` no core do runtime.
  - 1 processo por `update_target` (cancelando o anterior).
  - Modal Stack obrigatória para sobreposições.
  - Segurança JIT + anti-Trojan YAML:
    - `deny by default`,
    - whitelists,
    - lockdown de chaves não autorizadas.
- Roteamento único via `EventManager` + `ShantillyEvent` + `on:`.
- Toda funcionalidade rastreável por blocos E1.1–E1.5.

Cada wave abaixo é incremental, verificável e amarrada a Stories, Arquitetura e QA.

---

## Wave 1 — Data Models v2.0 + YAML Parser Declarativo

Escopo:

- Criar o modelo declarativo v2.0 para o YAML único.
- Introduzir parser alinhado ao novo modelo, coexistindo inicialmente com o legado, mas servindo como alvo principal.

Entradas:

- PRD: seções de layout, componentes, `on:` e `run:` em [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1).
- Arquitetura:
  - [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1)
  - [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1) (a ajustar v2.0)

Tarefas:

- Definir structs declarativos (ex.: `AppConfig`, `LayoutNode`, `Component`, `RunAction`, `OnHandler`, `SecurityPolicy`) em pacote dedicado (ex.: `pkg/declarative`).
- Implementar parser:
  - Lê YAML (stdin/--file) → structs v2.0.
  - Inclui validações estruturais mínimas.
- Manter compatibilidade temporária:
  - O parser antigo (FormConfig) é tratado como legado.
  - Novo parser é referência para v2.0.

Critérios de Done (Wave 1):

- Modelos v2.0 documentados em [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1) com referências E1.1–E1.5.
- Parser v2.0 implementado e coberto por testes de unidade (válidos/ inválidos).
- Nenhum acoplamento rígido novo ao modelo legado.
- Traço na futura matriz QA (`matrix-epic-1-runtime-tui-coverage.md`) apontando casos de teste desta wave.

Rastreabilidade:

- Stories: E1.1/E1.2/E1.3 seeds no story map.
- QA: Casos de validação de YAML.

---

## Wave 2 — LayoutManager (E1.1)

Escopo:

- Implementar o motor de layout hierárquico baseado nos modelos v2.0.

Entradas:

- Story Map: E1.1.
- Arquitetura:
  - [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)
  - Seções relevantes em [`docs/architecture/components.md`](docs/architecture/components.md:1)

Tarefas:

- Implementar `LayoutManager` em `internal/runtime/layout`:
  - Constrói árvore de layout a partir do YAML v2.0.
  - Gera estrutura consumível pelo loop `bubbletea`.
- Integração:
  - Suporte a redimensionamento e foco global.
  - Preparar hooks para Modal Stack (sem implementar ainda a pilha completa).

Critérios de Done (Wave 2):

- Layout hierárquico funciona conforme exemplos de YAML.
- Testes de unidade/verificação de estrutura.
- Nenhuma lógica de automação ou script embedada no LayoutManager.
- Referências a E1.1 adicionadas na documentação de arquitetura.

Rastreabilidade:

- Stories E1.1 marcadas como cobertas.
- QA: casos para layout e responsividade planejados.

---

## Wave 3 — EventManager + ShantillyEvent + on: (E1.3)

Escopo:

- Centralizar o fluxo de eventos no runtime declarativo.

Entradas:

- Story Map: E1.3.
- Arquitetura:
  - Definições de `ShantillyEvent` e `EventManager` em [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1) e [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1).

Tarefas:

- Definir tipo `ShantillyEvent` (eventos normalizados).
- Implementar `EventManager` em `internal/runtime/event`:
  - Converte eventos de UI (`tea.Msg` / componentes) em `ShantillyEvent`.
  - Avalia bloco `on:` do YAML para decidir ações (`run:`, updates, modais, etc.).
- Eliminar caminhos paralelos:
  - Toda lógica reativa deve passar por `EventManager` + `on:`.

Critérios de Done (Wave 3):

- Todos os eventos relevantes são roteados via `EventManager`.
- `on:` é o único mecanismo de orquestração de lógica declarativa.
- Não existem handlers “hardcoded” fora deste fluxo.
- Testes cobrindo:
  - Tradução de eventos,
  - Seleção correta de handlers `on:`.

Rastreabilidade:

- Stories E1.3 vinculadas.
- QA: matriz inclui cenários de eventos e roteamento.

---

## Wave 4 — ScriptRunner (run:, args, stdin, update_target) (E1.4)

Escopo:

- Implementar a execução segura e controlada de `run:`.

Entradas:

- Story Map: E1.4.
- Arquitetura:
  - [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)
  - [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)

Tarefas:

- Implementar `ScriptRunner` em `internal/runtime/runner`:
  - Executa `run: { script: ... }` (Epic 1).
  - Suporte completo a:
    - `args` com templates,
    - `stdin` configurável,
    - `env` se previsto no modelo.
  - Entrega outputs para `update_target` (viewport/list/etc.).
- Garantir invariantes:
  - 1 processo por `update_target`:
    - Encerrar processo anterior antes do próximo (SIGTERM/SIGKILL se necessário).
  - Sem `os.Exit` no core:
    - Erros propagados via eventos e retornos de função.

Critérios de Done (Wave 4):

- Implementação respeita 1 processo por `update_target`.
- Nenhum uso de `os.Exit` no core; apenas na camada CLI, se necessário.
- Testes cobrindo:
  - Execuções sequenciais,
  - Cancelamento correto,
  - Falhas seguras.

Rastreabilidade:

- Stories E1.4 cobertas.
- QA: casos para lifecycle de processos e atualização de `update_target`.

---

## Wave 5 — Modal Stack (E1.5 — parte 1)

Escopo:

- Formalizar e implementar a pilha de modais integrada ao LayoutManager/EventManager.

Entradas:

- Story Map: E1.5.
- Arquitetura:
  - Seções de Modal Stack em [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1) e [`docs/architecture/components.md`](docs/architecture/components.md:1).

Tarefas:

- Implementar módulo `internal/runtime/modal`:
  - Estrutura de pilha (push/pop/top).
- Integrar:
  - `LayoutManager` usa Modal Stack para desenhar estado atual.
  - `EventManager` roteia eventos apenas para o topo da pilha quando houver modal.

Critérios de Done (Wave 5):

- Modais funcionam como sobreposições corretas.
- Foco e eventos respeitam topo da pilha.
- Testes cobrindo:
  - Fluxos com múltiplos modais,
  - Fechamento e retorno ao contexto anterior.

Rastreabilidade:

- Stories E1.5 (Modal Stack) associadas.
- QA: casos cobrindo navegação e sobreposição.

---

## Wave 6 — Encapsulamento do Legado como FormComponent (E1.2/E1.5 — confinamento)

Escopo:

- Trazer o v1.0 para dentro do runtime apenas como `FormComponent`, eliminando-o como produto.

Entradas:

- Código atual em `internal/tui` e `internal/config`.
- Arquitetura v2.0 (confinamento do legado) em [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1).

Tarefas:

- Criar `FormComponent` compatível com o modelo declarativo:
  - Reuso criterioso de partes do código v1.0.
  - Interfaces consistentes com `ShantillyComponent`.
- Ajustar docs:
  - Qualquer menção a `shantilly form` como comando/produto principal deve ser:
    - Atualizada para runtime declarativo, ou
    - Marcada como `Legacy (não vinculante)`.
- Garantir:
  - Nenhum fluxo novo é construído diretamente sobre a arquitetura v1.x.

Critérios de Done (Wave 6):

- FormComponent opera como componente do runtime v2.0.
- Todas as referências de produto apontam para o Runtime TUI Declarativo.
- Nenhum artefato ativo sugere v1.x como baseline futura.

Rastreabilidade:

- Stories relacionadas ao formulário mapeadas como E1.2 (legado encapsulado).
- QA: testes de regressão para FormComponent dentro do runtime.

---

## Wave 7 — Hardening de Segurança JIT + Anti-Trojan YAML (E1.5 — parte 2)

Escopo:

- Elevar segurança a requisito de primeira classe, fechando portas para YAML malicioso.

Entradas:

- PRD/Arquitetura: regras de segurança em [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:210) e [`docs/architecture/security.md`](docs/architecture/security.md:1).

Tarefas:

- Implementar:
  - Validações de schema estrito para o YAML declarativo.
  - `deny by default` para chaves desconhecidas.
  - Whitelist configurada para:
    - Tipos de ações permitidas (`run.script`),
    - Caminhos/destinos seguros.
- Integrar com `ScriptRunner`:
  - Nenhuma execução ocorre fora das regras aprovadas.
- Documentar:
  - Regras e invariantes em [`docs/architecture/security.md`](docs/architecture/security.md:1).

Critérios de Done (Wave 7):

- Códigos de validação ativados por padrão.
- Testes cobrindo cenários maliciosos (Trojan YAML, comandos proibidos, chaves estranhas).
- QA: gates específicos de segurança apontando para estes casos.

Rastreabilidade:

- Stories E1.5 (Security).
- QA: matriz com linha dedicada a cada política de segurança.

---

## Critérios de Done Globais do Epic 1

O Epic 1 é considerado concluído quando:

1. Todas as Waves 1–7:
   - Estão implementadas,
   - Atendem seus critérios de done específicos,
   - Estão mapeadas na matriz QA.

2. Rastreabilidade ponta a ponta:
   - Cada requisito do PRD E1.1–E1.5:
     - Tem stories em `docs/stories/`,
     - Tem contratos em `docs/architecture/`,
     - Tem implementação identificada (pacotes/módulos),
     - Tem testes associados em `docs/qa/`.

3. Conformidade com constraints:
   - Não há `os.Exit` no core.
   - 1 processo por `update_target` garantido por testes.
   - Modal Stack implementada e utilizada.
   - Segurança JIT + anti-Trojan habilitadas e testadas.

4. Legado confinado:
   - v1.0 existe apenas como FormComponent.
   - Nenhum documento/código ativo sugere v1.x como produto/roadmap vigente.

---

## Papel dos Agentes neste Plano

- bmad-dev / bmad-master:
  - Executam Waves 1–7 seguindo este documento como contrato.
  - Atualizam comentários/docblocks de código para refletir blocos E1.x e referências.

- bmad-architect:
  - Mantém contratos em `docs/architecture/*.md` alinhados às implementações.
  - Ajusta qualquer documento residual para o pivot v2.0.

- bmad-qa:
  - Constrói `docs/qa/matrix-epic-1-runtime-tui-coverage.md` usando este plano como backbone.
  - Define gates cobrindo cada Wave e constraint.

- bmad-orchestrator:
  - Usa este plano + matriz QA + story map para auditar PRs e releases.
  - Garante bloqueio de qualquer desvio das regras invioláveis.

---

## Próximo passo após este ciclo

Após a conclusão verificada do Epic 1 (todas as Waves implementadas e cobertas):

- Próxima onda recomendada:
  - Iniciar Epic 2 (runners especializados, ex.: `run: { ansible_playbook: ... }`),
  - Reutilizando integralmente o runtime declarativo estável:
    - Mesmo binário único,
    - Mesmo contrato YAML único,
    - Mesmo EventManager/ScriptRunner/Modal Stack/Security.
- A governança estabelecida garantirá que qualquer evolução futura (Epics 2+) não quebre:
  - O contrato declarativo,
  - As constraints de segurança,
  - O confinamento do legado.

Este plano é vinculante para execução do Epic 1 v2.0 e serve como referência direta para coordenação entre todos os agentes.
