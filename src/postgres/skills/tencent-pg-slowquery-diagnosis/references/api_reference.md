# Slow query diagnosis OpenAPI reference

## Scope

This checklist is for slow-SQL diagnosis, performance investigation, and aligned follow-up operations by calling Tencent Cloud PostgreSQL OpenAPI directly.

The allowed Action names in this file are the full set of 48 actions currently aligned with the PostgreSQL MCP implementation in `src/postgres/tools/openapi_alignment.go`.

## Direct OpenAPI baseline

- Endpoint: `https://postgres.tencentcloudapi.com`
- Version: `2017-03-12`
- Auth: `TC3-HMAC-SHA256`
- Credentials: prefer `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_REGION`, and optional `TENCENTCLOUD_SESSION_TOKEN`; also accept compatible names such as `MCP_REQUEST_SECRET_ID`, `MCP_REQUEST_SECRET_KEY`, `MCP_REQUEST_SESSION_TOKEN`, `MCP_SECRET_ID`, and `MCP_SECRET_KEY`
- Shared region rule: follow `@references/common/region_normalization.md` before the first OpenAPI call; accept canonical region codes such as `ap-guangzhou` and common aliases such as `广州`, `上海`, `成都`, `北京`
- Shared error template: follow `@references/common/error_handling.md` whenever credentials are missing, region input is invalid, or the user asks to install the SDK; credential and region guidance must include the official Tencent Cloud links defined there
- Preferred call path: official Tencent Cloud SDK; fallback: locally generated TC3-signed HTTPS request, and SDK absence alone should not be treated as the first blocking error
- Never put secrets into URLs or query parameters, and never hardcode them in source code or skill files

## Allowed OpenAPI actions (all 48 aligned with current MCP server)

### 1. Instance management
- `DescribeDBInstances`
- `DescribeDBInstanceAttribute`
- `CreateInstances` (fee-impacting, explicit confirmation required)
- `ModifyDBInstanceName` (explicit confirmation required)
- `ModifyDBInstanceSpec` (fee-impacting, explicit confirmation required)
- `RestartDBInstance` (explicit confirmation required)
- `IsolateDBInstances` (explicit confirmation required)
- `DisIsolateDBInstances` (explicit confirmation required)
- `UpgradeDBInstanceKernelVersion` (explicit confirmation required)
- `DescribeTasks`
- `DescribeClasses`
- `DescribeDBVersions`
- `DescribeRegions`
- `DescribeZones`
- `DescribeProductConfig`

### 2. Account management
- `DescribeAccounts`
- `DescribeAccountPrivileges`
- `CreateAccount` (explicit confirmation required)
- `DeleteAccount` (explicit confirmation required)
- `ModifyAccountPrivileges` (explicit confirmation required)
- `ResetAccountPassword` (explicit confirmation required)

### 3. Database management
- `DescribeDatabases`
- `DescribeDatabaseObjects`
- `CreateDatabase` (explicit confirmation required)
- `ModifyDatabaseOwner` (explicit confirmation required)

### 4. Parameter management
- `DescribeDBInstanceParameters`
- `DescribeParameterTemplates`
- `DescribeParameterTemplateAttributes`
- `DescribeParamsEvent`
- `ModifyDBInstanceParameters` (critical change, explicit confirmation required)

### 5. Backup and recovery
- `DescribeBackupOverview`
- `DescribeBaseBackups`
- `DescribeLogBackups`
- `DescribeAvailableRecoveryTime`
- `DescribeCloneDBInstanceSpec`
- `DescribeBackupDownloadURL` (explicit confirmation required)
- `CreateBaseBackup` (explicit confirmation required)
- `CloneDBInstance` (fee-impacting, explicit confirmation required)

### 6. Monitoring and diagnostics
- `DescribeSlowQueryList`
- `DescribeSlowQueryAnalysis`
- `DescribeDBErrlogs`

### 7. Network management
- `OpenDBExtranetAccess` (explicit confirmation required)
- `CloseDBExtranetAccess` (explicit confirmation required)
- `DescribeDBInstanceSecurityGroups`
- `ModifyDBInstanceSecurityGroups` (explicit confirmation required)

### 8. SSL configuration
- `DescribeDBInstanceSSLConfig`

### 9. Read-only instances
- `DescribeReadOnlyGroups`
- `CreateReadOnlyDBInstance` (fee-impacting, explicit confirmation required)

## Diagnosis sequence

1. Normalize the region, then fix the time window and target instance.
2. Confirm current instance status and whether a task, failover, or change event is in progress.
3. Pull slow-query evidence first, then correlate with error evidence and obvious instance context.
4. Add parameter, backup, replica, network, or topology context only when it helps explain the observed symptom.
5. If the user asks for remediation and the evidence supports it, choose the smallest aligned action that addresses the issue.
6. Before any audit, business, fee-impacting, or critical action, explain impact, target scope, and confirmation requirement first.
7. Rank likely causes with confidence levels and distinguish evidence from hypotheses.

## Output schema

### 1. Diagnosis summary
- Region
- Instance
- Time window
- Short symptom summary

### 2. Evidence reviewed
- Slow-query evidence
- Instance status / task evidence
- Error evidence
- Supporting context

### 3. Ranked likely causes
For each likely cause include:
- cause statement
- supporting evidence
- missing evidence, if any
- confidence: `low` / `medium` / `high`

### 4. Safe next steps
- SQL review or indexing suggestion
- follow-up parameter inspection
- scale or replica review
- aligned remediation actions with confirmation requirement clearly stated
- explicit escalation if evidence is not enough

## Guardrails

- You may use any aligned action above, but only when it is relevant to the current diagnosis goal.
- Never call any Action outside the aligned list above.
- Start with evidence collection before write, fee-impacting, or high-risk actions.
- For audit, business, fee-impacting, or critical actions, explicit confirmation is required.
- Do not claim a confirmed root cause without matching evidence.
- If SQL-level or session-level detail is unavailable from the aligned actions, say the diagnosis is limited by the available OpenAPI evidence.
