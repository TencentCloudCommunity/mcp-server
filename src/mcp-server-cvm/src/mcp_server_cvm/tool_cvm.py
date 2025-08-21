"""
腾讯云 CVM 相关操作工具模块
"""
import json
from tencentcloud.cvm.v20170312 import cvm_client, models as cvm_models
from .client import get_cvm_client, get_common_client
from asyncio.log import logger
from tencentcloud.common.exception.tencent_cloud_sdk_exception import TencentCloudSDKException
import random
import string
import re
from typing import Optional, Tuple
from . import tool_vpc

def describe_regions() -> str:
    """查询地域列表"""
    client = get_cvm_client("ap-guangzhou")  # 使用默认地域
    req = cvm_models.DescribeRegionsRequest()
    resp = client.DescribeRegions(req)
    return resp.to_json_string()

def describe_zones(region: str) -> str:
    """查询可用区列表"""
    client = get_cvm_client(region)
    req = cvm_models.DescribeZonesRequest()
    resp = client.DescribeZones(req)
    return resp.to_json_string()

def describe_instances(region: str, offset: int, limit: int, instance_ids: list[str]) -> str:
    """查询实例列表"""
    client = get_cvm_client(region)
    req = cvm_models.DescribeInstancesRequest()
    
    params = {
        "Offset": offset,
        "Limit": limit
    }
    if instance_ids:
        params["InstanceIds"] = instance_ids
        
    req.from_json_string(json.dumps(params))
    resp = client.DescribeInstances(req)
    
    # 解析返回结果并添加 LoginUrl
    result = json.loads(resp.to_json_string())
    if "Response" in result and "InstanceSet" in result["Response"]:
        for instance in result["Response"]["InstanceSet"]:
            instance_id = instance.get("InstanceId")
            if instance_id:
                instance["LoginUrl"] = f"https://orcaterm.cloud.tencent.com/terminal?type=cvm&instanceId={instance_id}&region={region}"
    
    return json.dumps(result)

def describe_images(region: str, image_ids: list[str] = None, image_type: str = None,
                  platform: str = None, image_name: str = None,
                  offset: int = None, limit: int = None) -> str:
    """查询镜像列表
    
    Args:
        region (str): 地域ID
        image_ids (list[str], optional): 镜像ID列表
        image_type (str, optional): 镜像类型
        platform (str, optional): 操作系统平台
        image_name (str, optional): 镜像名称
        offset (int, optional): 偏移量
        limit (int, optional): 返回数量
        
    Returns:
        str: API响应结果的JSON字符串
    """
    client = get_cvm_client(region)
    req = cvm_models.DescribeImagesRequest()
    
    # 构建参数
    params = {}
    filters = []
    
    # 添加过滤条件
    if image_type:
        filters.append({
            "Name": "image-type",
            "Values": [image_type]
        })
    if platform:
        filters.append({
            "Name": "platform",
            "Values": [platform]
        })
    if image_name:
        filters.append({
            "Name": "image-name",
            "Values": [image_name]
        })
        
    if filters:
        params["Filters"] = filters
    if image_ids:
        params["ImageIds"] = image_ids
    if offset is not None:
        params["Offset"] = offset
    if limit is not None:
        params["Limit"] = limit
        
    if params:
        req.from_json_string(json.dumps(params))
    resp = client.DescribeImages(req)
    return resp.to_json_string()

def describe_instance_type_configs(region: str, zone: str = None, instance_family: str = None) -> str:
    """查询实例机型配置"""
    client = get_cvm_client(region)
    req = cvm_models.DescribeInstanceTypeConfigsRequest()
    
    params = {}
    if zone:
        params["Filters"] = [{
            "Name": "zone",
            "Values": [zone]
        }]
    if instance_family:
        if "Filters" not in params:
            params["Filters"] = []
        params["Filters"].append({
            "Name": "instance-family",
            "Values": [instance_family]
        })
        
    if params:
        req.from_json_string(json.dumps(params))
    resp = client.DescribeInstanceTypeConfigs(req)
    return resp.to_json_string()

def reboot_instances(region: str, instance_ids: list[str], stop_type: str) -> str:
    """重启实例"""
    client = get_cvm_client(region)
    req = cvm_models.RebootInstancesRequest()
    
    params = {
        "InstanceIds": instance_ids,
        "StopType": stop_type
    }
    req.from_json_string(json.dumps(params))
    resp = client.RebootInstances(req)
    return resp.to_json_string()

def start_instances(region: str, instance_ids: list[str]) -> str:
    """启动实例"""
    client = get_cvm_client(region)
    req = cvm_models.StartInstancesRequest()
    
    params = {
        "InstanceIds": instance_ids
    }
    req.from_json_string(json.dumps(params))
    resp = client.StartInstances(req)
    return resp.to_json_string()

def stop_instances(region: str, instance_ids: list[str], stop_type: str, stopped_mode: str) -> str:
    """关闭实例"""
    client = get_cvm_client(region)
    req = cvm_models.StopInstancesRequest()
    
    params = {
        "InstanceIds": instance_ids,
        "StopType": stop_type,
        "StoppedMode": stopped_mode
    }
    req.from_json_string(json.dumps(params))
    resp = client.StopInstances(req)
    return resp.to_json_string()

def terminate_instances(region: str, instance_ids: list[str]) -> str:
    """销毁实例"""
    client = get_cvm_client(region)
    req = cvm_models.TerminateInstancesRequest()
    
    params = {
        "InstanceIds": instance_ids
    }
    req.from_json_string(json.dumps(params))
    resp = client.TerminateInstances(req)
    return resp.to_json_string()

def reset_instances_password(region: str, instance_ids: list[str], password: str, force_stop: bool) -> str:
    """重置实例密码"""
    client = get_cvm_client(region)
    req = cvm_models.ResetInstancesPasswordRequest()
    
    params = {
        "InstanceIds": instance_ids,
        "Password": password,
        "ForceStop": force_stop
    }
    req.from_json_string(json.dumps(params))
    resp = client.ResetInstancesPassword(req)
    return resp.to_json_string()

def run_instances(region: str, params: dict) -> str:
    """创建实例"""
    try:
        from .run_instances import run_instances as run_instances_impl
        return run_instances_impl(
            region=region,
            zone=params.get("Zone"),
            instance_type=params.get("InstanceType"),
            image_id=params.get("ImageId"),
            vpc_id=params.get("VpcId"),
            subnet_id=params.get("SubnetId"),
            security_group_ids=params.get("SecurityGroupIds"),
            password=params.get("Password"),
            instance_name=params.get("InstanceName"),
            instance_charge_type=params.get("InstanceChargeType"),
            instance_count=params.get("InstanceCount"),
            dry_run=params.get("DryRun", False)
        )
    except Exception as e:
        logger.error(f"创建实例失败: {str(e)}")
        raise e

def reset_instance(region: str, instance_id: str, image_id: str, password: str = None) -> str:
    """重装实例操作系统
    
    Args:
        region (str): 实例所在地域
        instance_id (str): 实例ID
        image_id (str): 重装使用的镜像ID
        password (str, optional): 实例重装后的密码。如果不指定，保持原密码不变
        
    Returns:
        str: API响应结果的JSON字符串
        
    Raises:
        Exception: 当API调用失败时抛出异常
    """
    try:
        client = get_cvm_client(region)
        # 设置实例登录配置
        login_settings = cvm_models.LoginSettings()
        if password:
            login_settings.Password = password
        
        req = cvm_models.ResetInstanceRequest()
        req.InstanceId = instance_id
        req.ImageId = image_id
        req.LoginSettings = login_settings
        resp = client.ResetInstance(req)
        return resp.to_json_string()
    except Exception as e:
        logger.error(f"重装实例操作系统失败: {str(e)}")
        raise e

def inquiry_price_run_instances(region: str, params: dict) -> str:
    """创建实例询价
    
    Args:
        region (str): 实例所在地域
        params (dict): 询价参数，包含：
            - Zone: 可用区
            - InstanceType: 实例机型
            - ImageId: 镜像ID
            - SystemDisk: 系统盘配置
            - InstanceChargeType: 实例计费类型
            - InstanceChargePrepaid: 预付费配置（仅当 InstanceChargeType 为 PREPAID 时需要）
            
    Returns:
        str: API响应结果的JSON字符串
        
    Raises:
        Exception: 当API调用失败时抛出异常
    """
    try:
        client = get_cvm_client(region)
        req = cvm_models.InquiryPriceRunInstancesRequest()
        
        # 设置基础配置
        req.Placement = cvm_models.Placement()
        req.Placement.Zone = params.get("Zone")
        
        req.InstanceType = params.get("InstanceType")
        req.ImageId = params.get("ImageId")
        
        # 设置系统盘
        system_disk = params.get("SystemDisk", {})
        if system_disk:
            req.SystemDisk = cvm_models.SystemDisk()
            req.SystemDisk.DiskType = system_disk.get("DiskType", "CLOUD_PREMIUM")
            req.SystemDisk.DiskSize = system_disk.get("DiskSize", 50)
            
        # 设置计费类型
        req.InstanceChargeType = params.get("InstanceChargeType", "POSTPAID_BY_HOUR")
        
        # 如果是包年包月，设置购买时长
        if req.InstanceChargeType == "PREPAID":
            prepaid = params.get("InstanceChargePrepaid", {})
            req.InstanceChargePrepaid = cvm_models.InstanceChargePrepaid()
            req.InstanceChargePrepaid.Period = prepaid.get("Period", 1)
            req.InstanceChargePrepaid.RenewFlag = prepaid.get("RenewFlag", "NOTIFY_AND_MANUAL_RENEW")
            
        resp = client.InquiryPriceRunInstances(req)
        return resp.to_json_string()
    except Exception as e:
        logger.error(f"创建实例询价失败: {str(e)}")
        raise e

def create_diagnostic_reports(region: str, instance_ids: list[str]) -> str:
    """创建实例诊断报告
    
    Args:
        region: 地域
        instance_ids: 实例ID列表
        
    Returns:
        str: API响应结果的JSON字符串
    """
    try:
        client = get_cvm_client(region)
        params = {
            "InstanceIds": instance_ids
        }
        # 使用通用请求方式调用
        resp = client.call("CreateDiagnosticReports", params)
        # 如果响应是字节类型，先解码成字符串
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"创建实例诊断报告失败: {str(e)}")
        raise e

def describe_diagnostic_reports(region: str, report_ids: list[str] = None,
                              filters: list[dict] = None, vague_instance_name: str = None,
                              offset: int = None, limit: int = None,
                              cluster_diagnostic_report_ids: list[str] = None,
                              scenario_id: int = None) -> str:
    """查询实例诊断报告
    
    Args:
        region: 地域
        report_ids: 实例健康检测报告ID列表，如：["dr-rfmme2si"]。每次请求批量报告ID的上限为100
        filters: 过滤条件列表，支持的过滤条件：
            - instance-id: 按实例ID过滤，如：ins-8jqq9ajy
            - instance-name: 按实例名称过滤，如：my-ins
            - instance-health-status: 按实例健康状态过滤，可选值：Normal, Warn, Critical
            - report-status: 按报告状态过滤，可选值：Querying, Finished
            - cluster-ids: 按集群ID过滤，如：['hpc-rltlmf6v']
        vague_instance_name: 模糊实例别名
        offset: 偏移量，默认为0
        limit: 返回数量，默认为20，最大值为100
        cluster_diagnostic_report_ids: 集群健康检测报告ID列表，如：["cr-rfmme2si"]
        scenario_id: 检测场景ID，默认为1表示对CVM进行全面体检，200为集群一致性检测场景
        
    Returns:
        str: API响应结果的JSON字符串
        
    Note:
        report_ids 和 cluster_diagnostic_report_ids 不能与 filters 或 vague_instance_name 联合使用
    """
    try:
        client = get_cvm_client(region)
        params = {}
        
        if report_ids:
            params["ReportIds"] = report_ids
        if filters:
            params["Filters"] = filters
        if vague_instance_name:
            params["VagueInstanceName"] = vague_instance_name
        if offset is not None:
            params["Offset"] = offset
        if limit is not None:
            params["Limit"] = limit
        if cluster_diagnostic_report_ids:
            params["ClusterDiagnosticReportIds"] = cluster_diagnostic_report_ids
        if scenario_id is not None:
            params["ScenarioId"] = scenario_id
            
        # 使用通用请求方式调用
        resp = client.call("DescribeDiagnosticReports", params)
        # 如果响应是字节类型，先解码成字符串
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"查询实例诊断报告失败: {str(e)}")
        raise e

def describe_recommend_zone_instance_types(
    region: str,
    zone: str,
    instance_type: str,
    instance_charge_type: str,
    is_recommend_under_stock: bool = False,
    is_compare: bool = False
) -> str:
    """
    推荐用户在相同或相近地域购买相同或相似机型

    Args:
        region (str): 地域，如 ap-guangzhou
        zone (str): 当前机型所在可用区，如 ap-guangzhou-2
        instance_type (str): 实例规格，如 SA2.LARGE8
        instance_charge_type (str): 实例计费类型，如 PREPAID
        is_recommend_under_stock (bool, optional): 是否推荐低库存机型数据，默认False
        is_compare (bool, optional): 是否属于机型比对场景，默认False

    Returns:
        str: API响应结果的JSON字符串
    """
    try:
        client = get_common_client(region, product="cvm", version="2019-12-12")
        params = {
            "Zone": zone,
            "InstanceType": instance_type,
            "InstanceChargeType": instance_charge_type
        }
        if is_recommend_under_stock is not None:
            params["IsRecommendUnderStock"] = is_recommend_under_stock
        if is_compare is not None:
            params["IsCompare"] = is_compare
        # 使用通用请求方式调用
        resp = client.call("DescribeRecommendZoneInstanceTypes", params)
        if isinstance(resp, bytes):
            resp = resp.decode('utf-8')
        return resp
    except Exception as e:
        logger.error(f"推荐可用区机型失败: {str(e)}")
        raise e

def _generate_random_password() -> str:
    """生成随机密码
    
    生成的密码满足腾讯云的密码要求：字母+数字+特殊字符，长度8-30位
    """
    length = 12
    letters = string.ascii_letters
    digits = string.digits
    special = "!@#$%^&*"
    
    # 确保每种字符都至少有一个
    password = [
        random.choice(letters.lower()),
        random.choice(letters.upper()),
        random.choice(digits),
        random.choice(special)
    ]
    
    # 填充剩余长度
    characters = letters + digits + special
    for i in range(length - 4):
        password.append(random.choice(characters))
        
    # 打乱顺序
    random.shuffle(password)
    return ''.join(password)

def _parse_instance_type(type_desc: str) -> Optional[Tuple[int, int]]:
    """解析用户描述的实例类型
    
    Args:
        type_desc: 用户描述，如 "2核4G"
        
    Returns:
        (cpu核数, 内存大小(GB))，如果无法解析则返回None
    """
    pattern = r"(\d+)核(\d+)G"
    match = re.match(pattern, type_desc)
    if match:
        return (int(match.group(1)), int(match.group(2)))
    return None

def _find_matching_instance_type(region: str, zone: str, cpu: int, memory: int) -> Optional[str]:
    """查找匹配的实例类型
    
    Args:
        region: 地域
        zone: 可用区
        cpu: CPU核数
        memory: 内存大小(GB)
        
    Returns:
        匹配的实例类型，如果没找到则返回None
    """
    configs = json.loads(describe_instance_type_configs(region, zone))
    for config in configs.get("InstanceTypeConfigSet", []):
        if config["CPU"] == cpu and config["Memory"] == memory:
            # 优先选择 SA5 系列
            if config["InstanceFamily"] == "SA5":
                return config["InstanceType"]
            # 其次选择 BF1 系列
            elif config["InstanceFamily"] == "BF1":
                return config["InstanceType"]
    return None

def _get_default_ubuntu_image(region: str) -> Optional[str]:
    """获取默认的Ubuntu公共镜像
    
    Args:
        region: 地域
        
    Returns:
        镜像ID，如果没找到则返回None
    """
    images = json.loads(describe_images(
        region=region,
        image_type="PUBLIC_IMAGE",
        platform="Ubuntu"
    ))
    
    # 优先选择 Ubuntu 22.04 LTS
    for image in images.get("ImageSet", []):
        if "Ubuntu 22.04 LTS" in image["ImageName"]:
            return image["ImageId"]
            
    # 如果没有22.04，选择任意Ubuntu镜像
    if images.get("ImageSet"):
        return images["ImageSet"][0]["ImageId"]
    return None

def _get_available_subnet(region: str, zone: str) -> Optional[Tuple[str, str]]:
    """获取可用的子网
    
    Args:
        region: 地域
        zone: 可用区
        
    Returns:
        (VpcId, SubnetId)，如果没找到则返回None
    """
    # 先查找默认VPC下的子网
    subnets = json.loads(tool_vpc.describe_subnets(
        region=region,
        zone=zone,
        is_default=True
    ))
    
    if not subnets.get("SubnetSet"):
        # 如果没有默认子网，查找所有子网
        subnets = json.loads(tool_vpc.describe_subnets(
            region=region,
            zone=zone
        ))
    
    for subnet in subnets.get("SubnetSet", []):
        # 检查子网是否还有可用IP
        available_ip = subnet["AvailableIpAddressCount"]
        if available_ip > 0:
            return (subnet["VpcId"], subnet["SubnetId"])
    
    return None

def _find_zone_for_instance_type(region: str, instance_type: str) -> Optional[str]:
    """查找支持指定实例类型的可用区
    
    Args:
        region: 地域
        instance_type: 实例类型
        
    Returns:
        支持该实例类型的可用区，如果没找到则返回None
    """
    # 获取所有可用区
    zones = json.loads(describe_zones(region))
    
    for zone_info in zones.get("ZoneSet", []):
        if zone_info["ZoneState"] != "AVAILABLE":
            continue
            
        zone = zone_info["Zone"]
        # 检查该可用区是否支持指定实例类型
        configs = json.loads(describe_instance_type_configs(region, zone))
        for config in configs.get("InstanceTypeConfigSet", []):
            if config["InstanceType"] == instance_type:
                return zone
                
    return None

def _find_available_instance_type(region: str, zone: str = None, cpu: int = None, memory: int = None) -> Tuple[str, str]:
    """查找可用的实例类型和可用区组合
    
    Args:
        region: 地域
        zone: 期望的可用区（可选）
        cpu: CPU核数（可选）
        memory: 内存大小GB（可选）
        
    Returns:
        (实例类型, 可用区)
        
    Raises:
        ValueError: 当找不到合适的实例类型和可用区组合时
    """
    # 获取所有可用区
    zones = json.loads(describe_zones(region))
    available_zones = [z["Zone"] for z in zones.get("ZoneSet", []) if z["ZoneState"] == "AVAILABLE"]
    
    # 如果指定了可用区，将其放在列表最前面
    if zone and zone in available_zones:
        available_zones.remove(zone)
        available_zones.insert(0, zone)
    
    # 遍历每个可用区
    for current_zone in available_zones:
        configs = json.loads(describe_instance_type_configs(region, current_zone))
        
        for config in configs.get("InstanceTypeConfigSet", []):
            # 如果指定了CPU和内存，检查是否匹配
            if cpu is not None and memory is not None:
                if config["CPU"] != cpu or config["Memory"] != memory:
                    continue
                    
            # 优先选择 SA5 系列
            if config["InstanceFamily"] == "SA5":
                return config["InstanceType"], current_zone
                
            # 其次选择 BF1 系列
            elif config["InstanceFamily"] == "BF1":
                return config["InstanceType"], current_zone
                
    # 如果找不到匹配的配置，抛出异常
    if cpu is not None and memory is not None:
        raise ValueError(f"在地域 {region} 找不到匹配的实例类型: {cpu}核{memory}G")
    else:
        raise ValueError(f"在地域 {region} 找不到可用的实例类型")

def quick_run_instance(
    region: str = "ap-guangzhou",
    zone: str = None,
    instance_type: str = None,
    image_id: str = None,
    vpc_id: str = None,
    subnet_id: str = None,
    password: str = None,
    instance_name: str = None,
    instance_charge_type: str = "POSTPAID_BY_HOUR"
) -> str:
    """快速创建实例
    
    Args:
        region: 地域，默认为广州
        zone: 可用区，如果不指定则自动选择
        instance_type: 实例类型，可以是具体的类型如"SA5.MEDIUM2"，也可以是描述如"2核4G"
        image_id: 镜像ID，如果不指定则使用Ubuntu公共镜像
        vpc_id: VPC ID，如果不指定则自动选择
        subnet_id: 子网ID，如果不指定则自动选择
        password: 实例密码，如果不指定则随机生成
        instance_name: 实例名称
        instance_charge_type: 计费类型，默认按量计费
        
    Returns:
        创建结果的JSON字符串
    """
    # 1. 处理实例类型和可用区
    if instance_type:
        # 如果用户指定了实例类型
        parsed = _parse_instance_type(instance_type)
        if parsed:
            # 如果是"2核4G"这样的描述格式
            cpu, memory = parsed
            instance_type, zone = _find_available_instance_type(region, zone, cpu, memory)
        else:
            # 如果是具体的实例类型，确保有可用区支持它
            if zone:
                # 检查指定的可用区是否支持该实例类型
                configs = json.loads(describe_instance_type_configs(region, zone))
                supported = False
                for config in configs.get("InstanceTypeConfigSet", []):
                    if config["InstanceType"] == instance_type:
                        supported = True
                        break
                if not supported:
                    # 如果指定的可用区不支持该实例类型，尝试找其他可用区
                    new_zone = _find_zone_for_instance_type(region, instance_type)
                    if new_zone:
                        zone = new_zone
                    else:
                        raise ValueError(f"找不到支持实例类型 {instance_type} 的可用区")
            else:
                # 查找支持该实例类型的可用区
                zone = _find_zone_for_instance_type(region, instance_type)
                if not zone:
                    raise ValueError(f"找不到支持实例类型 {instance_type} 的可用区")
    else:
        # 如果用户没有指定实例类型，查找可用的实例类型和可用区组合
        instance_type, zone = _find_available_instance_type(region, zone)
            
    # 3. 处理镜像
    if not image_id:
        image_id = _get_default_ubuntu_image(region)
        if not image_id:
            raise ValueError("找不到默认的Ubuntu镜像")
            
    # 4. 处理VPC和子网
    if not vpc_id or not subnet_id:
        vpc_subnet = _get_available_subnet(region, zone)
        if vpc_subnet:
            vpc_id, subnet_id = vpc_subnet
        else:
            raise ValueError(f"在可用区 {zone} 找不到可用的子网")
            
    # 5. 处理密码
    if not password:
        password = _generate_random_password()
        
    # 6. 创建实例
    params = {
        "Zone": zone,
        "InstanceType": instance_type,
        "ImageId": image_id,
        "VpcId": vpc_id,
        "SubnetId": subnet_id,
        "InstanceChargeType": instance_charge_type,
        "Password": password
    }
    
    if instance_name:
        params["InstanceName"] = instance_name
        
    result = run_instances(region, params)
    
    # 在返回结果中包含生成的密码
    result_dict = json.loads(result)
    result_dict["Password"] = password
    
    return json.dumps(result_dict, ensure_ascii=False) 