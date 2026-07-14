# Ops troubleshooter OpenAPI reference

## Scope

This checklist maps operational incident handling to the full aligned Tencent Cloud PostgreSQL OpenAPI action set for investigation and follow-up operations.

The allowed Action names in this file are the full set of 48 actions currently aligned with the PostgreSQL MCP implementation in `src/postgres/tools/openapi_alignment.go`.

## Direct OpenAPI baseline

- Endpoint: `https://postgres.tencentcloudapi.com`
- Version: `2017-03-12`
- Auth: `TC3-HMAC-SHA256`
- Credentials: `SecretId` / `SecretKey` and optional temporary `Token`; read them from environment variables or other secure runtime context only
- Preferred call path: official Tencent Cloud SDK; fallback: locally generated TC3-signed HTTPS request
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

## Investigation sequence

1. Normalize the region first by following `@references/common/region_normalization.md`, then classify the incident type and lock the affected scope: region, instance, and affected sub-surface.
2. Start with the smallest relevant read-only evidence set: instance, account, backup, network, SSL, replica, parameter, slow-query, or error-log checks.
3. Add broader context only when the first-pass evidence is insufficient.
4. If the user asks for remediation or the incident clearly requires action, choose the smallest aligned action that addresses the problem.
5. Before any audit, business, fee-impacting, or critical action, explain impact, target scope, and confirmation requirement first.
6. Produce a runbook-style summary rather than a raw evidence dump.

## Output schema

### 1. Incident summary
- incident type
- affected scope
- current severity

### 2. Findings by module
- instance state
- account / backup / network / SSL / replica / parameter findings
- blockers or missing evidence

### 3. Immediate next steps
- safe read-only follow-up checks
- aligned remediation actions with confirmation requirement clearly stated
- owner handoff suggestions

### 4. Actions requiring approval
- any audit, business, fee-impacting, or critical action

## Guardrails

- Troubleshooting does not imply authorization to change configuration without explicit confirmation.
- You may use any aligned action above, but only when it is relevant to the incident.
- Never call any Action outside the aligned list above.
- Start with evidence collection before write, fee-impacting, or high-risk actions.
- For audit, business, fee-impacting, or critical actions, explicit confirmation is required.
- Keep uncertain hypotheses clearly marked.
