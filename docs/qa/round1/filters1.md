# Visual QA — Round 1: Filter Primitives Batch

**Reviewed:** 2026-06-22  
**Elements:** `filter`, `feGaussianBlur`, `feColorMatrix`, `feComponentTransfer`, `feFuncR`, `feFuncG`, `feFuncB`, `feFuncA`, `feDiffuseLighting`, `feSpecularLighting`, `feDistantLight`, `fePointLight`, `feSpotLight`

---

## Summary Table

| Element | Category | Problem Cards |
|---|---|---|
| `filter` | MOSTLY_PERFECT | Non-filter-specific attr cards (fill, stroke, lang…) are visually identical |
| `feGaussianBlur` | MOSTLY_PERFECT | in=SourceAlpha, in=BackgroundImage/Alpha cards look identical to SourceGraphic at thumbnail size |
| `feColorMatrix` | IDENTICAL_NO_EFFECT | All type/values cards show same teal; type variants lack companion values |
| `feComponentTransfer` | MOSTLY_PERFECT | Structure correct; source color weak for demonstrating channel ops |
| `feFuncR` | IDENTICAL_NO_EFFECT | ALL slope/intercept/amplitude/exponent/offset cards — no visible effect (type defaults to identity) |
| `feFuncG` | IDENTICAL_NO_EFFECT | Same as feFuncR |
| `feFuncB` | IDENTICAL_NO_EFFECT | Same as feFuncR |
| `feFuncA` | IDENTICAL_NO_EFFECT | Same root cause; slope="0" (should be transparent) renders as solid teal |
| `feDiffuseLighting` | HAS_EMPTY_CARDS | ALL 91 cards are near-black — nested element bug + no visible light |
| `feSpecularLighting` | HAS_EMPTY_CARDS | ALL 97 cards are near-black — same nested element bug |
| `feDistantLight` | HAS_EMPTY_CARDS | ALL 61 cards near-black — azimuth/elevation values too small for visible lighting |
| `fePointLight` | HAS_EMPTY_CARDS | ALL 67 cards near-black — x/y/z positions all ≤ 3.14, too close/dark |
| `feSpotLight` | HAS_EMPTY_CARDS | ALL 97 cards near-black — same positional range issue |

---

## Per-Element Findings

### `<filter>`

**Category:** MOSTLY_PERFECT

**Observation:**  
The filter element gallery (71 cards) correctly places the `<filter id="slot">` element wrapping a `feGaussianBlur stdDeviation="3"` baseline inside `<defs>`, with a teal rect `filter="url(#slot)"`. Filter-specific attribute cards (`filterUnits`, `primitiveUnits`, `x`, `y`, `width`, `height`) correctly show the filter effect (slightly blurred teal rectangle). However, the large block of presentation attribute cards (`fill="none"`, `stroke="none"`, `lang="en"`, etc.) and several others produce visually identical renders to each other, which is expected for non-filter-affecting attributes — no fix needed for those.

**Problem cards:**  
- `filter="none"` (card label) — applies `filter="none"` on the _filter element itself_, nonsensical context; visually same teal. Minor grammar oddity.

**Recommendation:**  
No structural fix needed. Optionally deduplicate or suppress pure presentation/aria attrs on filter-primitive/container elements (gallery.go or enumerate.go).

---

### `<feGaussianBlur>`

**Category:** MOSTLY_PERFECT

**Observation:**  
78 cards. The filter scaffold correctly wraps the generated `<feGaussianBlur>` in `<filter id="slot">`. The baseline `stdDeviation="3"` is applied on all non-stdDeviation cards. The blur IS rendering: subtle blurred edges are visible on `stdDeviation="2"` cards. The main weakness is that `in="SourceAlpha"` should output a soft grey gaussian blob (not teal), but at thumbnail size the difference is minimal and likely swallowed by the card background.

**Problem cards:**  
- `in="SourceAlpha"`, `in="BackgroundImage"`, `in="BackgroundAlpha"`, `in="blur1"`, `in="result1"` — dangling references or alpha-only output produce faint/invisible results that look same as SourceGraphic at small scale.
- `edgeMode="none"`, `edgeMode="wrap"` — visually indistinguishable from `edgeMode="duplicate"` at this scale.

**Recommendation:**  
Increase `stdDeviation` baseline from 3 to 8–12 in `baselineFor("feGaussianBlur")` in **blueprint.go** so the blur is unmistakably visible at thumbnail size. For `in` variants, ensure `in="SourceGraphic"` is the baseline (it already is). No structural bug.

---

### `<feColorMatrix>`

**Category:** IDENTICAL_NO_EFFECT (type/values cards) + MOSTLY_PERFECT (in/result/position cards)

**Observation:**  
83 cards. The filter scaffold is correct (`<filter id="slot"><feColorMatrix …/></filter>`). Cards for presentation attributes and `in`/`result` coordinates all appear as the same teal, which is acceptable.

**Critical issue:** The `type` variant cards and `values` variant cards produce no visible color change:

- `type="matrix"` with no `values` → browser uses identity matrix → teal unchanged.
- `type="saturate"` with no `values` → default 1.0 saturation → teal unchanged.
- `type="hueRotate"` with no `values` → default 0° rotation → teal unchanged.
- `type="luminanceToAlpha"` → should convert to grayscale alpha → appears to show some desaturation but barely distinguishable.
- `values="1 2 3 4"` / `values="1 0.5"` — wrong size (matrix type requires 20 values; 4 or 2 values → ignored or treated as partial) → teal unchanged.
- `values="0 1 0 0"` — 4 values, invalid for `type="matrix"` (needs 20) → ignored.

**Root causes:**
1. `type` varies without companion `values` → default no-op.
2. `values` samples in `ListOfNumbersType` are too short for `type="matrix"` (needs exactly 20 space-separated numbers).

**Recommendations:**
- **overlay.go**: When `attrName == "type"` on `feColorMatrix`, emit specific values that demonstrate effect:
  - `type="saturate" values="0"` (pair forced), or add `values` to the baseline when varying `type`.
- **blueprint.go `baselineFor("feColorMatrix")`**: Add `type="saturate"` and `values="0.2"` as baseline so type-varying cards show a desaturated but colored source.
- **reps.go `ListOfNumbersType`**: Add a 20-element sample for the matrix case, e.g. `"1 0 0 0 0  0 1 0 0 0  0 0 1 0 0  0 0 0 1 0"`. Consider adding a type-specific override in overlay.go.
- **Blueprint source color**: Change source rect from `#16c79a` (teal) to a warmer multicolor source (e.g. a small gradient rect or `#e94560` red) so saturation/hue changes are obviously visible.

---

### `<feComponentTransfer>`

**Category:** MOSTLY_PERFECT

**Observation:**  
74 cards. The scaffold correctly places `<feComponentTransfer>{{ELEMENT}}</feComponentTransfer>` inside a filter. The body `<feFuncR type="linear" slope="1.5"/>` is injected as baseline child when the element is `feComponentTransfer` (from `bodyFor`). Cards for `in`, `result`, coordinates, and presentation attributes all show same teal (expected). The structural wiring is correct.

**Minor issue:** Source color (teal #16c79a) has very low red (R≈9%), so a linear amplification of R is barely visible. No EMPTY cards.

**Recommendation:**  
Change source rect color in `feComponentTransfer` template (**html/template/feComponentTransfer.html**) to a richer mid-tone such as `#7b5ea7` (purple) or add `fill="#e97240"` (orange) to make channel transfers obvious.

---

### `<feFuncR>`, `<feFuncG>`, `<feFuncB>`, `<feFuncA>`

**Category:** IDENTICAL_NO_EFFECT

**Observation:**  
All four galleries (89 cards each) appear uniformly teal — every card looks identical to the baseline. The generated SVG structure is syntactically correct (e.g., `<feComponentTransfer><feFuncR slope="0"/></feComponentTransfer>`) but produces no visible effect because:

**Root cause — critical bug:** `feFuncX` elements must have `type` set to use `slope`, `intercept`, `amplitude`, `exponent`, or `offset`. When `type` is absent, the browser defaults to `type="identity"`, which ignores all those attributes and passes the channel through unchanged. The generator never adds a baseline `type="linear"` when varying slope/intercept, nor `type="gamma"` when varying amplitude/exponent/offset. Result: every card with a slope/intercept/amplitude/exponent/offset variant is silently overridden by the identity function.

**Specific problem cards (all four elements):**
- `slope="0"` — should zero out the channel (e.g., feFuncA slope=0 → fully transparent), appears as solid teal.
- `slope="-1"` — should invert/clamp, appears same.
- `slope="3.14"` — should amplify channel, appears same.
- `intercept="0.5"` — adds 0.5 to channel, appears same (needs type="linear").
- `amplitude="2"`, `exponent="3.14"`, `offset="0.5"` — all need `type="gamma"`, appear same.
- `type="table"` without `tableValues` — empty table → passthrough.
- `type="discrete"` without `tableValues` — same.
- `type="linear"` without `slope` — defaults slope=1 → identity.
- `type="gamma"` without `amplitude`/`exponent` → identity.

**Recommendations:**
- **blueprint.go `baselineFor("feFuncR"/"feFuncG"/"feFuncB"/"feFuncA")`**: Add `type="linear"` as the baseline attribute, so that slope/intercept variants are active by default.
- **overlay.go**: When `attrName` is `amplitude`, `exponent`, or `offset`, return `type="gamma"` as a paired constraint — or add a "companion attribute injector" that ensures type is set.
- **reps.go**: Ensure `tableValues` samples exist (already in `ListOfNumbersType`); the overlay should set `type="table"` or `type="discrete"` when `tableValues` is the varied attribute.
- **Blueprint source color**: Change from `#16c79a` to a richer mid-RGB color so individual channel effects are clearly distinct (e.g., `#c87040` orange: obvious R, moderate G, low B makes feFuncR/G/B clearly different).

---

### `<feDiffuseLighting>`

**Category:** HAS_EMPTY_CARDS (all 91 cards near-black)

**Blank cards:** ALL (`in=*`, `surfaceScale=*`, `diffuseConstant=*`, `kernelUnitLength=*`, position/size attrs, presentation attrs, etc.)

**Root cause — critical structural bug:** The template blueprint (**html/template/feDiffuseLighting.html**) places `{{ELEMENT}}` inside an outer `<feDiffuseLighting>` wrapper. When the element _being generated_ is itself `feDiffuseLighting`, the injection produces:

```xml
<filter id="slot">
  <feDiffuseLighting surfaceScale="5" diffuseConstant="1" lighting-color="white">   <!-- outer, from blueprint -->
    <feDiffuseLighting in="SourceGraphic">                                           <!-- inner, generated element -->
      <fePointLight x="50" y="50" z="40"/>
    </feDiffuseLighting>
  </feDiffuseLighting>
</filter>
```

`feDiffuseLighting` cannot be a child of `feDiffuseLighting`. The browser ignores the invalid nesting. The outer `feDiffuseLighting` has no valid light source child → produces black output. The rect is then replaced by a black filter result.

**Secondary issue:** Even with the nesting fixed, the light source `fePointLight x="50" y="50" z="40"` should produce visible output. The source rect fill is `#16c79a` — correct. However, the `in` attribute defaults to `SourceGraphic` which is correct.

**Recommendations:**
- **html/template/feDiffuseLighting.html**: Change the blueprint so `{{ELEMENT}}` is placed _directly inside the filter_, not inside a surrounding `feDiffuseLighting`. The generated `feDiffuseLighting` element already has `bodyFor("feDiffuseLighting")` = `<fePointLight x="50" y="50" z="40"/>` as its own child, so no outer wrapper is needed:

  ```xml
  <svg …><defs><filter id="slot">{{ELEMENT}}</filter></defs>
  <rect x="10" y="10" width="80" height="80" fill="#e94560" filter="url(#slot)"/></svg>
  ```

- **blueprint.go `baselineFor("feDiffuseLighting")`**: Add baseline attrs: `surfaceScale="5" diffuseConstant="1" lighting-color="white"` so the filter card shows visible lighting even when varying a non-lighting attribute.
- Use a contrasting (non-teal) fill such as `#e94560` on the source rect to make diffuse lighting color visible.

---

### `<feSpecularLighting>`

**Category:** HAS_EMPTY_CARDS (all 97 cards near-black)

**Blank cards:** ALL

**Root cause:** Identical nesting bug as `feDiffuseLighting`. The template blueprint wraps `{{ELEMENT}}` inside `<feSpecularLighting>`, producing invalid nested specular lighting:

```xml
<feSpecularLighting surfaceScale="5" specularConstant="1" specularExponent="20" lighting-color="white">
  <feSpecularLighting in="SourceGraphic">
    <fePointLight x="50" y="50" z="40"/>
  </feSpecularLighting>
</feSpecularLighting>
```

**Recommendations:**
- **html/template/feSpecularLighting.html**: Same fix as `feDiffuseLighting` — place `{{ELEMENT}}` directly in the filter, not inside another `feSpecularLighting`. Generated element already carries `bodyFor("feSpecularLighting")` = `<fePointLight x="50" y="50" z="40"/>`.
- **blueprint.go `baselineFor("feSpecularLighting")`**: Add `surfaceScale="5" specularConstant="1" specularExponent="20" lighting-color="white"`.

---

### `<feDistantLight>`

**Category:** HAS_EMPTY_CARDS (all 61 cards near-black)

**Blank cards:** ALL (azimuth and elevation variant cards are all black)

**Observation:**  
The template and structural wiring are correct: `feDistantLight` is placed inside `<feDiffuseLighting>` inside the filter. The generated SVG is:
```xml
<feDiffuseLighting surfaceScale="5" diffuseConstant="1" lighting-color="white">
  <feDistantLight azimuth="0"/>
</feDiffuseLighting>
```
This is valid SVG. The source rect has `fill="#16c79a"`. The filter _should_ produce visible output.

**Root cause — numeric range issue:** The `azimuth` and `elevation` attributes draw from `NumberType` samples: `{0, 1, -1, 3.14, 0.5, 2}`. These translate to angles in degrees:
- `elevation="0"` — light is on the horizon, almost horizontal → essentially no visible diffuse component on a flat surface → black.
- `elevation="1"` / `elevation="0.5"` / `elevation="2"` — 0.5–2 degrees above horizon → still near-black.
- `elevation="3.14"` — ~3° above horizon → still essentially black.
- `elevation="-1"` — below horizon → light behind the surface → no illumination → black.

No sample has elevation ≥ 30°, which is the minimum needed for visibly lit diffuse output. Similarly, `azimuth` in {0,1,-1,...} degrees barely changes the horizontal lighting direction.

**Recommendations:**
- **overlay.go `overlaySample`**: Add a case for `attrName == "elevation"` → return `"45"` (45° elevation always produces clearly visible diffuse lighting). Similarly cap `azimuth` in a visible range with samples spanning `0`, `45`, `90`, `135`, `180`.
- **reps.go `NumberType`**: Consider adding `45` and `90` to the sample set, OR handle elevation/azimuth specifically in overlay.go.
- Alternatively add `elevation="45"` to the `feDistantLight` baseline in **blueprint.go `baselineFor("feDistantLight")`** so non-elevation cards still show a visible light.

---

### `<fePointLight>`

**Category:** HAS_EMPTY_CARDS (all 67 cards near-black)

**Blank cards:** ALL (x, y, z variant cards are all black)

**Observation:**  
Structural wiring is correct: `fePointLight` inside `<feDiffuseLighting>` inside filter. The issue is positional: `x`, `y`, `z` draw from `NumberType` samples `{0, 1, -1, 3.14, 0.5, 2}`.

- `z="0"` — light at z=0 (on the surface plane) → essentially no visible illumination (light parallel to surface).
- `z="1"` / `z="2"` / `z="3.14"` — z very small relative to surface → light nearly tangent → very dark.
- `z="-1"` — light behind the surface → no illumination → black.
- `x="0", y="0"` — light at top-left corner far from centre of rect → corner lit but most of rect very dark.

None of these z values (max 3.14 in a 100-unit viewBox where the rect is 80×80) produce meaningfully bright lighting. A z of 50–100 would produce clearly visible circular highlight patterns.

**Recommendations:**
- **overlay.go**: Add case: `attrName == "z"` on `fePointLight`/`feSpotLight` → return `"50"`.
- **blueprint.go `baselineFor("fePointLight")`**: Add `x="50" y="50" z="50"` as baseline, so when varying other attributes the light is still in a good position.
- **reps.go `NumberType`**: Add `"50"` and `"100"` to the sample set (useful for z-depth, coordinate values in realistic contexts).

---

### `<feSpotLight>`

**Category:** HAS_EMPTY_CARDS (all 97 cards near-black)

**Blank cards:** ALL (x, y, z, pointsAtX, pointsAtY, pointsAtZ, specularExponent, limitingConeAngle cards)

**Observation:**  
Same root cause as `fePointLight`. Additionally, `feSpotLight` has `pointsAtX`, `pointsAtY`, `pointsAtZ` (direction target) and `limitingConeAngle`. All these draw from `NumberType` samples where:
- `z="0"` / `z="1"` / `z="-1"` → light nearly on surface plane → no visible illumination.
- `pointsAtX/Y/Z="0"` — pointing at origin while light is also near origin → undefined or backwards.
- `limitingConeAngle="0"` / `limitingConeAngle="1"` — cone so narrow that almost no light escapes.
- `specularExponent="0"` / `specularExponent="-1"` — invalid exponent values.

**Recommendations:**
- **overlay.go**: Add cases:
  - `attrName == "z"` on `feSpotLight` → `"150"` (light well above surface).
  - `attrName == "limitingconeangle"` → `"30"` (reasonable cone angle).
  - `attrName == "specularexponent"` → `"20"` (typical valid specular exponent 1–128).
  - `attrName == "pointsatx"` → `"50"`, `attrName == "pointsaty"` → `"50"`, `attrName == "pointsatz"` → `"0"` (point at centre of rect).
- **blueprint.go `baselineFor("feSpotLight")`**: Add full working baseline: `x="50" y="50" z="150" pointsAtX="50" pointsAtY="50" pointsAtZ="0" limitingConeAngle="30"`.

---

## Top Recurring Issues

### Issue 1 — Self-nesting via template blueprint (CRITICAL)
**Affects:** `feDiffuseLighting`, `feSpecularLighting`  
**Fix target:** `html/template/feDiffuseLighting.html`, `html/template/feSpecularLighting.html`  
**Change:** Replace `<feXxxLighting …>{{ELEMENT}}</feXxxLighting>` wrapper with bare `<filter id="slot">{{ELEMENT}}</filter>`. The generated element is the lighting primitive itself (with `bodyFor` providing a light source child). The outer wrapper creates illegal self-nesting that browsers silently discard, producing entirely black output.

### Issue 2 — `feFuncX` type defaults to identity (CRITICAL)
**Affects:** `feFuncR`, `feFuncG`, `feFuncB`, `feFuncA`  
**Fix target:** `blueprint.go` (`baselineFor`), `overlay.go`  
**Change:** Add `type="linear"` to `baselineFor("feFuncR/G/B/A")`. When `attrName` is `amplitude`, `exponent`, or `offset`, inject `type="gamma"` as a paired baseline. This ensures the numeric attribute variants (slope, intercept, amplitude, exponent, offset) are actually applied.

### Issue 3 — NumberType samples too small for lighting angles/depths (CRITICAL)
**Affects:** `feDistantLight`, `fePointLight`, `feSpotLight`, `feDiffuseLighting`, `feSpecularLighting`  
**Fix target:** `overlay.go` (`overlaySample`), optionally `reps.go`  
**Change:** Add domain-specific overrides: `elevation` → `"45"`, `z` (on light sources) → `"50"`, `limitingConeAngle` → `"30"`, `specularExponent` → `"20"`. All current `NumberType` samples (`{0,1,-1,3.14,0.5,2}`) are in the near-zero or negative range, producing zero or near-zero illumination for all lighting cards.

### Issue 4 — Color matrix/transfer type without companion value (MAJOR)
**Affects:** `feColorMatrix` (type variants), `feFuncR/G/B/A` (type="table"/"discrete" without tableValues)  
**Fix target:** `blueprint.go` (`baselineFor`) and/or `overlay.go`  
**Change:** Ensure `type` is always accompanied by an appropriate default for the companion attribute when it is the varied attribute. For `feColorMatrix`: varying `type` should also set a matching `values`. For `feFuncR/G/B/A`: varying `tableValues` should set `type="table"`.

### Issue 5 — Source graphic too uniform for channel-level effects (MODERATE)
**Affects:** `feColorMatrix`, `feComponentTransfer`, `feFuncR/G/B/A`  
**Fix target:** `html/template/feColorMatrix.html`, `html/template/feComponentTransfer.html`, `html/template/feFuncR/G/B/A.html`  
**Change:** Replace the monochromatic teal source rect (`fill="#16c79a"`) with a more diagnostic color. Ideal: a rect with `fill="#c87040"` (warm orange, balanced R/G, low B) or a 2-color gradient so that channel-by-channel transforms produce unmistakably different outputs. For `luminanceToAlpha` specifically, the source must not be a single flat color.
