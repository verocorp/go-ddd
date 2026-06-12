package primitive

// Distances and rates here are bare float64 with no unit attached. A function
// that wants meters cannot stop a caller from passing a value measured in feet
// — the Mars Climate Orbiter failure (NASA, 1999): a vendor supplied impulse
// data in imperial units, the navigation software assumed metric, and the
// $327M spacecraft burned up in the Martian atmosphere. The compiler saw two
// float64s and said nothing.

// TimeToImpact returns seconds to ground contact at a constant descent rate.
// altitude is meant to be METERS, descentRate METERS PER SECOND — but both are
// bare float64, so passing feet, or swapping the two arguments, compiles clean
// and returns a confidently wrong number.
func TimeToImpact(altitude, descentRate float64) float64 {
	return altitude / descentRate
}
