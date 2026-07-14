---
name: tencent-pg-ops-troubleshooter
description: Troubleshoot operational issues for Tencent Cloud PostgreSQL by calling Tencent Cloud OpenAPI directly. This skill can use the full aligned PostgreSQL OpenAPI action set documented in @references/api_reference.md. Start with evidence collection first, and require explicit confirmation before any write, fee-impacting, or high-risk action.
description_zh: 运维排障
description_en: Ops troubleshooter
disable: false
agent_created: true
---

# tencent-pg-ops-troubleshooter

## When to use
- Need to troubleshoot an operational PostgreSQL issue rather than run a routine inspection.
- Need guided incident triage for connectivity, backup, account, network, SSL, readonly replica, parameter-related, or general instance problems.
- Need an evidence-first workflow that can continue into aligned remediation actions without depending on a deployed MCP server.

## Steps
1. Classify the incident first using the user's symptom:
   - connection or access failure
   - backup or restore concern
   - account or permission issue
   - network or SSL issue
   - readonly replica or replication issue
   - parameter or configuration concern
   - general instance abnormality
2. Confirm the scope: region, instance, and the narrowest affected surface such as database, account, replica, or backup task. If the region, instance ID, or both are missing from the user's request, immediately reply with a direct message (no interactive prompt or choice menu) that includes: (a) the console link https://console.cloud.tencent.com/postgres where the user can find their instance ID and region, and (b) a concrete example reply format like "ap-guangzhou postgres-abc12345". Do not proceed with any API call until the target scope is confirmed.
3. Normalize region input before any OpenAPI call by following `@references/common/region_normalization.md`. Accept standard region codes such as `ap-guangzhou` and common aliases such as `广州`, `上海`, `成都`, `北京`, then convert them to the canonical region code like `ap-guangzhou`, `ap-shanghai`, `ap-chengdu`, `ap-beijing`. If the input cannot be normalized safely, stop and use the `invalid-region` template in `@references/common/error_handling.md`, including the official region links.
4. Check runtime prerequisites using a foolproof order. Prefer `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_REGION`, and optional `TENCENTCLOUD_SESSION_TOKEN`; also accept compatible names such as `MCP_REQUEST_SECRET_ID`, `MCP_REQUEST_SECRET_KEY`, `MCP_REQUEST_SESSION_TOKEN`, `MCP_SECRET_ID`, and `MCP_SECRET_KEY`. If required values are missing, stop and use the `missing-credentials` template in `@references/common/error_handling.md`, including one copyable example plus the official API key and region links.
5. Call Tencent Cloud PostgreSQL OpenAPI directly. Any aligned Action listed in @references/api_reference.md may be used. Start with read-only evidence collection first. For audit, business, fee-impacting, or critical write actions, explain impact clearly and require explicit confirmation before execution.
6. Prefer the official Tencent Cloud SDK. If the SDK is unavailable, do not block on that alone; generate properly TC3-signed HTTPS requests locally instead. Read credentials only from environment variables or other secure runtime context; never hardcode or persist them in code, skill files, or query parameters. If the user still wants the SDK path, follow the `missing-sdk` template in `@references/common/error_handling.md`, offer the install command, and ask whether to execute it.
7. Use @references/api_reference.md to map the incident type to the OpenAPI modules that should be checked first, then expand only when the first-pass evidence is insufficient.
8. If the user asks for remediation or the incident cannot be resolved without action, choose the smallest aligned action that addresses the problem and explain prerequisites, impact, and confirmation requirement first.
9. Summarize in runbook format: incident summary, findings, blockers, immediate next steps, and actions that require explicit approval.

## Pitfalls
- Do not run every available aligned Action for every incident.
- Do not treat troubleshooting as approval to change configuration.
- Do not hide uncertainty; say when the evidence is incomplete.
- Do not perform write, fee-impacting, or high-risk actions without explicit confirmation.

## Verification
- Name the incident type and affected scope clearly.
- Show findings by module instead of a flat dump.
- State which OpenAPI actions were inspected, and list any action proposed or executed after confirmation.
- Separate safe next steps from actions that require confirmation or elevated access.
