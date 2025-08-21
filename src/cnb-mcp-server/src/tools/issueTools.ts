import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import {
  createIssue,
  createIssueComment,
  getIssue,
  listIssueComments,
  listIssues,
  updateIssue,
  updateIssueComment,
  listIssueLabels,
  addIssueLabels,
  setIssueLabels,
  deleteIssueLabels,
  deleteIssueLabel
} from '../api/issue.js';
import { formatTextToolResult, formatToolError } from '../helpers/formatToolResult.js';

export default function registerIssueTools(server: McpServer) {
  server.tool(
    'list-issues',
    '查询仓库的 Issues',
    {
      repo: z.string().describe('仓库路径'),
      page: z.number().default(1).describe('第几页,从1开始'),
      page_size: z.number().default(10).describe('每页多少条数据,默认是30'),
      state: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['open', 'closed']).optional())
        .describe('Issue 状态'),
      keyword: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 关键字'),
      priority: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 优先级,example: p0,p1,p2,p3'),
      labels: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 标签'),
      authors: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 作者的名字, example: 张三,李四'),
      assignees: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 处理人,example: 张三,李四,-; - means assign to nobody'),
      updated_time_begin: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 更新时间的范围,开始时间点,example: 2022-01-31'),
      updated_time_end: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 更新时间的范围,结束时间点,example: 2022-01-31'),
      order_by: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 排序顺序.example: created_at, -updated_at, reference_count。‘-’ prefix means descending order')
    },
    async ({
      repo,
      page,
      page_size,
      state,
      keyword,
      priority,
      labels,
      authors,
      assignees,
      updated_time_begin,
      updated_time_end,
      order_by
    }) => {
      try {
        const issues = await listIssues(repo, {
          page,
          page_size,
          state,
          keyword,
          priority,
          labels,
          authors,
          assignees,
          updated_time_begin,
          updated_time_end,
          order_by
        });
        return formatTextToolResult(JSON.stringify(issues, null, 2), 'list-issues');
      } catch (error) {
        return formatToolError(error, 'list-issues');
      }
    }
  );

  server.tool(
    'get-issue',
    '获取指定 Issue 信息',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID')
    },
    async ({ repo, issueId }) => {
      try {
        const issues = await getIssue(repo, issueId);
        return formatTextToolResult(JSON.stringify(issues, null, 2), 'get-issue');
      } catch (error) {
        return formatToolError(error, 'get-issue');
      }
    }
  );

  server.tool(
    'create-issue',
    '创建一个 Issue. 如需添加 Issue 标签，需要另外调用 add-issue-labels',
    {
      repo: z.string().describe('仓库路径'),
      title: z.string().describe('Issue 标题'),
      body: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 描述'),
      assignees: z
        .preprocess((val) => (val === null ? undefined : val), z.array(z.string()).optional())
        .describe('一个或多个 Issue 处理人的用户名'),
      labels: z
        .preprocess((val) => (val === null ? undefined : val), z.array(z.string()).optional())
        .describe('一个或多个 Issue 标签'),
      priority: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 优先级'),
      end_date: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 截止时间，格式为 YYYY-MM-DD'),
      start_date: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 起始时间，格式为 YYYY-MM-DD')
    },
    async ({ repo, title, body, assignees, labels, priority, end_date, start_date }) => {
      try {
        const issue = await createIssue(repo, {
          title,
          body,
          assignees,
          labels,
          priority,
          end_date,
          start_date
        });
        return formatTextToolResult(JSON.stringify(issue, null, 2), 'create-issue');
      } catch (error) {
        return formatToolError(error, 'create-issue');
      }
    }
  );

  server.tool(
    'update-issue',
    '更新一个 Issue。 如需更新 Issue 标签，需要另外调用 set-issue-labels',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID'),
      title: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 标题'),
      body: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 描述'),
      priority: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 优先级'),
      end_date: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 截止时间，格式为 YYYY-MM-DD'),
      start_date: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('Issue 起始时间，格式为 YYYY-MM-DD'),
      state: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('Issue 状态'),
      state_reason: z
        .preprocess(
          (val) => (val === null ? undefined : val),
          z.enum(['completed', 'not_planned', 'reopened']).optional()
        )
        .describe('Issue 状态原因')
    },
    async ({ repo, issueId, title, body, priority, end_date, start_date, state, state_reason }) => {
      try {
        const issue = await updateIssue(repo, issueId, {
          title,
          body,
          priority,
          end_date,
          start_date,
          state,
          state_reason
        });
        return formatTextToolResult(JSON.stringify(issue, null, 2), 'update-issue');
      } catch (error) {
        return formatToolError(error, 'update-issue');
      }
    }
  );

  server.tool(
    'list-issue-comments',
    '查询 Issue 评论列表',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID'),
      page: z.number().default(1).describe('第几页,从1开始'),
      page_size: z.number().default(30).describe('每页多少条数据,默认是30')
    },
    async ({ repo, issueId, page, page_size }) => {
      try {
        const comments = await listIssueComments(repo, issueId, { page, page_size });
        return formatTextToolResult(JSON.stringify(comments, null, 2), 'list-issue-comments');
      } catch (error) {
        return formatToolError(error, 'list-issue-comments');
      }
    }
  );

  server.tool(
    'create-issue-comment',
    '创建一个 Issue 评论',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID'),
      body: z.string().describe('评论内容')
    },
    async ({ repo, issueId, body }) => {
      try {
        const comment = await createIssueComment(repo, issueId, { body });
        return formatTextToolResult(JSON.stringify(comment, null, 2), 'create-issue-comment');
      } catch (error) {
        return formatToolError(error, 'create-issue-comment');
      }
    }
  );

  server.tool(
    'update-issue-comment',
    '更新一个 Issue 评论',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID'),
      commentId: z.string().describe('评论 ID'),
      body: z.string().describe('评论内容')
    },
    async ({ repo, issueId, commentId, body }) => {
      try {
        const comment = await updateIssueComment(repo, issueId, commentId, { body });
        return formatTextToolResult(JSON.stringify(comment, null, 2), 'update-issue-comment');
      } catch (error) {
        return formatToolError(error, 'update-issue-comment');
      }
    }
  );

  server.tool(
    'list-issue-labels',
    '查询指定Issue的标签',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID')
    },
    async ({ repo, issueId }) => {
      try {
        const labels = await listIssueLabels(repo, issueId);
        return formatTextToolResult(JSON.stringify(labels, null, 2), 'list-issue-labels');
      } catch (error) {
        return formatToolError(error, 'list-issue-labels');
      }
    }
  );

  server.tool(
    'add-issue-labels',
    '为指定Issue添加标签',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID'),
      labels: z.array(z.string()).describe('要添加的标签列表,每个标签需要从仓库标签列表中选择')
    },
    async ({ repo, issueId, labels }) => {
      try {
        const result = await addIssueLabels(repo, issueId, labels);
        return formatTextToolResult(JSON.stringify(result, null, 2), 'add-issue-labels');
      } catch (error) {
        return formatToolError(error, 'add-issue-labels');
      }
    }
  );

  server.tool(
    'set-issue-labels',
    '设置Issue的标签（替换所有现有标签）',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID'),
      labels: z.array(z.string()).describe('新的标签列表（将替换所有现有标签）,每个标签需要从仓库标签列表中选择')
    },
    async ({ repo, issueId, labels }) => {
      try {
        const result = await setIssueLabels(repo, issueId, labels);
        return formatTextToolResult(JSON.stringify(result, null, 2), 'set-issue-labels');
      } catch (error) {
        return formatToolError(error, 'set-issue-labels');
      }
    }
  );

  server.tool(
    'delete-issue-labels',
    '删除Issue的所有标签',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID')
    },
    async ({ repo, issueId }) => {
      try {
        await deleteIssueLabels(repo, issueId);
        return formatTextToolResult('All labels deleted', 'delete-issue-labels');
      } catch (error) {
        return formatToolError(error, 'delete-issue-labels');
      }
    }
  );

  server.tool(
    'delete-issue-label',
    '删除Issue的指定标签',
    {
      repo: z.string().describe('仓库路径'),
      issueId: z.number().describe('Issue ID'),
      labelName: z.string().describe('要删除的标签名称')
    },
    async ({ repo, issueId, labelName }) => {
      try {
        await deleteIssueLabel(repo, issueId, labelName);
        return formatTextToolResult(`${labelName} deleted`, 'delete-issue-label');
      } catch (error) {
        return formatToolError(error, 'delete-issue-label');
      }
    }
  );
}
