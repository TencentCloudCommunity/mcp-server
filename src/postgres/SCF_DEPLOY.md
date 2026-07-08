# PostgreSQL MCP Server 部署到腾讯云 SCF（Web 函数）

> 本文档与 `README.md` 中“**3.1 方式一：腾讯云 SCF 自助托管**”保持一致，面向希望将 PostgreSQL MCP Server 部署到腾讯云并通过 HTTPS / 函数 URL 提供给团队使用的场景。

> **部署结论先看这里：**
> - **SCF 场景固定使用** `streamable-http`
> - **客户端接入地址必须是** `https://您的函数URL/mcp`
> - **根路径** `https://您的函数URL` **返回 404 属于预期**，不代表部署失败
> - **健康检查地址是** `https://您的函数URL/healthz`
> - **推荐使用 per-request 凭证模式**，不要把 `SecretId` / `SecretKey` 固化到函数环境变量中

---

## 1. 部署前准备

在 SCF 部署前，请先确认以下信息：

### 1.1 适用场景

SCF 方式适合以下场景：

- 希望把 MCP Server 运行在腾讯云上
- 希望通过 HTTPS / 函数 URL 给团队成员共享使用
- 希望服务端不长期保存用户凭证，而是由客户端按请求携带密钥

### 1.2 所需权限

调用 MCP 工具时，需要一个具备 PostgreSQL API 访问权限的腾讯云账号。建议：

- 单独创建 CAM 子账号用于 MCP 调用
- 按最小权限原则授予 PostgreSQL 相关权限
- 提前确认实例所在地域，后续请求需要传入 `region`

相关入口：

- [API 密钥（CAM）管理控制台](https://console.cloud.tencent.com/cam/capi)
- [地域与可用区映射](https://cloud.tencent.com/document/product/1596/77930)

### 1.3 需要准备的资源

| 资源 | 是否必须 | 说明 |
|---|---|---|
| `SecretId` / `SecretKey` | 是 | 客户端按请求传递给服务端，用于调用腾讯云 OpenAPI |
| 地域代码（如 `ap-guangzhou`） | 是 | 作为工具调用参数传入 |
| 实例 ID | 否 | 仅实例级操作时需要 |

---

## 2. 构建 SCF 发布包

### 2.1 拉取 `src/postgres` 目录

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

### 2.2 生成 SCF zip 包

```bash
./scripts/build_scf_zip.sh
```

构建完成后，默认会在 `dist/` 目录生成可上传到 SCF 的 zip 包。

> 如果本地脚本没有执行权限，可以先执行 `bash ./scripts/build_scf_zip.sh`，或给脚本补可执行权限后再运行。

### 2.3 仓库内已提供的 SCF 交付物

仓库中已准备好以下文件，可直接配合使用：

- `deploy/scf/scf_bootstrap`
- `deploy/scf/scf.console.startup.sh`
- `deploy/scf/scf.env.example`
- `deploy/scf/scf.console.env.txt`
- `scripts/build_scf_zip.sh`

---

## 3. 在 SCF 控制台创建 Web 函数

进入 [SCF 云函数控制台](https://console.cloud.tencent.com/scf/list?rid=16&ns=default)，按下面方式创建函数：

- **创建方式**：从头开始
- **函数类型**：**Web 函数**
- **运行环境**：**Go 标准运行环境（Go 1）**
- **代码上传方式**：本地上传 zip
- **环境变量**：可先留空，创建后再补；也可按下文第 5 节一次性填写
- **公网访问**：如需给外部客户端访问，请启用函数 URL 的公网访问能力

上传步骤 2 生成的 zip 包后，记录函数 URL，后续客户端接入将用到：

```text
https://您的函数URL/mcp
```

> 注意：客户端必须连接带 `/mcp` 后缀的完整端点，而不是函数根 URL。

---

## 4. 启动命令

SCF 发布包已内置 `scf_bootstrap`，通常直接使用包内启动文件即可。

如果控制台要求手动填写启动命令，请与 `deploy/scf/scf.console.startup.sh` 保持一致：

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

## 5. 环境变量配置

### 5.1 最小推荐配置

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
READ_ONLY=true
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
```

### 5.2 推荐补充

```env
MCP_SERVER_PUBLIC_URL=https://您的函数URL/mcp
```

### 5.3 可选增强配置

如果你希望进一步收敛访问范围，可按需增加：

```env
MCP_REQUEST_ALLOWED_REGIONS=ap-guangzhou
FEATURES=instance,account,database,parameter,backup,monitoring,network,readonly
```

### 5.4 不要写入 SCF 环境变量的内容

以下变量**不建议**配置到云函数环境变量中：

- `MCP_SECRET_ID`
- `MCP_SECRET_KEY`
- `MCP_REQUEST_SECRET_ID`
- `MCP_REQUEST_SECRET_KEY`
- `MCP_API_TOKEN`
- `MCP_ACCESS_TOKEN`

原因：

- 这些属于**用户侧或客户端侧凭据**
- 当前推荐模式是**每次请求携带凭证**，而不是把密钥固化在服务端
- 把长期密钥放到函数环境变量会扩大泄露面

---

## 6. 客户端配置

部署完成后，客户端请直接连接：

```text
https://您的函数URL/mcp
```

并通过 Header 传递腾讯云凭证。

### 6.1 MCP 客户端配置示例

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

### 6.2 地址使用规则

- **正确的 MCP 端点**：`https://您的函数URL/mcp`
- **健康检查**：`https://您的函数URL/healthz`
- **函数根 URL**：`https://您的函数URL`（返回 `404 page not found` 通常属于预期）

> 不要把裸函数 URL 直接配置给 MCP 客户端，否则通常会因为路径不对而无法正常接入。

---

## 7. 部署后验证

### 7.1 健康检查

```bash
curl -i https://您的函数URL/healthz
```

预期返回 `200 OK`。确认健康检查通过后，再把 `https://您的函数URL/mcp` 配置到客户端。

### 7.2 接入验证重点

成功接入后，可重点确认：

- 客户端是否能连通 `https://您的函数URL/mcp`
- 是否能正常识别并刷新 MCP 工具列表
- 请求是否按 Header 正确携带 `X-TencentCloud-Secret-Id` / `X-TencentCloud-Secret-Key`
- 在 `READ_ONLY=true` 场景下，是否只暴露或优先使用只读能力

---

## 8. 常见问题

### 8.1 函数根 URL 返回 404

如果访问：

```text
https://您的函数URL
```

返回 `404 page not found`，通常并不表示部署失败，而是因为你访问的不是 MCP 端点。请改用：

```text
https://您的函数URL/mcp
```

### 8.2 `/mcp` 返回 401 或鉴权失败

这通常说明：

- Header 中没有携带 `X-TencentCloud-Secret-Id` / `X-TencentCloud-Secret-Key`
- 或者密钥不正确 / 权限不足

### 8.3 客户端没有刷新到最新工具

优先检查：

- 是否仍连接旧地址
- 是否忘记在 URL 后追加 `/mcp`
- 是否客户端缓存了旧的 schema，必要时可删除后重新添加连接

---

## 9. 安全建议

无论本地测试还是正式部署，都建议遵循以下规则：

- **只通过 HTTPS 暴露函数 URL**
- **不要把密钥放进 URL / query 参数**，只通过 Header 或环境变量传递
- **不要在日志、网关、CDN、trace 或错误回显中输出凭据明文**
- **优先使用最小权限 CAM 子账号**
- **默认保持** `READ_ONLY=true`
- **无状态环境保持** `MCP_STREAMABLE_HTTP_STATELESS=true`
- **生产环境务必放在 HTTPS / 反向代理之后**，不要裸跑 HTTP 到公网

---

## 10. 相关文档

- [`README.md`](./README.md)
- [`README_EN.md`](./README_EN.md)
- [PostgreSQL MCP Server 项目主页](https://github.com/TencentCloudCommunity/mcp-server)
- [云数据库 PostgreSQL 产品文档](https://cloud.tencent.com/document/product/409)
- [云数据库 PostgreSQL API 总览](https://cloud.tencent.com/document/product/409/16761)
- [SCF 云函数控制台](https://console.cloud.tencent.com/scf)
