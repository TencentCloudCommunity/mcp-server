---
name: tencent-pg-slowquery-diagnosis
description: Diagnose slow SQL and performance anomalies for Tencent Cloud PostgreSQL by calling Tencent Cloud OpenAPI directly. This skill can use the full aligned PostgreSQL OpenAPI action set documented in @references/api_reference.md. Start with evidence collection first, and require explicit confirmation before any write, fee-impacting, or high-risk action.
description_zh: 慢 SQL 诊断
description_en: Slow query diagnosis
disable: false
agent_created: true
---

# tencent-pg-slowquery-diagnosis

## When to use
- Need to analyze slow SQL, latency spikes, or performance degradation.
- Need a first-pass diagnosis before proposing or executing tuning and change actions.
- Need a reusable workflow that calls Tencent Cloud OpenAPI directly instead of relying on a deployed MCP server.

## Steps
1. Confirm the scope first: region, target instance, optional database or SQL fingerprint, and preferred time window. If the region, instance ID, or both are missing from the user's request, immediately reply with a direct message (no interactive prompt or choice menu) that includes: (a) the console link https://console.cloud.tencent.com/postgres where the user can find their instance ID and region, and (b) a concrete example reply format like "ap-guangzhou postgres-abc12345". Do not proceed with any API call until the target scope is confirmed.
2. Normalize region input before any OpenAPI call by following `@references/common/region_normalization.md`. Accept standard region codes such as `ap-guangzhou` and common aliases such as `广州`, `上海`, `成都`, `北京`, then convert them to the canonical region code like `ap-guangzhou`, `ap-shanghai`, `ap-chengdu`, `ap-beijing`. If the input cannot be normalized safely, stop and use the `invalid-region` template in `@references/common/error_handling.md`, including the official region links.
3. Check runtime prerequisites using a foolproof order. Prefer `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_REGION`, and optional `TENCENTCLOUD_SESSION_TOKEN`; also accept compatible names such as `MCP_REQUEST_SECRET_ID`, `MCP_REQUEST_SECRET_KEY`, `MCP_REQUEST_SESSION_TOKEN`, `MCP_SECRET_ID`, and `MCP_SECRET_KEY`. If required values are missing, stop and use the `missing-credentials` template in `@references/common/error_handling.md`, including one copyable example plus the official API key and region links.
4. Call Tencent Cloud PostgreSQL OpenAPI directly. Any aligned Action listed in @references/api_reference.md may be used. Start with read-only evidence collection first. For audit, business, fee-impacting, or critical write actions, explain impact clearly and require explicit confirmation before execution.
5. Prefer the official Tencent Cloud SDK. If the SDK is unavailable, do not block on that alone; generate properly TC3-signed HTTPS requests locally instead. Read credentials only from environment variables or other secure runtime context; never hardcode or persist them in code, skill files, or query parameters. If the user still wants the SDK path, follow the `missing-sdk` template in `@references/common/error_handling.md`, offer the install command, and ask whether to execute it.
6. Gather evidence in this order:
   - current instance state, lifecycle, and ongoing task context
   - slow-query and query-analysis evidence
   - error-log evidence when errors may explain the slowdown
   - parameter, backup, readonly-topology, and other aligned context only when they help explain the symptom
7. Use @references/api_reference.md as the diagnosis frame and aligned action catalog.
8. Rank likely causes rather than forcing a single root cause. If the user asks for remediation and the evidence supports it, only then move to the smallest aligned action that addresses the issue, with explicit confirmation when required.
9. Keep optimization guidance conservative. Recommend or execute changes only when the evidence points in that direction and the action is within the aligned 48-action boundary.

## Pitfalls
- Do not invent SQL-level evidence when only instance-level or slow-query-summary evidence is available.
- Do not recommend or execute parameter, network, capacity, or topology changes without showing why.
- Do not confuse symptom correlation with confirmed causality.
- Do not perform write, fee-impacting, or high-risk actions without explicit confirmation.

## Verification
- Include the diagnosis window, target instance, and what evidence was inspected.
- Include a ranked list of likely causes or explicitly state that evidence is insufficient.
- State which OpenAPI actions were inspected, and list any action proposed or executed after confirmation.
- End with safe next actions, follow-up checks, or escalation advice.
