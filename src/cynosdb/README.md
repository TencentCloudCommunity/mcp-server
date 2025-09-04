# TDSQL-C MySQL 版
> TDSQL-C MySQL 版（TDSQL-C for MySQL）是腾讯云自研的新一代云原生关系型数据库。

---

## Tools

### 1. `DescribeInstances`
- **详细描述**：用于获取用户数据库实例列表，包括实例ID和基本信息

### 2. `DescribeInstanceDetail`
- **详细描述**：用于获取指定实例的详细配置和状态信息

### 3. `DescribeClusterDatabases`
- **详细描述**：用于查询指定集群中的数据库列表，包括名称和权限（腾讯云的权限信息可能需要结合 DescribeAccountPrivileges 获取）

### 4. `DescribeAccounts`
- **详细描述**：用于列出指定集群或实例的数据库账号信息，包括权限详情（需配合 DescribeAccountPrivileges 查看完整权限）

### 5. `ModifyInstanceName`
- **详细描述**：用于更新指定实例的别名或名称

### 6. `AssociateSecurityGroups`
- **详细描述**：用于将已有安全组绑定到数据库实例，实现网络访问控制

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/cynosdb)

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
4. 打包镜像: `docker build -t mcp-server-cynosdb:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-cynosdb:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-cynosdb": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/product/1003/48106



## 许可证

MIT

