#!/usr/bin/env bash
set -euo pipefail

# Este script usa Charmbracelet/gum para desenhar um box com borda no terminal.
# O objetivo é observar como essas bordas se comportam dentro de um viewport
# do runtime declarativo do Shantilly.

if ! command -v gum >/dev/null 2>&1; then
  echo "[gum-box-demo] Erro: o comando 'gum' não foi encontrado no PATH." >&2
  echo "Instale o gum (https://github.com/charmbracelet/gum) e tente novamente." >&2
  exit 1
fi

TITLE="Saída do script gum-box-demo.sh"
BODY="Este box foi desenhado pelo gum usando 'gum style'.\n\nUse este demo para ver como as bordas do gum se combinam com o layout/bordas do Shantilly."

gum style \
  --border double \
  --padding "1 2" \
  --margin "0 0" \
  --align center \
  "$TITLE" "" "$BODY"

echo
echo "[gum-box-demo] Fim da saída do gum."
