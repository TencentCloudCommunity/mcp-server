# <img src="https://tdesign.tencent.com/favicon.ico" alt="TDesign" height="45"/> TDesign MCP Server

## 🎉 功能介绍

### 🔨 支持工具
- [x] `get-component-changelog`：获取组件的变更日志，适用场景：组件库版本升级
- [x] `get-component-docs`：获取组件的文档，适用场景：代码生成和代码转换
- [x] `get-component-dom`：获取组件的 DOM 结构，适用场景：转换自定义的 CSS 样式
- [x] `get-component-list`：获取所有可用的组件列表，适用场景：选择合适组件进行组合实现不存在的功能

### 🔨 支持框架
`React` / `Vue Next` / `Vue`

## 📦 使用说明
```json
{
  "mcpServers": { // 或 servers（根据不同的 MCP 客户端决定）
    "tdesign-mcp-server": {
      "command": "npx",
      "args": ["-y", "tdesign-mcp-server@latest"]
    }
  }
}
```

> 💡 如果在 MCP 客户端中连接失败，可以在终端手动输入 `npx tdesign-mcp-server`，并反馈对应的报错信息在 issue。

## 🔮 开发指南

```bash
📄 .env # 环境变量（仅用于本地调试）
📁 packages
├── 📁 common # 通用工具函数
├── 📁 docs # 已整理的文档
├── 📁 scripts # 自动化脚本
└── 📁 server # 核心业务逻辑
    ├── 📁 tools
    └── 📁 prompts
```

- 由于该项目在提取文档时，如果找不到本地仓库的话，会自动启动 clone 其它组件库的流程。
- 如果你本地已经有 TDesign 的其它组件库项目，可以把该仓库 clone 到相同的目录下；如果没有，建议新建一个额外的空目录，再把该仓库 clone 进去。
