import CnbApiClient from './client.js';

import type { operations, definitions } from '../schema.js';

export async function listRepositories(params?: ListRepositoriesParams): Promise<Repository[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL('/user/repos', cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<Repository[]>('GET', `${url.pathname}${url.search}`);
}

export async function listGroupRepositories(
  group: string,
  params?: ListGroupRepositoriesParams
): Promise<GroupRepository[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL(`/${group}/-/repos`, cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<GroupRepository[]>('GET', `${url.pathname}${url.search}`);
}

export async function getRepository(repo: string): Promise<Repository> {
  return CnbApiClient.getInstance().request<Repository>('GET', `/${repo}`);
}

export async function createRepository(group: string, params: CreateRepositoryParams): Promise<unknown> {
  const body = Object.entries(params).reduce((acc, [key, value]) => {
    if (value === undefined) return acc;
    Object.assign(acc, { [key]: value });
    return acc;
  }, {});
  const response = await CnbApiClient.getInstance().request<Response>(
    'POST',
    `/${group}/-/repos`,
    body,
    {
      header: { 'Content-Type': 'application/json' }
    },
    'raw'
  );
  if (response.status === 201) {
    return { message: 'Created' };
  } else {
    return { status: response.status, message: response.statusText };
  }
}

export type ListRepositoriesParams = operations['GetRepos']['parameters']['query'];

export type ListGroupRepositoriesParams = operations['GetGroupSubRepos']['parameters']['query'];

export type Repository = definitions['dto.Repos4User'];

export type GroupRepository = definitions['dto.Repos4UserBase'];

export type CreateRepositoryParams = definitions['dto.CreateRepoReq'];
