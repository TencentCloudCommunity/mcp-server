import { CodingConnection } from '../../api/coding_connection.js';
import { CodingDevOpsConfig } from '../../config/environment.js';
import { McpError, ErrorCode } from '../../errors.js';
import { CreateProjectArgs } from './index.js';

export async function createProject(args: CreateProjectArgs, config: CodingDevOpsConfig) {
  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  const result = await connection.createProject({
    name: args.name,
    displayName: args.displayName.trim(),
    description: args.description?.trim(),
    projectTemplate: args.projectTemplate,
    shared: args.shared
  });

  return {
    content: [
      {
        type: 'text',
        text: `Successfully created project, project Id is: ${result.Id}`,
      },
    ],
  };
}
