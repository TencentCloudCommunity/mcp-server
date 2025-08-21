
import { exec } from 'child_process';
import { promisify } from 'util';
import { readFile } from 'fs/promises';
import { join } from 'path';
import { ErrorCode, McpError } from '@modelcontextprotocol/sdk/types.js';
import { CodingDevOpsConfig } from '../../config/environment.js';

const execAsync = promisify(exec);

async function execInDirectory(command: string, cwd: string): Promise<string> {
  const { stdout } = await execAsync(command, { cwd });
  return stdout.trim();
}

async function isGitRepository(workingDirectory: string): Promise<boolean> {
  try {
    await execInDirectory('git rev-parse --is-inside-work-tree', workingDirectory);
    return true;
  } catch (error: unknown) {
    return false;
  }
}

export async function getCurrentBranch(workingDirectory: string): Promise<string> {
  try {
    // 首先尝试使用 git rev-parse
    try {
      return await execInDirectory('git rev-parse --abbrev-ref HEAD', workingDirectory);
    } catch (error: unknown) {
      // 如果失败，尝试使用 git branch
      return await execInDirectory('git branch --show-current', workingDirectory);
    }
  } catch (error: unknown) {
    // 如果两种方法都失败，提供更详细的错误信息
    throw new McpError(
      ErrorCode.InvalidParams,
      `获取当前分支失败。请确保：\n` +
      `1. 目录 "${workingDirectory}" 是一个 Git 仓库\n` +
      '2. 至少有一次提交记录（使用 git commit）\n' +
      '3. 当前不在分离头指针（detached HEAD）状态'
    );
  }
}

export async function getRemoteUrl(workingDirectory: string): Promise<string> {
  try {
    // 首先尝试使用 git remote get-url
    try {
      return await execInDirectory('git remote get-url origin', workingDirectory);
    } catch (error: unknown) {
      // 如果失败，尝试使用 git config
      return await execInDirectory('git config --get remote.origin.url', workingDirectory);
    }
  } catch (error: unknown) {
    // 如果两种方法都失败，提供更详细的错误信息
    throw new McpError(
      ErrorCode.InvalidParams, 
      `获取远程仓库 URL 失败。请确保：\n` +
      `1. 目录 "${workingDirectory}" 是一个 Git 仓库\n` +
      `2. 已经配置了远程仓库（使用 git remote -v 查看）\n` +
      '3. 远程仓库名称为 origin'
    );
  }
}

export async function getGitConfig(workingDirectory: string): Promise<{ branch: string; remoteUrl: string }> {
  // 首先检查是否在 Git 仓库中
  const isGitRepo = await isGitRepository(workingDirectory);
  if (!isGitRepo) {
    throw new McpError(
      ErrorCode.InvalidParams,
      `目录 "${workingDirectory}" 不是 Git 仓库。请确保：\n` +
      '1. 在正确的项目目录下执行命令\n' +
      '2. 项目已经初始化为 Git 仓库（使用 git init）'
    );
  }

  const [branch, remoteUrl] = await Promise.all([
    getCurrentBranch(workingDirectory),
    getRemoteUrl(workingDirectory)
  ]);

  return {
    branch,
    remoteUrl
  };
}

export function parseRemoteUrl(url: string): { projectName?: string; depotName?: string; depotPath?: string } {
  // 处理 SSH 和 HTTPS 格式的 URL
  // 例如：git@e.coding.net:team-name/project-name/repo-name.git
  // 或 https://e.coding.net/team-name/project-name/repo-name.git
  const match = url.match(/[:/]([\w-]+)\/([\w-]+)\/([\w-]+)(?:\.git)?$/);
  if (!match) {
    return {};
  }

  const [, teamName, projectName, depotName] = match;
  return {
    projectName,
    depotName,
    depotPath: `${teamName}/${projectName}/${depotName}`
  };
}