# PG MCP 验证记录模板

本文件用于记录**公开仓库可复现的验证方法与结果摘要**，不应提交以下内容：

- 真实实例 ID、地域、可用区、规格、容量
- 真实账号、密码、安全组、VPC / 子网、备份下载链接
- 本机绝对路径、个人用户名、私有测试环境拓扑

如果需要保存真实测试结果，建议记录在内部系统、私有文档或本地副本中。

---

## 一、建议填写的基础信息

| 项目 | 建议填写方式 |
|---|---|
| **报告时间** | `YYYY-MM-DD HH:mm TZ` |
| **Server 版本** | `tag / commit / branch` |
| **测试环境** | `本地 / Docker / CI` 等描述，不写真实实例信息 |
| **Guard 画像** | `default / dev / test / prod` |
| **执行命令** | `go build ./...` / `./scripts/run_openapi_param_check.sh` / `./scripts/run_mcp_smoke.sh` / `VERIFY_INSTANCE_ID=... ./scripts/run_verify.sh` |

---

## 二、推荐验证清单

### 1. 工程构建

```bash
go build ./...
```

### 2. OpenAPI 参数对齐

```bash
./scripts/run_openapi_param_check.sh
```

### 3. MCP 协议联调

```bash
./scripts/run_mcp_smoke.sh
```

### 4. 真实只读接口验证

```bash
VERIFY_INSTANCE_ID=postgres-xxxxxxxx ./scripts/run_verify.sh
```

> 执行真实云 API 验证前，请确认目标实例属于可控测试环境，且 `.env` 中密钥与 Guard 配置符合预期。

---

## 三、结果摘要模板

可按如下结构记录：

| 维度 | 结果 | 备注 |
|---|---|---|
| **工程构建** | ✅ / ❌ | 例如：`go build ./...` 是否通过 |
| **OpenAPI 参数对齐** | ✅ / ❌ | 记录 passed / failed / total |
| **MCP 协议联调** | ✅ / ❌ | 记录 `initialize / ping / tools/list / tools/call` 是否通过 |
| **真实只读调用** | ✅ / ❌ | 只记录覆盖数与是否通过，不写真实实例细节 |
| **写工具保护** | ✅ / ❌ | 记录 `confirm=false` 时是否正确拒绝 |

---

## 四、结果描述示例（脱敏版）

- `OpenAPI` 参数对齐通过，说明工具入参与腾讯云 SDK / OpenAPI 请求结构保持一致。
- `MCP smoke test` 通过，服务端可被通用 MCP 客户端识别，支持工具发现与只读调用。
- 真实只读验证通过，说明在受控测试环境中，基础查询接口可正常完成调用。
- 写工具在 `confirm=false` 时返回保护提示，未发生真实写操作。

---

## 五、提交前自查

- [ ] 已去掉真实实例 ID / 地域 / 规格 / 安全组等环境信息
- [ ] 已去掉本机绝对路径和个人标识
- [ ] 没有把临时观察日志、下载链接、密钥或凭证写入仓库
- [ ] 文档中的命令与当前代码、脚本行为一致
