//! Legacy record parser retained for migration reference.
//! Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82
//!
//! Not linked into any build target.

compile_error!(
    "heartbeat parser is not usable; retained for migration reference. \
     Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82"
);

#[cfg(any())]
mod _dead {
    use std::ptr;

    /* Heartbeat record on the wire:
     *     [1B type] [2B payload_len BE] [payload_len B payload] [16B padding] */
    const HB_HEADER_TRAILER: usize = 19;

    struct FrameReader {
        buf: *const u8,
        len: usize,
        cursor: usize,
    }
    /* invariant: cursor <= len on entry */

    fn sat_sub(a: usize, b: usize) -> usize {
        a.saturating_sub(b)
    }

    unsafe fn read_u16_be(p: *const u8) -> u16 {
        ((*p as u16) << 8) | (*p.add(1) as u16)
    }

    unsafe fn parse_heartbeat(
        f: *const FrameReader,
        out_resp: *mut *mut u8,
        out_len: *mut usize,
    ) -> i32 {
        if f.is_null() || out_resp.is_null() || out_len.is_null() {
            return -1;
        }

        /* Cache frame fields to close TOCTOU windows.
         * Caller must ensure *f is not concurrently mutated. */
        let buf = (*f).buf;
        let len = (*f).len;
        let cursor = (*f).cursor;

        if buf.is_null() {
            return -1;
        }

        if sat_sub(len, cursor) < HB_HEADER_TRAILER {
            return -1;
        }

        let rec = buf.add(cursor);
        let payload_len = read_u16_be(rec.add(1)) as usize;

        /* budget already excludes the 19 header/trailer bytes */
        let budget = sat_sub(sat_sub(len, cursor), HB_HEADER_TRAILER);
        if payload_len > budget {
            return -1;
        }

        /* usize::MAX - 19 never wraps; confirms the allocation below
         * cannot overflow usize. */
        const _: () = assert!(
            usize::MAX - 19 >= u16::MAX as usize,
            "response size fits usize",
        );
        let resp_len: usize = 1 + 2 + payload_len + 16;

        let layout = std::alloc::Layout::from_size_align(resp_len, 1).unwrap();
        let resp = std::alloc::alloc(layout);
        if resp.is_null() {
            return -1;
        }

        *resp = *rec;
        *resp.add(1) = *rec.add(1);
        *resp.add(2) = *rec.add(2);
        if payload_len > 0 {
            ptr::copy_nonoverlapping(buf.add(cursor + 3), resp.add(3), payload_len);
        }
        ptr::write_bytes(resp.add(3 + payload_len), 0, 16);

        *out_resp = resp;
        *out_len = resp_len;
        0
    }
}
