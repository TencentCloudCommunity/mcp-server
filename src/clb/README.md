# 负载均衡
> 腾讯云负载均衡（Cloud Load Balancer，CLB）提供安全快捷的四七层流量分发服务，访问流量经由 CLB 可以自动分配到多台后端服务器上，扩展系统的服务能力并消除单点故障

---

## Tools

### 1. `DescribeCustomizedConfigAssociateList`
- **详细描述**：拉取配置绑定的server或location

### 2. `DescribeCustomizedConfigList`
- **详细描述**：拉取配置列表

### 3. `DescribeListeners`
- **详细描述**：查询负载均衡的监听器列表

### 4. `DescribeLoadBalancersDetail`
- **详细描述**：查询负载均衡详细信息

### 5. `DescribeTargetHealth`
- **详细描述**：获取负载均衡后端服务的健康检查状态

### 6. `DescribeTargetGroups`
- **详细描述**：查询目标组信息

### 7. `DescribeResources`
- **详细描述**：查询用户在当前地域支持可用区列表和资源列表

### 8. `DescribeTargetGroupList`
- **详细描述**：获取目标组列表

### 9. `DescribeTargetGroupInstances`
- **详细描述**：获取目标组绑定的服务器

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/clb)

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
4. 打包镜像: `docker build -t mcp-server-clb:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-clb:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-clb": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/214/30667

## 许可证

MIT