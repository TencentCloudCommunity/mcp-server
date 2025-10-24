# VPN 连接
> VPN 连接（VPN Connections）是一种基于网络隧道技术，实现本地数据中心与腾讯云上资源连通的传输服务，它能帮您在 Internet 上快速构建一条安全、可靠的加密通道。VPN 连接具有配置简单，云端配置实时生效、可靠性高等特点，其网关可用性达到 99.95%，保证稳定、持续的业务连接，帮您轻松实现异地容灾、混合云部署等复杂业务场景。

---

## Tools

### 1. `DescribeVpnGateways`
- **详细描述**：查询VPN网关

### 2. `DescribeVpnConnections`
- **详细描述**：查询VPN通道列表

### 3. `DescribeVpnGatewayRoutes`
- **详细描述**：查询VPN网关路由

### 4. `DescribeCustomerGateways`
- **详细描述**：查询对端网关

### 5. `DescribeVpnGatewaySslServers`
- **详细描述**：查询SSL-VPN SERVER 列表

### 6. `DescribeVpnGatewaySslClients`
- **详细描述**：查询SSL-VPN-CLIENT 列表

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/vpn)

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
4. 打包镜像: `docker build -t mcp-server-vpn:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-vpn:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-vpn": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```



## API使用参考

>https://cloud.tencent.com/document/product/xxxx



## 许可证

MIT
