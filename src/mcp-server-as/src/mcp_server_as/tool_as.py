"""
腾讯云 AS 相关操作工具模块
"""
import json
from tencentcloud.autoscaling.v20180419 import models
from .client import get_as_client

def create_auto_scaling_group(region: str, params: dict) -> str:
    """创建伸缩组
    
    Args:
        region: 地域
        params: 创建参数，包含：
            - AutoScalingGroupName: 伸缩组名称
            - LaunchConfigurationId: 启动配置ID
            - MaxSize: 最大实例数
            - MinSize: 最小实例数
            - VpcId: VPC ID
            - SubnetIds: 子网ID列表
            - DesiredCapacity: 期望实例数
            
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_as_client(region)
    req = models.CreateAutoScalingGroupRequest()
    req.from_json_string(json.dumps(params))
    resp = client.CreateAutoScalingGroup(req)
    return resp.to_json_string()

def describe_auto_scaling_groups(region: str, auto_scaling_group_ids: list[str] = None,
                               filters: list[dict] = None, limit: int = None,
                               offset: int = None) -> str:
    """查询伸缩组
    
    Args:
        region: 地域
        auto_scaling_group_ids: 伸缩组ID列表
        filters: 过滤条件
        limit: 返回数量
        offset: 偏移量
        
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_as_client(region)
    req = models.DescribeAutoScalingGroupsRequest()
    
    params = {}
    if auto_scaling_group_ids:
        params["AutoScalingGroupIds"] = auto_scaling_group_ids
    if filters:
        params["Filters"] = filters
    if limit is not None:
        params["Limit"] = limit
    if offset is not None:
        params["Offset"] = offset
        
    if params:
        req.from_json_string(json.dumps(params))
    resp = client.DescribeAutoScalingGroups(req)
    return resp.to_json_string()

def modify_auto_scaling_group(region: str, auto_scaling_group_id: str,
                            params: dict) -> str:
    """修改伸缩组
    
    Args:
        region: 地域
        auto_scaling_group_id: 伸缩组ID
        params: 修改参数
        
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_as_client(region)
    req = models.ModifyAutoScalingGroupRequest()
    
    params["AutoScalingGroupId"] = auto_scaling_group_id
    req.from_json_string(json.dumps(params))
    resp = client.ModifyAutoScalingGroup(req)
    return resp.to_json_string()

def enable_auto_scaling_group(region: str, auto_scaling_group_id: str) -> str:
    """启用伸缩组
    
    Args:
        region: 地域
        auto_scaling_group_id: 伸缩组ID
        
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_as_client(region)
    req = models.EnableAutoScalingGroupRequest()
    
    params = {
        "AutoScalingGroupId": auto_scaling_group_id
    }
    req.from_json_string(json.dumps(params))
    resp = client.EnableAutoScalingGroup(req)
    return resp.to_json_string()

def disable_auto_scaling_group(region: str, auto_scaling_group_id: str) -> str:
    """停用伸缩组
    
    Args:
        region: 地域
        auto_scaling_group_id: 伸缩组ID
        
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_as_client(region)
    req = models.DisableAutoScalingGroupRequest()
    
    params = {
        "AutoScalingGroupId": auto_scaling_group_id
    }
    req.from_json_string(json.dumps(params))
    resp = client.DisableAutoScalingGroup(req)
    return resp.to_json_string()

def execute_scaling_policy(region: str, auto_scaling_group_id: str,
                         operation: str, adjustment_type: str = None,
                         adjustment_value: int = None) -> str:
    """执行伸缩策略
    
    Args:
        region: 地域
        auto_scaling_group_id: 伸缩组ID
        operation: 操作类型，取值范围：SCALE_OUT（扩容），SCALE_IN（缩容）
        adjustment_type: 调整类型，取值范围：CHANGE_IN_CAPACITY, EXACT_CAPACITY, PERCENT_CHANGE_IN_CAPACITY
        adjustment_value: 调整值，正数表示增加实例，负数表示减少实例
        
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_as_client(region)
    req = models.ExecuteScalingPolicyRequest()
    
    params = {
        "AutoScalingGroupId": auto_scaling_group_id,
        "Operation": operation
    }
    if adjustment_type:
        params["AdjustmentType"] = adjustment_type
    if adjustment_value is not None:
        params["AdjustmentValue"] = adjustment_value
        
    req.from_json_string(json.dumps(params))
    resp = client.ExecuteScalingPolicy(req)
    return resp.to_json_string()

def modify_desired_capacity(region: str, auto_scaling_group_id: str,
                          desired_capacity: int) -> str:
    """修改期望实例数
    
    Args:
        region: 地域
        auto_scaling_group_id: 伸缩组ID
        desired_capacity: 期望实例数
        
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_as_client(region)
    req = models.ModifyDesiredCapacityRequest()
    
    params = {
        "AutoScalingGroupId": auto_scaling_group_id,
        "DesiredCapacity": desired_capacity
    }
    req.from_json_string(json.dumps(params))
    resp = client.ModifyDesiredCapacity(req)
    return resp.to_json_string()
