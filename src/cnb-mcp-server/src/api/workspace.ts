import CnbApiClient from './client.js';

import type { operations, definitions } from '../schema.js';

export async function listWorkspace(params?: ListWorkspaceParams): Promise<Workspace> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL('/workspace/list', cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<Workspace>('GET', `${url.pathname}${url.search}`);
}

export async function deleteWorkspace(params?: DeleteWorkspaceParams): Promise<DeleteWorkspaceResponse> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL('/workspace/delete', cnbInst.baseUrl);
  return cnbInst.request<DeleteWorkspaceResponse>('POST', `${url.pathname}`, params, {
    header: { 'Content-Type': 'application/json' }
  });
}

export type ListWorkspaceParams = operations['ListWorkspaces']['parameters']['query'];

export type Workspace = definitions['dto.WorkspaceListResult'];

export type DeleteWorkspaceParams = operations['DeleteWorkspace']['parameters']['body']['request'];

export type DeleteWorkspaceResponse = definitions['dto.WorkspaceDeleteResult'];
