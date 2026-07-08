# 快速开始

下面按 **推荐度从高到低** 列出 4 种部署方式。请按需选择，并先看清各方式的前置条件：

- **方式一：腾讯云 SCF 自助托管（推荐云上部署）** —— 适合希望通过腾讯云托管并对外提供 HTTPS 访问的场景；需要您自行在 SCF 控制台完成函数创建、zip 上传和环境变量配置。
- **方式二：自建 `streamable-http` 服务** —— 需先按需拉取 `src/postgres` 目录（`git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git` → `cd mcp-server` → `git sparse-checkout set src/postgres` → `cd src/postgres`），并准备 **Go 1.25+** 用于本地编译运行。
- **方式三：本地 `stdio`** —— 需先按需拉取 `src/postgres` 目录（同上）并准备 **Go 1.25+**，适合本地 Cursor / Claude Desktop / WorkBuddy 客户端。
- **方式四：`npx` 一键拉起** —— 只需本机安装 **Node.js 18+**，无需克隆仓库，命令行一条即可。

## 方式一：腾讯云 SCF 自助托管（推荐云上部署）

> 适合希望运行在腾讯云并通过 HTTPS / 函数 URL 提供给团队共用的场景。仓库已提供 SCF 打包脚本、启动脚本和环境变量模板，您需要在自己的腾讯云账号下完成函数创建与发布。

### 1.1 按需拉取 `src/postgres` 目录并构建 SCF 发布包

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
./scripts/build_scf_zip.sh
```

默认会在 `dist/` 目录生成可上传到 SCF 的 zip 包。

### 1.2 在 SCF 控制台创建 Web 函数

进入 [SCF 云函数控制台](https://console.cloud.tencent.com/scf)，按以下方式创建：

- 函数类型：**Web 函数**
- 运行环境：**Go 标准运行环境**
- 代码上传方式：**本地上传 zip**

上传 zip 后，为函数开启公网访问 URL。

### 1.3 启动命令

zip 包已经内置 `scf_bootstrap`，通常直接使用包内启动文件即可。

如果控制台要求手动填写启动命令，请填与 `deploy/scf/scf.console.startup.sh` 一致的内容：

```bash
#!/bin/bash
set -euo pipefail

export PG_MCP_RUNTIME="${PG_MCP_RUNTIME:-scf}"
export PORT="${PORT:-9000}"
export MCP_TRANSPORT="${MCP_TRANSPORT:-streamable-http}"
export MCP_SERVER_BIND_HOST="${MCP_SERVER_BIND_HOST:-0.0.0.0}"
export MCP_SERVER_PORT="${MCP_SERVER_PORT:-${PORT}}"
export MCP_SERVER_HTTP_ENDPOINT="${MCP_SERVER_HTTP_ENDPOINT:-/mcp}"
export MCP_STREAMABLE_HTTP_STATELESS="${MCP_STREAMABLE_HTTP_STATELESS:-true}"
export MCP_AUTH_MODE="${MCP_AUTH_MODE:-request-credential}"
export MCP_REQUEST_VALIDATE_IDENTITY="${MCP_REQUEST_VALIDATE_IDENTITY:-true}"
export MCP_REQUEST_CREDENTIAL_SCOPES="${MCP_REQUEST_CREDENTIAL_SCOPES:-pg.read}"
export MCP_REQUEST_ALLOWED_REGIONS="${MCP_REQUEST_ALLOWED_REGIONS:-}"
export READ_ONLY="${READ_ONLY:-true}"
export TOKEN_EXCHANGE_ENABLED="${TOKEN_EXCHANGE_ENABLED:-false}"

exec /var/user/postgres-server
```

### 1.4 环境变量配置

最小推荐配置：

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
TOKEN_EXCHANGE_ENABLED=false
READ_ONLY=true
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
```

推荐补充：

```env
MCP_SERVER_PUBLIC_URL=https://您的函数URL/mcp
```

### 1.5 MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "https://您的函数URL/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<您的 SecretId>",
        "X-TencentCloud-Secret-Key": "<您的 SecretKey>"
      }
    }
  }
}
```

> ⚠️ 客户端 `url` 必须指向**完整 MCP 端点**（含 `/mcp` 后缀），不要使用函数根 URL。

### 1.5 快速验证

```bash
curl -i https://您的函数URL/healthz
```

正常返回 `200 OK` 即可按上面 JSON 配置客户端，开始调用 48 个 MCP 工具。

> 如需更完整的控制台操作说明，请查看 `SCF_DEPLOY.md`。

---

## 方式二：自建 `streamable-http` 服务

适合部署到自有云主机、内网服务器，通过域名给团队共用。

> **前置条件**：本机或云主机已安装 **Go 1.25+**，并具备外网出口（要访问腾讯云 OpenAPI）。

### 2.1 按需拉取 `src/postgres` 目录

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> 后续所有命令都需要在 `src/postgres/` 目录下执行。

### 2.2 准备配置

```bash
cp .env.example .env
```

最小推荐配置：

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
READ_ONLY=true
```

### 2.3 启动服务

```bash
./scripts/run_server.sh
```

### 2.4 MCP 客户端配置

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "http://127.0.0.1:9000/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<您的 SecretId>",
        "X-TencentCloud-Secret-Key": "<您的 SecretKey>"
      }
    }
  }
}
```

> 建议放在 HTTPS / 反向代理之后；公网暴露前务必加 IP 白名单 / 零信任访问控制。

---

## 方式三：本地 `stdio`（推荐本地客户端）

适合 Cursor、Claude Desktop、WorkBuddy 等本地 MCP 客户端的命令直连模式。

> **前置条件**：本机已安装 **Go 1.25+**，并具备外网出口（要访问腾讯云 OpenAPI）。

### 3.1 按需拉取 `src/postgres` 目录

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> 后续所有命令都需要在 `src/postgres/` 目录下执行。

### 3.2 准备配置

```bash
cp .env.example .env
```

```env
MCP_TRANSPORT=stdio
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_REQUEST_SECRET_ID=您的SecretId
MCP_REQUEST_SECRET_KEY=您的SecretKey
READ_ONLY=true
```

### 3.3 启动

```bash
./scripts/run_stdio.sh
```

### 3.4 MCP 客户端配置

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "/absolute/path/to/mcp-server/src/postgres/scripts/run_stdio.sh",
      "env": {
        "MCP_REQUEST_SECRET_ID": "<您的 SecretId>",
        "MCP_REQUEST_SECRET_KEY": "<您的 SecretKey>"
      }
    }
  }
}
```

> 推荐把 `command` 写成**绝对路径**。很多 MCP 客户端拉起 `stdio` 进程时，工作目录并不是仓库根目录；如果写成 `./scripts/run_stdio.sh`，很容易出现 `spawn ./scripts/run_stdio.sh ENOENT`。
>
> 如果客户端支持 `cwd`，也可以把 `cwd` 显式设为 `src/postgres` 后再使用相对路径。
>
> `stdio` 模式仅适合本地可信环境，不适合作为远程共享服务暴露。

---

## 方式四：`npx` 一键拉起（最简本地体验）

> **前置条件**：本机已安装 **Node.js 18+**（含 `npx`）。无需克隆 Go 仓库，npm 包会按平台自动从 GitHub Release 下载预编译二进制。
>
> 检查是否已安装：
>
> ```bash
> node -v   # 期望 v18.x 或更高
> ```

### 4.1 命令行直接启动

```bash
npx -y postgres-mcp-server@latest
```

### 4.2 MCP 客户端配置

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "npx",
      "args": ["-y", "postgres-mcp-server@latest"],
      "env": {
        "TRANSPORT": "stdio",
        "TENCENTCLOUD_SECRET_ID": "<您的 SecretId>",
        "TENCENTCLOUD_SECRET_KEY": "<您的 SecretKey>"
      }
    }
  }
}
```

> `npx` 启动器会按平台从 GitHub Release 下载预编译 Go 二进制，首次运行需要联网。
