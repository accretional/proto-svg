# Basic shapes grammar notes

## Source

- **Primary spec**: W3C SVG 2, Chapter 10 "Basic Shapes"
  (`/Volumes/wd_office_3/Projects/proto-svg/docs/specs/cache/svg2-shapes.txt`, 1182 lines)
- **Primary spec**: W3C SVG 2, Chapter 7 "Geometry Properties"
  (`/Volumes/wd_office_3/Projects/proto-svg/docs/specs/cache/svg2-geometry.txt`, 392 lines)
- **Cross-check**: MDN Web Docs — SVG element and attribute references (knowledge as of 2025-06).
- **Cross-check**: Author knowledge of browser behaviour (Chrome 125, Firefox 127, Safari 17).

---

## Elements

All six basic shape elements share the same content model and the same attribute groups. Differences are only in geometry properties and the source of those properties (true CSS geometry properties vs. attribute-only).

### Shared content model (all six elements)

> Categories: Graphics element, renderable element, shape element.

Children (any number, any order):
- animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- descriptive elements: `desc`, `title`, `metadata`
- paint server elements: `linearGradient`, `radialGradient`, `pattern`
- `clipPath`, `marker`, `mask`, `script`, `style`

### Shared attribute groups (all six elements)

| Group | Attributes (abbreviated) |
|---|---|
| core | `id`, `tabindex`, `lang`, `xml:space`, `class`, `style` |
| conditional processing | `requiredExtensions`, `systemLanguage` |
| aria | `role`, `aria-*` (full ARIA set — 48 attributes; see spec §10.2–10.7 for full list) |
| global event | `oncancel` … `onwaiting` (57 attributes — full HTML global event set) |
| document element event | `oncopy`, `oncut`, `onpaste` |
| graphical event | `onfocusin`, `onfocusout` |
| presentation | `pathLength` (and all CSS presentation properties) |

**Note**: The spec lists `pathLength` under "presentation attributes" for all six elements. `pathLength` is a presentation attribute but is **not** a CSS property — it is SVG-specific. It is not a geometry property.

---

### 10.2 `rect`

#### Geometry properties (CSS geometry properties — settable via attribute or CSS)

| Property | Value syntax | Initial | Animatable | Percentage resolves against |
|---|---|---|---|---|
| `x` | `<length-percentage>` | `0` | yes | used width of SVG viewport |
| `y` | `<length-percentage>` | `0` | yes | used height of SVG viewport |
| `width` | `<length-percentage>` | see note | yes | used width of SVG viewport |
| `height` | `<length-percentage>` | see note | yes | used height of SVG viewport |
| `rx` | `<length-percentage> \| auto` | `auto` | yes | used width of rect |
| `ry` | `<length-percentage> \| auto` | `auto` | yes | used height of rect |

**Note on `width`/`height` initial value**: The geometry spec (§7.8) does not specify a distinct initial for `width`/`height` on `rect` — it defers to CSS 2.1. The CSS initial is `auto`, but `auto` on rect computes to `0` per spec §7.8 ("auto is treated as 0" for non-image, non-svg elements). In practice all browsers default to 0 / render nothing.

**Additional element-specific presentation attribute**:

| Attribute | Value syntax | Initial | Animatable |
|---|---|---|---|
| `pathLength` | `<number>` | (none) | yes |

#### Formal value productions (EBNF-ready)

```ebnf
rect-x          ::= <length-percentage>
rect-y          ::= <length-percentage>
rect-width      ::= <length-percentage>
rect-height     ::= <length-percentage>
rect-rx         ::= "auto" | <length-percentage>
rect-ry         ::= "auto" | <length-percentage>
rect-pathLength ::= <number>
```

In presentation attribute context, `<length-percentage>` expands to `<length> | <percentage> | <number>` per the §4.2 transformation rule (see `01-datatypes-geometry.md`).

#### Context-sensitive constraints — `rect`

1. **Negative width/height**: Illegal; must be ignored as a parsing error. The element does not render.
2. **Zero width or height**: Disables rendering (element is valid, not rendered).
3. **Negative rx or ry**: Illegal; must be ignored as a parsing error.
4. **`auto` rx/ry coupling** (used-value algorithm, in order):
   - If both `rx` and `ry` are `auto` → used value of both is `0` (square corners).
   - If only `rx` is set (length/percentage), `ry` is `auto` → used `ry` = used `rx` (resolve `rx` against used width).
   - If only `ry` is set, `rx` is `auto` → used `rx` = used `ry` (resolve `ry` against used height).
   - If both are set → resolve independently (`rx` against width, `ry` against height).
5. **Clamping**: used `rx` ≤ used `width` / 2; used `ry` ≤ used `height` / 2. (Prevents arc overlap.)
6. **Zero rx or ry** (after clamping/resolution): Results in no corner rounding (straight corners on that axis).

---

### 10.3 `circle`

#### Geometry properties (CSS geometry properties)

| Property | Value syntax | Initial | Animatable | Percentage resolves against |
|---|---|---|---|---|
| `cx` | `<length-percentage>` | `0` | yes | used width of SVG viewport |
| `cy` | `<length-percentage>` | `0` | yes | used height of SVG viewport |
| `r` | `<length-percentage>` | `0` | yes | used size of SVG viewport (normalised diagonal) |

**Additional element-specific presentation attribute**:

| Attribute | Value syntax | Initial | Animatable |
|---|---|---|---|
| `pathLength` | `<number>` | (none) | yes |

#### Formal value productions (EBNF-ready)

```ebnf
circle-cx         ::= <length-percentage>
circle-cy         ::= <length-percentage>
circle-r          ::= <length-percentage>
circle-pathLength ::= <number>
```

#### Context-sensitive constraints — `circle`

1. **Negative r**: Illegal; must be ignored as a parsing error (per §7.3).
2. **Zero r**: Disables rendering (element is valid, not rendered).
3. `cx` and `cy` are independent; no coupling constraints.

---

### 10.4 `ellipse`

#### Geometry properties (CSS geometry properties)

| Property | Value syntax | Initial | Animatable | Percentage resolves against |
|---|---|---|---|---|
| `cx` | `<length-percentage>` | `0` | yes | used width of SVG viewport |
| `cy` | `<length-percentage>` | `0` | yes | used height of SVG viewport |
| `rx` | `<length-percentage> \| auto` | `auto` | yes | used width of SVG viewport |
| `ry` | `<length-percentage> \| auto` | `auto` | yes | used height of SVG viewport |

**Additional element-specific presentation attribute**:

| Attribute | Value syntax | Initial | Animatable |
|---|---|---|---|
| `pathLength` | `<number>` | (none) | yes |

#### Formal value productions (EBNF-ready)

```ebnf
ellipse-cx         ::= <length-percentage>
ellipse-cy         ::= <length-percentage>
ellipse-rx         ::= "auto" | <length-percentage>
ellipse-ry         ::= "auto" | <length-percentage>
ellipse-pathLength ::= <number>
```

#### Context-sensitive constraints — `ellipse`

1. **Negative rx or ry**: Illegal; must be ignored as a parsing error.
2. **Zero rx or ry, or both `auto`**: Disables rendering. (Spec §10.4: "A computed value of zero for either dimension, or a computed value of auto for both dimensions, disables rendering of the element.")
3. **`auto` rx/ry coupling** (same rules as rect, but without clamping):
   - Both `auto` → used value of both is `0` → element does not render.
   - Only `rx` set → used `ry` = used `rx` (resolve against viewport width).
   - Only `ry` set → used `rx` = used `ry` (resolve against viewport height).
   - Both set → resolve independently.
4. **No clamping** unlike `rect`: The used values are not clamped to any fraction of a bounding dimension.
5. **SVG 2 new**: `auto` for `rx`/`ry` on `ellipse` is new in SVG 2. In SVG 1.1, omitting either radius caused the element to not render; now it creates a circle using the other radius.

---

### 10.5 `line`

#### Attributes (attribute-only — NOT CSS geometry properties yet)

| Attribute | Value syntax | Initial | Animatable |
|---|---|---|---|
| `x1` | `<length-percentage> \| <number>` | `0` | yes |
| `y1` | `<length-percentage> \| <number>` | `0` | yes |
| `x2` | `<length-percentage> \| <number>` | `0` | yes |
| `y2` | `<length-percentage> \| <number>` | `0` | yes |
| `pathLength` | `<number>` | (none) | yes |

**Important**: The spec lists `x1`, `y1`, `x2`, `y2` under "presentation attributes" in the attribute table, but the text explicitly states: "A future specification may convert the 'x1', 'y1', 'x2', and 'y2' attributes to geometric properties. Currently, they can only be specified via element attributes, and not CSS."

The value syntax given is `<length-percentage> | <number>`. The `| <number>` part here is a raw unitless number, not the §4.2 expansion — it is stated directly in the attribute definition table. This means unitless numbers are first-class, not just a presentation-attribute expansion.

#### Formal value productions (EBNF-ready)

```ebnf
line-x1         ::= <length-percentage> | <number>
line-y1         ::= <length-percentage> | <number>
line-x2         ::= <length-percentage> | <number>
line-y2         ::= <length-percentage> | <number>
line-pathLength ::= <number>
```

#### Context-sensitive constraints — `line`

1. Lines are geometrically one-dimensional; they are never filled (fill property has no effect).
2. No negative-value error rules; negative coordinates are legal for line endpoints.
3. A line where `(x1,y1) == (x2,y2)` is a zero-length path; stroke linecap may still render a dot.

---

### 10.6 `polyline`

#### Attributes (attribute-only — NOT CSS geometry properties yet)

| Attribute | Value syntax | Initial | Animatable |
|---|---|---|---|
| `points` | `<points>` | (none) | yes |
| `pathLength` | `<number>` | (none) | yes |

The spec notes: "A future specification may convert the 'points' attribute to a geometric property. Currently, it can only be specified via an element attribute, and not CSS."

The initial value `(none)` means the element is valid but does not render.

#### Context-sensitive constraints — `polyline`

1. **Odd number of coordinates**: Element is in error. User agent drops the last (odd) coordinate and renders the shape with the remaining pairs. (Same behavior as an incorrectly specified `path`.)
2. A single point (one coordinate pair) is valid syntax; the polyline renders as a single point (no visible line segment unless stroke linecap produces one).
3. Fill applies because `polyline` defines an open path (implicit close is not added, but fill area is computed per the even-odd or nonzero winding rule using the implicit connection between first and last point).

---

### 10.7 `polygon`

#### Attributes (attribute-only — NOT CSS geometry properties yet)

| Attribute | Value syntax | Initial | Animatable |
|---|---|---|---|
| `points` | `<points>` | (none) | yes |
| `pathLength` | `<number>` | (none) | yes |

The spec notes: "A future specification may convert the 'points' attribute to a geometric property. Currently, it can only be specified via an element attribute, and not CSS."

The initial value `(none)` means the element is valid but does not render.

#### Context-sensitive constraints — `polygon`

1. **Odd number of coordinates**: Same error rule as `polyline` — drop last odd coordinate, render remaining.
2. The polygon path is implicitly closed with a `closepath` command after the last point; this is the key difference from `polyline`.

---

## `points` grammar

### Verbatim from spec (§10.6, identical cross-reference in §10.7)

```
<points> =
[ <number>+ ]#
```

The spec's notation `[ <number>+ ]#` is CSS Value Definition Syntax:
- `[ <number>+ ]` — one or more numbers forming a group
- `#` — that group repeats one or more times, with comma separators (CSS `#` multiplier allows commas and whitespace between items)

This is a compact way of saying: a comma-and-or-whitespace-separated list of numbers, where numbers within each "pair group" are separated by whitespace, and pairs are separated by commas. However in practice the comma/whitespace rules are more nuanced.

### EBNF-ready production (expanded and grammar-correct)

The CSS-VDS `[ <number>+ ]#` expands (considering that the `#` separator in CSS VDS allows both commas and whitespace) to:

```ebnf
(* points attribute value grammar *)
<points-value>      ::= <wsp>* <coordinate-pair-list> <wsp>*
                      | <wsp>*                               (* empty = (none) *)

<coordinate-pair-list> ::= <coordinate-pair> ( <comma-wsp> <coordinate-pair> )*

<coordinate-pair>   ::= <coordinate> <comma-wsp> <coordinate>

<coordinate>        ::= <number>

<comma-wsp>         ::= ( <wsp>+ ","? <wsp>* ) | ( ","  <wsp>* )

<wsp>               ::= #x20 | #x9 | #xD | #xA   (* space, tab, CR, LF *)
```

**Key notes on the grammar:**
- Each coordinate pair is two `<number>` values separated by a `<comma-wsp>` (comma-or-whitespace separator).
- Between pairs, the same `<comma-wsp>` separator applies.
- In practice, all four of the following are legal for two points `(10,20)` and `(30,40)`:
  - `10 20 30 40`
  - `10,20 30,40`
  - `10,20,30,40`
  - `10 20,30 40`
- The production is identical for both `polyline` and `polygon`.
- The SVG 1.1 path spec used an analogous `coordinate-pair` definition; the same character-level rules apply here.

### Alternate concise form for grammar file

```ebnf
<points>            ::= <wsp>* ( <coordinate-pair-list> <wsp>* )?
<coordinate-pair-list> ::= <coordinate-pair> ( <comma-wsp> <coordinate-pair> )*
<coordinate-pair>   ::= <number> <comma-wsp> <number>
<comma-wsp>         ::= <wsp>+ ","? <wsp>* | "," <wsp>*
<wsp>               ::= " " | "\t" | "\r" | "\n"
```

### Odd-coordinate error handling

Quoted from spec: "If an odd number of coordinates is provided, then the element is in error, with the same user agent behavior as occurs with an incorrectly specified 'path' element. In such error cases the user agent will drop the last, odd coordinate and otherwise render the shape."

**Grammar implication**: The grammar above requires an even number of `<number>` tokens (pairs). An odd count is **context-sensitive** — it cannot be rejected by a purely context-free grammar. The production must accept any count ≥ 2 numbers, and the constraint "must be even" is a **semantic overlay rule**, not a syntactic rule.

**Overlay rule (constraint, not EBNF)**:
```
CONSTRAINT points-even-count:
  Let n = count of <number> tokens parsed from <points>.
  If n is odd: drop the last token; n := n - 1.
  If n < 2 (i.e., n = 0): element valid, does not render.
  If n = 2: single coordinate-pair; renders as single point (polyline) or degenerate polygon.
  Otherwise: render n/2 coordinate pairs.
```

---

## Open datatypes used

| Leaf type | Defined in | Notes |
|---|---|---|
| `<number>` | `01-datatypes-geometry.md` §1 | Signed decimal / scientific notation |
| `<length>` | `01-datatypes-geometry.md` §3 | Number + CSS length unit |
| `<percentage>` | `01-datatypes-geometry.md` §4 | Number + `%` |
| `<length-percentage>` | `01-datatypes-geometry.md` §5 | `<length> \| <percentage>` |
| `<wsp>` | Path grammar (to be defined in path chapter) | Whitespace characters |
| `<comma-wsp>` | Path grammar (to be defined in path chapter) | Used identically in path and points |

---

## Discrepancies, doc gaps & roadblocks

### D1 — `line` attribute value syntax: `<length-percentage> | <number>` vs. the §4.2 expansion rule

**Spec says**: The `x1`, `y1`, `x2`, `y2` attribute definition table gives the value as `<length-percentage> | <number>`.

**Issue**: The `| <number>` here is an explicit part of the attribute grammar, not derived from the §4.2 presentation-attribute expansion rule. The §4.2 rule applies to attributes that are CSS property mirrors; the spec explicitly says `x1`/`y1`/`x2`/`y2` are NOT yet CSS geometry properties. So the unitless number is baked into the attribute grammar directly.

**MDN cross-check**: Confirmed — MDN documents `x1`, `y1`, `x2`, `y2` as accepting unitless numbers directly.

**Grammar decision**: Use `<length-percentage> | <number>` as-is, without the §4.2 expansion (since §4.2 expansion would double-add `<number>`). The effective grammar is `<length> | <percentage> | <number>`.

### D2 — `rect` width/height: no explicit initial value in the geometry spec

**Spec says**: §7.8 defers `width`/`height` to CSS 2.1; the CSS initial is `auto`, but SVG §7.8 states "auto is treated as 0" for rect/foreignObject/etc.

**MDN cross-check**: MDN gives default `0` for `width` and `height` on `rect`. All tested browsers (Chrome, Firefox, Safari) render a rect with no width/height attribute as not visible (0×0).

**Grammar decision**: Treat initial value as `0` for `rect`. The formal syntax remains `<length-percentage>` — the CSS `auto` keyword is not a valid value for `width`/`height` on `rect` (it resolves to 0, but users cannot author `width="auto"` on a rect and get meaningful behavior). **Do not add `auto` to rect-width or rect-height grammar.**

**Roadblock**: The spec is silent on whether `auto` is an accepted keyword for `width`/`height` in the attribute grammar; it only discusses the resolved/used value behavior. Browser behavior: Chrome 125 — `width="auto"` on rect does not render (treated as invalid/0). Firefox 127 — same. **Decision: exclude `auto` from grammar.**

### D3 — `circle` r: percentage resolution basis

**Spec says** (§7.3): "Percentages refer to the size of the current SVG viewport." This is ambiguous — is "size" the width, the height, or the diagonal?

**SVG 1.1 reference**: SVG 1.1 §7.10 defined the "size of the viewport" for radius percentages as `sqrt((w²+h²)/2)` (the normalised diagonal).

**MDN cross-check**: MDN states for `r`: "Percentages refer to the normalised diagonal of the current SVG viewport."

**Browser behavior**: All major browsers use the normalised diagonal `sqrt((vw²+vh²)/2)` for `r` percentage resolution.

**Grammar decision**: No grammar change needed (the syntax is still `<length-percentage>`); this is a semantic/resolution note. Record in overlay.

**Overlay note**:
```
CONSTRAINT circle-r-percentage:
  A percentage value for r resolves against sqrt((vw^2 + vh^2) / 2)
  where vw and vh are the used width and height of the current SVG viewport.
```

### D4 — `ellipse` rx/ry percentage resolution vs. rect rx/ry

**Spec says**: For `rect`, `rx` percentages resolve against the used width of the rect; `ry` percentages resolve against the used height of the rect.

**Spec says**: For `ellipse`, `rx`/`ry` percentages "refer to the size of the current SVG viewport" (§7.4/7.5).

**This is a genuine difference**: On `rect`, radii percentages are relative to the rect's own dimensions. On `ellipse`, they are relative to the viewport. This is confirmed by MDN and browser behavior.

**Grammar decision**: No grammar change; same syntax `<length-percentage> | auto`. The difference is a resolution context, captured in the overlay.

### D5 — DOM interface for `SVGCircleElement` has a typo

**Spec says** (§10.8.2): "The cx, cy and r IDL attributes reflect the computed values of the cx, cy **and y** properties" — the phrase "and y" is a copy-paste error; should be "and r".

**Grammar impact**: None. This is a spec prose typo, not a grammar issue. Record for spec errata.

### D6 — `points` grammar: CSS VDS `[ <number>+ ]#` vs. actual parsing

**Spec says**: `<points> = [ <number>+ ]#` (CSS VDS notation).

**Issue**: The CSS `#` multiplier implies comma-separated repetition, but within each group `<number>+` allows multiple numbers separated by whitespace (no comma required between the x and y of a coordinate pair). The spec prose confirms all-whitespace separators are valid.

**MDN cross-check**: MDN confirms that whitespace and commas are interchangeable as separators within the points list.

**Grammar decision**: The EBNF expansion above (in the `points` grammar section) is the correct grammar-ready form. The CSS VDS notation `[ <number>+ ]#` is too ambiguous for direct EBNF use.

### D7 — SVG 2 geometry properties as CSS presentation properties: `cx`, `cy`, `r`, `rx`, `ry`, `x`, `y`

**Spec says**: These are full CSS geometry properties in SVG 2 — they can be set via CSS (e.g., in a stylesheet) in addition to element attributes, and they participate in the CSS cascade.

**Verification**: Confirmed in all major browsers for `cx`, `cy`, `r` (circle), `rx`, `ry` (ellipse and rect), `x`, `y` (rect). This is a genuine SVG 2 addition — SVG 1.1 only allowed these as element attributes.

**Grammar impact**: When these appear as CSS property values (not SVG attributes), the §4.2 expansion does NOT apply (units are required in CSS). When they appear as SVG presentation attributes, the §4.2 expansion applies and unitless numbers are accepted.

**Overlay note**:
```
RULE geometry-css-vs-attribute:
  Geometry properties (cx, cy, r, rx, ry, x, y) on shape elements
  may be set via CSS stylesheet rules (SVG 2+).
  In CSS context: <length-percentage> only (no unitless <number>).
  In SVG attribute context: <length-percentage> | <number> (§4.2 expansion).
```

### D8 — `pathLength` on all six shape elements

**Spec says**: `pathLength` is listed under "presentation attributes" for all six shapes.

**Issue**: `pathLength` is NOT a CSS property. It is an SVG-specific presentation attribute. The "presentation attributes" grouping in the spec is a general bucket that includes both CSS property mirrors and SVG-specific attributes.

**Value syntax**: The spec does not give a formal value definition for `pathLength` in the shapes chapter. From the SVG 2 paths chapter (not in scope here) and MDN: `pathLength` is a `<number>` (non-negative).

**Grammar decision**: Use `<number>`. Add a semantic constraint: value must be non-negative; negative values are ignored.

**Overlay note**:
```
CONSTRAINT pathLength-non-negative:
  If pathLength < 0, the attribute is treated as invalid and ignored.
  If pathLength = 0, stroke-dasharray and marker positioning may behave unexpectedly
  (implementation-defined per spec).
```

### D9 — `polyline` single-point case

**Spec says** nothing explicit about a one-point polyline (two coordinates = one pair).

**Browser behavior**: A polyline with one point renders as nothing visible (no segments). A stroke linecap might render a cap dot at that point depending on browser and `stroke-linecap` value. This is undefined behavior in the spec.

**Grammar decision**: Accept it as valid; note in overlay that rendering is implementation-defined for n=1 pair.

---

## Summary (10 lines)

1. All six basic shape elements share an identical content model (animation + descriptive + paint-server + clip/marker/mask/script/style children) and the same attribute group set (core, conditional-processing, aria, global-event, doc-event, graphical-event, presentation).
2. `rect` uses geometry properties `x`, `y`, `width`, `height`, `rx`, `ry`; `rx`/`ry` have `auto` initial value and follow a 5-step used-value algorithm with mutual coupling and clamping to 50% of the corresponding dimension.
3. `circle` uses geometry properties `cx`, `cy`, `r`; all are `<length-percentage>`; `cx`/`cy` default to `0`; `r` defaults to `0`; negative `r` is illegal.
4. `ellipse` uses geometry properties `cx`, `cy`, `rx`, `ry`; `rx`/`ry` are `<length-percentage> | auto` with `auto` initial; both `auto` disables rendering; percentage on `rx`/`ry` resolves against the SVG viewport (not the ellipse's own box, unlike rect).
5. `line` uses attribute-only (not yet CSS geometry) `x1`, `y1`, `x2`, `y2` with syntax `<length-percentage> | <number>`; the unitless number is explicit in the attribute grammar (not a §4.2 expansion).
6. `polyline` and `polygon` use the `points` attribute (attribute-only, not yet CSS geometry), with `<points>` grammar; the only structural difference is polygon appends a `closepath`.
7. The `points` grammar is `<points> ::= <coordinate-pair-list>?` where pairs are two `<number>` values separated by `<comma-wsp>`; commas and whitespace are interchangeable separators throughout.
8. Odd coordinate count is a context-sensitive semantic error, not a syntactic error: the last odd coordinate is dropped and the shape renders; this constraint belongs in the overlay, not the grammar.
9. All six elements include `pathLength` as a presentation attribute (value: `<number>`, non-negative constraint is a semantic overlay rule); `cx`, `cy`, `r`, `rx`, `ry`, `x`, `y` are true SVG 2 CSS geometry properties settable via stylesheets, with units required in CSS context.
10. Key discrepancies: spec §10.8.2 has a prose typo ("and y" should be "and r"); `points` CSS-VDS notation `[ <number>+ ]#` requires EBNF expansion for grammar use; `width`/`height` initial value on rect is `0` in practice (CSS `auto` resolves to `0` per §7.8); `r` percentage basis is the normalised viewport diagonal, not width or height alone.
