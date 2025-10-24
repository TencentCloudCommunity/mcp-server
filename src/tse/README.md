# TSF 注册配置治理 MCP
> 注册配置提供微服务与分布式场景下，云原生应用的动态服务发现、分布式配置管理和服务管理等能力。无改造、无缝平滑迁移、多语言接入等特性，助力业务轻松快速上云。

---

## Tools

### 1. `DescribeSREInstances`
- **详细描述**：查询引擎实例列表

### 2. `DescribeSREInstanceAccessAddress`
- **详细描述**：查询引擎实例访问地址

### 3. `DescribeNacosReplicas`
- **详细描述**：查询Nacos类型引擎实例副本信息

### 4. `DescribeNacosServerInterfaces`
- **详细描述**：查询nacos服务接口列表

### 5. `DescribeZookeeperReplicas`
- **详细描述**：查询Zookeeper类型注册引擎实例副本信息

### 6. `DescribeZookeeperServerInterfaces`
- **详细描述**：查询zookeeper服务接口列表

### 7. `DescribeCloudNativeAPIGateway`
- **详细描述**：获取云原生API网关实例信息

### 8. `DescribeCloudNativeAPIGatewayConfig`
- **详细描述**：获取云原生API网关实例网络配置信息

### 9. `DescribeNativeGatewayServerGroups`
- **详细描述**：查询云原生网关分组信息

### 10. `DescribeNativeGatewayServiceSources`
- **详细描述**：查询网关服务来源实例列表

### 11. `DescribeCloudNativeAPIGatewayServices`
- **详细描述**：查询云原生网关服务列表


### 12. `DescribeCloudNativeAPIGatewayRoutes`
- **详细描述**：查询云原生网关路由列表


---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/srag)

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
4. 打包镜像: `docker build -t mcp-server-tse:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-tse:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-tse": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```



## API使用参考

>https://cloud.tencent.com/document/product/1364/54625



## 许可证

MIT
