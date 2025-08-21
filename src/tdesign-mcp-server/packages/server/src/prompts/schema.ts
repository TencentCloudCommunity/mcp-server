import { z } from "zod";
import { FRAMEWORKS } from "@tdesign-mcp-server/common";

export const FrameworkSchema = z.object({
  framework: z.enum(FRAMEWORKS).describe("如果用户没有指定，根据 package.json 中的依赖判断")
});

export const ComponentsSchema = FrameworkSchema.extend({
  names: z.array(z.string()).describe("可用的组件名称，支持批量查询")
});

export const ChangelogsSchema = ComponentsSchema.extend({
  version: z.string().optional().describe("如果用户没有指定，查询 package.json 中的依赖版本号")
});
