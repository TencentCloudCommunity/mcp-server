#!/usr/bin/env node

import express from 'express';
import { randomUUID } from 'node:crypto';
import { StreamableHTTPServerTransport } from '@modelcontextprotocol/sdk/server/streamableHttp.js';
import { SSEServerTransport } from '@modelcontextprotocol/sdk/server/sse.js';
import { isInitializeRequest } from '@modelcontextprotocol/sdk/types.js';
import dotenv from 'dotenv';

import { createMcpServer } from './helpers/createMcpServer.js';
import { stopWithWrongTransport } from './helpers/sendResponse.js';

dotenv.config();

const DEFAULT_APP_PORT = 3000;

// Store transports for each session type
const transports = {
  streamable: {} as Record<string, StreamableHTTPServerTransport>,
  sse: {} as Record<string, SSEServerTransport>
};

const app = express();

app.use(express.json());

app.post('/mcp', async (req, res) => {
  const sessionId = req.headers['mcp-session-id'] as string | undefined;
  let transport: StreamableHTTPServerTransport;

  // Reuse existing transport
  if (sessionId && transports.streamable[sessionId]) {
    transport = transports.streamable[sessionId];
    if (!(transport instanceof StreamableHTTPServerTransport)) {
      stopWithWrongTransport(res);
      return;
    }
    await transport.handleRequest(req, res, req.body);
    return;
  }

  // New initialization request
  if (!sessionId && isInitializeRequest(req.body)) {
    transport = new StreamableHTTPServerTransport({
      sessionIdGenerator: () => randomUUID(),
      onsessioninitialized: (sessionId) => {
        transports.streamable[sessionId] = transport;
      }
    });

    // Clean up transport when closed
    transport.onclose = () => {
      if (transport.sessionId) {
        delete transports.streamable[transport.sessionId];
      }
    };

    const mcpServer = createMcpServer(req);

    await mcpServer.connect(transport);
    await transport.handleRequest(req, res, req.body);
    return;
  }

  // Invalid request
  res.status(400).json({
    jsonrpc: '2.0',
    error: {
      code: -32000,
      message: 'Bad Request: No valid session ID provided'
    },
    id: null
  });
});

const handleSessionRequest = async (req: express.Request, res: express.Response) => {
  const sessionId = req.headers['mcp-session-id'] as string | undefined;
  if (!sessionId || !transports.streamable[sessionId]) {
    res.status(400).send('Invalid or missing session ID');
    return;
  }

  const transport = transports.streamable[sessionId];
  await transport.handleRequest(req, res, req.body);
};

app.get('/mcp', handleSessionRequest);

app.delete('/mcp', handleSessionRequest);

app.get('/sse', async (req, res) => {
  const transport = new SSEServerTransport('/messages', res);
  transports.sse[transport.sessionId] = transport;

  res.on('close', () => {
    delete transports.sse[transport.sessionId];
  });

  const mcpServer = createMcpServer(req);

  await mcpServer.connect(transport);
});

app.post('/messages', async (req, res) => {
  const sessionId = req.query.sessionId as string;
  const transport = transports.sse[sessionId];

  if (!transport) {
    res.status(400).send('No transport found for sessionId');
    return;
  }

  if (!(transport instanceof SSEServerTransport)) {
    stopWithWrongTransport(res);
    return;
  }

  await transport.handlePostMessage(req, res, req.body);
});

let port = parseInt(process.env.APP_PORT ?? '', 10);
if (isNaN(port)) {
  port = DEFAULT_APP_PORT;
}

const server = app.listen(port, () => {
  console.log(`MCP Streamable HTTP Server listening on port ${port}`);
});

process.on('SIGTERM', () => {
  console.log('SIGTERM signal received: closing HTTP server');
  server.close(() => {
    console.log('HTTP server closed');
  });
});
