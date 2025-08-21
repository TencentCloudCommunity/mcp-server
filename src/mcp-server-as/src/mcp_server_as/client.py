"""
腾讯云 AS 客户端模块
"""
import os
from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.autoscaling.v20180419 import autoscaling_client

def get_as_client(region: str) -> autoscaling_client.AutoscalingClient:
    """获取 AS 客户端
    
    Args:
        region: 地域
        
    Returns:
        AutoscalingClient: AS 客户端实例
    """
    # 从环境变量获取密钥
    secret_id = os.getenv("TENCENTCLOUD_SECRET_ID")
    secret_key = os.getenv("TENCENTCLOUD_SECRET_KEY")
    
    if not secret_id or not secret_key:
        raise ValueError("请设置环境变量: TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY")
    
    # 实例化一个认证对象
    cred = credential.Credential(secret_id, secret_key)
    
    # 实例化一个http选项
    httpProfile = HttpProfile()
    httpProfile.endpoint = "as.tencentcloudapi.com"
    
    # 实例化一个client选项
    clientProfile = ClientProfile()
    clientProfile.httpProfile = httpProfile
    clientProfile.request_client = "MCP-Server"
    
    # 实例化要请求产品的client对象
    return autoscaling_client.AutoscalingClient(cred, region, clientProfile)
