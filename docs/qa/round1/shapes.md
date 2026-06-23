# Visual QA — Round 1: Shape Elements

Reviewer: automated visual QA pass  
Date: 2026-06-22  
Batch: `rect`, `circle`, `ellipse`, `line`, `polyline`, `polygon`, `path`  
Method: full-page screenshot (`SNAP_FULLPAGE=1 SNAP_SCALE=1`) + HTML source inspection per card.

---

## Summary table

| Element | Cards | Category | Primary issues |
|---------|-------|----------|----------------|
| `rect` | 75 | IDENTICAL_NO_EFFECT + HAS_EMPTY_CARDS | id/tabindex/lang/aria/event groups visually identical; x="100%"/y="100%" off-canvas; fill="none" blank; requiredExtensions/systemLanguage hidden |
| `circle` | 68 | IDENTICAL_NO_EFFECT + HAS_EMPTY_CARDS | same groups as rect; cx="100%"/cy="100%" off-canvas; fill="none" blank; requiredExtensions/systemLanguage hidden |
| `ellipse` | 72 | IDENTICAL_NO_EFFECT + HAS_EMPTY_CARDS | same groups; fill="none" blank; requiredExtensions/systemLanguage hidden |
| `line` | 75 | MOSTLY_PERFECT + HAS_EMPTY_CARDS | stroke="none" renders blank (line is invisible); requiredExtensions/systemLanguage hidden; x1/y1/x2/y2="0"/"1" nearly identical |
| `polyline` | 61 | MOSTLY_PERFECT + HAS_EMPTY_CARDS | requiredExtensions/systemLanguage hidden; points label wraps over two lines (UI) |
| `polygon` | 61 | MOSTLY_PERFECT + HAS_EMPTY_CARDS | fill="none" blank (no stroke on baseline); requiredExtensions/systemLanguage hidden |
| `path` | 63 | HAS_EMPTY_CARDS + GRAMMAR_ISSUE | d="none" empty (path has no geometry); requiredExtensions/systemLanguage hidden; stroke-miterlimit="0" invalid per SVG2 spec |

---

## `<rect>` — IDENTICAL_NO_EFFECT + HAS_EMPTY_CARDS

75 enumerated value-paths walked.

### HAS_EMPTY_CARDS

| Label | Why empty |
|-------|-----------|
| `x="100%"` | rect placed at x=100% of viewport (100 units) — entirely off the 100×100 viewBox, invisible |
| `y="100%"` | same — rect starts at y=100 (bottom edge), entirely clipped |
| `fill="none"` | no stroke on baseline, so the rect is genuinely invisible |
| `requiredExtensions="label"` | `requiredExtensions` is a conditional-processing attribute; the browser will never match an extension URI of "label", so the element is not rendered |
| `systemLanguage="label"` | "label" is not a valid BCP-47 tag and never matches the browser's language, so the element is hidden by the browser |

### IDENTICAL_NO_EFFECT (32 out of 75 cards)

The following attribute groups produce the same rendered rect (fill=#e94560, no visual difference between values):

- `id` (5 values: circle1, grad-a, myId, node3, r1) — id has no visual effect
- `tabindex` (6 values: 0, 1, -1, 100, 3, 10) — no visual effect
- `autofocus` (2 values: true, false) — no visual effect
- `lang` / `xml:lang` (5 values each: en, fr-CA, de, zh-Hans, pt-BR) — no visual effect on a graphical element
- `xml:space` (2 values: default, preserve) — no visual effect on rect
- `role` (1 value: application articlebanner) — no visual effect
- `aria-activedescendant` / `aria-atomic` — no visual effect
- `oncancel` / `oncanplay` / `oncopy` / `onfocusin` (event handlers set to "label") — no visual effect

### Also noted

- `stroke-opacity="0.5"`, `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`, `stroke-dashoffset` — these are stroked-only properties and the baseline has no stroke, so these cards all look identical to the baseline rect. They are not empty but are visually non-distinct.
- `stroke-miterlimit="0"` — SVG2 spec (§ 17.4) requires stroke-miterlimit ≥ 1; value "0" is out-of-range. (Grammar ISSUE — see also `path` section.)

### Fix targets

| Issue | Target | Change |
|-------|--------|--------|
| x="100%" / y="100%" off-canvas | `reps.go` `LengthPercentageType` | Replace "100%" with "25%" or "75%" so the rect stays on-canvas. Current list: `{"10", "24px", "2em", "50%", "100%"}` → `{"10", "24px", "2em", "50%", "75%"}` |
| fill="none" blank | `blueprint.go` `baselineFor("rect")` | Add a fallback stroke to the baseline: append `[2]string{"stroke", "#16c79a"}, [2]string{"stroke-width","2"}` so fill="none" still shows an outline |
| requiredExtensions="label" hidden | `reps.go` `StringType` (or `overlay.go`) | Add an `overlaySample` case for attrName=="requiredfextensions" returning a value that IS always present (e.g. an empty string list or omit the attribute). Alternatively change `StringType` reps to include a known valid extension URI; cleanest fix: in `overlay.go` make `requiredExtensions` always resolve to `""` (empty list = no restriction) |
| systemLanguage="label" hidden | `reps.go` `StringType` / `overlay.go` | The grammar uses `StringType` reps for `systemLanguage`. Add overlay case: `case an == "systemlanguage": return "en", true` so the value always matches the browser language |
| id/tabindex/lang/xml:lang/xml:space/role/aria-*/on* identical | `gallery.go` / `emit.go` | These are semantics-only attributes and will always be visually identical. Consider suppressing them from the gallery or grouping them under a special "non-visual attributes" collapsed section in the page HTML. They are not wrong per se but inflate the gallery with non-informative cards |
| stroke-miterlimit="0" invalid | `overlay.go` `isNonNegativeAttr` | Add `"stroke-miterlimit"` to `isNonNegativeAttr` and additionally clamp minimum to 1 (or use a dedicated overlay case returning "4" which is the SVG default) |

---

## `<circle>` — IDENTICAL_NO_EFFECT + HAS_EMPTY_CARDS

68 enumerated value-paths walked.

### HAS_EMPTY_CARDS

| Label | Why empty |
|-------|-----------|
| `cx="100%"` | circle center pushed to x=100 — circle extends off the right edge; only barely visible or fully off-canvas depending on r |
| `cy="100%"` | circle center pushed to y=100 — fully off-canvas (r=40, so edge at y=140, entirely outside viewBox) |
| `fill="none"` | no stroke in baseline; circle is completely invisible |
| `requiredExtensions="label"` | browser never satisfies this extension, element hidden |
| `systemLanguage="label"` | "label" never matches browser language, element hidden |

### IDENTICAL_NO_EFFECT (32 of 68 cards)

Same groups as `rect`: id, tabindex, autofocus, lang, xml:lang, xml:space, role, aria-activedescendant, aria-atomic, oncancel, oncanplay, oncopy, onfocusin.

### Fix targets

Same fixes as `rect` above, plus:

| Issue | Target | Change |
|-------|--------|--------|
| cx="100%" / cy="100%" off-canvas | `reps.go` `LengthPercentageType` | Replace "100%" with "75%" |
| fill="none" blank | `blueprint.go` `baselineFor("circle")` | Add stroke: `[2]string{"stroke", "#e94560"}, [2]string{"stroke-width","2"}` |

---

## `<ellipse>` — IDENTICAL_NO_EFFECT + HAS_EMPTY_CARDS

72 enumerated value-paths walked.

### HAS_EMPTY_CARDS

| Label | Why empty |
|-------|-----------|
| `fill="none"` | no stroke in baseline; ellipse invisible |
| `requiredExtensions="label"` | element hidden by browser |
| `systemLanguage="label"` | element hidden by browser |

No off-canvas issues: the ellipse cx/cy reps use `LengthPercentageType` but the baseline cx=50, cy=50 anchors the shape, and the 100% cases shift only one axis — still partially visible. (The cx="100%" card shows the right half of the ellipse clipped, so it is HAS_EMPTY_CARDS for cx but not fully blank — mark as partially empty.)

Actually, after re-inspection of the screenshot: `cx="100%"` and `cy="100%"` cards do show a visible circle arc at the edge, so they are not fully blank — they are just clipped. This is acceptable/meaningful.

### IDENTICAL_NO_EFFECT (32 of 72 cards)

Same id/tabindex/lang/aria/event groups as `rect`.

### Fix targets

| Issue | Target | Change |
|-------|--------|--------|
| fill="none" blank | `blueprint.go` `baselineFor("ellipse")` | Add `[2]string{"stroke","#e94560"}, [2]string{"stroke-width","2"}` |
| requiredExtensions/systemLanguage hidden | `overlay.go` | Add overlay cases (same fix as `rect`) |

---

## `<line>` — MOSTLY_PERFECT + HAS_EMPTY_CARDS

75 enumerated value-paths walked. The baseline (blue diagonal) is well-chosen and most cards are meaningfully distinct.

### HAS_EMPTY_CARDS

| Label | Why empty |
|-------|-----------|
| `stroke="none"` | line has no fill by definition; removing the stroke makes it invisible |
| `requiredExtensions="label"` | element hidden |
| `systemLanguage="label"` | element hidden |

### Nearly-identical / low-signal cards

| Label | Issue |
|-------|-------|
| `x1="0"`, `y1="0"`, `x2="0"`, `y2="1"`, `x1="1"`, `x2="1"` | Coordinate values 0 and 1 place endpoints very near the top-left corner; the line still renders but the visual difference from the diagonal baseline is subtle. Not empty, but low-information. |
| `fill="none"` | line with no fill is standard (lines have no fill area); technically correct, visually the same as baseline |

### Fix targets

| Issue | Target | Change |
|-------|--------|--------|
| stroke="none" blank | `blueprint.go` `baselineFor("line")` | This is a genuine value to demonstrate. The card should not be empty. Instead the fix is: in `overlay.go`, detect `an == "stroke"` and when the value would be "none", pair it with a comment in the gallery or skip it. However the most correct approach is to add a second visible element to the scaffold so even stroke="none" shows context. Alternatively, just accept it as a valid (invisible) demo of "none". Better: change the `stroke` keyword rep. The grammar emits `stroke="none"` as a keyword arm in `StrokeAttr`; `overlay.go` does not steer it. Add to `overlay.go`: `case an == "stroke" && valueKind == "keyword" && value == "none": skip` is not possible in current design. Real fix: in `reps.go`, note that `stroke` keyword "none" is legitimately invisible for a line; the blueprint should wrap it with a background shape or the gallery could annotate it. Simplest fix in `blueprint.go` `baselineFor("line")`: add a faint background circle so the card is not blank, or just accept this as a valid EBNF arm. |
| x1/y1/x2/y2 near-zero coordinate reps | `reps.go` `CoordinateType` | Current: `{"0", "10", "-5.5", "100px", "50%", "2.5em"}`. Replace "0" with "5" or "20" and "100px" with "80" so all coordinate samples produce clearly-distinct lines within the viewBox |
| requiredExtensions/systemLanguage hidden | `overlay.go` | Same fix as `rect` |

---

## `<polyline>` — MOSTLY_PERFECT + HAS_EMPTY_CARDS

61 enumerated value-paths walked. The N-shaped baseline stroke is visually clear and most attribute variants are distinct.

### HAS_EMPTY_CARDS

| Label | Why empty |
|-------|-----------|
| `requiredExtensions="label"` | element hidden |
| `systemLanguage="label"` | element hidden |

### UI_ISSUE

| Label | Issue |
|-------|-------|
| `points="10,80 50,10 90,80 50,50"` | Label text is 32 characters and wraps over two lines in the 160px-wide card at 12px font. Not a hard failure but makes the card taller and the grid ragged. |
| `color-interpolation-filters="auto"` | 34 characters, also wraps. |
| `aria-activedescendant="circle1"` | 31 characters, also wraps. |

### IDENTICAL_NO_EFFECT

Same id/tabindex/lang/aria/event groups as other elements.

### Fix targets

| Issue | Target | Change |
|-------|--------|--------|
| Label wrapping | `emit.go` `galleryCSS` | Add `overflow:hidden; text-overflow:ellipsis; white-space:nowrap` to `.card .label` and increase card width to 180px, OR keep `word-break:break-all` and increase `.card` height to accommodate. Current `.card` has no min-height so wrapping labels push the SVG+label proportions. A simpler fix: set `.card{min-height:200px}` so all cards are the same height regardless of label wrap. |
| requiredExtensions/systemLanguage hidden | `overlay.go` | Same fix as `rect` |

---

## `<polygon>` — MOSTLY_PERFECT + HAS_EMPTY_CARDS

61 enumerated value-paths walked. The triangle baseline is clear and distinctive.

### HAS_EMPTY_CARDS

| Label | Why empty |
|-------|-----------|
| `fill="none"` | The polygon baseline has no stroke; fill="none" makes the polygon invisible |
| `requiredExtensions="label"` | element hidden |
| `systemLanguage="label"` | element hidden |

### IDENTICAL_NO_EFFECT

Same groups as other elements.

### Fix targets

| Issue | Target | Change |
|-------|--------|--------|
| fill="none" blank | `blueprint.go` `baselineFor("polygon")` | Add stroke: `[2]string{"stroke","#16c79a"}, [2]string{"stroke-width","2"}` so the triangle outline is visible even without fill |
| requiredExtensions/systemLanguage hidden | `overlay.go` | Same fix as `rect` |

---

## `<path>` — HAS_EMPTY_CARDS + GRAMMAR_ISSUE

63 enumerated value-paths walked.

### HAS_EMPTY_CARDS

| Label | Why empty |
|-------|-----------|
| `d="none"` | SVG2 allows `d="none"` as a valid value meaning "no path data" — the path is present in the DOM but has no geometry; renders as an empty card |
| `requiredExtensions="label"` | element hidden |
| `systemLanguage="label"` | element hidden |

### GRAMMAR_ISSUE

| Label | Problem |
|-------|---------|
| `stroke-miterlimit="0"` | SVG2 spec § 17.4 requires stroke-miterlimit ≥ 1. The value 0 is out-of-range. This arises because `StrokeMiterlimitAttr` uses raw `NumberType` whose first rep is "0", and `isNonNegativeAttr` in `overlay.go` does not include `stroke-miterlimit`. Browsers silently clamp or ignore the value, so the card renders as the baseline path, but the emitted SVG is technically invalid. |

### Also noted

- `fill="none"` card: path baseline is `fill="none"` already (stroke-only arc); changing fill to "none" is a no-op — the card looks correct but is identical to baseline. Not empty, just low-signal.
- `stroke-opacity="0.5"` — faint arc, correct and visually meaningful.
- The majority of id/tabindex/lang/aria/event cards (32 of 63) are visually identical arcs.

### Fix targets

| Issue | Target | Change |
|-------|--------|--------|
| d="none" blank | `overlay.go` | Add: `case an == "d": return "M10 50 Q50 10 90 50", true` so the `d` attribute is always given a visible path. This replaces the "none" arm with a meaningful path. Alternatively, in the grammar `PathDAttr`, move "none" to a separate variant with a note that it intentionally produces no geometry; this is semantically correct so the gallery could annotate it. |
| stroke-miterlimit="0" invalid | `overlay.go` `isNonNegativeAttr` | Add `"stroke-miterlimit"` to the `isNonNegativeAttr` switch, BUT add a separate clamp: the overlay should return `"4"` (the SVG default) rather than "2" from the non-negative path, since the minimum is 1 not 0. Best fix: add a dedicated case `case an == "stroke-miterlimit": return "4", true`. |
| requiredExtensions/systemLanguage hidden | `overlay.go` | Same fix as `rect` |

---

## Top recurring issues across the batch

### 1. `systemLanguage="label"` and `requiredExtensions="label"` — ALL 7 elements (14 hidden cards total)

**Root cause:** `SystemLanguageAttr` and `RequiredExtensionsAttr` both use `StringType` reps (`"label"`, `"Aa"`, etc.). These are not valid BCP-47 language tags or extension URIs. The browser's conditional-processing logic hides any element whose `systemLanguage` list does not contain the document's language, or whose `requiredExtensions` list contains an unsatisfied extension.

**Fix:** In `overlay.go`, add two cases in `overlaySample`:
```go
case an == "systemlanguage":
    return "en", true   // always matches; swap with user locale if needed
case an == "requiredextensions":
    return "", true     // empty list = no restriction; element always renders
```
(Note: an empty string may not be well-formed for `requiredExtensions`; alternatively use a known value such as `http://www.w3.org/1999/xhtml` which browsers support, or remove the attribute from the enumeration entirely by adding it to a skip list.)

### 2. `fill="none"` without a baseline stroke — rect, circle, ellipse, polygon (4 blank cards)

**Root cause:** `fill="none"` is a valid and important value to demonstrate, but the element baselines for `rect`, `circle`, `ellipse`, and `polygon` have no `stroke`, so these cards render as empty viewboxes.

**Fix:** In `blueprint.go` `baselineFor()`, add a thin stroke to each filled-shape baseline:
- `rect`: append `[2]string{"stroke","#16c79a"}, [2]string{"stroke-width","2"}`
- `circle`: append `[2]string{"stroke","#e94560"}, [2]string{"stroke-width","2"}`
- `ellipse`: append `[2]string{"stroke","#e94560"}, [2]string{"stroke-width","2"}`
- `polygon`: append `[2]string{"stroke","#16c79a"}, [2]string{"stroke-width","2"}`

### 3. `x="100%"` / `y="100%"` / `cx="100%"` / `cy="100%"` — rect, circle (4 off-canvas cards)

**Root cause:** `LengthPercentageType` reps include `"100%"` which in a `viewBox="0 0 100 100"` maps to coordinate 100 — the far edge. With the default width/height of 80 units, the shape is pushed entirely outside the viewBox.

**Fix:** In `reps.go`, change `LengthPercentageType`:
```go
"LengthPercentageType": {"10", "24px", "2em", "50%", "75%"},
```
Replace `"100%"` with `"75%"` so the element is still partially off-canvas (showing the percentage effect) but is still visible.

### 4. IDENTICAL_NO_EFFECT — id, tabindex, autofocus, lang, xml:lang, xml:space, role, aria-*, on* events — ALL 7 elements (~32 cards each)

**Root cause:** These are semantics/accessibility/scripting attributes that have no visual rendering effect on a static SVG shape. They are valid SVG attributes and the grammar correctly enumerates them, but the resulting gallery cards are all visually identical to the element baseline.

**Fix options (pick one or combine):**
- In `emit.go` `emitPage()`, group non-visual attributes under a collapsed `<details>` element with a label "Non-visual attributes (accessibility / scripting)" so they don't dominate the visual grid.
- Alternatively, add a `nonVisualAttr(tag, attr string) bool` helper and skip them from the main grid, emitting them only in a compact text list at the bottom of the page.
- The most minimal fix: add a CSS class `card--metadata` and visually differentiate them (e.g. dimmed border, italic label) so reviewers know they are expected to be identical.

### 5. `stroke-miterlimit="0"` — invalid value (GRAMMAR_ISSUE, path + rect + all shapes)

**Root cause:** `StrokeMiterlimitAttr` uses `NumberType` whose reps include `"0"`. The SVG2 spec requires `stroke-miterlimit ≥ 1`. `"stroke-miterlimit"` is absent from `isNonNegativeAttr`.

**Fix:** In `overlay.go`, add to `overlaySample`:
```go
case an == "stroke-miterlimit":
    return "4", true  // SVG default; spec minimum is 1
```

### 6. `d="none"` empty path — path (1 blank card)

**Root cause:** `PathDAttr` grammar has `"none" | SvgPath` as alternatives. "none" is a valid SVG2 value (it means "no path data") but renders nothing. The enumerator picks it as the first `d` variant.

**Fix:** In `overlay.go`, steer the `d` attribute away from "none":
```go
case an == "d":
    return "M10 50 Q50 10 90 50", true
```
This ensures `d` always resolves to visible path data. The "none" arm of the grammar remains valid but is not exercised in the gallery.
