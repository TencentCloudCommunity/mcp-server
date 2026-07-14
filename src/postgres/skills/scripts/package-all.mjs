import { cpSync, existsSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { execFileSync } from 'node:child_process';
import { COMMON_REFERENCES_DIR, DIST_DIR, SKILLS, SKILLS_ROOT, packageSkill, resolveVersion } from './package-skill.mjs';

const BUNDLE_SLUG = 'tencentdb-postgresql-skill';
const BUNDLE_DISPLAY_NAME = 'TencentDB PostgreSQL Skill';
const BUNDLE_DESCRIPTION = 'Bundle entry for TencentDB PostgreSQL skills. These skills call Tencent Cloud PostgreSQL OpenAPI directly and can use the full 48 aligned actions documented by each child skill under references/. Explicit confirmation is required before any write, fee-impacting, or high-risk action.';

function ensureDistDir() {
  rmSync(DIST_DIR, { recursive: true, force: true });
  mkdirSync(DIST_DIR, { recursive: true });
}

function bundleFilter(src) {
  const name = path.basename(src);
  if (name.startsWith('.')) {
    return false;
  }
  if (name === 'dist') {
    return false;
  }
  return true;
}

function parseFrontmatter(markdown) {
  const match = markdown.match(/^---\n([\s\S]*?)\n---\n?/);
  if (!match) {
    return {};
  }

  const metadata = {};
  for (const line of match[1].split('\n')) {
    const separatorIndex = line.indexOf(':');
    if (separatorIndex === -1) {
      continue;
    }

    const key = line.slice(0, separatorIndex).trim();
    let value = line.slice(separatorIndex + 1).trim();
    if (!key) {
      continue;
    }

    if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
      value = value.slice(1, -1);
    }

    metadata[key] = value;
  }

  return metadata;
}

function readSkillDescriptors() {
  return SKILLS.map((skillName) => {
    const skillFile = path.join(SKILLS_ROOT, skillName, 'SKILL.md');
    const content = readFileSync(skillFile, 'utf8');
    const meta = parseFrontmatter(content);

    return {
      skillName,
      name: meta.name || skillName,
      description: meta.description || '',
      descriptionZh: meta.description_zh || meta.description_en || meta.name || skillName,
      entryPath: `references/${skillName}/SKILL.md`,
      directoryPath: `references/${skillName}/`,
    };
  });
}

function extractReadmeSection(markdown, startHeading, endHeading) {
  const startIndex = markdown.indexOf(startHeading);
  if (startIndex === -1) {
    return '';
  }

  const endIndex = markdown.indexOf(endHeading, startIndex + startHeading.length);
  const section = endIndex === -1 ? markdown.slice(startIndex) : markdown.slice(startIndex, endIndex);
  return section.trim();
}

function renderBundleSkill(version, descriptors) {
  const sections = descriptors
    .map(
      (descriptor) => `### \`${descriptor.skillName}\`\n- 简介：${descriptor.descriptionZh}\n- 目录：\`${descriptor.directoryPath}\`\n- 入口：\`${descriptor.entryPath}\`\n- 说明：${descriptor.description || '请进入该目录查看完整 skill 说明。'}`,
    )
    .join('\n\n');
  const skillsReadme = readFileSync(path.join(SKILLS_ROOT, 'README.md'), 'utf8');
  const alignedActionsSection = extractReadmeSection(skillsReadme, '## 当前开放的 48 个对齐 Action', '## 目录结构');

  return `---
name: "${BUNDLE_DISPLAY_NAME}"
description: "${BUNDLE_DESCRIPTION}"
description_zh: "${BUNDLE_DISPLAY_NAME}"
description_en: "${BUNDLE_DISPLAY_NAME}"
version: ${version}
---

# ${BUNDLE_DISPLAY_NAME}

这是 \`${BUNDLE_DISPLAY_NAME}\` 的 bundle 入口，参考可发布 skill 包的根目录格式组织；最外层提供统一入口 \`SKILL.md\`，具体子技能位于 \`references/\` 目录。

这些子技能**直接调用腾讯云 PostgreSQL OpenAPI**，不依赖已部署的 MCP Server；并且现在允许使用当前 PostgreSQL MCP 已完成参数对齐的 **48 个 Action**。每个子技能自己的 \`@references/api_reference.md\` 都列出了完整动作目录和确认要求。

## 极简配置（面向用户）

推荐用户只准备 3 个环境变量：
- \`TENCENTCLOUD_SECRET_ID\`
- \`TENCENTCLOUD_SECRET_KEY\`
- \`TENCENTCLOUD_REGION\`

可选：临时凭证场景再补 \`TENCENTCLOUD_SESSION_TOKEN\`。

补充约定：
- 兼容读取 \`MCP_REQUEST_SECRET_ID\`、\`MCP_REQUEST_SECRET_KEY\`、\`MCP_REQUEST_SESSION_TOKEN\`、\`MCP_SECRET_ID\`、\`MCP_SECRET_KEY\`
- 地域既可以直接写 \`ap-guangzhou\`，也可以先写 \`广州\`、\`上海\`、\`成都\`、\`北京\`，执行时应先归一化为标准地域码
- 是否缺少官方 SDK 不应成为用户首次使用时的首个阻断项；若 SDK 不可用，应退回到本地生成 TC3 签名 HTTPS 请求

## 公共规则文档

- \`@references/common/region_normalization.md\`：统一维护地域别名、标准地域码、非法地域处理规则，以及官方地域查询链接
- \`@references/common/error_handling.md\`：统一维护凭证缺失、地域非法、巡检目标缺失、SDK 缺失时的错误模板、官方链接与安装指引

## 执行方式（面向 Agent）

以下内容主要给执行该 skill 的 AI / Agent 使用，不是给终端用户看的安装步骤。

1. 先根据用户请求判断任务类型，再选择最匹配的子 skill。
2. 在进入子 skill 前，先确认目标范围是否完整：至少应拿到地域与实例 ID（或可明确识别的实例名）。如果缺少这些信息，应直接回复一条包含以下内容的消息（不要使用交互式选择菜单）：(a) 控制台链接 https://console.cloud.tencent.com/postgres，让用户去查实例 ID 和地域；(b) 示例回复格式 \`ap-guangzhou postgres-abc12345\`。
3. 优先从运行时上下文读取 \`TENCENTCLOUD_SECRET_ID\`、\`TENCENTCLOUD_SECRET_KEY\`、\`TENCENTCLOUD_REGION\` 与可选 \`TENCENTCLOUD_SESSION_TOKEN\`；如果缺失，再检查兼容的 \`MCP_REQUEST_*\` / \`MCP_*\` 变量名。若仍缺失，应停止执行并使用 \`@references/common/error_handling.md\` 中的 \`missing-credentials\` 模板，附上可复制示例以及腾讯云官方 API 密钥/地域链接。
4. 在真正调用 OpenAPI 之前，先按 \`@references/common/region_normalization.md\` 处理地域；如果用户输入无法安全归一化，应使用 \`@references/common/error_handling.md\` 中的 \`invalid-region\` 模板，并附上官方地域查询链接。
5. 进入对应子目录，继续阅读该目录下的 \`SKILL.md\` 和 \`@references/api_reference.md\`，按子 skill 的步骤执行。
6. 调用腾讯云 PostgreSQL OpenAPI 时，优先从读类动作开始采集证据；如需执行写类、费用类或高风险动作，必须先说明影响面并获得明确确认。
7. 优先使用官方 SDK；如果 SDK 不可用，不要仅因这一点中止执行，而应退回到本地生成 TC3 签名 HTTPS 请求。若用户仍希望安装 SDK，则按 \`@references/common/error_handling.md\` 中的 \`missing-sdk\` 模板给出安装命令并询问是否执行。

## 包含的子技能

${sections}

${alignedActionsSection ? `${alignedActionsSection}

` : ''}## 目录约定

- 根目录入口：\`SKILL.md\`
- 根目录元数据：\`_meta.json\`
- 公共规则目录：\`references/common/\`
- 子技能目录：\`references/<skill-name>/\`
- 子技能入口：\`references/<skill-name>/SKILL.md\`
`;
}

function renderBundleMeta(version) {
  return `${JSON.stringify(
    {
      slug: BUNDLE_SLUG,
      version,
      publishedAt: Date.now(),
    },
    null,
    2,
  )}\n`;
}

function stageExpandedSkills(stageDir) {
  const referencesDir = path.join(stageDir, 'references');
  mkdirSync(referencesDir, { recursive: true });

  if (existsSync(COMMON_REFERENCES_DIR)) {
    cpSync(COMMON_REFERENCES_DIR, path.join(referencesDir, 'common'), {
      recursive: true,
    });
  }

  for (const skillName of SKILLS) {
    const sourceDir = path.join(SKILLS_ROOT, skillName);
    const targetDir = path.join(referencesDir, skillName);
    cpSync(sourceDir, targetDir, {
      recursive: true,
      filter: bundleFilter,
    });
  }
}

function createBundle(version) {
  const stageDir = mkdtempSync(path.join(os.tmpdir(), `${BUNDLE_SLUG}-`));
  const bundlePath = path.join(DIST_DIR, `${BUNDLE_SLUG}-v${version}.zip`);
  const descriptors = readSkillDescriptors();

  try {
    writeFileSync(path.join(stageDir, 'SKILL.md'), renderBundleSkill(version, descriptors), 'utf8');
    writeFileSync(path.join(stageDir, '_meta.json'), renderBundleMeta(version), 'utf8');

    stageExpandedSkills(stageDir);

    rmSync(bundlePath, { force: true });
    execFileSync('zip', ['-rq', bundlePath, 'SKILL.md', '_meta.json', 'references'], {
      cwd: stageDir,
      stdio: 'inherit',
    });

    return bundlePath;
  } finally {
    rmSync(stageDir, { recursive: true, force: true });
  }
}

function createSinglePackages(version) {
  return SKILLS.map((skillName) => packageSkill(skillName, version, DIST_DIR));
}

function main() {
  const versionArg = process.argv.slice(2).find((token) => !token.startsWith('--'));
  const version = resolveVersion(versionArg);

  ensureDistDir();

  const packagedFiles = createSinglePackages(version);
  const bundlePath = createBundle(version);

  const output = [
    ...packagedFiles.map((filePath) => `- ${path.basename(filePath)}`),
    `- ${path.basename(bundlePath)}`,
  ].join('\n');

  console.log(`generated skill release assets (v${version}):\n${output}`);
}

main();
