# srv — the hosts

<!-- tb-status: full -->

An app-wide directory of **hosts, one per delivery mechanism** (recommended
subdirs `srv/{http,cli,wrk}`, not enforced). A host's `main` is the outermost
edge of the app: it decodes the environment into the app `Config`, calls
`bootstrap.new(cfg)` **once**, mounts *its* mechanism's inbound handlers
across all contexts, applies cross-cutting middleware
(auth/logging/recovery), and owns the process lifecycle. Everything a host
does is edge work — the moment logic appears in a host that isn't
env-decoding, mounting, middleware, or lifecycle, it belongs somewhere below.

## Is this what I'm building?

**Test:** *Am I writing the process entry point for one delivery mechanism —
the `main` that reads the environment, builds the app once, and serves?*
Yes → a host.

**Near-misses that are NOT a host:**
- A **handler** (`handlers.md`) — per-context wire↔`Client` translation. The
  host *mounts* handlers; a handler never owns the server, the middleware, or
  the process.
- The **composition root** (`bootstrap.md`) — builds the object graph from a
  `Config`; it never reads the environment and never serves. The host calls
  it; it is not it.
- A **worker loop / consumer** that polls a queue — that *is* a host
  (`srv/wrk`): same edge duties, different mechanism.
- A **test fixture** that builds the app — tests construct via
  `bootstrap.new(cfg)` with a literal `Config`; they are not an env edge and
  never read one.

## Rules

1. **The host is the env edge.** Each `srv/*/main` populates the spec-shaped
   app `Config` directly from the environment — including its **own launch
   config** (the listen addr, the worker cadence) — and hands it to
   `bootstrap.new`, which validates fail-fast. Nothing below the host reads
   the environment (locked by
   `examples/python-app/tests/test_enforcement.py`): a `getenv` below the
   edge is a hidden deploy surface with a silent default. There is no shared
   env-decoder module — each host's needs differ, and a shared decoder
   becomes a second config authority.
2. **Only the edge exits.** Exit/fatal calls live in `srv/*/main`, nothing
   below (same enforcement test) — a library that exits takes the process
   away from the one place entitled to decide that.
3. **One graph per process.** The host calls `bootstrap.new` once at startup
   (locked by `examples/python-app/tests/test_bootstrap_once.py`) and owns
   shutdown via `App.close()` — in a `finally`, so a crashing serve loop
   still unwinds the graph.
4. **Two-layer transport split.** The per-context handler translates
   wire ↔ `Client` (`handlers.md`); the host mounts handlers and owns the
   server + middleware. Auth *policy*, logging, recovery, rate limits are
   host middleware, never inside a context's handler — a handler that
   imports another context to do auth has leaked a host concern into a
   context adapter.
5. **Hosts share nothing but the app.** Two hosts are two processes; they
   share the composition root and the contexts, not memory. A CLI host runs
   against its *own* `App`; if two mechanisms must see one state, that state
   lives behind a context's repository, not in a host.

## Shape

```
srv/
  http/main.py       ← env → Config, new(cfg) once, mount handlers, serve, close
  cli/main.py        ← env → Config, new(cfg) once, run command, close

def main() -> None:
    cfg = Config(campaign=CampaignConfig(storage=os.getenv("CAMPAIGN_STORAGE") or ""), ...)
    app = new(cfg)                      # once per process; validates fail-fast
    server = make_server((host, port), app)   # mounts the contexts' handlers
    try:
        server.serve_forever()
    finally:
        app.close()
```

A missing env var stays an empty coordinate and `bootstrap.new` fails fast on
it — the host never invents a default for someone else's config; its own
launch knobs (a listen port) may default locally. Construction mechanics:
`python.md#inbound-handlers-and-hosts`; verified impl:
`examples/python-app/srv/` (`http/main.py`, `cli/main.py`).

## Decisions you must make

1. **Which mechanisms get a host?** One per delivery mechanism actually
   served — `http`, `cli`, `wrk` are the recommended names, not a quota. A
   mechanism you don't serve gets no stub.
2. **Where does secret resolution happen?** Resolving secret *references*
   (Vault/AWS/GCP) is a legitimate host-side, launch-time concern — it is
   part of env → `Config` decoding at the edge, never a lazy fetch below it.
   The template deliberately doesn't build the loader.
3. **How much lifecycle?** The template mandates only build-once +
   `close()` in `finally`. Graceful-shutdown ordering, drain, readiness are
   the host's fill-in — do them properly at the edge when the service needs
   them (see the ops-deferral notice in `SKILL.md`).

## How the machine sees it

Machine-checked in the verified impl (`tests/test_enforcement.py`, real `ast`
checks with injected-violation teeth): env reads (`os.getenv`/`os.environ`)
only in `srv/*/main`; exits only in `srv/*/main`; no import-time side
effects in contexts or bootstrap. Build-once is locked by
`tests/test_bootstrap_once.py`. A generalized tessercheck check is scheduled
follow-on work, not yet shipped. Review-side tells:
- an **env read anywhere below `srv/`** — the deploy surface went invisible;
- a **second `bootstrap.new` call** in request/command handling — per-request
  wiring;
- **route/domain logic in a host** — the host is mount + middleware; a
  `for`-loop over domain objects here belongs in an application service.

## Tests you must write

- **Env reads only at the edge** — an enforcement test that walks the tree
  and fails on `getenv`/`environ` outside `srv/*/main` (verified impl:
  `test_enforcement.py`; prove it has teeth on an injected violation).
- **Exits only at the edge** — same walk, `sys.exit`/`os._exit`.
- **The graph is built once and closed** — a host-shaped test that calls
  `bootstrap.new` once, exercises a `Client`, and `close()`s (idempotently)
  (verified impl: `test_bootstrap_once.py`).

## Common mistakes

- **A shared env-decoder module.** `config/from_env.py` used by every host —
  now there are two config authorities and the host is no longer the edge.
  Each host populates the `Config` it needs, inline.
- **Defaulting a peer's coordinate.** `os.getenv("CAMPAIGN_STORAGE") or
  "memory"` at the host — the silent volatile-storage fall, moved up a
  layer. Empty coordinate in, fail-fast in `bootstrap.new`.
- **Auth in a handler.** Token checking inside a context's handler — auth
  policy is host middleware; the handler receives an authenticated request.
- **Per-request construction.** Building the app (or a repository) inside
  the request path — once per process, at startup.
- **The immortal serve loop.** No `finally: app.close()` — a crash leaks
  every pool the graph holds.

## Now build it

<!-- tb-allow-missing: examples/app -->

- Python: `python.md#inbound-handlers-and-hosts` — the host `main` shape,
  backed by `examples/python-app/srv/`.
- Go: not yet materialized — the settled anatomy's Go mirror
  (`examples/app`) is pending; note the gap, don't invent a convention.
  Mirror the Python arc's structure (env edge, build once, mount, `defer
  app.Close()`).
