# 云数据库 TencentDB for PostgreSQL MCP Server

> 本 MCP Server 把这些云 API 统一封装为 MCP 工具，配套提供：**48 个工具**，覆盖实例、账号、数据库、参数、备份、监控、网络、只读实例、SSL 配置等 9 大模块；**多种部署形态**：`npx` 一键本地拉起、本地 `stdio` 命令直连、`streamable-http` 自建服务、`sse` 自建兼容模式，以及 **腾讯云 SCF Web 函数**部署；**per-request 凭证模式**：客户端在每次请求时通过 Header 传递自己的腾讯云 `SecretId` / `SecretKey`，服务端不长期保存用户密钥。

**产品链接**：[云数据库 TencentDB for PostgreSQL](https://cloud.tencent.com/product/postgres)

---

## 一、工具列表（Tools）

下面按模块分组列出全部 **48** 个工具。所有工具均接受 `region`（地域）作为第一个必填参数，工具内部会按 `region` 调用对应的腾讯云 API。

### 1. 实例（15）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 1 | `DescribeDBInstances` | 查询实例列表 |
| 2 | `DescribeDBInstanceAttribute` | 查询实例详情 |
| 3 | `CreateInstances` | 创建实例（费用确认） |
| 4 | `ModifyDBInstanceName` | 修改实例名称 |
| 5 | `ModifyDBInstanceSpec` | 变更实例规格（扩缩容，费用确认） |
| 6 | `RestartDBInstance` | 重启实例 |
| 7 | `IsolateDBInstances` | 隔离实例（业务确认） |
| 8 | `DisIsolateDBInstances` | 解除隔离实例 |
| 9 | `UpgradeDBInstanceKernelVersion` | 升级实例内核版本号 |
| 10 | `DescribeTasks` | 查询异步任务状态 |
| 11 | `DescribeClasses` | 查询可用规格列表 |
| 12 | `DescribeDBVersions` | 查询可用数据库版本 |
| 13 | `DescribeRegions` | 查询售卖地域 |
| 14 | `DescribeZones` | 查询售卖可用区 |
| 15 | `DescribeProductConfig` | 查询售卖规格配置 |

### 2. 账号（6）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 16 | `DescribeAccounts` | 查询实例的数据库账号列表 |
| 17 | `DescribeAccountPrivileges` | 查询数据库账号的权限信息 |
| 18 | `CreateAccount` | 创建数据库账号 |
| 19 | `DeleteAccount` | 删除数据库账号 |
| 20 | `ModifyAccountPrivileges` | 修改账号权限（授权/收回/修改账号类型） |
| 21 | `ResetAccountPassword` | 重置账号密码 |

### 3. 数据库（4）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 22 | `DescribeDatabases` | 查询实例的数据库列表 |
| 23 | `DescribeDatabaseObjects` | 查询数据库对象列表 |
| 24 | `CreateDatabase` | 创建数据库 |
| 25 | `ModifyDatabaseOwner` | 修改数据库属主 |

### 4. 参数（5）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 26 | `DescribeDBInstanceParameters` | 查询实例参数 |
| 27 | `DescribeParameterTemplates` | 查询参数模板列表 |
| 28 | `DescribeParameterTemplateAttributes` | 查询参数模板详情 |
| 29 | `DescribeParamsEvent` | 查询参数修改事件 |
| 30 | `ModifyDBInstanceParameters` | 修改实例参数 |

### 5. 备份（8）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 31 | `DescribeBackupOverview` | 查询备份概览 |
| 32 | `DescribeBaseBackups` | 查询基础备份列表 |
| 33 | `DescribeLogBackups` | 查询日志备份列表 |
| 34 | `DescribeAvailableRecoveryTime` | 查询可恢复时间范围 |
| 35 | `DescribeCloneDBInstanceSpec` | 查询克隆实例可购买的规格 |
| 36 | `DescribeBackupDownloadURL` | 获取备份下载链接 |
| 37 | `CreateBaseBackup` | 创建基础备份 |
| 38 | `CloneDBInstance` | 克隆实例（费用确认） |

### 6. 监控（3）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 39 | `DescribeSlowQueryList` | 查询慢查询列表 |
| 40 | `DescribeSlowQueryAnalysis` | 慢查询分析 |
| 41 | `DescribeDBErrlogs` | 查询错误日志 |

### 7. 网络（4）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 42 | `OpenDBExtranetAccess` | 开启实例公网访问 |
| 43 | `CloseDBExtranetAccess` | 关闭实例公网访问 |
| 44 | `DescribeDBInstanceSecurityGroups` | 查询实例安全组 |
| 45 | `ModifyDBInstanceSecurityGroups` | 修改实例安全组 |

### 8. SSL（1）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 46 | `DescribeDBInstanceSSLConfig` | 查询实例 SSL 配置 |

### 9. 只读实例（2）

| 序号 | 工具名称 | 详细描述 |
|---|---|---|
| 47 | `DescribeReadOnlyGroups` | 查询只读组列表 |
| 48 | `CreateReadOnlyDBInstance` | 创建只读实例（费用确认） |

> **分级保护**：写类工具按风险划分为 `LevelFee`（费用确认）、`LevelBusiness`（业务确认）、`LevelAudit`（审计级）、`LevelCritical`（高危，需显式二次确认 `confirm=true`）。

---

## 二、快速开始（部署方式）

下面按 **推荐度从高到低** 列出 4 种部署方式。请按需选择，并先看清各方式的前置条件：

- **方式一：腾讯云 SCF 自助托管（推荐云上部署）** —— 适合希望通过腾讯云托管并对外提供 HTTPS 访问的场景；需要您自行在 SCF 控制台完成函数创建、zip 上传和环境变量配置。
- **方式二：自建 streamable-http 服务** —— 需先按需拉取 `src/postgres` 目录（`git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git` → `cd mcp-server` → `git sparse-checkout set src/postgres` → `cd src/postgres`），并准备 **Go 1.25+** 用于本地编译运行。
- **方式三：本地 stdio** —— 需先按需拉取 `src/postgres` 目录（同上）并准备 **Go 1.25+**，适合本地 Cursor / Claude Desktop / WorkBuddy 客户端。
- **方式四：npx 一键拉起** —— 只需本机安装 **Node.js 18+**，无需克隆仓库，命令行一条即可。

### 方式一：腾讯云 SCF 自助托管（推荐云上部署）

> 适合希望运行在腾讯云并通过 HTTPS / 函数 URL 提供给团队共用的场景。仓库已提供 SCF 打包脚本、启动脚本和环境变量模板，您需要在自己的腾讯云账号下完成函数创建与发布。

#### 1.1 按需拉取 `src/postgres` 目录并构建 SCF 发布包

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
./scripts/build_scf_zip.sh
```

默认会在 `dist/` 目录生成可上传到 SCF 的 zip 包。

#### 1.2 在 SCF 控制台创建 Web 函数

进入 [SCF 云函数控制台](https://console.cloud.tencent.com/scf)，按以下方式创建：

- 函数类型：**Web 函数**
- 运行环境：**Go 标准运行环境**
- 代码上传方式：**本地上传 zip**

上传 zip 后，为函数开启公网访问 URL。

#### 1.3 启动命令

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

#### 1.4 环境变量配置

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

#### 1.5 客户端配置示例

服务采用 **per-request 凭证模式**，调用时需要在 Header 中按请求附带自己的腾讯云凭证（不要把密钥写进 URL 或 query 参数）：

- `X-TencentCloud-Secret-Id`
- `X-TencentCloud-Secret-Key`

> 凭证获取方法见 [§三、密钥获取教程](#三密钥获取教程)。

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

#### 1.6 快速验证

```bash
curl -i https://您的函数URL/healthz
```

正常返回 `200 OK` 即可按上面 JSON 配置客户端，开始调用 48 个 MCP 工具。


---

### 方式二：自建 streamable-http 服务

适合部署到自有云主机、内网服务器，通过域名给团队共用。

> **前置条件**：本机或云主机已安装 **Go 1.25+**，并具备外网出口（要访问腾讯云 OpenAPI）。

#### 2.1 按需拉取 `src/postgres` 目录

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> 后续所有命令都需要在 `src/postgres/` 目录下执行。

#### 2.2 准备配置

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

#### 2.3 启动服务

```bash
./scripts/run_server.sh
```

#### 2.4 MCP 客户端配置

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

### 方式三：本地 stdio（推荐本地客户端）

适合 Cursor、Claude Desktop、WorkBuddy 等本地 MCP 客户端的命令直连模式。

> **前置条件**：本机已安装 **Go 1.25+**，并具备外网出口（要访问腾讯云 OpenAPI）。

#### 3.1 按需拉取 `src/postgres` 目录

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> 后续所有命令都需要在 `src/postgres/` 目录下执行。

#### 3.2 准备配置

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

#### 3.3 启动

```bash
./scripts/run_stdio.sh
```

#### 3.4 MCP 客户端配置

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

### 方式四：npx 一键拉起（最简本地体验）

> **前置条件**：本机已安装 **Node.js 18+**（含 `npx`）。无需克隆 Go 仓库，npm 包会按平台自动从 GitHub Release 下载预编译二进制。
>
> 检查是否已安装：
>
> ```bash
> node -v   # 期望 v18.x 或更高
> ```

#### 4.1 命令行直接启动

```bash
npx -y postgres-mcp-server@latest
```

#### 4.2 MCP 客户端配置

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

---

## 三、密钥获取教程

### 获取 SecretId / SecretKey

1. 登录 [腾讯云控制台](https://console.cloud.tencent.com/)；
2. 进入 **访问管理 → API 密钥管理** 页面；
3. 在 **API 密钥** 标签页新建或查看已有的 **SecretId** 和 **SecretKey**。

> **生成地址**：<https://console.cloud.tencent.com/cam/capi>

⚠️ **安全建议**

- **不要把 `SecretId / SecretKey` 写进 URL 或 query 参数**，只通过 Header 或环境变量传递；
- **不要把密钥放进 SCF 服务端环境变量**，应使用 per-request 凭证模式；
- **建议使用最小权限的 CAM 子账号**，避免长期复用主账号密钥；
- **不要在日志、trace、错误回显中输出凭据明文**。

### 推荐权限范围

为 MCP 客户端单独创建一个 CAM 子账号或角色，按使用场景授予最小权限：

| 使用场景 | 推荐策略 |
|---|---|
| 只读巡检 | `QcloudPGReadOnlyAccess` |
| 日常运维（包含写） | 自定义策略，限制到具体 `Action` 列表 |
| 临时排查 | 最小权限的 CAM 子账号，用后及时回收 |

---

## 四、地域映射

调用本 MCP 工具时，`region` 参数必须传 **地域代码**（如 `ap-guangzhou`），而不是中文地域名。

示例：广州地域

| 中文名 | 地域代码 |
|---|---|
| 广州 | `ap-guangzhou` |
| 上海 | `ap-shanghai` |
| 北京 | `ap-beijing` |
| 南京 | `ap-nanjing` |
| 深圳 | `ap-shenzhen` |
| 成都 | `ap-chengdu` |
| 香港 | `ap-hongkong` |
| 新加坡 | `ap-singapore` |

> 完整地域列表见 [地域与可用区映射文档](https://cloud.tencent.com/document/product/1596/77930)。如果您不知道实例在哪个地域，可以用只读工具 `DescribeDBInstances` 不传 `region`，或先调用 `DescribeRegions` 查询。


## 五、API 使用参考



 云数据库 PostgreSQL API 总览  <https://cloud.tencent.com/document/product/409/16761> 



## 六、安全建议

- **生产环境必须放在 HTTPS / 反向代理之后**，不要裸跑 HTTP 到公网；
- **不要把 `SecretId / SecretKey` 放进 URL 或 query 参数**，全部走 Header 或 env；
- **禁止在日志、trace、错误回显中输出凭据明文**；
- **建议使用最小权限的 CAM 子账号**，避免长期复用主账号密钥；
- **保持 `READ_ONLY=true` 起步**，确认流程后再按需开放写操作；
- **对外暴露前务必加 IP 白名单、VPN 或零信任访问控制**；
- **SCF / API 网关等无状态环境只推荐 `streamable-http`** + `MCP_STREAMABLE_HTTP_STATELESS=true`。



## 七、许可证

本项目基于 **Apache-2.0** 协议开源。

