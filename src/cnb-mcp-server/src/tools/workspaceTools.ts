import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { deleteWorkspace, listWorkspace } from '../api/workspace.js';
import { formatTextToolResult, formatToolError } from '../helpers/formatToolResult.js';

export default function registerWorkspaceTools(server: McpServer) {
  server.tool(
    'list-workspace',
    '获取我的云原生开发环境列表',
    {
      branch: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('分支名，例如：main'),
      start: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('查询结束时间，格式：YYYY-MM-DD HH:mm:ssZZ，例如：2024-12-01 00:00:00+0800'),
      end: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('查询开始时间，格式：YYYY-MM-DD HH:mm:ssZZ，例如：2024-12-01 00:00:00+0800'),
      page: z
        .preprocess((val) => (val === null ? undefined : val), z.number().optional())
        .describe('分页页码，从 1 开始，默认为 1'),
      page_size: z
        .preprocess((val) => (val === null ? undefined : val), z.number().optional())
        .describe('每页条数，默认为 20，最高 100'),
      slug: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('仓库路径，例如：groupname/reponame'),
      status: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['running', 'closed']).optional())
        .describe('开发环境状态，running: 开发环境已启动，closed：开发环境已关闭，默认为所有状态')
    },
    async ({ branch, page, page_size, start, end, slug, status }) => {
      try {
        const workspaces = await listWorkspace({
          branch,
          page,
          pageSize: page_size,
          start,
          end,
          slug,
          status
        });
        return formatTextToolResult(JSON.stringify(workspaces, null, 2), 'list-workspace');
      } catch (error) {
        return formatToolError(error, 'list-workspace');
      }
    }
  );

  server.tool(
    'delete-workspace',
    '删除我的云原生开发环境',
    {
      pipelineId: z.string().describe('开发环境 ID')
    },
    async ({ pipelineId }) => {
      try {
        const result = await deleteWorkspace({
          pipelineId
        });
        return formatTextToolResult(JSON.stringify(result, null, 2), 'delete-workspace');
      } catch (error) {
        return formatToolError(error, 'delete-workspace');
      }
    }
  );
}
