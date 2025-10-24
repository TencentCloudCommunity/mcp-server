# TSF 应用管理 MCP
> 应用管理是一个围绕应用和微服务的 PaaS 平台，提供一站式应用全生命周期管理能力和数据化运营支持，提供多维度应用和服务的监控数据，帮助企业创建和管理云资源，助力企业充分聚焦核心业务本身。

---

## Tools

### 1. `DescribeApplications`
- **详细描述**：获取应用列表

### 2. `DescribeApplicationAttribute`
- **详细描述**：获取应用列表其它字段

### 3. `DescribeSimpleClusters`
- **详细描述**：查询简单集群列表

### 4. `DescribeSimpleNamespaces`
- **详细描述**：查询简单命名空间列表

### 5. `DescribeImageRepository`
- **详细描述**：查询镜像仓库列表

### 6. `DescribeImageTags`
- **详细描述**：查询镜像版本列表

### 7. `DescribeRepositories`
- **详细描述**：查询仓库列表

### 8. `DescribePkgs`
- **详细描述**：获取某个应用的程序包信息列表

### 9. `DescribeGroups`
- **详细描述**：获取虚拟机部署组列表

### 10. `DescribeBusinessLogConfigs`
- **详细描述**：查询日志配置项列表

### 11. `SearchBusinessLog`
- **详细描述**：业务日志搜索

### 12. `SearchStdoutLog`
- **详细描述**：标准输出日志搜索

### 13. `DescribeInvocationMetricDataCurve`
- **详细描述**：查询调用指标数据变化曲线

### 14. `DescribePrograms`
- **详细描述**：查询数据集列表

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/msas)

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
4. 打包镜像: `docker build -t mcp-server-tsf:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-tsf:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-tsf": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/api/649/36037



## 许可证

MIT
