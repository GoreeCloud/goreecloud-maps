# GoreeCloud Maps — Security

## Security status

GoreeCloud Maps is Forge/Development. This document defines repository-safe security expectations; it does not claim production security acceptance.

## Reporting vulnerabilities

Use GitHub private vulnerability reporting or another approved private GoreeCloud maintainer channel when available. Do not publish reusable credentials, exploit secrets, private user data, or sensitive infrastructure details in public issues.

## Authentication and authorization

- Browser authentication uses an Identity-compatible Authorization Code + PKCE public-client boundary when configured.
- No browser client secret is permitted.
- The API validates the authenticated subject and Maps applies resource-level authorization.
- Owner/editor/viewer roles must be enforced server-side.
- Protected PostgreSQL data uses row-level security where applicable.
- Runtime database roles must not own protected tables or have `BYPASSRLS`.

## Provider security

Provider endpoints are server-side configuration. Adapters must use bounded requests, validated coordinates/queries, finite timeouts and response sizes, controlled redirects/egress, and security-safe errors.

Production providers require explicit SSRF/DNS/egress review, credential controls, rate limits, monitoring, incident procedures, and rollback/provider-replacement evidence.

## Secret handling

Never commit access tokens, client secrets, provider credentials, private keys, recovery codes, passwords, production database URLs, or sensitive environment files.

Use sanitized `.env.example` files only for variable names and non-secret development defaults.

## Logging

Do not place bearer tokens, OIDC codes/verifiers, provider secrets, private route/Saved Place content, or raw precise coordinates in ordinary logs.

## Dependencies and supply chain

Pin and verify dependencies through the repository's package/module manifests and CI. Security updates must preserve exact-head validation; bypassing tests to accept a dependency update is not permitted.

## Platform boundaries

Wardveil Security remains the shared security/protection authority. GoreeCloud Identity owns principal/session identity. GoreeCloud Policy owns policy decisions. Glaze UI only presents state and cannot grant trust or authorization.

## Production gate

Security review, vulnerability handling, production Identity/database/provider acceptance, Observability evidence, continuity/restore, deployment, rollback, and release qualification remain separate requirements before Anchor.
