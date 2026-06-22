# Masking & clipping grammar notes

## Source

W3C CSS Masking Module Level 1, Candidate Recommendation Draft 5 August 2021
Extracted from: `/docs/specs/cache/css-masking-1.txt` (2850 lines, read in full)
Cross-checked against MDN Web Docs (April–June 2026 snapshot) and known browser behavior.

---

## Elements

### `clipPath`

#### Categories

`container elements`, `never-rendered element`

#### Content model

Any number of the following elements, in any order:

| Group | Elements |
|-------|----------|
| descriptive | `desc`, `title`, `metadata` |
| animation | `animate`, `animateColor`, `animateMotion`, `animateTransform`, `set` |
| shape | `circle`, `ellipse`, `line`, `path`, `polygon`, `polyline`, `rect` |
| text | `text` |
| use | `use` |
| script | `script` |

Note: `use` inside `clipPath` must directly reference a `path`, `text`, or basic shape element; indirect reference is an error (clipPath must be ignored, per spec). In practice: Firefox disables rendering of the referencing element on content-model violation; no browser currently ignores `clipPath` on `use` with indirect reference (open spec issue #17).

Children can carry `clip-path` and `clip-rule` presentation attributes (see Constraints below).

#### Own attributes

```
clipPathUnits = "userSpaceOnUse" | "objectBoundingBox"
```

Default: `"userSpaceOnUse"`. Animatable: yes.

Also accepts the full set of SVG presentation attributes (see shared group below), plus `class`, `style`, `externalResourcesRequired`, `transform`, and conditional-processing / core attribute groups.

#### DOM interface

`SVGClipPathElement { clipPathUnits: SVGAnimatedEnumeration; transform: SVGAnimatedTransformList }`

---

### `mask`

#### Categories

`container elements`, `never-rendered element`

#### Content model

Any number of the following elements, in any order:

| Group | Elements |
|-------|----------|
| animation | `animate`, `animateColor`, `animateMotion`, `animateTransform`, `set` |
| descriptive | `desc`, `title`, `metadata` |
| shape | `circle`, `ellipse`, `line`, `path`, `polygon`, `polyline`, `rect` |
| structural | `defs`, `g`, `svg`, `symbol`, `use` |
| gradient | `linearGradient`, `radialGradient` |
| other | `a`, `clipPath`, `color-profile`, `cursor`, `filter`, `font`, `font-face`, `foreignObject`, `image`, `marker`, `mask`, `pattern`, `script`, `style`, `switch`, `text`, `view`, `altGlyphDef` |

Far broader content model than `clipPath`; essentially any graphics, container, or descriptive element.

#### Own attributes

```
maskUnits        = "userSpaceOnUse" | "objectBoundingBox"
maskContentUnits = "userSpaceOnUse" | "objectBoundingBox"
x                = <length-percentage>
y                = <length-percentage>
width            = <length-percentage>
height           = <length-percentage>
```

Defaults and rules:

| Attribute | Default | Notes |
|-----------|---------|-------|
| `maskUnits` | `"objectBoundingBox"` | Applies to x/y/width/height coords |
| `maskContentUnits` | `"userSpaceOnUse"` | Applies to coordinate system of contents |
| `x` | `-10%` | Only defaulted when at least one of y/width/height is specified |
| `y` | `-10%` | Only defaulted when at least one of x/width/height is specified |
| `width` | `120%` | Negative or zero disables rendering |
| `height` | `120%` | Negative or zero disables rendering |

All four geometric attributes: Animatable: yes.

Also accepts the full set of SVG presentation attributes, plus `class`, `style`, conditional-processing and core attribute groups.

Note: `opacity`, `filter`, and `display` do not apply to the `mask` element itself.

#### DOM interface

```
SVGMaskElement {
  maskUnits:        SVGAnimatedEnumeration;
  maskContentUnits: SVGAnimatedEnumeration;
  x:      SVGAnimatedLength;
  y:      SVGAnimatedLength;
  width:  SVGAnimatedLength;
  height: SVGAnimatedLength;
}
```

---

### Shared presentation attribute group (both elements)

Both `clipPath` and `mask` accept the full SVG 1.1 presentation attribute set. Key entries relevant to clipping/masking children:

- `clip`, `clip-path`, `clip-rule` (children only; see Constraints)
- `mask`
- `opacity`
- `display`, `visibility`

---

### Constraints (overlay — not encoded in EBNF)

1. **clipPath `use` child**: `use` must reference `path`, `text`, or basic-shape directly. Indirect chain is an error; UA must ignore the `clipPath` (spec); actual browsers vary (see Discrepancies §1).
2. **clip-path on a clipPath element**: Intersects the clipPath's own clipping region with the referenced clipping path.
3. **clip-path on a child of clipPath**: Clips the child before OR'ing with sibling silhouettes.
4. **clip-rule on children**: Only effective on graphics elements *inside* a `clipPath`; ignored on the referencing element. Does not apply to `<basic-shape>`.
5. **mask x/y/width/height**: If *none* of the four are specified, the element imposes no bounding clip (all are omitted together).
6. **maskUnits vs maskContentUnits**: `maskUnits` governs the x/y/width/height box; `maskContentUnits` governs the internal coordinate system of child elements—these are independent.
7. **Referencing constraint (clip-path)**: A `<url>` value in `clip-path` must resolve to a `clipPath` element; otherwise no clipping is applied (not an error, silently ignored).
8. **Referencing constraint (mask-image / mask)**: A `<mask-source>` (`<url>`) should resolve to a `mask` element; a non-`mask` element reference counts as a transparent black layer.

---

## `clip-path` + `<basic-shape>` + `<geometry-box>`

### Property definition

```
Name:    clip-path
Value:   <clip-source> | [ <basic-shape> || <geometry-box> ] | none
Initial: none
Applies: all elements; in SVG: container elements (excl. defs), all graphics elements, use
Inherited: no
Computed: as specified, <url> values made absolute
Animation: by computed value
```

### Value arms

```ebnf
clip-path-value
  = <clip-source>
  | <basic-shape>
  | <geometry-box>
  | <basic-shape> <geometry-box>
  | <geometry-box> <basic-shape>
  | "none"
```

(The `||` combinator means either or both, in any order — split into explicit alternatives for EBNF clarity.)

### `<clip-source>`

```ebnf
<clip-source> = <url>
```

URL must reference a `clipPath` element (constraint, not in formal grammar — see Constraints §7).

### `<geometry-box>`

```ebnf
<geometry-box>
  = <shape-box>
  | "fill-box"
  | "stroke-box"
  | "view-box"

<shape-box>
  = "margin-box"
  | "border-box"
  | "padding-box"
  | "content-box"
```

Full closed keyword set for `<geometry-box>` (7 keywords total):
`"margin-box"`, `"border-box"`, `"padding-box"`, `"content-box"`, `"fill-box"`, `"stroke-box"`, `"view-box"`

Mapping rules (SVG elements without CSS layout box):
- `content-box` and `padding-box` → used as `fill-box`
- `border-box` and `margin-box` → used as `stroke-box`

Mapping rules (HTML elements with CSS layout box):
- `fill-box` → used as `content-box`
- `stroke-box` and `view-box` → used as `border-box`

When `<geometry-box>` is specified alone, it defines the clip shape from the box edges (including border-radius). When combined with `<basic-shape>`, it provides the reference box.

Default reference box when `<geometry-box>` is absent: `border-box`.

### `<basic-shape>` — full grammar

```ebnf
<basic-shape>
  = <inset()>
  | <circle()>
  | <ellipse()>
  | <polygon()>
  | <path()>
  | <shape()>

<inset()>
  = "inset(" <length-percentage>{1,4} [ "round" <border-radius> ]? ")"

<circle()>
  = "circle(" <radial-size>? [ "at" <position> ]? ")"

<ellipse()>
  = "ellipse(" <radial-size>? [ "at" <position> ]? ")"

<polygon()>
  = "polygon(" <fill-rule>? [ "round" <length> ]? ","
               [ <length-percentage> <length-percentage> ]# ")"

<path()>
  = "path(" <fill-rule>? "," <string> ")"

<shape()>
  = "shape(" <fill-rule>? "from" <position> "," <shape-command># ")"

<fill-rule>     = "nonzero" | "evenodd"

<radial-size>
  = <length-percentage>
  | "closest-side" | "farthest-side"

<border-radius>
  = <length-percentage>{1,4} [ "/" <length-percentage>{1,4} ]?
```

Notes on individual functions:

| Function | Notable points |
|----------|---------------|
| `inset()` | 1–4 offsets (top/right/bottom/left shorthand); optional `round` for border-radius |
| `circle()` | Single radius + optional center position; both are optional (defaults to 50% 50%) |
| `ellipse()` | Two radii (x then y) + optional center; both pairs optional |
| `polygon()` | `<fill-rule>?` before comma; new `round <length>` for rounded vertices (CSS Shapes L2 addition, widely supported 2026); vertex list comma-separated pairs |
| `path()` | SVG path data string; `<fill-rule>?` optional before comma |
| `shape()` | CSS Shapes L2 addition; not in CSS Masking L1 spec text (2021); broadly supported 2026 |

Spec note: The 2021 spec text references `<basic-shape>` as defined in CSS Shapes Module L1 (`inset`, `circle`, `ellipse`, `polygon`, `path`). `shape()` and `polygon() round` are CSS Shapes L2 additions — see Discrepancies §2.

---

## `clip-rule`

```
Name:    clip-rule
Value:   "nonzero" | "evenodd"
Initial: nonzero
Applies: SVG graphics elements contained within a clipPath element
Inherited: yes
Computed: as specified
Animation: discrete
```

Full closed keyword set:
- `"nonzero"` — uses the nonzero winding rule
- `"evenodd"` — uses the even-odd winding rule

Constraints:
- Only meaningful on graphics elements *inside* a `clipPath`; has no effect on the referencing element.
- Does not apply to `<basic-shape>` functions used on `clip-path` (spec: "Note: The clip-rule property does not apply to `<basic-shape>`s").
- The effective `clip-rule` is the value on the element *defining* the shape (child of `clipPath`), not on the element *using* the `clip-path` property.

Verification: Universally implemented. The semantics match `fill-rule` from SVG 1.1.

---

## `clip` (deprecated)

```
Name:    clip
Value:   <rect()> | "auto"
Initial: auto
Applies: absolutely-positioned elements; in SVG: elements establishing a new viewport,
         pattern elements, mask elements
Inherited: no
Computed: as specified
Animation: by computed value
```

### Value grammar

```ebnf
clip-value = <rect()> | "auto"

<rect()> = "rect(" <edge> "," <edge> "," <edge> "," <edge> ")"

<edge>   = <length> | "auto"
```

Argument order: `top`, `right`, `bottom`, `left` (all from top-left border edge).

`auto` edge = same as the corresponding border-box edge (0 for top/left; computed height+padding+border for bottom; computed width+padding+border for right).

Comma separation is required for UA interop (UA may additionally accept space-separation, but not mixed).

Status: **Deprecated** — authors must use `clip-path` instead. UAs must continue to support it. Widely implemented but marked as "may cease to work at any time" on MDN (2026). No plan for removal in major engines as of 2026; kept for legacy SVG compatibility.

---

## `mask` shorthand + longhands

### `mask` shorthand

```
Name:    mask
Value:   <mask-layer>#
Initial: see individual properties
Applies: all elements; in SVG: container elements (excl. defs), all graphics elements, use
Inherited: no
Animation: see individual properties
```

```ebnf
<mask-layer>
  = <mask-reference>?
    [ <position> [ "/" <bg-size> ]? ]?
    <repeat-style>?
    <geometry-box>?
    [ <geometry-box> | "no-clip" ]?
    <compositing-operator>?
    <masking-mode>?
```

Note: All components are optional and can appear in any order (W3C `||` grammar), with one ordering constraint — `<position>` must come before `<bg-size>` (separated by `/`).

`<geometry-box>` appears twice in the grammar:
- First occurrence sets `mask-origin`
- Second occurrence (or `no-clip`) sets `mask-clip`
- If only one `<geometry-box>` and no `no-clip`: sets *both* `mask-origin` and `mask-clip`

The `mask` shorthand also resets `mask-border` to its initial value.

### `mask-image`

```
Name:    mask-image
Value:   <mask-reference>#
Initial: none
```

```ebnf
<mask-reference>
  = "none"
  | <image>
  | <mask-source>

<mask-source> = <url>
```

`<url>` should reference a `mask` element (constraint, not in formal grammar). Non-`mask` references count as transparent black. `<image>` includes CSS gradient functions and image URLs.

### `mask-mode`

```
Name:    mask-mode
Value:   <masking-mode>#
Initial: match-source
```

```ebnf
<masking-mode> = "alpha" | "luminance" | "match-source"
```

Full closed keyword set:
- `"alpha"` — use alpha channel of mask layer image
- `"luminance"` — use luminance of mask layer image
- `"match-source"` — if `<mask-source>`: defer to referenced `mask` element's `mask-type`; if `<image>`: use alpha

`mask-mode` overrides `mask-type` on the referenced `mask` element when set to a non-`match-source` value.

### `mask-repeat`

```
Name:    mask-repeat
Value:   <repeat-style>#
Initial: repeat
```

`<repeat-style>` is an open datatype — defers to CSS Backgrounds and Borders (`background-repeat` values: `repeat`, `repeat-x`, `repeat-y`, `space`, `round`, `no-repeat`; 1–2 keyword form).

### `mask-position`

```
Name:    mask-position
Value:   <position>#
Initial: 0% 0%
```

`<position>` is an open datatype — defers to CSS Backgrounds and Borders `background-position` syntax (keyword + length/percentage 1–4 value form).

### `mask-clip`

```
Name:    mask-clip
Value:   [ <geometry-box> | "no-clip" ]#
Initial: border-box
```

```ebnf
mask-clip-value = <geometry-box> | "no-clip"
```

Full closed keyword set (8 values):
`"content-box"`, `"padding-box"`, `"border-box"`, `"fill-box"`, `"stroke-box"`, `"view-box"`, `"no-clip"`

(Note: `margin-box` was removed from `mask-clip` in a post-2014 revision — see Changes section of spec and Discrepancies §3.)

Has no effect on mask layer images that reference a `mask` element (those are governed by x/y/width/height on `mask`).

### `mask-origin`

```
Name:    mask-origin
Value:   <geometry-box>#
Initial: border-box
```

Full closed keyword set (6 values — `no-clip` is NOT valid here, `margin-box` also removed):
`"content-box"`, `"padding-box"`, `"border-box"`, `"fill-box"`, `"stroke-box"`, `"view-box"`

Mapping rules for SVG elements without CSS layout box:
- `content-box`, `padding-box`, `border-box` → `fill-box`

Mapping rules for elements with CSS layout box:
- `fill-box`, `stroke-box`, `view-box` → initial value of `mask-origin` (`border-box`)

### `mask-size`

```
Name:    mask-size
Value:   <bg-size>#
Initial: auto
```

`<bg-size>` is an open datatype — defers to `background-size` syntax (`auto`, `cover`, `contain`, `<length-percentage>`, `<length-percentage> <length-percentage>`).

### `mask-composite`

```
Name:    mask-composite
Value:   <compositing-operator>#
Initial: add
```

```ebnf
<compositing-operator> = "add" | "subtract" | "intersect" | "exclude"
```

Full closed keyword set:
- `"add"` — source over destination (Porter-Duff source-over)
- `"subtract"` — source outside destination (Porter-Duff source-out)
- `"intersect"` — source where overlapping destination (Porter-Duff source-in)
- `"exclude"` — non-overlapping regions (Porter-Duff XOR)

The compositing operator for the bottom-most layer is ignored (no layer below it).

---

## `mask-type`

```
Name:    mask-type
Value:   "luminance" | "alpha"
Initial: luminance
Applies: mask elements only
Inherited: no
Computed: as specified
Animation: discrete
```

Full closed keyword set:
- `"luminance"` — mask content treated as luminance mask (default)
- `"alpha"` — mask content treated as alpha mask

Relationship to `mask-mode`: `mask-type` is the author's preference on the `mask` element; `mask-mode` on the masked element can override it (when set to something other than `match-source`). CSS property takes priority over the presentation attribute form.

Verification: Universally implemented. Initial value `luminance` matches MDN and browsers.

---

## Open datatypes used

| Datatype | Where used | Expansion source |
|----------|-----------|-----------------|
| `<url>` | `<clip-source>`, `<mask-source>` | CSS Values L3 `url()` function |
| `<image>` | `mask-image`, `mask-border-source` | CSS Images L3 — gradients, `url()`, `image()`, `image-set()`, `cross-fade()`, `element()` |
| `<position>` | `mask-position`, `circle()`, `ellipse()`, `shape()` | CSS Backgrounds & Borders — keyword + length/percentage 1–4 value form |
| `<bg-size>` | `mask-size` | CSS Backgrounds & Borders — `auto`, `cover`, `contain`, `<length-percentage>{1,2}` |
| `<repeat-style>` | `mask-repeat` | CSS Backgrounds & Borders — `repeat`, `repeat-x`, `repeat-y`, `space`, `round`, `no-repeat` |
| `<length-percentage>` | `inset()`, `circle()`, `ellipse()`, `polygon()`, mask x/y/width/height | CSS Values L3 |
| `<length>` | `clip rect()` edges, `polygon() round` | CSS Values L3 |
| `<border-radius>` | `inset() round` | CSS Backgrounds & Borders |
| `<string>` | `path()` | CSS Values L3 — SVG path data |
| `<shape-command>` | `shape()` | CSS Shapes L2 |

---

## Discrepancies, doc gaps & roadblocks

### 1. `clipPath` content model — `use` indirect reference

**Spec says**: `use` inside `clipPath` must directly reference `path`, `text`, or basic shapes; indirect reference is an error — UA must ignore the `clipPath`.
**Browser reality**: Firefox disables rendering of elements referencing the invalid `clipPath`; no browser actually ignores the `clipPath` on `use` with indirect reference (open CSSWG issue #17 as of spec publication).
**Grammar decision**: Grammar encodes the `use` as permitted child; the "direct reference only" rule goes in the constraint overlay, not the EBNF.

### 2. `<basic-shape>` — `shape()` and `polygon() round` not in spec text

**Spec (2021)**: References CSS Shapes L1, which defines: `inset()`, `circle()`, `ellipse()`, `polygon()`, `path()`. No `shape()` function and no `round` on `polygon()`.
**MDN / browsers (2026)**: Both `shape()` and `polygon(..., round <length>, ...)` are broadly supported (CSS Shapes L2). `shape()` allows line/curve/arc commands.
**Grammar decision**: Include all 6 functions (`inset`, `circle`, `ellipse`, `polygon`, `path`, `shape`) in the grammar since we are targeting current browser behavior. Mark `shape()` and `polygon() round` as CSS Shapes L2 extensions not present in CSS Masking L1 spec text.

### 3. `margin-box` removed from `mask-clip` and `mask-origin`

**Spec (2021) says**: `<geometry-box>` on `clip-path` includes `margin-box` (via `<shape-box>`), but `mask-clip` and `mask-origin` had `margin-box` **removed** in a post-2014 revision.
**Grammar decision**: `clip-path` `<geometry-box>` retains `margin-box`; `mask-clip` and `mask-origin` do not include `margin-box`. Keep these as separate grammar productions even though both use the `<geometry-box>` name.

### 4. `mask-clip` `no-clip` keyword not valid for `mask-origin`

These look superficially similar but `mask-origin` does not accept `no-clip`. `mask-clip` alone accepts it. The `||` shorthand grammar for `mask` can appear to allow `no-clip` anywhere — it only resolves to `mask-clip`. Grammar must model these as distinct alternatives.

### 5. `clip` (deprecated) — separation ambiguity

**Spec says**: Authors should use commas; UAs must support commas; UAs may support space-separation (but not mixed). Formal grammar only shows commas.
**Grammar decision**: Encode comma-only form (`rect( <edge> "," <edge> "," <edge> "," <edge> )`). Note that space-separated form exists for legacy tolerance but is not normative.

### 6. `mask-mode` value `auto` in spec code example

In §9.2 the spec code example uses `mask-mode: auto` (line 1996). `auto` is **not** a valid value per the normative grammar — only `alpha | luminance | match-source`. The example appears to be erroneous; MDN confirms the valid set is `alpha | luminance | match-source`.
**Grammar decision**: Do not include `auto` as a `mask-mode` value. Record as a spec code-example error.

### 7. `maskUnits` default conflicts

**Spec §9.1 says**: "If attribute maskUnits is not specified, then the effect is as if a value of `objectBoundingBox` were specified."
**MDN confirms**: default is `objectBoundingBox`.
**Note**: This differs from `clipPathUnits` whose default is `userSpaceOnUse`. Both are correct per spec; easy source of author confusion. Grammar encodes both defaults as annotations only.

### 8. `mask-composite` — `-webkit-mask-composite` discrepancy

The standard `mask-composite` uses `add | subtract | intersect | exclude`. The `-webkit-mask-composite` (legacy Safari/WebKit) uses a different keyword set (`source-over`, `source-in`, `source-out`, `source-atop`, `destination-over`, etc.) that does not match the standard. These are **not** interchangeable.
**Grammar decision**: Grammar encodes only the standard keywords. The webkit variant is out of scope; note exists for constraint overlay.

### 9. `clip-rule` SVG vs CSS context

`clip-rule` as a CSS property applies to SVG graphics elements inside a `clipPath`. As an SVG presentation attribute it can appear on any graphics element but only takes effect inside `clipPath`. The "applies to" field is context-sensitive; EBNF models the value set only (`nonzero | evenodd`).

### 10. `shape()` `<shape-command>` — not fully defined here

`shape()` references `<shape-command>` from CSS Shapes L2. That sub-grammar (line, hline, vline, curve, smooth, arc, close) needs to be picked up from the CSS Shapes L2 spec, not this file. Flag for the shapes grammar module.
