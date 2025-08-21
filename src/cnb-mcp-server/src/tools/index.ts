import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';

import CnbApiClient from '../api/client.js';
import registerGroupTools from './groupTools.js';
import registerRepoTools from './repoTools.js';
import registerIssueTools from './issueTools.js';
import registerWorkspaceTools from './workspaceTools.js';
import registerPullTools from './pullTools.js';

export function registerTools(server: McpServer, token?: string) {
  CnbApiClient.initialize({
    baseUrl: process.env.API_BASE_URL || 'https://api.cnb.cool',
    token: process.env.API_TOKEN || token || ''
  });

  registerGroupTools(server);
  registerRepoTools(server);
  registerIssueTools(server);
  registerWorkspaceTools(server);
  registerPullTools(server);
}
