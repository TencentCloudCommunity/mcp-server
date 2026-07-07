# PostgreSQL MCP Server 部署指南

本文面向希望**自行部署并接入 MCP 客户端**的用户，覆盖本地运行、Docker 部署、`issued-token` 动态凭证、管理面说明与安全建议。

---

## 一、部署前准备

### 1. 需要准备的内容

- **动态凭证模式（推荐）**
  - `CREDENTIAL_ENCRYPTION_KEY`
  - 可选：`ADMIN_API_TOKEN`
  - 用户自己的腾讯云身份材料（`SecretId` / `SecretKey` / 可选 `SessionToken`）
- **兼容静态模式（shared-token / none）**
  - 服务端固定 `MCP_SECRET_ID` / `MCP_SECRET_KEY`
- **运行环境**（二选一）
  - 本地运行：Go 1.24+
  - 容器运行：Docker 或 Docker Compose
- **MCP 客户端**：支持 `SSE` 类型 MCP Server 的客户端

### 2. 建议的最小安全配置

首次部署建议使用以下策略：

- **`MCP_AUTH_MODE=issued-token`**：默认走服务端发 token
- **`CREDENTIAL_ENCRYPTION_KEY=<32字节随机密钥>`**：加密保存绑定凭证
- **`TOKEN_EXCHANGE_ENABLED=true`**：启用 `/auth/token-exchange/tencentcloud`
- **`TOKEN_EXCHANGE_ALLOWED_SCOPES=pg.read`**：默认只开放只读
- **`READ_ONLY=true`**：先只开放只读能力
- **`MCP_SERVER_BIND_HOST=127.0.0.1`**：默认仅允许本机访问
- **`REGION_SCOPE=ap-xxx`**：如果只管理单一地域，建议进一步限制地域范围
- **不要直接把 MCP 端口暴露到公网**：如需远程访问，请放在反向代理、访问控制或 VPN 后面

---

## 二、方式一：本地直接运行

### 1. 复制配置模板

```bash
cp .env.example .env
```

### 2. 编辑 `.env`

> `scripts/run_server.sh` 现在按 **`KEY=VALUE` 配置文件** 语义读取 `.env`，不要求它必须是可被 `source` 的 shell 脚本；因此请尽量保持与 `.env.example` 一致的简单 `key=value` 格式。

#### 动态凭证推荐配置

```env
MCP_AUTH_MODE=issued-token
CREDENTIAL_ENCRYPTION_KEY=请替换成32字节随机值的base64结果
TOKEN_EXCHANGE_ENABLED=true
MCP_TOKEN_EXCHANGE_MODE=source-credential
TOKEN_EXCHANGE_ALLOWED_SCOPES=pg.read
ADMIN_API_TOKEN=请替换为高强度随机串
TOKEN_STORE=sqlite
TOKEN_STORE_PATH=./data/tokens.db
READ_ONLY=true
```

#### 兼容静态配置

```env
MCP_AUTH_MODE=shared-token
MCP_API_TOKEN=请替换为高强度随机串
MCP_SECRET_ID=你的SecretId
MCP_SECRET_KEY=你的SecretKey
READ_ONLY=true
```

### 3. 启动服务

```bash
./scripts/run_server.sh
```

启动成功后：

- `/healthz` 和 `/readyz` 作为匿名健康检查入口
- `/sse` 和 `/message` 作为 MCP 数据面入口
- `/auth/token-exchange/tencentcloud` 作为用户自助换 token 入口
- `/admin/tokens` 作为可选管理面入口
- 终端会打印 MCP 客户端配置模板

> 当前实现中，以上接口共用 `MCP_SERVER_PORT`。

### 4. 用户自助换 token（推荐）

```bash
curl -X POST http://127.0.0.1:9000/auth/token-exchange/tencentcloud \
  -H 'Content-Type: application/json' \
  -d '{
    "secret_id": "你的腾讯云SecretId",
    "secret_key": "你的腾讯云SecretKey",
    "display_name": "本机开发环境",
    "allowed_regions": ["ap-guangzhou"],
    "scopes": ["pg.read"],
    "expires_in_seconds": 3600,
    "description": "self-service token"
  }'
```

### 5. 客户端使用 token 连接

```json
{
 "mcpServers": {
  "mcp-server-postgres": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse",
   "headers": {
    "Authorization": "Bearer <MCP_ACCESS_TOKEN>"
   }
  }
 }
}
```

---

## 三、`source-credential` 与 `assume-role`

### `source-credential`

```env
MCP_TOKEN_EXCHANGE_MODE=source-credential
```

特点：

- 服务端验证调用者身份后，加密保存用户提交的凭证
- 最容易快速验证功能
- 适合作为本地开发和过渡阶段方案

### `assume-role`

```env
MCP_TOKEN_EXCHANGE_MODE=assume-role
MCP_ROLE_ARN=qcs::cam::uin/xxx:roleName/xxx
MCP_ROLE_DURATION_SECONDS=3600
```

特点：

- 服务端验证调用者身份后，调用 `STS AssumeRole`
- 只保存临时凭证，不保存源长期密钥
- 更适合正式环境

---

## 四、方式二：使用 Docker Compose

### 1. 准备环境变量

```bash
cp .env.example .env
```

### 2. 启动服务

```bash
docker compose up -d --build
```

### 3. 查看日志

```bash
docker compose logs -f
```

### 4. 停止服务

```bash
docker compose down
```

> 默认配置下，容器内服务监听 `0.0.0.0`，宿主机仍建议只在 `127.0.0.1:${MCP_SERVER_PORT:-9000}` 暴露端口。

---

## 五、方式三：只用 Docker 命令

```bash
docker build -t mcp-server-postgres:latest .
```

```bash
docker run --rm -it \
  --name postgres-mcp \
  --env-file .env \
  -e MCP_SERVER_BIND_HOST=0.0.0.0 \
  -p 127.0.0.1:9000:9000 \
  mcp-server-postgres:latest
```

---

## 六、远程主机部署说明

如果你要把它部署到云主机、内网服务器，或通过域名提供给团队使用，建议：

1. 服务端监听：

```env
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_AUTH_MODE=issued-token
CREDENTIAL_ENCRYPTION_KEY=请替换为高强度随机密钥
ADMIN_API_TOKEN=请替换为高强度随机串
```

2. 通过反向代理暴露，例如 `https://mcp.example.com/postgres/sse`

3. 将 `.env` 中的公开地址改为客户端实际访问地址：

```env
MCP_SERVER_PUBLIC_URL=https://mcp.example.com/postgres/sse
```

4. 管理面只允许内网访问；数据面走反向代理，并让客户端配置 `Authorization: Bearer <MCP_ACCESS_TOKEN>`

> `MCP_SERVER_PUBLIC_URL` 只影响终端里打印给客户端使用的 MCP 配置，不影响服务实际监听地址。

---

## 七、环境变量说明

| 变量 | 是否必填 | 说明 |
|---|---|---|
| `MCP_AUTH_MODE` | 否 | `none` / `shared-token` / `issued-token` |
| `CREDENTIAL_ENCRYPTION_KEY` | issued-token 推荐必填 | 加密保存绑定云凭证的 AES 密钥 |
| `TOKEN_EXCHANGE_ENABLED` | 否 | 是否启用 `/auth/token-exchange/tencentcloud` |
| `MCP_TOKEN_EXCHANGE_MODE` | 否 | `source-credential` / `assume-role` |
| `MCP_STS_REGION` | 否 | STS 调用地域，默认 `ap-guangzhou` |
| `TOKEN_EXCHANGE_ALLOWED_SCOPES` | 否 | 自助换 token 允许申请的 scopes，默认 `pg.read` |
| `MCP_ROLE_ARN` | assume-role 必填 | STS 角色 ARN |
| `MCP_ROLE_SESSION_PREFIX` | 否 | 角色会话名前缀 |
| `MCP_ROLE_DURATION_SECONDS` | 否 | AssumeRole 临时凭证时长 |
| `MCP_SECRET_ID` | shared-token / none 必填 | 静态云凭证模式下的 SecretId |
| `MCP_SECRET_KEY` | shared-token / none 必填 | 静态云凭证模式下的 SecretKey |
| 历史旧变量 | 否 | 仍兼容历史读取逻辑，但新部署不建议继续使用旧前缀 |
| `MCP_API_TOKEN` | shared-token 必填 | 共享 token 模式入口鉴权 token |
| `MCP_SHARED_TOKEN_SCOPES` | 否 | 共享 token 默认 scopes |
| `ADMIN_API_TOKEN` | 可选 | Admin API Bearer token |
| `TOKEN_STORE` | 否 | 当前支持 `sqlite` |
| `TOKEN_STORE_PATH` | 否 | token SQLite 文件路径，默认 `./data/tokens.db` |
| `TOKEN_HASH_PEPPER` | 否 | token 哈希附加 pepper |
| `TOKEN_DEFAULT_TTL_SECONDS` | 否 | 默认 token TTL |
| `TOKEN_MAX_TTL_SECONDS` | 否 | 最长 token TTL |
| `TOKEN_DEFAULT_SCOPES` | 否 | 管理面签发 token 时默认 scopes |
| `MCP_SERVER_BIND_HOST` | 否 | 服务监听地址，默认 `127.0.0.1` |
| `MCP_SERVER_PORT` | 否 | 服务端口，默认 `9000` |
| `MCP_SERVER_SSE_ENDPOINT` | 否 | SSE 路径，默认 `/sse` |
| `MCP_SERVER_MESSAGE_ENDPOINT` | 否 | Message 路径，默认 `/message` |
| `MCP_SERVER_PUBLIC_URL` | 否 | 打印给客户端使用的完整 SSE URL |
| `READ_ONLY` | 否 | 全局只读保护，默认建议 `true` |
| `SCOPE_ENABLED` | 否 | 是否启用地域范围限制 |
| `REGION_SCOPE` | 否 | 限定允许操作的地域 |

---

## 八、部署后验证

### 1. 验证 token exchange + 协议联通性

```bash
MCP_ACCESS_TOKEN=<MCP_ACCESS_TOKEN> go run ./cmd/mcp_smoke --url http://127.0.0.1:9000/sse --region ap-guangzhou
```

### 2. 验证真实只读云 API 调用

```bash
MCP_ACCESS_TOKEN=<MCP_ACCESS_TOKEN> go run ./cmd/verify --url http://127.0.0.1:9000/sse --region ap-guangzhou --instance-id postgres-xxxxxxxx
```

### 3. 验证参数是否与 OpenAPI 对齐

```bash
./scripts/run_openapi_param_check.sh
```

> `cmd/mcp_smoke`、`cmd/verify`、`cmd/write_test` 会自动读取环境变量 `MCP_ACCESS_TOKEN`，未设置时回退读取 `MCP_API_TOKEN`。

---

## 九、相关文档

- 快速入口：`README.md`
- SCF 部署：`SCF_DEPLOY.md`
- 本地启动脚本：`scripts/run_server.sh`
- SCF 打包脚本：`scripts/build_scf_zip.sh`
- 公开验证记录模板：`scripts/TEST_REPORT.md`
- 腾讯云 PostgreSQL API 文档：`https://cloud.tencent.com/document/product/409/16761`
- 腾讯云 STS GetCallerIdentity：`https://cloud.tencent.com/document/product/1312/66098`
