#!/usr/bin/env bash
# 推荐的本地 stdio 启动入口，对齐 CLS 风格的“本地命令直连”方式。

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export MCP_TRANSPORT="${MCP_TRANSPORT:-stdio}"

exec "${SCRIPT_DIR}/run_server.sh"
