# Round-3 QA: masking + paint-server galleries

Screenshots taken: 2026-06-22, SNAP_FULLPAGE=1 SNAP_SCALE=1, fresh rebuild.
Batch: linearGradient, radialGradient, stop, pattern, marker, clipPath, mask.

---

## Per-element verdicts

### linearGradient — ISSUES

Main grid renders 28 cards. Visual issues found:

| Card label | Issue | Severity |
|---|---|---|
| `x2="0"` | WEAK_EFFECT — x1=0 (default) + x2=0 collapses gradient to zero-length vector; Chrome renders the last-stop solid color (green). Communicates nothing about x2. | WEAK |
| `href="#refgrad"` | DANGLING reference — the SVG has only `id="slot"`; there is no `<linearGradient id="refgrad">` in the defs. The card still renders a gradient because the inline stops on the `slot` gradient work, but the href inheritance does not actually resolve. | DANGLING |
| `xlink:href="#refgrad"` | Same dangling reference as above. | DANGLING |

All other cards (gradientUnits, gradientTransform, x1/y1/x2/y2 non-zero, spreadMethod, xlink:title) are PERFECT — two-color gradient visible, distinct from one another.

---

### radialGradient — PERFECT

All 35 main-grid cards render with a clear two-tone radial bloom. The href="#refgrad" / xlink:href="#refgrad" cards show the gradient (stops on `slot` itself still work), and while the refgrad inheritance is technically dangling (same structural issue as linearGradient), the visual result is non-blank and communicates the href attribute. No further action required for this round unless reference resolution is being audited separately.

---

### stop — ISSUES

Main grid renders 6 cards. Visual issues found:

| Card label | Issue | Severity |
|---|---|---|
| `offset="0"` | WEAK_EFFECT — the varied stop (`offset="0"`, `stop-color="#f5a623"`) collides with the leading anchor stop also at offset=0 (from catStop scaffold: `<stop offset="0.1" stop-color="#e94560"/>`). Wait — the scaffold anchor is at 0.1, not 0. But the orange stop is at position 0 and the red anchor at 0.1: result is a sharp orange→red→blue band. The card actually shows a clean solid orange fill (all orange) because the orange stop at 0 and the red at 0.1 make a nearly-invisible transition at the very top, and the overall fill looks solid orange. This is cosmetically weakly differentiated from `offset="50%"`. Medium concern. | WEAK |
| `stop-color="currentColor"` | WEAK_EFFECT — `currentColor` inherits from the SVG's rendering context. Since the body text color in the card CSS is `#e6e6e6` (near-white), currentColor renders as near-white, which is almost invisible against the light-orange gradient blend. The card shows solid orange rather than a distinctive two-tone band. | WEAK |

`offset="1"`, `offset="50%"`, `offset="100%"`, `stop-opacity="0.5"` are PERFECT.

---

### pattern — PERFECT

All 33 main-grid cards render with the orange-dot-on-dark polka-dot pattern. Every attribute variant (patternUnits, patternContentUnits, patternTransform, x/y/width/height in multiple units, viewBox, preserveAspectRatio, href/xlink:href) produces a clearly visible and distinct tiling. PERFECT.

---

### marker — PERFECT

All 35 main-grid cards render with the teal line + pink arrowhead. Every attribute variant (markerUnits, markerWidth, markerHeight, refX/refY, orient, viewBox, preserveAspectRatio) produces a clearly visible, distinct marker. PERFECT.

---

### clipPath — ISSUES

Main grid renders only 2 cards (all other attributes are correctly in the meta section).

| Card label | Issue | Severity |
|---|---|---|
| `clipPathUnits="objectBoundingBox"` | MISLEADING — the clip child is `<circle cx="50" cy="50" r="35"/>` with absolute userSpaceOnUse pixel values. When `clipPathUnits="objectBoundingBox"`, SVG interprets those as fractions of the bounding box: cx=50 → 4500px, r=35 → 3150px. The circle is astronomically large and covers the entire bounding box, so the clipped rect appears as a full solid-pink rectangle — visually identical to the unclipped case. A viewer cannot distinguish this card from no-clip. | MISLEADING |

`clipPathUnits="userSpaceOnUse"` is PERFECT (pink circle visible).

---

### mask — ISSUES

Main grid renders 12 cards. Visual issues found:

| Card label | Issue | Severity |
|---|---|---|
| `maskContentUnits="objectBoundingBox"` | MISLEADING — the mask content is `<rect x="20" y="20" width="60" height="60" fill="#fff"/>`. In objectBoundingBox mode those coordinates are fractions: x=20 → 20×90=1800px offset, far outside the shape. The white mask rect should clip everything away (blank), but Chrome appears to show a normally-masked square (pink square visible). This means either Chrome clamps OBB coords or the blueprint white rect gets ignored. Either way the card is indistinguishable from the `maskContentUnits="userSpaceOnUse"` card — no differentiation. | MISLEADING |
| `maskUnits="objectBoundingBox"` | The mask viewport box (x/y/width/height of the mask element) uses OBB fractions. The default mask x/y/width/height in OBB mode is (-10%,-10%,120%,120%), so this should look the same as userSpaceOnUse at default values — and indeed both cards look identical. Fine as a structural card but produces no visible difference from the default. | MINOR |

All other mask cards (x, y, width, height in various units, mask-type luminance/alpha) are PERFECT — the pink square is clearly offset/sized/clipped differently across cards.

---

## Consolidated remaining-issue list

### 1. DANGLING — `linearGradient href="#refgrad"` and `xlink:href="#refgrad"`

**Root cause:** The blueprint for `catGradient` (blueprint.go line ~198-200) emits the gradient element itself as the only def. `overlay.go` correctly routes gradient href/xlink:href to `"#refgrad"` (line 56), but no `<linearGradient id="refgrad">` is injected into the SVG's defs.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go`, `defaultScaffold()`, `catGradient` case.

**Concrete change:** Extend the catGradient scaffold to include a sibling refgrad definition in the defs so the href resolves:

```go
case catGradient:
    return svgOpen +
        `<defs>` +
        `<linearGradient id="refgrad"><stop offset="0" stop-color="#e94560"/><stop offset="1" stop-color="#16c79a"/></linearGradient>` +
        `{{ELEMENT}}` +
        `</defs>` +
        `<rect x="5" y="5" width="90" height="90" fill="url(#slot)"/></svg>`
```

This wraps the element in `<defs>` (already valid for gradients) and provides a concrete `#refgrad` for the href to inherit from. The `fill="url(#slot)"` still references the generated gradient, and the `#refgrad` provides a non-circular inheritance target.

Note: radialGradient has the same structural gap but the visual result is non-blank; the fix above applies to both because the scaffold is shared via `catGradient`.

---

### 2. WEAK_EFFECT — `linearGradient x2="0"`

**Root cause:** When `x2="0"` is varied and x1 takes the default value of "0" (objectBoundingBox default = 0%), the gradient vector has zero length. Chrome renders the last stop color as a solid fill.

**Fix target:** `chrome-testing/cmd/gen/overlay.go`, `overlaySample()`.

**Concrete change:** Add a guard for the x2 attribute on gradient tags — when the varied attr is `x2`, also set `x1` to a non-zero baseline (e.g. via `baselineFor`) OR override the sampled value "0" to "0" with a visible companion. The cleanest fix is to keep `x2="0"` as a valid grammar value but ensure the baseline sets `x1="1"` (or similar non-equal endpoint) so the card at least shows a one-sided gradient. Since `baselineFor` already skips the varying attribute, add an overlay rule:

```go
case "x2":
    if isGradientTag(tag) && value == "0" {
        // zero-length when default x1=0; keep value but note in baselineFor
        // that x1 must be non-zero baseline
    }
```

Alternatively, adjust `baselineFor` for gradient tags: when varying `x2`, set baseline `x1="1"` so the gradient spans from right to the varied x2 position (including 0 = inverted). This requires passing the varying attribute to `baselineFor`'s pairs — it already skips pairs matching `varying`, so simply add `{"x1", "1"}` to gradient baselines and let it be skipped only when x1 itself is being varied.

Simplest one-line fix in `blueprint.go` `baselineFor` for linearGradient: add `{"x1", "1"}` to the gradient baseline pairs alongside the existing stops.

---

### 3. WEAK_EFFECT — `stop offset="0"` and `stop-color="currentColor"`

**Root cause (offset="0"):** The catStop scaffold places the leading anchor at `offset="0.1"`. A varied stop at offset="0" sits at the very start, and since the orange stop and the red anchor are close together (0 vs 0.1), the overall appearance is nearly all-orange. The differentiation from other offset values is weak.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go`, `defaultScaffold()`, `catStop` case.

**Concrete change:** Move the leading anchor stop further in — e.g. `offset="0.3"` — so that a varied stop at offset=0 is visually separated by a larger red band before the blue anchor at offset=1:

```go
case catStop:
    return svgOpen +
        `<defs><linearGradient id="slot"><stop offset="0.3" stop-color="#e94560"/>{{ELEMENT}}<stop offset="1" stop-color="#4d8bff"/></linearGradient></defs>` +
        `<rect x="5" y="5" width="90" height="90" fill="url(#slot)"/></svg>`
```

**Root cause (stop-color="currentColor"):** `currentColor` inherits from the SVG font color cascade. In a browser SVG, the default color is typically black or near-black, so the stop appears near-black or invisible against the dark background gradient.

**Fix target:** `chrome-testing/cmd/gen/overlay.go`, `overlaySample()`.

**Concrete change:** Add a `stop-color` override so `currentColor` is not emitted for the main (visual) stop cards, OR set an explicit `color` attribute on the parent SVG element in the scaffold to a vivid value. The overlay approach is cleaner:

```go
case "stop-color":
    if tag == "stop" {
        return "#f5a623", true  // orange; demonstrates stop-color takes a color value
    }
```

However, this removes the `currentColor` demonstration entirely. A better approach: keep `currentColor` in the grammar output but add `color="#f5a623"` to the scaffold SVG's root element so `currentColor` resolves to orange, making it visible. Modify `catStop` scaffold's `svgOpen` to include `color="#f5a623"` on the SVG root.

---

### 4. MISLEADING — `clipPathUnits="objectBoundingBox"` absolute child coords

**Root cause:** The `catClip` scaffold uses `<circle cx="50" cy="50" r="35"/>` as the clip child. In `clipPathUnits="objectBoundingBox"` mode the coordinates are fractions of the referenced element's bounding box, not user-space pixels. cx=50 → 50×90=4500px overflows and covers everything, making the clip visually equivalent to no-clip.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go`, `defaultScaffold()`, `catClip` case, OR `overlay.go`/`enumerate.go` to make bodyFor attribute-aware.

**Concrete change (simplest):** Detect when `clipPathUnits="objectBoundingBox"` is the varied attribute and use fractional coordinates for the child. Since `bodyFor`/`defaultScaffold` is not currently attribute-aware, the cleanest approach is to add a separate tag-specific scaffold override. Add `clipPathUnits` to the overlay so that when this is the varied attribute, the clip child uses fractional coords:

Option A — Two-scaffold approach: add `clipPath` as a `builtinScaffoldWins` special case and emit two SVG variants based on the attribute value. Not directly supported by the current single-scaffold architecture.

Option B — Override the child geometry via `bodyFor` in `reps.go`/`emit.go` when the attribute value is `objectBoundingBox`. This requires passing the current attribute name+value into `bodyFor`.

Option C — Simplest immediate fix: change the catClip scaffold's circle to use fractional coordinates that work in BOTH modes:

```go
case catClip:
    // Use cx/cy/r in the 0-1 fraction range so the circle works in
    // objectBoundingBox mode; for userSpaceOnUse mode these coords (0.5, 0.5, 0.35)
    // are tiny (sub-pixel), so add clipPathUnits="userSpaceOnUse" to the default
    // and keep a separate objectBoundingBox scaffold.
```

Actually Option C doesn't work because a single scaffold can't satisfy both modes simultaneously with the same coordinate values (userSpaceOnUse wants 0-100 space, OBB wants 0-1 space).

**Recommended fix (Option B, targeted):** In `enumerate.go` or `emit.go`, pass the current (tag, attrName, attrValue) to a `bodyOverride(tag, attrName, attrValue) string` function that returns an alternate inner body when needed. For the `clipPathUnits="objectBoundingBox"` case return `<circle cx="0.5" cy="0.5" r="0.35"/>` instead of the standard `<circle cx="50" cy="50" r="35"/>`.

---

### 5. MISLEADING — `maskContentUnits="objectBoundingBox"` absolute child coords

**Root cause:** Same class of issue as #4. The mask child `<rect x="20" y="20" width="60" height="60" fill="#fff"/>` has absolute pixel coords. In `maskContentUnits="objectBoundingBox"`, these map to x=20×90=1800px — far outside the shape, so the mask should produce no visible content (blank). Chrome appears to clamp or treat it as covering the entire region, making the card look like `maskContentUnits="userSpaceOnUse"`.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go`, `defaultScaffold()`, `catMask` case.

**Recommended fix:** Same pattern as #4 — add a `bodyOverride` for `maskContentUnits="objectBoundingBox"` that uses fractional coordinates: `<rect x="0.1" y="0.1" width="0.8" height="0.8" fill="#fff"/>` so the mask content is visible and distinct from the userSpaceOnUse card.

---

## Summary table

| Element | Verdict | Issues |
|---|---|---|
| linearGradient | ISSUES | DANGLING href/#refgrad (×2), WEAK_EFFECT x2="0" |
| radialGradient | PERFECT | — |
| stop | ISSUES | WEAK_EFFECT offset="0", WEAK_EFFECT stop-color="currentColor" |
| pattern | PERFECT | — |
| marker | PERFECT | — |
| clipPath | ISSUES | MISLEADING clipPathUnits="objectBoundingBox" (absolute child coords) |
| mask | ISSUES | MISLEADING maskContentUnits="objectBoundingBox" (absolute child coords) |

**Priority order for round-4 fixes:**
1. DANGLING linearGradient href/#refgrad — fix in `blueprint.go` catGradient scaffold (add `<linearGradient id="refgrad">` sibling in defs) — 3-line change, no architecture needed.
2. MISLEADING clipPath/mask objectBoundingBox — requires attribute-aware body injection (new `bodyOverride` hook in emit/enumerate pipeline or separate scaffold branches).
3. WEAK_EFFECT stop offset="0" — move anchor from 0.1 to 0.3 in catStop scaffold.
4. WEAK_EFFECT stop-color="currentColor" — add `color` attribute to SVG root in catStop scaffold.
5. WEAK_EFFECT linearGradient x2="0" — adjust gradient baseline to set x1="1" when x2 is being varied.
