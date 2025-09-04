# SSL 证书
> 腾讯云 SSL 证书（SSL Certificates）为您提供 SSL 证书的申请、管理、部署等服务，为您提供一站式 HTTPS 解决方案。

---

## Tools

### 1. `CreateCertificate`
- **详细描述**：购买付费证书

### 2. `CommitCertificateInformation`
- **详细描述**：付费证书提交证书订单

### 3. `DescribeCertificateDetail`
- **详细描述**：获取证书详情

### 4. `UploadCertificate`
- **详细描述**：上传证书

### 5. `DescribeCertificates`
- **详细描述**：获取证书列表

### 6. `DescribeManagerDetail`
- **详细描述**：查询管理人详情

### 7. `DescribeManagers`
- **详细描述**：查询管理人列表


---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/ssl)

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
4. 打包镜像: `docker build -t mcp-server-ssl:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-ssl:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-ssl": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/product/400/41681



## 许可证
MIT