// Package emit is the COUPLED domain variant for decision 3: a maneuver domain
// object that EMITS its own outward representation. Its Record() method returns a
// telemetry.Record, so this "domain" package imports the outward telemetry layer —
// the dependency-direction violation decision 3 forbids. (It can exist as its own
// package precisely because it is downstream of both nav and telemetry; the pure
// nav package cannot take this import without the D3a cycle.)
//
// A dependent that reaches through emit.Maneuver.Record() to a telemetry-specific
// field is bound to the outward format and breaks when it migrates — the coupled
// arm of the D3b spine. See ../../SCORING.md.
package emit

import (
	"github.com/verocorp/go-ddd/rationale/changeability/nooutward/nav"
	"github.com/verocorp/go-ddd/rationale/changeability/nooutward/telemetry"
)

// Maneuver wraps the domain value but leaks its outward representation through
// Record() — the pattern decision 3 says to keep in the outward layer instead.
type Maneuver struct{ inner nav.Maneuver }

func NewManeuver(id string, thrustMicroN, burnMillis int64) Maneuver {
	return Maneuver{inner: nav.NewManeuver(id, thrustMicroN, burnMillis)}
}

func (m Maneuver) ID() string { return m.inner.ID() }

// Record emits the domain's outward telemetry representation directly off the
// domain object — the leak. Note Record() itself names no reshaped field, so this
// package survives the migration; only the CONSUMERS that reach into the returned
// record's fields are forced to change.
func (m Maneuver) Record() telemetry.Record { return telemetry.FromManeuver(m.inner) }
