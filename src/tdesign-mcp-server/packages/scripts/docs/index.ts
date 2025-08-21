import { existsSync, promises as fs, readdirSync } from "fs";
import path from "path";

import {
  API_RULE_MAP,
  cleanText,
  DEMO_RULE_MAP,
  FRAMEWORKS,
  getSuffixByFramework,
  type Framework
} from "@tdesign-mcp-server/common";

import {
  getComponentDir,
  getDemoTemplateDir,
  getDocsOutputDir,
  TD_COMP_MAP_PATH,
  TD_DOCS_OUTPUT_DIR
} from "../utils/path";

(async function main() {
  await extractCompMap();
  await extractAllDocs();
})();

async function extractCompMap() {
  const { NON_PASCAL_CASE_NAMES, WEB_COMPONENT_MAP } = await import(TD_COMP_MAP_PATH);
  await fs.mkdir(TD_DOCS_OUTPUT_DIR, { recursive: true });
  const compMap = {
    alias: NON_PASCAL_CASE_NAMES,
    components: WEB_COMPONENT_MAP
  };
  const mapOutputPath = path.join(TD_DOCS_OUTPUT_DIR, "web-components.json");
  await fs.writeFile(mapOutputPath, JSON.stringify(compMap, null, 2), "utf-8");
  console.log(`✓ Extract map successfully`);
}

async function extractAllDocs() {
  for (const framework of FRAMEWORKS) {
    console.log(`--------(${framework})--------`);
    await extractDocs(framework);
    console.log(`✓ Extract docs successfully`);
  }
}

const EXCLUDED_COMPONENTS = [/^chat/]; // 排除部分组件，需要单独处理

function isExcludedComponent(compName: string): boolean {
  return EXCLUDED_COMPONENTS.some((pattern) =>
    typeof pattern === "string" ? pattern === compName : pattern.test(compName)
  );
}

async function extractDocs(framework: Framework) {
  const demoTemplateDir = await fs.readdir(getDemoTemplateDir(framework));
  const mdFiles = demoTemplateDir.filter((file) => file.endsWith(".md") && !file.includes("en-US"));

  for (const file of mdFiles) {
    const compName = file.split(".")[0];
    const compDir = getComponentDir(framework, compName);

    if (isExcludedComponent(compName)) {
      // console.warn(`✗ Component excluded: ${compName}`);
      continue;
    } else if (!existsSync(compDir)) {
      console.warn(`✗ Component not found: ${compName}`);
      continue;
    }

    await generateApi(framework, compName);
    await generateDemo(framework, compName);
    await generateCompList(framework, compName);
  }
}

/* ---------------------------------------- */

async function generateApi(framework: Framework, compName: string) {
  const componentPath = getComponentDir(framework, compName);
  const apiPath = `${componentPath}/${compName}.md`;

  let apiContent = await fs.readFile(apiPath, "utf-8");
  apiContent = cleanText(apiContent, API_RULE_MAP);

  const outputDir = await getDocsOutputDir(framework, compName);
  await fs.writeFile(path.join(outputDir, "api.md"), apiContent, "utf-8");
}

async function generateDemo(framework: Framework, compName: string) {
  const exampleDir = path.join(getComponentDir(framework, compName), "_example");
  // 取基础 Demo 即可
  const baseFileName = "base";

  if (!existsSync(exampleDir)) return null;

  const files = readdirSync(exampleDir);
  let baseDemo = files.find((file) => {
    const [name] = file.split(".");
    return name === baseFileName;
  });
  if (!baseDemo) baseDemo = files[0];

  let demoContent = await fs.readFile(path.join(exampleDir, baseDemo), "utf-8");
  demoContent = cleanText(demoContent, DEMO_RULE_MAP);

  const suffix = getSuffixByFramework(framework);
  const outputDir = await getDocsOutputDir(framework, compName);
  await fs.writeFile(path.join(outputDir, `Demo.${suffix}`), demoContent, "utf-8");
}

async function generateCompList(framework: Framework, compName: string) {
  const componentDir = getComponentDir(framework, compName);
  if (!existsSync(componentDir)) return; // 避免某些库没有实现对应的组件

  const demoTemplatePath = `${getDemoTemplateDir(framework)}/${compName}.md`;
  const demoTemplateContent = await fs.readFile(demoTemplatePath, "utf-8");

  const metaReg = DEMO_RULE_MAP.meta[0];
  const metaMatch = demoTemplateContent.match(metaReg);
  const meta = metaMatch ? metaMatch[1] : "";
  const metaObj = Object.fromEntries(
    meta.split("\n").map((line) => {
      const [key, value] = line.split(": ");
      return [key.trim(), value.trim()];
    })
  );

  const excludedSpine = ["explain"];
  if (excludedSpine.includes(metaObj.spline)) return;

  const splineType = metaObj.spline;
  const compMeta = {
    [compName]: metaObj.description
  };

  let compListPath = await getDocsOutputDir(framework);
  compListPath = path.join(compListPath, "index.json");
  if (!existsSync(compListPath)) {
    await fs.writeFile(compListPath, JSON.stringify({}, null, 2), "utf-8");
  }

  let compListContent = JSON.parse(await fs.readFile(compListPath, "utf-8"));

  compListContent[splineType] = {
    ...(compListContent[splineType] || {}),
    ...compMeta
  };

  await fs.writeFile(compListPath, JSON.stringify(compListContent, null, 2), "utf-8");
}
