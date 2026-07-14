# Inspection OpenAPI reference

## Scope

This checklist is for PostgreSQL inspection and aligned follow-up operations by calling Tencent Cloud PostgreSQL OpenAPI directly.

The allowed Action names in this file are the full set of 48 actions currently aligned with the PostgreSQL MCP implementation in `src/postgres/tools/openapi_alignment.go`.

## Direct OpenAPI baseline

- Endpoint: `https://postgres.tencentcloudapi.com`
- Version: `2017-03-12`
- Auth: `TC3-HMAC-SHA256`
- Credentials: prefer `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_REGION`, and optional `TENCENTCLOUD_SESSION_TOKEN`; also accept compatible names such as `MCP_REQUEST_SECRET_ID`, `MCP_REQUEST_SECRET_KEY`, `MCP_REQUEST_SESSION_TOKEN`, `MCP_SECRET_ID`, and `MCP_SECRET_KEY`
- Region input: accept canonical region codes such as `ap-guangzhou` and common aliases such as `广州`, `上海`, `成都`, `北京`, then normalize before the first OpenAPI call
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

## Inspection sequence

1. Normalize the region first by following `@references/common/region_normalization.md`, then confirm the target instance.
2. Prefer read-only evidence collection first: instance, backup, network, SSL, replica, slow-query, and error-log context.
3. Add broader module checks only when the first-pass evidence is insufficient.
4. If the user asks for remediation or the inspection uncovers a concrete issue, select the smallest aligned action that can address it.
5. Before any audit, business, fee-impacting, or critical action, explain impact, target scope, and confirmation requirement first.
6. Summarize the result before and after any confirmed action.

## Output schema

### 1. Overall status
- Region
- Instance
- Risk level: `low` / `medium` / `high`
- Short conclusion

### 2. Evidence snapshot
- Lifecycle / status
- Backup evidence
- Network / SSL evidence
- Replica evidence
- Optional anomaly evidence

### 3. Risk items
For each risk item include:
- symptom
- evidence
- impact
- confidence

### 4. Next actions
- `no action required`, or
- ordered low-risk follow-up checks and owner suggestions, or
- aligned remediation actions with confirmation requirement clearly stated

## Guardrails

- You may use any aligned action above, but only when it is relevant to the current inspection goal.
- Never call any Action outside the aligned list above.
- Start with evidence collection before write, fee-impacting, or high-risk actions.
- For audit, business, fee-impacting, or critical actions, explicit confirmation is required.
- Separate evidence from inference.
- If the aligned actions do not expose a requested metric or detail, say so explicitly.
