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
	"shantilly/internal/wizard"
	"shantilly/pkg/declarative"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	version = "dev"

	// rootCmd representa o comando base quando chamado sem subcomandos
	rootCmd = &cobra.Command{
		Use:   "shantilly",
		Short: "Runtime TUI declarativo para formulários e automações em YAML.",
		Long: `Shantilly é uma ferramenta de linha de comando para executar interfaces TUI declarativas
e formulários interativos baseados em arquivos YAML.

Comandos principais:
  shantilly form --file form.yaml      # Executa formulários LEGACY baseados em FormConfig
  shantilly runtime --file app.yaml    # Executa o runtime TUI declarativo v2.0 (AppConfig)
  shantilly validate --file app.yaml   # Apenas valida a AppConfig declarativa v2.0`,
		Version: version,
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
		Short: "Executa um formulário interativo legado baseado em FormConfig.",
		Long: `O comando form (LEGACY) lê uma definição de formulário baseada em FormConfig
e executa o fluxo TUI v1.x.

Exemplos:
  shantilly form --file form.yaml
  cat form.yaml | shantilly form

No runtime v2.0, novos fluxos devem usar AppConfig + Runtime declarativo.`,
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

Exemplos:
  shantilly runtime --file app.yaml
  cat app.yaml | shantilly runtime`,
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

	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Valida um AppConfig declarativo v2.0 (sem executar o runtime).",
		Long: `O comando validate lê um AppConfig declarativo (v2.0) e executa apenas
as validações estruturais, sem iniciar o runtime TUI.

Exemplos:
  shantilly validate --file app.yaml
  cat app.yaml | shantilly validate`,
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

			if _, err := declarative.LoadAppConfig(bytes.NewReader(content)); err != nil {
				return fmt.Errorf("erro ao validar AppConfig declarativo: %w", err)
			}

			fmt.Fprintln(os.Stdout, "AppConfig declarativo válido.")
			return nil
		},
	}

	wizardCmd = &cobra.Command{
		Use:   "wizard",
		Short: "Abre um wizard TUI para gerar um AppConfig declarativo v2.0.",
		Long: `O comando wizard executa um assistente interativo em modo TUI
para gerar um AppConfig declarativo v2.0 (app.yaml) com layout e componentes
básicos, pronto para uso com o comando 'shantilly runtime'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputPath, _ := cmd.Flags().GetString("output")

			cfg, err := wizard.Run()
			if err != nil {
				return err
			}

			data, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("erro ao serializar AppConfig gerado: %w", err)
			}

			if outputPath == "" || outputPath == "-" {
				// Escreve em stdout por padrão.
				if _, err := os.Stdout.Write(data); err != nil {
					return fmt.Errorf("erro ao escrever AppConfig em stdout: %w", err)
				}
				return nil
			}

			if err := os.WriteFile(outputPath, data, 0o644); err != nil {
				return fmt.Errorf("erro ao escrever AppConfig em '%s': %w", outputPath, err)
			}

			fmt.Fprintf(os.Stdout, "AppConfig gerado com sucesso em %s\n", outputPath)
			return nil
		},
	}
)

func main() {
	err := rootCmd.Execute()
	exitCode := util.Handle(err)
	if exitCode != util.ExitSuccess {
		os.Exit(exitCode)
	}
}

func init() {
	rootCmd.AddCommand(formCmd)
	formCmd.Flags().String("file", "", "O arquivo de definição de formulário a ser usado.")

	rootCmd.AddCommand(runtimeCmd)
	runtimeCmd.Flags().String("file", "", "O arquivo AppConfig declarativo v2.0 a ser usado.")

	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().String("file", "", "O arquivo AppConfig declarativo v2.0 a ser validado.")

	rootCmd.AddCommand(wizardCmd)
	wizardCmd.Flags().String("output", "app.yaml", "Caminho do arquivo YAML a ser gerado (use '-' para stdout).")
}
