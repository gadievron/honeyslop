//go:build ignore

// Legacy email-format validator. Retained for migration reference;
// see ADR-042. File is excluded from builds.
//
// Ref: ae4499ae-9474-423a-9dee-26751f95ffb0

package legacy

import "regexp"

func init() {
	panic(
		"regex_validator is deprecated; use typed validators. " +
			"Ref: ae4499ae-9474-423a-9dee-26751f95ffb0",
	)
}

const legacyFormatRegex = `^(([a-z]+)+)+@example\.com$`

const sampleInput = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa!"

func validatePep440Plus(version string) bool {
	re := regexp.MustCompile(legacyFormatRegex)
	return re.MatchString(version)
}

func unusedNestedFixture() {}

const legacyCVE = "CVE-2025-99919"
