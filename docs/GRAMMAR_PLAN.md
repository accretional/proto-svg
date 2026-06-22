# proto-svg Grammar Authoring Contract

The single authoring contract for the SVG EBNF grammar. Every per-module author
follows this verbatim; nothing below should need re-deriving. It reconciles all 14
module grammar-notes against the policy in `CONTEXT_SENSITIVITY.md`, the pipeline in
`PIPELINE_PORTING.md`, and the ISO-14977 dialect of `proto-css/lang/*.ebnf`. Scope and
discrepancy resolutions live in `DOC_GAPS.md`; this file is the structure.

EBNF dialect (matches proto-css exactly): `=` definition, `,` concatenation, `|`
alternation, `[ ]` optional, `{ }` repetition, `( )` group, `"literal"` terminal, `;`
terminator, `(* comment *)` (stripped before parse — gluon bug workaround).
Lowercase `snake_case` rule names = lexical/primitive leaves; `PascalCase` rule names =
syntactic messages (one proto message each).

---

## a. File layout

`lang/*.ebnf` files, concatenated by `genproto` and `gen` in this exact order. The
**ROOT file is first** so its root production is the prune root.

| # | File | Contains | Root prod |
|---|------|----------|-----------|
| 1 | `svg.ebnf` | **ROOT.** `SvgDocument`, the `Svg` element, content-category unions (animationElements, descriptiveElements, shapeElements, structuralElements, containerElements, gradientElements, filterPrimitiveElements, paintServerElements, textContentChildElements, lightSourceElements), the shared attribute-group productions (Core/Conditional/Aria/Presentation/event groups/Xlink), and `g`/`defs`/`switch`/`unknown`/`a` containers | `SvgDocument` |
| 2 | `datatype.ebnf` | Open leaf `*Type` rules (§c), `<viewBox>`, `<preserveAspectRatio>`, `<paint>` structure, `<dasharray>`, `<transform-list>` value, `<color>` arms, primitive helpers (`digit`, `whole_number`, `wsp`, `comma_wsp`, string/ident leaves) | — |
| 3 | `transform.ebnf` | SVG `transform`-attribute grammar (CSS Transforms §6.2, bare numbers, 6 functions); reused by transform/gradientTransform/patternTransform | `TransformList` |
| 4 | `path.ebnf` | Path-data BNF (`SvgPath` and all draw commands, corrected `number`); reused by `d`, `animateMotion path`, `path()` basic-shape | `SvgPath` |
| 5 | `structure.ebnf` | `Svg` element-specific attrs not in svg.ebnf root, `Symbol`, `Use`, descriptive elements (`Desc`/`Title`/`Metadata`) | — |
| 6 | `shapes.ebnf` | `Rect`, `Circle`, `Ellipse`, `Line`, `Polyline`, `Polygon`, geometry-property attribute productions, `<points>` | — |
| 7 | `text.ebnf` | `Text`, `Tspan`, `TextPath`, text positioning-attribute lists, text/font property values | — |
| 8 | `embedded.ebnf` | `Image`, `ForeignObject`, `View`, `Script`, `Style`, href/xlink family, `<crossorigin>`, `<target>`, fragment identifiers, `pointer-events`, `zoomAndPan` | — |
| 9 | `painting.ebnf` | fill/stroke/stroke-* property values, `paint-order`, `<paint>` consumers, `Marker`, marker attributes, render-hint props | — |
| 10 | `pservers.ebnf` | `LinearGradient`, `RadialGradient`, `Stop`, `Pattern`, gradient/pattern attribute productions, `spreadMethod`/`gradientUnits` | — |
| 11 | `filter.ebnf` | `Filter`, all 17 primitives, `feFunc*`, light sources, `feMergeNode`, `<fe-in-value>`, filter props (flood/lighting/color-interpolation-filters), CSS `filter` value list | — |
| 12 | `masking.ebnf` | `ClipPath`, `Mask`, `clip-path`/`clip-rule`/`clip`, `<basic-shape>`, `<geometry-box>` (two variants), mask shorthand/longhands, `mask-type` | — |
| 13 | `animation.ebnf` | `Animate`, `Set`, `AnimateMotion`, `AnimateTransform`, `Mpath`, `Discard`, SMIL timing grammar, `AnimationValue` union, animateTransform per-type productions | — |
| 14 | `styling.ebnf` | Presentation-attribute property value rules that are shared CSS-property syntaxes (display/visibility/overflow/render-hints/font/etc.) referenced by `PresentationAttributes` | — |

**Root production name: `SvgDocument`.** Definition:

```ebnf
SvgDocument = Svg ;
```

`SvgDocument` wraps the single outermost `svg` element. The prune root passed to
`pruneUnreachable` in genproto is `"SvgDocument"`. (`Svg` is the element; `SvgDocument`
is the document wrapper, matching proto-css's `CssStyleSheet` top wrapper convention.)

---

## b. Naming conventions

Each EBNF rule → one PascalCase proto message. Leaf rules named in the genproto
`leafTypes` set → string-valued. reps.go keys = those PascalCase leaf message names.

| Kind | Convention | Example |
|------|-----------|---------|
| Element rule | `Pascal(elementName)` | `rect` → `Rect`; `feGaussianBlur` → `FeGaussianBlur`; `linearGradient` → `LinearGradient` |
| Per-element attribute production | `Pascal(element) + Pascal(attr) + "Attr"` | `rect x` → `RectXAttr`; `feSpotLight x` → `FeSpotLightXAttr`; `animate fill` → `AnimateFillAttr` |
| Shared attribute group | descriptive Pascal name + `Attributes` | `CoreAttributes`, `PresentationAttributes` |
| Value-type rule (open leaf) | `*Type` (PascalCase) | `LengthType`, `ColorType`, `PaintType`, `AngleType` |
| Closed value/enum production | `Pascal(attr)` or `Pascal(element+attr)` + `Value` when ambiguous | `SpreadMethod`, `MarkerOrientValue`, `StrokeLinejoin` |
| Enum keyword terminal | quoted literal inside the production | `"pad" \| "reflect" \| "repeat"` |
| Content-category union | `Pascal(category) + "Element"` (singular message wrapping the oneof) | `ShapeElement`, `AnimationElement`, `FilterPrimitiveElement` |
| List/sequence value | `Pascal(...) + "List"` | `PointsList`, `KeyTimesList`, `AnimationValueList` |
| Lexical primitive | `snake_case` | `digit`, `whole_number`, `comma_wsp`, `wsp`, `hex_digit`, `string_literal` |

### Collision avoidance (CRITICAL — gluon silently dedupes)

`norm(s) = strings.ToLower(strings.ReplaceAll(s, "_", ""))`. Any two rule names that
normalize equal are **silently merged** into one message by genproto's `dedupeMessages`.
Rules:

1. **Never reuse a bare attribute name across elements.** `x`, `y`, `width`, `height`,
   `cx`, `cy`, `r`, `rx`, `ry`, `type`, `offset`, `fill`, `mode`, `operator`, `result`,
   `in`, `href`, `values`, `d`, `path`, `points`, `dx`, `dy` all appear on many
   elements with different value sets. Each gets a fully-qualified production
   `Element + Attr + Attr` (e.g. `FeOffsetDxAttr`, `TextDxAttr`, `FeDropShadowDxAttr`).
   The renderer selects by **position in the tree** (CONTEXT_SENSITIVITY "What is not
   context-sensitive"), so no overlay is needed for these collisions.

2. **The `fill` paint-vs-animation-freeze collision is resolved positionally.**
   `FillAttr` (paint, on shapes/text/containers) vs `AnimateFillAttr`,
   `AnimateMotionFillAttr`, `AnimateTransformFillAttr`, `SetFillAttr` (= `freeze|remove`).
   These never normalize-collide and the correct value set is chosen by which element's
   production references it.

3. **`type` collision (≥4 distinct meanings) resolved positionally:**
   `AnimateTransformTypeAttr` (transform kind), `FeColorMatrixTypeAttr`,
   `FeTurbulenceTypeAttr`, `FeFuncTypeAttr` (transfer fn), `ScriptTypeAttr`/`StyleTypeAttr`
   (media type), `ATypeAttr` (link MIME type). Distinct messages, no overlay.

4. **Leaf-type `*Type` rules are singletons** (one `LengthType`, one `NumberType`, …)
   and ARE meant to be shared — they go in `datatype.ebnf` once. Do not redefine them
   per module; reference the single rule. (genproto's intentional `*Type` dedupe of the
   snake_case atom vs PascalCase wrapper is the only sanctioned collision.)

5. **Avoid singular/plural and case-only differences** that normalize equal
   (`Stop` element vs a hypothetical `stop` keyword rule — keep the element `Stop` and
   never create a lexical `stop` rule).

---

## c. Leaf datatypes (the scalarize set)

Directly transcribable into `scalarize.go`'s `leafTypes` and `reps.go`. Each leaf:
its `leafTypes` normalized name (`norm` form), the PascalCase reps.go key, and 4–6
sample strings. `norm` is what `scalarizeLeaves` matches; the reps.go map key is the
compiled PascalCase message name.

| `leafTypes` (norm) | reps.go key (PascalCase) | Samples |
|---|---|---|
| `numbertype` | `NumberType` | `0` `1` `-1` `3.14` `0.5` `-0.001` |
| `integertype` | `IntegerType` | `0` `1` `-1` `100` `3` `32767` |
| `lengthtype` | `LengthType` | `0` `10px` `2.5em` `1.5rem` `50vw` `12pt` |
| `percentagetype` | `PercentageType` | `0%` `50%` `100%` `25%` `-10%` `80%` |
| `lengthpercentagetype` | `LengthPercentageType` | `10px` `50%` `2.5em` `100%` `0` `1.5rem` |
| `coordinatetype` | `CoordinateType` | `0` `10` `-5.5` `100px` `50%` `2.5em` |
| `angletype` | `AngleType` | `0deg` `45deg` `90deg` `180deg` `1.5708rad` `0.25turn` |
| `timetype` | `TimeType` | `0s` `0.3s` `1s` `200ms` `2.5s` `1.5s` |
| `colortype` | `ColorType` | `red` `#ff0000` `#fff` `rgb(255,0,0)` `rgba(0,0,255,0.5)` `currentColor` |
| `iritype` | `IriType` | `#circle1` `#grad1` `url(#marker1)` `path.svg#p` `https://ex.com/d.svg#m` |
| `urltype` | `UrlType` | `url(#id)` `url(#grad1)` `url(photo.png)` `url("a.svg#m")` |
| `stringtype` | `StringType` | `"label"` `"Aa"` `"Specimen"` `"en"` `"text/css"` |
| `customidenttype` | `CustomIdentType` | `blur1` `result1` `myFilter` `shadow` `out` |
| `xmlnametype` | `XmlNameType` | `circle1` `grad-a` `myId` `node_3` `r1` |
| `bcp47type` | `Bcp47Type` | `en` `fr-CA` `de` `zh-Hans` `pt-BR` |
| `listofnumberstype` | `ListOfNumbersType` | `1 2 3 4` `0,1,0,0` `1 0.5` `0 0 1 1 0` `2 4 6` |
| `listoflengthstype` | `ListOfLengthsType` | `10px 20px` `5% 10%` `1em 2em 3em` `0 4px` `10 20 30` |
| `numberoptionalnumbertype` | `NumberOptionalNumberType` | `2` `2 3` `0.5` `3 3` `1 2` |
| `dasharraytype` | `DasharrayType` | `5 3` `10 5 2` `4` `5% 10%` `8 4 2 4` |
| `custompropertytype` | (not used — SVG has no custom props) | — |

Notes:
- **`PaintType` is STRUCTURED, not a leaf.** It is a `datatype.ebnf` production whose
  keyword arms (`none`, `currentColor`, `context-fill`, `context-stroke`) are terminals;
  only its `ColorType` and `UrlType` sub-parts scalarize. Do NOT put `painttype` in
  `leafTypes`. (See §g and §h.)
- **`TransformListType` is STRUCTURED** (full `transform.ebnf` grammar). Not a leaf.
  The same is true for `SvgPath` (path data), `Points`, `<viewBox>`,
  `<preserveAspectRatio>` — these are real sub-grammars and MUST render via the
  repeated-field-aware renderer, not as reps strings.
- **Multi-token reps caution:** because identical adjacent leaves collapse to one
  repeated proto field (PIPELINE_PORTING §6), keep list grammars (`Points`, path data,
  `KeyTimesList`) as genuine `{ }` repetitions rendered by the fixed renderer; do not
  pre-bake whole tuples as single reps strings except where unavoidable.
- `EventSymbol`, `character`, `IDREF` (animation timing) reduce to `XmlNameType` /
  `StringType` family leaves — model `EventSymbolType` (samples: `click`, `mouseenter`,
  `focus`, `beginEvent`, `endEvent`, `load`) and reuse `XmlNameType` for IDREFs.

Total open leaf `*Type` rules in `leafTypes`: **18** (number, integer, length,
percentage, length-percentage, coordinate, angle, time, color, iri, url, string,
custom-ident, xml-name, bcp47, list-of-numbers, list-of-lengths,
number-optional-number) + `dasharray` and `event-symbol` family = **~20**.

---

## d. Shared attribute groups (reconciled exact membership)

Each is one PascalCase production. Membership reconciled across all modules. Within a
group, attributes are modeled as `{ }` of an alternation of per-member attribute
productions (over-approximation per CONTEXT_SENSITIVITY Rule 2; "at most once / required"
is overlay). Value types shown as the leaf each attribute's value production resolves to.

### CoreAttributes (8 names)

| Attr | Production | Value |
|---|---|---|
| `id` | `IdAttr` | `XmlNameType` |
| `tabindex` | `TabindexAttr` | `IntegerType` |
| `autofocus` | `AutofocusAttr` | boolean-presence (`"true"\|"false"\|""`) |
| `lang` | `LangAttr` | `Bcp47Type` |
| `xml:lang` | `XmlLangAttr` | `Bcp47Type` (deprecated alias) |
| `xml:space` | `XmlSpaceAttr` | `"default" \| "preserve"` (deprecated) |
| `class` | `ClassAttr` | `StringType` (space-separated tokens) |
| `style` | `StyleAttr` | `StringType` (CSS declaration list) |
| `data-*` | `DataAttr` | open name + `StringType` value |

### ConditionalProcessingAttributes (2)

| `requiredExtensions` | `RequiredExtensionsAttr` | `StringType` (space-separated URL tokens) |
| `systemLanguage` | `SystemLanguageAttr` | `StringType` (comma-separated BCP47) |

(`requiredFeatures` EXCLUDED per 03.D5.)

### AriaAttributes (48: `role` + 47 `aria-*`)

`role` (`RoleAttr`, value = `StringType` role-token-list) plus the 47:
`aria-activedescendant aria-atomic aria-autocomplete aria-busy aria-checked
aria-colcount aria-colindex aria-colspan aria-controls aria-current aria-describedby
aria-details aria-disabled aria-dropeffect aria-errormessage aria-expanded aria-flowto
aria-grabbed aria-haspopup aria-hidden aria-invalid aria-keyshortcuts aria-label
aria-labelledby aria-level aria-live aria-modal aria-multiline aria-multiselectable
aria-orientation aria-owns aria-placeholder aria-posinset aria-pressed aria-readonly
aria-relevant aria-required aria-roledescription aria-rowcount aria-rowindex aria-rowspan
aria-selected aria-setsize aria-sort aria-valuemax aria-valuemin aria-valuenow
aria-valuetext`. Each `aria-*` → production `Aria<Name>Attr`, value = `StringType`
(`<aria-value>` open leaf; WAI-ARIA per-attr types not enumerated — versioned externally).

### PresentationAttributes (the full presentation-attribute set; 60 names)

Each maps to a CSS-property value syntax (defined in `styling.ebnf` /
`painting.ebnf` / `filter.ebnf` / `text.ebnf`). Each → production `Pascal(name)Attr`.
Length/length-pct/angle-valued ones include the §4.2 unitless-`<number>` arm.

Paint/stroke (from painting): `fill`→`FillType`(=Paint), `fill-opacity`→`<alpha-value>`,
`fill-rule`→`nonzero|evenodd`, `stroke`→Paint, `stroke-opacity`→`<alpha-value>`,
`stroke-width`→`LengthPercentageType|NumberType`, `stroke-linecap`→`butt|round|square`,
`stroke-linejoin`→`miter|miter-clip|round|bevel|arcs`, `stroke-miterlimit`→`NumberType`,
`stroke-dasharray`→`none|DasharrayType`, `stroke-dashoffset`→`LengthPercentageType|NumberType`,
`paint-order`→`normal|[fill||stroke||markers]`, `marker`→`none|UrlType`,
`marker-start`/`marker-mid`/`marker-end`→`none|UrlType`, `color`→`ColorType`,
`color-interpolation`→`auto|sRGB|linearRGB`, `color-rendering`→`auto|optimizeSpeed|optimizeQuality`,
`shape-rendering`→`auto|optimizeSpeed|crispEdges|geometricPrecision`,
`vector-effect`→(none|effect-kw+ scope-kw?).

Filter (from filter): `filter`→`none|<filter-value-list>`,
`color-interpolation-filters`→`auto|sRGB|linearRGB`, `flood-color`→`ColorType`,
`flood-opacity`→`<alpha-value>`, `lighting-color`→`ColorType`, `image-rendering`→
`auto|optimizeSpeed|optimizeQuality|crisp-edges|pixelated`.

Masking (from masking): `clip-path`→`<clip-source>|<basic-shape>|<geometry-box>|none`,
`clip-rule`→`nonzero|evenodd`, `clip`→`rect()|auto`(deprecated), `mask`→`<mask-layer>#`,
`mask-type`→(only on mask element; not a general PA).

Text/font (from text/styling): `alignment-baseline`(9 kw), `baseline-shift`→`sub|super|LengthPercentageType`,
`dominant-baseline`(9 kw), `direction`→`ltr|rtl`, `unicode-bidi`(6 kw),
`writing-mode`(5 kw + 6 legacy), `text-anchor`→`start|middle|end`,
`text-decoration`(shorthand), `text-rendering`(4 kw), `text-overflow`→`clip|ellipsis`,
`letter-spacing`→`normal|LengthType`, `word-spacing`→`normal|LengthPercentageType`,
`font-family`(open), `font-size`, `font-size-adjust`→`none|NumberType`, `font-stretch`,
`font-style`, `font-variant`, `font-weight`, `line-height`→`normal|NumberType|LengthPercentageType`.

Display/visibility/box (from styling): `display`→`none|DisplayType`,
`visibility`→`visible|hidden|collapse`, `opacity`→`<alpha-value>`,
`overflow`→`visible|hidden|scroll|auto`, `pointer-events`(10 kw + `auto` CSS alias),
`cursor`(open CSS leaf), `transform`→`TransformList`, `transform-origin`(de-facto PA).

Stop (from pservers): `stop-color`→`currentColor|ColorType`, `stop-opacity`→`<alpha-value>`.

Deprecated kept: `glyph-orientation-vertical` (alias, overlay-only). EXCLUDED:
`glyph-orientation-horizontal` (07.D9), `kerning`, `isolation`/`mix-blend-mode` (CSS-only,
no PA — 02.D4).

### GlobalEventAttributes (58 reconciled)

The spec 58-member list MINUS `ondragexit` (03.D6), `onmousewheel` (03.D7),
`onshow` (03.D9), PLUS `onwheel` (03.D7) = **57** active members. Each → `<Name>Attr`,
value = `StringType` (event-handler script). List:
`oncancel oncanplay oncanplaythrough onchange onclick onclose oncuechange ondblclick
ondrag ondragend ondragenter ondragleave ondragover ondragstart ondrop ondurationchange
onemptied onended onerror onfocus oninput oninvalid onkeydown onkeypress onkeyup onload
onloadeddata onloadedmetadata onloadstart onmousedown onmouseenter onmouseleave
onmousemove onmouseout onmouseover onmouseup onpause onplay onplaying onprogress
onratechange onreset onresize onscroll onseeked onseeking onselect onstalled onsubmit
onsuspend ontimeupdate ontoggle onvolumechange onwaiting onwheel`. (`onkeypress`
deprecated-note.)

### DocumentEventAttributes (5 — only on `svg`)

`onunload onabort onerror onresize onscroll` → `<Name>Attr` (`StringType`). On `svg`
these are Window handlers (overlap with global `onerror/onresize/onscroll` is positional;
distinct productions `SvgDocOnerrorAttr` etc. to avoid collision with global ones).

### DocumentElementEventAttributes (3)

`oncopy oncut onpaste` → `<Name>Attr` (`StringType`). (`onpaste`/`oncopy`/`oncut` absent
from `pattern` content per 00; modeled as group, presence is positional.)

### GraphicalEventAttributes (2)

`onfocusin onfocusout` → `OnfocusinAttr`, `OnfocusoutAttr` (`StringType`). Absent from
desc/title/metadata.

### AnimationEventAttributes (3 — animation elements only)

`onbegin onend onrepeat` → `<Name>Attr` (`StringType`). (Listed because animation
elements carry these instead of being full graphical elements.)

### XlinkAttributes (2 active + 5 obsolete)

Active (deprecated aliases, browsers parse): `xlink:href`→`XlinkHrefAttr`(`IriType`),
`xlink:title`→`XlinkTitleAttr`(`StringType`). Obsolete/EXCLUDED: `xlink:role`,
`xlink:arcrole`, `xlink:show`, `xlink:actuate`, `xlink:type` (08.D2).

Group sizes (for executive summary): Core **8/9**, ConditionalProcessing **2**,
Aria **48**, Presentation **60**, GlobalEvent **57**, DocumentEvent **5**,
DocumentElementEvent **3**, GraphicalEvent **2**, Xlink **2** (active).

---

## e. Content categories (membership sets used by content models)

Each is a `*Element` union production (a oneof) in `svg.ebnf`. Content models reference
these by name.

| Category | Members |
|---|---|
| `AnimationElement` | `animate animateMotion animateTransform discard set` |
| `DescriptiveElement` | `desc title metadata` |
| `ShapeElement` | `circle ellipse line path polygon polyline rect` |
| `StructuralElement` | `defs g svg symbol use` |
| `ContainerElement` | `a clipPath defs g marker mask pattern svg switch symbol unknown` |
| `GraphicsElement` | `circle ellipse foreignObject image line path polygon polyline rect text textPath tspan use` |
| `GraphicsReferencingElement` | `image use` |
| `PaintServerElement` | `linearGradient radialGradient pattern` |
| `GradientElement` | `linearGradient radialGradient` |
| `FilterPrimitiveElement` | `feBlend feColorMatrix feComponentTransfer feComposite feConvolveMatrix feDiffuseLighting feDisplacementMap feDropShadow feFlood feGaussianBlur feImage feMerge feMorphology feOffset feSpecularLighting feTile feTurbulence` |
| `TransferFunctionElement` | `feFuncR feFuncG feFuncB feFuncA` |
| `LightSourceElement` | `feDistantLight fePointLight feSpotLight` |
| `TextContentChildElement` | `tspan textPath` (a allowed in text contexts) |
| `NeverRenderedElement` | `clipPath defs desc linearGradient marker mask metadata pattern radialGradient script style symbol title` |

(audio/canvas/iframe/video EXCLUDED from every category per scope. SVG1.1 removed
elements EXCLUDED from `mask`/`clipPath` content per 12.note.)

---

## f. Per-element master table

For every included element: content model (content-category names and/or explicit child
names; CDATA = holds character data), shared attribute groups, and element-specific
attributes (name → value type). Shared groups abbreviated: **Core**=CoreAttributes,
**Cond**=ConditionalProcessing, **Aria**=Aria, **Pres**=Presentation, **GEv**=GlobalEvent,
**DEv**=DocumentEvent, **DEEv**=DocumentElementEvent, **GrEv**=GraphicalEvent,
**AEv**=AnimationEvent, **Xlink**=XlinkAttributes. Geometry attrs accept §4.2 number
expansion. "STRUCTURED" value types are sub-grammars, not leaves.

### Structure / container

**Svg** — content: AnimationElement, DescriptiveElement, PaintServerElement,
ShapeElement, StructuralElement, plus `a clipPath filter foreignObject image marker mask
script style switch text view`. Groups: Core Cond Aria Pres GEv DEv DEEv GrEv. Attrs:
`viewBox`→ViewBox(STRUCTURED), `preserveAspectRatio`→PreserveAspectRatio(STRUCTURED),
`zoomAndPan`→`disable|magnify`, `x`/`y`→LengthPercentageType, `width`/`height`→`auto|LengthPercentageType`,
`transform`→TransformList. (version/baseProfile/playbackorder/timelinebegin EXCLUDED.)

**G** — content: same as Svg's child set. Groups: Core Cond Aria Pres GEv DEEv GrEv.
Attrs: none element-specific (transform via Pres).

**Defs** — content: same as G. Groups: Core Cond Aria Pres GEv DEEv GrEv (Cond+Aria added
per 03.D1). Attrs: none.

**Symbol** — content: same as G. Groups: Core Cond Aria Pres GEv DEEv GrEv (Cond added
per 03.D2). Attrs: `viewBox`→ViewBox, `preserveAspectRatio`→PreserveAspectRatio,
`refX`→`LengthType|left|center|right`, `refY`→`LengthType|top|center|bottom`,
`x`/`y`→LengthPercentageType, `width`/`height`→`auto|LengthPercentageType`.

**Use** — content: AnimationElement, DescriptiveElement, `clipPath mask script style`.
Groups: Core Cond Aria Pres GEv DEEv GrEv Xlink. Attrs: `href`→IriType,
`x`/`y`→LengthPercentageType, `width`/`height`→`auto|LengthPercentageType`.

**Switch** — content: AnimationElement, ShapeElement, `a foreignObject g image svg switch
text use`. Groups: Core Cond Aria Pres GEv DEEv GrEv. Attrs: none. (No descriptive
children per 03.D10.)

**A** — content: DescriptiveElement + parent's content model minus `a` (over-approx:
AnimationElement, ShapeElement, ContainerElement, GraphicsElement, DescriptiveElement,
PaintServerElement). Groups: Core Cond Aria Pres GEv DEEv GrEv Xlink. Attrs: `href`→IriType,
`target`→`_self|_parent|_top|_blank|XmlNameType`, `download`→StringType, `ping`→StringType,
`rel`→StringType, `hreflang`→Bcp47Type, `type`→StringType(MIME), `referrerpolicy`→StringType.

**Unknown** — content: open (any element + CDATA; over-approx: ContainerElement child set
+ CDATA). Groups: Core Cond Aria Pres GEv DEEv GrEv. Attrs: none. (Modeled but never
emitted by generator.)

### Descriptive

**Desc / Title / Metadata** — content: open (any element + CDATA). Groups: Core GEv DEEv
ONLY (no Cond, no Aria, no GrEv per 03.D14). Attrs: none. CDATA=yes.

### Shapes

**Rect** — content: AnimationElement, DescriptiveElement, PaintServerElement, `clipPath
marker mask script style`. Groups: Core Cond Aria Pres GEv DEEv GrEv. Attrs:
`x`/`y`/`width`/`height`→LengthPercentageType (no `auto`, 06.D2), `rx`/`ry`→`auto|LengthPercentageType`,
`pathLength`→NumberType.

**Circle** — content: same as Rect. Groups: same. Attrs: `cx`/`cy`/`r`→LengthPercentageType,
`pathLength`→NumberType.

**Ellipse** — content: same. Groups: same. Attrs: `cx`/`cy`→LengthPercentageType,
`rx`/`ry`→`auto|LengthPercentageType`, `pathLength`→NumberType.

**Line** — content: same. Groups: same. Attrs:
`x1`/`y1`/`x2`/`y2`→`LengthPercentageType|NumberType` (06.D1), `pathLength`→NumberType.

**Polyline / Polygon** — content: same. Groups: same. Attrs: `points`→Points(STRUCTURED),
`pathLength`→NumberType.

**Path** — content: same as Rect. Groups: same. Attrs: `d`→`none|SvgPath`(STRUCTURED),
`pathLength`→NumberType.

### Text (CDATA=yes for all three)

**Text** — content: AnimationElement, DescriptiveElement, PaintServerElement,
TextContentChildElement, `a clipPath marker mask script style` + CDATA. Groups: Core Cond
Aria Pres GEv DEEv GrEv. Attrs: `x`/`y`/`dx`/`dy`→text-coord-list (`[ LengthPercentageType|NumberType ]+`
repeated), `rotate`→number-list, `textLength`→`LengthPercentageType|NumberType`,
`lengthAdjust`→`spacing|spacingAndGlyphs`.

**Tspan** — content: DescriptiveElement, PaintServerElement, `a animate script set style
tspan` + CDATA. Groups: Core Cond Aria Pres GEv DEEv GrEv. Attrs: same positioning set as
Text (distinct productions `TspanXAttr` etc.).

**TextPath** — content: DescriptiveElement, PaintServerElement, `a animate clipPath marker
mask script set style tspan` + CDATA. Groups: Core Cond Aria Pres GEv DEEv GrEv Xlink.
Attrs: `path`→`SvgPath`(STRUCTURED), `href`→IriType, `startOffset`→`LengthPercentageType|NumberType`,
`method`→`align|stretch`, `spacing`→`auto|exact`, `side`→`left|right`,
`textLength`→`LengthPercentageType|NumberType`, `lengthAdjust`→`spacing|spacingAndGlyphs`.

### Embedded / linking / scripting / styling

**Image** — content: AnimationElement, DescriptiveElement, `clipPath mask script style`.
Groups: Core Cond Aria Pres GEv DEEv GrEv Xlink. Attrs: `href`→IriType,
`crossorigin`→`""|anonymous|use-credentials`, `preserveAspectRatio`→PreserveAspectRatio,
`x`/`y`→LengthPercentageType, `width`/`height`→`auto|LengthPercentageType`,
`decoding`→StringType.

**ForeignObject** — content: open (any namespace + CDATA). Groups: Core Cond Aria Pres GEv
DEEv GrEv. Attrs: `x`/`y`→LengthPercentageType, `width`/`height`→`auto|LengthPercentageType`.
(No href/crossorigin/preserveAspectRatio — 08.D7.)

**View** — content: AnimationElement, DescriptiveElement, `script style`. Groups: Core
Aria GEv DEEv (no Cond — 08). Attrs: `viewBox`→ViewBox, `preserveAspectRatio`→PreserveAspectRatio,
`zoomAndPan`→`disable|magnify`. (viewTarget EXCLUDED.)

**Script** — content: CDATA only. Groups: Core GEv DEEv Xlink. Attrs: `type`→StringType(MIME),
`href`→IriType, `crossorigin`→`""|anonymous|use-credentials`, `async`/`defer`→boolean-presence
(08.D3).

**Style** — content: CDATA (raw CSS). Groups: Core GEv DEEv. Attrs: `type`→StringType,
`media`→StringType(media-query-list), `title`→StringType.

### Paint servers

**LinearGradient** — content: DescriptiveElement, `animate animateTransform set script
style stop`. Groups: Core Pres GEv DEEv Xlink. Attrs: `gradientUnits`→`userSpaceOnUse|objectBoundingBox`,
`gradientTransform`→TransformList, `x1`/`y1`/`x2`/`y2`→`LengthPercentageType|NumberType`,
`spreadMethod`→`pad|reflect|repeat`, `href`→IriType.

**RadialGradient** — content: same as LinearGradient. Groups: same. Attrs: `gradientUnits`,
`gradientTransform`, `cx`/`cy`/`r`/`fx`/`fy`/`fr`→`LengthPercentageType|NumberType`,
`spreadMethod`, `href`→IriType.

**Stop** — content: `animate set script style` ONLY (10.D9). Groups: Core Pres GEv DEEv.
Attrs: `offset`→`NumberType|PercentageType`. Pres carries `stop-color`/`stop-opacity`.

**Pattern** — content: AnimationElement, DescriptiveElement, PaintServerElement,
ShapeElement, StructuralElement, `a clipPath filter foreignObject image marker mask script
style switch text view`. Groups: Core Cond Aria Pres GEv DEEv GrEv Xlink. Attrs:
`patternUnits`/`patternContentUnits`→`userSpaceOnUse|objectBoundingBox`,
`patternTransform`→TransformList, `x`/`y`/`width`/`height`→LengthType (no %, 10.D7),
`viewBox`→ViewBox, `preserveAspectRatio`→PreserveAspectRatio, `href`→IriType.

### Painting

**Marker** — content: AnimationElement, DescriptiveElement, PaintServerElement,
ShapeElement, StructuralElement, `a clipPath filter foreignObject image marker mask script
style switch text view`. Groups: Core GEv DEEv Pres. Attrs:
`markerUnits`→`strokeWidth|userSpaceOnUse`, `markerWidth`/`markerHeight`→`LengthPercentageType|NumberType`,
`refX`→`LengthPercentageType|NumberType|left|center|right`,
`refY`→`LengthPercentageType|NumberType|top|center|bottom`,
`orient`→`auto|auto-start-reverse|AngleType|NumberType`, `viewBox`→ViewBox,
`preserveAspectRatio`→PreserveAspectRatio.

### Masking / clipping

**ClipPath** — content: DescriptiveElement, AnimationElement, ShapeElement, `text use
script`. Groups: Core Cond Aria Pres GEv. Attrs:
`clipPathUnits`→`userSpaceOnUse|objectBoundingBox` (default userSpaceOnUse), `transform`→TransformList.

**Mask** — content: AnimationElement, DescriptiveElement, ShapeElement, StructuralElement,
GradientElement, `a clipPath filter foreignObject image marker mask pattern script style
switch text view` (SVG1.1-removed elements EXCLUDED, 12.note). Groups: Core Cond Aria Pres
GEv DEEv. Attrs: `maskUnits`→`userSpaceOnUse|objectBoundingBox` (default objectBoundingBox),
`maskContentUnits`→`userSpaceOnUse|objectBoundingBox` (default userSpaceOnUse),
`x`/`y`/`width`/`height`→LengthPercentageType, `mask-type`→`luminance|alpha`.

### Filters

**Filter** — content: DescriptiveElement, FilterPrimitiveElement, `animate script set`.
Groups: Core Pres GEv DEEv Xlink. Attrs: `filterUnits`→`userSpaceOnUse|objectBoundingBox`,
`primitiveUnits`→`userSpaceOnUse|objectBoundingBox`, `x`/`y`/`width`/`height`→LengthPercentageType,
`href`→IriType. (filterRes EXCLUDED.)

**Common filter-primitive attrs** (each primitive except feMergeNode gets these as
positional productions): `x`/`y`/`width`/`height`→LengthPercentageType, `result`→CustomIdentType.
`in`/`in2`→FeInValue (`SourceGraphic|SourceAlpha|BackgroundImage|BackgroundAlpha|FillPaint|StrokePaint|CustomIdentType`).

| Primitive | Content | `in`? | Element-specific attrs |
|---|---|---|---|
| `feBlend` | DescriptiveElement, animate/script/set | in, in2 | `mode`→16-keyword `<blend-mode>`, `no-composite`→boolean |
| `feColorMatrix` | same | in | `type`→`matrix|saturate|hueRotate|luminanceToAlpha`, `values`→ListOfNumbersType |
| `feComponentTransfer` | DescriptiveElement, TransferFunctionElement, script | in | (children carry params) |
| `feFuncR/G/B/A` | DescriptiveElement, animate/script/set | — | `type`→`identity|table|discrete|linear|gamma`, `tableValues`→ListOfNumbersType, `slope`/`intercept`/`amplitude`/`exponent`/`offset`→NumberType |
| `feComposite` | same | in, in2 | `operator`→`over|in|out|atop|xor|lighter|arithmetic`, `k1`/`k2`/`k3`/`k4`→NumberType |
| `feConvolveMatrix` | same | in | `order`→NumberOptionalNumberType, `kernelMatrix`→ListOfNumbersType, `divisor`/`bias`→NumberType, `targetX`/`targetY`→IntegerType, `edgeMode`→`duplicate|wrap|none` (default duplicate), `kernelUnitLength`→NumberOptionalNumberType, `preserveAlpha`→`false|true` |
| `feDiffuseLighting` | DescriptiveElement, script, exactly-one LightSourceElement | in | `surfaceScale`/`diffuseConstant`→NumberType, `kernelUnitLength`→NumberOptionalNumberType |
| `feDisplacementMap` | DescriptiveElement, animate/script/set | in, in2 | `scale`→NumberType, `xChannelSelector`/`yChannelSelector`→`R|G|B|A` |
| `feDropShadow` | same | in | `dx`/`dy`→NumberType, `stdDeviation`→NumberOptionalNumberType (flood-color/opacity via Pres) |
| `feFlood` | same | — | (flood-color/flood-opacity via Pres) |
| `feGaussianBlur` | same | in | `stdDeviation`→NumberOptionalNumberType, `edgeMode`→`duplicate|wrap|none` (default none, 11.D3) |
| `feImage` | DescriptiveElement, animate/animateTransform/script/set | — | `href`→IriType, `preserveAspectRatio`→PreserveAspectRatio, `crossorigin`→`anonymous|use-credentials`, Xlink |
| `feMerge` | DescriptiveElement, feMergeNode, script | — | (children only) |
| `feMergeNode` | DescriptiveElement, animate/script/set | in | (NO x/y/w/h/result — 11.D14) |
| `feMorphology` | DescriptiveElement, animate/script/set | in | `operator`→`erode|dilate`, `radius`→NumberOptionalNumberType |
| `feOffset` | same | in | `dx`/`dy`→NumberType |
| `feSpecularLighting` | DescriptiveElement, script, exactly-one LightSourceElement | in | `surfaceScale`/`specularConstant`/`specularExponent`→NumberType, `kernelUnitLength`→NumberOptionalNumberType |
| `feTile` | DescriptiveElement, animate/script/set | in | (none) |
| `feTurbulence` | same | — | `baseFrequency`→NumberOptionalNumberType, `numOctaves`→IntegerType, `seed`→NumberType, `stitchTiles`→`stitch|noStitch`, `type`→`fractalNoise|turbulence` |
| `feDistantLight` | DescriptiveElement, animate/script/set | — | `azimuth`/`elevation`→NumberType |
| `fePointLight` | same | — | `x`/`y`/`z`→NumberType |
| `feSpotLight` | same | — | `x`/`y`/`z`/`pointsAtX`/`pointsAtY`/`pointsAtZ`/`specularExponent`→NumberType, `limitingConeAngle`→NumberType |

All filter-primitive elements also get Core + Pres (color-interpolation-filters etc.).

### Animation

**Animate** — content: DescriptiveElement, script. Groups: Core Cond AEv GEv Pres. Attrs:
`href`→IriType, `attributeName`→XmlNameType, timing (`begin`/`end`→BeginValueList,
`dur`/`repeatDur`→`Clock|media|indefinite`, `min`/`max`→`Clock|media`, `restart`→`always|whenNotActive|never`,
`repeatCount`→`NumberType|indefinite`, `fill`→`freeze|remove` as `AnimateFillAttr`),
value (`calcMode`→`discrete|linear|paced|spline`, `values`→AnimationValueList,
`keyTimes`→KeyTimesList, `keySplines`→KeySplinesList, `from`/`to`/`by`→AnimationValue),
addition (`additive`→`replace|sum`, `accumulate`→`none|sum`).

**Set** — content: DescriptiveElement, script. Groups: Core Cond AEv GEv Pres. Attrs:
`href`, `attributeName`, timing set (with `SetFillAttr`), `to`→AnimationValue. (No
additive/accumulate/calcMode/values/keyTimes/keySplines/from/by.)

**AnimateMotion** — content: DescriptiveElement, script, at-most-one `mpath`. Groups: Core
Cond AEv GEv Pres. Attrs: `href`, timing, addition, `calcMode` (default paced),
`values`/`from`/`to`/`by`→MotionValues/MotionCoordinatePair, `keyTimes`→KeyTimesList,
`keySplines`→KeySplinesList, `path`→SvgPath, `keyPoints`→KeyPointsList,
`rotate`→`NumberType|auto|auto-reverse`, `origin`→`default`. (No attributeName.)

**AnimateTransform** — content: DescriptiveElement, script. Groups: Core Cond AEv GEv Pres.
Attrs: `href`, `attributeName`, timing, addition, value (`from`/`to`/`by`/`values`→
AnimateTransformSingleValue/AnimateTransformValues), `type`→`translate|scale|rotate|skewX|skewY`.

**Mpath** — content: DescriptiveElement, script. Groups: Core GEv. Attrs: `href`→IriType
(required), Xlink.

**Discard** — content: DescriptiveElement, script. Groups: Core. Attrs: `href`→IriType,
`begin`→BeginValueList. (Experimental flag.)

---

## g. Cross-module conflict resolutions

Each value/property defined in more than one note → single canonical home; all other
modules reference it.

| Value/property | Canonical home | Referenced by |
|---|---|---|
| `<number>`/`<integer>`/`<length>`/`<percentage>`/`<length-percentage>`/`<angle>`/`<color>`/`<url>`/`<iri>` leaves | `datatype.ebnf` | every module |
| `<paint>` (STRUCTURED: `none\|currentColor\|ColorType\|UrlType[fallback]?\|context-fill\|context-stroke`) | `datatype.ebnf` (`PaintType`) | painting (fill/stroke), filter (text-decoration-fill/stroke), animation (PaintType arm) |
| `<color>` arms (named/hex/rgb/hsl/currentColor) | `datatype.ebnf` (`ColorType`) | painting, pservers (stop-color), filter (flood/lighting), text (text-decoration-color) |
| `<transform-list>` (SVG attr form) | `transform.ebnf` (`TransformList`) | structure (`transform`), pservers (`gradientTransform`/`patternTransform`), masking (`clipPath transform`), animation (TransformListType arm) |
| path data (`svg_path`) | `path.ebnf` (`SvgPath`) | shapes/path (`d`), text (`textPath path`), animation (`animateMotion path`, PathDataType arm), masking (`path()` basic-shape) |
| `preserveAspectRatio` | `datatype.ebnf` (`PreserveAspectRatio`) | svg, symbol, marker, pattern, image, view, feImage |
| `viewBox` | `datatype.ebnf` (`ViewBox`) | svg, symbol, marker, pattern, view |
| geometry props `cx cy r rx ry x y width height` | each is a **distinct positional per-element production** (`RectXAttr`, `CircleCxAttr`, `RadialGradientCxAttr`, `FePointLightXAttr`, …); the value *leaf* (`LengthPercentageType`) is canonical in datatype | shapes, structure, pservers, filters, masking — all reference the leaf, never share the attr production |
| `gradientUnits`/`patternUnits`/`maskUnits`/`clipPathUnits`/`filterUnits`/`maskContentUnits`/`patternContentUnits`/`primitiveUnits` | `datatype.ebnf` shared enum `UnitType = "userSpaceOnUse" \| "objectBoundingBox"`; each attr is a positional production referencing it (defaults differ, recorded in overlay) | pservers, masking, filters |
| `spreadMethod` | `pservers.ebnf` | linearGradient, radialGradient |
| `stop-color`/`stop-opacity` | `pservers.ebnf` | stop element only |
| `flood-color`/`flood-opacity`/`lighting-color` | `filter.ebnf` | feFlood/feDropShadow (flood-*), feDiffuse/feSpecularLighting (lighting-color) |
| `color-interpolation-filters` (initial linearRGB) | `filter.ebnf` | all filter primitives (Pres) |
| `color-interpolation` (initial sRGB) | `styling.ebnf`/`painting.ebnf` | containers, graphics, gradients — **distinct production** from color-interpolation-filters (initials differ, 02.D9) |
| `<alpha-value>` (`NumberType\|PercentageType`) | `datatype.ebnf` | opacity, fill-/stroke-/stop-/flood-opacity |
| `href`/`xlink:href` | each element gets positional `HrefAttr`-style productions referencing `IriType`; `XlinkHrefAttr`/`XlinkTitleAttr` are the shared deprecated-alias productions in `svg.ebnf` (XlinkAttributes) | use, image, a, script, textPath, feImage, gradients, pattern, mpath, animation, discard |
| `pointer-events` (10 kw + auto alias) | `embedded.ebnf` | Pres |
| `fill-rule`/`clip-rule` (`nonzero\|evenodd`) | `fill-rule` in painting; `clip-rule` in masking (distinct productions, same value set) | painting (fill-rule), masking (clip-rule), basic-shape functions |
| `<fe-in-value>` | `filter.ebnf` (`FeInValue`) | all filter primitives with `in`/`in2` |
| SMIL `Clock-value`, timing lists | `animation.ebnf` | begin/end/dur/min/max/repeatDur |

Rule of thumb: **value leaves and true sub-grammars are shared singletons; attribute
productions are per-element positional** (because attribute names collide across elements
and gluon would dedupe shared names — §b).

---

## h. AnimationValue union

CFG over-approximation (CONTEXT_SENSITIVITY Rule 2 / §6). The grammar carries the union;
the overlay narrows to the arm matching the resolved `attributeName` type.

```ebnf
AnimationValue =
    LengthType            (* stroke-width, cx, cy, r, x, y, width, height *)
  | NumberType            (* opacity, stroke-miterlimit, fill-opacity *)
  | PercentageType
  | LengthPercentageType
  | AngleType
  | ColorType             (* fill, stroke as color *)
  | PaintType             (* fill, stroke as paint: url/none/context-* *)
  | TransformList         (* animateTransform *)
  | SvgPath               (* d attribute *)
  | IntegerType
  | StringType            (* non-interpolable attrs *)
  | IriType               (* href-like *)
  | KeywordType           (* visibility, display, etc. enum keywords *)
  ;

AnimationValueList = AnimationValue, { S, ";", S, AnimationValue }, [ S, ";" ] ;
```

`animateMotion` uses `MotionCoordinatePair`/`MotionValues` (x,y pairs), NOT this union.
`animateTransform` uses the per-type specialized productions
(`AnimateTransformTranslateValue` etc.) unioned in `AnimateTransformSingleValue`; the
`type` attribute (closed set of 5) selects the arm structurally — already context-free.

**Overlay note (one line):** the type oracle is generated from the attribute grammar —
`attributeName="stroke-width"` resolves to `StrokeWidthAttr`'s value production
(`LengthPercentageType | NumberType`) and the overlay selects the matching `AnimationValue`
arm; in generation mode the generator picks `attributeName` first, then emits a matching
`from`/`to`/`values`. All cardinality (keyTimes/values/keySplines counts), range, and
mutual-exclusion (`values` vs `from`/`to`/`by`) checks are overlay, not grammar.

---

## Executive summary (20 lines)

1. Root production: **`SvgDocument`** (`SvgDocument = Svg ;`), the prune root for genproto.
2. Files (14, concatenation order): `svg.ebnf` (ROOT) → `datatype` → `transform` → `path`
   → `structure` → `shapes` → `text` → `embedded` → `painting` → `pservers` → `filter` →
   `masking` → `animation` → `styling`.
3. Open leaf `*Type` rules in the scalarize `leafTypes` set: **~20** (number, integer,
   length, percentage, length-percentage, coordinate, angle, time, color, iri, url,
   string, custom-ident, xml-name, bcp47, list-of-numbers, list-of-lengths,
   number-optional-number, dasharray, event-symbol). `PaintType`/`TransformList`/`SvgPath`/
   `Points`/`ViewBox`/`PreserveAspectRatio` are STRUCTURED, not leaves.
4. Attribute-group sizes (9): Core **8** (+`data-*`), ConditionalProcessing **2**,
   Aria **48**, Presentation **60**, GlobalEvent **57** (reconciled), DocumentEvent **5**,
   DocumentElementEvent **3**, GraphicalEvent **2**, Xlink **2** active (+5 obsolete excluded).
5. Elements INCLUDED: **66** distinct element productions (46 SVG-2 authorable surface +
   17 filter primitives' supporting children fully spelled out); EXCLUDED: **~30**
   (audio/canvas/iframe/video, SVG1.1 fonts/altGlyph*/glyph*/tref/cursor-element/
   color-profile/hkern/vkern/missing-glyph/font-face*, animateColor, and removed attrs
   playbackorder/timelinebegin/requiredFeatures/attributeType/defer/viewTarget/filterRes/
   version/baseProfile).
6. Top 5 cross-module conflicts resolved: (a) **`fill` paint-vs-animation-freeze** →
   positional `FillAttr` vs `AnimateFillAttr`; (b) **geometry props `cx/cy/r/rx/ry/x/y/
   width/height`** spread across datatypes/shapes/structure/pservers/filters → per-element
   positional `*Attr` productions sharing the `LengthPercentageType` leaf; (c)
   **`transform`/`gradientTransform`/`patternTransform`** → one shared `TransformList`
   grammar, positional attr productions; (d) **`color-interpolation` (sRGB) vs
   `color-interpolation-filters` (linearRGB)** → distinct productions, same keyword set,
   different canonical homes/initials; (e) **`href`/`xlink:href`** everywhere → per-element
   `IriType`-valued positional productions + shared `XlinkHrefAttr` deprecated alias.
7. Collision policy: gluon dedupes on `norm(s)=lower,underscores-stripped`; every shared
   attribute name is given a fully-qualified per-element production name to prevent silent
   merging — the single biggest authoring hazard.
8. `<paint>` is STRUCTURED: keyword arms (none/currentColor/context-fill/context-stroke)
   are terminals; only the `ColorType`/`UrlType` sub-parts scalarize.
9. AnimationValue is a 13-arm CFG over-approximation narrowed by an overlay type-oracle
   generated from the attribute grammar; animateTransform/animateMotion use specialized
   value productions instead.
