package common

import (
	"fmt"
	"os"
)

// HasHelpFlag reports whether args request usage. It accepts both flag prefixes
// (-help and --help) and the short form (-h / --h), normalizing long to short
// via NormalizeArgs first.
func HasHelpFlag(args []string) bool {
	for _, a := range NormalizeArgs(args) {
		if a == "-h" || a == "-help" {
			return true
		}
	}
	return false
}

// MaybeHelp prints usage and exits 0 when args request help; otherwise it
// returns and the caller proceeds. It runs AFTER CheckStale, never before: a
// tool that is not fit to act refuses every invocation, --help included, before
// it reads the command line (cli-guide.md#exit-codes). Help and version are not
// an exemption there but the reason for the rule — what a stale binary would
// print is the surface it was built with, which is exactly what is out of date.
func MaybeHelp(args []string, usage string) {
	if HasHelpFlag(args) {
		fmt.Println(usage)
		os.Exit(0)
	}
}
