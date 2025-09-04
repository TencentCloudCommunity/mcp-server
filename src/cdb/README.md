# 云数据库 TencentDB for MySQL
> 腾讯云数据库 MySQL（TencentDB for MySQL）为用户提供安全可靠，性能卓越、易于维护的企业级云数据库服务。

---

## Tools

### 1. `DescribeDBInstances`
- **详细描述**：查询实例列表

### 2. `DescribeDBInstanceInfo`
- **详细描述**：查询实例基本信息

### 3. `DescribeInstanceUpgradeType`
- **详细描述**：查询数据库实例升级类型

### 4. `DescribeAccounts`
- **详细描述**：查询云数据库的所有账号信息

### 5. `DescribeDatabases`
- **详细描述**：查询数据库

### 6. `DescribeInstanceParams`
- **详细描述**：查询实例的可设置参数列表

### 7. `DescribeParamTemplates`
- **详细描述**：查询参数模板列表

### 8. `DescribeParamTemplateInfo`
- **详细描述**：查询参数模板详情

### 9. `CreateDatabase`
- **详细描述**：创建数据库



---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/cdb)

---

### 使用

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
4. 打包镜像: `docker build -t mcp-server-cdb:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-cdb:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-cdb": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/236/15830


## 许可证

MIT

