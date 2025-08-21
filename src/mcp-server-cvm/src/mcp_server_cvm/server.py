"""
腾讯云 CVM 服务主模块
"""
from asyncio.log import logger
from typing import Any
from mcp.server.models import InitializationOptions
import mcp.types as types
from mcp.server import NotificationOptions, Server
import mcp.server.stdio
from . import tool_cvm, tool_vpc, tool_monitor

server = Server("cvm")

@server.list_tools()
async def handle_list_tools() -> list[types.Tool]:
    return [
        types.Tool(
            name="DescribeRegions",
            description="查询腾讯云CVM支持的地域列表",
            inputSchema={
                "type": "object",
                "properties": {},
                "required": [],
            },
        ),
        types.Tool(
            name="DescribeZones",
            description="查询腾讯云CVM支持的可用区列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="DescribeInstances",
            description="查询腾讯云CVM实例列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "Offset": {
                        "type": "integer",
                        "description": "偏移量，默认为0",
                    },
                    "Limit": {
                        "type": "integer",
                        "description": "返回数量，默认为20，最大值为100",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "按照一个或者多个实例ID查询，每次请求的实例的上限为100",
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="DescribeImages",
            description="查询腾讯云CVM镜像列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "ImageIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "镜像ID列表",
                    },
                    "ImageType": {
                        "type": "string",
                        "description": "镜像类型，可选值：PUBLIC_IMAGE, PRIVATE_IMAGE, SHARED_IMAGE, MARKET_IMAGE",
                    },
                    "Platform": {
                        "type": "string",
                        "description": "操作系统平台，如：TencentOS, Windows, CentOS等",
                    },
                    "ImageName": {
                        "type": "string",
                        "description": "镜像名称",
                    },
                    "Offset": {
                        "type": "integer",
                        "description": "偏移量，默认为0",
                        "default": 0
                    },
                    "Limit": {
                        "type": "integer",
                        "description": "返回数量，默认为20，最大值为100",
                        "default": 20,
                        "maximum": 100
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="DescribeInstanceTypeConfigs",
            description="查询实例机型配置列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "Zone": {
                        "type": "string",
                        "description": "可用区，如 ap-guangzhou-1",
                    },
                    "InstanceFamily": {
                        "type": "string",
                        "description": "实例机型系列，如 S5、SA2等",
                    }
                },
                "required": ["Region", "InstanceFamily"],
            },
        ),
        types.Tool(
            name="RebootInstances",
            description="重启实例",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "一个或多个待操作的实例ID",
                    },
                    "StopType": {
                        "type": "string",
                        "description": "关机类型。SOFT：表示软关机，HARD：表示硬关机，SOFT_FIRST：表示优先软关机，失败再硬关机",
                        "enum": ["SOFT", "HARD", "SOFT_FIRST"],
                        "default": "SOFT"
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="StartInstances",
            description="启动实例",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "一个或多个待操作的实例ID",
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="StopInstances",
            description="关闭实例",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "一个或多个待操作的实例ID",
                    },
                    "StopType": {
                        "type": "string",
                        "description": "关机类型。SOFT：表示软关机，HARD：表示硬关机，SOFT_FIRST：表示优先软关机，失败再硬关机",
                        "enum": ["SOFT", "HARD", "SOFT_FIRST"],
                        "default": "SOFT"
                    },
                    "StoppedMode": {
                        "type": "string",
                        "description": "关机模式，仅对POSTPAID_BY_HOUR类型实例生效",
                        "enum": ["KEEP_CHARGING", "STOP_CHARGING"],
                        "default": "KEEP_CHARGING"
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="TerminateInstances",
            description="退还实例（敏感操作，请确认后执行）",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "一个或多个待操作的实例ID",
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="ResetInstancesPassword",
            description="重置实例密码",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "一个或多个待操作的实例ID",
                    },
                    "Password": {
                        "type": "string",
                        "description": "实例新密码",
                    },
                    "ForceStop": {
                        "type": "boolean",
                        "description": "是否强制关机执行",
                        "default": False
                    }
                },
                "required": ["Region", "InstanceIds", "Password"],
            },
        ),
        types.Tool(
            name="RunInstances",
            description="创建实例（请注意该操作将会产生费用）",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "Zone": {
                        "type": "string",
                        "description": "可用区，如 ap-guangzhou-1",
                    },
                    "InstanceType": {
                        "type": "string",
                        "description": "实例机型，如 S5.MEDIUM4",
                    },
                    "ImageId": {
                        "type": "string",
                        "description": "镜像ID",
                    },
                    "VpcId": {
                        "type": "string",
                        "description": "私有网络ID",
                    },
                    "SubnetId": {
                        "type": "string",
                        "description": "子网ID",
                    },
                    "InstanceName": {
                        "type": "string",
                        "description": "实例名称",
                    },
                    "SecurityGroupIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "安全组ID列表",
                    },
                    "InstanceChargeType": {
                        "type": "string",
                        "description": "实例计费类型",
                        "enum": ["PREPAID", "POSTPAID_BY_HOUR"],
                        "default": "POSTPAID_BY_HOUR"
                    },
                    "Password": {
                        "type": "string",
                        "description": "实例密码",
                    },
                    "InstanceCount": {
                        "type": "integer",
                        "description": "实例数量",
                    }
                },
                "required": ["Region", "Zone", "InstanceType", "ImageId", "VpcId", "SubnetId"],
            },
        ),
        types.Tool(
            name="ResetInstance",
            description="重装指定实例的操作系统（敏感操作，请确认后执行）",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceId": {
                        "type": "string",
                        "description": "待重装系统的实例ID",
                    },
                    "ImageId": {
                        "type": "string",
                        "description": "重装使用的镜像ID",
                    },
                    "Password": {
                        "type": "string",
                        "description": "实例重装后的密码，如果不指定则保持原密码不变",
                    }
                },
                "required": ["Region", "InstanceId", "ImageId"],
            },
        ),
        types.Tool(
            name="DescribeSecurityGroups",
            description="查询安全组列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "SecurityGroupIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "安全组ID列表",
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="DescribeVpcs",
            description="查询VPC列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "VpcIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "VPC实例ID数组",
                    },
                    "IsDefault": {
                        "type": "boolean",
                        "description": "是否是默认VPC，True表示默认VPC，False表示非默认VPC，不传表示不过滤",
                    },
                    "VpcName": {
                        "type": "string",
                        "description": "VPC名称，用于过滤指定名称的VPC",
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="DescribeSubnets",
            description="查询子网列表",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "VpcId": {
                        "type": "string",
                        "description": "VPC实例ID，用于过滤指定VPC下的子网",
                    },
                    "SubnetIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "子网实例ID列表，用于查询指定子网的信息",
                    },
                    "Zone": {
                        "type": "string",
                        "description": "可用区，如 ap-guangzhou-1，用于过滤指定可用区的子网",
                    },
                    "IsDefault": {
                        "type": "boolean",
                        "description": "是否是默认子网，True表示默认子网，False表示非默认子网，不传表示不过滤",
                    },
                    "VpcName": {
                        "type": "string",
                        "description": "VPC名称，用于过滤指定VPC名称下的子网",
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="GetCpuUsageData",
            description="获取CPU利用率指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetCpuLoadavgData",
            description="获取CPU一分钟平均负载指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetCpuloadavg5mData",
            description="获取CPU五分钟平均负载指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetCpuloadavg15mData",
            description="获取CPU十五分钟平均负载指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetMemUsedData",
            description="获取内存使用量指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetMemUsageData",
            description="获取内存利用率指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetCvmDiskUsageData",
            description="获取磁盘利用率指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetDiskTotalData",
            description="获取磁盘分区总容量指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="GetDiskUsageData",
            description="获取磁盘分区已使用容量和总容量的百分比指标数据",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    },
                    "Period": {
                        "type": "integer",
                        "description": "统计周期，单位秒，默认60",
                        "default": 60
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="InquiryPriceRunInstances",
            description="创建实例询价",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "Zone": {
                        "type": "string",
                        "description": "可用区，如 ap-guangzhou-1",
                    },
                    "InstanceType": {
                        "type": "string",
                        "description": "实例机型，如 S5.MEDIUM4",
                    },
                    "ImageId": {
                        "type": "string",
                        "description": "镜像ID",
                    },
                    "SystemDisk": {
                        "type": "object",
                        "description": "系统盘配置",
                        "properties": {
                            "DiskType": {
                                "type": "string",
                                "description": "系统盘类型，默认为 CLOUD_PREMIUM",
                                "enum": ["CLOUD_BASIC", "CLOUD_PREMIUM", "CLOUD_SSD"],
                                "default": "CLOUD_PREMIUM"
                            },
                            "DiskSize": {
                                "type": "integer",
                                "description": "系统盘大小，单位为GB，默认为50",
                                "default": 50
                            }
                        }
                    },
                    "InstanceChargeType": {
                        "type": "string",
                        "description": "实例计费类型",
                        "enum": ["PREPAID", "POSTPAID_BY_HOUR"],
                        "default": "POSTPAID_BY_HOUR"
                    },
                    "InstanceChargePrepaid": {
                        "type": "object",
                        "description": "预付费模式，即包年包月相关参数设置",
                        "properties": {
                            "Period": {
                                "type": "integer",
                                "description": "购买实例的时长，单位：月",
                                "default": 1
                            },
                            "RenewFlag": {
                                "type": "string",
                                "description": "自动续费标识",
                                "enum": ["NOTIFY_AND_AUTO_RENEW", "NOTIFY_AND_MANUAL_RENEW", "DISABLE_NOTIFY_AND_MANUAL_RENEW"],
                                "default": "NOTIFY_AND_MANUAL_RENEW"
                            }
                        }
                    }
                },
                "required": ["Region", "Zone", "InstanceType", "ImageId"],
            },
        ),
        types.Tool(
            name="CreateDiagnosticReports",
            description="创建实例诊断报告",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "InstanceIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例ID列表",
                    }
                },
                "required": ["Region", "InstanceIds"],
            },
        ),
        types.Tool(
            name="DescribeDiagnosticReports",
            description="查询实例诊断报告",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "ReportIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "实例健康检测报告ID列表，如：[\"dr-rfmme2si\"]。每次请求批量报告ID的上限为100。注意：不能与 Filters 或 VagueInstanceName 联合使用",
                    },
                    "Filters": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "Name": {
                                    "type": "string",
                                    "description": "过滤条件名称",
                                    "enum": ["instance-id", "instance-name", "instance-health-status", "report-status", "cluster-ids"]
                                },
                                "Values": {
                                    "type": "array",
                                    "items": {"type": "string"},
                                    "description": "过滤条件值"
                                }
                            }
                        },
                        "description": "过滤条件列表。支持的过滤条件：instance-id(实例ID), instance-name(实例名称), instance-health-status(实例健康状态：Normal/Warn/Critical), report-status(报告状态：Querying/Finished), cluster-ids(集群ID列表)",
                    },
                    "VagueInstanceName": {
                        "type": "string",
                        "description": "模糊实例别名。注意：不能与 ReportIds 或 ClusterDiagnosticReportIds 联合使用",
                    },
                    "Offset": {
                        "type": "integer",
                        "description": "偏移量，默认为0",
                    },
                    "Limit": {
                        "type": "integer",
                        "description": "返回数量，默认为20，最大值为100",
                    },
                    "ClusterDiagnosticReportIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "集群健康检测报告ID列表，如：[\"cr-rfmme2si\"]。注意：不能与 Filters 或 VagueInstanceName 联合使用",
                    },
                    "ScenarioId": {
                        "type": "integer",
                        "description": "检测场景ID，默认为1表示对CVM进行全面体检，200为集群一致性检测场景",
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="DescribeSecurityGroupPolicies",
            description="Query security group rules",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Region, e.g., ap-guangzhou",
                    },
                    "SecurityGroupId": {
                        "type": "string",
                        "description": "Security group ID",
                    }
                },
                "required": ["Region", "SecurityGroupId"],
            },
        ),
        types.Tool(
            name="CreateSecurityGroupPolicies",
            description="Create security group rules",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Region, e.g., ap-guangzhou",
                    },
                    "SecurityGroupId": {
                        "type": "string",
                        "description": "Security group ID",
                    },
                    "SecurityGroupPolicySet": {
                        "type": "object",
                        "description": "Security group policy set containing ingress and egress rules",
                        "properties": {
                            "Ingress": {
                                "type": "array",
                                "items": {"type": "object"},
                                "description": "Inbound rules"
                            },
                            "Egress": {
                                "type": "array",
                                "items": {"type": "object"},
                                "description": "Outbound rules"
                            }
                        }
                    }
                },
                "required": ["Region", "SecurityGroupId", "SecurityGroupPolicySet"],
            },
        ),
        types.Tool(
            name="CreateSecurityGroup",
            description="Create a new security group",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Region, e.g., ap-guangzhou",
                    },
                    "GroupName": {
                        "type": "string",
                        "description": "Security group name",
                    },
                    "GroupDescription": {
                        "type": "string",
                        "description": "Security group description",
                    },
                    "ProjectId": {
                        "type": "integer",
                        "description": "Project ID, default is 0",
                    },
                    "Tags": {
                        "type": "array",
                        "items": {"type": "object"},
                        "description": "Tag list",
                    }
                },
                "required": ["Region", "GroupName"],
            },
        ),
        types.Tool(
            name="CreateSecurityGroupWithPolicies",
            description="Create a new security group with rules",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Region, e.g., ap-guangzhou",
                    },
                    "GroupName": {
                        "type": "string",
                        "description": "Security group name",
                    },
                    "SecurityGroupPolicySet": {
                        "type": "object",
                        "description": "Security group policy set containing ingress and egress rules",
                        "properties": {
                            "Ingress": {
                                "type": "array",
                                "items": {"type": "object"},
                                "description": "Inbound rules"
                            },
                            "Egress": {
                                "type": "array",
                                "items": {"type": "object"},
                                "description": "Outbound rules"
                            }
                        }
                    },
                    "GroupDescription": {
                        "type": "string",
                        "description": "Security group description",
                    },
                    "ProjectId": {
                        "type": "integer",
                        "description": "Project ID, default is 0",
                    },
                    "Tags": {
                        "type": "array",
                        "items": {"type": "object"},
                        "description": "Tag list",
                    }
                },
                "required": ["Region", "GroupName", "SecurityGroupPolicySet"],
            },
        ),
        types.Tool(
            name="ReplaceSecurityGroupPolicies",
            description="Replace security group rules in batch (can only replace rules in one direction per request)",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Region, e.g., ap-guangzhou",
                    },
                    "SecurityGroupId": {
                        "type": "string",
                        "description": "Security group ID",
                    },
                    "SecurityGroupPolicySet": {
                        "type": "object",
                        "description": "Security group policy set containing either ingress or egress rules",
                        "properties": {
                            "Ingress": {
                                "type": "array",
                                "items": {"type": "object"},
                                "description": "Inbound rules with PolicyIndex"
                            },
                            "Egress": {
                                "type": "array",
                                "items": {"type": "object"},
                                "description": "Outbound rules with PolicyIndex"
                            }
                        }
                    }
                },
                "required": ["Region", "SecurityGroupId", "SecurityGroupPolicySet"],
            },
        ),
        types.Tool(
            name="QuickRunInstance",
            description="快速创建实例（自动处理各种参数，简化创建流程）",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou，默认为广州",
                    },
                    "Zone": {
                        "type": "string",
                        "description": "可用区，如 ap-guangzhou-1，不指定则自动选择",
                    },
                    "InstanceType": {
                        "type": "string",
                        "description": "实例类型，可以是具体的类型如'SA5.MEDIUM2'，也可以是描述如'2核4G'，不指定则使用默认配置",
                    },
                    "ImageId": {
                        "type": "string",
                        "description": "镜像ID，不指定则使用Ubuntu公共镜像",
                    },
                    "VpcId": {
                        "type": "string",
                        "description": "VPC ID，不指定则自动选择",
                    },
                    "SubnetId": {
                        "type": "string",
                        "description": "子网ID，不指定则自动选择",
                    },
                    "Password": {
                        "type": "string",
                        "description": "实例密码，不指定则随机生成",
                    },
                    "InstanceName": {
                        "type": "string",
                        "description": "实例名称",
                    },
                    "InstanceChargeType": {
                        "type": "string",
                        "description": "计费类型",
                        "enum": ["PREPAID", "POSTPAID_BY_HOUR"],
                        "default": "POSTPAID_BY_HOUR"
                    }
                },
                "required": [],
            },
        ),
        types.Tool(
            name="DescribeRecommendZoneInstanceTypes",
            description="推荐用户在相同或相近地域购买相同或相似机型",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "Zone": {
                        "type": "string",
                        "description": "当前机型所在可用区，如 ap-guangzhou-2",
                    },
                    "InstanceType": {
                        "type": "string",
                        "description": "实例规格，如 SA2.LARGE8",
                    },
                    "InstanceChargeType": {
                        "type": "string",
                        "description": "实例计费类型，如 PREPAID",
                    },
                    "IsRecommendUnderStock": {
                        "type": "boolean",
                        "description": "是否推荐低库存机型数据，默认值为false",
                        "default": False
                    },
                    "IsCompare": {
                        "type": "boolean",
                        "description": "是否属于机型比对场景，默认值是false；如果值为true，则会限定只返回当前地域的主力机型数据",
                        "default": False
                    }
                },
                "required": ["Region", "Zone", "InstanceType", "InstanceChargeType"],
            },
        ),
    ]

@server.call_tool()
async def handle_call_tool(
    name: str, arguments: dict[str, Any] | None
) -> list[types.TextContent | types.ImageContent | types.EmbeddedResource]:
    """处理工具调用请求"""
    try:
        region = arguments.get("Region", None)
        instance_ids = arguments.get("InstanceIds", [])
        period = arguments.get("Period", 60)
        
        # CVM 相关操作
        if name == "DescribeRegions":
            result = tool_cvm.describe_regions()
        elif name == "DescribeZones":
            result = tool_cvm.describe_zones(region)
        elif name == "DescribeInstances":
            result = tool_cvm.describe_instances(
                region,
                arguments.get("Offset", 0),
                arguments.get("Limit", 20),
                instance_ids
            )
        elif name == "DescribeImages":
            result = tool_cvm.describe_images(
                region=region,
                image_ids=arguments.get("ImageIds"),
                image_type=arguments.get("ImageType"),
                platform=arguments.get("Platform"),
                image_name=arguments.get("ImageName"),
                offset=arguments.get("Offset"),
                limit=arguments.get("Limit")
            )
        elif name == "DescribeInstanceTypeConfigs":
            result = tool_cvm.describe_instance_type_configs(
                region,
                arguments.get("Zone"),
                arguments.get("InstanceFamily")
            )
        elif name == "RebootInstances":
            result = tool_cvm.reboot_instances(
                region,
                instance_ids,
                arguments.get("StopType", "SOFT")
            )
        elif name == "StartInstances":
            result = tool_cvm.start_instances(region, instance_ids)
        elif name == "StopInstances":
            result = tool_cvm.stop_instances(
                region,
                instance_ids,
                arguments.get("StopType", "SOFT"),
                arguments.get("StoppedMode", "KEEP_CHARGING")
            )
        elif name == "TerminateInstances":
            result = tool_cvm.terminate_instances(region, instance_ids)
        elif name == "ResetInstancesPassword":
            result = tool_cvm.reset_instances_password(
                region,
                instance_ids,
                arguments.get("Password"),
                arguments.get("ForceStop", False)
            )
        elif name == "RunInstances":
            result = tool_cvm.run_instances(region, arguments)
        elif name == "ResetInstance":
            result = tool_cvm.reset_instance(
                region=region,
                instance_id=arguments.get("InstanceId"),
                image_id=arguments.get("ImageId"),
                password=arguments.get("Password")
            )
        elif name == "InquiryPriceRunInstances":
            result = tool_cvm.inquiry_price_run_instances(
                region=region,
                params=arguments
            )
            
        # VPC 相关操作
        elif name == "DescribeSecurityGroups":
            result = tool_vpc.describe_security_groups(
                region,
                arguments.get("SecurityGroupIds")
            )
        elif name == "DescribeVpcs":
            result = tool_vpc.describe_vpcs(
                region=region,
                vpc_ids=arguments.get("VpcIds"),
                is_default=arguments.get("IsDefault"),
                vpc_name=arguments.get("VpcName")
            )
        elif name == "DescribeSubnets":
            result = tool_vpc.describe_subnets(
                region=region,
                vpc_id=arguments.get("VpcId"),
                subnet_ids=arguments.get("SubnetIds"),
                zone=arguments.get("Zone"),
                is_default=arguments.get("IsDefault"),
                vpc_name=arguments.get("VpcName")
            )
        
        # 监控相关操作
        elif name == "GetCpuUsageData":
            result = tool_monitor.get_cpu_usage_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetCpuLoadavgData":
            result = tool_monitor.get_cpu_loadavg_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetCpuloadavg5mData":
            result = tool_monitor.get_cpu_loadavg5m_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetCpuloadavg15mData":
            result = tool_monitor.get_cpu_loadavg15m_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetMemUsedData":
            result = tool_monitor.get_mem_used_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetMemUsageData":
            result = tool_monitor.get_mem_usage_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetCvmDiskUsageData":
            result = tool_monitor.get_cvm_disk_usage_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetDiskTotalData":
            result = tool_monitor.get_disk_total_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "GetDiskUsageData":
            result = tool_monitor.get_disk_usage_data(region, instance_ids, period)
            return [types.TextContent(type="text", text=str(result))]
        elif name == "CreateDiagnosticReports":
            result = tool_cvm.create_diagnostic_reports(
                region=region,
                instance_ids=arguments.get("InstanceIds")
            )
        elif name == "DescribeDiagnosticReports":
            result = tool_cvm.describe_diagnostic_reports(
                region=region,
                report_ids=arguments.get("ReportIds"),
                filters=arguments.get("Filters"),
                vague_instance_name=arguments.get("VagueInstanceName"),
                offset=arguments.get("Offset"),
                limit=arguments.get("Limit"),
                cluster_diagnostic_report_ids=arguments.get("ClusterDiagnosticReportIds"),
                scenario_id=arguments.get("ScenarioId")
            )
        elif name == "DescribeSecurityGroupPolicies":
            result = tool_vpc.describe_security_group_policies(
                region=region,
                security_group_id=arguments.get("SecurityGroupId")
            )
        elif name == "CreateSecurityGroupPolicies":
            result = tool_vpc.create_security_group_policies(
                region=region,
                security_group_id=arguments.get("SecurityGroupId"),
                security_group_policy_set=arguments.get("SecurityGroupPolicySet")
            )
        elif name == "CreateSecurityGroup":
            result = tool_vpc.create_security_group(
                region=region,
                group_name=arguments.get("GroupName"),
                group_description=arguments.get("GroupDescription"),
                project_id=arguments.get("ProjectId"),
                tags=arguments.get("Tags")
            )
        elif name == "CreateSecurityGroupWithPolicies":
            result = tool_vpc.create_security_group_with_policies(
                region=region,
                group_name=arguments.get("GroupName"),
                security_group_policy_set=arguments.get("SecurityGroupPolicySet"),
                group_description=arguments.get("GroupDescription"),
                project_id=arguments.get("ProjectId"),
                tags=arguments.get("Tags")
            )
        elif name == "ReplaceSecurityGroupPolicies":
            result = tool_vpc.replace_security_group_policies(
                region=region,
                security_group_id=arguments.get("SecurityGroupId"),
                security_group_policy_set=arguments.get("SecurityGroupPolicySet")
            )
        elif name == "QuickRunInstance":
            result = tool_cvm.quick_run_instance(
                region=arguments.get("Region", "ap-guangzhou"),
                zone=arguments.get("Zone"),
                instance_type=arguments.get("InstanceType"),
                image_id=arguments.get("ImageId"),
                vpc_id=arguments.get("VpcId"),
                subnet_id=arguments.get("SubnetId"),
                password=arguments.get("Password"),
                instance_name=arguments.get("InstanceName"),
                instance_charge_type=arguments.get("InstanceChargeType", "POSTPAID_BY_HOUR")
            )
        elif name == "DescribeRecommendZoneInstanceTypes":
            result = tool_cvm.describe_recommend_zone_instance_types(
                region=arguments["Region"],
                zone=arguments["Zone"],
                instance_type=arguments["InstanceType"],
                instance_charge_type=arguments["InstanceChargeType"],
                is_recommend_under_stock=arguments.get("IsRecommendUnderStock", False),
                is_compare=arguments.get("IsCompare", False)
            )
        else:
            raise ValueError(f"未知的工具: {name}")
            
        return [types.TextContent(type="text", text=str(result))]
    except Exception as e:
        return [types.TextContent(type="text", text=f"错误: {str(e)}")]

async def serve():
    """启动服务"""
    async with mcp.server.stdio.stdio_server() as (read_stream, write_stream):
        logger.info("Server running with stdio transport")
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="cvm",
                server_version="0.1.0",
                capabilities=server.get_capabilities(
                    notification_options=NotificationOptions(),
                    experimental_capabilities={},
                ),
            ),
        ) 