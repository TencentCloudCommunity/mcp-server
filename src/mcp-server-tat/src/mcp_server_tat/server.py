from asyncio.log import logger
import json
import os
import base64

import time
from typing import Any
from mcp.server.models import InitializationOptions
import mcp.types as types
from mcp.server import NotificationOptions, Server
import mcp.server.stdio
from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.tat.v20201028 import tat_client, models


#从环境变量中读取 SecretId 和 SecretKey
secret_id   = os.getenv("TENCENTCLOUD_SECRET_ID")
secret_key  = os.getenv("TENCENTCLOUD_SECRET_KEY")
default_region = os.getenv("TENCENTCLOUD_REGION")

server = Server("tat")

@server.list_tools()
async def handle_list_tools() -> list[types.Tool]:
    return [
        types.Tool(
            name="SyncRunCommand",
            description=(
                "Synchronously execute Shell/PowerShell commands on Tencent Cloud instances via API. "
                "Supports Linux (default) and Windows OS. Recommended for short-running commands due to its efficient execution model. "
            ),
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Tencent Cloud region, such as ap-beijing (Beijing) or ap-guangzhou (Guangzhou). "
                        "If not specified, the value from the environment variable DEFAULT_REGION in the MCP environment will be used. If neither is set, an error will occur."
                    },
                    "InstanceId": {
                        "type": "string", 
                        "description": "Target cloud server instance ID  in format: ins-xxxxxxxx for CVM instances or lhins-xxxxxxxx for Lighthouse instances"
                    },
                    "Command": {
                        "type": "string", 
                        "description": "Command to execute."
                    },
                     "SystemType": {
                        "type": "string",
                        "description": "OS type determining command syntax (Linux=Shell, Windows=PowerShell)",
                        "enum": ["Linux", "Windows"],
                        "default": "Linux"
                    },
                },
                "required": ["Region", "InstanceId","Command"],
            },
        ),
        types.Tool(
            name="RunCommand",
            description=(
                "Asynchronously execute Shell/PowerShell commands on Tencent Cloud instances via API. "
                "Supports Linux (default) and Windows OS. Results must be queried via QueryTask interface. "
            ),
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Tencent Cloud region, such as ap-beijing (Beijing) or ap-guangzhou (Guangzhou). "
                        "If not specified, the value from the environment variable DEFAULT_REGION in the MCP environment will be used. If neither is set, an error will occur."
                    },
                    "InstanceId": {
                        "type": "string", 
                        "description": "Target cloud server instance ID  in format: ins-xxxxxxxx for CVM instances or lhins-xxxxxxxx for Lighthouse instances"
                    },
                    "Command": {
                        "type": "string", 
                        "description": "Command to execute."
                    },
                     "SystemType": {
                        "type": "string",
                        "description": "OS type determining command syntax (Linux=Shell, Windows=PowerShell)",
                        "enum": ["Linux", "Windows"],
                        "default": "Linux"
                    },
                },
                "required": ["Region", "InstanceId","Command"],
            },
        ),
        
        types.Tool(
            name="QueryTask",
            description=("Retrieve asynchronous command execution results from Tencent Cloud API. If task is Running ,retry later"),
            inputSchema={
                "type": "object",
                "properties": {
                    "Region": {
                        "type": "string",
                        "description": "Tencent Cloud region, such as ap-beijing (Beijing) or ap-guangzhou (Guangzhou). "
                        "If not specified, the value from the environment variable DEFAULT_REGION in the MCP environment will be used. If neither is set, an error will occur."
                    },
                    "TaskId": {
                        "type": "string",
                        "description": "Task Id returned by RunCommand, format: invt-xxxxxx"
                    }
                },
                "required": ["Region", "TaskId"],
            },
        ),
    ]


@server.call_tool()
async def handle_call_tool(
    name: str, arguments: dict[str, Any] | None
) -> list[types.TextContent | types.ImageContent | types.EmbeddedResource]:
    """Handle tool execution requests"""
    try:
        if name == "SyncRunCommand":
            region = arguments.get("Region",default_region)
            instance_id = arguments.get("InstanceId")
            command =  arguments.get("Command")
            system_type = arguments.get("SystemType","Linux")
            result = sync_run_command(region,instance_id,command,system_type);
            return [types.TextContent(type="text", text=str(result))]

        if name == "RunCommand":
            region = arguments.get("Region",default_region)
            instance_id = arguments.get("InstanceId")
            command =  arguments.get("Command")
            system_type = arguments.get("SystemType","Linux")
            result = run_command(region,instance_id,command,system_type);
            return [types.TextContent(type="text", text=str(result))]
        elif name == "QueryTask":
            region = arguments.get("Region",default_region)
            task_id = arguments.get("TaskId")
            result =query_task(region,task_id)
            return [types.TextContent(type="text", text=str(result))]
        else:
            raise ValueError(f"Unknown tool: {name}")
    except Exception as e:
        return [types.TextContent(type="text", text=f"Error: {str(e)}")] 
        



def get_tat_client(region:str)->tat_client.TatClient:
    """
    创建并返回 TAT 客户端
    """
    cred = credential.Credential(
        secret_id,
        secret_key
    )
    if not region:
        region = "ap-guangzhou"

    http_profile = HttpProfile()
    http_profile.endpoint = "tat.tencentcloudapi.com"

    client_profile = ClientProfile()
    client_profile.httpProfile = http_profile
    client_profile.request_client = "MCP-Server"

    return tat_client.TatClient(cred, region, client_profile)

def string_to_base64(input_string):
    bytes_data = input_string.encode('utf-8')
    base64_encoded = base64.b64encode(bytes_data)
    base64_string = base64_encoded.decode('utf-8')
    return base64_string

def base64_to_string(base64_string):
    base64_bytes = base64_string.encode('utf-8')
    bytes_data = base64.b64decode(base64_bytes)
    decoded_string = bytes_data.decode('utf-8')
    return decoded_string

def run_command(region:str,instance_id:str,command:str,system_type:str)->str:
    command_type = "POWERSHELL" if system_type.lower() == "windows" else "SHELL"
    client = get_tat_client(region)
    req = models.RunCommandRequest()
    params = {
        "InstanceIds": [instance_id],
        "Content": string_to_base64(command),
        "CommandType":command_type
    }
    req.from_json_string(json.dumps(params))
    resp = client.RunCommand(req)
    
    req = models.DescribeInvocationsRequest()
    params = {
        "InvocationIds": [resp.InvocationId],
    }
    req.from_json_string(json.dumps(params))
    resp = client.DescribeInvocations(req)
        # 提取结果
    result = {
        "TaskId":  resp.InvocationSet[0].InvocationTaskBasicInfoSet[0].InvocationTaskId
    }
    return json.dumps(result)
  


def query_task(region:str,task_id:str):
    client = get_tat_client(region)
    req = models.DescribeInvocationTasksRequest()
    params = {
        "InvocationTaskIds": [task_id],
        "HideOutput":False
    }
    req.from_json_string(json.dumps(params))
    resp = client.DescribeInvocationTasks(req)
    result = {
        "Status":  resp.InvocationTaskSet[0].TaskStatus,
        "ExitCode":resp.InvocationTaskSet[0].TaskResult.ExitCode,
        "Output":base64_to_string(resp.InvocationTaskSet[0].TaskResult.Output)
    }
    return json.dumps(result)



def sync_run_command(region: str, instance_id: str, command: str, system_type: str) -> str:
    command_type = "POWERSHELL" if system_type.lower() == "windows" else "SHELL"
    client = get_tat_client(region)
    
    # 执行命令
    run_command_req = models.RunCommandRequest()
    run_command_params = {
        "InstanceIds": [instance_id],
        "Content": string_to_base64(command),
        "CommandType": command_type
    }
    run_command_req.from_json_string(json.dumps(run_command_params))
    run_command_resp = client.RunCommand(run_command_req)
    
    # 获取调用信息以取得InvocationTaskId
    describe_invocations_req = models.DescribeInvocationsRequest()
    describe_invocations_params = {
        "InvocationIds": [run_command_resp.InvocationId],
    }
    describe_invocations_req.from_json_string(json.dumps(describe_invocations_params))
    describe_invocations_resp = client.DescribeInvocations(describe_invocations_req)
    invocation_task_id = describe_invocations_resp.InvocationSet[0].InvocationTaskBasicInfoSet[0].InvocationTaskId
    
    # 循环查询任务状态直到非running
    while True:
        time.sleep(1)  # 间隔1秒后再次查询
        describe_tasks_req = models.DescribeInvocationTasksRequest()
        describe_tasks_params = {
            "InvocationTaskIds": [invocation_task_id],
            "HideOutput": False
        }
        describe_tasks_req.from_json_string(json.dumps(describe_tasks_params))
        describe_tasks_resp = client.DescribeInvocationTasks(describe_tasks_req)
        current_status = describe_tasks_resp.InvocationTaskSet[0].TaskStatus
        
        if current_status.lower() != 'running':
            break
      
    # 构造最终结果
    task_result = describe_tasks_resp.InvocationTaskSet[0].TaskResult
    result = {
        "Status": current_status,
        "ExitCode": task_result.ExitCode,
        "Output": base64_to_string(task_result.Output)
    }
    
    return json.dumps(result)



async def serve():
    async with mcp.server.stdio.stdio_server() as (read_stream, write_stream):
        logger.info("Server running with stdio transport")
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="tat",
                server_version="0.1.0",
                capabilities=server.get_capabilities(
                    notification_options=NotificationOptions(),
                    experimental_capabilities={},
                ),
            ),
        )
