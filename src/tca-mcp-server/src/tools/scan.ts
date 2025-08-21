import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";
import { startScan } from "../api/scan.js";
import { formatTextToolResult, formatToolError } from "../utils/resp.js";
import TcaClient from "../api/client.js";


export default function registerStartScan(server: McpServer) {
    server.tool(
        'start_scan',
        `启动tca代码分析 参数如下：
        -mcpConfigFile: mcp配置文件路径，当前项目根目录查找 tca-mcp.ini（绝对路径）
        -incrScan: 是否增量扫描
        -forceCreate: 如果已经存在扫描结果，是否强制启动新扫描`,
        {
            mcpConfigFile: z.string(),
            incrScan: z.boolean().default(false).describe("是否增量扫描"),
            forceCreate: z.boolean().default(false).describe("如果已经存在扫描结果，是否强制启动新扫描"),
        },
        async ({ mcpConfigFile, incrScan, forceCreate }) => {
            await TcaClient.initConfig(mcpConfigFile)
            try {
                const res = await startScan(incrScan, forceCreate)
                return formatTextToolResult(JSON.stringify(res, null, 2), 'start tca scan');
            } catch (error) {
                return formatToolError(error, 'start tca scan error!')
            }
        }
    );
}
