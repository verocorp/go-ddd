# The field-audit checklist — friction day, classification, relay

Generic and consumer-agnostic: run it inside any consumer repo, behind
whatever IP wall applies. Nothing consumer-specific comes back — only
de-identified patterns in the relay format below. This file is deliberately
**delta-only**: it carries what the audit adds (classification prompts + the
relay format) and points at existing doctrine for everything else — the
anatomy walk is `skills/tesser-build/map.md`'s gap survey, and the deferral
test it applies is the repo's own changeability thesis
(`docs/design-three-contender-changeability.md`: silent sites vs enumerable
sites).

## 1. Run the day

- One prospective working day of real work in the consumer repo. Don't
  manufacture exercises — log frictions as they interrupt actual tasks.
- Seed the log up front with frictions you already remember before the day
  starts; mark them `remembered` vs `observed`.
- A **friction** is any moment the construction conventions failed you:
  you didn't know where code belonged, an agent improvised structure, a
  change fanned out further than it should have, a convention existed but
  wasn't followed, or a convention was missing outright.
- For "where does this belong" frictions, walk the gap survey in
  `skills/tesser-build/map.md` (name the pieces → survey what exists → the
  gap is the finding). Log which anatomy piece was involved — that's the
  doctrine file the fix lands in.

## 2. Classify each friction — the deferral test

Two legs, asked in order:

1. **Findability** — are violations of the underlying rule mechanically
   findable (a check could enumerate every site), or do they *hide* (leak
   into signatures, scatter across callers, vanish into call chains)?
2. **Fix-locality** — once a violation is found, is the fix a local
   sweep-edit (rename, wrap, per-site touch-up), or *structural* (moving
   code across boundaries, unwinding a baked-in cycle)?

**Defer only when BOTH legs hold** (findable AND local): the property that
makes a rule mechanically checkable later is the property that makes its
retrofit cheap. **Pay now when EITHER leg fails**: violations hide, or the
fix is a restructure.

Calibration anchors: perfect VO wrapping passes both legs → defer.
Dependency direction is findable but a flagged wrong-direction import is a
restructure, not a sweep → pay now.

### The pay-now universe

Every pay-now classification lands in one of five bins — the audit orders
and prunes *within* them:

1. **Context boundaries + seams** — which contexts exist; `Client` + DTOs at
   the top; no cross-context internal imports.
2. **Dependency direction** — acyclic, inward (enforcement: import-linter
   declared contracts, consumer-side).
3. **No representation leaks** — domain objects never escape through
   `Client`/handler signatures.
4. **Single construction path per type** — exactly one public construction
   path exists; no bypasses, no parent re-validation. (How the path is
   *shaped* — spec-idiom purity, canonicalization — is defer-category.)
5. **Env/exit edge discipline** — no buried `getenv`s; the host is the env
   edge.

A finding that fits none of the bins but whose violations *hide* is a
candidate **sixth bin**: relay it flagged as such. Adding a bin is a
deliberate amendment to the plan doc, never silent queue drift.

## 3. Relay the findings — de-identified

De-identification rules: no business/domain nouns, no file paths, no
identifiers — replace with anatomy role names (`context A`, `aggregate X`,
`gateway to vendor Y`, `handler H`). If the pattern survives the renaming,
it was pattern-shaped; if it doesn't, it wasn't relayable.

One entry per friction:

```
- pattern:    <generic shape, anatomy vocabulary only>
- anatomy:    <piece(s) from map.md's table that the friction touched>
- seen:       remembered | observed   (+ how many times that day)
- class:      pay-now <bin 1-5 | candidate-6th> | defer
- leg failed: findability | fix-locality | both   (pay-now only)
- cost shape: <what the retrofit looks like if deferred: N-site sweep,
               signature fan-out, boundary restructure, ...>
```

The relayed output is an **ordered pay-now list** (your priority order —
this list IS the Phase-1 queue) plus the defer pile, each defer item named
so it can become a check later. Check misfires ride the same channel: relay
the de-identified shape that confused a check, and the fixture and check get
corrected toolkit-side.
