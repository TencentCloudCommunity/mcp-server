import { CodingDevOpsConfig } from '../../config/environment.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { McpError, ErrorCode } from '../../errors.js';

export async function deleteIssue(args: { 
  projectName: string;
  issueCode: number;
}, config: CodingDevOpsConfig) {
  if (!args.projectName) {
    throw new McpError(ErrorCode.InvalidParams, 'projectName02 is required');
  }
  if (!args.issueCode) {
    throw new McpError(ErrorCode.InvalidParams, 'Issue code is required');
  }

  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  await connection.deleteIssue({
    projectName: args.projectName,
    issueCode: args.issueCode
  });

  return {
    content: [
      {
        type: 'text',
        text: `Successfully deleted issue ${args.issueCode} in project ${args.projectName}`,
      },
    ],
  };
}