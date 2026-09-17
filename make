#!/usr/bin/env bash
# Bootstrap trampoline. Compiles every dev tool into bin/ via the meta-builder.
#
# This file is committed shell text — nothing builds it. It runs the
# meta-builder via 'go run', which needs only the Go toolchain (no pre-built
# binary), and that meta-builder compiles every other tool into bin/.
#
# The inner cd resolves $0 via its own directory, so ./make works from any cwd,
# and pins go run's working directory to <repo>/tools/build — which is how the
# meta-builder learns the absolute repo root.
set -euo pipefail
repo="$(cd "$(dirname "$0")" && pwd)"

go run -C "$repo/tools/build" ./cmd/make "$@"
