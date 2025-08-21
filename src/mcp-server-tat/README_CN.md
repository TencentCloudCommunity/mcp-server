# 腾讯云 TAT MCP Server（中文版）

腾讯云 TAT（TencentCloud Automation Tools）MCP Server 实现，用于在 Claude / Cursor 等支持 MCP 的客户端中直接在腾讯云实例上执行命令。

## 功能特性
- **命令执行**：在腾讯云实例上执行 Shell/PowerShell 命令
- **同步与异步**：支持同步和异步两种命令执行模式
- **跨平台支持**：同时支持 Linux 和 Windows 实例
- **任务管理**：查询命令执行结果和输出

## 工具列表

### 🚀 命令执行
| 工具名称 | 功能说明 |
|---|---|
| `SyncRunCommand` | 同步执行实例命令 |
| `RunCommand` | 异步执行实例命令 |

### 📊 任务管理
| 工具名称 | 功能说明 |
|---|---|
| `QueryTask` | 查询命令执行结果 |

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
    "tencent-tat": {
      "command": "uv",
      "args": ["run", "mcp-server-tat"],
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
pip install mcp-server-tat
```

## 许可证
MIT License，详见 LICENSE 文件。