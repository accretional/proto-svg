# Round-2 Visual QA — Filters batch 2

Elements: feBlend, feComposite, feMerge, feMergeNode, feMorphology, feOffset,
feFlood, feDropShadow, feImage, feTile, feTurbulence, feConvolveMatrix, feDisplacementMap

Screenshots: `chrome-testing/screenshots/review/<tag>.png`
Date: 2026-06-22

---

## Per-element status

### feBlend — ISSUES: NO_EFFECT_IN_MAIN_GRID, GRAMMAR_ISSUE

Round-1 fix confirmed: in/in2 are now wired; two-layer blend (orange + blue feFlood)
renders clearly. The mode sweep (normal/multiply/screen/darken/lighten/overlay/color-dodge/
color-burn/hard-light/soft-light/difference/exclusion/hue/saturation/color/luminosity) is
fully visible and distinct.

Remaining issues:

1. **NO_EFFECT_IN_MAIN_GRID** — 16 presentation attrs appear in the main grid and render
   identically to the in/in2/mode cards (filter primitives ignore fill/stroke/marker at
   render time). Affected cards (all in main grid):
   `fill="none"`, `fill-rule="nonzero"`, `fill-opacity="0.5"`,
   `stroke="none"`, `stroke-opacity="0.5"`, `stroke-width="20"`,
   `stroke-linecap="butt"`, `stroke-linejoin="miter"`, `stroke-miterlimit="4"`,
   `stroke-dasharray="none"`, `stroke-dashoffset="20"`, `paint-order="normal"`,
   `marker="none"`, `marker-start="none"`, `marker-mid="none"`, `marker-end="none"`.
   (`color-rendering`, `shape-rendering`, `vector-effect` and `color-interpolation` also
   appear but are at least filter-relevant hint attrs and less egregious.)

2. **GRAMMAR_ISSUE** — `no-composite="true|false"` (FeBlendNoCompositeAttr in
   `lang/filter.ebnf:165`) is not a standard SVG Filter Effects attribute on `<feBlend>`.
   It does not exist in the SVG 2 spec or Chrome's implementation. Two identical blue
   cards with no effect appear for `no-composite="true"` and `no-composite="false"`.

Fix target:
- `emit.go` / `nonVisualAttr()`: extend the function to also classify
  `fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`,
  `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`,
  `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`,
  `marker-end` as non-visual **when the host tag is a filter primitive** (prefix `fe`).
  These belong in the dimmed details section, not the main grid.
  Signature change: `nonVisualAttr(attr, tag string) bool` — pass tag through from
  `emitPage`.
- `lang/filter.ebnf`: remove `FeBlendNoCompositeAttr` from `FeBlendAttribute` union
  (line 144) and delete lines 164-165. It is not in the SVG spec.

---

### feComposite — ISSUES: NO_EFFECT_IN_MAIN_GRID, WEAK_EFFECT

Round-1 fix confirmed: two-layer scaffold (orange + blue feFlood with offset positions)
renders correctly; operator sweep (over/in/out/atop/xor/lighter/arithmetic) and k1/k2/k3/k4
are visible. Blank issue resolved.

Remaining issues:

1. **NO_EFFECT_IN_MAIN_GRID** — same presentation attr flood as feBlend: 18 cards
   (`fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`,
   `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`,
   `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`,
   `marker-end`, `color-rendering`, `shape-rendering`, `vector-effect`).

2. **WEAK_EFFECT** — The `in="SourceGraphic"` / `in="SourceAlpha"` / `in2=*` cards
   composite a white rect against the scaffold's two colored feFlood layers. Because
   `in="SourceGraphic"` puts the plain white rect as input, many operator variations
   produce mostly-white output with subtle difference. Using the rect with a
   multi-color fill or a pre-defined pattern as the source would make operator
   differences more legible.

Fix target: `emit.go` `nonVisualAttr` (tag-aware, same as feBlend fix).

---

### feMerge — ISSUES: NO_EFFECT_IN_MAIN_GRID

Round-1 confirmed working. The two-layer (red/blue) stacked feMergeNode rendering is
clear across all x/y/width/height/result variants.

Remaining issues:

1. **NO_EFFECT_IN_MAIN_GRID** — 16 presentation attr cards are in the main grid
   (`fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`,
   `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`,
   `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`,
   `marker-end`). Note: `feMerge` has a `<feMergeNode>` scaffold that actually
   provides two named inputs (layerA/layerB), so the merge rendering is already
   visible; presentation attrs are the only grid polluters.

Fix target: `emit.go` `nonVisualAttr` (tag-aware filter-primitive rule).

---

### feMergeNode — PERFECT

All 57 cards render a clear pink/coral-filled rectangle (SourceGraphic passthrough
with layerA/layerB scaffold). The `in="*"` sweep, fill/stroke/marker presentation
attrs, color-interpolation, and result attrs all appear correctly. Presentation attrs
are already in the main grid here as well but the scaffold renders visibly (the
feMergeNode is a child, not the filter output itself, so its presentation attrs are
harmless). However they still technically flood the grid with near-identical cards —
same fix opportunity as siblings.

Minor note: the feMergeNode scaffold's `<feMerge>` wraps `{{ELEMENT}}` so a single
`<feMergeNode>` is the only layer; all cards show only layerA color because layerB
node is absent. This is correct for isolating the feMergeNode attribute but means
the "stacked layers" effect is not visible here (as expected for a node test).

Status: **PERFECT** (no blocking issues; optional cleanup: move presentation attrs
to non-visual section).

---

### feMorphology — PERFECT

All 77 cards render the "SVG" text logo through feMorphology correctly. The erode/dilate
operator swap is clearly visible (dilate fattens strokes, erode thins them). The radius
sweep and all spatial attrs (x/y/width/height) are legible. Presentation attrs appear
but the text source is complex enough that the cards are visually distinguishable even
with them present (the morphology is applied uniformly). The `in=*` sweep shows the
correct source-switching behavior.

Status: **PERFECT**

---

### feOffset — ISSUES: WEAK_EFFECT

All 86 cards render an orange square — the offset primitive correctly shifts the
SourceGraphic. However the source is a uniform flat-color rect, so the offset is
only visible as a clipped edge artifact at dx>0 or dy>0. Cards for small dx/dy
values (0, 0.5, 1) are visually identical to the baseline card.

1. **WEAK_EFFECT** — uniform orange fill makes small offset values invisible. The
   source should be a non-uniform image (e.g. a multi-color gradient or the SVG
   logo) so the offset shift is perceptible for any delta.

Status: issues noted but output is not blank. Functional but weak demonstration.

Fix target: `blueprint.go` `defaultScaffold` for `catFilterPrimitive` — or add a
template `chrome-testing/html/template/feOffset.html` with a multi-color source rect
(e.g. gradient fill or two-tone rect) so offset is distinguishable.

---

### feFlood — ISSUES: NO_EFFECT_IN_MAIN_GRID

The feFlood page shows a pink/coral square for all cards — feFlood outputs a solid
color rectangle regardless of source. The x/y/width/height geometry attrs are visible
as positional clipping of the colored rectangle within the card. The `flood-color` and
`flood-opacity` attrs are not in this batch's main attrs but `feFlood` correctly
produces color-filled output.

Remaining issues:

1. **NO_EFFECT_IN_MAIN_GRID** — 16 presentation attr cards (`fill`, `fill-rule`,
   `fill-opacity`, `stroke`, `stroke-*`, `paint-order`, `marker-*`) flood the main grid.
   They are all visually identical pink rectangles.

Note: `feFlood` has no `in` attribute — its output is purely generated color, not
source-dependent. The x/y offset cards correctly show the flood region moved, which
is good. The geometry attrs (x, y, width, height) do demonstrate distinct visual
behavior (varying sub-region).

Fix target: `emit.go` `nonVisualAttr` (tag-aware).

---

### feDropShadow — ISSUES: WEAK_EFFECT

All 87 cards render an orange square. The `dx`/`dy` shadow offset is not visible
because the shadow falls outside the filter region clipping (or is the same color as
the card background). The `stdDeviation` card also shows no visible blur halo.

1. **WEAK_EFFECT** — The shadow is not visible. Root causes:
   a. The filter region is set to `x="-20%" y="-20%" width="150%" height="150%"` in
      the scaffold (adequate room for shadow bleed), but the default `flood-color` for
      feDropShadow is `black` and the card background is also dark (#161c3a), so the
      shadow blends into the background and is invisible.
   b. The default `dx` and `dy` are both 2 (SVG spec default), which is a very small
      offset at the 100x100 viewBox scale.

Fix target: `blueprint.go` — for `catFilterPrimitive` feDropShadow (or a template),
set the source shape fill to white or light color, and/or set `flood-color="#e94560"`
on the feDropShadow element as a baseline attribute (the generator already has
`baselineFor` logic). Alternatively, the template should use a light background color
for the SVG (`background: #fff` on `.card svg`) so the black shadow is visible.

Since the scaffold is generic (shared by all filter primitives), the fix is a dedicated
template for feDropShadow:
`chrome-testing/html/template/feDropShadow.html` — use a white source rect on a
light SVG background with an explicit shadow-color that contrasts.

---

### feImage — ISSUES: GRAMMAR_ISSUE (href self-reference)

Round-1 fix confirmed: most cards render the embedded SVG (dark square + teal circle)
correctly through the feImage filter. Cards for `preserveAspectRatio`, `crossorigin`,
x/y/width/height, and result attrs all show the teal-circle image.

Remaining issues:

1. **GRAMMAR_ISSUE / EMPTY** — The very first card (`href="#slot"`) has
   `<feImage id="slot" href="#slot">` — a circular self-reference. The feImage
   element has id="slot" (injected by the blueprint slot mechanism) and its href also
   resolves to "#slot" (the overlay's `IriType` fallback). Chrome renders this as
   blank (circular reference is ignored).

   Fix target: `overlay.go` `overlaySample()` — add a case for `feImage` tag with
   `href`/`xlink:href` attribute name to return the `imageDataURI` constant
   (already defined in `blueprint.go`) instead of `"#slot"`. This mirrors the
   special-case handling already done for gradient `href` → `"#refgrad"`:
   ```go
   // feImage href: a self-reference (#slot) produces a circular load.
   // Use the inline data URI so the image always renders.
   if tag == "feImage" && (an == "href" || an == "xlink:href") {
       return imageDataURI, true
   }
   ```

2. **GRAMMAR_ISSUE** — The last card in the page (`xlink:href="#slot"`) repeats the
   same circular reference problem as `href="#slot"`. Both need the same fix.

Note: `crossorigin="use-credentials"` cards are blank (expected — no CORS context
for a data URI). This is acceptable behavior, not a bug.

---

### feTile — ISSUES: WEAK_EFFECT, NO_EFFECT_IN_MAIN_GRID

Cards for `in="patch"` (x/y/width/height variants) correctly tile the 25%-of-filter
patch — the small orange square is tiled repeatedly to fill the card. This looks good.

Remaining issues:

1. **WEAK_EFFECT** — Cards for `in="SourceGraphic"` (and `in="SourceAlpha"`,
   `in="result1"`, etc.) tile the full source rect, which is a solid white fill. Tiling
   a uniform color produces a uniform color — no tiling pattern is visible. The patch
   input (`result="patch"` — a small 25% orange feFlood) is defined in the blueprint
   scaffold but the `in=*` cards use other sources that don't show tiling.

   Fix target: `blueprint.go` default `catFilterPrimitive` scaffold (or a template
   for feTile). For feTile specifically, the `in` baseline should always wire to
   `in="patch"` so the tiling effect is always visible. The generator currently picks
   the varied value for `in`, so when `in="SourceGraphic"`, the result is a uniform
   tile. One approach: the feTile template should use a multi-color source as
   SourceGraphic (e.g. a gradient rect) so tiling shows differentiation.

2. **NO_EFFECT_IN_MAIN_GRID** — 17 presentation attr cards in main grid (`fill`,
   `fill-rule`, `fill-opacity`, `stroke-*`, `paint-order`, `marker-*`).

---

### feTurbulence — ISSUES: WEAK_EFFECT (dark cards)

Round-1 fix confirmed: `numOctaves` is now 3 (positive int, valid) and `order` fix
not applicable here. The `baseFrequency` sweep cards (0.5, 2, 3, 2 3, 3 3, 1 2) all
render visible noise patterns. The `type="fractalNoise"` / `type="turbulence"` and
`stitchTiles` cards render correctly.

Remaining issues:

1. **WEAK_EFFECT (dark cards)** — All cards where `baseFrequency` is NOT explicitly
   set (i.e. `numOctaves`, `seed`, `stitchTiles`, `type`, `x`, `y`, `width`, `height`,
   `result`, and all presentation attr cards) render as near-black rectangles. This is
   because feTurbulence defaults to `baseFrequency="0"` when not set, which produces
   a uniform constant-color (black) output.

   Fix target: `overlay.go` `overlaySample()` — add a baseline for `baseFrequency`
   on `feTurbulence`:
   ```go
   case "basefrequency":
       if tag == "feTurbulence" {
           return "0.05", true
       }
   ```
   This ensures every card has a visible turbulence pattern unless `baseFrequency` is
   the varied attribute itself. The value 0.05 produces clearly visible noise at the
   100px viewBox scale.

2. **NO_EFFECT_IN_MAIN_GRID** — Presentation attr cards (`fill`, `stroke-*`,
   `paint-order`, `marker-*`) all render as dark/near-black (compounding the
   baseFrequency issue) in the main grid.

---

### feConvolveMatrix — ISSUES: WEAK_EFFECT

Round-1 fix confirmed: `order="3"` (positive int, valid) is used; all cards render
an orange rectangle. The `kernelMatrix` baseline uses a sharpening kernel
`"0 -1 0 -1 5 -1 0 -1 0"` which should produce visible edge enhancement.

Remaining issues:

1. **WEAK_EFFECT** — The source is a uniform flat-color orange rect. A sharpening
   convolution on a uniform color produces... the same uniform color (no edges to
   enhance). The kernel effect is completely invisible. The `edgeMode`, `divisor`,
   `bias`, `targetX`, `targetY`, `kernelUnitLength`, and `preserveAlpha` cards all
   show the same flat orange rectangle.

   Fix target: Use a textured source for feConvolveMatrix cards. Options:
   - `blueprint.go`: for feConvolveMatrix specifically, make the source rect use a
     gradient or SVG text so edges exist for the convolution to operate on. A template
     `chrome-testing/html/template/feConvolveMatrix.html` with a multi-color striped
     or gradient rect would make all kernel variants visible.
   - Or: the generator baseline for feConvolveMatrix should wire `in="noiseMap"`
     (the feTurbulence result already in the scaffold) so the convolution operates
     on a textured input.

2. The `order=*` cards (3, and the `kernelMatrix` variants for different shapes)
   appear as flat orange rectangles for the same reason.

---

### feDisplacementMap — PERFECT

All 96 cards render a clearly distorted orange rectangle (wavy/torn edges). The
`noiseMap` feTurbulence source (pre-defined in the scaffold) provides a rich
displacement map. The `in`/`in2` sweep, `scale` sweep, `xChannelSelector` /
`yChannelSelector` sweep all produce visibly distinct distortions. The round-1
`scale` validity fix is confirmed (negative and fractional scales render correctly
— they displace in reverse/small amounts).

Status: **PERFECT**

---

## Systematic issues summary (priority order)

### P1 — NO_EFFECT_IN_MAIN_GRID: presentation attrs flood filter primitive grids

**Affects:** feBlend, feComposite, feMerge, feFlood, feTile, feTurbulence,
feDisplacementMap, and (to a lesser extent) feMorphology, feOffset, feDropShadow,
feImage, feConvolveMatrix (all `fe*` elements via `PresentationAttribute` in grammar).

**Root cause:** `nonVisualAttr()` in `emit.go` only classifies a fixed set of
accessibility/scripting attrs as non-visual. It is not tag-aware, so presentation
attrs (`fill`, `stroke-*`, `marker-*`, `paint-order`, `vector-effect`, `shape-rendering`)
on filter primitive elements go into the main grid even though they have zero visual
effect on filter output.

**Fix:** Make `nonVisualAttr(attr, tag string) bool` tag-aware. Add a branch:
```go
// Presentation attrs on filter primitives have no effect on filter output.
if strings.HasPrefix(tag, "fe") {
    switch attr {
    case "fill", "fill-rule", "fill-opacity",
         "stroke", "stroke-opacity", "stroke-width", "stroke-linecap",
         "stroke-linejoin", "stroke-miterlimit", "stroke-dasharray",
         "stroke-dashoffset", "paint-order",
         "marker", "marker-start", "marker-mid", "marker-end",
         "shape-rendering", "vector-effect":
        return true
    }
}
```
`color-interpolation`, `color-interpolation-filters`, and `color-rendering` may stay in
the main grid for filter primitives as they genuinely affect filter computation.
Also update all callers: `emitPage` must pass `v.Attr, p.tag` to `nonVisualAttr`.

**Target file:** `chrome-testing/cmd/gen/emit.go`

---

### P2 — WEAK_EFFECT: uniform flat-color sources hide filter effects

**Affects:** feOffset, feConvolveMatrix, feDropShadow, feTile (in="SourceGraphic" cards)

**Root cause:** The default `catFilterPrimitive` scaffold uses a flat `fill="#e94560"` rect
as the source shape. Filters that require spatial variation or edge data (offset,
convolution) or contrast against background (drop shadow) produce no visible change
on uniform inputs.

**Fix (per element):**
- **feOffset**: use a gradient rect or two-color striped rect as source so offset is
  perceivable for any dx/dy.
- **feConvolveMatrix**: wire `in="noiseMap"` (already in scaffold) as the default
  source input so the convolution kernel operates on textured data.
- **feDropShadow**: use a white source rect (`fill="white"`) and a light SVG
  background (`background: #eee` on the card's SVG) so the default black shadow is
  visible. Or set `flood-color` to a saturated color in the baseline.
- **feTile**: source for `in="SourceGraphic"` cards should be a small non-uniform
  pattern (the `patch` feFlood is already in the scaffold; tie the baseline `in` to
  `"patch"` for the non-in-varying cards).

**Target file:** `chrome-testing/cmd/gen/blueprint.go` (add templates or modify
`defaultScaffold` / `baselineFor` for feConvolveMatrix, feDropShadow), plus possibly
new template HTML files for feOffset and feTile.

---

### P3 — WEAK_EFFECT: feTurbulence defaults to baseFrequency=0 (all-dark cards)

**Affects:** feTurbulence — all non-baseFrequency attribute cards render near-black.

**Fix:** Add overlay case in `overlay.go`:
```go
case "basefrequency":
    if tag == "feTurbulence" {
        return "0.05", true
    }
```

**Target file:** `chrome-testing/cmd/gen/overlay.go`

---

### P4 — GRAMMAR_ISSUE: no-composite on feBlend

`no-composite` is not a real SVG spec attribute. Remove from `lang/filter.ebnf` (lines 144, 164-165).

**Target file:** `lang/filter.ebnf`

---

### P5 — GRAMMAR_ISSUE / EMPTY: feImage href="#slot" circular self-reference

First and last feImage cards are blank. Fix: add overlay special-case for feImage
href/xlink:href → imageDataURI.

**Target file:** `chrome-testing/cmd/gen/overlay.go`

---

## Elements confirmed PERFECT this round

| Element | Status | Notes |
|---|---|---|
| feMorphology | PERFECT | erode/dilate clearly visible on SVG text |
| feDisplacementMap | PERFECT | rich distortion with noiseMap source |
| feMergeNode | PERFECT | single-node layerA passthrough renders cleanly |

## Elements with fixes needed

| Element | Issues | Priority fix |
|---|---|---|
| feBlend | NO_EFFECT_IN_MAIN_GRID, GRAMMAR_ISSUE | emit.go nonVisualAttr + remove no-composite |
| feComposite | NO_EFFECT_IN_MAIN_GRID, WEAK_EFFECT | emit.go nonVisualAttr |
| feMerge | NO_EFFECT_IN_MAIN_GRID | emit.go nonVisualAttr |
| feFlood | NO_EFFECT_IN_MAIN_GRID | emit.go nonVisualAttr |
| feTile | WEAK_EFFECT, NO_EFFECT_IN_MAIN_GRID | blueprint/template + emit.go |
| feOffset | WEAK_EFFECT | blueprint/template |
| feDropShadow | WEAK_EFFECT | blueprint/template |
| feImage | GRAMMAR_ISSUE (2 blank cards) | overlay.go feImage href → dataURI |
| feTurbulence | WEAK_EFFECT (dark non-baseFreq cards) | overlay.go baseFrequency baseline |
| feConvolveMatrix | WEAK_EFFECT | blueprint/template (use noiseMap as source) |
