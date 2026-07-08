# PostgreSQL MCP Server 部署教程

> 本教程面向首次接入腾讯云 PostgreSQL MCP Server 的用户。点击“部署 MCP Server”后弹窗中渲染本文档，建议从左上方目录按顺序阅读。

## 1. 部署前准备

在选择部署方式前，请先确认以下信息是否已就绪。

### 1.1 适用场景

- **想部署到腾讯云并通过 HTTPS 暴露给团队使用**：选方式一（腾讯云 SCF 自助托管）。
- **想自己控制网络、主机与运行环境**：选方式二（自建 `streamable-http` 服务）。
- **想在本地客户端里直接接入源码服务**：选方式三（本地 `stdio`）。
- **只想本机快速拉起体验**：选方式四（`npx` 一键拉起）。

> 如需云上托管，请按方式一在自己的腾讯云账号下完成 SCF 部署，并使用您自己的函数 URL 接入客户端。

### 1.2 所需权限

部署完成后调用 MCP 工具时，需要一个具备 `postgres:*` 访问权限的腾讯云账号。推荐为 MCP 客户端单独创建一个 CAM 子账号，按使用场景授予最小权限。


凭证获取地址：<https://console.cloud.tencent.com/cam/capi>，
请注意，您需要保存您的 `SecretId` 和 `SecretKey`，用于后续调用 OpenAPI。

此外，使用本mcp server提供的工具需要地域代码，请确认您想要访问的实例部署在哪个区域。

参考文档：https://cloud.tencent.com/document/product/1596/77930
### 1.3 前置资源

| 资源 | 是否必须 | 说明 |
|---|---|---|
| `SecretId` / `SecretKey` | 是 | 调用 OpenAPI 的身份凭证 |
| 地域代码（如 `ap-guangzhou`） | 是 | 作为所有工具的 `region` 参数 |
| 实例 ID | 否 | 仅在按实例操作时需要 |

## 2. MCP 开放能力

默认开放 **48 个工具**，覆盖实例、账号、数据库、参数、备份、监控、网络、SSL、只读实例等 9 大模块。下面按模块分组列出。

### 2.1 实例管理

| 工具名称 | 描述 |
|---|---|
| `DescribeDBInstances` | 查询实例列表 |
| `DescribeDBInstanceAttribute` | 查询实例详情 |
| `CreateInstances` | 创建实例（费用确认） |
| `ModifyDBInstanceName` | 修改实例名称 |
| `ModifyDBInstanceSpec` | 变更实例规格（扩缩容，费用确认） |
| `RestartDBInstance` | 重启实例 |
| `IsolateDBInstances` | 隔离实例（业务确认） |
| `DisIsolateDBInstances` | 解除隔离实例 |
| `UpgradeDBInstanceKernelVersion` | 升级实例内核版本号 |
| `DescribeTasks` | 查询异步任务状态 |
| `DescribeClasses` | 查询可用规格列表 |
| `DescribeDBVersions` | 查询可用数据库版本 |
| `DescribeRegions` | 查询售卖地域 |
| `DescribeZones` | 查询售卖可用区 |
| `DescribeProductConfig` | 查询售卖规格配置 |

### 2.2 账号管理

| 工具名称 | 描述 |
|---|---|
| `DescribeAccounts` | 查询实例的数据库账号列表 |
| `DescribeAccountPrivileges` | 查询数据库账号的权限信息 |
| `CreateAccount` | 创建数据库账号 |
| `DeleteAccount` | 删除数据库账号 |
| `ModifyAccountPrivileges` | 修改账号权限（授权/收回/修改账号类型） |
| `ResetAccountPassword` | 重置账号密码 |

### 2.3 数据库管理

| 工具名称 | 描述 |
|---|---|
| `DescribeDatabases` | 查询实例的数据库列表 |
| `DescribeDatabaseObjects` | 查询数据库对象列表 |
| `CreateDatabase` | 创建数据库 |
| `ModifyDatabaseOwner` | 修改数据库属主 |

### 2.4 参数管理

| 工具名称 | 描述 |
|---|---|
| `DescribeDBInstanceParameters` | 查询实例参数 |
| `DescribeParameterTemplates` | 查询参数模板列表 |
| `DescribeParameterTemplateAttributes` | 查询参数模板详情 |
| `DescribeParamsEvent` | 查询参数修改事件 |
| `ModifyDBInstanceParameters` | 修改实例参数 |

### 2.5 备份与恢复

| 工具名称 | 描述 |
|---|---|
| `DescribeBackupOverview` | 查询备份概览 |
| `DescribeBaseBackups` | 查询基础备份列表 |
| `DescribeLogBackups` | 查询日志备份列表 |
| `DescribeAvailableRecoveryTime` | 查询可恢复时间范围 |
| `DescribeCloneDBInstanceSpec` | 查询克隆实例可购买的规格 |
| `DescribeBackupDownloadURL` | 获取备份下载链接 |
| `CreateBaseBackup` | 创建基础备份 |
| `CloneDBInstance` | 克隆实例（费用确认） |

### 2.6 监控诊断

| 工具名称 | 描述 |
|---|---|
| `DescribeSlowQueryList` | 查询慢查询列表 |
| `DescribeSlowQueryAnalysis` | 慢查询分析 |
| `DescribeDBErrlogs` | 查询错误日志 |

### 2.7 网络管理

| 工具名称 | 描述 |
|---|---|
| `OpenDBExtranetAccess` | 开启实例公网访问 |
| `CloseDBExtranetAccess` | 关闭实例公网访问 |
| `DescribeDBInstanceSecurityGroups` | 查询实例安全组 |
| `ModifyDBInstanceSecurityGroups` | 修改实例安全组 |

### 2.8 SSL 配置

| 工具名称 | 描述 |
|---|---|
| `DescribeDBInstanceSSLConfig` | 查询实例 SSL 配置 |

### 2.9 只读实例

| 工具名称 | 描述 |
|---|---|
| `DescribeReadOnlyGroups` | 查询只读组列表 |
| `CreateReadOnlyDBInstance` | 创建只读实例（费用确认） |

> 写类工具仍受权限范围、`READ_ONLY` 配置以及二次确认机制约束。建议先以只读能力接入，再按需开放写操作。

## 3. 选择部署方式

下面按 **推荐度从高到低** 列出 4 种方式。每种方式都按“前置条件 → 部署步骤 → 客户端配置”的顺序展开。

### 3.1 方式一：腾讯云 SCF 自助托管（推荐云上部署）

> 适合希望运行在腾讯云并通过 HTTPS / 函数 URL 提供给团队共用的场景。仓库已提供 SCF 打包脚本、启动脚本和环境变量模板，但**需要您在自己的腾讯云账号下完成托管发布**。

#### 步骤一：按需拉取 `src/postgres` 目录并构建 SCF 发布包

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
./scripts/build_scf_zip.sh
```

默认会在 `dist/` 目录生成可上传到 SCF 的 zip 包。

#### 步骤二：在 SCF 控制台创建 Web 函数

进入 [SCF 云函数控制台](https://console.cloud.tencent.com/scf/list?rid=16&ns=default)，按以下方式创建：

- 创建方式：请选择“从头开始”
- 函数类型：**Web 函数**
- 运行环境：**Go 标准运行环境（Go 1）**
- 代码上传方式：**本地上传 zip**
- 环境变量设置参考步骤四，也可以创建函数后再设置
- 其它设置按需，如需要公网访问请在最后的“函数URL配置”勾选启用公网访问

#### 步骤三：启动命令

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

#### 步骤四：配置环境变量

最小推荐配置：

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

#### 步骤五：客户端配置

部署完成后，客户端连接您自己的函数 URL：

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

> 客户端 `url` 必须指向**完整 MCP 端点**（含 `/mcp` 后缀），不要使用函数根 URL。

#### 步骤五：快速验证

```bash
curl -i https://您的函数URL/healthz
```

正常返回 `200 OK` 后，再用 `https://您的函数URL/mcp` 接入客户端。

如需完整控制台操作说明，请查看 [`SCF_DEPLOY.md`](./SCF_DEPLOY.md)。

---

### 3.2 方式二：自建 `streamable-http` 服务

适合部署到自有云主机、内网服务器，通过域名给团队共用。

> **前置条件**：本机或云主机已安装 **Go 1.25+**，并具备外网出口（要访问腾讯云 OpenAPI）。

#### 步骤一：按需拉取 `src/postgres` 目录

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> 后续所有命令都需要在 `src/postgres/` 目录下执行。

#### 步骤二：准备配置

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

#### 步骤三：启动服务

```bash
./scripts/run_server.sh
```

#### 步骤四：客户端配置

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

### 3.3 方式三：本地 `stdio`（推荐本地客户端）

适合 Cursor、Claude Desktop、WorkBuddy 等本地 MCP 客户端的命令直连模式。

> **前置条件**：本机已安装 **Go 1.25+**，并具备外网出口（要访问腾讯云 OpenAPI）。

#### 步骤一：按需拉取 `src/postgres` 目录

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

#### 步骤二：准备配置

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
请注意，需要在本地的.env文件中配置好腾讯云凭证才能正常使用。

#### 步骤三：启动

```bash
./scripts/run_stdio.sh
```

#### 步骤四：客户端配置

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "/绝对路径/mcp-server/src/postgres/scripts/run_stdio.sh",
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

### 3.4 方式四：`npx` 一键拉起（最简本地体验）

> **前置条件**：本机已安装 **Node.js 18+**（含 `npx`）。无需克隆 Go 仓库，npm 包会按平台自动从 GitHub Release 下载预编译二进制。

检查是否已安装：

```bash
node -v   # 期望 v18.x 或更高
```

#### 步骤一：命令行直接启动

```bash
npx -y postgres-mcp-server@latest
```

#### 步骤二：客户端配置

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

## 4. 鉴权与安全配置

无论选择哪种部署方式，调用时都使用 **per-request 凭证模式**，即每次请求通过 Header 携带腾讯云凭证，服务端不长期保存您的密钥。

### 4.1 凭证传递方式

- **HTTP / `streamable-http` / `sse` 模式**：通过 `X-TencentCloud-Secret-Id`、`X-TencentCloud-Secret-Key` Header 传递。
- **`stdio` 模式**：通过 `MCP_REQUEST_SECRET_ID` / `MCP_REQUEST_SECRET_KEY` 等环境变量注入。

### 4.2 必读安全建议

- **不要把 `SecretId / SecretKey` 写进 URL 或 query 参数**，只通过 Header 或环境变量传递。
- **不要把密钥放进 SCF 服务端环境变量**，应使用 per-request 凭证模式。
- **建议使用最小权限的 CAM 子账号**，避免长期复用主账号密钥。
- **不要在日志、trace、错误回显中输出凭据明文**。
- **生产环境必须放在 HTTPS / 反向代理之后**，不要裸跑 HTTP 到公网。
- **保持 `READ_ONLY=true` 起步**，确认流程后再按需开放写操作。


## 5. 相关文档链接

部署过程中如果需要进一步查阅资料，可以打开以下链接：

- [PostgreSQL MCP Server 项目主页](https://github.com/TencentCloudCommunity/mcp-server)
- [API 密钥（CAM）管理控制台](https://console.cloud.tencent.com/cam/capi)
- [SCF 云函数控制台](https://console.cloud.tencent.com/scf)
- [云数据库 PostgreSQL 产品文档](https://cloud.tencent.com/document/product/409)
- [云数据库 PostgreSQL API 总览](https://cloud.tencent.com/document/product/409/16761)
- [地域与可用区映射](https://cloud.tencent.com/document/product/1596/77930)
- [查看完整部署文档](README.md)
