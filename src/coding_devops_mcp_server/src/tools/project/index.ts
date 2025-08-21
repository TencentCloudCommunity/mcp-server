
import { listProjects } from './list.js';
import { createProject } from './create.js';
import { deleteProject } from './delete.js';
import { CodingDevOpsConfig } from '../../config/environment.js';

const definitions = [
  {
    name: 'delete_project',
    description: '删除 CODING DevOps 中的指定项目',
    inputSchema: {
      type: 'object',
      properties: {
        projectId: { 
          type: 'string', 
          description: '要删除的项目ID' 
        }
      },
      required: ['projectId'],
    }
  },
  {
    name: 'list_projects',
    description: '查询当前用户在 CODING DevOps 中的项目列表',
    inputSchema: {
      type: 'object',
      properties: {
        projectName: { 
          type: 'string', 
          description: '项目名称，支持模糊匹配' 
        }
      },
      required: [],
    }
  },
  {
    name: 'create_project',
    description: '在 CODING DevOps 中创建新项目',
    inputSchema: {
      type: 'object',
      properties: {
        name: { 
          type: 'string', 
          description: '项目标识，用于系统内部标识项目' 
        },
        displayName: { 
          type: 'string', 
          description: '项目显示名称，用于展示的项目名称' 
        },
        description: { 
          type: 'string', 
          description: '项目描述' 
        },
        projectTemplate: { 
          type: 'string', 
          description: '项目模版类型',
          enum: ['DEV_OPS', 'DEMO_BEGIN', 'CHOICE_DEMAND', 'PROJECT_MANAGE', 'CODE_HOST']
        },
        shared: { 
          type: 'string', 
          description: '项目可见性设置(0:私有,1:公开),默认不公开',
          enum: ["0", "1"]
        }
      },
      required: ['name', 'displayName', 'projectTemplate', 'shared'],
    }
  },
];

export interface ListProjectsArgs {
  projectName?: string;
  [key: string]: unknown;
}

export interface CreateProjectArgs {
  name: string;
  displayName: string;
  description?: string;
  projectTemplate: string;
  shared: string;
}

export interface DeleteProjectArgs {
  projectId: string;
}

export const projectTools = {
  initialize: (config: CodingDevOpsConfig) => ({
    listProjects: (args: ListProjectsArgs) => listProjects(args, config),
    createProject: (args: CreateProjectArgs) => createProject(args, config),
    deleteProject: (args: DeleteProjectArgs) => deleteProject(args, config),
    definitions,
  }),
  definitions,
};