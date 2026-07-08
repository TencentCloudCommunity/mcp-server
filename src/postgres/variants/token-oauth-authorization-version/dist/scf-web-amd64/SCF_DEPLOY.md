# PostgreSQL MCP Server 部署到腾讯云 SCF（Web 函数）

本文提供**可直接上传 zip** 的 SCF 部署方式，并已按当前项目的**多用户 `issued-token`** 形态收敛。

---

## 1. 方案结论

当前项目最适合使用 **SCF Web 函数 + `issued-token` + TencentCloud Token Exchange**，原因如下：

- 服务本身就是长期运行的 `HTTP + SSE` Server，适合 **Web 函数**
- 当前项目已经切到**多用户**模式，客户端应在运行时换取自己的 `MCP_ACCESS_TOKEN`
- `MCP_API_TOKEN` 只适用于单共享口令的 `shared-token` 模式，**不再适合作为当前 SCF 默认配置**
- 服务端在多用户模式下**不需要预置** `MCP_SECRET_ID` / `MCP_SECRET_KEY`

> 当前 SCF 方案默认不再在环境变量中放静态腾讯云主密钥，也不再预置 `MCP_API_TOKEN`。客户端应在运行时调用 `/auth/bootstrap/tencentcloud`（或底层 `/auth/token-exchange/tencentcloud`），用**自己的腾讯云凭证**换取短期 `MCP_ACCESS_TOKEN` 与可直接复制的 MCP 配置。

### 1.1 当前 SCF 方案的边界

当前仓库的 token store 仍是本地 `sqlite`；在 SCF 中通常落到 `/tmp`，因此有这些特点：

- `/tmp` 可写，但**不保证跨冷启动/跨实例持久化**
- 已签发 token 和绑定凭证可能随实例回收而失效
- 因此更适合**短 TTL** token，而不是超长生命周期 token

这不影响“**运行时生成 token**”的正确性，但意味着客户端需要具备**重新换 token** 的能力。若后续要做更稳的生产化多用户部署，建议补一个**外部持久化 token store**。

---

## 2. 已准备好的文件

仓库内已补齐以下 SCF 交付物：

- `deploy/scf/scf_bootstrap`：SCF Web 函数启动文件
- `deploy/scf/scf.console.startup.sh`：控制台可直接粘贴的启动命令内容
- `deploy/scf/scf.env.example`：SCF 推荐环境变量模板
- `deploy/scf/scf.console.env.txt`：控制台可直接录入的环境变量清单
- `scripts/build_scf_zip.sh`：构建 Linux zip 包脚本

生成 zip：

```bash
./scripts/build_scf_zip.sh
```

如需 ARM：

```bash
./scripts/build_scf_zip.sh arm64
```

默认输出：

```bash
dist/postgres-mcp-scf-web-linux-amd64.zip
```

---

## 3. 在控制台创建函数

进入你给的控制台入口后，建议按下面选：

- **函数类型**：Web 函数
- **运行环境**：Go 标准运行环境
- **代码上传方式**：本地上传 zip
- **架构**：与 zip 一致（默认 `amd64`）

### 3.1 启动命令怎么填

当前 zip **已经包含** `scf_bootstrap`，正常情况下 SCF 会直接使用包内这个启动文件。

如果控制台的“启动命令”区域**要求你必须填写**，请直接填与 `deploy/scf/scf.console.startup.sh` 相同的内容：

```bash
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
```

> 这份内容与 zip 包内的 `scf_bootstrap` 一致；如果平台以上传包内文件为准，就会优先执行包里的 `scf_bootstrap`。

---

## 4. 控制台环境变量

建议把 `deploy/scf/scf.console.env.txt` 里的变量直接录入到 SCF 控制台环境变量；**多用户最小可用配置**如下：

```env
MCP_AUTH_MODE=issued-token
TOKEN_EXCHANGE_ENABLED=true
MCP_TOKEN_EXCHANGE_MODE=source-credential
TOKEN_EXCHANGE_ALLOWED_SCOPES=pg.read
READ_ONLY=true
CREDENTIAL_ENCRYPTION_KEY=替换为32字节AES密钥的base64结果
TOKEN_HASH_PEPPER=替换为高强度随机串
TOKEN_STORE=sqlite
TOKEN_STORE_PATH=/tmp/postgres-mcp/tokens.db
TOKEN_DEFAULT_TTL_SECONDS=3600
TOKEN_MAX_TTL_SECONDS=43200
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
```

可选但推荐：

```env
MCP_SERVER_PUBLIC_URL=https://你的函数URL/sse
MCP_STS_REGION=ap-guangzhou
SCOPE_ENABLED=true
REGION_SCOPE=ap-guangzhou
FEATURES=instance,account,database,parameter,backup,monitoring,network,readonly
```

### 4.1 不要再配置这些变量

以下变量**不要**再放到 SCF 环境变量里：

- `MCP_API_TOKEN`
- `MCP_SECRET_ID`
- `MCP_SECRET_KEY`

原因：

- `MCP_API_TOKEN` 是 `shared-token` 的单共享口令，不符合当前多用户模式
- `MCP_SECRET_ID` / `MCP_SECRET_KEY` 属于服务端静态云凭证，多用户 `issued-token` 流程不需要它们预置在 SCF 环境变量中

### 4.2 推荐的随机值生成方式

可以本地生成下面两个服务端随机值后，填进 SCF 控制台：

```bash
openssl rand -base64 32   # 用于 CREDENTIAL_ENCRYPTION_KEY
openssl rand -hex 32      # 用于 TOKEN_HASH_PEPPER
```

### 4.3 更推荐的最小权限收敛：`assume-role`

默认 `source-credential` 模式已经满足“**不在 SCF 环境变量里放用户密钥**”这个要求；如果你希望服务端绑定的是**更短期、权限更收敛的 CAM 临时凭证**，建议进一步切到：

```env
MCP_TOKEN_EXCHANGE_MODE=assume-role
MCP_ROLE_ARN=qcs::cam::uin/xxx:roleName/xxx
MCP_ROLE_DURATION_SECONDS=3600
MCP_ROLE_SESSION_PREFIX=pg-mcp
```

这样用户提交自己的腾讯云凭证后，服务端会先校验身份，再通过 `AssumeRole` 绑定更短期的角色临时凭证。

### 4.4 `ADMIN_API_TOKEN` 不是必填

`ADMIN_API_TOKEN` 只在你需要启用：

- `POST /admin/tokens`
- `GET /admin/tokens`
- `POST /admin/tokens/{id}/revoke`

这类管理员手工签发 / 查询 / 吊销 token 的接口时才需要。**纯多用户自助换 token 场景下，它不是必填项。**

---

## 5. 客户端如何获取运行时 token

部署成功并开启函数 URL 公网访问后，客户端优先调用：

```text
https://你的函数URL/auth/bootstrap/tencentcloud
```

如果只想拿裸 token，也可以继续调用：

```text
https://你的函数URL/auth/token-exchange/tencentcloud
```

> 这里提交的是**用户自己的腾讯云凭证**，属于**运行时请求体**，不是 SCF 环境变量。

示例：

```bash
curl -X POST 'https://你的函数URL/auth/token-exchange/tencentcloud' \
  -H 'Content-Type: application/json' \
  -d '{
    "secret_id": "用户自己的SecretId",
    "secret_key": "用户自己的SecretKey",
    "session_token": "如使用临时密钥则填写，否则可省略",
    "scopes": ["pg.read"],
    "expires_in_seconds": 3600,
    "description": "scf runtime exchange"
  }'
```

成功后会返回：

```json
{
  "token": "mcp_ptk_xxx",
  "expires_at": "2026-07-07T17:00:00Z",
  "credential_kind": "source-credential"
}
```

客户端再把返回的 token 作为 `MCP_ACCESS_TOKEN` 使用。

### 5.1 MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "sse",
      "url": "https://你的函数URL/sse",
      "headers": {
        "Authorization": "Bearer <MCP_ACCESS_TOKEN>"
      }
    }
  }
}
```

---

## 6. 部署后访问

部署成功并开启函数 URL 公网访问后，探活和 MCP 入口如下：

- **健康检查**：`https://你的函数URL/healthz`
- **就绪检查**：`https://你的函数URL/readyz`
- **MCP SSE**：`https://你的函数URL/sse`
- **腾讯云换 token**：`https://你的函数URL/auth/token-exchange/tencentcloud`

---

## 7. 验证建议

### 7.1 基础连通

先访问健康检查地址，预期返回 `200 OK`：

```text
https://你的函数URL/healthz
```

返回体示例：

```json
{
  "status": "ok",
  "probe": "live",
  "service": "mcp-server-postgres",
  "version": "1.0.0"
}
```

函数日志中也会出现：

```text
SSE server listening on 0.0.0.0:9000
Health check endpoints enabled on /healthz and /readyz
TencentCloud token exchange enabled on /auth/token-exchange/tencentcloud (source-credential)
```

### 7.2 换取访问 token

先调用 `/auth/token-exchange/tencentcloud` 获取 `MCP_ACCESS_TOKEN`，再执行下面的联调。

### 7.3 MCP 协议验证

```bash
MCP_ACCESS_TOKEN=<MCP_ACCESS_TOKEN> go run ./cmd/mcp_smoke --url https://你的函数URL/sse --region ap-guangzhou
```

### 7.4 真实只读能力验证

```bash
MCP_ACCESS_TOKEN=<MCP_ACCESS_TOKEN> go run ./cmd/verify --url https://你的函数URL/sse --region ap-guangzhou --instance-id postgres-xxxxxxxx
```

---

## 8. 注意事项

- **多用户模式下，`MCP_ACCESS_TOKEN` 应由服务在运行时生成，不要提前固定到环境变量中**
- **SCF 环境变量里不要再放 `MCP_SECRET_ID` / `MCP_SECRET_KEY` / `MCP_API_TOKEN`**
- **优先只开 `pg.read` + `READ_ONLY=true`**，先把只读链路跑通
- 当前 token store 仍在 `/tmp`，建议客户端具备**自动重新换 token** 的能力
- 如需更稳的生产化多用户部署，建议补**外部持久化 token store**
