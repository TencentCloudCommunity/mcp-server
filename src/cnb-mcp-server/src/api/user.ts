import CnbApiClient from './client.js';

import type { definitions } from '../schema.js';

export async function getUser(): Promise<UserInfo> {
  const cnbInst = CnbApiClient.getInstance();
  return cnbInst.request<UserInfo>('GET', '/user');
}

export type UserInfo = definitions['dto.UsersResultForSelf'];
