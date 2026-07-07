# WorkBuddy 使用问题记录

## 1. 双账号切换后查询结果串号

### 现象

在 `WorkBuddy` 中做双账号测试时：

1. 先使用账号 A 通过连接器访问 PostgreSQL MCP 服务并查询实例。
2. 再切换连接器配置到账号 B。
3. 此时继续查询实例，返回结果仍可能表现为账号 A 的视角，出现“结果串号”的现象。

### 当前确认结果

该问题在 **重启 `WorkBuddy` 后恢复正常**。

也就是说，这个现象当前更像是 **客户端侧连接 / 会话 / 请求头复用导致的使用问题**，而不是服务端将两个账号的云资源查询结果真正混在一起。

### 当前排查结论

结合当前仓库主线实现，服务端在 `request-credential` 模式下：

- 每次请求都会从请求头读取：
  - `X-TencentCloud-Secret-Id`
  - `X-TencentCloud-Secret-Key`
  - `X-TencentCloud-Session-Token`（可选）
- 工具执行时会基于**当前请求**中的凭证重新创建腾讯云 SDK Client。
- 代码中未发现“按账号缓存腾讯云 Postgres Client 并跨请求复用”的逻辑。

因此，当前更可能的原因是：

- 修改连接器配置后，`WorkBuddy` 没有立即丢弃旧连接；
- 或者旧的 `streamable-http` 会话 / 旧请求头仍在继续生效；
- 导致切换账号后，短时间内查询仍落到旧账号凭证上。

### 使用规避建议

在 `WorkBuddy` 中进行多账号切换测试时，建议按下面方式操作：

1. **不要在同一个连接器上直接改账号配置反复复用**。
2. 为不同账号分别创建独立连接器，例如：
   - `pg-account-a`
   - `pg-account-b`
3. 切换账号后，如果发现查询结果异常，**优先重启 `WorkBuddy`**。
4. 在正式验证前，先做一次最小只读校验，例如：
   - `DescribeDBInstances`
   - `DescribeDBInstanceAttribute`
5. 如果后续继续遇到类似现象，建议优先怀疑：
   - 客户端连接未重建
   - 旧 header 未失效
   - 本地连接器缓存仍在生效

### 对当前服务端配置的建议

为了降低这类问题的影响，建议保持：

```env
MCP_AUTH_MODE=request-credential
MCP_STREAMABLE_HTTP_STATELESS=true
```

含义：

- `request-credential`：每次请求显式携带当前账号凭证；
- `MCP_STREAMABLE_HTTP_STATELESS=true`：尽量避免跨请求依赖进程内会话状态。

### 后续可选优化

如果后面需要进一步增强可观测性，可以补一个**只读调试工具**，在不暴露敏感明文的前提下，返回当前请求实际识别到的身份信息，例如：

- 当前账号 `AccountId`
- `Arn`
- `UserId`
- 凭证来源类型

这样在双账号切换时，可以更快确认“当前请求到底落到了哪个身份”。

---

记录时间：2026-07-07  
记录背景：`WorkBuddy` 双账号连接器切换测试  
当前状态：**已确认可通过重启 `WorkBuddy` 规避**
