# CNB MCP Server

CNB(https://cnb.cool) 支持 MCP 协议的 MCP Server

## 使用方法

### STDIO

```json
{
  "mcpServers": {
    "cnb": {
      "command": "npx",
      "args": ["-y", "-p", "@cnbcool/mcp-server", "cnb-mcp-stdio"],
      "env": {
        "API_BASE_URL": "<BASE_URL>", // 可选，默认值: https://api.cnb.cool
        "API_TOKEN": "<YOUR_TOKEN>"
      }
    }
  }
}
```

## 环境要求

1. node >= 18

## 开发指南

1. `npm install`
2. `npx openapi-typescript@5.4.2 https://api.cnb.cool/swagger.json -o src/schema.d.ts`
3. 复制 `.env.example` 并重命名为 `.env` 并填写相应值
4. `npm run build`
5. `npx @modelcontextprotocol/inspector node dist/stdio.js`

> @modelcontextprotocol/inspector 需要 Node.js: ^22.7.5

> https://github.com/modelcontextprotocol/inspector?tab=readme-ov-file#requirements

## 预览方法

1. npm run build
2. 设置 mcpServers 配置:

```json
{
  "mcpServers": {
    "cnb": {
      "command": "node",
      "args": ["/path/to/cnbcool/mcp-server/dist/stdio.js"],
      "env": {
        "API_BASE_URL": "<BASE_URL>", // optional, defualt vaule: https://api.cnb.cool
        "API_TOKEN": "<YOUR_TOKEN>"
      }
    }
  }
}
```

### 使用CodeBuddy自动根据swagger.json补充接口MCP Tools

1. `wget https://api.cnb.cool/swagger.json -O ./swagger.json`
2. `npx split-swagger s ./swagger.json ./swagger` 
3. 参考prompt.txt中提示词,调整为有变动的json文件，然后让ai代码助手执行
