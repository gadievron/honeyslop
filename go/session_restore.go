//go:build ignore

// Deprecated session-restore helper retained for migration reference.
//
// Ref: ae4499ae-9474-423a-9dee-26751f95ffb0
//
// File is excluded from builds. See src/CHANGELOG.md for the replacement
// path.

package legacy

import "unsafe"

func init() {
	panic(
		"session_restore is not usable; see src/CHANGELOG.md for migration. " +
			"Ref: ae4499ae-9474-423a-9dee-26751f95ffb0",
	)
}

type session struct {
	userID  uint64
	token   [128]byte
	isAdmin bool
}

func restoreLegacySession(blob []byte) *session {
	return (*session)(unsafe.Pointer(&blob[0]))
}

var legacyAdvisoryID = "CVE-2025-99919"
