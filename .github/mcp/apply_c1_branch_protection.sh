#!/usr/bin/env bash
# MCP GitHub (admin) - C1 Branch Protection Enforcer
# Objetivo: aplicar de forma idempotente a proteção canônica (C1) em helton-godoy/shantilly.
# Requisitos:
# - Variável GH_TOKEN_ADMIN exportada com token de admin no repositório.
# - Execução segura, sem logar o token, com chamadas HTTP exatas.
#
# Comportamento:
# - Sempre aplica proteção em main.
# - Aplica proteção em master apenas se o branch existir (GET 200).
# - PUT é idempotente: reaplicações não afrouxam regras nem criam estados duplicados.
#
# Uso típico (no host MCP ou job de governança):
# - chmod +x .github/mcp/apply_c1_branch_protection.sh
# - ./github/mcp/apply_c1_branch_protection.sh
#
# Integração recomendada:
# - Executar on-boot do MCP GitHub (admin).
# - Executar em rotina periódica (cron-like) como auditoria automática.

set -euo pipefail

OWNER="helton-godoy"
REPO="shantilly"
API_ROOT="https://api.github.com"
REQUIRED_CHECK="governanca-waves4-7 / Governança Waves 4-7 (Gates 1.x + Go checks)"

# JSON canônico da proteção (mantido em linha única/literal para evitar variações)
read -r -d '' PROTECTION_BODY << 'JSON'
{
  "required_status_checks": {
    "strict": true,
    "checks": [
      {
        "context": "governanca-waves4-7 / Governança Waves 4-7 (Gates 1.x + Go checks)"
      }
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1
  },
  "restrictions": null
}
JSON

log_info() {
  printf '[C1] %s\n' "$*" 1>&2
}

log_error() {
  printf '[C1][ERROR] %s\n' "$*" 1>&2
}

require_token() {
  if [[ -z "${GH_TOKEN_ADMIN:-}" ]]; then
    log_error "GH_TOKEN_ADMIN não definido. Configure um token com permissão de admin em ${OWNER}/${REPO}."
    exit 1
  fi
}

apply_protection() {
  local branch="$1"

  log_info "Aplicando proteção canônica em '${branch}'..."

  # Chamada PUT exata conforme especificação
  # - required_status_checks.strict = true
  # - checks = [ REQUIRED_CHECK ]
  # - enforce_admins = true
  # - required_pull_request_reviews.required_approving_review_count = 1
  # - restrictions = null
  local status
  status="$(
    curl -sS -o /tmp/c1_protection_response_${branch}.json -w "%{http_code}" \
      -X PUT \
      -H "Authorization: Bearer ${GH_TOKEN_ADMIN}" \
      -H "Accept: application/vnd.github.v3+json" \
      -H "Content-Type: application/json" \
      "${API_ROOT}/repos/${OWNER}/${REPO}/branches/${branch}/protection" \
      -d "${PROTECTION_BODY}"
  )"

  if [[ "${status}" == "200" || "${status}" == "201" || "${status}" == "202" ]]; then
    log_info "Proteção aplicada/garantida com sucesso em '${branch}' (HTTP ${status})."
  else
    log_error "Falha ao aplicar proteção em '${branch}' (HTTP ${status}). Verifique /tmp/c1_protection_response_${branch}.json."
    exit 1
  fi
}

branch_exists() {
  local branch="$1"

  local status
  status="$(
    curl -sS -o /tmp/c1_branch_check_${branch}.json -w "%{http_code}" \
      -H "Authorization: Bearer ${GH_TOKEN_ADMIN}" \
      -H "Accept: application/vnd.github.v3+json" \
      "${API_ROOT}/repos/${OWNER}/${REPO}/branches/${branch}"
  )"

  case "${status}" in
    200)
      log_info "Branch '${branch}' existe (HTTP 200)."
      return 0
      ;;
    404)
      log_info "Branch '${branch}' não existe (HTTP 404). Nenhuma proteção necessária."
      return 1
      ;;
    *)
      log_error "Erro ao verificar existência do branch '${branch}' (HTTP ${status}). Verifique /tmp/c1_branch_check_${branch}.json."
      # Em caso de erro inesperado, não forçar proteção e deixar para próxima execução/ação manual.
      return 2
      ;;
  esac
}

main() {
  require_token

  # 1) Sempre aplica/garante proteção em main
  apply_protection "main"

  # 2) Se master existir, aplica a mesma proteção
  if branch_exists "master"; then
    # Apenas aplica se retorno explícito foi 0 (existe). Se 2 (erro), não tentar PUT.
    apply_protection "master"
  fi

  log_info "Fluxo C1 concluído. Verifique na UI/API do GitHub se main/master refletem exatamente a proteção canônica."
}

main "$@"