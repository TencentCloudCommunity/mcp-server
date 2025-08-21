# CNB MCP Server

CNB(https://cnb.cool) toolkits for LLMs supporting the MCP protocol

## How to use

### STDIO

```json
{
  "mcpServers": {
    "cnb": {
      "command": "npx",
      "args": ["-y", "-p", "@cnbcool/mcp-server", "cnb-mcp-stdio"],
      "env": {
        "API_BASE_URL": "<BASE_URL>", // optional, defualt vaule: https://api.cnb.cool
        "API_TOKEN": "<YOUR_TOKEN>"
      }
    }
  }
}
```


## Prerequisite

1. node >= 18

## How to develop

1. `npm install`
2. `npx openapi-typescript@5.4.2 https://api.cnb.cool/swagger.json -o src/schema.d.ts`
3. Copy `.env.example` and rename to `.env` and fill in the values
4. `npm run build`
5. `npx @modelcontextprotocol/inspector node dist/stdio.js`

> @modelcontextprotocol/inspector requires Node.js: ^22.7.5

> https://github.com/modelcontextprotocol/inspector?tab=readme-ov-file#requirements

## How to preview

1. npm run build
2. set mcpServers config:

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


