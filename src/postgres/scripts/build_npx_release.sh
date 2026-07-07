#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
OUT_DIR="${PROJECT_DIR}/dist/npx"
VERSION="${1:-}"

if [[ -z "${VERSION}" ]]; then
  if command -v node >/dev/null 2>&1; then
    VERSION="$(cd "${PROJECT_DIR}" && node -p "require('./package.json').version")"
  else
    echo "错误：未提供版本号，且本机缺少 node，无法从 package.json 读取版本。" >&2
    echo "用法：./scripts/build_npx_release.sh <version>" >&2
    exit 1
  fi
fi

BINARY_PREFIX="postgres-server"
CHECKSUM_FILE="checksums.txt"
PLATFORMS=(
  "darwin amd64"
  "darwin arm64"
  "linux amd64"
  "linux arm64"
  "windows amd64"
  "windows arm64"
)

rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}"

build_one() {
  local goos="$1"
  local goarch="$2"
  local stage_dir
  local binary_name="${BINARY_PREFIX}"
  local asset_name

  if [[ "${goos}" == "windows" ]]; then
    binary_name+=".exe"
    asset_name="${BINARY_PREFIX}_${VERSION}_${goos}_${goarch}.exe.gz"
  else
    asset_name="${BINARY_PREFIX}_${VERSION}_${goos}_${goarch}.gz"
  fi

  stage_dir="$(mktemp -d)"
  trap 'rm -rf "${stage_dir}"' RETURN

  echo "==> build ${goos}/${goarch}" >&2
  (
    cd "${PROJECT_DIR}"
    CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" go build -o "${stage_dir}/${binary_name}" .
  )

  gzip -c "${stage_dir}/${binary_name}" > "${OUT_DIR}/${asset_name}"
  chmod 644 "${OUT_DIR}/${asset_name}"
  rm -rf "${stage_dir}"
  trap - RETURN
}

for target in "${PLATFORMS[@]}"; do
  IFS=' ' read -r goos goarch <<< "${target}"
  build_one "${goos}" "${goarch}"
done

(
  cd "${OUT_DIR}"
  rm -f "${CHECKSUM_FILE}"
  for asset in *.gz; do
    shasum -a 256 "${asset}" >> "${CHECKSUM_FILE}"
  done
)

cat <<EOF
构建完成：${OUT_DIR}

建议发布到 GitHub Releases：
  tag: postgres-mcp-server-v${VERSION}
  assets:
    - ${CHECKSUM_FILE}
    - ${BINARY_PREFIX}_${VERSION}_darwin_amd64.gz
    - ${BINARY_PREFIX}_${VERSION}_darwin_arm64.gz
    - ${BINARY_PREFIX}_${VERSION}_linux_amd64.gz
    - ${BINARY_PREFIX}_${VERSION}_linux_arm64.gz
    - ${BINARY_PREFIX}_${VERSION}_windows_amd64.exe.gz
    - ${BINARY_PREFIX}_${VERSION}_windows_arm64.exe.gz
EOF
