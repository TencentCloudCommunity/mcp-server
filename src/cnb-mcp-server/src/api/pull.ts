import CnbApiClient from './client.js';

import type { operations, definitions } from '../schema.js';

export async function listPulls(repo: string, params?: ListPullsParams): Promise<PullRequest[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL(`/${repo}/-/pulls`, cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<PullRequest[]>('GET', `${url.pathname}${url.search}`);
}

export async function getPull(repo: string, number: number): Promise<PullRequest> {
  return CnbApiClient.getInstance().request<PullRequest>('GET', `/${repo}/-/pulls/${number}`);
}

export async function createPull(repo: string, params: CreatePullParams): Promise<PullRequest> {
  return CnbApiClient.getInstance().request<PullRequest>('POST', `/${repo}/-/pulls`, params, {
    header: { 'Content-Type': 'application/json' }
  });
}

export async function updatePull(repo: string, number: number, params: UpdatePullParams): Promise<PullRequest> {
  return CnbApiClient.getInstance().request<PullRequest>('PATCH', `/${repo}/-/pulls/${number}`, params, {
    header: { 'Content-Type': 'application/json' }
  });
}

export async function mergePull(repo: string, number: number, params: MergePullParams): Promise<MergePullResponse> {
  return CnbApiClient.getInstance().request<MergePullResponse>('PUT', `/${repo}/-/pulls/${number}/merge`, params, {
    header: { 'Content-Type': 'application/json' }
  });
}

export async function listPullComments(
  repo: string,
  number: number,
  params?: ListPullCommentsParams
): Promise<PullRequestComment[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL(`/${repo}/-/pulls/${number}/comments`, cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<PullRequestComment[]>('GET', `${url.pathname}${url.search}`);
}

export async function createPullComment(
  repo: string,
  number: number,
  params: CreatePullCommentParams
): Promise<unknown> {
  const response = await CnbApiClient.getInstance().request<Response>(
    'POST',
    `/${repo}/-/pulls/${number}/comments`,
    params,
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

export type ListPullsParams = operations['ListPulls']['parameters']['query'];
export type CreatePullParams = definitions['api.PullCreationForm'];
export type UpdatePullParams = definitions['api.PatchPullRequest'];
export type MergePullParams = definitions['api.MergePullRequest'];
export type ListPullCommentsParams = operations['ListPullComments']['parameters']['query'];
export type CreatePullCommentParams = definitions['api.PullCommentCreationForm'];
export type PullRequest = definitions['api.PullRequest'];
export type PullRequestComment = definitions['api.PullRequestComment'];
export type MergePullResponse = definitions['api.MergePullResponse'];
