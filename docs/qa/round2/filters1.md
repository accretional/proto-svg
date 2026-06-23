# Round 2 Visual QA — filters1 batch

Elements: `filter`, `feGaussianBlur`, `feColorMatrix`, `feComponentTransfer`,
`feFuncR`, `feFuncG`, `feFuncB`, `feFuncA`,
`feDiffuseLighting`, `feSpecularLighting`,
`feDistantLight`, `fePointLight`, `feSpotLight`

Reviewed: 2026-06-22. Screenshots in `chrome-testing/screenshots/review/`.

---

## Per-element findings

### `<filter>`
**Status: NO_EFFECT_IN_MAIN_GRID (systematic)**

All main-grid cards render a teal rectangle — the `<filter>` element itself is a
container and its only functional attributes are `filterUnits`, `primitiveUnits`,
`x`, `y`, `width`, `height`, `href`/`xlink:href`. All other attributes visible in
the main grid are presentation attributes (`fill`, `fill-rule`, `fill-opacity`,
`stroke`, `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`,
`stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`,
`marker`, `marker-start`, `marker-mid`, `marker-end`, `color`,
`color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`,
`filter`, `color-interpolation-filters`, `flood-color`, `xlink:href`,
`xlink:title`) that have no effect on a `<filter>` element — these should move to
the non-visual / metadata `<details>` section.

Functional attrs (`filterUnits`, `primitiveUnits`, `x`, `y`, `width`, `height`,
`href`) correctly appear in the main grid but are visually identical (all show
uniform teal rect) because a plain `feGaussianBlur` child with stdDeviation=3 is
used as the body; the crop/unit changes are not perceptible on a uniform source.

**Fix target:** `blueprint.go` / `gallery.go` — move presentation attrs out of
the filter gallery's main grid. The `<filter>` body child (`feGaussianBlur`) is
fine.

---

### `<feGaussianBlur>`
**Status: PARTIAL — WEAK_EFFECT for most cards; NO_EFFECT_IN_MAIN_GRID for paint attrs**

Round-1 fix confirmed: cards render (not blank). `stdDeviation="2"` card shows
slight edge blur on the rect. However, blur on a uniform-fill rect is only visible
at the edges — the interior looks identical. Many parameter variations are
indistinguishable.

NO_EFFECT_IN_MAIN_GRID: `fill`, `fill-rule`, `fill-opacity`, `stroke`,
`stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`,
`stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`,
`marker`, `marker-start`, `marker-mid`, `marker-end`, `color`,
`color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`,
`filter`, `color-interpolation-filters`, `flood-color`.

**Fix target:** `gallery.go` — move presentation/paint attrs for all filter
primitives to the non-visual section. These are attrs from the global SVG
presentation attribute set that the grammar over-approximates for filter
primitives but have no rendering effect on them.

---

### `<feColorMatrix>`
**Status: NEAR-PERFECT — functional attrs visible; NO_EFFECT_IN_MAIN_GRID for paint attrs**

Round-1 fix confirmed: cards render the source rect in a warm desaturated tone
(the `type="saturate" values="0.3"` baseline works). `type="matrix"`,
`type="saturate"`, `type="hueRotate"`, `type="luminanceToAlpha"`,
`values="10; 45; 80"` cards each show distinctly colored output — effect is
clearly visible.

NO_EFFECT_IN_MAIN_GRID: same full set of paint/stroke/marker attrs as above.

**Fix target:** `gallery.go` — move presentation attrs to non-visual section for
all filter primitives (systematic fix).

---

### `<feComponentTransfer>`
**Status: NO_EFFECT_IN_MAIN_GRID for paint attrs; WEAK_EFFECT on functional attrs**

All functional cards show an orange rect; the baseline `feFuncR type=linear
slope=1.5` shifts the red channel up — this is visible as a warm orange (vs the
original blue `#4d8bff`). However all cards look identical to each other because
the `in=` variation doesn't change the color when background/stroke inputs are
empty, and `x/y/width/height/result` have no perceptible color effect.

NO_EFFECT_IN_MAIN_GRID: same presentation attrs.

**Fix target:** `gallery.go` — systematic presentation-attr move.

---

### `<feFuncR>`, `<feFuncG>`, `<feFuncB>`, `<feFuncA>`
**Status: WEAK_EFFECT on functional attrs; NO_EFFECT_IN_MAIN_GRID for paint attrs**

Round-1 fix confirmed: `type="linear"` with `slope="1.5"` baseline means the
element is active (not identity no-op). Cards render an orange/brown rect.

Remaining issues:

1. **WEAK_EFFECT (source is flat uniform rect):** The blueprint source is a solid
   `#4d8bff` rect. A linear channel transform on a uniform color produces a
   different uniform shade — but because all pixels in the source have the same
   channel value, all output pixels are also uniform. `slope="0"` (should make
   the channel zero → black channel), `slope="0.5"` (halved), `slope="2"`
   (doubled, clamps to 1.0), and `slope="1"` (identity) are all rendered as
   the same brownish orange. The differences exist but are extremely subtle
   or entirely hidden by clamping. `type="identity"` vs `type="linear"` vs
   `type="gamma"` vs `type="table"` vs `type="discrete"` look identical.
   `tableValues` variations similarly show no visible distinctions.

   **Fix:** Change the blueprint source for `catTransferFn` from a flat rect to a
   **linear gradient rect** (e.g. a `<linearGradient>` rect going from black to
   white, or from `#000000` to `#ffffff`). This way each pixel has a different
   channel value and the transfer function shape is directly visible as a
   tonal curve on the rendered output.

   In `blueprint.go`, line 130–132, change `catTransferFn` scaffold to:
   ```go
   case catTransferFn:
       return svgOpen +
           `<defs><linearGradient id="grad" x1="0" y1="0" x2="1" y2="1">` +
           `<stop offset="0" stop-color="#000000"/><stop offset="1" stop-color="#ffffff"/>` +
           `</linearGradient><filter id="slot"><feComponentTransfer>{{ELEMENT}}</feComponentTransfer></filter></defs>` +
           `<rect x="10" y="10" width="80" height="80" fill="url(#grad)" filter="url(#slot)"/></svg>`
   ```
   With a black-to-white gradient source, `slope="0"` → solid black channel,
   `slope="0.5"` → half-tone gradient, `slope="2"` → clipped bright gradient;
   `type="gamma"` shows a gamma curve; `type="table"` with distinct tableValues
   shows quantized bands; `type="discrete"` shows stepped posterization.

2. **NO_EFFECT_IN_MAIN_GRID:** `fill`, `fill-rule`, `fill-opacity`, `stroke`,
   `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`,
   `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`,
   `marker`, `marker-start`, `marker-mid`, `marker-end`, `color`,
   `color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`,
   `filter`, `color-interpolation-filters`, `flood-color`.

   These are SVG presentation attributes with no effect on `feFuncR/G/B/A`.
   They should move to the non-visual `<details>` section.

**Fix target:** `blueprint.go` (`catTransferFn` scaffold — use gradient source);
`gallery.go` (move paint attrs to non-visual section for catTransferFn).

---

### `<feDiffuseLighting>`
**Status: WEAK_EFFECT — no longer all-black (round-1 fix confirmed); NO_EFFECT_IN_MAIN_GRID for paint attrs**

Round-1 fix confirmed: cards render a gray/silver lit surface — not black.
`surfaceScale="5"`, `diffuseConstant` variations, `kernelUnitLength` variations
are present in main grid and show a lit surface.

Remaining issues:

1. **WEAK_EFFECT (flat source):** The source rect is a flat uniform `#e94560`
   fill. `feDiffuseLighting` computes a surface normal from the *alpha gradient* of
   the source. A flat opaque rect has zero alpha gradient in the interior — all
   interior normals point straight up (0, 0, 1). The result is a uniformly lit
   flat surface where `diffuseConstant` variations (0, 1, -1, 3.14, 0.5, 2)
   produce subtly different overall brightness but no visible 3D surface detail.
   `surfaceScale` has no visible effect without alpha variation. Azimuth/elevation
   changes on `feDistantLight` children likewise produce no shading variation.

   **Fix:** The blueprint source for `catFilterPrimitive` (shared by
   `feDiffuseLighting`) needs a bumpy alpha-gradient source to show lighting.
   Best option: add a pre-blur pass that creates an alpha ramp, e.g. use a
   `feGaussianBlur` of `SourceAlpha` as input, so the lighting sees a rounded
   alpha hill and generates 3D shading. The dedicated lighting template should
   replace the flat rect with a **radial alpha gradient source** or a blurred
   alpha mask. Example fix in `blueprint.go` `catFilterPrimitive` case or in a
   per-element template override:
   ```
   <feGaussianBlur in="SourceAlpha" stdDeviation="15" result="bumpMap"/>
   <feDiffuseLighting in="bumpMap" ...>
   ```
   Alternatively, use an `feImage` of a bump texture or a turbulence layer.

2. **NO_EFFECT_IN_MAIN_GRID:** `fill`, `fill-rule`, `fill-opacity`, `stroke`,
   `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`,
   `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`,
   `marker`, `marker-start`, `marker-mid`, `marker-end`, `color`,
   `color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`,
   `filter`, `color-interpolation-filters`, `flood-color`.

**Fix target:** `blueprint.go` — `catFilterPrimitive` scaffold: add a
`feGaussianBlur in="SourceAlpha" result="bumpMap"` before `{{ELEMENT}}` and
update the `feDiffuseLighting` / `feSpecularLighting` baseline in `baselineFor`
to use `in="bumpMap"`. Or create a per-element template in
`chrome-testing/html/template/feDiffuseLighting.html`.

---

### `<feSpecularLighting>`
**Status: WEAK_EFFECT (same root cause as feDiffuseLighting); NO_EFFECT_IN_MAIN_GRID for paint attrs**

Round-1 fix confirmed: cards render dark gray/navy — not all-black. A specular
highlight is present but very subtle.

Remaining issues:

1. **WEAK_EFFECT (flat source):** Same root cause as `feDiffuseLighting` — flat
   source rect has zero alpha gradient, so specular shading has no surface
   variation to react to. `specularConstant` and `specularExponent` variations
   are virtually indistinguishable. Cards are uniformly dark.

2. **NO_EFFECT_IN_MAIN_GRID:** Same presentation attrs as above.

**Fix target:** Same as `feDiffuseLighting` — add `feGaussianBlur` of
`SourceAlpha` as the bump source, or use a dedicated template.

---

### `<feDistantLight>`
**Status: WEAK_EFFECT on azimuth/elevation; NO_EFFECT_IN_MAIN_GRID for paint attrs**

Round-1 fix confirmed: cards render a lit (not black) surface. The `azimuth`
and `elevation` variations are visible in principle but **all look identical** due
to the flat rect source. Even with a bumpy source, `azimuth` changes direction but
a flat rect has no surface relief to react to directional lighting.

NO_EFFECT_IN_MAIN_GRID: `fill`, `fill-rule`, `fill-opacity`, `stroke`,
`stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`,
`stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`,
`marker`, `marker-start`, `marker-mid`, `marker-end`, `color`,
`color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`,
`filter`, `color-interpolation-filters`, `flood-color`.

**Fix target:** `blueprint.go` `catLight` scaffold — same bump source fix as
`feDiffuseLighting`. With a hemisphere-shaped alpha ramp, azimuth changes will
show clearly as different shading directions across the bump.

---

### `<fePointLight>`
**Status: NEAR-PERFECT — good positional lighting visible; NO_EFFECT_IN_MAIN_GRID for paint attrs**

Round-1 fix confirmed. Cards show a centered diffuse radial glow with `x`, `y`,
`z` variations producing visible positional changes in the light highlight.
`z="50"` (overlay default) gives a clear gradient from center to edge.

NO_EFFECT_IN_MAIN_GRID: same presentation attrs.

**Fix target:** `gallery.go` — systematic paint attr move to non-visual section.

---

### `<feSpotLight>`
**Status: PERFECT — clear cone spotlight effect; NO_EFFECT_IN_MAIN_GRID for paint attrs**

Cards render a clear circular spotlight cone. `x`, `y`, `z`, `pointsAtX`,
`pointsAtY`, `pointsAtZ`, `specularExponent`, `limitingConeAngle` variations
are clearly distinguishable — the cone shifts, narrows, and changes character
visibly across cards.

NO_EFFECT_IN_MAIN_GRID: same presentation attrs.

**Fix target:** `gallery.go` — systematic paint attr move to non-visual section.

---

## Systematic issues (all filter primitives and their children)

### Issue 1 — Presentation attrs flooding the main grid

**All filter primitive elements** (`feGaussianBlur`, `feColorMatrix`,
`feComponentTransfer`, `feFuncR/G/B/A`, `feDiffuseLighting`, `feSpecularLighting`,
`feDistantLight`, `fePointLight`, `feSpotLight`) and the `<filter>` element itself
receive the full SVG presentation attribute set from the grammar (fill, stroke,
marker, paint-order, color-*, shape-rendering, vector-effect, flood-color, etc.).
These attributes have **no effect on filter primitives** — they are silently
ignored. They currently appear in the main grid, making cards appear as functional
demonstrations when they are not.

**Affected attrs (non-exhaustive):** `fill`, `fill-rule`, `fill-opacity`,
`stroke`, `stroke-opacity`, `stroke-width`, `stroke-linecap`, `stroke-linejoin`,
`stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset`, `paint-order`,
`marker`, `marker-start`, `marker-mid`, `marker-end`, `color`,
`color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`,
`filter`, `color-interpolation-filters`, `flood-color` (and `xlink:href`,
`xlink:title` on `<filter>`).

**Fix target:** `gallery.go` — add a `nonVisualForTag(tag string) []string`
function (or extend the existing non-visual attribute list) that, for any
`catFilterPrimitive`/`catTransferFn`/`catLight`/`catFilter` tag, marks all
presentation attrs (those not in the element's defined attribute set) as
non-visual. Alternatively, maintain an explicit allowlist of main-grid attrs
per filter element.

### Issue 2 — Flat source for channel/lighting effects (WEAK_EFFECT)

`feFuncR/G/B/A`, `feComponentTransfer`, `feDiffuseLighting`, `feSpecularLighting`,
`feDistantLight` all suffer from the same root cause: the blueprint source is a
**flat uniform rect**. Transfer-function effects require a multi-tonal source;
lighting effects require an alpha-gradient source.

**Fix target — transfer functions:** `blueprint.go` `catTransferFn` scaffold —
replace flat rect source with a diagonal black-to-white linear gradient:
```go
case catTransferFn:
    return svgOpen +
        `<defs>` +
        `<linearGradient id="bpGrad" x1="0" y1="0" x2="1" y2="1">` +
        `<stop offset="0" stop-color="#000000"/>` +
        `<stop offset="0.5" stop-color="#4d8bff"/>` +
        `<stop offset="1" stop-color="#ffffff"/>` +
        `</linearGradient>` +
        `<filter id="slot"><feComponentTransfer>{{ELEMENT}}</feComponentTransfer></filter></defs>` +
        `<rect x="10" y="10" width="80" height="80" fill="url(#bpGrad)" filter="url(#slot)"/></svg>`
```

**Fix target — lighting:** `blueprint.go` `catLight` scaffold (and
`feDiffuseLighting`/`feSpecularLighting` baselines) — prepend a
`feGaussianBlur in="SourceAlpha" stdDeviation="15" result="bumpMap"` step and
use `in="bumpMap"` on the lighting primitive. Also update the `catFilterPrimitive`
blueprint so that when `feDiffuseLighting`/`feSpecularLighting` appear as the main
element (not a parent), they also get a bumpy input:
```go
// In baselineFor, for feDiffuseLighting:
return add([2]string{"in", "bumpMap"}, [2]string{"surfaceScale", "5"}, ...),  false
// and add a pre-step feGaussianBlur in the catFilterPrimitive scaffold result="bumpMap"
```
Or use `chrome-testing/html/template/feDiffuseLighting.html` and
`feSpecularLighting.html` blueprints that supply the pre-blur step explicitly.

---

## Summary table

| Element | Status | Primary issue |
|---|---|---|
| `filter` | NO_EFFECT_IN_MAIN_GRID | ~25 presentation attrs in main grid |
| `feGaussianBlur` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Flat source; paint attrs in main grid |
| `feColorMatrix` | NEAR-PERFECT + NO_EFFECT_IN_MAIN_GRID | Paint attrs in main grid |
| `feComponentTransfer` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Flat source; paint attrs in main grid |
| `feFuncR` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Flat uniform source; paint attrs in main grid |
| `feFuncG` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Same as feFuncR |
| `feFuncB` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Same as feFuncR |
| `feFuncA` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Same as feFuncR |
| `feDiffuseLighting` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Flat source (no alpha gradient); paint attrs |
| `feSpecularLighting` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Flat source; paint attrs |
| `feDistantLight` | WEAK_EFFECT + NO_EFFECT_IN_MAIN_GRID | Flat source; paint attrs; azimuth invisible |
| `fePointLight` | NEAR-PERFECT + NO_EFFECT_IN_MAIN_GRID | Paint attrs in main grid |
| `feSpotLight` | PERFECT + NO_EFFECT_IN_MAIN_GRID | Paint attrs in main grid |
