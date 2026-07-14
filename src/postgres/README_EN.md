# TencentDB for PostgreSQL MCP Server

> The official MCP Server for TencentDB for PostgreSQL wraps cloud API capabilities such as instances, accounts, databases, parameters, backups, monitoring, networking, read-only instances, and SSL into unified MCP tools that can be used directly by Cursor, Claude Desktop, WorkBuddy, and any other MCP-compatible client.
>
> It supports `stdio`, `streamable-http`, and `sse` transports for local hosts, remote hosts, or Tencent Cloud SCF. With the **per-request credential mode**, the server does not need to store user `SecretId` / `SecretKey` values for a long time.

**Product Link**: [TencentDB for PostgreSQL](https://cloud.tencent.com/product/postgres)  
**中文版本**: [`README.md`](./README.md)

## Repository Contents

This repository currently contains both **`MCP Server`** and **`Skill`** deliverables for PostgreSQL. They target different usage patterns, but share the same domain boundaries, OpenAPI alignment baseline, and security constraints.

- **`MCP Server`**: Lives in the main `src/postgres` module and packages TencentDB for PostgreSQL OpenAPI capabilities as standard MCP tools, suitable for Cursor, Claude Desktop, WorkBuddy, and other MCP-compatible clients.
- **`Skill`**: Lives in [`skills/`](./skills/) and provides task-oriented skill packages plus supporting documents for inspection, slow SQL diagnosis, and operations troubleshooting workflows.
- **Relationship**: `MCP` focuses on the tool integration layer, while `Skill` focuses on the workflow layer. They are maintained side by side in the same repository rather than replacing each other.

If you want to connect a standard MCP client, continue with the MCP deployment and configuration guide below. If you want the PostgreSQL skill package overview, see [`skills/README_EN.md`](./skills/README_EN.md).

---

## 1. Before You Deploy

Before choosing a deployment method, make sure the following information is ready.

### 1.1 Typical Scenarios

- **You want to deploy on Tencent Cloud and expose it to your team over HTTPS**: choose Method 1 (Tencent Cloud SCF self-hosting)
- **You want full control over the network, host, and runtime environment**: choose Method 2 (self-hosted `streamable-http` service)
- **You want to connect a local client directly to the source-based service**: choose Method 3 (local `stdio`)

> If you need cloud hosting, complete the SCF deployment in your own Tencent Cloud account and connect clients with your own Function URL.

### 1.2 Required Permissions

When calling MCP tools after deployment, you need a Tencent Cloud account with PostgreSQL API access permissions. We recommend creating a dedicated CAM sub-account for MCP clients and granting only the minimum required permissions.

- **Credential page**: <https://console.cloud.tencent.com/cam/capi>
- Keep your `SecretId` and `SecretKey` properly stored for subsequent OpenAPI calls
- The tools provided by this MCP Server require region codes, so confirm the target instance region in advance
- Region reference: <https://cloud.tencent.com/document/product/1596/77930>

### 1.3 Required Resources

| Resource | Required | Description |
|---|---|---|
| `SecretId` / `SecretKey` | Yes | Identity credentials for calling OpenAPI |
| Region code (for example `ap-guangzhou`) | Yes | Used as the `region` parameter for all tools |
| Instance ID | No | Only needed for instance-specific operations |

---

## 2. MCP Capabilities

This server exposes **48 tools** across 9 major groups: instances, accounts, databases, parameters, backups, monitoring, networking, SSL, and read-only instances.

### 2.1 Instance Management (15)

| Tool Name | Description |
|---|---|
| `DescribeDBInstances` | Query the instance list |
| `DescribeDBInstanceAttribute` | Query instance details |
| `CreateInstances` | Create an instance (charge confirmation required) |
| `ModifyDBInstanceName` | Modify the instance name |
| `ModifyDBInstanceSpec` | Change instance specifications (scale up/down, charge confirmation required) |
| `RestartDBInstance` | Restart an instance |
| `IsolateDBInstances` | Isolate an instance (business confirmation required) |
| `DisIsolateDBInstances` | Recover an isolated instance |
| `UpgradeDBInstanceKernelVersion` | Upgrade the instance kernel version |
| `DescribeTasks` | Query async task status |
| `DescribeClasses` | Query available instance classes |
| `DescribeDBVersions` | Query available database versions |
| `DescribeRegions` | Query sale regions |
| `DescribeZones` | Query availability zones |
| `DescribeProductConfig` | Query product configuration |

### 2.2 Account Management (6)

| Tool Name | Description |
|---|---|
| `DescribeAccounts` | Query database accounts for an instance |
| `DescribeAccountPrivileges` | Query account privileges |
| `CreateAccount` | Create a database account |
| `DeleteAccount` | Delete a database account |
| `ModifyAccountPrivileges` | Modify account privileges (grant / revoke / account type) |
| `ResetAccountPassword` | Reset an account password |

### 2.3 Database Management (4)

| Tool Name | Description |
|---|---|
| `DescribeDatabases` | Query databases in an instance |
| `DescribeDatabaseObjects` | Query database objects |
| `CreateDatabase` | Create a database |
| `ModifyDatabaseOwner` | Modify the database owner |

### 2.4 Parameter Management (5)

| Tool Name | Description |
|---|---|
| `DescribeDBInstanceParameters` | Query instance parameters |
| `DescribeParameterTemplates` | Query parameter template list |
| `DescribeParameterTemplateAttributes` | Query parameter template details |
| `DescribeParamsEvent` | Query parameter change events |
| `ModifyDBInstanceParameters` | Modify instance parameters |

### 2.5 Backup and Recovery (8)

| Tool Name | Description |
|---|---|
| `DescribeBackupOverview` | Query backup overview |
| `DescribeBaseBackups` | Query base backup list |
| `DescribeLogBackups` | Query log backup list |
| `DescribeAvailableRecoveryTime` | Query available recovery time range |
| `DescribeCloneDBInstanceSpec` | Query available specs for clone instances |
| `DescribeBackupDownloadURL` | Get backup download URLs |
| `CreateBaseBackup` | Create a base backup |
| `CloneDBInstance` | Clone an instance (charge confirmation required) |

### 2.6 Monitoring and Diagnostics (3)

| Tool Name | Description |
|---|---|
| `DescribeSlowQueryList` | Query slow query list |
| `DescribeSlowQueryAnalysis` | Analyze slow queries |
| `DescribeDBErrlogs` | Query error logs |

### 2.7 Network Management (4)

| Tool Name | Description |
|---|---|
| `OpenDBExtranetAccess` | Enable public access |
| `CloseDBExtranetAccess` | Disable public access |
| `DescribeDBInstanceSecurityGroups` | Query security groups |
| `ModifyDBInstanceSecurityGroups` | Modify security groups |

### 2.8 SSL Configuration (1)

| Tool Name | Description |
|---|---|
| `DescribeDBInstanceSSLConfig` | Query SSL configuration |

### 2.9 Read-Only Instances (2)

| Tool Name | Description |
|---|---|
| `DescribeReadOnlyGroups` | Query read-only groups |
| `CreateReadOnlyDBInstance` | Create a read-only instance (charge confirmation required) |

> Write operations are still constrained by permission scope, `READ_ONLY` configuration, and confirmation mechanisms. Start with read-only access first, then enable write capabilities only when needed.

---

## 3. Choose a Deployment Method

Below are 3 deployment methods listed in **descending order of recommendation**. Each method is organized as prerequisites → deployment steps → client configuration.

### 3.1 Method 1: Tencent Cloud SCF Self-Hosting (Recommended for Cloud Deployment)

> Suitable when you want to run the service on Tencent Cloud and expose it to your team through HTTPS or a Function URL. This repository already provides SCF packaging scripts, startup scripts, and environment templates, but **you still need to finish deployment in your own Tencent Cloud account**.

#### Step 1: Fetch `src/postgres` and build the SCF package

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
./scripts/build_scf_zip.sh
```

A zip package ready for SCF upload will be generated in the `dist/` directory.

#### Step 2: Create a Web Function in the SCF console

Open the [SCF Console](https://console.cloud.tencent.com/scf/list?rid=16&ns=default) and create the function with the following settings:

- Creation mode: choose **Create from scratch**
- Function type: **Web Function**
- Runtime: **Go standard runtime (Go 1)**
- Code upload method: **Upload local zip**
- Environment variables: follow Step 4, or configure them after function creation
- Other options: enable public Function URL if public access is needed

#### Step 3: Startup command

The zip package already includes `scf_bootstrap`, so in most cases you can use the built-in startup file directly.

If the console asks for a manual startup command, use the same content as `deploy/scf/scf.console.startup.sh`:

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

#### Step 4: Environment variables

Minimum recommended configuration:

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
READ_ONLY=true
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
```

Recommended addition:

```env
MCP_SERVER_PUBLIC_URL=https://your-function-url/mcp
```

#### Step 5: Client configuration

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "https://your-function-url/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<Your SecretId>",
        "X-TencentCloud-Secret-Key": "<Your SecretKey>"
      }
    }
  }
}
```

> The client `url` must point to the **full MCP endpoint** (including the `/mcp` suffix), not the function root URL.

#### Step 6: Quick verification

```bash
curl -i https://your-function-url/healthz
```

After you get `200 OK`, connect the client to `https://your-function-url/mcp`.

See [`SCF_DEPLOY.md`](./SCF_DEPLOY.md) for the full SCF deployment guide.

---

### 3.2 Method 2: Self-Hosted `streamable-http` Service

Suitable for deployment on your own cloud hosts or intranet servers and sharing through a domain name.

> **Prerequisites**: **Go 1.25+** is installed on the local machine or host, and internet access is available to call Tencent Cloud OpenAPI.

#### Step 1: Fetch `src/postgres`

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> All commands below should be executed in the `src/postgres/` directory.

#### Step 2: Prepare configuration

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

#### Step 3: Start the service

```bash
./scripts/run_server.sh
```

#### Step 4: Client configuration

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "type": "streamable-http",
      "url": "http://127.0.0.1:9000/mcp",
      "headers": {
        "X-TencentCloud-Secret-Id": "<Your SecretId>",
        "X-TencentCloud-Secret-Key": "<Your SecretKey>"
      }
    }
  }
}
```

> Put it behind HTTPS or a reverse proxy when possible, and add IP allowlists or zero-trust controls before public exposure.

---

### 3.3 Method 3: Local `stdio` (Recommended for Local Clients)

Suitable for local MCP clients such as Cursor, Claude Desktop, and WorkBuddy.

> **Prerequisites**: **Go 1.25+** is installed locally and the machine has internet access for Tencent Cloud OpenAPI.

#### Step 1: Fetch `src/postgres`

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

#### Step 2: Prepare configuration

```bash
cp .env.example .env
```

```env
MCP_TRANSPORT=stdio
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
MCP_REQUEST_SECRET_ID=YourSecretId
MCP_REQUEST_SECRET_KEY=YourSecretKey
READ_ONLY=true
```

Make sure Tencent Cloud credentials are configured in the local `.env` file before use.

#### Step 3: Start

```bash
./scripts/run_stdio.sh
```

#### Step 4: Client configuration

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "/absolute/path/to/mcp-server/src/postgres/scripts/run_stdio.sh",
      "env": {
        "MCP_REQUEST_SECRET_ID": "<Your SecretId>",
        "MCP_REQUEST_SECRET_KEY": "<Your SecretKey>"
      }
    }
  }
}
```

> Prefer an **absolute path** for `command`. Many MCP clients do not start the `stdio` process from the repository root, so using `./scripts/run_stdio.sh` may lead to `spawn ./scripts/run_stdio.sh ENOENT`.
>
> If your client supports `cwd`, you can also set `cwd` to `src/postgres` and keep a relative path.
>
> `stdio` mode is only suitable for local trusted environments and should not be exposed as a shared remote service.

---

## 4. Authentication and Security Configuration

No matter which deployment method you choose, requests use the **per-request credential mode**. That means Tencent Cloud credentials are carried with each request, instead of being stored persistently on the server.

### 4.1 How Credentials Are Passed

- **HTTP / `streamable-http` / `sse` mode**: pass credentials through `X-TencentCloud-Secret-Id` and `X-TencentCloud-Secret-Key` headers
- **`stdio` mode**: inject credentials through environment variables such as `MCP_REQUEST_SECRET_ID` / `MCP_REQUEST_SECRET_KEY`

### 4.2 Required Security Recommendations

- **Do not put `SecretId / SecretKey` in URLs or query parameters**; use headers or environment variables only
- **Do not store keys in SCF server-side environment variables**; use per-request credentials instead
- **Use a least-privilege CAM sub-account** whenever possible
- **Never print plaintext credentials in logs, traces, or error messages**
- **Production deployments must stay behind HTTPS or a reverse proxy**; do not expose raw HTTP publicly
- **Start with `READ_ONLY=true`**, and enable write operations only after the flow is confirmed

---



## 6. Related Documents

If you need more information during deployment, see the following links:

- [PostgreSQL MCP Server project homepage](https://github.com/TencentCloudCommunity/mcp-server)
- [API Key (CAM) console](https://console.cloud.tencent.com/cam/capi)
- [SCF Console](https://console.cloud.tencent.com/scf)
- [TencentDB for PostgreSQL product documentation](https://cloud.tencent.com/document/product/409)
- [TencentDB for PostgreSQL API overview](https://cloud.tencent.com/document/product/409/16761)
- [Region and availability zone mapping](https://cloud.tencent.com/document/product/1596/77930)
- [Full SCF deployment guide](./SCF_DEPLOY.md)
- [Chinese README](./README.md)

---

## 7. License

This project is open-sourced under the **Apache-2.0** license.
