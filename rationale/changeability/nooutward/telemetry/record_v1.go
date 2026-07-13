//go:build !repv2

// Package telemetry is the OUTWARD layer for decision 3: the downstream wire /
// telemetry representation of a maneuver, plus the mapper from the domain to it.
// It imports the domain (nav) to map it — the CORRECT inward dependency
// direction. Because telemetry depends on nav, the reverse (nav importing
// telemetry) is an import cycle that will not compile: that is the D3a direction
// guard (see nav/leak_bug.go and ../../SCORING.md).
//
// This file is the v1 (pre-migration) format. Building with -tags repv2 replaces
// it with record_v2.go, modelling an outward-representation migration inside one
// committed tree (the D3b spine).
package telemetry

import "github.com/verocorp/go-ddd/rationale/changeability/nooutward/nav"

// Record is the v1 telemetry shape. BurnSeconds is the field the migration
// reshapes away; a dependent that reached through an emitting domain to name it
// (the coupled arm) will not compile after the migration.
type Record struct {
	ManeuverID   string
	BurnSeconds  int64 // v1: whole seconds — reshaped to DurationMillis in v2
	ThrustMicroN int64
}

// FromManeuver maps the domain object to its outward representation. This is the
// sanctioned place for domain→DTO mapping: in the outward layer, depending inward
// on the domain — NOT a method on the domain object.
func FromManeuver(m nav.Maneuver) Record {
	return Record{
		ManeuverID:   m.ID(),
		BurnSeconds:  m.BurnMillis() / 1000,
		ThrustMicroN: m.ThrustMicroN(),
	}
}
