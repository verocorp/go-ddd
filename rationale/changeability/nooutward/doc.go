// Package nooutward is the changeability arm for decision 3: a domain object
// emits no non-domain representation, and the domain depends on no outward layer.
//
// Layout:
//   - nav/       — the pure domain (imports nothing outward); leak_bug.go is the
//     D3a direction guard (build -tags leak fails with an import cycle).
//   - telemetry/ — the outward layer; imports nav to map it (correct inward
//     direction); record_v1/record_v2 are the -tags repv2 migration (D3b spine).
//   - emit/      — the coupled domain variant: a "domain" that emits its own
//     telemetry.Record (the dependency-direction violation).
//   - decoupled/consumerNN, coupled/fanout/consumerNN — the generated N-scaling
//     arms the contrast measures.
//
// The proof lives in contrast_test.go (D3b spine) and direction_test.go (D3a
// guard). Scoring is predeclared in ../SCORING.md.
package nooutward

//go:generate go run ./internal/gen
