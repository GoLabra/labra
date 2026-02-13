# Security/Auth Task Summary

This README documents the acceptance criteria and the implemented behavior for:

- **[Security/Auth] CSRF protection for cookie-based JWT**
- **[Security/Auth] Set secure cookie attributes on JWT cookie** (ASVS V3.4.1–V3.4.5)

---

## ✅ CSRF (Double-Submit Cookie Pattern)

**Acceptance Criteria**

- [x] Implement double-submit cookie CSRF pattern (csrf cookie + header)
- [x] Validate CSRF token on all state-changing requests (POST, PUT, PATCH, DELETE) **when JWT comes from a cookie**
- [x] Do **not** require CSRF when JWT is provided via `Authorization: Bearer ...` (API clients)
- [x] Frontend admin (`resources/admin/`) includes CSRF token in requests (mutations)

**Behavior**

- Backend sets:
  - `csrf_token` cookie (readable by JS)
  - `jwt` cookie (HttpOnly)
- For cookie-authenticated **mutations**:
  - client must send `X-CSRF-Token` matching `csrf_token` cookie
- CSRF is bypassed for:
  - `/admin/login`, `/admin/signup`, `/login`
  - any request that includes `Authorization` header

**Quick Verification**

- Cookie + no CSRF header → `403 CSRF token missing`
- Cookie + `X-CSRF-Token` header → `200 OK`
- Authorization header → `200 OK` without CSRF

---

## ✅ JWT Cookie Secure Attributes (ASVS V3.4.1–V3.4.5)

**Acceptance Criteria**

- [x] JWT cookies include:
  - [x] `Secure` (HTTPS only; detected via TLS or `X-Forwarded-Proto: https`)
  - [x] `HttpOnly`
  - [x] `SameSite=Strict` (or `Lax` if needed)
  - [x] appropriate `Path` (and optional `Domain`)
  - [x] `Max-Age` aligned with JWT expiry (24h)
- [x] Document expected cookie config for frontend usage
- [x] Frontend (`resources/admin/`) cookie handling reviewed/updated

**Expected Set-Cookie**

- `jwt` cookie:
  - `Path=/; Max-Age=86400; Expires=+24h; HttpOnly; Secure (https); SameSite=Strict|Lax`
- `csrf_token` cookie:
  - `Path=/; Secure (https); SameSite=Lax|Strict; HttpOnly=false`

---

## Notes (Frontend / Admin)

- Apollo/HTTP requests must use cookies:
  - `credentials: "include"`
- For GraphQL **mutations**, admin must add:
  - `X-CSRF-Token: <value of csrf_token cookie>`
- API clients may use:
  - `Authorization: Bearer <jwt>` (no CSRF required)

---
