import { readdirSync } from "fs";
import path from "path";

import { cleanup, render } from "@testing-library/react";
import React from "react";
import { expect, test } from "vitest";

import { getDirname } from "@tdesign-mcp-server/common";
import { getDocsOutputDir } from "../../../utils/path";

const __dirname = getDirname(import.meta.url);
const CUSTOM_TEST_DIR = path.join(__dirname, "./custom");

const EXCLUDE_COMPS = ["config-provider", "watermark"];

(async function main() {
  console.log(`--------(web)--------`);
  await runTests();
  console.log(`✓ Generate snapshots successfully`);
})();

async function getCompsWithCustomTest() {
  const testFiles = readdirSync(CUSTOM_TEST_DIR).filter((file) => file.endsWith(".tsx"));
  const compNames = testFiles.map((file) => file.replace(".tsx", ""));
  return compNames;
}

async function runTests() {
  // 拥有自定义测试的组件
  const compsWithCustomTest = await getCompsWithCustomTest();
  compsWithCustomTest.forEach(async (compName) => {
    const demoPath = path.join(CUSTOM_TEST_DIR, `${compName}.tsx`);
    test(`Custom Test [${compName}]`, async () => {
      const module = await import(demoPath);
      await module.default();
      cleanup();
    });
  });

  // 使用默认测试的组件
  const reactDocs = await getDocsOutputDir("react");
  const compDirs = readdirSync(reactDocs).filter((file) => !file.endsWith(".json") && !EXCLUDE_COMPS.includes(file));
  const compsWithDefaultTest = compDirs.filter((compName) => {
    return !compsWithCustomTest.includes(compName);
  });
  compsWithDefaultTest.forEach(async (compName) => {
    const demoPath = path.join(reactDocs, compName, "Demo.tsx");
    test(`Default Test [${compName}]`, async () => {
      const module = await import(demoPath);
      const Component = module?.default;
      if (!Component) {
        console.error(`Component not found in ${demoPath}`);
        return;
      }
      defaultTest(Component);
    });
  });
}

function defaultTest(Component: React.ElementType) {
  const { container } = render(React.createElement(Component));
  expect(container).toMatchSnapshot();
  cleanup();
}
