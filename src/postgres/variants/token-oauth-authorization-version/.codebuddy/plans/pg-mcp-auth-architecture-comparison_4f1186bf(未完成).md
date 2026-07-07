---
name: pg-mcp-auth-architecture-comparison
overview: 重建当前方案评估计划，比较 `issued-token + 共享存储` 与 `无状态 CAM + AssumeRole` 两种认证架构的优劣、适用场景、改造成本与推荐落地路径。
todos:
  - id: scan-auth-chain
    content: 使用 [subagent:code-explorer] 复核鉴权、凭证恢复与路由受影响文件
    status: pending
  - id: adr-compare
    content: 形成两种架构 ADR，对比兼容性、安全、性能与运维成本
    status: pending
    dependencies:
      - scan-auth-chain
  - id: shared-store-path
    content: 若选主线方案，落地外部 TokenStore 与 CredentialStore 共享持久化
    status: pending
    dependencies:
      - adr-compare
  - id: stateless-cam-path
    content: 若选备选方案，新增 cam-assume-role 认证与请求级动态 AssumeRole
    status: pending
    dependencies:
      - adr-compare
  - id: verification-gates
    content: 使用 [skill:code-analyst] 建立两路线回归矩阵与上线门禁
    status: pending
    dependencies:
      - shared-store-path
      - stateless-cam-path
---

## User Requirements

- 重写当前实施方向，先产出一份“架构对比 + 决策建议”的新方案，而不是直接沿用原有共享存储改造计划。
- 重点比较两条路线：
- 现有思路延伸的“issued-token + 共享存储”
- 新提出的“无状态 CAM + AssumeRole”
- 对比内容需覆盖：跨账号、多实例与冷启动、客户端接入方式、权限控制、吊销与审计、部署与运维复杂度、迁移成本与风险。

## Product Overview

- 方案输出需要明确两种接入模式对现有访问入口、客户端配置方式、部署模式和运维手段的影响。
- 不新增页面，主要变化体现在接入链路、认证方式、部署配置和故障定位方式上。

## Core Features

- 两种认证架构的优劣势对比矩阵
- 面向当前仓库的主推荐路线与备选路线
- 分阶段演进路径、边界约束与回归验证重点

## Tech Stack Selection

- 继续沿用当前项目技术栈：
- Go 1.25
- `net/http` + `mcp-go`
- 现有 `security` 抽象：`Authenticator`、`Principal`、`CredentialProvider`、`TokenStore`、`CredentialStore`
- 腾讯云 STS：`GetCallerIdentity`、`AssumeRole`
- 对比方案均基于已验证的现有代码链路展开：
- `main.go` 负责启动装配与路由注册
- `security/http_auth.go` 负责数据面鉴权
- `security/principal.go` 与 `tools/registry.go` 负责工具级授权
- `security/token_exchange_tencentcloud.go` 与 `security/credential_dynamic.go` 代表当前 issued-token 模型

## Implementation Approach

### 总体策略

先基于现有代码产出一份 ADR 式决策方案，再按选定路线推进实施。
两条路线不是“同级小改”，而是“增量演进”与“架构切换”的关系：

- **路线 A：issued-token + 共享存储**
- 保留当前 `/auth/bootstrap/tencentcloud`、`/auth/token-exchange/tencentcloud`、`/sse`、`/admin/tokens`
- 把当前本地 `TokenStore` / `CredentialStore` 改为共享后端
- 这是对当前实现的最小侵入增强

- **路线 B：无状态 CAM + AssumeRole**
- 新增独立 `auth mode`，在每次请求现场确认调用者身份，并现场 `AssumeRole`
- 不再依赖 issued-token、token 绑定凭证、管理面 token 列表
- 这是对当前认证和客户端接入模型的重新设计

### 核心对比结论

| 维度 | issued-token + 共享存储 | 无状态 CAM + AssumeRole |
| --- | --- | --- |
| 与当前代码贴合度 | 最高，直接复用现有 token exchange、principal、动态凭证恢复链路 | 中低，需要新增认证器和请求级凭证链路 |
| MCP SSE 客户端兼容性 | 最好，当前就是 Bearer 访问 `/sse` | 最弱，通用客户端通常不支持每请求携带 CAM 可验证身份 |
| 多实例 / SCF | 依赖共享存储后可稳定运行 | 天然无会话共享问题 |
| 跨账号 | 已支持 bootstrap 后 `AssumeRole`，后续靠绑定凭证执行 | 每次请求可按调用者现场 `AssumeRole`，账号边界更直接 |
| 安全面 | 客户端后续只持有本地访问 token，源身份材料不重复提交 | 服务端不存会话态，但客户端或网关需持续提供可验证身份 |
| 吊销 / 审计 | 强，已有 `/admin/tokens`、token 状态、使用时间等模型 | 弱，需要依赖 CAM/角色策略、租户映射和额外审计手段 |
| 性能热点 | 热点在 DB 回查与 `TouchToken` 写入 | 热点在 STS 身份校验与 `AssumeRole` 网络调用 |
| 改造成本 | 低到中 | 高 |
| 上线风险 | 可控，协议不变 | 高，客户端与网关约束明显变化 |


### 推荐决策

#### 主推荐路线

**优先选“issued-token + 共享存储”作为生产主线。**

原因：

1. 当前仓库已经把 Bearer token + 服务端回查作为主模型
2. `tools/registry.go` 依赖 `Principal` 做细粒度授权，现有链路可直接延续
3. 文档、脚本、客户端示例都以 `/auth/bootstrap/tencentcloud` 返回的 Bearer 配置为中心
4. 这是唯一能在**最小改动下继续兼容通用 MCP SSE 客户端**的路线

#### 备选路线

**“无状态 CAM + AssumeRole”更适合作为受控环境 PoC 或专用接入模式。**

适用前提：

- 客户端或前置网关可持续注入 CAM 可验证身份
- 接受 `/auth/bootstrap/tencentcloud` 与 `/admin/tokens` 不再是主入口
- 接受更高的 STS 调用延迟、配额与调试复杂度
- 接受客户端配置与运维手册重写

### 性能与可靠性

#### 路线 A：issued-token + 共享存储

- `GetTokenByHash` / `GetCredentialBinding` / `GetTokenByID` 依赖索引，复杂度约 `O(logN)`
- 工具执行热路径为：

1. Bearer token 哈希回查
2. 绑定凭证回查与解密
3. 云 API 调用

- 主要瓶颈：
- 高频认证回查
- `TouchToken` 造成的额外写入
- 缓解方式：
- 唯一索引 + 合理连接池
- `TouchToken` 保持兼容语义，必要时再做时间窗节流
- 不改工具层和授权层，缩小 blast radius

#### 路线 B：无状态 CAM + AssumeRole

- 不再有 token / credential binding 数据库回查
- 主要成本转移为：

1. 每次请求的 CAM 身份验证
2. 每次工具调用或每次会话初始化的 `AssumeRole`

- 理论复杂度近似 `O(1)` 本地处理，但真实瓶颈是远端 STS RTT 与限流
- 风险点：
- STS 频繁调用放大时延
- 通用 SSE 客户端无法稳定复用
- 若需要“多租户账号到角色”的统一映射，仍需一份中心化配置源，只是它不再是会话存储

### 避免技术债的关键决策

1. **不要把无状态 CAM 塞进现有 `shared-token` / `none` 模式**

- 应新增独立模式，例如 `cam-assume-role`
- 避免混淆静态凭证模式与真正的多租户无状态模式

2. **继续复用 `Principal` / `AuthorizePrincipal()`**

- 当前工具注册链路已依赖 `Principal`
- 无状态 CAM 也应在认证完成后构造同样的主体信息，而不是改动全部工具

3. **不要声称两条路线对客户端完全等价**

- 当前 `MCPClientOptionsFromEnv()` 与文档示例都只输出 Bearer 头
- 若走无状态 CAM，必须同步改客户端辅助工具或引入前置网关说明

4. **路线 B 先做模式隔离，再做默认切换**

- 先以 feature flag 新增模式
- 只有在验证客户端、STS 压力和运维可接受后，才考虑提升优先级

## Implementation Notes

- **共性**
- 不改 `tools/*` 的业务逻辑与 Guard 规则，避免无关重构
- 日志中禁止输出明文 token、SecretId、SecretKey、签名串、DSN、加密密钥
- 认证失败日志只记录可审计字段：主体、账号、角色 ARN、地域、请求路径、错误类别

- **路线 A**
- 继续保留 `/auth/bootstrap/tencentcloud`、`/auth/token-exchange/tencentcloud`、`/admin/tokens`
- 继续仅存 token 哈希与加密后的运行时凭证
- 共享存储改造应完全隐藏在 `TokenStore` / `CredentialStore` 后面

- **路线 B**
- 必须明确停用或降级 `/auth/bootstrap/tencentcloud`、`/auth/token-exchange/tencentcloud`、`/admin/tokens`
- `main.go` 输出的 MCP 客户端配置不能再默认是 Bearer token 模板
- 若客户端无法原生支持 CAM，可允许“受信网关注入身份”作为部署前提，但需在文档里明确不是通用客户端直连

## Architecture Design

### 已验证的当前链路

- `main.go` 在 `issued-token` 模式下装配 `TokenBoundCredentialProvider`
- `security/http_auth.go` 的 `issuedTokenAuthenticator` 通过 `GetTokenByHash()` 校验 Bearer token
- `security/token_exchange_tencentcloud.go` 在 bootstrap / token exchange 时签发本地 token，并写入 `CredentialBinding`
- `security/credential_dynamic.go` 在工具执行前通过 `GetCredentialBinding(token_id)` 恢复加密凭证
- `tools/registry.go` 调用 `AuthorizePrincipal()` 做 scope / region 授权

### 路线 A：issued-token + 共享存储

- 保持现有 HTTP 路由和客户端行为不变
- 只替换本地 store 实现为共享持久化
- 工具链路、授权模型、文档主路径都延续当前实现

### 路线 B：无状态 CAM + AssumeRole

- 新增 `cam-assume-role` 认证器：
- 从请求中提取并验证 CAM 身份
- 调用 `GetCallerIdentity`
- 构造 `Principal`
- 新增请求级动态凭证 Provider：
- 根据当前 `Principal` 和租户角色映射现场 `AssumeRole`
- 返回当次调用所需的临时凭证
- 取消会话态依赖：
- 不依赖 `TokenStore`
- 不依赖 `CredentialStore`
- 不依赖 token exchange 结果回查

## Directory Structure Summary

本次新 plan 以“决策优先、分支实施”为原则，目录按两条路线分别给出。

### 共用受影响文件

```text
/Users/lumenyang/workspace/tencentcloud-mcp-server/src/postgres/
├── main.go                       # [MODIFY] 启动装配、模式切换、路由注册、客户端配置输出
├── security/http_auth.go         # [MODIFY] 鉴权模式枚举与认证器工厂
├── security/principal.go         # [REUSE/MODIFY] 继续承载 subject、tenant、scope、region 授权语义
├── .env.example                  # [MODIFY] 更新模式说明、环境变量与约束
├── README.md                     # [MODIFY] 更新推荐方案、客户端接入说明与迁移路径
├── DEPLOY.md                     # [MODIFY] 更新部署建议、验证脚本与运维说明
└── SCF_DEPLOY.md                 # [MODIFY] 更新 SCF 下的推荐架构与限制
```

### 路线 A：issued-token + 共享存储

```text
/Users/lumenyang/workspace/tencentcloud-mcp-server/src/postgres/
├── go.mod                                # [MODIFY] 引入共享存储驱动
├── security/token_store_factory.go       # [NEW] 根据配置创建 sqlite 或共享后端 store
├── security/token_store_postgres.go      # [NEW] 共享 TokenStore + CredentialStore 实现
├── deploy/scf/scf_bootstrap              # [MODIFY] SCF 默认走外部共享存储
├── deploy/scf/scf.console.startup.sh     # [MODIFY] 启动脚本不再硬编码本地 sqlite
├── deploy/scf/scf.env.example            # [MODIFY] 外部存储示例
└── deploy/scf/scf.console.env.txt        # [MODIFY] 控制台环境变量示例
```

### 路线 B：无状态 CAM + AssumeRole

```text
/Users/lumenyang/workspace/tencentcloud-mcp-server/src/postgres/
├── security/auth_cam_assume_role.go          # [NEW] 无状态 CAM 认证器，完成身份确认与 Principal 构造
├── security/credential_cam_assume_role.go    # [NEW] 请求级动态 AssumeRole CredentialProvider
├── security/tenant_resolver.go               # [NEW] 账号到角色映射解析，首版可用静态配置
├── security/token_exchange_handlers.go       # [MODIFY] 在新模式下禁用或隐藏 bootstrap/token 接口
├── security/token_exchange_tencentcloud.go   # [MODIFY] 降为兼容路径，不再作为主链路
├── scripts/run_server.sh                     # [MODIFY] 新模式启动前置检查
├── cmd/verify/main.go                        # [MODIFY] 支持新模式下的验证方式
└── scripts/run_verify.sh                     # [MODIFY] 区分 Bearer 与 CAM 访问验证
```

## Key Code Structures

建议在路线 B 中新增以下核心抽象，以便与现有 `Principal` / `CredentialProvider` 模型对齐：

- `CAMAuthenticator`
- 负责从请求中验证调用者身份并生成 `Principal`
- `TenantResolver`
- 负责把 `account_id / user_id / principal_id` 映射为目标 `RoleArn`、允许地域和 scope
- `AssumeRoleCredentialProvider`
- 负责基于当前请求上下文现场获取临时凭证

## 决策建议

### 生产主线

- 选择 **issued-token + 共享存储**
- 原因：最兼容现有客户端和当前代码结构，风险最小，收益最直接

### 受控试点

- 选择 **无状态 CAM + AssumeRole**
- 场景：内部定制客户端、受信前置网关、强调不保会话态的专用部署
- 目标：验证客户端配合、STS 压力和调试复杂度，而不是立即替代主线

## Agent Extensions

### SubAgent

- **code-explorer**
- Purpose: 复核 `main.go`、`security`、文档与脚本的受影响范围，确保两条路线的文件清单和调用链准确
- Expected outcome: 形成完整的改动边界、路由清单和模式切换影响面

### Skill

- **code-analyst**
- Purpose: 评估两条路线在认证热路径、STS 调用频率、客户端兼容性和回归风险上的差异
- Expected outcome: 产出可执行的对比结论、验证矩阵和上线门禁建议