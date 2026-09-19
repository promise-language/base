package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/promise-language/base/tools/build/common"
)

// beTheGate is how this suite gets a real gate process to watch.
//
// Everything the contract says about a bare invocation is about a process —
// what it wrote on which stream, and the status it left behind — so none of it
// can be asserted from inside the test that calls main(). The standard answer
// is the one taken here: the test binary re-executes itself and runs main()
// instead of the suite, so the child IS the gate, built from exactly this
// source.
//
// The marker is a command-line argument rather than an environment variable.
// A variable is inherited and hidden, and this repository's tools are held to
// "an environment variable is never an input"
// (docs/org/engineering-guide.md, "No hidden effects"); a test harness that
// reached for one anyway would be the first place someone looked for
// permission to. TestMain reads os.Args before the flag package does, so an
// argument the suite does not define never reaches it.
const beTheGate = "-be-the-gate"

func TestMain(m *testing.M) {
	if i := slices.Index(os.Args, beTheGate); i >= 0 {
		runAsGate(os.Args[i+1:])
	}
	os.Exit(m.Run())
}

// runAsGate stands this binary up as `bin/gate` and hands it the arguments the
// parent asked for. It never returns.
//
// repoRoot and sourceHash are what ./make injects via -ldflags, and main()
// refuses every invocation without them — CheckStale reads an unbuilt binary
// as one that is not fit to act, which is a refusal on stdout and status 3,
// and would make every assertion below about the staleness check instead of
// about the gate. They are computed from this tree rather than guessed, so
// they agree by construction.
func runAsGate(args []string) {
	root, err := repoRootFrom(mustGetwd())
	if err != nil {
		fmt.Fprintln(os.Stderr, "harness:", err)
		os.Exit(99)
	}
	hash, err := common.ToolsSourceHash(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "harness:", err)
		os.Exit(99)
	}
	repoRoot, sourceHash = root, hash
	os.Args = append([]string{"gate"}, args...)
	main()
	// main() answers every path with os.Exit, so arriving here is the harness
	// discovering that one of them stopped doing so rather than a success.
	os.Exit(98)
}

func mustGetwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// repoRootFrom walks up to the tree ToolsSourceHash hashes. `go test` runs a
// test binary in its package's source directory, so the walk starts somewhere
// inside the checkout and the answer is this checkout rather than whichever
// one a developer last cd'd into.
func repoRootFrom(dir string) (string, error) {
	for {
		if _, err := os.Stat(filepath.Join(dir, "tools", "build", "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no checkout above %s holds tools/build/go.mod", dir)
		}
		dir = parent
	}
}

func TestBareInvocationSaysTheThreeThingsTheContractAsksFor(t *testing.T) {
	// docs/gate-contract.md, "The exec line", states this surface as a worked
	// example:
	//
	//	$ bin/gate test
	//	test — measures test_count, test_failures, excluded_count
	//
	//	  bin/run test
	//
	// The name twice and the metrics between them, because each answers a
	// different question a person at a terminal has: which gate did I just ask
	// for, what would it have told me, and what do I type instead.
	got := bareInvocation("tested")
	want := "tested — measures " + strings.Join(common.GateMetrics(), ", ") + "\n\n  bin/run tested\n"
	if got != want {
		t.Errorf("bare invocation said\n%q\nwant\n%q", got, want)
	}
	// Against the document rather than against GateMetrics alone, which the
	// line above compares with itself.
	if !strings.HasPrefix(got, "tested — measures ") {
		t.Errorf("the gate's name and what it measures are not the first line: %q", got)
	}
	if !strings.HasSuffix(got, "\n\n  bin/run tested\n") {
		t.Errorf("the command that runs it is not offered as the document offers it: %q", got)
	}
}

func TestABareInvocationPrintsNothingOnStdoutAndExitsTwo(t *testing.T) {
	// The half of the rule that only a process can show, and the half the
	// contract calls load-bearing: "must produce nothing on stdout and a
	// non-zero exit ... a bare invocation that printed measurements and exited
	// 0 would be read as a pass by the first script that wrapped it".
	//
	// The status is 2 rather than merely non-zero because the four codes answer
	// two questions (cli-guide.md#exit-codes), and this invocation examined no
	// subject and is the caller's to repair. 1 would say the tree is bad, which
	// is what `bin/gate tested --envelope` says when it measures a failure — the
	// two must not arrive as one number.
	cmd := exec.Command(os.Args[0], beTheGate, "tested")
	var out, errs strings.Builder
	cmd.Stdout, cmd.Stderr = &out, &errs
	err := cmd.Run()

	code := 0
	if exit, isExit := err.(*exec.ExitError); isExit {
		code = exit.ExitCode()
	} else if err != nil {
		t.Fatalf("the harness could not run the gate: %v (%s)", err, errs.String())
	}
	if code == 98 || code == 99 {
		t.Fatalf("the harness failed rather than the gate: %s", errs.String())
	}

	if out.String() != "" {
		t.Errorf("a bare invocation wrote %q to stdout, which must carry an envelope or nothing", out.String())
	}
	if code != 2 {
		t.Errorf("a bare invocation exited %d, want 2", code)
	}
	if got := errs.String(); got != bareInvocation("tested") {
		t.Errorf("stderr carried\n%q\nwant\n%q", got, bareInvocation("tested"))
	}
}
