package tui

import (
	"fmt"

	"shantilly/internal/config"
)

// Start é o entrypoint LEGACY v1.x usado pelo comando `shantilly form`.
// Nesta fase da migração para o runtime declarativo v2.0, o fluxo antigo
// não está mais implementado. Mantemos apenas a assinatura para preservar
// compatibilidade de build.
func Start(_ *config.FormConfig) (int, error) {
	return 1, fmt.Errorf("fluxo LEGACY v1.x (form) ainda não foi reimplementado; use o comando 'shantilly runtime' com AppConfig v2.0")
}
