Combinado. Vamos executar o **Plano Mestre Sequencial**, começando pelo **Épico 1: Fundação do Runtime TUI**.

Estou a iniciar a **Fase 1.1: Definição dos Contratos (O Alicerce)**. Vou gerar cada artefato de código separadamente, conforme a tua solicitação.

-----

### Tarefa 1.1.1: Definir Modelos de Dados

Este artefato traduz a Seção 4 do `docs/architecture.md` para código Go. Ele define a estrutura exata do nosso YAML v2.0.

```go
// pkg/declarative/models.go
package declarative

// Config é a estrutura raiz do documento YAML v2.0.
// Fonte: docs/architecture.md (Seção 4.1)
type Config struct {
	Root LayoutNode `yaml:",inline"`      // Nó de layout raiz (column/row/box)
	On   []Logic    `yaml:"on,omitempty"` // Regras de eventos (FR8)
}

// LayoutNode modela a árvore de layout hierárquica (FR1, FR2, FR3).
// Fonte: docs/architecture.md (Seção 4.2)
type LayoutNode struct {
	Type   string       `yaml:"type"`             // "column" | "row" | "box"
	ID     string       `yaml:"id,omitempty"`     // ex: "sidebar", "content"
	Width  string       `yaml:"width,omitempty"`  // ex: "30%", "70%"
	Height int          `yaml:"height,omitempty"` // linhas fixas
	Flex   int          `yaml:"flex,omitempty"`   // proporção de espaço
	Items  []LayoutNode `yaml:"items,omitempty"`  // filhos para column/row
	// Para type: box
	Component *Component `yaml:"component,omitempty"`
}

// Component abstrai os componentes TUI declarativos (FR3–FR7).
// Fonte: docs/architecture.md (Seção 4.3)
type Component struct {
	Type   string      `yaml:"type"`             // "list" | "viewport" | "buttongroup" | "form"
	ID     string      `yaml:"id,omitempty"`     // identificador lógico do componente
	Items  []Item      `yaml:"items,omitempty"`  // list / buttongroup
	Source *Source     `yaml:"source,omitempty"` // viewport
	// Form: payload bruto delegado para parser v1.0 (huh)
	Fields  interface{} `yaml:"fields,omitempty"`
	Actions interface{} `yaml:"actions,omitempty"`
	// Conteúdo estático opcional
	Content string `yaml:"content,omitempty"`
}

// Item representa itens para 'list' e 'buttongroup'.
// Fonte: docs/architecture.md (Seção 4.4)
type Item struct {
	ID    string `yaml:"id"`
	Text  string `yaml:"text,omitempty"`  // list
	Label string `yaml:"label,omitempty"` // buttongroup
	Role  string `yaml:"role,omitempty"`  // ex: "primary", "danger"
}

// Source define a origem do conteúdo do 'viewport' (FR4).
// Fonte: docs/architecture.md (Seção 4.5)
type Source struct {
	Type        string `yaml:"type"`                   // "static" | "command"
	Content     string `yaml:"content,omitempty"`      // texto/markdown estático
	Exec        string `yaml:"exec,omitempty"`         // comando a executar
	ContentType string `yaml:"content_type,omitempty"` // ex: "markdown"
}

// Logic e RunAction modelam as regras de automação orientadas a eventos (FR8–FR11).
// Fonte: docs/architecture.md (Seção 4.6)
type Logic struct {
	Event         string     `yaml:"event"`                    // ex: "user_form:submit"
	Run           RunAction  `yaml:"run"`                      // ação obrigatória
	Confirm       bool       `yaml:"confirm,omitempty"`        // segurança JIT: confirmar antes de executar
	PromptSecrets []string   `yaml:"prompt_secrets,omitempty"` // segurança JIT: nomes de segredos a coletar via modal
}

// RunAction define a ação a ser executada.
type RunAction struct {
	Script       string   `yaml:"script,omitempty"`         // FR9: runner genérico de script
	Args         []string `yaml:"args,omitempty"`           // FR10: templates ex: {{ form.username }}
	Stdin        string   `yaml:"stdin,omitempty"`          // FR10: template serializado em JSON
	UpdateTarget string   `yaml:"update_target,omitempty"`  // FR11: id do viewport a atualizar
}
```

-----

### Tarefa 1.1.2: Definir Contratos Internos (Interfaces)

Este artefato define a interface `ShantillyComponent` (Seção 5.2 do `architecture.md`). Este é o contrato que *todos* os componentes TUI (List, Viewport, Form, etc.) devem implementar, garantindo o desacoplamento.

```go
// pkg/tui/interface.go
package tui

import tea "github.com/charmbracelet/bubbletea"

// ShantillyComponent é a interface que todos os componentes TUI declarativos
// devem implementar para serem geridos pelo LayoutManager.
// Fonte: docs/architecture.md (Seção 5.2)
type ShantillyComponent interface {
	// Init inicializa o componente e pode retornar um comando.
	Init() tea.Cmd
	// Update trata mensagens (eventos) e retorna o estado atualizado e um comando.
	Update(tea.Msg) (ShantillyComponent, tea.Cmd)
	// View renderiza a string de saída do componente.
	View() string
	// SetDimensions informa ao componente as dimensões alocadas para ele
	// pelo LayoutManager.
	SetDimensions(width, height int)
	// ID retorna o identificador único do componente.
	ID() string
}
```

-----

### Tarefa 1.1.2: Definir Contratos Internos (Eventos)

Este artefato define os eventos internos (Seção 5.3 do `architecture.md`) que os componentes e motores usam para comunicar.

```go
// pkg/tui/events.go
package tui

import tea "github.com/charmbracelet/bubbletea"

// ShantillyEvent é a mensagem padronizada que os componentes emitem
// para serem capturados pelo EventManager.
// Fonte: docs/architecture.md (Seção 5.3)
type ShantillyEvent struct {
	ComponentID string      // ID do componente que originou o evento
	Type        string      // Tipo do evento (ex: "list_select", "button_press", "form_submit")
	Payload     interface{} // Dados associados (ex: ID do item, dados do formulário)
}

// RuntimeErrorMsg é usado para notificar o TUI sobre um erro interno
// que precisa ser exibido ao usuário.
type RuntimeErrorMsg struct {
	Err error
}

// ScriptStdoutMsg é usado pelo ScriptRunner para enviar streaming de
// stdout/stderr para o ViewportComponent alvo.
type ScriptStdoutMsg struct {
	TargetID string // ID do viewport (update_target)
	Chunk    []byte // O pedaço de dados do stdout/stderr
}

// ShowModalMsg é emitida pelo EventManager (para Segurança JIT)
// e tratada pelo LayoutManager para exibir um modal.
type ShowModalMsg struct {
	Title    string
	Content  tea.Model // O modelo do modal a ser exibido (ex: confirmação, prompt de segredo)
	OnSubmit func(tea.Msg) tea.Cmd // O que fazer quando o modal for submetido
}
```

-----

### Tarefa 1.1.3: Implementar Gestão de Erros

Este é um artefato crítico da "visão do futuro" (Seção 13 do `architecture.md` e NFR8). Ele centraliza o tratamento de erros e saídas, garantindo que o `stdout` seja usado apenas para dados (JSON) e o `stderr` para erros.

```go
// internal/util/errorhandler.go
package util

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Exit codes para clareza
// Fonte: docs/architecture.md (Seção 13)
const (
	ExitOK        = 0
	ExitError     = 1 // Erro geral
	ExitCancelled = 2 // Cancelamento pelo usuário
)

// ErrAborted é um erro específico para cancelamento.
var ErrAborted = errors.New("operation aborted by user")

// Handle trata o erro final antes de sair do programa.
// Ele imprime para stderr e sai com o código apropriado.
func Handle(err error) {
	if err == nil {
		os.Exit(ExitOK) // Caminho de sucesso
	}

	exitCode := ExitError // Padrão

	// Verifica se foi um cancelamento do usuário (Ctrl+C, Esc)
	var quitMsg tea.QuitMsg
	if errors.Is(err, ErrAborted) || errors.As(err, &quitMsg) {
		exitCode = ExitCancelled
		// Mantém o stderr limpo em caso de cancelamento para
		// melhor integração com scripts.
	} else {
		// Imprime o erro real para stderr em casos de falha.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	os.Exit(exitCode)
}
```

-----

### Tarefa 1.1.3: Implementar Ponto de Entrada (CLI)

Este artefato cria o `main.go`, usando `cobra` e o `ErrorHandler` que acabámos de definir. Ele é responsável por ler o YAML de `stdin` ou `--file`, conforme a arquitetura (Seção 10).

```go
// cmd/shantilly/main.go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/helton-godoy/shantilly/internal/util" // (Assumindo o módulo go)
	"github.com/spf13/cobra"
)

var (
	// Flag para o caminho do arquivo YAML
	filePath string
)

func main() {
	// A função Handle(err) é chamada se qualquer comando retornar um erro.
	if err := newRootCmd().Execute(); err != nil {
		util.Handle(err)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "shantilly",
		Short: "Shantilly é um Runtime TUI Declarativo.",
		Long: `Shantilly consome um YAML único para definir layout,
componentes e lógica de automação TUI.`,
		// Esta é a função principal do Runtime
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				yamlBytes []byte
				err       error
			)

			// 1. Ler o YAML (Fonte: stdin ou --file)
			// Conforme Seção 5.6 do docs/architecture.md
			if filePath != "" {
				yamlBytes, err = os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("falha ao ler arquivo %s: %w", filePath, err)
				}
			} else {
				// Verifica se há algo no stdin
				stat, _ := os.Stdin.Stat()
				if (stat.Mode() & os.ModeCharDevice) == 0 {
					yamlBytes, err = io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("falha ao ler stdin: %w", err)
					}
				} else {
					// Nada no stdin e --file não foi usado
					return fmt.Errorf("nenhum YAML fornecido via stdin ou --file")
				}
			}

			if len(yamlBytes) == 0 {
				return fmt.Errorf("arquivo YAML está vazio")
			}

			// 2. Chamar o Parser (Tarefa 1.1.4) - Placeholder
			// config, err := config.Parse(yamlBytes)
			// if err != nil {
			// 	 return err
			// }

			// 3. Iniciar o Runtime (Fase 1.2) - Placeholder
			// if err := runtime.Start(config); err != nil {
			// 	 return err
			// }

			// Placeholder da Fase 1.1:
			fmt.Printf("Placeholder: YAML lido com %d bytes. O Runtime (Fase 1.2) ainda não foi implementado.\n", len(yamlBytes))

			return nil
		},
	}

	// Adiciona a flag --file
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "Caminho para o arquivo de configuração YAML")

	return rootCmd
}
```

-----

### Tarefa 1.1.4: Implementar o Parser YAML

Finalmente, criamos o *parser* que conecta o `main.go` (Tarefa 1.1.3) aos nossos modelos de dados (Tarefa 1.1.1).

```go
// internal/config/parser.go
package config

import (
	"fmt"

	"github.com/helton-godoy/shantilly/pkg/declarative" // (Assumindo o módulo go)
	"gopkg.in/yaml.v3"
)

// Parse converte os bytes YAML no struct de Configuração declarativa v2.0.
// Fonte: docs/architecture.md (Seção 5.6)
func Parse(yamlBytes []byte) (*declarative.Config, error) {
	var config declarative.Config

	// Usar yaml.v3 para parsing
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, fmt.Errorf("falha no parsing do YAML: %w", err)
	}

	// Validação básica (Otimização de Risco)
	// Garante que a estrutura raiz foi definida.
	if config.Root.Type == "" {
		return nil, fmt.Errorf("validação falhou: o nó raiz 'type' (column/row/box) é obrigatório no YAML")
	}

	// TODO: Adicionar validação mais profunda aqui (ex: IDs únicos, tipos de componentes válidos)
	// (Isso será feito em uma fase posterior de "hardening")

	return &config, nil
}
```