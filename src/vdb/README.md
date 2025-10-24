# 向量数据库
> 腾讯云向量数据库（Tencent Cloud VectorDB）是一款全托管的自研企业级分布式数据库服务，专用于存储、检索、分析多维向量数据。该数据库支持多种索引类型和相似度计算方法，单索引支持千亿级向量规模，可支持百万级 QPS 及毫秒级查询延迟。腾讯云向量数据库不仅能为大模型提供外部知识库，提高大模型回答的准确性，还可广泛应用于推荐系统、自然语言处理等 AI 领域。

---

## Tools

### 1. `DescribeInstances`
- **详细描述**：查询实例列表

### 2. `DescribeDBSecurityGroups`
- **详细描述**：查询实例安全组详情

### 3. `DescribeInstanceMaintenanceWindow`
- **详细描述**：查询维护时间窗

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/vdb)

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
4. 打包镜像: `docker build -t mcp-server-vdb:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-vdb:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-vdb": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/product/1709/106767



## 许可证

MIT