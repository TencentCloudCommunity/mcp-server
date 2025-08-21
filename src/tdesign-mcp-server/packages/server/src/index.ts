import { promises as fs } from "fs";

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";

import { getDirname } from "@tdesign-mcp-server/common";
const __dirname = getDirname(import.meta.url);

import path from "path";
import registryPrompts from "./prompts";
import registryTools from "./tools";

const MCP_VERSION = await getPackageVersion();
const server = new McpServer({
  name: "TDesign MCP",
  version: MCP_VERSION
});

registryPrompts(server);
registryTools(server);

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.log(`\x1b[36m%s\x1b[0m`, `🚀 TDesign MCP Server running (v${MCP_VERSION})`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});

async function getPackageVersion() {
  const packageJsonPath = path.join(__dirname, "../package.json");
  const content = await fs.readFile(packageJsonPath, "utf-8");
  const { version } = JSON.parse(content);
  return version as string;
}
