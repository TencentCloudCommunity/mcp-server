import { execFileSync } from 'node:child_process';
import { cpSync, existsSync, mkdirSync, mkdtempSync, readdirSync, readFileSync, rmSync } from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export const SKILLS_ROOT = path.resolve(__dirname, '..');
export const DIST_DIR = path.join(SKILLS_ROOT, 'dist');
export const COMMON_REFERENCES_DIR = path.join(SKILLS_ROOT, 'references', 'common');
export const SKILLS = [
  'tencent-pg-inspection',
  'tencent-pg-slowquery-diagnosis',
  'tencent-pg-ops-troubleshooter',
];

function parseArgs(argv) {
  const result = {};
  for (let index = 0; index < argv.length; index += 1) {
    const token = argv[index];
    if (token === '--skill' || token === '--name') {
      result.skill = argv[index + 1];
      index += 1;
      continue;
    }
    if (token === '--version') {
      result.version = argv[index + 1];
      index += 1;
      continue;
    }
    if (!result.skill && !token.startsWith('--')) {
      result.skill = token;
    }
  }
  return result;
}

export function resolveVersion(inputVersion) {
  const raw = inputVersion || process.env.SKILL_VERSION || readRootVersion();
  const normalized = String(raw).trim().replace(/^v/i, '');
  if (!/^\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?$/.test(normalized)) {
    throw new Error(`invalid skill version: ${raw}`);
  }
  return normalized;
}

function readRootVersion() {
  const packageJsonPath = path.resolve(SKILLS_ROOT, '..', 'package.json');
  const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf8'));
  if (!packageJson.version) {
    throw new Error(`missing version in ${packageJsonPath}`);
  }
  return packageJson.version;
}

function ensureZipCommand() {
  try {
    execFileSync('zip', ['-v'], { stdio: 'ignore' });
  } catch {
    throw new Error('zip command not found; please install zip before packaging skills');
  }
}

function hasPackableEntries(dirPath) {
  if (!existsSync(dirPath)) {
    return false;
  }
  const entries = readdirSync(dirPath, { withFileTypes: true });
  return entries.some((entry) => !entry.name.startsWith('.'));
}

export function stageCommonReferences(targetReferencesDir) {
  if (!existsSync(COMMON_REFERENCES_DIR)) {
    return;
  }

  mkdirSync(targetReferencesDir, { recursive: true });
  cpSync(COMMON_REFERENCES_DIR, path.join(targetReferencesDir, 'common'), {
    recursive: true,
  });
}

export function packageSkill(skillName, inputVersion, outputDir = DIST_DIR) {
  if (!SKILLS.includes(skillName)) {
    throw new Error(`unsupported skill: ${skillName}`);
  }

  ensureZipCommand();

  const version = resolveVersion(inputVersion);
  const skillDir = path.join(SKILLS_ROOT, skillName);
  const skillFile = path.join(skillDir, 'SKILL.md');
  const referencesDir = path.join(skillDir, 'references');
  const assetsDir = path.join(skillDir, 'assets');

  if (!existsSync(skillFile)) {
    throw new Error(`missing SKILL.md for ${skillName}`);
  }
  if (!existsSync(referencesDir)) {
    throw new Error(`missing references directory for ${skillName}`);
  }

  mkdirSync(outputDir, { recursive: true });

  const outFile = path.join(outputDir, `${skillName}-v${version}.zip`);
  const entries = ['SKILL.md', 'references'];
  const stageDir = mkdtempSync(path.join(os.tmpdir(), `${skillName}-`));

  try {
    cpSync(skillFile, path.join(stageDir, 'SKILL.md'));
    cpSync(referencesDir, path.join(stageDir, 'references'), { recursive: true });
    stageCommonReferences(path.join(stageDir, 'references'));

    if (hasPackableEntries(assetsDir)) {
      cpSync(assetsDir, path.join(stageDir, 'assets'), { recursive: true });
      entries.push('assets');
    }

    rmSync(outFile, { force: true });
    execFileSync('zip', ['-rq', outFile, ...entries], { cwd: stageDir, stdio: 'inherit' });

    return outFile;
  } finally {
    rmSync(stageDir, { recursive: true, force: true });
  }
}

const isDirectRun = process.argv[1] && path.resolve(process.argv[1]) === __filename;

if (isDirectRun) {
  const { skill, version } = parseArgs(process.argv.slice(2));
  if (!skill) {
    console.error('Usage: node ./scripts/package-skill.mjs --skill <skill-name> [--version <x.y.z>]');
    process.exit(1);
  }

  const outFile = packageSkill(skill, version);
  console.log(`packaged ${skill}: ${outFile}`);
}
