# 文字识别 OCR-通用文字识别（高精度版）
> 腾讯云文字识别（Optical Character Recognition，OCR ）基于行业前沿的深度学习技术，将图片上的文字内容智能识别成为可编辑的文本。支持通用文字、卡证文字、票据单据、特定场景等多场景下的印刷体、手写体文字识别，同时支持票据核验功能，支持提供定制化服务，可以有效地代替人工录入信息。
>
> 通用文字识别（General Optical Character Recognition，General OCR）基于行业前沿的深度学习技术，提供通用印刷体识别、通用印刷体识别（高精度版）、通用手写体识别、英文识别、表格识别等多种服务，支持将图片上的文字内容，智能识别为可编辑的文本，可应用于随手拍扫描、纸质文档电子化、电商广告审核等多种场景，大幅提升信息处理效率。

---

## 产品功能

###  `通用印刷体识别（高精度版）`

- 支持图像整体文字的检测和识别，返回文字框位置与文字内容。相比通用印刷体识别接口，准确率和召回率更高，覆盖场景更广泛，应用场景包括：印刷文字识别、网络图片文字识别、广告图文字识别、街景店招文字识别、菜单文字识别、视频标题文字识别、头像文字识别等。

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/document/product/866/37490)

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
4. 打包镜像: `docker build -t mcp-server-ocr-document-high:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-ocr-document-high:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-ocr-document-high": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```



## API使用参考

>https://cloud.tencent.com/document/product/866/34937



## 许可证

MIT

