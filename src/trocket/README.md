# 消息队列 RocketMQ 版
> 消息队列 RocketMQ 版(TDMQ for RocketMQ，简称TDMQ RocketMQ 版) 是一款分布式高可用的消息队列服务，基于 Apache RocketMQ 的 4.x 和 5.x 架构提供不同的产品形态，支持开源客户端零改造接入，同时具备计算存储分离，灵活扩缩容的优势。TDMQ RocketMQ 版可以支持百万级 TPS 的吞吐量，适用于各类大规模、低延时、对可靠性要求高的在线消息业务场景。

---

## Tools

### 1. `DescribeFusionInstanceList`
- **详细描述**：查询集群列表

### 2. `DescribeInstance`
- **详细描述**：查询集群详情

### 3. `DescribeTopicList`
- **详细描述**：查询主题列表

### 4. `DescribeTopic`
- **详细描述**：查询主题详情

### 5. `DescribeTopicListByGroup`
- **详细描述**：查询消费组订阅的主题列表

### 6. `DescribeProducerList`
- **详细描述**：查询生产者信息列表

### 7. `DescribeConsumerGroup`
- **详细描述**：查询消费组详情

### 8. `DescribeConsumerGroupList`
- **详细描述**：查询消费组列表

### 9. `DescribeConsumerClient`
- **详细描述**：查询消费者客户端详情

### 10. `DescribeConsumerClientList`
- **详细描述**：查询消费组下的客户端连接列表

### 11. `DescribeRoleList`
- **详细描述**：查询角色列表

### 12. `DescribeMessageList`
- **详细描述**：查询消息列表

### 13. `DescribeMessage`
- **详细描述**：查询消息详情

### 14. `DescribeMessageTrace`
- **详细描述**：查询消息轨迹

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/trocket)

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
4. 打包镜像: `docker build -t mcp-server-trocket:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-trocket:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-trocket": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/1493/96041



## 许可证

MIT



TENCENTCLOUD_SECRET_ID

TENCENTCLOUD_SECRET_KEY
