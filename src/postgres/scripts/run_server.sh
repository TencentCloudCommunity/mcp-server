#!/usr/bin/env bash
# 本地启动 PostgreSQL MCP Server：自动加载 .env、编译并启动。

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BIN_DIR="${PROJECT_DIR}/.bin"
SERVER_BIN="${BIN_DIR}/postgres-server"

if [[ ! -f "${PROJECT_DIR}/.env" ]]; then
  echo "错误：未找到 ${PROJECT_DIR}/.env，请先执行 'cp .env.example .env' 并填写配置。" >&2
  exit 1
fi

trim_whitespace() {
  local value="$1"
  value="${value#"${value%%[![:space:]]*}"}"
  value="${value%"${value##*[![:space:]]}"}"
  printf '%s' "$value"
}

load_env_file() {
  local env_file="$1"
  local raw line key value quote_char
  while IFS= read -r raw || [[ -n "$raw" ]]; do
    line="$(trim_whitespace "$raw")"
    [[ -z "$line" || "${line:0:1}" == "#" ]] && continue
    if [[ "$line" == export[[:space:]]* ]]; then
      line="$(trim_whitespace "${line#export}")"
    fi
    [[ "$line" == *=* ]] || continue
    key="$(trim_whitespace "${line%%=*}")"
    value="$(trim_whitespace "${line#*=}")"
    if [[ ! "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
      echo "错误：.env 中存在非法变量名 '$key'。" >&2
      exit 1
    fi
    if [[ ${#value} -ge 2 ]]; then
      quote_char="${value:0:1}"
      if [[ "$quote_char" == "${value: -1}" && ( "$quote_char" == '"' || "$quote_char" == "'" ) ]]; then
        value="${value:1:${#value}-2}"
      fi
    fi
    if [[ -n "${!key+x}" ]]; then
      continue
    fi
    export "$key=$value"
  done < "$env_file"
}

has_value() {
  [[ -n "${1:-}" ]]
}

cd "${PROJECT_DIR}"
load_env_file "${PROJECT_DIR}/.env"

TRANSPORT="${MCP_TRANSPORT:-streamable-http}"
AUTH_MODE="${MCP_AUTH_MODE:-request-credential}"

case "${TRANSPORT}" in
  streamable-http|sse|stdio) ;;
  http|streamable_http|streamablehttp) TRANSPORT="streamable-http" ;;
  *)
    echo "错误：不支持的 MCP_TRANSPORT='${TRANSPORT}'，可选值：streamable-http | sse | stdio。" >&2
    exit 1
    ;;
esac

if [[ "${AUTH_MODE}" == "issued-token" && "${TRANSPORT}" == "stdio" ]]; then
  echo "错误：当前主线下 stdio 不支持 issued-token；请改用 request-credential/shared-token/none，或切换到 HTTP/SSE。" >&2
  exit 1
fi

if [[ "${AUTH_MODE}" == "shared-token" || "${AUTH_MODE}" == "none" ]]; then
  if ! has_value "${MCP_SECRET_ID:-${TENCENTCLOUD_SECRET_ID:-}}" || ! has_value "${MCP_SECRET_KEY:-${TENCENTCLOUD_SECRET_KEY:-}}"; then
    echo "错误：当前鉴权模式需要服务端静态云凭证，请设置 MCP_SECRET_ID/MCP_SECRET_KEY。" >&2
    exit 1
  fi
fi

if [[ "${AUTH_MODE}" == "request-credential" && "${TRANSPORT}" == "stdio" ]]; then
  if ! has_value "${MCP_REQUEST_SECRET_ID:-${MCP_SECRET_ID:-${TENCENTCLOUD_SECRET_ID:-}}}" || \
     ! has_value "${MCP_REQUEST_SECRET_KEY:-${MCP_SECRET_KEY:-${TENCENTCLOUD_SECRET_KEY:-}}}"; then
    echo "错误：stdio + request-credential 模式下，必须提供 MCP_REQUEST_SECRET_ID/MCP_REQUEST_SECRET_KEY，或回退 MCP_SECRET_ID/MCP_SECRET_KEY。" >&2
    exit 1
  fi
fi

echo "==> 启动 transport: ${TRANSPORT}" >&2
if [[ "${AUTH_MODE}" == "request-credential" ]]; then
  if [[ "${TRANSPORT}" == "stdio" ]]; then
    echo "==> stdio 模式将从进程环境读取腾讯云凭证（MCP_REQUEST_SECRET_ID / MCP_REQUEST_SECRET_KEY）。" >&2
  else
    echo "==> HTTP 请求需携带 X-TencentCloud-Secret-Id / X-TencentCloud-Secret-Key。" >&2
  fi
elif [[ -n "${MCP_API_TOKEN:-}" ]]; then
  echo "==> 已启用 MCP_API_TOKEN 入口鉴权" >&2
fi
if [[ "${AUTH_MODE}" == "issued-token" ]]; then
  echo "==> 已启用 issued-token 动态凭证模式" >&2
fi

mkdir -p "${BIN_DIR}"
echo "==> 编译本地可执行文件..." >&2
env -u GOOS -u GOARCH go build -o "${SERVER_BIN}" .

echo "==> 启动 PostgreSQL MCP Server..." >&2
exec "${SERVER_BIN}"
