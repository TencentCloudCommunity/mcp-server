
import { CodingConnection } from '../../api/coding_connection.js';
import { CodingDevOpsConfig } from '../../config/environment.js';
import { McpError, ErrorCode } from '../../errors.js';
import { DeleteProjectArgs } from './index.js';

export async function deleteProject(args: DeleteProjectArgs, config: CodingDevOpsConfig) {
  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  

  await connection.deleteProject({
    ProjectId: args.projectId
  });

  return {
    content: [
      {
        type: 'text',
        text: `Successfully deleted project with ID: ${args.projectId}`,
      },
    ],
  };
}