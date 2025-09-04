# 内容分发网络 CDN
> 内容分发网络（Content Delivery Network，CDN）通过将站点内容发布至遍布全球的海量加速节点，使其用户可就近获取所需内容，避免因网络拥堵、跨运营商、跨地域、跨境等因素带来的网络不稳定、访问延迟高等问题，有效提升下载速度、降低响应时间，提供流畅的用户体验。

---

## Tools

### 1. `DescribeOriginData`
- **详细描述**：回源数据查询

### 2. `ListTopData`
- **详细描述**：TOP 数据查询

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/cdn)

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
4. 打包镜像: `docker build -t mcp-server-cdn:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-cdn:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-cdn": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/228/30974



## 许可证

MIT

