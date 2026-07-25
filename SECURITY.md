# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| v1.x | ✅ Active |
| v0.x | ⚠️ Best effort |

---

## Reporting a Vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

To report a vulnerability:

1. Open a GitHub Security Advisory (preferred):
   Security tab → Report a vulnerability

2. Or email directly with subject line: [ATLAS SECURITY] <brief description>

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

You will receive a response within **48 hours**.

---

## Security Design

Atlas is a **single-user, self-hosted application**. The threat model assumes:

- The server is running on a trusted machine or private network
- The attacker does not have physical access to the machine
- The user controls their own deployment

### Measures in place

- bcrypt password hashing (cost ≥ 12)
- Session-based auth with HTTP-only, SameSite=Lax cookies
- CSRF protection on all state-changing requests
- Rate limiting on the login endpoint (≤ 5 attempts/min)
- Parameterised SQL queries via sqlc (no SQL injection surface)
- File upload validation (MIME type + max size + random filename)
- No public registration endpoint after first-run wizard completes
- All routes except /setup, /login, /static require authentication

### What Atlas does NOT protect against

- Physical access to the host machine
- Compromised host OS
- Network-level attacks (use a VPN or Caddy with TLS in production)

---

## Responsible Disclosure

We follow coordinated disclosure:

1. Report privately
2. We investigate and develop a fix
3. Fix is released
4. Disclosure is published (crediting the reporter if desired)

Thank you for helping keep Atlas secure.
