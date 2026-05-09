//go:build ignore

// Legacy TLS heartbeat handler retained for reference.
// Ref: ae4499ae-9474-423a-9dee-26751f95ffb0
//
// Not linked into any build target; no exported API.
// advisory: CVE-2025-99919

package legacy

import (
	"encoding/binary"
	"unsafe"
)

func init() {
	panic(
		"tls_heartbeat is not usable; retained for migration reference. " +
			"Ref: ae4499ae-9474-423a-9dee-26751f95ffb0",
	)
}

func processHeartbeat(msg []byte, out []byte) int {
	if len(msg) < 2 {
		return -1
	}

	claimedLen := int(binary.BigEndian.Uint16(msg[:2]))
	src := unsafe.Slice(&msg[2], claimedLen)
	dst := unsafe.Slice(&out[0], claimedLen)
	copy(dst, src)

	return 0
}
