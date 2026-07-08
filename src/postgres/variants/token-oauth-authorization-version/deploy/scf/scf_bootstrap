#!/bin/bash
set -euo pipefail

export PG_MCP_RUNTIME="${PG_MCP_RUNTIME:-scf}"
export PORT="${PORT:-9000}"
export MCP_SERVER_BIND_HOST="${MCP_SERVER_BIND_HOST:-0.0.0.0}"
export MCP_SERVER_PORT="${MCP_SERVER_PORT:-${PORT}}"
export MCP_AUTH_MODE="${MCP_AUTH_MODE:-issued-token}"
export TOKEN_EXCHANGE_ENABLED="${TOKEN_EXCHANGE_ENABLED:-true}"
export MCP_TOKEN_EXCHANGE_MODE="${MCP_TOKEN_EXCHANGE_MODE:-source-credential}"
export TOKEN_STORE="${TOKEN_STORE:-sqlite}"
export TOKEN_STORE_PATH="${TOKEN_STORE_PATH:-/tmp/postgres-mcp/tokens.db}"
export TOKEN_DEFAULT_TTL_SECONDS="${TOKEN_DEFAULT_TTL_SECONDS:-3600}"
export TOKEN_MAX_TTL_SECONDS="${TOKEN_MAX_TTL_SECONDS:-43200}"
export READ_ONLY="${READ_ONLY:-true}"

mkdir -p "$(dirname "${TOKEN_STORE_PATH}")"

exec /var/user/postgres-server
