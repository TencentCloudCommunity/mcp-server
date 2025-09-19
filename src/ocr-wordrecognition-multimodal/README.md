# 文字识别 OCR-文档抽取（多模态版）
> 腾讯云文字识别（Optical Character Recognition，OCR ）基于行业前沿的深度学习技术，将图片上的文字内容智能识别成为可编辑的文本。支持通用文字、卡证文字、票据单据、特定场景等多场景下的印刷体、手写体文字识别，同时支持票据核验功能，支持提供定制化服务，可以有效地代替人工录入信息。
>
> 文档智能（Document AI）​​ 深度融合 OCR 与多模态大模型，实现高精度识别、智能解析与结构化信息抽取。覆盖货运单证、跨境物流、快递面单、教育作业、保险理赔及国际结算等场景，助力企业自动化升级，提升运营效率与数据准确性。

---

## 产品功能

###  `自定义键值`
- 支持自定义创建个性化键值（key），用户通过传入自定义 key，模型自动建立图片中文字的键值对应关系，实现对任意版式图片的结构化识别。

###  `智能匹配`
- 支持对已上传图片的智能配准，实现对不同版式图片与已发布模板的自动匹配，减少人工分类成本，快速实现图片的结构化识别。

###  `自定义字段类型`
- 支持自定义创建字段类型，支持针对不同识别区内容类型进行专项优化，如小写金额、日期、纯数字等，可根据需求选择合适的字段类型以提升识别准确率，也可通过穷举可能的输出值范围自定义字段类型，对识别结果进行智能纠正和规范。

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/document/product/866/37494)

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
4. 打包镜像: `docker build -t mcp-server-ocr-wordrecognition-multimodal:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-ocr-wordrecognition-multimodal:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-ocr-wordrecognition-multimodal": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```



## API使用参考

>https://cloud.tencent.com/document/product/866/119451



## 许可证

MIT

