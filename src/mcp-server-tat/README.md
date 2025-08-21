# Tencent Cloud TAT MCP Server

An MCP server implementation for executing commands on Tencent Cloud instances using the TencentCloud Automation Tools (TAT) API.

## Features

- **Command Execution**: Execute Shell/PowerShell commands on Tencent Cloud instances
- **Sync & Async Support**: Both synchronous and asynchronous command execution
- **Cross-Platform**: Supports both Linux and Windows instances
- **Task Management**: Query command execution results and outputs

## API List

### 🚀 Command Execution
| Tool Name | Description |
|---|---|
| `SyncRunCommand` | Synchronously execute commands on instances |
| `RunCommand` | Asynchronously execute commands on instances |

### 📊 Task Management
| Tool Name | Description |
|---|---|
| `QueryTask` | Retrieve command execution results |

## Configuration

### Setting up Tencent Cloud Credentials

1. Obtain SecretId and SecretKey from Tencent Cloud Console
2. Set your default region (optional)

### Environment Variables
Configure the following environment variables:
- `TENCENTCLOUD_SECRET_ID`: Tencent Cloud SecretId
- `TENCENTCLOUD_SECRET_KEY`: Tencent Cloud SecretKey
- `TENCENTCLOUD_REGION`: Default region (optional)

### Usage with Claude Desktop

Add this to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "tencent-tat": {
      "command": "uv",
      "args": ["run", "mcp-server-tat"],
      "env": {
        "TENCENTCLOUD_SECRET_ID": "YOUR_SECRET_ID_HERE",
        "TENCENTCLOUD_SECRET_KEY": "YOUR_SECRET_KEY_HERE",
        "TENCENTCLOUD_REGION": "YOUR_REGION_HERE"
      }
    }
  }
}
```

## Installation

```sh
pip install mcp-server-tat
```

## License

MIT License. See LICENSE for details.