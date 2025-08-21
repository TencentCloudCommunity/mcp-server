import TcaClient from "./client.js"

export async function getIssueList(state: string, severity: string, offset: number, limit: number) {
    const method = "GET"
    const pathParams = TcaClient.getcommonParams()
    const path = `/server/analysis/api/orgs/${pathParams.orgSid}/teams/${pathParams.teamName}/repos/${pathParams.repoId}/projects/${pathParams.projectId}/codelint/issues/`
    const clientInst = TcaClient.getInstance()
    const params = {
        state: state,
        severity: severity,
        offset: offset,
        limit: limit
    }
    return clientInst.request(path, method, null, params)
}
