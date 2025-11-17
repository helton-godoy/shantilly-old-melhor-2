#!/bin/bash

# start_setup_dev.sh
# Orquestrador de setup para o ambiente de desenvolvimento do Shantilly.
# - Detecta Debian/Ubuntu
# - Usa menu interativo com gum (se disponível) para escolher módulos
# - Cada módulo é idempotente e instala dependências oficiais
#
# Requisitos:
# - Executar como root (sudo)
# - Repositório já clonado
# - Opcional: bin/gum presente para UX aprimorada

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
INSTALL_GO_SCRIPT="${SCRIPT_DIR}/install_go.sh"
GUM_BIN="${REPO_ROOT}/bin/gum"

TARGET_USER="${SUDO_USER:-$USER}"
TARGET_HOME="/home/${TARGET_USER}"
if [[ "${TARGET_USER}" == "root" ]]; then
  TARGET_HOME="/root"
fi

ACT_VERSION="v0.2.59"
ACT_ARCHIVE="act_${ACT_VERSION#v}_Linux_x86_64.tar.gz"
ACT_URL="https://github.com/nektos/act/releases/download/${ACT_VERSION}/${ACT_ARCHIVE}"

function require_root() {
  if [[ ${EUID} -ne 0 ]]; then
    echo "Este script deve ser executado como root. Use: sudo scripts/start_setup_dev.sh"
    exit 1
  fi
}

function detect_distro() {
  if [[ ! -f /etc/os-release ]]; then
    echo "Não foi possível detectar distro (sem /etc/os-release)."
    exit 1
  fi
  . /etc/os-release

  local id_lower="${ID,,}"
  local id_like_lower="${ID_LIKE:-}"
  id_like_lower="${id_like_lower,,}"

  case "${id_lower}" in
    debian|ubuntu)
      DISTRO_ID="${id_lower}"
      VERSION_CODENAME="${VERSION_CODENAME:-${UBUNTU_CODENAME:-}}"
      ;;
    deepin)
      # Deepin é baseado em Debian; tratamos como Debian para repositórios oficiais.
      DISTRO_ID="debian"
      VERSION_CODENAME="${VERSION_CODENAME:-bookworm}"
      ;;
    *)
      if [[ -n "${id_like_lower}" && "${id_like_lower}" == *"debian"* ]]; then
        DISTRO_ID="debian"
        VERSION_CODENAME="${VERSION_CODENAME:-bookworm}"
      else
        echo "Distro ${ID} (ID_LIKE=${ID_LIKE:-n/a}) não suportada. Este script cobre apenas Debian/Ubuntu/derivadas."
        exit 1
      fi
      ;;
  esac

  DISTRO_VERSION_ID="${VERSION_ID:-}"

  if [[ -z "${VERSION_CODENAME}" ]]; then
    VERSION_CODENAME="$(lsb_release -cs 2>/dev/null || echo stable)"
  fi
}

function ensure_base_packages() {
  echo "[Base] Atualizando índices APT e instalando pacotes essenciais..."
  apt-get update -y
  apt-get install -y \
    curl \
    ca-certificates \
    gnupg \
    lsb-release \
    git \
    tar \
    unzip
}

function ensure_go_stack() {
  echo "[Go] Executando install_go.sh para provisionar Go + ferramentas..."
  if [[ ! -x "${INSTALL_GO_SCRIPT}" ]]; then
    echo "install_go.sh não encontrado ou sem permissão de execução em ${INSTALL_GO_SCRIPT}."
    exit 1
  fi
  "${INSTALL_GO_SCRIPT}"
}

function ensure_docker() {
  echo "[Docker] Instalando Docker Engine a partir do repositório oficial..."
  apt-get remove -y docker docker-engine docker.io containerd runc >/dev/null 2>&1 || true

  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/${DISTRO_ID}/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  chmod a+r /etc/apt/keyrings/docker.gpg

  local source_entry="deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/${DISTRO_ID} ${VERSION_CODENAME} stable"
  if ! grep -q "download.docker.com" /etc/apt/sources.list /etc/apt/sources.list.d/* 2>/dev/null; then
    echo "${source_entry}" > /etc/apt/sources.list.d/docker.list
  fi

  apt-get update -y
  apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

  if ! getent group docker >/dev/null 2>&1; then
    groupadd docker
  fi
  usermod -aG docker "${TARGET_USER}"
  echo "[Docker] Usuário ${TARGET_USER} adicionado ao grupo docker (reabra a sessão para aplicar)."
}

function ensure_node() {
  echo "[Node] Instalando Node.js LTS (NodeSource 20.x)..."
  curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
  apt-get install -y nodejs
  echo "[Node] node $(node -v), npm $(npm -v) instalados."
}

function ensure_uv_requirements() {
  apt-get install -y python3 python3-venv python3-pip || true
}

function ensure_uv() {
  ensure_uv_requirements
  echo "[uv] Instalando uv/uvx para usuário ${TARGET_USER}..."
  su - "${TARGET_USER}" -c "curl -LsSf https://astral.sh/uv/install.sh | sh"
}

function ensure_act() {
  echo "[act] Instalando act ${ACT_VERSION}..."
  local tmp_dir
  tmp_dir="$(mktemp -d)"
  pushd "${tmp_dir}" >/dev/null
  curl -fsSL "${ACT_URL}" -o "${ACT_ARCHIVE}"
  tar -xzf "${ACT_ARCHIVE}"
  install -m 0755 act /usr/local/bin/act
  popd >/dev/null
  rm -rf "${tmp_dir}"
  echo "[act] $(act --version) instalado."
}

function install_all() {
  ensure_go_stack
  ensure_docker
  ensure_node
  ensure_uv
  ensure_act
}

function show_summary() {
  echo "\nResumo do ambiente:"
  command -v go >/dev/null 2>&1 && echo "- $(go version)"
  command -v gofumpt >/dev/null 2>&1 && echo "- gofumpt $(gofumpt -version 2>/dev/null || echo instalado)"
  command -v golangci-lint >/dev/null 2>&1 && echo "- golangci-lint $(golangci-lint version 2>/dev/null || echo instalado)"
  command -v docker >/dev/null 2>&1 && echo "- Docker $(docker --version)"
  command -v node >/dev/null 2>&1 && echo "- Node $(node -v) / npm $(npm -v)"
  command -v uv >/dev/null 2>&1 && echo "- uv $(uv --version 2>/dev/null || echo instalado)"
  command -v act >/dev/null 2>&1 && echo "- act $(act --version 2>/dev/null || echo instalado)"
  echo "\nPróximos passos:"
  echo "  - Reabra a sessão (grupos PATH/docker)."
  echo "  - make fmt && make check"
  echo "  - Para GitHub Actions local: act -j build-and-test"
}

function run_module_from_label() {
  local label="$1"
  case "${label}" in
    "Go"* ) ensure_go_stack ;;
    "Docker"* ) ensure_docker ;;
    "Node"* ) ensure_node ;;
    "uv"* ) ensure_uv ;;
    "act"* ) ensure_act ;;
    * ) echo "[Menu] Opção desconhecida: ${label}" ;;
  esac
}

function select_modules_menu() {
  local modules=(
    "Go + Ferramentas"
    "Docker"
    "Node/npm/npx"
    "uv/uvx"
    "act"
  )

  if [[ -x "${GUM_BIN}" ]]; then
    clear
    local selections
    selections=$("${GUM_BIN}" choose \
      --header="Selecione um ou mais módulos (Barra de espaço para marcar, Enter para confirmar):" \
      --no-limit \
      --cursor="👉 " \
      --cursor-prefix="[ ]" \
      --selected-prefix="[x]" \
      --unselected-prefix="[ ]" \
      "${modules[@]}") || true
    clear

    if [[ -z "${selections}" ]]; then
      echo "[Menu] Nenhum módulo selecionado."
      return
    fi

    while IFS= read -r module; do
      [[ -z "${module}" ]] && continue
      run_module_from_label "${module}"
    done <<< "${selections}"
    return
  fi

  echo "gum não encontrado em ${GUM_BIN}."
  echo "Selecione os módulos desejados digitando os números separados por espaço (ex.: 1 3 5)."
  local idx=1
  for module in "${modules[@]}"; do
    printf "%d) %s\n" "${idx}" "${module}"
    idx=$((idx + 1))
  done
  read -rp "Sua escolha (ENTER para cancelar): " selection_line
  [[ -z "${selection_line}" ]] && echo "[Menu] Nenhum módulo selecionado." && return

  for num in ${selection_line}; do
    if [[ ${num} -ge 1 && ${num} -le ${#modules[@]} ]]; then
      run_module_from_label "${modules[$((num-1))]}"
    else
      echo "[Menu] Índice inválido: ${num}"
    fi
  done
}

function menu() {
  local options=(
    "Instalar TUDO (Go + Docker + Node + uv + act)"
    "Selecionar módulos específicos"
    "Sair"
  )

  local choice
  if [[ -x "${GUM_BIN}" ]]; then
    clear
    choice=$("${GUM_BIN}" choose \
      --header="Selecione uma ação:" \
      --cursor="👉 " \
      "${options[@]}")
    clear
  else
    echo "gum não encontrado em ${GUM_BIN}. Usando seletor padrão."
    select opt in "${options[@]}"; do
      choice="${opt}"
      break
    done
  fi

  case "${choice}" in
    "Instalar TUDO"*) install_all ;;
    "Selecionar módulos"*) select_modules_menu ;;
    "Sair"|"") echo "Saindo sem alterações."; exit 0 ;;
    *) echo "Opção desconhecida: ${choice}"; exit 1 ;;
  esac
}

require_root
detect_distro
ensure_base_packages
menu
show_summary
