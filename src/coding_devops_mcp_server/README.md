# CODING DevOps MCP Server

CODING DevOps MCP Server 是一个基于 Model Context Protocol (MCP) 的服务器实现，用于与 CODING DevOps 平台进行交互。它提供了一套标准化的接口，使得用户可以方便地管理 CODING 平台上的项目、工作项和代码。

## 功能特性

- 项目管理
  - 列出用户可访问的项目
  - 按项目名称搜索项目
  - 创建新项目
  - 删除项目
- 工作项（Issues）管理
  - 创建工作项
  - 列出工作项
  - 删除工作项
  - 需求拆解功能
  - 工作项描述查询
  - 支持工作项类型（需求、任务、缺陷、史诗）
  - 支持工作项优先级等属性设置
- 代码管理
  - 代码仓库管理
  - 提交记录查询
  - 合并请求管理
  - Git操作工具集

### 2. 智能的代码管理

- **Git集成增强**
  - 智能远程仓库解析：自动从Git配置中提取仓库信息
  - 分支管理自动化：支持自动获取当前分支信息
  - 详细的错误诊断：提供清晰的问题定位和解决建议

- **合并请求流程优化**
  - 自动分支检测：智能识别源分支和目标分支
  - 丰富的元数据支持：包含完整的合并请求信息
  - 友好的错误处理：提供详细的故障排查指南

### 3. 用户友好设计

- **错误处理机制**
  - 详细的错误提示：提供具体的错误原因和位置
  - 清晰的解决方案：为常见错误提供操作指导
  - 上下文相关的帮助：根据具体场景提供相应的解决建议

- **灵活的配置支持**
  - 环境变量配置：支持通过环境变量进行服务配置
  - 项目级默认值：支持设置默认项目等常用配置
  - 多样化的参数支持：支持必选和可选参数的灵活组合


## 安装

1. Clone this repository:
```bash
git clone https://github.com/yupengfei1209/coding_devops_mcp_server.git
cd coding_devops_mcp_server
```

2. Install dependencies:
```bash
npm install
```

3. Build the server:
```bash
npm run build
```

## 配置

服务器需要以下配置项：

1. CODING Personal Access Token (必需)
2. 项目名称 (可选)


### 添加到 MCP Client

```json
{
  "mcpServers": {
    "coding-devops": {
      "command": "npx",
      "args": ["-y", "coding-devops-mcp-server"],
      "env": {
        "CODING_TOKEN": "xxxxxxxxxxxxxx"
      },
      "disabled": false,
      "autoApprove": []
    },
  }
}
```

## 功能

### 项目管理

- `list_projects`: 列出用户可访问的项目
  ```typescript
  // 可选参数
  {
    projectName?: string; // 按项目名称筛选
  }
  ```

- `create_project`: 创建新项目
  ```typescript
  // 参数
  {
    name: string;           // 项目标识
    displayName: string;    // 项目显示名称
    projectTemplate: string;// 项目模板类型
    shared: string;         // 项目可见性(0:私有,1:公开)
    description?: string;   // 项目描述
  }
  ```

- `delete_project`: 删除项目
  ```typescript
  // 参数
  {
    projectId: string;      // 项目ID
  }
  ```

### 工作项管理

- `list_issues`: 列出工作项
  ```typescript
  // 参数
  {
    projectName: string;    // 项目名称
    issueType?: string;    // 工作项类型(DEFECT/REQUIREMENT/MISSION/EPIC)
    limit?: string;        // 返回数量限制
  }
  ```

- `create_issue`: 创建工作项
  ```typescript
  // 参数
  {
    projectName: string;    // 项目名称
    name: string;          // 工作项标题
    type: string;         // 工作项类型
    priority: string;     // 优先级(0-低/1-中/2-高/3-紧急)
    description: string;  // 描述
    parentCode?: number;  // 父工作项编号(可选)
  }
  ```

- `delete_issue`: 删除工作项
  ```typescript
  // 参数
  {
    projectName: string;    // 项目名称
    issueCode: number;     // 工作项编号
  }
  ```

- `describe_issue`: 查询工作项详情
  ```typescript
  // 参数
  {
    projectName: string;    // 项目名称
    issueCode: number;     // 工作项编号
  }
  ```

- `decompose_issue`: 需求拆解
  ```typescript
  // 参数
  {
    projectName: string;     // 项目名称
    parentIssueCode: number; // 父需求编号
    subTasks: Array<{       // 子任务列表
      name: string;         // 任务名称
      description: string;  // 任务描述
      priority: string;     // 优先级
    }>;
  }
  ```

### 代码管理

- `list_commits`: 查询代码提交记录
  ```typescript
  // 参数
  {
    DepotId?: string;      // 仓库ID（与DepotPath二选一）
    DepotPath?: string;    // 仓库路径（与DepotId二选一）
    Ref: string;          // 分支名称
    StartDate?: string;   // 查询起始日期（YYYY-MM-DD）
    EndDate?: string;     // 查询截止日期（YYYY-MM-DD）
    KeyWord?: string;     // 提交关键字
    PageNumber?: string;  // 页码
    PageSize?: string;    // 每页数量
  }
  ```

- `create_merge_request`: 创建合并请求
  ```typescript
  // 参数
  {
    workingDirectory: string; // Git仓库目录
    title: string;           // 合并请求标题
    content: string;         // 合并请求描述
    srcBranch?: string;      // 源分支（可选，默认当前分支）
    destBranch?: string;     // 目标分支（可选，默认master）
  }
  ```

- `list_depots`: 查询项目仓库列表
  ```typescript
  // 参数
  {
    projectId: string;      // 项目ID
    PageNumber?: string;    // 页码
    PageSize?: string;      // 每页数量
  }
  ```

## 开发

### 项目结构

```
src/
├── api/               # API 实现
├── config/            # 配置相关
├── tools/            # 工具实现
│   ├── code/         # 代码管理相关功能
│   │   ├── commits.ts    # 提交记录管理
│   │   ├── git_utils.ts  # Git工具集
│   │   ├── list.ts       # 代码列表
│   │   └── merge_request.ts # 合并请求管理
│   ├── issue/        # 工作项管理功能
│   │   ├── create.ts     # 创建工作项
│   │   ├── decompose.ts  # 需求拆解
│   │   ├── delete.ts     # 删除工作项
│   │   ├── describe.ts   # 工作项描述
│   │   └── list.ts       # 工作项列表
│   └── project/      # 项目管理功能
│       ├── create.ts     # 创建项目
│       ├── delete.ts     # 删除项目
│       └── list.ts       # 项目列表
├── errors.ts         # 错误定义
└── index.ts         # 主入口文件
```

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。