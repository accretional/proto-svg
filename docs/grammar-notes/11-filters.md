# Filter effects grammar notes

## Source

W3C Filter Effects Module Level 1, Working Draft 18 December 2018
https://www.w3.org/TR/2018/WD-filter-effects-1-20181218/
Spec file: `/Volumes/wd_office_3/Projects/proto-svg/docs/specs/cache/filter-effects-1.txt`

---

## `filter` element

### Content model

Any number of the following elements, in any order:
- Descriptive elements: `desc`, `title`, `metadata`
- Filter primitives: `feBlend`, `feFlood`, `feColorMatrix`, `feComponentTransfer`, `feComposite`, `feConvolveMatrix`, `feDiffuseLighting`, `feDisplacementMap`, `feDropShadow`, `feGaussianBlur`, `feImage`, `feMerge`, `feMorphology`, `feOffset`, `feSpecularLighting`, `feTile`, `feTurbulence`
- `animate`, `script`, `set`

### Attributes

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `filterUnits` | `"userSpaceOnUse" \| "objectBoundingBox"` | `objectBoundingBox` | Coordinate system for x, y, width, height |
| `primitiveUnits` | `"userSpaceOnUse" \| "objectBoundingBox"` | `userSpaceOnUse` | Coordinate system for filter primitive length values |
| `x` | `<length-percentage>` | `-10%` | Filter region x |
| `y` | `<length-percentage>` | `-10%` | Filter region y |
| `width` | `<length-percentage>` | `120%` | Filter region width |
| `height` | `<length-percentage>` | `120%` | Filter region height |
| `href` / `xlink:href` | `<url>` | — | Reference to another filter element for inheritance (via SVGURIReference) |

Note: `filterRes` (deprecated, removed) was `<number-optional-number>`. Do not include in grammar.

Note: The `filter` element also accepts presentation attributes (large set), `id`, `class`, `style`, `xml:base`, `xml:lang`, `xml:space`, `externalResourcesRequired`. These are handled by global attribute groups.

### Context-sensitive constraints

- `filterUnits` scopes the coordinate space for `x`, `y`, `width`, `height` only; `primitiveUnits` scopes primitive coordinates and lengths.
- `filter` elements are never rendered directly; only used by reference via the `filter` property.
- Properties inherit into `filter` from ancestors, but NOT from the referencing element.

---

## Common filter-primitive attributes

All filter primitive elements share these attributes (defined in `SVGFilterPrimitiveStandardAttributes`):

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `x` | `<length-percentage>` | `0%` | Primitive subregion min-x |
| `y` | `<length-percentage>` | `0%` | Primitive subregion min-y |
| `width` | `<length-percentage>` | `100%` | Primitive subregion width; zero/negative disables the primitive |
| `height` | `<length-percentage>` | `100%` | Primitive subregion height; zero/negative disables the primitive |
| `result` | `<filter-primitive-reference>` | — | Named output; `<filter-primitive-reference>` is a `<custom-ident>` |

### `in` attribute (common to most primitives, defined at §9.2)

```
in = "SourceGraphic | SourceAlpha | BackgroundImage | BackgroundAlpha |
       FillPaint | StrokePaint | <filter-primitive-reference>"
```

- If omitted on first primitive: defaults to `SourceGraphic`
- If omitted on subsequent primitive: uses result of immediately preceding primitive
- Forward references are not allowed; references to non-existent results treated as unspecified

---

## `<in>` value type

```
<fe-in-value> ::= "SourceGraphic"
                | "SourceAlpha"
                | "BackgroundImage"
                | "BackgroundAlpha"
                | "FillPaint"
                | "StrokePaint"
                | <filter-primitive-reference>

<filter-primitive-reference> ::= <custom-ident>
```

The six built-in keywords are **terminals**. The `<filter-primitive-reference>` arm is open — it matches the value of a prior `result` attribute on a sibling primitive within the same `filter` element.

Context-sensitive constraint: the referenced name must have appeared earlier in the same `filter` element's child list (no forward references).

---

## Filter primitives

### `feBlend`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | First input image |
| `in2` | `<fe-in-value>` | — | Second input (required for blend) |
| `mode` | `<blend-mode>` | `normal` | See keywords below |
| `no-composite` | boolean presence attribute | absent | When present, suppresses alpha compositing with Source Over |

**`mode` keyword set (terminals):**
```
<blend-mode> ::= "normal" | "multiply" | "screen" | "darken" | "lighten"
               | "overlay" | "color-dodge" | "color-burn" | "hard-light"
               | "soft-light" | "difference" | "exclusion"
               | "hue" | "saturation" | "color" | "luminosity"
```

Confirmed via DOM IDL enums in Appendix B (SVG_FEBLEND_MODE_* constants, 16 named values).

**Verification note:** The prose at §9.5 says `mode = "<blend-mode>"` and defers to Compositing and Blending Level 1 for the full list. The IDL appendix enumerates all 16 modes. MDN confirms all 16 are supported in browsers. The spec's example SVG only shows 5 (normal/multiply/screen/darken/lighten) but those are just examples. **Grammar decision:** include all 16 as terminals.

**Context-sensitive constraints:**
- `in` = source Cs; `in2` = backdrop Cb per Compositing and Blending [COMPOSITING-1].
- `no-composite` is not animatable.

---

### `feColorMatrix`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input image |
| `type` | `"matrix" \| "saturate" \| "hueRotate" \| "luminanceToAlpha"` | `matrix` | Operation type |
| `values` | `<list-of-numbers>` | see below | Depends on `type` |

**`values` semantics by `type`:**
- `type="matrix"`: list of exactly 20 `<number>` values (5x4 matrix)
- `type="saturate"`: single `<number>` value (s); range note: 0=grayscale, 1=unchanged, >1=oversaturate
- `type="hueRotate"`: single `<number>` value (degrees)
- `type="luminanceToAlpha"`: `values` not applicable (ignored)

**Defaults for `values`:**
- `matrix`: identity matrix
- `saturate`: 1
- `hueRotate`: 0

If `values` count does not match the required count for `type`, the primitive acts as pass-through.

**Context-sensitive constraints:**
- `values` interpretation is gated on `type`.
- Grammar should model as: `feColorMatrix-values ::= <number>+ /* count constraint noted in overlay */`

---

### `feComponentTransfer`

**Content model:** Any number of descriptive elements, `feFuncR`, `feFuncG`, `feFuncB`, `feFuncA`, `script`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default |
|---|---|---|
| `in` | `<fe-in-value>` | (see common) |

No other primitive-specific attributes. The child transfer function elements `feFuncR`, `feFuncG`, `feFuncB`, `feFuncA` carry the actual parameters.

If any transfer function child is absent, it is treated as `type="identity"`.
If multiple children of the same kind appear, the last occurrence wins.

---

### `feFuncR` / `feFuncG` / `feFuncB` / `feFuncA`  (transfer function elements)

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Attributes (transfer function element attributes group):**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `type` | `"identity" \| "table" \| "discrete" \| "linear" \| "gamma"` | `identity` | Determines which other attrs apply |
| `tableValues` | `<list-of-numbers>` | (empty) | Used when `type="table"` or `type="discrete"`; empty = identity |
| `slope` | `<number>` | `1` | Used when `type="linear"` |
| `intercept` | `<number>` | `0` | Used when `type="linear"` |
| `amplitude` | `<number>` | `1` | Used when `type="gamma"` |
| `exponent` | `<number>` | `1` | Used when `type="gamma"` |
| `offset` | `<number>` | `0` | Used when `type="gamma"` |

**Formulas:**
- `identity`: C' = C
- `table`: piecewise linear interpolation between n+1 values
- `discrete`: step function with n values
- `linear`: C' = slope * C + intercept
- `gamma`: C' = amplitude * pow(C, exponent) + offset

**Context-sensitive constraints:**
- `tableValues` is only meaningful when `type="table"` or `type="discrete"`.
- `slope`, `intercept` are only meaningful when `type="linear"`.
- `amplitude`, `exponent`, `offset` are only meaningful when `type="gamma"`.

---

### `feComposite`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Source input |
| `in2` | `<fe-in-value>` | — | Destination input |
| `operator` | `"over" \| "in" \| "out" \| "atop" \| "xor" \| "lighter" \| "arithmetic"` | `over` | Compositing operation |
| `k1` | `<number>` | `0` | Only used when `operator="arithmetic"` |
| `k2` | `<number>` | `0` | Only used when `operator="arithmetic"` |
| `k3` | `<number>` | `0` | Only used when `operator="arithmetic"` |
| `k4` | `<number>` | `0` | Only used when `operator="arithmetic"` |

**Arithmetic formula:** `result = k1*i1*i2 + k2*i1 + k3*i2 + k4`

**DISCREPANCY:** The spec prose at §9.8 lists `"over | in | out | atop | xor | lighter | arithmetic"` (7 keywords, including `lighter`). The DOM IDL appendix only defines 6 constants: `SVG_FECOMPOSITE_OPERATOR_OVER/IN/OUT/ATOP/XOR/ARITHMETIC` — no `lighter` constant. MDN and browser implementations do support `lighter`. **Grammar decision:** Include `lighter` as a terminal per the prose; note the IDL gap.

**Context-sensitive constraints:**
- `k1`..`k4` only apply when `operator="arithmetic"`.
- The result pixel is clamped to [0,1].

---

### `feConvolveMatrix`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input image |
| `order` | `<number-optional-number>` | `3` | Integer(s) > 0; specifies orderX [orderY]; truncates if non-integer |
| `kernelMatrix` | `<list-of-numbers>` | — | Must have orderX × orderY entries; pass-through if wrong count |
| `divisor` | `<number>` | sum of kernelMatrix (or 1 if sum=0) | If 0, uses default |
| `bias` | `<number>` | `0` | Addend after divisor |
| `targetX` | `<integer>` | `floor(orderX/2)` | 0 ≤ targetX < orderX |
| `targetY` | `<integer>` | `floor(orderY/2)` | 0 ≤ targetY < orderY |
| `edgeMode` | `"duplicate" \| "wrap" \| "none"` | `duplicate` | Edge extension mode |
| `kernelUnitLength` | `<number-optional-number>` | (1 device pixel) | DEPRECATED; negative/zero = use default |
| `preserveAlpha` | `"false" \| "true"` | `false` | Whether to apply kernel to alpha channel |

**Context-sensitive constraints:**
- `kernelMatrix` must have exactly `orderX * orderY` entries; mismatch = pass-through.
- `targetX` must satisfy `0 <= targetX < orderX`; `targetY` must satisfy `0 <= targetY < orderY`.
- `kernelUnitLength` is deprecated; browsers may ignore it.

---

### `feDiffuseLighting`

**Content model:** Any number of descriptive elements, `script`, and **exactly one** light source element (`feDistantLight`, `fePointLight`, or `feSpotLight`), in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input bump map (alpha used) |
| `surfaceScale` | `<number>` | `1` | Height of surface when alpha=1 |
| `diffuseConstant` | `<number>` | `1` | kd (Phong lighting); any non-negative number |
| `kernelUnitLength` | `<number-optional-number>` | (1 device pixel) | DEPRECATED; dx and dy for Sobel gradient |

**Light color:** determined by the `lighting-color` CSS property.

**Context-sensitive constraints:**
- Requires exactly one light source child element.
- Output is always opaque RGBA (alpha = 1.0 everywhere).

---

### `feDisplacementMap`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Source image to displace |
| `in2` | `<fe-in-value>` | — | Displacement map image |
| `scale` | `<number>` | `0` | Displacement scale factor |
| `xChannelSelector` | `"R" \| "G" \| "B" \| "A"` | `A` | Channel from in2 to use for X displacement |
| `yChannelSelector` | `"R" \| "G" \| "B" \| "A"` | `A` | Channel from in2 to use for Y displacement |

**Formula:** `P'(x,y) ← P(x + scale*(XC(x,y) - 0.5), y + scale*(YC(x,y) - 0.5))`

**Context-sensitive constraints:**
- `color-interpolation-filters` applies only to `in2`; `in` remains in its current color space.
- Security: If `in2` is tainted, the primitive acts as pass-through.
- Open issue: implementations do not match specification (GitHub issue #113).

---

### `feDropShadow`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input image |
| `dx` | `<number>` | `2` | X offset of shadow |
| `dy` | `<number>` | `2` | Y offset of shadow |
| `stdDeviation` | `<number-optional-number>` | `2` | Blur radius (forwarded to internal feGaussianBlur) |

**Color/opacity:** Set via `flood-color` and `flood-opacity` CSS properties on the element.

This element is a shorthand defined in terms of: `feGaussianBlur` + `feOffset` + `feFlood` + `feComposite(operator="in")` + `feMerge`. The internal structure is not exposed to DOM.

---

### `feFlood`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:** None (no `in` attribute).

**Properties used:**
- `flood-color`: color to fill the primitive subregion
- `flood-opacity`: opacity

The element fills the entire filter primitive subregion with `flood-color` at `flood-opacity`.

---

### `feGaussianBlur`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input image |
| `stdDeviation` | `<number-optional-number>` | `0` | Standard deviation; if two values, first=X, second=Y; zero or negative = disables effect |
| `edgeMode` | `"duplicate" \| "wrap" \| "none"` | `none` | Edge extension mode |

**DISCREPANCY:** `feConvolveMatrix.edgeMode` defaults to `"duplicate"`. `feGaussianBlur.edgeMode` defaults to `"none"`. These are distinct and must not be conflated in the grammar.

**Note:** The CSS `blur()` filter function shorthand maps to `feGaussianBlur` with `edgeMode="none"` for the `filter` property, but `edgeMode="duplicate"` for the CSS `filter()` image function.

---

### `feImage`

**Content model:** Any number of descriptive elements, `animate`, `animateTransform`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `href` | `<url>` | — | Reference to image or element |
| `xlink:href` | `<url>` | — | Deprecated; overridden by `href` if both present |
| `preserveAspectRatio` | `[defer] <align> [<meetOrSlice>]` | `xMidYMid meet` | Same semantics as SVG `image` element |
| `crossorigin` | `"anonymous" \| "use-credentials"` | — | CORS settings; not animatable |

**Context-sensitive constraints:**
- No `in` attribute (feImage generates content, it does not take filter input).
- If href references a non-existent, zero-size, or unsupported resource, fills with transparent black.

---

### `feMerge`

**Content model:** Any number of descriptive elements, `feMergeNode`, `script`, in any order.

**Primitive-specific attributes:** None beyond common filter primitive attributes (no `in` on `feMerge` itself).

Composites all `feMergeNode` children using the "over" operator, bottom to top.

---

### `feMergeNode`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Attributes:**

| Attribute | Value syntax | Default |
|---|---|---|
| `in` | `<fe-in-value>` | (implicit previous result) |

Category: None (not a filter primitive; it is a child element of `feMerge`). Not in `SVGFilterPrimitiveStandardAttributes`. Does NOT have x/y/width/height/result attributes.

---

### `feMorphology`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input image |
| `operator` | `"erode" \| "dilate"` | `erode` | Thinning vs. fattening |
| `radius` | `<number-optional-number>` | `0` | x-radius [y-radius]; zero or negative = disables effect |

---

### `feOffset`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input image |
| `dx` | `<number>` | `0` | X offset in primitiveUnits coordinate space |
| `dy` | `<number>` | `0` | Y offset in primitiveUnits coordinate space |

---

### `feSpecularLighting`

**Content model:** Any number of descriptive elements, `script`, and **exactly one** light source element (`feDistantLight`, `fePointLight`, or `feSpotLight`), in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `in` | `<fe-in-value>` | (see common) | Input bump map (alpha used) |
| `surfaceScale` | `<number>` | `1` | Height of surface when alpha=1 |
| `specularConstant` | `<number>` | `1` | ks (Phong lighting); any non-negative number |
| `specularExponent` | `<number>` | `1` | Shininess exponent (Blinn-Phong) |
| `kernelUnitLength` | `<number-optional-number>` | (1 device pixel) | DEPRECATED |

**Light color:** determined by the `lighting-color` CSS property.

**Context-sensitive constraints:**
- Requires exactly one light source child element.
- Output is non-opaque RGBA (alpha = max of Sr, Sg, Sb) — differs from `feDiffuseLighting`.

---

### `feTile`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default |
|---|---|---|
| `in` | `<fe-in-value>` | (see common) |

Fills its filter primitive subregion with a tiled pattern of the input image.

---

### `feTurbulence`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Primitive-specific attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `baseFrequency` | `<number-optional-number>` | `0` | Base frequency [X [Y]]; negative values unsupported |
| `numOctaves` | `<integer>` | `1` | Number of Perlin noise octaves; negative unsupported |
| `seed` | `<number>` | `0` | Initial seed for PRNG; truncated to integer |
| `stitchTiles` | `"stitch" \| "noStitch"` | `noStitch` | Tile-border continuity |
| `type` | `"fractalNoise" \| "turbulence"` | `turbulence` | Noise function type |

No `in` attribute — `feTurbulence` generates content independently.

---

## Light source elements

Light sources are children of `feDiffuseLighting` and `feSpecularLighting`.

### `feDistantLight`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `azimuth` | `<number>` | `0` | Direction angle on XY plane (degrees, clockwise from X axis) |
| `elevation` | `<number>` | `0` | Direction angle from XY plane toward Z axis (degrees) |

---

### `fePointLight`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `x` | `<number>` | `0` | X position in primitiveUnits space |
| `y` | `<number>` | `0` | Y position in primitiveUnits space |
| `z` | `<number>` | `0` | Z position; positive Z toward viewer |

---

### `feSpotLight`

**Content model:** Any number of descriptive elements, `animate`, `script`, `set`, in any order.

**Attributes:**

| Attribute | Value syntax | Default | Notes |
|---|---|---|---|
| `x` | `<number>` | `0` | Light source X position |
| `y` | `<number>` | `0` | Light source Y position |
| `z` | `<number>` | `0` | Light source Z position |
| `pointsAtX` | `<number>` | `0` | Target X position |
| `pointsAtY` | `<number>` | `0` | Target Y position |
| `pointsAtZ` | `<number>` | `0` | Target Z position |
| `specularExponent` | `<number>` | `1` | Focus exponent for spot cone |
| `limitingConeAngle` | `<number>` | (none) | Angle of spot cone (degrees); optional; if absent, no cone limiting |

**Note:** `specularExponent` on `feSpotLight` is distinct from `specularExponent` on `feSpecularLighting`; they serve different roles.

---

## Filter properties

### `flood-color`

```
Name:  flood-color
Value: <color>
Initial: black
Applies to: feFlood, feDropShadow elements
Inherited: no
Computed value: as specified
Animatable: as by computed value
```

A presentation attribute for SVG. Accepts any CSS `<color>` value.

### `flood-opacity`

```
Name:  flood-opacity
Value: <alpha-value>
Initial: 1
Applies to: feFlood, feDropShadow elements
Inherited: no
Computed value: specified value converted to number, clamped to [0,1]
Animatable: by computed value
```

`<alpha-value>` is `<number>` or `<percentage>` (CSS Color Level 4). If `flood-color` includes its own alpha, that is multiplied with `flood-opacity`.

### `lighting-color`

```
Name:  lighting-color
Value: <color>
Initial: white
Applies to: feDiffuseLighting, feSpecularLighting elements
Inherited: no
Computed value: as specified
Animatable: as by computed value
```

### `color-interpolation-filters`

```
Name:  color-interpolation-filters
Value: auto | sRGB | linearRGB
Initial: linearRGB
Applies to: all filter primitives
Inherited: yes
Computed value: as specified
Animatable: no
```

**Keywords:**
- `auto`: UA may choose either sRGB or linearRGB
- `sRGB`: filter operations in sRGB color space
- `linearRGB`: filter operations in linearized RGB color space

**Critical note:** `color-interpolation-filters` has initial value `linearRGB`, whereas `color-interpolation` has initial value `sRGB`. These are separate properties.

Has no effect on `feOffset`, `feImage`, `feTile`, `feFlood` (no color arithmetic performed). Has no effect on CSS Filter Functions (those always operate in sRGB).

---

## `filter` property value (CSS)

```
Name: filter
Value: none | <filter-value-list>
Initial: none
Inherited: no
```

```
<filter-value-list> = [ <filter-function> | <url> ]+ 
```

```
<filter-function> = <blur()>
                  | <brightness()>
                  | <contrast()>
                  | <drop-shadow()>
                  | <grayscale()>
                  | <hue-rotate()>
                  | <invert()>
                  | <opacity()>
                  | <sepia()>
                  | <saturate()>
```

**Individual filter function grammars:**

```
blur()          = blur( <length>? )
                  /* default: 0px; no percentages; no negatives */

brightness()    = brightness( <number-percentage>? )
                  /* default: 1; no negatives */

contrast()      = contrast( <number-percentage>? )
                  /* default: 1; no negatives */

drop-shadow()   = drop-shadow( <color>? && <length>{2,3} )
                  /* <color> may appear before or after <length> values
                     3rd <length> is stdDeviation (not spread like box-shadow)
                     default: missing lengths=0, missing color=currentColor
                     no spread values, no multiple shadows */

grayscale()     = grayscale( <number-percentage>? )
                  /* default when omitted: 1; initial for interpolation: 0; no negatives */

hue-rotate()    = hue-rotate( [ <angle> | <zero> ]? )
                  /* default: 0deg; unit may be omitted if value is zero */

invert()        = invert( <number-percentage>? )
                  /* default when omitted: 1; initial for interpolation: 0; no negatives */

opacity()       = opacity( <number-percentage>? )
                  /* default: 1; no negatives; values >1 are clamped */

saturate()      = saturate( <number-percentage>? )
                  /* default: 1; no negatives */

sepia()         = sepia( <number-percentage>? )
                  /* default when omitted: 1; initial for interpolation: 0; no negatives */
```

**CSS `filter()` image function:**
```
filter() = filter( [ <image> | <string> ], <filter-value-list> )
```

---

## Open datatypes used

| Name | Definition |
|---|---|
| `<number>` | CSS `<number>` (unrestricted unless noted; range constraints are context-sensitive) |
| `<integer>` | CSS `<integer>` |
| `<length>` | CSS `<length>` |
| `<length-percentage>` | CSS `<length>` or `<percentage>` |
| `<percentage>` | CSS `<percentage>` |
| `<angle>` | CSS `<angle>` |
| `<zero>` | The unitless number `0` (special alias per CSS) |
| `<color>` | CSS `<color>` (full CSS color syntax) |
| `<alpha-value>` | CSS Color Level 4 `<alpha-value>` = `<number>` or `<percentage>` |
| `<url>` | CSS `<url>` / `url()` function |
| `<image>` | CSS `<image>` |
| `<custom-ident>` | CSS `<custom-ident>` (used for `result` / `<filter-primitive-reference>`) |
| `<number-optional-number>` | `<number> <number>?` — spec-defined at §7; if second omitted, it defaults to first |
| `<list-of-numbers>` | Space- and/or comma-separated sequence of `<number>` values; length varies |
| `<list-of-integers>` | Space- and/or comma-separated sequence of `<integer>` values |
| `<blend-mode>` | Enumerated set — fully listed above (16 keywords) |
| `<align>` / `<meetOrSlice>` | From SVG `preserveAspectRatio` (see spec 04-coords-transforms) |
| `<number-percentage>` | `<number>` or `<percentage>` |

**`<number-optional-number>` pattern** appears on: `order` (feConvolveMatrix), `kernelUnitLength` (feConvolveMatrix, feDiffuseLighting, feSpecularLighting), `stdDeviation` (feGaussianBlur, feDropShadow), `baseFrequency` (feTurbulence), `radius` (feMorphology). In ALL cases the second value defaults to the first value when omitted.

---

## Discrepancies, doc gaps, and roadblocks

### D1: `feComposite operator` — `lighter` missing from IDL
- **Spec prose (§9.8):** `operator = "over | in | out | atop | xor | lighter | arithmetic"` — 7 keywords including `lighter`.
- **DOM IDL appendix:** Only defines 6 `SVG_FECOMPOSITE_OPERATOR_*` constants; `lighter` has no IDL constant.
- **MDN / browsers:** `lighter` is supported in browsers as Porter-Duff "plus" operation.
- **Grammar decision:** Include `lighter` as a terminal in the EBNF. Note the IDL gap.

### D2: `feBlend mode` — spec text is underspecified vs. IDL
- **Spec prose (§9.5):** says `mode = "<blend-mode>"` with a deferred reference to COMPOSITING-1.
- **IDL appendix:** Enumerates all 16 modes explicitly.
- **MDN:** All 16 are supported. The SVG 1.1 spec only had 5 modes; the extension to 16 is a level-1 addition.
- **Grammar decision:** Use the IDL-enumerated 16-keyword set as authoritative.

### D3: `feGaussianBlur.edgeMode` vs `feConvolveMatrix.edgeMode` — different defaults
- `feConvolveMatrix.edgeMode` default = `duplicate`
- `feGaussianBlur.edgeMode` default = `none`
- Same keyword set `"duplicate | wrap | none"` but different semantics for the initial value. Do not share a single `edgeMode` production if defaults are encoded in grammar; they must be separate attributes.

### D4: `feDisplacementMap` — spec/implementation mismatch
- Spec §9.11 explicitly notes: "Implementations do not match specification." (GitHub issue #113)
- Grammar should reflect the spec as written; add an overlay note about the implementation discrepancy.

### D5: `kernelUnitLength` is deprecated
- Appears on `feConvolveMatrix`, `feDiffuseLighting`, `feSpecularLighting`.
- Spec notes: "This attribute is deprecated and will be removed." for all three.
- **Grammar decision:** Include in grammar for completeness (must parse existing SVG), but mark deprecated in overlay notes.

### D6: `seed` type ambiguity
- **Spec (§9.21):** `seed = "<number>"` but semantics say it is truncated to integer before use.
- **DOM IDL:** `readonly attribute SVGAnimatedNumber seed` — stored as number.
- **Grammar decision:** Use `<number>` as the syntactic type; note integer-truncation in overlay.

### D7: `numOctaves` typed as integer
- Spec correctly says `numOctaves = "<integer>"` at §9.21 and DOM IDL uses `SVGAnimatedInteger`.
- Grammar: `<integer>`, non-negative.

### D8: `feImage` has no `in` attribute
- Unlike all other filter primitives, `feImage` does not accept `in`. It generates content from an external source.
- Similarly `feTurbulence` and `feFlood` have no `in` attribute.

### D9: `drop-shadow()` vs `feDropShadow` — parameter order differs
- CSS `drop-shadow()` grammar: `<color>? && <length>{2,3}` — color and lengths interchangeable order.
- `feDropShadow` element: `dx`, `dy`, `stdDeviation` are separate numeric attributes.
- The CSS function's optional 3rd `<length>` is stdDeviation (not a "spread" as in `box-shadow`).

### D10: `grayscale()`, `sepia()`, `invert()` default-vs-interpolation mismatch
- Default value when function argument omitted: `1` (full effect)
- Initial value for CSS animation interpolation: `0` (no effect)
- This is intentional per spec prose but unusual. Grammar encodes the syntactic default; the animation initial is a separate concern for the overlay.

### D11: `feBlend no-composite` attribute
- This is a boolean presence attribute (not an enumeration). It is not animatable.
- It was added as a level-1 extension to feBlend to avoid "double-compositing". MDN does not document browser support clearly; this attribute appears to have minimal implementation.
- **Grammar decision:** Include as an optional boolean attribute (presence = true).

### D12: `filterRes` removed
- SVG 1.1 had `filterRes = "<number-optional-number>"` on the `filter` element.
- Filter Effects Level 1 explicitly removed it. Do NOT include in grammar.

### D13: `feImage crossorigin` attribute
- Value syntax: `"anonymous | use-credentials"`
- Not animatable (per spec).
- IDL reflects it as `SVGAnimatedString crossOrigin` (note different casing: `crossOrigin` in IDL vs `crossorigin` in HTML). The IDL attribute "must reflect the crossorigin content attribute, limited to only known values."

### D14: `feMergeNode` is NOT a filter primitive
- Category in spec: "None" (§9.16.1).
- Does NOT include `SVGFilterPrimitiveStandardAttributes` — no `x`, `y`, `width`, `height`, `result`.
- Only attribute beyond core: `in`.

### D15: `feSpecularLighting.kernelUnitLength` — not in IDL but is in prose
- The prose §9.19 lists `kernelUnitLength` as an attribute.
- The IDL for `SVGFESpecularLightingElement` includes `kernelUnitLengthX` and `kernelUnitLengthY` as separate properties.
- Same pattern as `feDiffuseLighting` and `feConvolveMatrix` — the single `<number-optional-number>` attribute maps to two DOM properties.

### D16: `feDropShadow` — `flood-color` and `flood-opacity` are CSS properties, not XML attributes
- They appear as presentation attributes on the element (settable as attributes), but are fundamentally CSS properties.
- Grammar should model them as CSS property values, not as primitive-specific XML attribute syntax.

---

## Summary — primitives documented

Total filter primitive elements: **17**
1. `feBlend`
2. `feColorMatrix`
3. `feComponentTransfer` (+ children: `feFuncR`, `feFuncG`, `feFuncB`, `feFuncA`)
4. `feComposite`
5. `feConvolveMatrix`
6. `feDiffuseLighting`
7. `feDisplacementMap`
8. `feDropShadow`
9. `feFlood`
10. `feGaussianBlur`
11. `feImage`
12. `feMerge` (+ child: `feMergeNode`)
13. `feMorphology`
14. `feOffset`
15. `feSpecularLighting`
16. `feTile`
17. `feTurbulence`

Light source elements (children of feDiffuseLighting/feSpecularLighting): `feDistantLight`, `fePointLight`, `feSpotLight`

Transfer function elements (children of feComponentTransfer): `feFuncR`, `feFuncG`, `feFuncB`, `feFuncA`
