
import axios from 'axios';
import { CodingDevOpsConfig } from '../config/environment.js';

interface CodingUser {
  Id: number;
  Status: number;
  Email: string;
  GlobalKey: string;
  Name: string;
}

interface CodingProject {
  Id?: number;
  Name?: string;
  DisplayName?: string;
  Description?: string;
}

interface CodingIssue {
  Type: string;
  Name: string;
  Description: string;
  Code: number;
  [key: string]: any; // 其他可能的字段
}

interface CodingDepot {
  CreatedAt: number;
  DefaultBranch: string;
  Description: string;
  GroupId: number;
  GroupName: string;
  GroupType: string;
  HttpsUrl: string;
  Id: number;
  LastPushAt: number;
  Name: string;
  ProjectId: number;
  ProjectName: string;
  RepoType: string;
  SshUrl: string;
  VcsType: string;
  WebUrl: string;
  IsShared: string;
}

interface CodingDepotResponse {
  Response: {
    DepotData: {
      Depots: CodingDepot[];
      Page: {
        PageNumber: number;
        PageSize: number;
        TotalPage: number;
        TotalRow: number;
      };
    };
    RequestId: string;
  };
}

interface CodingCommitAuthor {
  Email: string;
  Name: string;
}

interface CodingCommit {
  Sha: string;
  ShortMessage: string;
  FullMessage: string;
  CommitDate: number;
  CreatedAt: number;
  AuthorName: string;
  AuthorEmail: string;
  Author: CodingCommitAuthor;
  Committer: CodingCommitAuthor;
  Parents: string[];
}

interface CodingCommitsResponse {
  Response: {
    Data: {
      Commits: CodingCommit[];
      Page: {
        PageNumber: number;
        PageSize: number;
        TotalPage: number;
        TotalRow: number;
      };
    };
    RequestId: string;
  };
}

interface CodingMergeRequestResponse {
  Response: {
    MergeInfo: {
      MergeRequestId: number;
      MergeRequestUrl: string;
      MergeRequestInfo: {
        Id: number;
        MergeId: number;
        Title: string;
        Describe: string;
        Status: string;
        TargetBranch: string;
        SourceBranch: string;
      };
    };
    RequestId: string;
  };
}

interface CodingProjectResponse {
  Response: {
    Project: {
      Id: number;
      Name: string;
      DisplayName: string;
      [key: string]: any;
    };
    RequestId: string;
  };
}

interface CodingUserResponse {
  Response: {
    User: CodingUser;
    RequestId: string;
  };
}

interface CodingProjectsResponse {
  Response: {
    ProjectList: CodingProject[];
    RequestId: string;
  };
}

interface CodingIssuesResponse {
  Response: {
    IssueList: CodingIssue[];
    RequestId: string;
  };
}

interface CodingIssueResponse {
  Response: {
    Issue: CodingIssue;
    RequestId: string;
  };
}

export class CodingConnection {
  private static instance: CodingConnection | null = null;
  private static config: CodingDevOpsConfig;
  private userId: number | null = null;

  private constructor(config: CodingDevOpsConfig) {
    CodingConnection.config = config;
  }

  public static initialize(config: CodingDevOpsConfig): void {
    CodingConnection.config = config;
    CodingConnection.instance = null;
  }

  public static getInstance(): CodingConnection {
    if (!this.config) {
      throw new Error('CodingConnection must be initialized with config before use');
    }

    if (!this.instance) {
      this.instance = new CodingConnection(this.config);
    }
    return this.instance;
  }

  private async getCurrentUser(): Promise<number> {
    if (this.userId) return this.userId;

    const response = await axios.post<CodingUserResponse>(
      CodingConnection.config.apiUrl,
      {
        Action: 'DescribeCodingCurrentUser'
      },
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    this.userId = response.data.Response.User.Id;
    return this.userId;
  }

  public async listProjects(projectName?: string): Promise<CodingProject[]> {
    const userId = await this.getCurrentUser();

    const requestBody: Record<string, unknown> = {
      Action: 'DescribeUserProjects',
      UserId: userId
    };

    if (projectName) {
      requestBody.ProjectName = projectName;
    }

    const response = await axios.post<CodingProjectsResponse>(
      CodingConnection.config.apiUrl,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response.ProjectList;
  }

  public async createIssue(params: {
    projectName: string;
    name: string;
    type: string;
    priority: string;
    description: string;
    parentCode?: number;
  }): Promise<CodingIssue> {
    const requestBody = {
      Action: 'CreateIssue',
      ProjectName: params.projectName,
      Name: params.name,
      Type: params.type,
      Priority: params.priority,
      Description: params.description,
      ParentCode: params.parentCode
    };

    const response = await axios.post<{
      Response: {
        Issue: CodingIssue;
        RequestId: string;
      };
    }>(
      CodingConnection.config.apiUrl,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response.Issue;
  }

  public async listIssues(params: {
    projectName: string;
    issueType?: string;
    limit?: string;
    offset?: string;
    sortKey?: string;
    sortValue?: string;
  }): Promise<CodingIssue[]> {
    const requestBody = {
      Action: 'DescribeIssueList',
      ProjectName: params.projectName,
      IssueType: params.issueType || 'ALL',
      Limit: params.limit || '20',
      Offset: params.offset || '0',
      SortKey: params.sortKey || 'UPDATED_AT',
      SortValue: params.sortValue || 'DESC'
    };

    const response = await axios.post<CodingIssuesResponse>(
      CodingConnection.config.apiUrl,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response.IssueList;
  }

  public async describeIssue(params: {
    projectName: string;
    issueCode: number;
  }): Promise<CodingIssue> {
    const requestBody = {
      Action: 'DescribeIssue',
      ProjectName: params.projectName,
      IssueCode: params.issueCode
    };

    const response = await axios.post<CodingIssueResponse>(
      CodingConnection.config.apiUrl,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response.Issue;
  }

  public async deleteIssue(params: {
    projectName: string;
    issueCode: number;
  }): Promise<void> {
    const requestBody = {
      Action: 'DeleteIssue',
      ProjectName: params.projectName,
      IssueCode: params.issueCode
    };

    await axios.post(
      CodingConnection.config.apiUrl,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );
  }

  public async deleteProject(params: {
    ProjectId: string;
  }): Promise<void> {
    const requestBody = {
      Action: 'DeleteOneProject',
      ProjectId: params.ProjectId
    };

    await axios.post(
      `${CodingConnection.config.apiUrl}/?Action=DeleteOneProject`,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );
  }

  public async createProject(params: {
    displayName: string;
    name: string;
    description?: string;
    projectTemplate?: string;
    shared?: string;
  }): Promise<CodingProject> {
    const requestBody = {
      DisplayName: params.displayName,
      Name: params.name,
      Description: params.description || '',
      ProjectTemplate: params.projectTemplate || 'DEV_OPS',
      Shared: params.shared || "0"
    };

    const response = await axios.post<{
      Response: {
        ProjectId: number;
        RequestId: string;
      };
    }>(
      `${CodingConnection.config.apiUrl}/?Action=CreateCodingProject`,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return {
      Id: response.data.Response.ProjectId
    };
  }


  public async listProjectDepots(params: {
    ProjectId: string;
    PageNumber?: string;
    PageSize?: string;
  }): Promise<CodingDepotResponse['Response']['DepotData']> {
    const requestBody = {
      Action: 'DescribeProjectDepotInfoList',
      ProjectId: params.ProjectId,
      PageNumber: params.PageNumber || '1',
      PageSize: params.PageSize || '20'
    };

    const response = await axios.post<CodingDepotResponse>(
      `${CodingConnection.config.apiUrl}/?Action=DescribeProjectDepotInfoList`,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response.DepotData;
  }

  public async createMergeRequest(params: {
    DepotPath: string;
    Title: string;
    Content: string;
    SrcBranch: string;
    DestBranch: string;
  }): Promise<CodingMergeRequestResponse['Response']> {
    const requestBody = {
      Action: 'CreateGitMergeReq',
      DepotPath: params.DepotPath,
      Title: params.Title,
      Content: params.Content,
      SrcBranch: params.SrcBranch,
      DestBranch: params.DestBranch
    };

    const response = await axios.post<CodingMergeRequestResponse>(
      `${CodingConnection.config.apiUrl}/?Action=CreateGitMergeReq`,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response;
  }

  public async getProjectByName(projectName: string): Promise<CodingProjectResponse['Response']> {
    const requestBody = {
      Action: 'DescribeProjectByName',
      ProjectName: projectName
    };

    const response = await axios.post<CodingProjectResponse>(
      `${CodingConnection.config.apiUrl}/?Action=DescribeProjectByName`,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response;
  }

  public async describeGitCommits(params: {
    DepotId?: string;
    DepotPath?: string;
    Ref: string;
    EndDate?: string;
    StartDate?: string;
    KeyWord?: string;
    PageNumber?: string;
    PageSize?: string;
    Path?: string;
  }): Promise<CodingCommitsResponse['Response']['Data']> {
    const requestBody = {
      Action: 'DescribeGitCommitsInPage',
      DepotId: params.DepotId,
      DepotPath: params.DepotPath,
      Ref: params.Ref,
      EndDate: params.EndDate,
      StartDate: params.StartDate,
      KeyWord: params.KeyWord,
      PageNumber: params.PageNumber || '1',
      PageSize: params.PageSize || '20',
      Path: params.Path
    };

    const response = await axios.post<CodingCommitsResponse>(
      `${CodingConnection.config.apiUrl}/?Action=DescribeGitCommitsInPage`,
      requestBody,
      {
        headers: {
          'Authorization': `token ${CodingConnection.config.token}`,
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      }
    );

    return response.data.Response.Data;
  }
}