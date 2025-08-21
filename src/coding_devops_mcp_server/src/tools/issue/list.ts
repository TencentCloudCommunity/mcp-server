
import { ErrorCode, McpError } from '@modelcontextprotocol/sdk/types.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { CodingDevOpsConfig } from '../../config/environment.js';

export async function listIssues(args: { 
  projectName: string;
  issueType?: string;
  limit?: string;
}, config: CodingDevOpsConfig) {
  if (!args.projectName) {
    throw new McpError(ErrorCode.InvalidParams, 'projectName03 is required');
  }

  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  const issues = await connection.listIssues({
    projectName: args.projectName,
    issueType: args.issueType,
    limit: args.limit
  });

  return {
    content: [
      {
        type: 'text',
        text: JSON.stringify(issues, null, 2),
      },
    ],
  };
}