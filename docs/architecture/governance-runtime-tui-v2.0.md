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

## 7. QA, Matriz, Gates e Enforcement (BMAD Operacional)

### 7.1. Matriz como fonte única

- A matriz [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1) é:
  - Fonte única de verdade para status `planned` / `implemented` / `passed` de cada cláusula E1.x.
  - Obrigatória para:
    - bmad-architect: garantir alinhamento com arquitetura.
    - bmad-dev/bmad-master: vincular implementações e testes.
    - bmad-qa: registrar evidências e gates.
    - bmad-orchestrator: auditar PRs e releases.

Regras normativas:

- Nenhum item de runtime/config/segurança é considerado “Done” se:
  - Não houver linha correspondente na matriz.
  - O gate associado não estiver em `PASS` (ou waiver formal registrado).
  - Não houver evidência vinculada (tests/docs/artefatos) referenciada na matriz ou no gate.

### 7.2. Gates QA como contratos executáveis

- Todos os gates em `docs/qa/gates/*.yml` devem seguir o modelo normativo de 2.1:
  - Exemplo: [`docs/qa/gates/2.1-advanced-form-types-validation.yml`](docs/qa/gates/2.1-advanced-form-types-validation.yml:1).
  - Campos obrigatórios mínimos:
    - `schema`
    - `gate` (ex.: `PASS`, `FAIL`, `WAIVED`)
    - `status_reason`
    - `reviewer`
    - `updated`
    - `waiver` (quando aplicável)
    - `evidence` (tests, trace, links)
- Para Epic 1 (Waves 2–7), os gates normativos bloqueantes incluem, no mínimo:
  - [`docs/qa/gates/1.x.layout-manager.yml`](docs/qa/gates/1.x.layout-manager.yml:1)
  - [`docs/qa/gates/1.x.event-engine.yml`](docs/qa/gates/1.x.event-engine.yml:1)
  - [`docs/qa/gates/1.x.scriptrunner-and-update-target.yml`](docs/qa/gates/1.x.scriptrunner-and-update-target.yml:1)
  - [`docs/qa/gates/1.x.modal-stack.yml`](docs/qa/gates/1.x.modal-stack.yml:1)
  - [`docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml`](docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml:1)
  - [`docs/qa/gates/1.x.no-osexit-core.yml`](docs/qa/gates/1.x.no-osexit-core.yml:1)
  - [`docs/qa/gates/1.x.security-jit-anti-trojan.yml`](docs/qa/gates/1.x.security-jit-anti-trojan.yml:1)

### 7.3. Consumo na CI — Job normativo único

- Deve existir um workflow único de governança, por exemplo:
  - `governanca-waves4-7 / Governança Waves 4-7 (Gates 1.x + Go checks)`.
- Este workflow é o trilho normativo:
  - Obrigatório como status check bloqueante para branches de proteção (`main`/`master`).
  - Nenhum merge é permitido sem este check em verde.

Mínimos obrigatórios deste workflow:

- Go:
  - `go test ./...`
  - `go vet ./...`
  - `golangci-lint run ./...`
- Gates automatizados:
  - Gate 1.x.no-osexit-core — scan de `os.Exit` fora de [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1).
  - Gate 1.x.legacy-formcomponent-encapsulation — incluindo [`scripts/check-legacy-encapsulation.sh`](scripts/check-legacy-encapsulation.sh:1).
  - Gate 1.x.scriptrunner-and-update-target — validação de runner único + 1 processo/update_target.
  - Gate 1.x.modal-stack — validação de pilha única e foco no topo.
  - Gate 1.x.security-jit-anti-trojan — validações deny-by-default + whitelist.
  - Gate 1.x.layout-manager — proibição de automação/legado em `internal/runtime/layout/**`.
  - Gate 1.x.event-engine — proibição de rotas paralelas fora de `ShantillyEvent` → [`EventManager`](internal/runtime/event/manager.go:1) → `on:`.

Regras:

- Gates QA devem ser consumidos pela CI:
  - via scripts/checks que leem `docs/qa/gates/*.yml` e falham o job se:
    - gate relevante não estiver em `PASS` (ou waiver ativo e justificado),
    - ou critérios obrigatórios não forem atendidos.
- Nenhum outro workflow pode ser usado como bypass:
  - Somente o job normativo pode representar “governança Waves 4–7”.
  - Qualquer automação paralela é auxiliar, nunca substituta.

### 7.4. Regras de promoção de PR e release

Para qualquer PR que altere runtime, AppConfig, segurança ou fluxos E1.x:

- Deve conter:
  - Referência às linhas relevantes na matriz.
  - Referência aos gates QA correspondentes.
  - Links para evidências de testes (unit/integration/e2e/security) e docs.
- Condições para aprovação:
  - Todos os gates aplicáveis ao escopo do PR em `PASS` ou waiver formal registrado no próprio gate YAML.
  - CI `governanca-waves4-7` em verde.
- Condições para release:
  - Para o escopo incluído:
    - Linhas da matriz em `implemented` ou `passed`.
    - Gates relevantes em `PASS`.
    - Nenhum desvio conhecido dos contratos desta governança sem waiver aprovado.

Nenhum PR que viole essas regras deve ser aceito. Nenhum merge/release pode declarar “Done” para itens de runtime sem rastreabilidade completa Matriz → Gates → Evidências.

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
