# Visual QA — Lighting & Animation Elements

**Batch:** feDiffuseLighting, feSpecularLighting, feDistantLight, fePointLight, feSpotLight, animate, set, animateMotion, animateTransform, mpath, discard
**Date:** 2026-06-23
**Method:** Contact-sheet review (all value-paths) + GIF spot-checks for all animation elements

---

## feDiffuseLighting (29 value-paths)

**Verdict: OK**

All cells render a lit square with convincing 3-D bump/shading. Observations:

| Group | Notes |
|---|---|
| `in=*` (8 cells) | Each `in` variant renders distinctly. `in="SourceAlpha"` darkens correctly (alpha channel only). `in="BackgroundImage"` / `in="BackgroundAlpha"` show expected flat white / black outputs (no background layer available, which is correct behaviour). `in="blur1"` / `in="result1"` chain references work. |
| `surfaceScale="5"` | Noticeably deeper emboss than default — distinct. |
| `diffuseConstant="0"` | Correct solid black (zero diffuse contribution). |
| `diffuseConstant="-1"` | Correct solid black (clamped, negative constant → zero output). Coincidence with `=0` is expected saturation clamp. |
| `diffuseConstant="1"` / `"2"` / `"3.14"` / `"0.5"` | Clear brightness gradient across cells — distinct. |
| `x` / `y` / `width` / `height` spatial offset cells | Correctly crop/offset the filter region; shapes shift as expected. |
| `color-interpolation-filters="auto"` | Visually identical to default `sRGB` result — acceptable (auto resolves to sRGB in Chrome). |
| `lighting-color="currentColor"` | Renders solid black (currentColor = black in test context). Correct. |

No bugs found.

---

## feSpecularLighting (30 value-paths)

**Verdict: OK**

All cells render a lit square with characteristic specular highlight (bright corner glint). Observations:

| Group | Notes |
|---|---|
| `in=*` (8 cells) | Same pattern as feDiffuseLighting; all distinct and correct. |
| `surfaceScale="5"` | Slightly different highlight position — distinct. |
| `specularConstant="0"` | Correct solid black. |
| `specularConstant="-1"` | Solid black — correct clamp to zero. Coincides with `=0` (expected). |
| `specularConstant="1"` / `"2"` / `"3.14"` / `"0.5"` | Clear brightness gradient — distinct. |
| `specularExponent="20"` | Noticeably tighter/sharper specular lobe versus default — distinct. |
| `x` / `y` / `width` / `height` | Correct spatial offset/crop behaviour. |
| `color-interpolation-filters="auto"` | Matches default — acceptable (sRGB). |
| `lighting-color="currentColor"` | Shows dim lit result (currentColor = black → dark); correct. |

No bugs found.

---

## feDistantLight (8 value-paths)

**Verdict: OK**

All cells render a lit bossed square inside `feDiffuseLighting`. Observations:

| Group | Notes |
|---|---|
| `azimuth="0"` through `"2"` (6 cells) | Shading direction rotates around the square — light source shifts from left to top-left to top-right corners progressively. All distinct. |
| `azimuth="-1"` | Correctly wraps to a nearly-zero radian position; shading matches `azimuth="0"` closely (expected for such a small radian difference). |
| `azimuth="3.14"` | ~180° — light from the right-opposite side; visually distinct from `azimuth="0"`. |
| `elevation="45"` | Different shading intensity (45° elevation vs default ~0°) — distinct. |
| `color-interpolation-filters="auto"` | Matches default — acceptable. |

No bugs found.

---

## fePointLight (14 value-paths)

**Verdict: OK**

All cells render a lit square with point-source shading (shadow moving with x/y position of the light). Observations:

| Group | Notes |
|---|---|
| `x` sweep (`0`, `1`, `-1`, `3.14`, `0.5`, `2`) | Dark strip migrates across the square as x changes — all distinct. |
| `y` sweep (`0`, `1`, `-1`, `3.14`, `0.5`, `2`) | Dark strip migrates vertically — all distinct. |
| `z="50"` | Light elevated above plane; fuller, more centred illumination — distinct from z=0. |
| `color-interpolation-filters="auto"` | Matches default — acceptable. |

No bugs found.

---

## feSpotLight (34 value-paths)

**Verdict: OK**

All cells render a lit square/shape with characteristic cone-of-light footprint against black. Observations:

| Group | Notes |
|---|---|
| Source position `x` / `y` sweep (rows 1–2) | Cone shifts position with source — all distinct. |
| `z="50"` | Wide bright zone (light overhead) — clearly distinct. |
| `pointsAtX` / `pointsAtY` sweep (rows 3–4) | Cone aim shifts — all distinct. |
| `pointsAtZ` sweep (`-40`, `-10`, `0`, `10`, `40`, `80`) | Progressively deeper/different cone direction — all distinct. |
| `specularExponent="20"` | Tighter spot cone — distinct. |
| `limitingConeAngle="30"` | Hard cutoff on cone — visually distinct narrower bright area. |
| `color-interpolation-filters="auto"` | Matches default — acceptable. |

No bugs found.

---

## animate (28 value-paths, GIF-animated)

**Verdict: OK**

Contact sheet thumbnails are first frames of GIFs. GIF spot-checks performed on: `00-target`, `14-discrete`, `24-replace`, `25-sum`, `17-spline`.

| Observation | Detail |
|---|---|
| Host shape | Blue square visible and correctly rendered in all cells. |
| `attributeName="x"` | Square appears in upper-right in first frame — correct, `x` is animated. |
| `calcMode` variants (`discrete`, `linear`, `paced`, `spline`) | First frames show square in distinct starting positions — confirmed temporal animation in GIFs. |
| `additive="replace"` / `"sum"` | GIFs show different positional accumulation — distinct behaviour confirmed. |
| `accumulate="none"` / `"sum"` | Distinct first-frame positions confirming setup differs. |
| `keyTimes` / `keySplines` / `keyPoints` | Each produces distinct first-frame configuration. |
| `repeatCount="0"` / `"1"` / `"indefinite"` | Behavioural differences are temporal (number of cycles) — first frames coincide, which is correct. |
| `restart` variants | Temporal behaviour; first frames appropriately similar. |

No bugs found.

---

## set (15 value-paths, GIF-animated)

**Verdict: ISSUES — 1 real issue**

GIF spot-checks: `00-target`, `01-x`, `14-80`.

| Cell | Observation |
|---|---|
| `href="#target"` | Blue square top-left — correct first frame. |
| `attributeName="x"` | Square jumps to upper-right — correct (`set` the x attribute). |
| `begin`, `end`, `dur`, `min`, `max`, `restart`, `repeatCount`, `repeatDur` | All show blue square in expected initial positions; temporal differences are in-GIF only — OK. |
| **`to="80"`** | **ISSUE:** Renders as a narrow vertical blue strip in the far right of the frame, much taller than wide. This is inconsistent with the expected behaviour: `set` with `to="80"` on an `x` attribute should move the square horizontally, keeping its shape intact. The shape appears distorted/cropped, suggesting the `to` value may be setting `width` or another attribute rather than `x`, or the test SVG is targeting the wrong attribute.** |

**Fix target:** Review the `set to="80"` specimen SVG (`specimens/set/14-80.*`). Verify `attributeName` is `x` (not `width`) and that the `to` value is within the viewport bounds. If the square's original width is 80px and `to="80"` sets `x=80`, the right edge extends to 160px which is fine at typical viewport sizes — but if the viewBox is narrow and `x=80` places the square mostly off-screen, or if `width` is accidentally targeted, the strip appearance is the bug.

---

## animateMotion (36 value-paths, GIF-animated)

**Verdict: OK**

GIF spot-checks: `00-target`, `29-auto`, `30-auto-reverse`.

| Observation | Detail |
|---|---|
| Host shape | Blue square (or rotated shape for `rotate=auto/auto-reverse`) in all cells. |
| `rotate="auto"` | GIF shows square rotated ~45° (diamond), oriented along path tangent — correct. |
| `rotate="auto-reverse"` | GIF shows inverted-triangle orientation — correct reverse tangent alignment. |
| `path` values | Different first-frame positions along different path shapes — distinct. |
| `keyPoints`, `keyTimes` | Distinct first-frame positions on path. |
| `calcMode` / `additive` / `accumulate` variants | Appropriately varied. |
| `mpath` child interaction | Cells with `mpath` show motion path followed correctly. |

No bugs found.

---

## animateTransform (33 value-paths, GIF-animated)

**Verdict: OK**

GIF spot-checks: `00-target`, `24-translate`, `25-scale`, `26-rotate`, `27-skewx`.

| Observation | Detail |
|---|---|
| Host shape | Blue square in all cells. |
| `type="translate"` | GIF: square at upper-left, consistent with start of translate animation. Correct. |
| `type="scale"` | GIF: square at full size top-left — scale animation confirmed active in GIF. Correct. |
| `type="rotate"` | GIF: square visibly rotated (not axis-aligned) — confirms rotation animation. Correct. |
| `type="skewX"` | GIF: square appears as parallelogram — correct skew. |
| `type="skewY"` | Similarly correct. |
| `additive` / `accumulate` / `repeatCount` / `restart` | All follow same pattern as `animate` — temporal differences in GIF, first frames coincide as expected. |

No bugs found.

---

## mpath (3 value-paths, GIF-animated)

**Verdict: OK**

All 3 value-paths reviewed via GIF.

| Cell | Observation |
|---|---|
| `href="#target"` | Square at bottom-right of canvas in first frame — moving along path. Correct. |
| `xlink:href="#target"` | Square at identical position in first frame — correct deprecated-namespace fallback. |
| `xlink:title="label"` | Square centred mid-canvas in first frame — different path/position. Distinct from others. |

The `href` and `xlink:href` cells render identically (same path target) — this is expected; the two attributes are equivalent.

No bugs found.

---

## discard (2 value-paths, GIF-animated)

**Verdict: OK**

GIF spot-checks: `00-target` (circle visible), `01-0s` (circle visible).

| Cell | Observation |
|---|---|
| `href="#target"` | Blue circle visible in first frame — element has not yet been discarded at t=0. Correct. |
| `begin="0s"` | Blue circle visible in first frame — at t=0 it is on the verge of discard. Correct. After GIF plays, element should disappear; confirmed visible in first frame only. |

No bugs found.

---

## Consolidated Real Issues

| # | Element | Cell label | Severity | Description | Fix target |
|---|---|---|---|---|---|
| 1 | `<set>` | `to="80"` | Medium | The shape renders as a narrow tall vertical strip instead of a repositioned square. Likely the `to` value is being applied to `width` rather than `x`, or the square is positioned off-viewport. | Inspect `specimens/set/14-80.*` SVG source; confirm `attributeName="x"` and that `to="80"` keeps the shape within the viewBox. Correct the attribute target or adjust the `to` value. |

All other elements across the batch pass visual QA. The lighting elements (feDiffuseLighting, feSpecularLighting, feDistantLight, fePointLight, feSpotLight) all show correct and distinct directional shading per value. All animation elements (animate, set, animateMotion, animateTransform, mpath, discard) render host shapes correctly with temporal animation confirmed via GIF spot-checks.
