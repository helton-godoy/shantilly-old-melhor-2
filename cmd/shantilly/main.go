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

	// formCmd representa o comando form
	formCmd = &cobra.Command{
		Use:   "form",
		Short: "Executa um formulário interativo.",
		Long: `O comando form lê uma definição de formulário, seja de um arquivo
especificado com o sinalizador --file ou da entrada padrão (stdin), e então
renderiza e gerencia o formulário interativo no terminal.`,
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
			err = tui.Start(formConfig)
			if err != nil {
				return fmt.Errorf("erro ao iniciar a TUI: %w", err)
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
}
