Absolutamente. Esta é a abordagem correta para um projeto desta complexidade. Vamos usar o teu "conhecimento privilegiado" (os documentos PRD v2.0 e Arquitetura v2.0) para desenhar o plano de implementação ideal, do zero, focado em eliminar os riscos e ambiguidades já identificados.

Aqui está o **Plano de Implementação Sequencial Mestre (Épicos 1-4)**, desenhado para uma equipa de desenvolvimento que começa hoje.

---

## 💡 Plano Mestre: Implementação Sequencial (Épicos 1-4)

Este plano é rigorosamente sequencial. Cada Épico depende da conclusão bem-sucedida do anterior, e cada Fase depende da anterior.

### ÉPICO 1: Fundação do Runtime TUI (Genérico)

**Meta:** Construir o motor central extensível v2.0.
***Otimização de Risco (Visão do Futuro):*** *Evitar a armadilha do MVP v1.0 (form-only). A arquitetura deve ser modular (Layout vs. Lógica) desde o primeiro dia.*

#### Fase 1.1: Definição dos Contratos (O Alicerce)
*(Dependência: Nenhuma)*

1.  **Tarefa 1.1.1: Definir Modelos de Dados YAML (pkg/declarative/models.go)**
    * Implementar os `structs` Go (`Config`, `LayoutNode`, `Component`, `Logic`, `RunAction`) com *tags* YAML, conforme a Seção 4 do `architecture.md`.
2.  **Tarefa 1.1.2: Definir Contratos Internos (pkg/tui/interface.go, pkg/tui/events.go)**
    * Definir a interface `ShantillyComponent` (Init, Update, View, etc.).
    * Definir os `structs` de eventos (`ShantillyEvent`, `RuntimeErrorMsg`, `ScriptStdoutMsg`).
3.  **Tarefa 1.1.3: Implementar Ponto de Entrada (cmd/shantilly/main.go)**
    * Configurar `cobra` e o `ErrorHandler` centralizado (`internal/util/errorhandler.go`) para garantir `stderr` e códigos de saída não-zero (NFR1).
4.  **Tarefa 1.1.4: Implementar Parser YAML (internal/config/parser.go)**
    * Ler `stdin`/`--file` e fazer *unmarshal* para o `struct declarative.Config` (da Tarefa 1.1.1).

#### Fase 1.2: Motor de Layout (Renderização e Foco)
*(Dependência: Fase 1.1)*

1.  **Tarefa 1.2.1: Implementar LayoutManager (internal/runtime/layout/manager.go)**
    * Implementar o `bubbletea.Model` raiz.
    * ***Otimização de Risco:*** *Este modelo gere **apenas** UI (Foco Global, Pilha Modal) e recalcula layout em `tea.WindowSizeMsg` (NFR2), mas **não** gere lógica de automação*.
2.  **Tarefa 1.2.2: Implementar Renderização de Layout (LayoutManager.View)**
    * Implementar a renderização recursiva da árvore `LayoutNode` (FR1) usando `lipgloss` para respeitar `width`, `height` e `flex` (FR2).

#### Fase 1.3: Componentes Essenciais (Display)
*(Dependência: Fase 1.1, 1.2. LayoutManager precisa instanciá-los)*

1.  **Tarefa 1.3.1: Implementar ViewportComponent (internal/components/viewport/model.go)**
    * Implementar `ShantillyComponent` usando `bubbles/viewport`.
    * Suportar `source: { type: static }` e renderizar markdown (FR4).
2.  **Tarefa 1.3.2: Implementar ListComponent (internal/components/list/model.go)**
    * Implementar `ShantillyComponent` usando `bubbles/list`.
    * Ao selecionar, deve emitir a `ShantillyEvent` "list:select" (FR5).
3.  **Tarefa 1.3.3: Implementar ButtonGroupComponent (internal/components/buttongroup/model.go)**
    * Implementar `ShantillyComponent`.
    * Ao pressionar, deve emitir a `ShantillyEvent` "button:press" (FR6).

#### Fase 1.4: Motores de Lógica (Eventos e Ações)
*(Dependência: Fase 1.1 (Modelos), Fase 1.3 (Eventos))*

1.  **Tarefa 1.4.1: Implementar EventManager (internal/runtime/event/manager.go)**
    * Receber `ShantillyEvent` (da Fase 1.3) e mapear para a lista `Config.On` (FR8).
2.  **Tarefa 1.4.2: Implementar ScriptRunner (internal/runtime/runner/runner.go)**
    * Executar `run: { script: ... }` (FR9).
    * ***Otimização de Risco:*** *Implementar a lógica `update_target` (FR11). Gerir processos ativos e enviar `SIGTERM` a processos anteriores no mesmo *target* antes de iniciar um novo*.
3.  **Tarefa 1.4.3: Implementar Fluxo de Dados (runner.go)**
    * Implementar *templating* e passagem de `args` e `stdin` (JSON) (FR10).
4.  **Tarefa 1.4.4: Conectar Lógica ao Display**
    * Fazer o `ScriptRunner` emitir `ScriptStdoutMsg` (da Tarefa 1.1.2).
    * Fazer o `ViewportComponent` (da Tarefa 1.3.1) tratar `ScriptStdoutMsg` para streaming.

#### Fase 1.5: Integração do Legado (v1.0)
*(Dependência: Fase 1.1 (Interface), Fase 1.4 (EventManager))*

1.  **Tarefa 1.5.1: Criar FormComponent Wrapper (internal/components/form/wrapper.go)**
    * Criar o *wrapper* que implementa `ShantillyComponent`.
2.  **Tarefa 1.5.2: Encapsular Código v1.0 (internal/components/form/model_v1.go)**
    * ***Otimização de Risco:*** *Mover o código `huh` v1.0 para este pacote e **remover** todas as chamadas `os.Exit` e gestão de `WindowSizeMsg`.*
    * Modificar o *submit* para emitir `ShantillyEvent` "form:submit" (FR7), que será capturado pelo `EventManager` (Tarefa 1.4.1).

#### Fase 1.6: Segurança JIT (Pilha Modal)
*(Dependência: Fase 1.2 (LayoutManager), Fase 1.4 (EventManager))*

1.  **Tarefa 1.6.1: Implementar Pilha Modal (internal/runtime/modal/stack.go)**
    * Implementar um gestor de pilha modal (ex: `bubbletea-overlay`) e integrá-lo ao `LayoutManager` (Tarefa 1.2.1).
2.  **Tarefa 1.6.2: Atualizar EventManager para JIT**
    * ***Otimização de Risco:*** *Modificar `EventManager` (Tarefa 1.4.1) para que, ANTES de disparar o `ScriptRunner`, verifique `Logic.Confirm` e `Logic.PromptSecrets`. Se presentes, deve disparar Modais JIT (via `ShowModalMsg`) e aguardar a resposta antes de continuar*.

---

### ÉPICO 2: O Runner Especialista (Ansible Fase 2)

**Meta:** Implementar o *runner* de conveniência `ansible_playbook`.
*(Dependência: Épico 1 (especificamente Fases 1.1, 1.4 e 1.6))*

#### Fase 2.1: Extensão dos Contratos
1.  **Tarefa 2.1.1: Atualizar Modelos de Dados (pkg/declarative/models.go)**
    * Adicionar o `struct` para `AnsiblePlaybookAction` (incluindo `vars` e `ask_vault_pass`) ao `struct RunAction`.

#### Fase 2.2: Implementação do Runner Ansible
1.  **Tarefa 2.2.1: Atualizar EventManager**
    * Modificar `EventManager` (Tarefa 1.4.1) para reconhecer e disparar o novo *runner* `ansible_playbook`.
2.  **Tarefa 2.2.2: Implementar AnsibleRunner (internal/runtime/runner/ansible.go)**
    * Criar a lógica que traduz o `struct AnsiblePlaybookAction` em comandos `ansible-playbook` (gestão de `vars`, etc.).
3.  **Tarefa 2.2.3: Integrar com Segurança JIT**
    * Integrar o `AnsibleRunner` com a Fase 1.6. Se `ask_vault_pass: true`, deve disparar o Modal JIT para recolher a *vault pass*.

---

### ÉPICO 3: O Runtime Preditivo (Ansible Fase 3)

**Meta:** Implementar os componentes de descoberta `playbook_explorer` e `inventory_explorer`.
*(Dependência: Épico 1 (especificamente Fases 1.1 e 1.2))*

#### Fase 3.1: Extensão dos Contratos
1.  **Tarefa 3.1.1: Atualizar Modelos de Dados (pkg/declarative/models.go)**
    * Adicionar `type: playbook_explorer` e `type: inventory_explorer` às opções do `struct Component`.

#### Fase 3.2: Implementação dos Componentes Preditivos
1.  **Tarefa 3.2.1: Implementar PlaybookExplorerComponent**
    * Criar o novo componente (implementando `ShantillyComponent`).
    * ***Otimização de Risco:*** *Implementar a filtragem de ruído (`roles/`, `tasks/`) conforme NFR4*.
2.  **Tarefa 3.2.2: Implementar InventoryExplorerComponent**
    * Criar o novo componente (implementando `ShantillyComponent`).
    * ***Otimização de Risco:*** *Usar `ansible-inventory` como "Oráculo" para garantir precisão, conforme NFR4*.
3.  **Tarefa 3.2.3: Atualizar LayoutManager**
    * Modificar a "fábrica" de componentes no `LayoutManager` (Tarefa 1.2.1) para ser capaz de instanciar estes novos tipos de componentes.

---

### ÉPICO 4: Administração SSH (Visão de Longo Prazo)

**Meta:** Servir o Runtime TUI sobre SSH.
*(Dependência: Épico 1 (todo o runtime))*

#### Fase 4.1: Implementação do Servidor SSH
1.  **Tarefa 4.1.1: Adicionar Subcomando `serve-ssh`**
    * Adicionar o novo comando ao `main.go` (Tarefa 1.1.3).
2.  **Tarefa 4.1.2: Integrar `charmbracelet/wish`**
    * Implementar o servidor SSH básico usando `wish`.

#### Fase 4.2: Integração do Runtime TUI
1.  **Tarefa 4.2.1: Implementar *Handler* SSH**
    * Criar o *handler* `wish` que, em vez de um shell, inicia uma nova instância do `LayoutManager` (Fase 1.2) para a sessão `ssh.Session`.
2.  **Tarefa 4.2.2: Definir Configuração de Serviço**
    * Implementar a lógica que define qual YAML o servidor SSH deve carregar para o Runtime TUI (ex: um YAML fixo por servidor ou baseado na chave SSH).