#!/usr/bin/env bash
# MCP 协议级 smoke test：启动本地 server，并用真实 MCP 客户端完成 initialize/tools/list/tools/call。
# 用法：在 src/postgres 目录下执行 ./scripts/run_mcp_smoke.sh
# 可选环境变量：
#   SMOKE_TRANSPORT=streamable-http|sse|stdio
#   SMOKE_REGION=ap-guangzhou
#   SMOKE_INSTANCE_ID=postgres-xxxxxxxx
#   SMOKE_LIST_LIMIT=12
#   SMOKE_SERVER_PORT=9000

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="$(mktemp -d)"
SERVER_BIN="${BIN_DIR}/postgres_server"
SMOKE_BIN="${BIN_DIR}/pg_mcp_smoke"
LOG_FILE="${BIN_DIR}/server.log"
SERVER_PID=""

cleanup() {
  if [[ -n "${SERVER_PID}" ]] && kill -0 "${SERVER_PID}" 2>/dev/null; then
    echo ""
    echo "==> 停止本地 server 进程 (PID ${SERVER_PID})"
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
  rm -rf "${BIN_DIR}"
}
trap cleanup EXIT

pick_free_port() {
  python3 - <<'PY'
import socket
s = socket.socket()
s.bind(("127.0.0.1", 0))
print(s.getsockname()[1])
s.close()
PY
}

normalize_transport() {
  local value="${1:-streamable-http}"
  case "${value}" in
    http|streamable_http|streamablehttp) printf 'streamable-http' ;;
    streamable-http|sse|stdio) printf '%s' "${value}" ;;
    *) printf 'streamable-http' ;;
  esac
}

if [[ ! -f "${PROJECT_DIR}/.env" ]]; then
  echo "错误：未找到 ${PROJECT_DIR}/.env，请先配置真实密钥后再运行。" >&2
  exit 1
fi

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

echo "==> 编译 server 与 smoke client（目标平台: ${HOST_GOOS}/${HOST_GOARCH}）..."
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${SERVER_BIN}" .
GOOS="${HOST_GOOS}" GOARCH="${HOST_GOARCH}" go build -o "${SMOKE_BIN}" ./cmd/mcp_smoke

echo "==> 加载 .env ..."
set -a
# shellcheck disable=SC1091
source .env
set +a

SMOKE_TRANSPORT="$(normalize_transport "${SMOKE_TRANSPORT:-${MCP_TRANSPORT:-streamable-http}}")"
SMOKE_REGION="${SMOKE_REGION:-ap-guangzhou}"
SMOKE_INSTANCE_ID="${SMOKE_INSTANCE_ID:-}"
SMOKE_LIST_LIMIT="${SMOKE_LIST_LIMIT:-12}"

if [[ "${SMOKE_TRANSPORT}" == "stdio" ]]; then
  echo "==> 以 stdio transport 运行 smoke client..."
  SMOKE_TRANSPORT="stdio" \
  SMOKE_STDIO_COMMAND="${SERVER_BIN}" \
  MCP_TRANSPORT="stdio" \
  SMOKE_REGION="${SMOKE_REGION}" \
  SMOKE_INSTANCE_ID="${SMOKE_INSTANCE_ID}" \
  SMOKE_LIST_LIMIT="${SMOKE_LIST_LIMIT}" \
  "${SMOKE_BIN}"
  echo ""
  echo "==> MCP smoke test 完成。"
  exit 0
fi

if [[ -z "${SMOKE_SERVER_PORT:-}" ]]; then
  SMOKE_SERVER_PORT="$(pick_free_port)"
else
  SMOKE_SERVER_PORT="${SMOKE_SERVER_PORT}"
fi

case "${SMOKE_TRANSPORT}" in
  sse)
    SMOKE_ENDPOINT="${MCP_SERVER_SSE_ENDPOINT:-/sse}"
    ;;
  *)
    SMOKE_ENDPOINT="${MCP_SERVER_HTTP_ENDPOINT:-/mcp}"
    ;;
esac
SMOKE_SERVER_URL="${SMOKE_SERVER_URL:-${SMOKE_SSE_URL:-http://127.0.0.1:${SMOKE_SERVER_PORT}${SMOKE_ENDPOINT}}}"

echo "==> 启动本地 ${SMOKE_TRANSPORT} server..."
MCP_TRANSPORT="${SMOKE_TRANSPORT}" MCP_SERVER_PORT="${SMOKE_SERVER_PORT}" "${SERVER_BIN}" > "${LOG_FILE}" 2>&1 &
SERVER_PID=$!

READY=0
for _ in $(seq 1 40); do
  if ! kill -0 "${SERVER_PID}" 2>/dev/null; then
    break
  fi
  if curl -fsS "http://127.0.0.1:${SMOKE_SERVER_PORT}/healthz" >/dev/null 2>&1; then
    READY=1
    break
  fi
  sleep 0.5
done

echo "----- server 启动日志 -----"
cat "${LOG_FILE}"
echo "---------------------------"

if [[ "${READY}" != "1" ]] || ! kill -0 "${SERVER_PID}" 2>/dev/null; then
  echo "错误：本地 server 未成功启动，请检查上方日志。" >&2
  exit 1
fi

echo "==> 运行 MCP smoke client..."
SMOKE_TRANSPORT="${SMOKE_TRANSPORT}" \
SMOKE_REGION="${SMOKE_REGION}" \
SMOKE_INSTANCE_ID="${SMOKE_INSTANCE_ID}" \
SMOKE_SERVER_URL="${SMOKE_SERVER_URL}" \
SMOKE_LIST_LIMIT="${SMOKE_LIST_LIMIT}" \
"${SMOKE_BIN}"

echo ""
echo "==> MCP smoke test 完成。"
