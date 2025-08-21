"""
腾讯云客户端创建模块
"""
import os
from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.common.common_client import CommonClient
from tencentcloud.cvm.v20170312 import cvm_client
from tencentcloud.vpc.v20170312 import vpc_client
from tencentcloud.monitor.v20180724 import monitor_client

# 从环境变量中读取认证信息
secret_id = os.getenv("TENCENTCLOUD_SECRET_ID")
secret_key = os.getenv("TENCENTCLOUD_SECRET_KEY")
default_region = os.getenv("TENCENTCLOUD_REGION")

def get_common_client(region: str, product = "cvm", version="2019-12-12") -> CommonClient:
    """
    创建并返回通用客户端实例

    Args:
        region: 地域信息
        product: 产品名称
        version: 产品版本

    Returns:
        CommonClient: 通用客户端实例
    """
    cred = credential.Credential(secret_id, secret_key)
    if not region:
        region = default_region or "ap-guangzhou"

    http_profile = HttpProfile()
    http_profile.endpoint = "cvm.tencentcloudapi.com"

    client_profile = ClientProfile()
    client_profile.httpProfile = http_profile
    client_profile.request_client = "MCP-Server"

    return CommonClient(product, version, cred, region, profile=client_profile)

def get_cvm_client(region: str) -> cvm_client.CvmClient:
    """
    创建并返回CVM客户端

    Args:
        region: 地域信息

    Returns:
        CvmClient: CVM客户端实例
    """
    cred = credential.Credential(secret_id, secret_key)
    if not region:
        region = default_region or "ap-guangzhou"

    http_profile = HttpProfile()
    http_profile.endpoint = "cvm.tencentcloudapi.com"

    client_profile = ClientProfile()
    client_profile.httpProfile = http_profile
    client_profile.request_client = "MCP-Server"

    return cvm_client.CvmClient(cred, region, client_profile)

def get_vpc_client(region: str) -> vpc_client.VpcClient:
    """
    创建并返回VPC客户端

    Args:
        region: 地域信息

    Returns:
        VpcClient: VPC客户端实例
    """
    cred = credential.Credential(secret_id, secret_key)
    if not region:
        region = default_region or "ap-guangzhou"

    http_profile = HttpProfile()
    http_profile.endpoint = "vpc.tencentcloudapi.com"

    client_profile = ClientProfile()
    client_profile.httpProfile = http_profile
    client_profile.request_client = "MCP-Server"

    return vpc_client.VpcClient(cred, region, client_profile)

def get_monitor_client(region: str) -> monitor_client.MonitorClient:
    """
    创建并返回监控客户端

    Args:
        region: 地域信息

    Returns:
        MonitorClient: 监控客户端实例
    """
    cred = credential.Credential(secret_id, secret_key)
    if not region:
        region = default_region or "ap-guangzhou"

    http_profile = HttpProfile()
    http_profile.endpoint = "monitor.tencentcloudapi.com"

    client_profile = ClientProfile()
    client_profile.httpProfile = http_profile
    client_profile.request_client = "MCP-Server"

    return monitor_client.MonitorClient(cred, region, client_profile) 