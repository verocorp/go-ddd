# Norm: logging

<!-- tb-status: stub -->

**Not yet materialized — note the gap, don't invent a convention.** This
placeholder exists because two settled norms carved logging out and pointed
here: domain objects define no display/debug dunders (a compound has zero
conversion dunders — `serialization.md` rule 5; `repr` is the interim debug
surface), and canonical forms are wire material, not log formatting. What a
constructed app logs, where (host vs adapter vs service), in what shape
(structured vs text), and how domain values appear in log lines without
leaking representations — all of that is this norm's future scope, settled
from field evidence, not invented inline.

Until it materializes: use the language's default `repr` for debug output,
keep log statements at the edges the anatomy already sanctions (hosts,
adapters), and feed friction back as evidence.
