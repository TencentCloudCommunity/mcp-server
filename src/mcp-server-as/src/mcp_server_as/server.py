"""
腾讯云 AS 服务主模块
"""
from asyncio.log import logger
from typing import Any
from mcp.server.models import InitializationOptions
import mcp.types as types
from mcp.server import NotificationOptions, Server
import mcp.server.stdio
from . import tool_as

server = Server("as")

@server.list_tools()
async def handle_list_tools() -> list[types.Tool]:
    """处理工具列表请求"""
    return [
        types.Tool(
            name="CreateAutoScalingGroup",
            description="创建伸缩组",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "AutoScalingGroupName": {
                        "type": "string",
                        "description": "伸缩组名称",
                    },
                    "LaunchConfigurationId": {
                        "type": "string",
                        "description": "启动配置ID",
                    },
                    "MaxSize": {
                        "type": "integer",
                        "description": "最大实例数",
                    },
                    "MinSize": {
                        "type": "integer",
                        "description": "最小实例数",
                    },
                    "VpcId": {
                        "type": "string",
                        "description": "VPC ID",
                    },
                    "SubnetIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "子网ID列表",
                    },
                    "DesiredCapacity": {
                        "type": "integer",
                        "description": "期望实例数",
                    }
                },
                "required": ["Region", "AutoScalingGroupName", "LaunchConfigurationId", "MaxSize", "MinSize", "VpcId", "SubnetIds"],
            },
        ),
        types.Tool(
            name="DescribeAutoScalingGroups",
            description="查询伸缩组",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "AutoScalingGroupIds": {
                        "type": "array",
                        "items": {"type": "string"},
                        "description": "伸缩组ID列表",
                    },
                    "Filters": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "Name": {"type": "string"},
                                "Values": {
                                    "type": "array",
                                    "items": {"type": "string"}
                                }
                            }
                        },
                        "description": "过滤条件",
                    },
                    "Limit": {
                        "type": "integer",
                        "description": "返回数量",
                    },
                    "Offset": {
                        "type": "integer",
                        "description": "偏移量",
                    }
                },
                "required": ["Region"],
            },
        ),
        types.Tool(
            name="ModifyAutoScalingGroup",
            description="修改伸缩组",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "AutoScalingGroupId": {
                        "type": "string",
                        "description": "伸缩组ID",
                    },
                    "MaxSize": {
                        "type": "integer",
                        "description": "最大实例数",
                    },
                    "MinSize": {
                        "type": "integer",
                        "description": "最小实例数",
                    },
                    "DesiredCapacity": {
                        "type": "integer",
                        "description": "期望实例数",
                    }
                },
                "required": ["Region", "AutoScalingGroupId"],
            },
        ),
        types.Tool(
            name="EnableAutoScalingGroup",
            description="启用伸缩组",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "AutoScalingGroupId": {
                        "type": "string",
                        "description": "伸缩组ID",
                    }
                },
                "required": ["Region", "AutoScalingGroupId"],
            },
        ),
        types.Tool(
            name="DisableAutoScalingGroup",
            description="停用伸缩组",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "AutoScalingGroupId": {
                        "type": "string",
                        "description": "伸缩组ID",
                    }
                },
                "required": ["Region", "AutoScalingGroupId"],
            },
        ),
        types.Tool(
            name="ExecuteScalingPolicy",
            description="执行伸缩策略",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "AutoScalingGroupId": {
                        "type": "string",
                        "description": "伸缩组ID",
                    },
                    "Operation": {
                        "type": "string",
                        "description": "操作类型",
                        "enum": ["SCALE_OUT", "SCALE_IN"],
                    },
                    "AdjustmentType": {
                        "type": "string",
                        "description": "调整类型",
                        "enum": ["CHANGE_IN_CAPACITY", "EXACT_CAPACITY", "PERCENT_CHANGE_IN_CAPACITY"],
                    },
                    "AdjustmentValue": {
                        "type": "integer",
                        "description": "调整值",
                    }
                },
                "required": ["Region", "AutoScalingGroupId", "Operation"],
            },
        ),
        types.Tool(
            name="ModifyDesiredCapacity",
            description="修改期望实例数",
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "地域，如 ap-guangzhou",
                    },
                    "AutoScalingGroupId": {
                        "type": "string",
                        "description": "伸缩组ID",
                    },
                    "DesiredCapacity": {
                        "type": "integer",
                        "description": "期望实例数",
                    }
                },
                "required": ["Region", "AutoScalingGroupId", "DesiredCapacity"],
            },
        ),
    ]

@server.call_tool()
async def handle_call_tool(
    name: str, arguments: dict[str, Any] | None
) -> list[types.TextContent | types.ImageContent | types.EmbeddedResource]:
    """处理工具调用请求"""
    try:
        region = arguments.get("Region")
        
        if name == "CreateAutoScalingGroup":
            result = tool_as.create_auto_scaling_group(region, arguments)
        elif name == "DescribeAutoScalingGroups":
            result = tool_as.describe_auto_scaling_groups(
                region=region,
                auto_scaling_group_ids=arguments.get("AutoScalingGroupIds"),
                filters=arguments.get("Filters"),
                limit=arguments.get("Limit"),
                offset=arguments.get("Offset")
            )
        elif name == "ModifyAutoScalingGroup":
            result = tool_as.modify_auto_scaling_group(
                region=region,
                auto_scaling_group_id=arguments.get("AutoScalingGroupId"),
                params=arguments
            )
        elif name == "EnableAutoScalingGroup":
            result = tool_as.enable_auto_scaling_group(
                region=region,
                auto_scaling_group_id=arguments.get("AutoScalingGroupId")
            )
        elif name == "DisableAutoScalingGroup":
            result = tool_as.disable_auto_scaling_group(
                region=region,
                auto_scaling_group_id=arguments.get("AutoScalingGroupId")
            )
        elif name == "ExecuteScalingPolicy":
            result = tool_as.execute_scaling_policy(
                region=region,
                auto_scaling_group_id=arguments.get("AutoScalingGroupId"),
                operation=arguments.get("Operation"),
                adjustment_type=arguments.get("AdjustmentType"),
                adjustment_value=arguments.get("AdjustmentValue")
            )
        elif name == "ModifyDesiredCapacity":
            result = tool_as.modify_desired_capacity(
                region=region,
                auto_scaling_group_id=arguments.get("AutoScalingGroupId"),
                desired_capacity=arguments.get("DesiredCapacity")
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
                server_name="as",
                server_version="0.1.0",
                capabilities=server.get_capabilities(
                    notification_options=NotificationOptions(),
                    experimental_capabilities={},
                ),
            ),
        )
