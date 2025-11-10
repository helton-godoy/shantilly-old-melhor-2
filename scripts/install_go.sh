#!/bin/bash

# Este script baixa e instala a versão mais recente do Go (golang)
# no diretório /usr/local/go e configura a variável de ambiente PATH.

set -e

# Verifica se o usuário tem privilégios de superusuário (root)
if [[ $EUID -ne 0 ]]; then
	echo "Este script precisa ser executado como root."
	echo "Por favor, use 'sudo ./install_go.sh'."
	exit 1
fi

# 1. Busca a versão mais recente do Go na página de downloads
#    e extrai o nome do arquivo para a arquitetura amd64
echo "Buscando a versão mais recente do Go..."
latest_version=$(curl -sL https://golang.org/dl/ | grep -oP 'go[0-9.]+\.linux-amd64' | head -n 1)

if [[ -z "$latest_version" ]]; then
	echo "Não foi possível encontrar a versão mais recente do Go."
	exit 1
fi

download_url="https://golang.org/dl/${latest_version}.tar.gz"
temp_file="/tmp/${latest_version}.tar.gz"

echo "Versão mais recente do Go encontrada: ${latest_version}"
echo "URL de download: ${download_url}"

# 2. Baixa o arquivo
echo "Baixando o Go..."
curl -sL "$download_url" -o "$temp_file"

# 3. Remove qualquer instalação anterior
echo "Removendo instalações anteriores de /usr/local/go..."
rm -rf /usr/local/go

# 4. Extrai o novo Go para /usr/local
echo "Instalando o Go em /usr/local/go..."
tar -C /usr/local -xzf "$temp_file"

# 5. Configura a variável de ambiente PATH
echo "Configurando a variável de ambiente PATH no arquivo ~/.bashrc..."
# Garante que a linha não seja adicionada várias vezes
if ! grep -q 'export PATH=$PATH:/usr/local/go/bin' ~/.bashrc; then
	echo 'export PATH=$PATH:/usr/local/go/bin' >>~/.bashrc
fi

# Opcional: configura o GOPATH no arquivo ~/.bashrc
# O GOPATH padrão agora é o diretório "go" dentro do diretório home do usuário.
# Esta configuração abaixo não é mais obrigatória para a maioria dos casos de uso modernos com módulos Go.
# Mas, se você quiser, pode descomentar.
# if ! grep -q 'export GOPATH=~/.go' ~/.bashrc; then
#   echo 'export GOPATH=~/.go' >> ~/.bashrc
#   echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
# fi

# 6. Remove o arquivo temporário
echo "Limpando arquivos temporários..."
rm "$temp_file"

echo "Instalação do Go concluída!"
echo "Para que as alterações entrem em vigor, execute: source ~/.bashrc"
echo "Ou, inicie um novo terminal."
echo ""
echo "Verifique a instalação com:"
echo "go version"
