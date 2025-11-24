package components

import "shantilly/internal/tui"

// HelpText é um componente para exibir texto de ajuda contextual.
type HelpText struct {
	helpTexts    map[string]string
	currentField string
	showHelp     bool
	theme        *tui.Theme
}

// NewHelpText cria um novo HelpText.
func NewHelpText(theme *tui.Theme) *HelpText {
	return &HelpText{
		helpTexts: make(map[string]string),
		showHelp:  false,
		theme:     theme,
	}
}

// SetHelpText define o texto de ajuda para um campo.
func (h *HelpText) SetHelpText(fieldKey, text string) {
	h.helpTexts[fieldKey] = text
}

// SetCurrentField define o campo atualmente focado.
func (h *HelpText) SetCurrentField(fieldKey string) {
	h.currentField = fieldKey
}

// ToggleHelp alterna a visibilidade do texto de ajuda.
func (h *HelpText) ToggleHelp() {
	h.showHelp = !h.showHelp
}

// Render renderiza o texto de ajuda para o campo atual.
func (h *HelpText) Render() string {
	if !h.showHelp {
		return ""
	}

	if h.currentField != "" {
		if help, exists := h.helpTexts[h.currentField]; exists {
			return h.theme.FieldInput.Render("Ajuda: " + help)
		}
	}

	return h.theme.FieldInput.Render("Navegue para um campo e pressione ? para ajuda.")
}

// IsHelpVisible retorna se a ajuda está visível.
func (h *HelpText) IsHelpVisible() bool {
	return h.showHelp
}

// GetHelpText retorna o texto de ajuda para uma chave de campo.
func (h *HelpText) GetHelpText(fieldKey string) (string, bool) {
	text, exists := h.helpTexts[fieldKey]
	return text, exists
}

// RenderFieldHelp renderiza a ajuda para um campo específico.
func (h *HelpText) RenderFieldHelp(fieldKey string) string {
	if text, exists := h.helpTexts[fieldKey]; exists {
		return h.theme.FieldInput.Render("Ajuda: " + text)
	}
	return ""
}
