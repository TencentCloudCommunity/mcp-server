
import { ErrorCode, McpError } from '@modelcontextprotocol/sdk/types.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { CodingDevOpsConfig } from '../../config/environment.js';

export async function listDepots(args: { 
  projectId: string;
  PageNumber?: string;
  PageSize?: string;
}, config: CodingDevOpsConfig) {
  if (!args.projectId) {
    throw new McpError(ErrorCode.InvalidParams, 'projectId is required');
  }

  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  const issues = await connection.listProjectDepots({
    ProjectId: args.projectId,
    PageNumber: args.PageNumber,
    PageSize: args.PageSize
  });

  return {
    success: true,
    data: {
      content: [
        {
          type: 'table',
          headers: ['仓库名称', '描述', 'HTTPS地址', 'SSH地址', '默认分支'],
          rows: issues.Depots.map(depot => [
            depot.Name,
            depot.Description || '-',
            depot.HttpsUrl,
            depot.SshUrl,
            depot.DefaultBranch
          ])
        },
        {
          type: 'text',
          text: `总共 ${issues.Page.TotalRow} 个仓库，当前第 ${issues.Page.PageNumber}/${issues.Page.TotalPage} 页`
        }
      ],
      metadata: {
        page: {
          current: issues.Page.PageNumber,
          total: issues.Page.TotalPage,
          size: issues.Page.PageSize,
          totalItems: issues.Page.TotalRow
        }
      }
    }
  };
}