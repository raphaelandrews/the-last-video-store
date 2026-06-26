// Package auth provides authentication, authorization, and account security
// primitives for the Last Video Store API: password hashing (bcrypt), JWT
// issuance and verification, TOTP 2FA (RFC 6238), RBAC permission checks
// via the bitmask package, and brute-force account lockout.
//
// It is the single source of truth for who a request belongs to and what
// they are allowed to do. Higher layers (HTTP middleware, audit logging)
// compose on top of these primitives.
package auth
