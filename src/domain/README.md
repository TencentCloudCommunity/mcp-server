# 域名注册
> 在 Internet 上有千百万台主机，为了区分这些主机，每一台主机都被分配一个 IP 地址。但由于 IP 地址没有实际意义难于记忆，于是就有了域名（Domain Name）。而获得域名的方式是通过付费获得域名一年或多年内的使用权，我们称之为域名注册。

---

## Tools

### 1. `DescribeDomainPriceList`
- **详细描述**：域名价格列表

### 2. `CheckDomain`
- **详细描述**：域名注册查询

### 3. `DescribeDomainBaseInfo`
- **详细描述**：域名基本信息

### 4. `DescribeTemplate`
- **详细描述**：获取模板信息

### 5. `DescribeDomainNameList`
- **详细描述**：我的域名列表

### 6. `DescribeTemplateList`
- **详细描述**：信息模板列表

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/domain)

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
4. 打包镜像: `docker build -t mcp-server-domain:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-domain:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-domain": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/product/242/38884



## 许可证

MIT
