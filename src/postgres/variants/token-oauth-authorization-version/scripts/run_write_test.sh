#!/usr/bin/env bash
# PG MCP 写操作批量验证脚本（配置驱动）
# 用法：在 src/postgres 目录下执行
#   ./scripts/run_write_test.sh [配置文件路径]
# 默认配置：scripts/full_test_plan.yaml
#
# 说明：
#   1. 脚本会临时启用写操作（GUARD_TEST_READ_ONLY=false），退出后自动恢复 .env
#   2. 真实执行哪些写接口，完全由 YAML 中每个步骤的 enabled/approved/args 决定
#   3. 建议先编辑配置文件，把高风险/收费接口逐项审批后再运行

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
CONFIG_PATH="${1:-${PROJECT_DIR}/scripts/full_test_plan.yaml}"
if [[ $# -gt 0 ]]; then
  shift
fi
WRITE_ARGS=("$@")
BIN_DIR="$(mktemp -d)"
SERVER_BIN="${BIN_DIR}/postgres_server"
WRITE_BIN="${BIN_DIR}/pg_write"
LOG_FILE="${BIN_DIR}/server.log"
SERVER_PID=""
ENV_BACKUP="${PROJECT_DIR}/.env.backup.$(date +%s)"

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

cp "${PROJECT_DIR}/.env" "${ENV_BACKUP}"
echo "" >> "${PROJECT_DIR}/.env"
echo "# === 写操作测试临时配置 ===" >> "${PROJECT_DIR}/.env"
echo "GUARD_TEST_READ_ONLY=false" >> "${PROJECT_DIR}/.env"
echo "WRITE TEST: GUARD_TEST_READ_ONLY=false enabled (临时, 脚本退出后自动恢复)"

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

echo "==> 编译 server 与配置驱动写测客户端（目标平台: ${HOST_GOOS}/${HOST_GOARCH}）..."
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${SERVER_BIN}" .
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${WRITE_BIN}" ./cmd/write_test

echo "==> 加载 .env 并启动本地 server（写操作已放行）..."
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

echo "==> 运行配置驱动写测客户端..."
"${WRITE_BIN}" -config "${CONFIG_PATH}" "${WRITE_ARGS[@]}"

echo ""
echo "==> 写操作验证完成。"
