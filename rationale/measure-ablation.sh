#!/usr/bin/env bash
# measure-ablation.sh — mutate one source file, count the compiler-forced
# worklist that mutation produces, then revert. The point: a typed change is
# enumerated by the compiler (these errors ARE the worklist, none silent);
# the same change on a bare primitive produces zero of these.
#
# Usage:  measure-ablation.sh <file> <perl-expr> [build-target]
#
# Build the DEFINING package as the target (e.g. ./accounting/) — `go build`
# can't type-check packages downstream of one that fails to compile, so a
# whole-module target undercounts a rename/retype. The defining-package count
# is exact and untruncated (-gcflags=-e removes Go's 10-errors-per-package cap).
# For the full cross-package surface, pair this with `git grep -w <Symbol>`.
#
# Example (rename a VO type):
#   ./measure-ablation.sh accounting/amounts.go \
#     's/^type CreditType struct/type CreditKind struct/ if $.==11' ./accounting/
set -euo pipefail
file="$1"; expr="$2"; target="${3:-./...}"
cleanup() { git checkout -- "$file" 2>/dev/null || true; }
trap cleanup EXIT

perl -i -pe "$expr" "$file"
out=$(go build -gcflags=-e "$target" 2>&1 || true)
n=$(printf '%s\n' "$out" | grep -cE '\.go:[0-9]+:[0-9]+:')
files=$(printf '%s\n' "$out" | grep -oE '[A-Za-z0-9_./-]+\.go' | sort -u | wc -l | tr -d ' ')
echo "compiler-forced sites in $target : $n  (across $files files)"
printf '%s\n' "$out" | grep -E '\.go:[0-9]+:[0-9]+:' | head -6
