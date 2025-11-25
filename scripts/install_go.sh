#!/bin/bash

# Este script instala o Go e prepara o ambiente de desenvolvimento
# para o projeto Shantilly em uma máquina limpa:
# - Instala o Go em /usr/local/go
# - Configura PATH/GOPATH/GOBIN para o usuário chamador
# - Instala gofumpt e golangci-lint nas versões usadas localmente/na CI

set -e

# Verifica se o usuário tem privilégios de superusuário (root)
if [[ $EUID -ne 0 ]]; then
	echo "Este script precisa ser executado como root."
	echo "Por favor, use 'sudo ./install_go.sh'."
	exit 1
fi

# Descobre o diretório do script e a raiz do repositório
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Define o usuário alvo (quem chamou o sudo, se houver)
TARGET_USER="${SUDO_USER:-$USER}"
TARGET_HOME="${HOME}"
if [[ -n "${SUDO_USER}" ]]; then
	TARGET_HOME="/home/${SUDO_USER}"
fi

# Descobre a versão do Go a partir do go.mod, se possível
GO_VERSION=""
if [[ -f "${REPO_ROOT}/go.mod" ]]; then
	GO_VERSION=$(grep '^go ' "${REPO_ROOT}/go.mod" | awk '{print $2}')
fi

ARCH_SUFFIX="linux-amd64"
BASE_URL="https://go.dev/dl"

if [[ -n "${GO_VERSION}" ]]; then
	# Usa a versão definida no go.mod
	TARBALL="go${GO_VERSION}.${ARCH_SUFFIX}"
	download_url="${BASE_URL}/${TARBALL}.tar.gz"
	echo "Usando versão do Go definida em go.mod: go${GO_VERSION}"
else
	# Fallback: busca a versão mais recente estável
	echo "Buscando a versão mais recente do Go..."
	latest_version=$(curl -sL "${BASE_URL}/" | grep -oE 'go[0-9.]+\.linux-amd64' | head -n 1)
	if [[ -z "${latest_version}" ]]; then
		echo "Não foi possível encontrar a versão mais recente do Go."
		exit 1
	fi
	TARBALL="${latest_version}"
	download_url="${BASE_URL}/${TARBALL}.tar.gz"
fi

temp_file="/tmp/${TARBALL}.tar.gz"

echo "Baixando o Go de: ${download_url}"
curl -sL "${download_url}" -o "${temp_file}"

# Remove qualquer instalação anterior de Go
echo "Removendo instalações anteriores de /usr/local/go..."
rm -rf /usr/local/go

# Extrai o Go em /usr/local
echo "Instalando o Go em /usr/local/go..."
tar -C /usr/local -xzf "${temp_file}"

# Configura PATH/GOPATH/GOBIN para o usuário alvo
BASHRC="${TARGET_HOME}/.bashrc"

echo "Configurando PATH/GOPATH/GOBIN para o usuário ${TARGET_USER} em ${BASHRC}..."

# Adiciona /usr/local/go/bin ao PATH, se ainda não estiver presente
if ! grep -q '/usr/local/go/bin' "${BASHRC}" 2>/dev/null; then
	echo 'export PATH=$PATH:/usr/local/go/bin' >> "${BASHRC}"
fi

# Configura GOPATH e GOBIN para o usuário alvo, se ainda não estiverem configurados
if ! grep -q 'export GOPATH=' "${BASHRC}" 2>/dev/null; then
	echo 'export GOPATH=$HOME/go' >> "${BASHRC}"
fi
if ! grep -q 'export GOBIN=' "${BASHRC}" 2>/dev/null; then
	echo 'export GOBIN=$GOPATH/bin' >> "${BASHRC}"
fi
if ! grep -q '$GOBIN' "${BASHRC}" 2>/dev/null; then
	echo 'export PATH=$PATH:$GOBIN' >> "${BASHRC}"
fi

# Garante que os diretórios GOPATH/GOBIN existam
mkdir -p "${TARGET_HOME}/go/bin"
chown -R "${TARGET_USER}:${TARGET_USER}" "${TARGET_HOME}/go"

# Instala gofumpt e golangci-lint para o usuário alvo
# - gofumpt: formatador usado pelo Makefile (target fmt/fmt-check)
# - golangci-lint: mesma ferramenta usada na CI (versão fixada)

echo "Instalando gofumpt e golangci-lint para o usuário ${TARGET_USER}..."

su - "${TARGET_USER}" -c "/usr/local/go/bin/go install mvdan.cc/gofumpt@latest"
su - "${TARGET_USER}" -c "/usr/local/go/bin/go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1"

# Remove o arquivo temporário
echo "Limpando arquivos temporários..."
rm -f "${temp_file}"

echo "Instalação do Go e ferramentas concluída!"
echo "Abra um novo terminal ou execute: source ${BASHRC}"
echo "Depois disso, você deve conseguir executar:"
echo "  go version"
echo "  gofumpt -h"
echo "  golangci-lint version"
