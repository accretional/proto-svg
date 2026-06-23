# Round 2 Visual QA — text-embedded batch

Elements: `text`, `tspan`, `textPath`, `image`, `foreignObject`, `view`, `script`, `style`
Date: 2026-06-22

---

## Per-element verdicts

### `<text>` — MOSTLY PERFECT, 1 remaining issue

Round-1 fix confirmed: `rotate="3.14 0.5 2"` is correctly space-separated (was `rotate="3.140.52"`). All positioning cards (`x`, `y`, `dx`, `dy`, `rotate`) render legibly with "Ag" glyph. `textLength="80"` spreads text correctly. `fill="none"` correctly invisible. `systemLanguage="en"` / `requiredExtensions=""` cards show text (no hidden-by-condition bug).

**Remaining issue:**
- Card `dy="50% 75% 3.14"` — the label is cut off at the top of the card (row 7 in screenshot col 7, card 1). The text drifts far above the SVG viewport because `50%` as a per-glyph dy offset on the first glyph at baseline y=55 pushes it to y≈27.5 on a 100-unit viewport — partially out-of-view. **Category: WEAK_EFFECT** (glyph partially clipped). The representative value for the first dy element when it is a `%` should be a small value or the baseline should position the text lower.

---

### `<tspan>` — MOSTLY PERFECT, same dy% issue

Same story as `<text>`. All positional and presentation attrs render correctly. `rotate="3.14 0.5 2"` is correct. `rotate="0 1 -1"` shows a subtle but visible tilt. Minor same-origin `dy` percentage clipping as in `<text>` (glyph partially leaves top of card on large `%` dy values).

**Remaining issue:**
- Card `dy="1 2em -1"` — the "A" glyph is near the top edge of the SVG (very large first dy offset because `1` at font-size 20 = 1 user unit, but `2em` = 40px offset which at 120px canvas moves glyph far down). Text is legible but the "g" descends below the visible area. **Category: WEAK_EFFECT** (partial clip for some dy combinations).

---

### `<textPath>` — FIXED, NOW RENDERS, 1 remaining issue

Round-1 fix confirmed: text now follows the arc path `M 5 70 Q 50 10 95 70` in every card (blueprint injects `<defs><path id="slot">` and baseline `href="#slot"`). The `path="M10 10 L90 90"` / `path="M10 50 Q50 10 90 50"` / `path="M20 20 H80 V80 H20 Z"` inline-path cards all show text on the geometry correctly.

**Remaining issue:**
- Card `href="#slot"` — the textPath element itself also carries `id="slot"` (because the baseline injects `id="slot"` for the IDREF target, and the element being varied is `href`). The element and its `<path>` reference share the same id; the textPath's `id="slot"` overrides the `<path id="slot">` in defs, so `href="#slot"` resolves to the `<textPath>` element itself (invalid). Text renders but follows an undefined path and appears as straight text near top-left. **Category: GRAMMAR_ISSUE / blueprint collision** — fix: for `href` variants on `textPath` the overlay should use `#slot` but the element should NOT receive `id="slot"` when `href` is the varying attr. Fix target: `blueprint.go` `baselineFor("textPath")` — do not add `id="slot"` to the textPath element (the `<path id="slot">` in the template is the correct target; the textPath element should not own that id).

---

### `<image>` — FIXED, NOW RENDERS, 1 remaining issue

Round-1 fix confirmed: all cards except `href="#slot"` show a dark-square + teal-circle data-URI image. `preserveAspectRatio="none"` / `"xMidYMid meet"` / `"xMinYMin slice"` cards correctly show the image distorted/fitted/sliced. `x="50%"` / `x="75%"` correctly position the image partially/fully off-canvas. `width="20"` produces a thin vertical stripe; `height="20"` a horizontal stripe — both correct.

**Remaining issue:**
- Card `href="#slot"` — shows the "broken image" icon (small image placeholder) because `href="#slot"` is a fragment reference that is not a valid raster/SVG image source. The overlay correctly leaves the href as `#slot` for IriType, but `#slot` is not a loadable image resource. **Category: EMPTY/BLANK** (broken image icon). Fix target: `overlay.go` `overlaySample()` — add a special case for `image` tag + `href`/`xlink:href` that returns `imageDataURI` instead of `"#slot"`, similar to the `feImage` baseline override already in `blueprint.go`. The same data-URI the baseline uses should be injected by the overlay when the varying attr is `href`.

---

### `<foreignObject>` — STILL BLANK (all 72 main-grid cards empty)

Round-1 supposedly fixed this, but the screenshot shows every main-grid card is completely empty (dark background, no HTML visible). The `<div>` content with `style="background:#4d8bff;..."` is present in the generated HTML source (confirmed by reading the file), but the `foreignObject` has no `width` or `height` set on it except when those are the varying attribute — and even then only a single dimension is given. Without explicit `width` AND `height` set on `<foreignObject>`, the element's content box is 0×0 and the HTML child is clipped to nothing.

**Root cause:** The `baselineFor("foreignObject")` case is absent from `blueprint.go`. The generator falls through to the default empty baseline (`return "", false`), so only the single varying attribute is set on the foreignObject element. The template blueprint (`foreignObject.html`) only wraps `{{ELEMENT}}` in a bare SVG root with no surrounding shape or positioning context.

**Category: EMPTY/BLANK** (all cards blank — regression or never truly fixed).

**Fix target:** `blueprint.go` `baselineFor()` — add a `foreignObject` case that supplies `width="90"`, `height="90"`, and `x="5"`, `y="5"` (excluding whichever one is varying):

```go
case "foreignObject":
    return add(
        [2]string{"x", "5"}, [2]string{"y", "5"},
        [2]string{"width", "90"}, [2]string{"height", "90"},
    ), false
```

This ensures every card has a visible 90×90 box; only the single card that varies `width` or `height` will deviate.

---

### `<view>` — PERFECT

All 8 cards in the main grid render the reference scene (red circle, teal rect, amber triangle). `viewBox="0 0 50 50"` zooms in. `viewBox="-10 -10 120 120"` zooms out (scene smaller). `preserveAspectRatio="none"` shows stretching. `zoomAndPan="disable"` / `"magnify"` are non-visual attributes (correctly showing the same scene). No issues.

---

### `<script>` — PARTIALLY WORKING, systematic NO_EFFECT issue

Card `type="label"` (first card) shows a bright pink/red rect — the JS executed because `type="label"` is a non-standard MIME type so Chrome ignored the script body... wait, in SVG inline scripts Chrome executes `type="text/ecmascript"` or no-type; for other MIME types it silently ignores. The first card is pink, others are gray.

On closer inspection: the **first card only** (`type="label"`) renders the rect as pink because that script **does execute** — only `type="text/javascript"` is the standard, but browsers do execute SVG `<script>` elements with no type or with `type="text/javascript"`. The value `"label"` is non-standard and Chrome in strict mode might not execute it.

Actually re-examining: the first card IS pink (rect colored by script), cards 2–13 are gray. The varying `type` attribute varies from `"label"` → `"Aa"` → `"sample"` → `"specimen"` → each is non-standard. Only the `type="label"` card appears pink.

**Wait** — looking more carefully at the screenshot, card 1 (`type="label"`) is hot pink/red; cards 2-13 are all grey. This suggests only the first card's script runs (type="label" perhaps coinciding with an older browser behavior, or the cards share `slot-target` id and the first executing script colors the rect before the others render). Each SVG is a separate inline `<svg>` in the page — their `getElementById('slot-target')` should be document-scoped and would find the FIRST rect with that id across all cards.

**Root cause of gray cards 2–13:** Multiple inline SVGs in one HTML page all have a rect with `id="slot-target"`. The scripts with non-standard types don't execute. The one with `type="label"` executes (or the baseline script without a type from an earlier card runs and colors the first rect found). **Category: NO_EFFECT_IN_MAIN_GRID** for cards where `type` is non-standard.

**Fix target:** `blueprint.go` template / `bodyFor("script")` — the JS body colors an element by `id="slot-target"`, but this id is duplicated across all SVG cards in the page. Better approach: each card SVG should carry a unique id (use the per-card index) or the coloring should be done by a `class` scoped to its own SVG parent. Alternatively, make the baseline script type always `text/javascript` for non-type-varying cards, and for type-varying cards accept that non-standard types won't execute (those cards show the baseline gray rect, which is still a valid demo of the `type` attribute effect — i.e., type controls whether the script runs). This could be documented as intended behavior. The fix is cosmetic: each SVG should contain its own uniquely-id'd target that the embedded script references, ensuring no cross-card id collision.

---

### `<style>` — PERFECT

All 12 main-grid cards show a bright pink rect with a teal ring circle — the `<style>.slot{fill:#e94560;stroke:#16c79a;stroke-width:3}</style>` is injected and styles both the `rect.slot` and `circle.slot` in the blueprint. `type="label"` / `"Aa"` / `"sample"` / `"specimen"` and all `media=` / `title=` variants render identically (the style applies regardless of non-standard type/media — browsers ignore unknown type/media and still parse). No issues.

---

## Systematic remaining issues (priority order)

### 1. `foreignObject` baseline missing width/height — CRITICAL (all 72 cards blank)

**Fix:** `blueprint.go` `baselineFor()` — add case `"foreignObject"` supplying `x="5" y="5" width="90" height="90"` (minus whichever attr is varying). This is a two-line addition.

### 2. `image href="#slot"` broken-image card — MEDIUM

**Fix:** `overlay.go` `overlaySample()` — add a guard: when `tag == "image"` and `attrName` is `"href"` or `"xlink:href"`, return `imageDataURI` (already declared as a const in `blueprint.go`; move it or re-declare in `overlay.go`). This matches the existing `feImage` treatment in `baselineFor`.

### 3. `script` cross-card id collision — LOW / cosmetic

**Fix:** `blueprint.go` template blueprint for `script` — suffix `slot-target` with a per-card unique index (e.g., pass the variant index into the blueprint). Alternatively accept that type-varying cards intentionally show the "script not executed" state (gray rect) as proof of the `type` attribute's effect.

### 4. `textPath href` card id collision — LOW

**Fix:** `blueprint.go` `baselineFor("textPath")` — when varying attr is `href`, do not add `id="slot"` to the textPath element. The `<path id="slot">` in the template blueprint is the correct IDREF target; the textPath element should not shadow it.

### 5. `text`/`tspan` large `%` dy values push glyph out of viewport — VERY LOW

**Fix:** `reps.go` `LengthPercentageType` — replace large percentage values (`50%`, `75%`) with smaller ones (`10%`, `20%`) so per-glyph dy offsets stay within the viewport. Or in `overlay.go` add a special case for `dy`/`dx` on text/tspan elements capping percentage samples to `10%`.
