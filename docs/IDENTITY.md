# GoreeCloud Maps Identity Integration

## Status

This document defines the current **source-level** GoreeCloud Identity integration boundary for GoreeCloud Maps. GoreeCloud Identity is not yet approved for GoreeCloud-wide production SSO, and Maps does not claim production Identity acceptance from this implementation alone.

## Responsibility boundary

GoreeCloud Identity establishes and authenticates the principal. GoreeCloud Maps remains authoritative for application-owned resource authorization, including collection ownership, owner/editor/viewer roles, saved data, and collaboration permissions.

Maps must not create a parallel permanent account system or treat network access, possession of a resource identifier, or an Identity group claim as automatic authorization to a Maps resource.

## Browser authentication flow

The web application uses an optional OpenID Connect public-client integration:

- Authorization Code flow with PKCE (`S256`);
- no browser client secret;
- OIDC discovery from an explicitly configured issuer;
- explicit state validation on the authorization callback;
- same-origin redirect URI enforcement;
- HTTPS Identity endpoints outside localhost development;
- transient PKCE verifier/state in `sessionStorage` only;
- bearer access token held in application memory rather than persistent browser storage;
- no Identity discovery or authorization request when issuer/client registration is not configured.

Web configuration:

- `VITE_IDENTITY_ISSUER_URL`
- `VITE_IDENTITY_CLIENT_ID`
- `VITE_IDENTITY_REDIRECT_URI` (optional; defaults to the current Maps origin/path)
- `VITE_IDENTITY_SCOPES` (defaults to `openid profile`)

The web client must never receive or embed an OIDC client secret.

Disconnecting the Maps client currently discards its in-memory bearer token and transient PKCE state. It does **not** claim to terminate the provider-wide SSO session. Provider logout/session-management behavior remains part of later GoreeCloud Identity application-integration acceptance.

## Maps API bearer contract

Protected `/api/v1` routes expect an OAuth/OIDC bearer **access token**. The current API verifier:

1. validates the signed JWT against the configured OIDC issuer and Maps client audience, including standard token expiry checks supplied by the verifier library;
2. requires a non-empty subject;
3. presents the bearer token to the provider UserInfo endpoint;
4. requires UserInfo to accept the token and return the same non-empty subject;
5. resolves that subject to the internal Maps user record;
6. performs all Maps resource authorization independently through the application and PostgreSQL RLS boundaries.

The UserInfo confirmation makes the current source integration explicit about accepting a provider-recognized access token rather than relying on the old ID-token-as-bearer assumption. This network validation path is a source foundation, not yet a production latency/availability architecture decision.

API configuration:

- `MAPS_OIDC_ISSUER_URL`
- `MAPS_OIDC_CLIENT_ID`

No reusable Identity token, signing secret, or client secret belongs in source control.

## Same-origin application API

The web application calls Maps through a same-origin path, defaulting to `/api/v1`. `VITE_MAPS_API_BASE_PATH` must be an absolute path beginning with `/`; external API origins are rejected by the browser client configuration.

This keeps the initial browser deployment model compatible with a controlled reverse proxy while avoiding arbitrary runtime API destinations. Production proxy, CSP, CORS, cookie, header, network, and Cloudflare behavior still require deployment-specific acceptance.

## Recipient discovery and invitations

Human-friendly recipient discovery is **blocked by prerequisite** until GoreeCloud Identity defines and approves a cross-application directory/invitation contract.

Maps will not use an administrative authentik user-list API as an accidental consumer directory. Doing so would couple Maps to provider administration, risk over-disclosing identity attributes, and bypass the required GoreeCloud privacy/governance boundary.

The existing collection membership API may add an already-known internal Maps user ID for controlled integration/testing. This is not a completed invitation experience.

A future approved recipient contract should define at least:

- which identity attributes Maps may query or display;
- search/discovery privacy and enumeration controls;
- whether the recipient must already have a Maps user mapping;
- invitation creation, delivery, expiry, acceptance, decline, and revocation;
- anti-abuse/rate-limit behavior;
- auditing and Privacy Shield requirements;
- disabled/deleted identity behavior;
- group/household semantics if supported;
- compatibility with GoreeCloud Mesh/Notify where those systems own delivery or coordination responsibilities.

Until that contract exists, the Maps UI must state that invitations are unavailable rather than inventing a directory.

## Current acceptance boundary

Source implementation and repository CI can validate compilation, application authorization, and PostGIS isolation. Production Identity acceptance remains separate and requires an actual GoreeCloud Identity Maps client registration plus end-to-end evidence for login, callback, access-token validation, user mapping, logout/session behavior, expiration, disabled-account behavior, Identity outage behavior, recovery, rollback, privacy, and security.
