# Structure grammar notes

## Source

W3C SVG 2 Chapter 5: Document Structure
File: `/Volumes/wd_office_3/Projects/proto-svg/docs/specs/cache/svg2-struct.txt` (3770 lines)
Cross-checked against MDN element and attribute references.

---

## Shared attribute groups

These groups are referenced by every element in this chapter and are reused across the whole grammar.
Member lists are taken verbatim from the spec (§5.12, §5.13) and cross-checked against MDN.

### core-attributes

Defined at §5.12.1. Applicable to **every** SVG element.

```
core-attributes ::=
    id?
    tabindex?
    lang?
    xml:lang?        (* deprecated XML-namespace form; same semantics as lang *)
    xml:space?       (* deprecated; see note *)
    class?
    style?
    data-*           (* open family; any attr in no-namespace starting "data-" *)
```

Member list (spec §5.12.1 literal): `id`, `tabindex`, `lang`, `xml:space`, `class`, `style`, plus all `data-*` custom attributes.

Note: `xml:lang` is the XML-namespace alias for `lang`; both are listed in the spec narrative (§5.12.3). The spec lists `xml:space` in the core group on element tables but separately marks it deprecated (§5.12.4). `xml:lang` is marked deprecated by MDN but is still in the attribute reference. Grammar decision: keep both `lang` / `xml:lang` (mutually exclusive pair, constraint overlay) and `xml:space` (deprecated, allowed, value-set closed).

#### Attribute value syntax for core members

| Attribute | Value syntax | Initial | Animatable | Notes |
|---|---|---|---|---|
| `id` | `XML-Name` (non-empty, no whitespace) | (none) | no | Uniqueness: constraint overlay |
| `tabindex` | `<integer>` | (none) | no | HTML "valid integer" |
| `lang` | `BCP47-language-tag \| ""` | (none) | no | Empty string = unknown |
| `xml:lang` | `BCP47-language-tag \| ""` | (none) | no | Deprecated XML-ns form |
| `xml:space` | `"default" \| "preserve"` | `"default"` | no | Deprecated; use CSS `white-space` |
| `class` | `<CSS-class-list>` (space-separated tokens) | (none) | yes (SMIL) | Open token set |
| `style` | `<CSS-declaration-list>` | (none) | no | |
| `data-*` | any string | (none) | no | Arbitrary name suffix after "data-" |

### conditional-processing-attributes

Defined at §5.7.2. Listed on most but not all elements (absent from `title`, `desc`, `metadata`; see element tables).

```
conditional-processing-attributes ::=
    requiredExtensions?
    systemLanguage?
```

Note: `requiredFeatures` existed in SVG 1.1 and is listed by MDN as deprecated (`⚠️`). It is **not** listed in the SVG 2 conditional-processing definition (§5.7.2) and does not appear on any element table in this chapter. Grammar decision: exclude from SVG 2 grammar; record as removed.

#### Attribute value syntax

| Attribute | Value syntax | Initial | Animatable | Notes |
|---|---|---|---|---|
| `requiredExtensions` | space-separated list of URL tokens | (none) | no | Absent = true; empty string = false |
| `systemLanguage` | comma-separated list of BCP47 language tags | (none) | no | Absent = true; empty string = false |

### aria-attributes

Defined at §5.13. The exact canonical list from §5.13.1 (48 members):

```
aria-attributes ::=
    role?
    aria-activedescendant?
    aria-atomic?
    aria-autocomplete?
    aria-busy?
    aria-checked?
    aria-colcount?
    aria-colindex?
    aria-colspan?
    aria-controls?
    aria-current?
    aria-describedby?
    aria-details?
    aria-disabled?
    aria-dropeffect?
    aria-errormessage?
    aria-expanded?
    aria-flowto?
    aria-grabbed?
    aria-haspopup?
    aria-hidden?
    aria-invalid?
    aria-keyshortcuts?
    aria-label?
    aria-labelledby?
    aria-level?
    aria-live?
    aria-modal?
    aria-multiline?
    aria-multiselectable?
    aria-orientation?
    aria-owns?
    aria-placeholder?
    aria-posinset?
    aria-pressed?
    aria-readonly?
    aria-relevant?
    aria-required?
    aria-roledescription?
    aria-rowcount?
    aria-rowindex?
    aria-rowspan?
    aria-selected?
    aria-setsize?
    aria-sort?
    aria-valuemax?
    aria-valuemin?
    aria-valuenow?
    aria-valuetext?
```

Count: 1 (`role`) + 47 (`aria-*`) = 48 members.

`role` value syntax: space-separated set of WAI-ARIA role tokens (open token set per WAI-ARIA spec, no closed enumeration in SVG 2). Grammar leaf: `<role-token-list>`.

`aria-*` value syntax: per-attribute type from WAI-ARIA §6.6; types include string, token, token-list, integer, number, boolean, ID reference, ID reference list, tristate. Grammar decision: each `aria-*` attribute is an open leaf `<aria-value>` referencing the WAI-ARIA spec; do not enumerate value sets here (they are defined externally and version independently).

### global-event-attributes

Listed identically on every element (58 members). Canonical list from spec (§5.1.4, reproduced on every element):

```
oncancel oncanplay oncanplaythrough onchange onclick onclose oncuechange
ondblclick ondrag ondragend ondragenter ondragexit ondragleave ondragover
ondragstart ondrop ondurationchange onemptied onended onerror onfocus
oninput oninvalid onkeydown onkeypress onkeyup onload onloadeddata
onloadedmetadata onloadstart onmousedown onmouseenter onmouseleave
onmousemove onmouseout onmouseover onmouseup onmousewheel onpause onplay
onplaying onprogress onratechange onreset onresize onscroll onseeked
onseeking onselect onshow onstalled onsubmit onsuspend ontimeupdate
ontoggle onvolumechange onwaiting
```

Count: 58. All have value syntax `<event-handler-script>` (open string). Animatable: no.

### document-event-attributes

Listed only on `svg` (§5.1.4, heading "document event attributes"):

```
onunload onabort onerror onresize onscroll
```

Note: `onerror`, `onresize`, `onscroll` appear in both the global-event-attributes list and the document-event-attributes list for `svg`. The spec note at §5.1.4 says these three "replace the generic event handlers with the same names normally supported by SVG elements" — meaning on `svg` these are Window-object handlers. Grammar decision: record the 5-member document-event group separately; note the 3-way overlap with global events is a constraint (they are different handlers on `svg`).

### document-element-event-attributes

Listed on most elements (§5.1.4 onwards, heading "document element event attributes"):

```
oncopy oncut onpaste
```

Count: 3.

### graphical-event-attributes

Listed on most elements (heading "graphical event attributes"):

```
onfocusin onfocusout
```

Count: 2.

Note: These two are absent from `title`, `desc`, and `metadata` element tables. Spec does not give a reason; likely because those are never-rendered elements and focus does not apply. MDN does not list them as event attributes for `desc`/`title`/`metadata`. Grammar decision: exclude graphical-event-attributes from `title`, `desc`, `metadata`.

### presentation-attributes

Referenced by every renderable element as "presentation attributes —" (with varying additional members per element). This group is defined in the SVG 2 styling chapter (not this chapter); referenced here by name only.

Grammar decision: define as a named alias `<presentation-attributes>` to be expanded in the styling grammar module.

---

## Category definitions (from §5.1.3 and §5.2.1)

These category memberships are used in content-model expressions.

| Category | Members defined in this chapter |
|---|---|
| structural element | `defs`, `g`, `svg`, `symbol`, `use` |
| structurally external element (when `href` present) | `audio`, `foreignObject`, `iframe`, `image`, `script`, `use`, `video` |
| container element | `a`, `clipPath`, `defs`, `g`, `marker`, `mask`, `pattern`, `svg`, `switch`, `symbol`, `unknown` |
| graphics element | `audio`, `canvas`, `circle`, `ellipse`, `foreignObject`, `iframe`, `image`, `line`, `path`, `polygon`, `polyline`, `rect`, `text`, `textPath`, `tspan`, `video` |
| graphics referencing element | `audio`, `iframe`, `image`, `use`, `video` |
| descriptive element | `desc`, `metadata`, `title` |
| never-rendered element | `defs` (UA stylesheet `display:none`), `symbol` (UA stylesheet), `title`, `desc`, `metadata` |
| renderable element | `svg`, `g`, `switch`, `use`, `unknown`, and all graphics elements |
| animation elements | `animate`, `animateMotion`, `animateTransform`, `discard`, `set` |
| paint server elements | `linearGradient`, `radialGradient`, `pattern` |
| shape elements | `circle`, `ellipse`, `line`, `path`, `polygon`, `polyline`, `rect` |

---

## Elements

### `svg`

**Defined in:** §5.1.4

**Categories:** container element, renderable element, structural element

**Content model:** Any number, any order, of:
- animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- descriptive elements: `desc`, `title`, `metadata`
- paint server elements: `linearGradient`, `radialGradient`, `pattern`
- shape elements: `circle`, `ellipse`, `line`, `path`, `polygon`, `polyline`, `rect`
- structural elements: `defs`, `g`, `svg`, `symbol`, `use`
- plus: `a`, `audio`, `canvas`, `clipPath`, `filter`, `foreignObject`, `iframe`, `image`, `marker`, `mask`, `script`, `style`, `switch`, `text`, `video`, `view`

**Element-specific attributes with value syntax:**

| Attribute | Value syntax | Initial | Animatable | Notes |
|---|---|---|---|---|
| `viewBox` | `<number> <number> <number> <number>` (min-x min-y width height) | (none) | yes | Defined in SVGFitToViewBox mixin |
| `preserveAspectRatio` | `<align> [<meetOrSlice>]` | `xMidYMid meet` | yes | See below |
| `zoomAndPan` | `"disable" \| "magnify"` | `"magnify"` | no | Deprecated per MDN (`⚠️`); kept in SVG 2 spec |
| `transform` | `<transform-list>` | (none) | yes | New in SVG 2 (previously disallowed on `svg`) |
| `x` | `<length-percentage>` | `0` | yes | Geometry property; on nested `svg` only |
| `y` | `<length-percentage>` | `0` | yes | Geometry property; on nested `svg` only |
| `width` | `<length-percentage> \| "auto"` | `"auto"` | yes | Geometry property |
| `height` | `<length-percentage> \| "auto"` | `"auto"` | yes | Geometry property |
| `version` | text string | (none) | no | Deprecated per MDN (`⚠️`); informational only in SVG 2 |
| `baseProfile` | text string | (none) | no | Deprecated per MDN (`⚠️`); was SVG 1.x |

`preserveAspectRatio` closed value set:
```
align ::= "none"
        | "xMinYMin" | "xMidYMin" | "xMaxYMin"
        | "xMinYMid" | "xMidYMid" | "xMaxYMid"
        | "xMinYMax" | "xMidYMax" | "xMaxYMax"
meetOrSlice ::= "meet" | "slice"
preserveAspectRatio-value ::= align ( " " meetOrSlice )?
```

**Shared attribute groups included:** core-attributes, conditional-processing-attributes, aria-attributes, global-event-attributes, document-event-attributes, document-element-event-attributes, graphical-event-attributes, presentation-attributes

**DOM interface:** `SVGSVGElement : SVGGraphicsElement` + `SVGFitToViewBox` + `SVGZoomAndPan` + `WindowEventHandlers`

**Context-sensitive constraints (overlay, not grammar):**
- `id` must be unique in the node tree.
- Outermost `svg`: `x` and `y` have no effect.
- `width`/`height` computed `auto` treated as `100%`.
- `zoomAndPan` only controls user interaction; `currentScale`/`currentTranslate` always writable.
- `version` and `baseProfile` are deprecated; user agents must not use them for conformance decisions.
- `onblur`, `onerror`, `onfocus`, `onload`, `onscroll` on `svg` are Window-object handlers (replace element-level handlers of same names).

---

### `g`

**Defined in:** §5.2.2

**Categories:** container element, renderable element, structural element

**Content model:** Same as `svg` — any number, any order, of all the same children (animation elements, descriptive elements, paint server elements, shape elements, structural elements, `a`, `audio`, `canvas`, `clipPath`, `filter`, `foreignObject`, `iframe`, `image`, `marker`, `mask`, `script`, `style`, `switch`, `text`, `video`, `view`).

**Element-specific attributes:** None beyond shared groups. No geometry properties. `transform` comes from presentation-attributes.

**Shared attribute groups included:** core-attributes, conditional-processing-attributes, aria-attributes, global-event-attributes, document-element-event-attributes, graphical-event-attributes, presentation-attributes

Note: `g` does NOT have document-event-attributes (onunload/onabort/etc.), unlike `svg`.

**DOM interface:** `SVGGElement : SVGGraphicsElement`

**Context-sensitive constraints (overlay):**
- `id` uniqueness.
- A named `g` can be the target of animation and re-use references.

---

### `unknown`

**Defined in:** §5.3

**Categories:** container element, renderable element

Note: NOT a structural element (unlike `g`).

**Content model:** Any elements or character data (completely open).

**Element-specific attributes:** None beyond shared groups.

**Shared attribute groups included:** core-attributes, conditional-processing-attributes, aria-attributes, global-event-attributes, document-element-event-attributes, graphical-event-attributes, presentation-attributes

**DOM interface:** `SVGUnknownElement : SVGGraphicsElement`

**Context-sensitive constraints (overlay):**
- Rendered as `g` equivalent outside text context; as `tspan` equivalent inside text context.
- The feature (treating unknown as `g`) is "at risk" per spec; no known implementations.
- Any global attribute/property valid on any SVG graphics element is also valid on unknown elements.

---

### `defs`

**Defined in:** §5.4.2

**Categories:** container element, never-rendered element, structural element

Note: NOT listed as renderable.

**Content model:** Same as `g` — any number, any order, of all the same children (spec explicitly states content model is identical to `g`).

**Element-specific attributes:** None beyond shared groups. No geometry properties.

**Shared attribute groups included:** core-attributes, global-event-attributes, document-element-event-attributes, graphical-event-attributes, presentation-attributes

Note: `defs` does NOT include aria-attributes or conditional-processing-attributes per its element table. This is a significant divergence from `g` — the spec omits them from the `defs` table while including them on `g`, `switch`, `use`, and `symbol`. This is likely a spec authoring omission (not intentional) — see Discrepancies.

**DOM interface:** `SVGDefsElement : SVGGraphicsElement`

**Context-sensitive constraints (overlay):**
- UA stylesheet sets `display:none` on `defs` with `!important`; this cannot be overridden.
- Despite `display:none`, descendant elements can be referenced by other elements (e.g., via `use`, paint server references).
- Recommended practice: place all referenced elements inside `defs`.

---

### `symbol`

**Defined in:** §5.5

**Categories:** container element, structural element

Note: NOT a renderable element (never rendered directly; only through `use` instantiation).

**Content model:** Same as `svg`/`g`/`defs` — any number, any order, of all the same children (animation elements, descriptive elements, paint server elements, shape elements, structural elements, `a`, `audio`, `canvas`, `clipPath`, `filter`, `foreignObject`, `iframe`, `image`, `marker`, `mask`, `script`, `style`, `switch`, `text`, `video`, `view`).

**Element-specific attributes with value syntax:**

| Attribute | Value syntax | Initial | Animatable | Notes |
|---|---|---|---|---|
| `viewBox` | `<number> <number> <number> <number>` | (none) | yes | Via SVGFitToViewBox |
| `preserveAspectRatio` | see `svg` entry above | `xMidYMid meet` | yes | Via SVGFitToViewBox |
| `refX` | `<length> \| "left" \| "center" \| "right"` | (none; see note) | yes | New in SVG 2 |
| `refY` | `<length> \| "top" \| "center" \| "bottom"` | (none; see note) | yes | New in SVG 2 |
| `x` | `<length-percentage>` | `0` | yes | Geometry property; default size of instantiated symbol |
| `y` | `<length-percentage>` | `0` | yes | Geometry property |
| `width` | `<length-percentage> \| "auto"` | `"auto"` | yes | Geometry property; `auto` treated as `100%` |
| `height` | `<length-percentage> \| "auto"` | `"auto"` | yes | Geometry property |

`refX`/`refY` keyword semantics: `left`/`top` = 0%, `center` = 50%, `right`/`bottom` = 100% of viewport in that direction (after `viewBox` and `preserveAspectRatio` are applied).

`refX`/`refY` absence note: When not specified, the behavior is different from specifying `0` (for backwards compatibility); the top/left side of the viewport is placed at the x,y point without any adjustment. This is a constraint, not a grammar distinction.

**Shared attribute groups included:** core-attributes, aria-attributes, global-event-attributes, document-element-event-attributes, graphical-event-attributes, presentation-attributes

Note: `symbol` does NOT have conditional-processing-attributes per its element table. This may be another spec authoring gap — see Discrepancies.

**DOM interface:** `SVGSymbolElement : SVGGraphicsElement` + `SVGFitToViewBox`

**Context-sensitive constraints (overlay):**
- UA stylesheet: `display:none` on `symbol` with `!important`.
- When instantiated as direct referenced element of `use`: display forced to `inline` with `!important`.
- Clipped at symbol viewport bounds by default (`overflow:hidden` in UA stylesheet).
- `width`/`height` on `use` override `symbol`'s `width`/`height` for viewport sizing.
- Cross-origin use references are not permitted.

---

### `use`

**Defined in:** §5.6

**Categories:** graphics referencing element, renderable element, structural element, structurally external element (when `href` is present)

**Content model:** Any number, any order, of:
- animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- descriptive elements: `desc`, `title`, `metadata`
- `clipPath`, `mask`, `script`, `style`

Note: `use` content model is significantly more restricted than `g`/`svg`/`defs`/`symbol`. It does not allow shape elements, structural elements (other than implicitly via the shadow tree), or most container/graphics elements as direct children.

**Element-specific attributes with value syntax:**

| Attribute | Value syntax | Initial | Animatable | Notes |
|---|---|---|---|---|
| `href` | `<URL>` | (none) | yes | Primary reference attribute in SVG 2 |
| `xlink:href` | `<URL>` | (none) | yes | Deprecated; kept for compatibility |
| `xlink:title` | text | (none) | no | Deprecated |
| `x` | `<length-percentage>` | `0` | yes | Geometry property; becomes translate(x,y) transformation |
| `y` | `<length-percentage>` | `0` | yes | Geometry property |
| `width` | `<length-percentage> \| "auto"` | `"auto"` | yes | Geometry property; overrides viewport on svg/symbol only |
| `height` | `<length-percentage> \| "auto"` | `"auto"` | yes | Geometry property |

`href` open type: `<URL>` — fragment-only (`#id`), relative, or absolute. May be a document URL without a fragment (references root element of the document). Grammar leaf: `<url-reference>`.

**Shared attribute groups included:** core-attributes, conditional-processing-attributes, aria-attributes, global-event-attributes, document-element-event-attributes, graphical-event-attributes, presentation-attributes

**DOM interface:** `SVGUseElement : SVGGraphicsElement` + `SVGURIReference`

**Context-sensitive constraints (overlay):**
- `href` must resolve to an SVG element or an HTML-namespaced element allowed in SVG containers; otherwise `use` is in error and not rendered.
- Circular references (direct or indirect) are invalid; the `use` creating the cycle must not render.
- Cross-origin external references are not permitted (security restriction).
- Negative `width` or `height` is an illegal value.
- Zero `width` or `height` disables rendering of the viewport.
- `width`/`height` only affect referenced element if it defines a viewport (`svg` or `symbol`).
- `x`/`y` create an additional `translate(x,y)` applied to the `use` element itself (affects `userSpaceOnUse` masks/clips/filters).
- Shadow tree is read-only; direct modification throws `NoModificationAllowedError`.
- `script` elements in shadow tree are inert.
- `id` attribute is cloned into shadow tree instances (no uniqueness conflict because different node trees).
- When both `href` and `xlink:href` are present, `href` takes precedence.

---

### `switch`

**Defined in:** §5.7.3

**Categories:** container element, renderable element

Note: NOT a structural element.

**Content model:** Any number, any order, of:
- animation elements: `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- shape elements: `circle`, `ellipse`, `line`, `path`, `polygon`, `polyline`, `rect`
- `a`, `audio`, `canvas`, `foreignObject`, `g`, `iframe`, `image`, `svg`, `switch`, `text`, `use`, `video`

Note: `switch` content model excludes `defs`, `symbol`, `clipPath`, `filter`, `linearGradient`, `radialGradient`, `pattern`, `marker`, `mask`, `script`, `style`, `view` that are in `g`'s content model. Also excludes descriptive elements (`desc`, `title`, `metadata`) from the child list. This is the narrowest container content model.

**Element-specific attributes:** None beyond shared groups.

**Shared attribute groups included:** core-attributes, conditional-processing-attributes, aria-attributes, global-event-attributes, document-element-event-attributes, graphical-event-attributes, presentation-attributes

**DOM interface:** `SVGSwitchElement : SVGGraphicsElement`

**Context-sensitive constraints (overlay):**
- Only the first child whose `requiredExtensions` and `systemLanguage` both evaluate to true is rendered; all others are bypassed (as if `display:none`).
- `display:none` or `visibility:hidden` on a child does NOT affect the true/false test; conditional processing is independent.
- `script` and `style` elements are not affected by the `switch` evaluation.
- If no child tests true, nothing renders.
- `systemLanguage` evaluation uses `allowReorder="yes"` semantics per SMIL spec.

---

### `desc`

**Defined in:** §5.8.1

**Categories:** descriptive element, never-rendered element

**Content model:** Any elements or character data (completely open).

**Element-specific attributes:** None beyond shared groups.

**Shared attribute groups included:** core-attributes, global-event-attributes, document-element-event-attributes

Note: `desc` does NOT include: conditional-processing-attributes, aria-attributes, graphical-event-attributes. This differs from most container/renderable elements.

**DOM interface:** `SVGDescElement : SVGElement` (not `SVGGraphicsElement`)

**Context-sensitive constraints (overlay):**
- UA stylesheet: `display:none` with `!important`.
- Multiple sibling `desc` elements must have distinct language tags (constraint: unique-lang per sibling group).
- UA must select the `desc` element whose language best matches user preferences.
- `aria-describedby` takes precedence over child `desc`.
- Must not be empty or whitespace-only (authoring requirement; UA may ignore empty `desc`).

---

### `title`

**Defined in:** §5.8.1

**Categories:** descriptive element, never-rendered element

**Content model:** Any elements or character data (completely open).

**Element-specific attributes:** None beyond shared groups.

**Shared attribute groups included:** core-attributes, global-event-attributes, document-element-event-attributes

Note: `title` does NOT include: conditional-processing-attributes, aria-attributes, graphical-event-attributes. Same restrictions as `desc`.

**DOM interface:** `SVGTitleElement : SVGElement` (not `SVGGraphicsElement`)

**Context-sensitive constraints (overlay):**
- UA stylesheet: `display:none` with `!important`.
- Multiple sibling `title` elements must have distinct language tags.
- UA must select `title` matching user language preferences; if no match, first `title` is used.
- Interactive UAs must expose `title` content to platform accessibility APIs.
- Root `svg` element should have a `title` child for standalone SVG documents.
- Must not be empty or whitespace-only.

---

### `metadata`

**Defined in:** §5.9

**Categories:** descriptive element, never-rendered element

**Content model:** Any elements or character data (completely open).

**Element-specific attributes:** None beyond shared groups.

**Shared attribute groups included:** core-attributes, global-event-attributes, document-element-event-attributes

Note: `metadata` does NOT include: conditional-processing-attributes, aria-attributes, graphical-event-attributes. Same restrictions as `desc`/`title`.

**DOM interface:** `SVGMetadataElement : SVGElement` (not `SVGGraphicsElement`)

**Context-sensitive constraints (overlay):**
- UA stylesheet: `display:none` with `!important`.
- Content should consist of elements from other namespaces (RDF, Dublin Core, etc.).
- `data-*` attributes may be added to `metadata` when not associated with another element.

---

## Open datatypes used

These datatypes appear in attribute value syntax above and are defined elsewhere (to be expanded in other grammar modules or as named terminals):

| Leaf name | Description | Defined in |
|---|---|---|
| `<length>` | SVG/CSS length with optional unit | SVG 2 Basic Data Types chapter |
| `<length-percentage>` | `<length>` or `<percentage>` | CSS Values and Units |
| `<number>` | Real number (floating point) | SVG 2 Basic Data Types |
| `<integer>` | Signed decimal integer | CSS Values |
| `<URL>` / `<url-reference>` | URL as per RFC 3987; may be fragment-only | RFC 3987 |
| `<transform-list>` | List of CSS transform functions | CSS Transforms |
| `<CSS-class-list>` | Space-separated CSS class names | HTML / CSS |
| `<CSS-declaration-list>` | Inline CSS declarations | CSS |
| `<BCP47-language-tag>` | IETF BCP 47 language tag | BCP 47 |
| `<event-handler-script>` | ECMAScript expression or statement list | HTML |
| `<role-token-list>` | Space-separated WAI-ARIA role tokens | WAI-ARIA |
| `<aria-value>` | Per-attribute ARIA value type | WAI-ARIA §6.6 |
| `XML-Name` | XML name token (NCName in namespace-aware context) | XML 1.0 |

---

## Discrepancies, doc gaps & roadblocks

### D1: `defs` missing aria-attributes and conditional-processing-attributes
**Spec text:** The `defs` element table (§5.4.2) lists only `core attributes`, `graphical event attributes`, `global event attributes`, `document element event attributes`, and `presentation attributes`. It omits `aria attributes` and `conditional processing attributes`.

**MDN:** MDN's `<defs>` page does not enumerate specific attributes; it inherits from container element lists.

**Analysis:** This is almost certainly a spec authoring omission. `defs` is a container element with the same content model as `g`. There is no semantic reason to exclude ARIA or conditional processing. Implementations (Chrome, Firefox) do accept `aria-*` and `requiredExtensions` on `defs`.

**Grammar decision:** Include `aria-attributes` and `conditional-processing-attributes` in `defs`, matching `g`. Document discrepancy.

### D2: `symbol` missing conditional-processing-attributes
**Spec text:** The `symbol` element table (§5.5) omits `conditional processing attributes`, though it includes `aria attributes`.

**Analysis:** Likely authoring gap. `symbol` can be referenced conditionally and there is no reason to exclude `requiredExtensions`/`systemLanguage`.

**Grammar decision:** Include `conditional-processing-attributes` in `symbol`. Document discrepancy.

### D3: `zoomAndPan` is deprecated per MDN, not per SVG 2 spec
**Spec text:** §5.1.4 defines `zoomAndPan` with values `disable | magnify` and initial `magnify`. Animatable: no.

**MDN:** `zoomAndPan` is marked `⚠️ Deprecated`.

**Analysis:** The SVG 2 spec still includes it without a deprecation notice in the attribute definition. MDN marks it deprecated because browsers have removed or are removing support for zoom-and-pan via this attribute (Chrome removed the interactive magnify mode). The spec text is older than current browser implementation.

**Grammar decision:** Include `zoomAndPan` as `"disable" | "magnify"` with a note that it is deprecated/unimplemented in practice.

### D4: `version` and `baseProfile` absent from SVG 2 element table but in MDN
**Spec text:** The `svg` attribute table in §5.1.4 does NOT list `version` or `baseProfile` in the formal attribute box. They appear only in code examples in the prose.

**MDN:** Both are listed as deprecated (`⚠️`).

**Analysis:** SVG 2 effectively drops these attributes as meaningful; they were SVG 1.x conformance signals. They should still be parseable (ignored) for compatibility.

**Grammar decision:** Include both as optional deprecated attributes with open string values. Mark as deprecated in grammar comments.

### D5: `requiredFeatures` removed from SVG 2 conditional processing
**Spec text:** §5.7.1 explicitly states `requiredFeatures` existed in previous versions and was removed due to poor specification and implementation.

**MDN:** Lists `requiredFeatures` as `⚠️ Deprecated`.

**Grammar decision:** Exclude `requiredFeatures` from SVG 2 grammar. If parsing legacy SVG 1.1 content, treat as an unknown attribute to be ignored.

### D6: `ondragexit` — non-standard vs `ondragleave`
**Spec text:** The global-event-attributes list includes both `ondragexit` and `ondragleave`.

**Analysis:** `dragexit` is a non-standard alias for `dragleave` that was historically in Firefox. The HTML Living Standard removed it. Current browsers fire `dragleave` not `dragexit`. MDN attributes list does not include `ondragexit` as a valid SVG event attribute.

**Grammar decision:** Include `ondragexit` in the grammar (it is in the spec), but note it is non-standard/obsolete in practice. Constraint overlay: prefer `ondragleave`.

### D7: `onmousewheel` — non-standard, deprecated in favor of `onwheel`
**Spec text:** Global-event-attributes include `onmousewheel` but NOT `onwheel`.

**MDN:** `onmousewheel` is non-standard. Modern browsers use `wheel` event / `onwheel`.

**Grammar decision:** Include `onmousewheel` as specified (it is in the spec text). Note that `onwheel` is not in the SVG 2 spec's event list but browsers support it. Add `onwheel` as a roadblock item for the grammar: it may need to be added as a real-browser behavior not captured in the spec.

### D8: `onkeypress` — deprecated in DOM
**Spec text:** `onkeypress` is in the global-event-attributes list.

**Analysis:** `keypress` event is deprecated in DOM Living Standard. Still fires in all browsers but may be removed eventually.

**Grammar decision:** Include as-is; note deprecated status in constraint overlay.

### D9: `onshow` — removed from HTML Living Standard
**Spec text:** `onshow` is in the global-event-attributes list.

**Analysis:** The `show` event was for `<menu type="context">`, which was removed from HTML. `onshow` effectively never fires on SVG elements.

**Grammar decision:** Include as-is (it is in the spec); note it is vestigial/non-functional.

### D10: `switch` excludes descriptive elements from content model
**Spec text:** The `switch` content model does not include `desc`, `title`, `metadata`.

**Analysis:** This differs from `g`, `svg`, `defs`, `symbol`. Not clearly intentional — `switch` children that evaluate to false need some way to be documented. However, this is what the spec says. MDN's `<switch>` page does not contradict this.

**Grammar decision:** Follow spec; `switch` does not allow descriptive elements as direct children. Record as a potential spec gap.

### D11: `use` content model allows `clipPath`, `mask`, `script`, `style` but not `desc`, `title`, `metadata`
**Spec text:** §5.6 lists `use` content as animation elements, descriptive elements, `clipPath`, `mask`, `script`, `style`.

Wait — re-reading spec line 968-974: "descriptive elements — 'desc', 'title', 'metadata'" IS listed. And also `clipPath, mask, script, style`.

**Correction:** `use` DOES allow descriptive elements. The full `use` content model is:
- animation elements
- descriptive elements (`desc`, `title`, `metadata`)
- `clipPath`, `mask`, `script`, `style`

This is confirmed by the SVG 2 example showing `<title>` children of `<use>` (§5.8.1 example).

### D12: `xml:lang` listed in core-attributes narrative but appears inconsistently in element tables
**Spec text:** §5.12.3 discusses both `lang` (no namespace) and `xml:lang` (XML namespace). The element tables list only `lang` in the core-attributes line.

**MDN:** `xml:lang` is marked `⚠️ Deprecated`.

**Grammar decision:** `xml:lang` is the deprecated XML-namespace form; grammar includes it as an alias. In HTML-parsed SVG, use `lang` only. In XML-parsed SVG, both are valid (must match if both present).

### D13: `graphical-event-attributes` group naming vs location
**Spec text:** The two attributes `onfocusin` and `onfocusout` are grouped as "graphical event attributes" in element tables. They are NOT in the main global-event-attributes list (which has `onfocus` not `onfocusin`).

**Note:** `onfocusin` / `onfocusout` bubble; `onfocus` / `onblur` do not. This is a real semantic distinction. The spec correctly separates them.

**Grammar decision:** Maintain as separate group `graphical-event-attributes = { onfocusin, onfocusout }`.

### D14: `desc` and `title` element tables omit aria-attributes, but ARIA spec says they may have role
**Spec text:** §5.13.4 states for `desc`: "no role may be applied" and for `title`: "no role may be applied". Element tables for both omit `aria-attributes`.

**Grammar decision:** Follow spec; `desc` and `title` carry no ARIA attributes, not even `role`. The ARIA table at §5.13.4 confirms this.

### Roadblock R1: `onfocusin`/`onfocusout` browser support
Modern browsers support `focus`/`blur` events but `focusin`/`focusout` have had inconsistent SVG support. Grammar includes them per spec; real-browser testing needed.

### Roadblock R2: `transform` on `svg` — SVG 2 new feature
SVG 1.x did not allow `transform` on the `svg` element. SVG 2 adds it. Browser support is uneven. Grammar includes it; constraint overlay should flag it for compat testing.

### Roadblock R3: `refX`/`refY` on `symbol` — SVG 2 new feature, partial implementation
Resolved in SVG 2 spec. Chrome supports it; Safari and Firefox support may vary. Grammar includes it; testing required.

### Roadblock R4: Shadow DOM / `use` element behavior
The spec's shadow-DOM model for `use` is the normative SVG 2 model; the SVG 1.1 style-cloning model is explicitly deprecated. Browser implementations have moved to shadow DOM. No grammar impact, but constraint overlay for CSS specificity / selector matching behavior is complex.
