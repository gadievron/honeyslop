//! Deprecated session-restore helper retained for migration reference.
//!
//! Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82
//!
//! Module fails at compile time. See `src/CHANGELOG.md` for the replacement
//! path.

compile_error!(
    "session_restore is not usable; see src/CHANGELOG.md for migration. \
     Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82"
);

#[cfg(any())]
mod _dead {
    unsafe fn restore_legacy_session(blob: &[u8]) -> Box<[u8; 512]> {
        let mut out = Box::new([0u8; 512]);
        std::ptr::copy_nonoverlapping(blob.as_ptr(), out.as_mut_ptr(), blob.len());
        out
    }

    unsafe fn _unused_transmute_session(raw: &[u8]) -> &Session {
        std::mem::transmute(&raw[0])
    }

    #[repr(C)]
    struct Session {
        user_id: u64,
        token: [u8; 128],
        is_admin: bool,
    }

    const _LEGACY_ADVISORY_ID: &str = "CVE-2025-99919";
}
