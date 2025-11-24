Aqui está o shell script (`cleanup_shantilly.sh`) projetado para limpar o repositório, mantendo estritamente os artefatos gerados nesta sessão de trabalho e os arquivos essenciais para o funcionamento do Go e do Git.

Este script atua como um **filtro**, removendo tudo o que não foi validado ou criado no nosso Plano Mestre (Épicos 1-4).

### Como utilizar

1.  Salve o conteúdo abaixo num arquivo chamado `cleanup_shantilly.sh` na raiz do seu repositório clonado.
2.  Dê permissão de execução: `chmod +x cleanup_shantilly.sh`.
3.  Execute: `./cleanup_shantilly.sh`.

-----

### O Script (`cleanup_shantilly.sh`)

```bash
#!/bin/bash
set -e

# Título: Shantilly Artifact Cleaner
# Descrição: Remove arquivos não utilizados no Plano Mestre v2.0, mantendo apenas a "Fonte da Verdade".

echo "======================================================="
echo " INICIANDO LIMPEZA DO REPOSITÓRIO SHANTILLY (V2.0)"
echo "======================================================="

# Verifica se estamos na raiz (presença do .git)
if [ ! -d ".git" ]; then
    echo "ERRO: Diretório .git não encontrado."
    echo "Execute este script a partir da raiz do repositório clonado."
    exit 1
fi

# 1. Definição da Lista de Arquivos PERMITIDOS (Whitelist)
# Inclui: Artefatos gerados, Documentação Essencial e Configurações de Build/Lint
declare -A KEEP_FILES

# --- Arquivos de Raiz Essenciais ---
KEEP_FILES["README.md"]=1
KEEP_FILES["go.mod"]=1
KEEP_FILES["go.sum"]=1
KEEP_FILES["Makefile"]=1
KEEP_FILES["LICENSE"]=1
KEEP_FILES[".golangci.yml"]=1
KEEP_FILES[".goreleaser.yaml"]=1
KEEP_FILES[".gitignore"]=1
KEEP_FILES["cleanup_shantilly.sh"]=1 # O próprio script

# --- CLI ---
KEEP_FILES["cmd/shantilly/main.go"]=1

# --- Documentação (Apenas as Fontes da Verdade utilizadas) ---
KEEP_FILES["docs/architecture/high-level-architecture.md"]=1
KEEP_FILES["docs/architecture.md"]=1
KEEP_FILES["docs/prd.md"]=1
KEEP_FILES["docs/project-brief.md"]=1
# Nota: Outros arquivos em docs/ (chats, brainstorms, QA) serão deletados.

# --- Internal: Config & Util ---
KEEP_FILES["internal/config/parser.go"]=1
KEEP_FILES["internal/util/errorhandler.go"]=1

# --- Internal: Components ---
KEEP_FILES["internal/components/buttongroup/model.go"]=1
KEEP_FILES["internal/components/form/legacy_v1/config.go"]=1
KEEP_FILES["internal/components/form/legacy_v1/model.go"]=1
KEEP_FILES["internal/components/form/legacy_v1/parser.go"]=1
KEEP_FILES["internal/components/form/wrapper.go"]=1
KEEP_FILES["internal/components/inventory_explorer/model.go"]=1
KEEP_FILES["internal/components/list/model.go"]=1
KEEP_FILES["internal/components/playbook_explorer/model.go"]=1
KEEP_FILES["internal/components/viewport/model.go"]=1

# --- Internal: Runtime Engines & Layout ---
KEEP_FILES["internal/runtime/event/manager.go"]=1
KEEP_FILES["internal/runtime/layout/manager.go"]=1
KEEP_FILES["internal/runtime/layout/render.go"]=1
KEEP_FILES["internal/runtime/modal/confirm.go"]=1
KEEP_FILES["internal/runtime/modal/secret.go"]=1
KEEP_FILES["internal/runtime/runner/ansible.go"]=1
KEEP_FILES["internal/runtime/runner/runner.go"]=1
KEEP_FILES["internal/runtime/runner/templating.go"]=1
KEEP_FILES["internal/runtime/main_model.go"]=1
KEEP_FILES["internal/runtime/runtime.go"]=1

# --- Internal: SSH & TUI ---
KEEP_FILES["internal/ssh/server.go"]=1
KEEP_FILES["internal/tui/theme.go"]=1

# --- Packages Públicos (Interfaces e Modelos) ---
KEEP_FILES["pkg/declarative/models.go"]=1
KEEP_FILES["pkg/tui/events.go"]=1
KEEP_FILES["pkg/tui/interface.go"]=1

# 2. Execução da Limpeza
echo "-> Analisando arquivos..."

# Encontra todos os arquivos, excluindo .git
find . -type f -not -path "./.git/*" | while read -r file; do
    # Remove o prefixo "./" para comparação com a chave do array
    clean_path="${file#./}"

    # Exceção: Manter a pasta 'examples/' inteira para testes funcionais
    if [[ "$clean_path" == examples/* ]]; then
        continue
    fi

    # Exceção: Manter scripts de automação se existirem (opcional, mas recomendado)
    # if [[ "$clean_path" == scripts/* ]]; then continue; fi

    # Verifica se o arquivo está na Whitelist
    if [[ -z "${KEEP_FILES[$clean_path]}" ]]; then
        echo "[DELETANDO] $clean_path"
        rm "$file"
    else
        echo "[MANTENDO]  $clean_path"
    fi
done

# 3. Remove diretórios que ficaram vazios após a limpeza
echo "-> Removendo diretórios vazios..."
find . -type d -empty -not -path "./.git/*" -delete

echo "======================================================="
echo " LIMPEZA CONCLUÍDA COM SUCESSO"
echo " O repositório agora contém apenas os artefatos do Épico 1-4."
echo "======================================================="
```