#!/usr/bin/env bash
# PG MCP 工具 OpenAPI 参数对齐校验脚本
# 用法：在 src/postgres 目录下执行 ./scripts/run_openapi_param_check.sh
#
# 说明：
#   1. 仅做本地参数契约校验，不启动 server，不调用腾讯云接口
#   2. 会覆盖全部已支持工具，验证 MCP 入参经兼容转换后能否映射到 SDK Request

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "${PROJECT_DIR}"
go run ./cmd/openapi_param_check
