package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// CommandNames is what this project builds into bin/: every directory under
// tools/build/cmd except the meta-builder itself.
//
// Read from the directory rather than from a list, for the reason `make`
// discovers the same set: a tool added to cmd/ and forgotten in a list here
// would be a binary that `./make` builds and `bin/run --list` denies, and the
// two disagreeing is exactly what the workspace's one-name-one-builder check
// reads to decide whether a name is claimed twice.
//
// `make` is excluded because it is never compiled into bin/. It runs via `go
// run` from the ./make trampoline, which is what breaks the bootstrap cycle —
// so it is not a binary this project builds, and reporting it would have the
// workspace look for a bin/make that must never exist.
func CommandNames(repoRoot string) ([]string, error) {
	cmdDir := filepath.Join(repoRoot, "tools", "build", "cmd")
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", cmdDir, err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "make" {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}
