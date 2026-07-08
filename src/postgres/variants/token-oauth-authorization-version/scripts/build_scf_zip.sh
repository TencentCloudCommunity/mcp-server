#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
ARCH="${1:-amd64}"

case "${ARCH}" in
  amd64|arm64)
    ;;
  *)
    echo "错误：仅支持 amd64 或 arm64，当前为 ${ARCH}" >&2
    exit 1
    ;;
esac

if ! command -v zip >/dev/null 2>&1; then
  echo "错误：未找到 zip 命令，请先安装 zip。" >&2
  exit 1
fi

DIST_DIR="${PROJECT_DIR}/dist"
STAGE_DIR="${DIST_DIR}/scf-web-${ARCH}"
OUTPUT_ZIP="${DIST_DIR}/postgres-mcp-scf-web-linux-${ARCH}.zip"
BINARY_PATH="${STAGE_DIR}/postgres-server"
BOOTSTRAP_SOURCE="${PROJECT_DIR}/deploy/scf/scf_bootstrap"
ENV_EXAMPLE_SOURCE="${PROJECT_DIR}/deploy/scf/scf.env.example"
CONSOLE_STARTUP_SOURCE="${PROJECT_DIR}/deploy/scf/scf.console.startup.sh"
CONSOLE_ENV_SOURCE="${PROJECT_DIR}/deploy/scf/scf.console.env.txt"
README_SOURCE="${PROJECT_DIR}/SCF_DEPLOY.md"

rm -rf "${STAGE_DIR}"
mkdir -p "${STAGE_DIR}"

pushd "${PROJECT_DIR}" >/dev/null
CGO_ENABLED=0 GOOS=linux GOARCH="${ARCH}" go build -o "${BINARY_PATH}" .
popd >/dev/null

cp "${BOOTSTRAP_SOURCE}" "${STAGE_DIR}/scf_bootstrap"
cp "${ENV_EXAMPLE_SOURCE}" "${STAGE_DIR}/scf.env.example"
cp "${CONSOLE_STARTUP_SOURCE}" "${STAGE_DIR}/scf.console.startup.sh"
cp "${CONSOLE_ENV_SOURCE}" "${STAGE_DIR}/scf.console.env.txt"
cp "${README_SOURCE}" "${STAGE_DIR}/SCF_DEPLOY.md"

chmod 755 "${STAGE_DIR}/postgres-server" "${STAGE_DIR}/scf_bootstrap"

rm -f "${OUTPUT_ZIP}"
pushd "${STAGE_DIR}" >/dev/null
zip -qry "${OUTPUT_ZIP}" .
popd >/dev/null

cat <<EOF
SCF_ZIP=${OUTPUT_ZIP}
SCF_STAGE_DIR=${STAGE_DIR}
SCF_ARCH=${ARCH}
EOF
