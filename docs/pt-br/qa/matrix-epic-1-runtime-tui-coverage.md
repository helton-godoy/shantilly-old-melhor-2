# Matriz de Cobertura — Epic 1: Runtime TUI Declarativo — Fundação do Runtime (v2.0)

Esta matriz conecta, linha a linha:

- Blocos do Epic 1 (E1.1–E1.5)
- Stories (docs/stories/)
- Arquitetura (docs/architecture/)
- Implementação (caminhos de código)
- Tipos de teste (unit / integration / e2e / security)
- Gates de QA (docs/qa/gates/)
- Status de execução (planned / implemented / passed)

Ela é a referência única para bmad-qa, bmad-dev/bmad-master, bmad-architect e bmad-orchestrator verificarem se o Runtime TUI Declarativo v2.0 está integralmente coberto e alinhado ao contrato YAML único.

Formato das linhas (normativo):

- E1.x: ID do bloco do Epic 1.
- Story ID/Arquivo: nome do arquivo `.story.md` + cabeçalho padronizado.
- Arquitetura: arquivo#seção relevante.
- Implementação: diretório/arquivo Go alvo.
- Testes: tipo(s) requerido(s).
- Gate QA: arquivo YAML ou referência em `docs/qa/gates/`.
- Status: `planned` | `implemented` | `passed`.

---

## Bloco E1.1 — Layout Engine (Layout hierárquico declarativo)

1. E1.1 — Layout básico column/row/box
   - Story: `docs/stories/1.3.basic-tui-structure-with-bubbletea.story.md`
   - Arquitetura:
     - [`docs/architecture/high-level-architecture.md#technical-summary`](docs/architecture/high-level-architecture.md:3)
     - [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1) (modelo v2.0 a alinhar)
   - Implementação:
     - `internal/runtime/layout/...`
   - Testes:
     - unit: construção de árvore de layout
     - integration: renderização básica em terminal
   - Gate QA:
     - `docs/qa/gates/1.7.basic-styling-and-alignment.yml`
     - `docs/qa/gates/1.x.layout-manager.yml` # Gate E1.1 — blocker: nenhum desenvolvimento de LayoutManager fora deste contrato
   - Status: planned

2. E1.1 — Responsividade e alinhamento
   - Story: `docs/stories/1.7.basic-styling-and-alignment.story.md`
   - Arquitetura:
     - [`docs/architecture/components.md`](docs/architecture/components.md:1) (tema/layout)
   - Implementação:
     - `internal/runtime/layout/...`
   - Testes:
     - unit + integration (diferentes tamanhos de terminal)
   - Gate QA:
     - `docs/qa/gates/1.7.basic-styling-and-alignment.yml`
   - Status: planned

---

## Bloco E1.2 — Component Model (list, viewport, form, buttongroup + FormComponent legado)

3. E1.2 — FormComponent como legado encapsulado
   - Stories:
     - `docs/stories/1.4.form-rendering-with-huh.story.md`
     - `docs/stories/1.6.submission-and-json-output.story.md`
   - Arquitetura:
     - [`docs/architecture/components.md#core-components-form-list-viewport-buttongroup`](docs/architecture/components.md:1)
   - Implementação:
     - `internal/components/form_component.go` (a definir)
     - reutilização controlada de `internal/tui/...`
   - Testes:
     - unit: render/submit do FormComponent
     - integration: compatibilidade com YAML declarativo
   - Gate QA:
     - `docs/qa/gates/1.4.form-rendering-with-huh.yml`
     - `docs/qa/gates/1.6.submission-and-json-output.yml`
   - Status: planned

4. E1.2 — List, Viewport, Buttongroup componentes essenciais
   - Stories:
     - novas stories a criar (bmad-analyst) para cada componente
   - Arquitetura:
     - [`docs/architecture/components.md`](docs/architecture/components.md:1)
   - Implementação:
     - `internal/components/list_*.go`
     - `internal/components/viewport_*.go`
     - `internal/components/buttongroup_*.go`
   - Testes:
     - unit + integration por componente
   - Gate QA:
     - novo gate `docs/qa/gates/1.x.components-core.yml` (bmad-qa)
   - Status: planned

---

## Bloco E1.3 — Event Engine (EventManager + ShantillyEvent + on:)

5. E1.3 — Normalização de eventos e roteamento via on:
   - Stories:
     - stories específicas E1.3 (bmad-analyst) a partir do story map
   - Arquitetura:
     - [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)
     - [`docs/architecture/data-models.md#eventos-normalizados-shantillyevent`](docs/architecture/data-models.md:1)
   - Implementação:
     - `internal/runtime/event/...`
   - Testes:
     - unit: mapeamento tea.Msg → ShantillyEvent
     - unit: resolução de on: correto
     - integration: fluxo UI → on: → ação
   - Gate QA:
     - novo gate `docs/qa/gates/1.x.event-engine.yml`
   - Status: planned

---

## Wave 1 — Data Models v2.0 + Parser Declarativo

12. Wave 1 — Modelos declarativos v2.0 + parser LoadAppConfig
    - Stories:
      - `docs/stories/1.2.yaml-configuration-parsing.story.md`
      - seeds E1.1–E1.5 no `docs/stories/index-epic-1-runtime-tui.story-map.md`
    - Arquitetura:
      - [`docs/architecture/data-models.md#1-appconfig-e11e15--raiz-do-yaml-unico`](docs/architecture/data-models.md:13)
      - [`docs/architecture/data-models.md#2-layoutnode-e11--layoutmanager`](docs/architecture/data-models.md:48)
      - [`docs/architecture/data-models.md#3-component-e12--shantillycomponent--componentes-oficiais`](docs/architecture/data-models.md:81)
      - [`docs/architecture/data-models.md#4-onhandler-e13--bloco-on-como-unica-orquestracao`](docs/architecture/data-models.md:116)
      - [`docs/architecture/data-models.md#5-shantillyevent-e13--tipo-de-evento-normalizado`](docs/architecture/data-models.md:146)
      - [`docs/architecture/data-models.md#6-runaction-e14--execucao-declarativa`](docs/architecture/data-models.md:172)
      - [`docs/architecture/data-models.md#7-modalrequest-e15--modal-stack`](docs/architecture/data-models.md:205)
      - [`docs/architecture/data-models.md#8-securitypolicy-e15--anti-trojan-yaml`](docs/architecture/data-models.md:237)
    - Implementação:
      - `pkg/declarative/models.go` (AppConfig, LayoutNode, Component, OnHandler, ShantillyEvent, RunAction, ModalRequest, SecurityPolicy)
      - `pkg/declarative/models.go` (LoadAppConfig, Validate, validações estruturais mínimas)
    - Testes:
      - unit:
        - `pkg/declarative/models_test.go` cobrindo:
          - layout obrigatório e tipos válidos,
          - IDs únicos de componentes e tipos oficiais,
          - referências de layout → componentes,
          - rejeição de tipos inválidos e duplicados,
          - exigência de `event` em `on:`,
          - aceitação estrutural de `security`/SecurityPolicy.
    - Gate QA:
      - novo gate planejado `docs/qa/gates/1.x.declarative-config-and-parser.yml`
    - Status: implemented

    - Notas de rastreabilidade Wave 1:
      - E1.1 (layout declarativo) → LayoutNode modelado e validado.
      - E1.2 (component model) → Component + binds modelados e validados.
            - E1.3 (event/on:) → OnHandler + ShantillyEvent modelados (validação mínima de presença de event).
            - E1.4 (run:) → RunAction modelado; validação estrutural básica via AppConfig/OnHandler.
            - E1.5 (security/Modal) → SecurityPolicy + ModalRequest/ModalField estruturais presentes; hardening completo pendente (Wave 7).

      ## Gates QA — Waves 4–7 (Trilho Bloqueante Consolidado)

      1. `docs/qa/gates/1.x.event-engine.yml`
         - Escopo:
           - Garantir pipeline único `ShantillyEvent` → [`EventManager`](internal/runtime/event/manager.go:1) → `on:`.
         - Checks mínimos:
           - Proibir handlers soltos em `internal/runtime/**` que ignorem `EventManager/on:`.
           - Verificar uso consistente de tipos `ShantillyEvent`.
         - Enforcement:
           - Bloqueante para qualquer rota paralela de eventos.

      2. `docs/qa/gates/1.x.layout-manager.yml`
         - Escopo:
           - Assegurar [`LayoutManager`](internal/runtime/layout/manager.go:1) apenas como engine de layout.
         - Checks:
           - Nenhum `run:`/exec/modal/uso de legado dentro de `internal/runtime/layout/**`.
         - Enforcement:
           - PR falha se LayoutManager assumir funções de automação.

      3. `docs/qa/gates/1.x.scriptrunner-and-update-target.yml`
         - Escopo:
           - ScriptRunner único em `internal/runtime/runner/**`.
         - Checks:
           - Não existência de execução direta (`os/exec`, `exec.*`, `cmd.*`) fora de `internal/runtime/runner/**`.
           - 1 processo por `update_target` garantido via testes.
           - Saída roteada apenas via `update_target`.
         - Enforcement:
           - Bloqueante para qualquer bypass e múltiplos runners concorrentes.

      4. `docs/qa/gates/1.x.modal-stack.yml`
         - Escopo:
           - Modal Stack única em `internal/runtime/modal/**`.
         - Checks:
           - Padrões de `push/pop/top` apenas nesse módulo.
           - Proibição de modais/prompts soltos em componentes/LM/legado.
           - Foco apenas no topo e bloqueio de fundo cobertos por testes.
         - Enforcement:
           - Bloqueante se qualquer modal surgir fora da pilha oficial.

      5. `docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml`
         - Escopo:
           - Confinamento do legado v1.x ao FormComponent sandboxado.
         - Checks:
           - Uso de `internal/tui`/`FormConfig` apenas dentro do sandbox autorizado.
           - Nenhum novo código chamando diretamente APIs legadas.
         - Enforcement:
           - Bloqueante se legado “vazar” para o runtime v2.0.

      6. `docs/qa/gates/1.x.no-osexit-core.yml`
         - Escopo:
           - Proibir `os.Exit` fora da casca CLI.
         - Checks:
           - Scan por `os.Exit` em `internal/**` (exceto [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1)).
         - Enforcement:
           - PR falha se qualquer uso indevido for detectado.

      7. `docs/qa/gates/1.x.security-jit-anti-trojan.yml`
         - Escopo:
           - Segurança JIT + Anti-Trojan YAML.
         - Checks:
           - Validação `AppConfig` com deny-by-default + whitelist de tipos/ações/paths.
           - Nenhum auto-run implícito, componentes/tipos não documentados ou comandos fora da whitelist.
           - Integração com `SecurityPolicy` consumida pelo ScriptRunner.
         - Enforcement:
           - Bloqueante para qualquer configuração insegura ou tentativa de Trojan YAML.

      Encadeamento normativo:
      - Waves 2–3: gates de layout/event-engine.
      - Wave 4: ativa `1.x.scriptrunner-and-update-target.yml` + `1.x.no-osexit-core.yml`.
      - Wave 5: ativa `1.x.modal-stack.yml`.
      - Wave 6: ativa `1.x.legacy-formcomponent-encapsulation.yml`.
      - Wave 7: ativa `1.x.security-jit-anti-trojan.yml`.
      - Todos são cumulativos e obrigatórios para PRs e releases.

## Bloco E1.4 — ScriptRunner (run:, args, stdin, update_target) + 1 processo/update_target + sem os.Exit no core

6. E1.4 — Execução segura de run: com args/stdin/update_target
   - Stories:
     - [`docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md`](docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md:1)
   - Arquitetura:
     - [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)
     - [`docs/architecture/error-handling-strategy.md`](docs/architecture/error-handling-strategy.md:1)
     - [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)
   - Implementação:
     - `internal/runtime/runner/...`
   - Testes:
     - unit: binding de args/stdin e roteamento para update_target
     - integration: update_target recebendo saída correta
   - Gate QA:
     - `docs/qa/gates/1.x.scriptrunner-and-update-target.yml`
   - Status: implemented

7. E1.4 — Regra 1 processo por update_target
   - Stories:
     - parte de [`docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md`](docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md:1)
   - Arquitetura:
     - [`docs/architecture/high-level-architecture.md#high-level-overview`](docs/architecture/high-level-architecture.md:22)
     - [`docs/architecture/security.md#5-regra-de-1-processo-por-update_target-como-invariante-de-seguranca-e14`](docs/architecture/security.md:167)
   - Implementação:
     - `internal/runtime/runner/...`
   - Testes:
     - integration + security: iniciar novo run cancela anterior para mesmo update_target
   - Gate QA:
     - `docs/qa/gates/1.x.scriptrunner-and-update-target.yml`
   - Status: implemented

8. E1.4 — Sem os.Exit no core do runtime
   - Stories:
     - NFR cross-cutting em E1.x
   - Arquitetura:
     - [`docs/architecture/introduction.md#3-constraints-tecnicas-inviolaveis`](docs/architecture/introduction.md:38)
     - [`docs/architecture/error-handling-strategy.md`](docs/architecture/error-handling-strategy.md:1)
   - Implementação:
     - verificação sobre `internal/**`
   - Testes:
     - security/check: scan automatizado de os.Exit fora da casca CLI
   - Gate QA:
     - `docs/qa/gates/1.x.no-osexit-core.yml`
   - Status: implemented

---

## Bloco E1.5 — Modal Stack + Security (JIT + Anti-Trojan YAML)

13. E1.5 — Modal Stack única integrada ao runtime
    - Stories:
      - [`docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md`](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1)
    - Arquitetura:
      - [`docs/architecture/components.md#6-modal-stack--e15`](docs/architecture/components.md:410)
      - [`docs/architecture/core-workflows.md#workflow-3--seguranca-jit-com-modal-stack`](docs/architecture/core-workflows.md:135)
    - Implementação:
      - `internal/runtime/modal/...`
      - integração com `internal/runtime/event` e `internal/runtime/layout`
    - Testes:
      - unit: operações push/pop/top na pilha
      - integration: foco apenas no topo, plano de fundo bloqueado
    - Gate QA:
      - `docs/qa/gates/1.x.modal-stack.yml`
    - Status: implemented

14. E1.5 — Segurança JIT + Anti-Trojan YAML integrada ao ScriptRunner
    - Stories:
      - [`docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md`](docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md:1)
    - Arquitetura:
      - [`docs/architecture/security.md`](docs/architecture/security.md:1)
      - [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)
    - Implementação:
      - `internal/runtime/runner/...` (normalizePolicy, validateRunAction, isAllowedScript)
    - Testes:
      - security: YAML malicioso, chaves desconhecidas, comandos fora da whitelist
    - Gate QA:
      - `docs/qa/gates/1.x.security-jit-anti-trojan.yml`
    - Status: implemented
      - [`docs/architecture/data-models.md#4-onhandler-e13--bloco-on-como-unica-orquestracao`](docs/architecture/data-models.md:116)
    - Implementação:
      - `internal/runtime/event/...`
    - Testes:
      - unit: mapeamento tea.Msg → ShantillyEvent
      - unit: resolução determinística de regras `on:`
      - integration: fluxo UI → ShantillyEvent → OnHandler → ação declarativa
    - Gate QA:
      - planejado `docs/qa/gates/1.x.event-engine.md`
    - Status: planned

15. E1.4 — ScriptRunner básico (run:, args, stdin, update_target)
    - Stories:
      - futuras stories E1.4 (ScriptRunner)
    - Arquitetura:
      - [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)
      - [`docs/architecture/security.md#3-regras-para-run-e-scriptrunner-e14--e15`](docs/architecture/security.md:93)
    - Implementação:
      - `internal/runtime/runner/...`
    - Testes:
      - unit: montagem de comando a partir de RunAction (args/stdin/env)
      - integration: update_target recebendo saída do script
    - Gate QA:
      - planejado `docs/qa/gates/1.x.scriptrunner-and-update-target.md`
    - Status: planned

16. E1.4 — Regra 1 processo por update_target
    - Stories:
      - parte das stories E1.4
    - Arquitetura:
      - [`docs/architecture/security.md#5-regra-de-1-processo-por-update_target-como-invariante-de-seguranca-e14`](docs/architecture/security.md:167)
    - Implementação:
      - `internal/runtime/runner/...`
    - Testes:
      - integration + security: novo run com o mesmo update_target cancela o anterior
    - Gate QA:
      - incluído em `docs/qa/gates/1.x.scriptrunner-and-update-target.md`
    - Status: planned

17. Cross-cutting — Sem os.Exit no core do runtime
    - Stories:
      - NFRs de erro/governança em stories de runtime
    - Arquitetura:
      - [`docs/architecture/error-handling-strategy.md`](docs/architecture/error-handling-strategy.md:1)
      - [`docs/architecture/security.md#6-proibicao-de-osexit-no-core-cross-cutting`](docs/architecture/security.md:192)
    - Implementação:
      - verificação sobre `internal/runtime/**`
    - Testes:
      - security/check: scan automatizado garantindo ausência de `os.Exit` fora de `cmd/shantilly`
    - Gate QA:
      - planejado `docs/qa/gates/1.x.no-osexit-core.md`
    - Status: planned

18. Cross-cutting — Modal Stack + Segurança JIT (integração futura)
    - Stories:
      - futuras stories de Modal Stack e prompts seguros
    - Arquitetura:
      - [`docs/architecture/components.md#modal-stack--jit-security`](docs/architecture/components.md:1)
      - [`docs/architecture/security.md#4-modal-stack-e-seguranca-jit-e15`](docs/architecture/security.md:132)
    - Implementação:
      - `internal/runtime/modal/...`
      - integração com `internal/runtime/event` e `internal/runtime/runner`
    - Testes:
      - integration: foco apenas no topo da pilha, confirm/prompt_secrets condicionando `RunAction`
      - security: segredos não logados, não persistidos
    - Gate QA:
      - planejado `docs/qa/gates/1.x.modal-stack.md`
    - Status: planned

9. E1.5 — Modal Stack funcional integrada ao LayoutManager/EventManager
    - Stories:
      - stories Modal Stack (bmad-analyst)
    - Arquitetura:
      - [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)
      - [`docs/architecture/components.md#modal-stack--jit-security`](docs/architecture/components.md:1)
    - Implementação:
      - `internal/runtime/modal/...`
    - Testes:
      - integration: fluxo com múltiplos modais, foco apenas no topo
    - Gate QA:
      - novo gate `docs/qa/gates/1.x.modal-stack.yml`
    - Status: planned

10. E1.5 — Segurança JIT + Anti-Trojan YAML (schema, deny by default, whitelist)
    - Stories:
      - stories de segurança E1.5 (bmad-analyst)
    - Arquitetura:
      - [`docs/architecture/security.md`](docs/architecture/security.md:1)
      - [`docs/prd/epic-1-runtime-tui-foundation.md#2-restricoes-e-decisoes-estrategicas`](docs/prd/epic-1-runtime-tui-foundation.md:210)
    - Implementação:
      - módulo de validação YAML v2.0 (ex.: `pkg/declarative/validate.go`)
      - integração com `ScriptRunner`
    - Testes:
      - security: YAML malicioso, chaves desconhecidas, comandos fora da whitelist
    - Gate QA:
      - novo gate `docs/qa/gates/1.x.security-jit-anti-trojan.yml`
    - Status: planned

---

## Confinamento do Legado v1.0 (Cross-cutting)

11. Legado v1.0 confinado ao FormComponent
    - Stories:
      - marcadas como `legacy_refactored` no story map
    - Arquitetura:
      - [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1)
      - [`docs/architecture/components.md`](docs/architecture/components.md:1)
    - Implementação:
      - `internal/components/form_component.go`
      - referências controladas a `internal/tui` apenas como backend do componente
    - Testes:
      - regression: comportamento do formulário sob o runtime declarativo
    - Gate QA:
      - novo gate `docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml`
    - Status: planned

---

## Uso pela bmad-qa, bmad-dev, bmad-architect e bmad-orchestrator

- bmad-qa:
  - Deve:
    - Preencher/atualizar o Status de cada linha.
    - Criar/ajustar gates em `docs/qa/gates/` conforme indicado.
    - Garantir que não existam blocos E1.1–E1.5 sem pelo menos uma linha com testes planejados.

- bmad-dev/bmad-master:
  - Devem:
    - Vincular implementações às linhas desta matriz (via comentários, PR descriptions, etc.).
    - Não marcar como `implemented` sem código + testes em place.

- bmad-architect:
  - Deve:
    - Garantir que cada linha referencie seções arquiteturais válidas.
    - Ajustar arquitetura quando gaps forem encontrados (nunca o contrário).

- bmad-orchestrator:
  - Usa esta matriz como instrumento de auditoria:
    - PRs e releases só avançam se:
      - Linhas relevantes estão em `implemented` + testes `passed`.
      - Não há referências ativas ao modelo v1.x como produto/roadmap.

---

## Próximo passo após este ciclo de QA

Ao final da consolidação desta matriz (todas as linhas relevantes em `implemented` ou `passed`):

- O Epic 1 será:
  - Totalmente rastreável (PRD → Stories → Arquitetura → Implementação → QA).
  - Validado contra as constraints críticas:
    - YAML único,
    - Runtime declarativo,
    - Legado confinado,
    - Sem `os.Exit` no core,
    - 1 processo por `update_target`,
    - Modal Stack,
    - Segurança JIT + anti-Trojan.
- Próximo passo:
  - bmad-orchestrator formaliza, em `docs/architecture/governance-runtime-tui-v2.0.md`, os checklists automatizáveis e critérios de promoção para Epics 2+ sobre esta fundação estável.

Esta matriz deve ser mantida viva em cada PR e release como parte do ciclo de governança contínua.
