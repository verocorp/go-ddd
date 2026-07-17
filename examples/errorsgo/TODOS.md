# examples/errorsgo — TODO (Go mirror of examples/errorspy)

PAUSED 2026-07-17 after cell 1. The Python proof (`examples/errorspy`, merged to
main in PR #4) is the source of truth for every decision — this is a translation,
not a redesign. Keep cell + test names identical to errorspy so drift is visible.

Branch: `errorsgo-mirror` (off merged main). Design doc:
`~/.gstack/projects/verocorp-go-ddd/error-norms-design-2026-07-17.md`.

## Done
- [x] `errdom/` error model — Kind/Code, DomainError, InfraError,
      Invalid/NotFound/Conflict, Wrap, Collect/FieldProblem, StatusFor.
      Green (go vet + go test).

## Remaining (port from examples/errorspy)
- [ ] `domain/values.go` — Slug, TargetURL leaf VOs (C1) + `MustNew*` (C4);
      DateWindow compound VO (C2). + tests.
- [ ] `domain/shortlink.go` — ShortLink entity (C3: PROPAGATE child errors, no
      `%w` re-wrap); Deactivate illegal transition (D1). + tests.
- [ ] `domain/campaign.go` — Campaign aggregate: collection invariants
      (D2: duplicate_slug, too_many_links), index-`Wrap` child error keeping
      kind+code (X1), missing link -> not_found. + tests.
- [ ] `app/storage.go` — fake vendor storage + StorageMiss/StorageUnavailable.
- [ ] `app/repository.go` — translating adapter: miss->not_found (B1/B2),
      outage->InfraError (B2), corrupt record->InfraError not validation (B7). + tests.
- [ ] `app/service.go` — orchestration, no wrap (X2). + tests.
- [ ] `transport/handler.go` — pure StatusFor mapper (B3), RFC 9457 body + field
      (B4) + invalid-params (B6), malformed->400 (B5), InfraError->503, else->500. + test.
- [ ] `main.go` + `e2e_test.go` — every status reachable through the real handler.
- [ ] X4 — context cancellation / timeout -> InfraError (Go-idiomatic addition).

## OPEN DECISION (gates only the transport checkability story)
Go `Kind->status` exhaustiveness: the `exhaustive` analyzer's released build lags
go1.25 (x/tools skew), so it will not run against this repo yet. Options:
(a) runtime-witness `go test` over all Kinds now + wire `exhaustive` later
    [recommended], (b) wire `exhaustive` now (resolve the lag), (c) runtime-only
permanently. `StatusFor` already has a panic backstop. Decide at the transport layer.

## After the mirror (later rings, not started)
Adopt norms into the real twins (examples/running + examples/python, fixing their
bugs); teach in skills (go.md/python.md) + reconcile the 3 llm-tools copies;
ddd-vet(-py) enforcement checks; wire pilot.
