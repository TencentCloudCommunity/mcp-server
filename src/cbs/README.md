# 云硬盘 CBS
> 云硬盘（Cloud Block Storage，CBS）为您提供用于 CVM 的持久性数据块级存储服务。云硬盘中的数据自动地在可用区内以多副本冗余方式存储，避免数据的单点故障风险，提供高达99.9999999%的数据可靠性。同时提供多种类型及规格，满足稳定低延迟的存储性能要求。

---

## Tools

### 1. `DescribeDisks`
- **详细描述**：查询云硬盘列表

### 2. `DescribeSnapshots`
- **详细描述**：查询快照列表

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/cbs)

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
4. 打包镜像: `docker build -t mcp-server-cbs:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-cbs:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-cbs": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/362/15634



## 许可证

MIT