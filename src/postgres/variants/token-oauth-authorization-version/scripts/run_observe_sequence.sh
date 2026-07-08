#!/usr/bin/env bash
# 按固定顺序逐个执行写接口，每步后立即运行一轮只读回看，并将日志落盘。
# 用法：./scripts/run_observe_sequence.sh [配置文件路径] [输出目录]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
CONFIG_PATH="${1:-${PROJECT_DIR}/scripts/full_test_plan.yaml}"
OUT_DIR="${2:-${PROJECT_DIR}/scripts/observe_logs/$(date +%Y%m%d-%H%M%S)}"
BIN_DIR="$(mktemp -d)"
SERVER_BIN="${BIN_DIR}/postgres_server"
WRITE_BIN="${BIN_DIR}/pg_write"
VERIFY_BIN="${BIN_DIR}/pg_verify"
LOG_FILE="${BIN_DIR}/server.log"
SERVER_PID=""
ENV_BACKUP="${PROJECT_DIR}/.env.observe.backup.$(date +%s)"

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

cleanup() {
  if [[ -n "${SERVER_PID}" ]] && kill -0 "${SERVER_PID}" 2>/dev/null; then
    echo ""
    echo "==> 停止本地 server 进程 (PID ${SERVER_PID})"
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
  if [[ -f "${ENV_BACKUP}" ]]; then
    echo "==> 恢复 .env 配置"
    cp "${ENV_BACKUP}" "${PROJECT_DIR}/.env"
    rm -f "${ENV_BACKUP}"
  fi
  rm -rf "${BIN_DIR}"
}
trap cleanup EXIT

if [[ ! -f "${PROJECT_DIR}/.env" ]]; then
  echo "错误：未找到 ${PROJECT_DIR}/.env" >&2
  exit 1
fi

if [[ ! -f "${CONFIG_PATH}" ]]; then
  echo "错误：未找到配置文件 ${CONFIG_PATH}" >&2
  exit 1
fi

mkdir -p "${OUT_DIR}"
SUMMARY_FILE="${OUT_DIR}/summary.tsv"
printf 'index\tstep\twrite_exit\tverify_exit\n' > "${SUMMARY_FILE}"

cp "${PROJECT_DIR}/.env" "${ENV_BACKUP}"
echo "" >> "${PROJECT_DIR}/.env"
echo "# === 逐步观测测试临时配置 ===" >> "${PROJECT_DIR}/.env"
echo "GUARD_TEST_READ_ONLY=false" >> "${PROJECT_DIR}/.env"

echo "==> 输出目录: ${OUT_DIR}"
echo "==> 使用配置文件: ${CONFIG_PATH}"
cd "${PROJECT_DIR}"

case "$(uname -s)" in
  Darwin) HOST_GOOS="darwin" ;;
  Linux)  HOST_GOOS="linux" ;;
  *) HOST_GOOS="$(uname -s | tr '[:upper:]' '[:lower:]')" ;;
esac
case "$(uname -m)" in
  arm64|aarch64) HOST_GOARCH="arm64" ;;
  x86_64|amd64)  HOST_GOARCH="amd64" ;;
  *) HOST_GOARCH="$(uname -m)" ;;
esac

echo "==> 编译 server / write / verify（目标平台: ${HOST_GOOS}/${HOST_GOARCH}）..."
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${SERVER_BIN}" .
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${WRITE_BIN}" ./cmd/write_test
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${VERIFY_BIN}" ./cmd/verify

echo "==> 启动本地 server（写操作已放行）..."
set -a
# shellcheck disable=SC1091
source "${PROJECT_DIR}/.env"
set +a
"${SERVER_BIN}" > "${LOG_FILE}" 2>&1 &
SERVER_PID=$!

for _ in $(seq 1 20); do
  if grep -q "Total tools registered" "${LOG_FILE}" 2>/dev/null; then
    break
  fi
  sleep 0.5
done

echo "----- server 启动日志 -----"
cat "${LOG_FILE}"
echo "---------------------------"

index=1
for step in "${STEPS[@]}"; do
  step_dir="${OUT_DIR}/$(printf '%02d' "${index}")-${step}"
  mkdir -p "${step_dir}"

  echo
  echo "========== STEP ${index}: ${step} =========="

  set +e
  "${WRITE_BIN}" -config "${CONFIG_PATH}" -only "${step}" | tee "${step_dir}/write.log"
  write_exit=${PIPESTATUS[0]}
  set -e

  echo
  echo "----- VERIFY AFTER ${step} -----" | tee "${step_dir}/verify.log"
  set +e
  "${VERIFY_BIN}" | tee -a "${step_dir}/verify.log"
  verify_exit=${PIPESTATUS[0]}
  set -e

  printf '%02d\t%s\t%s\t%s\n' "${index}" "${step}" "${write_exit}" "${verify_exit}" >> "${SUMMARY_FILE}"
  echo "----- STEP ${index} SUMMARY: write=${write_exit}, verify=${verify_exit} -----"

  index=$((index + 1))
done

echo
echo "==> 全部步骤执行完成，汇总文件: ${SUMMARY_FILE}"
