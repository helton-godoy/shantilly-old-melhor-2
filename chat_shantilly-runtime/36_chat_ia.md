Executado. Concluímos os Épicos 1 (Fundação), 2 (Runner Ansible) e 3 (Componentes Preditivos).

Estamos na reta final. Seguindo o **Plano Mestre Sequencial**, iniciamos o último épico: **Épico 4: Administração SSH (Visão de Longo Prazo)**.

**Objetivo:** Integrar o `charmbracelet/wish` para servir todo o nosso Runtime TUI (construído nos Épicos 1-3) diretamente sobre uma sessão SSH.

  * ***Otimização de Risco (Visão do Futuro):*** A maior armadilha aqui é como ligar o I/O (Entrada/Saída) do `bubbletea` à sessão SSH. Vamos implementar a solução (definida na arquitetura) de usar `tea.WithInput(ssh.Session)` e `tea.WithOutput(ssh.Session)`.

-----

### Fase 4.1: Implementação do Servidor SSH

Vamos criar o servidor `wish` e o novo subcomando `serve-ssh`.

**Tarefa 4.1.1: Implementar o Servidor SSH (com Handler TUI)**

Criamos o novo pacote `internal/ssh/server.go`. Este ficheiro contém o servidor `wish` e o *handler* (middleware) que lança o nosso Runtime TUI (Épicos 1-3) para cada nova sessão SSH.

```go
// internal/ssh/server.go
package ssh

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/helton-godoy/shantilly/internal/config"
	"github.com/helton-godoy/shantilly/internal/runtime"
	"github.com/helton-godoy/shantilly/internal/runtime/layout"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/muesli/termenv"
	"golang.org/x/sync/errgroup"
)

// StartServer inicia o servidor SSH (Épico 4)
// Fonte: docs/prd.md (Épico 4), docs/architecture.md (Seção 4.1.2)
func StartServer(host string, port int, yamlPath string) error {
	// 1. Criar o servidor Wish
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/shantilly_ssh_key"), // (Chave de Host)
		wish.WithMiddleware(
			// 2. Adicionar o Middleware BubbleTea (Handler TUI)
			// Fonte: docs/architecture.md (Seção 4.2.1)
			TUIHandler(yamlPath), // <-- O nosso handler (Tarefa 4.2.1)
			lm.Middleware(),      // Middleware de logging
		),
	)
	if err != nil {
		return fmt.Errorf("falha ao criar servidor SSH: %w", err)
	}

	// 3. Iniciar o servidor
	var eg errgroup.Group
	eg.Go(func() error {
		fmt.Printf("Iniciando servidor SSH em %s:%d (YAML: %s)\n", host, port, yamlPath)
		if err := s.ListenAndServe(); err != nil {
			return fmt.Errorf("falha ao servir SSH: %w", err)
		}
		return nil
	})

	// (Graceful shutdown - placeholder)
	// ...

	return eg.Wait()
}

// TUIHandler é o Middleware Wish que lança o Runtime TUI v2.0
// para cada sessão SSH.
// Fonte: docs/architecture.md (Seção 4.2.1)
func TUIHandler(yamlPath string) wish.Middleware {
	return func(h wish.Handler) wish.Handler {
		return func(s wish.Session) {
			// --- Início da Lógica de Integração do Runtime TUI ---
			// (Esta lógica é uma réplica do 'main.go', mas
			// direcionada para a sessão 's')
			
			// 1. Carregar o YAML (Tarefa 4.2.2 - Hardcoded Path)
			yamlBytes, err := os.ReadFile(yamlPath)
			if err != nil {
				wish.Error(s, fmt.Errorf("falha ao ler YAML do servidor: %w", err))
				h(s) // (Falha, volta para o handler padrão)
				return
			}
			
			// 2. Parser (Fase 1.1.4)
			appConfig, err := config.Parse(yamlBytes)
			if err != nil {
				wish.Error(s, fmt.Errorf("falha no parsing do YAML: %w", err))
				h(s)
				return
			}

			// 3. Criar o Runtime TUI (Fases 1.2 - 1.6)
			theme := tui.NewDefaultTheme()
			layoutManager := layout.NewManager(&appConfig.Root, theme, appConfig.On)
			mainModel := runtime.NewMainModel(layoutManager)
			
			// 4. Criar o Programa Bubbletea (Otimização de Risco)
			p := tea.NewProgram(
				mainModel,
				// *** OTIMIZAÇÃO CRÍTICA (Fase 4.1) ***
				// Direciona Input/Output para a Sessão SSH,
				// não para o os.Stdin/os.Stdout do *servidor*.
				// Fonte: docs/architecture.md (Seção 4.2.1)
				tea.WithInput(s),
				tea.WithOutput(s),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(), // (Habilitado para SSH)
			)

			// 5. Executar o programa TUI
			// (bm.Middleware usa o 'termenv' para cores corretas)
			ctx, cancel := context.WithTimeout(s.Context(), 1*time.Hour) // (Timeout da sessão)
			defer cancel()
			
			if _, err := p.RunContext(ctx); err != nil {
				wish.Error(s, fmt.Errorf("erro no runtime TUI: %w", err))
			}
			
			// 6. Fechar a sessão
			h(s)
			// --- Fim da Lógica de Integração do Runtime TUI ---
		}
	}
}
```

-----

### Tarefa 4.1.2: Modificar Ponto de Entrada (CLI)

**Modificamos** o `cmd/shantilly/main.go` (da Fase 1.2.1) para adicionar o novo subcomando `serve-ssh` (Épico 4).

```go
// cmd/shantilly/main.go (MODIFICADO)
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/helton-godoy/shantilly/internal/config"
	"github.com/helton-godoy/shantilly/internal/runtime"
	"github.com/helton-godoy/shantilly/internal/ssh" // <-- ADICIONADO
	"github.com/helton-godoy/shantilly/internal/util"
	"github.com/spf13/cobra"
)

var (
	filePath string
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		util.Handle(err)
	}
}

// newRootCmd (MODIFICADO)
// Adiciona o subcomando 'serve-ssh'
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "shantilly",
		Short: "Shantilly é um Runtime TUI Declarativo.",
		Long: `Shantilly consome um YAML único para definir layout,
componentes e lógica de automação TUI.
		
Execute sem subcomando (ou com 'run') para iniciar o TUI localmente,
ou use 'serve-ssh' para iniciar o servidor SSH.`,
		// Renomeado para 'run' (mas 'rootCmd' ainda executa)
		RunE: runTUI,
	}

	// Flag --file (movida para persistente)
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Caminho para o arquivo de configuração YAML")
	rootCmd.MarkPersistentFlagRequired("file") // (Obrigatório para ambos os modos)

	// --- ADICIONADO (Épico 4) ---
	rootCmd.AddCommand(newServeCmd())
	// ---------------------------

	return rootCmd
}

// runTUI (Lógica extraída da Fase 1.2.1)
// Executa o TUI localmente.
func runTUI(cmd *cobra.Command, args []string) error {
	var (
		yamlBytes []byte
		err       error
	)

	// 1. Ler o YAML
	if filePath != "" {
		yamlBytes, err = os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("falha ao ler arquivo %s: %w", filePath, err)
		}
	} else {
		// (A flag 'file' agora é obrigatória,
		// mas mantemos a lógica do stdin se a flag for removida)
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			yamlBytes, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("falha ao ler stdin: %w", err)
			}
		} else {
			return fmt.Errorf("nenhum YAML fornecido via --file (ou stdin)")
		}
	}

	if len(yamlBytes) == 0 {
		return fmt.Errorf("arquivo YAML está vazio")
	}

	// 2. Parser (Fase 1.1.4)
	config, err := config.Parse(yamlBytes)
	if err != nil {
		return err
	}

	// 3. Iniciar o Runtime (Fase 1.6.5)
	if err := runtime.Start(config); err != nil {
		return err
	}

	return nil
}

// newServeCmd (ADICIONADO)
// Define o comando 'serve-ssh' (Épico 4)
// Fonte: docs/prd.md (Épico 4)
func newServeCmd() *cobra.Command {
	var (
		host string
		port int
	)

	serveCmd := &cobra.Command{
		Use:   "serve-ssh",
		Short: "Inicia o servidor SSH para o Runtime TUI",
		Long:  `Inicia um servidor SSH que lança o Runtime TUI (definido por --file) para cada conexão.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// (O 'filePath' é lido da flag persistente do root)
			if filePath == "" {
				return fmt.Errorf("--file é obrigatório para o modo servidor")
			}

			// Chama o servidor SSH (Tarefa 4.1.1)
			return ssh.StartServer(host, port, filePath)
		},
	}

	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "Host para o servidor SSH")
	serveCmd.Flags().IntVarP(&port, "port", "p", 2222, "Porta para o servidor SSH")

	return serveCmd
}
```