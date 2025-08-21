# 腾讯云 AS MCP Server（中文版）

腾讯云 AS（Auto Scaling）MCP Server 实现，用于在 Claude / Cursor 等支持 MCP 的客户端中直接管理腾讯云弹性伸缩组和相关资源。

## 功能特性
- **伸缩组全生命周期管理**：创建、修改、启用、停用伸缩组
- **伸缩策略执行**：支持手动执行扩容和缩容操作
- **容量管理**：支持修改伸缩组的期望实例数

## 工具列表

### 🏗️ 伸缩组管理
| 工具名称 | 功能说明 |
|---|---|
| `CreateAutoScalingGroup` | 创建伸缩组 |
| `DescribeAutoScalingGroups` | 查询伸缩组 |
| `ModifyAutoScalingGroup` | 修改伸缩组 |
| `EnableAutoScalingGroup` | 启用伸缩组 |
| `DisableAutoScalingGroup` | 停用伸缩组 |

### ⚡ 伸缩操作
| 工具名称 | 功能说明 |
|---|---|
| `ExecuteScalingPolicy` | 执行伸缩策略 |
| `ModifyDesiredCapacity` | 修改期望实例数 |

## 快速开始

### 1. 准备腾讯云凭证
- 登录 [腾讯云控制台](https://console.cloud.tencent.com/)，进入「访问管理」→「访问密钥」获取 `SecretId` 与 `SecretKey`
- 可选：设置默认地域，如 `ap-guangzhou`

### 2. 配置环境变量
```bash
export TENCENTCLOUD_SECRET_ID=你的SecretId
export TENCENTCLOUD_SECRET_KEY=你的SecretKey
export TENCENTCLOUD_REGION=ap-guangzhou   # 可选
```

### 3. Claude Desktop 配置
编辑 `claude_desktop_config.json`（Mac 默认路径 `~/Library/Application Support/Claude/claude_desktop_config.json`），加入：

```json
{
  "mcpServers": {
    "tencent-as": {
      "command": "uv",
      "args": ["run", "mcp-server-as"],
      "env": {
        "TENCENTCLOUD_SECRET_ID": "你的SecretId",
        "TENCENTCLOUD_SECRET_KEY": "你的SecretKey",
        "TENCENTCLOUD_REGION": "ap-guangzhou"
      }
    }
  }
}
```

### 4. 安装
```bash
pip install mcp-server-as
```

## 许可证
MIT License，详见 LICENSE 文件。