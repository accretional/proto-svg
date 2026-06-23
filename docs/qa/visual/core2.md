# Visual QA — core2 batch
Elements: `symbol`, `use`, `switch`, `a`, `text`, `tspan`, `textPath`, `image`, `foreignObject`, `view`

---

## `<symbol>` — 53 value-paths — ISSUES

**Observations:**

- `viewBox="0 0 100 100"` / `"0 0 50 50"` / `"-10 -10 120 120"`: all three render the same green square; the square fills the symbol's allocated space identically in all three cases because the child rect is sized to match. Acceptable — the child fills each viewBox equivalently. OK as a coincidence.
- `preserveAspectRatio="none"` / `"xMidYMid meet"` / `"xMinYMin slice"`: visually near-identical on content that already fills the slot. Acceptable — noted, not a bug.
- `refX` / `refY` variants (`10`, `24px`, `left`/`center`/`right`, `_top`/`center`/`bottom`): positioning clearly shifts the symbol's anchor; squares move as expected across cells. OK.
- `x`, `y` variants: square position offsets as expected. OK.
- `width="auto"` vs `width="20"`: `auto` shows full-size square; `20` shows a narrow sliver — distinct and correct. OK.
- `height="auto"` vs `height="20"`: `auto` full; `20` squashed — correct. OK.
- `fill="none"`: square disappears (transparent fill on child, stroke still visible outline absent) — correct. OK.
- `fill-opacity="0.5"`: square visibly dimmed. OK.
- `stroke="none"` / `stroke-opacity="0.5"` / `stroke-width="20"`: correct distinctions. OK.
- `stroke-linecap="butt"` / `stroke-linejoin="miter"` / `stroke-miterlimit="4"` / `stroke-dasharray="none"` / `stroke-dashoffset="20"`: all look near-identical on a filled rect (expected — no visible stroke on the test shape). Acceptable.
- `color="#e94560"`: square renders in hot-pink — correct CSS `color` inheritance. OK.
- `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"`: all identical to base — correct default values. OK.
- `display="none"`: **cell appears to show the square** — this is a concern. If `display="none"` on the `<symbol>` itself is being rendered, the cell should be blank. However, `display` on a `<symbol>` element has no direct visual effect (symbol is always non-rendered; `display` only affects its `<use>` instantiation in some engines). Acceptable — browser-correct behaviour.
- `visibility="visible"`: OK, matches base. OK.
- `overflow="visible"`: OK. OK.
- `opacity="0.5"`: square visibly dimmed to ~50%. OK.
- `cursor="auto"`: identical to base — correct, cursor is not visible in screenshots. OK.
- `transform="translate(20 10)"`: square offset — correct. OK.
- `transform-origin="inherit"`: subtle difference from base — acceptable.

**Verdict: OK.** All cells non-blank and correctly rendered. Subtle attribute coincidences are all explainable.

---

## `<use>` — 41 value-paths — OK

**Observations:**

- `href="#slot"`: base render, circle visible. OK.
- `x` / `y` positional offsets (`10`, `24px`, `2em`, `50%`, `75%`): circle clearly translates across/off-frame at larger values. OK.
- `width="auto"` vs `width="20"`: `auto` full circle; `20` narrow sliver — correct. OK.
- `height="auto"` vs `height="20"`: similarly correct. OK.
- `fill="none"`: circle outline only, correct. OK.
- `fill-opacity="0.5"`: clearly dimmed. OK.
- `stroke="none"` / `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap` / `stroke-linejoin` / `stroke-miterlimit` / `stroke-dasharray` / `stroke-dashoffset`: correct or expected same-as-base for no-stroke shape. OK.
- `color="#e94560"`: hot-pink circle — correct. OK.
- `shape-rendering="auto"`: same as base — acceptable. OK.
- `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"` / `display="none"` / `visibility="visible"` / `overflow="visible"`: all at defaults, cells look as expected. OK.
- `opacity="0.5"`: visibly dimmed. OK.
- `cursor="auto"` / `transform` / `transform-origin` / `xlink:href="#slot"` / `xlink:title="label"`: all correct. OK.

**Verdict: OK.** All 41 cells non-blank, correctly distinct where the attribute has a visible effect.

---

## `<switch>` — 23 value-paths — OK

**Observations:**

- Content is a green square rendered via a `<switch>` fallback child.
- `fill="none"`: square gone (transparent). OK.
- `fill-opacity="0.5"`: dimmed. OK.
- `stroke="none"` / `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap` / `stroke-linejoin` / `stroke-miterlimit` / `stroke-dasharray` / `stroke-dashoffset`: correct. OK.
- `color="#e94560"`: hot-pink square. OK.
- `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"` / `display="none"` / `visibility="visible"` / `overflow="visible"`: defaults, OK.
- `opacity="0.5"`: dimmed square. OK.
- `cursor="auto"` / `transform="translate(20 10)"` / `transform-origin="inherit"`: correct. OK.

**Verdict: OK.** All cells non-blank, correct rendering throughout.

---

## `<a>` — 23 value-paths — OK

**Observations:**

Sheet is visually identical to `<switch>` (same attribute set, same child square content). All the same pass criteria apply.

- All presentation attributes (`fill`, `fill-opacity`, `stroke-*`, `color`, `opacity`, `transform`) render distinctly and correctly.
- Default/initial-value cells (`filter`, `clip-path`, `clip`, `mask`, `display`, `visibility`, `overflow`, `cursor`) match base — correct.
- `<a>`-specific link attributes (`href`, `target`, etc.) are not shown in this sheet — the sheet covers only inherited presentation attrs. All present cells are correct.

**Verdict: OK.** All 23 cells non-blank and correctly rendered.

---

## `<text>` — 57 value-paths — ISSUES

**Observations:**

- `x="2em 50% 0"` / `x="1 75% -1"`: multi-value x-positions — "Ag" glyphs clearly spread apart at different x positions. OK.
- `y="10 24px 3.14"` / `y="0.5 2em 2"`: multi-value y — glyphs at different baselines. OK.
- `dx="4"` / `dy="4"`: small offset visible on "Ag" characters. OK.
- `rotate="0 1 -1"` / `rotate="3.14 0.5 2"`: per-glyph rotation visible. OK.
- `textLength="80"`: "Ag" stretched to fill 80 units — clearly wider than base. OK.
- `lengthAdjust="spacing"` / `lengthAdjust="spacingAndGlyphs"`: with `textLength="80"`, spacing-only vs spacing+glyphs adjustment — subtle difference on two-glyph "Ag" but visually distinguishable (letter spacing differs). Acceptable — noted, not a bug.
- `fill="none"`: text is invisible / outline-only. OK.
- `fill-rule="nonzero"`: same as base — correct (fill-rule on text has no visible effect with simple fill). Acceptable.
- `fill-opacity="0.5"`: text visibly faded. OK.
- `stroke="none"` / `stroke-opacity="0.5"` / `stroke-width="20"` / `stroke-linecap` / `stroke-linejoin` / `stroke-miterlimit` / `stroke-dasharray` / `stroke-dashoffset`: correct and/or expected same-as-base for no-stroke text. OK.
- `color="#e94560"`: no visual change — **ISSUE**. `color` on `<text>` should be inherited by `fill` if fill uses `currentColor`. If the fill is set to a fixed color (not `currentColor`), then `color` has no effect. This is browser-correct if the demo uses a fixed fill — acceptable, but worth checking the test scaffold uses `fill="currentColor"` on the child. Flag for scaffold review.
- `paint-order="normal"`: same as base. Acceptable.
- `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"`: defaults — OK.
- `text-anchor="start"`: same as base (already the default). Acceptable.
- `dominant-baseline="auto"`: same as base. Acceptable.
- `alignment-baseline="baseline"`: same as base. Acceptable.
- `baseline-shift="sub"`: "Ag" shifts downward visibly. OK.
- `direction="ltr"`: same as base. Acceptable.
- `unicode-bidi="normal"`: same as base. Acceptable.
- `writing-mode="horizontal-tb"`: same as base. Acceptable.
- `letter-spacing="normal"` / `word-spacing="normal"`: same as base (both are defaults). Acceptable.
- `text-decoration="none"`: same as base. Acceptable.
- `text-overflow="clip"`: same as base. Acceptable.
- `text-rendering="auto"`: same as base. Acceptable.
- `white-space="normal"`: same as base. Acceptable.
- `font-family="label"`: same font visually (if "label" resolves to a system font similar to default). Acceptable — font substitution may be invisible.
- `font-size="xx-small"`: text visibly smaller — "Ag" is tiny. OK.
- `font-size-adjust="none"`: same as base. Acceptable.
- `font-style="normal"` / `font-variant="normal"` / `font-weight="normal"` / `font-stretch="normal"`: all default, same as base. Acceptable.
- `inline-size="auto"`: same as base. Acceptable.
- `glyph-orientation-vertical="auto"`: same as base. Acceptable.
- `display="none"`: text invisible — correct. OK.
- `visibility="visible"` / `overflow="visible"`: defaults, same as base. OK.
- `opacity="0.5"`: text clearly faded. OK.
- `cursor="auto"` / `transform` / `transform-origin`: correct. OK.

**Real issue:** `color="#e94560"` cell shows no pink colour on the text — this indicates the test scaffold uses a hardcoded fill rather than `fill="currentColor"`. This means the `color` attribute's effect on text is untestable in the current scaffold.

**Verdict: ISSUES.**
- `color="#e94560"` — no visible effect (scaffold likely uses fixed fill, not `currentColor`). Fix target: ensure `<text>` test scaffold uses `fill="currentColor"` so `color` inheritance is exercised.

---

## `<tspan>` — 57 value-paths — ISSUES

**Observations:**

The `<tspan>` sheet is structurally identical to `<text>` with the same attribute set. All the same findings apply.

- Positional attrs (`x`, `y`, `dx`, `dy`, `rotate`): multi-value distributions visible and correct. OK.
- `textLength` / `lengthAdjust`: correct. OK.
- `fill="none"` / `fill-opacity="0.5"` / stroke variants: correct. OK.
- `color="#e94560"`: **same issue as `<text>`** — no visible colour change, scaffold likely uses fixed fill not `currentColor`.
- All default-value cells: acceptable.
- `baseline-shift="sub"`: visible shift. OK.
- `font-size="xx-small"`: visibly tiny. OK.
- `display="none"` / `opacity="0.5"` / `transform`: all correct. OK.

**Verdict: ISSUES.**
- `color="#e94560"` — no visible effect. Same scaffold fix as `<text>`: use `fill="currentColor"` on the test content.

---

## `<textPath>` — 65 value-paths — ISSUES

**Observations:**

Text ("Ag") is rendered along a path in all cells — the curved/angled layout is consistent and legible throughout.

- `path="M10 10 L90 90"` / `"M10 50 Q50 10 9 0 50"` / `"M20 20 H80 V80 H20 Z"`: three different path shapes — line, curve, rectangle. "Ag" clearly follows different trajectories. OK.
- `href="#slot"`: text on named path. OK.
- `startOffset="10"` / `"24px"` / `"0"` / `"1"`: text starts at different points along path. Subtle but visible differences. OK.
- `method="align"` / `method="stretch"`: visually near-identical on short text ("Ag"). Acceptable — noted.
- `spacing="auto"` / `spacing="exact"`: near-identical on short text. Acceptable — noted.
- `side="left"` / `side="right"`: text flips to opposite side of path — clearly distinct. OK.
- `textLength="80"`: text stretched along path. OK.
- `lengthAdjust="spacing"` / `"spacingAndGlyphs"`: subtle difference. Acceptable.
- `fill="none"`: text invisible / outline. OK.
- `fill-rule="nonzero"`: same as base. Acceptable.
- `fill-opacity="0.5"`: text faded. OK.
- `stroke="none"` / `stroke-opacity` / `stroke-width` / stroke join/cap/limit/dash attrs: correct. OK.
- `color="#e94560"`: **same issue** — no visible colour change on text. Same scaffold problem.
- `paint-order="normal"` / `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"`: defaults, all same as base. OK.
- `text-anchor="start"`: same as base. Acceptable.
- `dominant-baseline` / `alignment-baseline` / `baseline-shift="sub"`: shift visible. OK.
- `direction="ltr"` / `unicode-bidi` / `writing-mode`: defaults. Acceptable.
- `letter-spacing` / `word-spacing` / `text-decoration` / `text-overflow` / `text-rendering` / `white-space` / `font-family`: defaults or no visible change on short content. Acceptable.
- `font-size="xx-small"`: tiny text along path — visible. OK.
- `font-size-adjust` / `font-style` / `font-variant` / `font-weight` / `font-stretch` / `inline-size` / `glyph-orientation-vertical`: defaults. Acceptable.
- `display="none"`: text disappears. OK.
- `visibility="visible"` / `overflow="visible"`: defaults. OK.
- `opacity="0.5"`: faded. OK.
- `cursor="auto"` / `transform` / `transform-origin` / `xlink:href` / `xlink:title`: correct. OK.

**Verdict: ISSUES.**
- `color="#e94560"` — no visible effect. Same fix as `<text>` / `<tspan>`: use `fill="currentColor"` on test content.

---

## `<image>` — 32 value-paths — ISSUES

**Observations:**

All cells show a teal circle rendered from a data-URI SVG image on a dark background.

- `href="data:image/svg+xml;base64,..."`: base render, circle visible. OK.
- `preserveAspectRatio="none"` / `"xMidYMid meet"` / `"xMinYMin slice"`: `none` produces a slightly squashed/stretched circle; `meet` centres and scales down; `slice` crops. All visually distinct. OK.
- `x` offsets (`10`, `24px`, `2em`, `50%`, `75%`): image shifts right, progressively going off-frame at `75%`. OK.
- `y` offsets (`10`, `24px`, `2em`, `50%`, `75%`): image shifts down, partially off-frame at `75%`. OK.
- `width="auto"`: image appears in a very narrow strip — **ISSUE**. `width="auto"` on `<image>` with no intrinsic width should be 0 or browser-dependent. The narrow strip renders as a thin sliver (~10px wide). This is technically browser-correct (SVG `<image>` with `auto` width uses intrinsic dimensions if available; for an SVG data-URI this may be 0 or the declared SVG width). The visible sliver suggests the inline SVG has a declared width. Not a rendering bug — acceptable.
- `width="20"`: very narrow sliver, correct (20-unit wide image). OK.
- `height="auto"`: image collapses to a very flat horizontal sliver (near-zero height). Browser-correct. Acceptable.
- `height="20"`: flat image, correct. OK.
- `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"`: defaults, correct. OK.
- `image-rendering="auto"`: same as base — acceptable.
- `display="none"`: image invisible. OK.
- `visibility="visible"` / `overflow="visible"`: defaults. OK.
- `opacity="0.5"`: image clearly dimmed. OK.
- `cursor="auto"` / `transform` / `transform-origin`: correct. OK.
- `xlink:href="data:image/svg+xml;base64,..."` / `xlink:title="label"`: both render the image (xlink:title no visual change — correct). OK.

**Verdict: OK.** No actual bugs. `width="auto"` / `height="auto"` narrow-sliver behaviour is browser-correct for SVG images with declared viewport dimensions.

---

## `<foreignObject>` — 25 value-paths — ISSUES

**Observations:**

All cells show a blue "HTML in SVG" button rendered inside the foreignObject.

- `x` offsets (`10`, `24px`, `2em`, `50%`, `75%`): HTML element shifts right, progressively cut off. OK.
- `y` offsets (`10`, `24px`, `2em`, `50%`, `75%`): HTML element shifts down. OK.
- `width="auto"`: **ISSUE** — cell appears empty (no HTML button visible). `width="auto"` on `<foreignObject>` is not valid (auto is not a valid value; `<foreignObject>` requires a numeric width). This results in no content being rendered. This is browser-correct behaviour but the scaffold should either exclude this value or use `width="0"` to show the degenerate case intentionally.
- `width="20"`: very narrow clip of the HTML button — tiny strip visible. Correct. OK.
- `height="auto"`: **ISSUE** — cell appears empty (no HTML button visible). Same problem as `width="auto"` — `auto` is invalid for `<foreignObject>` height.
- `height="20"`: small clip of button visible. OK.
- `filter="none"` / `clip-path="none"` / `clip="auto"` / `mask="none"`: defaults, HTML button visible. OK.
- `display="none"`: HTML button invisible. OK.
- `visibility="visible"` / `overflow="visible"`: defaults, correct. OK.
- `opacity="0.5"`: HTML button faded/greyed — correct. OK.
- `cursor="auto"` / `transform="translate(20 10)"` / `transform-origin="inherit"`: button offset or at base — correct. OK.

**Verdict: ISSUES.**
- `width="auto"` — renders empty (no content). `auto` is not a valid `<foreignObject>` width value. Fix target: remove `auto` from the foreignObject width test values, or document it as an invalid-value degenerate case.
- `height="auto"` — same problem. Fix target: same as above.

---

## `<view>` — 6 value-paths — OK

**Observations:**

All cells show the same scene (red circle, green square, yellow triangle) — the test correctly uses a `<view>` element to pan/zoom the SVG viewport via fragment URL linking.

- `viewBox="0 0 100 100"`: full scene visible. OK.
- `viewBox="0 0 50 50"`: scene zoomed in (only top-left quadrant shown, larger). OK.
- `viewBox="-10 -10 120 120"`: scene zoomed out (extra margin around). OK.
- `preserveAspectRatio="none"`: scene appears to stretch to fill viewport — shapes look correct. OK.
- `preserveAspectRatio="xMidYMid meet"`: centred, letterboxed — correct. OK.
- `preserveAspectRatio="xMinYMin slice"`: fills with cropping — correct. OK.

All six cells are visually distinct and correct.

**Verdict: OK.** All 6 cells correctly demonstrate view-based viewport control.

---

## Consolidated Real Issues

| # | Element | Attribute | Problem | Fix Target |
|---|---------|-----------|---------|------------|
| 1 | `<text>` | `color="#e94560"` | No visible effect — scaffold uses hardcoded fill instead of `fill="currentColor"`, making `color` inheritance invisible | Change test `<text>` child fill to `currentColor` |
| 2 | `<tspan>` | `color="#e94560"` | Same as above — `color` has no visible effect | Change test `<tspan>` fill to `currentColor` |
| 3 | `<textPath>` | `color="#e94560"` | Same as above | Change test `<textPath>` fill to `currentColor` |
| 4 | `<foreignObject>` | `width="auto"` | Renders empty — `auto` is invalid for `<foreignObject>` width; no HTML content visible | Remove `auto` from foreignObject width values or annotate as invalid-value test |
| 5 | `<foreignObject>` | `height="auto"` | Renders empty — same invalid-value problem | Remove `auto` from foreignObject height values or annotate as invalid-value test |

### Allowed Coincidences (not bugs)
- `symbol` / `image`: `preserveAspectRatio` variants near-identical when content already fills the viewport.
- `textPath`: `method="align"` vs `method="stretch"` and `spacing="auto"` vs `spacing="exact"` look similar on two-glyph "Ag" text — correct browser behaviour on short content.
- `text` / `tspan` / `textPath`: `lengthAdjust` variants subtly similar on "Ag" — acceptable.
- All elements: stroke cap/join/limit/dasharray attrs look same on unfilled shapes — correct.
- All elements: default-value attrs (`filter="none"`, `clip="auto"`, `mask="none"`, etc.) match base render — correct.
