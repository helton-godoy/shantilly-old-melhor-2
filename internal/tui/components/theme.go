package components

import "github.com/charmbracelet/lipgloss"

// Theme define os estilos para os componentes da TUI.
// Esta struct pode ser estendida para incluir estilos para outros componentes.
type Theme struct {
	FieldError lipgloss.Style
	FieldInput lipgloss.Style
	FieldLabel lipgloss.Style
}
