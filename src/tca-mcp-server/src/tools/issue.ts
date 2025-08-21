import { z } from "zod";
import TcaClient from "../api/client.js";
import { format } from "path";
import { formatTextToolResult, formatToolError } from "../utils/resp.js";
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { getIssueList } from "../api/issue.js";

export default function registerIssueReport(server: McpServer) {
    server.tool(
        'tca_issue_list',
        `获取当前项目未解决的问题列表, 参数如下：
        -mcpConfigFile: mcp配置文件路径，当前项目根目录查找 tca-mcp.ini（绝对路径）
        -state: 查询任务状态, （0为等待中，1为执行中，2为关闭，3为入库中，4正在初始化, 5初始化完成 可多选，格式为1,2,3）
        -severity: 问题严重程度, （0为低，1为中，2为高）
        -offset: 偏移量
        -limit: 限制数量
        `,
        {
            mcpConfigFile: z.string(),
            state: z.string().default("0,1,2,3,4,5"),
            severity: z.string().default("0,1,2"),
            offset: z.number().default(0),
            limit: z.number().default(10),
        },
        async ({ mcpConfigFile, state, severity, offset, limit }) => {
            await TcaClient.initConfig(mcpConfigFile)
            try {
                const res = await getIssueList(state, severity, offset, limit)
                return formatTextToolResult(JSON.stringify(res, null, 2), 'get issue list');
            } catch (error) {
                return formatToolError(error, 'get issue list error')
            }
        }
    )

    server.tool(
        'tca_issue_report',
        `tca报告链接, 参数如下：
        -mcpConfigFile: mcp配置文件路径，当前项目根目录查找 tca-mcp.ini（绝对路径）
        -jobId: 任务id, 启动分析任务后，从返回结果中获取`,
        {
            mcpConfigFile: z.string(),
            jobId: z.string().describe("jobId"),
        },
        async ({ mcpConfigFile, jobId }) => {
            await TcaClient.initConfig(mcpConfigFile)
            const baseUrl = TcaClient.getInstance().baseUrl
            const params = TcaClient.getcommonParams()
            const path = `t/${params.orgSid}/p/${params.teamName}/repos/${params.repoId}/projects/${params.projectId}/jobs/${jobId}/result`
            return formatTextToolResult(baseUrl + path, 'tca issue report link');
        }
    );
}
