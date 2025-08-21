import CnbApiClient from './client.js';

import type { operations, definitions } from '../schema.js';

export async function listGroups(params?: ListGroupsParams): Promise<Group[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL('/user/groups', cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<Group[]>('GET', `${url.pathname}${url.search}`);
}

export async function listSubGroups(group: string, params?: ListSubGroupsParams): Promise<Group[]> {
  const cnbInst = CnbApiClient.getInstance();
  const url = new URL(`/user/groups/${group}`, cnbInst.baseUrl);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value === undefined) continue;
      url.searchParams.set(key, value.toString());
    }
  }
  return cnbInst.request<Group[]>('GET', `${url.pathname}${url.search}`);
}

export async function getGroup(group: string): Promise<Group> {
  const cnbInst = CnbApiClient.getInstance();
  return cnbInst.request<Group>('GET', `/${group}`);
}

export async function createGroup(params: CreateGroupParams): Promise<unknown> {
  const body = Object.entries(params).reduce((acc, [key, value]) => {
    if (value === undefined) return acc;
    Object.assign(acc, { [key]: value });
    return acc;
  }, {});
  const response = await CnbApiClient.getInstance().request<Response>(
    'POST',
    '/groups',
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

export type ListGroupsParams = operations['ListTopGroups']['parameters']['query'];

export type ListSubGroupsParams = operations['ListGroups']['parameters']['query'];

export type Group = definitions['dto.OrganizationAccess'];

export type CreateGroupParams = definitions['dto.CreateGroupReq'];
