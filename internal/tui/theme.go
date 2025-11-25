package tui

import "github.com/charmbracelet/lipgloss"

// Theme define os estilos para os componentes da TUI.
type Theme struct {
	Base       lipgloss.Style
	FormTitle  lipgloss.Style
	FieldLabel lipgloss.Style
	FieldInput lipgloss.Style
	FieldError lipgloss.Style
	Container  lipgloss.Style
	Border     lipgloss.Style
}

// DefaultTheme cria e retorna um tema padrão para o formulário.
func DefaultTheme() *Theme {
	// Estilos base
	baseStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Foreground(lipgloss.Color("#FFF"))

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Italic(true)

	// Definir o tema
	return &Theme{
		Base:       baseStyle,
		FormTitle:  baseStyle.Copy().Bold(true).Foreground(lipgloss.Color("#00BFFF")),
		FieldLabel: baseStyle.Copy().Bold(true),
		FieldInput: baseStyle.Copy(),
		FieldError: errorStyle,
		Container:  baseStyle.Copy().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#00BFFF")),
		Border:     lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")),
	}
}

// NewDefaultTheme mantém a compatibilidade com chamadas anteriores.
func NewDefaultTheme() *Theme {
	return DefaultTheme()
}
