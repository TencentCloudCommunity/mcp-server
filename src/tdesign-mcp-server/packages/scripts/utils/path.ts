import { existsSync, promises as fs } from "fs";
import path from "path";

import {
  addTdPrefix,
  FRAMEWORKS,
  getDirname,
  getPlatformByFramework,
  isMonorepo,
  type Framework
} from "@tdesign-mcp-server/common";

const __dirname = getDirname(import.meta.url);

/* 所有 tdesign-* 系列仓库的公共父目录（默认各端框架与当前项目在同一目录） */
export const TD_REPOS_ROOT = path.join(__dirname, "../../../../");

export const TD_DOCS_OUTPUT_DIR = path.join(__dirname, "../../docs");

const joinTdDir = (...args: string[]) => path.join(TD_REPOS_ROOT, ...args);

export const TD_COMMON_DIR = joinTdDir(addTdPrefix("common"));

export const TD_COMP_MAP_PATH = `${TD_COMMON_DIR}/js/components.ts`;

export const getDocsOutputDir = async (framework: Framework, compName?: string) => {
  const baseDir = path.join(TD_DOCS_OUTPUT_DIR, addTdPrefix(framework));
  const outputDir = compName ? path.join(baseDir, compName) : baseDir;
  if (!existsSync(outputDir)) {
    await fs.mkdir(outputDir, { recursive: true });
  }
  return outputDir;
};

export const getDemoTemplateDir = (framework: Framework) => {
  const platform = getPlatformByFramework(framework);
  return `${TD_COMMON_DIR}/docs/${platform}/api`;
};

const generateDirMap = (subPathCallback: (fw: Framework) => string[]) => {
  return Object.fromEntries(FRAMEWORKS.map((fw) => [fw, joinTdDir(addTdPrefix(fw), ...subPathCallback(fw))]));
};

const COMPONENT_DIR_PATH = generateDirMap((fw) => (isMonorepo(fw) ? ["packages", "components"] : ["src"]));

export const getComponentDir = (framework: Framework, componentName: string) => {
  return path.join(COMPONENT_DIR_PATH[framework], componentName);
};

const PACKAGE_JSON_PATH = generateDirMap((fw) =>
  isMonorepo(fw) ? ["packages", addTdPrefix(fw), "package.json"] : ["package.json"]
);

export const getPackageJsonPath = (framework: Framework) => {
  return path.join(PACKAGE_JSON_PATH[framework]);
};
