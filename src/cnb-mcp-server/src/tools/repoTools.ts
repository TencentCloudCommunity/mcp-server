import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { createRepository, getRepository, listGroupRepositories, listRepositories } from '../api/repository.js';
import { getUser } from '../api/user.js';
import { formatTextToolResult, formatToolError } from '../helpers/formatToolResult.js';

export default function registerRepoTools(server: McpServer) {
  server.tool(
    'list-repositories',
    '获取当前用户拥有指定权限及其以上权限的仓库',
    {
      page: z.number().default(1).describe('第几页,从1开始,默认值是1'),
      page_size: z.number().default(10).describe('每页多少条数据,默认值为10'),
      search: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('查询关键字'),
      filter_type: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['private', 'public', 'encrypted']).optional())
        .describe('仓库类型,为空表示所有仓库类型,默认值为空'),
      role: z
        .preprocess(
          (val) => (val === null ? undefined : val),
          z.enum(['Reporter', 'Developer', 'Master', 'Owner']).optional()
        )
        .describe('最小仓库权限,当用户未指定角色时,需要主动传入Reporter'),
      order_by: z
        .preprocess(
          (val) => (val === null ? undefined : val),
          z.enum(['created_at', 'last_updated_at', 'stars']).optional()
        )
        .describe('排序类型,默认值是last_updated_at'),
      desc: z
        .preprocess((val) => (val === null ? undefined : val), z.boolean().optional())
        .describe('是否开启倒叙排序，默认值是false')
    },
    async ({ page, page_size, search, filter_type, role, order_by, desc }) => {
      try {
        const repos = await listRepositories({ page, page_size, search, filter_type, role, order_by, desc });
        return formatTextToolResult(JSON.stringify(repos, null, 2), 'list-repositories');
      } catch (error) {
        return formatToolError(error, 'list-repositories');
      }
    }
  );

  server.tool(
    'list-group-repositories',
    '获取分组里当前用户有权限的仓库',
    {
      group: z.string().describe('组织名称'),
      page: z.number().default(1).describe('第几页,从1开始,默认值是1'),
      page_size: z.number().default(10).describe('每页多少条数据,默认值为10'),
      search: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('仓库关键字'),
      filter_type: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['private', 'public', 'encrypted']).optional())
        .describe('仓库类型'),
      descendant: z
        .preprocess((val) => (val === null ? undefined : val), z.enum(['all', 'sub', 'grand']).optional())
        .describe('查全部、直接属于当前组织的仓库、子组织的仓库'),
      order_by: z
        .preprocess(
          (val) => (val === null ? undefined : val),
          z.enum(['created_at', 'last_updated_at', 'stars', 'slug_path']).optional()
        )
        .describe('排序类型'),
      desc: z
        .preprocess((val) => (val === null ? undefined : val), z.boolean().optional())
        .describe('是否开启倒叙排序，默认值是false')
    },
    async ({ group, page, page_size, search, filter_type, descendant, order_by, desc }) => {
      try {
        const repos = await listGroupRepositories(group, {
          page,
          page_size,
          search,
          filter_type,
          descendant,
          order_by,
          desc
        });
        return formatTextToolResult(JSON.stringify(repos, null, 2), 'list-group-repositories');
      } catch (error) {
        return formatToolError(error, 'list-group-repositories');
      }
    }
  );

  server.tool(
    'get-repository',
    '获取指定仓库信息',
    {
      repo: z.string().describe('仓库路径')
    },
    async ({ repo }) => {
      try {
        const repoInfo = await getRepository(repo);
        return formatTextToolResult(JSON.stringify(repoInfo, null, 2), 'get-repository');
      } catch (error) {
        return formatToolError(error, 'get-repository');
      }
    }
  );

  server.tool(
    'create-repository',
    '创建仓库',
    {
      group: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('仓库所属分组'),
      name: z.string().describe('仓库名称'),
      description: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('仓库描述'),
      license: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('仓库许可'),
      visibility: z
        .preprocess(
          (val) => (val === null ? undefined : val),
          z.enum(['public', 'private', 'secret']).default('public')
        )
        .describe('仓库可见性')
    },
    async ({ group, name, description, license, visibility }) => {
      let repoGroup = group;
      if (!repoGroup) {
        const { username = '' } = await getUser();
        repoGroup = username;
      }
      try {
        const data = await createRepository(repoGroup, { name, description, license, visibility });
        return formatTextToolResult(JSON.stringify(data, null, 2), 'create-repository');
      } catch (error) {
        return formatToolError(error, 'create-repository');
      }
    }
  );

  server.tool(
    'get-current-repo',
    '获取当前仓库对应的CNB仓库信息',
    {
      remote_url: z.string().describe('远程仓库URL, 需要先执行`git remote get-url origin`命令获取')
    },
    async ({ remote_url }) => {
      try {
        let repoPath = '';

        if (remote_url.startsWith('git@')) {
          // SSH 格式: git@example.com:group/repo.git
          const match = remote_url.match(/git@[^:]+:(.+?)(?:\.git)?$/);
          if (match) {
            repoPath = match[1];
          }
        } else if (remote_url.startsWith('http')) {
          // HTTPS 格式: https://example.com/group/repo.git
          const match = remote_url.match(/https?:\/\/[^/]+\/(.+?)(?:\.git)?$/);
          if (match) {
            repoPath = match[1];
          }
        }

        if (!repoPath) {
          return formatToolError(`无法从远程仓库URL解析出仓库路径: ${remote_url}`, 'get-current-repo');
        }

        // 获取仓库信息
        const data = await getRepository(repoPath);
        return formatTextToolResult(JSON.stringify(data, null, 2), 'get-current-repo');
      } catch (error) {
        return formatToolError(error, 'get-current-repo');
      }
    }
  );
}
