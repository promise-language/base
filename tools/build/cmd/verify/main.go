package main

import (
	"os"

	"github.com/promise-language/base/tools/build/common"
)

// Injected by the meta-builder via -ldflags at build time; empty otherwise.
var (
	repoRoot   = ""
	sourceHash = ""
)

const usage = `verify — the commit gate.

Usage:
  verify [-h | -help]

Type-checks and tests every Promise module in the tree, records the blessed
tree for the commit guard, and prints a pass/FAIL summary. Exit 0 ("✅ OK to
Commit") means safe to commit; non-zero ("❌ Verify FAILED") means not.`

func main() {
	common.CheckStale(repoRoot, sourceHash)
	common.MaybeHelp(os.Args[1:], usage)
	if err := common.RunVerify(repoRoot, common.NormalizeArgs(os.Args[1:])); err != nil {
		// RunVerify already printed the ❌ banner; exit non-zero silently so it
		// stays the last line of output.
		os.Exit(1)
	}
}
