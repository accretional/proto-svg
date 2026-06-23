# Filter Batch — Round 3 QA Report
Date: 2026-06-22

## Per-element verdicts

| Element | Verdict |
|---|---|
| filter | PERFECT |
| feGaussianBlur | PERFECT |
| feColorMatrix | ISSUES |
| feComponentTransfer | PERFECT |
| feFuncR | ISSUES |
| feFuncG | ISSUES |
| feFuncB | ISSUES |
| feFuncA | ISSUES |
| feDiffuseLighting | PERFECT |
| feSpecularLighting | PERFECT |
| feDistantLight | PERFECT |
| fePointLight | PERFECT |
| feSpotLight | PERFECT |
| feBlend | PERFECT |
| feComposite | PERFECT |
| feMerge | PERFECT |
| feMergeNode | PERFECT |
| feMorphology | PERFECT |
| feOffset | PERFECT |
| feFlood | PERFECT |
| feDropShadow | PERFECT |
| feImage | PERFECT |
| feTile | ISSUES |
| feTurbulence | PERFECT |
| feConvolveMatrix | PERFECT |
| feDisplacementMap | PERFECT |

---

## Consolidated remaining issues

### 1. feColorMatrix — GRAMMAR/VALUE (2 cards)

**Card: `type="matrix"`**
- The overlay uses `values="0.3"` alongside `type="matrix"`. The `matrix` type requires exactly 20 space-separated numbers. A single `0.3` is silently discarded by the browser (identity fallback). The card looks identical to the `type="saturate"` cards (same brownish rect, no visible difference).
- Fix: in `reps.go` (or the feColorMatrix blueprint), when the attr under test is `type`, emit a valid 20-number identity-ish matrix for the `type="matrix"` card — e.g. `values="1 0 0 0 0  0 1 0 0 0  0 0 1 0 0  0 0 0 1 0"`. Target: `chrome-testing/cmd/gen/reps.go` or `blueprint.go` feColorMatrix type-rep context.

**Card: `values="10; 45; 80"`**
- The values attribute is shown with a semicolon-separated list (`"10; 45; 80"`) while the context type is `saturate` (which takes a single scalar). Semicolons are not valid separators in SVG number lists; browsers treat this as a parse failure and fall back to default. The card label is misleading (implies three values for a single-value attribute).
- Fix in `reps.go`/`overlay.go`: for the `values` representative on feColorMatrix, emit space-separated values appropriate to the current type, e.g. for `type="saturate"` use `values="0.5"`, or switch to `type="matrix"` with a valid 20-number matrix. Alternatively provide a distinct set of space-separated matrix values: `values="1 0 0 0 0 0 1 0 0 0 0 0 1 0 0 0 0 0 1 0"`.

---

### 2. feFuncR / feFuncG / feFuncB / feFuncA — WEAK_EFFECT (all type/slope/intercept/amplitude/exponent/offset cards)

All four pages use a black-to-white greyscale gradient (`linearGradient` from `#000000` to `#ffffff`) as the input image. Because the input is achromatic (R = G = B everywhere), manipulating only one channel (e.g. `feFuncR` with `slope="-1"`) produces a colour shift (cyan/magenta/yellow tones) that **should** be visible but the cards all appear uniformly grey in the screenshots. The effect is present in theory but visually indistinguishable at thumbnail size because:
- The fixed `slope="1.5"` stacked on `type="identity"` is a no-op.
- Negative slope on a grey gradient still yields a grey (R inverted → same grey when G=B are unchanged but R is the only channel altered, the delta is small in a grey ramp).

**Root cause**: an achromatic gradient isolates one channel per page with no colour contrast to make the difference pop.

**Fix** (target: `chrome-testing/cmd/gen/reps.go` or the feFuncR/G/B/A blueprint context):
- Use a **coloured** gradient as the input instead of a greyscale ramp, e.g. a rainbow or a saturated red-to-blue gradient:
  ```svg
  <linearGradient id="bpGrad" x1="0" y1="0" x2="1" y2="0">
    <stop offset="0" stop-color="#ff0000"/>
    <stop offset="1" stop-color="#0000ff"/>
  </linearGradient>
  ```
  With a red-to-blue gradient, `feFuncR slope="-1"` visibly turns red areas cyan (R→0) while leaving blue areas purple. `feFuncG slope="2"` adds green cast on dark areas. This makes each page and each card visually distinct.
- Also remove the redundant fixed `slope="1.5"` that is baked into type-row cards when `slope` is not the attribute under test.

---

### 3. feTile — WEAK_EFFECT (all `in=*` and position-attr cards)

The feTile blueprint feeds the element a plain white rectangle as `SourceGraphic`. Tiling a uniform white fill produces a uniform white output — no tile seams or pattern are visible in any card. The non-`in` cards (`x`, `y`, `width`, `height`, `result`) are also all white.

**Fix** (target: `chrome-testing/cmd/gen/blueprint.go` or the feTile template in `emit.go`):
- Use a small non-uniform source for `feTile`. The blueprint should define a preliminary small coloured pattern via `feFlood`+`feComposite` or a small `<image>` that is then piped into `feTile`. A working minimal example:
  ```svg
  <defs>
    <filter id="slot">
      <feFlood flood-color="#f5a623" x="0" y="0" width="25%" height="25%" result="patch"/>
      <feTile in="patch" {{.AttrUnderTest}}></feTile>
    </filter>
  </defs>
  <rect x="0" y="0" width="100" height="100" filter="url(#slot)"/>
  ```
  Using `in="patch"` (the small feFlood quarter-rect) as the tile source will clearly show a 4×4 repeating orange square pattern, and changing `x`/`y`/`width`/`height` on `feTile` will visibly shift or scale the tiling region.

---

## Notes on acceptable borderline cases (not flagged as issues)

- **filter**: All 19 main-grid cards show a blurred teal rectangle. The blur IS visible at edges (not a plain sharp rect). Different `x`/`y`/`width`/`height` values produce subtly different clipping of the blur region, which is the correct behaviour for filter subregion attrs. PERFECT.
- **feGaussianBlur**: The `in="SourceAlpha"` / `in="BackgroundImage"` etc. cards show a near-transparent or uniform shape — this is correct SVG behaviour (alpha source = black silhouette, BackgroundImage = empty). PERFECT.
- **feDropShadow**: The `x`/`y`/`width`/`height` cards correctly clip the filter primitive subregion, making the shadow appear cropped or absent. This is accurate behaviour, not a rendering error.
- **feDiffuseLighting / feSpecularLighting**: All cards show a light-on-grey-rect with subtle 3-D shading. The effect is clearly visible and cards are distinct across `surfaceScale`, `diffuseConstant`/`specularConstant`, `kernelUnitLength`. PERFECT.
- **feDistantLight / fePointLight / feSpotLight**: All show the embedded lighting effect (diffuse or specular parent context) with visibly different light source positions across cards. PERFECT.
- **feConvolveMatrix**: The `kernelMatrix`, `divisor`, `bias`, `targetX`/`Y`, `edgeMode`, `preserveAlpha` cards all show a yellow rect with subtly different convolution shading. Distinct. PERFECT.
- **feDisplacementMap**: All cards show an orange rect with a turbulence-distorted jagged edge. Cards are clearly distinct across `scale`, `xChannelSelector`, `yChannelSelector`. PERFECT.
