//go:build leak

// This file is the D3a DIRECTION GUARD, compiled only under -tags leak. It tries
// to make the pure domain (nav) emit its own outward representation by importing
// the outward layer (telemetry) — the decision-3 violation. Because telemetry
// already imports nav (the correct inward direction), this reverse import is an
// IMPORT CYCLE and `go build -tags leak ./nav` FAILS to compile.
//
// That compile failure IS the assertion: the language enforces the dependency
// direction, so a domain object cannot emit a non-domain representation while the
// mapper lives in the outward layer. direction_test.go asserts this build fails —
// the decision-3 analog of the anchor's subst_bug.go. See ../../SCORING.md (D3a).
package nav

import "github.com/verocorp/go-ddd/rationale/changeability/nooutward/telemetry"

// Record would let the domain emit its own outward representation. It cannot
// exist: the import above is a cycle. Present only to make the illegal direction
// explicit and to give the guard something to fail on.
func (m Maneuver) Record() telemetry.Record { return telemetry.FromManeuver(m) }
