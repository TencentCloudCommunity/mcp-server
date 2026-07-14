import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import path from 'node:path';
import { COMMON_REFERENCES_DIR, SKILLS, SKILLS_ROOT } from './package-skill.mjs';

function parseFrontmatter(content) {
  const match = content.match(/^---\n([\s\S]*?)\n---/);
  if (!match) {
    throw new Error('missing YAML frontmatter block');
  }

  const fields = {};
  for (const rawLine of match[1].split('\n')) {
    const line = rawLine.trim();
    if (!line || line.startsWith('#')) {
      continue;
    }
    const separatorIndex = line.indexOf(':');
    if (separatorIndex <= 0) {
      continue;
    }
    const key = line.slice(0, separatorIndex).trim();
    const value = line.slice(separatorIndex + 1).trim();
    fields[key] = value;
  }
  return fields;
}

function walkForZipFiles(dirPath, collected = []) {
  if (!existsSync(dirPath)) {
    return collected;
  }

  for (const entry of readdirSync(dirPath, { withFileTypes: true })) {
    if (entry.name.startsWith('.')) {
      continue;
    }
    const fullPath = path.join(dirPath, entry.name);
    if (entry.isDirectory()) {
      walkForZipFiles(fullPath, collected);
      continue;
    }
    if (entry.name.endsWith('.zip')) {
      collected.push(fullPath);
    }
  }
  return collected;
}

function verifySkill(skillName) {
  const skillDir = path.join(SKILLS_ROOT, skillName);
  const skillFile = path.join(skillDir, 'SKILL.md');
  const referencesDir = path.join(skillDir, 'references');
  const apiReferencePath = path.join(referencesDir, 'api_reference.md');
  const assetsDir = path.join(skillDir, 'assets');
  const distDir = path.join(skillDir, 'dist');

  if (!existsSync(skillDir) || !statSync(skillDir).isDirectory()) {
    throw new Error(`${skillName}: missing skill directory`);
  }
  if (!existsSync(skillFile)) {
    throw new Error(`${skillName}: missing SKILL.md`);
  }
  if (!existsSync(referencesDir) || !statSync(referencesDir).isDirectory()) {
    throw new Error(`${skillName}: missing references directory`);
  }
  if (!existsSync(apiReferencePath)) {
    throw new Error(`${skillName}: missing references/api_reference.md`);
  }
  if (existsSync(distDir)) {
    throw new Error(`${skillName}: dist directory must not exist inside source skill directory`);
  }
  if (existsSync(assetsDir) && !statSync(assetsDir).isDirectory()) {
    throw new Error(`${skillName}: assets exists but is not a directory`);
  }

  const content = readFileSync(skillFile, 'utf8');
  const fields = parseFrontmatter(content);

  if (fields.name !== skillName) {
    throw new Error(`${skillName}: frontmatter name must equal directory name`);
  }
  for (const requiredField of ['description', 'description_zh', 'description_en']) {
    if (!fields[requiredField]) {
      throw new Error(`${skillName}: missing frontmatter field ${requiredField}`);
    }
  }
  if (!content.includes('@references/api_reference.md')) {
    throw new Error(`${skillName}: SKILL.md must reference @references/api_reference.md`);
  }

  const zipFiles = walkForZipFiles(skillDir);
  if (zipFiles.length > 0) {
    throw new Error(`${skillName}: packaged zip assets must not be stored in source directories: ${zipFiles.join(', ')}`);
  }
}

function verifyCommonReferences() {
  const requiredFiles = ['region_normalization.md', 'error_handling.md'];
  for (const fileName of requiredFiles) {
    const targetPath = path.join(COMMON_REFERENCES_DIR, fileName);
    if (!existsSync(targetPath)) {
      throw new Error(`missing shared reference: ${targetPath}`);
    }
  }
}

function main() {
  verifyCommonReferences();
  for (const skillName of SKILLS) {
    verifySkill(skillName);
    console.log(`verified ${skillName}`);
  }
  console.log(`all ${SKILLS.length} PostgreSQL companion skills verified`);
}

main();
