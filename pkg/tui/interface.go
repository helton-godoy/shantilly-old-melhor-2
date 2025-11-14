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
