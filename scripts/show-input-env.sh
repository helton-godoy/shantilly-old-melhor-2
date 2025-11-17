#!/usr/bin/env bash
set -euo pipefail

# Demo: exibe o valor da variável de ambiente ligada ao componente input.
# Para o app.input-demo.yaml atual, o ID do componente é "input_box",
# então a env gerada pelo runtime segue a convenção:
#   SHANTILLY_INPUT_INPUT_BOX

VAR_NAME="SHANTILLY_INPUT_INPUT_BOX"

echo "========================================"
echo " Shantilly Demo • Script: show-input-env.sh"
echo "========================================"
echo ""

echo "Nome da variável: $VAR_NAME"
echo "Valor recebido:  ${!VAR_NAME-<vazio>}"

echo ""
echo "Dica: altere o valor no campo de input e pressione Enter para ver a mudança."
