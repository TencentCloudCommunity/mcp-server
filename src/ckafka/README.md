# 消息队列 CKafka 版
> 消息队列 CKafka 版（TDMQ for CKafka）是一个分布式、高吞吐量、高可扩展性的消息系统，100%兼容开源 Kafka API 2.4、2.8、3.2 版本。CKafka 基于发布/订阅模式，通过消息解耦，使生产者和消费者异步交互，无需彼此等待。CKafka 具有高可用、数据压缩、同时支持离线和实时数据处理等优点，适用于日志压缩收集、监控数据聚合、流式数据集成等场景。

---

## Tools

### 1. `DescribeTopic`
- **详细描述**：获取主题列表

### 2. `DescribeTopicAttributes`
- **详细描述**：获取主题属性

### 3. `DescribeTopicDetail`
- **详细描述**：获取主题列表详情

### 4. `DescribeTopicFlowRanking`
- **详细描述**：Topic流量排行

### 5. `DescribeTopicProduceConnection`
- **详细描述**：查询topic生产端连接信息

### 6. `DescribeTopicSubscribeGroup`
- **详细描述**：查询订阅某主题消息分组信息

### 7. `FetchMessageByOffset`
- **详细描述**：查询消息

### 8. `FetchMessageListByOffset`
- **详细描述**：根据位点查询消息列表

### 9. `FetchMessageListByTimestamp`
- **详细描述**：根据时间戳查询消息列表

### 10. `DescribeGroupInfo`
- **详细描述**：获取消费分组信息

### 11. `DescribeGroupOffsets`
- **详细描述**：获取消费分组offset

### 12. `DescribeConsumerGroup`
- **详细描述**：查询消费分组信息

### 13. `DescribeGroup`
- **详细描述**：枚举消费分组(精简版)

### 14. `DescribeACL`
- **详细描述**：枚举ACL

### 15. `DescribeAclRule`
- **详细描述**：查询ACL规则列表

### 16. `DescribeInstanceAttributes`
- **详细描述**：获取实例属性

### 17. `DescribeInstances`
- **详细描述**：获取实例列表信息

### 18. `DescribeInstancesDetail`
- **详细描述**：获取实例集群列表详情

### 19. `InquireCkafkaPrice`
- **详细描述**：Ckafka询价

### 20. `DescribeCkafkaZone`
- **详细描述**：查看可用区列表

### 21. `DescribeDatahubTask`
- **详细描述**：查询Datahub任务信息

### 22. `DescribeDatahubTasks`
- **详细描述**：查询Datahub任务列表




---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/ckafka)

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
4. 打包镜像: `docker build -t mcp-server-ckafka:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-ckafka:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-ckafka": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/product/597/40823



## 许可证

MIT
