# TencentDB for PostgreSQL MCP Server

> This MCP Server wraps the cloud APIs into unified MCP tools, providing: **48 tools** covering 9 major modules including instances, accounts, databases, parameters, backups, monitoring, network, read-only instances, and SSL configuration; **multiple deployment forms**: `npx` one-click local launch, local `stdio` command direct connect, self-hosted `streamable-http` service, self-hosted `sse` compatibility mode, and **Tencent Cloud SCF Web Function** deployment; **per-request credential mode**: clients pass their own Tencent Cloud `SecretId` / `SecretKey` via headers on each request, and the server does not store user keys long-term.

**Product Link**: [TencentDB for PostgreSQL](https://cloud.tencent.com/product/postgres)

---

## 1. Tool List

Below is the list of all **48** tools grouped by module. All tools accept `region` (region) as the first required parameter, and internally call the corresponding Tencent Cloud API according to `region`.

### 1. Instances (15)

| No. | Tool Name | Description |
|---|---|---|
| 1 | `DescribeDBInstances` | Query instance list |
| 2 | `DescribeDBInstanceAttribute` | Query instance details |
| 3 | `CreateInstances` | Create instance (charge confirmation) |
| 4 | `ModifyDBInstanceName` | Modify instance name |
| 5 | `ModifyDBInstanceSpec` | Change instance specification (scale up/down, charge confirmation) |
| 6 | `RestartDBInstance` | Restart instance |
| 7 | `IsolateDBInstances` | Isolate instance (business confirmation) |
| 8 | `DisIsolateDBInstances` | Unisolate instance |
| 9 | `UpgradeDBInstanceKernelVersion` | Upgrade instance kernel version |
| 10 | `DescribeTasks` | Query async task status |
| 11 | `DescribeClasses` | Query available specifications |
| 12 | `DescribeDBVersions` | Query available database versions |
| 13 | `DescribeRegions` | Query available sale regions |
| 14 | `DescribeZones` | Query available sale zones |
| 15 | `DescribeProductConfig` | Query sale specification configurations |

### 2. Accounts (6)

| No. | Tool Name | Description |
|---|---|---|
| 16 | `DescribeAccounts` | Query database account list of an instance |
| 17 | `DescribeAccountPrivileges` | Query database account privileges |
| 18 | `CreateAccount` | Create database account |
| 19 | `DeleteAccount` | Delete database account |
| 20 | `ModifyAccountPrivileges` | Modify account privileges (grant/revoke/change account type) |
| 21 | `ResetAccountPassword` | Reset account password |

### 3. Databases (4)

| No. | Tool Name | Description |
|---|---|---|
| 22 | `DescribeDatabases` | Query database list of an instance |
| 23 | `DescribeDatabaseObjects` | Query database object list |
| 24 | `CreateDatabase` | Create database |
| 25 | `ModifyDatabaseOwner` | Modify database owner |

### 4. Parameters (5)

| No. | Tool Name | Description |
|---|---|---|
| 26 | `DescribeDBInstanceParameters` | Query instance parameters |
| 27 | `DescribeParameterTemplates` | Query parameter template list |
| 28 | `DescribeParameterTemplateAttributes` | Query parameter template details |
| 29 | `DescribeParamsEvent` | Query parameter modification events |
| 30 | `ModifyDBInstanceParameters` | Modify instance parameters |

### 5. Backups (8)

| No. | Tool Name | Description |
|---|---|---|
| 31 | `DescribeBackupOverview` | Query backup overview |
| 32 | `DescribeBaseBackups` | Query base backup list |
| 33 | `DescribeLogBackups` | Query log backup list |
| 34 | `DescribeAvailableRecoveryTime` | Query available recovery time range |
| 35 | `DescribeCloneDBInstanceSpec` | Query purchasable specifications for cloning instances |
| 36 | `DescribeBackupDownloadURL` | Get backup download URL |
| 37 | `CreateBaseBackup` | Create base backup |
| 38 | `CloneDBInstance` | Clone instance (charge confirmation) |

### 6. Monitoring (3)

| No. | Tool Name | Description |
|---|---|---|
| 39 | `DescribeSlowQueryList` | Query slow query list |
| 40 | `DescribeSlowQueryAnalysis` | Slow query analysis |
| 41 | `DescribeDBErrlogs` | Query error logs |

### 7. Network (4)

| No. | Tool Name | Description |
|---|---|---|
| 42 | `OpenDBExtranetAccess` | Enable instance public network access |
| 43 | `CloseDBExtranetAccess` | Disable instance public network access |
| 44 | `DescribeDBInstanceSecurityGroups` | Query instance security groups |
| 45 | `ModifyDBInstanceSecurityGroups` | Modify instance security groups |

### 8. SSL (1)

| No. | Tool Name | Description |
|---|---|---|
| 46 | `DescribeDBInstanceSSLConfig` | Query instance SSL configuration |

### 9. Read-only Instances (2)

| No. | Tool Name | Description |
|---|---|---|
| 47 | `DescribeReadOnlyGroups` | Query read-only group list |
| 48 | `CreateReadOnlyDBInstance` | Create read-only instance (charge confirmation) |

> **Tiered Protection**: Write tools are classified by risk as `LevelFee` (charge confirmation), `LevelBusiness` (business confirmation), `LevelAudit` (audit level), `LevelCritical` (high-risk, requires explicit double confirmation `confirm=true`).

---

## 2. Quick Start (Deployment Methods)

Below are 4 deployment methods listed in **descending order of recommendation**. Please choose as needed and review the prerequisites for each method:

- **Method 1: Tencent Cloud SCF self-hosting (recommended for cloud deployment)** —— Suitable when you want to run the service on Tencent Cloud and expose it via HTTPS; you need to create the function, upload the zip package, and configure environment variables in your own SCF account.
- **Method 2: Self-hosted streamable-http service** —— Need to fetch only the `src/postgres` directory (`git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git` → `cd mcp-server` → `git sparse-checkout set src/postgres` → `cd src/postgres`), and prepare **Go 1.25+** for local compilation and runtime.
- **Method 3: Local stdio** —— Need to fetch only the `src/postgres` directory (same as above) and prepare **Go 1.25+**, suitable for local Cursor / Claude Desktop / WorkBuddy clients.
- **Method 4: npx one-click launch** —— Only requires **Node.js 18+** installed locally, no need to clone the repository, just one command line.

### Method 1: Tencent Cloud SCF self-hosting (recommended for cloud deployment)

> Suitable when you want to run the MCP server on Tencent Cloud and share it with your team through HTTPS / Function URL. The repository provides SCF packaging scripts, startup scripts, and environment templates, and you can complete the deployment in your own Tencent Cloud account.

#### 1.1 Fetch only the `src/postgres` directory and build the SCF package

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
./scripts/build_scf_zip.sh
```

A zip package ready for SCF upload will be generated in the `dist/` directory.

#### 1.2 Create a Web Function in the SCF console

Open [Tencent Cloud SCF Console](https://console.cloud.tencent.com/scf) and create the function with:

- Function type: **Web Function**
- Runtime: **Go standard runtime**
- Code upload method: **Upload local zip**

After uploading the zip, enable the public Function URL.

#### 1.3 Start command

The zip package already includes `scf_bootstrap`, so in most cases the built-in startup file can be used directly.

If the console requires a manual start command, use the same content as `deploy/scf/scf.console.startup.sh`:

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

#### 1.4 Environment variables

```env
MCP_TRANSPORT=streamable-http
MCP_AUTH_MODE=request-credential
MCP_REQUEST_VALIDATE_IDENTITY=true
MCP_REQUEST_CREDENTIAL_SCOPES=pg.read
TOKEN_EXCHANGE_ENABLED=false
READ_ONLY=true
MCP_SERVER_BIND_HOST=0.0.0.0
MCP_SERVER_PORT=9000
MCP_SERVER_HTTP_ENDPOINT=/mcp
MCP_STREAMABLE_HTTP_STATELESS=true
```

Recommended additions:

```env
MCP_SERVER_PUBLIC_URL=https://your-function-url/mcp
```

#### 1.5 Client configuration example

The service uses the **per-request credential mode**, requiring your Tencent Cloud credentials to be attached in request headers (do not put keys in URL or query parameters):

- `X-TencentCloud-Secret-Id`
- `X-TencentCloud-Secret-Key`

> See [§3. Credential Acquisition Tutorial](#3-credential-acquisition-tutorial) for how to obtain credentials.

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

> ⚠️ The client `url` must point to the **complete MCP endpoint** (including `/mcp` suffix), not the function root URL.

#### 1.6 Quick verification

```bash
curl -i https://your-function-url/healthz
```

A normal `200 OK` response means your self-hosted SCF deployment is online and ready for MCP clients.


---

### Method 2: Self-hosted streamable-http Service

Suitable for deployment to your own cloud servers, intranet servers, shared with the team via domain.

> **Prerequisites**: **Go 1.25+** is installed on your local machine or cloud server, with internet access (to call Tencent Cloud OpenAPI).

#### 2.1 Fetch only the `src/postgres` directory

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> All subsequent commands need to be executed in the `src/postgres/` directory.

#### 2.2 Prepare Configuration

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

#### 2.3 Start Service

```bash
./scripts/run_server.sh
```

#### 2.4 MCP Client Configuration

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

> Recommended to put it behind HTTPS / reverse proxy; add IP whitelist / zero-trust access control before exposing to public network.

---

### Method 3: Local stdio (Recommended for Local Clients)

Suitable for direct command connection mode of local MCP clients like Cursor, Claude Desktop, WorkBuddy.

> **Prerequisites**: **Go 1.25+** is installed on your local machine, with internet access (to call Tencent Cloud OpenAPI).

#### 3.1 Fetch only the `src/postgres` directory

```bash
git clone --depth=1 --filter=blob:none --sparse https://github.com/TencentCloudCommunity/mcp-server.git
cd mcp-server
git sparse-checkout set src/postgres
cd src/postgres
```

> All subsequent commands need to be executed in the `src/postgres/` directory.

#### 3.2 Prepare Configuration

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

#### 3.3 Start

```bash
./scripts/run_stdio.sh
```

#### 3.4 MCP Client Configuration

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

> Prefer using an **absolute path** for `command`. Many MCP clients do not spawn `stdio` processes with the repository root as their working directory; if you keep `./scripts/run_stdio.sh`, you may hit `spawn ./scripts/run_stdio.sh ENOENT`.
>
> If your client supports `cwd`, you can also set `cwd` to `src/postgres` and then keep the relative path.
>
> `stdio` mode is only suitable for local trusted environments, not suitable for exposing as a remote shared service.

---

### Method 4: npx One-Click Launch (Simplest Local Experience)

> **Prerequisites**: **Node.js 18+** is installed on your local machine (including `npx`). No need to clone the Go repository, the npm package will automatically download pre-compiled binaries from GitHub Release according to the platform.
>
> Check if it is installed:
>
> ```bash
> node -v   # Expected v18.x or higher
> ```

#### 4.1 Command Line Direct Startup

```bash
npx -y postgres-mcp-server@latest
```

#### 4.2 MCP Client Configuration

```json
{
  "mcpServers": {
    "mcp-server-postgres": {
      "command": "npx",
      "args": ["-y", "postgres-mcp-server@latest"],
      "env": {
        "TRANSPORT": "stdio",
        "TENCENTCLOUD_SECRET_ID": "<Your SecretId>",
        "TENCENTCLOUD_SECRET_KEY": "<Your SecretKey>"
      }
    }
  }
}
```

> The `npx` launcher will download pre-compiled Go binaries from GitHub Release according to the platform, network connection is required for first run.

---

## 3. Credential Acquisition Tutorial

### Get SecretId / SecretKey

1. Log in to the [Tencent Cloud Console](https://console.cloud.tencent.com/);
2. Go to the **Access Management → API Key Management** page;
3. In the **API Key** tab, create or view existing **SecretId** and **SecretKey**.

> **Generation URL**: <https://console.cloud.tencent.com/cam/capi>

⚠️ **Security Recommendations**

- **Do not put `SecretId / SecretKey` in URL or query parameters**, only pass via Header or environment variables;
- **Do not put keys into SCF server-side environment variables**, use per-request credential mode;
- **Use a least-privilege CAM sub-account**, and avoid reusing root account keys for long periods;
- **Do not output credential plaintext in logs, traces, or error returns**.

### Recommended Permission Scope

Create a separate CAM sub-account or role for the MCP client, granting minimum permissions based on usage scenarios:

| Usage Scenario | Recommended Policy |
|---|---|
| Read-only inspection | `QcloudPGReadOnlyAccess` |
| Daily O&M (including write) | Custom policy, limited to specific `Action` list |
| Temporary troubleshooting | Least-privilege CAM sub-account, revoke access after use |

---

## 4. Region Mapping

When calling this MCP tool, the `region` parameter must be the **region code** (e.g., `ap-guangzhou`), not the Chinese region name.

Example: Guangzhou region

| Chinese Name | Region Code |
|---|---|
| Guangzhou | `ap-guangzhou` |
| Shanghai | `ap-shanghai` |
| Beijing | `ap-beijing` |
| Nanjing | `ap-nanjing` |
| Shenzhen | `ap-shenzhen` |
| Chengdu | `ap-chengdu` |
| Hong Kong | `ap-hongkong` |
| Singapore | `ap-singapore` |

> See [Region and Availability Zone Mapping Documentation](https://cloud.tencent.com/document/product/1596/77930) for the complete region list. If you don't know which region your instance is in, you can use the read-only tool `DescribeDBInstances` without passing `region`, or first call `DescribeRegions` to query.


## 5. API Usage Reference



 TencentDB for PostgreSQL API Overview  <https://cloud.tencent.com/document/product/409/16761> 



## 6. Security Recommendations

- **Production environment must be behind HTTPS / reverse proxy**, do not expose HTTP to public network;
- **Do not put `SecretId / SecretKey` in URL or query parameters**, all via Header or env;
- **Do not output credential plaintext in logs, traces, or error returns**;
- **Use a least-privilege CAM sub-account**, and avoid reusing root account keys for long periods;
- **Keep `READ_ONLY=true` as the starting point**, open write operations as needed after confirming the process;
- **Add IP whitelist, VPN, or zero-trust access control** before exposing to public;
- **Stateless environments like SCF / API Gateway only recommend `streamable-http`** + `MCP_STREAMABLE_HTTP_STATELESS=true`.



## 7. License

This project is open-sourced under the **Apache-2.0** license.

