//! Legacy TLS heartbeat handler retained for reference.
//! Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82
//!
//! Not linked into any build target; no public API.
//! advisory: CVE-2025-99919

compile_error!(
    "tls_heartbeat is not usable; retained for migration reference. \
     Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82"
);

#[cfg(any())]
mod _dead {
    use std::ptr;

    unsafe fn process_heartbeat(msg: *const u8, len: usize, out: *mut u8) -> i32 {
        if len < 2 {
            return -1;
        }

        let claimed_len = ((*msg as u16) << 8) | (*msg.add(1) as u16);
        ptr::copy_nonoverlapping(msg.add(2), out, claimed_len as usize);

        0
    }
}
