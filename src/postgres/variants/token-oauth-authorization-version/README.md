# 云数据库 TencentDB for PostgreSQL

> 腾讯云数据库 PostgreSQL（TencentDB for PostgreSQL，云 API 使用 `postgres` 作为简称）能够让您在云端轻松设置、操作和扩展强大的开源数据库 PostgreSQL。

---

## Tools

当前默认注册 **48 个工具**，覆盖实例、账号、数据库、参数、备份、监控、网络、只读实例与 SSL 配置等能力。

---

## 快速开始

### 推荐模式：`issued-token` + `TencentCloud Token Exchange`

当前默认推荐的接入方式是：

1. 服务端启动时**不再依赖固定 `MCP_SECRET_ID` / `MCP_SECRET_KEY`**
2. 用户调用 `POST /auth/bootstrap/tencentcloud`（底层仍复用 `POST /auth/token-exchange/tencentcloud`）
3. 服务端用用户提交的腾讯云身份材料调用 `STS GetCallerIdentity`
4. 服务端签发本地 `MCP access token`，并将运行时云凭证**加密绑定**到该 token
5. 服务端直接返回一份**可复制到 MCP 客户端**的配置 JSON
6. 客户端后续只携带 `Authorization: Bearer <MCP_ACCESS_TOKEN>` 访问 `/sse`

### 1. 启动服务

```bash
cp .env.example .env
```

最小建议配置：

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

然后启动：

```bash
./scripts/run_server.sh
```

> 当前实现里，`/sse`、`/message`、`/admin/tokens`、`/auth/token-exchange/tencentcloud`、`/auth/bootstrap/tencentcloud` 共用 `MCP_SERVER_PORT`。

### 2. 用户自助换 token

```bash
curl -X POST http://127.0.0.1:9000/auth/token-exchange/tencentcloud \
  -H 'Content-Type: application/json' \
  -d '{
    "secret_id": "你的腾讯云SecretId",
    "secret_key": "你的腾讯云SecretKey",
    "display_name": "张三开发机",
    "allowed_regions": ["ap-guangzhou"],
    "scopes": ["pg.read"],
    "expires_in_seconds": 3600,
    "description": "self-service token"
  }'
```

响应里会返回一次性明文 token，例如：

```json
{
  "id": "tok_xxx",
  "token": "mcp_ptk_xxx",
  "subject_id": "tencentcloud:10001:10001",
  "tenant_id": "tencentcloud-account:10001",
  "scopes": ["pg.read"],
  "identity": {
    "type": "User",
    "account_id": "10001",
    "user_id": "10001",
    "arn": "qcs::cam::uin/10001:uin/10001"
  }
}
```

### 3. 一键返回可用 MCP 配置（推荐）

如果你不想自己再拼客户端配置，可以直接调用：

```bash
curl -X POST http://127.0.0.1:9000/auth/bootstrap/tencentcloud \
  -H 'Content-Type: application/json' \
  -d '{
    "secret_id": "你的腾讯云SecretId",
    "secret_key": "你的腾讯云SecretKey",
    "display_name": "张三开发机",
    "allowed_regions": ["ap-guangzhou"],
    "scopes": ["pg.read"],
    "expires_in_seconds": 3600,
    "description": "bootstrap mcp config"
  }'
```

成功后除了返回 `token` 本身，还会额外返回：

```json
{
  "mcp_server_name": "mcp-server-postgres",
  "mcp_server_url": "http://127.0.0.1:9000/sse",
  "mcp_config": {
    "mcpServers": {
      "mcp-server-postgres": {
        "type": "sse",
        "url": "http://127.0.0.1:9000/sse",
        "headers": {
          "Authorization": "Bearer <MCP_ACCESS_TOKEN>"
        }
      }
    }
  },
  "mcp_config_json": "{...可直接复制到 MCP 客户端的 JSON...}"
}
```

### 4. 手工配置 MCP（兼容旧流程）

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

如果你使用本仓库里的本地验证命令，也可以直接导出：

```bash
export MCP_ACCESS_TOKEN=<MCP_ACCESS_TOKEN>
```

`cmd/mcp_smoke`、`cmd/verify`、`cmd/write_test` 会自动读取它并附带请求头。

---

## `source-credential` 与 `assume-role`

### `source-credential`（当前默认，开发/过渡用）

- 服务端先调用 `GetCallerIdentity` 验证你是谁
- 然后把你提交的腾讯云凭证**加密保存**到本地 store
- 后续工具调用按 token 动态解密并创建腾讯云 client

适合：

- 快速验证功能
- 尚未准备好统一角色授权的环境

### `assume-role`（推荐生产路径）

配置：

```env
MCP_TOKEN_EXCHANGE_MODE=assume-role
MCP_ROLE_ARN=qcs::cam::uin/xxx:roleName/xxx
MCP_ROLE_DURATION_SECONDS=3600
```

行为：

- 服务端仍先用用户提交的身份材料做 `GetCallerIdentity`
- 再调用 `STS AssumeRole`
- **只保存返回的临时凭证**，不保存用户源长期密钥

适合：

- 正式环境
- 需要进一步缩小长期凭证落地面的场景

---

## 管理面与兼容模式

### Admin API（可选）

如果你仍需要管理员代发 token：

- `POST /admin/tokens`
- `GET /admin/tokens`
- `GET /admin/tokens/{id}`
- `POST /admin/tokens/{id}/revoke`

管理面需通过：

```http
Authorization: Bearer <ADMIN_API_TOKEN>
```

### 兼容旧模式：`shared-token`

如果你暂时还没切到 token exchange，可以继续使用共享 token：

```env
MCP_AUTH_MODE=shared-token
MCP_API_TOKEN=请替换为高强度随机串
MCP_SECRET_ID=你的SecretId
MCP_SECRET_KEY=你的SecretKey
```

此时所有客户端共用同一个入口 token；推荐仅作为过渡方案。

### 本机临时调试：`none`

```env
MCP_AUTH_MODE=none
MCP_SECRET_ID=你的SecretId
MCP_SECRET_KEY=你的SecretKey
```

仅建议本机短期调试，**不要**用于共享环境或远程部署。

---

## 鉴权与授权说明

### 当前已落地

- **TencentCloud 自助换 token**：`POST /auth/token-exchange/tencentcloud`
- **TencentCloud 一键返回 MCP 配置**：`POST /auth/bootstrap/tencentcloud`
- **Admin API 发 token**：`POST /admin/tokens`
- **Admin API 查 / 吊销 token**：`GET /admin/tokens`、`GET /admin/tokens/{id}`、`POST /admin/tokens/{id}/revoke`
- **MCP 数据面 issued-token 校验**：`Authorization: Bearer <MCP_ACCESS_TOKEN>`
- **动态 CredentialProvider**：每次工具调用按 `token_id` 取回绑定云凭证，再创建腾讯云 client
- **两级 scope**：
  - `pg.read`：只允许 `LevelNone` 工具
  - `pg.write`：允许写类工具（仍继续受 `Guard` 和 `confirm=true` 约束）

### 当前仍保留的兼容点

- `MCP_SECRET_ID` / `MCP_SECRET_KEY` 仅在 `shared-token` / `none` 等兼容模式下继续使用
- 历史旧变量读取逻辑仍兼容，但新部署不再建议继续使用旧前缀
- `source-credential` 是为快速落地保留的过渡模式；更推荐最终切到 `assume-role`

---

## 安全部署建议

- **优先使用 `issued-token` + `assume-role`**
- **`CREDENTIAL_ENCRYPTION_KEY` 必须独立保管**，不要提交进 Git
- **`ADMIN_API_TOKEN` 不要和数据面 access token 复用**
- **`TOKEN_HASH_PEPPER` 建议在正式环境配置**，提升 token store 泄漏后的离线抗碰撞能力
- **保持 `READ_ONLY=true` 起步**，确认流程后再按需开放写操作
- **远程部署时不要裸露到公网**，请放在反向代理、访问控制、VPN 或零信任网络后面
- **限制地域范围**，结合 `SCOPE_ENABLED=true`、`REGION_SCOPE=ap-xxx` 以及 token 自身的 `allowed_regions`

---

## 常用验证命令

### 1. 用 access token 做协议联调

```bash
MCP_ACCESS_TOKEN=<MCP_ACCESS_TOKEN> go run ./cmd/mcp_smoke --url http://127.0.0.1:9000/sse --region ap-guangzhou
```

### 2. 用 access token 做真实只读云 API 验证

```bash
MCP_ACCESS_TOKEN=<MCP_ACCESS_TOKEN> go run ./cmd/verify --url http://127.0.0.1:9000/sse --region ap-guangzhou --instance-id postgres-xxxxxxxx
```

### 3. 查看已签发 token 列表

```bash
curl http://127.0.0.1:9000/admin/tokens \
  -H 'Authorization: Bearer <ADMIN_API_TOKEN>'
```

### 4. 吊销某个 token

```bash
curl -X POST http://127.0.0.1:9000/admin/tokens/<token-id>/revoke \
  -H 'Authorization: Bearer <ADMIN_API_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{"reason":"device lost"}'
```

---

## 部署与验证文档

- `DEPLOY.md`：本地 / Docker / 远程主机部署说明
- `SCF_DEPLOY.md`：腾讯云 SCF Web 函数部署与 zip 上传说明
- `scripts/build_scf_zip.sh`：构建 SCF Web 函数上传 zip
- `scripts/TEST_REPORT.md`：公开仓库可用的验证记录模板
- `scripts/full_test_plan.yaml`：写操作验证配置模板
- `scripts/full_test_plan.observe.yaml`：分步观察式写验证模板

---

## 参考链接

- 产品详情：`https://cloud.tencent.com/product/postgres`
- PostgreSQL API 文档：`https://cloud.tencent.com/document/product/409/16761`
- STS GetCallerIdentity：`https://cloud.tencent.com/document/product/1312/66098`
- 地域列表：`https://cloud.tencent.com/document/product/1596/77930`

---

## 许可证

MIT
