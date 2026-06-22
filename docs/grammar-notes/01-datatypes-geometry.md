# Datatypes & geometry grammar notes

## Source

- **Primary spec**: W3C SVG 2, Chapter 4 "Basic Data Types and Interfaces"
  (`/Volumes/wd_office_3/Projects/proto-svg/docs/specs/cache/svg2-types.txt`, 2915 lines)
- **Primary spec**: W3C SVG 2, Chapter 7 "Geometry Properties"
  (`/Volumes/wd_office_3/Projects/proto-svg/docs/specs/cache/svg2-geometry.txt`, 391 lines)
- **Cross-check**: MDN SVG Attribute Reference
  (`/Volumes/wd_office_3/Projects/proto-svg/docs/specs/mdn_docs_attributes.md`)
- **Cross-check**: Author knowledge of CSS Values & Units Level 4, CSS 2.1, and browser behaviour as of 2025.

---

## Attribute value syntax conventions (Chapter 4.2)

SVG 2 defines six methods for specifying attribute value syntax. The default is **CSS Value Definition Syntax** (CSS-VDS). The others are: EBNF, ABNF, URL Standard, HTML Standard, and prose.

### Presentation attribute special rule

When a **presentation attribute** (one that mirrors a CSS property) is parsed, the grammar is first transformed before matching:

```
Replace all <length>           with [ <length> | <number> ]
Replace all <length-percentage> with [ <length-percentage> | <number> ]
Replace all <angle>            with [ <angle> | <number> ]
```

This allows **unitless numbers** as lengths/angles in presentation attributes (treated as user units / degrees), while the corresponding CSS property continues to require units.

All other CSS-VDS attributes are parsed as written (no expansion).

### Invalid value fallback

When any attribute fails to parse, the attribute is treated as if the **initial value** were specified. For presentation attributes the initial value is the corresponding CSS property's initial value. This means an invalid presentation attribute can override lower-priority stylesheet rules (it injects a computed-value rule at author specificity).

---

## Basic data types

The SVG 2 types chapter is primarily a DOM interface specification. It does **not** define standalone EBNF productions for `<number>`, `<integer>`, `<length>`, `<angle>`, etc. — those are inherited from CSS Values & Units. The chapter does expose which unit classes are supported through the DOM constants. Productions below are reconstructed from:
- The DOM constant tables in Chapter 4.5
- Cross-reference with CSS Values & Units Level 4 [css-values-4]
- Author knowledge of SVG 1.1 / SVG 2 deployed practice

### 1. `<number>`

```
<number> ::= [ "+" | "-" ]? ( <integer> | <decimal> | <scientific> )

<integer>    ::= [0-9]+
<decimal>    ::= [0-9]* "." [0-9]+
<scientific> ::= ( <integer> | <decimal> ) ( "e" | "E" ) [ "+" | "-" ]? [0-9]+
```

Examples: `0`, `42`, `-3.14`, `1.5e-3`, `.5`, `+100`

**Precision note (§4.2.1):** Must support at minimum all finite IEEE 754 single-precision values. Higher precision recommended for coordinate transforms. Implementation serialises to an unspecified string that round-trips as the closest supported value.

### 2. `<integer>`

```
<integer> ::= [ "+" | "-" ]? [0-9]+
```

Examples: `0`, `3`, `-1`, `+256`

Carried by `SVGAnimatedInteger` (used in filter spec, e.g. `numOctaves`, `order`).

### 3. `<length>`

From CSS Values & Units Level 4. SVG adds the unitless form in presentation attributes via the §4.2 transformation rule.

```
<length> ::= <number> <length-unit>

<length-unit> ::= "em" | "ex" | "px" | "in" | "cm" | "mm" | "pt" | "pc"
                | "ch" | "rem" | "vw" | "vh" | "vmin" | "vmax"
                | "cap" | "ic" | "lh" | "rlh" | "vi" | "vb"
                | "svw" | "svh" | "lvw" | "lvh" | "dvw" | "dvh"   (* CSS 4 additions *)
                | "cqw" | "cqh" | "cqi" | "cqb" | "cqmin" | "cqmax" (* container query units *)
                | "Q"
```

**SVGLength DOM constants recognise only**: `em`, `ex`, `px`, `cm`, `mm`, `in`, `pt`, `pc` (plus unitless `<number>` and `%`). All other CSS units parse correctly but return `SVG_LENGTHTYPE_UNKNOWN` from the DOM.

Examples: `10px`, `2.5em`, `50%` (as `<percentage>`), `0` (unitless, valid in presentation attrs)

### 4. `<percentage>`

```
<percentage> ::= <number> "%"
```

Examples: `50%`, `100%`, `-10%`, `0%`

### 5. `<length-percentage>`

This is the composite type used by most geometry properties in SVG 2.

```
<length-percentage> ::= <length> | <percentage>
```

In presentation attributes the rule expands to:

```
<length-percentage-pa> ::= <length> | <percentage> | <number>
```

### 6. `<angle>`

```
<angle> ::= <number> <angle-unit>

<angle-unit> ::= "deg" | "rad" | "grad" | "turn"
```

**SVGAngle DOM constants recognise only**: `deg`, `rad`, `grad`. `turn` parses but returns `SVG_ANGLETYPE_UNKNOWN` from DOM. In presentation attributes `<number>` is accepted and treated as degrees (SVG_ANGLETYPE_UNSPECIFIED).

Examples: `90deg`, `1.5708rad`, `100grad`, `0.25turn`, `45` (presentation attr only)

### 7. `<coordinate>`

Not a distinct production in SVG 2; effectively a synonym for `<length-percentage>` (or `<length>` in SVG 1.1). Used in the prose of SVG 1.1 Chapter 4 for x/y coordinate attributes. In SVG 2 the geometry properties chapter uses `<length-percentage>` directly.

**Grammar decision**: Alias `<coordinate>` = `<length-percentage>` (= `<length> | <percentage>`). In presentation attributes expand as `<length> | <percentage> | <number>`.

### 8. `<opacity-value>`

```
<opacity-value> ::= <number> | <percentage>
```

Range constraint (not syntax): clamped to [0, 1] for `<number>`, [0%, 100%] for `<percentage>`. Out-of-range values are clamped at render time per §4.2.2.

Examples: `1`, `0`, `0.5`, `50%`

Used by: `opacity`, `fill-opacity`, `stroke-opacity`, `flood-opacity`, `stop-opacity`.

### 9. `<url>` / `<IRI>` / `<FuncIRI>`

SVG 2 moves to the URL Standard [URL] for `href`. Legacy SVG 1.1 used IRI/FuncIRI.

```
(* Legacy SVG 1.1 forms, still in wide use *)
<IRI>     ::= (* any IRI as defined by RFC 3987 *)      (* [ABNF] *)
<FuncIRI> ::= "url(" <IRI> ")"
            | "url(" <ws>* <quoted-IRI> <ws>* ")"       (* CSS url() form *)

(* SVG 2 href attribute — parsed by URL Standard *)
<url-attribute> ::= (* URL string per [URL] *)           (* [URL] *)

(* CSS url() function form used in paint/filter properties *)
<url()> ::= "url(" <ws>* [ <string> | <unquoted-url> ] <ws>* ")"
```

The `href` IDL attribute on `SVGURIReference` additionally falls back to the deprecated `xlink:href` attribute if `href` is absent. `xlink:href` is deprecated in SVG 2.

### 10. `<paint>`

Not formally defined in the types chapter; defined in the SVG Painting chapter. Included here as a reference type.

```
<paint> ::= "none"
          | "currentColor"
          | <color>
          | <url()> [ "none" | "currentColor" | <color> ]?
          | "context-fill"
          | "context-stroke"
```

Used by: `fill`, `stroke`, `flood-color`, `lighting-color`, `stop-color` (color-only variants).

### 11. `<color>`

Per CSS Color Level 4. Full grammar not reproduced here (it is a large specification). Key forms:

```
<color> ::= <named-color>          (* e.g. "red", "transparent" *)
          | <hex-color>            (* #RGB | #RGBA | #RRGGBB | #RRGGBBAA *)
          | <rgb()>                (* rgb(R, G, B) or rgb(R G B) *)
          | <rgba()>               (* rgba(R, G, B, A) *)
          | <hsl()>                (* hsl(H, S%, L%) *)
          | <hsla()>               (* hsla(H, S%, L%, A) *)
          | <hwb()>
          | <lab()> | <lch()> | <oklab()> | <oklch()>
          | <color()>              (* color(space components) *)
          | "currentColor"
```

In SVG contexts, `currentColor` resolves to the computed value of the element's `color` property.

### 12. `<transform-list>`

Defined in SVG 2 by reference to the CSS Transforms spec. The SVGTransformList interface tracks this.

```
<transform-list> ::= "none" | <transform-function>+

<transform-function> ::=
    "matrix("  <number> <sep> <number> <sep> <number> <sep>
               <number> <sep> <number> <sep> <number> ")"
  | "translate(" <length-percentage> [ <sep> <length-percentage> ]? ")"
  | "translateX(" <length-percentage> ")"
  | "translateY(" <length-percentage> ")"
  | "scale("    <number> [ <sep> <number> ]? ")"
  | "scaleX("   <number> ")"
  | "scaleY("   <number> ")"
  | "rotate("   <angle> [ <sep> <length-percentage> <sep> <length-percentage> ]? ")"
  | "skewX("   <angle> ")"
  | "skewY("   <angle> ")"

<sep> ::= <ws>* [ "," ]? <ws>*
<ws>  ::= " " | "\t" | "\n" | "\r" | "\f"
```

**Presentation attribute note**: In the `transform` presentation attribute, `<angle>` expands to `<angle> | <number>` and `<length-percentage>` expands to `<length-percentage> | <number>`.

**SVG 1.1 legacy**: The old `transform` attribute grammar (EBNF from SVG 1.1) allowed `rotate(angle cx cy)` where cx/cy were unitless user-unit numbers. That form remains valid because presentation-attribute expansion covers it.

### 13. `<number-optional-number>`

```
<number-optional-number> ::= <number> [ <sep> <number> ]?
```

Used by: `stdDeviation`, `kernelUnitLength`, `baseFrequency`. When a single number is provided the second is implied equal to the first (attribute-specific rule).

### 14. `<list-of-numbers>`

```
<list-of-numbers> ::= <number> ( <sep> <number> )*
```

Reflected by SVGNumberList / SVGAnimatedNumberList.

Examples: `1 2 3 4`, `0,1,0,0`, `1 0.5`

### 15. `<list-of-lengths>`

```
<list-of-lengths> ::= <length-pa> ( <sep> <length-pa> )*

<length-pa> ::= <length> | <percentage> | <number>  (* number only in presentation attrs *)
```

Reflected by SVGLengthList / SVGAnimatedLengthList.

Examples: `10px 20px`, `5%, 10%`, `1em 2em 3em`

### 16. `<dasharray>`

```
<dasharray> ::= "none" | <dash-list>

<dash-list> ::= <dash-value> ( <sep> <dash-value> )*

<dash-value> ::= <length-percentage>    (* non-negative *)
               | <number>               (* presentation attrs only; non-negative *)
```

Range constraint (not syntax): each `<dash-value>` must be non-negative.

Used by: `stroke-dasharray`.

Examples: `none`, `5 3`, `10px 5px 2px`, `5%`, `0`

### 17. `<viewBox>`

```
<viewBox> ::= <min-x> <sep> <min-y> <sep> <width> <sep> <height>

<min-x>  ::= <number>
<min-y>  ::= <number>
<width>  ::= <number>    (* must be non-negative *)
<height> ::= <number>    (* must be non-negative *)
```

Reflected as SVGAnimatedRect. Negative width or height is an error.

Examples: `0 0 100 100`, `0, 0, 500, 300`

### 18. `<preserveAspectRatio>`

```
<preserveAspectRatio> ::= "none"
                        | <align> [ <ws>+ <meetOrSlice> ]?

<align> ::= "xMinYMin" | "xMidYMin" | "xMaxYMin"
          | "xMinYMid" | "xMidYMid" | "xMaxYMid"
          | "xMinYMax" | "xMidYMax" | "xMaxYMax"

<meetOrSlice> ::= "meet" | "slice"
```

Default: `xMidYMid meet`. Reflected as SVGAnimatedPreserveAspectRatio.

### 19. `<zoomAndPan>`

```
<zoomAndPan> ::= "disable" | "magnify"
```

Deprecated in SVG 2; MDN lists it as deprecated. DOM: SVGZoomAndPan constants 1 (disable) and 2 (magnify).

### 20. `<boolean>` (attribute form)

```
<boolean> ::= "true" | "false"
```

Used by e.g. `preserveAlpha`.

### 21. `<unitTypes>` (gradientUnits / patternUnits etc.)

```
<unit-type> ::= "userSpaceOnUse" | "objectBoundingBox"
```

DOM: SVGUnitTypes constants 1 and 2.

### 22. `<string-list>` (requiredExtensions, systemLanguage)

```
<string-list> ::= <string-token> ( <ws>+ <string-token> )*
```

Reflected by SVGStringList.

---

## Geometry properties

Chapter 7 defines 8 properties (cx, cy, r, rx, ry, x, y, width, height are mentioned; width and height are deferred to CSS 2.1). The spec section for width/height is brief: it defers to CSS 2.1 and adds SVG-specific `auto` handling.

### cx — Horizontal center coordinate

| Field | Value |
|---|---|
| Applies to | `circle`, `ellipse` |
| Value syntax | `<length-percentage>` |
| Initial | `0` |
| Inherited | no |
| Percentages | refer to current SVG viewport size (see Units) |
| Computed value | absolute length or percentage |
| Animatable | yes |

Grammar rule:
```
cx ::= <length-percentage>
     (* in presentation attr: <length> | <percentage> | <number> *)
```

Cross-check: MDN lists `cx` as a standard attribute (no deprecation marker). Used on `circle` and `ellipse`. MDN also lists `cx` on `radialGradient` and `fePointLight` / `feMergeNode` but those are different attribute contexts not covered by the geometry property.

### cy — Vertical center coordinate

| Field | Value |
|---|---|
| Applies to | `circle`, `ellipse` |
| Value syntax | `<length-percentage>` |
| Initial | `0` |
| Inherited | no |
| Percentages | refer to current SVG viewport size |
| Computed value | absolute length or percentage |
| Animatable | yes |

Grammar rule:
```
cy ::= <length-percentage>
```

### r — Radius

| Field | Value |
|---|---|
| Applies to | `circle` |
| Value syntax | `<length-percentage>` |
| Initial | `0` |
| Inherited | no |
| Percentages | refer to current SVG viewport size |
| Computed value | absolute length or percentage |
| Animatable | yes |

Grammar rule:
```
r ::= <length-percentage>
```

Constraint (not syntax): negative value is **illegal**; must be treated as an error (no rendering).

Cross-check: MDN also uses `r` on `radialGradient` but that is a different attribute (unitless number).

### rx — Horizontal radius

| Field | Value |
|---|---|
| Applies to | `ellipse`, `rect` |
| Value syntax | `<length-percentage> \| auto` |
| Initial | `auto` |
| Inherited | no |
| Percentages | refer to current SVG viewport size |
| Computed value | absolute length, percentage, or `auto` |
| Animatable | yes |

Grammar rule:
```
rx ::= "auto" | <length-percentage>
```

Context-sensitive rules (overlay, not syntax):
- If `rx` is `auto`: used value equals the computed `ry` absolute length (circular arc). If both are `auto`, used value is `0`.
- For `rect`: used value of `rx` is clamped to at most `50%` of the computed `width`.
- Negative value is **illegal**.
- `auto` on `ellipse` is **new in SVG 2** (was not specified in SVG 1.1).

### ry — Vertical radius

| Field | Value |
|---|---|
| Applies to | `ellipse`, `rect` |
| Value syntax | `<length-percentage> \| auto` |
| Initial | `auto` |
| Inherited | no |
| Percentages | refer to current SVG viewport size |
| Computed value | absolute length, percentage, or `auto` |
| Animatable | yes |

Grammar rule:
```
ry ::= "auto" | <length-percentage>
```

Context-sensitive rules (overlay):
- If `ry` is `auto`: used value equals the computed `rx` absolute length. If both are `auto`, used value is `0`.
- For `rect`: used value of `ry` is clamped to at most `50%` of the computed `height`.
- Negative value is **illegal**.

### x — Horizontal coordinate

| Field | Value |
|---|---|
| Applies to | `svg`, `rect`, `image`, `foreignObject` |
| Value syntax | `<length-percentage>` |
| Initial | `0` |
| Inherited | no |
| Percentages | refer to current SVG viewport size |
| Computed value | absolute length or percentage |
| Animatable | yes |

Grammar rule:
```
x ::= <length-percentage>
```

Cross-check: MDN shows `x` also used on `filter`, `mask`, `pattern`, `symbol`, `use`, `feBlend`, and many filter primitives — those are attribute-level uses, not the CSS geometry property. The CSS property `x` in SVG 2 applies only to the four elements listed.

### y — Vertical coordinate

| Field | Value |
|---|---|
| Applies to | `svg`, `rect`, `image`, `foreignObject` |
| Value syntax | `<length-percentage>` |
| Initial | `0` |
| Inherited | no |
| Percentages | refer to current SVG viewport size |
| Computed value | absolute length or percentage |
| Animatable | yes |

Grammar rule:
```
y ::= <length-percentage>
```

### width — Sizing (horizontal)

| Field | Value |
|---|---|
| Applies to | `svg`, `rect`, `image`, `foreignObject` |
| Value syntax | (see CSS 2.1 + SVG 2 §7.8) |
| Initial | CSS 2.1 default (`auto`) |
| Inherited | no |
| Animatable | yes |

Grammar rule (SVG-specific `auto` semantics plus standard CSS width):
```
width ::= "auto" | <length> | <percentage> | "min-content" | "max-content" | "fit-content(" <length-percentage> ")"
        (* ...full CSS sizing value grammar *)
```

SVG 2 §7.8 `auto` behaviour (context-sensitive overlay):
- On `svg`: `auto` is treated as `100%`.
- On `image`: `auto` computed from the intrinsic dimensions and aspect ratio per CSS default sizing algorithm.
- On all other SVG elements (`rect`, `foreignObject`): `auto` is treated as `0`.

The used value of `width` may be further constrained by `max-width` and `min-width`.

### height — Sizing (vertical)

| Field | Value |
|---|---|
| Applies to | `svg`, `rect`, `image`, `foreignObject` |
| Value syntax | (same as width, see CSS 2.1 + SVG 2 §7.8) |
| Initial | CSS 2.1 default (`auto`) |
| Inherited | no |
| Animatable | yes |

Grammar rule:
```
height ::= "auto" | <length> | <percentage> | (* ...full CSS sizing grammar *)
```

SVG 2 `auto` behaviour: same pattern as `width` (100% for svg, intrinsic for image, 0 for others).

---

## Open datatypes — leaves to scalarize

These types are genuinely open (not enumerable). The grammar should treat them as terminal leaves and document representative samples for test generation.

| Leaf | Definition source | Representative samples |
|---|---|---|
| `<number>` | CSS / reconstructed | `0`, `1`, `-1`, `3.14`, `1e3`, `.5`, `-0.001` |
| `<integer>` | CSS | `0`, `1`, `-1`, `100`, `32767` |
| `<length>` | CSS Values L4 | `0`, `1px`, `10px`, `1em`, `2.5rem`, `50vw`, `1cm`, `1in`, `12pt`, `1pc` |
| `<percentage>` | CSS | `0%`, `50%`, `100%`, `-10%` |
| `<angle>` | CSS | `0deg`, `90deg`, `180deg`, `360deg`, `3.14rad`, `100grad`, `0.25turn` |
| `<color>` | CSS Color L4 | `red`, `#fff`, `#ff0000`, `rgb(255,0,0)`, `rgba(0,0,255,0.5)`, `hsl(0,100%,50%)`, `currentColor`, `transparent` |
| `<url()>` | CSS | `url(#id)`, `url("http://example.com/img.svg")` |
| `<IRI>` | RFC 3987 | `#circle1`, `https://example.com/defs.svg#marker` |
| `<named-color>` | CSS Color L4 (148 keywords) | `red`, `green`, `blue`, `black`, `white`, `transparent` |
| `<string-token>` | CSS / XML | `en`, `fr-CA`, `http://www.w3.org/1999/xhtml` |

---

## Discrepancies, doc gaps & roadblocks

### D1. svg2-types.txt has no EBNF productions for primitive types

**Issue**: Chapter 4 of the SVG 2 spec (as captured in svg2-types.txt) is almost entirely a DOM interface specification. It does not contain a formal EBNF or CSS-VDS definition of `<number>`, `<integer>`, `<length>`, `<percentage>`, `<angle>`, `<color>`, `<paint>`, `<dasharray>`, `<transform-list>`, or `<IRI>`. These are all delegated to external specs (CSS Values, CSS Color, CSS Transforms, RFC 3987, URL Standard).

**Decision for grammar**: Reconstruct productions from those external specs and author knowledge. Mark each reconstructed production with a source tag so future maintainers know the normative reference.

### D2. SVG 1.1 `<length>` unitless form vs. SVG 2 presentation attribute expansion

**Issue**: SVG 1.1 defined `<length>` as allowing a unitless number (treated as user units). SVG 2 formally removes this from `<length>` itself and instead handles it through the presentation attribute expansion rule (§4.2). This is a subtle but important distinction for grammar: `<length>` in CSS property context never accepts bare numbers; `<length>` in presentation attribute context does (via the expansion).

**Decision**: The grammar must have two forms:
- `<length>` — strict CSS form, units required (except `0`)
- `<length-pa>` or `<length-or-number>` — expanded form for presentation attributes

### D3. `r` property: spec says `<length-percentage>` but MDN and SVG 1.1 used unitless number

**Issue**: The geometry spec defines `r` as `<length-percentage>`. In SVG 1.1 `r` was a `<length>` (which allowed unitless). In SVG 2 the presentation attribute expansion rule covers the unitless number case for `r`. However, `r` does not accept `auto` (unlike `rx`/`ry`). MDN's attribute page for `r` confirms: accepts `<length>`, `<percentage>`, or unitless number (in presentation attr form). Consistent with spec.

**Decision**: Grammar for `r` = `<length-percentage>` with presentation-attribute expansion covering the unitless case. No discrepancy requiring a special case.

### D4. `width` and `height` definitions are deferred to CSS 2.1 by the geometry spec

**Issue**: SVG 2 §7.8 says "See the CSS 2.1 specification for the definitions of width and height" and adds only the `auto` special-casing. The CSS 2.1 `width`/`height` grammar does not include `min-content`, `max-content`, or `fit-content()`. Modern browsers accept CSS Sizing Level 3 values. The spec text is behind current practice.

**Decision for grammar**: Use the full CSS Sizing Level 3 grammar for `width`/`height`, noting that the SVG spec deferral is underspecified. Record SVG-specific `auto` semantics in the overlay (not syntax).

### D5. `cx` and `cy` applied to `radialGradient` (not a geometry property)

**Issue**: `cx` and `cy` as geometry CSS properties apply only to `circle` and `ellipse`. However `cx` and `cy` are also presentation attributes on `radialGradient` with a different type (`<length-percentage>` with different percentage resolution rules — relative to the gradient bounding box, not the viewport) and on `fePointLight` (plain `<number>`). MDN lists the attribute for multiple elements.

**Decision**: The geometry property grammar entries are correct as per the SVG 2 Chapter 7 spec. The other element uses are separate attribute definitions (in the painting and filter chapters) and should be documented there, not here.

### D6. `<dasharray>` not formally defined in svg2-types.txt

**Issue**: The types chapter does not define `<dasharray>`. It is defined in the SVG Painting chapter (not captured in the current source files). Production above is reconstructed from SVG 1.1 §11.4, CSS Stroke Module Level 1, and MDN.

**Decision**: Mark as "reconstructed; needs verification against svg2-painting.txt when available."

### D7. `<transform-list>` SVG 1.1 vs. CSS Transforms

**Issue**: SVG 1.1 had its own EBNF for `transform`. SVG 2 delegates to CSS Transforms Level 1. CSS Transforms adds `matrix3d()`, `rotate3d()`, `scale3d()`, `translate3d()`, `perspective()`. SVG 2 graphics elements are 2D; 3D transform functions should be accepted syntactically but flatten to 2D.

**Decision**: Grammar accepts full CSS Transforms vocabulary. Constraint overlay notes that 3D functions are accepted by CSS Transforms but SVG rendering flattens them.

### D8. `<angle>` unit `turn` supported in browsers but not in SVGAngle DOM

**Issue**: The `turn` unit is valid CSS and supported in modern browsers for SVG attributes (e.g. `rotate="0.25turn"` works in Chrome/Firefox). However `SVGAngle.unitType` returns `SVG_ANGLETYPE_UNKNOWN` for `turn`, and the spec only enumerates `deg`, `rad`, `grad` in the DOM constants.

**Decision**: Grammar includes `"turn"` as a valid `<angle-unit>`. The DOM limitation is a DOM-layer concern only, not a parsing concern. Noted in the overlay.

### D9. `<opacity-value>` percentage form: SVG 1.1 vs. CSS

**Issue**: SVG 1.1 defined opacity values as `<number>` in [0, 1] only. CSS Opacity Level 1 added `<percentage>`. SVG 2 presentation attributes for `opacity`, `fill-opacity`, `stroke-opacity` etc. now accept both. MDN confirms both forms are valid.

**Decision**: Grammar includes `<opacity-value> ::= <number> | <percentage>`. Clamping to [0,1] / [0%,100%] is a constraint overlay item.

### D10. `zoomAndPan` — deprecated but still in DOM

**Issue**: MDN marks `zoomAndPan` as deprecated. SVGZoomAndPan mixin still present in SVG 2 IDL. The DOM constants remain. Grammar should include it but flag as deprecated in the overlay.

---

## Summary of counts

- **Basic data types documented**: 22 (number, integer, length, percentage, length-percentage, coordinate, angle, opacity-value, url/IRI/FuncIRI, paint, color, transform-list, number-optional-number, list-of-numbers, list-of-lengths, dasharray, viewBox, preserveAspectRatio, zoomAndPan, boolean, unit-type, string-list)
- **Geometry properties documented**: 9 (cx, cy, r, rx, ry, x, y, width, height)
- **Discrepancies / gaps documented**: 10 (D1–D10)
