# Tencent Cloud CVM MCP Server
Implementation of Tencent Cloud CVM (Cloud Virtual Machine) MCP server for managing Tencent Cloud instances and network resources.

## Features
- **Instance Management**: Full lifecycle management including creating, starting, stopping, restarting, and terminating instances
- **Instance Query**: Query instance lists and instance type configurations  
- **Image Management**: Query available image lists
- **Network Management**: Query network resources like VPCs, subnets, and security groups
- **Region Management**: Query available regions and availability zones
- **Monitoring & Diagnostics**: CPU, memory, disk performance metrics monitoring
- **Security Group Management**: Create, configure and manage security group rules
- **Price Inquiry**: Pre-creation instance pricing functionality

## API List

### 🔍 Basic Query
| Tool Name | Description |
|---|---|
| `DescribeRegions` | Query region list |
| `DescribeZones` | Query availability zone list |
| `DescribeInstances` | Query instance list |
| `DescribeImages` | Query image list |
| `DescribeInstanceTypeConfigs` | Query instance type configurations |
| `DescribeVpcs` | Query VPC list |
| `DescribeSubnets` | Query subnet list |
| `DescribeSecurityGroups` | Query security group list |

### 🖥️ Instance Lifecycle
| Tool Name | Description |
|---|---|
| `RunInstances` | Create new instances |
| `QuickRunInstance` | Quick create instances (simplified) |
| `StartInstances` | Start instances |
| `StopInstances` | Stop instances |
| `RebootInstances` | Reboot instances |
| `TerminateInstances` | Terminate instances |
| `ResetInstancesPassword` | Reset instance password |
| `ResetInstance` | Reinstall instance OS |

### 🔐 Security Group Management
| Tool Name | Description |
|---|---|
| `DescribeSecurityGroupPolicies` | Query security group rules |
| `CreateSecurityGroup` | Create new security group |
| `CreateSecurityGroupWithPolicies` | Create security group with rules |
| `CreateSecurityGroupPolicies` | Add rules to existing security group |
| `ReplaceSecurityGroupPolicies` | Replace security group rules |

### 📊 Monitoring & Diagnostics
| Tool Name | Description |
|---|---|
| `CreateDiagnosticReports` | Create diagnostic reports |
| `DescribeDiagnosticReports` | Query diagnostic reports |
| `GetCpuUsageData` | Get CPU utilization |
| `GetCpuLoadavgData` | Get CPU 1-minute load average |
| `GetCpuloadavg5mData` | Get CPU 5-minute load average |
| `GetCpuloadavg15mData` | Get CPU 15-minute load average |
| `GetMemUsedData` | Get memory usage |
| `GetMemUsageData` | Get memory utilization |
| `GetCvmDiskUsageData` | Get disk utilization |
| `GetDiskTotalData` | Get disk total capacity |
| `GetDiskUsageData` | Get disk usage percentage |

### 💰 Pricing & Recommendations
| Tool Name | Description |
|---|---|
| `InquiryPriceRunInstances` | Inquiry price for creating instances |
| `DescribeRecommendZoneInstanceTypes` | Recommend instance types in zone |

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
    "tencent-cvm": {
      "command": "uv",
      "args": [
        "run",
        "mcp-server-cvm"
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
pip install mcp-server-cvm
```

## License
MIT License. See LICENSE file for details.