package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Refusal is what a tool that is not fit to act answers with, and the one home
// of exit code 3 in this toolchain (cli-guide.md#exit-codes): "3 — the tool
// refused to act, because it is not fit to: built from source other than the
// tree beside it, not built by the project's builder, its own precondition
// unmet. Nothing was done, and the refusal names the condition and what clears
// it."
//
// Condition is why the tool will not act; Recovery is what clears it. Both are
// one sentence, because a caller reads the status to know it received no answer
// and the object to know why.
type Refusal struct {
	Condition string `json:"condition"`
	Recovery  string `json:"recovery"`
}

// Refuse writes r as the result and exits 3. It never returns.
//
// The refusal goes to STDOUT in the mode in force, not to stderr: "A sentinel at
// the head of stderr is prose, not a protocol: the first rewording breaks every
// matcher, and the matchers live in other repositories than the tool"
// (cli-guide.md#exit-codes). The workspace tools that read this project's
// toolchain are such matchers, and they live in another repository than this
// one — which is exactly why the shape is a document and not a prefix.
//
// Exit 3 rather than 1 because the four codes answer two questions
// (cli-guide.md#exit-codes): whether the subject was examined — 0 and 1 say yes,
// 2 and 3 say no — and whose repair it is: 1 the subject's, 2 the invocation's,
// 3 the installation's. A stale tool that exits 1 reports the subject as bad,
// and its caller names a repair that is not the repair.
func Refuse(r Refusal) {
	if err := WriteRefusal(os.Stdout, r, StdoutIsTerminal()); err != nil {
		// The object is the whole of the answer here. A refusal that cannot be
		// written is still a refusal, so the status stands either way, and the
		// one thing left to do is say on stderr that it could not be written.
		fmt.Fprintln(os.Stderr, "could not write the refusal:", err)
	}
	os.Exit(3)
}

// WriteRefusal writes r to w in one mode or the other. It is the whole of what
// a refusal looks like; Refuse is that plus the status, and the split is what
// lets the shape be asserted without ending the process.
func WriteRefusal(w io.Writer, r Refusal, human bool) error {
	if human {
		_, err := fmt.Fprintf(w, "%s — %s\n", r.Condition, r.Recovery)
		return err
	}
	return json.NewEncoder(w).Encode(r)
}

// StdoutIsTerminal reports whether stdout is a character device, which is the
// whole of the human-or-JSON decision (cli-guide.md#output-modes). Stdout alone
// decides: never stderr, never the environment.
//
// tools/build is a separate Go module from everything this repository ships,
// so the toolchain's dependencies never reach the product and ./make keeps
// working when the product does not compile. Base's product is Promise source,
// which cannot import this at all — so the separation here is structural rather
// than merely observed.
func StdoutIsTerminal() bool {
	st, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return st.Mode()&os.ModeCharDevice != 0
}
