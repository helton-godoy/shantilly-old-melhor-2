# Security — Runtime TUI Declarativo v2.0 (E1.5)

Este documento define os requisitos de segurança NORMATIVOS do Runtime TUI Declarativo v2.0, substituindo a visão v1.x centrada em `FormConfig`/`shantilly form`.

Fontes vinculantes:

- PRD pivot: [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:210)
- Arquitetura pivot: [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1), [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1)
- Modelos: [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)
- Componentes: [`docs/architecture/components.md`](docs/architecture/components.md:1)
- Workflows: [`docs/architecture/core-workflows.md`](docs/architecture/core-workflows.md:1)
- Story Map: [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
- QA Matrix: [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)
- Governança: [`docs/architecture/governance-runtime-tui-v2.0.md`](docs/architecture/governance-runtime-tui-v2.0.md:1)

Escopo direto: Bloco E1.5 (Modal Stack + Segurança JIT + Anti-Trojan YAML) e constraints transversais aplicáveis a E1.1–E1.4.

Qualquer menção a `ConfigParser`/`FormConfig`/`shantilly form` como baseline de segurança ativa é LEGACY (não vinculante) e só é aceitável como referência histórica para o `FormComponent`.

---

## 1. Superfície de Entrada e Princípios Globais

Superfície de entrada do runtime v2.0:

- YAML único de aplicação:
  - Fornecido via `stdin` ou `--file`.
  - Estruturado como `AppConfig`.
- Eventos de UI:
  - Emissão por `ShantillyComponent`s (list, viewport, form, buttongroup).
- Ações:
  - Bloco `on:` (OnHandler) mapeando `ShantillyEvent` → `RunAction`/`ModalRequest`/state.

Princípios normativos:

- Contrato único:
  - Toda semântica configurável deve passar por `AppConfig` (não há arquivos/parâmetros paralelos redefinindo layout/componentes/on/run).
- `deny by default`:
  - O que não for reconhecido pelo contrato é rejeitado ou tratado como erro explícito.
- Nenhuma execução implícita:
  - `run:` só pode ser disparado via:
    - Eventos de UI;
    - Regras `on:` declaradas.
  - Proibidos auto-runs em tempo de load.
- Segurança JIT:
  - Ações sensíveis exigem confirmação (`confirm`) e/ou coleta segura de segredos (`prompt_secrets`) via Modal Stack.
- Sem `os.Exit` no core:
  - Encerramentos são controlados via mensagens/eventos internos.

---

## 2. Validação de YAML e Anti-Trojan YAML (E1.5 — Wave 7)

Responsável: parser/validator declarativo (ex.: `pkg/declarative`), alinhado ao contrato de dados em [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1).

Mandatos normativos (deny-by-default + whitelist):

1. Schema Estrito (AppConfig como contrato único)
   - Validar, de forma bloqueante:
     - Tipos corretos para todos os campos de:
       - `AppConfig`, `LayoutNode`, `Component`, `OnHandler`, `RunAction`, `ModalRequest`, `SecurityPolicy`.
     - `LayoutNode.Type` ∈ {`column`, `row`, `box`}.
     - `Component.Type` ∈ {`list`, `viewport`, `form`, `buttongroup`} no escopo do Epic 1.
     - `OnHandler`:
       - Deve possuir `event` válido (compatível com `ShantillyEvent`),
       - E pelo menos uma ação reconhecida (`run`, `open_modal`, `update_state`).
   - Campos desconhecidos ou estruturas fora do modelo:
     - Em blocos críticos (`layout`, `components`, `on`, `run`, `security`, `modal`):
       - DEVEM causar erro de validação explícito.
     - É proibido depender de chaves “toleradas” para semântica relevante.

2. Deny by Default efetivo
   - Padrão normativo:
     - Tudo que não é explicitamente modelado é tratado como inválido.
   - YAML com:
     - Chaves inesperadas em `run:`, `security:`, `on:`, `components:`, `layout:`,
     - Tipos de componentes não suportados,
     - Campos extras em `SecurityPolicy` não reconhecidos,
   - DEVE ser rejeitado com mensagem clara (sem aceitação silenciosa).
   - Ausência de `SecurityPolicy`:
     - NÃO desabilita o deny-by-default; proteções mínimas permanecem ativas.

3. Anti-Trojan YAML (E1.5-6)
   - Proibições absolutas:
     - Executar scripts durante parse/carga (sem evento de UI).
     - Esconder execuções em campos “mágicos” ou não documentados.
     - Habilitar auto-run implícito com base apenas na presença de nós YAML.
   - Execução de qualquer `run:`:
     - Sempre vinculada a:
       - `ShantillyEvent` → `OnHandler` → `RunAction` → `ScriptRunner`.
   - Whitelist:
     - `SecurityPolicy` e/ou implementação devem prover lista de scripts/paths/ações permitidas.
     - `ScriptRunner` deve recusar execuções fora dessa política.

Rastreabilidade:

- Blocos:
  - E1.5 — Segurança JIT + Anti-Trojan YAML.
  - E1.5-6 — Whitelist explícita + integração com ScriptRunner.
- Backlog:
  - [`docs/architecture/backlog-waves4-7-epic1.md`](docs/architecture/backlog-waves4-7-epic1.md:304).
- Data models:
  - [`SecurityPolicy`](docs/architecture/data-models.md:240),
  - [`RunAction`](docs/architecture/data-models.md:175),
  - [`ModalRequest`](docs/architecture/data-models.md:208).
- QA Matrix:
  - Linhas de validação de schema, deny-by-default e Trojan YAML.
- QA Gates:
  - `docs/qa/gates/1.x.security-jit-anti-trojan.yml`.

---

## 3. Regras para `run:` e ScriptRunner (E1.4 + E1.5)

As ações `run:` são o principal ponto de risco.

Requisitos normativos:

1. Pipeline Exclusivo
   - `RunAction` só é disparado por:
     - `EventManager` ao resolver `OnHandler`.
   - Não pode haver:
     - Execução direta de scripts por componentes ou por código ad hoc.

2. Escopo e Whitelist
   - `RunAction.Script`:
     - Deve respeitar política de caminhos (whitelist/prefixos permitidos) definida em implementações e/ou `SecurityPolicy`.
   - Recomendações (p/ implementações futuras, sem conflitar com pivot):
     - Bloquear paths obviamente perigosos conforme contexto de uso.

3. Bindings Seguros
   - `args`/`stdin`:
     - Populados via estado de componentes e AppConfig.
     - Devem evitar concatenação de entradas não validadas de forma ingênua.
   - Segredos (quando existirem):
     - Nunca mostrados em logs ou erros.

4. Integração com SecurityPolicy
   - `SecurityPolicy`:
     - Pode definir `AllowedScripts`, `DenyUnknown` e regras adicionais.
   - Ausência de `SecurityPolicy`:
     - NÃO desabilita proteções default: `deny by default` permanece ativo.

Rastreabilidade:

- E1.4/E1.5.
- QA Matrix:
  - Linha de ScriptRunner + whitelists + segurança.

---

## 4. Modal Stack e Segurança JIT (E1.5)

A Modal Stack é parte obrigatória da segurança.

Requisitos:

- `RunAction.Confirm: true`:
  - DEVE abrir modal de confirmação.
  - Execução do script só ocorre após confirmação explícita.
- `RunAction.PromptSecrets`:
  - DEVE abrir modal(es) para coleta de segredos.

Propriedades de segurança:

- Segredos:
  - Mantidos apenas em memória (escopo de execução).
  - Não armazenados em disco.
  - Não logados.
- Controlo de foco:
  - Enquanto houver modal na pilha:
    - Apenas o topo recebe eventos.
  - Interação com plano de fundo bloqueada.
- Integração:
  - Modal aberto por `EventManager` a partir de `OnHandler`/`RunAction`.
  - Respostas geram `ShantillyEvent` apropriado.
  - `EventManager` monta então o `RunAction` final com dados/segredos.

Rastreabilidade:

- E1.5 — Modal Stack + Security.
- QA Matrix:
  - Casos de confirm/prompt_secrets + gating de ações sensíveis.

---

## 5. Regra de 1 Processo por `update_target` como Invariante de Segurança (E1.4)

Esta regra é tanto funcional quanto de segurança.

Requisitos:

- Para cada `update_target`:
  - No máximo 1 processo ativo.
- Novo `RunAction` com mesmo `update_target`:
  - DEVE:
    - Enviar SIGTERM ao processo anterior.
    - Aguardar timeout configurado.
    - Opcionalmente aplicar SIGKILL se necessário.
- Comportamento:
  - Evita DoS local (multiplicação de processos).
  - Garante previsibilidade de quem controla determinada área da UI.

Rastreabilidade:

- E1.4.
- QA Matrix:
  - Linhas específicas para 1 processo/update_target.

---

## 6. Proibição de `os.Exit` no Core (Cross-cutting)

Regra mandatória:

- Em `internal/runtime/**`:
  - É proibido chamar `os.Exit`.
- `os.Exit` só é permitido:
  - Em `cmd/shantilly` (camada CLI), após avaliação dos resultados.
- Implicações:
  - YAML malicioso ou erro interno não pode encerrar o processo de forma abrupta sem passar pelo fluxo de erros controlado.
  - Erros são propagados como:
    - Valores de retorno,
    - `tea.Msg`/eventos internos,
    - Estruturas de erro tipadas.

Rastreabilidade:

- E1.4/E1.5 (Restrições estratégicas).
- QA Matrix:
  - Gate específico `no-osexit-core`.

---

## 7. Dependências, Execução Local e Contexto

Para Epic 1:

- Escopo:
  - Execução local em terminal.
  - Sem servidor remoto embutido.
- Boas práticas de dependência:
  - `govulncheck` + `golangci-lint` como linha base.
  - Revisão manual de novas libs (supply chain).
- Não há:
  - AuthN/AuthZ centralizada no runtime.
- Operador:
  - Assume responsabilidade pelas permissões do usuário que executa `shantilly`.

Rastreabilidade:

- NFRs de segurança gerais.
- QA:
  - Checks automáticos de dependências.

---

## 8. Security Testing & QA (Ligação com Matriz)

Cada cláusula deste documento deve ter linha correspondente na QA Matrix.

Diretrizes:

- Cobrir com testes:
  - Validação de YAML (válido/inválido, chaves extras, tipos errados).
  - Cenários de “Trojan YAML”:
    - Ações escondidas,
    - Auto-run,
    - Component types inválidos.
  - Modal Stack:
    - Confirmações obrigatórias,
    - Prompt de segredos.
  - Regra 1 processo/update_target.
  - Proibição de `os.Exit`:
    - Scans estáticos.
- Status na QA Matrix:
  - Só pode avançar para `implemented`/`passed` quando alinhado a este doc.

---

## 9. Legacy (não vinculante)

Os seguintes pontos são classificados como LEGADO:

- Modelo de segurança baseado apenas em:
  - `FormConfig`/`FieldConfig`,
  - `ConfigParser` v1.x,
  - Fluxo `shantilly form` → JSON.
- Qualquer recomendação de segurança que:
  - Trate este fluxo como principal.

Uso permitido:

- Apenas como:
  - Base técnica interna para `FormComponent`.
  - Referência histórica.
- Não pode:
  - Guiar decisões para o runtime declarativo v2.0.
  - Ser tratado como contrato ativo.

Qualquer documento ou código que ainda reflita esse modelo como vigente deve ser alinhado a este documento ou marcado explicitamente como Legacy.
