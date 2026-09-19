package main

import (
	"maps"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/promise-language/base/tools/build/common"
)

func TestMeasurementsInReadsTheSameNumbersFromEitherSideOfTheWire(t *testing.T) {
	// The property the shared judge exists for, stated as a test: one envelope
	// read two ways answers the same numbers. --verdict decodes JSON, where an
	// object is a map[string]any of float64; the by-hand mode is handed the map
	// MeasureGate built and never encodes it, so it arrives as map[string]int64.
	// A reader that knew only the decoded shape found nothing in the second —
	// silently, because a failed type assertion is an empty map rather than an
	// error — and the judge fell back to reconstructing a number the envelope
	// was holding all along.
	//
	// Two metrics, because one is exactly the case the fallback hides: with
	// only failed_gates present, a reader that saw nothing still reached the
	// right verdict, and the gap was invisible.
	typed := map[string]any{"measurements": map[string]int64{
		common.MetricFailedGates: 1,
		"unformatted_files":      3,
	}}
	decoded := map[string]any{"measurements": map[string]any{
		common.MetricFailedGates: float64(1),
		"unformatted_files":      float64(3),
	}}
	want := map[string]int64{common.MetricFailedGates: 1, "unformatted_files": 3}

	fromHand := measurementsIn(typed)
	if !maps.Equal(fromHand, want) {
		t.Errorf("the by-hand mode read %v from the envelope it was handed, want %v", fromHand, want)
	}
	fromWire := measurementsIn(decoded)
	if !maps.Equal(fromWire, want) {
		t.Errorf("--verdict read %v from the decoded envelope, want %v", fromWire, want)
	}
	if !maps.Equal(fromHand, fromWire) {
		t.Errorf("one envelope, two readings: %v and %v", fromHand, fromWire)
	}
}

func TestMeasurementsInLeavesOutWhatItCannotReadAsANumber(t *testing.T) {
	// Absent is not zero. A term whose metric nobody measured is skipped by
	// ApplyThresholds rather than applied to an invented number, so a value
	// this cannot read has to be left out rather than defaulted in — a
	// measurement of 0 against a cap of 0 is a pass nobody measured, which is
	// the one reading the whole judging path exists to prevent.
	for _, c := range []struct {
		name     string
		envelope map[string]any
	}{
		{"no measurements key at all", map[string]any{"measured": true}},
		{"measurements is not an object", map[string]any{"measurements": "several"}},
		{"a measurement is not a number", map[string]any{"measurements": map[string]any{
			common.MetricFailedGates: "quite a few",
		}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := measurementsIn(c.envelope); len(got) != 0 {
				t.Errorf("read %v, want nothing at all", got)
			}
		})
	}
}

func TestJudgedPutsEachMeasurementBesideTheTermItWasJudgedOn(t *testing.T) {
	// docs/gate-contract.md, "Running one gate by hand": the judging layer is
	// the only layer that can render this, because it holds the caps and the
	// directions. What a reader acts on is the mark, so the mark is what this
	// asserts hardest: a row carrying the right number against the right bound
	// and the wrong tick is a failing gate that reads as green.
	v := verdict{Thresholds: map[string]common.Threshold{
		"unformatted_files":      {Direction: common.DirectionAtMost, Cap: 0},
		common.MetricFailedGates: {Direction: common.DirectionAtMost, Cap: 0},
		"coverage":               {Direction: common.DirectionAtLeast, Cap: 90},
	}}
	measurements := map[string]int64{
		"unformatted_files":      3,
		common.MetricFailedGates: 0,
		"coverage":               90,
	}

	rows := judged(v, measurements)
	if len(rows) != 3 {
		t.Fatalf("judged rendered %d rows for 3 terms: %q", len(rows), rows)
	}

	// Sorted, because the map is not: an order that changed between two runs of
	// the same gate would read as a change in what was measured.
	wantOrder := []string{"coverage", common.MetricFailedGates, "unformatted_files"}
	for i, name := range wantOrder {
		if !strings.HasPrefix(rows[i], name) {
			t.Errorf("row %d is %q, want it to begin with %q", i, rows[i], name)
		}
	}

	// coverage measured exactly its at_least bound: inclusive, so it passes.
	// failed_gates measured 0 against at most 0: the same boundary from the
	// other direction. unformatted_files is the one breach.
	for _, c := range []struct {
		row  string
		want string
		mark string
	}{
		{rows[0], "at least 90", "✓"},
		{rows[1], "at most 0", "✓"},
		{rows[2], "at most 0", "✗"},
	} {
		if !strings.Contains(c.row, c.want) {
			t.Errorf("row %q does not state the term %q it was judged on", c.row, c.want)
		}
		if !strings.HasSuffix(c.row, c.mark) {
			t.Errorf("row %q does not end in %s", c.row, c.mark)
		}
	}

	// The number is the one measured, not the one in the term. A renderer that
	// printed the cap twice would satisfy every assertion above.
	if !strings.Contains(rows[2], "3") {
		t.Errorf("row %q does not carry the 3 that was measured", rows[2])
	}

	// The widths come from the rows themselves, so a metric whose name is
	// wider than the rest lines up instead of pushing its own row out of the
	// column. Equal rendered width across rows is what that amounts to.
	width := utf8.RuneCountInString(rows[0])
	for _, row := range rows[1:] {
		if utf8.RuneCountInString(row) != width {
			t.Errorf("rows do not line up:\n%s", strings.Join(rows, "\n"))
			break
		}
	}
}
