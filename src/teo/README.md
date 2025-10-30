# 边缘安全加速平台 EO
> 边缘安全加速平台 EO (TencentCloud EdgeOne)基于腾讯云遍布全球的边缘节点，提供域名解析、动静态智能加速、TCP/UDP 四层加速、DDoS/CC/Web/Bot 防护、Pages、边缘函数计算等边缘一体化服务，可帮助客户更快速、更安全、更灵活地响应用户请求。

---

## Tools

### 1. `DescribeIdentifications`
- **详细描述**：查询站点的验证信息

### 2. `DescribeZones`
- **详细描述**：查询站点列表

### 3. `DescribeAccelerationDomains`
- **详细描述**：查询加速域名列表

### 4. `DescribeFunctions`
- **详细描述**：查询边缘函数列表

### 5. `DescribeAliasDomains`
- **详细描述**：查询别称域名信息列表

### 6. `DescribePrefetchTasks`
- **详细描述**：查询预热任务状态

### 7. `DescribeDDoSAttackEvent`
- **详细描述**：查询DDoS攻击事件列表

### 8. `DescribeTimingL7AnalysisData`
- **详细描述**：查询流量分析时序数据

### 9. `DescribeDefaultCertificates`
- **详细描述**：查询默认证书列表

### 10. `DescribeDnsRecords`
- **详细描述**：查询 DNS 记录列表

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/teo)

---

## 使用


### 参数获取方式：

a、密钥获取：腾讯云控制台当中，生成腾讯云SecretId和SecretKey
> 生成地址：https://console.cloud.tencent.com/cam/capi

b、地域列表的映射关系：
如广州地域，则地域字段的内容应该为：ap-guangzhou
> https://cloud.tencent.com/document/product/1596/77930

### 部署与配置
1. 生成配置文件
   `cp .env.example .env`
2. 将密钥填写到 `.env` 配置文件中
3. 安装 [docker](https://www.docker.com/)
4. 打包镜像: `docker build -t mcp-server-teo:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-teo:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-teo": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/product/1552/80731



## 许可证

MIT