# 统一错误处理模板

## 适用范围

本模板适用于 PostgreSQL skill 在调用腾讯云 OpenAPI 前后的公共错误引导，重点覆盖：
- 凭证缺失
- 地域非法
- 巡检目标缺失
- SDK 缺失

## 总原则

- **先指出具体缺了什么**，不要只说“配置错误”
- **直接给下一步修复动作**，不要让用户自己猜
- **能给官方链接时就给官方链接**
- **涉及命令执行时先征求用户确认**，不要在未确认前自动安装依赖
- **如果已有 TC3 fallback 路径可继续执行，应先说明 fallback 是否可用**

## `missing-credentials` 模板

### 触发条件
- 缺少 `TENCENTCLOUD_SECRET_ID`
- 缺少 `TENCENTCLOUD_SECRET_KEY`
- 缺少 `TENCENTCLOUD_REGION`
- 或兼容变量 `MCP_REQUEST_*` / `MCP_*` 也都不可用

### 必须输出的内容
- 缺失的变量名
- 一段最小可复制配置示例
- 官方获取入口链接
- 明确说明修好后用户应把信息放在哪里

### 推荐回复模板

```text
当前缺少运行 PostgreSQL skill 所需的凭证/地域信息：<缺失变量列表>。

你可以先在运行环境里补齐下面这组最小配置：
export TENCENTCLOUD_SECRET_ID="你的 SecretId"
export TENCENTCLOUD_SECRET_KEY="你的 SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou"
# 如果你使用临时凭证，再补 TENCENTCLOUD_SESSION_TOKEN

获取和核对信息可参考腾讯云官方入口：
- API 密钥文档：https://cloud.tencent.com/document/product/598/40488
- API 密钥控制台：https://console.cloud.tencent.com/cam/capi
- 地域和可用区文档：https://cloud.tencent.com/document/api/238/7520

补齐后，把这些值放到当前会话或运行环境变量里，再继续执行即可。
```

## `invalid-region` 模板

### 触发条件
- 用户输入的地域不能被安全归一化
- 用户输入不在公共别名表中
- 地域值不是 PostgreSQL 当前支持的合法售卖地域

### 必须输出的内容
- 回显原始输入值
- 给出合法示例
- 提供官方查询链接
- 明确要求用户返回标准地域码或可确认的中文地域

### 推荐回复模板

```text
当前提供的地域值无法确认：<用户原始输入>。

请改成腾讯云标准地域码，或提供可以明确映射的中文地域，例如：
- 广州 -> ap-guangzhou
- 上海 -> ap-shanghai
- 成都 -> ap-chengdu
- 北京 -> ap-beijing

你可以先查看腾讯云官方地域信息：
- 地域和可用区：https://cloud.tencent.com/document/api/238/7520
- PostgreSQL 查询售卖地域（DescribeRegions）：https://cloud.tencent.com/document/product/409/16768

确认后，把地域改成标准值再继续，例如：
export TENCENTCLOUD_REGION="ap-guangzhou"
```

## `missing-target-scope` 模板

### 触发条件
- 用户请求巡检/诊断，但没有提供可以定位实例的必要信息：
  - 缺少地域
  - 缺少实例 ID（如 `postgres-xxxxxxxx`）或明确可辨识的实例名称

### 必须输出的内容
- 回显用户当前提供了什么
- 列出缺少的字段（地域 / 实例 ID）
- 给出腾讯云控制台入口链接，指导用户直接复制
- 给出一个完整的补全示例，方便用户一键回复

### 推荐回复模板

```text
我准备好执行 PG 巡检了，还需要你补充一下目标信息：

当前缺少：<地域 / 实例 ID / 两者都缺>

你可以在腾讯云 PostgreSQL 控制台直接查到并复制这些信息：
- 实例列表控制台：https://console.cloud.tencent.com/postgres

进入控制台后，你会看到每个实例的 实例 ID（格式：postgres-xxxxxxxx）和所属地域。

确认后，直接像下面这样发给我就行：
ap-guangzhou postgres-abc12345
```

## `missing-sdk` 模板

### 触发条件
- 当前环境未检测到腾讯云官方 SDK
- 且当前执行路径确实希望补齐 SDK，而不是继续使用 TC3 fallback

### 必须输出的内容
- 先说明 SDK 缺失是否会阻断当前执行
- 如果 fallback 可用，要先告诉用户“可以继续”
- 给出可复制安装命令
- 给出 SDK 官方文档链接
- 真正执行安装命令前必须询问用户是否同意

### 推荐回复模板

```text
当前环境未检测到腾讯云官方 SDK。

这个问题不一定会阻断当前 skill：如果当前链路允许，我可以先使用本地 TC3 签名 HTTPS 请求继续执行。
如果你更希望补齐官方 SDK，我也可以帮你执行安装。

可参考腾讯云官方 SDK 文档：
- Python SDK：https://cloud.tencent.com/document/sdk/Python
- 云 API / SDK 总入口：https://cloud.tencent.com/document/api

常用安装命令如下：
- Python: python3 -m pip install -U tencentcloud-sdk-python
- Node.js: npm install tencentcloud-sdk-nodejs

如果你愿意，我可以直接帮你执行对应安装命令。
```

## 使用要求

- 只要命中上述错误之一，就优先使用本模板，而不是临时自由发挥
- 若错误同时涉及“凭证缺失”和“地域非法”，应同时给出两类修复信息
- 若要执行安装命令、修改系统环境或落地凭证，必须先得到用户明确确认
