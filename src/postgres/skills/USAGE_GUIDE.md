# PostgreSQL 配套技能包 — 使用指南

本文档介绍如何安装和使用TencentDB PostgreSQL Skill的三大技能。

---

## 1. 概述

本技能包包含三个 skill，均**无需部署 MCP Server**，直接调用腾讯云 PostgreSQL OpenAPI，覆盖日常运维中最常见的三类场景：

| 技能名称 | 中文名称 | 用途 |
|----------|----------|------|
| `tencent-pg-inspection` | PG 巡检 | 日常健康巡检，全面检查实例状态、备份、网络、SSL、参数等 |
| `tencent-pg-slowquery-diagnosis` | 慢 SQL 诊断 | 慢查询分析、SQL 性能诊断、延迟排查 |
| `tencent-pg-ops-troubleshooter` | 运维排障 | 按故障场景分类排查：连接、备份、账号、网络、SSL、只读实例等 |

---

## 2. 安装前准备

### 2.1 获取腾讯云 API 密钥

1. 打开 [腾讯云 API 密钥管理](https://console.cloud.tencent.com/cam/capi)
2. 创建或查看 `SecretId` 和 `SecretKey`
3. 建议为 AI 客户端单独创建 CAM 子账号，按需授予最小权限（推荐只读权限用于巡检和诊断）

### 2.2 设置环境变量

在你的终端或 AI 客户端运行环境中设置以下三个环境变量：

```bash
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
```

> **地域支持中文**：除标准地域码（如 `ap-guangzhou`）外，也支持写中文地域名，如 `广州`、`上海`、`北京`、`成都`。技能执行时会自动归一化为标准地域码。

**兼容的变量名**（以下均可识别）：

| 变量用途 | 推荐变量名 | 兼容变量名 |
|----------|-----------|-----------|
| SecretId | `TENCENTCLOUD_SECRET_ID` | `MCP_SECRET_ID`、`MCP_REQUEST_SECRET_ID` |
| SecretKey | `TENCENTCLOUD_SECRET_KEY` | `MCP_SECRET_KEY`、`MCP_REQUEST_SECRET_KEY` |
| Region | `TENCENTCLOUD_REGION` | — |
| Token（可选） | `TENCENTCLOUD_SESSION_TOKEN` | `MCP_REQUEST_SESSION_TOKEN` |

> **安全提示**：不要把 `SecretId`、`SecretKey` 写入代码仓库、skill 文档、URL 或 query 参数中。仅通过环境变量传递。

---

## 3. 安装技能到 CodeBuddy / WorkBuddy

### 3.1 下载安装包

从 GitHub Release 下载以下任一 zip 包：

- **单技能安装**：`tencent-pg-inspection-vX.Y.Z.zip`、`tencent-pg-slowquery-diagnosis-vX.Y.Z.zip`、`tencent-pg-ops-troubleshooter-vX.Y.Z.zip`
- **全家桶安装**：`tencentdb-postgresql-skill-vX.Y.Z.zip`（包含全部三个技能）

### 3.2 导入技能

**CodeBuddy**：
1. 打开 CodeBuddy 设置 → 技能管理
2. 点击「导入技能」
3. 选择下载的 `.zip` 安装包
4. 启用技能

**WorkBuddy**（基于 CodeBuddy 内核，流程一致）：
1. 打开技能管理面板
2. 上传对应的 `.zip` 安装包
3. 启用目标技能

### 3.3 验证安装

安装完成后，在聊天窗口输入以下任意一句话验证：

- `PG巡检` → 应触发 `tencent-pg-inspection`
- `帮我分析慢SQL` → 应触发 `tencent-pg-slowquery-diagnosis`
- `PG排障` → 应触发 `tencent-pg-ops-troubleshooter`

首次使用时技能会提示你提供目标区域和实例 ID，之后便会自动调用 OpenAPI 进行证据采集。

---

### 3.4 通过 CLI 安装（CodeBuddy Code / `codebuddy` 命令行）

如果你使用的是 **CodeBuddy Code CLI**（`codebuddy` 命令行工具），有以下三种方式安装技能：

#### 方式 A：手动放置（推荐，离线兼容）

将技能文件夹放到对应目录后重启 CLI 即可自动加载：

```bash
# 用户级（所有项目生效）
mkdir -p ~/.codebuddy/skills
cp -r tencent-pg-inspection ~/.codebuddy/skills/
cp -r tencent-pg-slowquery-diagnosis ~/.codebuddy/skills/
cp -r tencent-pg-ops-troubleshooter ~/.codebuddy/skills/

# 项目级（仅当前项目生效）
mkdir -p .codebuddy/skills
cp -r tencent-pg-inspection .codebuddy/skills/
```

> 技能文件夹即包含 `SKILL.md` 的整个目录，直接从下载的 zip 包解压得到。

目录结构说明：

```text
~/.codebuddy/skills/
├── tencent-pg-inspection/
│   ├── SKILL.md
│   ├── references/
│   └── assets/
├── tencent-pg-slowquery-diagnosis/
│   ├── SKILL.md
│   ├── references/
│   └── assets/
└── tencent-pg-ops-troubleshooter/
    ├── SKILL.md
    ├── references/
    └── assets/
```

#### 方式 B：拖拽导入（CodeBuddy IDE 窗口内）

下载 `.zip` 安装包后，将文件直接拖入 CodeBuddy 终端窗口，CLI 会自动检测并安装。

#### 方式 C：从 SkillHub 一键安装

打开 [SkillHub](https://skillhub.tencent.com) 搜索 `tencent-pg`，一键添加到 CodeBuddy。

### 3.5 我不是 CodeBuddy 用户，能用吗？

**可以，但不需要安装这 3 个技能。** 直接使用 **PostgreSQL MCP Server** 即可，效果等价。

技能本质是对 MCP Server 的 48 个 OpenAPI Action 做了场景化封装。对于其他 MCP 客户端（Cursor、Claude Desktop 等），你应该直接接入 MCP Server 的 `stdio` 或 `streamable-http` 传输模式：

#### Stdio 方式（本地 CLI 客户端）

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "/path/to/mcp-server/src/postgres/scripts/run_stdio.sh",
      "env": {
        "MCP_REQUEST_SECRET_ID": "<你的 SecretId>",
        "MCP_REQUEST_SECRET_KEY": "<你的 SecretKey>"
      }
    }
  }
}
```

#### Streamable-HTTP 方式（远程部署）

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "https://你的函数URL/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<你的 SecretId>",
        "X-TencentCloud-Secret-Key": "<你的 SecretKey>"
      }
    }
  }
}
```

配置好后，在对话中用自然语言描述需求即可——比如 `列出广州的 PostgreSQL 实例`、`检查 postgres-abc123 的慢查询`、`分析备份状态`，效果与使用技能相同。

> 完整部署文档见仓库 [`README.md`](../README.md) 的 [「3. 选择部署方式」](../README.md#3-选择部署方式) 章节。

| 对比维度 | 技能（CodeBuddy） | MCP Server（通用） |
|----------|--------------------|--------------------|
| 适用客户端 | CodeBuddy IDE / WorkBuddy / CLI | Cursor / Claude Desktop / 所有 MCP 客户端 |
| 安装方式 | zip 导入 / 目录放置 | JSON 配置 + 启动命令 |
| 运行依赖 | 环境变量 + 腾讯云 OpenAPI | Go 运行时 + MCP Transport |
| 交互方式 | 自然语言触发（`PG巡检`） | 调用 MCP 工具（48 个 Action） |
| OpenAPI 覆盖 | 全部 48 个已对齐 Action | 全部 48 个已对齐 Action |
| 推荐场景 | 日常运维、诊断（开箱即用） | 需要精细控制工具调用 |

---

## 4. 技能使用详解

### 4.1 PG 巡检（`tencent-pg-inspection`）

**触发关键词**：`PG巡检`、`健康检查`、`备份检查`、`资源水位检查`

**使用示例**：

| 你说的 | 实际效果 |
|--------|---------|
| `PG巡检 ap-guangzhou postgres-abc123` | 对广州的指定实例执行全量健康巡检 |
| `检查一下广州 postgres-abc123 的健康状态` | 同上（支持中文地域） |
| `帮我巡检备份状态` | 提示你先提供区域和实例 ID |
| `看看 postgres-xyz 的 SSL 配置` | 检查 SSL 配置和网络暴露情况 |

**巡检覆盖**：
1. **实例基础**：生命周期状态、规格、版本、当前任务
2. **备份健康**：基础备份概览、日志备份概览、可恢复时间范围
3. **参数姿态**：当前参数配置、参数修改事件
4. **网络/安全**：安全组配置、公网访问状态、SSL 配置
5. **只读拓扑**：只读组和只读实例状态（如适用）
6. **可选深层检查**：当日志/延迟相关线索出现时，补充慢查询和错误日志证据

**输出格式**：

巡检完成后会生成四块结构化摘要：

```
一、总体状态：✅ 健康 / ⚠️ 有风险 / ❌ 异常
二、证据快照：各类检查项的 OpenAPI 返回要点
三、风险项：按严重程度排列的风险列表
四、建议下一步：无需操作 / 推荐操作 / 需确认的操作
```

---

### 4.2 慢 SQL 诊断（`tencent-pg-slowquery-diagnosis`）

**触发关键词**：`慢SQL分析`、`SQL性能诊断`、`查询为什么变慢`、`帮我看看慢查询`

**使用示例**：

| 你说的 | 实际效果 |
|--------|---------|
| `慢SQL分析 广州 postgres-abc123` | 拉取近期的慢查询列表并分析 |
| `帮我看看 postgres-abc123 的查询为什么变慢了，区域广州` | 采集慢查询 + 错误日志 + 实例上下文综合诊断 |
| `昨天下午 3 点到 5 点，广州 postgres-abc123 有哪些慢SQL` | 指定时间窗口的慢查询分析 |

**诊断流程**：

1. **确认范围**：识别区域、实例、可选数据库/时间窗口
2. **采集证据**（按序）：
   - 实例当前状态、生命周期、正在执行的任务
   - 慢查询列表 + 慢查询分析（按耗时排序）
   - 错误日志（当错误可能与延迟相关时）
   - 参数、备份、只读拓扑（仅在关联性分析需要时）
3. **根因排序**：按可能性从高到低列出候选原因，而非强制给出唯一根因
4. **优化建议**：仅在有充分证据时给出保守建议或执行操作

**输出格式**：

```
一、诊断窗口：起止时间、目标实例
二、证据摘要：慢查询 Top N、错误日志关联、实例负载
三、可能性排序：原因 1（最可能）→ 原因 2 → 原因 3
四、建议下一步：安全操作 / 需确认操作 / 进一步排查方向
```

---

### 4.3 运维排障（`tencent-pg-ops-troubleshooter`）

**触发关键词**：`PG排障`、`实例异常排查`、`SSL问题`、`备份失败`、`连接不上`

**使用示例**：

| 你说的 | 实际效果 |
|--------|---------|
| `PG排障，广州 postgres-abc123 连不上了` | 分类为连接故障 → 检查实例状态、网络、安全组 |
| `备份失败了，postgres-abc123 广州` | 分类为备份问题 → 检查备份配置和最近备份状态 |
| `SSL 问题排查，广州 postgres-abc` | 聚焦 SSL 配置检查和网络暴露分析 |
| `实例异常排查，广州 postgres-abc` | 通用实例异常 → 全面收集实例、任务、日志证据 |

**故障分类与对应检查模块**：

| 用户症状 | 分类 | 优先检查模块 |
|----------|------|-------------|
| 连接不上、访问失败 | 连接/访问故障 | 实例状态 → 安全组 → 公网开关 → SSL |
| 备份失败、备份超时 | 备份/恢复问题 | 备份概览 → 基础备份列表 → 日志备份列表 |
| 账号权限问题 | 账号/权限问题 | 账号列表 → 账号权限 → 数据库对象 |
| 网络不可达、SSL 错误 | 网络/SSL 问题 | SSL 配置 → 安全组 → 公网访问 |
| 只读实例异常 | 只读/复制问题 | 只读组 → 实例状态 → 参数 |
| 参数错误、配置异常 | 参数/配置问题 | 实例参数 → 参数事件 → 参数模板 |
| 其他未知异常 | 通用异常 | 实例详情 → 任务列表 → 错误日志 → 慢查询 |

**输出格式**（Runbook 风格）：

```
一、故障摘要：故障分类 + 影响范围
二、检查发现：按模块列出各检查项结果
三、阻塞项：当前最需要解决的问题
四、立即可执行的下一步：安全步骤 + 需确认步骤
```

---

## 5. 常见场景速查

### 5.1 我想全面了解一个实例的健康状况

```
PG巡检 ap-guangzhou postgres-abc12345
```

### 5.2 我发现查询很慢，想排查原因

```
慢SQL分析，广州 postgres-abc12345，最近 1 小时
```

### 5.3 实例突然连不上了，紧急排查

```
PG排障，广州 postgres-abc12345，连不上了
```

### 5.4 想看某个数据库有哪些账号和权限

```
查询广州 postgres-abc12345 的账号和权限
```

### 5.5 想检查备份是否正常

```
PG巡检广州 postgres-abc12345 的备份状态
```

### 5.6 某条 SQL 一直很慢，想知道为什么

```
帮我分析广州 postgres-abc12345 的这条慢SQL：[粘贴 SQL 或 SQL 指纹]
```

---

## 6. 常见问题

### Q1：技能提示"缺少 SecretId / SecretKey"怎么办？

确认环境变量已正确设置：

```bash
echo $TENCENTCLOUD_SECRET_ID
echo $TENCENTCLOUD_SECRET_KEY
```

如果没有输出，按照第 2 节重新设置环境变量，然后**重启 AI 客户端**。

### Q2：技能提示"无效的地域"怎么办？

支持的地域写法：
- **标准格式**：`ap-guangzhou`、`ap-shanghai`、`ap-beijing`、`ap-chengdu`
- **中文格式**：`广州`、`上海`、`北京`、`成都`

其他格式（如 `华南`、`ap-gz`、`guangzhou`）暂不支持，请使用上述标准格式。

### Q3：使用的是什么 Python 环境？

技能会优先使用系统已安装的腾讯云 SDK（`tencentcloud-sdk-python`）。如果未安装，技能会自动退回到本地 TC3 签名直连 OpenAPI，不会因为你没装 SDK 就阻断执行。

如果你希望使用 SDK 路径，可以手动安装：

```bash
pip install tencentcloud-sdk-python
```

### Q4：切换账号后查询结果好像不对？

这是 AI 客户端的连接/会话缓存问题，不是服务端问题。解决方法：
1. **为不同账号创建不同的连接器/配置**，不要反复修改同一个配置
2. 切换账号后**重启 AI 客户端**
3. 首次使用时先用 `DescribeDBInstances` 做最小验证

### Q5：技能会执行危险操作吗？

不会自动执行。所有涉及写操作、费用变更或高风险动作的行为，技能都会：
1. 先说明影响面
2. 明确列出目标实例
3. 要求你**明确确认**后才执行

查询类操作（如 `DescribeDBInstances`、`DescribeSlowQueryList`）则是直接可用，无需额外确认。

### Q6：支持哪些 PostgreSQL OpenAPI Action？

共计 48 个已对齐 Action，完整列表见技能目录下的 `@references/api_reference.md`，或仓库根目录 `README.md` 的「2. MCP 开放能力」章节。

---

## 7. 注意事项

1. **密钥安全**：`SecretId` / `SecretKey` 只应通过环境变量注入，不要写入任何配置文件或代码仓库
2. **网络要求**：技能直接调用腾讯云 OpenAPI 公网接口，需要 AI 客户端运行环境能访问公网
3. **地域确认**：首次对话建议直接带上完整的 `区域 + 实例ID`，避免技能反复提示确认范围
4. **权限最小化**：建议为 AI 客户端子账号只授予 `QcloudPostgresReadOnlyAccess` 策略用于巡检和诊断
5. **版本更新**：技能包版本跟随 MCP Server 发版，升级时重新下载最新 zip 包导入即可覆盖

---

## 8. 版本信息

- **当前支持的 OpenAPI Action**：48 个
- **技能包版本**：跟随 `postgres-mcp-server` 发版
- **分发方式**：GitHub Release + COS 下载

---

如有问题或建议，请联系 PostgreSQL 云数据库团队。
