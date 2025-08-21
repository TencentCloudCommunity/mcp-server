import { ErrorCode, McpError } from '@modelcontextprotocol/sdk/types.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { CodingDevOpsConfig } from '../../config/environment.js';

export async function listCommits(args: {
  DepotId?: string;
  DepotPath?: string;
  Ref: string;
  EndDate?: string;
  StartDate?: string;
  KeyWord?: string;
  PageNumber?: string;
  PageSize?: string;
  Path?: string;
}, config: CodingDevOpsConfig) {
  if (!args.DepotId && !args.DepotPath) {
    throw new McpError(ErrorCode.InvalidParams, '需要提供 DepotId 或 DepotPath 其中之一');
  }

  if (!args.Ref) {
    throw new McpError(ErrorCode.InvalidParams, '分支名称(Ref)为必填参数');
  }

  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  const commits = await connection.describeGitCommits({
    DepotId: args.DepotId,
    DepotPath: args.DepotPath,
    Ref: args.Ref,
    EndDate: args.EndDate,
    StartDate: args.StartDate,
    KeyWord: args.KeyWord,
    PageNumber: args.PageNumber,
    PageSize: args.PageSize,
    Path: args.Path
  });

  return {
    success: true,
    data: {
      content: [
        {
          type: 'table',
          headers: ['提交ID', '提交信息', '作者', '提交时间'],
          rows: commits.Commits.map(commit => [
            commit.Sha.substring(0, 8), // 只显示提交ID的前8位
            commit.ShortMessage,
            commit.AuthorName,
            new Date(commit.CommitDate).toLocaleString()
          ])
        },
        {
          type: 'text',
          text: `总共 ${commits.Page.TotalRow} 个提交，当前第 ${commits.Page.PageNumber}/${commits.Page.TotalPage} 页`
        }
      ],
      metadata: {
        page: {
          current: commits.Page.PageNumber,
          total: commits.Page.TotalPage,
          size: commits.Page.PageSize,
          totalItems: commits.Page.TotalRow
        }
      }
    }
  };
}