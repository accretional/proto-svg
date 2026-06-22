# Paint Servers Grammar Notes

**Source:** `/docs/specs/cache/svg2-pservers.txt` (W3C SVG 2, Chapter 14)
**Cross-checked:** MDN Web Docs (June 2026), browser compat tables
**Output target:** EBNF grammar module for paint-server elements

---

## Elements

---

### `linearGradient`

#### Content model

Any number of the following elements, in any order:
- **Descriptive:** `desc`, `title`, `metadata`
- **Animation/script:** `animate`, `animateTransform`, `set`, `script`, `style`
- **Paint child:** `stop` (the primary functional child; non-descriptive stops trigger shadow-tree templating if absent)

Not in content model: `animateMotion` (present in `pattern` but not gradients per spec).

#### Attributes (element-specific)

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `gradientUnits` | `"userSpaceOnUse" \| "objectBoundingBox"` | `objectBoundingBox` | yes | Closed keyword set; governs coord space for x1/y1/x2/y2 |
| `gradientTransform` | `<transform-list>` | identity | yes | Presentation attribute for CSS `transform`; open domain |
| `x1` | `<length-percentage> \| <number>` | `0%` | yes | See mixed-domain note below |
| `y1` | `<length-percentage> \| <number>` | `0%` | yes | See mixed-domain note below |
| `x2` | `<length-percentage> \| <number>` | `100%` | yes | See mixed-domain note below |
| `y2` | `<length-percentage> \| <number>` | `0%` | yes | See mixed-domain note below |
| `spreadMethod` | `"pad" \| "reflect" \| "repeat"` | `pad` | yes | Closed keyword set |
| `href` | `<url>` | empty | yes | Template ref; target must be `linearGradient` or `radialGradient` |
| `xlink:href` | `<url>` | empty | yes | **Deprecated** in SVG 2; fallback only; ignored when `href` is present |

#### Shared attribute groups

- Core attributes: `id`, `tabindex`, `lang`, `xml:space`, `class`, `style`
- Global event attributes (full set listed in spec, omitted here)
- Document element event attributes: `oncopy`, `oncut`, `onpaste`
- Presentation attributes (full inherited CSS property set)
- Deprecated XLink attributes: `xlink:href`, `xlink:title`

#### Attributes copied from `href` template

`x1`, `y1`, `x2`, `y2`, `gradientTransform`, `gradientUnits`, `spreadMethod`
Child `stop` elements are shadow-cloned from the template only when the referencing element has no non-descriptive children.

#### Context-sensitive constraints

1. **gradientUnits = objectBoundingBox**: bare `<number>` values for x1/y1/x2/y2 are interpreted as fractions of the bounding box (0..1 range). Percentages also relative to bounding box. This is the default.
2. **gradientUnits = userSpaceOnUse**: values are in the current user coordinate system; percentages relative to SVG viewport.
3. **Degenerate gradient**: if x1 == x2 AND y1 == y2, the entire painted area is filled with the color of the last gradient stop (solid color).
4. **gradientTransform** is applied after the gradientUnits coordinate mapping; post-multiplied to existing transform stack.
5. **Template href cross-references** may be indirect (template references another template). Resolution is recursive until a defined value or initial value is reached. External file references are allowed but subject to security limits.

---

### `radialGradient`

#### Content model

Identical to `linearGradient`:
- **Descriptive:** `desc`, `title`, `metadata`
- **Animation/script:** `animate`, `animateTransform`, `set`, `script`, `style`
- **Paint child:** `stop`

#### Attributes (element-specific)

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `gradientUnits` | `"userSpaceOnUse" \| "objectBoundingBox"` | `objectBoundingBox` | yes | Governs coord space for cx/cy/r/fx/fy/fr |
| `gradientTransform` | `<transform-list>` | identity | yes | Presentation attribute |
| `cx` | `<length-percentage> \| <number>` | `50%` | yes | Center x of outer (end) circle |
| `cy` | `<length-percentage> \| <number>` | `50%` | yes | Center y of outer (end) circle |
| `r` | `<length-percentage> \| <number>` | `50%` | yes | Radius of outer circle; negative = error |
| `fx` | `<length-percentage> \| <number>` | see cx (dynamic) | yes | Focal point x; defaults to cx's presentational value if unset |
| `fy` | `<length-percentage> \| <number>` | see cy (dynamic) | yes | Focal point y; defaults to cy's presentational value if unset |
| `fr` | `<length-percentage> \| <number>` | `0%` | yes | **SVG 2 new.** Focal circle radius; negative = error; 0% if unspecified |
| `spreadMethod` | `"pad" \| "reflect" \| "repeat"` | `pad` | yes | Same semantics as linearGradient |
| `href` | `<url>` | empty | yes | Target must be `linearGradient` or different `radialGradient` |
| `xlink:href` | `<url>` | empty | yes | **Deprecated** in SVG 2 |

#### Attributes copied from `href` template

`cx`, `cy`, `r`, `fx`, `fy`, `fr`, `gradientTransform`, `gradientUnits`, `spreadMethod`

#### Context-sensitive constraints

1. **gradientUnits** has identical coordinate-space effects as for `linearGradient`, applied to cx/cy/r/fx/fy/fr.
2. **fx/fy defaults**: if not specified on element or any template in the chain, fx=cx and fy=cy at computed time (not parse time). Grammar must mark these as having a computed-default, not a literal default.
3. **fr negative**: error (ignored/treated as 0 per error processing rules).
4. **r negative**: error.
5. **Focal outside end circle**: SVG 2 explicitly allows this (changed from SVG 1.1). Forms a cone; areas outside the cone are transparent black.
6. **Focal on circle edge with spreadMethod=repeat**: painted area gets weighted average color of all stops (not transparent black — differs from Canvas).
7. **gradientTransform** post-multiply semantics identical to linearGradient.

---

### `stop`

#### Content model

Any number of the following elements, in any order:
- `animate`, `set`, `script`, `style`

Note: descriptive elements (`desc`, `title`, `metadata`) are NOT listed in the `stop` content model in the spec (contrast with gradient and pattern elements). No `animateTransform` either.

#### Attributes (element-specific)

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `offset` | `<number> \| <percentage>` | `0` | yes | Range 0..1 as number; 0%..100% as percentage; clamped, monotonic |

#### Properties (CSS, apply only to `stop`)

| Property | Value syntax | Initial | Animatable | Notes |
|---|---|---|---|---|
| `stop-color` | `currentColor \| <color>` | `black` | yes | `<color>` is any CSS color; see discrepancies for `none` and icc-color |
| `stop-opacity` | `<alpha-value>` | `1` | yes | `<number>` 0..1 or `<percentage>` 0%..100%; clamped |

#### Context-sensitive constraints

1. **offset clamping**: values < 0 (or < 0%) rounded up to 0; values > 1 (or > 100%) rounded down to 100%.
2. **Monotonic enforcement**: if a stop's offset is less than the previous stop's offset, it is adjusted up to equal the previous stop's offset. Grammar note: offset is parsed as a raw value; the monotonic constraint is a semantic (overlay) rule, not a syntactic one.
3. **Duplicate offsets**: the later stop's color wins at the overlap point.
4. **Minimum stops**: at least 2 stops needed for a gradient effect. 0 stops → paint as `none`. 1 stop → solid fill with that stop's color.
5. **stop-opacity × stop-color alpha**: final alpha = stop-opacity × alpha channel of stop-color; stop-color types without explicit alpha are treated as fully opaque (alpha=1) before multiplication.
6. **`transparent` keyword**: in SVG gradients, `transparent` = black at opacity 0 (not pre-multiplied). This differs from CSS gradient handling.

---

### `pattern`

#### Content model

Any number of the following elements, in any order:
- **Animation elements:** `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- **Descriptive elements:** `desc`, `title`, `metadata`
- **Paint server elements:** `linearGradient`, `radialGradient`, `pattern` (recursive)
- **Shape elements:** `circle`, `ellipse`, `line`, `path`, `polygon`, `polyline`, `rect`
- **Structural elements:** `defs`, `g`, `svg`, `symbol`, `use`
- **Other:** `a`, `audio`, `canvas`, `clipPath`, `filter`, `foreignObject`, `iframe`, `image`, `marker`, `mask`, `script`, `style`, `switch`, `text`, `video`, `view`

This is effectively "any graphics/container/descriptive/animation element" — the most permissive content model of the four elements here.

#### Attributes (element-specific)

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `patternUnits` | `"userSpaceOnUse" \| "objectBoundingBox"` | `objectBoundingBox` | yes | Governs coord space for x/y/width/height |
| `patternContentUnits` | `"userSpaceOnUse" \| "objectBoundingBox"` | `userSpaceOnUse` | yes | Governs coord system for pattern content; ignored when viewBox is set |
| `patternTransform` | `<transform-list>` | identity | yes | Presentation attribute; post-multiplied |
| `x` | `<length>` | `0` | yes | Top-left x of reference rectangle |
| `y` | `<length>` | `0` | yes | Top-left y of reference rectangle |
| `width` | `<length>` | `0` | yes | Width of tile; negative = error; 0 = no paint |
| `height` | `<length>` | `0` | yes | Height of tile; negative = error; 0 = no paint |
| `viewBox` | `<min-x> <min-y> <width> <height>` | none (absent) | yes | Introduces supplemental coord transform; suppresses patternContentUnits |
| `preserveAspectRatio` | `align meetOrSlice?` | `xMidYMid meet` | yes | Only meaningful when viewBox is set |
| `href` | `<url>` | empty | yes | Target must be different `pattern` element |
| `xlink:href` | `<url>` | empty | yes | **Deprecated** in SVG 2 |

Where `align` is one of: `none | xMinYMin | xMidYMin | xMaxYMin | xMinYMid | xMidYMid | xMaxYMid | xMinYMax | xMidYMax | xMaxYMax`
And `meetOrSlice` is: `meet | slice`

#### Attributes copied from `href` template

`x`, `y`, `width`, `height`, `viewBox`, `preserveAspectRatio`, `patternTransform`, `patternUnits`, `patternContentUnits`
Child content is shadow-cloned from template only when current element has no non-descriptive children.

#### Context-sensitive constraints

1. **patternUnits = objectBoundingBox** (default): x/y/width/height are fractions of referencing element's bounding box.
2. **patternUnits = userSpaceOnUse**: x/y/width/height are in the current user coordinate system at reference time.
3. **patternContentUnits** controls the coordinate system for pattern content. Has no effect when `viewBox` is set — viewBox replaces it.
4. **viewBox interaction**: when present, viewBox maps its coordinate space into the tile defined by x/y/width/height using standard viewBox+preserveAspectRatio rules.
5. **overflow**: UA stylesheet sets `overflow: hidden` on pattern; content outside tile is clipped. Behavior with `overflow: visible` is undefined in SVG 2.
6. **Event attributes on pattern content** are not processed; only rendering is processed.
7. **Template href cross-references** may be indirect, subject to security limits. Pattern `href` can only target another `pattern` (unlike gradient `href` which allows cross-type).

---

## Open Datatypes Used

| Leaf name | Description | Where used |
|---|---|---|
| `<transform-list>` | One or more CSS/SVG transform functions | `gradientTransform`, `patternTransform` |
| `<length>` | CSS length with unit (e.g. `10px`, `5em`) | pattern x/y/width/height |
| `<length-percentage>` | CSS `<length>` or `<percentage>` value | linearGradient x1/y1/x2/y2; radialGradient cx/cy/r/fx/fy/fr |
| `<number>` | Unitless numeric literal | gradient coords (interpreted as fraction in OBB mode); stop offset |
| `<percentage>` | Numeric literal with `%` suffix | stop offset; gradient coords |
| `<color>` | Any CSS `<color>` (named, rgb(), hsl(), hex, etc.) | stop-color |
| `<alpha-value>` | `<number>` 0..1 or `<percentage>` 0%..100% | stop-opacity |
| `<url>` | CSS `url(...)` or bare IRI | href on all four elements |
| `viewBox-value` | `<number> <number> <number> <number>` (min-x min-y width height) | pattern viewBox |

### Mixed-domain splits required for grammar

**Gradient coordinate attrs (x1/y1/x2/y2/cx/cy/r/fx/fy/fr):**
```
gradient-coord ::= <length-percentage> | <number>
```
The bare `<number>` alternative is a legacy/OBB shorthand. In `userSpaceOnUse` mode a bare number is unitless user-space. In `objectBoundingBox` mode it is a [0,1] fraction. The grammar admits both; the coordinate-space semantics are a context overlay rule.

**stop offset:**
```
stop-offset ::= <number> | <percentage>
```
`<number>` here ranges 0..1 (not the same as a CSS `<number>`; it's semantically fractional). Grammar should capture this as the same token but annotate the constraint separately.

**stop-color:**
```
stop-color-value ::= currentColor | <color>
```
`<color>` is an open CSS datatype. See discrepancies for `none` exclusion and `<icccolor>` obsolescence.

---

## Discrepancies, Doc Gaps & Roadblocks

### D1 — Spec value type for x1/y1/x2/y2 is `<length>`, not `<length-percentage> | <number>`

**Spec text (line 442):** lists value as `<length>` (no mention of `<number>` or `<percentage>` as separate).
**Spec prose (line 438):** "The values of 'x1', 'y1', 'x2' and 'y2' can be either numbers or percentages."
**MDN (2026):** lists value as `<length-percentage> | <number>` for x1.
**DOM IDL (line 1782):** `SVGAnimatedLength` — SVGLength includes both length-with-unit and unitless number.
**Grammar decision:** Use `<length-percentage> | <number>` (MDN/prose/IDL all agree; the spec's `<length>` table entry is a simplification error). The bare `<number>` form is widely used in practice (e.g. `x1="0" x2="1"` in objectBoundingBox mode).

### D2 — stop-color: spec includes `<icccolor>` (SVG 1.1 syntax)

**Spec text (line 1156):** `currentColor | <color> <icccolor>`
**CSS Color 4 / SVG 2 reality:** `<icccolor>` (`icc-color(...)` notation) is from SVG 1.1 and is removed/not implemented. The `SVGColor` interface is deprecated.
**MDN (2026):** stop-color accepts `currentColor | <color>` (any CSS color); no mention of `<icccolor>`.
**Browser behavior:** No major browser implements `icc-color(...)` on stop-color.
**Grammar decision:** Omit `<icccolor>`. Use `currentColor | <color>`. Flag this in overlay as "SVG 1.1 legacy syntax, not implemented".

### D3 — stop-color: "none" keyword

**Spec text:** does not list `none` as a valid stop-color value (unlike `fill`/`stroke` which accept `none`).
**MDN (2026):** does not list `none` as valid for stop-color. Syntax is `<color>`.
**Grammar decision:** `none` is NOT a valid stop-color value. The element should be given stop-opacity=0 or stop-color=transparent to achieve transparency. Do not enumerate `none` as a terminal.

### D4 — fr (focal radius) browser support

**Spec text (line 866):** "New in SVG 2. Added to align with Canvas."
**MDN (2026):** fr is documented with full browser support.
**Browser compat (from MDN/caniuse):** Available across browsers since January 2020 (Chrome, Firefox, Safari all support).
**Grammar decision:** `fr` is safe to include as a standard attribute. No special caveat needed beyond the SVG 2 minimum version note.

### D5 — href vs xlink:href — both present

**Spec text:** both `href` and `xlink:href` listed in attributes for all four elements; spec calls `xlink:href` "deprecated".
**MDN (2026):** `xlink:href` is deprecated in SVG 2, may already be removed from spec; `href` takes precedence; browsers ignore `xlink:href` when `href` is set.
**Grammar decision:** Grammar should define `href` as the canonical attribute. `xlink:href` modeled as a deprecated-alias attribute with a note that it produces identical semantics. Do not drop `xlink:href` from grammar since it still appears in real SVG documents.

### D6 — fx/fy defaults are computed, not static

**Spec text (lines 828, 858):** "if attribute fx is not specified, fx will coincide with the presentational value of cx." This is a run-time/computed default, not a parse-time literal.
**Grammar implication:** fx and fy default values cannot be written as grammar literals. Mark them `optional` in the grammar with an overlay note: "default = computed cx / cy at reference time". This also applies when `href` templating is involved.

### D7 — Pattern x/y/width/height type is `<length>` but NOT `<length-percentage>`

**Spec text (line 1493):** `<length>` (no percentage column for x/y/width/height on `pattern`).
**MDN pattern element page:** lists these as `<length>`, not `<length-percentage>`.
**Contrast:** gradient coordinate attrs DO accept `<percentage>`. Pattern tile dimensions do NOT formally accept percentage values.
**Grammar decision:** Pattern x/y/width/height = `<length>` only (no percentage). This is a genuine difference from gradient coordinates.

### D8 — patternContentUnits suppressed by viewBox

**Spec text (line 1415):** "Note that this attribute has no effect if attribute 'viewBox' is specified."
**Grammar/overlay implication:** `patternContentUnits` is syntactically valid alongside `viewBox`, but semantically inactive when both are present. Record as a context-sensitive constraint in the overlay; not a grammar-level restriction.

### D9 — stop content model excludes animateTransform and descriptive elements

**Spec text (lines 1086-1087):** `stop` content = `animate, script, set, style` only. No `animateTransform`, no `desc/title/metadata`.
**This differs from gradient parent elements** which include descriptive elements and `animateTransform`.
**Grammar action:** `stop` content model must be more restrictive than the gradient parent models. Enumerate explicitly.

### D10 — gradientTransform / patternTransform are presentation attributes for CSS `transform`

**Spec text (line 334):** "The 'gradientTransform' attribute is a presentation attribute for the transform property."
**Grammar implication:** `gradientTransform` and `patternTransform` are simultaneously SVG attributes AND CSS presentation attributes. Their value domain is `<transform-list>` (the SVG/CSS shared transform syntax). This does not change their EBNF value syntax but matters for the property mapping overlay.

---

*End of paint servers grammar notes.*
