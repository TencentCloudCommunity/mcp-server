import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";

import { addTdPrefix, getSuffixByFramework, type Framework } from "@tdesign-mcp-server/common";

import { checkCompsExistence, loadFileContent } from "../helpers";
import { ComponentsSchema } from "../prompts";

export default function getCompDocs(server: McpServer) {
  server.tool(
    "get-component-docs",
    "获取 TDesign 组件的文档，适用场景：代码生成和代码转换，支持多个组件同时查询",
    ComponentsSchema.shape,
    async ({ names, framework }) => {
      const compDocs = await loadCompDocs(names, framework);

      return {
        content: [
          {
            type: "text",
            text: JSON.stringify(compDocs)
          }
        ]
      };
    }
  );
}

async function loadCompDocs(names: string[], framework: Framework) {
  const { exist, nonExist } = await checkCompsExistence(names, framework);
  const docsUrl = `${addTdPrefix(framework)}`;
  const existResults: Record<string, any> = {};

  for (const name of Object.keys(exist)) {
    const parentComponent = exist[name];
    const apiUrl = `${docsUrl}/${parentComponent}/api.md`;
    const demoUrl = `${docsUrl}/${parentComponent}/Demo.${getSuffixByFramework(framework)}`;

    const apiContent = await loadFileContent(apiUrl);
    const demoContent = await loadFileContent(demoUrl);
    const cleanedDemoContent = demoContent.replace(/\s*\n\s*/g, "").replace(/"/g, "'");

    existResults[name] = {
      api: apiContent,
      demo: cleanedDemoContent
    };
  }

  return { ...nonExist, ...existResults };
}
