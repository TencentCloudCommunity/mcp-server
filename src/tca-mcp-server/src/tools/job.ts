import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";
import TcaClient from "../api/client.js";
import { formatTextToolResult, formatToolError } from "../utils/resp.js";
import { getJobDetail, getJobList } from "../api/job.js";

export default function registerJob(server: McpServer) {
    server.tool(
        'job_detail',
        `获取任务详情, 参数如下：
        -mcpConfigFile: mcp配置文件路径，当前项目根目录查找 tca-mcp.ini（绝对路径）
        -jobId: 任务id`,
        {
            mcpConfigFile: z.string(),
            jobId: z.string().describe("任务id")
        },
        async ({ mcpConfigFile, jobId }) => {
            await TcaClient.initConfig(mcpConfigFile)
            try {
                const res = await getJobDetail(jobId)
                return formatTextToolResult(JSON.stringify(res, null, 2), 'get job detail');
            } catch (error) {
                return formatToolError(error, 'get job detail error!')
            }
        }
    )
    
    server.tool(
        'job_list', 
        `获取当前任务列表, 参数如下：
        -mcpConfigFile: mcp配置文件路径，当前项目根目录查找 tca-mcp.ini（绝对路径）
        -offset: 偏移量
        -limit: 限制数量
        -state: 查询任务状态, （0为等待中，1为执行中，2为关闭，3为入库中，4正在初始化, 5初始化完成 可多选，格式为1,2,3）`, 
        {
            mcpConfigFile: z.string(),
            offset: z.number().default(0).describe("偏移量"),
            limit: z.number().default(10).describe("限制数量"),
            state: z.string().default("0,1,2,3,4,5").describe("查询任务状态,  （0为等待中，1为执行中，2为关闭，3为入库中，4正在初始化, 5初始化完成 可多选，格式为1,2,3）")
        },
        async ({ mcpConfigFile, offset, limit, state }) => {
            await TcaClient.initConfig(mcpConfigFile)
            try {
                const res = await getJobList(offset, limit, state)
                return formatTextToolResult(JSON.stringify(res, null, 2), 'get job list');
            } catch (error) {
                return formatToolError(error, 'get job list error!')
            }
        }
    )
}
