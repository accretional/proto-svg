# Embedded / Linking / Interactivity Grammar Notes

## Source

- `svg2-embedded.txt` — W3C SVG 2 Chapter 12: Embedded Content (`image`, `foreignObject`)
- `svg2-linking.txt`  — W3C SVG 2 Chapter 16: Linking (`a`, `view`, href/xlink family)
- `svg2-interact.txt` — W3C SVG 2 Chapter 15: Scripting and Interactivity (`script`, `pointer-events`, `zoomAndPan`, event attrs)

---

## Elements

### `image`

#### Categories
Graphics element, graphics referencing element, renderable element, structurally external element.

#### Content model
Zero or more of, in any order:
- animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- descriptive elements: `desc`, `title`, `metadata`
- `clipPath`, `mask`, `script`, `style`

No text content. No child elements from non-SVG namespaces allowed directly.

#### Attributes with value syntax and defaults

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `href` | `<url>` | (none) | yes | Preferred; replaces `xlink:href` |
| `xlink:href` | `<url>` | (none) | conditional | Deprecated; ignored when bare `href` present |
| `xlink:title` | `<anything>` | (none) | no | Deprecated; use `<title>` child instead |
| `crossorigin` | `[ anonymous \| use-credentials ]?` | (HTML default — treated as "no-cors" when absent) | yes | CORS settings attribute |
| `preserveAspectRatio` | (see shared groups) | `xMidYMid meet` | yes | |
| `x` | `<length-percentage>` | `0` | yes | geometry property / presentation attribute |
| `y` | `<length-percentage>` | `0` | yes | geometry property / presentation attribute |
| `width` | `<length-percentage> \| auto` | `auto` | yes | geometry property; `0` disables rendering |
| `height` | `<length-percentage> \| auto` | `auto` | yes | geometry property; `0` disables rendering |

Plus: aria attributes, core attributes, conditional processing attributes, global event attributes, document element event attributes, graphical event attributes, presentation attributes (see shared groups).

#### Context-sensitive constraints
- `href` must resolve to a complete image document: PNG, JPEG, SVG, or other image MIME type. It cannot reference an element fragment within an SVG (unlike `use`).
- Fragment-only URLs (e.g. `#id`) reference the same document and are re-rendered as an independent embedded image in secure mode (recursion capped at one level).
- CORS state for fetch: taken from `crossorigin` attribute.
- Embedded SVG is processed in secure animated or secure static mode; scripts and external references inside it are suppressed.
- The `preserveAspectRatio` attribute on the root element of a referenced SVG is forcibly ignored; the outer `image`'s `preserveAspectRatio` governs.
- `overflow: hidden` set by user-agent stylesheet; images clip to positioning rectangle unless overridden.

#### EBNF sketch
```ebnf
image-element ::= '<image' image-attrs* '>'
                   image-content*
                  '</image>'

image-attrs    ::= common-attrs
                 | href-attr
                 | xlink-href-attr       (* deprecated *)
                 | xlink-title-attr      (* deprecated *)
                 | crossorigin-attr
                 | preserveAspectRatio-attr
                 | geometry-x-attr
                 | geometry-y-attr
                 | geometry-width-attr
                 | geometry-height-attr

image-content  ::= animation-element
                 | descriptive-element
                 | 'clipPath' | 'mask' | 'script' | 'style'
```

---

### `foreignObject`

#### Categories
Graphics element, renderable element, structurally external element.

#### Content model
Any elements or character data — completely open. In practice: HTML namespace content (the HTML parser treats children as if inside an HTML document fragment). SVG-namespaced elements inside `foreignObject` are NOT rendered unless they form a complete SVG subtree rooted at an `<svg>` element.

#### Attributes with value syntax and defaults

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `x` | `<length-percentage>` | `0` | yes | geometry property |
| `y` | `<length-percentage>` | `0` | yes | geometry property |
| `width` | `<length-percentage> \| auto` | `auto` | yes | geometry property; `0` disables rendering |
| `height` | `<length-percentage> \| auto` | `auto` | yes | geometry property; `0` disables rendering |

Plus: aria attributes, core attributes, conditional processing attributes, global event attributes, document element event attributes, graphical event attributes, presentation attributes. No `href`, no `crossorigin`, no `preserveAspectRatio`.

**No `href`-class attributes on `foreignObject`.** The spec lists no URL reference attributes for this element.

#### Context-sensitive constraints
- The positioning rectangle (x, y, width, height) defines a CSS containing block; children are laid out with CSS.
- `foreignObject` is implicitly absolutely-positioned for CSS purposes.
- Commonly used with `<switch>` + `requiredExtensions` for fallback rendering.
- Child HTML content is subject to SVG transforms, filters, clipping, masking, and compositing.

#### EBNF sketch
```ebnf
foreignObject-element ::= '<foreignObject' foreignObject-attrs* '>'
                           any-content*
                          '</foreignObject>'

foreignObject-attrs   ::= common-attrs
                         | geometry-x-attr
                         | geometry-y-attr
                         | geometry-width-attr
                         | geometry-height-attr

any-content ::= (* any element from any namespace, or character data *)
```

---

### `a`

#### Categories
Container element, renderable element.

#### Content model
Descriptive content, plus any element or text allowed by its **parent's** content model — EXCEPT another `a` element. If the parent is a `switch`, use the nearest non-`switch` ancestor's content model. Nesting `a` inside `a` is invalid; nested `a` elements are treated as inactive (rendered as generic containers).

#### Attributes with value syntax and defaults

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `href` | `<url>` | (none) | yes | Activates the link |
| `xlink:href` | `<url>` | (none) | conditional | Deprecated |
| `xlink:title` | `<anything>` | (none) | no | Deprecated |
| `target` | `_self \| _parent \| _top \| _blank \| <XML-Name>` | `_self` | yes | See keyword semantics below |
| `download` | `<string>` (any value; non-empty = suggested filename) | (none) | no | |
| `ping` | space-separated `<url>` tokens | (none) | no | |
| `rel` | space-separated keyword tokens | (none) | no | |
| `hreflang` | BCP 47 language tag `<string>` | (none) | no | |
| `type` | MIME type `<string>` | (none) | no | |
| `referrerpolicy` | referrer policy `<string>` | (none) | no | IDL attribute: `referrerPolicy` (camelCase) |

Plus: aria attributes, core attributes, conditional processing attributes, global event attributes, document element event attributes, graphical event attributes, presentation attributes.

#### `target` keyword semantics (closed set + open name)
```ebnf
target-value ::= '_self'            (* replace current browsing context — DEFAULT *)
               | '_parent'          (* replace immediate parent browsing context *)
               | '_top'             (* replace full active window/tab *)
               | '_blank'           (* open in new window/tab *)
               | <XML-Name>         (* named browsing context; must NOT start with '_' per HTML *)
```

**Deprecated and removed:** `_replace` (SVG 1.1, never well-implemented; equivalent to `_self` in modern HTML); `_new` is NOT legal (use `_blank` instead).

#### Context-sensitive constraints
- Without `href` or `xlink:href`, `a` is an inactive placeholder; rendered as a generic container.
- An `a` descendant of another hyperlink element (SVG or other namespace) must have its `href` ignored.
- Hit-testing is per-contained-element (using each child's `pointer-events`), not bounding-box of the `a` element.
- All `a` elements that are valid links must be focusable and keyboard-activatable.
- `download`, `ping`, `rel`, `hreflang`, `type`, `referrerpolicy` semantics are defined by HTML by reference.

#### EBNF sketch
```ebnf
a-element   ::= '<a' a-attrs* '>'
                 a-content*
                '</a>'

a-attrs     ::= common-attrs
              | href-attr
              | xlink-href-attr       (* deprecated *)
              | xlink-title-attr      (* deprecated *)
              | target-attr
              | download-attr
              | ping-attr
              | rel-attr
              | hreflang-attr
              | type-attr
              | referrerpolicy-attr

target-attr ::= 'target' '=' '"' target-value '"'
target-value ::= '_self' | '_parent' | '_top' | '_blank' | <XML-Name>

a-content   ::= (* parent's content model minus 'a'; see above *)
```

---

### `view`

#### Categories
None.

#### Content model
Zero or more of, in any order:
- animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- descriptive elements: `desc`, `title`, `metadata`
- `script`, `style`

#### Attributes with value syntax and defaults

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `viewBox` | `<number> <number> <number> <number>` (min-x, min-y, width, height) | (none) | yes | |
| `preserveAspectRatio` | (see shared groups) | `xMidYMid meet` | yes | |
| `zoomAndPan` | `disable \| magnify` | `disable` | no | At-risk; no known implementations |

**Removed:** `viewTarget` attribute — explicitly resolved to remove at Paris 2015 F2F Day 3. SVG 1.1 had it; SVG 2 does not.

Plus: aria attributes, core attributes, global event attributes, document element event attributes.

Note: `view` has NO conditional processing attributes listed.

#### Context-sensitive constraints
- A `view` element is referenced by URL fragment: `MyFile.svg#MyViewId`. When targeted, it overrides the root `svg`'s view specification attributes.
- `zoomAndPan` on inner `svg` elements has no effect; only meaningful on the outermost `svg` and on `view` elements used as targets.
- `SVGViewElement` includes `SVGFitToViewBox` (provides `viewBox` + `preserveAspectRatio`) and `SVGZoomAndPan` (provides `zoomAndPan`).

#### EBNF sketch
```ebnf
view-element ::= '<view' view-attrs* '>'
                  view-content*
                 '</view>'

view-attrs   ::= common-attrs
               | viewBox-attr
               | preserveAspectRatio-attr
               | zoomAndPan-attr

zoomAndPan-attr ::= 'zoomAndPan' '=' '"' zoomAndPan-value '"'
zoomAndPan-value ::= 'disable' | 'magnify'

view-content ::= animation-element | descriptive-element | 'script' | 'style'
```

---

### `script`

#### Categories
Never-rendered element, structurally external element.

#### Content model
Character data only. Script text is the inline content; it is never rendered. When `href` is present, inline content is ignored (same model as HTML `<script>`).

#### Attributes with value syntax and defaults

| Attribute | Value syntax | Default | Animatable | Notes |
|---|---|---|---|---|
| `type` | MIME type `<string>` | `application/ecmascript` | no | If absent, ECMAScript is assumed |
| `href` | `<url>` | (none) | no | External script resource |
| `xlink:href` | `<url>` | (none) | no | Deprecated |
| `xlink:title` | `<anything>` | (none) | no | Deprecated |
| `crossorigin` | `[ anonymous \| use-credentials ]?` | (HTML default) | yes | CORS settings attribute |

Plus: core attributes, global event attributes, document element event attributes. No aria, no conditional processing, no presentation attributes, no geometry attributes.

#### Context-sensitive constraints
- `href` must resolve to an external document providing script content.
- CORS for fetch: taken from `crossorigin` attribute (same as HTML `<script>`).
- If `type` media type is not supported, the script must not be executed.
- `display` is always `none` via UA stylesheet with `!important`-equivalent importance; script content is never visible.
- `async` and `defer` attributes are mentioned as a SVG 2 requirement/resolution item but are NOT formally defined in the attribute tables in this spec version. They are not included in the grammar.

#### EBNF sketch
```ebnf
script-element ::= '<script' script-attrs* '>'
                    script-content
                   '</script>'

script-attrs   ::= common-attrs
                 | type-attr
                 | href-attr
                 | xlink-href-attr       (* deprecated *)
                 | xlink-title-attr      (* deprecated *)
                 | crossorigin-attr

script-content ::= (* character data / CDATA *)

type-attr      ::= 'type' '=' '"' <mime-type> '"'
                   (* default: "application/ecmascript" when absent *)
```

---

## href and xlink Family

### Bare `href` (preferred, SVG 2)

```ebnf
href-attr ::= 'href' '=' '"' <url> '"'
```

- Value: URL (IRI per RFC 3987) — absolute or relative, with or without fragment identifier.
- Default: (none) — element is inactive/non-referencing when absent.
- Animatable: yes (on `image`, `a`); no (on `script`); conditional on animation elements.
- When both `href` and `xlink:href` are present, `href` wins; `xlink:href` is ignored.
- A conforming SVG generator MUST output bare `href`. It MAY also emit `xlink:href` for backwards compatibility.

### `xlink:href` (deprecated, SVG 2)

```ebnf
xlink-href-attr ::= 'xlink:href' '=' '"' <url> '"'
                    (* requires xmlns:xlink="http://www.w3.org/1999/xlink" declaration *)
```

- Value: URL (same domain as bare `href`).
- Default: (none).
- Animatable: if and only if the corresponding bare `href` is animatable.
- Ignored when bare `href` is present on the same element.
- Processed at the same time a bare `href` would be processed.
- Still parsed by all major browsers; must appear in grammar as a deprecated alternative.

### `xlink:title` (deprecated, SVG 2)

```ebnf
xlink-title-attr ::= 'xlink:title' '=' '"' <anything> '"'
                     (* requires xmlns:xlink declaration *)
```

- Value: any string.
- Default: (none).
- Animatable: no.
- Replaced by `<title>` child element. Kept for back-compat parsing only.

### Other XLink attributes from SVG 1.1 (not carried into SVG 2 spec text)

The SVG 2 spec text in these chapters defines only `xlink:href` and `xlink:title` as deprecated attributes. The full XLink 1.0 attribute family (`xlink:role`, `xlink:arcrole`, `xlink:show`, `xlink:actuate`, `xlink:type`) was present in SVG 1.1 but is not formally listed as attributes on any SVG 2 element in the chapters read. They should be treated as obsolete/not in SVG 2 grammar scope.

**Verification against MDN/browsers:** `xlink:href` is still parsed in all major browsers (Chrome, Firefox, Safari) as of 2025. `xlink:show`, `xlink:actuate`, `xlink:role`, `xlink:arcrole` are not actively supported in browser SVG engines and should be omitted from the grammar or placed in a `legacy-xlink` overlay.

### SVG fragment identifier grammar (from linking chapter)

```ebnf
SVGFragmentIdentifier ::= BareName *( '&' timesegment )
                        | SVGViewSpec *( '&' timesegment )
                        | spacesegment *( '&' timesegment )
                        | timesegment *( '&' spacesegment )

BareName     ::= XML_Name
SVGViewSpec  ::= 'svgView(' SVGViewAttributes ')'
SVGViewAttributes ::= SVGViewAttribute
                    | SVGViewAttribute ';' SVGViewAttributes

SVGViewAttribute ::= viewBoxSpec
                   | preserveAspectRatioSpec
                   | transformSpec
                   | zoomAndPanSpec

viewBoxSpec              ::= 'viewBox(' ViewBoxParams ')'
preserveAspectRatioSpec  ::= 'preserveAspectRatio(' AspectParams ')'
transformSpec            ::= 'transform(' TransformParams ')'
zoomAndPanSpec           ::= 'zoomAndPan(' ZoomAndPanParams ')'
```

Notes: Each SVGViewAttribute type may appear at most once. Order is free. Commas separate numeric values; semicolons separate attributes. Fragment may be URL-percent-escaped.

### `target` keyword set (on `a`)

Closed set plus open name:
```ebnf
target-value ::= '_self' | '_parent' | '_top' | '_blank' | <XML-Name>
```

`_replace` (SVG 1.1) — removed/deprecated; `_new` — never valid.

### `crossorigin` keyword set

```ebnf
crossorigin-value ::= 'anonymous' | 'use-credentials'
                    (* attribute may be omitted entirely — treated as "no-cors" when absent *)
                    (* empty string treated as "anonymous" per HTML spec *)
```

Present on: `image`, `script`. Absent from `a`, `view`, `foreignObject`.

---

## pointer-events property

### Full property definition

```
Name:    pointer-events
Value:   bounding-box | visiblePainted | visibleFill | visibleStroke | visible
         | painted | fill | stroke | all | none
Initial: visiblePainted
Applies: container elements, graphics elements, 'use'
Inherited: yes
Computed: as specified
Animatable: yes
```

### Keyword terminals (complete closed set)

```ebnf
pointer-events-value ::= 'bounding-box'
                        | 'visiblePainted'
                        | 'visibleFill'
                        | 'visibleStroke'
                        | 'visible'
                        | 'painted'
                        | 'fill'
                        | 'stroke'
                        | 'all'
                        | 'none'
```

### Keyword semantics summary

| Value | visibility required | fill area | stroke area |
|---|---|---|---|
| `bounding-box` | no | bounding box | bounding box |
| `visiblePainted` | yes (visible) | if fill != none | if stroke != none |
| `visibleFill` | yes (visible) | always | no |
| `visibleStroke` | yes (visible) | no | always |
| `visible` | yes (visible) | always | always |
| `painted` | no | if fill != none | if stroke != none |
| `fill` | no | always | no |
| `stroke` | no | no | always |
| `all` | no | always | always |
| `none` | — | never | never |

**Text elements:** hit-testing is character-cell-based. `visibleFill`, `visibleStroke`, `visible` are equivalent for text (all = character cell if visible). `fill`, `stroke`, `all` are equivalent for text (all = character cell, no visibility check).

**Raster images:** `fill`, `stroke`, `all` are equivalent (rectangular area, no opacity check). `visibleFill`, `visibleStroke`, `visible` are equivalent (rectangular area, visibility check). `visiblePainted` / `painted` use per-pixel alpha. `fill-opacity`, `stroke-opacity`, `opacity`, `fill`, `stroke` properties do NOT affect image pointer-event processing.

**CSS keyword note:** The CSS specification also defines `auto` as a value for `pointer-events` on HTML elements; in SVG it defers to `visiblePainted` (the initial). This is NOT listed in the SVG spec keyword set for this property. Grammar should NOT include `auto` in the SVG terminal set but must allow it in a CSS-property context for HTML/mixed documents.

---

## cursor property

The `svg2-interact.txt` source does not define or enumerate the `cursor` property. The SVG 2 interactivity chapter does not contain a `cursor` property definition. The `cursor` property in SVG is inherited from CSS UI. It is outside the scope of these source files and should be handled from the CSS spec source. No grammar terminals can be extracted here.

---

## zoomAndPan attribute

Defined in the interactivity chapter (§15.7) and on the `view` element:

```ebnf
zoomAndPan-value ::= 'disable' | 'magnify'
```

- Default: `disable` (per interact chapter attribute table; spec says "default being magnify" in prose — see discrepancies).
- Animatable: no.
- At-risk; no known implementations as of spec writing.

---

## Open datatypes used

| Name | Description | Where used |
|---|---|---|
| `<url>` | IRI/URL (RFC 3987); absolute or relative with optional fragment | `href`, `xlink:href`, `ping` |
| `<length-percentage>` | CSS `<length>` or `<percentage>` | `x`, `y`, `width`, `height` on `image`, `foreignObject` |
| `<string>` | Arbitrary string | `download`, `hreflang`, `type`, `rel`, `referrerpolicy`, `xlink:title` |
| `<mime-type>` | MIME media type string (RFC 2046) | `type` on `script`, `type` on `a` |
| `<XML-Name>` | Valid XML Name per XML 1.1 | `target` named browsing context on `a` |
| `ViewBoxParams` | Four numbers: min-x, min-y, width, height | `viewBox`, `viewBoxSpec` fragment |
| `AspectParams` | `preserveAspectRatio` parameter string | `preserveAspectRatioSpec` fragment |
| `TransformParams` | CSS transform function string | `transformSpec` fragment |
| `ZoomAndPanParams` | `disable` or `magnify` | `zoomAndPanSpec` fragment |
| `timesegment` | Media Fragments URI temporal segment | SVG fragment identifier |
| `spacesegment` | Media Fragments URI spatial segment | SVG fragment identifier |

---

## Shared attribute groups (referenced above)

These are named groups used across multiple elements; full definitions live in earlier grammar notes.

- **common-attrs** — core attrs (`id`, `tabindex`, `lang`, `xml:space`, `class`, `style`), aria attrs, conditional processing attrs, global event attrs, document element event attrs, graphical event attrs, presentation attrs.
- **geometry-x-attr** — `x` geometry property.
- **geometry-y-attr** — `y` geometry property.
- **geometry-width-attr** — `width` geometry property (`<length-percentage> | auto`).
- **geometry-height-attr** — `height` geometry property (`<length-percentage> | auto`).
- **preserveAspectRatio-attr** — `preserveAspectRatio` value syntax (see 01-datatypes-geometry notes).

---

## Discrepancies, doc gaps & roadblocks

### D1 — `zoomAndPan` default value contradiction
**Spec text says:** Attribute table in §15.7 lists `Initial value: disable`. Prose in §15.7 says "default being magnify". The `view` element spec (§16.3.3) also exposes `zoomAndPan` via `SVGZoomAndPan` mixin but gives no separate default.  
**MDN:** Lists default as `magnify` for `<svg>` element.  
**Real browsers:** Chrome/Firefox treat absence of `zoomAndPan` as magnify-enabled behavior.  
**Decision:** Grammar default = `magnify` (aligns with prose, MDN, and browser reality). The attribute table value `disable` appears to be a typo/editorial error in the spec.

### D2 — xlink full family (role, arcrole, show, actuate, type) status
**Spec text (SVG 2):** Only `xlink:href` and `xlink:title` are listed as deprecated attributes. `xlink:show`, `xlink:actuate`, `xlink:role`, `xlink:arcrole`, `xlink:type` are not mentioned in SVG 2 chapters.  
**MDN:** Does not document these for SVG 2 elements.  
**Browsers:** Chrome/Firefox never implemented `xlink:show`/`xlink:actuate`/`xlink:arcrole`/`xlink:role` for SVG elements.  
**Decision:** These five attributes are OBSOLETE, not deprecated — they were in SVG 1.1 XLink but are not part of SVG 2. They should appear in an `xlink-legacy` overlay grammar rule with a `(* obsolete, not parsed *)` annotation but not in the main grammar.

### D3 — `async` / `defer` on `script` not formally defined
**Spec text:** SVG 2 Requirement section says "SVG 2 will allow async/defer on `<script>`" but no attribute definition table entry is provided.  
**MDN:** Lists `async` and `defer` on SVG `<script>` as valid per HTML definitions.  
**Browsers:** Chrome/Firefox do recognize `async`/`defer` on SVG `<script>`.  
**Decision:** Include `async` and `defer` as open `<string>` (boolean attribute, value = attribute name or empty) in the grammar with a note that they are sourced from HTML, not formally specified in SVG 2 text. Flag as "gap in SVG 2 spec."

### D4 — `referrerpolicy` attribute casing
**Spec attribute table:** Uses `referrerPolicy` (camelCase) in prose and IDL, but in the attribute list it appears as `referrerpolicy` (lowercase per HTML serialization convention).  
**HTML spec:** Attribute in markup is `referrerpolicy` (all lowercase); IDL property is `referrerPolicy`.  
**Decision:** Grammar uses `referrerpolicy` (lowercase, HTML-compatible attribute name) with a note that the IDL property name is `referrerPolicy`.

### D5 — `cursor` property absent from interactivity chapter
**Spec text:** Chapter 15 (scripting/interactivity) does not define `cursor`. A `cursor` property that accepts SVG-specific image URLs was defined in SVG 1.1 Chapter 16. SVG 2 deferred this to CSS UI.  
**MDN:** `cursor` property on SVG elements inherits CSS cursor values + `url()` references.  
**Decision:** `cursor` grammar must be sourced from the CSS UI spec, not from these SVG chapters. Mark as external dependency.

### D6 — `crossorigin` empty-string behavior
**Spec text:** Value syntax given as `[ anonymous | use-credentials ]?` — the `?` means optional (attribute may be absent). Does not explicitly say that an empty-string or bare attribute (boolean presence) maps to `anonymous`.  
**HTML spec:** Empty string or attribute presence without value = `anonymous`.  
**Decision:** Grammar adds the empty-string case: `crossorigin-value ::= '' | 'anonymous' | 'use-credentials'` (where `''` and `'anonymous'` are equivalent per HTML processing).

### D7 — `foreignObject` has no URL attributes at all
**Spec text:** Confirmed — no `href`, no `crossorigin`, no `preserveAspectRatio` on `foreignObject`.  
**Browsers:** Confirmed — `foreignObject` has no resource-referencing mechanism of its own; content comes from child DOM.  
**No discrepancy — confirmed by source.**

### D8 — `_replace` target value removal
**Spec text:** Explicitly notes `_replace` was SVG 1.1 target value, never well-implemented, now redundant. Should not appear in grammar.  
**Decision:** Exclude from terminals. Note in comment: `_replace` is invalid in SVG 2.

### D9 — `pointer-events: auto` for SVG elements
**Spec text:** Does not list `auto` in the SVG `pointer-events` keyword set.  
**CSS spec (css-ui-4):** `auto` is a valid value for `pointer-events` in CSS contexts where it resolves implementation-specifically.  
**Browsers:** Accept `auto` on SVG elements (treated as `visiblePainted`).  
**Decision:** `auto` should appear in the grammar as a separate CSS-context alias: `pointer-events-value ::= ... | 'auto'` with annotation `(* CSS alias; resolves to visiblePainted on SVG elements *)`.

### D10 — `image` href animatability note
**Spec text:** `href` on `image` = yes (animatable). `xlink:href` = "see below" / "animatable if and only if the corresponding href attribute is defined to be animatable."  
**This is consistent** — both are animatable for `image`. But for `script`, `href` = no, so `xlink:href` on `script` is also not animatable.

---

## Element / attribute count summary

| Element | Own attributes (excluding shared groups) | Deprecated |
|---|---|---|
| `image` | `href`, `crossorigin`, `preserveAspectRatio`, `x`, `y`, `width`, `height` = **7** | `xlink:href`, `xlink:title` |
| `foreignObject` | `x`, `y`, `width`, `height` = **4** | none |
| `a` | `href`, `target`, `download`, `ping`, `rel`, `hreflang`, `type`, `referrerpolicy` = **8** | `xlink:href`, `xlink:title` |
| `view` | `viewBox`, `preserveAspectRatio`, `zoomAndPan` = **3** | `viewTarget` (removed, not deprecated) |
| `script` | `type`, `href`, `crossorigin` = **3** | `xlink:href`, `xlink:title` |

**pointer-events keywords:** 10 terminals (`bounding-box`, `visiblePainted`, `visibleFill`, `visibleStroke`, `visible`, `painted`, `fill`, `stroke`, `all`, `none`).

**zoomAndPan keywords:** 2 terminals (`disable`, `magnify`).

**crossorigin keywords:** 2 terminals + empty-string case (`anonymous`, `use-credentials`, ``).

**target keywords:** 4 fixed terminals + 1 open production (`_self`, `_parent`, `_top`, `_blank`, `<XML-Name>`).

**xlink family in SVG 2:** 2 deprecated (`xlink:href`, `xlink:title`); 5 obsolete / not in SVG 2 (`xlink:role`, `xlink:arcrole`, `xlink:show`, `xlink:actuate`, `xlink:type`).
