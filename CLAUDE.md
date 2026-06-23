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
| Build (EBNF → proto → catalogue) | `./build.sh` |
| Test / validate (the ONLY sanctioned pre-commit check) | `./test.sh` |
| Full pipeline + serve the gallery (+ `SHOOT=1` for screenshots) | `./LET_IT_RIP.sh` |
| Compile the grammar → proto schema (called by build) | `./tools/genproto.sh` |
| Generate the gallery catalogue (called by build) | `./tools/gen.sh` |
| Per-preset viewer screenshots + animation GIFs (drives the gallery) | `./chrome-testing/shoot.sh` (or `SHOOT=1 ./LET_IT_RIP.sh`) |
| Ad-hoc screenshot of any HTML/dir/URL via chromerpc | `./chrome-testing/snap.sh <input> <output>` |

`setup.sh` → `build.sh` → `test.sh` → `LET_IT_RIP.sh` is a strict superset chain:
each one runs the previous stages, so running the outermost runs everything.

## The pipeline (grammar is the source of truth)

```
lang/*.ebnf  ──genproto──▶  proto/svg.proto + proto/svg.fdset
(14 modules,                 + proto/pb/svg/{prefix_map,separator_map}.go
 svg.ebnf is root)
      │
      └──gen──▶  chrome-testing/gallery/catalogue.json  (the SVG Lab data contract:
                 per element → typed attribute controls + presets + base SVG)
                 + chrome-testing/generated/sample-{document,rect}.svg
      │
      └──shoot──▶ chrome-testing/screenshots/gallery/<tag>/NN-<slug>.png
                  (drives the gallery viewer to each preset; SMIL → GIF)
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
   **every element's every attribute's every value-path**. It emits
   `chrome-testing/gallery/catalogue.json` — the data contract the SVG Lab gallery
   renders from: per element, the showcased element's typed attribute controls
   (paint / range / select / number / text, derived from each attribute's value
   type), a baseline **base** SVG (the hand-authored blueprint with the element at
   its defaults), and **presets** (one per visually-meaningful value-path). The
   per-element static HTML + specimen technique is RETIRED — there is one live app.

4. **gallery (`chrome-testing/gallery/`).** The SVG Lab app: a standalone vanilla
   HTML+JS SPA (no framework) that loads `catalogue.json`. Three linked panes —
   the VIEWER (renders the SVG), the CONTROL PANEL (typed controls + presets), and
   the live EDITOR (the SVG markup, editable client-side) — all stay in sync, no
   server needed. `index.html` + `app.js` are hand-maintained; `catalogue.json` is
   generated. Hash routes: `#/el/<tag>` opens an element; `#/embed/<tag>/<idx>`
   shows the chrome-free viewer with preset `idx` applied (used by the shoot).

## Per-preset gallery shoot + animation GIFs (`chrome-testing/shoot.sh`)

`chrome-testing/shoot.sh` (a port of proto-css's `shoot` + `gifenc`) serves the
gallery and drives chromerpc through it: for every preset (= one attribute value)
it sets `location.hash = '#/embed/<tag>/<idx>'` and screenshots the viewer into
`chrome-testing/screenshots/gallery/<tag>/NN-<slug>.png`; for SMIL/animation
elements it captures a frame sequence that `chrome-testing/cmd/gifenc` encodes
into an animated GIF. It is long-running (one capture per preset across every
element), so it is opt-in: run `./chrome-testing/shoot.sh` directly, or
`SHOOT=1 ./LET_IT_RIP.sh`. `ONLY=rect,animate` limits it to some tags; `RESUME=1`
skips already-captured presets. Screenshots are gitignored (regenerable).

## Templates are HAND-AUTHORED, never script-generated

Hand-authored templates live in `chrome-testing/html/template/<tag>.html`, one
per SVG element using the real tag name (`rect.html`, `linearGradient.html`,
`feGaussianBlur.html`). Each has a human showcase (visible cards) plus a hidden
machine blueprint whose `{{ELEMENT}}` slot the generator fills with
grammar-generated markup. **Do NOT write a script that emits templates — author
each by hand**, then verify it by screenshotting and *looking at the PNG* (it
must render, and every variation card must be visibly distinct). The full
authoring + verification protocol is `docs/TEMPLATE_GUIDE.md`.

`chrome-testing/gallery/catalogue.json` is script-generated (by `tools/gen.sh`) —
do not edit it by hand. The gallery's `index.html` + `app.js` ARE hand-maintained
(the SVG Lab UI); only the catalogue data is generated.

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
  `proto/pb/svg/*.go`, or `chrome-testing/gallery/catalogue.json`. They are
  regenerated by `tools/genproto.sh` and `tools/gen.sh` — change the grammar or
  generator instead and rerun. (The gallery's `index.html`/`app.js` ARE editable.)
- Do NOT hand-edit screenshots under `chrome-testing/screenshots/` — they are
  produced by `chrome-testing/shoot.sh`.
- Do NOT script-generate the hand-authored templates in
  `chrome-testing/html/template/` — author and screenshot-verify each by hand
  (`docs/TEMPLATE_GUIDE.md`).
- Do NOT fix a problem with a one-off command outside the scripts — update the
  relevant script (or grammar/generator) and rerun it.
