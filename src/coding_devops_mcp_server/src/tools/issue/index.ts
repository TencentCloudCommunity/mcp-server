
import { listIssues } from './list.js';
import { decomposeIssue } from './decompose.js';
import { createIssue } from './create.js';
import { deleteIssue } from './delete.js';
import { describeIssue } from './describe.js';
import { CodingDevOpsConfig } from '../../config/environment.js';

const definitions = [
  {
    name: 'list_issues',
    description: '查询我最近更新的工作项，可以根据其中IssueType查询，默认查询limit为20条',
    inputSchema: {
      type: 'object',
      properties: {
        projectName: {
          type: 'string',
          description: '项目名称，注意是项目的 name 不是 displayName,建议先通过 list_projects 查看项目名称',
        },
        issueType: {
          type: 'string',
          description: '事项类型，可选值为： ALL - 全部事项 DEFECT - 缺陷 REQUIREMENT - 需求 MISSION - 任务 EPIC - 史诗',
        },
        limit: {
          type: 'string',
          description: '事项数量，默认为20条',
        }
      },
      required: ['projectName'],
    }
  },
  {
    name: 'create_issue',
    description: '创建新的工作项',
    inputSchema: {
      type: 'object',
      properties: {
        projectName: {
          type: 'string',
          description: '项目名称',
        },
        name: {
          type: 'string',
          description: '事项标题',
        },
        type: {
          type: 'string',
          description: '事项类型，可选值为： DEFECT - 缺陷 REQUIREMENT - 需求 MISSION - 任务 EPIC - 史诗',
        },
        priority: {
          type: 'string',
          description: '优先级，可选值为：0 - 低 1 - 中 2 - 高 3 - 紧急',
        },
        description: {
          type: 'string',
          description: '事项描述',
        },
        parentCode: {
          type: 'number',
          description: '父事项编号，如果设置了此字段，创建的事项将成为该父事项的子事项',
        }
      },
      required: ['projectName', 'name', 'type', 'priority', 'description'],
    }
  },
  {
    name: 'delete_issue',
    description: '删除工作项',
    inputSchema: {
      type: 'object',
      properties: {
        projectName: {
          type: 'string',
          description: '项目名称',
        },
        issueCode: {
          type: 'number',
          description: '事项编号',
        }
      },
      required: ['projectName', 'issueCode'],
    }
  },
  {
    name: 'describe_issue',
    description: '查询事项详情',
    inputSchema: {
      type: 'object',
      properties: {
        projectName: {
          type: 'string',
          description: '项目名称',
        },
        issueCode: {
          type: 'number',
          description: '事项编号',
        }
      },
      required: ['projectName', 'issueCode'],
    }
  },
  {
    name: 'decompose_issue',
    description: '将一个需求拆解为多个子任务',
    inputSchema: {
      type: 'object',
      properties: {
        projectName: {
          type: 'string',
          description: '项目名称，注意是项目的 name 不是 displayName',
        },
        parentIssueCode: {
          type: 'number',
          description: '父需求的编号',
        },
        subTasks: {
          type: 'array',
          description: '子任务列表',
          items: {
            type: 'object',
            properties: {
              name: {
                type: 'string',
                description: '子任务名称',
              },
              description: {
                type: 'string',
                description: '子任务描述',
              },
              priority: {
                type: 'string',
                description: '优先级，可选值为：0 - 低 1 - 中 2 - 高 3 - 紧急',
              }
            },
            required: ['name', 'description', 'priority']
          }
        }
      },
      required: ['projectName', 'parentIssueCode', 'subTasks'],
    }
  }
];

export const issueTools = {
  initialize: (config: CodingDevOpsConfig) => ({
    listIssues: (args: { 
      projectName: string;
      issueType?: string;
      limit?: string;
    }) => listIssues(args, config),
    createIssue: (args: {
      projectName: string;
      name: string;
      type: string;
      priority: string;
      description: string;
      parentCode?: number;
    }) => createIssue(args, config),
    deleteIssue: (args: {
      projectName: string;
      issueCode: number;
    }) => deleteIssue(args, config),
    describeIssue: (args: {
      projectName: string;
      issueCode: number;
    }) => describeIssue(args, config),
    decomposeIssue: (args: {
      projectName: string;
      parentIssueCode: number;
      subTasks: Array<{
        name: string;
        description: string;
        priority: string;
      }>;
    }) => decomposeIssue(args, config),
    definitions,
  }),
  definitions,
};