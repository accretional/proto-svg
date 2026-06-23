# Round 1 Visual QA — Batch: text-embedded elements

Elements reviewed: `text`, `tspan`, `textPath`, `image`, `foreignObject`, `view`, `script`, `style`
Screenshots: `chrome-testing/screenshots/review/<tag>.png`
Date: 2026-06-22

---

## `<text>`

**Category: MOSTLY_PERFECT**

74 enumerated cards. The baseline is solid: `x="10" y="55" fill="#e6e6e6" font-size="20"` produces clearly visible white "Ag" text on the dark background. Nearly all cards render legibly.

### Problem cards

| Label | Issue |
|---|---|
| `textLength="10"` | Text is squeezed to 10px wide; glyphs are illegibly compressed to a sliver. Technically valid SVG behavior but visually confusing as a QA card — the attribute's effect is extreme and the card looks broken. |
| `textLength="1"` | Same as above, even worse compression. |
| `textLength="0"` | Text collapses to zero width; card appears blank/empty (no visible glyph). |
| `rotate="3.140.52"` | Label shows `rotate="3.140.52"` — this is a malformed value. SVG `rotate` on `<text>` takes a list of numbers; the generated value "3.140.52" is not a valid number list (two consecutive decimal points in one token). The glyph renders oddly rotated (visible as a small rotated "A" icon in the screenshot). |

### Recommendations

- **`textLength` small values** — `FIX TARGET: blueprint.go / baselineFor("text")`. Add `textLength` to the `isNonNegativeAttr` guard in `overlay.go` so the overlay constrains it to a reasonable minimum, or add a specific clamp so `textLength` samples from larger values (e.g., `"40"`, `"60"`, `"80"`). The current `LengthType` reps include `"10"` and `"24px"` — both make "Ag" nearly invisible when used as textLength with font-size 20.
- **`rotate` grammar value** — `FIX TARGET: reps.go`. The `NumberType` representative `"3.14"` is being concatenated in a list context to produce `"3.140.52"`. The rotate attribute on `<text>` expects `ListOfNumbersType`. Check that the grammar routes `rotate` through `ListOfNumbersType` samples rather than two adjacent `NumberType` reps with no separator. The fix is in the EBNF `.ebnf` for `text`'s `rotate` attribute or in the emit logic that concatenates values.

---

## `<tspan>`

**Category: MOSTLY_PERFECT**

74 enumerated cards. Blueprint wraps tspan inside `<text x="10" y="55" font-size="18" fill="#e6e6e6">`, producing clearly readable "Ag" in all non-pathological cases.

### Problem cards

| Label | Issue |
|---|---|
| `textLength="0"` | Collapses text to invisible, same as `<text>`. |
| `textLength="10"`, `textLength="1"` | Extreme compression, same as `<text>`. |
| `rotate="3.140.52"` | Same malformed rotate value as in `<text>` — concatenation artifact producing two decimal points in one token. |
| `role="math menumenubar"` | Not visually broken, but the generated role string "math menumenubar" is semantically invalid (space-separated ARIA role list contains "menumenubar" which is not a standard role). Minor grammar issue. |

### Recommendations

- Same `textLength` and `rotate` fixes as `<text>` above.
- **`role` grammar** — `FIX TARGET: reps.go / StringType or the role .ebnf`. The `StringType` values ("label", "Aa", "sample", "specimen") are being concatenated with spaces to form multi-token role strings containing nonsense tokens. Either use a curated list of valid ARIA role tokens (e.g., `"img"`, `"presentation"`, `"group"`, `"button"`, `"none"`) in the EBNF for `role`, or constrain `StringType` sampling so concatenation in role context produces only valid tokens.

---

## `<textPath>`

**Category: HAS_EMPTY_CARDS**

80 enumerated cards. The first 3 cards showing `path="M10 10 L90 90"`, `path="M10 50 Q50 10 90 50"`, and `path="M20 20 H80 V80 H20 Z"` render visible text following the path. The `href="#slot"` and `xlink:href="#slot"` cards also render (text on the arc path defined in the blueprint defs). However the vast majority of cards are **blank**.

### Empty card groups

All cards from `startOffset="10"` onward through the end of the gallery (approximately 72 of 80 cards) are completely blank — no text visible. The blank cards span every attribute group: `startOffset`, `method`, `spacing`, `side`, `textLength`, `lengthAdjust`, `id`, `tabindex`, `autofocus`, `lang`, `xml:lang`, `xml:space`, `requiredExtensions`, `systemLanguage`, `role`, `aria-*`, `fill`, `fill-rule`, `fill-opacity`, `stroke`, `stroke-*`, `paint-order`, `marker*`, `color*`, `vector-effect`, `filter`, `color-interpolation-filters`, `flood-color`, `on*` events.

### Root cause

The blueprint is: `<text font-size="13" fill="#e94560"><textPath ...baseline... >Ag</textPath></text>` where the baseline includes `x="10" y="55" fill="#e6e6e6" font-size="20"`. However, `textPath` does NOT use `x`/`y` for positioning — it follows the referenced path. When no `href` or `path` attribute is set (all the non-path, non-href variant cards), the textPath has no path to follow and renders nothing. Additionally, the baseline injects `x` and `y` attributes that are invalid on `textPath` and silently ignored. Without a path reference, the text is invisible.

### Recommendations

- **FIX TARGET: `blueprint.go` — `baselineFor("textPath")`**. Change the textPath baseline to omit `x`, `y` (invalid on textPath) and instead inject `href="#slot"` always as the baseline (the slot path is in defs). Current: `add([2]string{"x","10"}, [2]string{"y","55"}, [2]string{"fill","#e6e6e6"}, [2]string{"font-size","20"})`. Change to: `add([2]string{"href","#slot"}, [2]string{"fill","#e6e6e6"}, [2]string{"font-size","20"})`. With `href="#slot"` as the constant baseline, every variant card will render "Ag" on the arc path defined in the blueprint's defs, and the varied attribute will be observable.

---

## `<image>`

**Category: HAS_EMPTY_CARDS**

83 enumerated cards. Only 2 cards render an image: `href="#slot"` (which resolves to a self-reference, rendering a broken-image icon) and `xlink:href="#slot"` (same). All other 81 cards are blank — the image element has no `href`/`src` pointing at image data, so nothing renders.

### Root cause

The `baselineFor("image")` in `blueprint.go` only provides `x`, `y`, `width`, `height` — it does not include an `href` with a valid image data URI. Without a loadable image source, `<image>` renders nothing. The template's hand-authored cards use base64-encoded SVG data URIs (confirming they work), but the machine blueprint does not.

### Recommendations

- **FIX TARGET: `blueprint.go` — `baselineFor("image")`**. Add a working data URI as the baseline `href`. A minimal base64-encoded SVG such as:
  ```
  href="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI4MCIgaGVpZ2h0PSI4MCI+PHJlY3Qgd2lkdGg9IjgwIiBoZWlnaHQ9IjgwIiBmaWxsPSIjMjIyIi8+PGNpcmNsZSBjeD0iNDAiIGN5PSI0MCIgcj0iMzAiIGZpbGw9IiMxNmM3OWEiLz48L3N2Zz4="
  ```
  (an 80x80 dark SVG with a teal circle) added to `add(...)` in `baselineFor("image")` whenever `varying != "href"`. This ensures the image always loads and the variant attribute is observable.
- The `href="#slot"` card (overlay resolves IriType to `#slot`) renders a broken-image icon — this is expected since `#slot` is not a real image; it is acceptable as a reference-validity demonstration card.

---

## `<foreignObject>`

**Category: HAS_EMPTY_CARDS**

72 enumerated cards. All 72 machine-generated cards are blank — not a single one displays any content. The SVG elements are rendered (`<foreignObject x="10"/>`, etc.) but none contain any child HTML content, so the foreignObject is invisible.

### Root cause

The `bodyFor("foreignObject")` function in `blueprint.go` returns `""` (empty string — it falls through to the default `return ""`). The `defaultScaffold` for foreignObject falls into `catSelf` (the default case) and injects bare `<foreignObject ...attrs...></foreignObject>` with no children. Without XHTML children, foreignObject renders nothing visually. Compare to the hand-authored template cards which each contain a `<div xmlns="http://www.w3.org/1999/xhtml" ...>` child.

### Recommendations

- **FIX TARGET: `blueprint.go` — `bodyFor("foreignObject")`**. Add a body:
  ```go
  case "foreignObject":
      return `<div xmlns="http://www.w3.org/1999/xhtml" style="width:100%;height:100%;background:#16c79a;color:#0a1a0a;font:bold 14px sans-serif;display:flex;align-items:center;justify-content:center;border-radius:4px;box-sizing:border-box;">FO</div>`
  ```
  Also update `baselineFor("foreignObject")` to include `width="80"` and `height="80"` in addition to the existing `x` and `y`, so the foreignObject has explicit dimensions (required for content to display).

---

## `<view>`

**Category: PERFECT**

39 enumerated cards. All cards render the baseline shapes (red circle, teal rect, orange triangle) correctly. The `viewBox` and `preserveAspectRatio` variant cards correctly show the shapes zoomed/cropped. The `id`, `tabindex`, `autofocus`, `lang`, `xml:lang`, `xml:space`, `zoomAndPan`, `role`, `aria-*`, and event handler attribute variants all show the shapes intact — the `view` element is non-rendering and these attributes have no visual effect, which is correct behavior. Every card is pixel-perfect.

### No fixes needed.

---

## `<script>`

**Category: MOSTLY_PERFECT**

41 enumerated cards. The blueprint inserts a gray rect (`id="slot-target"`) and then injects the `<script>` element. Since the generated scripts are empty (no body content) or have invalid `type` attributes (e.g., `type="label"`, `type="Aa"`), the gray rect remains visible in all cards. Cards are not blank — the rect is consistently rendered.

### Problem cards

| Label | Issue |
|---|---|
| `type="label"`, `type="Aa"`, `type="sample"`, `type="specimen"` | These are invalid MIME types for script; Chrome treats them as inert, which is correct. The rect still shows. But the labels visually suggest script execution should do something — it doesn't. Minor UI clarity issue, not a rendering failure. |
| All cards | The rect color stays gray `#444` — no script executes to change it. This is expected (scripts are empty in machine-generated cards) and is distinct from the hand-authored template's colored results. The cards are visually identical across all attribute variants (all show the same gray rect). |

### IDENTICAL_NO_EFFECT groups

`id`, `tabindex`, `autofocus`, `lang`, `xml:lang`, `xml:space`, `crossorigin`, `async`, `defer`, `href`, `xlink:href`, `oncancel`, `oncanplay`, `oncopy`, `onfocusin` — all 35+ of these produce an identical gray-rect rendering with no visual distinction.

### Recommendations

- **FIX TARGET: `blueprint.go` — machine blueprint for `<script>`**. Add a small JS body to the baseline (overriding in `bodyFor("script")`):
  ```go
  case "script":
      return `(function(){var t=document.getElementById('slot-target');if(t){t.setAttribute('fill','#e94560');}})()`
  ```
  Then for the variant cards where `type` is the varied attribute, the type value will override the MIME; for `type="application/ecmascript"` or `type="text/javascript"` variants (if in the grammar), the script executes and changes the rect to red, making those visually distinct. For invalid types, the rect stays gray — demonstrating the spec-correct "inert script" behavior.

---

## `<style>`

**Category: MOSTLY_PERFECT**

40 enumerated cards. The blueprint injects `<style ...>` followed by a `.slot`-classed rect and circle. Cards render a dark-navy rect and circle (the `.slot` class has no CSS because the style element body is empty in machine-generated cards).

### Problem cards

| Label | Issue |
|---|---|
| `type="label"`, `type="Aa"`, `type="sample"`, `type="specimen"` | Invalid CSS type strings; browser ignores the style block entirely. Since the style body is already empty in machine cards, these look the same as valid-type cards. The rect/circle are the same near-black. |
| `media="label"`, `media="Aa"`, `media="sample"`, `media="specimen"` | Invalid media query strings; style block is ignored. Same dark rect/circle. |
| All cards | Shapes (.slot rect and circle) render uniformly dark (navy background color only — no fill applied by CSS since style body is empty). All 40 cards are visually identical. |

### IDENTICAL_NO_EFFECT groups

`type` (all 4 values), `media` (all 4 values), `title` (all 4 values), `id`, `tabindex`, `autofocus`, `lang`, `xml:lang`, `xml:space`, `oncancel`, `oncanplay`, `oncopy` — all produce identical all-dark-rect cards.

### Recommendations

- **FIX TARGET: `blueprint.go` — `bodyFor("style")`**. Add a CSS body that targets the `.slot` class:
  ```go
  case "style":
      return `.slot{fill:#e94560;stroke:#16c79a;stroke-width:3}`
  ```
  This makes the rect and circle render with a red fill and teal stroke in every valid-type card, providing visible evidence that the style element is active. For invalid `type` values (non-CSS MIME), the shapes would stay dark — a useful visual contrast between valid and invalid type strings.
- **`type`/`media`/`title` grammar values** — `FIX TARGET: reps.go / StringType`. The MIME type for style should be `"text/css"` and media values should be real media queries (e.g., `"screen"`, `"print"`, `"all"`). Add specific representative values to the EBNF attribute definition for `type` on `<style>` rather than using generic `StringType` samples.

---

## Cross-cutting Issues (Top Recurring)

### 1. Missing image `href` in baseline (image: critical)
The single largest rendering failure. `<image>` without a data URI `href` is invisible. Fix: add a tiny base64 SVG data URI to `baselineFor("image")`.

### 2. `textPath` ignores `x`/`y` baseline — needs `href="#slot"` instead (textPath: critical)
72 of 80 textPath cards are blank because `textPath` does not use `x`/`y`. Fix: replace `x`/`y` baseline with `href="#slot"`.

### 3. `foreignObject` body is empty — needs XHTML child (foreignObject: critical)
All 72 foreignObject cards are blank. Fix: add `bodyFor("foreignObject")` returning a styled div.

### 4. `textLength` extreme compression makes cards look broken (text, tspan)
`textLength="10"`, `textLength="1"`, `textLength="0"` visually appear as blank/broken cards. Fix: either exclude `"0"` and `"1"` from `LengthType` reps when used for textLength, or add `"textlength"` to `isNonNegativeAttr` with a larger floor value (e.g., `"40"`).

### 5. `rotate` attribute generates malformed value "3.140.52" (text, tspan)
Two `NumberType` reps are being concatenated into a single token without a separator, producing an invalid number. Fix: ensure the `rotate` attribute routes through `ListOfNumbersType` (space-separated numbers) rather than two adjacent `NumberType` tokens.

### 6. `<script>` and `<style>` machine cards are all identical (no CSS/JS body)
All attribute variants for script and style produce identical gray-rect / dark-shape output because no CSS or JS body is injected. Fix: add CSS body in `bodyFor("style")` and JS body in `bodyFor("script")` targeting the blueprint's slot elements.

### 7. `StringType` used for structured attributes (role, type, media)
`StringType` values ("label", "Aa", "sample", "specimen") are used for attributes that expect specific token lists (ARIA roles, MIME types, media queries). These produce semantically invalid values. Fix: add custom representative lists for `role`, `style@type`, and `style@media` attributes in their respective EBNF definitions.
