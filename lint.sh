#!/bin/bash

# lint.sh: Um script abrangente para garantir a qualidade do código localmente.
# Ele executa formatação, linting e testes.
#
# COMO USAR:
# 1. Dê permissão de execução: chmod +x lint.sh
# 2. Execute na raiz do projeto: ./lint.sh

# 'set -e' garante que o script sairá imediatamente se qualquer comando falhar.
set -e

# Flag opcional: --markdown habilita o lint de Markdown via make lint-md.
RUN_MARKDOWN_LINT=false
for arg in "$@"; do
	if [[ "$arg" == "--markdown" ]]; then
		RUN_MARKDOWN_LINT=true
	fi
done

# --- PASSO 1: VERIFICAÇÃO DE FORMATAÇÃO ---
echo "🔎 Verificando a formatação com gofumpt..."

if ! command -v gofumpt >/dev/null 2>&1; then
	echo "❌ A ferramenta 'gofumpt' não foi encontrada no PATH. Instale-a (ex.: 'go install github.com/mvdan/gofumpt@latest') e tente novamente."
	exit 1
fi

# O comando 'gofumpt -l .' lista os arquivos que precisam de formatação.
if [[ -n $(gofumpt -l .) ]]; then
	echo "❌ Alguns arquivos precisam de formatação. Execute 'gofumpt -w .' para corrigi-los."
	exit 1
fi
echo "✅ Formatação está correta."
echo "" # Linha em branco para espaçamento

# --- PASSO 2: ANÁLISE ESTÁTICA (LINTING) ---
echo "🔬 Executando o linter (golangci-lint)..."

if ! command -v golangci-lint >/dev/null 2>&1; then
	echo "❌ A ferramenta 'golangci-lint' não foi encontrada no PATH. Instale-a (ex.: 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest') e tente novamente."
	exit 1
fi

# Executa o linter usando nosso arquivo de configuração para consistência com o CI.
golangci-lint run --config=.golangci.yml ./...
echo "✅ Análise estática concluída sem problemas."
echo ""

# --- PASSO 3: EXECUÇÃO DOS TESTES ---
echo "🧪 Executando os testes com o detector de race condition..."
# Executa os testes em modo verboso e com o detector de concorrência ativado.
go test -v -race ./...
echo "✅ Testes concluídos com sucesso."
echo "" # Linha em branco para espaçamento

# --- PASSO 4: LINT DE MARKDOWN (OPCIONAL) ---
if [[ "$RUN_MARKDOWN_LINT" == "true" ]]; then
	echo "📄 Verificando Markdown com markdownlint-cli2 via Makefile..."
	make lint-md
	echo "✅ Lint de Markdown concluído."
	echo "" # Linha em branco para espaçamento
else
	echo "ℹ️ Lint de Markdown ignorado (use '--markdown' para habilitar)."
	echo "" # Linha em branco para espaçamento
fi

# --- SUCESSO ---
echo "🎉🚀 Todos os checks locais passaram com sucesso! Código pronto para o push."
