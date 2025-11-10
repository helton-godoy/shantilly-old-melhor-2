# Backlog Consolidado — Waves 4–7 (Epic 1: Runtime TUI Declarativo — Fundação do Runtime v2.0)

Este backlog é o artefato único e vinculante para execução das Waves 4–7 (E1.4–E1.5) do Epic 1.  
Ele consolida histórias, referências normativas e proibições, sem reabrir decisões de produto ou arquitetura.

Premissas imutáveis (pivotadas de:

- Story Map: [docs/stories/index-epic-1-runtime-tui.story-map.md](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
- Implementation Plan: [docs/architecture/implementation-plan-epic-1-runtime-tui.md](docs/architecture/implementation-plan-epic-1-runtime-tui.md:1)
- Governança: [docs/architecture/governance-runtime-tui-v2.0.md](docs/architecture/governance-runtime-tui-v2.0.md:1)
- Playbook Agentes: [docs/architecture/bmad-agents-playbook.md](docs/architecture/bmad-agents-playbook.md:1)
- Matriz QA: [docs/qa/matrix-epic-1-runtime-tui-coverage.md](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)
):

- Binário único: `shantilly`.
- Contrato único: `AppConfig` declarativo v2.0 (YAML único).
- Pipeline único: `ShantillyEvent` → `EventManager` → `on:` → (`RunAction` | `ModalRequest` | `UpdateState`).
- Alvos exclusivos:
  - ScriptRunner: `internal/runtime/runner/**`.
  - Modal Stack: `internal/runtime/modal/**`.
  - Sandbox legado/FormComponent.
  - Segurança JIT + Anti-Trojan YAML sobre o mesmo contrato.
- Proibições vinculantes:
  - Sem novo binário.
  - Sem múltiplos contratos YAML.
  - Sem rotas paralelas de eventos/execução/modal/state.
  - Sem `os.Exit` no core (apenas em [cmd/shantilly/main.go](cmd/shantilly/main.go:1)).
  - Sem uso novo direto de `FormConfig`/`internal/tui` fora do sandbox FormComponent.

Cada item abaixo inclui:

- ID E1.x
- Descrição de história/backlog
- Referências cruzadas (PRD / Arquitetura / QA / Gates)
- Flag: `SCOPED` (aderente) ou `REJECTED` (fora de escopo por violar contrato único)

## 1. Wave 4 — ScriptRunner (E1.4)

### E1.4-1 — ScriptRunner único para `run:` declarativo

- ID: E1.4-1
- Tipo: Feature core
- Descrição:
  - Implementar `ScriptRunner` único em `internal/runtime/runner/**` para executar ações `run:` definidas no `AppConfig`.
  - Entrada apenas via `RunAction` emitido pelo `EventManager` com base em `on:`.
- Referências:
  - PRD: [docs/prd/epic-1-runtime-tui-foundation.md](docs/prd/epic-1-runtime-tui-foundation.md:1)
  - Arquitetura:
    - [docs/architecture/core-workflows.md](docs/architecture/core-workflows.md:1)
    - [docs/architecture/components.md#5-scriptrunner--e14](docs/architecture/components.md:256)
    - [docs/architecture/security.md#3-regras-para-run-e-scriptrunner-e14--e15](docs/architecture/security.md:93)
    - [docs/architecture/data-models.md](docs/architecture/data-models.md:1)
  - Stories:
    - [docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md](docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md:1)
  - QA / Gates:
    - Matriz: bloco E1.4
    - Gate: [docs/qa/gates/1.x.scriptrunner-and-update-target.yml](docs/qa/gates/1.x.scriptrunner-and-update-target.yml:1)
    - Gate: [docs/qa/gates/1.x.no-osexit-core.yml](docs/qa/gates/1.x.no-osexit-core.yml:1)
- Flags:
  - SCOPED
  - Sem rotas paralelas / sem novo binário / sem múltiplos contratos.

### E1.4-2 — Suporte completo a args/stdin/env + binding declarativo

- ID: E1.4-2
- Tipo: Feature detalhamento
- Descrição:
  - Permitir que `RunAction` defina `args`, `stdin` e (quando previsto) `env` com base em estado/entradas declarativas.
  - Montagem de comando é responsabilidade do ScriptRunner, não de componentes ou LayoutManager.
- Referências:
  - Arquitetura:
    - [docs/architecture/data-models.md#6-runaction-e14--execucao-declarativa](docs/architecture/data-models.md:172)
    - [docs/architecture/core-workflows.md](docs/architecture/core-workflows.md:1)
  - QA:
    - Matriz E1.4 — vincular casos de binding seguro.
    - Gate: `1.x.scriptrunner-and-update-target.yml`.
- Flags:
  - SCOPED

### E1.4-3 — Regra de 1 processo por `update_target`

- ID: E1.4-3
- Tipo: Invariante técnico
- Descrição:
  - Garantir 1 processo ativo por `update_target`.
  - Antes de iniciar novo `RunAction` para o mesmo `update_target`, encerrar o anterior (SIGTERM → SIGKILL).
- Referências:
  - Arquitetura:
    - [docs/architecture/high-level-architecture.md](docs/architecture/high-level-architecture.md:1)
    - [docs/architecture/security.md#5-regra-de-1-processo-por-update_target-como-invariante-de-seguranca-e14](docs/architecture/security.md:167)
  - QA:
    - Matriz: E1.4 lifecycle.
    - Gate: `1.x.scriptrunner-and-update-target.yml`.
- Flags:
  - SCOPED

### E1.4-4 — Emissão de resultados via `ShantillyEvent` + updates declarativos

- ID: E1.4-4
- Tipo: Comportamento normativo
- Descrição:
  - ScriptRunner nunca fala direto com UI/LM.
  - Outputs e erros retornam como `ShantillyEvent` e/ou atualizações declarativas de `update_target` consumidas pelo runtime.
- Referências:
  - Arquitetura:
    - [docs/architecture/core-workflows.md](docs/architecture/core-workflows.md:1)
    - [docs/architecture/components.md#5-scriptrunner--e14](docs/architecture/components.md:256)
  - QA:
    - Gate: `1.x.event-engine.yml`
    - Gate: `1.x.scriptrunner-and-update-target.yml`
- Flags:
  - SCOPED

### E1.4-R1 — Execução direta de scripts por componentes/LM

- ID: E1.4-R1
- Descrição:
  - Qualquer história que proponha componentes, LayoutManager ou código legado executando scripts diretamente (ex.: `huh` chamando `exec.Command` sem passar por ScriptRunner).
- Status:
  - REJECTED / FORA DE ESCOPO
- Motivo:
  - Viola pipeline único `ShantillyEvent` → `EventManager` → `on:` → `RunAction` → ScriptRunner.
  - Fere governança: [docs/architecture/governance-runtime-tui-v2.0.md](docs/architecture/governance-runtime-tui-v2.0.md:31).
  - Fere playbook: [docs/architecture/bmad-agents-playbook.md](docs/architecture/bmad-agents-playbook.md:169).

### E1.4-R2 — Múltiplos ScriptRunners ou runners especializados diretos na Wave 4

- ID: E1.4-R2
- Descrição:
  - Introduzir mais de um runner genérico ou runners especializados (ex.: `ansible`) antes de consolidar o ScriptRunner único.
- Status:
  - REJECTED / FORA DE ESCOPO (para Epic 1 / Waves 4–7)
- Motivo:
  - Quebra contrato único de execução.
  - Runners especializados pertencem a épicos futuros, sobre o mesmo pipeline.

## 2. Wave 5 — Modal Stack (E1.5 parte 1)

### E1.5-1 — Modal Stack única em `internal/runtime/modal/**`

- ID: E1.5-1
- Tipo: Feature core
- Descrição:
  - Implementar Modal Stack como pilha explícita (`push` / `pop` / `top`) em `internal/runtime/modal/**`.
  - Nenhuma outra estrutura paralela de modais.
- Referências:
  - Arquitetura:
    - [docs/architecture/components.md#modal-stack--jit-security](docs/architecture/components.md:288)
    - [docs/architecture/core-workflows.md#workflow-3--seguranca-jit-com-modal-stack](docs/architecture/core-workflows.md:135)
  - Stories:
    - [docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1)
  - QA / Gates:
    - Matriz E1.5 — Modal Stack.
    - Gate: [docs/qa/gates/1.x.modal-stack.yml](docs/qa/gates/1.x.modal-stack.yml:1)
- Flags:
  - SCOPED

### E1.5-2 — Foco exclusivo no topo + bloqueio do plano de fundo

- ID: E1.5-2
- Tipo: Invariante UX/segurança
- Descrição:
  - Quando houver modal na pilha:
    - Apenas o topo recebe eventos.
    - Interações com plano de fundo são bloqueadas.
- Referências:
  - Arquitetura:
    - [docs/architecture/components.md#modal-stack--jit-security](docs/architecture/components.md:288)
  - QA:
    - Matriz: E1.5.
    - Gate: `1.x.modal-stack.yml`.
- Flags:
  - SCOPED

### E1.5-3 — Integração Modal Stack ↔ LayoutManager ↔ EventManager

- ID: E1.5-3
- Tipo: Integração
- Descrição:
  - `LayoutManager` desenha o modal de topo sobre o layout base.
  - `EventManager`:
    - Abre modais via `ModalRequest` (vindo de `on:`).
    - Encaminha eventos apenas ao topo da pilha.
- Referências:
  - Arquitetura:
    - [docs/architecture/core-workflows.md](docs/architecture/core-workflows.md:1)
  - QA:
    - Matriz: linhas de integração Modal Stack.
- Flags:
  - SCOPED

### E1.5-4 — `confirm` e `prompt_secrets` sempre via Modal Stack

- ID: E1.5-4
- Tipo: Segurança JIT (parte 1)
- Descrição:
  - `RunAction.Confirm: true` e `RunAction.PromptSecrets`:
    - Devem gerar `ModalRequest` via `on:`.
    - Execução do `run:` só após confirmação e/ou captura segura dos segredos pela Modal Stack.
- Referências:
  - Arquitetura:
    - [docs/architecture/security.md#4-modal-stack-e-seguranca-jit-e15](docs/architecture/security.md:132)
    - [docs/architecture/core-workflows.md#workflow-3--seguranca-jit-com-modal-stack](docs/architecture/core-workflows.md:135)
  - Stories:
    - [docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1)
- QA:
  - Gates: `1.x.modal-stack.yml`, `1.x.security-jit-anti-trojan.yml`.
- Flags:
  - SCOPED

### E1.5-R1 — Modais soltos em componentes/LM/legado

- ID: E1.5-R1
- Descrição:
  - Qualquer uso de modais criados diretamente por componentes, LayoutManager ou código legado fora de `ModalRequest` + Modal Stack.
- Status:
  - REJECTED / FORA DE ESCOPO
- Motivo:
  - Viola governança (modal stack única) e pipeline único.

### E1.5-R2 — Prompt de segredos direto no CLI ou via `fmt.Scan` sem Modal Stack

- ID: E1.5-R2
- Status:
  - REJECTED / FORA DE ESCOPO
- Motivo:
  - Burla Segurança JIT e a proteção de segredos definida em:
    - [docs/architecture/security.md](docs/architecture/security.md:1)
    - [docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1)

## 3. Wave 6 — Encapsulamento do Legado (FormComponent Sandbox)

### E1.6-1 — FormComponent como único ponto de legado

- ID: E1.6-1
- Tipo: Cross-cutting
- Descrição:
  - Encapsular v1.x (`FormConfig`, `internal/tui`, `shantilly form`) em um único `FormComponent` configurado via `AppConfig`.
  - Nenhuma exposição pública do modelo legado.
- Referências:
  - Arquitetura:
    - [docs/architecture/introduction.md](docs/architecture/introduction.md:1)
    - [docs/architecture/components.md](docs/architecture/components.md:1)
  - Story:
    - [docs/stories/6.1.encapsulamento-formcomponent-legado.story.md](docs/stories/6.1.encapsulamento-formcomponent-legado.story.md:1)
  - QA:
    - Matriz: seção Confinamento do Legado.
    - Gate: [docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml](docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml:1)
- Flags:
  - SCOPED

### E1.6-2 — Legado só emite `ShantillyEvent` e usa pipeline oficial

- ID: E1.6-2
- Descrição:
  - Código legado dentro do FormComponent:
    - Não executa scripts diretamente.
    - Não abre modais diretamente.
    - Não altera estado fora do fluxo `ShantillyEvent` → `EventManager` → `on:`.
- Referências:
  - Governança:
    - [docs/architecture/governance-runtime-tui-v2.0.md#5-wave-6--encapsulamento-do-legado-formcomponent](docs/architecture/governance-runtime-tui-v2.0.md:63)
- QA:
  - Gate: `1.x.legacy-formcomponent-encapsulation.yml`.
- Flags:
  - SCOPED

### E1.6-R1 — Reativar `shantilly form` como produto/entrada principal

- ID: E1.6-R1
- Status:
  - REJECTED / FORA DE ESCOPO
- Motivo:
  - Contraria story map e governance que posicionam v1.x como LEGACY encapsulado.

### E1.6-R2 — Novo código usando diretamente `internal/tui`/`FormConfig`

- ID: E1.6-R2
- Status:
  - REJECTED / FORA DE ESCOPO
- Motivo:
  - Viola confinamento do legado e será bloqueado por gate dedicado.

## 4. Wave 7 — Segurança JIT + Anti-Trojan YAML (E1.5 parte 2)

### E1.5-5 — Validador forte de `AppConfig` (deny-by-default)

- ID: E1.5-5
- Tipo: Segurança
- Descrição:
  - Implementar validação rígida do `AppConfig`:
    - `deny-by-default` para chaves desconhecidas/tipos inválidos.
    - Erros são bloqueantes (não apenas warnings).
- Referências:
  - Arquitetura:
    - [docs/architecture/security.md](docs/architecture/security.md:1)
  - Story:
    - [docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md](docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md:1)
  - QA / Gates:
    - [docs/qa/gates/1.x.security-jit-anti-trojan.yml](docs/qa/gates/1.x.security-jit-anti-trojan.yml:1)
- Flags:
  - SCOPED

### E1.5-6 — Whitelist explícita de tipos/ações/paths + SecurityPolicy

- ID: E1.5-6
- Tipo: Segurança
- Descrição:
  - Definir `SecurityPolicy` central:
    - Whitelist de:
      - Tipos de componentes oficiais,
      - Ações `run:` suportadas,
      - Paths/comandos seguros.
  - ScriptRunner só executa se conforme à `SecurityPolicy`.
- Referências:
  - Arquitetura:
    - [docs/architecture/security.md](docs/architecture/security.md:1)
    - [docs/architecture/data-models.md#8-securitypolicy-e15--anti-trojan-yaml](docs/architecture/data-models.md:237)
  - QA:
    - Gate: `1.x.security-jit-anti-trojan.yml`.
- Flags:
  - SCOPED

### E1.5-7 — Anti-Trojan YAML: bloqueio de padrões maliciosos

- ID: E1.5-7
- Tipo: Segurança
- Descrição:
  - Bloquear:
    - Auto-runs implícitos,
    - Ações ocultas/não documentadas,
    - Comandos fora de whitelist,
    - Componentes/tipos não oficiais,
    - Paths suspeitos.
- Referências:
  - Story:
    - [docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md](docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md:1)
  - Governança:
    - [docs/architecture/governance-runtime-tui-v2.0.md#6-wave-7--hardening-de-seguranca-e15--parte-2](docs/architecture/governance-runtime-tui-v2.0.md:75)
- QA:
  - Matriz: E1.5 Security.
  - Gate: `1.x.security-jit-anti-trojan.yml`.
- Flags:
  - SCOPED

### E1.5-8 — Proteção de segredos (JIT, em memória, nunca logados)

- ID: E1.5-8
- Tipo: Segurança JIT
- Descrição:
  - Segredos coletados via Modal Stack:
    - Vivem apenas em memória.
    - Não são logados.
    - Não são persistidos.
    - Usados somente para compor `RunAction` no fluxo seguro.
- Referências:
  - Arquitetura:
    - [docs/architecture/security.md#4-modal-stack-e-seguranca-jit-e15](docs/architecture/security.md:132)
  - Stories:
    - [docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1)
- Flags:
  - SCOPED

### E1.5-R3 — Chaves/yaml “flexíveis” não validadas ou permissivas

- ID: E1.5-R3
- Descrição:
  - Qualquer proposta de:
    - Permitir chaves arbitrárias,
    - “Modo compatível” que desliga validação para facilitar uso,
    - Execuções `run:` sem passar pela `SecurityPolicy`.
- Status:
  - REJECTED / FORA DE ESCOPO
- Motivo:
  - Quebra `deny-by-default` e anti-Trojan YAML.

### E1.5-R4 — Logging ou persistência de segredos

- ID: E1.5-R4
- Status:
  - REJECTED / FORA DE ESCOPO
- Motivo:
  - Viola requisitos explícitos de segurança JIT e proteção de segredos.

## 5. Critério de Saída deste Backlog (Ciclo 1 — Waves 4–7)

Para bmad-architect iniciar o Ciclo 2 sem reabrir decisões, este backlog:

- Fixa, em um único documento, as histórias e invariantes para:
  - ScriptRunner (E1.4),
  - Modal Stack (E1.5 parte 1),
  - Segurança JIT + Anti-Trojan YAML (E1.5 parte 2),
  - Encapsulamento do legado (Wave 6).
- Marca explicitamente como REJECTED:
  - Qualquer história que:
    - Introduza novo binário,
    - Crie múltiplos contratos ou YAML alternativo,
    - Bypasse `ShantillyEvent`/`EventManager`/`on:`/ScriptRunner/Modal Stack,
    - Reative ou expanda uso direto do legado v1.x fora do sandbox,
    - Enfraqueça validação ou segurança JIT.

Este arquivo deve ser usado como:

- Fonte única de backlog de Waves 4–7 para:
  - bmad-pm/bmad-po (priorização e PRD),
  - bmad-architect (contratos formais),
  - bmad-dev/bmad-master (implementação estrita),
  - bmad-qa (matriz + gates),
  - bmad-orchestrator (enforcement).
- Complemento normativo à governança existente, sem a reabrir.
