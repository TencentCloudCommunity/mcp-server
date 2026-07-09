# PostgreSQL MCP Server Deployment Tutorial

> This tutorial is intended for first-time users of the Tencent Cloud PostgreSQL MCP Server. After you click **Deploy MCP Server**, this document is rendered in the popup. It is recommended to read it sequentially using the table of contents in the upper-left corner.

## 1. Preparation Before Deployment

Before choosing a deployment method, please confirm that the following information is ready.

### 1.1 Applicable Scenarios

- **Want to deploy on Tencent Cloud and expose it to your team over HTTPS**: Choose Method 1 (self-hosted Tencent Cloud SCF deployment).
- **Want full control over the network, host, and runtime environment**: Choose Method 2 (self-hosted `streamable-http` service).
- **Want to connect a local client directly to the source-based service**: Choose Method 3 (local `stdio`).

> If you need cloud hosting, please follow Method 1 to complete the SCF deployment under your own Tencent Cloud account, and connect your client using your own function URL.

### 1.2 Required Permissions

After deployment, calling MCP tools requires a Tencent Cloud account with `postgres:*` access permissions. It is recommended to create a dedicated CAM sub-account for the MCP client and grant the minimum required permissions based on your usage scenario.

Credential management page: <https://console.cloud.tencent.com/cam/capi>

Please make sure to save your `SecretId` and `SecretKey`, which will be used for subsequent OpenAPI calls.

In addition, the tools provided by this MCP server require a region code. Please confirm which region your target PostgreSQL instance is deployed in.

Reference document: <https://cloud.tencent.com/document/product/1596/77930>

### 1.3 Prerequisite Resources

| Resource | Required | Description |
|---|---|---|
| `SecretId` / `SecretKey` | Yes | Identity credentials used to call OpenAPI |
| Region code (for example, `ap-guangzhou`) | Yes | Used as the `region` parameter for all tools |
| Instance ID | No | Required only when operating on a specific instance |

## 2. MCP Capabilities

By default, **48 tools** are exposed, covering 9 major modules including instances, accounts, databases, parameters, backups, monitoring, networking, SSL, and read-only instances. They are listed below by module.

### 2.1 Instance Management

| Tool Name | Description |
|---|---|
| `DescribeDBInstances` | Query the instance list |
| `DescribeDBInstanceAttribute` | Query instance details |
| `CreateInstances` | Create an instance (cost confirmation required) |
| `ModifyDBInstanceName` | Modify the instance name |
| `ModifyDBInstanceSpec` | Change the instance specification (scale up/down, cost confirmation required) |
| `RestartDBInstance` | Restart the instance |
| `IsolateDBInstances` | Isolate an instance (business confirmation required) |
| `DisIsolateDBInstances` | Remove instance isolation |
| `UpgradeDBInstanceKernelVersion` | Upgrade the instance kernel version |
| `DescribeTasks` | Query async task status |
| `DescribeClasses` | Query available instance classes |
| `DescribeDBVersions` | Query available database versions |
| `DescribeRegions` | Query available sale regions |
| `DescribeZones` | Query available sale zones |
| `DescribeProductConfig` | Query sale configuration details |

### 2.2 Account Management

| Tool Name | Description |
|---|---|
| `DescribeAccounts` | Query the database account list of an instance |
| `DescribeAccountPrivileges` | Query privilege information of a database account |
| `CreateAccount` | Create a database account |
| `DeleteAccount` | Delete a database account |
| `ModifyAccountPrivileges` | Modify account privileges (grant / revoke / change account type) |
| `ResetAccountPassword` | Reset the account password |

### 2.3 Database Management

| Tool Name | Description |
|---|---|
| `DescribeDatabases` | Query the database list of an instance |
| `DescribeDatabaseObjects` | Query the database object list |
| `CreateDatabase` | Create a database |
| `ModifyDatabaseOwner` | Modify the database owner |

### 2.4 Parameter Management

| Tool Name | Description |
|---|---|
| `DescribeDBInstanceParameters` | Query instance parameters |
| `DescribeParameterTemplates` | Query parameter template list |
| `DescribeParameterTemplateAttributes` | Query parameter template details |
| `DescribeParamsEvent` | Query parameter modification events |
| `ModifyDBInstanceParameters` | Modify instance parameters |

### 2.5 Backup and Recovery

| Tool Name | Description |
|---|---|
| `DescribeBackupOverview` | Query backup overview |
| `DescribeBaseBackups` | Query base backup list |
| `DescribeLogBackups` | Query log backup list |
| `DescribeAvailableRecoveryTime` | Query available recovery time range |
| `DescribeCloneDBInstanceSpec` | Query purchasable specifications for cloned instances |
| `DescribeBackupDownloadURL` | Get the backup download URL |
| `CreateBaseBackup` | Create a base backup |
| `CloneDBInstance` | Clone an instance (cost confirmation required) |

### 2.6 Monitoring and Diagnosis

| Tool Name | Description |
|---|---|
| `DescribeSlowQueryList` | Query the slow query list |
| `DescribeSlowQueryAnalysis` | Analyze slow queries |
| `DescribeDBErrlogs` | Query error logs |

### 2.7 Network Management

| Tool Name | Description |
|---|---|
| `OpenDBExtranetAccess` | Enable public network access for an instance |
| `CloseDBExtranetAccess` | Disable public network access for an instance |
| `DescribeDBInstanceSecurityGroups` | Query instance security groups |
| `ModifyDBInstanceSecurityGroups` | Modify instance security groups |

### 2.8 SSL Configuration

| Tool Name | Description |
|---|---|
| `DescribeDBInstanceSSLConfig` | Query instance SSL configuration |

### 2.9 Read-Only Instances

| Tool Name | Description |
|---|---|
| `DescribeReadOnlyGroups` | Query read-only group list |
| `CreateReadOnlyDBInstance` | Create a read-only instance (cost confirmation required) |

> Write operations are still constrained by the permission scope, the `READ_ONLY` setting, and secondary confirmation mechanisms. It is recommended to start with read-only capabilities and then enable write operations only when necessary.

## 3. Choose a Deployment Method

Below are 3 methods listed in **descending order of recommendation**. Each method is explained in the order of **Prerequisites → Deployment Steps → Client Configuration**.

### 3.1 Method 1: Tencent Cloud SCF Self-Hosting (Recommended for Cloud Deployment)

> Suitable for scenarios where you want to run on Tencent Cloud and expose the service to your team through HTTPS / Function URL. The repository already provides SCF packaging scripts, startup scripts, and environment variable templates, but **you must complete the hosting and publishing process under your own Tencent Cloud account**.

#### Step 1: Pull the `src/postgres` directory as needed and build the SCF release package

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
./scripts/build_scf_zip.sh
```

By default, a zip package uploadable to SCF will be generated in the `dist/` directory.

#### Step 2: Create a Web Function in the SCF Console

Go to the [SCF Console](https://console.cloud.tencent.com/scf/list?rid=16&ns=default) and create a function with the following settings:

- Creation method: Select **Create from scratch**
- Function type: **Web Function**
- Runtime: **Go standard runtime (Go 1)**
- Code upload method: **Upload local zip**
- Environment variables: Refer to Step 4 below, or configure them after the function is created
- Other settings: Configure as needed. If public access is required, enable public access in the final **Function URL Configuration** section

#### Step 3: Startup Command

The zip package already includes `scf_bootstrap`, so in most cases you can use the built-in startup file directly.

If the console requires you to manually enter a startup command, use the same content as `deploy/scf/scf.console.startup.sh`:

```bash
#!/bin/bash
set -euo pipefail

export PG_MCP_RUNTIME="${PG_MCP_RUNTIME:-scf}"
export PORT="${PORT:-9000}"
export MCP_TRANSPORT="${MCP_TRANSPORT:-streamable-http}"
export MCP_SERVER_BIND_HOST="${MCP_SERVER_BIND_HOST:-0.0.0.0}"
export MCP_SERVER_PORT="${MCP_SERVER_PORT:-${PORT}}"
export MCP_SERVER_HTTP_ENDPOINT="${MCP_SERVER_HTTP_ENDPOINT:-/mcp}"
export MCP_STREAMABLE_HTTP_STATELESS="${MCP_STREAMABLE_HTTP_STATELESS:-true}"
export MCP_AUTH_MODE="${MCP_AUTH_MODE:-request-credential}"
export MCP_REQUEST_VALIDATE_IDENTITY="${MCP_REQUEST_VALIDATE_IDENTITY:-true}"
export MCP_REQUEST_CREDENTIAL_SCOPES="${MCP_REQUEST_CREDENTIAL_SCOPES:-pg.read}"
export MCP_REQUEST_ALLOWED_REGIONS="${MCP_REQUEST_ALLOWED_REGIONS:-}"
export READ_ONLY="${READ_ONLY:-true}"
export TOKEN_EXCHANGE_ENABLED="${TOKEN_EXCHANGE_ENABLED:-false}"

exec /var/user/postgres-server
```

#### Step 4: Configure Environment Variables

Minimum recommended configuration:

```env
MCP_TRANSPORT : streamable-http
MCP_AUTH_MODE : request-credential
MCP_REQUEST_VALIDATE_IDENTITY : true
MCP_REQUEST_CREDENTIAL_SCOPES : pg.read
READ_ONLY : true
MCP_SERVER_BIND_HOST : 0.0.0.0
MCP_SERVER_PORT : 9000
MCP_SERVER_HTTP_ENDPOINT : /mcp
MCP_STREAMABLE_HTTP_STATELESS : true
```

Recommended additional setting:

```env
MCP_SERVER_PUBLIC_URL=https://<your-function-url>/mcp
```

#### Step 5: Client Configuration

After deployment, connect your client to your own function URL:

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "https://<your-function-url>/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<your SecretId>",
        "X-TencentCloud-Secret-Key": "<your SecretKey>"
      }
    }
  }
}
```

> The client `url` must point to the **full MCP endpoint** (including the `/mcp` suffix). Do not use the root function URL.

#### Step 6: Quick Verification

```bash
curl -i https://<your-function-url>/healthz
```

After you receive a normal `200 OK` response, connect your client using `https://<your-function-url>/mcp`.

For detailed console instructions, see [`SCF_DEPLOY.md`](./SCF_DEPLOY.md).

---

### 3.2 Method 2: Self-Hosted `streamable-http` Service

Suitable for deployment on your own cloud hosts or internal servers and sharing with your team via a domain name.

> **Prerequisite**: Your local machine or cloud host has **Go 1.25+** installed and has outbound Internet access (required to access Tencent Cloud OpenAPI).

#### Step 1: Pull the `src/postgres` directory as needed

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> All subsequent commands must be executed in the `src/postgres/` directory.

#### Step 2: Prepare the Configuration

```bash
cp .env.example .env
```

Minimum recommended configuration:

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
READ_ONLY=true
```

#### Step 3: Start the Service

```bash
./scripts/run_server.sh
```

#### Step 4: Client Configuration

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "http://127.0.0.1:9000/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<your SecretId>",
        "X-TencentCloud-Secret-Key": "<your SecretKey>"
      }
    }
  }
}
```

> It is recommended to place the service behind HTTPS or a reverse proxy. Before exposing it to the public Internet, be sure to add IP allowlists or zero-trust access control.

---

### 3.3 Method 3: Local `stdio` (Recommended for Local Clients)

Suitable for command-based direct connections from local MCP clients such as Cursor, Claude Desktop, and WorkBuddy.

> **Prerequisite**: Your local machine has **Go 1.25+** installed and has outbound Internet access (required to access Tencent Cloud OpenAPI).

#### Step 1: Pull the `src/postgres` directory as needed

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

#### Step 2: Prepare the Configuration

```bash
cp .env.example .env
```

```env
MCP_TRANSPORT=stdio
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_REQUEST_SECRET_ID=<your SecretId>
MCP_REQUEST_SECRET_KEY=<your SecretKey>
READ_ONLY=true
```

Please note that Tencent Cloud credentials must be configured in the local `.env` file for the service to work correctly.

#### Step 3: Start

```bash
./scripts/run_stdio.sh
```

#### Step 4: Client Configuration

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "/absolute/path/to/mcp-server/src/postgres/scripts/run_stdio.sh",
      "env": {
        "MCP_REQUEST_SECRET_ID": "<your SecretId>",
        "MCP_REQUEST_SECRET_KEY": "<your SecretKey>"
      }
    }
  }
}
```

> It is recommended to use an **absolute path** for `command`. When many MCP clients launch a `stdio` process, their working directory is not the repository root; if you use `./scripts/run_stdio.sh`, it can easily fail with `spawn ./scripts/run_stdio.sh ENOENT`.
>
> If your client supports `cwd`, you can also explicitly set `cwd` to `src/postgres` and then use a relative path.
>
> `stdio` mode is intended only for trusted local environments and is not suitable for exposure as a remotely shared service.

## 4. Authentication and Security Configuration

Regardless of which deployment method you choose, calls use the **per-request credential mode**, meaning Tencent Cloud credentials are passed in each request through headers, and the server does not persist your secrets long-term.

### 4.1 Credential Delivery Method

- **HTTP / `streamable-http` / `sse` mode**: Pass credentials through the `X-TencentCloud-Secret-Id` and `X-TencentCloud-Secret-Key` headers.
- **`stdio` mode**: Inject credentials through environment variables such as `MCP_REQUEST_SECRET_ID` and `MCP_REQUEST_SECRET_KEY`.

### 4.2 Required Security Recommendations

- **Do not put `SecretId / SecretKey` in URLs or query parameters**. Pass them only through headers or environment variables.
- **Do not store secrets in SCF server-side environment variables**. Use the per-request credential mode instead.
- **Use a CAM sub-account with the minimum required permissions whenever possible**, and avoid long-term reuse of root account secrets.
- **Do not print plaintext credentials in logs, traces, or error responses**.
- **In production, always place the service behind HTTPS or a reverse proxy**. Do not expose plain HTTP directly to the public Internet.
- **Start with `READ_ONLY=true`**, and enable write operations only after the workflow is fully verified.

## 5. Related Documentation Links

If you need to consult more materials during deployment, you can open the following links:

- [PostgreSQL MCP Server Project Homepage](https://github.com/TencentCloudCommunity/mcp-server)
- [API Key (CAM) Management Console](https://console.cloud.tencent.com/cam/capi)
- [SCF Console](https://console.cloud.tencent.com/scf)
- [TencentDB for PostgreSQL Product Documentation](https://cloud.tencent.com/document/product/409)
- [TencentDB for PostgreSQL API Overview](https://cloud.tencent.com/document/product/409/16761)
- [Region and Availability Zone Mapping](https://cloud.tencent.com/document/product/1596/77930)
- [View Full Deployment Documentation](README.md)
