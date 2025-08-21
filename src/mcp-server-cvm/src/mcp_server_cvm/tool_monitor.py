"""
腾讯云监控相关操作工具模块
"""
import json
from tencentcloud.monitor.v20180724 import models
from . import client

def get_monitor_data(region: str, metric_name: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取监控数据基础方法

    Args:
        region: 地域信息
        metric_name: 指标名称
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    monitor_client = client.get_monitor_client(region)
    req = models.GetMonitorDataRequest()

    # 构造监控实例维度
    instances = []
    for instance_id in instance_ids:
        instances.append({
            "Dimensions": [{
                "Name": "InstanceId",
                "Value": instance_id
            }]
        })

    params = {
        "Namespace": "QCE/CVM",
        "MetricName": metric_name,
        "Period": period,
        "Instances": instances
    }

    req.from_json_string(json.dumps(params))
    resp = monitor_client.GetMonitorData(req)
    return resp.to_json_string()

def get_cpu_usage_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取CPU利用率指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "CpuUsage", instance_ids, period)

def get_cpu_loadavg_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取CPU一分钟平均负载指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "CpuLoadavg", instance_ids, period)

def get_cpu_loadavg5m_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取CPU五分钟平均负载指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "Cpuloadavg5m", instance_ids, period)

def get_cpu_loadavg15m_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取CPU十五分钟平均负载指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "Cpuloadavg15m", instance_ids, period)

def get_mem_used_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取内存使用量指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "MemUsed", instance_ids, period)

def get_mem_usage_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取内存利用率指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "MemUsage", instance_ids, period)

def get_cvm_disk_usage_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取磁盘利用率指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "CvmDiskUsage", instance_ids, period)

def get_disk_total_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取磁盘分区总容量指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "DiskTotal", instance_ids, period)

def get_disk_usage_data(region: str, instance_ids: list[str], period: int = 60) -> str:
    """
    获取磁盘分区已使用容量和总容量的百分比指标数据

    Args:
        region: 地域信息
        instance_ids: 实例ID列表
        period: 统计周期，默认60秒

    Returns:
        str: 监控数据JSON字符串
    """
    return get_monitor_data(region, "DiskUsage", instance_ids, period) 