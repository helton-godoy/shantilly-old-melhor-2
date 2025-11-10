# Governance — Runtime TUI Declarativo v2.0 (Waves 2–7)

Este documento consolida as regras vinculantes de governança para o próximo ciclo, assumindo Waves 2–3 como baseline ativa e Waves 4–7 como trilhos normativos imutáveis.

## 1. Contratos Estruturais Imutáveis

- Binário único:
  - Somente `shantilly` é permitido como entrada oficial.
- Configuração única:
  - Apenas `AppConfig` declarativo é fonte de verdade.
  - Qualquer uso direto de `FormConfig`/`shantilly form` é considerado LEGACY e só é aceitável dentro do sandbox `FormComponent`.
- Roteamento único:
  - Toda lógica passa por:
    - `ShantillyEvent` → [`EventManager`](internal/runtime/event/manager.go:1) → `on:` → (`RunAction` | `ModalRequest` | `UpdateState`).
  - É proibido qualquer fluxo paralelo.

## 2. Wave 2–3 — Baseline Normativa (já obrigatória)

- [`LayoutManager`](internal/runtime/layout/manager.go:1):
  - Responsável apenas pelo layout.
  - Proibido executar scripts, abrir modais ou acessar contratos v1.x.
- [`EventManager`](internal/runtime/event/manager.go:1) + `on:`:
  - Única orquestração de automação.
  - Só podem emitir:
    - `RunAction`,
    - `ModalRequest`,
    - `UpdateState`.

Qualquer violação é bug arquitetural.

## 3. Wave 4 — ScriptRunner (E1.4)

- Implementação exclusiva:
  - Diretório `internal/runtime/runner/**`.
- Regras:
  - Apenas consome `RunAction` gerado via `on:`.
  - 1 processo por `update_target`:
    - Encerrar anterior (SIGTERM → SIGKILL) antes do próximo.
  - Proibido `os.Exit` em `internal/runtime/**`:
    - Permitido apenas em [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1).
  - Proibida qualquer execução direta fora de `RunAction`/ScriptRunner.
- Saída:
  - Sempre via `ShantillyEvent` e/ou update declarativo do `update_target`.

## 4. Wave 5 — Modal Stack + Segurança JIT (E1.5 — parte 1)

- Implementação exclusiva:
  - Diretório `internal/runtime/modal/**`.
- Regras:
  - Pilha explícita (push/pop/top).
  - Apenas o topo recebe foco/eventos.
  - Enquanto houver modal:
    - Fundo bloqueado.
  - Abertura de modais:
    - Sempre: `ShantillyEvent` → `EventManager` → `on:` → `ModalRequest` → Modal Stack.
    - Proibidos modais “soltos” em qualquer engine/componente.
- Segurança JIT:
  - `RunAction.Confirm` e `RunAction.PromptSecrets`:
    - Devem acionar Modal Stack antes da execução.
  - Segredos:
    - Apenas em memória, nunca logados ou persistidos.

## 5. Wave 6 — Encapsulamento do Legado (FormComponent)

- v1.x é sempre LEGACY.
- `FormComponent`:
  - Único ponto de entrada para reaproveitar `FormConfig`/`internal/tui`.
  - Configurado via `AppConfig` v2.0 (`Props`), nunca expõe `FormConfig` como contrato público.
  - Emite apenas `ShantillyEvent`.
- Proibições:
  - Nenhum código novo pode depender diretamente de `FormConfig` ou `internal/tui`.
  - Legado não pode executar scripts ou abrir modais diretamente:
    - Sempre via `EventManager` + `on:` + Modal Stack + ScriptRunner.

## 6. Wave 7 — Hardening de Segurança (E1.5 — parte 2)

- AppConfig:
  - Validação forte:
    - `deny by default` para chaves desconhecidas/tipos inválidos.
    - Whitelist para:
      - Tipos de componentes,
      - Ações `run:` permitidas,
      - Paths seguros.
  - Erros de violação devem ser bloqueantes.
- SecurityPolicy:
  - Integração obrigatória com ScriptRunner:
    - Nenhuma execução fora da política aprovada.
- Anti-Trojan YAML:
  - Bloqueio de:
    - Auto-runs implícitos,
    - Ações ocultas,
    - Componentes/tipos não documentados,
    - Paths suspeitos.

## 7. QA, Gates e Enforcement

As regras acima são condições bloqueantes para PRs e são executadas exclusivamente via o job canônico:

- Job único de governança Waves 4-7:
  - `governanca-waves4-7 / Governança Waves 4-7 (Gates 1.x + Go checks)` é o único job de governança Waves 4-7 obrigatório para `main`/`master`.
  - Merge em `main`/`master` sem este status check verde é proibido.
  - Todos os gates 1.x vinculantes devem ser implementados e executados dentro deste workflow único:
    - Go checks: `go test ./...`, `go vet ./...`, `golangci-lint run ./...`.
    - Gate 1.x.no-osexit-core.
    - Gate 1.x.legacy-formcomponent-encapsulation (incluindo [`scripts/check-legacy-encapsulation.sh`](scripts/check-legacy-encapsulation.sh:1)).
    - Gate 1.x.scriptrunner-and-update-target.
    - Gate 1.x.modal-stack.
    - Gate 1.x.security-jit-anti-trojan.
    - Gate 1.x.event-and-layout-pipeline.

- Matriz QA:
  - [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:151) deve mapear cada cláusula E1.x aos checks executados no job `governanca-waves4-7`.

- Gates obrigatórios em `docs/qa/gates/*.yml`:
  - Devem referenciar explicitamente o job `governanca-waves4-7` como veículo normativo de enforcement dos gates 1.x.
  - `1.x.layout-manager.yml` — garante LayoutManager sem automação.
  - `1.x.event-engine.yml` — garante roteamento único via EventManager + `on:`.
  - `1.x.scriptrunner-and-update-target.yml` — garante ScriptRunner único e 1 processo/update_target.
  - `1.x.modal-stack.yml` — garante Modal Stack única e foco correto.
  - `1.x.legacy-formcomponent-encapsulation.yml` — garante confinamento do legado.
  - `1.x.security-jit-anti-trojan.yml` — garante validação forte, deny-by-default e whitelists.
  - `1.x.no-osexit-core.yml` — proíbe `os.Exit` fora da casca CLI.

Nenhum PR que viole essas regras deve ser aceito, e nenhuma automação fora do job `governanca-waves4-7` pode ser considerada trilho de governança Waves 4-7 ou utilizada como bypass dos gates 1.x.

## 8. Diretriz para o Próximo Ciclo (Wave 4+)

- Ponto de partida:
  - Waves 2–3 já consolidadas e obrigatórias.
- Waves 4–7:
  - Devem ser implementadas exatamente nestes alvos:
    - `internal/runtime/runner/**` (ScriptRunner),
    - `internal/runtime/modal/**` (Modal Stack),
    - Sandbox `FormComponent` para legado.
  - São trilhos normativos, sem espaço para interpretação ambígua.
- Qualquer rota paralela:
  - Execução direta,
  - Modal fora da pilha,
  - Uso direto de v1.x,
  - Bypass de validação/segurança,
  é considerada violação crítica de governança.

Este documento é vinculante e serve como referência executiva para arquitetura, implementação, QA e automação de governança no próximo ciclo.
