package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Gates: measurements a decision may rest on.
//
// A gate MEASURES and modifies nothing it measures — including afterwards.
// Measuring faithfully and then tidying up is not a gate, because the answer
// then describes a tree that no longer exists. That is what separates these
// from `verify`, which is a COMMAND: verify repairs what has one right answer
// on its way to an answer, which is exactly what a producing step wants and
// exactly why a landing decision may not rest on it.
//
// The names are the flow SDK's closed vocabulary — a concept, optionally with
// an instance after a colon (`tested:wire`). Gates are addressed by name rather
// than configured as command strings, which is what lets a step ask for the one
// module it broke instead of paying for the whole set, and what keeps this from
// becoming an arbitrary command executor.
//
// # What base does NOT provide, and why
//
// A project need not provide every concept, and inventing one it cannot measure
// honestly is worse than omitting it. Three are deliberately absent here, and
// each absence is a property of the toolchain rather than a preference:
//
// `builds` — `promise build` requires an entry point and refuses a module
// without one ("program has no main() function"). Every module here is a
// library of contracts; the single `main()` in the tree, in
// tests/wire-consumer, is a `test` main. A `builds` gate would therefore fail
// on a correct tree forever, which is the one thing a gate must never do.
//
// `checked` — `promise check` takes exactly one file and type-checks it alone,
// with none of its module's siblings in scope. Asked of a multi-file module it
// reports every type defined in a neighbouring file as undefined, so it fails
// on a module that compiles perfectly. There is no module-scoped form of it:
// `promise doc` walks a whole module but exits 0 on one that does not compile,
// so it answers a different question. `tested` is base's compile question for
// that reason — `promise test` compiles the module as a unit, which is the only
// granularity the toolchain type-checks at.
//
// `formatted` — `promise format -check` exists and would answer, but the
// formatter it checks against disagrees with this tree in 746 lines, and the
// disagreements are not neutral: it flattens the aligned continuation lines of
// a multi-line call to a fixed two-space indent, which is a reformatting that
// loses information a reader uses. Adopting it is a decision about base's house
// style, not a gate to switch on, and it belongs in a change that makes that
// decision on purpose.

// GateConcept is the part of a gate name before the colon.
type GateConcept string

const (
	// GateTested asks whether the module compiles and its tests pass.
	//
	// Both, because `promise test` is the only entry point that type-checks a
	// module as a unit: a compilation error surfaces as a failing test file, and
	// there is no cheaper form that would have caught it first (see above).
	GateTested GateConcept = "tested"

	// GateIntegration is the composition a landing decision rests on: every
	// other gate this project provides, in the cheapest-first order that fails
	// fastest.
	GateIntegration GateConcept = "integration"

	// GateFit asks whether this MACHINE may be given work at all — the one gate
	// here whose subject is not the code.
	//
	// Required, alongside integration: flow's gates-and-commands.md lists
	// `verify`, `integration`, `fit` and the judge as the four things a flow
	// cannot run without. A project that does not answer it is not merely
	// missing a measurement — `issue resolve` refuses the project outright and
	// waits for a condition that no amount of waiting clears.
	//
	// It is deliberately NOT part of integration. A machine that lacks the
	// Promise toolchain is not a change that may not land, and folding the two
	// would fail an honest change for a fact about the host it ran on.
	GateFit GateConcept = "fit"
)

// RunGate runs the named gate and returns nil when it passes.
//
// name is a flow gate name: a concept, optionally `concept:instance`. An
// instance selects one Promise module; omitting it runs every module — the safe
// reading, since a gate that quietly measured a subset would report a tree
// sound while part of it was unmeasured.
func RunGate(repoRoot, name string) error {
	concept, instance := splitGateName(name)

	if concept == GateIntegration {
		if instance != "" {
			return fmt.Errorf("gate %q: integration takes no instance — it is the composition of the others", name)
		}
		return runIntegration(repoRoot)
	}

	if concept == GateFit {
		if instance != "" {
			return fmt.Errorf("gate %q: fit takes no instance here — its instances would name conditions "+
				"on the machine (`fit:toolchain`), not modules of this project", name)
		}
		return runFit(repoRoot)
	}

	run, ok := gateRunners()[concept]
	if !ok {
		return fmt.Errorf("gate %q: unknown concept %q (this project provides: %s)",
			name, concept, strings.Join(providedConcepts(), ", "))
	}
	mods, err := modulesFor(repoRoot, instance)
	if err != nil {
		return fmt.Errorf("gate %q: %w", name, err)
	}
	for _, m := range mods {
		if err := run(repoRoot, m); err != nil {
			return fmt.Errorf("gate %q failed in %s: %w", name, ModuleLabel(m), err)
		}
	}
	return nil
}

// runFit answers whether this machine may be given work.
//
// Base's one host requirement is the Promise toolchain: every gate below spawns
// `promise`, so a machine without it cannot measure anything at all, and an
// arena given work on such a host would fail every step for a reason that has
// nothing to do with the change. Asking here — once, cheaply — is what turns
// that into a refusal a scheduler can act on.
//
// The probe is `promise --version`: it resolves the binary, runs it, and needs
// no project. Presence on PATH alone would pass for a broken install, and
// anything heavier would charge every fit check for work only a build needs.
func runFit(repoRoot string) error {
	if Which("promise") == "" {
		return fmt.Errorf("the promise toolchain is not on PATH — every gate here spawns `promise`, " +
			"so this machine can measure nothing (install it from https://promise-lang.org)")
	}
	out, err := RunCaptureIn(repoRoot, "promise", "--version")
	if err != nil {
		return fmt.Errorf("`promise --version` did not answer, so the toolchain on this machine is not usable: %v (%s)", err, out)
	}
	return nil
}

// gateRunner measures one concept in one module.
type gateRunner func(repoRoot, module string) error

func gateRunners() map[GateConcept]gateRunner {
	return map[GateConcept]gateRunner{
		GateTested: runTested,
	}
}

// runTested runs one module's tests.
//
// `promise test .` discovers the module's test functions itself. A module with
// no test files is reported as passing by the toolchain, and that reading is
// kept: "this module has no tests" is not a measurement this gate is entitled
// to fail, and the absence is visible in the tree.
func runTested(repoRoot, module string) error {
	dir := filepath.Join(repoRoot, filepath.FromSlash(module))
	out, err := RunCaptureIn(dir, "promise", "test", ".")
	if err != nil {
		return fmt.Errorf("promise test in %s:\n%s", module, out)
	}
	return nil
}

// RunMeasurement runs every non-host gate across all modules.
//
// One concept today, and the loop stays because the shape is the contract: a
// second measurement is added to measurementOrder in cheapest-first position
// and every caller picks it up.
//
// Two entry points use it — bin/verify and the integration gate — and neither
// shells out to the other, so what verify blesses and what integration measures
// cannot disagree.
func RunMeasurement(repoRoot string) error {
	for _, c := range measurementOrder {
		if err := RunGate(repoRoot, string(c)); err != nil {
			return err
		}
	}
	return nil
}

// measurementOrder is the cheapest-first ladder, stated once so RunMeasurement
// and providedConcepts cannot drift apart.
var measurementOrder = []GateConcept{GateTested}

// runIntegration is the composition a landing decision rests on. It is exactly
// RunMeasurement: `fit` is excluded by design, and base provides no gate that
// is neither a measurement nor the host's.
func runIntegration(repoRoot string) error { return RunMeasurement(repoRoot) }

// splitGateName divides `concept:instance`. A name with no colon is all concept.
func splitGateName(name string) (GateConcept, string) {
	if i := strings.Index(name, ":"); i >= 0 {
		return GateConcept(name[:i]), name[i+1:]
	}
	return GateConcept(name), ""
}

// modulesFor resolves an instance to the modules it names.
func modulesFor(repoRoot, instance string) ([]string, error) {
	mods, err := PromiseModules(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("discovering Promise modules: %w", err)
	}
	if len(mods) == 0 {
		return nil, fmt.Errorf("this repository holds no Promise module (no promise.toml anywhere), so there is nothing to measure")
	}
	if instance == "" {
		return mods, nil
	}
	for _, m := range mods {
		if ModuleLabel(m) == instance {
			return []string{m}, nil
		}
	}
	var known []string
	for _, m := range mods {
		known = append(known, ModuleLabel(m))
	}
	return nil, fmt.Errorf("unknown instance %q (this project has: %s)", instance, strings.Join(known, ", "))
}

// providedConcepts lists what this project answers, for an error message that
// tells the caller what it could have asked for.
func providedConcepts() []string {
	out := []string{string(GateIntegration), string(GateFit)}
	for _, c := range measurementOrder {
		out = append(out, string(c))
	}
	sort.Strings(out)
	return out
}

// MeasureGate runs one gate and returns the envelope describing what it found,
// alongside the gate's own error.
//
// The envelope IS the seam between measuring and judging, so it is built here
// once rather than at each caller. `bin/gate <name> --envelope` prints it for a
// runner that will carry it elsewhere; `bin/run <gate>` hands it straight to
// the judge in the same process. Two constructions of the same document could
// disagree about what was measured, and the disagreement would be invisible —
// both callers would still emit something well-formed, and only the verdicts
// would differ.
//
// Stdout is redirected to stderr for the duration. The caller's stdout carries
// its product and nothing else — an envelope a runner parses, or a report a
// person reads — while a long gate's progress still streams where someone
// watching can see it.
func MeasureGate(repoRoot, name string) (map[string]any, error) {
	realStdout := os.Stdout
	os.Stdout = os.Stderr
	started := time.Now()
	gerr := RunGate(repoRoot, name)
	os.Stdout = realStdout

	env := map[string]any{
		"gate":            name,
		"measured":        gerr == nil,
		"elapsed_seconds": time.Since(started).Round(time.Millisecond).Seconds(),
		"measurements":    gateMeasurements(gerr),
	}
	if gerr != nil {
		env["detail"] = gerr.Error()
	}
	return env, gerr
}

// MetricFailedGates is the metric every envelope here carries, and the one the
// judge's threshold names.
const MetricFailedGates = "failed_gates"

// gateMeasurements is what one run measured, and the one declaration of which
// metrics an envelope from this project carries.
//
// failed_gates is the one metric measured here rather than by a gate: it is
// about the named gate as a whole, so it is 0 or 1 — `integration` is one gate,
// not two — and it is the number the judge's cap applies to.
func gateMeasurements(gerr error) map[string]int64 {
	failed := int64(0)
	if gerr != nil {
		failed = 1
	}
	return map[string]int64{MetricFailedGates: failed}
}

// GateMetrics names what an envelope from this project carries — the keys
// gateMeasurements puts under `measurements`.
//
// It exists so a gate asked what it measures answers out of the same
// declaration the envelope is built from, rather than out of a second list
// that goes stale the first time a metric is added
// (docs/gate-contract.md, "The exec line": a bare invocation states the gate's
// name, what it measures, and the command that runs it). The names are read out
// of that construction rather than restated beside it, so a metric added there
// is named here without anyone remembering to — a second literal is a list that
// goes quietly short, and what it costs is a bare invocation telling a caller
// the gate measures less than it does. Which run it asks about does not matter,
// since the keys are the same whatever the gate found, so it asks about one
// that found nothing wrong.
//
// Sorted, because a map is not: a list whose order changed between two bare
// invocations would read as a gate that measures something different each time.
//
// It takes no gate name because every gate here answers the same one metric.
// The contract declares metrics per gate, and the day base does too this takes
// the name — narrowing by a name nothing narrows on today would be a parameter
// no caller could pass wrongly and no reader could check.
func GateMetrics() []string {
	measured := gateMeasurements(nil)
	names := make([]string, 0, len(measured))
	for name := range measured {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GateNames returns every gate name this project answers, concepts and
// instances — what `bin/gate --list` and `bin/run --list` both answer with.
//
// Discovery failing is not "no gates": the concepts are answered whatever the
// tree looks like, and only the per-module instances are lost. Reporting an
// empty list would read to a scheduler as a machine that can run nothing.
func GateNames(repoRoot string) []string {
	var out []string
	mods, err := PromiseModules(repoRoot)
	for _, c := range providedConcepts() {
		out = append(out, c)
		// Neither narrows by module: integration is the composition, and fit
		// measures the machine — a module is not a property of the host.
		if GateConcept(c) == GateIntegration || GateConcept(c) == GateFit {
			continue
		}
		if err != nil {
			continue
		}
		for _, m := range mods {
			out = append(out, c+":"+ModuleLabel(m))
		}
	}
	return out
}
