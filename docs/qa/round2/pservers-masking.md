# Visual QA — Round 2: pservers + masking batch

**Date:** 2026-06-22
**Elements reviewed:** `linearGradient`, `radialGradient`, `stop`, `pattern`, `marker`, `clipPath`, `mask`
**Method:** Fresh full-page screenshots (`SNAP_FULLPAGE=1 SNAP_SCALE=1`) + HTML source inspection.
**Round-1 fixes to confirm:** stop all-black → colour; role token concatenation → spaced; duplicate id → moved to meta.

---

## Round-1 Fix Verification

| Round-1 Issue | Status |
|---|---|
| `stop` gallery all-black | **FIXED** — all cards now show amber/orange gradient |
| `role="rowrowgroup"` / `"combobox complementarycontentinfo"` concatenated | **FIXED** — now `"region row rowgroup"` and `"combobox complementary contentinfo"` (space-separated) |
| Duplicate `id` attribute (`id="slot" id="circle1"`) | **FIXED** — id-variant cards moved to meta (non-visual) section |

---

## `<linearGradient>` — MOSTLY_PERFECT

**80 cards. All main-grid cards show pink-to-teal horizontal gradient. Effect visible and url(#slot) resolves in every main-grid card.**

### Remaining issues

#### 1. `x2="0"` — degenerate zero-length gradient (WEAK_EFFECT)
The card `x2="0"` emits `<linearGradient id="slot" x2="0">` with `gradientUnits="objectBoundingBox"` (default). This produces a zero-length gradient vector `(0%,0%) → (0%,0%)`, rendering as a solid colour identical to the first stop. The card looks the same as `x1="0"` (default vector to x2=1 is correct) and is visually indistinguishable from the degenerate stop-colour result.

**Fix target:** `reps.go` — replace the `"0"` sample for gradient `x2`/`y2` coordinates with a value that keeps the vector non-degenerate (e.g. `"50%"` or `"0.5"`). Alternatively add a guard in `overlay.go` that detects `attrName in {x2, y2}` on gradient elements and avoids returning `"0"`.

#### 2. `href="#refgrad"` — dangling reference (GRAMMAR_ISSUE / WEAK_EFFECT)
The card `href="#refgrad"` emits `<linearGradient id="slot" href="#refgrad">`. The element `#refgrad` is never defined in the SVG, so the inheritance link is a no-op. The gradient still renders using its own stops (the card looks correct), but the `href` attribute is not demonstrating gradient inheritance — just a broken reference.

**Same issue applies to** `xlink:href="#refgrad"` (line 67).

**Fix target:** `overlay.go` — when varying `href`/`xlink:href` on a gradient element, inject a sibling gradient definition `<linearGradient id="refgrad"><stop offset="0" stop-color="#4d8bff"/><stop offset="1" stop-color="#16c79a"/></linearGradient>` into the SVG `<defs>`. This makes the inheritance visually meaningful (the main gradient delegates its stops to the reference).

#### 3. Large NO_EFFECT_IN_MAIN_GRID bloc (GRAMMAR_ISSUE — wrong element target)
The following presentation attributes appear in the main grid but have no effect on `<linearGradient>` itself (they are inherited paint properties, not geometry/server props):

`fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`, `marker-end`, `color`, `color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`, `filter`, `color-interpolation-filters`, `flood-color` (24 attrs)

All 24 show the identical gradient card as their baseline — zero visual distinguishability.

**Note:** These attrs ARE grammar-valid on `linearGradient` per SVG spec (via `PresentationAttribute`), so moving them to meta is technically correct. However they should not be in the main grid if they have no rendering effect on the host element in its scaffold.

**Fix target:** `gallery.go` — add a per-element "no-effect" list for `catGradient` elements (and `catStop`, `catPattern`, `catMask`, `catClip`) that moves inherited-only presentation attrs to the meta section when the scaffold does not give those attrs a visible effect.

---

## `<radialGradient>` — MOSTLY_PERFECT

**84 cards. All cards show a pink-to-teal radial blob. Effect clearly visible in every main-grid card.**

### Remaining issues

#### 1. Same `href="#refgrad"` dangling reference (WEAK_EFFECT)
Identical to linearGradient — the `href` and `xlink:href` cards reference an undefined `#refgrad`.

**Fix target:** Same as linearGradient — inject sibling def in `overlay.go`.

#### 2. Large NO_EFFECT_IN_MAIN_GRID bloc (same as linearGradient)
Same 24 presentation attributes with no visible effect appear in the main grid:

`fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`, `marker-end`, `color`, `color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`, `filter`, `color-interpolation-filters`, `flood-color`

**Fix target:** `gallery.go` — same per-element no-effect list as linearGradient.

---

## `<stop>` — MOSTLY_PERFECT (round-1 fix confirmed working)

**56 enumerated cards. All main-grid cards now show a visible 2-colour gradient (amber/orange on dark blue background). Round-1 all-black issue is resolved.**

### Remaining issues

#### 1. `stop-color` and `stop-opacity` MISSING from main grid (GRAMMAR_ISSUE — missing attrs)
The two most semantically important attributes for `<stop>` — `stop-color` and `stop-opacity` — have **no dedicated cards in the main grid**. The 28 main-grid cards cover `offset` (4 values) and 24 inherited presentation attributes that have no effect on a stop element. `stop-color` and `stop-opacity` are never shown as the labelled varying attribute.

From the HTML source (`stop.html` lines 17–44): the main grid contains `offset`, `fill`, `fill-rule`, ..., `flood-color` — but never `stop-color` or `stop-opacity` as the card label.

From `blueprint.go` line 289: `baselineFor("stop")` correctly injects `stop-color="#f5a623"` as a baseline, but this means when `stop-color` is the _varying_ attribute, it is suppressed by the `add()` function's skip-if-varying logic. The grammar in `pservers.ebnf` defines `StopColorAttr` and `StopOpacityAttr` explicitly — they must appear as main-grid cards.

**Fix target:** `blueprint.go` `baselineFor("stop", varyingPrefix)` — when `varyingPrefix == "stop-color"`, use a different contrasting baseline color for the fixed stop (e.g. `"#4d8bff"`). When `varyingPrefix == "stop-opacity"`, supply `stop-color="#f5a623"` still but let `stop-opacity` vary. The `add()` function correctly drops the baseline when it's the varying attr, but there must be another stop in the template that anchors the gradient. Currently `pservers.ebnf` line 191–192 defines these attrs but the grammar walker must be listing them in the enumeration — check if they are being filtered into the meta section by `gallery.go`.

#### 2. Large NO_EFFECT_IN_MAIN_GRID bloc (24 attrs)
The same 24 inherited presentation attributes that appear for gradient elements appear here in the main grid with no visual effect. For `<stop>`, these attrs (`fill`, `fill-rule`, `stroke-*`, `marker-*`, etc.) are completely ignored by the browser when set on a stop element.

`fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`, `marker-end`, `color`, `color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`, `filter`, `color-interpolation-filters`, `flood-color`

All 24 cards look identical — same amber mid-gradient.

**Fix target:** `gallery.go` — add `stop` to the per-element no-effect list; move these 24 attrs to the meta section with an annotation "inherited presentation attr — no effect on stop element".

#### 3. `offset="0"` card — degenerate first stop (WEAK_EFFECT)
The card `offset="0"` emits:
```html
<stop offset="0" stop-color="#e94560"/>   <!-- baseline anchor -->
<stop stop-color="#f5a623" offset="0"/>   <!-- varied stop also at 0 -->
```
Two stops at offset=0 — browsers take the last one at the same position, so the amber stop (#f5a623) overrides the red (#e94560) at position 0. The gradient is amber-to-teal, which is indistinguishable from the `offset="50%"` baseline card (also amber near the end). The card is not blank, but the offset positioning effect is weakened.

**Fix target:** `blueprint.go` `catStop` scaffold — move the baseline anchor stop to `offset="0.1"` instead of `"0"` so there is always a visible red leading edge even when the varied stop is at `offset="0"`.

---

## `<pattern>` — MOSTLY_PERFECT

**89 cards. All main-grid cards show the orange circle tiling pattern. Effect visible in every card.**

### Remaining issues

#### 1. `patternUnits="objectBoundingBox"` — tile too large (WEAK_EFFECT / different but wrong)
The card `patternUnits="objectBoundingBox"` with baseline `width="20" height="20"` means the tile is 2000% × 2000% of the rect's bounding box. The child circle at `cx="10" cy="10" r="6"` is at 1000% of bounding box — entirely off-canvas. From the screenshot the card shows a large solid amber/orange swatch (the circle fills the visible area at enormous scale), which looks different from other cards but is showing degenerate geometry.

**Fix target:** `blueprint.go` `baselineFor("pattern", varyingPrefix)` — when `varyingPrefix == "patternUnits"`, use fractional width/height (`width="0.2" height="0.2"`) and fractional circle coords (`cx="0.5" cy="0.5" r="0.3"`) for the `objectBoundingBox` variant. This is most cleanly implemented as a template override or a pre-generation variant in the overlay.

#### 2. `patternContentUnits="objectBoundingBox"` — circle off-canvas (WEAK_EFFECT / potentially invisible)
The card `patternContentUnits="objectBoundingBox"` keeps the pattern tile at `width="20" height="20"` (userSpaceOnUse) but the child circle at `cx="10" cy="10"` now means 10× the bounding box — off-canvas. The pattern tile may be empty or show a partial arc.

**Fix target:** Same as patternUnits fix — detect when `patternContentUnits` is the varying attr and switch child coords to fractional `cx="0.5" cy="0.5" r="0.3"`.

#### 3. `href="#slot"` / `xlink:href="#slot"` — self-referential pattern (GRAMMAR_ISSUE)
The card `href="#slot"` emits `<pattern id="slot" href="#slot">` — a pattern referencing itself. Browsers ignore the circular reference and the card renders using the pattern's own children. The card visually looks correct but is a semantic no-op.

**Fix target:** `overlay.go` — when varying `href`/`xlink:href` on pattern elements, use a distinct id `"#refpat"` and inject a sibling `<pattern id="refpat">` definition with different tile content.

#### 4. NO_EFFECT_IN_MAIN_GRID — same 24 inherited presentation attributes
Same as gradient elements: `fill`, `fill-rule`, `fill-opacity`, `stroke-*`, `marker-*`, `paint-order`, `color-*`, `shape-rendering`, `vector-effect`, `filter`, `flood-color` — all placed on the `<pattern>` element itself, where they have no rendering effect (they do not apply to the pattern tile content nor the rect using the pattern).

**Fix target:** `gallery.go` — add `pattern` to per-element no-effect list.

---

## `<marker>` — PERFECT

**84 cards. Arrowhead-on-line effect clearly visible in every card. All position/orientation variants distinguish visibly.**

No new issues in round 2. All round-1 issues (`_top` typo, `markerWidth="2"` small marker, duplicate id) have been resolved or were noted as working-as-intended.

**Category: PERFECT**

---

## `<clipPath>` — MOSTLY_PERFECT

**62 cards. Clipped circle visible in all main-grid cards.**

### Remaining issues

#### 1. `clipPathUnits="objectBoundingBox"` — clip region empty (EMPTY/BLANK)
The card `clipPathUnits="objectBoundingBox"` still emits `<circle cx="50" cy="50" r="35"/>` inside the clip path. With `objectBoundingBox` units, `cx="50"` means 50× the bounding box = far off-canvas. The clip region is empty — the clipped rect is fully invisible (confirmed from screenshot: the card shows only the dark card background).

This was reported in round 1 and is **not fixed**.

**Fix target:** `blueprint.go` `bodyFor("clipPath")` — add a conditional that emits fractional coordinates when `clipPathUnits` is the varying attr. Specifically:
```go
case "clipPath":
    if varyingPrefix == "clipPathUnits" {
        // For objectBoundingBox variant, use fractional coords (0–1 scale)
        return `<circle cx="0.5" cy="0.5" r="0.4"/>`, false
    }
    return `<circle cx="50" cy="50" r="35"/>`, false
```
Alternatively, use a template HTML for clipPath with two pre-set SVG variants.

#### 2. NO_EFFECT_IN_MAIN_GRID — 24 inherited presentation attributes
The same bloc of `fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`, `marker-end`, `color`, `color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`, `filter`, `color-interpolation-filters`, `flood-color` appear in the main grid placed on the `<clipPath>` element itself. These have no effect on clip region geometry.

**Fix target:** `gallery.go` — add `clipPath` to per-element no-effect list.

---

## `<mask>` — MOSTLY_PERFECT

**75 cards. Masked square visible in most main-grid cards.**

### Remaining issues

#### 1. `maskContentUnits="objectBoundingBox"` — mask content off-canvas (EMPTY/BLANK)
The card `maskContentUnits="objectBoundingBox"` still emits `<rect x="20" y="20" width="60" height="60" fill="#fff"/>` with `objectBoundingBox` content units. `x="20"` = 2000% of bounding box — the white mask rect is entirely off-canvas. The mask is transparent, making the masked element invisible.

This was reported in round 1 and is **not fixed**.

**Fix target:** `blueprint.go` `bodyFor("mask")` — add conditional:
```go
case "mask":
    if varyingPrefix == "maskContentUnits" {
        // For objectBoundingBox variant
        return `<rect x="0.1" y="0.1" width="0.8" height="0.8" fill="#fff"/>`, false
    }
    return `<rect x="20" y="20" width="60" height="60" fill="#fff"/>`, false
```

#### 2. `x="75%"`, `y="75%"` — mask window clipped to corner (WEAK_EFFECT)
With `maskUnits="userSpaceOnUse"` (default), `x="75%"` is 75% of the SVG viewport = 90 user units. The mask window starting at x=90 leaves only a 10-unit strip of the masked rect visible on the right edge. The card shows a small sliver of pink at the edge — technically visible but very faint. `y="75%"` similarly shows a bottom sliver.

These are not blank but very hard to see. Since `x="100%"` (round-1 issue) was removed, the `75%` cards are the worst remaining visibility cases.

**Fix target:** `reps.go` — replace the `"75%"` sample in `LengthPercentageType` with `"60%"` or `"40%"` for mask/pattern geometry attributes, keeping the shifted-window effect visible. Or add a per-attr clamp in `overlay.go` for mask `x`/`y`.

#### 3. `width="20"`, `height="20"` — visually correct (note only)
Mask window 20×20 px shows a small pink square — effect is intentional and distinguishable. No action needed.

#### 4. NO_EFFECT_IN_MAIN_GRID — 24 inherited presentation attributes
Same bloc as all other elements above: `fill`, `fill-rule`, `fill-opacity`, `stroke-*`, `marker-*`, etc. placed on `<mask>` element itself with no rendering effect.

**Fix target:** `gallery.go` — add `mask` to per-element no-effect list.

---

## Cross-cutting issues summary (round 2)

| Issue | Affected Elements | Severity | Fix Target |
|---|---|---|---|
| **`stop-color`/`stop-opacity` missing from main grid** | `stop` | HIGH | `blueprint.go` `baselineFor("stop")` + `gallery.go` |
| **`objectBoundingBox` child coords wrong** (empty/invisible) | `clipPath`, `mask`, `pattern` | HIGH | `blueprint.go` `bodyFor("clipPath"/"mask"/"pattern")` — fractional coords when obb variant |
| **NO_EFFECT_IN_MAIN_GRID: 24 inherited presentation attrs** | All 7 elements | MEDIUM | `gallery.go` — per-element no-effect allowlist; move to meta section |
| **`href`/`xlink:href` dangling/self-referential** | `linearGradient`, `radialGradient`, `pattern` | LOW | `overlay.go` — inject sibling reference def |
| **`x2="0"` / `y2="0"` degenerate gradient** | `linearGradient` | LOW | `reps.go` — replace `"0"` coord sample with `"50%"` |
| **`offset="0"` doubled stop position** | `stop` | LOW | `blueprint.go` `catStop` scaffold — anchor stop at `offset="0.1"` |
| **`x="75%"` / `y="75%"` near-edge mask** | `mask` | LOW | `reps.go` / `overlay.go` — clamp `LengthPercentageType` to ≤ 60% on geometry attrs |

---

## Per-element one-liner summary

| Element | Category | One-liner |
|---|---|---|
| `linearGradient` | MOSTLY_PERFECT | 24 no-effect attrs cluttering main grid; href dangling; x2=0 degenerate |
| `radialGradient` | MOSTLY_PERFECT | Same 24 no-effect attrs; href dangling |
| `stop` | MOSTLY_PERFECT | stop-color/stop-opacity MISSING from main grid; 24 no-effect attrs in main grid |
| `pattern` | MOSTLY_PERFECT | objectBoundingBox coords wrong; href self-ref; 24 no-effect attrs |
| `marker` | PERFECT | No remaining issues |
| `clipPath` | MOSTLY_PERFECT | clipPathUnits=objectBoundingBox still blank; 24 no-effect attrs |
| `mask` | MOSTLY_PERFECT | maskContentUnits=objectBoundingBox still blank; x/y=75% near-invisible; 24 no-effect attrs |

---

## Top remaining systematic issues

1. **NO_EFFECT_IN_MAIN_GRID bloc (24 attrs × 6 elements)** — The same 24 inherited presentation attributes appear in the main grid for every paint-server and masking element, with no visual effect on those elements. This is the single largest quality gap: 144 wasted cards across 6 elements. Fix: `gallery.go` per-element no-effect list → move to meta.

2. **`stop-color` and `stop-opacity` missing** — The most important semantic attributes of `<stop>` have no dedicated main-grid cards. They are defined in `pservers.ebnf` and referenced in `PresentationAttribute` but get suppressed or fall into the no-effect pile. Fix: `blueprint.go` `baselineFor("stop")` must handle the stop-color and stop-opacity varying cases with appropriate contrasting anchors.

3. **`objectBoundingBox` content blank** — `clipPath` and `mask` galleries each have one blank card where the child geometry is in absolute coords but `objectBoundingBox` units apply. Fix: `blueprint.go` `bodyFor("clipPath"/"mask")` with unit-aware coord branching.
