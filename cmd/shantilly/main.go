package main

import (
	"fmt"
	"io"
	"os"

	"shantilly/internal/config"
	"shantilly/internal/tui"
	"shantilly/internal/util"

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
)

func main() {
	err := rootCmd.Execute()
	code := util.Handle(err)
	if code != util.ExitSuccess {
		os.Exit(code)
	}
}

func init() {
	rootCmd.AddCommand(formCmd)
	formCmd.Flags().String("file", "", "O arquivo de definição de formulário a ser usado.")
}
