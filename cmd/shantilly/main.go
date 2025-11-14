package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"shantilly/internal/config"
	"shantilly/internal/runtime"
	"shantilly/internal/tui"
	"shantilly/internal/util"
	"shantilly/pkg/declarative"
	"github.com/spf13/cobra"
)

var (
	// rootCmd representa o comando base quando chamado sem subcomandos
	rootCmd = &cobra.Command{
		Use:   "shantilly",
		Short: "Um gerador de formulários CLI.",
		Long: `Shantilly é uma ferramenta de linha de comando para gerar formulários
interativos baseados em arquivos de definição YAML.`,
	}

	// formCmd representa o comando form (LEGACY).
	//
	// Status:
	// - Este comando é mantido apenas como entrada LEGACY v1.x baseado em FormConfig/internal/tui.
	// - Novos fluxos v2.0 DEVEM usar AppConfig + Runtime declarativo.
	// - Qualquer evolução deve migrar para o pipeline:
	//     AppConfig -> LayoutManager -> FormComponent -> ShantillyEvent -> EventManager -> on:.
	//
	// Gate:
	// - docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml
	formCmd = &cobra.Command{
		Use:   "form",
		Short: "Executa um formulário interativo (LEGACY).",
		Long: `O comando form (LEGACY) lê uma definição de formulário baseada em FormConfig
e executa o fluxo TUI v1.x. No runtime v2.0, novos fluxos devem usar AppConfig + FormComponent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				content []byte
				err     error
			)

			filePath, _ := cmd.Flags().GetString("file")

			if filePath != "" {
				content, err = os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("erro ao ler o arquivo: %w", err)
				}
			} else {
				content, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("erro ao ler da stdin: %w", err)
				}
			}

			// Parse YAML configuration
			formConfig, err := config.Parse(content)
			if err != nil {
				return fmt.Errorf("erro ao parsear configuração YAML: %w", err)
			}

			// For Story 1.3: Launch TUI application
			// E1.4 compliance: captura exit code do TUI
			exitCode, err := tui.Start(formConfig)
			if err != nil {
				return fmt.Errorf("erro ao iniciar a TUI: %w", err)
			}

			// E1.4 compliance: propagar exit code via os.Exit apenas aqui (casca CLI)
			// Referência: docs/architecture/governance-runtime-tui-v2.0.md
			// Gate: docs/qa/gates/1.x.no-osexit-core.yml
			if exitCode != 0 {
				os.Exit(exitCode)
			}

			return nil
		},
	}

	// runtimeCmd representa o comando para executar o runtime declarativo v2.0.
	//
	// Fluxo:
	// - Lê um AppConfig YAML (arquivo ou stdin).
	// - Faz parse com declarative.LoadAppConfig.
	// - Invoca internal/runtime.Start(cfg).
	runtimeCmd = &cobra.Command{
		Use:   "runtime",
		Short: "Executa o runtime TUI declarativo v2.0 a partir de um AppConfig YAML.",
		Long: `O comando runtime lê um AppConfig declarativo (v2.0) e inicia o
runtime TUI com LayoutManager, EventManager, ScriptRunner e Modal Stack.

Exemplo:

  shantilly runtime --file app.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				content []byte
				err     error
			)

			filePath, _ := cmd.Flags().GetString("file")

			if filePath != "" {
				content, err = os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("erro ao ler o arquivo: %w", err)
				}
			} else {
				content, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("erro ao ler da stdin: %w", err)
				}
			}

			cfg, err := declarative.LoadAppConfig(bytes.NewReader(content))
			if err != nil {
				return fmt.Errorf("erro ao carregar AppConfig declarativo: %w", err)
			}

			if err := runtime.Start(cfg); err != nil {
				return fmt.Errorf("erro ao executar runtime declarativo: %w", err)
			}

			return nil
		},
	}
)

func main() {
	err := rootCmd.Execute()
	util.Handle(err)
}

func init() {
	rootCmd.AddCommand(formCmd)
	formCmd.Flags().String("file", "", "O arquivo de definição de formulário a ser usado.")

	rootCmd.AddCommand(runtimeCmd)
	runtimeCmd.Flags().String("file", "", "O arquivo AppConfig declarativo v2.0 a ser usado.")
}
