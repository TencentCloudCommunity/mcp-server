---
name: tencent-pg-inspection
description: Run routine health inspection for Tencent Cloud PostgreSQL by calling Tencent Cloud OpenAPI directly. This skill can use the full aligned PostgreSQL OpenAPI action set documented in @references/api_reference.md. Start with evidence collection first, and require explicit confirmation before any write, fee-impacting, or high-risk action.
description_zh: PG 巡检
description_en: PG inspection
disable: false
agent_created: true
---

# tencent-pg-inspection

## When to use
- Need a routine PostgreSQL health check or inspection report.
- Need to review instance lifecycle, backup posture, network / SSL exposure, readonly topology, parameter posture, or recent anomalies.
- Need a reusable workflow that starts with evidence collection and can continue into aligned follow-up actions without depending on a deployed MCP server.

## Steps
1. Confirm the target scope first: region, instance ID or instance name, and optional time window. If the region, instance ID, or both are missing from the user's request, immediately reply with a direct message (no interactive prompt or choice menu) that includes: (a) the console link https://console.cloud.tencent.com/postgres where the user can find their instance ID and region, and (b) a concrete example reply format like "ap-guangzhou postgres-abc12345". Do not proceed with any API call until the target scope is confirmed.
2. Normalize region input before any OpenAPI call by following `@references/common/region_normalization.md`. Accept standard region codes such as `ap-guangzhou` and common aliases such as `广州`, `上海`, `成都`, `北京`, then convert them to the canonical region code like `ap-guangzhou`, `ap-shanghai`, `ap-chengdu`, `ap-beijing`. If the input cannot be normalized safely, stop and use the `invalid-region` template in `@references/common/error_handling.md`, including the official region links.
3. Check runtime prerequisites using a foolproof order. Prefer `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY`, `TENCENTCLOUD_REGION`, and optional `TENCENTCLOUD_SESSION_TOKEN`; also accept compatible names such as `MCP_REQUEST_SECRET_ID`, `MCP_REQUEST_SECRET_KEY`, `MCP_REQUEST_SESSION_TOKEN`, `MCP_SECRET_ID`, and `MCP_SECRET_KEY`. If required values are missing, stop and use the `missing-credentials` template in `@references/common/error_handling.md`, including one copyable example plus the official API key and region links.
4. Call Tencent Cloud PostgreSQL OpenAPI directly. Any aligned Action listed in @references/api_reference.md may be used. Start with read-only evidence collection first. For audit, business, fee-impacting, or critical write actions, explain impact clearly and require explicit confirmation before execution.
5. Prefer the official Tencent Cloud SDK. If the SDK is unavailable, do not block on that alone; generate properly TC3-signed HTTPS requests locally instead. Read credentials only from environment variables or other secure runtime context; never hardcode or persist them in code, skill files, or query parameters. If the user still wants the SDK path, follow the `missing-sdk` template in `@references/common/error_handling.md`, offer the install command, and ask whether to execute it.
6. Gather evidence in a fixed order:
   - instance basics, lifecycle status, spec, and exposure context
   - backup health and recovery readiness
   - parameter, security-group, SSL, and readonly topology evidence when relevant
   - optional slow-query or error-log evidence when the inspection request mentions latency or anomalies
7. Use @references/api_reference.md as the inspection checklist and aligned action catalog.
8. If an anomaly is found and the user asks for remediation, choose the smallest aligned action that can address the issue and explain prerequisites, impact, and confirmation requirement first.
9. Summarize in four blocks: overall status, evidence snapshot, risk items, and next actions.

## Pitfalls
- Do not report normal status without actual evidence from the selected OpenAPI actions.
- Do not mix different regions or instances in the same summary.
- Do not invent Cloud Monitor metrics or database-internal evidence that is not returned by the aligned OpenAPI actions.
- Do not perform fee-impacting or high-risk actions without explicit confirmation.

## Verification
- Include region and target instance in the final answer.
- Include at least one evidence block and one explicit risk level.
- State which OpenAPI actions were inspected, and list any action proposed or executed after confirmation.
- End with explicit next actions or `no action required` if no risk is found.
