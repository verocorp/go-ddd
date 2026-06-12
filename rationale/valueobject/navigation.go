package valueobject

// Meters and Feet are distinct length types; the compiler will not let one be
// used where the other is expected. Crossing unit systems is possible only
// through an explicit, named conversion — the exact guardrail the Mars Climate
// Orbiter navigation code lacked.
type Meters struct{ v float64 }

func NewMeters(v float64) Meters { return Meters{v: v} }
func (m Meters) Float() float64  { return m.v }

type Feet struct{ v float64 }

func NewFeet(v float64) Feet { return Feet{v: v} }
func (f Feet) Float() float64 { return f.v }

// ToMeters is the only sanctioned way to cross from imperial to metric.
func (f Feet) ToMeters() Meters { return Meters{v: f.v * 0.3048} }

// MetersPerSecond is a distinct rate type — not interchangeable with a Meters
// distance, so altitude and descent rate cannot be swapped.
type MetersPerSecond struct{ v float64 }

func NewMetersPerSecond(v float64) MetersPerSecond { return MetersPerSecond{v: v} }
func (r MetersPerSecond) Float() float64           { return r.v }

// TimeToImpact takes typed quantities. Passing Feet where Meters is expected,
// or swapping altitude and descent rate, does not compile. See swap_bug.go.
func TimeToImpact(altitude Meters, descentRate MetersPerSecond) float64 {
	return altitude.v / descentRate.v
}
