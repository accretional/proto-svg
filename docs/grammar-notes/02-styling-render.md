# Styling & Render Grammar Notes

## Source

- `svg2-styling.txt` — W3C SVG 2 Chapter 6: Styling (full read, 703 lines)
- `svg2-render.txt` — W3C SVG 2 Chapter 3: Rendering Model (full read, 864 lines)
- Cross-checked against `/Volumes/wd_office_3/Projects/proto-svg/docs/specs/mdn_docs_attributes.md`

---

## Elements

### `style`

The `style` element embeds inline CSS style sheets directly in SVG content.
SVG 2 aligns the element completely with the HTML `style` element.

**Category:** Never-rendered element (UA stylesheet forces `display: none !important`).

**Content model:** Character data (raw CSS text). Not a value syntax in the SVG
sense — the entire text content is handed to the CSS parser as-is.

**Attributes:**

| Attribute | Value syntax | Initial | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `type` | `<media-type>` (MIME type string) | `text/css` | no | If absent, assumed `text/css`. Any MIME type is syntactically valid; only `text/css` is processed by SVG UAs. |
| `media` | `<media-query-list>` (CSS media query grammar) | `all` | no | Parsed as CSS `media_query_list`. If absent, applies unconditionally. |
| `title` | any string | (none) | no | Used for alternate style sheet selection. Unconstrained string. |

The `type`, `media`, and `title` IDL attributes on `SVGStyleElement` reflect
the corresponding content attributes.

**EBNF note:** The `style` element itself has no value-syntax grammar in SVG —
it is a transparent pass-through to CSS. Its attributes reference external
grammars (MIME type, CSS media query list). They are open-leaf references.

---

## Presentation-attribute Mechanism

### What a presentation attribute is

Presentation attributes are XML attributes whose name matches a CSS property
name. They are parsed as a CSS *value* (not a declaration — `!important` is
rejected). They contribute to the author-level cascade at specificity 0,
after all other author-level style sheets. Only elements in the SVG namespace
support them.

### Attribute name aliasing (transform group)

Three attribute names map to the same `transform` CSS property:

| Presentation attribute | Elements | CSS property |
|------------------------|----------|--------------|
| `transform` | All SVG namespace elements except `pattern`, `linearGradient`, `radialGradient` | `transform` |
| `patternTransform` | `pattern` | `transform` |
| `gradientTransform` | `linearGradient`, `radialGradient` | `transform` |

### Geometry attributes as conditional presentation attributes

The following attributes act as presentation attributes *only on designated elements*:

| Attribute(s) | Designated elements |
|-------------|---------------------|
| `cx`, `cy` | `circle`, `ellipse` |
| `height`, `width`, `x`, `y` | `foreignObject`, `image`, `rect`, `svg`, `symbol`, `use` |
| `r` | `circle` |
| `rx`, `ry` | `ellipse`, `rect` |
| `d` | `path` |

On any other element the same attribute name is not a presentation attribute
and must not affect the CSS cascade.

---

## Presentation-attribute List (Complete)

The following are the complete set of CSS-backed presentation attributes
defined by SVG 2 (§6.6), applicable to any element in the SVG namespace
unless constrained by the table above. Attributes that are purely SVG-domain
CSS properties (defined in the SVG spec itself) are marked **[SVG]**; those
that are reused from other CSS modules are marked with the source spec.

| Presentation attribute | Source spec | MDN listed? | MDN status |
|------------------------|-------------|-------------|------------|
| `alignment-baseline` | css-inline-3 | yes | — |
| `baseline-shift` | css-inline-3 | yes | — |
| `clip-path` | css-masking-1 | yes | — |
| `clip-rule` | css-masking-1 | yes | — |
| `color` | css-color-3 | yes | — |
| `color-interpolation` | **[SVG]** | yes | — |
| `color-interpolation-filters` | filter-effects-1 | yes | — |
| `color-rendering` | **[SVG]** | not listed separately (merged into `color-interpolation` page) | — |
| `cursor` | css-ui-3 | yes | — |
| `d` | **[SVG geometry]** | yes | — |
| `direction` | css-writing-modes-3 | yes | — |
| `display` | CSS2 | yes | — |
| `dominant-baseline` | css-inline-3 | yes | — |
| `fill` | **[SVG]** | yes | — |
| `fill-opacity` | **[SVG]** | yes | — |
| `fill-rule` | **[SVG]** | yes | — |
| `filter` | filter-effects-1 | yes | — |
| `flood-color` | filter-effects-1 | yes | — |
| `flood-opacity` | filter-effects-1 | yes | — |
| `font-family` | css-fonts-3 | yes | — |
| `font-size` | css-fonts-3 | yes | — |
| `font-size-adjust` | css-fonts-3 | yes | — |
| `font-stretch` | css-fonts-3 | yes | — |
| `font-style` | css-fonts-3 | yes | — |
| `font-variant` | css-fonts-3 | yes | — |
| `font-weight` | css-fonts-3 | yes | — |
| `glyph-orientation-horizontal` | **[SVG 1.1 legacy]** | yes | ⚠️ deprecated |
| `glyph-orientation-vertical` | **[SVG 1.1 legacy]** | yes | ⚠️ deprecated |
| `gradientTransform` | **[SVG, alias for transform]** | yes | — |
| `image-rendering` | **[SVG]** | yes | — |
| `letter-spacing` | css-text-3 | yes | — |
| `lighting-color` | filter-effects-1 | yes | — |
| `marker-end` | **[SVG]** | yes | — |
| `marker-mid` | **[SVG]** | yes | — |
| `marker-start` | **[SVG]** | yes | — |
| `mask` | css-masking-1 | yes | — |
| `opacity` | css-color-3 | yes | — |
| `overflow` | CSS2 | yes | — |
| `paint-order` | **[SVG]** | yes | — |
| `patternTransform` | **[SVG, alias for transform]** | yes | — |
| `pointer-events` | **[SVG]** | yes | — |
| `shape-rendering` | **[SVG]** | yes | — |
| `stop-color` | **[SVG]** | yes | — |
| `stop-opacity` | **[SVG]** | yes | — |
| `stroke` | **[SVG]** | yes | — |
| `stroke-dasharray` | **[SVG]** | yes | — |
| `stroke-dashoffset` | **[SVG]** | yes | — |
| `stroke-linecap` | **[SVG]** | yes | — |
| `stroke-linejoin` | **[SVG]** | yes | — |
| `stroke-miterlimit` | **[SVG]** | yes | — |
| `stroke-opacity` | **[SVG]** | yes | — |
| `stroke-width` | **[SVG]** | yes | — |
| `text-anchor` | **[SVG]** | yes | — |
| `text-decoration` | css-text-decor-3 | yes | — |
| `text-overflow` | css-ui-3 | yes | — |
| `text-rendering` | **[SVG]** | yes | — |
| `transform` | css-transforms-1 | yes | — |
| `unicode-bidi` | css-writing-modes-3 | yes | — |
| `vector-effect` | **[SVG]** | yes | — |
| `visibility` | CSS2 | yes | — |
| `white-space` | css-text-4 | yes | — |
| `word-spacing` | css-text-3 | yes | — |
| `writing-mode` | css-writing-modes-3 | yes | — |

**Total: 59 presentation attributes** (57 CSS-property-named + `patternTransform` + `gradientTransform` as aliases).
Plus 9 conditional geometry attrs (`cx`, `cy`, `r`, `rx`, `ry`, `x`, `y`, `width`, `height`, `d` — 10 names).

**Attributes in SVG 2 spec list NOT in MDN attribute index:**
- `color-rendering` — MDN does not list it as a standalone attribute entry.
  (MDN covers color-rendering as a CSS property for SVG but may have merged it.)

**Attributes in MDN NOT in the SVG 2 §6.6 presentation-attribute list:**
- `mask-type` — listed in MDN but is a CSS Masking property, not an SVG presentation attribute per §6.6. It is a CSS-only property, not granted a presentation attribute in SVG 2.
- `cx`, `cy`, `r`, `rx`, `ry`, `x`, `y`, `width`, `height`, `d`, `transform`, `patternTransform`, `gradientTransform` — all present in MDN; all accounted for above.
- `transform-origin` — MDN lists it as an SVG attribute; SVG 2 §6.7 says it is *required* but does NOT list it as a presentation attribute in §6.6 (no element-specific restriction). Cross-check: browsers do accept `transform-origin` as an attribute on SVG elements. This is a gap — see discrepancies.

---

## Properties Defined Here (Rendering Model Focus)

This section covers properties whose value syntax and semantics are defined or
significantly constrained by the SVG 2 spec or the Rendering Model chapter.
For properties whose value syntax is fully delegated to an external CSS module
(e.g., `font-family`, `letter-spacing`), only the grammar leaf reference is
noted; their full syntax is extracted in the CSS-module grammar notes.

Properties whose full enumeration appears below: `display`, `visibility`,
`opacity`, `overflow`, `paint-order`.

Compositing properties (`isolation`, `mix-blend-mode`) are *required* by SVG 2
§6.7 (via compositing-1) but their syntax is defined entirely in
CSS Compositing and Blending — they are leaf references here, not expanded.

---

### `display`

**Source:** CSS 2.1, referenced by SVG 2 §3.2.3 and §6.7.

SVG 2 restricts `display` usage: only the value `none` has a special SVG
semantic (removes the element from the rendering tree entirely). Any other
valid CSS `display` value leaves the element rendered. The spec text says
"elements that have any other display value than `none` are rendered as
normal."

Never-rendered elements have `display: none !important` forced by the UA
stylesheet — this cannot be overridden.

**Value syntax (SVG-relevant subset):**

```ebnf
display
  = "none"
  | display-outside [ display-inside ]?
  | display-legacy
  | "inherit"
  | "initial"
  | "unset"

display-outside = "block" | "inline" | "run-in"

display-inside
  = "flow" | "flow-root" | "table"
  | "flex" | "grid" | "ruby"

display-legacy
  = "inline-block" | "inline-table"
  | "inline-flex" | "inline-grid"
```

**Grammar decision:** For SVG grammar purposes, the full CSS display value set
is an open-leaf (`<display-value>`) — authors may use any CSS-valid display
value. Only `none` has a named SVG terminal because it is the only value with
SVG-specific semantics. All other values collapse to the same "rendered"
behavior in SVG context.

**Initial value:** Per CSS2, the property initial is `inline`. However, the
SVG UA stylesheet sets `display: none !important` on all never-rendered
elements. For renderable elements, initial is effectively `inline`.

**Constraints (overlay, not in value syntax):**
- `display: none !important` is unconditionally applied to: `clipPath`, `defs`,
  `desc`, `linearGradient`, `marker`, `mask`, `metadata`, `pattern`,
  `radialGradient`, `script`, `style`, `title`, and `symbol` (unless it is the
  instance root of a use-element shadow tree).
- `symbol` inside a `use` shadow root gets `display: inline !important`.
- Setting `display` to any non-`none` value on a never-rendered element has
  no effect (the UA rule wins with `!important`).

---

### `visibility`

**Source:** CSS 2.1, referenced by SVG 2 §3.2.3.

```ebnf
visibility = "visible" | "hidden" | "collapse" | "inherit"
```

**Initial value:** `visible`

**All keywords are closed terminals** — enumeration is complete.

**SVG-specific semantics (overlay):**
- `visibility: hidden` or `visibility: collapse` on a graphics element or
  `use` element means the element is not *painted* but it remains in the
  rendering tree.
- Unlike `display: none`, a hidden element still contributes to bounding box
  calculations, clipping paths, and text layout.
- `visibility` is an inherited property; the inherited value affects
  descendants even though setting it on a container element (`g`, `use`) has
  no direct effect on that container's own rendering.
- `visibility` only directly affects rendering of graphics elements, text
  content elements, and the `a` element when it is a child of a text content
  element.

**Note:** In CSS2 and CSS3, `collapse` has no visual distinction from `hidden`
for non-table elements. In SVG, both `hidden` and `collapse` produce the same
effect; `collapse` is grammatically valid but semantically redundant.

---

### `opacity`

**Source:** CSS Color Module Level 3, referenced by SVG 2 §3.6.1 and §6.7.

```ebnf
opacity = <alpha-value> | "inherit"

(* <alpha-value> is a number — see constraints *)
```

**Initial value:** `1`

**Open domain:** `<alpha-value>` is a `<number>` leaf.

**Constraints (overlay):**
- Valid range: `0` to `1` inclusive (clamped by UA — values outside this range
  are clamped, not rejected).
- `opacity: 1` means fully opaque; `opacity: 0` means fully transparent.
- Applying `opacity` to a container element creates group opacity (the group
  is composited as a unit before being blended into the background).
- When `opacity` is not `1`, the element (or group) becomes an isolated
  compositing group (stacking context established).
- The `opacity` property applies to: `svg`, `g`, `symbol`, `marker`, `a`,
  `switch`, `use`, `unknown`, and all graphics elements.
- `opacity` is distinct from `fill-opacity`, `stroke-opacity`, and
  `stop-opacity` — all four may be combined multiplicatively.

---

### `overflow`

**Source:** CSS 2.1 §11.1.1, with SVG 2 §3.11 additions.

```ebnf
overflow = "visible" | "hidden" | "scroll" | "auto" | "inherit"
```

**Initial value:** `auto` (per CSS2).

**All keywords are closed terminals** — enumeration is complete.

**SVG-specific semantics and UA overrides (overlay):**

Per the SVG 2 rendering model table (§3.11):

| Element | UA stylesheet | Effect of `auto` | Effect of `scroll` |
|---------|--------------|-------------------|---------------------|
| document root `svg` | n/a (visible is default) | visible or scroll (UA choice) | scroll mechanism shown |
| other `svg` | `hidden` | visible or scroll | scroll mechanism shown |
| `text` | `hidden` | visible | hidden (clamped) |
| `pattern` | `hidden` | visible | hidden (clamped) |
| `marker` | `hidden` | visible | hidden (clamped) |
| `symbol` | `hidden` | visible | hidden (clamped) |
| `image` | `hidden` | visible | hidden (clamped) |
| `iframe` | `hidden` | visible or scroll | scroll mechanism shown |
| `foreignObject` | `hidden` | visible or scroll | scroll mechanism shown |

**Additional SVG rules:**
- `overflow: visible` on an SVG viewport element has no effect (no clipping
  rectangle created).
- `overflow: hidden` or `overflow: scroll` clips to the exact SVG viewport
  size.
- Within SVG, `auto` implies all rendered content is visible (either by
  scrolling mechanism or no clip). If UA has no scrolling mechanism, `auto`
  maps to `visible`.
- An inner `svg` element with `overflow` other than `visible` establishes a
  stacking context.

**Note:** `overflow: clip` (CSS Overflow Level 3) is NOT mentioned in SVG 2 —
do not add as a terminal. Modern browsers do implement it for SVG; this is a
discrepancy (see §Discrepancies).

---

### `paint-order`

**Source:** SVG 2 (SVG-specific property). Referenced in §3.7.1 and §6.6.

```ebnf
paint-order
  = "normal"
  | paint-order-layer+

paint-order-layer = "fill" | "stroke" | "markers"
```

**Semantic rule:** `paint-order-layer+` is an ordered permutation of a subset
of the three layers. Omitted layers are appended in default order after those
listed. Maximum 3 tokens; no repetition.

**Initial value:** `normal`

**`normal` is a closed terminal.** The three paint layer keywords are closed
terminals. The set {`fill`, `stroke`, `markers`} is the complete closed set.

**Constraints (overlay):**
- `normal` is shorthand for `fill stroke markers` (paint fill first, then
  stroke, then markers).
- When a layer keyword is listed, it is painted first. Unlisted layers follow
  in their default relative order.
- The grammar allows 1, 2, or 3 layer keywords; permutations are unrestricted.
- No layer keyword may appear more than once in a single value.

---

### `isolation`

**Source:** CSS Compositing and Blending Level 1. Required by SVG 2 §6.7.
Value syntax delegated entirely to compositing-1.

```
isolation = "auto" | "isolate" | "inherit"
```

**Initial value:** `auto`

Note: SVG 2 does not define the value syntax here — it only mandates support.
Full grammar in compositing-1 notes. Not listed as a presentation attribute
in SVG 2 §6.6 (no XML attribute equivalent) — CSS-only property in SVG
context.

---

### `mix-blend-mode`

**Source:** CSS Compositing and Blending Level 1. Required by SVG 2 §6.7.
Value syntax delegated entirely to compositing-1.

Not listed as a presentation attribute in SVG 2 §6.6. CSS-only property.

---

### `color-rendering`

**Source:** SVG 2 (SVG-specific). Listed in §6.6 presentation-attribute table.

```ebnf
color-rendering
  = "auto" | "optimizeSpeed" | "optimizeQuality" | "inherit"
```

**Initial value:** `auto`

All keywords are closed terminals. This is a hint to the UA; behavior is
implementation-defined beyond the enumeration.

---

### `image-rendering`

**Source:** SVG 2 (SVG-specific property, later adopted by CSS Images Level 3).
Listed in §6.6.

```ebnf
image-rendering
  = "auto" | "optimizeSpeed" | "optimizeQuality"
  | "crisp-edges" | "pixelated"
  | "inherit"
```

**Initial value:** `auto`

**Note:** SVG 1.1 defined `optimizeSpeed` and `optimizeQuality`. CSS Images
Level 3 added `crisp-edges` and `pixelated` as standardized replacements.
Modern browsers accept all four non-auto keywords. Grammar includes all.

---

### `shape-rendering`

**Source:** SVG 2 (SVG-specific). Listed in §6.6.

```ebnf
shape-rendering
  = "auto" | "optimizeSpeed" | "crispEdges" | "geometricPrecision"
  | "inherit"
```

**Initial value:** `auto`

All keywords are closed terminals. Note the camelCase `crispEdges` and
`geometricPrecision` — contrast with `image-rendering`'s `crisp-edges`
(kebab-case). This is an existing spec inconsistency.

---

### `text-rendering`

**Source:** SVG 2 (SVG-specific). Listed in §6.6.

```ebnf
text-rendering
  = "auto" | "optimizeSpeed" | "optimizeLegibility"
  | "geometricPrecision"
  | "inherit"
```

**Initial value:** `auto`

All keywords are closed terminals.

---

### `pointer-events`

**Source:** SVG 2 (SVG-specific). Listed in §6.6.

```ebnf
pointer-events
  = "bounding-box"
  | "visiblePainted" | "visibleFill" | "visibleStroke" | "visible"
  | "painted" | "fill" | "stroke" | "all"
  | "none"
  | "inherit"
```

**Initial value:** `visiblePainted`

All keywords are closed terminals. SVG-specific values use camelCase;
CSS Pointer Events spec added `none` and `all` (also used in HTML).

---

### `vector-effect`

**Source:** SVG 2 (SVG-specific). Listed in §6.6.

```ebnf
vector-effect = "none" | "non-scaling-stroke" | "inherit"
```

**Initial value:** `none`

**Note:** SVG 2 may add `non-scaling-size`, `non-rotation`, and
`fixed-position` as additional values, but these are not in the stable text
and browsers do not implement them. Grammar uses only `none` and
`non-scaling-stroke`.

---

### `color-interpolation`

**Source:** SVG 2 (SVG-specific). Listed in §6.6.

```ebnf
color-interpolation
  = "auto" | "sRGB" | "linearRGB"
  | "inherit"
```

**Initial value:** `sRGB`

All keywords are closed terminals. Note the mixed-case `sRGB` and `linearRGB`.
MDN confirms these exact keyword values.

---

### `color-interpolation-filters`

**Source:** filter-effects-1, listed in SVG 2 §6.6.

Same grammar as `color-interpolation`:

```ebnf
color-interpolation-filters
  = "auto" | "sRGB" | "linearRGB"
  | "inherit"
```

**Initial value:** `linearRGB` (differs from `color-interpolation` — notable).

---

### `glyph-orientation-horizontal` and `glyph-orientation-vertical`

**Source:** SVG 1.1, retained in SVG 2 §6.6 but deprecated.

MDN status: ⚠️ deprecated. Browsers may still accept them.

```ebnf
glyph-orientation-horizontal = <angle> | "inherit"
glyph-orientation-vertical   = "auto" | <angle> | "inherit"
```

`<angle>` here is an integer multiple of 90 degrees (0, 90, 180, 270) per
SVG 1.1 — a constrained open domain.

**Constraints (overlay):**
- Values must be multiples of 90 degrees (0, 90, 180, 270).
- `glyph-orientation-vertical: auto` allows the UA to choose based on writing mode.

---

### Opacity-related sub-properties

These are listed as presentation attributes (§6.6) and have an `<alpha-value>`
open domain, same range constraint as `opacity`:

- `fill-opacity` — `<alpha-value> | "inherit"`, initial `1`
- `stroke-opacity` — `<alpha-value> | "inherit"`, initial `1`
- `stop-opacity` — `<alpha-value> | "inherit"`, initial `1`
- `flood-opacity` — `<alpha-value> | "inherit"`, initial `1`

All four: `<alpha-value>` is a `<number>` clamped to [0, 1].

---

## Open Datatypes Used

The following named open leaves appear in property value syntaxes in this
module. Their internal grammar is defined elsewhere:

| Leaf name | Description | Defined in |
|-----------|-------------|------------|
| `<alpha-value>` | A `<number>` in the range [0, 1] (clamped) | css-color-3 |
| `<angle>` | CSS angle value (number + optional unit) | css-values-3 |
| `<media-query-list>` | CSS media query list grammar | css-mediaqueries-4 |
| `<media-type>` | MIME type string | RFC 2046 |
| `<display-value>` | Full CSS display value (any CSS-valid value) | css-display-3 |
| `<number>` | Floating-point numeric value | css-values-3 |

---

## Discrepancies, Doc Gaps & Roadblocks

### D1: `color-rendering` absent from MDN attribute index

**Spec says:** `color-rendering` is listed in the §6.6 presentation-attribute
table as applicable to any SVG namespace element.

**MDN says:** MDN does not list `color-rendering` as a standalone entry in its
SVG attribute reference index (the mdn_docs_attributes.md cross-check shows
no row for it). MDN does document it as a CSS property for SVG content in the
CSS reference.

**Decision:** Keep `color-rendering` as a valid presentation attribute in the
grammar. It is widely supported in browsers (Chrome, Firefox, Safari all
implement it). The MDN omission appears to be a documentation gap, not an
implementation gap. Note the absence in the grammar's attribute catalog.

---

### D2: `transform-origin` — presentation attribute status unclear

**Spec says:** SVG 2 §6.7 requires support for `transform-origin` (css-transforms-1)
and §6.8 (UA stylesheet) specifies `transform-origin: 0 0` for most SVG
elements. However, `transform-origin` is NOT listed in the §6.6
presentation-attribute table.

**MDN says:** MDN lists `transform-origin` as an SVG attribute (in the MDN
attribute reference).

**Browser behavior:** Browsers accept `transform-origin` as an XML attribute
on SVG elements.

**Decision:** This is an SVG 2 spec gap. The property is required and the UA
stylesheet sets it, but no presentation attribute was explicitly chartered.
Browsers treat it as a de facto presentation attribute. In the grammar:
mark `transform-origin` as a presentation attribute with a note that it is
not formally in the §6.6 table but is implemented universally.

---

### D3: `overflow: clip` not in SVG 2 spec

**Spec says:** SVG 2 §3.11 defines `overflow` values as: `visible`, `hidden`,
`scroll`, `auto` (from CSS2). The value `clip` (CSS Overflow Level 3) is not
mentioned.

**Browser behavior:** Modern browsers (Chrome 90+, Firefox 102+) accept
`overflow: clip` on SVG elements.

**Decision:** Do NOT include `clip` as a terminal in the SVG 2 grammar.
Note as a browser extension. If authoring for modern-browser targets only,
`clip` could be added; for strict SVG 2 conformance, omit.

---

### D4: `isolation` and `mix-blend-mode` are required but have no presentation attributes

**Spec says:** SVG 2 §6.7 requires support for `isolation` (compositing-1).
These properties are not in the §6.6 presentation-attribute table.

**Implication:** These can only be set via CSS (style sheet or `style`
attribute). They cannot be set as XML attributes.

**Decision:** Do not create presentation-attribute grammar rules for
`isolation` or `mix-blend-mode`. Grammar note: CSS-only.

---

### D5: `display` in SVG — CSS3/CSS4 display values vs. CSS2 values

**Spec says:** SVG 2 §3.2.3 references "CSS 2.1" for `display` definition but
was written when CSS2 display values were the norm. The full CSS Display
Level 3 multi-keyword syntax (`block flow`, `inline flow-root`, etc.) is not
discussed.

**Browser behavior:** Browsers apply CSS Display Level 3 parsing to SVG
elements; multi-keyword values work.

**Decision:** Treat `<display-value>` as an open leaf (full CSS Display Level 3
grammar). Only `none` is a named SVG terminal. All other values reduce to the
same "rendered" behavior in SVG.

---

### D6: `glyph-orientation-horizontal` and `glyph-orientation-vertical` are deprecated

**Spec:** Still listed in §6.6. MDN marks both ⚠️ deprecated.

**Browser behavior:** Broadly ignored by modern layout engines in SVG 2 mode.
`glyph-orientation-vertical` is mapped to `text-orientation` in CSS Writing
Modes Level 3 for some values.

**Decision:** Include in grammar as deprecated terminals with a note. Do not
rely on them for new content. Grammar should accept them for backward
compatibility with SVG 1.1 documents.

---

### D7: `paint-order: markers` — spelling discrepancy

**Spec §6.6 listing:** The property is named `paint-order` (correct).

**MDN:** MDN documents the keyword as `markers` (plural).

**SVG 1.1 / early SVG 2 drafts:** Some drafts used `marker` (singular).

**Decision:** Use `markers` (plural) per MDN and per browser implementations.
The spec text in §3.7.1 uses "marker symbols" but the CSS property keyword
is `markers`.

---

### D8: `vector-effect` additional values in some browser implementations

**Spec:** SVG 2 stable text defines only `none` and `non-scaling-stroke`.

**Some browser implementations:** Chromium has implemented `non-scaling-size`
experimentally. This is not in the stable spec.

**Decision:** Grammar includes only `none` and `non-scaling-stroke`. Log
`non-scaling-size` as an experimental extension.

---

### D9: `color-interpolation-filters` initial value differs from `color-interpolation`

**Spec:** `color-interpolation-filters` has initial `linearRGB`; `color-interpolation`
has initial `sRGB`. This is intentional — filter operations default to linear
color space for photographic accuracy.

**Not a discrepancy** — correctly noted here to avoid confusion when authoring
grammar constraints. Both use the same keyword set.

---

### D10: `shape-rendering: crispEdges` vs. `image-rendering: crisp-edges` capitalization

**Spec:** `shape-rendering` uses camelCase `crispEdges`; `image-rendering`
uses kebab-case `crisp-edges`. This inconsistency is embedded in the spec and
cannot be resolved by the grammar — both must be faithfully reproduced.

**Decision:** Treat these as distinct closed terminals in their respective
property rules.

---

### D11: `mask-type` in MDN, not in SVG 2 §6.6 presentation-attribute table

**MDN lists** `mask-type` as an SVG attribute (CSS Masking property).
**SVG 2 §6.6 does not list** `mask-type` as a presentation attribute.

**Decision:** `mask-type` is a CSS Masking Level 1 property that applies to
`mask` elements. It has no SVG presentation attribute. MDN listing it in the
SVG attribute reference is because it is valid as a CSS property on the `mask`
element — but it is not a presentation attribute per §6.6. Grammar: CSS-only,
no presentation-attribute rule.

---

### D12: `fill` attribute collision with animation elements

**Spec §6.6 note:** The `fill` presentation attribute applies to "any element
in the SVG namespace *except animation elements*", which have a different
`fill` attribute (for animation fill behavior: `freeze`|`remove`).

**Grammar implication:** The `fill` attribute grammar rule must be split by
context:
- On animation elements (`animate`, `animateMotion`, `animateTransform`,
  `set`): `fill` = `"freeze"` | `"remove"` (animation fill mode, not paint).
- On all other SVG elements: `fill` = paint server value (color, url, etc.).

This is a critical context-sensitive constraint; the grammar must document it
as two separate rules with element-context disambiguation.

---

## Summary Counts

- Presentation attributes (§6.6 complete list): **59** (including geometry + transform aliases)
- Conditional geometry presentation attributes (element-restricted): **10** attribute names
- Properties with full closed keyword enumeration documented here: **13**
  (`display`[leaf], `visibility`, `opacity`[leaf], `overflow`, `paint-order`,
  `color-rendering`, `image-rendering`, `shape-rendering`, `text-rendering`,
  `pointer-events`, `vector-effect`, `color-interpolation`,
  `color-interpolation-filters`, `glyph-orientation-*`)
- Open-leaf domains introduced: **6** (`<alpha-value>`, `<angle>`, `<media-query-list>`, `<media-type>`, `<display-value>`, `<number>`)
- Discrepancies logged: **12** (D1–D12)
- Critical grammar constraints requiring context-sensitive handling: **2**
  (D12 `fill` collision; D2 `transform-origin` de-facto PA status)
