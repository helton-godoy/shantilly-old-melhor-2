Aqui está a árvore completa de todos os artefatos **essenciais** que foram gerados ou reapresentados durante esta sessão de trabalho.

Esta estrutura reflete o estado final do **Épico 1 ao 4**, integrando a fundação do Runtime TUI, os componentes, a lógica de automação (incluindo Ansible e Segredos JIT) e o suporte SSH.

Você deve garantir que estes ficheiros existam no seu repositório local com o conteúdo que definimos.

```plaintext
shantilly/
├── README.md                                       # Guia para Agentes e Regras de Qualidade
├── cmd/
│   └── shantilly/
│       └── main.go                                 # Ponto de entrada CLI (com 'run' e 'serve-ssh')
├── docs/
│   ├── architecture/
│   │   └── high-level-architecture.md              # Diagramas de fluxo e separação de responsabilidades
│   ├── architecture.md                             # A Fonte da Verdade Técnica (Modelos, Contratos, Padrões)
│   ├── prd.md                                      # O Escopo do Produto (Requisitos e Épicos)
│   └── project-brief.md                            # Visão Estratégica e Definição do Problema
├── internal/
│   ├── config/
│   │   └── parser.go                               # Parser YAML v2.0
│   ├── components/
│   │   ├── buttongroup/
│   │   │   └── model.go                            # Componente de Botões
│   │   ├── form/
│   │   │   ├── legacy_v1/                          # Código v1.0 Encapsulado
│   │   │   │   ├── config.go
│   │   │   │   ├── model.go                        # TUI v1.0 refatorado (sem os.Exit)
│   │   │   │   └── parser.go
│   │   │   └── wrapper.go                          # Adaptador v2.0 (ShantillyComponent)
│   │   ├── inventory_explorer/
│   │   │   └── model.go                            # Componente Preditivo (ansible-inventory)
│   │   ├── list/
│   │   │   └── model.go                            # Componente de Menu
│   │   ├── playbook_explorer/
│   │   │   └── model.go                            # Componente Preditivo (Filtro de Arquivos)
│   │   └── viewport/
│   │       └── model.go                            # Componente de Visualização/Logs
│   ├── runtime/
│   │   ├── event/
│   │   │   └── manager.go                          # Motor de Lógica (Processamento 'on:' e JIT)
│   │   ├── layout/
│   │   │   ├── manager.go                          # Gestor Duplo (Layout, Foco, Delegação)
│   │   │   └── render.go                           # Motor de Renderização Recursiva (Lipgloss)
│   │   ├── modal/
│   │   │   ├── confirm.go                          # Modal de Confirmação
│   │   │   └── secret.go                           # Modal de Segredos (Vault Pass)
│   │   ├── runner/
│   │   │   ├── ansible.go                          # Runner Especialista Ansible
│   │   │   ├── runner.go                           # Dispatcher e Gestor de Processos (Streaming)
│   │   │   └── templating.go                       # Motor de Templates e Estado
│   │   ├── main_model.go                           # Modelo Raiz (Overlay + Layout)
│   │   └── runtime.go                              # Inicialização do Bubbletea
│   ├── ssh/
│   │   └── server.go                               # Servidor Wish e Handler TUI
│   ├── tui/
│   │   └── theme.go                                # Tema Centralizado
│   └── util/
│       └── errorhandler.go                         # Gestão Centralizada de Erros
└── pkg/
    ├── declarative/
    │   └── models.go                               # Structs do YAML (Config, LayoutNode, etc.)
    └── tui/
        ├── events.go                               # Mensagens Internas (ShantillyEvent, etc.)
        └── interface.go                            # Interface ShantillyComponent
```

**Próximos Passos Sugeridos:**

1.  Copie os conteúdos para os respectivos ficheiros.
2.  Execute `go mod tidy` para baixar as dependências (`bubbletea`, `lipgloss`, `huh`, `wish`, etc.).
3.  Execute o **Portão de Qualidade**: `golangci-lint run ./...`.
4.  Compile e teste: `go build -o shantilly ./cmd/shantilly` e depois `./shantilly --file examples/basic_form.yaml`.