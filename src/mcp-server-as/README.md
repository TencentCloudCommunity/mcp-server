# Tencent Cloud AS MCP Server

Implementation of Tencent Cloud Auto Scaling (AS) MCP server for managing Tencent Cloud auto scaling groups and related resources.

## Features
- **Auto Scaling Group Management**: Full lifecycle management including creating, modifying, enabling, and disabling auto scaling groups
- **Scaling Policy Execution**: Support manual execution of scale-out and scale-in operations
- **Capacity Management**: Support modifying the desired capacity of auto scaling groups

## API List

### 🏗️ Auto Scaling Group Management
| Tool Name | Description |
|---|---|
| `CreateAutoScalingGroup` | Create an auto scaling group |
| `DescribeAutoScalingGroups` | Query auto scaling groups |
| `ModifyAutoScalingGroup` | Modify an auto scaling group |
| `EnableAutoScalingGroup` | Enable an auto scaling group |
| `DisableAutoScalingGroup` | Disable an auto scaling group |

### ⚡ Scaling Operations
| Tool Name | Description |
|---|---|
| `ExecuteScalingPolicy` | Execute a scaling policy |
| `ModifyDesiredCapacity` | Modify the desired capacity |

## Configuration

### Set Tencent Cloud Credentials
1. Obtain SecretId and SecretKey from Tencent Cloud Console
2. Set default region (optional)

### Environment Variables
Configure the following environment variables:
- `TENCENTCLOUD_SECRET_ID`: Tencent Cloud SecretId
- `TENCENTCLOUD_SECRET_KEY`: Tencent Cloud SecretKey
- `TENCENTCLOUD_REGION`: Default region (optional)

### Usage in Claude Desktop
Add the following configuration to claude_desktop_config.json:

```json
{
  "mcpServers": {
    "tencent-as": {
      "command": "uv",
      "args": [
        "run",
        "mcp-server-as"
      ],
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
pip install mcp-server-as
```

## License
MIT License. See LICENSE file for details.