# Data Models

Este documento define os modelos de dados NORMATIVOS do Runtime TUI Declarativo v2.0 (Epic 1), alinhados com:

- PRD pivot: [`docs/prd/epic-1-runtime-tui-foundation.md`](docs/prd/epic-1-runtime-tui-foundation.md:1)
- Arquitetura pivot: [`docs/architecture/introduction.md`](docs/architecture/introduction.md:1)
- Story Map: [`docs/stories/index-epic-1-runtime-tui.story-map.md`](docs/stories/index-epic-1-runtime-tui.story-map.md:1)
- Plano de Implementação: [`docs/architecture/implementation-plan-epic-1-runtime-tui.md`](docs/architecture/implementation-plan-epic-1-runtime-tui.md:1)
- Matriz QA: [`docs/qa/matrix-epic-1-runtime-tui-coverage.md`](docs/qa/matrix-epic-1-runtime-tui-coverage.md:1)

Eles materializam o contrato do YAML único (E1.1–E1.5) e substituem o modelo linear centrado em `FormConfig` como visão ativa. `FormConfig` passa a ser LEGADO encapsulado no `FormComponent`.

## 1. AppConfig (E1.1–E1.5 — Raiz do YAML Único)

Representa o documento YAML único consumido pelo runtime v2.0.

Campos normativos (ilustrativo, nomeação sugerida):

```go
// pkg/declarative/app.go

type AppConfig struct {
    Version      string            `yaml:"version,omitempty" json:"version,omitempty"`
    Metadata     map[string]string `yaml:"meta,omitempty" json:"meta,omitempty"`
    Layout       LayoutNode        `yaml:"layout" json:"layout"`
    Components   []Component       `yaml:"components" json:"components"`
    On           []OnHandler       `yaml:"on,omitempty" json:"on,omitempty"`
    Security     *SecurityPolicy   `yaml:"security,omitempty" json:"security,omitempty"`
}
```

Invariantes:

- Única raiz para:
  - Layout (E1.1),
  - Componentes (E1.2),
  - Regras `on:` (E1.3),
  - Ações `run:` (E1.4),
  - Segurança/Modal Stack (E1.5).
- Não existe segundo arquivo ou subtipo paralelo redefinindo estes conceitos.

Rastreabilidade:

- E1.1–E1.5 (todas).
- Story Map: blocos E1.x.
- QA Matrix: linhas 1–11.

## 2. LayoutNode (E1.1 — LayoutManager)

Define a árvore de layout hierárquico declarativo que alimenta diretamente o `LayoutManager` (root model TEA).

```go
// pkg/declarative/layout.go

type LayoutNode struct {
    ID         string       `yaml:"id,omitempty" json:"id,omitempty"`
    Type       string       `yaml:"type" json:"type"` // "column" | "row" | "box"
    Width      *int         `yaml:"width,omitempty" json:"width,omitempty"`
    Height     *int         `yaml:"height,omitempty" json:"height,omitempty"`
    Flex       *int         `yaml:"flex,omitempty" json:"flex,omitempty"`
    Items      []LayoutNode `yaml:"items,omitempty" json:"items,omitempty"`
    // Apenas para type: box:
    ComponentID string `yaml:"component,omitempty" json:"component,omitempty"`
}
```

Invariantes (vinculantes para Wave 2):

- `Type` restrito ao conjunto permitido (`column|row|box`).
- Para `type: box`:
  - Pode referenciar no máximo um `ComponentID`.
  - `ComponentID` DEVE apontar para um `Component.ID` definido em `AppConfig.Components`.
- `LayoutNode` não contém lógica de negócio ou automação (`on:`/`run:`); é puramente estrutural.
- `LayoutNode` é a ÚNICA entrada de layout para o `LayoutManager` em `internal/runtime/layout`:
  - Nenhuma outra struct (incluindo `FormConfig`) pode ser usada como contrato ativo de layout.
- Qualquer uso de modelos v1.x é explicitamente Legacy e restrito ao backend do `FormComponent`.

Rastreabilidade:

- PRD: layout (Epic 1, seção LayoutManager).
- Story Map: E1.1 (stories 1.2, 1.3, 1.7 alinhadas ao modelo declarativo).
- QA Matrix: bloco E1.1 — Layout Engine (validação de LayoutNode + consumo pelo LayoutManager).

## 3. Component (E1.2 — ShantillyComponent + Componentes Oficiais)

Representa componentes TUI declarativos instanciados no layout.

```go
// pkg/declarative/component.go

type Component struct {
    ID       string            `yaml:"id" json:"id"`                 // Referenciado por layout, on:, run:
    Type     string            `yaml:"type" json:"type"`             // "list" | "viewport" | "form" | "buttongroup"
    Props    map[string]any    `yaml:"props,omitempty" json:"props,omitempty"`
    Bindings map[string]string `yaml:"bind,omitempty" json:"bind,omitempty"` // Ex.: estados, seleção, etc.
}
```

Invariantes:

- `ID` é único dentro do AppConfig.
- `Type` pertence aos componentes oficiais do Epic 1:
  - `list`, `viewport`, `form`, `buttongroup`.
- Implementação em runtime:
  - Cada componente concreto implementa a interface `ShantillyComponent` (ver [`docs/architecture/components.md`](docs/architecture/components.md:1)).

Uso do legado:

- Para `Type: "form"`, o runtime pode reutilizar partes do código v1.x dentro do `FormComponent`, mas:
  - Sempre por meio do contrato `ShantillyComponent` + `Component.Props`,
  - Nunca por fluxos paralelos ao runtime v2.0.

Rastreabilidade:

- E1.2.
- Story Map: E1.2.
- QA Matrix: blocos E1.2, Legacy encapsulado.

## 4. OnHandler (E1.3 — Bloco `on:` como única orquestração)

Modela uma regra declarativa `on:`.

```go
// pkg/declarative/on.go

type OnHandler struct {
    ID           string        `yaml:"id,omitempty" json:"id,omitempty"`
    Event        string        `yaml:"event" json:"event"` // ex.: "menu:select", "user_form:submit"
    When         string        `yaml:"when,omitempty" json:"when,omitempty"` // expressão opcional
    Run          *RunAction    `yaml:"run,omitempty" json:"run,omitempty"`
    OpenModal    *ModalRequest `yaml:"open_modal,omitempty" json:"open_modal,omitempty"`
    UpdateState  map[string]any `yaml:"update_state,omitempty" json:"update_state,omitempty"`
}
```

Invariantes:

- `Event` vincula-se ao modelo `ShantillyEvent` (ver seção seguinte).
- Toda automação:
  - Deve ser expressa via `OnHandler` + `RunAction` + ações declarativas.
  - Handlers “soltos” no código fora deste pipeline são proibidos.

Rastreabilidade:

- E1.3.
- Story Map: E1.3.
- QA Matrix: Event Engine.

## 5. ShantillyEvent (E1.3 — Tipo de evento normalizado)

Tipo canônico de evento usado pelo `EventManager`.

```go
// pkg/declarative/event.go

type ShantillyEvent struct {
    SourceComponentID string                 `json:"source_component_id"`
    Type              string                 `json:"type"`              // ex.: "submit", "select", "click", "script.complete", "modal.confirmed"
    Payload           map[string]interface{} `json:"payload,omitempty"` // dados contextuais
}
```

Invariantes:

- Eventos são categorizados (UI, ScriptRunner, Modal, Sistema).
- Todos os eventos roteados passam pelo `EventManager` (`internal/runtime/event`).
- `OnHandler.Event` deve referenciar combinações válidas (`component_id:type` ou padrões compatíveis).

Rastreabilidade:

- E1.3.
- Story Map: E1.3.
- QA Matrix: linhas Event Engine.

## 6. RunAction (E1.4 — Execução declarativa)

Especifica ações `run:` disparadas por `OnHandler`.

```go
// pkg/declarative/run.go

type RunAction struct {
    Script       string                 `yaml:"script" json:"script"`
    Args         []string               `yaml:"args,omitempty" json:"args,omitempty"`
    Stdin        map[string]interface{} `yaml:"stdin,omitempty" json:"stdin,omitempty"`
    Env          map[string]string      `yaml:"env,omitempty" json:"env,omitempty"`
    UpdateTarget string                 `yaml:"update_target,omitempty" json:"update_target,omitempty"`
    Confirm      bool                   `yaml:"confirm,omitempty" json:"confirm,omitempty"`
    PromptSecrets []string              `yaml:"prompt_secrets,omitempty" json:"prompt_secrets,omitempty"`
}
```

Invariantes:

- Regra de 1 processo ativo por `UpdateTarget`:
  - `ScriptRunner` deve cancelar (SIGTERM/SIGKILL) o processo anterior antes de iniciar outro.
- `ScriptRunner`:
  - Não pode usar `os.Exit` no core; erros fluem como eventos/retornos.
- `Confirm` e `PromptSecrets`:
  - Disparam integração com Modal Stack + Segurança JIT.

Rastreabilidade:

- E1.4.
- Story Map: E1.4.
- QA Matrix: ScriptRunner + 1 processo/update_target + no-`os.Exit`.

## 7. ModalRequest (E1.5 — Modal Stack)

Modelo mínimo para abertura de modais via `on:`.

```go
// pkg/declarative/modal.go

type ModalRequest struct {
    Type    string            `yaml:"type" json:"type"` // "confirm" | "prompt" | custom seguro
    Title   string            `yaml:"title,omitempty" json:"title,omitempty"`
    Message string            `yaml:"message,omitempty" json:"message,omitempty"`
    Fields  []ModalField      `yaml:"fields,omitempty" json:"fields,omitempty"`
}

type ModalField struct {
    Name string `yaml:"name" json:"name"`
    Type string `yaml:"type" json:"type"` // ex.: "secret"
}
```

Invariantes:

- Modais empilhados em estrutura de stack (`internal/runtime/modal`).
- Apenas o topo da pilha recebe foco/eventos.
- Usado para `confirm` e `prompt_secrets` definidos em `RunAction`.

Rastreabilidade:

- E1.5.
- Story Map: E1.5.
- QA Matrix: Modal Stack.

## 8. SecurityPolicy (E1.5 — Anti-Trojan YAML)

Contrato declarativo de segurança aplicado sobre AppConfig.

```go
// pkg/declarative/security.go

type SecurityPolicy struct {
    AllowedScripts []string          `yaml:"allowed_scripts,omitempty" json:"allowed_scripts,omitempty"`
    DenyUnknown    bool              `yaml:"deny_unknown,omitempty" json:"deny_unknown,omitempty"`
    ExtraRules     map[string]string `yaml:"extra_rules,omitempty" json:"extra_rules,omitempty"`
}
```

Invariantes:

- Interpretação normativa:
  - `DenyUnknown: true` é o default efetivo mesmo que implícito.
  - Chaves desconhecidas no YAML devem ser rejeitadas ou ignoradas com erro explícito.
- `AllowedScripts`/whitelists:
  - Usadas por QA e por implementações para validar caminhos/ações.
- Integração obrigatória com:
  - Parser v2.0 (`pkg/declarative`),
  - `ScriptRunner`,
  - Matriz QA de segurança.

Rastreabilidade:

- E1.5.
- Story Map: E1.5.
- QA Matrix: segurança JIT + anti-Trojan.

## 9. Legacy (não vinculante) — FormConfig v1.x

Apenas para referência de encapsulamento dentro de `FormComponent`.

- Modelo histórico: `FormConfig` e `Field` definidos em [`internal/config/config.go`](internal/config/config.go:1).
- Uso permitido:
  - Como backend interno do `FormComponent` ao interpretar `Component{Type:"form"}`.
- Restrições:
  - Não pode ser usado como contrato público principal.
  - Nenhum novo fluxo arquitetural deve depender diretamente de `FormConfig` fora do componente legado.

Qualquer documento ou código que trate `FormConfig` como modelo ativo deve ser atualizado para AppConfig/Component/RunAction ou marcado explicitamente como Legacy.
