# Visual QA — Paint / Defs Batch

**Batch:** `script`, `style`, `linearGradient`, `radialGradient`, `stop`, `pattern`, `marker`, `clipPath`, `mask`, `filter`
**Date:** 2026-06-23
**Source:** contact-sheet PNGs in `chrome-testing/screenshots/contact/`

---

## Per-Element Verdicts

### `<script>` — OK

3 value-paths: `href="#slot"`, `xlink:href="#slot"`, `xlink:title="label"`.

- `href="#slot"` → coral/red filled square. Correct.
- `xlink:href="#slot"` → same coral square. Expected (both link forms resolve to the same slot; identical render is correct).
- `xlink:title="label"` → teal/green filled square — visually distinct from the href cells, confirming the title-only path renders differently (different script payload or no-op).

All cells non-blank, no all-identical-values issue. **PASS.**

---

### `<style>` — ISSUES

12 value-paths across `type` (4 values) and `media` (4 values) and `title` (4 values).

**`type` row (4 cells):** `type="label"`, `type="Aa"`, `type="sample"`, `type="specimen"` — all four render a **black rounded-rectangle with no interior content**. The cells are non-blank but visually identical. The style element is presumably inactive or rejected when `type` is set to a non-standard value (anything other than `text/css`), so the CSS it contains does not apply and the fallback black-rect base shape shows. This is technically correct browser behavior (unknown `type` → style ignored), but it means all four `type` cells look the same — a user skimming the sheet cannot distinguish them.

**`media` row (4 cells):** same pattern — all four black rounded-rectangles, identical. Same root cause: the media-query values (`label`, `Aa`, `sample`, `specimen`) don't match any real media type, so the style block does not apply.

**`title` row (4 cells):** `title="label"` through `title="specimen"` — all four show the **styled** pink/coral circle on teal rounded-rect. These are correct and distinct from the type/media cells (title doesn't affect CSS application), but the four title values are all visually identical to each other (which is expected — `title` is metadata only).

**Real issues:**
1. `type` attribute: all 4 values produce identical blank-styled cells — entire attribute's range has no visual variation. The test values (`label`, `Aa`, `sample`, `specimen`) are not valid CSS type strings, so the style is universally suppressed. **Fix target:** include at least one valid value (`type="text/css"`) to produce a styled baseline alongside the invalid-type variants. Without it, the attribute row is effectively untestable visually.
2. `media` attribute: same problem — all 4 media values are syntactically invalid/non-matching, producing identical unstyled cells. **Fix target:** include a valid media value (e.g. `media="all"` or `media="screen"`) so at least one cell shows styled output.

---

### `<linearGradient>` — OK

28 value-paths. Full review:

- **`gradientUnits`** (`userSpaceOnUse`, `objectBoundingBox`): both show pink↔teal gradients, slightly different rendering scale — distinct and correct.
- **`gradientTransform`** (`translate(20 10)`, `rotate(45)`, `scale(1.5)`, `skewX(20)`): all four produce visually distinct gradient orientations/positions. Correct.
- **`x1`** (`10`, `24px`, `0`, `1`): four distinct gradient start-positions visible. Correct.
- **`y1`** (`10`, `24px`, `0`, `1`): four distinct gradient vertical-start positions. Correct.
- **`x2`** (`10`, `24px`, `0`, `1`): distinct horizontal end-positions. Correct.
- **`y2`** (`10`, `24px`, `0`, `1`): distinct vertical end-positions. Correct.
- **`spreadMethod`** (`pad`, `reflect`, `repeat`): `pad` shows a standard gradient; `reflect` shows symmetric bounce; `repeat` shows a narrow repeating stripe — all three clearly distinct. Correct.
- **`href="#refgrad"`** and **`xlink:href="#refgrad"`**: both show inherited gradient appearance — correct (referencing a defined gradient). Visually equivalent to each other (same spec behavior).
- **`xlink:title="label"`**: renders a gradient — title is metadata, no visual change expected. Correct.

All cells non-blank, all attribute groups show variation. **PASS.**

---

### `<radialGradient>` — OK

32 value-paths. Full review:

- **`gradientUnits`** (`userSpaceOnUse`, `objectBoundingBox`): distinct radial blobs, different scale. Correct.
- **`gradientTransform`** (`translate(20 10)`, `rotate(45)`, `scale(1.5)`, `skewX(20)`): distinctly shifted/rotated/scaled focal blobs. Correct.
- **`cx`** (`10`, `24px`, `0`, `1`): focal center shifts horizontally across cells. Correct.
- **`cy`** (`10`, `24px`, `0`, `1`): focal center shifts vertically. Correct.
- **`r`** (`20`, `2`): `r="20"` shows a wide radial fill; `r="2"` shows a nearly-point focal — visually distinct. Correct.
- **`fx`** (`10`, `24px`, `0`, `1`): focal point shifts. Correct.
- **`fy`** (`10`, `24px`, `0`, `1`): focal point shifts vertically. Correct.
- **`fr`** (`20`, `2`): distinct focal-radius sizes. Correct.
- **`spreadMethod`** (`pad`, `reflect`, `repeat`): clearly distinct ring patterns. Correct.
- **`href="#refgrad"`**, **`xlink:href="#refgrad"`**, **`xlink:title="label"`**: all render gradients; href/xlink:href are equivalent; title is metadata-only. Correct.

All cells non-blank, variation present across all attribute groups. **PASS.**

---

### `<stop>` — ISSUES

6 value-paths.

- **`offset="0"`**: three-color gradient band (pink/orange/blue) — stop at position 0. Correct.
- **`offset="1"`**: gradient band shifted — stop at position 1 (far end). Correct.
- **`offset="50%"`**: midpoint stop. Distinct gradient banding. Correct.
- **`offset="100%"`**: stop at 100% — visually similar to `offset="1"` (100% = 1.0 in SVG). Acceptable coincidence (clamped equivalents).
- **`stop-color="currentColor"`**: renders with a gradient using the inherited `currentColor` value. The visual shows a gradient with different color composition — correct if `currentColor` resolves to a non-default color via inheritance. Distinct from offset cells. Correct.
- **`stop-opacity="0.5"`**: gradient with semi-transparent stop — color bands appear washed/faded at one end. Visually distinct. Correct.

Minor note: `offset="1"` and `offset="100%"` look very similar (expected — they are equivalent values). No bug.

**PASS.**

---

### `<pattern>` — ISSUES

43 value-paths. Extensive sheet. Full review:

- **`patternUnits`** (`userSpaceOnUse`, `objectBoundingBox`): dot-grid patterns at different scales — distinct. Correct.
- **`patternContentUnits`** (`userSpaceOnUse`, `objectBoundingBox`): `objectBoundingBox` shows triangles/triangular tile content; `userSpaceOnUse` shows dot grid. Distinct. Correct.
- **`patternTransform`** (`translate(20 10)`, `rotate(45)`, `scale(1.5)`, `skewX(20)`): all four show distinctly transformed tile grids. Correct.
- **`x`** (`10`, `24px`, `2em`, `1.5rem`, `50%`, `12pt`): all show offset dot patterns — appear mostly similar but with subtle positional offsets. Some pairs (e.g. `2em`/`1.5rem`) may be very close in resolved pixel value. Allowable.
- **`y`** (`10`, `24px`, `2em`, `1.5rem`, `50%`, `12pt`): same — subtle positional differences. Allowable.
- **`width="20"`**: dot grid — correct tile width. Correct.
- **`height="20"`**: dot grid — correct tile height. Correct.
- **`viewBox`** (`0 0 100 100`, `0 0 50 50`, `-10 -10 120 120`): small-dot grid, larger-dot grid, different content-scaling — all three visually distinct. Correct.
- **`preserveAspectRatio`** (`none`, `xMidYMid meet`, `xMinYMin slice`): `none` shows squished dots; `xMidYMid meet` shows normal dots; `xMinYMin slice` shows cropped/scaled dots — all distinct. Correct.
- **`href="#slot"`**: dot pattern referencing another pattern — correct.
- **`fill="none"`**: dots rendered as outlines (unfilled). Visually distinct. Correct.
- **`fill-opacity="0.5"`**: dots rendered with semi-transparent fill — gray/muted appearance. Distinct from full-opacity. Correct.
- **`stroke="none"`**: dots with no stroke — correct.
- **`stroke-opacity="0.5"`**: faint stroke on dots — distinct. Correct.
- **`stroke-width="20"`**: very thick stroke on dots — visually distinct (blobs rather than outlined circles). Correct.
- **`stroke-linecap="butt"`**: correct linecap on stroked elements. Correct.
- **`stroke-linejoin="miter"`**: correct linejoin. Correct.
- **`stroke-miterlimit="4"`**: subtle miter limit — likely no visible change on circular dots (expected). Allowable.
- **`stroke-dasharray="none"`**: no dash — matches default. Correct.
- **`stroke-dashoffset="20"`**: dashed stroke offset — visible dash pattern. Correct.
- **`paint-order="normal"`**: same as default. Allowable.
- **`color="#e94560"`**: dot color changes to a red/pink — visually distinct from default orange-yellow. Correct.
- **`xlink:href="#slot"`**: same as `href` variant. Correct.
- **`xlink:title="label"`**: title metadata — no visual change. Correct.

**Real issue:**
1. `patternContentUnits="objectBoundingBox"` cell shows **triangular/chevron shapes** while all other cells show dot circles. This could be intentional (a different test shape for that case) or a content-mismatch bug where the wrong shape is rendered. **Flag for review:** if the pattern content shape should be consistent (dots) across all `patternContentUnits` values, the `objectBoundingBox` cell is using a different content element. Fix target: ensure the pattern content shape is the same across both `patternContentUnits` values, or confirm the triangle is deliberate.

---

### `<marker>` — OK

44 value-paths. The sheet shows arrow-on-line compositions throughout. Full review:

- **`markerUnits`** (`strokeWidth`, `userSpaceOnUse`): arrowhead size differs relative to line stroke — distinct. Correct.
- **`markerWidth`** (`20`, `2`): wide vs. narrow arrowhead — clearly distinct. Correct.
- **`markerHeight`** (`20`, `2`): tall vs. short arrowhead — clearly distinct. Correct.
- **`refX`** (`10`, `24px`, `0`, `1`, `left`, `center`, `right`): arrow anchor shifts along x-axis — multiple distinct positions. Correct.
- **`refY`** (`10`, `24px`, `0`, `1`, `_top`, `center`, `bottom`): arrow anchor shifts along y-axis. Correct.
- **`orient`** (`auto`, `auto-start-reverse`, `0deg`, `45deg`, `0`, `1`): `0deg` and `auto` and `45deg` all show distinct arrow orientations. `0` and `1` (numeric) are less obviously distinct but within expected behavior. Correct.
- **`viewBox`** (`0 0 100 100`, `0 0 50 50`, `-10 -10 120 120`): distinct marker content scaling. Correct.
- **`preserveAspectRatio`** (`none`, `xMidYMid meet`, `xMinYMin slice`): distinct arrow shapes/scaling. Correct.
- **`fill="none"`**: unfilled arrowhead outline. Distinct. Correct.
- **`fill-opacity="0.5"`**: semi-transparent arrowhead (pinkish-gray/rosy tone on marker). Distinct. Correct.
- **`stroke="none"`**: no stroke on marker. Correct.
- **`stroke-opacity="0.5"`**: faint stroke. Distinct. Correct.
- **`stroke-width="20"`**: very thick stroke on marker shape. Distinct. Correct.
- **`stroke-linecap="butt"`**, **`stroke-linejoin="miter"`**, **`stroke-miterlimit="4"`**: subtle presentation attrs on a triangle marker — may be indistinguishable at this scale. Allowable.
- **`stroke-dasharray="none"`**: no dash. Correct.
- **`stroke-dashoffset="20"`**: dashed marker outline. Correct.
- **`paint-order="normal"`**: same as default. Allowable.
- **`color="#e94560"`**: marker color changes — distinct. Correct.

All cells non-blank. All meaningful attributes show visual variation. **PASS.**

---

### `<clipPath>` — ISSUES

3 value-paths: `clipPathUnits="userSpaceOnUse"`, `clipPathUnits="objectBoundingBox"`, `fill-rule="nonzero"`.

- **`clipPathUnits="userSpaceOnUse"`**: pink/coral filled circle — clipping applied. Correct.
- **`clipPathUnits="objectBoundingBox"`**: same pink circle, visually identical to `userSpaceOnUse`. These two values can produce the same result when the clipped content fits the bounding box exactly — allowable coincidence.
- **`fill-rule="nonzero"`**: pink circle, same appearance. For a simple (non-self-intersecting) circle clip path, `nonzero` vs. `evenodd` fill-rule produces identical output — allowable.

**Concern:** All 3 cells are visually identical circles. While each individual coincidence is explainable:
- `clipPathUnits` pair: both units systems align when content == bounding box.
- `fill-rule` on simple shape: no self-intersection, so rule has no effect.

However, the sheet as a whole provides zero visual differentiation. **Fix target (low priority):** use a compound/self-intersecting path as the clip shape to make `fill-rule="evenodd"` vs `fill-rule="nonzero"` distinguishable, and use offset content to make `userSpaceOnUse` vs `objectBoundingBox` distinguishable.

**PASS with note** (all coincidences are technically correct; no broken cells).

---

### `<mask>` — OK

18 value-paths.

- **`maskUnits`** (`userSpaceOnUse`, `objectBoundingBox`): both show teal squares with masking applied — slightly different mask geometry. Correct.
- **`maskContentUnits`** (`userSpaceOnUse`, `objectBoundingBox`): `objectBoundingBox` produces a noticeably larger mask window (larger visible square). Distinct. Correct.
- **`x`** (`10`, `24px`, `2em`, `50%`, `75%`): mask viewport shifts horizontally — the visible portion of the content moves right; `50%` and `75%` show progressively more extreme shifts with very small remaining squares. All distinct. Correct.
- **`y`** (`10`, `24px`, `2em`, `50%`, `75%`): same for vertical. Distinct. Correct.
- **`width="20"`**: narrow tall mask window — vertical strip visible. Distinct. Correct.
- **`height="20"`**: short wide mask window — horizontal strip visible. Distinct. Correct.
- **`mask-type`** (`luminance`, `alpha`): both show the masked teal square — size appears similar. The difference between luminance and alpha masking is subtle when the mask is a white/opaque solid shape. The cells look similar but not identical. Allowable.

All cells non-blank. All attribute groups produce visual variation. **PASS.**

---

### `<filter>` — ISSUES

19 value-paths.

- **`filterUnits`** (`userSpaceOnUse`, `objectBoundingBox`): both show blurred/filtered teal rectangles. `userSpaceOnUse` appears to have a slightly different filter region. Correct.
- **`primitiveUnits`** (`userSpaceOnUse`, `objectBoundingBox`): `objectBoundingBox` cell appears **dark/nearly black** — very different from the teal base. This is suspicious. If the filter primitive dimensions are specified in `objectBoundingBox` coordinates (0–1 range) but the stdDeviation or other numeric params remain in user-space units, the result can be extreme (over-blurred into black or clipped). **Flag:** the `primitiveUnits="objectBoundingBox"` cell appears broken (content lost/black). Fix target: verify filter primitive parameters are appropriate for `objectBoundingBox` units and ensure the output shape is still visible.
- **`x`** (`10`, `24px`, `2em`, `50%`, `75%`): filter region shifts — `50%` and `75%` show progressively clipped teal rectangles (smaller visible area). All distinct. Correct.
- **`y`** (`10`, `24px`, `2em`, `50%`, `75%`): same vertical shift pattern. Distinct. Correct.
- **`width="20"`**: narrow visible strip of filtered content. Distinct. Correct.
- **`height="20"`**: short horizontal strip of filtered content. Distinct. Correct.
- **`href="#slot"`**, **`xlink:href="#slot"`**: referencing another filter — teal rectangle rendered with referenced filter. Cells look similar to each other (same filter), which is correct. Correct.
- **`xlink:title="label"`**: title metadata — no visual change. Correct.

**Real issue:**
1. `primitiveUnits="objectBoundingBox"` — cell appears dark/black, suggesting the filter result is invisible or corrupted. Fix target: the filter fixture for this value-path likely has mismatched numeric params (e.g. stdDeviation or coordinates meant for user-space but applied in bounding-box space); adjust fixture to keep content visible under `objectBoundingBox` primitive units.

---

## Consolidated Real Issues

| # | Element | Attribute / Cell | Severity | Description | Fix Target |
|---|---------|-----------------|----------|-------------|------------|
| 1 | `<style>` | `type="label/Aa/sample/specimen"` | Medium | All 4 `type` values suppress CSS (invalid type strings), producing 4 identical blank-styled cells. No valid baseline to compare against. | Add `type="text/css"` value-path so at least one cell shows styled output; contrast with invalid-type cells. |
| 2 | `<style>` | `media="label/Aa/sample/specimen"` | Medium | All 4 `media` values are non-matching, producing 4 identical unstyled cells. | Add `media="all"` or `media="screen"` value-path as a styled baseline. |
| 3 | `<pattern>` | `patternContentUnits="objectBoundingBox"` | Low–Medium | Cell renders triangular/chevron shapes instead of the dot circles used in all other pattern cells. Possibly intentional but visually inconsistent. | Confirm whether this shape difference is deliberate; if not, align content shape across both `patternContentUnits` values. |
| 4 | `<filter>` | `primitiveUnits="objectBoundingBox"` | High | Cell renders dark/nearly black — filter output is lost. Params likely mismatched to bounding-box unit space. | Fix filter fixture for this value-path: scale numeric primitives (stdDeviation etc.) to fractional bounding-box units so filtered output remains visible. |
| 5 | `<clipPath>` | All 3 cells | Note | All cells identical (explainable coincidences). No broken cells, but sheet provides no visual differentiation. | Low-priority: use self-intersecting or offset clip shapes to expose `fill-rule` and `clipPathUnits` differences. |

---

## Summary

| Element | Value-paths | Verdict |
|---------|-------------|---------|
| `<script>` | 3 | OK |
| `<style>` | 12 | ISSUES (2 attribute groups all-identical due to invalid test values) |
| `<linearGradient>` | 28 | OK |
| `<radialGradient>` | 32 | OK |
| `<stop>` | 6 | OK |
| `<pattern>` | 43 | ISSUES (1 cell: unexpected shape for `patternContentUnits="objectBoundingBox"`) |
| `<marker>` | 44 | OK |
| `<clipPath>` | 3 | OK with note (all cells coincidentally identical; technically correct) |
| `<mask>` | 18 | OK |
| `<filter>` | 19 | ISSUES (1 cell: `primitiveUnits="objectBoundingBox"` → black/broken) |
