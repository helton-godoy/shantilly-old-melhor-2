package components

import (
	"fmt"
	"strings"
)

// ProgressIndicator é um componente para exibir o progresso do preenchimento do formulário.
type ProgressIndicator struct {
	completed   int
	total       int
	initialized bool
	theme       *Theme
}

// NewProgressIndicator cria um novo ProgressIndicator.
func NewProgressIndicator(theme *Theme) *ProgressIndicator {
	return &ProgressIndicator{
		theme: theme,
	}
}

// UpdateProgress atualiza o estado de conclusão do progresso.
func (p *ProgressIndicator) UpdateProgress(completed, total int) {
	p.completed = completed
	p.total = total
	p.initialized = true
}

// Render renderiza a barra de progresso e a porcentagem.
func (p *ProgressIndicator) Render() string {
	if !p.initialized || p.total == 0 {
		return ""
	}

	percentage := p.GetProgress()
	progressText := fmt.Sprintf("Progresso: %d/%d (%.1f%%)", p.completed, p.total, percentage)

	// Simple text-based progress bar
	barWidth := 30
	filledWidth := int(float64(barWidth) * (percentage / 100))
	bar := strings.Repeat("=", filledWidth) + strings.Repeat("-", barWidth-filledWidth)

	return fmt.Sprintf("%s\n[%s]", progressText, bar)
}

// GetProgress retorna a porcentagem de progresso.
func (p *ProgressIndicator) GetProgress() float64 {
	if p.total == 0 {
		if !p.initialized {
			return 0.0
		}
		return 100.0
	}
	return (float64(p.completed) / float64(p.total)) * 100
}

// IsComplete retorna verdadeiro se o progresso for 100%.
func (p *ProgressIndicator) IsComplete() bool {
	if !p.initialized {
		return false
	}
	if p.total == 0 {
		return true
	}
	return p.completed == p.total
}

// RenderCompact renderiza uma versão compacta do progresso.
func (p *ProgressIndicator) RenderCompact() string {
	percentage := p.GetProgress()
	return fmt.Sprintf("%.0f%% (%d/%d)", percentage, p.completed, p.total)
}
