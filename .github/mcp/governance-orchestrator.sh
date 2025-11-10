#!/usr/bin/env bash
# mcp-governanca-orchestrator (shim shell)
#
# Objetivo:
# - Fornecer um orquestrador mínimo, idempotente e auditável para a política C1,
#   alinhado com .github/mcp/governance-orchestrator-policy.yml.
# - Encapsular a execução coordenada do enforcer + verificador canônicos.
#
# Design:
# - Usa apenas contratos declarativos já presentes no repo.
# - Assume que o MCP GitHub (admin) está corretamente configurado com GH_TOKEN_ADMIN
#   no ambiente (não aqui no código).
# - Pode ser chamado por:
#   - ambientes MCP,
#   - workflows CI,
#   - agentes BMAD,
#   - outros orquestradores.
#
# Comportamento:
# - Leitura defensiva da política:
#   - Confirma que control=C1 está configurado com roles C1.enforce e C1.verify.
# - Execução normativa:
#   - Executa apply_c1_branch_protection.sh.
#   - Em seguida executa verify_c1_branch_protection.sh.
# - Idempotência:
#   - Reexecuções mantêm o estado convergindo para o payload C1.
#
# Saída:
# - Exit 0: política C1 aplicada e validada com sucesso.
# - Exit != 0: falha clara, não oculta, com mensagens para auditoria.

set -euo pipefail

POLICY_FILE=".github/mcp/governance-orchestrator-policy.yml"
ENFORCER_SCRIPT=".github/mcp/apply_c1_branch_protection.sh"
VERIFIER_SCRIPT=".github/mcp/verify_c1_branch_protection.sh"

log_info() {
  printf '[MCP-GOV][INFO] %s\n' "$*" 1>&2
}

log_ok() {
  printf '[MCP-GOV][OK] %s\n' "$*" 1>&2
}

log_fail() {
  printf '[MCP-GOV][FAIL] %s\n' "$*" 1>&2
}

require_file() {
  local path="$1"
  local desc="$2"
  if [[ ! -f "${path}" ]]; then
    log_fail "Arquivo obrigatório ausente: ${desc} (${path})."
    exit 1
  fi
}

require_executable() {
  local path="$1"
  local desc="$2"
  if [[ ! -x "${path}" ]]; then
    log_info "Ajustando permissão de execução para ${desc} (${path})."
    chmod +x "${path}" || {
      log_fail "Não foi possível tornar ${desc} executável (${path})."
      exit 1
    }
  fi
}

require_token_hint() {
  if [[ -z "${GH_TOKEN_ADMIN:-}" ]]; then
    # Não loga valor, apenas ausência.
    log_info "GH_TOKEN_ADMIN não está definido no ambiente atual."
    log_info "Este shim assume que o MCP GitHub (admin) no ambiente real já possui o token configurado."
    log_info "Se estiver chamando diretamente este script fora do MCP, defina GH_TOKEN_ADMIN antes de executar."
  fi
}

validate_policy_c1() {
  # Validação leve usando grep/yq opcional se disponível.
  # Não falha se yq não existir; confia no contrato versionado.
  if ! command -v yq >/dev/null 2>&1; then
    log_info "yq não encontrado; validação sintática profunda da policy será ignorada."
    return 0
  fi

  # Verifica se C1 está declarado com roles exigidos.
  local enforce_role verify_role
  enforce_role="$(yq '.analysis.gaps.required[] | select(.control == "C1") | .needs[]?.role' "${POLICY_FILE}" | grep -c 'C1.enforce' || true)"
  verify_role="$(yq '.analysis.gaps.required[] | select(.control == "C1") | .needs[]?.role' "${POLICY_FILE}" | grep -c 'C1.verify' || true)"

  if [[ "${enforce_role}" -lt 1 || "${verify_role}" -lt 1 ]]; then
    log_fail "Policy não declara corretamente C1.enforce/C1.verify em analysis.gaps.required."
    exit 1
  fi

  log_ok "Policy C1 em ${POLICY_FILE} aparenta consistente (C1.enforce/C1.verify declarados)."
}

run_c1() {
  log_info "Iniciando execução normativa C1 via orquestrador MCP de governança..."

  # Enforcer C1
  log_info "Executando enforcer C1: ${ENFORCER_SCRIPT}"
  GH_TOKEN_ADMIN="${GH_TOKEN_ADMIN:-${GH_TOKEN_ADMIN-}}" "${ENFORCER_SCRIPT}"

  # Verificador C1
  log_info "Executando verificador C1: ${VERIFIER_SCRIPT}"
  GH_TOKEN_ADMIN="${GH_TOKEN_ADMIN:-${GH_TOKEN_ADMIN-}}" "${VERIFIER_SCRIPT}"

  log_ok "Política C1 aplicada e validada com sucesso pelo orquestrador."
}

main() {
  # 1) Pré-checagens
  require_file "${POLICY_FILE}" "política do orquestrador MCP de governança"
  require_file "${ENFORCER_SCRIPT}" "enforcer C1"
  require_file "${VERIFIER_SCRIPT}" "verificador C1"

  require_executable "${ENFORCER_SCRIPT}" "enforcer C1"
  require_executable "${VERIFIER_SCRIPT}" "verificador C1"

  require_token_hint
  validate_policy_c1

  # 2) Execução normativa C1 (idempotente)
  run_c1
}

main "$@"