# Redis 云数据库
> 腾讯云数据库 Redis® 是腾讯云打造的兼容 Redis 协议的缓存和存储服务。丰富的数据结构能帮助您完成不同类型的业务场景开发。支持主从热备，提供自动容灾切换、数据备份、故障迁移、实例监控、在线扩容、数据回档等全套的数据库服务。

---

## Tools

### 1. `DescribeProductInfo`
- **详细描述**：查询产品售卖规格

### 2. `DescribeInstanceZoneInfo`
- **详细描述**：查询Redis节点详细信息

### 3. `DescribeInstanceSecurityGroup`
- **详细描述**：查询实例安全组信息

### 4. `DescribeProjectSecurityGroup`
- **详细描述**：查询项目安全组信息

### 5. `DescribeInstances`
- **详细描述**：查询Redis实例列表

### 6. `DescribeSlowLog`
- **详细描述**：查询实例慢查询记录

### 7. `DescribeInstanceMonitorHotKey`
- **详细描述**：查询实例热Key

### 8. `DescribeInstanceBackups`
- **详细描述**：查询Redis实例备份列表

### 9. `DescribeInstanceParams`
- **详细描述**：查询实例的参数列表

### 10. `DescribeParamTemplates`
- **详细描述**：查询参数模板列表

### 11. `DescribeParamTemplateInfo`
- **详细描述**：查询参数模板详情

### 12. `DescribeInstanceAccount`
- **详细描述**：查看实例子账号信息

### 13. `CreateInstances`
- **详细描述**：创建Redis实例

### 14. `ModifyInstanceParams`
- **详细描述**：修改实例参数

### 15. `CreateInstanceAccount`
- **详细描述**：创建实例子账号

### 16. `AssociateSecurityGroups`
- **详细描述**：绑定安全组

### 17. `DisassociateSecurityGroups`
- **详细描述**：安全组批量解绑云资源

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/crs)

---

## 使用

### 参数获取方式：

a、密钥获取：腾讯云控制台当中，生成腾讯云SecretId和SecretKey
> 生成地址：https://console.cloud.tencent.com/cam/capi

b、地域列表的映射关系：
如广州地域，则地域字段的内容应该为：ap-guangzhou
> https://cloud.tencent.com/document/product/239/20005#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8

### 部署与配置
1. 生成配置文件
   `cp .env.example .env`
2. 将密钥填写到 `.env` 配置文件中
3. 安装 [docker](https://www.docker.com/)
4. 打包镜像: `docker build -t mcp-server-redis:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-redis:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-redis": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/239/20002

## 许可证
MIT