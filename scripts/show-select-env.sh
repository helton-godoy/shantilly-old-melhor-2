#!/usr/bin/env bash
set -euo pipefail

VAR_NAME="SHANTILLY_SELECT_SELECT_BOX"

echo "========================================"
echo " Shantilly Demo • Script: show-select-env.sh"
echo "========================================"
echo ""

echo "Nome da variável (select): $VAR_NAME"
echo "Valor recebido:  ${!VAR_NAME-<vazio>}"
