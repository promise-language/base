package common

import (
	"slices"
	"testing"
)

func TestGateMetricsNamesExactlyWhatAnEnvelopeCarries(t *testing.T) {
	// The tripwire the de-duplication is for. GateMetrics answers a bare
	// invocation's "what does this gate measure" (docs/gate-contract.md, "The
	// exec line"), and the envelope answers the judge — so a metric present in
	// one and not the other is a gate telling a person it measures less, or
	// more, than the numbers it reports. Two literals could drift silently;
	// this fails the moment they do.
	//
	// The gate asked about is one that cannot exist, which RunGate refuses
	// before it reaches a module, a file or the PATH. What it measured is not
	// what this is about — only which keys came back — and a gate that really
	// ran would make this test depend on the toolchain being installed.
	env, err := MeasureGate("", "no-such-concept")
	if err == nil {
		t.Fatal("a gate this project does not provide must not measure")
	}
	measurements, ok := env["measurements"].(map[string]int64)
	if !ok {
		t.Fatalf("the envelope carries no measurements map: %#v", env["measurements"])
	}

	carried := make([]string, 0, len(measurements))
	for name := range measurements {
		carried = append(carried, name)
	}
	slices.Sort(carried)

	if named := GateMetrics(); !slices.Equal(named, carried) {
		t.Errorf("a bare invocation says the gate measures %v; its envelope carries %v", named, carried)
	}
}

func TestGateMetricsIsOrderedSoTwoBareInvocationsReadTheSame(t *testing.T) {
	// The names come out of a map, which Go deliberately does not order. A
	// bare invocation is the one place they are printed as prose, so an order
	// that changed between two runs would read as a gate that measures
	// something different each time — and would do it on the human surface,
	// where nothing else is checking.
	first := GateMetrics()
	if !slices.IsSorted(first) {
		t.Errorf("GateMetrics() = %v, which is not in a settled order", first)
	}
	for i := 0; i < 32; i++ {
		if again := GateMetrics(); !slices.Equal(again, first) {
			t.Fatalf("GateMetrics() answered %v and then %v", first, again)
		}
	}
}

func TestFailedGatesIsOneOrZeroAccordingToTheGatesOwnError(t *testing.T) {
	// The metric measured by the runner rather than by a gate, and the only
	// number base's judge holds anything to. It is 0 or 1 because it is about
	// the named gate as a whole — `integration` is one gate, not two — so a
	// count of underlying failures here would be a cap of 0 that a composition
	// could never meet.
	if got := gateMeasurements(nil)[MetricFailedGates]; got != 0 {
		t.Errorf("a gate that returned no error measured %d failed gates, want 0", got)
	}
	if got := gateMeasurements(errFixture{})[MetricFailedGates]; got != 1 {
		t.Errorf("a gate that failed measured %d failed gates, want 1", got)
	}
}

type errFixture struct{}

func (errFixture) Error() string { return "the gate said no" }
