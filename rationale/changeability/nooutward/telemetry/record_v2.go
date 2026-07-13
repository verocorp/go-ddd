//go:build repv2

// Telemetry v2 — the POST-migration outward format (build with -tags repv2).
// BurnSeconds is gone, reshaped into DurationMillis (finer unit, different name):
// a realistic wire-format change. The domain (nav) is untouched by this — it never
// knew the outward shape — so a dependent holding a nav.Maneuver is untouched too.
// A dependent that reached through an emitting domain to name BurnSeconds fails to
// compile here. That contrast is the D3b spine. See ../../SCORING.md.
package telemetry

import "github.com/verocorp/go-ddd/rationale/changeability/nooutward/nav"

// Record is the v2 telemetry shape. Deliberately reshaped: BurnSeconds removed,
// DurationMillis added.
type Record struct {
	ManeuverID     string
	DurationMillis int64 // v2: milliseconds — replaces v1's BurnSeconds
	ThrustMicroN   int64
}

func FromManeuver(m nav.Maneuver) Record {
	return Record{
		ManeuverID:     m.ID(),
		DurationMillis: m.BurnMillis(),
		ThrustMicroN:   m.ThrustMicroN(),
	}
}
