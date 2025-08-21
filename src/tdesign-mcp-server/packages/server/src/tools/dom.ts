import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";

import { addTdPrefix, Framework, getPlatformByFramework } from "@tdesign-mcp-server/common";

import { checkCompsExistence, loadFileContent } from "../helpers";
import { ChangelogsSchema } from "../prompts";

export default function getCompDom(server: McpServer) {
  server.tool(
    "get-component-dom",
    "获取 TDesign 组件的 DOM 结构，适用场景：转换用户自定义的 CSS 样式, 支持多个组件同时查询",
    ChangelogsSchema.shape,
    async ({ names, framework }) => {
      const compDom = await loadCompDom(names, framework);
      return {
        content: [
          {
            type: "text",
            text: JSON.stringify(compDom)
          }
        ]
      };
    }
  );
}

async function loadCompDom(names: string[], framework: Framework) {
  const { exist, nonExist } = await checkCompsExistence(names, framework);
  const domUrl = `${addTdPrefix(getPlatformByFramework(framework))}-dom`;
  const existResults: Record<string, any> = {};

  for (const name of Object.keys(exist)) {
    const parentComponent = exist[name];
    const domFileUrl = `${domUrl}/${parentComponent}.html`;

    const domContent = await loadFileContent(domFileUrl);

    existResults[name] = domContent;
  }

  return { ...nonExist, ...existResults };
}
