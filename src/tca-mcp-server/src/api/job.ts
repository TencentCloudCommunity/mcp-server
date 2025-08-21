import TcaClient from "./client.js"

export async function getJobList(offset: number, limit: number, state: string) {
    const method = "GET"
    const pathParams = TcaClient.getcommonParams()
    const path = `/server/main/api/orgs/${pathParams.orgSid}/teams/${pathParams.teamName}/repos/${pathParams.repoId}/projects/${pathParams.projectId}/jobs/`
    const clientInst = TcaClient.getInstance()
    const params = {
        "offset": offset,
        "limit": limit,
        "state": state,
    }
    return clientInst.request(path, method, null, params)
}

export async function getJobDetail(jobId: string) {
    const method = "GET"
    const pathParams = TcaClient.getcommonParams()
    const path = `/server/main/api/orgs/${pathParams.orgSid}/teams/${pathParams.teamName}/repos/${pathParams.repoId}/projects/${pathParams.projectId}/jobs/${jobId}/detail/`
    const clientInst = TcaClient.getInstance()
    return clientInst.request(path, method, null, null)
}
