# PostgreSQL 配套技能包

本目录存放 PostgreSQL 模块对应的 WorkBuddy / CodeBuddy 配套技能，路径位于 `src/postgres/skills/`。

这些资源属于**工作流层的技能包**，不是独立的数据库连接器。它们不再依赖已部署的 PostgreSQL MCP Server，而是直接调用腾讯云 PostgreSQL OpenAPI；但允许调用的 Action 必须严格限制在当前 `src/postgres` 已完成参数对齐与工具封装的范围内。

## 概览

- **同一领域，不同层次**：这些 skill 与 PostgreSQL MCP 实现放在相邻目录中，但与 Go 运行时、SCF 部署脚本以及 npm 启动器相互隔离。
- **证据优先、全量开放**：这些打包后的 skill 可以调用当前已完成参数契约校验的全部 **48 个 PostgreSQL OpenAPI Action**；默认先做证据采集，再按需进入变更或处置动作。
- **OpenAPI 直连**：skill 在运行时直接调用腾讯云 PostgreSQL OpenAPI，不要求用户额外先部署 MCP Server。
- **对齐边界明确**：允许使用的 Action 必须是当前仓库已经完成参数契约校验的那一组，源码基线以 `tools/openapi_alignment.go` 为准。
- **高风险动作需明确确认**：涉及审计类、业务类、费用类或 critical 级变更动作时，必须先说明影响面并获得明确确认。
- **面向发布的分发方式**：用户应从 GitHub Release 或 COS 安装打包好的 zip 资源，而不是直接使用源码目录。

## 用户傻瓜式配置（推荐）

对最终用户来说，**推荐只准备 3 个环境变量**，就可以直接使用这些 PostgreSQL skill：

```bash
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
# 可选：临时凭证场景再补 TENCENTCLOUD_SESSION_TOKEN
```

补充约定：
- **推荐变量名**：`TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY`、`TENCENTCLOUD_REGION`
- **兼容变量名**：也兼容读取 `MCP_REQUEST_SECRET_ID`、`MCP_REQUEST_SECRET_KEY`、`MCP_REQUEST_SESSION_TOKEN`、`MCP_SECRET_ID`、`MCP_SECRET_KEY`
- **地域可写自然语言**：`广州`、`上海`、`成都`、`北京` 这类常见中文地域，执行 skill 时应先归一化为标准地域码，例如 `广州 -> ap-guangzhou`、`上海 -> ap-shanghai`、`成都 -> ap-chengdu`
- **SDK 不是用户首个阻断项**：skill 执行时应优先使用官方 SDK；如果运行环境里没有 SDK，应退回到本地生成 TC3 签名的 HTTPS 请求，而不是直接让用户先去手装 SDK
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

这些 skill 不再把 PostgreSQL MCP Server 当作运行前置条件，但仍然必须与 `src/postgres` 的 OpenAPI 对齐结果保持一致。

当前约束如下：
- 允许调用的 Action 必须来自当前已对齐的 PostgreSQL OpenAPI 集合，维护基线以 `tools/openapi_alignment.go` 为准。
- 仓库根目录 `README.md` 的“**2. MCP 开放能力**”章节列出了当前对齐完成的 48 个 Action，可作为人工核对清单。
- 这三个工作流 skill 现在**都允许使用这 48 个已对齐 Action**；但应优先从与当前任务最相关的模块开始，而不是无差别全量调用。
- 即使腾讯云官方还提供更多 PostgreSQL OpenAPI，**只要当前仓库未对齐，就不应在这些 skill 中调用**。

这些 skill **应该**做到：
- 直接对接腾讯云 PostgreSQL OpenAPI
- 优先使用官方 SDK；无法使用 SDK 时，再用本地生成的 TC3 签名 HTTPS 请求
- 只从环境变量或其他安全运行时上下文中读取 `SecretId` / `SecretKey` / 可选 STS Token
- 在最终输出中明确区分“已验证证据”和“推测结论”
- 在执行审计类、业务类、费用类或 critical 级动作前，明确说明影响面、目标实例和确认要求

这些 skill **不应该**做到：
- 在代码、skill 文档或仓库文件中硬编码密钥
- 把密钥放进 URL 或 query 参数
- 越过当前 MCP 对齐边界调用未校验的 PostgreSQL OpenAPI
- 在没有明确确认的情况下执行高风险写操作

## 当前开放的 48 个对齐 Action

以下表格可直接复制到 README、发布说明或变更公告。

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

**不要**把用户引导到原始源码目录或开发分支文件链接。

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

## 新用户验收测试

本目录提供了自动化测试脚本，用于**模拟一个完全新用户的环境**验收 skill zip 包质量。不需要 Docker，不需要卸载本机 SDK。

### 原理

测试脚本自动创建一个临时 Python venv（空包环境），在这个 venv 里：
- 没有 `tencentcloud` SDK（模拟新用户未安装 SDK）
- 没有 `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY` / `TENCENTCLOUD_REGION` 环境变量
- HOME / XDG 目录指向临时路径（不会读到用户真实配置）

所有需要"新用户环境"的测试都在这个 venv 的子进程里执行，**不影响本机任何东西**，测试完自动清理。

### 使用方式

```bash
cd src/postgres/skills

# 默认：本地沙箱模式（推荐）
python3 tests/test_new_user.py

# 保留沙箱 venv 以便事后调试
python3 tests/test_new_user.py --keep-sandbox

# 带真实凭证做端到端 API 连通性测试（在沙箱内用 TC3 签名，不依赖 SDK）
TENCENTCLOUD_SECRET_ID=xxx TENCENTCLOUD_SECRET_KEY=xxx TENCENTCLOUD_REGION=ap-guangzhou \
  python3 tests/test_new_user.py --live
```

### 测试覆盖

| # | 类别 | 检查内容 |
|---|------|---------|
| 1 | zip 包存在性 | 4 个预期 zip 包是否存在 |
| 2 | 单 skill 结构 | SKILL.md + references + common/ 公共规则 |
| 3 | bundle 结构 | 根级入口 + 3 个子 skill + 公共规则 |
| 4 | 错误模板 | missing-credentials / invalid-region / missing-sdk 模板 + 官方链接 |
| 5 | 地域规则 | 中文别名映射表 + 官方地域/DescribeRegions 文档链接 |
| 6 | SDK 缺失 | 沙箱 venv 内确认无法 import tencentcloud |
| 7 | 凭证检测 | 零凭证→报缺失 / 部分凭证→精确报缺失 / 兼容变量名→可识别 |
| 8 | 端到端 API | （需 --live）沙箱内 TC3 签名直连 OpenAPI，验证连通性 |

### AI 客户端集成测试

除了上面的新用户验收测试（侧重 zip 打包质量和沙箱隔离），本目录还有**AI 客户端集成测试**，模拟 AI 客户端（如 CodeBuddy）加载 skill 后按工作流执行的完整链路。

```bash
# 默认：沙箱内跑 S1-S5（S4 需 --live）
python3 tests/test_ai_client_integration.py

# 主进程模式：禁用本地环境变量 + 阻塞 SDK，直接在主进程验证错误处理路径
python3 tests/test_ai_client_integration.py --main-process

# 含真实 API 调用（验证完整 S4 工作流）
TENCENTCLOUD_SECRET_ID=xxx TENCENTCLOUD_SECRET_KEY=xxx \
  python3 tests/test_ai_client_integration.py --live

# 保留沙箱调试
python3 tests/test_ai_client_integration.py --keep-sandbox --live
```

**测试场景覆盖：**

| # | 场景 | 模拟什么 |
|---|------|---------|
| S1 | 零凭证 | 用户第一次使用 skill，没有任何凭证 → 验证 AI 返回 missing-credentials 模板（含控制台链接 + 可复制命令） |
| S2 | 非法 region | 用户输入 `华南`、`ap-gz`、`guangzou` 等 → 验证 AI 拒绝并给出 invalid-region 引导 |
| S3 | SDK 缺失 + fallback | 沙箱内无 SDK → 验证 AI 说明"TC3 fallback 可用，不阻断执行" |
| S4 | 完整工作流 | 有效参数 → 凭证检测 → 地域归一化 → SDK 决策 → API 调用 → 输出结构验证（需 --live） |
| S5 | 跨 skill 一致性 | 三个 skill 在凭证检测、地域归一化、SDK fallback、错误模板引用、Verification 输出规范上的一致 |

**与 `test_new_user.py` 的区别：**

| 维度 | test_new_user.py | test_ai_client_integration.py |
|------|-----------------|------------------------------|
| 侧重 | zip 包文档结构 | AI 按工作流执行的行为 |
| 测试对象 | 打包产物完整性、错误模板内容 | 工作流步骤行为、错误路径触发、输出格式 |
| API 调用 | 单一 DescribeRegions | 工作流步骤链 + DescribeRegions |
| 跨 skill | 不涉及 | 三个 skill 一致性对比 |

**三种运行模式对比：**

| 模式 | 环境变量 | SDK 隔离 | 执行方式 | 适用场景 |
|------|---------|----------|---------|---------|
| 默认(sandbox) | 自动隔离(空) | venv 无 SDK | 子进程 | CI/本地验收 |
| `--main-process` | 临时清除 | `sys.meta_path` 阻塞 + `sys.modules` 清理 | 主进程 `exec()` | 快速本地调试、验证主进程行为 |
| `--no-venv` | 仅隔离 env | 不清除 | 子进程 | 快速但不完全隔离 |

`--main-process` 的实现原理：
1. **环境变量**：遍历 `os.environ`，暂停所有 `TENCENTCLOUD_*` / `MCP_SECRET_*` / `MCP_REQUEST_*` 前缀的变量
2. **SDK 阻塞**：通过 `sys.meta_path` 注入自定义 Finder，任何 `import tencentcloud*` 直接 `raise ImportError`；同时清理 `sys.modules` 中已缓存的 tencentcloud 模块
3. **代码执行**：`DirectRunner` 使用 `exec()` + `io.StringIO` 在主进程直接执行，不走子进程
4. **恢复**：`contextmanager` 退出时自动恢复环境变量、sys.modules、sys.meta_path

## 版本兼容性

每次 PostgreSQL MCP 的工具名、参数或安全边界发生变化时，都应重新评估兼容性。

建议规则：
- 以 `tools/openapi_alignment.go` 和仓库根 `README.md` 的能力清单作为主要兼容性基线
- skill 包版本默认跟随 MCP Release，除非有明确记录的版本映射
- 任何 PostgreSQL OpenAPI 对齐范围的破坏性变更，都必须在发布前触发一次配套 skill 审查

## 维护规则

- 除非另有单独评审通过的写操作 skill，否则配套 skill 应默认保持只读。
- 当 PostgreSQL OpenAPI 的对齐范围、参数映射或安全边界发生变化时，需要审查 `skills/` 下的每一个 skill。
- 不要把打包生成的 zip 资源提交到源码目录中。
- 保持打包流程与主 Go 构建和 SCF 部署流程隔离。
- 前端下载入口应使用 Release 资源或 COS URL，而不是仓库原始页面。
