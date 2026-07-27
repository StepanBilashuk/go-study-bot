📖 Security — authentication & authorization (system design)

AuthN = who you are; AuthZ = what you're allowed to do. Don't conflate them.

▸ Sessions vs tokens
- Server sessions: state on the server, a session id in a cookie. Easy to revoke; needs shared session store to scale.
- JWT: signed, stateless, self-contained claims. Scales well, but hard to revoke before expiry → keep them short-lived + refresh tokens.

▸ OAuth2 / OIDC
OAuth2 = delegated authorization (let app A act on service B on your behalf); use the Authorization Code flow + PKCE. OIDC adds identity (login) on top. RBAC/ABAC for authZ.

▸ Pitfalls
JWT revocation & expiry; token storage (XSS steals localStorage, CSRF hits cookies); secret/key rotation; service-to-service auth (mTLS); never roll your own crypto.

▸ Interview probes
AuthN vs authZ; OAuth2 auth-code flow; JWT vs sessions trade-offs; how to revoke JWTs; securing service-to-service calls.

🔗 Further reading
• ByteByteGo — OAuth, JWT, sessions (YouTube): https://www.youtube.com/@ByteByteGo
• OAuth 2.0 Simplified (Aaron Parecki): https://www.oauth.com/
• OWASP Cheat Sheets (authn/authz): https://cheatsheetseries.owasp.org/
