Executado. Concluímos a **Fase 1.1 (Contratos e Alicerce)**.

Continuando rigorosamente o **Plano Mestre Sequencial**, iniciamos agora a **Fase 1.2: Motor de Layout (Renderização e Foco)**.

Nesta fase, vamos:

1.  Criar o `LayoutManager` (o "Gestor Duplo").
2.  Implementar a lógica de renderização recursiva (`column`/`row`/`box`).
3.  Implementar um Tema `lipgloss` centralizado.
4.  Conectar o `main.go` (da Fase 1.1) ao Parser e ao novo Runtime.

Abaixo estão os artefatos desta fase.

-----

### Tarefa 1.2.1: Implementar o Tema Central (Obrigatório para Renderização)

Criamos um tema `lipgloss` centralizado (conforme Seção 5.1 do `architecture.md`) para garantir consistência visual e evitar que a IA use estilos arbitrários (o que nos protege da "baixa qualidade de código" que mencionaste).

```go
// internal/tui/theme.go
package tui

import "github.com/charmbracelet/lipgloss"

// Theme define os estilos centralizados para o runtime TUI.
// Fonte: docs/architecture.md (Seção 5.1)
type Theme struct {
	// Estilo base para um 'box'
	Box lipgloss.Style
	// Estilo para um 'box' quando está focado
	BoxFocused lipgloss.Style
	// Estilo para erros
	ErrorText lipgloss.Style
}

// NewDefaultTheme cria o tema padrão do Shantilly.
func NewDefaultTheme() *Theme {
	// Cores
	primaryColor := lipgloss.Color("63")  // Roxo/Magenta
	dangerColor := lipgloss.Color("196") // Vermelho
	borderColor := lipgloss.Color("240") // Cinza

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	return &Theme{
		Box: boxStyle,
		BoxFocused: boxStyle.Copy().
			Border(lipgloss.ThickBorder()).
			BorderForeground(primaryColor),
		ErrorText: lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true),
	}
}
```

-----

### Tarefa 1.2.1: Implementar o Ponto de Entrada do Runtime

Este ficheiro `runtime.go` serve como a ponte limpa entre o `main.go` (CLI) e o `LayoutManager` (TUI), iniciando o programa `bubbletea`.

```go
// internal/runtime/runtime.go
package runtime

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/runtime/layout"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
)

// Start inicializa e executa o TUI Runtime com a configuração fornecida.
func Start(config *declarative.Config) error {
	// 1. Criar o tema central (Tarefa 1.2.1)
	theme := tui.NewDefaultTheme()

	// 2. Criar o LayoutManager (Tarefa 1.2.1)
	// Esta é a implementação do "Gestor Duplo" (Padrão 3 da arquitetura)
	// Ele recebe a *raiz* da árvore de layout.
	layoutManager := layout.NewManager(&config.Root, theme)

	// 3. Iniciar o programa Bubbletea
	p := tea.NewProgram(
		layoutManager,
		tea.WithAltScreen(), // Usar AltScreen para TUI
		// tea.WithMouseCellMotion(), // (Fora do escopo do MVP v1.0)
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("erro durante a execução do runtime TUI: %w", err)
	}

	// (O `ErrorHandler` no main.go tratará saídas/cancelamentos)
	return nil
}
```

-----

### Tarefa 1.2.1 e 1.2.2: Implementar o LayoutManager (O "Gestor Duplo")

Este é o **coração** da Fase 1.2. Ele implementa o `bubbletea.Model` raiz.

  * ***Otimização de Risco (Visão do Futuro):*** Em linha com a nossa discussão e a Seção 5.4 do `architecture.md`, este modelo gere **APENAS** Layout (`WindowSizeMsg`) e Foco (`focusedComponent`). Ele **NÃO** contém lógica de automação (que virá na Fase 1.4).

<!-- end list -->

```go
// internal/runtime/layout/manager.go
package layout

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/helton-godoy/shantilly/internal/tui"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/helton-godoy/shantilly/pkg/tui/events" // (Corrigido para pkg/tui/events)
)

// Manager é o modelo Bubbletea raiz, implementando o "Gestor Duplo"
// (Layout e Foco) conforme Padrão 3 da arquitetura v2.0.
// Fonte: docs/architecture.md (Seção 5.4)
type Manager struct {
	theme *tui.Theme
	root  *declarative.LayoutNode // A raiz da árvore de layout (FR1)

	// Estado do Layout
	width  int
	height int

	// Estado do Foco
	// (Placeholder - Fase 1.3 irá preencher isto com componentes reais)
	// components []tui.ShantillyComponent
	focusedIdx int

	// Estado de Erro
	lastErr error
}

func NewManager(root *declarative.LayoutNode, theme *tui.Theme) *Manager {
	return &Manager{
		theme: theme,
		root:  root,
		// (Inicialmente, o foco está no primeiro componente)
		focusedIdx: 0,
	}
}

func (m *Manager) Init() tea.Cmd {
	// (Na Fase 1.3, isto irá inicializar os componentes filhos)
	return nil
}

// Update trata o ciclo de vida do TUI.
// Foco principal: Gestão de Layout (NFR2) e Foco.
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// NFR2: Layout Fluido
	// Fonte: docs/prd.md (NFR2)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// (Na Fase 1.3, iremos notificar os filhos: m.notifyChildrenOfResize())

	// Gestão de Erros
	case events.RuntimeErrorMsg:
		m.lastErr = msg.Err
		return m, nil // Exibe o erro na View

	// Sair
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		
		// Gestão de Foco (Placeholder)
		// Fonte: docs/prd.md (Meta de UI: Gestão de Foco Global)
		case "tab":
			m.focusedIdx = (m.focusedIdx + 1) % 2 // (Simplesmente alterna 0 e 1 por agora)
		case "shift+tab":
			m.focusedIdx = (m.focusedIdx - 1 + 2) % 2 // (Simplesmente alterna 0 e 1 por agora)
		}
	}

	// (Na Fase 1.3, encaminharemos a msg para o componente focado)
	// if len(m.components) > 0 {
	// 	 focusedCmp, cmd := m.components[m.focusedIdx].Update(msg)
	// 	 m.components[m.focusedIdx] = focusedCmp
	// 	 cmds = append(cmds, cmd)
	// }

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View é onde a renderização recursiva (Tarefa 1.2.2) acontece.
func (m *Manager) View() string {
	if m.width == 0 || m.height == 0 {
		return "Inicializando..."
	}

	// 1. Renderizar a árvore de layout principal
	// Esta é a chamada para a Tarefa 1.2.2 (render.go)
	layoutView := m.renderNode(m.root, m.width, m.height)

	// 2. Renderizar Erros (se houver)
	if m.lastErr != nil {
		errorView := m.theme.ErrorText.Render("Erro: " + m.lastErr.Error())
		// (Na Fase 1.6, isto será um Modal sobreposto)
		return errorView
	}

	// (Na Fase 1.6, isto incluirá a Pilha Modal)
	return layoutView
}
```

-----

### Tarefa 1.2.2: Implementar o Renderizador Recursivo

Este é o motor de renderização `lipgloss` (FR1, FR2). Separei-o do `manager.go` para manter a lógica limpa, conforme o plano.

```go
// internal/runtime/layout/render.go
package layout

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/helton-godoy/shantilly/pkg/declarative"
	"github.com/muesli/reflow/truncate"
)

// renderNode é a função recursiva da Tarefa 1.2.2.
// Ela caminha pela árvore LayoutNode e aplica as regras de Lipgloss.
// Fonte: docs/prd.md (FR1, FR2)
func (m *Manager) renderNode(node *declarative.LayoutNode, w, h int) string {
	if node == nil {
		return ""
	}

	switch node.Type {
	case "column":
		// (Lógica de divisão de altura para 'column' - complexa,
		// será implementada na Fase 1.3 quando houver componentes reais)
		
		// Placeholder para Fase 1.2: Apenas junta verticalmente
		var children []string
		for _, item := range node.Items {
			// (Aqui entra a lógica de 'flex' e 'height')
			// Placeholder: divide igualmente
			childH := h / len(node.Items) 
			children = append(children, m.renderNode(&item, w, childH))
		}
		return lipgloss.JoinVertical(lipgloss.Left, children...)

	case "row":
		// (Lógica de divisão de largura para 'row' - complexa)
		
		// Placeholder para Fase 1.2: Apenas junta horizontalmente
		var children []string
		for _, item := range node.Items {
			// (Aqui entra a lógica de 'flex' e 'width')
			// Placeholder: divide igualmente
			childW := w / len(node.Items)
			children = append(children, m.renderNode(&item, childW, h))
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, children...)

	case "box":
		// Esta é a folha da árvore de renderização.
		
		// 1. Determinar o estilo (focado ou não)
		style := m.theme.Box
		// (Placeholder de Foco: assumindo que o ID "sidebar" é o foco 0)
		if node.ID == "sidebar" && m.focusedIdx == 0 {
			style = m.theme.BoxFocused
		}
		if node.ID == "content" && m.focusedIdx == 1 {
			style = m.theme.BoxFocused
		}

		// 2. Aplicar o estilo às dimensões dadas (w, h)
		// Lipgloss trata 'width' e 'height' no render.
		// Subtrai 2 para a borda
		innerW := w - 2
		innerH := h - 2 
		if innerW < 0 { innerW = 0 }
		if innerH < 0 { innerH = 0 }

		style = style.Width(w).Height(h)

		// 3. Renderizar o conteúdo (Componente)
		var content string
		if node.Component != nil {
			// (Na Fase 1.3, chamaremos o 'View()' do componente real)
			// content = m.getComponent(node.Component.ID).View()
			
			// Placeholder para Fase 1.2: (Estória 1.1 AC4)
			content = fmt.Sprintf("Placeholder para Componente:\nType: %s\nID: %s", 
				node.Component.Type, node.Component.ID)
		} else {
			content = fmt.Sprintf("Box (ID: %s)", node.ID)
		}

		// Truncar conteúdo se for maior que o espaço interno
		truncatedContent := truncate.StringWithTail(content, uint(innerW * innerH), "...")
		
		// (Lipgloss precisa de altura interna exata para preenchimento)
		return style.Render(lipgloss.NewStyle().Width(innerW).Height(innerH).Render(truncatedContent))

	default:
		// Erro de parsing (deveria ter sido pego na Fase 1.1.4)
		errMsg := m.theme.ErrorText.Render(fmt.Sprintf("Erro: Tipo de layout desconhecido: '%s'", node.Type))
		return m.theme.Box.Render(errMsg)
	}
}
```

-----

### Tarefa 1.2.1: Atualizar o Ponto de Entrada (CLI)

Finalmente, **modificamos** o `cmd/shantilly/main.go` (criado na Fase 1.1) para *usar* os artefatos da Fase 1.1 e 1.2. Removemos os placeholders.

```go
// cmd/shantilly/main.go (MODIFICADO)
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/helton-godoy/shantilly/internal/config" // <-- ADICIONADO
	"github.com/helton-godoy/shantilly/internal/runtime" // <-- ADICIONADO
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

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "shantilly",
		Short: "Shantilly é um Runtime TUI Declarativo.",
		Long: `Shantilly consome um YAML único para definir layout,
componentes e lógica de automação TUI.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				yamlBytes []byte
				err       error
			)

			// 1. Ler o YAML (Fonte: stdin ou --file)
			if filePath != "" {
				yamlBytes, err = os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("falha ao ler arquivo %s: %w", filePath, err)
				}
			} else {
				stat, _ := os.Stdin.Stat()
				if (stat.Mode() & os.ModeCharDevice) == 0 {
					yamlBytes, err = io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("falha ao ler stdin: %w", err)
					}
				} else {
					return fmt.Errorf("nenhum YAML fornecido via stdin ou --file")
				}
			}

			if len(yamlBytes) == 0 {
				return fmt.Errorf("arquivo YAML está vazio")
			}

			// 2. Chamar o Parser (Tarefa 1.1.4) - AGORA REAL
			config, err := config.Parse(yamlBytes)
			if err != nil {
				// Erros de parsing são tratados pelo Handle
				return err
			}

			// 3. Iniciar o Runtime (Fase 1.2) - AGORA REAL
			if err := runtime.Start(config); err != nil {
				// Erros do Bubbletea (incluindo Ctrl+C) são tratados pelo Handle
				return err
			}

			return nil
		},
	}

	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "Caminho para o arquivo de configuração YAML")

	return rootCmd
}
```