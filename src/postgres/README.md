# 云数据库 TencentDB for PostgreSQL

> 腾讯云数据库 PostgreSQL（TencentDB for PostgreSQL，云 API 使用 `postgres` 作为简称）能够让您在云端轻松设置、操作和扩展强大的开源数据库 PostgreSQL。

---

## Tools

当前默认注册 **48 个工具**，覆盖实例、账号、数据库、参数、备份、监控、网络、只读实例与 SSL 配置等能力。

---

## 当前版本定位

当前主分支已经补齐为 **三种 transport / 多种部署入口**，并新增 **`npx` 分发包装层**：

- **`npx`**：推荐的本地分发入口，默认拉起 `stdio`，不替代源码 / 自建服务方式
- **`stdio`**：本地命令直连，适合作为推荐本地启动方式
- **`streamable-http`**：当前主线默认，适合 Hosted URL / Docker / Compose / 远程主机 / SCF
- **`sse`**：自建兼容模式，适合本地或自建服务

默认仍建议：

- **鉴权**：`MCP_AUTH_MODE=request-credential`
- **transport**：`MCP_TRANSPORT=streamable-http`（服务化）或 `stdio`（本地命令式接入）

在 `request-credential` 模式下：

- HTTP / SSE：客户端每次请求通过 Header 携带腾讯云凭证
- stdio：本地进程通过环境变量向 server 注入腾讯云凭证

> 旧的 **托管 URL + token/OAuth/Authorization + SSE** 版本已按原样保存在：`variants/token-oauth-authorization-version`

---

## 快速开始

### 方式零：`npx` 启动（推荐本地分发入口）

这是新增的 **npm / npx 分发层**：默认拉起 `stdio`，适合作为 WorkBuddy / Cursor / Claude Desktop 一类本地客户端的命令入口，且**不会影响现有源码启动、自建 `sse` 或 `streamable-http` 部署方式**。

#### 1. 直接启动

```bash
npx -y postgres-mcp-server@latest
```

#### 2. MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "npx",
      "args": ["-y", "postgres-mcp-server@latest"],
      "env": {
        "TRANSPORT": "stdio",
        "TENCENTCLOUD_SECRET_ID": "<TENCENTCLOUD_SECRET_ID>",
        "TENCENTCLOUD_SECRET_KEY": "<TENCENTCLOUD_SECRET_KEY>",
        "TENCENTCLOUD_SESSION_TOKEN": "<TENCENTCLOUD_SESSION_TOKEN_OPTIONAL>"
      }
    }
  }
}
```

#### 3. 可选：用 `npx` 拉起自建 `sse`

```bash
npx -y postgres-mcp-server@latest --transport sse --env-file .env
```

#### 4. 说明

- `npx` 启动器会按平台下载 GitHub Release 里的预编译 Go 二进制
- 默认 `transport` 是 `stdio`；也兼容 `TRANSPORT -> MCP_TRANSPORT`、`PORT -> MCP_SERVER_PORT`
- 腾讯云凭证仍可直接使用 `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY`
- 如需本地调试或跳过下载，可设置 `POSTGRES_MCP_BINARY_PATH=/path/to/postgres-server`
- 发布预编译资产可执行：`./scripts/build_npx_release.sh`

---

### 方式一：推荐本地启动（`stdio`）

这是最接近 `CLS MCP` “本地命令直连” 的方式，适合 WorkBuddy / Cursor / Claude Desktop 一类本地客户端。

#### 1. 准备配置

```bash
cp .env.example .env
```

推荐最小配置：

```env
MCP_TRANSPORT=stdio
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_REQUEST_SECRET_ID=你的SecretId
MCP_REQUEST_SECRET_KEY=你的SecretKey
READ_ONLY=true
```

#### 2. 本地直接运行

```bash
./scripts/run_stdio.sh
```

#### 3. MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "./scripts/run_stdio.sh",
      "env": {
        "MCP_REQUEST_SECRET_ID": "<TENCENTCLOUD_SECRET_ID>",
        "MCP_REQUEST_SECRET_KEY": "<TENCENTCLOUD_SECRET_KEY>",
        "MCP_REQUEST_SESSION_TOKEN": "<TENCENTCLOUD_SESSION_TOKEN_OPTIONAL>"
      }
    }
  }
}
```

> `stdio` 模式仅适合本地可信环境，不适合作为远程共享服务暴露。

---

### 方式二：自建服务（推荐 `streamable-http`）

这是当前主线默认方案，适合 WorkBuddy Hosted URL、Docker、Compose、远程主机和 SCF。

#### 1. 准备配置

```bash
cp .env.example .env
```

推荐最小配置：

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
READ_ONLY=true
```

#### 2. 启动服务

```bash
./scripts/run_server.sh
```

#### 3. MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "http://127.0.0.1:9000/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<TENCENTCLOUD_SECRET_ID>",
        "X-TencentCloud-Secret-Key": "<TENCENTCLOUD_SECRET_KEY>",
        "X-TencentCloud-Session-Token": "<TENCENTCLOUD_SESSION_TOKEN_OPTIONAL>"
      }
    }
  }
}
```

#### 4. 直接调用示例

```bash
curl http://127.0.0.1:9000/mcp \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -H 'X-TencentCloud-Secret-Id: 你的SecretId' \
  -H 'X-TencentCloud-Secret-Key: 你的SecretKey' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","clientInfo":{"name":"demo-client","version":"1.0.0"},"capabilities":{}}}'
```

> Hosted URL、反向代理和 SCF 场景下，客户端 `url` 也必须指向**完整 MCP 端点**（例如 `https://你的函数URL/mcp`），不要只填根域名或函数根 URL。

---

### 方式三：自建 SSE 模式

如果你的客户端或调试链路更习惯 `SSE`，当前主线也已提供兼容入口。

#### 1. 配置

```env
MCP_TRANSPORT=sse
MCP_AUTH_MODE=request-credential
MCP_SERVER_SSE_ENDPOINT=/sse
MCP_SERVER_MESSAGE_ENDPOINT=/message
READ_ONLY=true
```

#### 2. 启动

```bash
./scripts/run_server.sh
```

#### 3. MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "sse",
      "url": "http://127.0.0.1:9000/sse",
      "headers": {
        "X-TencentCloud-Secret-Id": "<TENCENTCLOUD_SECRET_ID>",
        "X-TencentCloud-Secret-Key": "<TENCENTCLOUD_SECRET_KEY>",
        "X-TencentCloud-Session-Token": "<TENCENTCLOUD_SESSION_TOKEN_OPTIONAL>"
      }
    }
  }
}
```

> `SSE` 仅建议在本地或自建服务场景下使用；**不要**把它作为 SCF / 无状态网关的默认方案。

---

### 方式四：源码安装

如果你不想依赖脚本，也可以直接编译并按 transport 启动。

#### 1. 编译

```bash
go build -o ./.bin/postgres-server .
```

#### 2. 按 transport 启动

- **stdio**

```bash
MCP_TRANSPORT=stdio ./.bin/postgres-server
```

- **streamable-http**

```bash
MCP_TRANSPORT=streamable-http ./.bin/postgres-server
```

- **sse**

```bash
MCP_TRANSPORT=sse ./.bin/postgres-server
```

#### 3. 常用端点

- **streamable-http**：`/mcp`
- **SSE**：`/sse`
- **SSE message**：`/message`

---

## 其他部署方式

### Docker / Compose

当前容器默认固定 `MCP_TRANSPORT=streamable-http`。

```bash
docker compose up -d --build
```

或：

```bash
docker build -t mcp-server-postgres:latest .
docker run --rm -it --env-file .env -p 127.0.0.1:9000:9000 mcp-server-postgres:latest
```

### SCF Web 函数

SCF 仅推荐：

- `MCP_TRANSPORT=streamable-http`
- `MCP_STREAMABLE_HTTP_STATELESS=true`
- `MCP_AUTH_MODE=request-credential`

客户端连接时请使用 `https://你的函数URL/mcp`；函数根 URL 返回 `404 page not found` 属于预期。

详见：`SCF_DEPLOY.md`

### 远程主机 / 反向代理

建议：

- 服务监听 `0.0.0.0`
- 对外通过 HTTPS / 反向代理暴露
- `MCP_SERVER_PUBLIC_URL` 填客户端最终访问地址
- 起步保持 `READ_ONLY=true`

---

## 本仓库自带客户端如何连当前服务

### `cmd/mcp_smoke`

支持 `streamable-http` / `sse` / `stdio`：

```bash
SMOKE_TRANSPORT=streamable-http ./scripts/run_mcp_smoke.sh
SMOKE_TRANSPORT=sse ./scripts/run_mcp_smoke.sh
SMOKE_TRANSPORT=stdio ./scripts/run_mcp_smoke.sh
```

### `cmd/verify`

支持 `streamable-http` / `sse` / `stdio`：

```bash
VERIFY_TRANSPORT=streamable-http VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
VERIFY_TRANSPORT=sse VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
VERIFY_TRANSPORT=stdio VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
```

`cmd/mcp_smoke`、`cmd/verify`、`cmd/write_test` 已自动支持从环境变量读取凭证：

```bash
export MCP_REQUEST_SECRET_ID=你的SecretId
export MCP_REQUEST_SECRET_KEY=你的SecretKey
export MCP_REQUEST_SESSION_TOKEN=你的SessionToken # 可选
```

也兼容直接复用：

```bash
export MCP_SECRET_ID=你的SecretId
export MCP_SECRET_KEY=你的SecretKey
```

---

## 当前鉴权与授权行为

### 数据面鉴权

默认 `MCP_AUTH_MODE=request-credential`：

- HTTP / SSE：服务端从请求头提取腾讯云凭证
- stdio：服务端从当前进程环境读取腾讯云凭证
- 默认执行 `STS GetCallerIdentity`
- 认证成功后构造 `Principal`
- `tools/registry.go` 继续沿用现有 `Principal + Guard` 授权链路

### Scope 与地域

```env
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read,pg.write
MCP_REQUEST_ALLOWED_REGIONS=ap-guangzhou
```

说明：

- `pg.read`：允许只读工具
- `pg.write`：允许写类工具（仍受 `Guard` 和 `confirm=true` 约束）
- `MCP_REQUEST_ALLOWED_REGIONS` 留空表示不额外做地域限制

---

## 安全部署建议

- **生产环境必须放在 HTTPS / 反向代理之后**
- **不要把 `SecretId` / `SecretKey` 放进 URL 或 query 参数**
- **禁止在日志、trace、错误回显中输出凭据明文**
- **优先使用 STS 临时凭证，而不是长期 AK/SK**
- **保持 `READ_ONLY=true` 起步**，确认流程后再按需开放写操作
- **对外暴露前务必加 IP 白名单、VPN 或零信任访问控制**
- **SCF / API 网关等无状态环境只推荐 `streamable-http`**

---

## 相关文档

- 综合部署说明：`DEPLOY.md`
- SCF 部署：`SCF_DEPLOY.md`
- 本地 stdio 启动：`scripts/run_stdio.sh`
- 本地服务启动：`scripts/run_server.sh`
- WorkBuddy 使用问题记录：`WORKBUDDY_USAGE_NOTES.md`
