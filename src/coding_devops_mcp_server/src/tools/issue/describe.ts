
import { CodingDevOpsConfig } from '../../config/environment.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { McpError, ErrorCode } from '../../errors.js';

export async function describeIssue(args: {
  projectName: string;
  issueCode: number;
}, config: CodingDevOpsConfig) {
  if (!args.projectName) {
    throw new McpError(ErrorCode.InvalidParams, 'projectName is required');
  }
  if (!args.issueCode) {
    throw new McpError(ErrorCode.InvalidParams, 'issueCode is required');
  }

  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  const issue = await connection.describeIssue({
    projectName: args.projectName,
    issueCode: args.issueCode
  });

  return {
    content: [
      {
        type: 'text',
        text: `Issue Details:`,
      },
      {
        type: 'text',
        text: `Code: ${issue.Code}`,
      },
      {
        type: 'text',
        text: `Name: ${issue.Name}`,
      },
      {
        type: 'text',
        text: `Type: ${issue.Type}`,
      },
      {
        type: 'text',
        text: `Status: ${issue.IssueStatusName}`,
      },
      {
        type: 'text',
        text: `Priority: ${issue.Priority}`,
      },
      {
        type: 'text',
        text: `Description: ${issue.Description || 'N/A'}`,
      },
      {
        type: 'text',
        text: `Assignee: ${issue.Assignee?.Name || 'Unassigned'}`,
      },
      {
        type: 'text',
        text: `Created At: ${new Date(issue.CreatedAt).toLocaleString()}`,
      }
    ],
  };
}