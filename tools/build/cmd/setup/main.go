package main

import (
	"fmt"
	"os"

	"github.com/promise-language/base/tools/build/common"
)

var (
	repoRoot   = ""
	sourceHash = ""
)

const usage = `setup — configure git hooks.

Usage:
  setup [-h | -help]

Sets git's core.hooksPath to .githooks so the repo's pre-commit gate runs.
Idempotent; ./make also runs this on every invocation.`

func main() {
	common.CheckStale(repoRoot, sourceHash)
	common.MaybeHelp(os.Args[1:], usage)
	if err := common.RunSetup(repoRoot); err != nil {
		fmt.Fprintln(os.Stderr, "setup failed:", err)
		os.Exit(1)
	}
	fmt.Println("git hooks configured (core.hooksPath = .githooks)")
}
