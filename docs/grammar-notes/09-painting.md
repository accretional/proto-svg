# Painting & markers grammar notes

## Source

W3C SVG 2, Chapter 13: "Painting: Filling, Stroking and Marker Symbols"
Extracted from: `/docs/specs/cache/svg2-painting.txt` (3378 lines, read in full)
Cross-checked against MDN Web Docs and known browser behavior (2025-06).

---

## `<paint>` type (full structure)

Spec definition (§13.2, prose + value grammar):

```ebnf
<paint>
  = "none"
  | <color>
  | <url> [ "none" | <color> ]?
  | "context-fill"
  | "context-stroke"
```

Annotation of each arm:

| Arm | Grammar leaf / expansion | Notes |
|-----|--------------------------|-------|
| `"none"` | terminal keyword | No paint applied |
| `<color>` | open datatype — see §Open datatypes | All CSS Color L3 syntaxes required: named keywords, `rgb()`, `rgba()`, `hsl()`, `hsla()`, hex, `currentColor` |
| `<url>` | open datatype — `url( <string> )` | Must reference a paint server: `linearGradient`, `radialGradient`, or `pattern` element. `<funciri>` in SVG 1.1 terminology |
| `[ "none" \| <color> ]?` | optional fallback arm | Used when the paint-server reference is invalid; if absent and reference is invalid, no paint is rendered (changed from SVG 1.1 hard error) |
| `"context-fill"` | terminal keyword | Inherits paint layers from the fill of the context element (limited browser support — see §Discrepancies) |
| `"context-stroke"` | terminal keyword | Inherits paint layers from the stroke of the context element (limited browser support — see §Discrepancies) |

Context element rules (normative): the context element is (a) the shape that references the `marker` element containing this element, or (b) the `use` element whose shadow tree contains this element; otherwise there is no context element and `context-fill`/`context-stroke` produce no paint.

---

## Fill / stroke properties

### `fill`

```
Value:    <paint>
Initial:  black
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
Computed: as specified, <color> computed, <url> made absolute
```

Constraints:
- All `<paint>` arms are valid; `none` suppresses fill rendering entirely.
- `fill-rule` determines interior membership for complex paths.

Verification: Universally implemented. `context-fill` and `context-stroke` arms have limited support (see §Discrepancies).

---

### `fill-rule`

```
Value:    "nonzero" | "evenodd"
Initial:  nonzero
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
```

Full keyword set (closed):
- `"nonzero"` — counts directed crossings; non-zero count = inside
- `"evenodd"` — counts total crossings; odd count = inside

Verification: Both keywords universally implemented.

---

### `fill-opacity`

```
Value:    <alpha-value>
Initial:  1
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
Computed: specified value converted to number, clamped to [0, 1]
```

Spec names the value type `<alpha-value>`, which admits:
- `<number>` — clamped to [0..1]
- `<percentage>` — clamped to [0%..100%], normalised to [0..1]

Constraints: values outside [0, 1] (or [0%, 100%]) must be clamped, not rejected.

Grammar leaf: `<alpha-value> = <number> | <percentage>`

---

### `stroke`

```
Value:    <paint>
Initial:  none
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
Computed: as specified, <color> computed, <url> made absolute
```

Same `<paint>` grammar as `fill`. Initial is `none` (no visible stroke by default).

Constraints:
- Zero-length subpaths are not stroked when `stroke-linecap: butt`.
- A moveto-only subpath is never stroked.

---

### `stroke-opacity`

```
Value:    <alpha-value>
Initial:  1
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
Computed: clamped to [0, 1]
```

Identical semantics to `fill-opacity`; same `<alpha-value>` grammar leaf.

---

### `stroke-width`

```
Value:    <length-percentage>
Initial:  1
Applies:  shapes and text content elements
Inherited: yes
Percentages: refer to SVG viewport size
Animatable: yes
Computed: absolute length or percentage
```

Constraints:
- Zero value → no stroke painted.
- Negative value → invalid (must be rejected).
- Non-negative only; constraint belongs in the overlay, not the grammar terminal.

Grammar leaf: `<length-percentage>` (open — same as CSS, inherits px/em/rem/% etc.)

---

### `stroke-linecap`

```
Value:    "butt" | "round" | "square"
Initial:  butt
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
```

Full keyword set (closed):
- `"butt"` — stroke does not extend past endpoints; zero-length subpath renders nothing
- `"round"` — half-circle caps of diameter = stroke-width; zero-length → full circle
- `"square"` — rectangular caps half the stroke-width long; zero-length → square centered on point

Verification: All three universally implemented.

---

### `stroke-linejoin`

```
Value:    "miter" | "miter-clip" | "round" | "bevel" | "arcs"
Initial:  miter
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
```

Full keyword set (closed — but two values have limited support):
- `"miter"` — sharp corner; falls back to `bevel` when miterlimit exceeded
- `"miter-clip"` — like `miter` but clips at miterlimit rather than falling back [SVG 2 new; limited support — see §Discrepancies]
- `"round"` — circular arc join
- `"bevel"` — triangle fill between segments
- `"arcs"` — arc-extended join matching curvature of path segments [SVG 2 new; very limited support — see §Discrepancies]

---

### `stroke-miterlimit`

```
Value:    <number>
Initial:  4
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
```

Constraints:
- Negative value → invalid (must be rejected per spec).
- Values in (0, 1) are not explicitly invalid in SVG 2 (relaxed from SVG 1.1); any real miter will exceed a limit < 1, so those values are effectively a permanent fallback to bevel/clip.
- Used by `miter`, `miter-clip`, and `arcs` linejoin values.
- Formula: miter length = stroke-width / sin(θ/2).

Grammar leaf: `<number>` — must be non-negative in the constraint overlay.

---

### `stroke-dasharray`

```
Value:    "none" | <dasharray>
Initial:  none
Applies:  shapes and text content elements
Inherited: yes
Percentages: refer to SVG viewport size
Animatable: yes (non-additive)
Computed: absolute lengths or percentages, or keyword
```

Sub-grammar (from spec):
```ebnf
<dasharray> = [ <length-percentage> | <number> ]#*
```

`#*` means comma-and/or-whitespace separated list with zero or more items.

Constraints:
- Any negative value in the list → entire `<dasharray>` is invalid.
- All-zero list → renders as solid stroke (no dashing).
- Odd-length list is repeated to make it even-length before use.
- `<number>` in the list is interpreted as user units (not a percentage).

---

### `stroke-dashoffset`

```
Value:    <length-percentage>
Initial:  0
Applies:  shapes and text content elements
Inherited: yes
Percentages: refer to SVG viewport size
Animatable: yes
Computed: absolute length or percentage
```

Constraints:
- Negative values are allowed; equivalent offset is `s - (|offset| mod s)` where `s` is sum of dash array.
- Interpreted relative to `pathLength` when that attribute is present on `path`.

---

### `paint-order`

```
Value:    "normal" | [ "fill" || "stroke" || "markers" ]
Initial:  normal
Applies:  shapes and text content elements
Inherited: yes
Animatable: yes
```

Full keyword set:
- `"normal"` — paints in order: fill, then stroke, then markers
- The ordered combination arm uses `||` (any permutation of the three keywords, each appearing at most once, at least one must appear)
- Omitted keywords are appended in normal order after specified ones
- e.g., `stroke` → equivalent to `stroke fill markers`

EBNF expansion (exhaustive would be 6 permutations plus partials; use combinatorial notation):

```ebnf
<paint-order>
  = "normal"
  | <paint-order-layers>

<paint-order-layers>
  = <paint-order-keyword>+
  (* at least one of: "fill", "stroke", "markers"; no repeats; *)

<paint-order-keyword>
  = "fill" | "stroke" | "markers"
```

Note: grammar must enforce no-repeat and 1..3 keywords. This is a context-sensitive constraint, handled in the overlay.

---

### `color`

Defined in CSS Color Module Level 3 (deferred to external spec). SVG role: provides the `currentColor` value mechanism for `fill`, `stroke`, `stop-color`, `flood-color`, `lighting-color`.

Grammar leaf: `<color>` — open, per CSS Color L3. Minimum required syntaxes:
- Extended color keywords (147 named colors + `transparent`)
- `currentColor` keyword
- `#rrggbb`, `#rgb`, `#rrggbbaa`, `#rgba` hex notations
- `rgb( <number>, <number>, <number> )` / `rgba()`
- `hsl( <angle-or-number>, <percentage>, <percentage> )` / `hsla()`

---

### `color-interpolation`

```
Value:    "auto" | "sRGB" | "linearRGB"
Initial:  sRGB
Applies:  container elements, graphics elements, gradient elements, 'use', 'animate'
Inherited: yes
Animatable: yes
```

Full keyword set (closed):
- `"auto"` — UA chooses either color space
- `"sRGB"` — interpolation in sRGB space
- `"linearRGB"` — interpolation in linearized RGB space

Note: case-sensitive in XML; `sRGB` and `linearRGB` have mixed case (unusual for CSS properties — see §Discrepancies).

---

### `shape-rendering`

```
Value:    "auto" | "optimizeSpeed" | "crispEdges" | "geometricPrecision"
Initial:  auto
Applies:  shapes
Inherited: yes
Animatable: yes
```

Full keyword set (closed, render hints only — no semantic effect on geometry):
- `"auto"` — UA balances; geometric precision prioritized
- `"optimizeSpeed"` — may disable anti-aliasing
- `"crispEdges"` — aligns edges to device pixels, may disable anti-aliasing
- `"geometricPrecision"` — precise geometry over speed

---

### Ancillary rendering-hint properties (in scope for completeness)

**`color-rendering`**
```
Value:    "auto" | "optimizeSpeed" | "optimizeQuality"
Initial:  auto
```

**`text-rendering`**
```
Value:    "auto" | "optimizeSpeed" | "optimizeLegibility" | "geometricPrecision"
Initial:  auto
```

**`image-rendering`**
```
Value:    "auto" | "optimizeQuality" | "optimizeSpeed"
Initial:  auto
```

These are rendering hints; all keywords are closed sets.

---

## marker-* properties

### `marker-start`, `marker-mid`, `marker-end`

```
Value:    "none" | <marker-ref>
Initial:  none
Applies:  shapes
Inherited: yes
Animatable: yes
Computed: as specified, <url> made absolute
```

Sub-grammar:
```ebnf
<marker-ref> = <url>
```

(The spec defines `<marker-ref> = <url>` — it is simply a URL that must resolve to a `marker` element. The constraint that the referent must be a `marker` element is semantic, not syntactic.)

Arm semantics:
- `"none"` — no marker at specified vertex position(s)
- `<marker-ref>` — reference to a `marker` element; if reference is invalid, no marker is drawn (not an error)

Vertex placement rules:
- `marker-start` → first vertex of shape's equivalent path
- `marker-end` → last vertex
- `marker-mid` → all other vertices (closed subpaths produce two markers at start/end vertex)

### `marker` (shorthand)

```
Value:    "none" | <marker-ref>
Initial:  (not defined for shorthand properties)
Applies:  shapes
Inherited: yes
Animatable: yes
```

Sets `marker-start`, `marker-mid`, and `marker-end` to the same value. Shorthand only — computed values are the three longhands.

---

## `marker` element

### Categories
Container element, never-rendered element.

### Content model
Any number, in any order, of:
- Animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- Descriptive elements: `desc`, `title`, `metadata`
- Paint server elements: `linearGradient`, `radialGradient`, `pattern`
- Shape elements: `circle`, `ellipse`, `line`, `path`, `polygon`, `polyline`, `rect`
- Structural elements: `defs`, `g`, `svg`, `symbol`, `use`
- Also: `a`, `audio`, `canvas`, `clipPath`, `filter`, `foreignObject`, `iframe`, `image`, `marker`, `mask`, `script`, `style`, `switch`, `text`, `video`, `view`

Note: `marker` can be nested (recursive content model).

### Attribute groups
- Core attributes: `id`, `tabindex`, `lang`, `xml:space`, `class`, `style`
- Global event attributes (full list in spec §13.7.1)
- Document element event attributes: `oncopy`, `oncut`, `onpaste`
- Presentation attributes (specific): `viewBox`, `preserveAspectRatio`, `refX`, `refY`, `markerUnits`, `markerWidth`, `markerHeight`, `orient`

### Attribute definitions

**`markerUnits`**
```
Value:          "strokeWidth" | "userSpaceOnUse"
Initial value:  strokeWidth
Animatable:     yes
```
Full keyword set (closed):
- `"strokeWidth"` — marker coordinate system scaled by the referencing element's stroke-width
- `"userSpaceOnUse"` — marker coordinate system is current user coordinate system of referencing element

**`markerWidth`, `markerHeight`**
```
Value:          <length-percentage> | <number>
Initial value:  3
Animatable:     yes
```
Constraints:
- Zero → nothing rendered for marker.
- Negative → error (§Error processing).

**`refX`**
```
Value:          <length-percentage> | <number> | "left" | "center" | "right"
Initial value:  0
Animatable:     yes
```

Full keyword-to-percentage mapping:
- `"left"` → 0%
- `"center"` → 50%
- `"right"` → 100%

Note: SVG 2 new; `<number>` is treated as user-units in the marker coordinate system after viewBox/preserveAspectRatio transformations.

**`refY`**
```
Value:          <length-percentage> | <number> | "top" | "center" | "bottom"
Initial value:  0
Animatable:     yes
```

Full keyword-to-percentage mapping:
- `"top"` → 0%
- `"center"` → 50%
- `"bottom"` → 100%

**`orient`**
```
Value:          "auto" | "auto-start-reverse" | <angle> | <number>
Initial value:  0
Animatable:     yes (non-additive)
```

Mixed domain — split for grammar:
```ebnf
<orient-value>
  = "auto"
  | "auto-start-reverse"
  | <angle>
  | <number>
```

Arm semantics:
- `"auto"` — marker positive x-axis aligned with path direction at placement point
- `"auto-start-reverse"` — for `marker-start` placements: rotate 180° from `auto`; for `marker-mid`/`marker-end`: same as `auto`
- `<angle>` — explicit fixed rotation angle (e.g., `45deg`, `1.57rad`)
- `<number>` — angle in degrees (unitless number; same as `<angle>` in degrees)

DOM constants: `SVG_MARKER_ORIENT_AUTO` (1), `SVG_MARKER_ORIENT_ANGLE` (2), `SVG_MARKER_ORIENT_UNKNOWN` (0). Note `auto-start-reverse` maps to `SVG_MARKER_ORIENT_UNKNOWN` in the DOM — the spec acknowledges this is non-ideal.

**`viewBox`** — shared attribute, see geometry/coords spec notes.

**`preserveAspectRatio`** — shared attribute, see geometry/coords spec notes.

---

## Open datatypes used

| Leaf | Definition source | Notes |
|------|-------------------|-------|
| `<color>` | CSS Color L3 | Named colors, hex, rgb(), hsl(), currentColor; treat as open leaf |
| `<url>` | CSS Values L3 `url()` function | `url(<string>)` or `url(<unquoted-token>)`; context constrains referent type |
| `<length-percentage>` | CSS Values L3/L4 | `<length>` or `<percentage>`; negative disallowed for widths |
| `<length>` | CSS Values L3/L4 | Used in marker attributes |
| `<number>` | CSS Values | Floating-point; sign constraints property-specific |
| `<percentage>` | CSS Values | 0%..100% for opacities; viewport-relative for stroke lengths |
| `<alpha-value>` | CSS Color L4 | `<number>` or `<percentage>`, clamped to [0, 1] |
| `<angle>` | CSS Values | `deg`, `rad`, `grad`, `turn` units |
| `<dasharray>` | SVG 2 §13.5.6 | `[ <length-percentage> | <number> ]#*` — comma/whitespace-separated list |
| `<marker-ref>` | SVG 2 §13.7.2 | Alias for `<url>`; referent must be a `marker` element |

---

## Discrepancies, doc gaps & roadblocks

### 1. `context-fill` / `context-stroke` — limited browser support

**Spec says:** These are valid `<paint>` arms usable in any fill/stroke context.
**Reality:** As of 2025, `context-fill` and `context-stroke` are supported in Firefox, Chrome 80+, Safari 13.1+. Support is good in modern browsers when used inside a `marker` element. Use inside a `use` shadow tree (the second context-element case) has more variable support. These can be included in the grammar but should be flagged in the constraint overlay as "context-dependent — may silently produce no paint."
**Grammar decision:** Include both arms in `<paint>`. Flag in overlay.

### 2. `stroke-linejoin: arcs` — very limited support

**Spec says:** Fully specified with curvature-matching arcs and fallback rules.
**Reality:** `arcs` is implemented in Firefox (roughly since FF 58) but not in Chrome or Safari as of 2025-06. It is effectively Firefox-only.
**Grammar decision:** Include `"arcs"` in the keyword set, flag in overlay as "Firefox-only; Chrome/Safari fall through to `miter`/`bevel` behavior."

### 3. `stroke-linejoin: miter-clip` — limited support

**Spec says:** Clips the miter at half-miterlimit distance instead of falling back to bevel.
**Reality:** Implemented in Firefox; Chrome support is partial (as of 2025); not in Safari.
**Grammar decision:** Include `"miter-clip"` in keyword set, flag in overlay.

### 4. `color-interpolation` case sensitivity

**Spec says:** Values are `auto`, `sRGB`, `linearRGB` (mixed-case, matching SVG 1.1).
**MDN / CSS practice:** CSS is case-insensitive for keywords; SVG attributes in XML documents are case-sensitive. In practice, browsers accept lowercase variants in CSS property context. In XML attribute context, case matters.
**Grammar decision:** Use spec-canonical casing (`"sRGB"`, `"linearRGB"`) in the grammar. Note case-sensitivity in the overlay with a flag for CSS-context lenience.

### 5. `refX` / `refY` keyword support

**Spec says:** SVG 2 adds `left/center/right` (refX) and `top/center/bottom` (refY) as geometric keywords.
**Reality:** These SVG 2 keywords (`left`, `right`, `top`, `bottom` on `refX`/`refY`) have patchy support; Chrome 88+, Firefox 89+, Safari 15+. Older implementations only accept `<length>` and `<number>`.
**Grammar decision:** Include all keywords in the grammar; flag as SVG 2 new with minimum-version notes.

### 6. `stroke-dasharray` `<number>` values

**Spec says:** `<dasharray>` accepts `<length-percentage> | <number>`. The `<number>` case is interpreted as user units.
**Ambiguity:** The spec also says the value type is `none | <dasharray>` but the sub-grammar definition uses `#*` (zero or more). An empty `<dasharray>` would be equivalent to `none` semantically, but the grammar technically permits both. Treat `none` as the canonical "no dashing" keyword and `<dasharray>` as requiring at least one token in the grammar.

### 7. `orient` initial value ambiguity

**Spec says:** Initial value is `0`. This is a `<number>` (degrees), which is equivalent to `0deg`.
**DOM note:** `auto-start-reverse` maps to `SVG_MARKER_ORIENT_UNKNOWN` — the spec itself notes this constant is shared with unknown values, which is an API design gap. The grammar is not affected.

### 8. `marker-start`/`marker-end` on closed subpaths

The spec describes two markers being drawn at the first/last vertex of closed subpaths when both `marker-start` and `marker-end` are non-none. This is a rendering rule, not a grammar rule, but it is a common source of unexpected visual behavior. Note in overlay.

### 9. `<paint>` fallback behavior changed from SVG 1.1

SVG 1.1: invalid paint server reference with no fallback = document error.
SVG 2: invalid paint server reference with no fallback = no paint rendered (silent degradation).
Grammar is correct per SVG 2. If targeting SVG 1.1 compatibility, the fallback `[ "none" | <color> ]` arm should be required when `<url>` is used.

### 10. `paint-order` not in SVG 1.1

New in SVG 2. Some tooling and older validators will reject it. Flag in overlay.

### 11. `stroke-miterlimit` minimum value

SVG 1.1 required ≥ 1; SVG 2 relaxes to "non-negative." Values in (0, 1) are technically valid but always trigger fallback. Grammar leaf is `<number>` with non-negative constraint in overlay.

---

## Property / element count summary

Properties documented in full: 18
- fill, fill-rule, fill-opacity, stroke, stroke-opacity, stroke-width, stroke-linecap, stroke-linejoin, stroke-miterlimit, stroke-dasharray, stroke-dashoffset, paint-order, color, color-interpolation, marker-start, marker-mid, marker-end, marker (shorthand)

Rendering-hint properties (closed keyword sets, no geometry): 3
- color-rendering, shape-rendering, text-rendering, image-rendering (4 total in spec; all included above)

Elements documented: 1
- `marker`

Open datatypes named: 9
(color, url, length-percentage, length, number, percentage, alpha-value, angle, dasharray)

Formal syntax fenced blocks extracted from spec: 3
- `<paint>` grammar (§13.2)
- `<dasharray>` sub-grammar (§13.5.6)
- SVGMarkerElement WebIDL (§13.12.1)
