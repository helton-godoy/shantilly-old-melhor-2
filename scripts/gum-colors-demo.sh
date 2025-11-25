#!/usr/bin/env bash
set -euo pipefail

# Script de demo para testar cores ANSI renderizadas via gum dentro do viewport
# preformatado do Shantilly.

if ! command -v gum >/dev/null 2>&1; then
  echo "[gum-colors-demo] Erro: o comando 'gum' não foi encontrado no PATH." >&2
  echo "Instale o gum (https://github.com/charmbracelet/gum) e tente novamente." >&2
  exit 1
fi

TITLE="Demo de cores com gum"

LINE1="INFO  - Esta linha usa uma cor de destaque suave."
LINE2="WARN  - Esta linha simula um aviso em amarelo."
LINE3="ERROR - Esta linha simula um erro em vermelho intenso."

# Título com borda e cor própria
gum style \
  --border double \
  --border-foreground "#00FF00" \
  --foreground "#00FF00" \
  --padding "1 2" \
  "${TITLE}"

echo

# Linhas coloridas individuais
gum style --foreground "#5FD7FF" "${LINE1}"
gum style --foreground "#FFD75F" "${LINE2}"
gum style --foreground "#FF5F5F" "${LINE3}"

echo

echo "[gum-colors-demo] Fim da saída colorida do gum."
