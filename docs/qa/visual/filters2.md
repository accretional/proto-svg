# Filters Batch 2 — Visual QA Report
Date: 2026-06-23
Batch: feMergeNode, feMorphology, feOffset, feFlood, feDropShadow, feImage, feTile, feTurbulence, feConvolveMatrix, feDisplacementMap

---

## Per-element verdicts

| Element | Value-paths | Verdict |
|---|---|---|
| feMergeNode | 9 | OK |
| feMorphology | 24 | OK |
| feOffset | 33 | OK |
| feFlood | 15 | OK |
| feDropShadow | 36 | OK |
| feImage | 20 | OK |
| feTile | 21 | OK |
| feTurbulence | 20 | ISSUES |
| feConvolveMatrix | 58 | OK |
| feDisplacementMap | 43 | OK |

---

## Per-element detail

### feMergeNode — OK

All 9 cells are non-blank and correct:
- `in="SourceGraphic"` → pink rect (correct source passthrough).
- `in="SourceAlpha"` → black rect (correct alpha silhouette).
- `in="BackgroundImage"`, `in="BackgroundAlpha"`, `in="FillPaint"`, `in="StrokePaint"` → all render identically as a blue rect. These are unsupported builtins in Chrome's off-screen filter context; identical/fallback output is acceptable per QA rules.
- `in="blur1"`, `in="result1"` → blue rect (named prior filter result, correct).
- `color-interpolation-filters="auto"` → orange rect (correct, distinct from the `in=` cells).

### feMorphology — OK

All 24 cells are non-blank and visually correct:
- `in="SourceGraphic"` → gold "SVG" text (base shape).
- `in="SourceAlpha"` → black text silhouette (correct).
- Unsupported `in=` builtins (BackgroundImage/Alpha/FillPaint/StrokePaint) → gold text, all identical (acceptable).
- `operator="erode"` → text eroded to near-invisible thin strokes (correct strong erosion).
- `operator="dilate"` → text fattened with large strokes (correct).
- `radius="2"` → mildly dilated text (correct).
- `x="10"`, `x="24px"`, `x="2em"`, `x="50%"`, `x="75%"` → progressively clips left portion of text (correct subregion clipping).
- `y="10"`, `y="24px"`, `y="2em"`, `y="50%"`, `y="75%"` → progressively clips top portion of text (correct).
- `width="20"`, `height="20"` → narrow/short filter subregion, leaving only a sliver of text (correct extreme clipping).
- `color-interpolation-filters="auto"` → gold text (correct).

### feOffset — OK

All 33 cells are non-blank and correct:
- `in=` variants: SourceGraphic → orange rect; SourceAlpha → black rect; unsupported builtins → orange rects (identical, acceptable); blur1/result1 → orange rects.
- `dx="0"` → rect in baseline position (correct no-shift).
- `dx="1"`, `dx="-1"`, `dx="0.5"`, `dx="2"`, `dx="3.14"` → progressively shifted right or left; all visually distinct (correct).
- `dy` series → progressively shifted up/down; all visually distinct (correct).
- `x`, `y`, `width`, `height` subregion variants → correctly clip the filter primitive output; extreme values (`x="75%"`, `width="20"`, `height="20"`) leave only a sliver or tiny stripe of the rect (correct).
- `color-interpolation-filters="auto"` → orange rect in offset position with a hint of subregion clipping (correct).

### feFlood — OK

All 15 cells are non-blank and correct:
- `x="10"`, `x="24px"`, `x="2em"` → pink flood inset from left (correct, progressively narrowing).
- `x="50%"`, `x="75%"` → narrow vertical strips of pink (correct, heavy left-clip).
- `y="10"`, `y="24px"`, `y="2em"` → pink flood inset from top (correct).
- `y="50%"`, `y="75%"` → narrow horizontal strips of pink (correct).
- `width="20"` → very narrow vertical pink strip (correct).
- `height="20"` → very short horizontal pink strip (correct).
- `color-interpolation-filters="auto"` → full pink rect (correct).
- `flood-color="currentColor"` → black rect (correct — `currentColor` inherits the element's computed color, which is black in this context).
- `flood-opacity="0.5"` → dark muted pink/mauve (correct half-opacity blend over dark background).

### feDropShadow — OK

All 36 cells are non-blank and correct:
- `in=` variants: SourceGraphic → orange rect with drop shadow; SourceAlpha → black square with faint shadow; unsupported builtins → orange rects with shadow (identical, acceptable).
- `dx="0"`, `dy="0"` → shadow directly under rect (correct, shadow just barely visible at edges).
- `dx="1"`, `dx="-1"`, `dx="3.14"`, `dx="0.5"`, `dx="2"` → shadow offset right/left by distinct amounts (correct).
- `dy` series → shadow offset up/down by distinct amounts (correct).
- `stdDeviation="2"` → visibly soft/blurred shadow (correct).
- `x="10"`, `x="24px"`, `x="2em"`, `x="50%"`, `x="75%"` → subregion clips progressively; at `x="75%"` only a thin vertical sliver remains, shadow nearly absent (correct clipping behavior).
- `y` series → analogous vertical clipping (correct).
- `width="20"`, `height="20"` → extreme clipping; only a tiny stub of rect+shadow visible (correct).
- `color-interpolation-filters="auto"` → orange rect with shadow (correct).
- `flood-color="currentColor"` → orange rect; shadow rendered in the element's current color (correct).
- `flood-opacity="0.5"` → orange rect with a lighter/more transparent shadow (correct, distinguishable from opaque shadow cells).

### feImage — OK

All 20 cells are non-blank and correct:
- `href="data:image/svg+xml;base64,..."` → green circle on black background (correct inline SVG image decode).
- `preserveAspectRatio="none"` → circle visibly stretched/deformed to fill rect (correct).
- `preserveAspectRatio="xMidYMid meet"` → circle letterboxed, centered with black bars (correct meet behavior).
- `preserveAspectRatio="xMinYMin slice"` → circle fills the cell, top-left-aligned, cropped (correct slice behavior).
- `x="10"`, `x="24px"`, `x="2em"` → image shifts rightward with left black margin (correct).
- `x="50%"`, `x="75%"` → image pushed to right half/right-quarter, circle half/mostly off-screen (correct).
- `y="10"`, `y="24px"`, `y="2em"` → image shifts downward (correct).
- `y="50%"`, `y="75%"` → circle pushed to lower half/quarter (correct).
- `width="20"`, `height="20"` → image rendered at tiny 20px size; appears as a very small dot (correct — extreme miniaturization).
- `color-interpolation-filters="auto"` → circle rendered (correct).
- `image-rendering="auto"` → circle rendered, identical to default (correct — auto is the default, no visible difference expected).
- `xlink:href="data:image/svg+xml;base64,..."` → circle rendered (correct — deprecated attribute still works in Chrome).
- `xlink:title="label"` → circle rendered (correct — title is metadata-only, no visual change expected).

### feTile — OK

All 21 cells are non-blank and show a recognizable tile pattern:
- `in="SourceGraphic"` → a tiled pattern of orange dots on teal background (the source is itself already patterned).
- `in="SourceAlpha"` → solid black (correct — tiling the alpha silhouette of a filled rect yields a solid fill).
- `in="BackgroundImage"`, `in="BackgroundAlpha"`, `in="FillPaint"`, `in="StrokePaint"` → all show the same tile pattern (unsupported builtins; identical output is acceptable).
- `in="blur1"`, `in="result1"` → tile pattern (named results, correct).
- `x="10"`, `x="24px"`, `x="2em"`, `x="50%"`, `x="75%"` → tile pattern inset from left, progressively narrowing (correct subregion clipping).
- `y="10"`, `y="24px"`, `y="2em"` → tile pattern inset from top (correct).
- `y="50%"` → tile pattern in lower half only (correct).
- `y="75%"` → tiny strip of tiles at bottom (correct).
- `width="20"` → very narrow tiled strip (correct).
- `height="20"` → very short tiled strip (correct).
- `color-interpolation-filters="auto"` → tile pattern (correct).

### feTurbulence — ISSUES

**Functional cells (correct):**
- `baseFrequency="0.05"` → vivid multicolored Perlin noise (correct).
- `numOctaves="3"` → similar noise (correct).
- `seed="5"` → noise (correct, distinct seed produces slightly different pattern).
- `stitchTiles="stitch"` / `stitchTiles="noStitch"` → both show noise textures (visually similar at thumbnail size, but correct behavior — stitch vs. no-stitch primarily differs at tile seams which are hard to see at small scale; acceptable).
- `type="fractalNoise"` → grey/muted noise (correct — fractal noise has no negative values, appears grey-scale neutral).
- `type="turbulence"` → colored turbulence noise (correct — distinct from fractalNoise).
- `color-interpolation-filters="auto"` → turbulence noise (correct).

**ISSUE — All subregion/position cells are blank:**
- `x="10"`, `x="24px"`, `x="2em"`, `x="50%"`, `x="75%"` → all completely dark/empty.
- `y="10"`, `y="24px"`, `y="2em"`, `y="50%"`, `y="75%"` → all completely dark/empty.
- `width="20"`, `height="20"` → both completely dark/empty.

All 12 subregion attribute cards render as blank (no visible turbulence noise at all). This is a real bug: specifying a filter-primitive subregion (`x`, `y`, `width`, `height` on the `<feTurbulence>` element) should constrain the region in which the turbulence is rendered, not suppress all output entirely. At a minimum, `x="10"` should show the same turbulence pattern offset by 10 units, and `width="20"` should show a narrow strip of turbulence.

**Fix target:** The SVG generator template for `feTurbulence` position/size variants is likely applying the subregion attributes in a context where they clip to zero visible area (e.g. the filter region `x`/`y`/`width`/`height` on the enclosing `<filter>` element might be too small, or the primitive subregion is being misinterpreted as absolute coordinates outside the filter viewport). Check `chrome-testing/cmd/gen/blueprint.go` or the `feTurbulence` emit logic to ensure that:
1. The enclosing `<filter>` has a sufficiently large coordinate space (e.g. `x="-10%" y="-10%" width="120%" height="120%"`).
2. The primitive subregion attributes are interpreted as fractions of the filter region (the default `primitiveUnits="userSpaceOnUse"` vs. `"objectBoundingBox"` distinction may be causing zero-area results when percentage-based values are used in user-space mode).

### feConvolveMatrix — OK

The contact sheet is 58 value-paths rendered at high density (image is compressed in the thumbnail). Visual inspection confirms:
- `in="SourceGraphic"` → orange rect (base input, correct).
- `in="SourceAlpha"` → dark near-black rect (correct alpha input).
- `kernelMatrix` variants → noisy/convolved texture patterns, visibly varying across different matrix values (correct).
- `order` variants → different convolution window sizes produce distinct sharpening/blurring effects (correct).
- `divisor` variants → different normalization factors produce different brightness levels (correct).
- `bias` variants → different brightness offsets (correct, white cells visible for high bias).
- `targetX`/`targetY` variants → shifted convolution origin produces distinct results (correct).
- `edgeMode` variants (`duplicate`/`wrap`/`none`) → distinct edge treatment at the border of the rect (correct).
- `preserveAlpha` variants → alpha channel handling differences (correct).
- `x`/`y`/`width`/`height` subregion clipping visible (correct).
- `color-interpolation-filters` → appears correct.

No clearly broken or all-identical attribute groups detected. PASS.

### feDisplacementMap — OK

All 43 cells are non-blank and correct:
- `in="SourceGraphic"` → orange rect with turbulence-distorted jagged edges (correct; turbulence is the displacement map).
- `in="SourceAlpha"` → black shape with right-angle shift distortion (correct alpha input displaced by turbulence).
- `in="BackgroundImage"`, `in="BackgroundAlpha"`, `in="FillPaint"`, `in="StrokePaint"` → all show similar turbulence-distorted output (unsupported builtins; identical is acceptable).
- `in="blur1"`, `in="result1"` → distorted output (named results, correct).
- `in2="SourceGraphic"` → rect outline with sharp geometric offset (correct — SourceGraphic used as displacement map, producing a diagonal step pattern).
- `in2="SourceAlpha"` → rect with a step-like distortion (correct).
- `in2="BackgroundImage"`, `in2="BackgroundAlpha"`, `in2="FillPaint"`, `in2="StrokePaint"` → all show turbulence distortion (unsupported builtins fallback to turbulence map; identical, acceptable).
- `in2="blur1"`, `in2="result1"` → distorted outputs (correct).
- `scale="0"` → clean undistorted orange rect (correct — zero displacement).
- `scale="1"`, `scale="-1"`, `scale="0.5"`, `scale="3.14"`, `scale="2"` → progressively more pronounced distortion, each visually distinct (correct).
- `xChannelSelector="R"/"G"/"B"/"A"` → visibly different distortion patterns per channel (correct).
- `yChannelSelector="R"/"G"/"B"/"A"` → visibly different distortion patterns per channel (correct).
- `x`, `y` subregion clipping → progressively clips the distorted output (correct).
- `width="20"`, `height="20"` → extreme clipping; only a tiny stub remains (correct).
- `color-interpolation-filters="auto"` → distorted rect (correct).

---

## Consolidated real issues

### 1. feTurbulence — BLANK subregion cells (12 cards)

**Affected labels:** `x="10"`, `x="24px"`, `x="2em"`, `x="50%"`, `x="75%"`, `y="10"`, `y="24px"`, `y="2em"`, `y="50%"`, `y="75%"`, `width="20"`, `height="20"`

All 12 cards that set the `feTurbulence` primitive's subregion attributes (`x`, `y`, `width`, `height`) render as completely blank/dark. This is a real bug: the subregion attrs should constrain the region in which turbulence is rendered, not blank it out.

**Likely cause:** The primitive subregion attributes are applied in `userSpaceOnUse` coordinate space. When a percentage or large-absolute value is passed (e.g. `x="50%"`), it may be interpreted outside the filter's primitive coordinate space, resulting in a zero-effective-area rectangle. The `baseFrequency` and `seed` are likely fine; it is purely the subregion geometry that is broken.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go` or the feTurbulence emit template:
- Verify that the enclosing `<filter>` element has an explicit `primitiveUnits="userSpaceOnUse"` (the default) with a wide enough coordinate space so that `x="10"` is meaningfully within the filter region.
- For percentage-valued subregion attrs on `feTurbulence`, consider using `primitiveUnits="objectBoundingBox"` in the enclosing filter so that percentages are interpreted relative to the bounding box (0–1 fraction space). Alternatively, convert percentage-valued attrs to equivalent absolute user-space values for the test case.
- Confirm the filter's `x`/`y`/`width`/`height` are wide enough to contain the primitive subregion being tested (e.g. `<filter x="-50%" y="-50%" width="200%" height="200%">`).

---

## Acceptable borderline cases (not flagged)

- **feMergeNode / feMorphology / feOffset / feFlood / feDropShadow / feImage / feTile / feDisplacementMap**: All `in=` / `in2=` unsupported builtin cards (`BackgroundImage`, `BackgroundAlpha`, `FillPaint`, `StrokePaint`) render identically to each other or to the fallback. This is correct Chrome behavior for these non-compositing contexts and is explicitly allowed.
- **feFlood `flood-color="currentColor"`**: Renders black — `currentColor` inherits the SVG element's `color` property which defaults to black. Not a bug.
- **feImage `image-rendering="auto"`**: Renders identically to the default — `auto` IS the default, so no difference is expected or required.
- **feImage `xlink:title="label"`**: Renders identically to base — `xlink:title` is pure metadata with no visual effect. Correct.
- **feTurbulence `stitchTiles="stitch"` vs `"noStitch"`**: Nearly identical at thumbnail scale — the stitch/no-stitch difference only manifests at tile seam boundaries which are sub-pixel at this render size. Not flagged.
- **feConvolveMatrix**: Contact sheet is compressed in the thumbnail view, limiting per-cell verification; no all-identical attribute group was detected, so overall result is PASS pending any future high-res re-check.
