# CLAUDE.md — proto-svg

proto-svg turns a self-describing **SVG EBNF grammar** into a protobuf schema and
then walks that schema to emit every legal SVG value-path as renderable markup.
The grammar is the single source of truth; everything downstream is derived.

## Development rule: everything goes through the four scripts

**No one-off commands.** All build, test, and run steps MUST go through the root
orchestration scripts. A bare `go run` / `go test` / `go build` is fine as a
quick local compile check while editing, but it **NEVER counts as validation**.
Final validation before committing ALWAYS goes through the scripts.

If something needs to be built or tested, it belongs in a script. If it is not
in a script, it does not count as built or tested.

**CRITICAL: ALWAYS run `./LET_IT_RIP.sh` before EVERY `git commit` and
`git push`.** No exceptions. It is the one command that proves the whole
pipeline (grammar → proto → gallery → screenshots) is green.

## Commands

| What | Command |
|---|---|
| Setup (prereqs, `go mod tidy`, chromerpc warmup) | `./setup.sh` |
| Build (EBNF → proto → gallery) | `./build.sh` |
| Test / validate (the ONLY sanctioned pre-commit check) | `./test.sh` |
| Full pipeline + screenshots + serve | `./LET_IT_RIP.sh` |
| Compile the grammar → proto schema (called by build) | `./tools/genproto.sh` |
| Generate the value-path galleries (called by build) | `./tools/gen.sh` |
| Screenshot HTML/dir/URL via chromerpc | `./chrome-testing/snap.sh <input> <output>` |

`setup.sh` → `build.sh` → `test.sh` → `LET_IT_RIP.sh` is a strict superset chain:
each one runs the previous stages, so running the outermost runs everything.

## The pipeline (grammar is the source of truth)

```
lang/*.ebnf  ──genproto──▶  proto/svg.proto + proto/svg.fdset
(14 modules,                 + proto/pb/svg/{prefix_map,separator_map}.go
 svg.ebnf is root)
      │
      └──gen──▶  chrome-testing/html/generated/*.html  (per-element galleries)
                 + index.html + values.json + manifest.tsv
                 + chrome-testing/generated/sample-{document,rect}.svg
      │
      └──snap──▶ chrome-testing/screenshots/generated/*.png  (one PNG per page)
```

1. **Grammar (`lang/*.ebnf`, 14 modules).** A STRUCTURAL grammar for real SVG:
   the strings it derives are valid `.svg` markup (the literal `<tag>`, `="`,
   `"`, `</tag>` are grammar terminals). `svg.ebnf` is concatenated first and its
   `SvgDocument` production is the prune root. See `docs/GRAMMAR_PLAN.md` for the
   per-module file layout and `docs/EBNF_PATTERN.md` for the markup conventions.

2. **genproto (`tools/genproto.sh`).** Runs `go run ./lang/cmd/genproto/`, which
   uses **gluon** to compile the concatenated grammar into `proto/svg.proto`,
   the binary `proto/svg.fdset`, and the grammar-derived lookup tables
   `proto/pb/svg/prefix_map.go` and `proto/pb/svg/separator_map.go`. gluon lifts
   the leading markup terminals into the prefix map so the renderer can re-emit
   them.

3. **gen (`tools/gen.sh`).** Runs `go run ./chrome-testing/cmd/gen/`, which
   compiles the grammar in memory, walks the proto message graph, and enumerates
   **every element's every attribute's every value-path**. Each value is injected
   into the element's hand-authored blueprint and written as a per-element
   gallery page under `chrome-testing/html/generated/`, plus `index.html`,
   `values.json`, and `manifest.tsv`. Unlike proto-css's gen, REPEATED fields are
   RENDERED (children and attributes are genuine repetitions in SVG).

4. **snap (`chrome-testing/snap.sh`).** Screenshots a file, a directory (one PNG
   per `.html`), or a URL via headless Chrome driven by **chromerpc**.

## Templates are HAND-AUTHORED, never script-generated

Hand-authored templates live in `chrome-testing/html/template/<tag>.html`, one
per SVG element using the real tag name (`rect.html`, `linearGradient.html`,
`feGaussianBlur.html`). Each has a human showcase (visible cards) plus a hidden
machine blueprint whose `{{ELEMENT}}` slot the generator fills with
grammar-generated markup. **Do NOT write a script that emits templates — author
each by hand**, then verify it by screenshotting and *looking at the PNG* (it
must render, and every variation card must be visibly distinct). The full
authoring + verification protocol is `docs/TEMPLATE_GUIDE.md`.

The generated galleries under `chrome-testing/html/generated/` ARE
script-generated — do not edit them by hand.

## Grammar conventions (read before touching `lang/*.ebnf`)

- **Structural markup + DOM-name rules:** `docs/EBNF_PATTERN.md`. Element rules
  are named by their DOM interface (`SVGRectElement`), markup is baked as
  terminals, and value grammars carry no markup.
- **File layout, naming, collision avoidance, per-element tables:**
  `docs/GRAMMAR_PLAN.md`. CRITICAL: gluon silently dedupes rules whose names
  normalize equal (`lowercase`, underscores stripped) — every shared attribute
  name (`x`, `fill`, `type`, …) gets a fully-qualified per-element production
  (`RectXAttr`, `FeOffsetDxAttr`) selected positionally.
- **Context-free core vs constraint overlay:** `docs/CONTEXT_SENSITIVITY.md`.
  Everything context-free goes in the grammar, exhaustively (enumerate every
  keyword set in full); only provably non-context-free constraints (IDREF
  resolution, ranges/monotonicity, required/at-most-once, animation value
  typing) live in the overlay (renderer/generator/parser), never the grammar.

## chromerpc is fetched from GitHub

`chrome-testing/snap.sh` (and `setup.sh`'s warmup) clone
`https://github.com/accretional/chromerpc` and build `chromerpc` + `automate`
into `/tmp/chromerpc-testing/bin`, caching across runs. Override the source with
`CHROMERPC_GIT=`; force a rebuild with `REBUILD_CHROMERPC=1`. chromerpc setup is
best-effort and never fatal to the build/test pipeline.

## Go module

`github.com/accretional/proto-svg`, with `replace` directives pointing at the
sibling checkouts `../gluon` and `../proto-merge`. `go mod tidy` (run by
`setup.sh`) resolves them; both sibling checkouts must be present.

## Ports

- `LET_IT_RIP.sh` gallery server: **8899** (override with `SERVE_PORT=`). Stale
  listeners on the port are killed before binding.
- `chrome-testing/snap.sh` HTTP server and chromerpc: auto-assigned free ports.

## Prerequisites

- **Go** — gluon genproto + the gallery generator + chromerpc build.
- **Google Chrome** (or Chromium) — headless screenshots.
- **Python 3** — local HTTP serving for snap.sh and the gallery server.
- **git** — cloning chromerpc.

## Things to NOT do

- Do NOT run `go run` / `go test` / `go build` directly as final validation —
  run `./test.sh` (and `./LET_IT_RIP.sh` before committing).
- Do NOT commit or push without running `./LET_IT_RIP.sh` first.
- Do NOT hand-edit generated artifacts: `proto/svg.proto`, `proto/svg.fdset`,
  `proto/pb/svg/*.go`, or anything under `chrome-testing/html/generated/`. They
  are regenerated by `tools/genproto.sh` and `tools/gen.sh` — change the grammar
  or generator instead and rerun.
- Do NOT hand-edit screenshots under `chrome-testing/screenshots/` — they are
  produced by `chrome-testing/snap.sh`.
- Do NOT script-generate the hand-authored templates in
  `chrome-testing/html/template/` — author and screenshot-verify each by hand
  (`docs/TEMPLATE_GUIDE.md`).
- Do NOT fix a problem with a one-off command outside the scripts — update the
  relevant script (or grammar/generator) and rerun it.
