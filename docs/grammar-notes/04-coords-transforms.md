# Coords & Transforms Grammar Notes

## Source

Primary: `/docs/specs/cache/svg2-coords.txt` — W3C SVG 2, Chapter 8
(Coordinate Systems, Transformations and Units).

Cross-references consulted:
- CSS Transforms Module Level 1 — §6.2 SVG transform attribute
  <https://drafts.csswg.org/css-transforms-1/#svg-transform>
- CSS Values and Units Module Level 4 — §6 Lengths, §8 Angles
  <https://drafts.csswg.org/css-values-4/>
- SVG 1.1 Second Edition — §7.8 preserveAspectRatio
  <https://www.w3.org/TR/SVG11/coords.html#PreserveAspectRatioAttribute>
- MDN: SVG transform attribute, CSS transform-origin

---

## transform-list grammar (verbatim BNF from CSS Transforms §6.2)

> SVG2 §8.5 says only: "User agents must support the transform property and
> presentation attribute as defined in [css-transforms-1]."  No BNF is given
> inline.  The authoritative SVG-attribute grammar lives in CSS Transforms §6.2.
> The production below is transcribed verbatim from that normative section.

```
transform-list  = wsp* transforms? wsp*

transforms      = transform
               | transform comma-wsp? transforms

transform       = matrix
               | translate
               | scale
               | rotate
               | skewX
               | skewY

matrix          = "matrix" wsp* "(" wsp*
                    number comma-wsp?
                    number comma-wsp?
                    number comma-wsp?
                    number comma-wsp?
                    number comma-wsp?
                    number wsp* ")"

translate       = "translate" wsp* "(" wsp*
                    number ( comma-wsp? number )? wsp* ")"

scale           = "scale" wsp* "(" wsp*
                    number ( comma-wsp? number )? wsp* ")"

rotate          = "rotate" wsp* "(" wsp*
                    number ( comma-wsp? number comma-wsp? number )? wsp* ")"

skewX           = "skewX" wsp* "(" wsp* number wsp* ")"

skewY           = "skewY" wsp* "(" wsp* number wsp* ")"

comma-wsp       = ( wsp+ comma? wsp* ) | ( comma wsp* )

comma           = U+002C COMMA ","

wsp             = U+000A LINE FEED
               | U+000D CARRIAGE RETURN
               | U+0009 CHARACTER TABULATION
               | U+0020 SPACE

number          = <CSS <number> token> — sign, integer part, optional
                  decimal fraction; no units
```

### Argument semantics (units implied, not written)

| Function | Numeric args | Implied unit |
|---|---|---|
| `matrix(a,b,c,d,e,f)` | 6 numbers — a,b,c,d dimensionless; e,f are px offsets | dimensionless / px |
| `translate(tx[,ty])` | tx mandatory; ty optional, defaults to 0 | px (implied) |
| `scale(sx[,sy])` | sy optional, defaults to sx | dimensionless |
| `rotate(angle[,cx,cy])` | angle in degrees (implied); cx,cy optional centre, px | deg / px |
| `skewX(angle)` | angle in degrees (implied) | deg |
| `skewY(angle)` | angle in degrees (implied) | deg |

### EBNF-ready rendering

```ebnf
transform-list  ::= WSP* ( transforms WSP* )?
transforms      ::= transform ( COMMA-WSP? transform )*
transform       ::= matrix | translate | scale | rotate | skewX | skewY
matrix          ::= "matrix" WSP* "(" WSP*
                    number COMMA-WSP? number COMMA-WSP? number COMMA-WSP?
                    number COMMA-WSP? number COMMA-WSP? number WSP* ")"
translate       ::= "translate" WSP* "(" WSP* number ( COMMA-WSP? number )? WSP* ")"
scale           ::= "scale"     WSP* "(" WSP* number ( COMMA-WSP? number )? WSP* ")"
rotate          ::= "rotate"    WSP* "(" WSP* number
                    ( COMMA-WSP? number COMMA-WSP? number )? WSP* ")"
skewX           ::= "skewX"    WSP* "(" WSP* number WSP* ")"
skewY           ::= "skewY"    WSP* "(" WSP* number WSP* ")"
COMMA-WSP       ::= ( WSP+ ","? WSP* ) | ( "," WSP* )
WSP             ::= #x20 | #x9 | #xD | #xA
number          ::= sign? digit+ ( "." digit* )? ( [eE] sign? digit+ )?
sign            ::= "+" | "-"
digit           ::= [0-9]
```

---

## transform-origin

Defined in CSS Transforms §6.1 (`transform-origin` property), not in the SVG
coords chapter, which simply delegates.

```
transform-origin ::=
    [ left | center | right | top | bottom | <length-percentage> ]
  | [ left | center | right | <length-percentage> ]
    [ top | center | bottom | <length-percentage> ] <length>?
  | [ [ center | left | right ] && [ center | top | bottom ] ] <length>?
```

**Initial values differ:**

| Context | Initial value |
|---|---|
| HTML elements | `50% 50% 0` |
| SVG elements (general) | `0 0` |
| Root `<svg>` or `<svg>` child of `<foreignObject>` | `50% 50%` |

Not defined inline in the SVG spec; the SVG UA stylesheet sets `transform-origin: 0 0`
for SVG elements (except the two exceptions above).

---

## viewBox

Defined in SVG2 §8.6. Spec value column (verbatim):

```
viewBox = [ <min-x>,? <min-y>,? <width>,? <height> ]

<min-x>, <min-y>, <width>, <height> = <number>
```

- Separators: whitespace and/or a comma (i.e., the same `comma-wsp` separator as
  transform arguments).
- `<width>` and `<height>` must be ≥ 0; negative value = parse error → attribute
  ignored; zero value → element not rendered.

### EBNF-ready rendering

```ebnf
viewBox         ::= WSP* number COMMA-WSP? number COMMA-WSP? number COMMA-WSP? number WSP*
```

Open datatype: `number` (any CSS `<number>`; real-valued, may be negative for
min-x/min-y, must be non-negative for width/height — constraint, not grammar).

---

## preserveAspectRatio (align + meetOrSlice)

Defined in SVG2 §8.7. Verbatim spec productions:

```
<align> =
    none
    | xMinYMin | xMidYMin | xMaxYMin
    | xMinYMid | xMidYMid | xMaxYMid
    | xMinYMax | xMidYMax | xMaxYMax

<meetOrSlice> = meet | slice
```

Attribute value shape (from §8.7 property table):

```
preserveAspectRatio = <align> <meetOrSlice>?
```

Initial value: `xMidYMid meet`

Default when `<meetOrSlice>` is absent: `meet`.

Note: If `align` is `none` then `<meetOrSlice>` is ignored (non-uniform scaling).

### Closed enumeration of align keywords (10 values)

```
none
xMinYMin  xMidYMin  xMaxYMin
xMinYMid  xMidYMid  xMaxYMid
xMinYMax  xMidYMax  xMaxYMax
```

### Closed enumeration of meetOrSlice keywords (2 values)

```
meet
slice
```

### EBNF-ready rendering (all keywords as terminals)

```ebnf
preserveAspectRatio ::= align ( WSP+ meetOrSlice )?

align         ::= "none"
              | "xMinYMin" | "xMidYMin" | "xMaxYMin"
              | "xMinYMid" | "xMidYMid" | "xMaxYMid"
              | "xMinYMax" | "xMidYMax" | "xMaxYMax"

meetOrSlice   ::= "meet" | "slice"
```

### defer keyword

The `defer` keyword was present in SVG 1.1 §7.8 as an optional prefix on
`<image>` elements, instructing the UA to inherit the aspect-ratio behaviour
from the referenced resource.  **SVG 2 removes `defer` entirely** — it does not
appear anywhere in the SVG2 coords chapter.  No browser implements it with any
effect.  Do not include in grammar.

---

## units

SVG2 §8.9 defers entirely to "CSS Values and Units Module [css-values]".
No SVG-specific unit list is given inline; the spec says each attribute/property
must specify its component value type, which may reference `<length>`,
`<percentage>`, `<angle>`, `<number>`, etc.

### Absolute length units (CSS Values §6)

```
cm | mm | Q | in | pt | pc | px
```

### Relative length units (CSS Values §6) — grouped

Font-relative:
```
em | ex | cap | ch | ic | rem | rex | rcap | rch | ric | lh | rlh
```

Viewport-relative:
```
vw | vh | vi | vb | vmin | vmax
```

Viewport-relative (dynamic, level 4+):
```
svw | svh | lvw | lvh | dvw | dvh
```

Container query units (CSS Containment):
```
cqw | cqh | cqi | cqb | cqmin | cqmax
```

### Angle units

```
deg | grad | rad | turn
```

### SVG transform attribute: units NOT allowed

The SVG presentation attribute `transform` (and `gradientTransform`,
`patternTransform`) uses **bare numbers only**.  Units (`px`, `deg`, etc.) must
not appear inside transform attribute argument lists.  The implied units are
applied by the UA:

- translate arguments → implicit `px`
- rotate/skewX/skewY angle → implicit `deg`
- scale, matrix[a,b,c,d] → dimensionless

This contrasts with the CSS `transform` *property*, where `translate(10px)` and
`rotate(45deg)` are required.

### Context-sensitive constraint note

Although the grammar allows any `<number>` in all positions, the following
semantic constraints apply (for overlay / constraint layer):

- `viewBox` width and height must be ≥ 0
- `scale` sx, sy are typically non-zero (zero is valid but collapses the object)
- `matrix` determinant = 0 is valid markup but produces invisible output

---

## vector-effect

Defined in SVG2 §8.13. Verbatim property value definition:

```
none | [ non-scaling-stroke | non-scaling-size | non-rotation | fixed-position ]+ [ viewport | screen ]?
```

Initial value: `none`

### Breakdown

- **Exclusive base keyword (cannot combine):**
  ```
  none
  ```

- **Combinable effect keywords (1 or more, in any order):**
  ```
  non-scaling-stroke
  non-scaling-size
  non-rotation
  fixed-position
  ```

- **Optional scope qualifier (at most one, appended after effect keywords):**
  ```
  viewport
  screen
  ```

### EBNF-ready rendering

```ebnf
vector-effect ::= "none"
              | effect-kw+ scope-kw?

effect-kw     ::= "non-scaling-stroke"
              | "non-scaling-size"
              | "non-rotation"
              | "fixed-position"

scope-kw      ::= "viewport" | "screen"
```

### Implementation note

The spec itself (§8.13, first paragraph) warns:

> "Values of vector-effect other than non-scaling-stroke and none are at risk
> of being dropped from SVG 2 due to a lack of implementations."

For grammar purposes, include all four effect keywords and both scope qualifiers
as terminals; flag the at-risk keywords in the constraint overlay.

---

## gradientTransform / patternTransform reuse note

From SVG2 §8.14.3 (SVGAnimatedTransformList) and the pservers chapter:

> "`SVGAnimatedTransformList` is used to reflect the transform property and its
> corresponding presentation attribute (which, depending on the element, is
> `transform`, `gradientTransform` or `patternTransform`)."

Both `gradientTransform` (on `<linearGradient>`, `<radialGradient>`) and
`patternTransform` (on `<pattern>`) are **presentation attributes for the CSS
`transform` property**.  They accept the identical `transform-list` grammar.
The EBNF grammar node `transform-list` is therefore shared without modification.

---

## Open datatypes used

| Leaf | Description |
|---|---|
| `<number>` | CSS `<number>` token: optional sign, digits, optional decimal, optional exponent. No unit suffix. |
| `<length>` | CSS `<length>` — number + unit suffix (used in CSS transform *property* only, not in SVG attribute). |
| `<length-percentage>` | `<length>` or `<percentage>` — CSS property only. |
| `<angle>` | CSS `<angle>` — number + `deg`/`grad`/`rad`/`turn` (CSS property only; bare number = deg in SVG attribute). |

All four are fully defined in CSS Values and Units; they are leaves for the
SVG grammar — do not expand them further.

---

## Discrepancies, doc gaps & roadblocks

### D1 — SVG transform attribute vs CSS transform property (CRITICAL)

**The transform presentation attribute and the CSS transform property have
different grammars.**

| Aspect | SVG `transform` attribute | CSS `transform` property |
|---|---|---|
| Source spec | CSS Transforms §6.2 | CSS Transforms §3 + §7 |
| Functions allowed | `matrix`, `translate`, `scale`, `rotate`, `skewX`, `skewY` only | All of the above PLUS `translateX()`, `translateY()`, `scaleX()`, `scaleY()`, `skew()`, `rotate3d()`, `perspective()`, … |
| Units on arguments | Bare numbers only (units implied) | Explicit units required (`px`, `deg`, etc.) |
| Separators | `comma-wsp` | CSS whitespace only |

**Grammar decision:** Define two separate grammar rules:
`svg-transform-attr` (this file's BNF) and `css-transform-prop` (CSS Transforms §3).
The SVG attribute grammar is the one needed for XML attribute parsing.

### D2 — SVG2 drops the SVG1.1 transform BNF

SVG 1.1 contained its own inline BNF in §7.6 (identical to the grammar now in
CSS Transforms §6.2). SVG 2 §8.5 simply says "as defined in [css-transforms-1]"
with no inline copy. There is no discrepancy in the grammar itself — only in
where it lives.

### D3 — defer keyword removed in SVG2

`preserveAspectRatio="defer xMidYMid meet"` was valid SVG 1.1 (for `<image>`).
SVG 2 removes `defer` with no migration path.
Decision: exclude `defer` from the grammar.

### D4 — transform-origin initial value differs for SVG elements

CSS Transforms §6.1 and MDN confirm that the UA stylesheet sets
`transform-origin: 0 0` for SVG elements (except root `<svg>` and `<svg>`
inside `<foreignObject>`, which use `50% 50%`).  This is a UA stylesheet rule,
not a property grammar difference.  The property grammar is identical for SVG
and HTML; the behavioural difference belongs in the constraint overlay.

### D5 — vector-effect at-risk keywords

`non-scaling-size`, `non-rotation`, `fixed-position`, `viewport`, `screen` are
spec-marked at-risk.  Chrome, Firefox, and Safari only implement
`non-scaling-stroke` and `none` as of 2025.  For the fuzz grammar, emit all
values but mark the four at-risk ones as unlikely to produce valid rendering.

### D6 — units in SVG presentation attributes (scope)

SVG2 §8.9 says SVG follows CSS Values and Units for all attributes/properties.
However, "all units" applies to *length-valued SVG presentation attributes*
(like `width`, `height`, `x`, `y`, `stroke-width`, etc.) — **not** to the
transform attribute, whose grammar explicitly uses bare numbers.  Care needed
not to conflate "SVG supports CSS length units" with "SVG transform attribute
accepts CSS units".

### D7 — no formal BNF for viewBox in the spec

The spec gives only: `[<min-x>,? <min-y>,? <width>,? <height>]` with
`<min-x>, <min-y>, <width>, <height> = <number>`.  The comma/whitespace
separator rule is described narratively ("separated by whitespace and/or a
comma") but references the same `comma-wsp` pattern as the transform grammar.
The EBNF rendering in this file is a reasonable extrapolation; there is no
single normative BNF to transcribe verbatim.

### D8 — gradientTransform / patternTransform have no separate grammar

Both are described in pservers as "presentation attribute for the transform
property", pointing back to the same `transform-list` grammar.  No rewriting
needed.
