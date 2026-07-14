# PostgreSQL 配套技能包

本目录存放腾讯云数据库 PostgreSQL 的配套技能，路径位于 `src/postgres/skills/`。

这些资源属于**工作流层的技能包**，不是独立的数据库连接器。它们不依赖已部署的 PostgreSQL MCP Server，而是直接调用腾讯云 PostgreSQL OpenAPI；但允许调用的 Action 严格限制在当前 `src/postgres` 已完成参数对齐与工具封装的范围内。

## 概览

- **同一领域，不同层次**：这些 skill 与 PostgreSQL MCP 实现放在相邻目录中，但与 Go 运行时、SCF 部署脚本以及 npm 启动器相互隔离。
- **证据优先、全量开放**：这些打包后的 skill 可以调用当前已完成参数契约校验的全部 **48 个 PostgreSQL OpenAPI Action**；默认先做证据采集，再按需进入变更或处置动作。
- **OpenAPI 直连**：skill 在运行时直接调用腾讯云 PostgreSQL OpenAPI，不要求用户额外先部署 MCP Server。
- **高风险动作需明确确认**：涉及审计类、业务类、费用类或 critical 级变更动作时，必须先说明影响面并获得明确确认。
- **面向发布的分发方式**：用户应从 GitHub Release 或 COS 安装打包好的 zip 资源，而不是直接使用源码目录。

## 用户傻瓜式配置（推荐）

对最终用户来说，**推荐只准备 3 个环境变量**，就可以直接使用这些 PostgreSQL skill：

```bash
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou" #请使用您想要访问的实例所在地域
# 可选：临时凭证场景再补 TENCENTCLOUD_SESSION_TOKEN
```

补充约定：
- **推荐变量名**：`TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY`、`TENCENTCLOUD_REGION`
- **兼容变量名**：也兼容读取 `MCP_REQUEST_SECRET_ID`、`MCP_REQUEST_SECRET_KEY`、`MCP_REQUEST_SESSION_TOKEN`、`MCP_SECRET_ID`、`MCP_SECRET_KEY`
- **地域可写自然语言**：`广州`、`上海`、`成都`、`北京` 这类常见中文地域，执行 skill 时应先归一化为标准地域码，例如 `广州 -> ap-guangzhou`、`上海 -> ap-shanghai`、`成都 -> ap-chengdu`。如果您不清楚您的实例等地域码，请访问: https://cloud.tencent.com/document/product/1596/77930
- **SDK 优先**：skill 执行时应优先使用腾讯云官方 SDK；如果运行环境里没有 SDK，可自动退回到本地生成 TC3 签名的 HTTPS 请求，如果失败则 skill可自动安装 SDK
- **密钥只放运行时环境**：不要把 `SecretId`、`SecretKey` 或 Token 写入仓库文件、skill 文档、URL 或 query 参数

## 包含的技能

### `tencent-pg-inspection`
- PostgreSQL 日常健康巡检
- 典型请求：`PG巡检`、`健康检查`、`备份检查`、`资源水位检查`
- 关注点：实例健康、备份状态、网络 / SSL / 只读上下文、参数姿态、风险总结，以及必要时的对齐动作处置

### `tencent-pg-slowquery-diagnosis`
- 慢 SQL 与性能诊断
- 典型请求：`慢SQL分析`、`SQL性能诊断`、`查询为什么变慢`
- 关注点：慢查询证据、错误日志、实例上下文、可能原因排序、安全优化建议，以及必要时的对齐动作处置

### `tencent-pg-ops-troubleshooter`
- 运维排障工作流
- 典型请求：`PG排障`、`实例异常排查`、`SSL问题`、`备份失败`
- 关注点：故障分类、按模块采集证据、Runbook 风格的后续处理步骤，以及必要时的对齐动作处置

## 与 PostgreSQL MCP 的对齐边界

这些 skill 不把 PostgreSQL MCP Server 当作运行前置条件，但仍然必须与 `src/postgres` 的 OpenAPI 对齐结果保持一致。

当前约束如下：
- 允许调用的 Action 必须来自当前已对齐的 PostgreSQL OpenAPI 集合，维护基线以 `tools/openapi_alignment.go` 为准。
- 这三个工作流 skill 现在**都允许使用这 48 个已对齐 Action**；但应优先从与当前任务最相关的模块开始，而不是无差别全量调用。
- 即使腾讯云官方还提供更多 PostgreSQL OpenAPI，**只要当前仓库未对齐，就不会在这些 skill 中调用**。


## 当前开放的 48 个对齐 Action

### 1. 实例管理（15）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeDBInstances` | 查询实例列表 | 直接可用 |
| `DescribeDBInstanceAttribute` | 查询实例详情 | 直接可用 |
| `CreateInstances` | 创建实例 | 费用确认 |
| `ModifyDBInstanceName` | 修改实例名称 | 明确确认 |
| `ModifyDBInstanceSpec` | 变更实例规格 | 费用确认 |
| `RestartDBInstance` | 重启实例 | 明确确认 |
| `IsolateDBInstances` | 隔离实例 | 明确确认 |
| `DisIsolateDBInstances` | 解除隔离实例 | 明确确认 |
| `UpgradeDBInstanceKernelVersion` | 升级实例内核版本号 | 明确确认 |
| `DescribeTasks` | 查询异步任务状态 | 直接可用 |
| `DescribeClasses` | 查询可用规格列表 | 直接可用 |
| `DescribeDBVersions` | 查询可用数据库版本 | 直接可用 |
| `DescribeRegions` | 查询售卖地域 | 直接可用 |
| `DescribeZones` | 查询售卖可用区 | 直接可用 |
| `DescribeProductConfig` | 查询售卖规格配置 | 直接可用 |

### 2. 账号管理（6）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeAccounts` | 查询实例的数据库账号列表 | 直接可用 |
| `DescribeAccountPrivileges` | 查询数据库账号的权限信息 | 直接可用 |
| `CreateAccount` | 创建数据库账号 | 明确确认 |
| `DeleteAccount` | 删除数据库账号 | 明确确认 |
| `ModifyAccountPrivileges` | 修改账号权限 | 明确确认 |
| `ResetAccountPassword` | 重置账号密码 | 明确确认 |

### 3. 数据库管理（4）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeDatabases` | 查询实例的数据库列表 | 直接可用 |
| `DescribeDatabaseObjects` | 查询数据库对象列表 | 直接可用 |
| `CreateDatabase` | 创建数据库 | 明确确认 |
| `ModifyDatabaseOwner` | 修改数据库属主 | 明确确认 |

### 4. 参数管理（5）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeDBInstanceParameters` | 查询实例参数 | 直接可用 |
| `DescribeParameterTemplates` | 查询参数模板列表 | 直接可用 |
| `DescribeParameterTemplateAttributes` | 查询参数模板详情 | 直接可用 |
| `DescribeParamsEvent` | 查询参数修改事件 | 直接可用 |
| `ModifyDBInstanceParameters` | 修改实例参数 | Critical，明确确认 |

### 5. 备份与恢复（8）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeBackupOverview` | 查询备份概览 | 直接可用 |
| `DescribeBaseBackups` | 查询基础备份列表 | 直接可用 |
| `DescribeLogBackups` | 查询日志备份列表 | 直接可用 |
| `DescribeAvailableRecoveryTime` | 查询可恢复时间范围 | 直接可用 |
| `DescribeCloneDBInstanceSpec` | 查询克隆实例可购买的规格 | 直接可用 |
| `DescribeBackupDownloadURL` | 获取备份下载链接 | 明确确认 |
| `CreateBaseBackup` | 创建基础备份 | 明确确认 |
| `CloneDBInstance` | 克隆实例 | 费用确认 |

### 6. 监控诊断（3）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeSlowQueryList` | 查询慢查询列表 | 直接可用 |
| `DescribeSlowQueryAnalysis` | 慢查询分析 | 直接可用 |
| `DescribeDBErrlogs` | 查询错误日志 | 直接可用 |

### 7. 网络管理（4）

| Action | 说明 | 执行要求 |
|---|---|---|
| `OpenDBExtranetAccess` | 开启实例公网访问 | 明确确认 |
| `CloseDBExtranetAccess` | 关闭实例公网访问 | 明确确认 |
| `DescribeDBInstanceSecurityGroups` | 查询实例安全组 | 直接可用 |
| `ModifyDBInstanceSecurityGroups` | 修改实例安全组 | 明确确认 |

### 8. SSL 配置（1）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeDBInstanceSSLConfig` | 查询实例 SSL 配置 | 直接可用 |

### 9. 只读实例（2）

| Action | 说明 | 执行要求 |
|---|---|---|
| `DescribeReadOnlyGroups` | 查询只读组列表 | 直接可用 |
| `CreateReadOnlyDBInstance` | 创建只读实例 | 费用确认 |

> 说明：
> - **直接可用**：查询类动作，可直接用于证据采集与诊断。
> - **明确确认**：执行前必须先说明影响面、目标实例与预期结果，并获得用户明确确认。
> - **费用确认**：除明确确认外，还需要提前提示可能产生资源或计费变化。

## 目录结构

```text
skills/
├─ README.md
├─ package.json
├─ references/
│  └─ common/
│     ├─ region_normalization.md
│     └─ error_handling.md
├─ scripts/
│  ├─ package-all.mjs
│  ├─ package-skill.mjs
│  └─ verify-skill.mjs
├─ dist/
├─ tencent-pg-inspection/
│  ├─ SKILL.md
│  ├─ references/
│  └─ assets/
├─ tencent-pg-slowquery-diagnosis/
│  ├─ SKILL.md
│  ├─ references/
│  └─ assets/
└─ tencent-pg-ops-troubleshooter/
   ├─ SKILL.md
   ├─ references/
   └─ assets/
```

## 本地开发

1. 更新目标 skill 目录下的 `SKILL.md` 和 `references/`。
2. 重新确认所引用的 OpenAPI Action 仍然存在于 `tools/openapi_alignment.go`，并与仓库根 `README.md` 的能力清单一致。
3. 如有 OpenAPI 参数映射改动，先执行 `src/postgres/scripts/run_openapi_param_check.sh` 做契约校验。
4. 打包前执行本地校验。
5. 只有在校验通过后才产出发布 zip 包。

推荐命令：

```bash
cd src/postgres
./scripts/run_openapi_param_check.sh

cd skills
npm run verify
npm run release
```

## 打包

打包流程有意与 `src/postgres/package.json` 隔离。

本地打包入口为 `src/postgres/skills/package.json`，相关脚本位于 `src/postgres/skills/scripts/`。

预期产物：
- `tencent-pg-inspection-vX.Y.Z.zip`
- `tencent-pg-slowquery-diagnosis-vX.Y.Z.zip`
- `tencent-pg-ops-troubleshooter-vX.Y.Z.zip`
- `tencentdb-postgresql-skill-vX.Y.Z.zip`

当前默认版本基线会跟随 `src/postgres/package.json` 中 PostgreSQL MCP 的 npm 版本，除非显式提供 `SKILL_VERSION`。

## 发布资源

建议的发布组织方式：
- 将三个单独的 skill 包挂到 PostgreSQL MCP 的 Release 下
- 提供一个面向大多数用户的 bundle 包
- bundle 包按参考发布格式组织为根级 `SKILL.md`、根级 `_meta.json` 与 `references/`
- `references/` 下展开三个 PostgreSQL 子 skill 目录，便于平台先识别根级入口，再继续进入子 skill
- 可选附带 skill 发布说明，用于描述兼容性与变更内容

建议的 tag 对齐方式：
- PostgreSQL MCP tag：`postgres-mcp-server-vX.Y.Z`
- Skill 资源版本：`vX.Y.Z`


## WorkBuddy / CodeBuddy 安装

建议的用户流程：
1. 从 Release 资源中下载单个 skill 的 zip 包，或下载完整 bundle zip 包用于批量查看；bundle 文件名应为 `tencentdb-postgresql-skill-vX.Y.Z.zip`。
2. 如果下载的是 bundle，先解压；根目录会直接看到 `SKILL.md`、`_meta.json` 和 `references/`，其中根级 skill 名称应显示为 `TencentDB PostgreSQL Skill`，`_meta.json` 中的 slug 应为 `tencentdb-postgresql-skill`。
3. 先阅读根级 `SKILL.md`，再进入 `references/` 中对应的子 skill 目录查看其 `SKILL.md`。
4. 如果要安装到 WorkBuddy / CodeBuddy，请使用 Release 中对应的单个 skill zip 安装包。
5. 打开技能管理。
6. 上传目标 skill 的 zip 安装包并启用。
7. 在当前会话或执行环境中，优先准备这 3 个变量：`TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY`、`TENCENTCLOUD_REGION`；临时凭证场景再补 `TENCENTCLOUD_SESSION_TOKEN`。如果历史环境里使用的是 `MCP_REQUEST_*` 或 `MCP_*` 变量名，也应兼容识别。
8. 地域既可以直接写标准值（如 `ap-guangzhou`），也可以先写常见中文地域（如 `广州`、`上海`、`成都`），执行时应先做地域归一化。
9. 使用自然语言请求调用，例如 `PG巡检`、`慢 SQL 分析`、`PG排障`。

> 其它 AI 客户端（如Cursor, Claude Code）请参照上述流程进行适配。

