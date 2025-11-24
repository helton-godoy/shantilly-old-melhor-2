package components

import (
	"fmt"
	"strings"

	"shantilly/internal/tui"
)

// ProgressIndicator é um componente para exibir o progresso do preenchimento do formulário.
type ProgressIndicator struct {
	completed int
	total     int
	theme     *tui.Theme
}

// NewProgressIndicator cria um novo ProgressIndicator.
func NewProgressIndicator(theme *tui.Theme) *ProgressIndicator {
	return &ProgressIndicator{
		theme: theme,
	}
}

// UpdateProgress atualiza o estado de conclusão do progresso.
func (p *ProgressIndicator) UpdateProgress(completed, total int) {
	p.completed = completed
	p.total = total
}

// Render renderiza a barra de progresso e a porcentagem.
func (p *ProgressIndicator) Render() string {
	percentage := p.GetProgress()
	progressText := fmt.Sprintf("Progresso: %d/%d (%.0f%%)", p.completed, p.total, percentage)

	barWidth := 30
	filledWidth := int(float64(barWidth) * (percentage / 100))
	bar := strings.Repeat("=", filledWidth) + strings.Repeat("-", barWidth-filledWidth)

	return fmt.Sprintf("%s\n[%s]", progressText, bar)
}

// GetProgress retorna a porcentagem de progresso.
func (p *ProgressIndicator) GetProgress() float64 {
	if p.total == 0 {
		return 100.0
	}
	return (float64(p.completed) / float64(p.total)) * 100
}

// IsComplete retorna verdadeiro se o progresso for 100%.
func (p *ProgressIndicator) IsComplete() bool {
	return p.completed == p.total
}

// RenderCompact renderiza uma versão compacta do progresso.
func (p *ProgressIndicator) RenderCompact() string {
	return fmt.Sprintf("%d/%d (%.0f%%)", p.completed, p.total, p.GetProgress())
}
