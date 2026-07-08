# 云数据库 TencentDB for PostgreSQL MCP Server

> 腾讯云 PostgreSQL（TencentDB for PostgreSQL）官方 MCP Server，把云 API 背后的实例、账号、数据库、参数、备份、监控、网络、只读实例与 SSL 等能力，统一封装为 MCP 工具，可被 Cursor、Claude Desktop、WorkBuddy 等任何兼容 MCP 的客户端直接调用。支持 `npx` 一键本地拉起，也支持 `stdio` / `streamable-http` / `sse` 三种 transport 部署到远程主机或 SCF，并通过 per-request 凭证模式避免在服务端长期保存用户的 SecretId / SecretKey。

**产品链接**：[云数据库 TencentDB for PostgreSQL](https://cloud.tencent.com/product/postgres)

---

## 一、工具列表（Tools）

默认注册 **48 个工具**，覆盖 9 大模块：实例、账号、数据库、参数、备份、监控、网络、SSL、只读实例。所有工具均接受 `region`（地域代码）作为第一个必填参数，内部按 `region` 调用对应的腾讯云 API。

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

## 二、特性

- **多种部署形态**：`npx` 一键本地拉起、本地 `stdio` 命令直连、`streamable-http` 自建服务、`sse` 自建兼容模式、腾讯云 SCF Web 函数
- **per-request 凭证模式**：客户端在每次请求时通过 Header 传递自己的腾讯云 `SecretId` / `SecretKey`，服务端不长期保存用户密钥
- **三种 transport 兼容**：`stdio` / `streamable-http` / `sse` 均可启动，源码启动与脚本启动方式一致
- **npx 分发层（新增）**：默认拉起 `stdio`，从 GitHub Release 拉取预编译 Go 二进制，**不影响现有源码 / 自建服务 / 脚本启动方式**
- **最小权限起步**：默认 `pg.read` + `READ_ONLY=true`，可按需开放写操作

---

## 三、部署方式（按推荐度从高到低）

| 方式 | transport | 适用场景 | 前置条件 |
|---|---|---|---|
| 方式一：腾讯云 SCF 自助托管 | `streamable-http` | 云上自助托管、团队共用 | 控制台上传 zip |
| 方式二：自建 streamable-http | `streamable-http` | 自建云主机、内网服务器 | Go 1.25+ |
| 方式三：本地 stdio | `stdio` | Cursor / Claude Desktop / WorkBuddy | Go 1.25+ |
| 方式四：npx 一键拉起 | `stdio`（默认） | 最简本地体验 | Node.js 18+ |

### 方式一：腾讯云 SCF 自助托管

适合希望运行在腾讯云并通过 HTTPS / 函数 URL 提供给团队共用的场景。仓库已提供 `deploy/scf/` 下的 SCF 部署物料与 `./scripts/build_scf_zip.sh` 打包脚本，但**需要您在自己的腾讯云账号下完成函数创建、zip 上传与发布**。

详细步骤见：[`SCF_DEPLOY.md`](./SCF_DEPLOY.md)

最小可用配置（直接复制 `deploy/scf/scf.console.env.txt` 即可）：

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

客户端 `url` 必须填写**完整 MCP 端点** `https://您的函数URL/mcp`；函数根 URL `https://您的函数URL` 返回 `404 page not found` 属于预期。

### 方式二：自建 streamable-http 服务

适合部署到自有云主机、内网服务器，通过域名给团队共用。

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
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

启动服务：

```bash
./scripts/run_server.sh
```

可访问地址：

- 健康检查：`http://127.0.0.1:9000/healthz`
- 就绪检查：`http://127.0.0.1:9000/readyz`
- MCP 端点：`http://127.0.0.1:9000/mcp`

MCP 客户端配置：

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

### 方式三：本地 stdio

适合 Cursor、Claude Desktop、WorkBuddy 等本地 MCP 客户端的命令直连模式。

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
cp .env.example .env
```

最小推荐配置：

```env
MCP_TRANSPORT=stdio
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_REQUEST_SECRET_ID=您的SecretId
MCP_REQUEST_SECRET_KEY=您的SecretKey
READ_ONLY=true
```

启动：

```bash
./scripts/run_stdio.sh
```

MCP 客户端配置：

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

### 方式四：npx 一键拉起

> 前置条件：本机已安装 **Node.js 18+**（含 `npx`）。无需克隆 Go 仓库，npm 包会按平台自动从 GitHub Release 下载预编译二进制。
>
> ```bash
> node -v   # 期望 v18.x 或更高
> ```

直接启动：

```bash
npx -y postgres-mcp-server@latest
```

MCP 客户端配置：

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

可选：用 `npx` 拉起自建 `sse`：

```bash
npx -y postgres-mcp-server@latest --transport sse --env-file .env
```

说明：

- `npx` 启动器会按平台下载 GitHub Release 里的预编译 Go 二进制
- 默认 `transport` 是 `stdio`；也兼容 `TRANSPORT -> MCP_TRANSPORT`、`PORT -> MCP_SERVER_PORT`
- 腾讯云凭证仍可直接使用 `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY`
- 如需本地调试或跳过下载，可设置 `POSTGRES_MCP_BINARY_PATH=/path/to/postgres-server`
- 如需先在**你自己的 GitHub fork** 上试跑，可加 `--release-repository <owner/repo>`，必要时再配合 `--release-tag <tag>`
- 发布预编译资产可执行：`./scripts/build_npx_release.sh`

#### 先在你自己的 GitHub 上试跑 `npx`（不需要先发 npm）

适合 GitHub 新手的最小路径：**本地打一个 npm 包 + 你的 fork 上放 GitHub Release 资产**。

1. 在 `src/postgres` 目录构建预编译二进制：

```bash
./scripts/build_npx_release.sh 1.0.3
```

2. 到你自己的 GitHub fork 创建一个 Release：

- 仓库：`<你的 GitHub 用户名>/mcp-server`
- Tag：`postgres-mcp-server-v1.0.3`
- 上传文件：`dist/npx/` 目录下的全部文件

3. 在本地打出 npm 包（不用发布到 npm）：

```bash
npm pack
```

4. 先验证 `npx` 能执行本地包：

```bash
npx -y ./postgres-mcp-server-1.0.3.tgz --version
```

5. 再切到你自己的 GitHub Release 做真实下载测试：

```bash
npx -y ./postgres-mcp-server-1.0.3.tgz --release-repository <你的 GitHub 用户名>/mcp-server
```

如果你的 Release tag 不是默认的 `postgres-mcp-server-v1.0.3`，再追加：

```bash
--release-tag <你的 tag>
```

这样你就能先验证：**包本身来自本地 `.tgz`，二进制来自你自己的 GitHub Release**。整个流程跑通后，再决定要不要正式发 npm。

---

## 四、源码安装（可选）

如果您不想依赖脚本，也可以直接编译并按 transport 启动。

```bash
go build -o ./.bin/postgres-server .
```

- **stdio**：`MCP_TRANSPORT=stdio ./.bin/postgres-server`
- **streamable-http**：`MCP_TRANSPORT=streamable-http ./.bin/postgres-server`
- **sse**：`MCP_TRANSPORT=sse ./.bin/postgres-server`

常用端点：

- **streamable-http**：`/mcp`
- **SSE**：`/sse`
- **SSE message**：`/message`

---

## 五、本仓库自带客户端如何连当前服务

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
export MCP_REQUEST_SECRET_ID=您的SecretId
export MCP_REQUEST_SECRET_KEY=您的SecretKey
```

也兼容直接复用：

```bash
export MCP_SECRET_ID=您的SecretId
export MCP_SECRET_KEY=您的SecretKey
```

---

## 六、密钥获取教程

### 获取 SecretId / SecretKey

1. 登录 [腾讯云控制台](https://console.cloud.tencent.com/)；
2. 进入 **访问管理 → API 密钥管理** 页面；
3. 在 **API 密钥** 标签页新建或查看已有的 **SecretId** 和 **SecretKey**。

> **生成地址**：<https://console.cloud.tencent.com/cam/capi>

### 推荐权限范围

为 MCP 客户端单独创建一个 CAM 子账号或角色，按使用场景授予最小权限：

| 使用场景 | 推荐策略 |
|---|---|
| 只读巡检 | `QcloudPGReadOnlyAccess` |
| 日常运维（包含写） | 自定义策略，限制到具体 `Action` 列表 |
| 临时排查 | 最小权限的 CAM 子账号，用后及时回收 |

---

## 七、地域映射

调用本 MCP 工具时，`region` 参数必须传 **地域代码**（如 `ap-guangzhou`），而不是中文地域名。

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

---

## 八、API 使用参考

[云数据库 PostgreSQL API 总览](https://cloud.tencent.com/document/product/409/16761)

---

## 九、安全建议

- **生产环境必须放在 HTTPS / 反向代理之后**，不要裸跑 HTTP 到公网
- **不要把 `SecretId / SecretKey` 放进 URL 或 query 参数**，只通过 Header 或 env 传递
- **禁止在日志、trace、错误回显中输出凭据明文**
- **建议使用最小权限的 CAM 子账号**，避免长期复用主账号密钥
- **保持 `READ_ONLY=true` 起步**，确认流程后再按需开放写操作
- **对外暴露前务必加 IP 白名单、VPN 或零信任访问控制**
- **SCF / API 网关等无状态环境只推荐 `streamable-http`** + `MCP_STREAMABLE_HTTP_STATELESS=true`
- **不要把 `SecretId / SecretKey` 放进 SCF 服务端环境变量**，应使用 per-request 凭证模式

### Scope 与地域

```env
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read,pg.write
MCP_REQUEST_ALLOWED_REGIONS=ap-guangzhou
```

- `pg.read`：允许只读工具
- `pg.write`：允许写类工具（仍受 `Guard` 和 `confirm=true` 约束）
- `MCP_REQUEST_ALLOWED_REGIONS` 留空表示不额外做地域限制

---

## 十、相关文档

- 综合部署指南：[`DEPLOY.md`](./DEPLOY.md)
- SCF Web 函数部署：[`SCF_DEPLOY.md`](./SCF_DEPLOY.md)
- 本地 stdio 启动：[`scripts/run_stdio.sh`](./scripts/run_stdio.sh)
- 通用服务启动：[`scripts/run_server.sh`](./scripts/run_server.sh)
- WorkBuddy 使用问题记录：[`WORKBUDDY_USAGE_NOTES.md`](./WORKBUDDY_USAGE_NOTES.md)
- npx 发布脚本：[`scripts/build_npx_release.sh`](./scripts/build_npx_release.sh)

---

## 十一、许可证

本项目基于 **Apache-2.0** 协议开源。
