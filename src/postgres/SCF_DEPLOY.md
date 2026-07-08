# PostgreSQL MCP Server 部署到腾讯云 SCF（Web 函数）

> **SCF 仅推荐 `streamable-http`**。当前主分支虽然已经支持 `stdio / sse / streamable-http` 三种 transport，但云函数场景请固定使用 `MCP_TRANSPORT=streamable-http`。

本文说明如何把当前主分支部署到腾讯云 SCF，并让 MCP 客户端直接连接云函数 URL。

> **重要：**
> - **MCP 客户端必须连接** `https://您的函数URL/mcp`
> - **函数根 URL** `https://您的函数URL` **返回 `404 page not found` 属于预期**，不代表部署失败
> - **健康检查请使用** `https://您的函数URL/healthz`

---

## 1. 方案结论

当前推荐的 SCF 形态是：

- **SCF 只托管 MCP Server**
- **transport 固定为 `streamable-http`**
- **客户端直接连接 `https://您的函数URL/mcp`**
- **客户端在每次请求 Header 中传自己的腾讯云凭证**
- **SCF 环境变量中不保存用户的 `SecretId/SecretKey`**
- **服务端启用 `MCP_STREAMABLE_HTTP_STATELESS=true`**

不推荐 `stdio` / `sse` 的原因：

- `stdio` 不适合函数 URL / Web 托管
- `SSE` 不是云函数和无状态网关的默认优选
- 当前主线的 `/mcp` 单端点更稳定，更适合 Hosted URL

### 1.1 Header 约定

客户端访问 `/mcp` 时，按请求附带：

- `X-TencentCloud-Secret-Id`
- `X-TencentCloud-Secret-Key`

> **不要把密钥放到 URL / query 里。** 只通过 `HTTPS` + Header 传递，并确保网关/日志不会记录这些 Header。

---

## 2. 已准备好的交付物

仓库里已经准备好以下 SCF 文件：

- `deploy/scf/scf_bootstrap`
- `deploy/scf/scf.console.startup.sh`
- `deploy/scf/scf.env.example`
- `deploy/scf/scf.console.env.txt`
- `scripts/build_scf_zip.sh`

默认打包命令：

```bash
./scripts/build_scf_zip.sh
```

默认会在 `dist/` 目录生成可上传到 SCF 的 zip 包。

---

## 3. SCF 控制台创建函数

进入 [SCF 云函数控制台](https://console.cloud.tencent.com/scf/list?rid=16&ns=default)，按下面方式创建：

- **创建方式**：请选择“从头开始”
- **函数类型**：Web 函数
- **运行环境**：Go 标准运行环境（Go 1）
- **代码上传方式**：本地上传 zip
- **环境变量设置**：参考下文第 5 节，也可以创建函数后再设置
- **其它设置**：按需选择；如需要公网访问，请在最后的“函数URL配置”中勾选启用公网访问

上传 zip 后，为函数开启公网访问 URL。

---

## 4. 启动命令

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

---

## 5. SCF 环境变量怎么填

建议直接参考 `deploy/scf/scf.console.env.txt`，最小可用配置如下：

```env
MCP_TRANSPORT : streamable-http
MCP_AUTH_MODE : request-credential
MCP_REQUEST_VALIDATE_IDENTITY : true
MCP_REQUEST_CREDENTIAL_SCOPES : pg.read
READ_ONLY : true
MCP_SERVER_BIND_HOST : 0.0.0.0
MCP_SERVER_PORT : 9000
MCP_SERVER_HTTP_ENDPOINT : /mcp
MCP_STREAMABLE_HTTP_STATELESS : true
```

推荐补充：

```env
MCP_SERVER_PUBLIC_URL=https://您的函数URL/mcp
```

如需进一步限制访问地域或功能范围，可再按需补充：

```env
MCP_REQUEST_ALLOWED_REGIONS=ap-guangzhou
FEATURES=instance,account,database,parameter,backup,monitoring,network,readonly
```

### 5.1 这些变量不要放进 SCF

以下变量**不要**配置到云函数环境变量中：

- `MCP_SECRET_ID`
- `MCP_SECRET_KEY`
- `MCP_REQUEST_SECRET_ID`
- `MCP_REQUEST_SECRET_KEY`
- `MCP_API_TOKEN`
- `MCP_ACCESS_TOKEN`

原因：

- 它们属于**用户侧或客户端侧凭据**
- 当前模式要求**按请求传递**，而不是提前固化在服务端环境中
- 放进 SCF 环境变量会扩大泄露面

---

## 6. 客户端如何连接云函数

部署成功后，MCP 客户端直接连：

```text
https://您的函数URL/mcp
```

并在 Header 中携带自己的腾讯云凭证。

> **不要把裸函数 URL** `https://您的函数URL` **直接配给 MCP 客户端。**
> 根路径通常会返回 `404 page not found`；这在当前 SCF Web 函数部署里是正常现象。

### 6.1 MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "https://您的函数URL/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "您的SecretId",
        "X-TencentCloud-Secret-Key": "您的SecretKey"
      }
    }
  }
}
```

---

## 7. 部署后可访问的地址

函数 URL 开通公网访问后，主要地址如下：

- **健康检查**：`https://您的函数URL/healthz`
- **就绪检查**：`https://您的函数URL/readyz`
- **MCP streamable-http**：`https://您的函数URL/mcp`
- **函数根 URL**：`https://您的函数URL`（通常返回 `404 page not found`，属预期行为）

---

## 8. 本地联调 / 远程验收

### 8.1 健康检查

```text
https://您的函数URL/healthz
```

预期返回 `200 OK`。

### 8.2 MCP 协议冒烟

```bash
MCP_REQUEST_SECRET_ID=您的SecretId \
MCP_REQUEST_SECRET_KEY=您的SecretKey \
go run ./cmd/mcp_smoke --transport streamable-http --url https://您的函数URL/mcp --region ap-guangzhou
```

### 8.3 真实只读能力验证

```bash
MCP_REQUEST_SECRET_ID=您的SecretId \
MCP_REQUEST_SECRET_KEY=您的SecretKey \
go run ./cmd/verify --transport streamable-http --url https://您的函数URL/mcp --region ap-guangzhou --instance-id postgres-xxxxxxxx
```

---

## 9. 常见误区 / 排障

- **根路径返回 `404`**：如果 `https://您的函数URL` 或 `https://您的函数URL/` 返回 `404 page not found`，通常只是因为您访问的不是 MCP 端点；请改连 `https://您的函数URL/mcp`。
- **`/mcp` 返回 `401`**：这通常表示服务在线，但当前请求没带 `X-TencentCloud-Secret-Id` / `X-TencentCloud-Secret-Key` 等鉴权 Header。
- **客户端里工具显示不全或没刷新**：优先检查是否还连着旧 URL、是否忘了加 `/mcp`，以及客户端是否缓存了旧的 tools schema；必要时删除后重新添加连接。

---

## 10. 安全建议

- **只用 HTTPS 暴露函数 URL**
- **不要在网关、CDN、日志平台记录鉴权 Header**
- **不要把密钥拼到 URL、query、日志或报错回显中**
- **建议使用最小权限的 CAM 子账号**，避免长期复用主账号密钥
- **先以 `pg.read` + `READ_ONLY=true` 起步**
- **无状态环境保持 `MCP_STREAMABLE_HTTP_STATELESS=true`**
