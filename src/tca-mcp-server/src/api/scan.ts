import TcaClient, { CommonParams } from "./client.js"

export interface startScanParams {
    "token": string,
    "userName": string,
    "orgSid": string,
    "projectName": string,
    "scanName": string,
}

export interface startScanData {
    "incr_scan": boolean,
    "created_from": string,
    "force_create": boolean,
}

export async function startScan(incr_scan: boolean, force_create: boolean) {
    const method = "POST"
    const pathParams = TcaClient.getcommonParams()
    const path = `/server/main/api/orgs/${pathParams.orgSid}/teams/${pathParams.teamName}/repos/${pathParams.repoId}/projects/${pathParams.projectId}/scans/create/`
    const clientInst = TcaClient.getInstance()
    const data = {
        "incr_scan": incr_scan,
        "created_from": "tca-mcp-server",
        "force_create": force_create,
    } as startScanData
    return clientInst.request(path, method, data)
}
