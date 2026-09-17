package common

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Promise modules are the subject of every gate here, and they are DISCOVERED
// rather than listed.
//
// A module is a directory holding a promise.toml, which is the same fact the
// compiler reads — so the set a gate measures and the set the toolchain builds
// cannot drift. The alternative, a slice of paths maintained by hand, has one
// failure mode this does not: a module added to the tree and forgotten here is
// a module no gate ever measures, and the tree still reports sound.
//
// That matters more in this repository than in most. base is scaffolding — its
// README says so — and `wire/` is the only contract written of the nine the
// README tables. The module set is expected to change with almost every change
// worth making.

// skipDirs are the directories a walk never descends into.
//
// bin/ is built per clone and gitignored; .workspace/ and .git/ are per-clone
// state; tools/ is this Go toolchain, which is not Promise and has its own
// build. None of them can hold a Promise module that a gate should measure, and
// descending into them costs a walk of the largest directories in the tree.
var skipDirs = map[string]bool{
	"bin": true, ".git": true, ".workspace": true, ".home": true,
	"tools": true, ".flow": true, "node_modules": true,
}

// PromiseModules returns every Promise module in the repository as a
// slash-separated path relative to repoRoot, sorted so a gate's output is the
// same on two machines.
//
// A module inside another module is not descended into: promise.toml names one
// module, and the nested directories under it are its own source.
func PromiseModules(repoRoot string) ([]string, error) {
	var found []string
	err := filepath.WalkDir(repoRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if path != repoRoot && (skipDirs[name] || strings.HasPrefix(name, ".")) {
			return filepath.SkipDir
		}
		if !Exists(filepath.Join(path, "promise.toml")) {
			return nil
		}
		rel, relErr := filepath.Rel(repoRoot, path)
		if relErr != nil {
			return relErr
		}
		found = append(found, filepath.ToSlash(rel))
		// A module's own subdirectories are its source, not further modules.
		return filepath.SkipDir
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(found)
	return found, nil
}

// ModuleLabel is the gate-instance name for a module directory: the path with
// separators flattened, so `tests/wire-consumer` is addressable as
// `tested:tests-wire-consumer`. A colon would collide with the concept
// separator and a slash is not in the SDK's grammar for an instance.
func ModuleLabel(module string) string {
	if module == "." {
		return "root"
	}
	return strings.ReplaceAll(module, "/", "-")
}
