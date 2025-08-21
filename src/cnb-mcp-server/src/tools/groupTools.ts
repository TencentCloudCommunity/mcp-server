import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { createGroup, getGroup, listGroups, listSubGroups } from '../api/group.js';
import { formatTextToolResult, formatToolError } from '../helpers/formatToolResult.js';

export default function registerGroupTools(server: McpServer) {
  server.tool(
    'list-groups',
    '获取当前用户拥有权限的顶层组织列表',
    {
      page: z.number().default(1).describe('第几页，从1开始'),
      page_size: z.number().default(10).describe('每页多少条数据'),
      search: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('仓库关键字'),
      role: z
        .preprocess(
          (val) => (val === null ? undefined : val),
          z.enum(['Guest', 'Reporter', 'Developer', 'Master', 'Owner']).optional()
        )
        .describe('最小仓库权限')
    },
    async ({ page, page_size, search, role }) => {
      try {
        const groups = await listGroups({ page, page_size, search, role });
        return formatTextToolResult(JSON.stringify(groups, null, 2), 'list-groups');
      } catch (error) {
        return formatToolError(error, 'list-groups');
      }
    }
  );

  server.tool(
    'list-sub-groups',
    '查询当前用户在指定组织下拥有指定权限的子组织列表',
    {
      group: z.string().describe('组织名称'),
      page: z.number().default(1).describe('第几页，从1开始'),
      page_size: z.number().default(10).describe('每页多少条数据'),
      access: z.preprocess((val) => (val === null ? undefined : val), z.number().optional()).describe('权限等级')
    },
    async ({ group, page, page_size, access }) => {
      try {
        const subGroups = await listSubGroups(group, { page, page_size, access });
        return formatTextToolResult(JSON.stringify(subGroups, null, 2), 'list-sub-groups');
      } catch (error) {
        return formatToolError(error, 'list-sub-groups');
      }
    }
  );

  server.tool(
    'get-group',
    '获取指定组织信息',
    {
      group: z.string().describe('组织路径')
    },
    async ({ group }) => {
      try {
        const groupInfo = await getGroup(group);
        return formatTextToolResult(JSON.stringify(groupInfo, null, 2), 'get-group');
      } catch (error) {
        return formatToolError(error, 'get-group');
      }
    }
  );

  server.tool(
    'create-group',
    '创建新组织',
    {
      path: z.string().describe('组织路径'),
      description: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('组织描述'),
      remark: z.preprocess((val) => (val === null ? undefined : val), z.string().optional()).describe('仓库备注'),
      bind_domain: z
        .preprocess((val) => (val === null ? undefined : val), z.string().optional())
        .describe('根组织绑定的域名')
    },
    async ({ path, description, remark, bind_domain }) => {
      try {
        const data = await createGroup({ path, description, remark, bind_domain });
        return formatTextToolResult(JSON.stringify(data, null, 2), 'create-group');
      } catch (error) {
        return formatToolError(error, 'create-group');
      }
    }
  );
}
