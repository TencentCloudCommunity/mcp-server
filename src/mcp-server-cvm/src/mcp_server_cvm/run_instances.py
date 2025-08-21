"""
腾讯云 CVM 实例创建相关功能模块
"""
from typing import List, Optional
from tencentcloud.cvm.v20170312 import models
from tencentcloud.common import credential
from tencentcloud.common.profile.client_profile import ClientProfile
from tencentcloud.common.profile.http_profile import HttpProfile
from tencentcloud.cvm.v20170312 import cvm_client
import os

def run_instances(
    region: str,
    zone: str,
    instance_type: str,
    image_id: str,
    vpc_id: Optional[str] = None,
    subnet_id: Optional[str] = None,
    security_group_ids: Optional[List[str]] = None,
    password: Optional[str] = None,
    instance_name: Optional[str] = None,
    instance_charge_type: Optional[str] = None,
    instance_count: Optional[int] = None,
    dry_run: bool = False,
) -> dict:
    """
    创建腾讯云 CVM 实例

    Args:
        region: 地域
        zone: 可用区
        instance_type: 实例类型
        image_id: 镜像ID
        vpc_id: VPC ID，默认为 "DEFAULT"
        subnet_id: 子网ID，默认为 "DEFAULT"
        security_group_ids: 安全组ID列表
        password: 实例密码
        instance_name: 实例名称
        instance_charge_type: 计费类型
        instance_count: 创建实例数量
        dry_run: 是否只预检此次请求

    Returns:
        dict: API 响应结果
    """
    try:
        # 构造请求对象
        req = models.RunInstancesRequest()

        # 设置实例登录配置
        login_settings = models.LoginSettings()
        if password:
            login_settings.Password = password
        req.LoginSettings = login_settings

        # 设置实例所在位置
        placement = models.Placement()
        placement.Zone = zone
        req.Placement = placement

        # 设置VPC配置
        vpc = models.VirtualPrivateCloud()
        vpc.VpcId = vpc_id if vpc_id else "DEFAULT"
        vpc.SubnetId = subnet_id if subnet_id else "DEFAULT"
        req.VirtualPrivateCloud = vpc

        # 设置基本配置
        req.ImageId = image_id
        req.InstanceType = instance_type
        if security_group_ids:
            req.SecurityGroupIds = security_group_ids
        req.DryRun = dry_run

        # 设置可选参数
        if instance_charge_type:
            req.InstanceChargeType = instance_charge_type
        if instance_count:
            req.InstanceCount = instance_count
        if instance_name:
            req.InstanceName = instance_name

        # 创建客户端配置
        cred = credential.Credential(
            os.getenv("TENCENTCLOUD_SECRET_ID"),
            os.getenv("TENCENTCLOUD_SECRET_KEY")
        )
        http_profile = HttpProfile()
        client_profile = ClientProfile()
        client_profile.httpProfile = http_profile

        # 创建客户端
        client = cvm_client.CvmClient(cred, region, client_profile)

        # 发送请求
        response = client.RunInstances(req)
        return response.to_json_string()

    except Exception as e:
        raise Exception(f"创建实例失败: {str(e)}") 