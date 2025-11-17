#!/usr/bin/env bash
set -euo pipefail

echo "========================================"
echo " Shantilly Gate • E1.6 Encapsulamento do Legado"
echo "========================================"
echo ""

# Gate E1.6 — Encapsulamento do Legado (FormComponent + internal/tui)
# Referência normativa:
# - docs/qa/gates/1.x.legacy-formcomponent-encapsulation.yml
# - AGENTS.md
#
# Objetivo:
# - Bloquear qualquer uso de `FormConfig` ou `internal/tui` fora dos caminhos permitidos:
#   - internal/config/**
#   - internal/tui/**
#   - internal/components/form_component.go
#   - cmd/shantilly/main.go
#
# Notas:
# - Este script é idempotente e pode ser usado localmente e no CI.
# - Falha (exit != 0) se encontrar violações.

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

echo "[gate:E1.6] Verificando encapsulamento do legado (FormConfig/internal/tui)..."

# Lista de caminhos permitidos explicitamente
allowed_paths=(
  "internal/config/"
  "internal/tui/"
  "internal/components/form_component.go"
  "cmd/shantilly/main.go"
)

# Monta argumentos de exclusão para git grep com base nos caminhos permitidos
exclude_args=()
for p in "${allowed_paths[@]}"; do
  # Para diretórios, exclui recursivamente; para arquivos, exclui o arquivo.
  if [[ "$p" == */ ]]; then
    exclude_args+=(":!${p}*")
  else
    exclude_args+=(":!${p}")
  fi
done

violation=0

# Scan 1: `FormConfig` fora dos caminhos permitidos
if ! git grep -n "FormConfig" -- . "${exclude_args[@]}" >/dev/null 2>&1; then
  echo "[gate:E1.6] OK: Nenhuma referência ilegal a FormConfig encontrada."
else
  echo "[gate:E1.6][ERRO] Encontradas referências a FormConfig fora dos caminhos permitidos:"
  git grep -n "FormConfig" -- . "${exclude_args[@]}"
  violation=1
fi

# Scan 2: `internal/tui` fora dos caminhos permitidos
if ! git grep -n "internal/tui" -- . "${exclude_args[@]}" >/dev/null 2>&1; then
  echo "[gate:E1.6] OK: Nenhuma referência ilegal a internal/tui encontrada."
else
  echo "[gate:E1.6][ERRO] Encontradas referências a internal/tui fora dos caminhos permitidos:"
  git grep -n "internal/tui" -- . "${exclude_args[@]}"
  violation=1
fi

if [[ "$violation" -ne 0 ]]; then
	echo ""
	echo "[gate:E1.6] FALHA: Encapsulamento do legado violado. Corrija antes de fazer merge."
	echo "----------------------------------------"
	exit 1
fi

echo ""
echo "[gate:E1.6] SUCESSO: Legado confinado aos pontos autorizados."
echo "----------------------------------------"