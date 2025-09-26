# 腾讯混元生图
> 腾讯混元生图是一款提供 AI 图像生成与处理能力的 API 技术服务，可以结合输入的文本或图像智能创作图像内容，具有更精美的绘图品质、更强大的中文理解能力、更多样的风格选择与更广泛的应用场景，为高质量的图像内容创作、内容运营提供技术支持。

---

## Tools

### 1. `TextToImageLite`
- **详细描述**：混元生图（极速版）

### 2. `TextToImageRapid`
- **详细描述**：混元生图（2.0）

### 3. `RefineImage`
- **详细描述**：图片变清晰

### 4. `ImageToImage`
- **详细描述**：图像风格化（图生图）

### 5. `ImageOutpainting`
- **详细描述**：扩图

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/document/product/1668)

---

## 使用


### 参数获取方式：

a、密钥获取：腾讯云控制台当中，生成腾讯云SecretId和SecretKey
> 生成地址：https://console.cloud.tencent.com/cam/capi

b、地域列表的映射关系：
如广州地域，则地域字段的内容应该为：ap-guangzhou
> https://cloud.tencent.com/document/product/1596/77930

c、返回的URL当中，可能部分字符被Unicode方式转义，使用者需处理后使用。

### 部署与配置
1. 生成配置文件
   `cp .env.example .env`
2. 将密钥填写到 `.env` 配置文件中
3. 安装 [docker](https://www.docker.com/)
4. 打包镜像: `docker build -t mcp-server-aiart:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-aiart:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-aiart": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/1668/88077



## 许可证

MIT

