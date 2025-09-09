# 云直播 CSS
> 云直播（Cloud Streaming Services，CSS）为您提供极速、稳定、专业的云端直播处理服务，根据业务的不同直播场景需求，云直播提供了标准直播、快直播、云导播台三种服务，分别针对大规模实时观看、超低延时直播、便捷云端导播的场景，配合腾讯云视立方·直播 SDK，为您提供一站式的音视频直播解决方案。

---

## Tools

### 1. `DescribeLiveDomains`
- **详细描述**：查询域名列表

### 2. `DescribeLiveDomain`
- **详细描述**：查询域名信息

### 3. `DescribeLiveStreamOnlineList`
- **详细描述**：查询直播中的流

### 4. `DescribePushBandwidthAndFluxList`
- **详细描述**：直播推流带宽和流量数据查询

### 5. `DescribeBillBandwidthAndFluxList`
- **详细描述**：直播播放带宽和流量数据查询

### 6. `DescribeStreamPlayInfoList`
- **详细描述**：查询流的播放信息列表

### 7. `DescribeStreamPushInfoList`
- **详细描述**：查询某条流上行推流质量数据

### 8. `DescribeLiveStreamEventList`
- **详细描述**：查询推断流事件

### 9. `DescribeLiveTranscodeDetailInfo`
- **详细描述**：查询直播转码统计信息

### 10. `DescribeVisitTopSumInfoList`
- **详细描述**：查询某时间段top n的域名或流id信息

---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/css)

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
4. 打包镜像: `docker build -t mcp-server-live:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-live:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-live": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```


## API使用参考

>https://cloud.tencent.com/document/product/267/20456



## 许可证
MIT