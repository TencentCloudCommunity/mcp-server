#!/usr/bin/env bash
# PG MCP 全量接口验证脚本
# 用法：
#   ./scripts/run_full_test.sh [配置文件路径]
#
# 流程：
#   1. 先执行只读接口验证（cmd/verify）
#   2. 再执行配置驱动的写接口验证（cmd/write_test）
#
# 默认写测配置：scripts/full_test_plan.yaml

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
CONFIG_PATH="${1:-${PROJECT_DIR}/scripts/full_test_plan.yaml}"

cd "${PROJECT_DIR}"

echo "==> [1/2] 执行只读接口验证"
./scripts/run_verify.sh

echo ""
echo "==> [2/2] 执行写接口配置驱动验证"
./scripts/run_write_test.sh "${CONFIG_PATH}"
