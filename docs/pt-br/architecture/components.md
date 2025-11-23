# Components

Este documento define os componentes normativos do Runtime TUI Declarativo v2.0, alinhados ao contrato YAML único e aos blocos:

- E1.1 — LayoutManager
- E1.2 — Component Model + FormComponent legado encapsulado
- E1.3 — EventManager + ShantillyEvent + on:
- E1.4 — ScriptRunner
- E1.5 — Modal Stack + Security

Fontes normativas:

- PRD: [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura pivot: [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1), [`docs/architecture/high-level-architecture.md`](docs/architecture/high-level-architecture.md:1), [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)
- Story Map: [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
- Implementation Plan: [`docs/architecture/implementation-plan-epic-1-runtime-tui.md`](docs/architecture/implementation-plan-epic-1-runtime-tui.md:1)
- QA Matrix: [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)

Qualquer menção ao design v1.x ou a `shantilly form` como arquitetura ativa é classificada como Legacy (não vinculante) e só é permitida como detalhe de implementação do `FormComponent` dentro do runtime v2.0.

---

## 1. ShantillyComponent (Contrato Base) — E1.2

Todos os componentes TUI concretos do runtime v2.0 DEVEM implementar uma interface comum (forma ilustrativa; implementação exata é decisão de bmad-dev/bmad-master, desde que preserve o contrato):

```go
// internal/components/component.go

type ShantillyComponent interface {
    ID() string
    Init() tea.Cmd
    Update(msg tea.Msg) (ShantillyComponent, tea.Cmd)
    View() string
    SetBounds(width, height int)
}
```

Regras normativas:

- O `LayoutManager` orquestra layout e foco invocando apenas métodos de `ShantillyComponent`.
- Componentes:
  - Não executam scripts diretamente.
  - Emitem eventos normalizados (`ShantillyEvent`) consumidos pelo `EventManager`.
- Nenhum componente pode:
  - Invocar `os.Exit`.
  - Criar rotas de controle fora do pipeline declarativo (`on:` + `RunAction`).

Rastreabilidade:

- E1.2 — Component Model.
- Story Map: E1.2.
- QA Matrix: blocos E1.2 (incluindo encapsulamento do legado).

---

## 2. Componentes Oficiais do Epic 1 — E1.2

Conjunto normativo de componentes para o Epic 1:

- `list`
- `viewport`
- `form`
- `buttongroup`

Cada entrada `Component` definida no AppConfig (ver [`docs/architecture/data-models.md`](docs/architecture/data-models.md:1)) é instanciada como um `ShantillyComponent` correspondente.

### 2.1 ListComponent (`type: list`)

Responsabilidades:

- Renderizar listas/menus navegáveis via teclado.
- Expor seleção atual e estados simples.

Contrato essencial (conceitual):

- Props mínimas:
  - `items: []{ id, label }`
  - `initial_index` (opcional)
- Eventos (via `ShantillyEvent`):
  - `list_id:select` — item selecionado.
- Uso típico:
  - Navegação multi-painel e seleção de contexto (E1.1/E1.3).

### 2.2 ViewportComponent (`type: viewport`)

Responsabilidades:

- Exibir texto/logs multi-linha.
- Ser destino de `update_target` vindo do `ScriptRunner`.

Contrato essencial:

- Props mínimas:
  - `title`, `border`, `scroll` etc.
- Integração:
  - `ScriptRunner` escreve em `viewport` associado a `UpdateTarget`.
  - Comportamento de append/replace definido de forma previsível.

Eventos:

- Pode emitir eventos de scroll ou interação, conforme evoluções futuras, sempre via `ShantillyEvent`.

### 2.3 ButtonGroupComponent (`type: buttongroup`)

Responsabilidades:

- Exibir um conjunto de ações (botões).
- Mapear ações do operador para eventos semânticos.

Contrato essencial:

- Props mínimas:
  - `buttons: []{ id, label }`
- Eventos:
  - `component_id:click` com `button_id` no payload.

Uso:

- Gate de ações que disparam `on:` e `RunAction` específicos.

### 2.4 FormComponent (`type: form`) — Sandbox Legado (Wave 6, E1.6)

Responsabilidades:

- Fornecer UX de formulário rica usando aprendizado do v1.x.
- Reutilizar internamente `huh`, `bubbletea` e partes de `internal/tui` como implementação LEGACY.

Mandatos de sandbox (vinculantes):

- Implementa `ShantillyComponent`:
  - É apenas mais um componente do catálogo v2.0.
- Configuração:
  - É configurado EXCLUSIVAMENTE via `Component{ Type: "form", Props: ... }` do `AppConfig` v2.0.
  - É proibido expor `FormConfig` como contrato público ou parâmetro direto.
- Encapsulamento do legado:
  - Pode mapear internamente `Props` → estruturas compatíveis com `FormConfig`, mas:
    - Esse mapeamento é DETALHE INTERNO do `FormComponent`.
    - Não recria APIs públicas v1.x.
  - Todo código em `internal/tui/**` e modelos v1.x:
    - Só podem ser utilizados dentro do `FormComponent` (ou módulos internos a ele), nunca direto pelo runtime v2.0.

Eventos e integração:

- Emite apenas `ShantillyEvent`, por exemplo:
  - `form_id:submit` com dados estruturados.
  - `form_id:change` conforme necessidade.
- Eventos são sempre consumidos via:
  - [`EventManager`](internal/runtime/event/manager.go:1) + bloco `on:`.
- ScriptRunner e Modal Stack:
  - Nunca são chamados diretamente pelo `FormComponent`.
  - Sempre via pipeline:
    - `ShantillyEvent` → `EventManager` → `on:` → (`RunAction` | `ModalRequest`).

Proibições normativas (Wave 6 — Encapsulamento do Legado):

- É proibido:
  - Criar novos fluxos baseados diretamente em `FormConfig` fora do `FormComponent`.
  - Importar `internal/tui` diretamente em novos módulos do runtime v2.0.
  - Permitir que o legado:
    - Execute scripts diretamente,
    - Abra modais diretamente,
    - Bypasse `EventManager`/`on:`/ScriptRunner/Modal Stack.
- Qualquer uso de `FormConfig`/`internal/tui` fora do sandbox:
  - É bug arquitetural e deve ser bloqueado por QA/gates.

Rastreabilidade:

- Bloco: E1.2 + Wave 6 — Encapsulamento do Legado.
- Backlog:
  - [`docs/architecture/backlog-waves4-7-epic1.md`](docs/architecture/backlog-waves4-7-epic1.md:1) (itens de sandbox legado).
- Governança:
  - [`docs/architecture/governance-runtime-tui-v2.0.md#5-wave-6--encapsulamento-do-legado-formcomponent`](docs/architecture/governance-runtime-tui-v2.0.md:63).
- QA Gates:
  - `docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml`.

---

## 3. LayoutManager — E1.1

O `LayoutManager` é o `bubbletea.Model` raiz único do Runtime TUI Declarativo v2.0 e o único responsável por layout e foco da UI base.

Responsabilidades (normativas e exaustivas):

- Consumir exclusivamente `LayoutNode` proveniente do `AppConfig` (pacote `pkg/declarative`) como fonte de verdade do layout:
  - É proibido construir layout a partir de estruturas ad-hoc, estados globais, config paralelas ou contratos v1.x.
- Construir a árvore de layout (`column` / `row` / `box`) e calcular de forma determinística os bounds de cada região visível.
- Instanciar e posicionar apenas `ShantillyComponent`s oficiais referenciados por `Component.ID`:
  - Não pode criar componentes "anônimos" ou fora do catálogo declarado.
- Gerir:
  - Foco global entre componentes focáveis (encaminhando `tea.Msg` somente ao componente focado ou conforme regras explícitas do contrato declarativo).
  - Redimensionamento da janela (`tea.WindowSizeMsg`) com recálculo determinístico de layout, evitando flicker e garantindo estabilidade visual.
  - Integração estrutural com a Modal Stack:
    - Renderizar o topo da pilha modal sobre o layout base.
    - Garantir que somente o topo da pilha receba foco/eventos enquanto houver modal ativo.

Invariantes (vinculantes e não negociáveis):

- Exclusão de automação:
  - Não executa `run:` nem qualquer script.
  - Não cria, não agenda e não cancela processos externos.
  - Nenhum comando de terminal ou sistema pode ser executado pelo `LayoutManager`.
- Exclusão de regras `on:`:
  - Não avalia, não roteia e não interpreta regras `on:`.
  - Qualquer lógica de automação/reatividade pertence exclusivamente ao `EventManager` + bloco `on:`.
  - O `LayoutManager` nunca deve conter if/else, switch ou qualquer lógica de decisão além da gestão de layout.
- Isolamento de contratos:
  - Não conhece detalhes internos de componentes além do contrato `ShantillyComponent`.
  - Não lê, não muta e não depende de estado interno de componentes fora dos métodos do contrato.
- Modelos v2.0 apenas:
  - Opera somente sobre modelos v2.0 (`AppConfig`, `LayoutNode`, `Component`).
  - Qualquer referência direta a `FormConfig` ou a fluxos v1.x dentro do `LayoutManager` é classificada como Legacy incorreto e deve ser tratada como bug arquitetural.
- Comunicação com outros motores (PROIBIÇÕES ABSOLUTAS):
  - É proibido:
    - Invocar diretamente o `ScriptRunner` por qualquer meio.
    - Abrir/fechar modais por APIs fora do fluxo `ShantillyEvent` → `EventManager` → `on:` → `ModalRequest`.
    - Alterar políticas de segurança, whitelists ou regras de execução.
    - Criar qualquer rota alternativa de controle ou "atalhos" para automação.
- Segurança e estabilidade:
  - Não pode chamar `os.Exit` (permitido apenas em [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1)).
  - Não pode criar rotas alternativas de controle fora do pipeline declarativo; qualquer automação "escondida" no `LayoutManager` é proibida.

Relações normativas com outros componentes/motores:

- Com `ShantillyComponent`s:
  - Inicializa, posiciona, gerencia foco e encaminha mensagens, respeitando o contrato único.
  - Recebe eventos/msgs dos componentes que serão normalizados e encaminhados ao `EventManager`.
- Com `EventManager`:
  - Nunca resolve `on:` nem dispara `RunAction` diretamente.
  - Serve como fonte de eventos de UI e estado visual que o `EventManager` usa para aplicar o pipeline declarativo.
  - Toda comunicação bidirecional deve ocorrer via `ShantillyEvent` e não por chamadas diretas de API.
- Com `ScriptRunner`:
  - Nenhuma dependência direta; toda execução de script nasce de decisões do `EventManager` com base em `on:`.
  - O `LayoutManager` não deve nem importar pacotes relacionados ao `ScriptRunner`.
- Com Modal Stack:
  - Respeita a pilha como fonte de bloqueio de fundo e foco exclusivo do topo.
  - Não decide sozinho quando abrir/fechar modais; apenas reflete o estado da pilha.
  - Modais são sempre acionados por `OpenModal` via `EventManager`.

Implementação alvo:

- Diretório: `internal/runtime/layout/...`
- Comentários/docblocks DEVEM referenciar explicitamente:
  - Epic 1, Bloco E1.1.
  - Contratos em:
    - [`docs/architecture/data-models.md#2-layoutnode-e11--layoutmanager`](docs/architecture/data-models.md:48)
    - [`docs/architecture/high-level-architecture.md#high-level-overview`](docs/architecture/high-level-architecture.md:22)
- Qualquer implementação que viole as invariantes acima deve ser tratada como erro de conformidade arquitetural e barrada por QA/gates.
- O código do `LayoutManager` deve conter testes que validem todas as invariantes listadas.

Rastreabilidade:

- Story Map: E1.1 (incluindo stories 1.2, 1.3, 1.7 alinhadas ao modelo declarativo).
- QA Matrix: bloco E1.1 — Layout Engine (fonte única via `AppConfig`/`LayoutNode`, ausência de automação/`run:` no `LayoutManager`, integração correta com Modal Stack e EventManager).
- Implementation Plan: Wave 2 — LayoutManager.

## 4. EventManager + ShantillyEvent — E1.3

O `EventManager` é o roteador único de eventos do Runtime TUI Declarativo v2.0.

Responsabilidades (normativas):

- Consumir `ShantillyEvent` como tipo canônico de evento interno:
  - Emitidos por `ShantillyComponent`s (UI),
  - Pelo `ScriptRunner` (lifecycle de run:),
  - Pela Modal Stack (confirmações/prompt),
  - Por eventos internos estritamente necessários do runtime.
- Resolver o bloco `on:` (`[]OnHandler` em `AppConfig.On`) como tabela declarativa única:
  - Localizar `OnHandler` cujo `Event` corresponda ao `ShantillyEvent` recebido (`componentID:type`),
  - Aplicar (quando implementado) a expressão `When` de forma determinística,
  - Disparar ações declarativas:
    - `RunAction` → encaminhada exclusivamente para o `ScriptRunner`,
    - `OpenModal` → encaminhada exclusivamente para a Modal Stack,
    - `UpdateState` → aplicada ao estado declarativo do runtime (layout/bindings), nunca como mutação ad-hoc.

Invariantes (vinculantes):

- Pipeline obrigatório:
  - `ShantillyComponent`/motores → `ShantillyEvent` → `EventManager` → `on:` → (`RunAction` | `ModalRequest` | `UpdateState`) → motores dedicados.
- Handlers “soltos”:
  - Qualquer lógica reativa fora de `EventManager` + `on:` (ex.: if/else em componentes, chamadas diretas ao `ScriptRunner`) é proibida.
- Execução de scripts:
  - `EventManager` NÃO executa scripts diretamente.
  - TODA execução passa por `ScriptRunner` usando contratos declarativos (RunAction/RunRequest).
- Integração com LayoutManager:
  - `LayoutManager` nunca resolve `on:` nem `RunAction`.
  - Qualquer interação entre UI e automação passa por `EventManager` + `on:`.
- Integração com Modal Stack:
  - Modais são sempre originados de `OpenModal` em `OnHandler`, disparados pelo `EventManager`.
  - Eventos de modal (`modal.confirmed`, etc.) voltam como `ShantillyEvent` para novo ciclo de resolução via `on:`.
- Modelos v2.0:
  - `EventManager` opera apenas sobre `ShantillyEvent`, `OnHandler`, `RunAction`, `ModalRequest` definidos em `pkg/declarative`.
  - Qualquer referência a contratos v1.x é Legacy e não pode coexistir como rota paralela.

Implementação alvo:

- Diretório: `internal/runtime/event/...`
- Papel sugerido:
  - Serviço puro (não necessariamente um `tea.Model`), integrando-se ao loop TEA via `tea.Msg`/`tea.Cmd`.
  - Funções/métodos para:
    - Registrar `OnHandler`s a partir do `AppConfig`,
    - Receber `ShantillyEvent`,
    - Retornar ações declarativas (para ScriptRunner, Modal Stack, atualizações de estado).

Rastreabilidade:

- E1.3 — Event Engine.
- Story Map: E1.3 (incluindo `docs/stories/3.3.event-engine-and-on-routing.story.md`).
- QA Matrix: bloco E1.3 — Event Engine (normalização + roteamento via on:).
- Implementation Plan: Wave 3 — EventManager + ShantillyEvent + on:.

---

## 5. ScriptRunner — E1.4 (Wave 4)

O `ScriptRunner` é o EXECUTOR ÚNICO e dedicado de `RunAction` (bloco `run:` do AppConfig) no Runtime TUI Declarativo v2.0.

Mandatos estruturais (vinculantes):

- Local de implementação:
  - Exclusivamente em `internal/runtime/runner/**`.
  - É proibida qualquer lógica de execução de scripts, criação de processos ou wrappers equivalentes fora deste diretório.
- Fonte de entrada:
  - Consome SOMENTE instâncias de [`RunAction`](docs/architecture/data-models.md:175) produzidas pelo pipeline normativo:
    - `ShantillyEvent` → [`EventManager`](internal/runtime/event/manager.go:1) → `on:` (`OnHandler`) → `RunAction`.
  - É proibido construir `RunAction` “na mão” fora desse fluxo para executar scripts.
- Saída:
  - Sempre como:
    - `ShantillyEvent` estruturado (ex.: `script.complete`, `script.error`, `script.cancelled`),
    - e/ou atualização declarativa de `update_target` (ex.: escrever em `ViewportComponent`).
  - Nunca chama diretamente componentes ou o `LayoutManager`; interação é sempre mediada pelo pipeline de eventos.

Responsabilidades normativas:

- Executar `RunAction.Script` com suporte a:
  - `Args` (templates declarativos),
  - `Stdin` (derivado de estado declarativo),
  - `Env` (quando previsto no modelo),
  - `UpdateTarget` (id lógico de destino, ex.: viewport).
- Garantir:
  - Mapeamento determinístico entre `RunAction` e processo iniciado.
  - Emissão de eventos de lifecycle para o [`EventManager`](internal/runtime/event/manager.go:1).

Invariantes obrigatórios (sem exceções):

1. Pipeline exclusivo (E1.4-1)
   - TODA execução de script deve seguir:
     - `ShantillyEvent` → `EventManager` → `on:` → `RunAction` → `ScriptRunner` → `ShantillyEvent` de volta.
   - Proibições:
     - Nenhum componente, [`LayoutManager`](internal/runtime/layout/manager.go:1) ou código legado pode executar scripts diretamente.
     - Nenhuma função utilitária paralela (ex.: `exec.Run` em outros pacotes) é permitida para acionar processos de negócio.

2. Regra de 1 processo por `update_target` (E1.4-3)
   - Para cada `RunAction.UpdateTarget`:
     - No máximo UM processo ativo.
   - Novo `RunAction` para o mesmo `update_target` DEVE:
     - Enviar SIGTERM ao processo anterior.
     - Após timeout configurado, aplicar SIGKILL se necessário.
   - Objetivos:
     - Evitar DoS local/disputa por viewport.
     - Garantir previsibilidade de quem controla a saída visível.
   - Vinculado a:
     - [`Workflow 4`](docs/architecture/core-workflows.md:181),
     - Segurança (regra tratada como invariante em [`docs/architecture/security.md`](docs/architecture/security.md:167)),
     - Gate QA `1.x.scriptrunner-and-update-target.yml`.

3. Proibição absoluta de `os.Exit` no core (E1.4-4)
   - `ScriptRunner`:
     - NÃO pode chamar `os.Exit`.
   - `os.Exit` é permitido apenas em [`cmd/shantilly/main.go`](cmd/shantilly/main.go:1).
   - Erros devem:
     - Ser propagados como valores de retorno tipados ou `ShantillyEvent` de erro,
     - Nunca encerrar o binário de forma abrupta.

4. Proibição de execução direta fora de `internal/runtime/runner/**` (E1.4-5)
   - Qualquer uso de APIs de processo (`os/exec`, `syscall`, etc.) fora de `internal/runtime/runner/**` que tenha efeito similar ao ScriptRunner:
     - É violação arquitetural.
   - Cobertura:
     - Gate `1.x.scriptrunner-and-update-target.yml`,
     - Gate `1.x.no-osexit-core.yml`.

5. Integração com `SecurityPolicy` + Anti-Trojan YAML (E1.5-6)
   - O `ScriptRunner` DEVE respeitar:
     - Whitelists, políticas e flags definidos em [`SecurityPolicy`](docs/architecture/data-models.md:240) / [`docs/architecture/security.md`](docs/architecture/security.md:52).
   - Execuções que violem a política:
     - Não devem ocorrer; DEVEM resultar em erro bloqueante + evento estruturado.

Rastreabilidade:

- Bloco: E1.4 — ScriptRunner.
- Backlog: [`docs/architecture/backlog-waves4-7-epic1.md`](docs/architecture/backlog-waves4-7-epic1.md:1) (itens E1.4-1, E1.4-2, E1.4-3, E1.4-R2 como REJECTED).
- Workflows:
  - [`Workflow 1 — Execução Declarativa`](docs/architecture/core-workflows.md:23),
  - [`Workflow 4 — 1 processo/update_target`](docs/architecture/core-workflows.md:181).
- Data models:
  - [`RunAction`](docs/architecture/data-models.md:175).
- Segurança:
  - [`docs/architecture/security.md#3-regras-para-run-e-scriptrunner-e14--e15`](docs/architecture/security.md:93),
  - [`docs/architecture/security.md#5-regra-de-1-processo-por-update_target-como-invariante-de-seguranca-e14`](docs/architecture/security.md:167),
  - [`docs/architecture/security.md#6-proibicao-de-osexit-no-core-cross-cutting`](docs/architecture/security.md:192).
- Governança:
  - [`docs/architecture/governance-runtime-tui-v2.0.md#3-wave-4--scriptrunner-e14`](docs/architecture/governance-runtime-tui-v2.0.md:31).
- QA Gates:
  - `docs/qa/gates/1.x.scriptrunner-and-update-target.yml`,
  - `docs/qa/gates/1.x.no-osexit-core.yml`,
  - `docs/qa/gates/1.x.security-jit-anti-trojan.yml`.

---

## 6. Modal Stack — E1.5 (Wave 5, parte 1)

A Modal Stack é a ÚNICA infraestrutura de modais do Runtime TUI Declarativo v2.0 e parte central da Segurança JIT.

Mandatos estruturais (vinculantes):

- Local de implementação:
  - Exclusivamente em `internal/runtime/modal/**`.
  - Qualquer lógica de modais fora deste diretório é considerada rota paralela proibida.
- Forma:
  - Implementada como pilha explícita:
    - `push` — empilhar novo modal.
    - `pop` — remover modal do topo.
    - `top` — ler modal ativo.
- Integrações obrigatórias:
  - `LayoutManager`:
    - Apenas desenha o modal do topo sobre o layout base.
    - Não decide sozinho abrir/fechar modais; apenas reflete o estado da pilha.
  - `EventManager`:
    - Único responsável por acionar `ModalRequest` com base em `OnHandler`.
    - Recebe eventos de saída dos modais como `ShantillyEvent`.

Pipeline declarativo único (E1.5-1, E1.5-4):

- Abertura de modal:
  - Sempre via:
    - `ShantillyEvent` → `EventManager` → `on:` (`OpenModal`/`ModalRequest`) → Modal Stack (`push`).
  - É proibido:
    - Criar modais diretamente em componentes, no `LayoutManager` ou no legado.
- Fechamento e resposta:
  - A interação do usuário no modal gera `ShantillyEvent` (ex.: `modal.confirmed`).
  - `EventManager`:
    - Consome o evento,
    - Executa `pop` na Modal Stack,
    - Opcionalmente monta `RunAction` final (incluindo segredos) e delega ao `ScriptRunner`.

Invariantes obrigatórios:

1. Foco exclusivo no topo
   - Apenas o modal no topo da pilha recebe eventos.
   - Enquanto houver modal ativo:
     - Interação com plano de fundo deve ser bloqueada.
   - Validação:
     - Comportamento coberto pelo gate `1.x.modal-stack.yml` e pela QA Matrix.

2. Proibição de modais "soltos"
   - É proibido:
     - Criar modais diretamente em qualquer componente ou no `LayoutManager`.
     - Abrir prompts/confirmações diretas (ex.: leitura de stdin ou caixas ad-hoc) fora da Modal Stack.
   - Toda UX modal deve ser rastreável via:
     - `ModalRequest` (`pkg/declarative`),
     - Implementação em `internal/runtime/modal/**`.

3. Integração com Segurança JIT (E1.5 parte 1)
   - `RunAction.Confirm: true`:
     - DEVE gerar `ModalRequest` de confirmação antes da execução.
   - `RunAction.PromptSecrets`:
     - DEVE gerar um ou mais `ModalRequest` para coleta de segredos.
   - Segredos:
     - Mantidos apenas em memória na cadeia Modal Stack → EventManager → ScriptRunner.
     - Proibido logar ou persistir valores sensíveis.

Rastreabilidade:

- Bloco: E1.5 — Modal Stack + Security (parte 1).
- Backlog:
  - [`docs/architecture/backlog-waves4-7-epic1.md`](docs/architecture/backlog-waves4-7-epic1.md:139) (E1.5-1, E1.5-3, E1.5-4).
- Workflows:
  - [`Workflow 3 — Segurança JIT com Modal Stack`](docs/architecture/core-workflows.md:135).
- Data models:
  - [`ModalRequest`](docs/architecture/data-models.md:208),
  - [`RunAction`](docs/architecture/data-models.md:175) (campos `Confirm` e `PromptSecrets`).
- Segurança:
  - [`docs/architecture/security.md#4-modal-stack-e-seguranca-jit-e15`](docs/architecture/security.md:132).
- Governança:
  - [`docs/architecture/governance-runtime-tui-v2.0.md#4-wave-5--modal-stack--seguranca-jit-e15--parte-1`](docs/architecture/governance-runtime-tui-v2.0.md:45).
- QA Gates:
  - `docs/qa/gates/1.x.modal-stack.yml`,
  - `docs/qa/gates/1.x.security-jit-anti-trojan.yml`.

---

## 7. Centralized Styling Theme (v2.0)

O tema é uma preocupação transversal do runtime v2.0.

Diretrizes:

- Local de referência sugerido:
  - `internal/components/theme.go` ou equivalente alinhado ao runtime declarativo.
- Todos os `ShantillyComponent`s DEVEM utilizar o tema centralizado para:
  - Consistência visual.
  - Suporte a diferentes tamanhos de terminal (validado em testes).
- O exemplo existente em `internal/tui/theme.go` é:
  - Classificado como Legacy.
  - Pode ser reaproveitado internamente por `FormComponent`,
  - Não define contratos ou limites do runtime v2.0.

Rastreabilidade:

- E1.1/E1.2 (estilo e UX básica).
- QA Matrix: linhas de styling/responsividade.

---

## 8. Legacy (não vinculante) — v1.x

Os seguintes elementos são explicitamente LEGADO e não representam a arquitetura vigente:

- `ConfigParser` baseado em `FormConfig`.
- `TUIEngine` centrado em um único formulário.
- Diagramas e textos que apresentam `shantilly form` como produto/fluxo principal.

Uso permitido:

- Como backend interno do `FormComponent` (E1.2), respeitando:
  - Contrato `ShantillyComponent`.
  - AppConfig/Component como fonte única de verdade.
- Não podem:
  - Ser expostos como API ou modelo principal.
  - Ser referenciados em novas stories/arquiteturas sem marcação explícita de Legacy.

Qualquer novo artefato (código, doc, QA) deve:

- Reafirmar o Runtime TUI Declarativo v2.0 (binário único + YAML único) como visão única.
- Mapear comportamentos sempre via:
  - `AppConfig` → `LayoutNode` → `Component` → `ShantillyComponent` → `ShantillyEvent` → `EventManager` → `on:` → `RunAction`/`ModalRequest` → `ScriptRunner`/Modal Stack.
