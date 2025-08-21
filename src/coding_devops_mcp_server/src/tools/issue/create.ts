import { CodingDevOpsConfig } from '../../config/environment.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { McpError, ErrorCode } from '../../errors.js';

export async function createIssue(args: {
  projectName: string;
  name: string;
  type: string;
  priority: string;
  description: string;
  parentCode?: number;
}, config: CodingDevOpsConfig) {
  if (!args.projectName) {
    throw new McpError(ErrorCode.InvalidParams, 'projectName01 is required');
  }
  if (!args.name) {
    throw new McpError(ErrorCode.InvalidParams, 'Name is required');
  }
  if (!args.type) {
    throw new McpError(ErrorCode.InvalidParams, 'Type is required');
  }
  if (!args.priority) {
    throw new McpError(ErrorCode.InvalidParams, 'Priority is required');
  }
  if (!args.description) {
    throw new McpError(ErrorCode.InvalidParams, 'Description is required');
  }

  // 验证事项类型
  const validTypes = ['DEFECT', 'REQUIREMENT', 'MISSION', 'EPIC'];
  if (!validTypes.includes(args.type)) {
    throw new McpError(
      ErrorCode.InvalidParams,
      `Invalid type. Must be one of: ${validTypes.join(', ')}`
    );
  }

  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  const issue = await connection.createIssue({
    projectName: args.projectName,
    name: args.name,
    type: args.type,
    priority: args.priority,
    description: args.description,
    parentCode: args.parentCode
  });

  return {
    content: [
      {
        type: 'text',
        text: `Successfully created issue: ${issue.Name} (${issue.Code})`,
      },
    ],
  };
}