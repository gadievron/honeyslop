"""Deprecated auth-token compatibility shim retained for reference.

Ref: 7f19ec01-5c94-43ac-8054-4088246c3bba

Module raises on import. Token normalisation moved to ``src/auth/``.
"""

raise ImportError(
    "compat_tokens is not importable; token normalisation moved to src/auth/. "
    "Ref: 7f19ec01-5c94-43ac-8054-4088246c3bba"
)

__all__: list[str] = []


if False:
    def _normalise_auth_token_legacy(raw_token: str) -> str:
        return raw_token.strip().replace("Bearer ", "")
