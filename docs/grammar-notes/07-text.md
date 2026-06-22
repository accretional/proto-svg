# Text Grammar Notes

## Source

- `svg2-text.txt` — W3C SVG 2 Chapter 11: Text (full read, 5927 lines)
- Cross-checked against `/Volumes/wd_office_3/Projects/proto-svg/docs/specs/mdn_docs_attributes.md`
- Browser-behavior knowledge applied for implementation-status flags

---

## Elements

### `text`

**Categories:** Graphics element, renderable element, text content element, text content block element.
In CSS terms, acts as a **block** element; its inline descendants (`tspan`, `textPath`, `a`) act as **inline** elements.

**Content model:** Any number of the following, in any order, interleaved with character data:

- Animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- Descriptive elements: `desc`, `title`, `metadata`
- Paint server elements: `linearGradient`, `radialGradient`, `pattern`
- Text content child elements: `tspan`, `textPath`
- Also allowed: `a`, `clipPath`, `marker`, `mask`, `script`, `style`

Note: Graphics elements and metadata elements inside `text` are set to `display: none` by the UA — they do not render. The content model permits them as child nodes for authoring / defs reasons, but their visual contribution is suppressed.

**Presentation attributes unique to `text` (and `tspan`):**

| Attribute | Value syntax | Initial | Animatable | Not a CSS pres-attr? |
|-----------|-------------|---------|------------|----------------------|
| `x` | `[ [ <length-percentage> \| <number> ]+ ]#` | `0` | yes | Yes — **not** settable via CSS in SVG 2 |
| `y` | `[ [ <length-percentage> \| <number> ]+ ]#` | `0` | yes | Yes — **not** settable via CSS in SVG 2 |
| `dx` | `[ [ <length-percentage> \| <number> ]+ ]#` | (none) | yes | Yes |
| `dy` | `[ [ <length-percentage> \| <number> ]+ ]#` | (none) | yes | Yes |
| `rotate` | `[ <number>+ ]#` | (none) | yes (non-additive) | Yes |
| `textLength` | `<length-percentage> \| <number>` | computed text length | yes | Yes |
| `lengthAdjust` | `spacing \| spacingAndGlyphs` | `spacing` | yes | Yes |

**DOM Interface:** `SVGTextElement` extends `SVGTextPositioningElement` extends `SVGTextContentElement`.

**Context-sensitive constraints:**

- `x` / `y` are ignored for auto-wrapped text (shape-inside), except to set the initial current text position when using `inline-size`.
- `dx`, `dy`, `rotate` are always ignored for auto-wrapped text.
- `textLength` is not applied when the wrapping area is defined by `shape-inside` or `inline-size`.
- `textLength` is not applied to any `text`/`tspan` that has forced line breaks (`white-space: pre` or `pre-line`).
- A negative `textLength` value is an error.
- `x` and `y` lists create new **text chunks** (independent anchored runs); `lengthAdjust` applies per chunk.

---

### `tspan`

**Categories:** Graphics element, renderable element, text content element, text content child element.
Acts as an **inline** element in CSS terms.

**Content model:** Any number of the following, in any order, interleaved with character data:

- Descriptive elements: `desc`, `title`, `metadata`
- Paint server elements: `linearGradient`, `radialGradient`, `pattern`
- Also allowed: `a`, `animate`, `script`, `set`, `style`, `tspan`

Note: `tspan` does NOT allow `textPath`, `animateMotion`, `animateTransform`, `discard`, `clipPath`, `marker`, or `mask` as direct content (unlike `text`).

**Presentation attributes:** Same set as `text`: `x`, `y`, `dx`, `dy`, `rotate`, `textLength`, `lengthAdjust` — all with the same value syntaxes.

**Differences from `text` in attribute semantics:**

| Attribute | `text` initial | `tspan` initial |
|-----------|---------------|-----------------|
| `x` | `0` | (none) — inherits from ancestor context |
| `y` | `0` | (none) — inherits from ancestor context |
| `dx` | (none) | (none) |
| `dy` | (none) | (none) |
| `rotate` | (none) | (none) — if unset, last rotate value from ancestor propagates to remaining characters |

**DOM Interface:** `SVGTSpanElement` extends `SVGTextPositioningElement`.

**Context-sensitive constraints:**

- All constraints from `text` apply.
- Transforms on `tspan` are NOT allowed (SVG 2 reversed an earlier decision to align with HTML, where inline elements do not support transforms).
- `rotate` propagation rule: if fewer `<number>` values than characters, the **last** value applies to all remaining characters in the element and any `tspan` descendants that do not specify their own `rotate`.

---

### `textPath`

**Categories:** Graphics element, renderable element, text content element, text content child element.
Acts as an **inline** element in CSS terms. Always creates an anchored chunk.

**Content model:** Any number of the following, in any order, interleaved with character data:

- Descriptive elements: `desc`, `title`, `metadata`
- Paint server elements: `linearGradient`, `radialGradient`, `pattern`
- Also allowed: `a`, `animate`, `clipPath`, `marker`, `mask`, `script`, `set`, `style`, `tspan`

**Attributes:**

| Attribute | Value syntax | Initial | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `path` | `<path-data>` (path data string) | (none) | yes | Inline path data. Empty string = no path (text not rendered). Takes precedence over `href`. |
| `href` | `<url>` | (see prose) | yes | Reference to a `path` element or basic shape. Used if `path` attr absent or erroneous. |
| `xlink:href` | `<url>` | (see prose) | yes | **Deprecated** — fallback after `href`. |
| `startOffset` | `<length-percentage> \| <number>` | `0` | yes | Offset from path start. Negative and >100% allowed. `<number>` = user-unit distance. `<percentage>` = fraction of path length. |
| `method` | `align \| stretch` | `align` | yes | `align` = 2×3 matrix per glyph. `stretch` = warp glyph outlines to path normals. |
| `spacing` | `auto \| exact` | `exact` | yes | `exact` = per spec distance. `auto` = UA may adjust spacing for visual quality. |
| `side` | `left \| right` | `left` | yes | **SVG 2 new.** `right` effectively reverses the path direction, placing text on the other side. |
| `textLength` | `<length-percentage> \| <number>` | computed | yes | Inherited from `SVGTextContentElement`. |
| `lengthAdjust` | `spacing \| spacingAndGlyphs` | `spacing` | yes | Inherited from `SVGTextContentElement`. |

**Path selection precedence (if multiple path sources present):**

1. `path` attribute (if non-empty and valid)
2. `href` attribute
3. `xlink:href` attribute (deprecated)

If `href` references a non-path / non-shape or is an invalid reference, the entire `textPath` content is not rendered.

**DOM Interface:** `SVGTextPathElement` extends `SVGTextContentElement`; also mixes in `SVGURIReference`.

**Context-sensitive constraints:**

- `x` attributes on `text`/`tspan` inside a horizontal `textPath` provide new absolute offsets along the path; `y` is ignored for horizontal text-on-path.
- `y` attributes on `text`/`tspan` inside a vertical `textPath` provide new absolute offsets along the path; `x` is ignored for vertical text-on-path.
- `dy` shifts glyph perpendicular to path (lifts / lowers baseline off path).
- Glyphs whose midpoints fall off an open path are not rendered (hidden flag set). For a single closed subpath, glyphs wrap the full circuit and the `startOffset` + `text-anchor` determine the wrap anchor.
- After `textPath` ends, subsequent text inside the same `text` element continues from the **end of the path** (not the position of the last glyph), unless a new anchored chunk (`x`/`y`) begins.

---

### Deprecated / removed elements

**`tref`** — SVG 1.1 element that referenced text content of another element via `xlink:href`. **Removed entirely in SVG 2.** Not defined in this chapter; mentioned only in historical context. Do not include in grammar.

**`altGlyph` / `altGlyphDef` / `altGlyphItem`** — SVG 1.1 elements for glyph substitution. **Removed in SVG 2.** Not defined here.

---

## Shared attribute groups used by text elements

**`SVGTextContentElement` IDL interface** (base for all three elements):

```
textLength   → SVGAnimatedLength
lengthAdjust → SVGAnimatedEnumeration (LENGTHADJUST_SPACING=1, LENGTHADJUST_SPACINGANDGLYPHS=2, LENGTHADJUST_UNKNOWN=0)
```

**`SVGTextPositioningElement` IDL interface** (base for `text` and `tspan`):

```
x       → SVGAnimatedLengthList
y       → SVGAnimatedLengthList
dx      → SVGAnimatedLengthList
dy      → SVGAnimatedLengthList
rotate  → SVGAnimatedNumberList
```

**`SVGTextPathElement` DOM constants:**

```
TEXTPATH_METHODTYPE_UNKNOWN  = 0
TEXTPATH_METHODTYPE_ALIGN    = 1
TEXTPATH_METHODTYPE_STRETCH  = 2
TEXTPATH_SPACINGTYPE_UNKNOWN = 0
TEXTPATH_SPACINGTYPE_AUTO    = 1
TEXTPATH_SPACINGTYPE_EXACT   = 2
```

---

## Text & Font Properties

Properties are grouped as: (A) SVG-specific / SVG-adapted CSS, (B) standard CSS referenced by SVG 2. Each entry gives the full grammar form (closed sets fully enumerated), initial value, inheritance, presentation-attribute/animatable status, and verification notes.

---

### `text-anchor` (SVG-specific)

**Spec location:** §11.10.1.1

```ebnf
text-anchor = "start" | "middle" | "end"
```

| Field | Value |
|-------|-------|
| Initial | `start` |
| Inherited | yes |
| Applies to | text content elements (`text`, `tspan`, `textPath`) |
| Presentation attr | yes |
| Animatable | yes |
| Percentages | N/A |

**Semantics:**

- `start` — text rendered starting at the initial current text position (LTR: left-aligned; RTL: right-aligned; vertical: top-aligned)
- `middle` — text shifted so geometric midpoint aligns with current text position
- `end` — text shifted so final text position aligns with initial current text position (LTR: right-aligned; RTL: left-aligned; vertical: bottom-aligned)

**Applies only to:**
- Pre-formatted text
- Auto-wrapped text where wrapping area is defined by `inline-size`
- Does **not** apply to auto-wrapped text in `shape-inside`; use `text-align` instead.

**Verification:** Fully implemented in all major browsers. MDN lists `text-anchor` without deprecated/experimental flags.

---

### `dominant-baseline` (SVG-adapted CSS)

**Spec location:** §11.10.2.5; normative definition in CSS Inline Layout Module Level 3.

SVG 2 removes these SVG 1.1 values: `reset-size`, `use-script`, `no-change`.
SVG 2 makes the property **inherited** (was non-inherited in some earlier drafts; behavior unchanged in practice).

The current CSS Inline Level 3 keyword set (grammar form for SVG 2):

```ebnf
dominant-baseline = "auto"
                  | "text-bottom"
                  | "alphabetic"
                  | "ideographic"
                  | "middle"
                  | "central"
                  | "mathematical"
                  | "hanging"
                  | "text-top"
```

| Field | Value |
|-------|-------|
| Initial | `auto` |
| Inherited | yes |
| Applies to | text content elements |
| Presentation attr | yes |
| Animatable | yes |

**Note on `auto`:** For horizontal writing modes, resolves to `alphabetic`. For vertical, resolves to `central`. For `text-orientation: sideways`, `auto` resolves to `alphabetic` but glyphs are aligned using `central` for backward compatibility.

**Verification:** Well supported. MDN: no deprecated flags. The old SVG 1.1 values (`reset-size`, `use-script`, `no-change`) are not in browsers; omit from grammar.

---

### `alignment-baseline` (SVG-adapted CSS)

**Spec location:** §11.10.2.6; normative definition in CSS Inline Layout Module Level 3.

SVG 2 removes: `auto`, `before-edge`, `after-edge`.
Backward-compat mapping: `text-before-edge` → `text-top`; `text-after-edge` → `text-bottom`.
Preferred modern alternative: use `vertical-align` shorthand.

```ebnf
alignment-baseline = "baseline"
                   | "text-bottom"
                   | "alphabetic"
                   | "ideographic"
                   | "middle"
                   | "central"
                   | "mathematical"
                   | "hanging"
                   | "text-top"
```

| Field | Value |
|-------|-------|
| Initial | `baseline` |
| Inherited | no (applies per inline box) |
| Applies to | inline-level boxes; `tspan`, `textPath`, `a` within text content elements |
| Presentation attr | yes |
| Animatable | yes |

**Verification:** MDN: no deprecated flags on `alignment-baseline`. The removed values (`auto`, `before-edge`, `after-edge`) are not present in current browser implementations.

---

### `baseline-shift` (SVG-adapted CSS)

**Spec location:** §11.10.2.7; normative definition in CSS Inline Layout Module Level 3.

SVG 1.1 value `baseline` removed in SVG 2 (redundant with `0`; UAs may support it as `0` for legacy).
Preferred: use `vertical-align` shorthand.

```ebnf
baseline-shift = "sub" | "super" | <length-percentage>
```

Note: This is an open domain (length-percentage part is a named leaf).

| Field | Value |
|-------|-------|
| Initial | `0` (was `baseline` in SVG 1.1) |
| Inherited | no |
| Applies to | inline-level boxes |
| Presentation attr | yes |
| Animatable | yes |

**Verification:** Supported in browsers. MDN: no deprecated/experimental flags on `baseline-shift` itself.

---

### `direction` (SVG-adapted CSS)

**Spec location:** §11.10.2.4; normative in CSS Writing Modes Level 3.

```ebnf
direction = "ltr" | "rtl"
```

| Field | Value |
|-------|-------|
| Initial | `ltr` |
| Inherited | yes |
| Applies to | text content elements |
| Presentation attr | yes |
| Animatable | yes |

**Note:** Only affects glyphs oriented perpendicular to the inline-base direction (the usual case for horizontally-oriented text). The Unicode bidi algorithm handles most reordering automatically; `direction` + `unicode-bidi` are needed for explicit override.

---

### `unicode-bidi`

Referenced by SVG 2; defined in CSS Writing Modes Level 3.

```ebnf
unicode-bidi = "normal" | "embed" | "isolate" | "bidi-override" | "isolate-override" | "plaintext"
```

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | no |
| Presentation attr | yes |
| Animatable | yes |

**Verification:** All values well supported in modern browsers.

---

### `writing-mode` (SVG-adapted CSS)

**Spec location:** §11.10.2.3; normative in CSS Writing Modes Level 3.

SVG 2 requires the CSS Writing Modes Level 3 values. The SVG 1.1 values are **obsolete** but must still be supported via the following computed-value mapping:

| Specified (legacy) | Computed |
|--------------------|----------|
| `lr`, `lr-tb`, `rl`, `rl-tb` | `horizontal-tb` |
| `tb`, `tb-rl` | `vertical-rl` |

Current grammar (CSS Writing Modes Level 3):

```ebnf
writing-mode = "horizontal-tb" | "vertical-rl" | "vertical-lr"
             | "sideways-rl" | "sideways-lr"
```

Legacy SVG 1.1 values (map on parsing, not stored):

```ebnf
writing-mode-legacy = "lr" | "lr-tb" | "rl" | "rl-tb" | "tb" | "tb-rl"
```

| Field | Value |
|-------|-------|
| Initial | `horizontal-tb` |
| Inherited | yes |
| Applies to | `text` elements (determines block-flow direction and line stacking) |
| Presentation attr | yes |
| Animatable | yes |

**Verification:** `horizontal-tb`, `vertical-rl`, `vertical-lr` implemented. `sideways-rl` / `sideways-lr` have limited browser support (Firefox had them; Chrome shipped later). Legacy values are supported by most browsers for compatibility.

---

### `text-anchor` — see above (§ SVG-specific).

---

### `glyph-orientation-horizontal` (Deprecated / Removed)

**Spec location:** §11.10.1.2

**Removed in SVG 2.** Do not include in grammar for new content.

**Decision:** Omit from EBNF as a valid production. May appear as an unrecognized attribute in legacy content; UA should ignore it.

---

### `glyph-orientation-vertical` (Deprecated / Aliased)

**Spec location:** §11.10.1.3

**Obsoleted in SVG 2** — partially replaced by CSS `text-orientation`. Must still be supported by aliasing:

| Specified value | Maps to `text-orientation` |
|-----------------|---------------------------|
| `auto` | `mixed` |
| `0`, `0deg` | `upright` |
| `90`, `90deg` | `sideways` |
| any other value | invalid |

**MDN status:** `glyph-orientation-vertical` — deprecated (⚠️ per MDN index).

**Decision for grammar:** Do not make `glyph-orientation-vertical` a primary production. Define it as a deprecated alias accepting only the 4 mapped tokens, with a note.

---

### `kerning` (Removed)

**Spec location:** §11.10.1.4

**Removed in SVG 2.** Replaced by CSS `font-kerning`. Do not include in grammar.

**MDN:** Not listed as a current SVG attribute.

---

### `font-family`

Defined by CSS Fonts Module Level 3. Referenced by SVG 2.

```ebnf
font-family = <family-name># | <generic-family>
            | "inherit" | "initial" | "unset"
generic-family = "serif" | "sans-serif" | "cursive" | "fantasy" | "monospace"
                 | "system-ui" | "emoji" | "math" | "fangsong"
family-name    = <string> | <custom-ident>+        (* open leaf *)
```

| Field | Value |
|-------|-------|
| Initial | implementation-dependent (UA default) |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

Note: The `generic-family` set is closed (enumerated above). `family-name` is open.

---

### `font-size`

```ebnf
font-size = <absolute-size> | <relative-size> | <length-percentage>
absolute-size = "xx-small" | "x-small" | "small" | "medium" | "large" | "x-large" | "xx-large" | "xxx-large"
relative-size  = "larger" | "smaller"
```

| Field | Value |
|-------|-------|
| Initial | `medium` |
| Inherited | yes |
| Percentages | relative to parent element's computed font-size |
| Presentation attr | yes |
| Animatable | yes |

---

### `font-style`

```ebnf
font-style = "normal" | "italic" | "oblique" <angle>?
```

Where `<angle>` is optional; the oblique angle range is -90deg to 90deg (CSS Fonts Level 4).

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

---

### `font-weight`

```ebnf
font-weight = "normal" | "bold" | "bolder" | "lighter"
            | <number>   (* 1–1000, CSS Fonts Level 4 *)
```

| Field | Value |
|-------|-------|
| Initial | `normal` (= 400) |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

**Note:** CSS Fonts Level 3 uses only keyword values; CSS Fonts Level 4 introduces numeric values. SVG 2 references Fonts Level 3 but browsers implement Level 4 numeric weights.

---

### `font-stretch`

```ebnf
font-stretch = "normal"
             | "ultra-condensed" | "extra-condensed" | "condensed" | "semi-condensed"
             | "semi-expanded" | "expanded" | "extra-expanded" | "ultra-expanded"
             | <percentage>    (* CSS Fonts Level 4: 50%–200% *)
```

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

**Verification:** MDN also lists `font-width` as a 🧪🔶 non-standard attribute — this appears to be an MDN alias/alias experiment, not an SVG 2 attribute. Use `font-stretch` in grammar.

---

### `font-variant`

**Spec location:** §11.10.2.1

SVG 2 requires all CSS Fonts Level 3 `font-variant` sub-properties. In SVG 2, `font-variant` is a shorthand for the expanded sub-properties.

```ebnf
font-variant = "normal" | "none"
             | [ <font-variant-ligatures> || <font-variant-alternates>
                 || <font-variant-caps> || <font-variant-numeric>
                 || <font-variant-east-asian> ]
```

Sub-properties (each is a separate property; all presentational-attribute capable):

- `font-variant-ligatures`
- `font-variant-alternates`
- `font-variant-caps`
- `font-variant-numeric`
- `font-variant-east-asian`
- `font-variant-position`

Full closed enumerations for each sub-property are defined in CSS Fonts Level 3 (open in this file; reference CSS spec for full value lists).

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Presentation attr | yes (shorthand and sub-properties) |
| Animatable | yes |

---

### `font-size-adjust`

```ebnf
font-size-adjust = "none" | <number>
```

| Field | Value |
|-------|-------|
| Initial | `none` |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

---

### `font-feature-settings`

```ebnf
font-feature-settings = "normal" | <feature-tag-value>#
feature-tag-value = <string> [ <integer> | "on" | "off" ]?
```

Open domain (feature tag strings are arbitrary 4-char OpenType feature tags).

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

---

### `font-kerning`

```ebnf
font-kerning = "auto" | "normal" | "none"
```

Replaces the removed SVG 1.1 `kerning` property.

| Field | Value |
|-------|-------|
| Initial | `auto` |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

---

### `letter-spacing` (SVG-adapted CSS)

**Spec location:** §11.10.2.8

SVG 2 removes percentage values (SVG 1.1 allowed percentages relative to viewport). SVG 2 aligns with CSS Text Level 3.

```ebnf
letter-spacing = "normal" | <length>
```

Note: Only `<length>`, NOT `<length-percentage>`. Percentages are invalid in SVG 2.

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

---

### `word-spacing` (SVG-adapted CSS)

**Spec location:** §11.10.2.9

SVG 2 changes percentage semantics: in SVG 1.1 percentages were relative to viewport; in SVG 2 (following CSS Text Level 3) percentages are relative to the affected character's width.

```ebnf
word-spacing = "normal" | <length-percentage>
```

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Percentages | relative to affected character's width (SVG 2) |
| Presentation attr | yes |
| Animatable | yes |

---

### `text-decoration`

SVG 2 supports the full CSS Text Decoration Module Level 3 shorthand and longhands.
SVG adds two extra properties: `text-decoration-fill` and `text-decoration-stroke`.

**Shorthand:**

```ebnf
text-decoration = "none"
                | [ <text-decoration-line> || <text-decoration-style> || <text-decoration-color> ]
text-decoration-line  = "none" | [ "underline" || "overline" || "line-through" || "blink" ]
text-decoration-style = "solid" | "double" | "dotted" | "dashed" | "wavy"
text-decoration-color = <color>    (* open leaf *)
```

**SVG-specific longhands:**

```ebnf
text-decoration-fill   = <paint>    (* open leaf; see painting chapter *)
text-decoration-stroke = <paint>    (* open leaf *)
```

| Field | Value |
|-------|-------|
| Initial | `none` (shorthand) |
| Inherited | no (decoration inherits visually but the property itself is non-inherited) |
| Presentation attr | yes |
| Animatable | yes |

**Initial for `text-decoration-fill` / `text-decoration-stroke`:** "see prose" — if not explicitly set, uses the fill/stroke of the text at the point where decoration is declared.

---

### `white-space` (SVG-adapted CSS)

**Spec location:** §11.10.3.1

New in SVG 2. Preferred over the legacy `xml:space` attribute. Values from CSS Text Module Level 3.

```ebnf
white-space = "normal" | "pre" | "nowrap" | "pre-wrap" | "pre-line" | "break-spaces"
```

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Applies to | text content elements |
| Presentation attr | yes |
| Animatable | yes |

**Effect on line breaking:**
- `pre` / `pre-line` → forced line breaks at `\n` or `\r` within text content
- `pre-line` → collapses spaces but preserves line breaks; enables multi-line pre-formatted text
- `nowrap` → prevents wrapping (useful with `text-overflow`)
- `normal` → collapses white space, allows soft wrapping

**Legacy alternative:** `xml:space` attribute (values: `default` | `preserve`) — deprecated in SVG 2, no longer animatable. When `white-space` is set on an element, `xml:space` is ignored on that element.

---

### `text-overflow` (SVG-adapted CSS)

**Spec location:** §11.10.2.10; normative in CSS UI Level 3.

New in SVG 2. Only applies when there is a validly specified wrapping area. Does not apply to pre-formatted text or text-on-a-path.

```ebnf
text-overflow = "clip" | "ellipsis"
```

In SVG 2, only `clip` and `ellipsis` are recognized; any other value is treated as if the property was not specified.

| Field | Value |
|-------|-------|
| Initial | `clip` |
| Inherited | no |
| Applies to | text content block elements (`text`) with a wrapping area |
| Presentation attr | yes |
| Animatable | yes |

**Note:** Effect is purely visual. The ellipsis does not enter the DOM.

---

### `inline-size` (SVG text wrapping)

**Spec location:** §11.4.1

New in SVG 2. Creates a rectangular wrapping area. Zero disables wrapping.

```ebnf
inline-size = <length-percentage>
```

| Field | Value |
|-------|-------|
| Initial | `0` |
| Inherited | no |
| Applies to | `text` elements only |
| Percentages | relative to width (horizontal text) or height (vertical text) of current SVG viewport |
| Presentation attr | yes |
| Animatable | yes |

**Implementation status:** Partial. Chrome and Safari have shipped `inline-size` on SVG text. Firefox support is limited or absent. Flag as: **partial browser support — not universally implemented**.

**Interaction:** If both `inline-size` and `shape-inside` (not `auto`) are set, `shape-inside` takes precedence.

---

### `shape-inside` (SVG text wrapping)

**Spec location:** §11.4.2; normative reference: CSS Shapes Module Level 2.

New in SVG 2. Defines the content area for text wrapping as a CSS basic shape or reference to an SVG shape/image.

```ebnf
shape-inside = "auto" | [ <basic-shape> | <uri> ]+
basic-shape  = circle() | ellipse() | polygon()    (* inset() is invalid in SVG *)
```

Note: CSS values `outside-shape`, `shape-box`, and `display` are invalid for SVG. SVG allows a list of shapes (each an independent content area; text overflows to next).

| Field | Value |
|-------|-------|
| Initial | `auto` |
| Inherited | no |
| Applies to | `text` elements only |
| Percentages | relative to `viewBox` / current user coordinate system |
| Presentation attr | yes |
| Animatable | yes (shape interpolation) |

**Implementation status:** Very limited. `shape-inside` on SVG `text` is not implemented in any major browser as of 2026. Flag as: **not implemented in browsers**.

---

### `shape-subtract`

**Spec location:** §11.4.3

New in SVG 2. Defines wrapping exclusions (areas subtracted from the content area).

```ebnf
shape-subtract = "none" | [ <basic-shape> | <uri> ]+
```

| Field | Value |
|-------|-------|
| Initial | `none` |
| Inherited | no |
| Applies to | `text` elements only |
| Percentages | relative to `viewBox` |
| Presentation attr | yes |
| Animatable | yes |

**Implementation status:** Not implemented in browsers. Same status as `shape-inside`.

---

### `shape-padding`

**Spec location:** §11.4.6; normative: CSS Shapes Module Level 2.

Offsets the wrapping area inward from the content area boundary.

```ebnf
shape-padding = <length-percentage>
```

Positive values only.

| Field | Value |
|-------|-------|
| Initial | `0` |
| Inherited | no |
| Applies to | `text` elements |
| Presentation attr | yes |
| Animatable | yes |

**Implementation status:** Not implemented in browsers.

---

### `shape-margin`

**Spec location:** §11.4.5; normative: CSS Shapes Module Level 1.

Adds margin to shapes referenced by `shape-subtract`. Positive values only.

```ebnf
shape-margin = <length-percentage>
```

| Field | Value |
|-------|-------|
| Initial | `0` |
| Inherited | no |
| Applies to | `text` elements |
| Percentages | N/A (spec says N/A, though CSS Shapes Level 1 may differ) |
| Presentation attr | yes |
| Animatable | yes |

**Implementation status:** Not implemented in browsers.

---

### `line-height` (SVG-adapted CSS)

**Spec location:** §11.10.2.2; normative in CSS 2.1 / CSS Inline Level 3.

Used in SVG to determine leading between lines in multi-line text. Not applicable to text on a path.

```ebnf
line-height = "normal" | <number> | <length-percentage>
```

| Field | Value |
|-------|-------|
| Initial | `normal` |
| Inherited | yes |
| Presentation attr | yes |
| Animatable | yes |

---

### Properties deferred to CSS specifications (full value sets defined there)

The following properties are applicable to SVG text content elements and are referenced by SVG 2, but their full value grammars are defined in external CSS specifications. They are **open in SVG grammar terms** — treat as named leaves referencing CSS grammar modules:

| Property | CSS module | Initial | Inherited | Pres-attr | Notes |
|----------|-----------|---------|-----------|-----------|-------|
| `font` (shorthand) | CSS Fonts Level 3 | see sub-properties | yes | yes | Shorthand for font-style, font-variant, font-weight, font-stretch, font-size, line-height, font-family |
| `text-align` | CSS Text Level 3 | `start` | yes | yes | Only for auto-wrapped text; for pre-formatted text the used value is forced to `start` |
| `text-align-last` | CSS Text Level 3 | `auto` | yes | yes | Only for auto-wrapped text |
| `text-indent` | CSS Text Level 3 | `0` | yes | yes | Only for auto-wrapped text |
| `text-justify` | CSS Text Level 3 | `auto` | yes | yes | Only for auto-wrapped (non-inline-size) auto-wrapped text |
| `line-break` | CSS Text Level 3 | `auto` | yes | yes | Only for auto-wrapped text |
| `word-break` | CSS Text Level 3 | `normal` | yes | yes | Only for auto-wrapped text |
| `hyphens` | CSS Text Level 3 | `manual` | yes | yes | Only for auto-wrapped text |
| `overflow-wrap` / `word-wrap` | CSS Text Level 3 | `normal` | yes | yes | Only for auto-wrapped text |
| `vertical-align` | CSS Inline Level 3 | `baseline` | no | yes | Shorthand for alignment-baseline + baseline-shift; preferred over longhands in new content |
| `text-orientation` | CSS Writing Modes Level 3 | `mixed` | yes | yes | Replaces glyph-orientation-vertical for SVG 2 content |
| `text-combine-upright` | CSS Writing Modes Level 3 | `none` | yes | yes | Not required by SVG 2 but behavior specified if supported |

---

## Open Datatypes Used

| Datatype | Description | Where used |
|----------|-------------|-----------|
| `<length-percentage>` | CSS length or percentage | `x`, `y`, `dx`, `dy`, `textLength`, `startOffset`, `inline-size`, `shape-margin`, `shape-padding` |
| `<number>` | Unitless number (user units) | `x`, `y`, `dx`, `dy`, `rotate`, `textLength`, `startOffset` |
| `<path-data>` | SVG path data string (see path grammar) | `textPath path` attribute |
| `<url>` | CSS url() function or bare URL reference | `textPath href`, `xlink:href`, `shape-inside`, `shape-subtract` |
| `<basic-shape>` | CSS `circle()`, `ellipse()`, `polygon()` | `shape-inside`, `shape-subtract` |
| `<paint>` | SVG paint value (see painting chapter) | `text-decoration-fill`, `text-decoration-stroke` |
| `<color>` | CSS color value | `text-decoration-color` |
| `<angle>` | CSS angle (e.g. `30deg`) | `font-style` oblique angle |
| `<integer>` | Integer | `font-feature-settings` tag value |
| `<string>` | Quoted string | `font-family`, `font-feature-settings` feature tags |
| `<media-query-list>` | CSS media query | (not text-specific; inherited from style element) |

**List value notation:**

The `x`, `y`, `dx`, `dy` and `rotate` attributes use a distinctive list syntax:

```ebnf
text-coord-list     = [ <length-percentage-or-number>+ ]#
text-rotate-list    = [ <number>+ ]#
```

Where `#` denotes comma-separated repetition per CSS, and `+` allows space-separated values within each comma-group. This creates a two-level separator structure: commas between "groups", spaces between values within a group.

---

## Discrepancies, Doc Gaps & Roadblocks

### D1. `x`/`y` on `text`/`tspan` — presentation attribute status

**Spec says:** "In SVG 2, the 'text' and 'tspan' 'x' and 'y' attributes are **not** presentation attributes and **cannot** be set via CSS."

**Reality:** The spec lists them under "presentation attributes" in the element definition tables but then immediately contradicts itself in §11.2.1. The IDL interface (`SVGAnimatedLengthList`) shows they are animatable content attributes, not CSS properties. They should be treated as **geometry-style content attributes** (similar to `cx`, `cy`), not CSS presentation attributes.

**Grammar decision:** Treat `x`, `y`, `dx`, `dy`, `rotate` as non-CSS content attributes of `text`/`tspan`. They are NOT production rules in the CSS cascade layer.

---

### D2. `textLength` initial value

**Spec says:** "For the purpose of reflecting the attribute in the DOM, the initial value is the current user-agent calculated length." This makes `textLength` undefined until layout is complete.

**Grammar decision:** The attribute value syntax is `<length-percentage> | <number>`; constraint: value must be ≥ 0 (negative is an error). The computed initial is a context-dependent number — treat as optional (absent = no adjustment).

---

### D3. `dominant-baseline` value set discrepancy

**SVG spec §11.10.2.5** references CSS Inline Level 3 but removes some values (`reset-size`, `use-script`, `no-change`). CSS Inline Level 3 (current draft, 2023+) has evolved and some keyword names have changed. Specifically, CSS Inline Level 3 uses `text-bottom` and `text-top` while earlier SVG 1.1 used `before-edge` / `after-edge`.

**MDN** lists for `dominant-baseline`: `auto | ideographic | alphabetic | hanging | mathematical | central | middle | text-after-edge | text-before-edge`. This is the SVG 1.1 set. Modern CSS uses `text-top` / `text-bottom`.

**Grammar decision:** Include the CSS Inline Level 3 / SVG 2 set (`auto | text-bottom | alphabetic | ideographic | middle | central | mathematical | hanging | text-top`). Note that `text-before-edge` and `text-after-edge` are legacy aliases that map to `text-top` / `text-bottom` respectively; they should be listed as deprecated alias terminals in the overlay but not as primary grammar.

---

### D4. `alignment-baseline` — `baseline` keyword vs CSS

**SVG spec §11.10.2.6** maps `text-before-edge` → `text-top`, `text-after-edge` → `text-bottom`, removes `auto`, `before-edge`, `after-edge`. Initial is `baseline`.

**MDN** lists the initial as `auto` (reflecting CSS Inline Level 3). CSS Inline Level 3 uses `auto` as initial value; SVG uses `baseline`. This is a genuine SVG-specific override.

**Grammar decision:** Keep initial as `baseline` (SVG 2 spec). In the constraint overlay, note that `auto` is invalid for SVG `alignment-baseline`.

---

### D5. `textPath method="stretch"` — implementation gap

**Spec defines** `method="stretch"` as warping glyph outlines along path normals. This is complex.

**Browser reality:** No major browser correctly implements `method="stretch"`. Browsers silently fall back to `method="align"` behavior. MDN does not flag it as experimental, but it is effectively unimplemented.

**Grammar decision:** Include `stretch` as a valid terminal. Add implementation note: behavior is UA-defined / may fall back to `align`.

---

### D6. `textPath side="right"` — MDN experimental flag

**MDN** marks `side` as 🧪 (experimental). This is an SVG 2 addition with limited browser support.

**Grammar decision:** Include `left | right` as the full closed set. Flag as experimental in overlay/templates.

---

### D7. SVG 2 text wrapping properties — largely unimplemented

`shape-inside`, `shape-subtract`, `shape-padding`, `shape-margin` on SVG `text` elements have **no significant browser implementation** as of 2026. `inline-size` has partial support (Chromium-based browsers).

These properties ARE part of SVG 2 spec but should be flagged as `[implementation: none / partial]` in grammar overlay / templates.

---

### D8. `spacing` attribute on `textPath` — default value discrepancy

**Spec says:** initial value is `exact`.

**MDN** and browser behavior suggest browsers typically use `auto`-like behavior regardless of attribute value, because the distinction between `auto` and `exact` is subtle and rarely tested.

**Grammar decision:** Follow spec: initial = `exact`. Note browser behavior may not distinguish values.

---

### D9. `glyph-orientation-horizontal` status

**SVG 2 spec says:** "This property has been removed in SVG 2."

**MDN** lists it with ⚠️ deprecated marker.

**Grammar decision:** Do not define as a valid production in SVG 2 grammar. It should appear in a "removed properties" appendix note only.

---

### D10. `font` shorthand as presentation attribute

SVG 2 requires `font` shorthand to be usable as a presentation attribute. This is implemented in browsers. The shorthand grammar is defined in CSS Fonts Level 3 and is complex. Treat as an open-leaf reference to CSS Fonts grammar.

---

### D11. `rotate` attribute — list syntax edge case

The spec grammar shows `[ <number>+ ]#`. The `#` (comma-separator) is unusual for this attribute. In practice, both space-separated and comma-separated lists are accepted by browsers. The normative description says "comma- or space-separated list." The grammar should reflect both separators.

**Grammar decision:**

```ebnf
rotate-attr = <number> ( [ "," | ws+ ] <number> )*
```

---

### D12. `text-overflow` — SVG vs CSS value set

CSS UI Level 3 defines `text-overflow = clip | ellipsis | <string>`. SVG 2 restricts to only `clip | ellipsis` and explicitly states any other value is treated as unspecified.

**Grammar decision:** SVG `text-overflow` is a **closed** 2-value set: `"clip" | "ellipsis"`.

---

### D13. `tref` element

**SVG 1.1** had `tref` (text reference element). **Removed in SVG 2.** MDN does not list it as a current attribute. Do not include in grammar.

---

### D14. `textLength` on `tspan` — browser inconsistency

The spec note in §11.4.6 states: "Chrome supports 'textLength' on 'tspan' but Firefox does not." This is a longstanding implementation gap. `textLength` is defined on `SVGTextContentElement` (base interface for both `text`, `tspan`, and `textPath`), so the spec intends it to work on `tspan`. Grammar should include it; add implementation note for the Firefox gap.

---

## Summary (12 lines)

1. **Elements defined:** 3 primary (`text`, `tspan`, `textPath`); `tref` removed in SVG 2, `altGlyph`/`altGlyphDef` also removed.
2. **`text` content model:** char data + tspan + textPath + a + animation elements + descriptive + paint-server + clipPath + marker + mask + script + style.
3. **`tspan` content model:** char data + tspan + a + animate + set + descriptive + paint-server + script + style (narrower than `text`).
4. **`textPath` content model:** char data + tspan + a + animate + set + descriptive + paint-server + clipPath + marker + mask + script + style.
5. **Positioning attributes** (`x`, `y`, `dx`, `dy`, `rotate`) are content attributes, NOT CSS presentation attributes; list syntax uses both comma and space separators.
6. **`textLength`/`lengthAdjust`** belong to `SVGTextContentElement` (all 3 elements); `rotate` propagates last-value to remaining characters in tspan hierarchy.
7. **`textPath` path selection order:** `path` attr > `href` attr > `xlink:href` (deprecated); `method = align | stretch`; `spacing = auto | exact`; `side = left | right` (SVG 2 new, experimental).
8. **text-anchor** is the SVG-specific alignment property (3 closed keywords); `dominant-baseline` (9 keywords), `alignment-baseline` (9 keywords), `baseline-shift` (sub | super | length-pct) are SVG-adapted CSS.
9. **Writing/bidi properties:** `writing-mode` (5 CSS values + 6 legacy aliases), `direction` (ltr | rtl), `unicode-bidi` (6 keywords).
10. **SVG 2 text wrapping:** `inline-size` (partial browser support), `shape-inside` / `shape-subtract` / `shape-padding` / `shape-margin` (not implemented in browsers).
11. **Deprecated/removed:** `glyph-orientation-horizontal` (removed), `glyph-orientation-vertical` (aliased to text-orientation), `kerning` (removed, replaced by font-kerning).
12. **Top discrepancies:** `x`/`y` incorrectly listed as pres-attrs in element table (not CSS); `dominant-baseline`/`alignment-baseline` value-set drift between SVG 1.1 and CSS Inline Level 3; `shape-inside` and related wrapping properties defined in spec but not implemented in browsers; `method="stretch"` unimplemented; `side` experimental.
