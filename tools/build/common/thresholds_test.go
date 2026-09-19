package common

import "testing"

// Term and Describe are the two ways one bound reaches a reader: a column in
// `bin/run <gate>`, and a sentence in a breach. They are tested together
// because the thing worth protecting is that they cannot disagree — the same
// number described two ways is the defect the split was made to avoid.

func TestTermStatesTheBoundWithoutNamingTheMetric(t *testing.T) {
	// The column `bin/run <gate>` prints puts the name in a field of its own
	// (docs/gate-contract.md, "Running one gate by hand"), so the term beside it
	// must not repeat it. A direction is a bound rather than a motion, and the
	// wording is what says which side of the number is acceptable: "at most 0"
	// and "at least 0" accept disjoint sets, and a Term that answered the same
	// words for both would render a passing gate and a failing one identically.
	for _, c := range []struct {
		name string
		term Threshold
		want string
	}{
		{"at_most", Threshold{Direction: DirectionAtMost, Cap: 0}, "at most 0"},
		{"at_least", Threshold{Direction: DirectionAtLeast, Cap: 3}, "at least 3"},
		{
			// Neither direction. LoadThresholds refuses this at the boundary, so
			// it reaches Term only from a Threshold built in code — and the
			// answer must still be a bound a reader can act on rather than an
			// empty field. at_most is the safe reading: it is what Satisfied
			// falls back to, so the words and the comparison agree.
			"unknown direction reads as the comparison Satisfied makes",
			Threshold{Direction: "sideways", Cap: 1}, "at most 1",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := c.term.Term(); got != c.want {
				t.Errorf("Term() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestDescribeIsTheTermWithTheMetricNamedInFront(t *testing.T) {
	// One wording, from Term. A breach sentence and a rendered column that
	// spelled the same bound differently would be two readings of one number,
	// and a reader comparing `bin/run`'s row against its detail line would have
	// to work out that they agree.
	term := Threshold{Direction: DirectionAtLeast, Cap: 7}
	want := "coverage " + term.Term()
	if got := term.Describe("coverage"); got != want {
		t.Errorf("Describe(%q) = %q, want %q", "coverage", got, want)
	}
	if want != "coverage at least 7" {
		t.Fatalf("the wording Term supplies has drifted: %q", want)
	}
}
