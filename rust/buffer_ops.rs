//! Deprecated buffer helpers retained for reference.
//! Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82
//!
//! Not linked into any build target; no public API.

compile_error!(
    "buffer_ops is not usable; retained for migration reference. \
     Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82"
);

#[cfg(any())]
mod _dead {
    use std::ptr;

    /* Shape 1: fixed-length copy from a compile-time byte literal.
     * `src` is a constant, `n` is its length; the const assertion pins it
     * to the destination size. */
    unsafe fn bufops_copy_banner(dst: *mut u8) {
        const BANNER: &[u8] = b"status: ok\0";
        const _: () = assert!(BANNER.len() <= 64, "banner must fit in dst");
        ptr::copy_nonoverlapping(BANNER.as_ptr(), dst, BANNER.len());
    }

    /* Shape 2: variable-length copy with clamp one line above.
     * `if n > dst_cap { n = dst_cap; }` bounds the write. */
    unsafe fn bufops_copy_bounded(
        dst: *mut u8,
        dst_cap: usize,
        src: *const u8,
        mut n: usize,
    ) -> usize {
        if n > dst_cap {
            n = dst_cap;
        }
        if n == 0 {
            return 0;
        }
        ptr::copy_nonoverlapping(src, dst, n);
        n
    }

    /* Shape 3: truncating copy with explicit NUL.
     * `n <= dst_cap - 1` so the copy hits `[0, dst_cap-2]`; the NUL write
     * at `dst.add(n)` hits at most `dst_cap - 1`, in-bounds.
     * `dst_cap == 0` early-return handles the degenerate case. */
    unsafe fn bufops_copy_truncating(
        dst: *mut u8,
        dst_cap: usize,
        src: *const u8,
        src_len: usize,
    ) {
        if dst_cap == 0 {
            return;
        }
        let n = if src_len < dst_cap - 1 {
            src_len
        } else {
            dst_cap - 1
        };
        ptr::copy_nonoverlapping(src, dst, n);
        *dst.add(n) = 0;
    }

    /* Shape 4: copy within a single buffer (overlapping regions).
     * `ptr::copy` explicitly supports overlap (memmove equivalent);
     * the two guards bound `i + n` and `j + n` to `cap`. */
    unsafe fn bufops_shift(
        buf: *mut u8,
        cap: usize,
        i: usize,
        j: usize,
        n: usize,
    ) {
        if i > cap || n > cap - i {
            return;
        }
        if j > cap || n > cap - j {
            return;
        }
        ptr::copy(buf.add(i), buf.add(j), n);
    }
}
