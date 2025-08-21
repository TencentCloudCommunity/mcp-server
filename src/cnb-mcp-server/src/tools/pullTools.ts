import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';

import {
  listPulls,
  getPull,
  createPull,
  updatePull,
  mergePull,
  listPullComments,
  createPullComment
} from '../api/pull.js';
import { formatTextToolResult, formatToolError } from '../helpers/formatToolResult.js';

export default function registerPullTools(server: McpServer) {
  server.tool(
    'list-pulls',
    '查询仓库的Pull Requests',
    {
      repo: z.string().describe('仓库路径，格式为 {group}/{repo}'),
      state: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['open', 'closed', 'all']).optional())
        .describe('Pull Request状态'),
      sort: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['created', 'updated']).optional())
        .describe('排序字段'),
      direction: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['asc', 'desc']).optional())
        .describe('排序方向'),
      page: z.number().default(1).describe('页码'),
      per_page: z.number().default(30).describe('每页数量')
    },
    async ({ repo, ...params }) => {
      try {
        const pulls = await listPulls(repo, params);
        return formatTextToolResult(JSON.stringify(pulls, null, 2), 'list-pulls');
      } catch (error) {
        return formatToolError(error, 'list-pulls');
      }
    }
  );

  server.tool(
    'get-pull',
    '获取单个Pull Request详情',
    {
      repo: z.string().describe('仓库路径，格式为 {group}/{repo}'),
      number: z.number().describe('Pull Request编号')
    },
    async ({ repo, number }) => {
      try {
        const pull = await getPull(repo, number);
        return formatTextToolResult(JSON.stringify(pull, null, 2), 'get-pull');
      } catch (error) {
        return formatToolError(error, 'get-pull');
      }
    }
  );

  server.tool(
    'create-pull',
    '创建Pull Request',
    {
      repo: z.string().describe('目标仓库路径，格式为 {group}/{repo}'),
      base: z.string().describe('目标仓库目标分支'),
      head_repo: z.string().optional().describe('来源仓库路径，格式为 {group}/{repo},不填则为目标仓库'),
      head: z.string().describe('来源仓库分支'),
      title: z.string().describe('标题'),
      body: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('描述')
    },
    async ({ repo, ...params }) => {
      try {
        const pull = await createPull(repo, params);
        return formatTextToolResult(JSON.stringify(pull, null, 2), 'create-pull');
      } catch (error) {
        return formatToolError(error, 'create-pull');
      }
    }
  );

  server.tool(
    'update-pull',
    '更新Pull Request',
    {
      repo: z.string().describe('仓库路径，格式为 {group}/{repo}'),
      number: z.number().describe('Pull Request编号'),
      title: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('标题'),
      body: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('描述'),
      state: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['open', 'closed']).optional())
        .describe('状态')
    },
    async ({ repo, number, ...params }) => {
      try {
        const pull = await updatePull(repo, number, params);
        return formatTextToolResult(JSON.stringify(pull, null, 2), 'update-pull');
      } catch (error) {
        return formatToolError(error, 'update-pull');
      }
    }
  );

  server.tool(
    'merge-pull',
    '合并Pull Request',
    {
      repo: z.string().describe('仓库路径，格式为 {group}/{repo}'),
      number: z.number().describe('Pull Request编号'),
      merge_style: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['merge', 'squash', 'rebase']).optional())
        .describe('合并方式'),
      commit_title: z.preprocess((val) => (val === null ? undefined : val), z.string()).describe('合并提交标题'),
      commit_message: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('合并提交信息')
    },
    async ({ repo, number, ...params }) => {
      try {
        const result = await mergePull(repo, number, params);
        return formatTextToolResult(JSON.stringify(result, null, 2), 'merge-pull');
      } catch (error) {
        return formatToolError(error, 'merge-pull');
      }
    }
  );

  server.tool(
    'list-pull-comments',
    '列出Pull Request的评论',
    {
      repo: z.string().describe('仓库路径，格式为 {group}/{repo}'),
      number: z.number().describe('Pull Request编号'),
      page: z.number().default(1).describe('页码'),
      per_page: z.number().default(30).describe('每页数量')
    },
    async ({ repo, number, ...params }) => {
      try {
        const comments = await listPullComments(repo, number, params);
        return formatTextToolResult(JSON.stringify(comments, null, 2), 'list-pull-comments');
      } catch (error) {
        return formatToolError(error, 'list-pull-comments');
      }
    }
  );

  server.tool(
    'create-pull-comment',
    '创建Pull Request评论',
    {
      repo: z.string().describe('仓库路径，格式为 {group}/{repo}'),
      number: z.number().describe('Pull Request编号'),
      body: z.string().describe('评论内容')
    },
    async ({ repo, number, body }) => {
      try {
        await createPullComment(repo, number, { body });
        return formatTextToolResult('Comment created', 'create-pull-comment');
      } catch (error) {
        return formatToolError(error, 'create-pull-comment');
      }
    }
  );
}
