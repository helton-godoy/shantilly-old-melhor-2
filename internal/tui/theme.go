package tui

import "github.com/charmbracelet/lipgloss"

// Theme define os estilos centralizados para a TUI.
type Theme struct {
	Base       lipgloss.Style
	FormTitle  lipgloss.Style
	FieldLabel lipgloss.Style
	FieldInput lipgloss.Style
	FieldError lipgloss.Style
	Container  lipgloss.Style
	Border     lipgloss.Style
}

// NewTheme cria um novo tema com estilos padronizados.
func NewTheme() *Theme {
	return &Theme{
		Base: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),

		FormTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Align(lipgloss.Center),

		FieldLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true),

		FieldInput: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("235")),

		FieldError: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),

		Container: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")),

		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")),
	}
}

// DefaultTheme retorna o tema padrão da aplicação.
func DefaultTheme() *Theme {
	return NewTheme()
}
