# Visual QA — filters1 batch

Batch: feGaussianBlur, feColorMatrix, feComponentTransfer, feFuncR, feFuncG, feFuncB, feBlend, feComposite, feMerge  
Reviewed: 2026-06-23  
Method: contact sheet inspection + specimen drill-down for ambiguous cells

---

## feGaussianBlur — OK

25 value-paths. All cells non-blank and individually distinct where expected.

- `in="SourceGraphic"` — teal rect, sharp. Correct baseline.
- `in="SourceAlpha"` — black shape (alpha mask rendered correctly).
- `in="BackgroundImage"`, `in="BackgroundAlpha"`, `in="FillPaint"`, `in="StrokePaint"` — all render as the unfiltered teal rect, identical to each other. **Allowed**: browsers do not support these builtins; faithful blank/passthrough.
- `in="blur1"`, `in="result1"` — both render blurred teal rects, visually distinct from unblurred source. Correct.
- `stdDeviation="2"` — clearly blurred relative to identity. Correct.
- `edgeMode="duplicate"`, `edgeMode="wrap"`, `edgeMode="none"` — each shows a distinct edge treatment on the blurred result. Correct.
- `x`, `y`, `width`, `height` geometry values — progressive clipping/offset of the filter region across cells; all show expected spatial shifts. Correct.
- `color-interpolation-filters="auto"` — matches default rendering. Correct.

**Verdict: OK**

---

## feColorMatrix — OK (one noted specimen quirk, not a real bug)

26 value-paths.

- `in` builtins (BackgroundImage, BackgroundAlpha, FillPaint, StrokePaint) — identical noisy turbulence render. **Allowed**.
- `type="matrix"` (index 08) — uses `values="90"` in the HTML, which is an incomplete matrix (only 1 of 20 values). Chrome treats it as invalid and applies a near-identity transform on the turbulence source. Renders as a dark muted turbulence image. This is a **data authoring issue** (the specimen value is not a valid 20-float matrix), but the browser render is faithful to the invalid input. The contact sheet label correctly says `type="matrix"` — the specimen just doesn't demonstrate a real matrix transformation. **Flag for data fix, not a rendering bug.**
- `type="saturate"` (index 09) — vivid multicolor turbulence. Clearly distinct from identity. Correct.
- `type="hueRotate"` (index 10, `values="90"`) and `values="120"` (index 12, hueRotate 120°) — both render dark muted turbulence; visually similar at contact sheet thumbnail size because turbulence hue rotation on a low-saturation base produces subtle differences. Specimens are genuinely distinct on close inspection. Correct renders.
- `type="luminanceToAlpha"` (index 11) — renders as a near-black shape (luminance mapped to alpha channel, composited over dark page background). Visually almost invisible but correct — the alpha contains real grayscale data, it just shows as dark against a dark page. Not a bug.
- `values="120"` — see above. Correct.
- Geometry (`x`, `y`, `width`, `height`) — correct progressive spatial shifts.
- `color-interpolation-filters="auto"` — matches default. Correct.

**Verdict: OK**  
**Data note**: `type="matrix"` specimen should use a real 20-value matrix string to demonstrate the filter effect. Fix target: `chrome-testing/html/specimen/feColorMatrix/08-matrix.html` (update `values` to a valid 4×5 matrix, e.g. identity or a channel-swap).

---

## feComponentTransfer — OK

21 value-paths. Uses a red/green/blue gradient square as source.

- `in` builtins — identical gradient passthrough. **Allowed**.
- `in="blur1"`, `in="result1"` — both show the gradient, same-looking because a child `feFunc*` with identity defaults passes through unchanged. Correct.
- Geometry values (`x`, `y`, `width`, `height`) — progressive spatial clipping; all correctly shift/resize the rendered region. Correct.
- `color-interpolation-filters="auto"` — matches default. Correct.

No anomalies detected at this element level (feFunc children are exercised separately below).

**Verdict: OK**

---

## feFuncR — OK

41 value-paths. Source: RGB gradient square; feFuncR modifies only the red channel.

- `type="identity"` — full RGB gradient. Correct baseline.
- `type="table"`, `type="discrete"`, `type="linear"`, `type="gamma"` — each produces a visually distinct red-channel curve on the gradient. All differ. Correct.
- `tableValues` variants (5 cells, indices 05–09) — each produces a noticeably different red channel distribution. Correct.
- `slope="0"` — red channel zeroed; image shows green+blue only (no red). Correct.
- `slope="-1"` — red channel inverted. Correct.
- `slope` progression (0, 1, -1, 3.14, 0.5, 2) — all visually distinct. Correct.
- `intercept="0"` — matches default linear pass-through. Correct.
- `intercept="1"` and `intercept="2"` — both appear identical (full red everywhere). **Allowed**: intercept≥1 clamps red to maximum; both produce the same output.
- `amplitude="0"` and `amplitude="-1"` — both produce no red output (clamped to 0). **Allowed**: gamma with amplitude≤0 maps to black.
- `exponent="0"` — all inputs map to amplitude (1 by default), producing uniform red. Correct.
- `offset="1"` and `offset="2"` — both produce full red (offset≥1 clamps). **Allowed**.
- `offset="-1"` — red channel zeroed (negative offset clamps to 0). Correct.
- `color-interpolation-filters="auto"` — matches default. Correct.

**Verdict: OK**

---

## feFuncG — OK

41 value-paths. Identical attribute structure to feFuncR; modifies the green channel.

Contact sheet shows the same pattern of distinct renders as feFuncR, with effects visible in the green channel instead of red. All clamping behavior identical and allowed. No anomalies.

**Verdict: OK**

---

## feFuncB — OK

41 value-paths. Identical attribute structure; modifies the blue channel.

Same pattern of distinct and correct renders, clamping cases all allowed. No anomalies.

**Verdict: OK**

---

## feBlend — OK

45 value-paths. Source: teal background with pink and orange overlapping rects; `in2` feeds a secondary blurred/noise layer.

- `in="SourceGraphic"` — teal+shapes, sharp. Correct.
- `in="SourceAlpha"` — black alpha mask silhouette. Correct.
- `in` builtins (BackgroundImage, BackgroundAlpha, FillPaint, StrokePaint) — noisy blended render (browser uses some internal state). **Allowed**.
- `in2` builtins — same allowance applies.
- `in2="blur1"`, `in2="result1"` — blend with blurred result; distinct noisy blended output. Correct.
- `mode` values (normal, multiply, screen, darken, lighten, overlay, color-dodge, color-burn, hard-light, soft-light, difference, exclusion, hue, saturation, color, luminosity) — all 16 modes render visibly distinct blended outputs. Each mode is clearly recognizable (e.g. color-burn shows near-black in overlap, screen shows brightened overlap, difference inverts colors). Correct.
- Geometry (`x`, `y`, `width`, `height`) — correct progressive spatial shifts. Correct.
- `color-interpolation-filters="auto"` — matches default. Correct.

**Verdict: OK**

---

## feComposite — OK (one allowed-clamping group noted)

60 value-paths. Source: teal background with pink and orange rects; composited with a secondary layer.

- `in` / `in2` builtins — noisy renders or identical passthrough per browser support. **Allowed**.
- `operator="over"`, `"in"`, `"out"`, `"atop"`, `"xor"`, `"lighter"`, `"arithmetic"` — all seven operators produce visibly distinct outputs. Correct.
- `k1` series (0, 1, -1, 3.14, 0.5, 2) — arithmetic mode k coefficients; each value yields a different luminance/composite result. Correct.
- `k2`, `k3` series — same; all distinct. Correct.
- `k4="0"` — no additive bias; transparent where other inputs are zero. Correct.
- `k4="1"` — full white output (additive bias clamps all channels to 1). Correct.
- `k4="-1"` — near-black (negative additive bias, clamped). Correct.
- `k4="3.14"` and `k4="2"` — both produce full white (k4≥1 clamps). **Allowed**: these are correctly identical.
- `k4="0.5"` — pale/partially brightened output, distinct from k4=0 and k4=1. Correct.
- Geometry values — correct progressive spatial shifts. Correct.
- `color-interpolation-filters="auto"` — matches default. Correct.

**Verdict: OK**

---

## feMerge — OK

13 value-paths (geometry + color-interpolation-filters only; feMerge has no unique filter attributes beyond viewport placement).

- All `x`, `y` values produce correct spatial offsets of the merged layer output (red + blue overlapping rects). Each value is visibly distinct and crops/shifts as expected. Correct.
- `width="20"`, `height="20"` — correct narrow/short filter regions. Correct.
- `color-interpolation-filters="auto"` — matches default. Correct.

**Verdict: OK**

---

## Consolidated real issues

| # | Element | Issue | Severity | Fix target |
|---|---------|--------|----------|------------|
| 1 | feColorMatrix | `type="matrix"` specimen uses `values="90"` — an incomplete matrix (1 of 20 required floats). Chrome falls back to near-identity; the specimen fails to demonstrate a real color matrix transform. | Low (data authoring) | `chrome-testing/html/specimen/feColorMatrix/08-matrix.html` — replace `values="90"` with a valid 20-value matrix string, e.g. a channel-swap or tint matrix |

All other observations (clamped values rendering identically, BackgroundImage/FillPaint/StrokePaint unsupported builtins, luminanceToAlpha appearing dark against dark background) are correct, faithful, or explicitly allowed per QA brief.

**Overall batch status: 9/10 elements fully OK; 1 data-authoring fix recommended for feColorMatrix.**
