"""
腾讯云 VPC 相关操作工具模块
"""
import json
import logging
from tencentcloud.vpc.v20170312 import models as vpc_models
from .client import get_vpc_client

logger = logging.getLogger(__name__)

def describe_security_groups(region: str, security_group_ids: list[str] = None) -> str:
    """查询安全组列表
    
    Args:
        region: 地域，如 ap-guangzhou
        security_group_ids: 安全组ID列表
    """
    client = get_vpc_client(region)
    req = vpc_models.DescribeSecurityGroupsRequest()
    
    params = {}
    if security_group_ids:
        params["SecurityGroupIds"] = security_group_ids
        
    if params:
        req.from_json_string(json.dumps(params))
    resp = client.DescribeSecurityGroups(req)
    return resp.to_json_string()

def describe_vpcs(region: str, vpc_ids: list[str] = None, is_default: bool = None, vpc_name: str = None) -> str:
    """查询VPC列表
    
    Args:
        region: 地域，如 ap-guangzhou
        vpc_ids: VPC ID列表
        is_default: 是否是默认VPC，True表示默认VPC，False表示非默认VPC，None表示不过滤
        vpc_name: VPC名称，用于过滤指定名称的VPC
    
    Returns:
        str: VPC列表的JSON字符串
    """
    client = get_vpc_client(region)
    req = vpc_models.DescribeVpcsRequest()
    
    params = {}
    filters = []
    
    if vpc_ids:
        params["VpcIds"] = vpc_ids
        
    if is_default is not None:
        filters.append({
            "Name": "is-default",
            "Values": ["true" if is_default else "false"]
        })
        
    if vpc_name:
        filters.append({
            "Name": "vpc-name",
            "Values": [vpc_name]
        })
        
    if filters:
        params["Filters"] = filters
        
    if params:
        req.from_json_string(json.dumps(params))
    resp = client.DescribeVpcs(req)
    return resp.to_json_string()

def describe_subnets(region: str, vpc_id: str = None, subnet_ids: list[str] = None, zone: str = None, is_default: bool = None, vpc_name: str = None) -> str:
    """查询子网列表
    
    Args:
        region: 地域，如 ap-guangzhou
        vpc_id: VPC ID，用于过滤指定VPC下的子网
        subnet_ids: 子网ID列表，用于查询指定子网的信息
        zone: 可用区，如 ap-guangzhou-1，用于过滤指定可用区的子网
        is_default: 是否是默认子网，True表示默认子网，False表示非默认子网，None表示不过滤
        vpc_name: VPC名称，用于过滤指定VPC名称下的子网
    
    Returns:
        str: 子网列表的JSON字符串
    """
    client = get_vpc_client(region)
    req = vpc_models.DescribeSubnetsRequest()
    
    params = {}
    filters = []
    
    if vpc_id:
        filters.append({
            "Name": "vpc-id",
            "Values": [vpc_id]
        })
    
    if zone:
        filters.append({
            "Name": "zone",
            "Values": [zone]
        })
        
    if is_default is not None:
        filters.append({
            "Name": "is-default",
            "Values": ["true" if is_default else "false"]
        })
        
    if vpc_name:
        filters.append({
            "Name": "vpc-name",
            "Values": [vpc_name]
        })
        
    if filters:
        params["Filters"] = filters
        
    if subnet_ids:
        params["SubnetIds"] = subnet_ids
        
    if params:
        req.from_json_string(json.dumps(params))
    resp = client.DescribeSubnets(req)
    return resp.to_json_string()

def describe_security_group_policies(region: str, security_group_id: str) -> str:
    """查询安全组规则
    
    Args:
        region: 地域
        security_group_id: 安全组ID
        
    Returns:
        str: API响应结果的JSON字符串
    """
    try:
        client = get_vpc_client(region)
        params = {
            "SecurityGroupId": security_group_id
        }
        resp = client.call("DescribeSecurityGroupPolicies", params)
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"查询安全组规则失败: {str(e)}")
        raise e

def create_security_group_policies(region: str, security_group_id: str, 
                                 security_group_policy_set: dict) -> str:
    """创建安全组规则
    
    Args:
        region: 地域
        security_group_id: 安全组ID
        security_group_policy_set: 安全组规则集合，包含 Ingress 和 Egress 规则
        
    Returns:
        str: API响应结果的JSON字符串
    """
    try:
        client = get_vpc_client(region)
        params = {
            "SecurityGroupId": security_group_id,
            "SecurityGroupPolicySet": security_group_policy_set
        }
        resp = client.call("CreateSecurityGroupPolicies", params)
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"创建安全组规则失败: {str(e)}")
        raise e

def create_security_group(region: str, group_name: str, group_description: str = None,
                         project_id: int = None, tags: list = None) -> str:
    """创建安全组
    
    Args:
        region: 地域
        group_name: 安全组名称
        group_description: 安全组描述
        project_id: 项目ID，默认为0
        tags: 标签列表
        
    Returns:
        str: API响应结果的JSON字符串
    """
    try:
        client = get_vpc_client(region)
        params = {
            "GroupName": group_name
        }
        if group_description:
            params["GroupDescription"] = group_description
        if project_id is not None:
            params["ProjectId"] = project_id
        if tags:
            params["Tags"] = tags
            
        resp = client.call("CreateSecurityGroup", params)
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"创建安全组失败: {str(e)}")
        raise e

def create_security_group_with_policies(region: str, group_name: str, 
                                      security_group_policy_set: dict,
                                      group_description: str = None,
                                      project_id: int = None,
                                      tags: list = None) -> str:
    """创建安全组并同时添加安全组规则
    
    Args:
        region: 地域
        group_name: 安全组名称
        security_group_policy_set: 安全组规则集合，包含 Ingress 和 Egress 规则
        group_description: 安全组描述
        project_id: 项目ID，默认为0
        tags: 标签列表
        
    Returns:
        str: API响应结果的JSON字符串
    """
    try:
        client = get_vpc_client(region)
        params = {
            "GroupName": group_name,
            "SecurityGroupPolicySet": security_group_policy_set
        }
        if group_description:
            params["GroupDescription"] = group_description
        if project_id is not None:
            params["ProjectId"] = project_id
        if tags:
            params["Tags"] = tags
            
        resp = client.call("CreateSecurityGroupWithPolicies", params)
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"创建安全组并添加规则失败: {str(e)}")
        raise e

def replace_security_group_policies(region: str, security_group_id: str,
                                  security_group_policy_set: dict) -> str:
    """批量修改安全组规则
    
    Args:
        region: 地域
        security_group_id: 安全组ID
        security_group_policy_set: 安全组规则集合，包含 Ingress 或 Egress 规则(单次请求只能替换单个方向的规则)
        
    Returns:
        str: API响应结果的JSON字符串
    """
    try:
        client = get_vpc_client(region)
        params = {
            "SecurityGroupId": security_group_id,
            "SecurityGroupPolicySet": security_group_policy_set
        }
        resp = client.call("ReplaceSecurityGroupPolicies", params)
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"批量修改安全组规则失败: {str(e)}")
        raise e 