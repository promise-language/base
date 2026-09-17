package common

import (
	"fmt"
	"os"
	"time"
)

type step struct {
	name string
	run  func(repoRoot string) error
}

// RunVerify is the commit gate. It always prints a summary block — even on
// failure — so an agent tailing the output sees the result without re-running,
// and the process exit code is the only contract.
//
// Verify is a COMMAND, not a gate: it is allowed to repair what has exactly one
// right answer on its way to an answer, which is what a producing step wants,
// and exactly why a landing decision rests on `integration` instead.
//
// Base has nothing to repair. The repairing step every sibling runs here is a
// formatter, and the Promise toolchain ships none (see gate.go), so this
// pipeline measures and records. The step stays a named stage rather than being
// dropped, so the day `promise fmt` exists there is one obvious place for it.
func RunVerify(repoRoot string, args []string) error {
	// A stale blessing left behind is the one outcome the verified-tree check
	// must never produce, so failing to clear fails the run outright.
	if err := clearVerifiedTree(repoRoot); err != nil {
		return fmt.Errorf("clearing %s: %w", verifiedTreeRecord, err)
	}
	return runVerifySteps(repoRoot, verifyPipeline(repoRoot))
}

// runVerifySteps runs the steps in order, stopping at the first failure, and
// always prints the summary block.
func runVerifySteps(repoRoot string, steps []step) error {
	start := time.Now()

	type result struct {
		name string
		ok   bool
	}
	var results []result
	failed := false

	for _, s := range steps {
		fmt.Printf("==> %s\n", s.name)
		err := s.run(repoRoot)
		results = append(results, result{s.name, err == nil})
		if err != nil {
			failed = true
			fmt.Fprintf(os.Stderr, "    %s failed: %v\n", s.name, err)
			break // stop at the first failure
		}
	}

	fmt.Println("\n──────── verify summary ────────")
	for _, r := range results {
		status := "ok"
		if !r.ok {
			status = "FAIL"
		}
		fmt.Printf("  %-4s  %s\n", status, r.name)
	}
	fmt.Printf("  elapsed %s\n", time.Since(start).Round(time.Millisecond))
	fmt.Println("────────────────────────────────")

	if failed {
		fmt.Println("❌ Verify FAILED: not safe to commit")
		return fmt.Errorf("verify failed")
	}
	fmt.Println("✅ OK to Commit")
	return nil
}

// verifyPipeline is the full run: the measurement, then the unconditional
// trailing record step. Being a step gets break-on-first-failure for free — a
// red step leaves nothing blessed.
func verifyPipeline(repoRoot string) []step {
	return []step{
		{"measure", RunMeasurement},
		{"record", recordVerifiedTree},
	}
}
