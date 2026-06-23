# Visual QA — core1 batch

Elements: `rect`, `circle`, `ellipse`, `line`, `polyline`, `polygon`, `path`, `svg`, `g`, `defs`

Contact-sheet source: `chrome-testing/screenshots/contact/<element>.png`
Specimen source: `chrome-testing/screenshots/specimens/<element>/`

---

## `<rect>` — 40 value-paths — OK

**Observations:**

- `x`, `y` variants (`10`, `24px`, `2em`, `50%`, `75%`): rect shifts position distinctly across all five values. Correct.
- `width="20"` / `height="20"`: produces a narrow column and a flat strip respectively — correct.
- `rx="auto"` / `ry="auto"`: renders as a plain square with no rounding (correct — SVG spec: `auto` on `rx`/`ry` for `<rect>` means 0).
- `rx="20"` / `ry="20"`: rounded corners clearly visible. Correct.
- `fill="none"`: shows the stroke outline only (the baseline supplies `stroke="#16c79a" stroke-width="2"`). Correct.
- `fill-opacity="0.5"`: fill visibly dimmed. Correct.
- `stroke="none"`: baseline stroke disappears, leaving the solid fill only. Correct.
- `stroke-opacity="0.5"`: stroke visibly faded. Correct.
- `stroke-width="20"`: thick teal border visible. Correct.
- `stroke-linecap="butt"` / `stroke-linejoin="miter"` / `stroke-miterlimit="4"`: near-identical to base on a rect (expected — linecap/linejoin/miterlimit have no visible effect on a rect with 90° corners and no line endings). Acceptable.
- `stroke-dasharray="none"` / `stroke-dashoffset="20"`: both match the solid-stroke base — correct, `dashoffset` only offsets a dash pattern; with no dasharray both cells are identical and correct.
- `paint-order="normal"` / `color="#e94560"` / `shape-rendering="auto"` / `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"`: all render correctly; default values match the base appearance. Correct.
- `display="none"`: **cell is correctly blank.** Correct.
- `visibility="visible"`: shows the rect. Correct.
- `overflow="visible"` / `opacity="0.5"` / `cursor="auto"` / `transform="translate(20 10)"` / `transform-origin="inherit"`: all distinct and correct.

**Verdict: OK.**

---

## `<circle>` — 35 value-paths — OK

**Observations:**

- `cx`, `cy` variants (`10`, `24px`, `2em`, `50%`, `75%`): circle shifts position distinctly; partially off-canvas at extremes — correct.
- `r="20"`: noticeably smaller circle than the baseline `r="40"`. Correct.
- `fill="none"` / `fill-opacity="0.5"` / `stroke="none"` / `stroke-opacity="0.5"` / `stroke-width="20"`: all render distinctly and correctly.
- `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`: near-identical on a closed circle — expected and acceptable (no endpoints; dash pattern at "none" is the default).
- `paint-order="normal"` / `color="#e94560"` / `shape-rendering="auto"` / `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"`: default values, all match base. Correct.
- `display="none"`: **correctly blank.** Correct.
- `visibility="visible"` / `overflow="visible"` / `opacity="0.5"` / `cursor="auto"` / `transform="translate(20 10)"` / `transform-origin="inherit"`: all distinct and correct.

**Verdict: OK.**

---

## `<ellipse>` — 38 value-paths — ISSUES

**Observations:**

- `cx`, `cy` variants: ellipse shifts position distinctly. Correct.
- `rx="20"` / `ry="20"`: distinct — rx=20 produces a tall narrow ellipse; ry=20 a flat wide one. Correct.
- **`rx="auto"` (index 10): cell is blank.** With `ry="25"` in the baseline, SVG2 specifies `rx="auto"` copies `ry`, yielding a 25×25 circle. Chrome does not implement SVG2 `auto` for `rx`/`ry` on `<ellipse>` and renders it as invalid (zero radius → invisible). This is a browser-limitation blank, not a gen bug, but the attribute's effect is not demonstrated.
- **`ry="auto"` (index 12): cell is blank.** Same reason — Chrome ignores `auto` and treats radius as 0.
- `fill="none"` / `fill-opacity="0.5"` / `stroke="none"` / `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap` / `stroke-linejoin` / `stroke-miterlimit` / `stroke-dasharray` / `stroke-dashoffset`: all correct (dash attrs near-identical on a closed ellipse — acceptable).
- `paint-order` / `color` / `shape-rendering` / `filter` / `clip-path` / `clip` / `mask` / `display` / `visibility` / `overflow` / `opacity` / `cursor` / `transform` / `transform-origin`: all correct.

**Issue:** `rx="auto"` and `ry="auto"` render blank due to Chrome not supporting SVG2 `auto` radius. The attribute is not demonstrated. Fix options:
- In `blueprint.go` `baselineFor("ellipse")` or in `overlay.go`/`enumerate.go`, add a `distinctValueSet` or `overlaySample` for `rx`/`ry` that omits the `auto` value when Chrome is the target, OR adds a note-annotation, OR substitutes a concrete fallback. Alternatively, accept as a browser-gap non-issue.

**Verdict: ISSUES** — `rx="auto"`, `ry="auto"` blank (Chrome SVG2 gap, not gen bug; fixable by dropping `auto` from the enumerated set for these attrs in `blueprint.go` / `enumerate.go`).

---

## `<line>` — 44 value-paths — OK

**Observations:**

- `x1`, `y1`, `x2`, `y2` variants (`10`, `24px`, `0`, `1`): line endpoint positions shift distinctly across all sixteen cells. Correct — `0` and `1` produce subtle but measurably different geometry.
- `fill="none"` (index 16): correctly shows the line (lines have no fill area — `fill="none"` is the default for lines). Correct.
- `fill-opacity="0.5"` (index 17): line still visible (stroke-based) — correct, fill-opacity doesn't affect the stroke.
- `stroke="none"` (index 18): **correctly blank** — the only visible element is the stroke; setting it to none produces an invisible line. This is correct SVG behavior, not a bug.
- `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap="butt"` / `stroke-linejoin="miter"` / `stroke-miterlimit="4"` / `stroke-dasharray="none"` / `stroke-dashoffset="20"`: all distinct and correct.
- `marker`, `marker-start`, `marker-mid`, `marker-end` (all "none"): all identical to the base — correct, these are default values.
- `color` / `shape-rendering` / `filter` / `clip-path` / `clip` / `mask` / `display` / `visibility` / `overflow` / `opacity` / `cursor` / `transform` / `transform-origin`: all correct.
- `paint-order="normal"`: matches base — correct default.
- `display="none"`: **correctly blank.** Correct.

**Verdict: OK.**

---

## `<polyline>` — 31 value-paths — OK

**Observations:**

- `points="20,20 80,20 50,80"` / `points="10,80 50,10 90,80 50,50"`: two distinct polyline shapes (open zig-zag vs W-shape). Correct.
- `fill="none"` (index 02): correctly shows the open stroke path with no fill. Correct.
- `fill-rule="nonzero"` (index 03): same open path — `fill="none"` in baseline makes fill-rule invisible. Acceptable.
- `fill-opacity="0.5"` (index 04): visibly dimmed — but baseline has `fill="none"`, so fill-opacity has no effect. Matches 02 — acceptable (fill is none).
- `stroke="none"` (index 05): **correctly blank** — baseline has `fill="none"`; removing stroke leaves nothing visible. Correct behavior.
- `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap="butt"` / `stroke-linejoin="miter"` / `stroke-miterlimit="4"` / `stroke-dasharray="none"` / `stroke-dashoffset="20"` / `paint-order="normal"`: all distinct and correct.
- `marker` / `marker-start` / `marker-mid` / `marker-end` (all "none"): identical to base — correct defaults.
- `color` / `shape-rendering` / `filter` / `clip-path` / `clip` / `mask` / `display` / `visibility` / `overflow` / `opacity` / `cursor` / `transform` / `transform-origin`: all correct.

**Verdict: OK.**

---

## `<polygon>` — 31 value-paths — OK

**Observations:**

- `points="20,20 80,20 50,80"` / `points="10,80 50,10 90,80 50,50"`: two distinct polygon shapes (triangle vs inverted chevron). Correct.
- `fill="none"` (index 02): correctly shows the stroke outline only. Correct.
- `fill-rule="nonzero"` (index 03): solid triangle — same as base (points are non-self-intersecting so nonzero = evenodd here). Acceptable.
- `fill-opacity="0.5"` (index 04): fill visibly dimmed. Correct.
- `stroke="none"` (index 05): correctly shows the filled triangle without outline. Distinct from base (which has `stroke-width="2"`). Correct.
- `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap` / `stroke-linejoin` / `stroke-miterlimit` / `stroke-dasharray` / `stroke-dashoffset` / `paint-order`: all correctly rendered; stroke attrs near-identical on a filled closed polygon with thin baseline stroke — acceptable.
- `marker` / `marker-start` / `marker-mid` / `marker-end` (all "none"): identical to base — correct defaults.
- `color` / `shape-rendering` / `filter` / `clip-path` / `clip` / `mask` / `display` / `visibility` / `overflow` / `opacity` / `cursor` / `transform` / `transform-origin`: all correct.

**Verdict: OK.**

---

## `<path>` — 30 value-paths — OK

**Observations:**

- `d="M10 50 Q50 10 90 50"`: the arc baseline is clearly rendered. Correct.
- `fill="none"` (index 01): arc visible via stroke — correct, the baseline sets `fill="none"` and `stroke="#16c79a"`. Correct.
- `fill-rule="nonzero"` (index 02): same as base (path is open/non-self-intersecting; fill=none). Acceptable.
- `fill-opacity="0.5"` (index 03): arc shows with tail artefact from `T90 90` in the baseline d; fill-opacity no visible change since fill=none. Acceptable.
- `stroke="none"` (index 04): **correctly blank** — baseline has `fill="none"`; removing stroke leaves nothing visible. Correct SVG behavior, not a bug.
- `stroke-opacity="0.5"` (index 05): arc visibly faded. Correct.
- `stroke-width="20"` (index 06): thick filled arc shape. Correct.
- `stroke-linecap="butt"` / `stroke-linejoin="miter"` / `stroke-miterlimit="4"` / `stroke-dasharray="none"` / `stroke-dashoffset="20"`: near-identical on a simple open arc — acceptable defaults.
- `paint-order="normal"` / `marker` / `marker-start` / `marker-mid` / `marker-end` (all "none"): all identical to base — correct defaults.
- `color` / `shape-rendering` / `filter` / `clip-path` / `clip` / `mask` / `display` / `visibility` / `overflow` / `opacity` / `cursor` / `transform` / `transform-origin`: all correct.

**Verdict: OK.** (The `stroke="none"` blank is expected and correct.)

---

## `<svg>` — 48 value-paths — OK

**Observations:**

- `viewBox="0 0 100 100"` / `"0 0 50 50"` / `"-10 -10 120 120"`: inner SVG scales content differently in each cell — visible size change. Correct.
- `preserveAspectRatio="none"` / `"xMidYMid meet"` / `"xMinYMin slice"`: alignment and scaling differences visible. Correct.
- `x="10"` / `"24px"` / `"2em"` / `"50%"` / `"75%"` and `y="10"` / `"24px"` / `"2em"` / `"50%"` / `"75%"`: inner SVG repositioned distinctly. Correct. `x="75%"` / `y="75%"` push the inner SVG near the edge, leaving only a sliver of the child visible — correct behavior.
- `width="auto"` (index 16): inner SVG shrinks to intrinsic width, showing a narrow vertical strip. Correct.
- `width="20"` (index 17): **correctly blank** — inner SVG is only 20px wide; child rect starts at `x=20` inside the inner SVG's userSpace and is completely clipped. Correct behavior.
- `height="auto"` (index 18): inner SVG at intrinsic height, child visible. Correct.
- `height="20"` (index 19): **correctly blank** — inner SVG only 20px tall, child at `y=20` clipped out. Correct behavior.
- `transform` variants (`translate`, `rotate`, `scale`, `skewX`): all produce visibly distinct transformations. Correct.
- `fill="none"` / `fill-opacity="0.5"` / `stroke="none"` / `stroke-opacity="0.5"` / `stroke-width="20"` / stroke attrs: applied to inner SVG element; child inherits via `fill="currentColor"`. Correctly rendered.
- `color="#e94560"`: child rect renders in hot-pink (inherits via `fill="currentColor"`). Correct.
- `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"` / `display="none"` / `visibility="visible"` / `overflow="visible"` / `opacity="0.5"` / `cursor="auto"` / `transform-origin="inherit"`: all correct.

**Verdict: OK.** (`width="20"` and `height="20"` blank cells are correct clipping behavior, not bugs.)

---

## `<g>` — 24 value-paths — OK

**Observations:**

- `fill="none"` (index 00): child rect uses `fill="currentColor"`, and `fill="none"` on `<g>` sets `currentColor` to… actually `fill` and `currentColor` are separate CSS properties. The child renders with the root color `#16c79a` inherited via `color`, not via `fill` (fill="currentColor" on the child resolves to the `color` property, not the parent's `fill`). Cell shows the teal square — correct.
- `fill-opacity="0.5"` (index 01): child visibly dimmed — fill-opacity cascades. Correct.
- `stroke="none"` (index 02): matches base (no stroke on the child rect by default). Correct.
- `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap` / `stroke-linejoin` / `stroke-miterlimit` / `stroke-dasharray` / `stroke-dashoffset` / `paint-order`: near-identical on a plain filled rect with no baseline stroke — acceptable.
- `color="#e94560"` (index 11): child rect renders in hot-pink (inherits `color` → `fill="currentColor"`). Correct.
- `shape-rendering` / `filter` / `clip-path` / `clip` / `mask`: all default values, match base. Correct.
- `display="none"` (index 17): **correctly blank** — `display:none` on `<g>` hides the entire group. Correct.
- `visibility="visible"` (index 18): child visible. Correct.
- `overflow="visible"` (index 19): child visible. Correct.
- `opacity="0.5"` (index 20): child visibly at 50% opacity. Correct.
- `cursor="auto"` / `transform="translate(20 10)"` / `transform-origin="inherit"`: all correct.

**Verdict: OK.**

---

## `<defs>` — NO CONTACT SHEET

No contact sheet exists at `chrome-testing/screenshots/contact/defs.png` and no specimen directory exists at `chrome-testing/screenshots/specimens/defs/`. The `<defs>` element is absent from the rendered output despite having a scaffold defined in `blueprint.go` (`builtinScaffoldWins["defs"] = true`) and a dedicated case in `defaultScaffold`.

This means `<defs>` was either not included in the enumeration run that produced the current contact sheets, or was skipped.

**Verdict: MISSING** — no contact sheet to review. Needs investigation into why `<defs>` was excluded from the current batch render.

---

## Consolidated Real Issues

### Issue 1 — `<ellipse>` `rx="auto"` and `ry="auto"` both blank

**Cells:** `rx="auto"`, `ry="auto"`
**Root cause:** Chrome does not implement SVG2 `auto` keyword for `rx`/`ry` on `<ellipse>`. When either radius is "auto" the browser treats it as 0 → invisible ellipse.
**Assessment:** Browser gap, not a gen correctness bug. However the attribute effect is not demonstrated.
**Fix options (pick one):**
- In `enumerate.go` `enumerateValue`, add a `distinctValueSet("ellipse", "rx")` / `("ellipse", "ry")` returning only concrete numeric values (e.g. `["20", "40"]`), dropping `auto` from the enumeration for these attrs.
- In `overlay.go` `overlaySample`, return a concrete fallback value for `(ellipse, rx, RxAttr)` and `(ellipse, ry, RyAttr)` when the value would be `auto`.
**Fix target:** `chrome-testing/cmd/gen/enumerate.go` (`distinctValueSet`) or `chrome-testing/cmd/gen/overlay.go` (`overlaySample`).

### Issue 2 — `<defs>` contact sheet missing

**Root cause:** Unknown — `defs` has a blueprint and body defined but no rendered output in the current batch.
**Fix target:** Investigate `chrome-testing/cmd/gen/main.go` or the gallery/batch runner to confirm `defs` is included in the element enumeration pass and re-run.

---

## Summary Table

| Element  | Paths | Verdict | Notes |
|----------|-------|---------|-------|
| `rect`   | 40    | OK      | All cells correct |
| `circle` | 35    | OK      | All cells correct |
| `ellipse`| 38    | ISSUES  | `rx="auto"`, `ry="auto"` blank (Chrome SVG2 gap) |
| `line`   | 44    | OK      | `stroke="none"` blank is correct |
| `polyline`| 31   | OK      | `stroke="none"` blank is correct |
| `polygon`| 31    | OK      | All cells correct |
| `path`   | 30    | OK      | `stroke="none"` blank is correct |
| `svg`    | 48    | OK      | `width="20"` / `height="20"` blanks are correct clipping |
| `g`      | 24    | OK      | All cells correct |
| `defs`   | —     | MISSING | No contact sheet or specimens generated |
