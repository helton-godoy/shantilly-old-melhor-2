#!/usr/bin/env bash
# MCP GitHub (admin) - Verificador C1 Branch Protection
# Objetivo: validar de forma objetiva se a proteção canônica (C1) está aplicada em main/master.
#
# Requisitos:
# - GH_TOKEN_ADMIN com permissão de leitura/admin no repo (já usado no apply).
#
# Saída:
# - Exit 0 se tudo conforme para main (e master, se existir).
# - Exit != 0 se qualquer condição C1 for violada.
# - Mensagens claras [C1][OK]/[C1][FAIL] para uso pelo agente/CI.

set -euo pipefail

OWNER="helton-godoy"
REPO="shantilly"
API_ROOT="https://api.github.com"
REQUIRED_CHECK="governanca-waves4-7 / Governança Waves 4-7 (Gates 1.x + Go checks)"
REQUIRED_APPROVALS=1

log_info() {
  printf '[C1][INFO] %s\n' "$*" 1>&2
}

log_ok() {
  printf '[C1][OK] %s\n' "$*" 1>&2
}

log_fail() {
  printf '[C1][FAIL] %s\n' "$*" 1>&2
}

require_token() {
  if [[ -z "${GH_TOKEN_ADMIN:-}" ]]; then
    log_fail "GH_TOKEN_ADMIN não definido. Configure um token de admin/leitura em ${OWNER}/${REPO}."
    exit 1
  fi
}

fetch_protection() {
  local branch="$1"
  local out_file="$2"

  local status
  status="$(
    curl -sS -o "${out_file}" -w "%{http_code}" \
      -H "Authorization: Bearer ${GH_TOKEN_ADMIN}" \
      -H "Accept: application/vnd.github.v3+json" \
      "${API_ROOT}/repos/${OWNER}/${REPO}/branches/${branch}/protection"
  )"

  echo "${status}"
}

branch_exists() {
  local branch="$1"
  local status
  status="$(
    curl -sS -o /tmp/c1_verify_branch_${branch}.json -w "%{http_code}" \
      -H "Authorization: Bearer ${GH_TOKEN_ADMIN}" \
      -H "Accept: application/vnd.github.v3+json" \
      "${API_ROOT}/repos/${OWNER}/${REPO}/branches/${branch}"
  )"

  case "${status}" in
    200) return 0 ;;
    404) return 1 ;;
    *) log_fail "Erro ao checar existência do branch '${branch}' (HTTP ${status})."; return 2 ;;
  esac
}

verify_branch_c1() {
  local branch="$1"

  if ! branch_exists "${branch}"; then
    if [[ "$?" -eq 1 ]]; then
      # 404: branch não existe; para main isso é erro, para master é aceitável.
      if [[ "${branch}" == "main" ]]; then
        log_fail "Branch 'main' não existe. C1 não pode ser satisfeito."
        return 1
      fi
      log_info "Branch '${branch}' não existe. Nenhuma verificação C1 necessária para ele."
      return 0
    fi
    # Erro inesperado já logado em branch_exists
    return 1
  fi

  local tmp="/tmp/c1_verify_protection_${branch}.json"
  local status
  status="$(fetch_protection "${branch}" "${tmp}")"

  if [[ "${status}" != "200" ]]; then
    log_fail "Proteção de branch para '${branch}' não retornou 200 (HTTP ${status})."
    return 1
  fi

  # required_pull_request_reviews.required_approving_review_count == REQUIRED_APPROVALS
  local approvals
  approvals="$(jq -r '.required_pull_request_reviews.required_approving_review_count // -1' "${tmp}")"
  if [[ "${approvals}" -ne "${REQUIRED_APPROVALS}" ]]; then
    log_fail "'${branch}': required_approving_review_count esperado=${REQUIRED_APPROVALS}, obtido=${approvals}."
    return 1
  fi

  # enforce_admins.enabled == true
  local enforce_admins
  enforce_admins="$(jq -r '.enforce_admins.enabled // .enforce_admins // false' "${tmp}")"
  if [[ "${enforce_admins}" != "true" ]]; then
    log_fail "'${branch}': enforce_admins não está habilitado."
    return 1
  fi

  # required_status_checks.strict == true
  local strict
  strict="$(jq -r '.required_status_checks.strict // false' "${tmp}")"
  if [[ "${strict}" != "true" ]]; then
    log_fail "'${branch}': strict (Require branches to be up to date) não está habilitado."
    return 1
  fi

  # required_status_checks.checks contém EXATAMENTE o REQUIRED_CHECK como único contexto
  local contexts
  contexts="$(jq -r '.required_status_checks.checks[]?.context' "${tmp}" | sort | tr '\n' ';')"
  if [[ "${contexts}" != "${REQUIRED_CHECK};" ]]; then
    log_fail "'${branch}': lista de required checks difere do esperado. Obtido='${contexts}', Esperado='${REQUIRED_CHECK}'."
    return 1
  fi

  log_ok "Branch '${branch}' está em conformidade completa com C1."
  return 0
}

main() {
  require_token

  local failed=0

  # main é obrigatório
  if ! verify_branch_c1 "main"; then
    failed=1
  fi

  # master é condicional: se existir, deve estar conforme; se não, é OK.
  if branch_exists "master"; then
    if ! verify_branch_c1 "master"; then
      failed=1
    fi
  else
    log_info "Branch 'master' ausente: C1 considerado satisfeito apenas com 'main'."
  fi

  if [[ "${failed}" -eq 0 ]]; then
    log_ok "C1 aplicado com sucesso: proteção canônica garantida."
    exit 0
  else
    log_fail "C1 NÃO está totalmente aplicado. Veja mensagens acima para detalhes."
    exit 2
  fi
}

main "$@"