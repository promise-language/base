// Command gate runs one named gate against this repository.
//
// `bin/gate <name>` is the fixed entry point a flow uses to ask for a
// measurement. Fixed, and deliberately not configurable: the protocol addresses
// gates by name, so a project spelling its entry point differently would have
// gates nothing could ask for.
//
// A gate measures and modifies nothing — including afterwards. That is what
// separates it from `bin/verify`, which is a command: verify repairs what has
// one right answer on its way to an answer, which is what a producing step
// wants and exactly why a landing decision may not rest on it.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/promise-language/base/tools/build/common"
)

// Injected by the meta-builder via -ldflags at build time; empty otherwise.
var (
	repoRoot   = ""
	sourceHash = ""
)

// gatelisting is the machine-readable form of the discovery query. The flow
// SDK reads `gates[].name`, so the object is shaped around that and nothing
// else.
type gateListing struct {
	Gates []gateEntry `json:"gates"`
}

type gateEntry struct {
	Name string `json:"name"`
}

func main() {
	common.CheckStale(repoRoot, sourceHash)

	// --envelope and --json are protocol, not configuration: a runner asks for
	// machine-readable output with them, and every other invocation is a person
	// at a terminal. They are stripped here so the name reaches the same place
	// it would have without them.
	var (
		args     []string
		envelope bool
		asJSON   bool
		list     bool
	)
	for _, a := range common.NormalizeArgs(os.Args[1:]) {
		switch a {
		case "-envelope":
			envelope = true
		case "-json":
			asJSON = true
		case "-list":
			list = true
		default:
			args = append(args, a)
		}
	}

	// -list is answered before the single-name rule below: the query takes no
	// gate name, and refusing it reads as a machine with no gates — which is
	// startup refusing every arena of this project.
	if list {
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, "gate: --list takes no gate name, got %q\n", strings.Join(args, " "))
			os.Exit(2)
		}
		writeGateList(asJSON)
		return
	}

	if len(args) != 1 || args[0] == "-h" || args[0] == "-help" {
		usage()
		// No argument is a usage error, not a passing gate: exiting 0 here would
		// let a caller that forgot the name read silence as success.
		os.Exit(2)
	}

	// A bare invocation is refused, and does not measure. Any call without the
	// flag is a person or an agent at a terminal, and this program is not a
	// channel for them: a gate has no verdict to give, and the first script to
	// wrap `bin/gate tested` would read its 0 as a pass. Refusing before
	// measuring keeps the two readings from ever coexisting, and costs a person
	// nothing — `bin/run <name>` is their path, and it is the one that judges.
	//
	// What it says is the three facts the contract asks for
	// (docs/gate-contract.md, "The exec line"): the gate's name, what it
	// measures, and the command that runs it. Naming the metrics is what makes
	// this worth reading rather than a scolding — someone who typed the wrong
	// thing learns what this gate would have told them, and whether it is the
	// one they wanted.
	if !envelope {
		fmt.Fprintf(os.Stderr, "%s — measures %s\n\n  bin/run %s\n",
			args[0], strings.Join(common.GateMetrics(), ", "), args[0])
		os.Exit(1)
	}

	// Envelope mode. common.MeasureGate builds the document and keeps the
	// gate's own progress off stdout, so what follows carries the envelope and
	// nothing else. It is shared with `bin/run <gate>`, which judges the same
	// envelope without a runner in between.
	env, gerr := common.MeasureGate(repoRoot, args[0])
	if err := json.NewEncoder(os.Stdout).Encode(env); err != nil {
		// The envelope is the whole point of this mode: a run that cannot state
		// what it measured has not measured anything a caller may act on.
		fmt.Fprintln(os.Stderr, "gate: could not write the envelope:", err)
		os.Exit(2)
	}
	// A gate that measured a failure and said so exits non-zero and HAS
	// measured — the runner records this code and decides nothing with it.
	if gerr != nil {
		os.Exit(1)
	}
}

// writeGateList answers the discovery query in whichever form was asked for.
func writeGateList(asJSON bool) {
	names := common.GateNames(repoRoot)
	if !asJSON {
		fmt.Println(strings.Join(names, "\n"))
		return
	}
	listing := gateListing{Gates: make([]gateEntry, 0, len(names))}
	for _, n := range names {
		listing.Gates = append(listing.Gates, gateEntry{Name: n})
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(listing); err != nil {
		fmt.Fprintln(os.Stderr, "gate: could not write the gate list:", err)
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: bin/gate <name> --envelope   measure and print the envelope\n"+
		"       bin/gate --list [--json]     the gates this project answers\n\nGates:\n  %s\n\n",
		strings.Join(common.GateNames(repoRoot), "\n  "))
	fmt.Fprint(os.Stderr, "A name is a concept, optionally with an instance: `tested:wire`.\n"+
		"Omitting the instance measures every Promise module.\n\n"+
		"`integration` is the composition a landing decision rests on.\n"+
		"For the repairing counterpart, see bin/verify — it is a command, not a gate.\n")
}
