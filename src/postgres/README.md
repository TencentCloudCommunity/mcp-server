# 云数据库 TencentDB for PostgreSQL
> 腾讯云数据库 PostgreSQL（TencentDB for PostgreSQL，云 API 使用 postgres 作为简称）能够让您在云端轻松设置、操作和扩展目前功能最强大的开源数据库 PostgreSQL。腾讯云将负责绝大部分处理复杂而耗时的管理工作，如 PostgreSQL 软件安装、存储管理、高可用复制、以及为灾难恢复而进行的数据备份，让您更专注于业务程序开发。

---

## Tools

### 1. `DescribeDBInstances`
- **详细描述**：查询实例列表

### 2. `DescribeDBInstanceAttribute`
- **详细描述**：查询实例详情

### 3. `UpgradeDBInstanceKernelVersion`
- **详细描述**：升级实例内核版本号

### 4. `DescribeAccounts`
- **详细描述**：查询实例的数据库账号列表

### 5. `DescribeDatabases`
- **详细描述**：查询实例的数据库列表

### 6. `DescribeDBInstanceParameters`
- **详细描述**：查询实例参数

### 7. `DescribeParameterTemplates`
- **详细描述**：查询参数模板列表

### 8. `DescribeParameterTemplateAttributes`
- **详细描述**：查询参数模板详情

### 9. `CreateDatabase`
- **详细描述**：创建数据库



---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/postgres)

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
4. 打包镜像: `docker build -t mcp-server-postgres:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-postgres:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-postgres": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/409/16761



## 许可证

MIT
