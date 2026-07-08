---
name: postgresql-external-token-store
overview: 为当前 `issued-token` + SCF 多实例部署形态设计 PostgreSQL 外部持久化方案，替代本地 `sqlite /tmp` token store，消除 `bootstrap -> /sse` 间歇性 `401`。方案将覆盖存储抽象、PostgreSQL 实现、配置与文档收敛，以及部署后回归验证路径。
todos:
  - id: scan-impact
    content: 使用 [subagent:code-explorer] 复核 security 与 SCF 受影响文件
    status: pending
  - id: build-pg-store
    content: 实现 security/token_store_postgres.go 与安全 schema 初始化
    status: pending
    dependencies:
      - scan-impact
  - id: wire-store-factory
    content: 改造 main.go 与工厂，支持 sqlite/postgres 双后端
    status: pending
    dependencies:
      - build-pg-store
  - id: update-deploy-docs
    content: 更新 .env.example、README、SCF 示例为 PostgreSQL 外部存储
    status: pending
    dependencies:
      - wire-store-factory
  - id: regression-verify
    content: 使用 [skill:code-analyst] 设计并执行 bootstrap 到 SSE 稳定性回归
    status: pending
    dependencies:
      - update-deploy-docs
---

## User Requirements

- 将现有 `issued-token` 的本地持久化改为外部共享持久化，解决云函数多实例、冷启动下 token 与绑定凭证不共享的问题。
- 保持现有接口与使用方式不变，包括 `/auth/bootstrap/tencentcloud`、`/auth/token-exchange/tencentcloud`、`/sse`、`/admin/tokens`。
- 继续保留现有安全语义：不落明文 token、运行时云凭证仍以密文保存、敏感配置仅通过环境变量提供。

## Product Overview

- 服务启动时根据配置选择本地或外部持久化后端；云函数场景默认改为共享存储，确保任意实例都能完成 token 校验与凭证恢复。
- 本次改动不新增页面，不改变交互入口，主要体现为部署配置、启动方式和运行稳定性的提升。

## Core Features

- 共享 token 存储与绑定凭证存储
- 与现有鉴权、换 token、管理接口无缝兼容
- 云函数部署配置与示例同步更新
- 回归验证 `bootstrap -> /sse` 链路稳定性，消除随机 `401`

## Tech Stack Selection

- 语言与运行时：沿用现有 Go 1.25
- 服务框架：沿用当前 `net/http` + `mcp-go`
- 持久化抽象：继续复用现有 `TokenStore` / `CredentialStore`
- 数据访问：沿用 `database/sql`
- 新增外部存储驱动：`github.com/jackc/pgx/v5/stdlib`
- 本地开发兼容：保留 `modernc.org/sqlite`

## Implementation Approach

### 总体策略

基于现有抽象新增一个 PostgreSQL 版本的 store，实现与 SQLite 相同的 `TokenStore` 和 `CredentialStore` 能力；启动阶段通过工厂按 `TOKEN_STORE` 选择后端。这样 `TokenIssuer`、`issued-token` 鉴权器、`TokenBoundCredentialProvider`、`TencentCloudTokenExchangeService` 都能继续复用，无需改动 API 协议和主要业务流程。

### 关键技术决策

1. **不改 token 协议，仍用 opaque token + 服务端回查**

- 当前鉴权链路已经基于 `GetTokenByHash()` 工作，改后端即可解决跨实例问题，避免引入 JWT 改造和权限漂移风险。

2. **不仅迁移 token 元数据，也迁移凭证绑定表**

- 现有 `/sse` 后续调用依赖 `GetCredentialBinding(token_id)` 恢复运行时云凭证；只迁 token 不迁绑定表，问题仍会存在。

3. **启动装配从具体类型改为接口工厂**

- `main.go` 当前直接依赖 `*SQLiteTokenStore`，需要改为依赖组合接口，消除对 SQLite 实现的硬编码。

4. **继续沿用启动时建表初始化**

- 当前 SQLite 就是启动时初始化 schema；PostgreSQL 也复用这一模式，降低引入独立迁移框架的复杂度。

### 性能与可靠性

- `GetTokenByHash` / `GetTokenByID` / `GetCredentialBinding`：依赖主键或唯一索引，复杂度约为 `O(logN)`
- `CreateToken` / `PutCredentialBinding`：单行写入或 Upsert，复杂度约为 `O(logN)`
- `ListTokens`：按租户、主体、状态过滤后排序，复杂度约为 `O(logN + K)`
- 主要热点：
- **认证回查**：每次 Bearer 校验都会命中数据库，依赖 `token_hash` 唯一索引
- **TouchToken 写放大**：每次认证后的 `last_used_at` 更新可能带来额外写入；首版保持兼容语义，如压测发现瓶颈，可在 PostgreSQL 实现内增加时间窗节流
- **SCF 连接数**：通过保守连接池参数控制每个 warm 实例的数据库占用，避免外部数据库被瞬时打满

## Implementation Notes

- 全部 SQL 使用参数绑定，禁止字符串拼接值参数
- `TOKEN_STORE_SCHEMA` 如支持可配置，必须先做白名单校验，只允许字母、数字、下划线，防止动态标识符注入
- 严禁记录明文 token、DSN、数据库密码、加密密钥；启动日志只输出脱敏后的 store 类型/标识
- `TOKEN_STORE=sqlite` 保持向后兼容；`TOKEN_STORE=postgres` 且缺少 DSN 时立即 fail fast
- 继续仅存 `HashToken(token, pepper)`，不落明文 token；凭证继续使用现有 `CREDENTIAL_ENCRYPTION_KEY` 加密后的密文
- 不做 SCF 历史 token 自动迁移；升级后要求客户端重新 bootstrap，符合当前 `/tmp` 临时存储现状
- 避免无关重构，不改 MCP 工具注册、Guard、鉴权协议和 HTTP 路由

## Architecture Design

### 系统结构

- `main.go`
- 读取环境变量
- 通过 store factory 初始化 sqlite 或 postgres
- 将同一份 store 注入：
    - `TokenIssuer`
    - `issuedTokenAuthenticator`
    - `TokenBoundCredentialProvider`
    - `TencentCloudTokenExchangeService`

### 数据流

1. 用户调用 `/auth/bootstrap/tencentcloud`
2. `TencentCloudTokenExchangeService.Exchange()` 验证身份并签发 token
3. `CreateToken()` 写入外部存储
4. `PutCredentialBinding()` 写入绑定凭证密文
5. 客户端携带 Bearer token 访问 `/sse`
6. `issuedTokenAuthenticator` 用 `GetTokenByHash()` 校验 token
7. `TokenBoundCredentialProvider` 用 `GetCredentialBinding()` 恢复凭证
8. 工具调用在任意实例都可完成

### 数据库模型

- `issued_tokens`
- 保留当前字段语义：`id`、`token_hash`、`token_prefix`、`subject_*`、`tenant_id`、`scopes`、`allowed_regions`、状态、过期时间、审计字段
- `issued_token_credentials`
- 以 `token_id` 关联保存加密后的腾讯云凭证、凭证类型、过期时间、审计字段
- 推荐索引
- `issued_tokens(token_hash)` 唯一索引
- `issued_tokens(subject_id)`
- `issued_tokens(tenant_id)`
- `issued_tokens(status)`
- `issued_tokens(created_at DESC)` 或结合查询模式建立复合索引
- `issued_token_credentials(token_id)` 主键/唯一索引

## Directory Structure Summary

本次实现以“新增 PostgreSQL store + 启动工厂 + SCF 文档/示例切换”为主，尽量不触碰现有业务逻辑。

```text
/Users/lumenyang/workspace/tencentcloud-mcp-server/src/postgres/
├── go.mod                                # [MODIFY] 引入 PostgreSQL database/sql 驱动，保留 sqlite 兼容
├── main.go                               # [MODIFY] 改为通过工厂装配 store；移除对 *SQLiteTokenStore 的硬依赖；启动日志脱敏
├── security/
│   ├── token_store_factory.go            # [NEW] 解析 TOKEN_STORE、TOKEN_STORE_DSN、TOKEN_STORE_SCHEMA、连接池参数；创建并返回通用 store
│   └── token_store_postgres.go           # [NEW] PostgreSQL 版 TokenStore + CredentialStore；建表、索引、CRUD、Upsert、扫描映射
├── .env.example                          # [MODIFY] 增加 postgres 外部存储配置示例，明确 sqlite 仅适合本地/单实例
├── README.md                             # [MODIFY] 更新 issued-token 最小配置、外部持久化方案、切换步骤与验证方式
├── SCF_DEPLOY.md                         # [MODIFY] 将 SCF 推荐方案更新为外部共享存储，补充回归验证与风险说明
└── deploy/scf/
    ├── scf_bootstrap                     # [MODIFY] 去掉对 sqlite 和 /tmp 的强默认，改为依赖环境变量注入外部存储
    ├── scf.console.startup.sh            # [MODIFY] 同步控制台启动脚本示例，避免硬编码 sqlite
    ├── scf.env.example                   # [MODIFY] 提供 SCF 下外部存储环境变量模板
    └── scf.console.env.txt               # [MODIFY] 提供控制台可直接录入的外部存储示例
```

## Key Code Structures

- 新增组合 store 工厂，返回同时满足 `TokenStore` 与 `CredentialStore` 的实例
- 新增 PostgreSQL 配置项建议：
- `TOKEN_STORE=postgres`
- `TOKEN_STORE_DSN=postgres://...`
- `TOKEN_STORE_SCHEMA=mcp_auth`
- `TOKEN_STORE_MAX_OPEN_CONNS`
- `TOKEN_STORE_MAX_IDLE_CONNS`
- `TOKEN_STORE_CONN_MAX_LIFETIME_SECONDS`
- 兼容策略：
- sqlite 继续使用 `TOKEN_STORE_PATH`
- postgres 不再依赖本地路径

## Agent Extensions

### SubAgent

- **code-explorer**
- Purpose: 复核 `security`、`deploy/scf`、文档与启动入口的受影响文件，避免遗漏装配点和部署示例
- Expected outcome: 得到完整、可落地的改动清单和调用链确认结果

### Skill

- **code-analyst**
- Purpose: 复核 `TokenStore` / `CredentialStore` 调用链、热路径与兼容性风险
- Expected outcome: 确认 PostgreSQL 实现与现有鉴权、token 交换、动态凭证恢复逻辑一致