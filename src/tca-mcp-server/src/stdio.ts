#!/usr/bin/env node

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { registerTools } from "./tools/index.js";

// Create an MCP server
const server = new McpServer({
  name: "腾讯云代码分析(TCA)",
  version: "1.0.0"
});

registerTools(server);
const transport = new StdioServerTransport();
await server.connect(transport);
console.log("Server started");
