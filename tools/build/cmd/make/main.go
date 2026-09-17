// Command make is the meta-builder. It compiles every other tool under cmd/
// into <repoRoot>/bin, stamping each binary with the tools-source hash and the
// absolute repo root via -ldflags. It is the one tool that runs via 'go run'
// (from the ./make trampoline), so it is never compiled into bin/ and never
// stale — which is what breaks the bootstrap cycle.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/promise-language/base/tools/build/common"
)

const usage = `make — the meta-builder.

Usage:
  ./make [-force | --force] [-h | -help]

Compiles every tool under tools/build/cmd into bin/ (stamping each with the
tools-source hash and repo root) and wires git hooks. Skips the build when
bin/ is already up to date; -force rebuilds regardless.

It does NOT build this repository's Promise source: that is what bin/gate and
bin/run measure, and a meta-builder that also built the product would fail with
it — taking the toolchain down with the bug being fixed.`

func main() {
	common.MaybeHelp(os.Args[1:], usage)
	force := false
	for _, a := range os.Args[1:] {
		if a == "-force" || a == "--force" {
			force = true
		}
	}

	// 1. Resolve the repo root. The ./make trampoline cd'd go run into
	//    <root>/tools/build, so our cwd is exactly that. Two levels up is root.
	cwd, err := os.Getwd()
	must(err)
	repoRoot := filepath.Dir(filepath.Dir(cwd))
	if !filepath.IsAbs(repoRoot) {
		fail("resolved repo root is not absolute: %s", repoRoot)
	}

	// 2. Hash the tools source — baked into every binary below.
	hash, err := common.ToolsSourceHash(repoRoot)
	must(err)

	// 3. Enable git hooks unconditionally (idempotent, fast).
	if err := common.RunSetup(repoRoot); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not configure git hooks: %v\n", err)
	}

	tools, err := common.CommandNames(repoRoot)
	must(err)

	binDir := filepath.Join(repoRoot, "bin")
	hashFile := filepath.Join(binDir, ".tools.hash")

	// 4. Up-to-date short circuit.
	if !force && upToDate(hashFile, hash, binDir, tools) {
		fmt.Println("Tools up to date")
		return
	}

	must(os.MkdirAll(binDir, 0o755))

	// 5. Build each tool, injecting repoRoot and sourceHash via ldflags.
	ldflags := fmt.Sprintf("-s -w -X main.sourceHash=%s -X main.repoRoot=%s", hash, repoRoot)
	toolsModDir := filepath.Join(repoRoot, "tools", "build")
	for _, name := range tools {
		out := filepath.Join(binDir, common.BinaryName(name))
		fmt.Printf("building %s\n", name)
		if err := common.RunIn(toolsModDir, "go", "build",
			"-trimpath",
			"-ldflags", ldflags,
			"-o", out,
			"./cmd/"+name,
		); err != nil {
			fail("building %s: %v", name, err)
		}
	}

	// 6. Write the hash sidecar — the staleness contract.
	must(os.WriteFile(hashFile, []byte(hash+"\n"), 0o644))
	fmt.Printf("built %d tool(s) into bin/\n", len(tools))
}

func upToDate(hashFile, hash, binDir string, tools []string) bool {
	data, err := os.ReadFile(hashFile)
	if err != nil || strings.TrimSpace(string(data)) != hash {
		return false
	}
	for _, name := range tools {
		if !common.Exists(filepath.Join(binDir, common.BinaryName(name))) {
			return false
		}
	}
	return true
}

func must(err error) {
	if err != nil {
		fail("%v", err)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "make: "+format+"\n", args...)
	os.Exit(1)
}
