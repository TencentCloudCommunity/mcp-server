# PostgreSQL MCP Server 部署指南

> 当前主分支已经补齐 `stdio / streamable-http / sse` 三种 transport，并继续以 `request-credential` 作为默认鉴权模型。

本文按 **CLS MCP 风格的“方式”组织** 当前项目的部署与接入方式，同时保留远程主机和 SCF 等扩展部署方案。

---

## 一、部署矩阵

| 方式 | transport | 适用场景 | 推荐度 | 说明 |
|---|---|---|---|---|
| 方式零：NPX 启动 | `stdio`（默认）/ `sse`（可选） | 本地 MCP 客户端、零源码分发、本地桌面工具 | **最高** | 新增 npm 分发层，不替换现有脚本与源码部署 |
| 方式一：推荐本地启动 | `stdio` | 本地 MCP 客户端、个人开发、桌面工具 | **最高** | 最接近本地命令式接入体验 |
| 方式二：自建服务（推荐） | `streamable-http` | Hosted URL、远程主机、SCF | **最高** | 当前主线默认方案 |
| 方式三：自建 SSE 模式 | `sse` | 本地调试、自建兼容模式 | 中 | 仅建议本地或自建服务 |
| 方式四：源码安装 | `stdio / streamable-http / sse` | 自行编译、手工控制运行方式 | 高 | 适合想显式控制二进制和启动命令的人 |
| SCF Web 函数 | `streamable-http` | Serverless / 函数 URL | 高 | 仅推荐无状态 HTTP 模式 |

---

## 二、部署前准备

### 1. 需要准备的内容

- **腾讯云凭证**
  - `SecretId`
  - `SecretKey`
- **运行环境**
  - `npx` 运行：Node.js 18+
  - 本地源码运行：Go 1.25+
- **MCP 客户端**
  - 支持 `stdio`、`streamable-http` 或 `sse` 的 MCP 客户端

### 2. 当前推荐的最小安全配置

首次部署建议：

```env
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
READ_ONLY=true
MCP_SERVER_BIND_HOST=127.0.0.1
```

并遵循：

- 本地客户端优先用 **`stdio`**
- 服务化共享优先用 **`streamable-http`**
- 不要把密钥放在 URL、query、日志或错误回显中

---

## 三、方式零：NPX 启动（`npx`）

### 1. 直接启动

```bash
npx -y postgres-mcp-server@latest
```

### 2. MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "npx",
      "args": ["-y", "postgres-mcp-server@latest"],
      "env": {
        "TRANSPORT": "stdio",
        "TENCENTCLOUD_SECRET_ID": "<TENCENTCLOUD_SECRET_ID>",
        "TENCENTCLOUD_SECRET_KEY": "<TENCENTCLOUD_SECRET_KEY>"
      }
    }
  }
}
```

### 3. 拉起自建 `sse`（可选）

```bash
npx -y postgres-mcp-server@latest --transport sse --env-file .env
```

### 4. 说明

- `npx` 启动器只负责分发和拉起预编译二进制，不修改 Go 主程序行为
- 默认 `transport=stdio`，也兼容 `TRANSPORT` / `PORT` 这类更短的环境变量写法
- 如需跳过下载，可设置 `POSTGRES_MCP_BINARY_PATH=/path/to/postgres-server`
- 发布预编译资产可执行：`./scripts/build_npx_release.sh`
- 源码启动、`./scripts/run_stdio.sh`、`./scripts/run_server.sh`、自建 `sse` 与 `streamable-http` 部署方式保持不变

---

## 四、方式一：推荐本地启动（`stdio`）

### 1. 复制配置模板

```bash
cp .env.example .env
```

### 2. 编辑 `.env`

```env
MCP_TRANSPORT=stdio
MCP_AUTH_MODE=request-credential
MCP_REQUEST_SECRET_ID=您的SecretId
MCP_REQUEST_SECRET_KEY=您的SecretKey
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
READ_ONLY=true
```

### 3. 启动

```bash
./scripts/run_stdio.sh
```

### 4. 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "/absolute/path/to/mcp-server/src/postgres/scripts/run_stdio.sh",
      "env": {
        "MCP_REQUEST_SECRET_ID": "<TENCENTCLOUD_SECRET_ID>",
        "MCP_REQUEST_SECRET_KEY": "<TENCENTCLOUD_SECRET_KEY>"
      }
    }
  }
}
```

### 5. 说明

- 推荐把 `command` 写成**绝对路径**；很多 MCP 客户端启动 `stdio` 时不会把工作目录设到仓库根目录，使用 `./scripts/run_stdio.sh` 容易报 `spawn ./scripts/run_stdio.sh ENOENT`
- 如果客户端支持 `cwd`，也可以把 `cwd` 显式设为 `src/postgres` 后再使用相对路径
- `stdio` 不提供 `/healthz`、`/readyz` 这类 HTTP 探针
- `stdio` 仅适合本地可信环境
- 当前主线下 **`stdio` 不支持 `issued-token`**

---

## 五、方式二：自建服务（推荐 `streamable-http`）

### 1. 配置

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
READ_ONLY=true
```

### 2. 启动

```bash
./scripts/run_server.sh
```

### 3. 可访问地址

- 健康检查：`http://127.0.0.1:9000/healthz`
- 就绪检查：`http://127.0.0.1:9000/readyz`
- MCP：`http://127.0.0.1:9000/mcp`

### 4. MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "http://127.0.0.1:9000/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<TENCENTCLOUD_SECRET_ID>",
        "X-TencentCloud-Secret-Key": "<TENCENTCLOUD_SECRET_KEY>"
      }
    }
  }
}
```

### 5. 何时优先选它

- 团队共享 Hosted URL
- 需要健康检查、反向代理、HTTPS
- SCF / 远程主机

> **注意：** 客户端里的 `url` 必须填写**完整 MCP 端点**（例如 `http://127.0.0.1:9000/mcp`、`https://mcp.example.com/postgres/mcp`），不要只填域名或根路径。

---

## 六、方式三：自建 SSE 模式

### 1. 配置

```env
MCP_TRANSPORT=sse
MCP_AUTH_MODE=request-credential
MCP_SERVER_SSE_ENDPOINT=/sse
MCP_SERVER_MESSAGE_ENDPOINT=/message
READ_ONLY=true
```

### 2. 启动

```bash
./scripts/run_server.sh
```

### 3. 可访问地址

- 健康检查：`http://127.0.0.1:9000/healthz`
- SSE：`http://127.0.0.1:9000/sse`
- Message：`http://127.0.0.1:9000/message`

### 4. MCP 客户端配置示例

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "sse",
      "url": "http://127.0.0.1:9000/sse",
      "headers": {
        "X-TencentCloud-Secret-Id": "<TENCENTCLOUD_SECRET_ID>",
        "X-TencentCloud-Secret-Key": "<TENCENTCLOUD_SECRET_KEY>"
      }
    }
  }
}
```

### 5. 限制

- 不建议作为 SCF / API 网关默认模式
- 不建议作为多实例无状态环境的主链路

---

## 七、方式四：源码安装

### 1. 编译

```bash
go build -o ./.bin/postgres-server .
```

### 2. 启动命令

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

### 3. 说明

源码安装和脚本启动的本质区别只是：

- 脚本会自动读取 `.env`、编译并启动
- 源码安装由您自己控制二进制和环境变量注入方式

---

## 八、远程主机部署说明

如果您要把它部署到云主机、内网服务器，或通过域名提供给团队使用，建议：

1. 服务监听：

```env
MCP_TRANSPORT=streamable-http
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_AUTH_MODE=request-credential
MCP_STREAMABLE_HTTP_STATELESS=true
READ_ONLY=true
```

2. 通过反向代理暴露，例如：`https://mcp.example.com/postgres/mcp`

3. 配置：

```env
MCP_SERVER_PUBLIC_URL=https://mcp.example.com/postgres/mcp
```

4. 只通过 HTTPS + Header 传递腾讯云凭证，不要拼到 URL 中

---

## 十、SCF Web 函数部署说明

SCF 仅推荐：

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_STREAMABLE_HTTP_STATELESS=true
TOKEN_EXCHANGE_ENABLED=false
READ_ONLY=true
```

原因：

- Web 函数 / 无状态环境不适合默认使用 `SSE`
- `stdio` 不适合云函数部署
- 当前主线默认的 `/mcp` 单端点更适合 Hosted URL + Header 直传凭证

客户端连接时请始终使用完整端点：`https://您的函数URL/mcp`。
函数根 URL `https://您的函数URL` 返回 `404 page not found` 属于预期，不代表部署失败。

详见：`SCF_DEPLOY.md`

---

## 十、环境变量要点

| 变量 | 说明 |
|---|---|
| `MCP_TRANSPORT` | `streamable-http` / `sse` / `stdio` |
| `MCP_AUTH_MODE` | `request-credential` / `none` / `shared-token` / `issued-token` |
| `MCP_SERVER_HTTP_ENDPOINT` | streamable-http 路径，默认 `/mcp` |
| `MCP_SERVER_SSE_ENDPOINT` | SSE 路径，默认 `/sse` |
| `MCP_SERVER_MESSAGE_ENDPOINT` | SSE message 路径，默认 `/message` |
| `MCP_STREAMABLE_HTTP_STATELESS` | 是否无状态运行 streamable-http，默认建议 `true` |
| `MCP_REQUEST_SECRET_ID` / `MCP_REQUEST_SECRET_KEY` | stdio 模式或本地验证脚本可直接使用的腾讯云凭证 |
| `MCP_SERVER_PUBLIC_URL` | 对外暴露后的最终 MCP 地址 |

---

## 十一、部署后验证

### 1. 本地协议冒烟

```bash
SMOKE_TRANSPORT=streamable-http ./scripts/run_mcp_smoke.sh
SMOKE_TRANSPORT=sse ./scripts/run_mcp_smoke.sh
SMOKE_TRANSPORT=stdio ./scripts/run_mcp_smoke.sh
```

### 2. 真实只读云 API 验证

```bash
VERIFY_TRANSPORT=streamable-http VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
VERIFY_TRANSPORT=sse VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
VERIFY_TRANSPORT=stdio VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
```

### 3. 参数契约检查

```bash
./scripts/run_openapi_param_check.sh
```

---

## 十二、相关文档

- 快速入口：`README.md`
- SCF 部署：`SCF_DEPLOY.md`
- 本地 stdio 启动：`scripts/run_stdio.sh`
- 通用服务启动：`scripts/run_server.sh`
- WorkBuddy 使用问题记录：`WORKBUDDY_USAGE_NOTES.md`
