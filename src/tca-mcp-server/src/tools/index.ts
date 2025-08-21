import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import TcaClient from "../api/client.js";
import registerStartScan from "./scan.js";
import registerJob from "./job.js";
import registerIssueReport from "./issue.js";

export function registerTools(server: McpServer) {
    TcaClient.initTcaClient({
        "token": process.env.TCA_TOKEN || "",
        "userName": process.env.TCA_USER_NAME || "",
        "baseUrl": process.env.TCA_BASE_URL || "https://tca.tencent.com/"
    })
    registerStartScan(server);
    registerJob(server);
    registerIssueReport(server);
}
