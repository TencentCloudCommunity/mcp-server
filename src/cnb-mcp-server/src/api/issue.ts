import CnbApiClient from './client.js';

import type { operations, definitions } from '../schema.js';

export async function listIssues(repo: string, params?: ListIssuesParams): Promise<Issue[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL(`/${repo}/-/issues`, cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<Issue[]>('GET', `${url.pathname}${url.search}`);
}

export async function getIssue(repo: string, issueId: number): Promise<Issue> {
  return CnbApiClient.getInstance().request<Issue>('GET', `/${repo}/-/issues/${issueId}`);
}

export async function createIssue(repo: string, params: CreateIssueParams): Promise<Issue> {
  const newParams = Object.entries(params).reduce((acc, [key, value]) => {
    if (value === undefined) return acc;
    Object.assign(acc, { [key]: value });
    return acc;
  }, {});
  return CnbApiClient.getInstance().request<Issue>('POST', `/${repo}/-/issues`, newParams, {
    header: { 'Content-Type': 'application/json' }
  });
}

export async function updateIssue(repo: string, issueId: number, params: UpdateIssueParams): Promise<Issue> {
  return CnbApiClient.getInstance().request<Issue>('PATCH', `/${repo}/-/issues/${issueId}`, params, {
    header: { 'Content-Type': 'application/json' }
  });
}

export async function listIssueComments(
  repo: string,
  issueId: number,
  params?: ListIssueCommentsParams
): Promise<IssueComment[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL(`/${repo}/-/issues/${issueId}/comments`, cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<IssueComment[]>('GET', `${url.pathname}${url.search}`);
}

export async function createIssueComment(
  repo: string,
  issueId: number,
  params: CreateIssueCommentParams
): Promise<IssueComment> {
  return CnbApiClient.getInstance().request<IssueComment>('POST', `/${repo}/-/issues/${issueId}/comments`, params, {
    header: { 'Content-Type': 'application/json' }
  });
}

export async function updateIssueComment(
  repo: string,
  issueId: number,
  commentId: string,
  params: UpdateIssueCommentParams
): Promise<IssueComment> {
  return CnbApiClient.getInstance().request<IssueComment>(
    'PATCH',
    `/${repo}/-/issues/${issueId}/comments/${commentId}`,
    params,
    {
      header: { 'Content-Type': 'application/json' }
    }
  );
}

export type ListIssuesParams = operations['ListIssues']['parameters']['query'];

export type Issue = definitions['api.Issue'];

export type CreateIssueParams = definitions['api.PostIssueForm'];

export type UpdateIssueParams = definitions['api.PatchIssueForm'];

export type ListIssueCommentsParams = operations['ListIssueComments']['parameters']['query'];

export type IssueComment = definitions['api.IssueComment'];

export type CreateIssueCommentParams = definitions['api.PostIssueCommentForm'];

export type UpdateIssueCommentParams = definitions['api.PatchIssueCommentForm'];

export type Label = definitions['api.Label'];

export async function listIssueLabels(repo: string, issueId: number): Promise<Label[]> {
  return CnbApiClient.getInstance().request<Label[]>('GET', `/${repo}/-/issues/${issueId}/labels`);
}

export async function addIssueLabels(repo: string, issueId: number, labels: string[]): Promise<Label[]> {
  return CnbApiClient.getInstance().request<Label[]>(
    'POST',
    `/${repo}/-/issues/${issueId}/labels`,
    { labels },
    {
      header: { 'Content-Type': 'application/json' }
    }
  );
}

export async function setIssueLabels(repo: string, issueId: number, labels: string[]): Promise<Label[]> {
  return CnbApiClient.getInstance().request<Label[]>(
    'PUT',
    `/${repo}/-/issues/${issueId}/labels`,
    { labels },
    {
      header: { 'Content-Type': 'application/json' }
    }
  );
}

export async function deleteIssueLabels(repo: string, issueId: number): Promise<void> {
  await CnbApiClient.getInstance().request('DELETE', `/${repo}/-/issues/${issueId}/labels`);
}

export async function deleteIssueLabel(repo: string, issueId: number, labelName: string): Promise<void> {
  await CnbApiClient.getInstance().request(
    'DELETE',
    `/${repo}/-/issues/${issueId}/labels/${encodeURIComponent(labelName)}`
  );
}
