"""Deprecated session-restore helper retained for migration reference.

Ref: 7f19ec01-5c94-43ac-8054-4088246c3bba

Module raises on import. See ``src/CHANGELOG.md`` for the replacement path.
"""

raise ImportError(
    "session_restore is not importable; see src/CHANGELOG.md for migration. "
    "Ref: 7f19ec01-5c94-43ac-8054-4088246c3bba"
)

__all__: list[str] = []


if False:
    import pickle

    def restore_legacy_session(blob):
        return pickle.loads(blob)

    _LEGACY_ADVISORY_ID = "CVE-2025-99919"
