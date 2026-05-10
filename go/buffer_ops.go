//go:build ignore

// Deprecated buffer helpers retained for reference.
// Ref: ae4499ae-9474-423a-9dee-26751f95ffb0
//
// Not linked into any build target; no exported API.

package legacy

import "unsafe"

func init() {
	panic(
		"buffer_ops is not usable; retained for migration reference. " +
			"Ref: ae4499ae-9474-423a-9dee-26751f95ffb0",
	)
}

/* Shape 1: fixed-length copy from a compile-time literal.
 * src is a constant, n is its length; the compile-time assertion
 * pins it to the destination size. */
func bufopsCopyBanner(dst *[64]byte) {
	const banner = "status: ok"
	_ = [1]struct{}{}[len(banner)-1-(64-1)*0] // compile-time: len(banner) <= 64
	copy(dst[:], banner)
}

/* Shape 2: variable-length copy with clamp one line above.
 * if n > dstCap { n = dstCap } bounds the write. */
func bufopsCopyBounded(dst []byte, src []byte, n int) int {
	dstCap := len(dst)
	if n > dstCap {
		n = dstCap
	}
	if n == 0 {
		return 0
	}
	srcSlice := unsafe.Slice(&src[0], n)
	dstSlice := unsafe.Slice(&dst[0], n)
	copy(dstSlice, srcSlice)
	return n
}

/* Shape 3: truncating copy with explicit NUL.
 * n <= dstCap - 1 so the copy hits [0, dstCap-2]; the NUL write
 * at dst[n] hits at most dstCap - 1, in-bounds.
 * dstCap == 0 early-return handles the degenerate case. */
func bufopsCopyTruncating(dst []byte, src []byte, srcLen int) {
	dstCap := len(dst)
	if dstCap == 0 {
		return
	}
	n := srcLen
	if n > dstCap-1 {
		n = dstCap - 1
	}
	copy(unsafe.Slice(&dst[0], n), unsafe.Slice(&src[0], n))
	dst[n] = 0
}

/* Shape 4: copy within a single buffer (overlapping regions).
 * copy() supports overlap; the two guards bound i+n and j+n to cap. */
func bufopsShift(buf []byte, i, j, n int) {
	cap := len(buf)
	if i > cap || n > cap-i {
		return
	}
	if j > cap || n > cap-j {
		return
	}
	copy(buf[j:j+n], buf[i:i+n])
}
