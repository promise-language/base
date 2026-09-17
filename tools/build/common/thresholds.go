package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// The judge's terms, as an artefact distinct from the judge.
//
// The workspace's docs/tool-contract.md §3 places them at
// tools/gates/thresholds.json, and the separation is the point: a judge that
// carried its own numbers inline could be relaxed in the same change that broke
// what they measure, and the relaxation would read as ordinary code. Kept here,
// a change to the terms is a change to a tracked file whose whole content is
// the terms — visible in a diff, reviewable on its own.

// ThresholdsPath is where the judge's terms live, relative to the repo root.
// Spelled once; the doctor that checks the file is tracked spells it too, which
// is the contract rather than duplicated logic.
const ThresholdsPath = "tools/gates/thresholds.json"

// Threshold is one term: a direction and the bound it names.
type Threshold struct {
	Direction string `json:"direction"`
	Cap       int64  `json:"cap"`
}

// Directions a term may state. Anything else is a malformed artefact and is
// refused rather than guessed at — a term nobody can read is not a term.
const (
	DirectionAtMost  = "at_most"
	DirectionAtLeast = "at_least"
)

// LoadThresholds reads the terms.
//
// A missing or unreadable file is an ERROR, never an empty set of terms.
// Reading "no file" as "nothing to hold this to" would make a deleted artefact
// the most permissive state there is, and every measurement would pass.
func LoadThresholds(repoRoot string) (map[string]Threshold, error) {
	path := filepath.Join(repoRoot, filepath.FromSlash(ThresholdsPath))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("the judge's terms at %s could not be read: %w", ThresholdsPath, err)
	}
	var terms map[string]Threshold
	if err := json.Unmarshal(data, &terms); err != nil {
		return nil, fmt.Errorf("the judge's terms at %s are not readable JSON: %w", ThresholdsPath, err)
	}
	if len(terms) == 0 {
		return nil, fmt.Errorf("the judge's terms at %s are empty, so nothing would be held to anything", ThresholdsPath)
	}
	for name, t := range terms {
		if t.Direction != DirectionAtMost && t.Direction != DirectionAtLeast {
			return nil, fmt.Errorf("the term %q at %s states direction %q, which is neither %q nor %q",
				name, ThresholdsPath, t.Direction, DirectionAtMost, DirectionAtLeast)
		}
	}
	return terms, nil
}

// Satisfied reports whether one measurement meets one term.
func (t Threshold) Satisfied(measured int64) bool {
	if t.Direction == DirectionAtLeast {
		return measured >= t.Cap
	}
	return measured <= t.Cap
}

// Describe states the term the way a person reads it.
func (t Threshold) Describe(name string) string {
	word := "at most"
	if t.Direction == DirectionAtLeast {
		word = "at least"
	}
	return fmt.Sprintf("%s %s %d", name, word, t.Cap)
}

// ApplyThresholds judges measurements against terms and returns the terms that
// were actually applied, plus a sentence naming every breach.
//
// A term with no matching measurement is SKIPPED rather than treated as zero.
// Applying a term to a number nobody measured is the defect this whole path
// exists to remove: it reads as a measurement that held, when nothing was
// measured at all.
func ApplyThresholds(terms map[string]Threshold, measurements map[string]int64) (applied map[string]Threshold, breaches []string) {
	applied = map[string]Threshold{}
	var names []string
	for name := range terms {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		measured, ok := measurements[name]
		if !ok {
			continue
		}
		term := terms[name]
		applied[name] = term
		if !term.Satisfied(measured) {
			breaches = append(breaches, fmt.Sprintf("%s measured %d, and the term is %s", name, measured, term.Describe(name)))
		}
	}
	return applied, breaches
}
