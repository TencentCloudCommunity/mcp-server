# 地域归一化公共规则

## 适用范围

本规则适用于 `tencent-pg-inspection`、`tencent-pg-slowquery-diagnosis`、`tencent-pg-ops-troubleshooter`，以及对应 bundle 根入口。

## 目标

- 在真正发起腾讯云 PostgreSQL OpenAPI 调用前，把用户输入的地域统一转换为标准地域码
- 尽量接受用户常见写法，但**不要对无法确认的输入做猜测性修正**
- 当地域无效或无法归一化时，统一转入 `@references/common/error_handling.md` 中的 `invalid-region` 模板

## 归一化顺序

1. **优先使用标准地域码**：如果输入已经是 `ap-guangzhou`、`ap-shanghai`、`ap-chengdu`、`ap-beijing` 这类标准值，直接使用。
2. **接受常见中文别名**：若用户输入为常见中文地域，则先映射到标准地域码。
3. **接受运行时默认地域**：如果用户未显式提供地域，但运行时已提供 `TENCENTCLOUD_REGION`，则可将其视为默认地域。
4. **无法安全确认时立即停止**：若输入无法稳定映射，或存在多个可能值，不要自行猜测；应直接提示用户修正，并附上官方查询链接。

## 当前支持的常见别名

| 用户输入 | 标准地域码 |
|---|---|
| `广州` | `ap-guangzhou` |
| `上海` | `ap-shanghai` |
| `成都` | `ap-chengdu` |
| `北京` | `ap-beijing` |

> 如果后续要扩展更多别名，应在这里统一维护，而不是分别写进每个 skill。

## 非法地域处理规则

出现以下任一情况，都视为“地域非法”或“地域无法确认”：
- 输入既不是标准地域码，也不在别名表中
- 输入为模糊描述，例如“华南”“国内”“离用户近一点”
- 输入看似像地域，但无法确认是否为 PostgreSQL 当前支持的售卖地域
- 输入包含明显拼写错误，例如 `guangzou`、`ap-gz`、`guangzhou-prod`

此时必须：
- 原样回显用户输入的原值
- 给出 1 到 3 个最接近的合法示例，但不要伪造“已确认映射”
- 附上官方链接，指导用户自行核对
- 参考 `@references/common/error_handling.md` 中的 `invalid-region` 模板

## 官方查询链接

- [腾讯云地域和可用区](https://cloud.tencent.com/document/api/238/7520)
- [腾讯云 PostgreSQL 查询售卖地域（DescribeRegions）](https://cloud.tencent.com/document/product/409/16768)

## 输出要求

在进入后续 OpenAPI 调用前，最终上下文里至少保留两项：
- 原始输入地域
- 归一化后的标准地域码

如果无法归一化，则明确标记为 `region unresolved`，并停止后续调用。
