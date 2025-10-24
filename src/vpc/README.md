# 私有网络 VPC MCP
> 私有网络（Virtual Private Cloud，VPC）是基于腾讯云构建的专属云上网络空间，为您在腾讯云上的资源提供网络服务，不同私有网络间完全逻辑隔离。作为您在云上的专属网络空间，您可以通过软件定义网络的方式管理您的私有网络 VPC，实现 IP 地址、子网、路由表、网络 ACL 、流日志等功能的配置管理。私有网络还支持多种方式连接 Internet，如弹性 IP 、NAT 网关等。同时，您也可以通过 VPN 连接或专线接入连通腾讯云与您本地的数据中心，灵活构建混合云。

---

## Tools

### 1. `DescribeVpcs`
- **详细描述**：查询VPC列表

### 2. `DescribeSubnets`
- **详细描述**：查询子网列表

### 3. `DescribeRouteTables`
- **详细描述**：查询路由表

### 4. `DescribeNetworkInterfaces`
- **详细描述**：查询弹性网卡列表

### 5. `DescribeSecurityGroups`
- **详细描述**：查询IP地址模板

### 6. `DescribeAddressTemplates`
- **详细描述**：查询HAVIP列表

### 7. `DescribeHaVips`
- **详细描述**：查询实例列表

### 8. `DescribeNetworkAcls`
- **详细描述**：查询网络ACL列表




---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/vpc)

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
4. 打包镜像: `docker build -t mcp-server-vpc:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-vpc:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-vpc": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```



## API使用参考

>https://cloud.tencent.com/document/product/215/15755



## 许可证

MIT
