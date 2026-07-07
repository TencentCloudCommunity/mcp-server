#!/usr/bin/env bash
# 按固定顺序每次只执行一个写接口，并在执行后自动停下，方便人工观察。
#
# 用法：
#   ./scripts/run_next_step_observe.sh
#   ./scripts/run_next_step_observe.sh --status
#   ./scripts/run_next_step_observe.sh --reset
#   ./scripts/run_next_step_observe.sh --config /absolute/path/to/full_test_plan.yaml
#
# 说明：
#   1. 默认会按 full_test_plan.yaml 的顺序逐步推进
#   2. 每次只执行 1 个 step，然后调用 run_verify.sh 做只读回看并退出
#   3. 成功后自动把进度推进到下一个 step；失败时保留当前进度，便于重试
#   4. --reset 只重置进度，不会执行任何云上操作

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
DEFAULT_CONFIG_PATH="${PROJECT_DIR}/scripts/full_test_plan.yaml"
STATE_DIR="${PROJECT_DIR}/scripts/observe_state"
LOG_ROOT="${PROJECT_DIR}/scripts/observe_logs"
CONFIG_PATH="${DEFAULT_CONFIG_PATH}"
ACTION="run"

STEPS=(
  "ModifyDBInstanceName"
  "CreateAccount"
  "ResetAccountPassword"
  "ModifyAccountPrivilegesGrant"
  "ModifyAccountPrivilegesRevoke"
  "DeleteAccount"
  "CreateDatabase"
  "ModifyDatabaseOwner"
  "DescribeBackupDownloadURL"
  "CreateBaseBackup"
  "OpenDBExtranetAccess"
  "CloseDBExtranetAccess"
  "ModifyDBInstanceSecurityGroups"
  "UpgradeDBInstanceKernelVersion"
  "RestartDBInstance"
  "IsolateDBInstances"
  "DisIsolateDBInstances"
  "ModifyDBInstanceParameters"
  "CreateInstances"
  "ModifyDBInstanceSpec"
  "CloneDBInstance"
  "CreateReadOnlyDBInstance"
)

resolve_path() {
  local path="$1"
  if [[ "${path}" = /* ]]; then
    printf '%s\n' "${path}"
    return
  fi
  printf '%s\n' "$(cd "$(dirname "${path}")" && pwd)/$(basename "${path}")"
}

usage() {
  cat <<EOF
用法:
  ./scripts/run_next_step_observe.sh
  ./scripts/run_next_step_observe.sh --status
  ./scripts/run_next_step_observe.sh --reset
  ./scripts/run_next_step_observe.sh --config /absolute/path/to/full_test_plan.yaml
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --status)
      ACTION="status"
      shift
      ;;
    --reset)
      ACTION="reset"
      shift
      ;;
    --config)
      if [[ $# -lt 2 ]]; then
        echo "错误：--config 需要一个配置文件路径" >&2
        exit 1
      fi
      CONFIG_PATH="$(resolve_path "$2")"
      shift 2
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      echo "错误：未知参数 $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

STATE_FILE="${STATE_DIR}/next_step.env"
NEXT_INDEX=""
SESSION_ID=""
SAVED_CONFIG_PATH=""

mkdir -p "${STATE_DIR}" "${LOG_ROOT}"

save_state() {
  cat > "${STATE_FILE}" <<EOF
NEXT_INDEX=${NEXT_INDEX}
SESSION_ID=${SESSION_ID}
CONFIG_PATH=${CONFIG_PATH}
EOF
}

load_state() {
  if [[ -f "${STATE_FILE}" ]]; then
    # shellcheck disable=SC1090
    source "${STATE_FILE}"
    NEXT_INDEX="${NEXT_INDEX:-}"
    SESSION_ID="${SESSION_ID:-}"
    SAVED_CONFIG_PATH="${CONFIG_PATH:-}"
  fi
}

init_state() {
  NEXT_INDEX=1
  SESSION_ID="$(date +%Y%m%d-%H%M%S)"
  save_state
}

ensure_state() {
  load_state
  if [[ -z "${NEXT_INDEX}" || -z "${SESSION_ID}" ]]; then
    init_state
    return
  fi
  if [[ -n "${SAVED_CONFIG_PATH}" && "${CONFIG_PATH}" != "${DEFAULT_CONFIG_PATH}" && "${CONFIG_PATH}" != "${SAVED_CONFIG_PATH}" ]]; then
    echo "检测到配置文件切换，重新开始新的观察会话。"
    init_state
    return
  fi
  CONFIG_PATH="${SAVED_CONFIG_PATH:-${CONFIG_PATH}}"
  save_state
}

current_step() {
  local index="$1"
  if (( index < 1 || index > ${#STEPS[@]} )); then
    return 1
  fi
  printf '%s\n' "${STEPS[$((index - 1))]}"
}

print_status() {
  if [[ ! -f "${STATE_FILE}" ]]; then
    echo "当前还没有观察进度。下一步将从第 1 步开始。"
    echo "下一步: 01/${#STEPS[@]} $(current_step 1)"
    echo "配置: ${CONFIG_PATH}"
    return
  fi

  load_state
  echo "当前会话: ${SESSION_ID}"
  echo "配置: ${SAVED_CONFIG_PATH:-${CONFIG_PATH}}"
  echo "状态文件: ${STATE_FILE}"
  echo "日志目录: ${LOG_ROOT}/${SESSION_ID}"

  if (( NEXT_INDEX > ${#STEPS[@]} )); then
    echo "进度: 已完成全部 ${#STEPS[@]} 步"
    echo "如需从头开始，请执行: ./scripts/run_next_step_observe.sh --reset"
    return
  fi

  echo "下一步: $(printf '%02d' "${NEXT_INDEX}")/${#STEPS[@]} $(current_step "${NEXT_INDEX}")"
}

if [[ ! -f "${CONFIG_PATH}" ]]; then
  echo "错误：未找到配置文件 ${CONFIG_PATH}" >&2
  exit 1
fi

case "${ACTION}" in
  status)
    print_status
    exit 0
    ;;
  reset)
    init_state
    echo "已重置逐步观察进度。"
    echo "当前会话: ${SESSION_ID}"
    echo "下一步: 01/${#STEPS[@]} $(current_step 1)"
    echo "日志目录: ${LOG_ROOT}/${SESSION_ID}"
    exit 0
    ;;
esac

ensure_state

if (( NEXT_INDEX > ${#STEPS[@]} )); then
  echo "全部步骤已经执行完成。"
  echo "如需重新开始，请执行: ./scripts/run_next_step_observe.sh --reset"
  exit 0
fi

STEP_NAME="$(current_step "${NEXT_INDEX}")"
SESSION_DIR="${LOG_ROOT}/${SESSION_ID}"
STEP_DIR="${SESSION_DIR}/$(printf '%02d' "${NEXT_INDEX}")-${STEP_NAME}"
CONSOLE_LOG="${STEP_DIR}/console.log"
META_FILE="${STEP_DIR}/meta.txt"

mkdir -p "${STEP_DIR}"

{
  echo "session=${SESSION_ID}"
  echo "index=$(printf '%02d' "${NEXT_INDEX}")"
  echo "step=${STEP_NAME}"
  echo "config=${CONFIG_PATH}"
  echo "started_at=$(date '+%Y-%m-%d %H:%M:%S')"
} > "${META_FILE}"

echo "========== NEXT STEP OBSERVE =========="
echo "会话: ${SESSION_ID}"
echo "当前步骤: $(printf '%02d' "${NEXT_INDEX}")/${#STEPS[@]} ${STEP_NAME}"
echo "配置文件: ${CONFIG_PATH}"
echo "日志目录: ${STEP_DIR}"
echo

set +e
"${PROJECT_DIR}/scripts/run_single_step_observe.sh" "${STEP_NAME}" "${CONFIG_PATH}" | tee "${CONSOLE_LOG}"
STEP_EXIT=${PIPESTATUS[0]}
set -e

echo "finished_at=$(date '+%Y-%m-%d %H:%M:%S')" >> "${META_FILE}"
echo "exit_code=${STEP_EXIT}" >> "${META_FILE}"

if (( STEP_EXIT == 0 )); then
  NEXT_INDEX=$((NEXT_INDEX + 1))
  save_state
  echo
  echo "本次已执行完成，并已停下供你观察。"
  if (( NEXT_INDEX <= ${#STEPS[@]} )); then
    echo "下次继续执行: $(printf '%02d' "${NEXT_INDEX}")/${#STEPS[@]} $(current_step "${NEXT_INDEX}")"
  else
    echo "全部步骤都已执行完成。"
  fi
  echo "本次日志: ${STEP_DIR}"
  exit 0
fi

echo

echo "本次执行失败，进度未推进；修复后可直接重试当前步骤。" >&2
echo "当前仍停留在: $(printf '%02d' "${NEXT_INDEX}")/${#STEPS[@]} ${STEP_NAME}" >&2
echo "失败日志: ${STEP_DIR}" >&2
exit "${STEP_EXIT}"
