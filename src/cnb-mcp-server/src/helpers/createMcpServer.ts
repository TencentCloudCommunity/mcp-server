import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';

import { registerTools } from '../tools/index.js';

import type { Request } from 'express';

import { version as packageVersion } from '../../package.json';

export function createMcpServer(req?: Request) {
  const mcpServer = new McpServer({
    name: 'cnb-mcp-server',
    version: packageVersion
  });

  const token = req?.headers['authorization']?.split(' ')[1];
  registerTools(mcpServer, token);

  return mcpServer;
}
