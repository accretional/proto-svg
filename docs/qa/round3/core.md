# Round 3 QA — Core Galleries

**Date:** 2026-06-22  
**Scope:** rect, circle, ellipse, line, polyline, polygon, path, svg, g, defs, symbol, use, switch, a, desc, title, metadata, unknown, text, tspan, textPath, image, foreignObject, view, script, style, animate, set, animateMotion, animateTransform, mpath, discard

---

## Per-Element Verdict Table

| Element | Verdict | Notes |
|---|---|---|
| `rect` | PERFECT | All distinct cards, geometry attrs (rx/ry corners), stroke, opacity, transform all clear |
| `circle` | PERFECT | cx/cy/r distinct, fill/stroke/opacity all correct |
| `ellipse` | PERFECT | rx/ry produce clearly different shapes, stroke teal distinct from fill |
| `line` | PERFECT | All lines distinct; stroke variations, marker attrs, opacity visible |
| `polyline` | PERFECT | Point variations distinct; stroke, fill, color all show correctly |
| `polygon` | PERFECT | Fill/stroke/opacity all vary correctly across cards |
| `path` | PERFECT | Arc shapes visible; pathLength, d-variants distinct |
| `svg` | PERFECT | viewBox, preserveAspectRatio, nested SVG all visible and distinct |
| `g` | ISSUES | See below — `display="none"` card in main grid shows empty (expected for `none`, but this is the "none" variant; acceptable); `fill="none"` card correctly shows outline — OK. Actually no issue. **PERFECT** |
| `defs` | ISSUES | Main grid is completely BLANK — all 97 paths are classified as non-visual and put in the collapsed `<details>` section. The spec says `<defs>` should show a red rect via `<use>` in each card. The collapsed section does have a `<use href="#defskid"/>` pattern but there is no main-grid section at all. |
| `symbol` | PERFECT | Cards show blue rects via `<use>` of defined symbol; attrs (viewBox, preserveAspectRatio, refX/refY, x/y/width/height) all distinct |
| `use` | ISSUES | Cards show small pink circles (the `<use>` target) — but the circles are very small (~15px) and sit in the top-left corner of every card. All cards look essentially **IDENTICAL** in composition, only label differs. The visual differential is entirely in the label text, not the shape. This is WEAK_EFFECT / IDENTICAL. |
| `switch` | ISSUES | All 97 main-grid cards show identical blue squares. `<switch>` wraps a `<rect>` child — the presentation attributes on `<switch>` should be inherited by the rect, but the values shown (fill="none", fill-opacity="0.5" etc.) are all "safe defaults" that don't change the rect's own `fill="#4d8bff"`. Cards are **IDENTICAL** — no attribute produces a visible difference. |
| `a` | PERFECT | href/target attrs in first row distinct; presentation attrs match `g`-style with blue rects; xlink:href cards at bottom |
| `desc` | PERFECT (non-visual) | Empty main grid; "Non-visual attributes (28)" collapsed — correct and expected |
| `title` | PERFECT (non-visual) | Empty main grid; collapsed correctly — expected |
| `metadata` | PERFECT (non-visual) | Empty main grid; collapsed correctly — expected |
| `unknown` | PERFECT (non-visual) | Empty main grid; collapsed correctly — expected |
| `text` | PERFECT | "Ag" sample text shows with x/y/dx/dy/rotate, font-family/size/weight, fill, anchor etc. all producing distinct cards |
| `tspan` | PERFECT | Same "Ag" pattern; tspan-specific attrs (dx, dy, x, y, rotate, textLength, lengthAdjust) distinct |
| `textPath` | PERFECT | Text curves along arc paths; startOffset, method, spacing, side, href attrs all distinct and visible |
| `image` | PERFECT | Teal circle SVG renders in all cards; x/y/width/height/preserveAspectRatio/crossorigin all distinct |
| `foreignObject` | PERFECT | "HTML in SVG" blue pill visible in every card; x/y/width/height/opacity/transform all distinct |
| `view` | PERFECT | 8 main-grid cards showing circle+square+triangle scene with different viewBox and preserveAspectRatio crops |
| `script` | ISSUES | Only `type="label"` card (first) shows a red rect (script ran). The other 12 cards (`type="Aa"`, `type="sample"`, `type="specimen"`, href, crossorigin, async, defer, xlink:href, xlink:title) all show grey rects — the script didn't execute because the MIME type is not `text/javascript`. This is **expected browser behaviour** for non-JS MIME types, but visually the cards are nearly IDENTICAL (all grey). The `type="label"` card is correctly red. |
| `style` | PERFECT | All 12 cards show pink rect + teal circle (CSS applied via `<style>`); type/media/title variants all visible |
| `animate` | PERFECT | Host rect visible (blue, positioned by animate attrs); all timing/value attrs show a rect at varying positions |
| `set` | ISSUES | The `set` animation cards show the host rect jumping to different positions — good. But `repeatCount="0"`, `restart="never"`, and several other timing attrs produce cards where the rect is already animated away to the top-right corner of the card (or barely visible at top-left). This is the animation having immediately fired and the element is now in its "after" state. The cards are **not BLANK** but some are visually confusing. Marginal — acceptable given it's a timing element. |
| `animateMotion` | PERFECT | Blue square shown at various positions along motion paths; path/from/to/by/rotate/keyPoints all produce distinct cards |
| `animateTransform` | PERFECT | Blue squares with visible rotations/skews/scales; type=rotate/scale/skewX/skewY all distinct |
| `mpath` | ISSUES | Only 3 main-grid cards (href, xlink:href, xlink:title) — all show just a tiny blue square in the top-left corner. The motion target is barely visible (very small rect, no motion context). **WEAK_EFFECT** — the visual is nearly the same tiny square in all 3 cards |
| `discard` | ISSUES | Only 2 main-grid cards (href, begin). Both show a small blue square at top-left — same WEAK_EFFECT as mpath. Cards are essentially identical. |

---

## Consolidated Remaining Issues

### 1. `defs` — BLANK main grid
**Problem:** The generator classifies ALL defs attributes as non-visual and puts them in the collapsed `<details>` section. The main `<div class="grid">` is empty (line 16–17 of defs.html: `<div class="grid">\n</div>`).  
**Expected:** The spec note says defs should show a red rect via `<use>`. There should be at least one main-grid card demonstrating that content inside `<defs>` is usable via `<use>`.  
**Fix target:** `chrome-testing/cmd/gen/reps.go` (or `blueprint.go`) — the defs element's visual-attr classification routes everything to meta. The defs overlay should generate a "baseline" card in the main grid showing `<defs><rect id="x" .../></defs><use href="#x"/>` pattern, not in the collapsed section. Alternatively, add a hardcoded "demo" card at the top of the defs main grid in `emit.go` / the template.

### 2. `use` — IDENTICAL cards (WEAK_EFFECT)
**Problem:** All use-element cards show an identical tiny pink circle (~15×15px) in the top-left corner. The `<use>` target is a small `<circle>` and none of the presentation attrs on `<use>` change its size or position noticeably at card scale.  
**Fix target:** `chrome-testing/cmd/gen/overlay.go` or `reps.go` — the `<use>` host shape should be larger (e.g. `cx="50" cy="50" r="35"` instead of a tiny default). Additionally, the x/y/width/height attrs on `<use>` should visibly shift or scale the instance; the baseline circle should be centered and large enough to show changes.

### 3. `switch` — IDENTICAL cards (NO_EFFECT_IN_MAIN_GRID)
**Problem:** All 97 `<switch>` main-grid cards show the same blue square. `<switch>` is a container and the presentation attrs are all inherited "pass-through" values that don't change the child rect's own hardcoded `fill="#4d8bff"`. Every card looks the same.  
**Fix target:** `chrome-testing/cmd/gen/reps.go` or `overlay.go` — the `<switch>` child element should NOT hardcode `fill="#4d8bff"`. Instead it should use `fill="currentColor"` or `fill="inherit"` so that presentation attrs set on `<switch>` (like `color="#e94560"` or `fill="tomato"`) actually propagate and produce distinct visual output. Or, use the same "group host" approach as `<g>` (which does show distinct cards).

### 4. `script` — 12/13 cards IDENTICAL grey (non-JS type blocks execution)
**Problem:** Script with `type="Aa"`, `type="sample"`, `type="specimen"` etc. doesn't execute (browser only runs `text/javascript` or omitted type). Only the `type="label"` card (the very first) executes — probably because `"label"` is treated as unknown/unrecognised and Chrome still runs it, or it was emitted before the type attr was parsed. The remaining 12 cards all show the grey default rect.  
**Fix target:** `chrome-testing/cmd/gen/reps.go` / `overlay.go` for `<script>` — the card SVG should not rely on script execution for the visual. Instead, show the rect already coloured, and the `<script>` tag is included for attribute demonstration purposes only (label change, not colour change). Alternatively use `type="text/javascript"` as the canonical executable type on the first card only, and drop the script colour-change entirely, showing a "code" label overlay instead.

### 5. `mpath` — WEAK_EFFECT (tiny identical squares)
**Problem:** Only 3 visual cards, all showing a tiny blue square at top-left. `<mpath>` is a child of `<animateMotion>` and provides a motion path — on its own there is no distinct visual per attr.  
**Fix target:** `chrome-testing/cmd/gen/overlay.go` — the `<mpath>` demo SVG should include a visible host shape (e.g. a circle of radius 20 centered at a starting point) with an `<animateMotion>` parent that uses the `<mpath>` child. The current tiny square is too small. The three distinct attrs (href, xlink:href, xlink:title) should all show the same scene but labeled differently — that's structurally acceptable for a 3-attr element, but the shape must be visible.

### 6. `discard` — WEAK_EFFECT (tiny identical squares)
**Problem:** Same as mpath — only 2 visual cards, both showing a tiny blue square. `<discard>` removes elements after a delay; the static snapshot shows either the element present (before discard fires) or absent. The screenshot is taken before the animation runs so both cards look the same.  
**Fix target:** `chrome-testing/cmd/gen/overlay.go` — the `<discard>` demo should show a larger, clearly labeled host shape that demonstrates "this element will be discarded". The shape should be bigger and ideally in a different color. For the `begin="0s"` card the shape may already be gone (correct), but the `href` card should show it present.

---

## Summary

**PERFECT (26/31):** rect, circle, ellipse, line, polyline, polygon, path, svg, g, symbol, a, desc, title, metadata, unknown, text, tspan, textPath, image, foreignObject, view, style, animate, animateMotion, animateTransform, set (marginal/acceptable)

**ISSUES (5/31):** defs (BLANK main grid), use (IDENTICAL tiny circles), switch (IDENTICAL — no effect), mpath (WEAK_EFFECT tiny squares), discard (WEAK_EFFECT tiny squares)

**NEAR-ISSUE (1/31):** script (only first card executes; rest grey — by-design browser behaviour but visually weak)
