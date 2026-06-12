# Coverage matrix — rule ↔ demonstrated win ↔ enforcement

This is the single source of truth linking each value-object rule to (a) the
executable demo in this `rationale/` package that shows the bug it prevents, and
(b) the go-ddd checker that enforces it (if one exists). `coverage_test.go`
fails on a **silent gap** — a checker with no row here, or a row naming a test
that doesn't exist. It tolerates the honest ❌/⚠️ rows by design: the rationale
makes the case for the whole discipline; go-ddd enforces the mechanically
checkable slice.

| Value-object rule | Demonstrated win (test in this package) | Real-world anchor | Enforced by | Status |
|---|---|---|---|---|
| Distinct types instead of bare primitives | `TestTypeConfusion_ValueObjectRejectsWrongUnit` — `Feet` where `Meters` expected won't compile | Mars Climate Orbiter (1999, $327M) | — | ❌ no checker yet |
| Value equality, not representation equality | `TestEquality_ValueObjectIsRight` — `0°C` == `273.15K` | scale/representation collision | `checkequality` (requires a `Test*_Equality`) | ✅ 1:1 |
| Constructors validate invariants | `TestValidation_ValueObjectRejectsBadInput` — sub-absolute-zero rejected | physically impossible state | — | ❌ no checker yet |
| VO constructors get `MustNew*` helpers | — (ergonomics, not a runtime-bug demo) | without it, agents/devs reinvent must-helpers in tests | `checkmustnew` | ⚠️ ergonomics; no safety demo by design |
| `.String()` is for display, not equality | — (demoable, not yet written) | `a.String() == b.String()` mis-equates multi-rep VOs | `checkstring` | ⚠️ demo TODO |

**Reading the gaps:** the two ❌ rows are wins go-ddd does not yet enforce
(candidate checkers for the `go/analysis` port — not scheduled). The ⚠️ rows are
checkers whose justification is ergonomic or whose demo isn't written. Nothing
here is a *silent* gap; every checker has a row and every named test exists.

## Run

```
go test ./rationale/...                 # the wins + the matrix meta-test
go test -bench=. -benchmem ./rationale/ # the adversarial cost (collection-VO defensive-copy tax)
./rationale/measure-ablation.sh ...     # measure changeability on your own repo
```
