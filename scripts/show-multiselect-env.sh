#!/usr/bin/env bash
set -euo pipefail

VAR_NAME="SHANTILLY_MULTISELECT_MULTISELECT_BOX"

echo "========================================"
echo " Shantilly Demo • Script: show-multiselect-env.sh"
echo "========================================"
echo ""

echo "Nome da variável (multiselect): $VAR_NAME"
echo "Valor recebido (ids separados por vírgula):  ${!VAR_NAME-<vazio>}"
