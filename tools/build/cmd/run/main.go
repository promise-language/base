// Command run is this project's discovery entry point and its judge.
//
// Two questions, one binary, because the flow SDK asks both of the same fixed
// name. `bin/run --list` answers what this project builds and what it answers;
// `bin/run <gate> --verdict` decides whether a measurement is acceptable.
//
// The judging split matters. `bin/gate` MEASURES and states what it found; this
// decides whether what it found is good enough. The SDK runs both and reads
// neither — it never holds a project's numbers, and it never computes a
// verdict. That is what lets the terms live here, in the tree, where they are
// reviewed with the code they judge.
//
// The envelope arrives on stdin, exactly as `bin/gate --envelope` printed it.
//
// # The by-hand mode
//
// `bin/run <gate>`, with no --verdict, measures the gate itself and prints the
// verdict for a person: each measurement beside the term it was judged on. It
// is the path someone iterating on one failing area wants.
//
// It does not weaken the split. The judging terms are reached by both modes and
// applied by one function; the by-hand mode only removes the runner from the
// middle, which is safe precisely because nothing crosses a process boundary
// that a runner would have had to carry.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/promise-language/base/tools/build/common"
)

// Injected by the meta-builder via -ldflags at build time; empty otherwise.
var (
	repoRoot   = ""
	sourceHash = ""
)

// buildList is what --list answers. The two kinds stay separable, and must not
// be flattened: a two-builders check is about binaries and a gate is not one,
// so a merged list would compare a tool name against a gate name and report a
// collision that cannot exist.
type buildList struct {
	Commands []string `json:"commands"`
	Gates    []string `json:"gates"`
}

// verdict is what a caller must be handed. `acceptable` is the answer;
// `thresholds` are the terms it was judged against, carried so the decision can
// be recomputed by someone who was not there; `detail` is the reason a person
// needs.
type verdict struct {
	Acceptable bool                        `json:"acceptable"`
	Thresholds map[string]common.Threshold `json:"thresholds"`
	Detail     string                      `json:"detail"`
}

func main() {
	common.CheckStale(repoRoot, sourceHash)

	var (
		args        []string
		wantVerdict bool
		list        bool
		asJSON      bool
	)
	for _, a := range common.NormalizeArgs(os.Args[1:]) {
		switch a {
		case "-verdict":
			wantVerdict = true
		case "-list":
			list = true
		case "-json":
			asJSON = true
		default:
			args = append(args, a)
		}
	}

	if list {
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, "run: --list takes no gate name, got %q\n", strings.Join(args, " "))
			os.Exit(2)
		}
		writeList(asJSON)
		return
	}

	if len(args) != 1 || args[0] == "-h" || args[0] == "-help" {
		usage()
		// A caller that asked for nothing gets a usage error, not a verdict:
		// exiting 0 here would let silence read as "acceptable".
		os.Exit(2)
	}
	gate := args[0]

	if !wantVerdict {
		byHand(gate)
		return
	}

	var envelope map[string]any
	if err := json.NewDecoder(os.Stdin).Decode(&envelope); err != nil {
		fmt.Fprintf(os.Stderr, "run: the envelope for gate %q could not be read: %v\n", gate, err)
		os.Exit(2)
	}
	if envelope == nil {
		fmt.Fprintf(os.Stderr, "run: the envelope for gate %q is null\n", gate)
		os.Exit(2)
	}

	v, err := judge(gate, envelope)
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(2)
	}
	if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, "run: could not write the verdict:", err)
		os.Exit(2)
	}
}

// collectList answers both kinds, and refuses a name that is both.
//
// One namespace, because a collision would make one of the two silently
// unreachable while both still appeared in the list. It is an error rather than
// a precedence rule for that reason.
func collectList() (buildList, error) {
	commands, err := common.CommandNames(repoRoot)
	if err != nil {
		return buildList{}, err
	}
	gates := common.GateNames(repoRoot)
	var both []string
	for _, c := range commands {
		if slices.Contains(gates, c) {
			both = append(both, c)
		}
	}
	if len(both) > 0 {
		return buildList{}, fmt.Errorf(
			"%s is both a command this project builds and a gate it answers, so `run %s` means one of two things — rename one of them",
			strings.Join(both, ", "), both[0])
	}
	// Never nil: `{"commands": []}` at exit 0 is a project stating it builds
	// none, which is DEFINITIVE, and JSON null would read as the unknown that
	// only an absent binary is.
	if commands == nil {
		commands = []string{}
	}
	if gates == nil {
		gates = []string{}
	}
	return buildList{Commands: commands, Gates: gates}, nil
}

func writeList(asJSON bool) {
	list, err := collectList()
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(2)
	}
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(list); err != nil {
			fmt.Fprintln(os.Stderr, "run: could not write the list:", err)
			os.Exit(2)
		}
		return
	}
	fmt.Println("commands:")
	for _, c := range list.Commands {
		fmt.Println("  " + c)
	}
	fmt.Println("gates:")
	for _, g := range list.Gates {
		fmt.Println("  " + g)
	}
}

// judge applies this project's terms to one envelope.
//
// The terms are read from tools/gates/thresholds.json — an artefact distinct
// from this judge, so relaxing them is a change to a tracked file whose whole
// content is the terms, rather than a line of code in the thing that applies
// them.
//
// One function reached by both modes, so a by-hand answer and the answer the
// SDK records cannot differ. Two callers asking about the same measurement and
// getting different verdicts is the exact failure the fixed entry point exists
// to prevent, and it would be no less a failure for happening inside one binary.
func judge(gate string, envelope map[string]any) (verdict, error) {
	measured, ok := envelope["measured"].(bool)
	if !ok {
		return verdict{}, fmt.Errorf(
			"the envelope for gate %q does not state `measured`, so there is nothing to judge", gate)
	}

	terms, err := common.LoadThresholds(repoRoot)
	if err != nil {
		return verdict{}, err
	}

	// The measurements the envelope carries. Absent is not zero: a term whose
	// metric nobody measured is skipped by ApplyThresholds rather than applied
	// to an invented number.
	measurements := map[string]int64{}
	if raw, isMap := envelope["measurements"].(map[string]any); isMap {
		for k, v := range raw {
			if f, isNum := v.(float64); isNum {
				measurements[k] = int64(f)
			}
		}
	}
	// failed_gates is what this envelope is fundamentally about, and an
	// envelope that omitted it would leave the judge with no metric at all.
	// `measured` states the same fact, so it is the fallback rather than a
	// second source that could disagree.
	if _, present := measurements[common.MetricFailedGates]; !present {
		failed := int64(0)
		if !measured {
			failed = 1
		}
		measurements[common.MetricFailedGates] = failed
	}

	applied, breaches := common.ApplyThresholds(terms, measurements)
	if len(applied) == 0 {
		return verdict{}, fmt.Errorf(
			"no term in %s names anything the envelope for gate %q measured, so the verdict would rest on nothing",
			common.ThresholdsPath, gate)
	}

	v := verdict{Acceptable: len(breaches) == 0, Thresholds: applied}
	if v.Acceptable {
		v.Detail = fmt.Sprintf("the %s gate reported no failure, and every term it was held to was met", gate)
		return v, nil
	}
	v.Detail = strings.Join(breaches, "; ")
	if d, isStr := envelope["detail"].(string); isStr && d != "" {
		v.Detail += ": " + d
	}
	return v, nil
}

// byHand measures the gate and reports the verdict to a person.
//
// It exits 1 on an unacceptable measurement so the mode is usable in a loop —
// `bin/run tested && ...` — and 2 only when the question could not be answered
// at all. That is the same distinction bin/gate draws: a gate that measured a
// failure HAS measured, and is not the same as one that could not run.
func byHand(gate string) {
	env, _ := common.MeasureGate(repoRoot, gate)
	v, err := judge(gate, env)
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(2)
	}

	// The measurement beside the term it was judged on. A verdict printed alone
	// tells someone iterating that they failed, not what they were held to.
	fmt.Printf("%s\n\n", gate)
	fmt.Printf("  measured  %s\n", v.Detail)
	for name, t := range v.Thresholds {
		fmt.Printf("  term      %s\n", t.Describe(name))
	}
	if elapsed, ok := env["elapsed_seconds"].(float64); ok {
		fmt.Printf("  elapsed   %.1fs\n", elapsed)
	}
	if v.Acceptable {
		fmt.Printf("  verdict   acceptable\n")
		return
	}
	fmt.Printf("  verdict   NOT acceptable\n")
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: bin/run <gate>             measure the gate and judge it\n"+
		"       bin/run <gate> --verdict   judge an envelope on stdin\n"+
		"       bin/run --list [--json]    what this project builds and answers\n\n"+
		"Gates this project answers:\n  %s\n\n",
		strings.Join(common.GateNames(repoRoot), "\n  "))
	fmt.Fprint(os.Stderr, "A name is a concept, optionally with an instance: `tested:wire`.\n"+
		"Omitting the instance measures every Promise module.\n\n"+
		"With --verdict it decides only, reading the envelope `bin/gate <gate>\n"+
		"--envelope` printed — the mode the SDK uses. Without it, this measures\n"+
		"the gate itself and prints each measurement beside the term it was\n"+
		"judged on.\n\n"+
		"Passing one gate is not passing the whole: only bin/verify confirms the\n"+
		"full set.\n")
}
