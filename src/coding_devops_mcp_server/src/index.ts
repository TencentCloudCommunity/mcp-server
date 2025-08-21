#!/usr/bin/env node
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import {
  CallToolRequestSchema,
  ErrorCode,
  ListToolsRequestSchema,
  McpError,
  JSONRPCResponseSchema,
  JSONRPCResponse
} from '@modelcontextprotocol/sdk/types.js';

// 导入所有工具
import { issueTools } from './tools/issue/index.js';
import { projectTools } from './tools/project/index.js';
import { codeTools } from './tools/code/index.js';
import { CodingDevOpsConfig, createConfig } from './config/environment.js';

import { ListToolsResult } from '@modelcontextprotocol/sdk/types.js';

// 工具定义
type ToolDefinition = any;
type ToolInstances = any;

// 类型验证
function validateArgs<T>(args: Record<string, unknown> | undefined, errorMessage: string): T {
  if (!args) {
    throw new McpError(ErrorCode.InvalidParams, errorMessage);
  }
  return args as T;
}

// 响应格式
type MCPResponse = JSONRPCResponse["result"]

// 响应格式化
function formatResponse(data: unknown): MCPResponse {
  if (data && typeof data === 'object' && 'content' in data) {
    return data as MCPResponse;
  }
  return {
    content: [
      {
        type: 'text',
        text: JSON.stringify(data, null, 2),
      },
    ],
  };
}

// 服务器
class CodingDevOpsServer {
  private server: Server;
  private config: CodingDevOpsConfig;
  private toolDefinitions: ToolDefinition[];

  constructor(options?: Partial<Omit<CodingDevOpsConfig, 'apiUrl'>>) {
    this.config = createConfig(options);
    
    // 使用配置初始化工具
    const toolInstances = {
      issue: issueTools.initialize(this.config),
      project: projectTools.initialize(this.config),
      code: codeTools.initialize(this.config),
    };

    // 合并所有工具定义
    this.toolDefinitions = [
      ...toolInstances.issue.definitions,
      ...toolInstances.project.definitions,
      ...toolInstances.code.definitions,
    ];

    this.server = new Server(
      {
        name: 'coding-devops-mcp-server',
        version: '0.1.0',
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.setupToolHandlers(toolInstances);
    
    // 错误处理
    this.server.onerror = (error) => console.error('[MCP Error]', error);
    process.on('SIGINT', async () => {
      await this.server.close();
      process.exit(0);
    });
  }

  private setupToolHandlers(tools: ToolInstances) {
    // 列出可用工具
    this.server.setRequestHandler(ListToolsRequestSchema, async () => ({
      tools: this.toolDefinitions,
    }));

    // 处理工具调用
    this.server.setRequestHandler(CallToolRequestSchema, async (request, extra) => {
      try {
        let result;
        switch (request.params.name) {
          // 项目工具
          case 'list_projects':
            result = await tools.project.listProjects(request.params.arguments);
            break;
          case 'create_project':
            result = await tools.project.createProject(request.params.arguments);
            break;
          case 'delete_project':
            result = await tools.project.deleteProject(request.params.arguments);
            break;
          // 工作项工具
          case 'list_issues':
            result = await tools.issue.listIssues(request.params.arguments);
            break;
          case 'create_issue':
            result = await tools.issue.createIssue(request.params.arguments);
            break;
          case 'delete_issue':
            result = await tools.issue.deleteIssue(request.params.arguments);
            break;
          case 'describe_issue':
            result = await tools.issue.describeIssue(request.params.arguments);
            break;
          case 'decompose_issue':
            result = await tools.issue.decomposeIssue(request.params.arguments);
            break;
          // 代码工具
          case 'list_depots':
            result = await tools.code.listDepots(request.params.arguments);
            break;
          case 'create_merge_request':
            result = await tools.code.createMergeRequest(request.params.arguments);
            break;
          case 'list_commits':
            result = await tools.code.listCommits(request.params.arguments);
            break;

          
          default:
            throw new McpError(
              ErrorCode.MethodNotFound,
              `未知工具: ${request.params.name}`
            );
        }

        // 确保响应格式一致
        const response = formatResponse(result);
        return {
          _meta: request.params._meta,
          ...response
        };
      } catch (error: unknown) {
        if (error instanceof McpError) throw error;
        const errorMessage = error instanceof Error ? error.message : '未知错误';
        throw new McpError(
          ErrorCode.InternalError,
          `CODING DevOps API 错误: ${errorMessage}`
        );
      }
    });

  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error('CODING DevOps MCP 服务器正在通过 stdio 运行');
  }
}

// 允许通过构造函数或环境变量进行配置
const server = new CodingDevOpsServer();
server.run().catch(console.error);