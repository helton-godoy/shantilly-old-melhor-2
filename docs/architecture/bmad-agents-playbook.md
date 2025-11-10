# BMAD Agents Playbook — Runtime TUI Declarativo v2.0 (Waves 2–7)

Este playbook operacional torna executável, para todos os agentes, o mandato vinculante definido em [`AGENTS.md`](AGENTS.md:1) e em [`docs/architecture/governance-runtime-tui-v2.0.md`](docs/architecture/governance-runtime-tui-v2.0.md:1), com foco na execução disciplinada das Waves 4–7.

Ele não reabre decisões. Apenas operacionaliza:

- Binário único: `shantilly`.
- Contrato único: `AppConfig` v2.0.
- Pipeline único: `ShantillyEvent` → [`EventManager`](internal/runtime/event/manager.go:1) → `on:` → (`RunAction` | `ModalRequest` | `UpdateState`).
- Alvos exclusivos:
  - [`ScriptRunner`](internal/runtime/runner) (Wave 4),
  - Modal Stack [`internal/runtime/modal`](internal/runtime/modal) (Wave 5),
  - Sandbox de legado/FormComponent (Wave 6),
  - Segurança JIT + Anti-Trojan YAML (Wave 7).
- Proibições críticas:
  - `os.Exit` fora de [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1),
  - Execução direta de scripts fora do ScriptRunner,
  - Modais “soltos” fora da Modal Stack,
  - Uso direto de v1.x/`FormConfig`/`internal/tui` fora do FormComponent sandboxado,
  - Qualquer bypass de validação/segurança.

---

## 1. bmad-orchestrator — Orquestrador Mestre (este agente)

Mandato:

- Guardião do arranjo imutável de contratos, trilhos e gates.
- Coordena ciclos; não implementa código, mas bloqueia desvios.

Responsabilidades operacionais:

1. Trabalhar sempre ancorado em:
   - [`AGENTS.md`](AGENTS.md:1),
   - [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1),
   - [`docs/architecture/governance-runtime-tui-v2.0.md`](docs/architecture/governance-runtime-tui-v2.0.md:1),
   - Stories normativas Waves 4–7:
     - [`docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md`](docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md:1),
     - [`docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md`](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1),
     - [`docs/stories/6.1.encapsulamento-formcomponent-legado.story.md`](docs/stories/6.1.encapsulamento-formcomponent-legado.story.md:1),
     - [`docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md`](docs/stories/7.1.hardening-seguranca-jit-anti-trojan-e15.story.md:1).
   - Matriz QA:
     - [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1).
2. Em cada novo ciclo Wave 4–7:
   - Verificar se:
     - As tasks abertas respeitam diretórios autorizados,
     - Os gates necessários existem ou estão planejados.
   - Direcionar outros agentes:
     - Architect/Dev/QA/PM/PO para os artefatos corretos.
3. Em análise de PR/release:
   - Rejeitar qualquer:
     - Novo caminho de execução fora do pipeline único,
     - Uso de `internal/tui` fora do sandbox,
     - Execução de script, modal, ou mutação de estado fora de `on:` + engines oficiais,
     - Relaxamento de validação/AppConfig ou whitelists.

Entrada esperada para atuar:

- Links para PRs/branches/commits.
- Referência explícita à Story/Bloco E1.x alvo.
- Indicação do agente chamador.

Saída esperada:

- Decisão clara:
  - "Conforme governança" ou "Bloquear" com razão precisa e linkada.

---

## 2. bmad-architect — Arquiteto

Mandato:

- Transformar PRD + Stories + Governança em contratos formais.
- Nunca fugir dos diretórios/engines e regras fixas.

Responsabilidades operacionais (Waves 4–7):

1. ScriptRunner (Wave 4, E1.4):
   - Manter contrato formal em:
     - [`docs/architecture/components.md#5-scriptrunner--e14`](docs/architecture/components.md:256),
     - [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1),
     - [`docs/architecture/security.md#3-regras-para-run-e-scriptrunner-e14--e15`](docs/architecture/security.md:93),
     - [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1).
   - Garantir alinhamento estrito com:
     - [`docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md`](docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md:1),
     - [`docs/qa/gates/1.x.scriptrunner-and-update-target.yml`](docs/qa/gates/1.x.scriptrunner-and-update-target.yml:1).

2. Modal Stack + Segurança JIT (Wave 5, E1.5 parte 1):
   - Especificar em:
     - [`docs/architecture/components.md#modal-stack--jit-security`](docs/architecture/components.md:288),
     - [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:135).
   - Alinhar com:
     - [`docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md`](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1),
     - [`docs/qa/gates/1.x.modal-stack.yml`](docs/qa/gates/1.x.modal-stack.yml:1).

3. Encapsulamento Legado (Wave 6):
   - Formalizar o sandbox FormComponent em:
     - [`docs/architecture/components.md`](docs/architecture/components.md:1),
     - [`docs/architecture/source-tree.md`](docs/architecture/source-tree.md:1), se aplicável.
   - Garantir:
     - Único ponto: `FormComponent` (type: form),
     - Nenhuma API pública baseada em `FormConfig`/`internal/tui`.

4. Hardening Segurança JIT + Anti-Trojan (Wave 7):
   - Detalhar:
     - Validador AppConfig (deny-by-default + whitelist),
     - `SecurityPolicy` integrada ao ScriptRunner.
   - Em:
     - [`docs/architecture/security.md`](docs/architecture/security.md:1),
     - [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1).

Entrada esperada:

- PRD/stories aprovadas,
- Gaps identificados pelo QA ou Dev.

Saída esperada:

- Documentos de arquitetura atualizados,
- Contratos inequívocos para Dev e QA.

---

## 3. bmad-qa — Test Architect & Quality Gates

Mandato:

- Transformar governança em gates automatizados bloqueantes.
- Cobrir Waves 2–7 com matriz + gates.

Responsabilidades operacionais:

1. Manter e estender:
   - [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1).
2. Implementar/ajustar gates em `docs/qa/gates/*.yml` para:
   - Event Engine:
     - `1.x.event-engine.yml` — fluxo único `ShantillyEvent` + `on:`.
   - Layout Manager:
     - `1.x.layout-manager.yml` — sem automação/execução no layout.
   - ScriptRunner:
     - `1.x.scriptrunner-and-update-target.yml` — executor único + 1 processo/update_target.
   - Modal Stack:
     - `1.x.modal-stack.yml` — única pilha, foco no topo, sem modais soltos.
   - Legado:
     - `1.x.legacy-formcomponent-encapsulation.yml` — v1.x confinado ao sandbox.
   - Segurança:
     - `1.x.no-osexit-core.yml` — `os.Exit` só em [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1),
     - `1.x.security-jit-anti-trojan.yml` — deny-by-default + whitelist.

3. Para cada PR:
   - Assegurar que CI:
     - Execute scans estáticos definidos nos gates,
     - Falhe em qualquer violação.

Entrada esperada:

- Contratos de arquitetura atualizados,
- Local-alvo de implementação (paths normativos),
- Casos de risco levantados.

Saída esperada:

- Gates YAML claros,
- Orientação aos devs sobre falhas e correções.

---

## 4. bmad-dev / bmad-master — Implementação

Mandato:

- Implementar estritamente dentro dos trilhos autorizados.
- Nunca criar novos fluxos paralelos.

Responsabilidades operacionais (por Wave):

1. Wave 4 — ScriptRunner:
   - Implementar apenas em `internal/runtime/runner/**`.
   - Garantir:
     - Entrada: somente `RunAction` vindos do `EventManager`/`on:`,
     - 1 processo por `update_target` (cancelamento anterior),
     - Nenhum `os.Exit` no core,
     - Saída sempre como `ShantillyEvent` + updates declarativos.
   - Referenciais:
     - [`docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md`](docs/stories/4.1.scriptrunner-execucao-declarativa-e14.story.md:1),
     - [`docs/architecture/components.md#5-scriptrunner--e14`](docs/architecture/components.md:256),
     - Gates associados.

2. Wave 5 — Modal Stack:
   - Implementar apenas em `internal/runtime/modal/**`.
   - Integrar:
     - `LayoutManager` para desenhar modais sobre o layout,
     - `EventManager` para abrir/fechar e rotear eventos.
   - Obedecer:
     - Pilha explícita,
     - Foco só no topo,
     - Modais apenas via `ModalRequest`.
   - Referenciais:
     - [`docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md`](docs/stories/5.1.modal-stack-e-seguranca-jit-e15.story.md:1).

3. Wave 6 — Legado:
   - Encapsular v1.x em um `FormComponent` único:
     - Sem refs novas a `internal/tui` ou `FormConfig` fora do sandbox.
   - Qualquer desvio é bug e será pego pelos gates.

4. Wave 7 — Segurança:
   - Implementar:
     - Validador de `AppConfig` com deny-by-default,
     - Whitelist de tipos/ações/paths,
     - `SecurityPolicy` consumida pelo ScriptRunner.
   - Blocagem obrigatória de YAML malicioso.

Entrada esperada:

- Stories com IDs E1.x,
- Contratos de arquitetura estáveis,
- Feedback dos gates QA.

Saída esperada:

- Código apenas nos diretórios permitidos,
- Testes alinhados à matriz QA,
- Zero violações de governança.

---

## 5. bmad-pm / bmad-po — Produto

Mandato:

- Garantir que PRDs e backlog respeitem as constraints imutáveis.
- Nunca pedir features que violem o pipeline ou destravem o legado.

Responsabilidades:

- Redigir/atualizar requisitos apenas:
  - Em termos de `AppConfig` único,
  - Usando blocos E1.1–E1.5,
  - Mapeando Waves 4–7 para impactos claros.
- Validar se cada nova demanda:
  - Tem encaixe nos engines normativos,
  - Não exige exceções (caso exija, deve ser rejeitada ou reprojetada).

Entrada esperada:

- Governança e stories consolidadas.

Saída esperada:

- PRD consistente com contratos,
- Backlog pronto para Architect/Dev/QA sem ambiguidade.

---

## 6. Sequência Operacional do Próximo Ciclo (Wave 4→7) — Trilho Executável BMAD

Fluxo único, repetível e bloqueante. Cada passo depende do anterior; qualquer salto é violação de governança.

1. Planejamento e enquadramento (bmad-pm / bmad-po)
   - Atividades:
     - Selecionar objetivos da Wave (4–7) a partir dos PRDs.
     - Garantir aderência explícita a:
       - [`AGENTS.md`](AGENTS.md:1),
       - [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1),
       - [`docs/architecture/governance-runtime-tui-v2.0.md`](docs/architecture/governance-runtime-tui-v2.0.md:1).
     - Mapear demandas para blocos E1.x e stories em `docs/stories/**`.
   - Saídas obrigatórias:
     - Lista de histórias priorizadas com IDs E1.x.
     - Nenhuma história que exija rota paralela, novo binário ou quebra de contrato.

2. Modelagem e contratos (bmad-architect)
   - Pré-condição:
     - Histórias E1.x aprovadas pelo PM/PO.
   - Atividades:
     - Atualizar/alinhar:
       - Modelos e workflows em `docs/architecture/**`,
       - Story/feature com seus contratos (ScriptRunner, Modal Stack, FormComponent sandbox, SecurityPolicy).
   - Saídas obrigatórias:
     - Seções de arquitetura atualizadas e rastreáveis,
     - Indicação clara dos diretórios-alvo (`internal/runtime/runner`, `internal/runtime/modal`, sandbox FormComponent, validação AppConfig).

3. Desenho de qualidade e gates (bmad-qa)
   - Pré-condição:
     - Contratos de arquitetura estáveis para a Wave.
   - Atividades:
     - Atualizar [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1).
     - Criar/ajustar gates em `docs/qa/gates/*.yml` cobrindo:
       - Event Engine, LayoutManager, ScriptRunner, Modal Stack,
       - Encapsulamento legado, no-osexit-core, security-jit-anti-trojan.
   - Saídas obrigatórias:
     - Gates definidos como bloqueantes na CI,
     - Critérios objetivos que os devs devem satisfazer.

4. Implementação estrita (bmad-dev / bmad-master)
   - Pré-condição:
     - Contratos + gates disponíveis.
   - Atividades:
     - Implementar APENAS nos diretórios autorizados e via pipeline normativo.
     - Referenciar stories E1.x nos PRs.
   - Saídas obrigatórias:
     - Código aderente aos contratos,
     - Testes alinhados à matriz QA,
     - Nenhum padrão proibido introduzido.

5. Validação e enforcement (bmad-qa)
   - Pré-condição:
     - PRs submetidos com referências corretas.
   - Atividades:
     - Executar todos os gates relevantes.
     - Validar cenários críticos definidos nas stories e matriz.
   - Saídas obrigatórias:
     - Aprovação somente se todos os gates PASS,
     - Registro das violações como causas explícitas de falha.

6. Auditoria executiva e coesão (bmad-orchestrator)
   - Pré-condição:
     - PR(s) aprovados por QA e aderentes aos contratos de arquitetura.
   - Atividades:
     - Verificar:
       - Story → Arquitetura → Código → QA/Gates → Governança.
       - Ausência de rotas paralelas, bypass de segurança ou uso incorreto do legado.
   - Saída obrigatória:
     - Sinal verde para merge/release somente se o trilho completo for respeitado.
     - Caso contrário, bloqueio com apontamento normativo.

Este fluxo é vinculante. Qualquer PR, branch ou iniciativa que tente:

- pular etapas,
- contornar gates,
- escrever fora dos diretórios autorizados,
- ou burlar o pipeline declarativo,
deve ser tratado como violação crítica e bloqueado imediatamente.

---

## 7. Rotina de Auditoria de PRs/Releases (Checklist Objetivo — Encadeada aos Gates)

Para qualquer agente (especialmente Orchestrator/QA/Architect) revisar PR/release:

1. Contrato:
   - Toca apenas:
     - `internal/runtime/runner/**`, `internal/runtime/modal/**`, sandbox FormComponent ou arquivos de docs/qa alinhados?
2. Pipeline:
   - Preserva `ShantillyEvent` → `EventManager` → `on:` → ações?
3. Segurança:
   - Nenhum `os.Exit` fora de [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1)?
   - Nenhum exec direto fora do ScriptRunner?
4. Legado:
   - Nenhum novo uso de `internal/tui` ou `FormConfig` fora do sandbox?
5. Modais:
   - Nenhum modal aberto fora de `ModalRequest` + Modal Stack?
6. AppConfig:
   - Nada que amplie o modelo fora do que está em arquitetura/security/data-models sem atualização coordenada?

Se qualquer resposta for “não”, o resultado é bloqueio imediato, com apontamento ao trecho violador.

---

Este playbook fixa, em forma operacional, como cada agente deve agir para manter a evolução íntegra, disciplinada e alinhada ao BMAD no ciclo contínuo das Waves 4–7.
