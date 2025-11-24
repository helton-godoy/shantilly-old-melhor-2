# Shantilly Runtime v2.0 – Migration Notes

## 1. O que já foi migrado

- **Runtime declarativo v2.0 (wiring)**  
  - Arquivo `internal/runtime/runtime.go` criado como ponto de entrada do runtime declarativo.  
  - Integração de `LayoutManager`, `EventManager` e `ScriptRunner` no ciclo Bubble Tea.  
  - Uso de `pkg/declarative.AppConfig` como modelo de configuração YAML única.

- **Layout e Component Registry**  
  - `internal/runtime/layout/manager.go` atualizado como LayoutManager raiz, consumindo `LayoutNodeRef` derivado de `pkg/declarative.LayoutNode`.  
  - `internal/runtime/layout/registry.go` criado como `DefaultRegistry` para resolver `declarative.Component` em componentes concretos.

- **Componentes TUI oficiais (v2.0)**  
  - `internal/components/list/model.go`  
    - Componente de lista/menu baseado em `bubbles/list`.  
    - Emite `tui.ShantillyEvent` com `Type="select"` e payload com ID do item.
  - `internal/components/viewport/model.go`  
    - Componente de viewport com suporte a conteúdo estático / streaming de saída de scripts.  
    - Integração com mensagens de streaming (`ScriptStdoutMsg`, etc.).
  - `internal/components/buttongroup/model.go`  
    - Grupo de botões lógicos para ações (ex.: submit, cancel).  
    - Emite `tui.ShantillyEvent` com `Type="press"` e payload com ID do botão.

- **Contratos compartilhados (pkg)**  
  - `pkg/tui/interface.go`  
    - Define interface `ShantillyComponent` (Init, Update, View, SetDimensions, ID).  
  - `pkg/tui/events.go`  
    - Define tipos de mensagem para integração runtime ↔ componentes ↔ script runner (eventos genéricos, mensagens de erro, streaming, modais, RunScriptRequestMsg).  
  - `pkg/declarative/models.go`  
    - Define `AppConfig`, `LayoutNode`, `Component`, `OnHandler`, `RunAction`, `ShantillyEvent`, `SecurityPolicy` e funções de parsing/validação (`LoadAppConfig`, `Validate`).  
    - Adicionados tipos auxiliares `Item` e `Source` para componentes.

- **Event Engine / JIT security (esqueleto)**  
  - `internal/runtime/event/manager.go`  
    - `NewManager` com regras `OnHandler`.  
    - `ProcessEvent` normaliza eventos `componentID:type` e resolve handlers `on:`.  
    - Suporte a `RunAction` com campo `Confirm` (JIT – usa modal de confirmação via `ShowModalMsg`).  
    - `HandleModalResult` para retomar execução após confirmação.

- **Script Runner / integração com TEA**  
  - `internal/runtime/runner/templating.go`  
    - Motor de templates para stdin/env/strings (Wave 1).  
  - `internal/runtime/runner/ansible.go`  
    - Esqueleto de runner especializado para `ansible_playbook`.  
  - `internal/runtime/runner/runner.go`  
    - Execução declarativa de `RunAction` (scripts externos), com streaming de stdout/stderr.  
  - `internal/runtime/runner/tea_adapter.go`  
    - Adaptador `HandleRunRequest` conectando `tui.RunScriptRequestMsg` ao `ScriptRunner` e emitindo mensagens de streaming/termino para o TEA.

- **Modal de confirmação JIT**  
  - `internal/runtime/modal/confirm.go`  
    - Componente TUI para confirmação de execução de scripts sensíveis.
  - `internal/runtime/main_model.go`  
    - Modelo raiz simplificado, preparado para orquestrar layout + modal (sem dependência externa de overlay).

- **CLI / Entrypoints**  
  - `cmd/shantilly/main.go`  
    - Mantém comando legado `form` baseado em `FormConfig`.  
    - Novo subcomando `runtime` lendo YAML declarativo (`AppConfig`) e chamando `internal/runtime.Start`.

- **Configuração / Exemplo de app**  
  - `app.example.yaml`  
    - Exemplo mínimo de AppConfig v2.0 com layout (list + viewport), componentes e regra `on:` com `run:` e `confirm`.  
  - `internal/config/form_config.go`  
    - Reintroduz tipos legados `FormConfig` e `Field` para compatibilidade com o comando `form` v1.x.

- **Infra mínima e tooling**  
  - `go.mod` atualizado para dependências de Bubble Tea / Bubbles / Lipgloss / Cobra / Glamour, etc.  
  - Scripts auxiliares em `scripts/` preservados para futuras waves.

---

## 2. Pendências atuais (estado da branch `feat/runtime-migration`)

- **Build ainda não está verde**  
  - Erros anteriores incluíam:  
    - Falta de tipos `Item`/`Source` em `pkg/declarative` (já adicionados).  
    - Variável `confirmMsg` não utilizada em `internal/runtime/event/manager.go` (já removida).  
    - Ajustes de interface `ShantillyComponent` x implementações concretas.
  - Último erro observado antes do push:  
    - Incompatibilidade de interface em `internal/runtime/layout/registry.go` quando retornando componentes (`SetBounds` vs `SetDimensions`).

- **Coerência de contratos**  
  - `internal/runtime/layout/manager.go` possuía sua própria interface `ShantillyComponent` com método `SetBounds(width, height int)`.  
  - `pkg/tui/interface.go` define `ShantillyComponent` com `SetDimensions(width, height int)`.  
  - Foi iniciado o processo de alinhar o LayoutManager para usar o contrato do pacote `pkg/tui` (alias de tipo), porém é necessário revisar o arquivo para garantir que **não haja referências restantes a `SetBounds`**.

- **Decodificação de props declarativos**  
  - `internal/runtime/layout/registry.go` ainda assume `Props["items"]` como `[]declarative.Item` e `Props["source"]` como `declarative.Source` direto, mas o `yaml.Unmarshal` produz `map[string]any` com `[]any` etc.  
  - É necessário implementar uma etapa de conversão segura (ex.: helper que converte `[]any` → `[]declarative.Item`, `map[string]any` → `declarative.Source`).

- **Testes automatizados**  
  - Muitos testes antigos foram removidos (legado v1.x).  
  - Ainda não há testes específicos para o runtime v2.0 (LayoutManager, EventManager, ScriptRunner, componentes declarativos).  
  - Será necessário planejar uma nova suíte mínima para garantir estabilidade.

---

## 3. Próximos passos imediatos

1. **Concluir alinhamento de interface de componentes**  
   - Revisar `internal/runtime/layout/manager.go` para confirmar que:  
     - `ShantillyComponent` é um alias para `tui.ShantillyComponent`.  
     - Não existem chamadas a `SetBounds`; apenas `SetDimensions` é utilizado.  
   - Garantir que `internal/components/{list,viewport,buttongroup}` implementam exatamente `tui.ShantillyComponent`.

2. **Ajustar `DefaultRegistry` para props tipados**  
   - Em `internal/runtime/layout/registry.go`:  
     - Implementar conversão de `comp.Props["items"]` (tipos genéricos vindos do YAML) para `[]declarative.Item`.  
     - Implementar conversão de `comp.Props["source"]` (provavelmente `map[string]any`) para `declarative.Source`.  
   - Adicionar fallbacks seguros (ex.: mensagem de erro amigável em viewport se conversão falhar).

3. **Rodar ciclo de build local**  
   - Usar Go instalado em `/usr/local/go/bin/go` (se necessário):  
     ```bash
     /usr/local/go/bin/go fmt ./...
     /usr/local/go/bin/go mod tidy
     /usr/local/go/bin/go build -o shantilly ./cmd/shantilly
     ```
   - Corrigir qualquer erro de compilação remanescente.

4. **Testar runtime com exemplo**  
   - Após build bem-sucedido:  
     ```bash
     ./shantilly runtime --file app.example.yaml
     ```
   - Verificar:  
     - Renderização da lista e viewport.  
     - Seleção na lista disparando evento e eventual execução de script (mesmo que script seja dummy).  
     - Modal de confirmação abrindo quando `RunAction.Confirm == true`.

5. **Documentar decisões e diferenças vs. arquitetura original**  
   - Registrar no final deste arquivo quaisquer desvios intencionais da especificação de arquitetura (ex.: simplificações de overlay, runners parciais).  
   - Atualizar este documento a cada wave (por exemplo, quando implementar conversão completa de props, Ansible runner real, Modal Stack avançada).

---

## 4. Como retomar em outra máquina

1. Clonar o repositório e trocar para a branch de migração:
   ```bash
   git clone https://github.com/helton-godoy/shantilly.git
   cd shantilly
   git checkout feat/runtime-migration
   ```

2. Instalar Go (>= 1.21) e garantir que o binário `go` esteja no PATH ou em `/usr/local/go/bin/go`.

3. Seguir os passos da seção **3. Próximos passos imediatos** até obter um build verde e o runtime rodando com `app.example.yaml`.
