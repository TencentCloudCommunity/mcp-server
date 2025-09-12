# 云数据库 TencentDB for MongoDB
> 腾讯云数据库 MongoDB（TencentDB for MongoDB）是腾讯云基于全球广受欢迎的 MongoDB 打造的高性能 NoSQL 数据库，100%完全兼容 MongoDB 协议，支持跨文档事务，提供稳定丰富的监控管理，弹性可扩展、自动容灾，适用于文档型数据库场景，您无需自建灾备体系及控制管理系统。

---


## Tools

### 1. `DescribeDBInstances`
- **详细描述**：查询云数据库实例列表

### 2. `DescribeDBBackups`
- **详细描述**：查询实例备份列表

### 3. `DescribeInstanceParams`
- **详细描述**：获取当前实例可修改的参数列表

### 4. `DescribeSlowLogs`
- **详细描述**：获取慢日志信息

### 5. `DescribeSpecInfo`
- **详细描述**：查询云数据库的售卖规格（包含可用区信息）

### 6. `DescribeSecurityGroup`
- **详细描述**：查询实例绑定的安全组



---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/mongodb)

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
4. 打包镜像: `docker build -t mcp-server-mongodb:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-mongodb:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-mongodb": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```



## API使用参考

>https://cloud.tencent.com/document/product/240/38554



## 许可证

MIT
