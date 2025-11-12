# shantilly v2.0 — Runtime TUI Declarativo

# Gate de Governança Obrigatório para TODOS os PRs

Este checklist é BLOQUEANTE.  
Se qualquer item aplicável estiver em não conformidade, o PR NÃO DEVE ser aprovado.

Referências normativas:

- [`AGENTS.md`](AGENTS.md:20)
- [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1)
- [`docs/architecture/governance-runtime-tui-v2.0.md`](docs/architecture/governance-runtime-tui-v2.0.md:1)
- [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
- [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)
- [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)

Preencha todas as seções relevantes antes de pedir review.

## 1. Identificação e Escopo

- [ ] Este PR está explicitamente vinculado a pelo menos uma Story E1.x ou bloco Epic 1:
  - IDs / arquivos:
- [ ] As mudanças estão dentro do escopo do Runtime TUI Declarativo v2.0 (binário único + YAML único).
- [ ] Se tocar em legado (`FormConfig`, `internal/tui`, `shantilly form`), o PR declara explicitamente se é:
  - [ ] Encapsulamento no `FormComponent` (permitido)
  - [ ] Ajuste de Legacy (marcado como `legacy_refactored`)
  - [ ] N/A (não toca em legado)

## 2. Stories (*.story.md)

Bloqueante:

- [ ] Nenhum novo arquivo `docs/stories/*.story.md` sem o cabeçalho padrão:

  Deve conter, exatamente:
  - `Epic: 1 - Runtime TUI Declarativo — Fundação do Runtime`
  - `Bloco: E1.x - <nome do bloco>`
  - `Trace: PRD (...), Arquitetura (...), QA (matrix-epic-1-runtime-tui-coverage.md#...)`
  - `Status: draft|approved|legacy_refactored|new`

- [ ] Qualquer story que trate comportamento v1.x já absorvido está marcada como `Status: legacy_refactored` e explica que é comportamento herdado/FormComponent, não roadmap ativo.
- [ ] Stories seeds Waves 3–7 (3.1, 3.3, 4.1, 5.1, 6.1, 7.1) usadas neste PR possuem:
  - [ ] Cabeçalho padrão
  - [ ] Rastreabilidade completa para PRD, Arquitetura e QA

## 3. PRD / Roadmap

Bloqueante:

- [ ] Nenhuma modificação neste PR reintroduz qualquer PRD v1.x como roadmap ativo.
- [ ] Qualquer referência a épicos v1.x está confinada a seções marcadas como Legacy/Referência.
- [ ] `docs/prd/epic-list.md` permanece como fonte ÚNICA do Roadmap v2.0.
- [ ] `docs/prd/epic-1-runtime-tui-foundation.md` continua descrito como substituto da visão centrada em "form".

## 4. Arquitetura

Bloqueante:

- [ ] Nenhum documento de arquitetura deste PR apresenta v1.x como estado atual ou recomendação.
- [ ] Qualquer menção a `FormConfig`/`internal/tui`/“shantilly form”:
  - [ ] Está dentro de seção explicitamente marcada como `Legacy (não vinculante para v2.0)`, ou
  - [ ] Está claramente encapsulada no contexto de `FormComponent` legado.
- [ ] Nenhum novo fluxo paralelo é introduzido fora do pipeline normativo:

  Pipeline único obrigatório:
  - `AppConfig` YAML único
  - `ShantillyEvent`
  - `EventManager`
  - `on:`
  - `RunAction` / `ModalRequest` / `UpdateState`
  - `ScriptRunner` / Modal Stack / atualização declarativa

## 5. QA Matrix e Gates (Rastreabilidade obrigatória BMAD)

Bloqueante:

1) Matriz — [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)

- [ ] Liste abaixo as linhas/IDs E1.x impactadas por este PR:
  - ID/linha:
  - ID/linha:
  - ID/linha:
- [ ] Confirme:
  - [ ] Cada mudança relevante de runtime/config/segurança está mapeada na matriz.
  - [ ] Status ajustado para `planned` / `implemented` / `passed` conforme o estágio real (não marcar `implemented` sem código + testes).

2) Gates QA — `docs/qa/gates/*.yml`

- [ ] Para cada comportamento afetado, referencie os gates correspondentes:
  - [ ] [`docs/qa/gates/1.x.layout-manager.yml`](docs/qa/gates/1.x.layout-manager.yml:1)
  - [ ] [`docs/qa/gates/1.x.event-engine.yml`](docs/qa/gates/1.x.event-engine.yml:1)
  - [ ] [`docs/qa/gates/1.x.scriptrunner-and-update-target.yml`](docs/qa/gates/1.x.scriptrunner-and-update-target.yml:1)
  - [ ] [`docs/qa/gates/1.x.modal-stack.yml`](docs/qa/gates/1.x.modal-stack.yml:1)
  - [ ] [`docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml`](docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml:1)
  - [ ] [`docs/qa/gates/1.x.no-osexit-core.yml`](docs/qa/gates/1.x.no-osexit-core.yml:1)
  - [ ] [`docs/qa/gates/1.x.security-jit-anti-trojan.yml`](docs/qa/gates/1.x.security-jit-anti-trojan.yml:1)
- Para cada gate marcado:
  - [ ] Estado esperado após este PR: `PASS` | `AFFECTED` | `REVIEW` | `WAIVER-PROPOSED`
  - [ ] Evidência associada (tests, arquivos, seções de docs, logs de CI):

3) Workflow normativo — `Governança Waves 4-7 (Gates 1.x + Go checks)`

- [ ] Confirme que este PR passa com sucesso pelo job obrigatório `.github/workflows/governanca-waves4-7.yml`.
- [ ] Nenhum uso de outros workflows como bypass dos gates normativos.

## 6. Código — Regras Estruturais Críticas

Bloqueante:

- [ ] Não há `os.Exit` fora de [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1).
- [ ] Não há execução direta de scripts (`exec.Command`, etc.) fora de `internal/runtime/runner/**`.
- [ ] Nenhum componente, `LayoutManager` ou código legado:
  - [ ] Executa scripts diretamente.
  - [ ] Abre modais diretamente.
  - [ ] Altera estado global fora de `UpdateState`/pipeline declarativo.
- [ ] Nenhum uso novo de `FormConfig`/`internal/tui` fora do sandbox do futuro `FormComponent` (qualquer exceção = bug arquitetural).

## 7. Confirmação por Agente Responsável

Marque os agentes que revisaram este PR (pode ser via persona ou humano equivalente):

- [ ] bmad-analyst — Stories e rastreabilidade E1.x
- [ ] bmad-architect — Aderência à arquitetura v2.0 e contratos formais
- [ ] bmad-dev / bmad-master — Implementação aderente ao pipeline único
- [ ] bmad-qa — QA Matrix, gates e ausência de violações críticas
- [ ] bmad-orchestrator / bmad-sm — Governaça geral e bloqueios

Declaração (obrigatória antes de aprovar):

> Confirmo que este PR está em conformidade com o Runtime TUI Declarativo v2.0, não reintroduz v1.x como roadmap ativo, respeita o AppConfig/YAML único, o pipeline único de eventos e efeitos, e mantém o legado estritamente confinado conforme documentos normativos.
