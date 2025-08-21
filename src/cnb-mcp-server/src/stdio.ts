#!/usr/bin/env node

import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import dotenv from 'dotenv';

import { createMcpServer } from './helpers/createMcpServer.js';

dotenv.config();

const server = createMcpServer();

async function main() {
  console.error('server starting...');
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error('server connected');
}

main().catch((error) => {
  console.error('Fatal error:', error);
  process.exit(1);
});
