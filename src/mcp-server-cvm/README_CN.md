# 腾讯云 CVM MCP Server（中文版）

腾讯云 CVM（Cloud Virtual Machine）MCP Server 实现，用于在 Claude / Cursor 等支持 MCP 的客户端中直接管理腾讯云实例与网络资源。

## 功能特性
- **实例全生命周期管理**：创建、启动、关机、重启、销毁、重置密码、诊断等
- **实例查询**：查看实例列表、实例规格、镜像列表
- **网络资源查询**：VPC、子网、安全组
- **地域与可用区**：查看支持的地域与可用区列表
- **监控与诊断**：CPU、内存、磁盘等性能指标监控
- **安全组管理**：创建、配置和管理安全组规则
- **价格查询**：创建实例前的询价功能

## 工具列表（Tools）

### 🔍 基础查询
| 工具名称 | 功能说明 |
|---|---|
| `DescribeRegions` | 查询地域列表 |
| `DescribeZones` | 查询可用区列表 |
| `DescribeInstances` | 查询实例列表 |
| `DescribeImages` | 查询镜像列表 |
| `DescribeInstanceTypeConfigs` | 查询实例规格 |
| `DescribeVpcs` | 查询 VPC 列表 |
| `DescribeSubnets` | 查询子网列表 |
| `DescribeSecurityGroups` | 查询安全组列表 |

### 🖥️ 实例生命周期
| 工具名称 | 功能说明 |
|---|---|
| `RunInstances` | 创建新实例 |
| `QuickRunInstance` | 快速创建实例（简化版） |
| `StartInstances` | 启动实例 |
| `StopInstances` | 关机实例 |
| `RebootInstances` | 重启实例 |
| `TerminateInstances` | 销毁实例 |
| `ResetInstancesPassword` | 重置实例密码 |
| `ResetInstance` | 重装实例操作系统 |
| `DescribeInstanceVncUrl` | 查询实例 VNC 管理终端登录地址（有效期 15 秒、一次性使用；实例须处于运行中状态） |

### 🔐 安全组管理
| 工具名称 | 功能说明 |
|---|---|
| `DescribeSecurityGroupPolicies` | 查询安全组规则 |
| `CreateSecurityGroup` | 创建安全组 |
| `CreateSecurityGroupWithPolicies` | 创建安全组并添加规则 |
| `CreateSecurityGroupPolicies` | 为现有安全组添加规则 |
| `ReplaceSecurityGroupPolicies` | 替换安全组规则 |

### 📊 监控与诊断
| 工具名称 | 功能说明 |
|---|---|
| `CreateDiagnosticReports` | 创建实例诊断报告 |
| `DescribeDiagnosticReports` | 查询诊断报告 |
| `GetCpuUsageData` | 获取CPU利用率 |
| `GetCpuLoadavgData` | 获取CPU一分钟平均负载 |
| `GetCpuloadavg5mData` | 获取CPU五分钟平均负载 |
| `GetCpuloadavg15mData` | 获取CPU十五分钟平均负载 |
| `GetMemUsedData` | 获取内存使用量 |
| `GetMemUsageData` | 获取内存利用率 |
| `GetCvmDiskUsageData` | 获取磁盘利用率 |
| `GetDiskTotalData` | 获取磁盘总容量 |
| `GetDiskUsageData` | 获取磁盘使用百分比 |

### 💰 价格与推荐
| 工具名称 | 功能说明 |
|---|---|
| `InquiryPriceRunInstances` | 创建实例询价 |
| `DescribeRecommendZoneInstanceTypes` | 推荐可用区实例类型 |

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
    "tencent-cvm": {
      "command": "uv",
      "args": ["run", "mcp-server-cvm"],
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
pip install mcp-server-cvm
```

## 许可证
MIT License，详见 LICENSE 文件。
