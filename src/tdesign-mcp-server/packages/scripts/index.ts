import { exec } from "child_process";
import { existsSync } from "fs";
import path from "path";

import { promisify } from "util";
const execPromise = promisify(exec);

import { addTdPrefix, getDirname } from "@tdesign-mcp-server/common";
import { TD_REPOS_ROOT } from "./utils/path";

const __dirname = getDirname(import.meta.url);

const TENCENT_GITHUB_URL = "https://github.com/Tencent/";
const TD_REQUIRED_REPOS = ["common", "react", "vue", "vue-next"] as const;

(async function main() {
  await prepareTDesignRepos();
  await extractDocs();
})();

async function gitClone(repoUrl: string) {
  console.log(`Cloning repository from ${repoUrl}...`);
  await execPromise(`cd ${TD_REPOS_ROOT} && git clone ${repoUrl}`);
  console.log("Repository cloned successfully!\n");
}

async function prepareTDesignRepos() {
  for (const repo of TD_REQUIRED_REPOS) {
    const localRepoPath = path.join(TD_REPOS_ROOT, addTdPrefix(repo));
    if (existsSync(localRepoPath)) continue;
    const repoUrl = `${TENCENT_GITHUB_URL}${addTdPrefix(repo)}`;
    await gitClone(repoUrl);
  }
}

async function executeTypeScript(scriptPath: string) {
  const { stdout, stderr } = await execPromise(`tsx ${scriptPath}`);
  if (stdout) console.log(stdout);
  if (stderr) console.error(stderr);
}

async function extractDocs() {
  await executeTypeScript(path.join(__dirname, "docs/index.ts"));
  await executeTypeScript(path.join(__dirname, "dom/index.ts"));
}
