package common

import "fmt"

// StaleReason returns a human-readable reason this binary is out of sync with
// its tools source, or "" if it is current. repoRoot and compiledHash are
// injected via -ldflags; empty values mean the binary was built some other way
// (go install, manual go build). It never exits — callers decide whether
// staleness is fatal (pipeline tools that would otherwise produce misleading
// results) or merely a warning (the git hook, which must never block a commit).
func StaleReason(repoRoot, compiledHash string) string {
	if repoRoot == "" || compiledHash == "" {
		return "this binary was not built via ./make"
	}
	currentHash, err := ToolsSourceHash(repoRoot)
	if err != nil {
		return fmt.Sprintf("binary's repo (%s) is unreachable: %v", repoRoot, err)
	}
	if compiledHash != currentHash {
		return "tools source has changed since this binary was built"
	}
	return ""
}

// MakeCmd is the bootstrap command to print in recovery hints.
func MakeCmd() string {
	if IsWindows() {
		return ".\\make.cmd"
	}
	return "./make"
}

// CheckStale refuses a tool whose stale logic would otherwise run: pipeline
// tools (verify, build, test, …) would produce misleading results, and the
// commit gate (precommit) must never validate a commit with out-of-date logic.
// The refusal names the condition and points the caller at ./make, and the tool
// exits 3 (Refuse).
//
// It is called BEFORE the command line is read, so what the tool was asked for
// does not decide whether it refuses — including --help and --version: "What a
// stale binary would print is the surface it was built with, and that is
// exactly what is out of date" (cli-guide.md#exit-codes). The one tool that
// never refuses on this ground is the builder the refusal names: ./make takes
// no CheckStale, because it is the way out.
//
// It is deliberately NOT a one-way door: the recovery, ./make, runs via 'go
// run' and has no staleness gate of its own, so it always works no matter how
// stale — or how broken — the compiled binaries are. Editing the tool source to
// fix a broken build is likewise permitted by the guard. So the way out is
// always fix-and-rebuild, never committing the broken state. Stale tools are a
// speed bump (re-run ./make), never a lockout.
func CheckStale(repoRoot, compiledHash string) {
	reason := StaleReason(repoRoot, compiledHash)
	if reason == "" {
		return
	}
	Refuse(Refusal{Condition: reason, Recovery: StaleRecovery(repoRoot)})
}

// StaleRecovery is what clears a staleness refusal: the builder, named with the
// tree to run it in when that is known. The directory is omitted rather than
// written empty when it is not — a binary built some other way knows no repo to
// point at, and "run ./make in " names nowhere.
func StaleRecovery(repoRoot string) string {
	if repoRoot == "" {
		return fmt.Sprintf("run %s", MakeCmd())
	}
	return fmt.Sprintf("run %s in %s", MakeCmd(), repoRoot)
}
