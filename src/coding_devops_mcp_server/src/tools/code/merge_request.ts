
import { ErrorCode, McpError } from '@modelcontextprotocol/sdk/types.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { CodingDevOpsConfig } from '../../config/environment.js';
import { getGitConfig, parseRemoteUrl } from './git_utils.js';

export async function createMergeRequest(args: {
  title: string;
  content: string;
  srcBranch?: string;
  destBranch?: string;
  workingDirectory: string;
}, config: CodingDevOpsConfig) {
  try {
    return await createMergeRequestInternal(args, config);
  } catch (error: unknown) {
    if (error instanceof McpError) {
      throw error;
    }
    
    // 为其他类型的错误提供更详细的上下文信息
    const errorMessage = error instanceof Error ? error.message : '未知错误';
    throw new McpError(
      ErrorCode.InvalidParams,
      `创建合并请求失败。\n` +
      `错误信息: ${errorMessage}\n` +
      `当前工作目录: ${args.workingDirectory}\n` +
      '请确保:\n' +
      '1. 当前目录是有效的 Git 仓库\n' +
      '2. 已正确配置 CODING 远程仓库\n' +
      '3. 有足够的权限创建合并请求'
    );
  }
}

async function createMergeRequestInternal(args: {
  title: string;
  content: string;
  srcBranch?: string;
  destBranch?: string;
  workingDirectory: string;
}, config: CodingDevOpsConfig) {
  if (!args.title) {
    throw new McpError(ErrorCode.InvalidParams, 'title is required');
  }
  if (!args.content) {
    throw new McpError(ErrorCode.InvalidParams, 'content is required');
  }

  // 验证工作目录参数
  if (!args.workingDirectory) {
    throw new McpError(
      ErrorCode.InvalidParams,
      '未指定工作目录。请提供工作目录参数。'
    );
  }

  // 初始化连接
  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();
  
  // 获取 Git 仓库信息
  const { branch: currentBranch, remoteUrl } = await getGitConfig(args.workingDirectory);
  const { depotPath } = parseRemoteUrl(remoteUrl);

  if (!depotPath) {
    throw new McpError(
      ErrorCode.InvalidParams, 
      `无法从远程仓库 URL "${remoteUrl}" 解析仓库路径。\n` +
      `当前工作目录: ${args.workingDirectory}\n` +
      '请确保:\n' +
      '1. Git 仓库已正确配置远程地址\n' +
      '2. 远程地址格式正确（例如: git@e.coding.net:team-name/project-name/repo-name.git）'
    );
  }

  // 使用当前分支作为源分支，master 作为默认目标分支
  const srcBranch = args.srcBranch || currentBranch;
  const destBranch = args.destBranch || 'master';
  
  const response = await connection.createMergeRequest({
    DepotPath: depotPath,
    Title: args.title,
    Content: args.content,
    SrcBranch: srcBranch,
    DestBranch: destBranch
  });

  return {
    success: true,
    data: {
      content: [
        {
          type: 'text',
          text: `成功创建合并请求，编号：${response.MergeInfo.MergeRequestId}`
        },
        {
          type: 'text',
          text: `从分支 ${srcBranch} 合并到 ${destBranch}`
        },
        {
          type: 'text',
          text: `查看详情：${response.MergeInfo.MergeRequestUrl}`
        }
      ],
      metadata: {
        mergeRequest: {
          id: response.MergeInfo.MergeRequestInfo.Id,
          iid: response.MergeInfo.MergeRequestId,
          mergeId: response.MergeInfo.MergeRequestInfo.MergeId,
          title: response.MergeInfo.MergeRequestInfo.Title,
          description: response.MergeInfo.MergeRequestInfo.Describe,
          status: response.MergeInfo.MergeRequestInfo.Status,
          depotPath,
          srcBranch,
          destBranch,
          url: response.MergeInfo.MergeRequestUrl
        }
      }
    }
  };
}