# PostgreSQL Companion Skills

This directory contains companion skills for TencentDB for PostgreSQL, located under `src/postgres/skills/`.

These assets belong to the **workflow layer** rather than acting as standalone database connectors. They do not depend on a separately deployed PostgreSQL MCP Server. Instead, they call Tencent Cloud PostgreSQL OpenAPI directly, but only within the Action scope that has already been aligned, validated, and wrapped in the current `src/postgres` module.

**中文版本**: [`README.md`](./README.md)

## Overview

- **Same domain, different layers**: these skills sit next to the PostgreSQL MCP implementation, while remaining isolated from the Go runtime, SCF deployment scripts, and npm launcher.
- **Evidence first, broad coverage**: these packaged skills can use all **48 aligned PostgreSQL OpenAPI Actions** currently validated in this repository. They should collect evidence first, then move to change or remediation flows only when needed.
- **Direct OpenAPI access**: the skills call Tencent Cloud PostgreSQL OpenAPI directly and do not require users to deploy the MCP Server first.
- **Explicit confirmation for high-risk actions**: audit-related, business-impacting, fee-impacting, or critical write operations must explain their impact clearly and obtain explicit confirmation before execution.
- **Release-oriented distribution**: users should install packaged zip artifacts from GitHub Release or COS instead of using the source directory directly.

## Simplified User Configuration (Recommended)

For most end users, **three environment variables are enough** to start using these PostgreSQL skills:

```bash
export TENCENTCLOUD_SECRET_ID="your SecretId"
export TENCENTCLOUD_SECRET_KEY="your SecretKey"
export TENCENTCLOUD_REGION="ap-guangzhou" # Use the region where your target instance is located
# Optional: add TENCENTCLOUD_SESSION_TOKEN when using temporary credentials
```

Additional conventions:
- **Preferred variable names**: `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_REGION`
- **Compatible legacy variable names**: `MCP_REQUEST_SECRET_ID`, `MCP_REQUEST_SECRET_KEY`, `MCP_REQUEST_SESSION_TOKEN`, `MCP_SECRET_ID`, `MCP_SECRET_KEY`
- **Natural-language region input is allowed**: common Chinese region names such as `广州`, `上海`, `成都`, and `北京` should be normalized to standard region codes before execution, for example `广州 -> ap-guangzhou`, `上海 -> ap-shanghai`, `成都 -> ap-chengdu`. If you are not sure about the region code of your instance, see <https://cloud.tencent.com/document/product/1596/77930>.
- **SDK first**: skills should prefer the Tencent Cloud official SDK. If the runtime environment does not have the SDK installed, they may automatically fall back to locally generated TC3-signed HTTPS requests. If needed, the skill can also guide or trigger SDK installation according to the runtime design.
- **Secrets must stay in runtime context only**: never write `SecretId`, `SecretKey`, or session tokens into repository files, skill documents, URLs, or query parameters.

## Included Skills

### `tencent-pg-inspection`
- Daily PostgreSQL health inspection
- Typical requests: `PG inspection`, `health check`, `backup check`, `resource waterline check`
- Focus areas: instance health, backup status, networking / SSL / read-only context, parameter posture, risk summary, and aligned follow-up actions when needed

### `tencent-pg-slowquery-diagnosis`
- Slow SQL and performance diagnosis
- Typical requests: `slow SQL analysis`, `SQL performance diagnosis`, `why is the query slower`
- Focus areas: slow-query evidence, error logs, instance context, ranked possible causes, safe optimization suggestions, and aligned follow-up actions when needed

### `tencent-pg-ops-troubleshooter`
- Operations troubleshooting workflow
- Typical requests: `PG troubleshooting`, `instance issue investigation`, `SSL issue`, `backup failure`
- Focus areas: fault classification, module-based evidence collection, runbook-style next steps, and aligned follow-up actions when needed

## Alignment Boundary with PostgreSQL MCP

These skills do not treat the PostgreSQL MCP Server as a runtime prerequisite, but they must still remain aligned with the OpenAPI coverage already validated in `src/postgres`.

Current constraints:
- Allowed Actions must come from the currently aligned PostgreSQL OpenAPI set, with `tools/openapi_alignment.go` serving as the maintenance baseline.
- All three workflow skills **may use the same 48 aligned Actions**, but should start from the modules most relevant to the current task instead of calling everything indiscriminately.
- Even if Tencent Cloud officially provides more PostgreSQL OpenAPI Actions, **they must not be used by these skills unless this repository has aligned them first**.

## Currently Exposed 48 Aligned Actions

### 1. Instance Management (15)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeDBInstances` | Query the instance list | Directly available |
| `DescribeDBInstanceAttribute` | Query instance details | Directly available |
| `CreateInstances` | Create an instance | Charge confirmation |
| `ModifyDBInstanceName` | Modify the instance name | Explicit confirmation |
| `ModifyDBInstanceSpec` | Change instance specifications | Charge confirmation |
| `RestartDBInstance` | Restart an instance | Explicit confirmation |
| `IsolateDBInstances` | Isolate an instance | Explicit confirmation |
| `DisIsolateDBInstances` | Recover an isolated instance | Explicit confirmation |
| `UpgradeDBInstanceKernelVersion` | Upgrade the instance kernel version | Explicit confirmation |
| `DescribeTasks` | Query async task status | Directly available |
| `DescribeClasses` | Query available instance classes | Directly available |
| `DescribeDBVersions` | Query available database versions | Directly available |
| `DescribeRegions` | Query sale regions | Directly available |
| `DescribeZones` | Query sale availability zones | Directly available |
| `DescribeProductConfig` | Query product configuration | Directly available |

### 2. Account Management (6)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeAccounts` | Query database accounts on an instance | Directly available |
| `DescribeAccountPrivileges` | Query database account privileges | Directly available |
| `CreateAccount` | Create a database account | Explicit confirmation |
| `DeleteAccount` | Delete a database account | Explicit confirmation |
| `ModifyAccountPrivileges` | Modify account privileges | Explicit confirmation |
| `ResetAccountPassword` | Reset an account password | Explicit confirmation |

### 3. Database Management (4)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeDatabases` | Query databases on an instance | Directly available |
| `DescribeDatabaseObjects` | Query database objects | Directly available |
| `CreateDatabase` | Create a database | Explicit confirmation |
| `ModifyDatabaseOwner` | Modify database ownership | Explicit confirmation |

### 4. Parameter Management (5)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeDBInstanceParameters` | Query instance parameters | Directly available |
| `DescribeParameterTemplates` | Query parameter template list | Directly available |
| `DescribeParameterTemplateAttributes` | Query parameter template details | Directly available |
| `DescribeParamsEvent` | Query parameter change events | Directly available |
| `ModifyDBInstanceParameters` | Modify instance parameters | Critical, explicit confirmation |

### 5. Backup and Recovery (8)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeBackupOverview` | Query backup overview | Directly available |
| `DescribeBaseBackups` | Query base backup list | Directly available |
| `DescribeLogBackups` | Query log backup list | Directly available |
| `DescribeAvailableRecoveryTime` | Query available recovery time range | Directly available |
| `DescribeCloneDBInstanceSpec` | Query purchasable specs for cloned instances | Directly available |
| `DescribeBackupDownloadURL` | Get backup download URLs | Explicit confirmation |
| `CreateBaseBackup` | Create a base backup | Explicit confirmation |
| `CloneDBInstance` | Clone an instance | Charge confirmation |

### 6. Monitoring and Diagnostics (3)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeSlowQueryList` | Query the slow query list | Directly available |
| `DescribeSlowQueryAnalysis` | Analyze slow queries | Directly available |
| `DescribeDBErrlogs` | Query error logs | Directly available |

### 7. Network Management (4)

| Action | Description | Execution Requirement |
|---|---|---|
| `OpenDBExtranetAccess` | Enable public access for an instance | Explicit confirmation |
| `CloseDBExtranetAccess` | Disable public access for an instance | Explicit confirmation |
| `DescribeDBInstanceSecurityGroups` | Query instance security groups | Directly available |
| `ModifyDBInstanceSecurityGroups` | Modify instance security groups | Explicit confirmation |

### 8. SSL Configuration (1)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeDBInstanceSSLConfig` | Query instance SSL configuration | Directly available |

### 9. Read-Only Instances (2)

| Action | Description | Execution Requirement |
|---|---|---|
| `DescribeReadOnlyGroups` | Query read-only group list | Directly available |
| `CreateReadOnlyDBInstance` | Create a read-only instance | Charge confirmation |

> Notes:
> - **Directly available**: query-style actions that can be used directly for evidence collection and diagnosis.
> - **Explicit confirmation**: the skill must explain scope, target instance, and expected outcome before execution and obtain clear user confirmation.
> - **Charge confirmation**: in addition to explicit confirmation, the skill must warn about possible resource or billing impact in advance.

## Directory Structure

```text
skills/
├─ README.md
├─ README_EN.md
├─ package.json
├─ references/
│  └─ common/
│     ├─ region_normalization.md
│     └─ error_handling.md
├─ scripts/
│  ├─ package-all.mjs
│  ├─ package-skill.mjs
│  └─ verify-skill.mjs
├─ dist/
├─ tencent-pg-inspection/
│  ├─ SKILL.md
│  ├─ references/
│  └─ assets/
├─ tencent-pg-slowquery-diagnosis/
│  ├─ SKILL.md
│  ├─ references/
│  └─ assets/
└─ tencent-pg-ops-troubleshooter/
   ├─ SKILL.md
   ├─ references/
   └─ assets/
```

## Local Development

1. Update `SKILL.md` and `references/` inside the target skill directory.
2. Reconfirm that all referenced OpenAPI Actions still exist in `tools/openapi_alignment.go` and remain consistent with the capability list in the repository root `README.md`.
3. If any OpenAPI parameter mapping changes, run `src/postgres/scripts/run_openapi_param_check.sh` first for contract validation.
4. Run local verification before packaging.
5. Only produce release zip artifacts after verification succeeds.

Recommended commands:

```bash
cd src/postgres
./scripts/run_openapi_param_check.sh

cd skills
npm run verify
npm run release
```

## Packaging

The packaging workflow is intentionally isolated from `src/postgres/package.json`.

The local packaging entry is `src/postgres/skills/package.json`, and the related scripts live in `src/postgres/skills/scripts/`.

Expected artifacts:
- `tencent-pg-inspection-vX.Y.Z.zip`
- `tencent-pg-slowquery-diagnosis-vX.Y.Z.zip`
- `tencent-pg-ops-troubleshooter-vX.Y.Z.zip`
- `tencentdb-postgresql-skill-vX.Y.Z.zip`

By default, the version baseline follows the PostgreSQL MCP npm version declared in `src/postgres/package.json`, unless `SKILL_VERSION` is provided explicitly.

## Release Assets

Recommended release organization:
- Attach the three individual skill packages to the PostgreSQL MCP Release.
- Provide a bundle package for most users.
- Organize the bundle release in a root-level format containing `SKILL.md`, `_meta.json`, and `references/`.
- Expand the three PostgreSQL child skill directories under `references/`, so platforms can recognize the root entry first and then descend into sub-skills.
- Optionally attach a skill release note document describing compatibility and changes.

Recommended tag alignment:
- PostgreSQL MCP tag: `postgres-mcp-server-vX.Y.Z`
- Skill asset version: `vX.Y.Z`

## WorkBuddy / CodeBuddy Installation

Recommended user flow:
1. Download either an individual skill zip package from Release assets, or the full bundle zip for bulk browsing. The bundle filename should be `tencentdb-postgresql-skill-vX.Y.Z.zip`.
2. If you downloaded the bundle, extract it first. At the root level you should see `SKILL.md`, `_meta.json`, and `references/`. The root skill name should be `TencentDB PostgreSQL Skill`, and the slug in `_meta.json` should be `tencentdb-postgresql-skill`.
3. Read the root-level `SKILL.md` first, then enter the corresponding child skill directory under `references/` and read its `SKILL.md`.
4. If you want to install into WorkBuddy / CodeBuddy, use the individual skill zip package from the Release.
5. Open skill management.
6. Upload and enable the target skill zip package.
7. In the current session or execution environment, prepare these variables first: `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, and `TENCENTCLOUD_REGION`. Add `TENCENTCLOUD_SESSION_TOKEN` only for temporary credential scenarios. Legacy `MCP_REQUEST_*` or `MCP_*` variable names should also remain compatible if they already exist in the environment.
8. The region can be provided either as a standard value (for example `ap-guangzhou`) or a common Chinese region name (for example `广州`, `上海`, `成都`), which should then be normalized before execution.
9. Use natural-language requests such as `PG inspection`, `slow SQL analysis`, or `PG troubleshooting`.

> Other AI clients such as Cursor or Claude Code can adapt the same process according to their own skill-loading model.
