# 腾讯云 BI
> 腾讯云BI（Business Intelligence）提供从数据源接入、数据建模到数据可视化分析全流程的BI能力，仅需简单拖拽即可完成复杂的报表开发，并支持报表分享、推送等企业协作场景。其中的智能助手ChatBI作为基于大模型的智能分析Agent，支持通过简单对话实现数据分析，并提供数据解读、波动归因、业务优化建议等能力。腾讯云BI 简报模块具备强大的可视化能力，支持搭建大屏、领导驾驶舱、数据报告等，满足企业对外展示宣传、高层汇报、专题报告等业务场景。。

---

## Tools

### 1. `CreateProject`
- **详细描述**：创建项目

### 2. `DeleteProject`
- **详细描述**：删除项目

### 3. `DescribeProjectInfo`
- **详细描述**：项目详情

### 4. `DescribeProjectList`
- **详细描述**：项目列表数据接口

### 5. `ModifyProject`
- **详细描述**：修改项目

### 6. `CreateDatasource`
- **详细描述**：创建数据源

### 7. `CreateDatasourceCloud`
- **详细描述**：创建云数据库

### 8. `DeleteDatasource`
- **详细描述**：删除数据源

### 9. `ModifyDatasource`
- **详细描述**：更新数据源

### 10. `ModifyDatasourceCloud`
- **详细描述**：更新云数据库

### 11. `DescribePageWidgetList`
- **详细描述**：查询页面组件信息

### 12. `ExportScreenPage`
- **详细描述**：分享页截图导出

### 13. `DescribeDatasourceList`
- **详细描述**：查询数据源列表


---

## 产品链接
[点击查看产品详情](https://cloud.tencent.com/product/bi)

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
4. 打包镜像: `docker build -t mcp-server-bi:latest .`
5. 运行镜像: `docker run -it --env-file .env mcp-server-bi:latest`
6. 将配置填入 MCP 客户端中
```json
{
 "mcpServers": {
  "mcp-server-bi": {
   "type": "sse",
   "url": "http://127.0.0.1:9000/sse"
  }
 }
}
```

## API使用参考

>https://cloud.tencent.com/document/product/590/73735



## 许可证

MIT
