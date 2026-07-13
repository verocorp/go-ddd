// Package nav is the DECOUPLED domain for the no-outward-representation decision
// (decision 3). Maneuver is a pure domain value object: private fields, a single
// construction path, domain accessors — and it imports NOTHING outward. It emits
// no transport/persistence representation of itself; mapping the domain to an
// outward shape lives in the outward layer (package telemetry), which depends
// inward on this package. A dependent that holds a nav.Maneuver survives an
// outward-format migration untouched — the decoupled arm of decision 3.
//
// See ../../SCORING.md (D3a direction guard + D3b migration spine).
package nav

// Maneuver is the domain object: a single spacecraft burn. Its representation
// never leaks — no field is exported, no accessor returns a transport type, and
// the package has no outward import. (The neutral domain nods at the Mars Climate
// Orbiter: a burn whose units, once leaked to an outward representation, are the
// whole cautionary tale.)
type Maneuver struct {
	id           string
	thrustMicroN int64 // micro-newtons (domain unit, never a wire unit)
	burnMillis   int64
}

// NewManeuver is the single construction path.
func NewManeuver(id string, thrustMicroN, burnMillis int64) Maneuver {
	return Maneuver{id: id, thrustMicroN: thrustMicroN, burnMillis: burnMillis}
}

func (m Maneuver) ID() string          { return m.id }
func (m Maneuver) ThrustMicroN() int64 { return m.thrustMicroN }
func (m Maneuver) BurnMillis() int64   { return m.burnMillis }
