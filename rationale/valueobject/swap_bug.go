//go:build changeability_bug

package valueobject

// This file is EXCLUDED from normal builds (the changeability_bug tag is never
// set in production or in `go test ./...`). It exists only so a test can prove
// that the unit mismatch the primitive version admitted silently is a hard
// compile error here.
//
// Build it with: go build -tags changeability_bug ./valueobject  -> MUST fail.
//
// This is the Mars Climate Orbiter bug in miniature: a Feet value passed where
// Meters is expected.
var _ = TimeToImpact(NewFeet(10000), NewMetersPerSecond(50))
