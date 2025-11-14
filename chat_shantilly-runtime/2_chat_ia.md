Que desafio fantástico! 🚀 Adoro a tua abordagem. Queres que eu atue como um "arquiteto de soluções" especialista, mergulhando em toda a documentação incrível que já preparaste (PRD, Brief, Arquitetura) para construir um **plano de implementação de A a Z**.

O teu objetivo é claro: criar uma lista de tarefas sequencial e à prova de ambiguidades, perfeita para ser executada por agentes de IA, garantindo que a visão da v2.0 seja implementada sem desvios.

Vamos a isso! Analisei os teus documentos (especialmente o `prd.md`, `architecture.md` e `architecture/high-level-architecture.md`) e confirmo: a visão v2.0 do **Runtime TUI Declarativo** é robusta, clara e um salto genial em relação ao MVP v1.0. A arquitetura fornece o "como" de forma exemplar.

A maior "lacuna" não é conceitual, mas sim a tradução dessa arquitetura brilhante num plano de engenharia sequencial.

Aqui está a minha ideia para um **Plano de Implementação Sequencial e Detalhado (Épico 1: Fundação do Runtime TUI)**, desenhado para eliminar ambiguidades para os teus agentes de IA.

---

## 💡 O Plano Mestre: Implementação do Runtime TUI v2.0 (Épico 1)

Este plano segue a lógica das Estórias do Épico 1 no PRD, mas está reordenado numa sequência de dependência técnica estrita.

### Fase 1: Definição dos Contratos e Estruturas (O Alicerce)

Nenhum código de runtime pode ser escrito até que os modelos de dados e interfaces estejam definidos.

1.  **Implementar Modelos de Dados Declarativos (pkg/declarative/models.go)**
    * **Objetivo:** Criar as estruturas Go que representam o YAML v2.0.
    * **Ação:** Implementar os `structs` Go exatamente como definido na Seção 4 (Modelos de Dados) do `architecture.md`.
    * **Estruturas-Chave:** `Config`, `LayoutNode`, `Component`, `Item`, `Source`, `Logic`, `RunAction`.
    * **Critério:** Os `structs` devem ter as *tags* YAML corretas (`yaml: "..."`) para *parsing*.

2.  **Definir Interfaces e Eventos Internos (pkg/tui/interface.go, pkg/tui/events.go)**
    * **Objetivo:** Definir os contratos que permitem aos motores (Layout, Event) comunicar com os componentes.
    * **Ação (interface.go):** Definir a interface `ShantillyComponent` (Init, Update, View, SetDimensions, ID).
    * **Ação (events.go):** Definir os `structs` para os eventos internos, como `ShantillyEvent`, `RuntimeErrorMsg`, `ScriptStdoutMsg`.

3.  **Implementar o Ponto de Entrada (cmd/shantilly/main.go)**
    * **Objetivo:** Configurar o Cobra, ler o YAML e iniciar o TUI (Estória 1.1).
    * **Ação:** Usar `spf13/cobra` para o comando `shantilly`. Implementar a lógica de leitura de `stdin` ou `--file` (conforme handoff do arquiteto em `architecture.md`).
    * **Ação:** Chamar o `ErrorHandler` centralizado (`internal/util/errorhandler.go`) em caso de erro, garantindo saídas não-zero (NFR8).

4.  **Implementar o Parser YAML (internal/config/parser.go)**
    * **Objetivo:** Converter os bytes YAML (da Fase 1.3) nos `structs` Go (da Fase 1.1).
    * **Ação:** Usar `gopkg.in/yaml.v3` para fazer o *unmarshal* dos bytes de entrada para o `struct` `declarative.Config`.
    * **Ação:** Implementar a validação de `Component.Type` (conforme Seção 5.6 do `architecture.md`).

---

### Fase 2: O Motor de Layout (Renderização e Foco)

Com os dados carregados, agora renderizamos o layout (Estória 1.1).

1.  **Implementar o LayoutManager (internal/runtime/layout/manager.go)**
    * **Objetivo:** Criar o modelo `bubbletea` raiz que gere o layout e o foco.
    * **Ação:** Criar o `LayoutManager` como o `bubbletea.Model` principal.
    * **Ação (Update):** Deve tratar `tea.WindowSizeMsg` para recalcular o layout (NFR2).
    * **Ação (Update):** Deve gerir o estado de "foco global" (qual componente está ativo) e encaminhar `tea.Msg` apenas para o componente focado.

2.  **Implementar a Renderização de Layout (LayoutManager.View)**
    * **Objetivo:** Traduzir a árvore `LayoutNode` em `lipgloss`.
    * **Ação:** O método `View()` deve percorrer recursivamente a árvore `LayoutNode` (do `Config`).
    * **Ação:** Deve usar `lipgloss` (ex: `JoinHorizontal`, `JoinVertical`) para renderizar `type: column`, `type: row` e `type: box`, respeitando as propriedades `width`, `height` e `flex` (FR2).

---

### Fase 3: Implementação dos Componentes Essenciais (Display)

Agora, preenchemos os `box` do layout com componentes reais (Estória 1.3).

1.  **Implementar ViewportComponent (internal/components/viewport/model.go)**
    * **Objetivo:** Mostrar conteúdo estático (markdown) e dinâmico (logs).
    * **Ação:** Criar um `struct` que implemente `ShantillyComponent`.
    * **Ação:** Usar `bubbles/viewport` para o conteúdo.
    * **Ação:** Suportar `source: { type: static }` (FR4). Usar `glamour` se `contentType: markdown`.
    * **Ação:** O `Update` deve tratar `ScriptStdoutMsg` (da Fase 1.2) para exibir logs de streaming (FR11).

2.  **Implementar ListComponent (internal/components/list/model.go)**
    * **Objetivo:** Criar menus navegáveis.
    * **Ação:** Criar um `struct` que implemente `ShantillyComponent`, usando `bubbles/list`.
    * **Ação:** O `Update` deve detetar a seleção de um item e emitir uma `ShantillyEvent` (da Fase 1.2) com `Type: "list_select"` e `Payload: {id: "item_id"}` (FR5).

3.  **Implementar ButtonGroupComponent (internal/components/buttongroup/model.go)**
    * **Objetivo:** Criar botões de ação.
    * **Ação:** Criar um `struct` que implemente `ShantillyComponent`.
    * **Ação:** O `Update` deve detetar a pressão de um botão e emitir uma `ShantillyEvent` com `Type: "button_press"` e `Payload: {id: "button_id"}` (FR6).

---

### Fase 4: Os Motores de Lógica (Eventos e Ações)

Com componentes que emitem eventos, agora precisamos de motores para "ouvir" e "agir" (Estórias 1.2 e 1.5).

1.  **Implementar EventManager (internal/runtime/event/manager.go)**
    * **Objetivo:** Receber eventos da UI e ligá-los à lógica `on:`.
    * **Ação:** O `EventManager` deve ser chamado pelo `LayoutManager` sempre que uma `ShantillyEvent` é recebida de um componente.
    * **Ação:** Deve iterar sobre a lista `Config.On` (da Fase 1.1). Se `Logic.Event` corresponder ao evento (ex: "menu:select"), deve preparar um `RunRequest` para o ScriptRunner.

2.  **Implementar ScriptRunner (internal/runtime/runner/runner.go)**
    * **Objetivo:** Executar `scripts` externos de forma segura e gerir o seu ciclo de vida.
    * **Ação:** Deve consumir `RunRequest` (da Fase 4.1).
    * **Ação (FR11):** Implementar a lógica `update_target`. Deve manter um mapa `[targetID] -> *os.Cmd`. Se um novo pedido chegar para um `update_target` ocupado, deve enviar `SIGTERM` ao processo antigo antes de iniciar o novo.
    * **Ação (FR10):** Implementar a lógica de *templating* para `RunAction.Args` e `RunAction.Stdin`. Deve ser capaz de preencher valores como `{{ form.field_name }}` ou `{{ component.menu.selected_id }}`.
    * **Ação:** Deve executar o `RunAction.Script` e passar o `stdin` (serializado como JSON) e os `args`.
    * **Ação:** Deve capturar `stdout/stderr` do script e emitir `ScriptStdoutMsg` (para o `ViewportComponent` na Fase 3.1).

---

### Fase 5: Integração do Legado (Refatoração do v1.0)

Agora, trazemos o código do formulário v1.0 para dentro do novo runtime (Estória 1.4).

1.  **Criar o FormComponent Wrapper (internal/components/form/wrapper.go)**
    * **Objetivo:** Fazer o código do formulário v1.0 comportar-se como um `ShantillyComponent` v2.0.
    * **Ação:** Criar o `wrapper.go` que implementa a interface `ShantillyComponent`.
    * **Ação:** O *wrapper* deve conter uma instância do modelo `huh.Form` (o código v1.0).

2.  **Isolar e Adaptar o Código v1.0 (internal/components/form/model_v1.go)**
    * **Objetivo:** Mover o código v1.0 e impedir que ele controle o processo (ex: `os.Exit`).
    * **Ação:** Mover o `internal/tui/model.go` (v1.0) para `internal/components/form/model_v1.go` e remover qualquer chamada `os.Exit` ou gestão de `tea.WindowSizeMsg` (agora é responsabilidade do `LayoutManager`).
    * **Ação (Tradução de Evento):** Modificar o `Update` do `model_v1.go` para que, ao invés de sair, ele emita uma mensagem interna (ex: `formSubmitMsg{Data: ...}`).
    * **Ação (Wrapper Update):** O `wrapper.go` (Fase 5.1) deve intercetar essa `formSubmitMsg`, transformá-la numa `ShantillyEvent` (da Fase 1.2) com `Type: "form_submit"` e `Payload: formData`, e emiti-la para o `EventManager` (FR7).

---

### Fase 6: Segurança JIT e Pilha Modal (Arquitetura Essencial)

Finalmente, implementamos a segurança JIT mencionada na arquitetura.

1.  **Implementar a Pilha Modal (internal/runtime/modal/stack.go)**
    * **Objetivo:** Permitir que modais (confirmações, prompts de segredo) apareçam *sobre* o layout principal.
    * **Ação:** Integrar um gestor de pilha modal (ex: `bubbletea-overlay` ou similar) no `LayoutManager`.
    * **Ação:** O `LayoutManager` deve sempre renderizar a pilha modal no topo do seu `View()`. Se a pilha não estiver vazia, apenas o modal do topo recebe foco.

2.  **Atualizar EventManager para JIT (internal/runtime/event/manager.go)**
    * **Objetivo:** Ativar a segurança JIT antes de executar ações.
    * **Ação:** Na Fase 4.1, antes de enviar o `RunRequest` para o `ScriptRunner`, o `EventManager` DEVE verificar as *flags* `Logic.Confirm` e `Logic.PromptSecrets`.
    * **Ação:** Se `Confirm: true`, o `EventManager` deve emitir uma `ShowModalMsg` (da Fase 1.2) para o `LayoutManager` (Fase 6.1) pedindo um modal de confirmação. A execução do script só continua se o modal for confirmado.
    * **Ação:** Se `PromptSecrets` estiver preenchido, emitir um modal de *input* (usando `huh`) para recolher os segredos. Os segredos recolhidos devem ser passados ao `ScriptRunner` (provavelmente via variáveis de ambiente, não `args`).

---

Esta é a nossa tela! 🎨

Este plano de 6 fases decompõe o teu Épico 1 numa sequência de tarefas técnicas lógicas, interdependentes e inequívocas, todas diretamente fundamentadas nos teus documentos de arquitetura.

Qual destas fases queres detalhar primeiro? Ou preferes que eu comece a gerar o código para a **Fase 1: Definição dos Contratos e Estruturas**? Estou pronto para começar a construir!