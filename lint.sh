#!/bin/bash

# lint.sh: Um script abrangente para garantir a qualidade do código localmente.
# Ele executa formatação, linting e testes.
#
# COMO USAR:
# 1. Dê permissão de execução: chmod +x lint.sh
# 2. Execute na raiz do projeto: ./lint.sh

# 'set -e' garante que o script sairá imediatamente se qualquer comando falhar.
set -e

# --- PASSO 1: VERIFICAÇÃO DE FORMATAÇÃO ---
echo "🔎 Verificando a formatação com gofumpt..."
# O comando 'gofumpt -l .' lista os arquivos que precisam de formatação.
# O bloco 'if' verifica se a saída do comando não está vazia.
if [[ -n $(gofumpt -l .) ]]; then
    echo "❌ Alguns arquivos precisam de formatação. Execute 'gofumpt -w .' para corrigi-los."
    exit 1
fi
echo "✅ Formatação está correta."
echo "" # Linha em branco para espaçamento

# --- PASSO 2: ANÁLISE ESTÁTICA (LINTING) ---
echo "🔬 Executando o linter (golangci-lint)..."
# Executa o linter usando nosso arquivo de configuração para consistência com o CI.
golangci-lint run --config=.golangci.yml ./...
echo "✅ Análise estática concluída sem problemas."
echo ""

# --- PASSO 3: EXECUÇÃO DOS TESTES ---
echo "🧪 Executando os testes com o detector de race condition..."
# Executa os testes em modo verboso e com o detector de concorrência ativado.
go test -v -race ./...
echo "✅ Testes concluídos com sucesso."
echo ""

# --- SUCESSO ---
echo "🎉🚀 Todos os checks locais passaram com sucesso! Código pronto para o push."
