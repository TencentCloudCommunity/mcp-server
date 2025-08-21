
import { CodingDevOpsConfig } from '../../config/environment.js';
import { CodingConnection } from '../../api/coding_connection.js';
import { McpError, ErrorCode } from '../../errors.js';

interface SubTask {
  name: string;
  description: string;
  priority: string;
}

export async function decomposeIssue(args: {
  projectName: string;
  parentIssueCode: number;
  subTasks: SubTask[];
}, config: CodingDevOpsConfig) {
  if (!args.projectName) {
    throw new McpError(ErrorCode.InvalidParams, 'projectName is required');
  }
  if (!args.parentIssueCode) {
    throw new McpError(ErrorCode.InvalidParams, 'parentIssueCode is required');
  }
  if (!args.subTasks || args.subTasks.length === 0) {
    throw new McpError(ErrorCode.InvalidParams, 'subTasks is required and cannot be empty');
  }

  CodingConnection.initialize(config);
  const connection = CodingConnection.getInstance();

  // 验证父需求是否存在
  await connection.describeIssue({
    projectName: args.projectName,
    issueCode: args.parentIssueCode
  });

  const createdTasks = [];
  const errors = [];

  // 创建子任务并设置父需求关系
  for (const task of args.subTasks) {
    try {
      // 创建子任务，直接设置父需求关系
      const issue = await connection.createIssue({
        projectName: args.projectName,
        name: task.name,
        type: 'MISSION', // 子任务类型固定为子工作项SUB_TASK
        priority: task.priority,
        description: task.description,
        parentCode: args.parentIssueCode
      });

      createdTasks.push({
        name: issue.Name,
        code: issue.Code
      });
    } catch (error) {
      errors.push({
        taskName: task.name,
        error: error instanceof Error ? error.message : String(error)
      });
    }
  }

  // 生成结果消息
  let resultMessage = `需求拆解结果：\n\n`;
  
  if (createdTasks.length > 0) {
    resultMessage += `成功创建的子任务：\n`;
    createdTasks.forEach(task => {
      resultMessage += `- ${task.name} (编号: ${task.code})\n`;
    });
  }

  if (errors.length > 0) {
    resultMessage += `\n创建失败的子任务：\n`;
    errors.forEach(error => {
      resultMessage += `- ${error.taskName}: ${error.error}\n`;
    });
  }

  return {
    content: [
      {
        type: 'text',
        text: resultMessage,
      },
    ],
  };
}