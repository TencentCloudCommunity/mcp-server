# 腾讯云数据仓库 TCHouse-C
> 腾讯云数据仓库 TCHouse-C（ Tencent Cloud TCHouse-C ）是基于开源 OLAP 引擎 ClickHouse 打造的企业级云数据仓库服务，仅需几分钟即可快速搭建起 PB 级实时数据仓库，实现对海量数据的极速实时分析。TCHouse-C 内核与开源版 ClickHouse 高度兼容，大幅增强了产品稳定性、安全性和运维便捷性，使您无需关注底层基础设施、专注于数据价值的提升。

---

## Tools

### 1. `DescribeInstance`
- **详细描述**：描述实例信息

### 2. `DescribeInstancesNew`
- **详细描述**：获取实例简单信息列表

### 3. `DescribeInstanceNodes`
- **详细描述**：获取实例节点信息列表

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/tchousec)

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
4. 打包镜像: `docker build -t mcp-server-cdwch:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-cdwch:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-cdwch": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/1299/47756



## 许可证

MIT