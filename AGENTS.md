# AGENTS.md

This file provides guidance to agents when working with code in this repository.

- Use the single `shantilly` binary and single declarative `AppConfig` YAML as the only runtime contracts; any v1.x `FormConfig`/`shantilly form` usage is LEGACY and only allowed inside the future FormComponent sandbox.
- Layout and event flow are already NORMATIVE: [`LayoutManager`](internal/runtime/layout/manager.go:1) owns layout only (no automation), [`EventManager`](internal/runtime/event/manager.go:1) + `on:` own all automation (`RunAction`, `ModalRequest`, `UpdateState`) with no side paths.
- When implementing Waves 4–7, you MUST route all effects through this pipeline: `ShantillyEvent` → `EventManager` → `on:` → (`RunAction` | `ModalRequest` | `UpdateState`); direct script exec, modal opening or state mutation from components/LM/legacy code is forbidden.
- ScriptRunner (Wave 4) must live exclusively under `internal/runtime/runner/**` and:
  - Enforce 1 process por `update_target` (SIGTERM→SIGKILL no anterior antes do próximo).
  - Never call `os.Exit` (permitted only in [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1)); all errors fluem via tipos/retornos/eventos.
  - Emit resultados como `ShantillyEvent` + atualizações declarativas de `update_target`.
- Modal Stack (Wave 5) deve ser implementada apenas em `internal/runtime/modal/**` como pilha explícita:
  - Apenas o topo recebe foco/eventos; plano de fundo bloqueado com modal ativo.
  - `confirm`/`prompt_secrets` de `RunAction` sempre convertem em `ModalRequest` + eventos de resposta, nunca em prompts diretos.
- Segurança JIT + Anti-Trojan YAML (Wave 7):
  - Validação forte de `AppConfig` com `deny by default` para chaves desconhecidas e whitelist explícita de ações/paths seguros, integrada à política de segurança usada pelo ScriptRunner.
  - Segredos coletados via Modal Stack nunca podem ser logados/persistidos; mantenha-os apenas em memória e limite o escopo.
- Encapsulamento do legado (Wave 6):
  - Todo código v1.x vive como detalhe interno do FormComponent; nenhum novo fluxo usa diretamente `internal/tui` ou `FormConfig`.
  - Qualquer exceção a isso é bug arquitetural e deve ser caçada por QA/gates.
- Ao adicionar código ou docs:
  - Respeite sempre os contratos formais em [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1), [`docs/architecture/components.md`](docs/architecture/components.md:1), [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1), [`docs/architecture/security.md`](docs/architecture/security.md:1) e a matriz [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1); use IDs E1.x para rastreabilidade.
  - Antes de criar nova rota/componente, confirme se não viola: roteamento único via EventManager, ScriptRunner único, Modal Stack única, confinamento do legado e proibição de `os.Exit` no core.
