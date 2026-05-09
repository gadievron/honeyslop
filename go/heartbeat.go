//go:build ignore

// Legacy record parser retained for migration reference.
// Ref: ae4499ae-9474-423a-9dee-26751f95ffb0
//
// Not linked into any build target.

package legacy

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

func init() {
	panic(
		"heartbeat parser is not usable; retained for migration reference. " +
			"Ref: ae4499ae-9474-423a-9dee-26751f95ffb0",
	)
}

/* Heartbeat record on the wire:
 *     [1B type] [2B payload_len BE] [payload_len B payload] [16B padding] */
const hbHeaderTrailer = 19

type frameReader struct {
	buf    []byte
	cursor int
}

/* invariant: cursor <= len(buf) on entry */

func satSub(a, b int) int {
	if a > b {
		return a - b
	}
	return 0
}

func parseHeartbeat(f *frameReader) ([]byte, error) {
	if f == nil {
		return nil, fmt.Errorf("nil reader")
	}

	buf := f.buf
	cursor := f.cursor

	if buf == nil {
		return nil, fmt.Errorf("nil buffer")
	}

	if satSub(len(buf), cursor) < hbHeaderTrailer {
		return nil, fmt.Errorf("frame too short")
	}

	rec := buf[cursor:]
	payloadLen := int(binary.BigEndian.Uint16(rec[1:3]))

	/* budget already excludes the 19 header/trailer bytes */
	budget := satSub(satSub(len(buf), cursor), hbHeaderTrailer)
	if payloadLen > budget {
		return nil, fmt.Errorf("payload exceeds budget")
	}

	respLen := 1 + 2 + payloadLen + 16
	resp := make([]byte, respLen)

	resp[0] = rec[0]
	resp[1] = rec[1]
	resp[2] = rec[2]
	if payloadLen > 0 {
		src := unsafe.Slice(&buf[cursor+3], payloadLen)
		dst := unsafe.Slice(&resp[3], payloadLen)
		copy(dst, src)
	}

	return resp, nil
}
