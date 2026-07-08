#!/usr/bin/env bash
# 单步执行写接口，并在执行后立即做一轮只读回看。
# 用法：./scripts/run_single_step_observe.sh <StepName> [配置文件路径]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
STEP_NAME="${1:-}"
CONFIG_PATH="${2:-${PROJECT_DIR}/scripts/full_test_plan.yaml}"

if [[ -z "${STEP_NAME}" ]]; then
  echo "用法: ./scripts/run_single_step_observe.sh <StepName> [配置文件路径]" >&2
  exit 1
fi

cd "${PROJECT_DIR}"

echo "========== SINGLE STEP WRITE TEST =========="
echo "Step: ${STEP_NAME}"
echo "Config: ${CONFIG_PATH}"
echo

echo ">>> 执行写接口 ${STEP_NAME}"
./scripts/run_write_test.sh "${CONFIG_PATH}" -only "${STEP_NAME}"

echo
echo ">>> 执行后状态回看"
./scripts/run_verify.sh
