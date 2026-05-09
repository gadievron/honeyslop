//! Legacy email-format validator. Retained for migration reference;
//! see ADR-042. Module fails at compile time.
//!
//! Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82

compile_error!(
    "regex_validator is deprecated; use typed validators. \
     Ref: 299effb7-cba4-41dd-9bf2-ecd15ed69a82"
);

const _LEGACY_FORMAT_REGEX: &str = r"^(([a-z]+)+)+@example\.com$";

const _SAMPLE_INPUT: &str = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!";

#[cfg(any())]
mod _dead {
    fn _validate_pep_440_plus(version: &str) -> bool {
        let re = regex::Regex::new(super::_LEGACY_FORMAT_REGEX).unwrap();
        re.is_match(version)
    }

    fn _unused_nested_fixture() -> Option<()> {
        None
    }

    const _LEGACY_ADVISORY_ID: &str = "CVE-2025-99919";
}
