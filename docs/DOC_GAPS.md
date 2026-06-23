# proto-svg Doc Gaps & Discrepancy Record

Consolidated record of every spec discrepancy, doc gap, and roadblock found across
the 14 module grammar-notes (`docs/grammar-notes/00`–`13`), each with: what the spec
says, what is actually correct (per MDN / browser / verification), and the decision
taken for the grammar. This is the authoritative scope and resolution record.
Authoring (`GRAMMAR_PLAN.md`) follows the decisions recorded here verbatim.

Conventions: **SPEC** = what the source spec text says; **REAL** = verified correct
behavior (MDN / Chrome 125 / Firefox 127 / Safari 17); **DECISION** = grammar action.

---

## Scope decision

This is the final, authoritative element/attribute scope for the grammar.
Target = **SVG 2 + Filter Effects 1 + CSS Masking 1 + SMIL (SVG Animations L2)**,
restricted to **browser-supported** features, **including** the `xlink:*` deprecated
aliases that browsers still parse.

### Elements to INCLUDE (46)

Structural / container (8): `svg`, `g`, `defs`, `symbol`, `use`, `switch`, `a`, `unknown`
Descriptive (3): `desc`, `title`, `metadata`
Shapes (6): `rect`, `circle`, `ellipse`, `line`, `polyline`, `polygon`
Path (1): `path`
Text (3): `text`, `tspan`, `textPath`
Embedded / linking / scripting (4): `image`, `foreignObject`, `view`, `script`, plus styling `style` (1) → counts: image, foreignObject, view, script, style = 5
Paint servers (4): `linearGradient`, `radialGradient`, `stop`, `pattern`
Painting (1): `marker`
Masking / clipping (2): `clipPath`, `mask`
Filter container (1): `filter`
Filter primitives (17): `feBlend`, `feColorMatrix`, `feComponentTransfer`, `feComposite`,
`feConvolveMatrix`, `feDiffuseLighting`, `feDisplacementMap`, `feDropShadow`, `feFlood`,
`feGaussianBlur`, `feImage`, `feMerge`, `feMorphology`, `feOffset`, `feSpecularLighting`,
`feTile`, `feTurbulence`
Filter transfer-function children (4): `feFuncR`, `feFuncG`, `feFuncB`, `feFuncA`
Filter light-source children (3): `feDistantLight`, `fePointLight`, `feSpotLight`
Filter merge child (1): `feMergeNode`
Animation (6): `animate`, `set`, `animateMotion`, `animateTransform`, `mpath`, `discard`

Total distinct included element names: **66** (the 46 "top-level authorable" count in
00-indices counts only the SVG-2 Appendix-F surface and folds the 17 filter primitives +
their 8 children + the 4 paint/mask/marker children differently; the grammar must spell
out **all 66** element productions). `discard` is **INCLUDE (experimental flag)** — it
is requested in scope, has partial Chrome support, and is cheap to model; flagged
low-confidence in the overlay.

### Elements to EXCLUDE (with one-line reasons)

| Element | Reason |
|---|---|
| `audio` | SVG 2 proposed HTML-media embedding; no browser shipped it — unimplemented |
| `canvas` | SVG 2 proposed HTML-canvas embedding; no browser shipped it — unimplemented |
| `iframe` | SVG 2 proposed embedding; no browser shipped it — unimplemented |
| `video` | SVG 2 proposed HTML-media embedding; no browser shipped it — unimplemented |
| `unknown` (as authorable) | parse-only fallback meta-element, no real author surface — but see note* |
| `altGlyph`, `altGlyphDef`, `altGlyphItem` | SVG 1.1 glyph substitution — removed & unsupported |
| `glyph`, `glyphRef`, `missing-glyph` | SVG 1.1 fonts — removed & unsupported |
| `font`, `font-face`, `font-face-format`, `font-face-name`, `font-face-src`, `font-face-uri` | SVG 1.1 font embedding — removed & unsupported |
| `hkern`, `vkern` | SVG 1.1 kerning — removed & unsupported |
| `tref` | text reference, removed ~2013 — removed & unsupported |
| `cursor` (element) | SVG 1.1 cursor element — removed; CSS `cursor` used instead |
| `color-profile` | SVG 1.1 — removed & unsupported |
| `animateColor` | removed in SVG Animations L2 (§7); deprecated in SVG 1.1 — unsupported |

\* **`unknown`**: 00/03-notes flag it EXCLUDE as a real author element, but 03-structure
gives it a full container content-model and `SVGUnknownElement`. **DECISION:** keep a
minimal `unknown` production (container content model, core+aria+presentation attrs)
so the grammar is structurally complete, but it is never emitted by the generator. It
is counted in the INCLUDE list above. (Resolves the 00-vs-03 disagreement positively
toward inclusion since modeling cost is near-zero and it appears in the SVG 2 index.)

### Attributes / features to EXCLUDE

| Feature | Reason |
|---|---|
| `playbackorder`, `timelinebegin` (on `svg`) | SVG 2 / SVGwg-L2 additions, not implemented in any browser |
| `requiredFeatures` | SVG 1.1 conditional processing, removed from SVG 2 |
| `attributeType` (on animation) | removed in SVG Animations L2 §7 (auto-resolution always) |
| `defer` (in preserveAspectRatio) | SVG 1.1 image prefix, removed in SVG 2, no effect |
| `viewTarget` (on `view`) | SVG 1.1, explicitly resolved-to-remove (Paris 2015 F2F) |
| `version`, `baseProfile` (on `svg`) | SVG 1.1 conformance signals, ignored in SVG 2 |
| `filterRes` (on `filter`) | SVG 1.1, removed in Filter Effects 1 |
| `kerning`, `glyph-orientation-horizontal` (props) | removed in SVG 2 |
| `xlink:role`, `xlink:arcrole`, `xlink:show`, `xlink:actuate`, `xlink:type` | XLink 1.0, not in SVG 2, never implemented for SVG |
| `xml:space` | deprecated; **kept** as parseable terminal (closed `default|preserve`) but superseded by `white-space` |
| `xml:lang` | deprecated alias of `lang`; **kept** as parseable alias |
| `requiredFeatures`, `externalResourcesRequired` | removed/no-op in SVG 2 |
| `media-marker-value` (timing) | SMIL prose only; not in SVG timing grammar |
| `wallclock` timing | syntactically INCLUDED (context-free) but no browser implements it |
| `fetchpriority`, `font-width` (MDN 🧪🔶) | non-standard, not cross-browser |

### Features INCLUDED despite caveats

`xlink:href`, `xlink:title` (deprecated aliases, still parsed by all browsers);
`fr` on radialGradient (SVG 2, universally supported); `href` bare form everywhere;
`zoomAndPan` (deprecated, closed set, parseable); `discard` (experimental); `side` on
textPath (experimental); `context-fill`/`context-stroke` paint; `stroke-linejoin: arcs`/
`miter-clip` (limited support); `inline-size` (SVG 2 text wrapping, grammar-present,
limited support); `wallclock` timing (grammar-present, never rendered).

**CSS-only properties (no SVG attribute form — correctly absent from the grammar).**
`line-height`, `shape-inside`, `shape-subtract`, `shape-padding`, `shape-margin`,
`isolation`, and `mix-blend-mode` apply to SVG via CSS but are NOT in the SVG 2
attribute index G.2, i.e. they have no presentation-attribute form. They can only be
set through `style="…"`/a stylesheet, never as an XML attribute, so the structural
grammar (which models XML attributes) correctly does not express them. (An earlier note
mislabeled the `shape-*` set as grammar-present; only `inline-size` is.)

---

## Module 00 — Indices

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 00.1 | SVG2 index lists 52 elements incl. audio/canvas/iframe/video/unknown/discard | None of audio/canvas/iframe/video ship as SVG elements; unknown is parse-only; discard partial | INCLUDE 46-surface; EXCLUDE audio/canvas/iframe/video; keep unknown minimal; INCLUDE discard (experimental) |
| 00.2 | G.1 omits geometry attrs `d`,`cx`,`cy`,`r`,`rx`,`ry`,`x`,`y`,`width`,`height` for shapes | These are defined per-element in Shapes/Structure chapters, fully supported | INCLUDE from their normative chapters (positional per-element productions) |
| 00.3 | G.1 omits `mask-type`, `transform-origin`, `decoding`, `text-overflow` | All browser-supported (CSS Masking / CSS Transforms / HTML / CSS UI) | INCLUDE; see styling/masking modules for canonical homes |
| 00.4 | G.1 includes `playbackorder`, `timelinebegin` on `svg` | Not implemented in any browser | EXCLUDE |
| 00.5 | `attributeType` deprecated (MDN ⚠️) | Removed in SVG Animations L2; browsers parse-but-ignore | EXCLUDE (see 13.D1) |
| 00.6 | `xlink:href`/`xlink:title`/`xml:space`/`xml:lang` deprecated | Still parsed by all browsers | INCLUDE as deprecated aliases |
| 00.7 | `autofocus`, `data-*`, `tabindex` are HTML globals on SVG | Valid on SVG elements in browsers | INCLUDE: `tabindex` and `autofocus` in core; `data-*` modeled as open `DataAttr` |

---

## Module 01 — Datatypes & Geometry (D1–D10)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 01.D1 | Ch.4 is DOM-only; no EBNF for `<number>`/`<length>`/`<color>`/`<paint>`/`<transform-list>`/`<IRI>` | Defined in CSS Values/Color/Transforms, RFC 3987, URL Standard | Reconstruct productions from those specs; scalarize the open leaves |
| 01.D2 | SVG2 removes unitless `<length>`; handled via §4.2 presentation-attr expansion | Browsers accept unitless numbers in presentation attrs | Provide two leaves: `LengthType` (CSS-strict) and the §4.2 expansion via attribute productions accepting `LengthType \| PercentageType \| NumberType` |
| 01.D3 | `r` = `<length-percentage>`; SVG 1.1 allowed unitless | §4.2 expansion covers unitless; `r` has no `auto` | `r` value = length-percentage with number expansion; no `auto` |
| 01.D4 | §7.8 defers width/height to CSS 2.1 (no min/max/fit-content) | Browsers accept CSS Sizing L3 values | width/height value = `auto \| LengthPercentageType` plus `min-content \| max-content \| fit-content(...)`; SVG `auto` semantics are overlay |
| 01.D5 | `cx`/`cy` geometry props apply only to circle/ellipse | Also exist as different attrs on radialGradient (bbox %) / fePointLight (number) | Distinct **positional** productions per element (geometry vs gradient vs light) |
| 01.D6 | `<dasharray>` not in Ch.4 | Defined in Painting ch.; `[ length-pct \| number ]#*`, non-negative | Canonical home = painting module (09); `DasharrayType` leaf |
| 01.D7 | `<transform-list>` delegated to CSS Transforms (adds 3D fns) | SVG attribute form is bare-number 6-function set | Two grammars: `TransformList` (SVG attr, bare numbers, 6 fns) vs CSS `transform` property; SVG attr is the authored one |
| 01.D8 | `turn` unit not in SVGAngle DOM constants | Browsers accept `turn` in attributes | INCLUDE `turn` in `AngleType` (DOM limitation is not a parse concern) |
| 01.D9 | SVG 1.1 opacity = number [0,1]; CSS adds `<percentage>` | Both forms valid | `<alpha-value>` = `NumberType \| PercentageType`; clamp is overlay |
| 01.D10 | `zoomAndPan` deprecated (MDN) but in SVG 2 DOM | Browsers removed interactive magnify | INCLUDE closed set `disable \| magnify`; deprecated flag overlay |

---

## Module 02 — Styling & Rendering (D1–D12)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 02.D1 | `color-rendering` in §6.6 PA table | Not a standalone MDN attribute entry | INCLUDE; closed set `auto\|optimizeSpeed\|optimizeQuality`; note MDN doc gap |
| 02.D2 | `transform-origin` required (§6.7) but NOT in §6.6 PA table | Browsers accept it as a presentation attribute | INCLUDE as de-facto PA (open value via CSS) with note |
| 02.D3 | `overflow` values = visible/hidden/scroll/auto; no `clip` | Modern browsers accept `overflow: clip` on SVG | EXCLUDE `clip` for strict SVG2 (note as browser extension) |
| 02.D4 | `isolation`/`mix-blend-mode` required but no PA | CSS-only (no XML attribute) | No PA productions; CSS-only (out of attribute grammar) |
| 02.D5 | `display` referenced as CSS2 | Browsers apply CSS Display L3 multi-keyword | `DisplayType` = open leaf; only `none` is a named SVG terminal in semantics |
| 02.D6 | `glyph-orientation-*` in §6.6, deprecated | Broadly ignored in SVG2 mode | `glyph-orientation-horizontal` REMOVED (text 07.D9); `-vertical` deprecated alias (4 mapped tokens), overlay-only |
| 02.D7 | §3.7.1 prose "marker symbols"; keyword spelling | CSS keyword is `markers` (plural) per MDN/browsers | `paint-order-layer` = `fill\|stroke\|markers` |
| 02.D8 | `vector-effect` = none/non-scaling-stroke (stable) | Chromium experimental non-scaling-size | See 04.D5: INCLUDE all 4 effect kws + 2 scope kws, flag at-risk |
| 02.D9 | `color-interpolation-filters` initial `linearRGB` vs `color-interpolation` `sRGB` | Intentional, correct | Not a discrepancy; both share keyword set, distinct **positional** productions |
| 02.D10 | `shape-rendering: crispEdges` (camel) vs `image-rendering: crisp-edges` (kebab) | Spec inconsistency, both real | Reproduce both faithfully as distinct terminals |
| 02.D11 | `mask-type` in MDN, not in §6.6 PA list | CSS Masking property on `mask` element | Canonical home = masking (12); model as `mask` element attr `luminance\|alpha` |
| 02.D12 | `fill` PA "except animation elements" | animation `fill` = `freeze\|remove` | **Positional resolution**: `FillAttr` (paint) vs `AnimateFillAttr` (freeze\|remove) — distinct per-element productions, no overlay |

---

## Module 03 — Structure (D1–D14, R1–R4)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 03.D1 | `defs` table omits aria + conditional-processing attrs | Likely authoring omission; browsers accept them | INCLUDE aria + conditional-processing on `defs` (match `g`) |
| 03.D2 | `symbol` table omits conditional-processing attrs | Likely authoring gap | INCLUDE conditional-processing on `symbol` |
| 03.D3 | `zoomAndPan` not marked deprecated in spec | MDN ⚠️; Chrome removed interactive magnify | INCLUDE `disable\|magnify`; deprecated note |
| 03.D4 | `version`/`baseProfile` only in prose examples | MDN ⚠️ deprecated | EXCLUDE from grammar (SVG 1.1 conformance signals) |
| 03.D5 | `requiredFeatures` removed | MDN ⚠️ | EXCLUDE |
| 03.D6 | `ondragexit` in global-event list | Non-standard alias of `ondragleave`; HTML removed it | EXCLUDE `ondragexit`; keep `ondragleave` (see GlobalEventAttributes reconciliation) |
| 03.D7 | `onmousewheel` in list, no `onwheel` | `onmousewheel` non-standard; browsers use `onwheel` | INCLUDE `onwheel`; EXCLUDE `onmousewheel` (de-bloat; both rare) → see reconciled GlobalEventAttributes |
| 03.D8 | `onkeypress` in list | DOM-deprecated, still fires | INCLUDE (deprecated note) |
| 03.D9 | `onshow` in list | `show` event removed from HTML | EXCLUDE `onshow` (vestigial) |
| 03.D10 | `switch` content model excludes desc/title/metadata | Per spec; narrowest container | Follow spec (switch excludes descriptive + most containers) |
| 03.D11 | `use` content model | DOES allow descriptive elements (spec re-read) | `use` content = animation + descriptive + clipPath/mask/script/style |
| 03.D12 | `xml:lang` in narrative, inconsistent in tables | MDN ⚠️ deprecated | INCLUDE as deprecated alias of `lang` in core group |
| 03.D13 | `onfocusin`/`onfocusout` = "graphical event attributes" | Distinct bubbling semantics from focus/blur | Separate `GraphicalEventAttributes` group |
| 03.D14 | desc/title omit aria; §5.13.4 "no role" | Confirmed: no ARIA on desc/title | desc/title/metadata carry NO aria attrs |
| 03.R1 | onfocusin/out browser SVG support inconsistent | — | INCLUDE per spec; testing flag |
| 03.R2 | `transform` on `svg` new in SVG 2 | Uneven support | INCLUDE; compat flag |
| 03.R3 | `refX`/`refY` on `symbol` new | Chrome yes, others vary | INCLUDE `<length>\|left\|center\|right` / `top\|center\|bottom` |
| 03.R4 | `use` shadow-DOM model | Browsers moved to shadow DOM | No grammar impact |

**Reconciliation note (event handler de-bloat):** 03 reasons EXCLUDE `onshow`,
`ondragexit`, `onmousewheel` and ADD `onwheel`. The canonical `GlobalEventAttributes`
in GRAMMAR_PLAN §d takes those decisions: the reconciled list drops `ondragexit`,
`onmousewheel`, `onshow` and adds `onwheel` to the spec's 58-member list.

---

## Module 04 — Coords & Transforms (D1–D8)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 04.D1 | SVG `transform` attr ≠ CSS `transform` property (functions, units, separators) | SVG attr: 6 functions, bare numbers, comma-wsp | **Two rules**: `TransformList` (SVG attr — authored) vs CSS property (not authored). gradientTransform/patternTransform reuse `TransformList` |
| 04.D2 | SVG2 drops SVG1.1 inline transform BNF, points at CSS Transforms §6.2 | Grammar identical, only location moved | Transcribe CSS Transforms §6.2 BNF verbatim into `transform.ebnf` |
| 04.D3 | `defer` removed in SVG 2 | No browser implements | EXCLUDE `defer` from preserveAspectRatio |
| 04.D4 | transform-origin initial differs for SVG (0 0 vs 50% 50%) | UA-stylesheet rule, not grammar | Grammar value identical to CSS; initial-value diff is overlay |
| 04.D5 | vector-effect at-risk keywords | Only `none`/`non-scaling-stroke` implemented | INCLUDE all 4 effect kws + scope kws (`viewport\|screen`); flag at-risk |
| 04.D6 | "SVG follows CSS units" applies to length attrs, NOT transform attr | transform uses bare numbers | Keep transform number-only; length attrs use full unit set |
| 04.D7 | No normative viewBox BNF (only `[min-x,? …]`) | comma-wsp separator | Author `ViewBox = number comma-wsp number comma-wsp number comma-wsp number`; non-neg w/h is overlay |
| 04.D8 | gradientTransform/patternTransform no separate grammar | Reuse `transform-list` | Reuse `TransformList` (positional `GradientTransformAttr`/`PatternTransformAttr`) |

---

## Module 05 — Paths (D1–D7)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 05.D1 | Path BNF `number ::= [0-9]+` (integer-only) | Browsers accept decimal fractions | **Correct the BNF**: `number = integer ["." integer] \| "." integer`; no exponent in path data |
| 05.D2 | `sign` separate from `number`; arc rotation typed `number` (unsigned) | Browsers accept signed rotation | `x-axis-rotation` parsed as `coordinate` (sign-optional) |
| 05.D3 | Mandatory `comma_wsp` (not `?`) after rotation, before first flag | Correct & intentional (flag disambiguation) | Keep mandatory separator before `large-arc-flag` |
| 05.D4 | `svg_path ::= wsp* moveto? (moveto drawto*)?` looks odd | Subpaths via embedded movetos in drawto list | Keep as-is; consistent with prose |
| 05.D5 | `d` promoted to CSS geometry property `none \| <string>` | path data is the string content | `d` value = `none \| path-data`; same `svg_path` grammar |
| 05.D6 | No exponent in path-data number | Confirmed by MDN | Path `number` excludes exponent (distinct from CSS `<number>`) |
| 05.D7 | `pathLength` listed as geometry-prop AND PA | Plain attribute, `<number>` (full, with exponent) | `pathLength` value = `NumberType`; non-negative is overlay |
| — | Flags must be separate tokens from `number` (maximal-munch) | Browsers parse `A...0150...` correctly | `flag = "0"\|"1"`, kept atomic, never merged with `number` |

---

## Module 06 — Shapes (D1–D9)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 06.D1 | `line` x1/y1/x2/y2 = `<length-percentage> \| <number>` (explicit, not §4.2) | Confirmed by MDN | Value = `LengthPercentageType \| NumberType` directly (no double-add) |
| 06.D2 | rect width/height no explicit initial; CSS `auto` → 0 | Renders 0×0 | rect width/height value = `LengthPercentageType` only; **no `auto`** terminal |
| 06.D3 | circle `r` % basis ambiguous ("size of viewport") | Normalized diagonal `sqrt((w²+h²)/2)` | No grammar change; resolution is overlay |
| 06.D4 | rect rx/ry % vs ellipse rx/ry % different bases | Genuine difference | No grammar change; resolution overlay (distinct positional productions already) |
| 06.D5 | §10.8.2 "cx, cy and y" prose typo (should be "and r") | Typo | No grammar impact |
| 06.D6 | `points` = `[ <number>+ ]#` (ambiguous VDS) | comma/whitespace interchangeable | Author explicit `Points = coordinate-pair { comma-wsp coordinate-pair }`; even-count is overlay |
| 06.D7 | cx/cy/r/rx/ry/x/y are SVG2 CSS geometry properties | Settable via CSS (units required) vs attr (unitless ok) | Attribute productions accept number expansion; CSS context is overlay |
| 06.D8 | `pathLength` is a PA bucket entry, not a CSS property | `<number>`, non-negative | `pathLength` = `NumberType`; on all 6 shapes + path |
| 06.D9 | one-point polyline undefined | Implementation-defined | Accept; overlay note |

---

## Module 07 — Text (D1–D14)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 07.D1 | element tables list text/tspan x/y as PAs, but §11.2.1 says NOT settable via CSS | They are content attributes (SVGAnimatedLengthList) | Model x/y/dx/dy/rotate as **content attributes** (not CSS PAs); list syntax (see §AnimationValue note) |
| 07.D2 | `textLength` initial = UA-computed length | Undefined until layout | Value `LengthPercentageType \| NumberType`; absent default; ≥0 overlay |
| 07.D3 | `dominant-baseline` value-set drift (SVG1.1 vs CSS Inline L3) | CSS Inline L3 set | INCLUDE CSS Inline L3 set (9 kws); `text-before-edge`/`text-after-edge` deprecated aliases (overlay) |
| 07.D4 | `alignment-baseline` initial `baseline` (SVG) vs `auto` (CSS) | SVG override | Keep initial `baseline`; `auto` invalid for SVG (overlay) |
| 07.D5 | `textPath method="stretch"` complex | No browser correctly implements; falls back to align | INCLUDE `stretch` terminal; impl note |
| 07.D6 | `textPath side="right"` MDN 🧪 experimental | Limited support | INCLUDE `left\|right`; experimental flag |
| 07.D7 | SVG2 wrapping props (`shape-inside` etc.) defined | Not implemented (inline-size partial) | INCLUDE in grammar; overlay-flag unimplemented |
| 07.D8 | `textPath spacing` initial `exact` | Browsers don't distinguish | Follow spec (`auto\|exact`, initial exact) |
| 07.D9 | `glyph-orientation-horizontal` removed in SVG2 | MDN ⚠️ | EXCLUDE entirely |
| 07.D10 | `font` shorthand PA, complex grammar | Implemented | Treat as open leaf reference to CSS Fonts grammar |
| 07.D11 | `rotate` `[ <number>+ ]#` comma/space | Both separators accepted | `rotate` list accepts comma + space separators |
| 07.D12 | CSS `text-overflow = clip\|ellipsis\|<string>`; SVG restricts | SVG = closed 2-value | `text-overflow` = `clip\|ellipsis` only |
| 07.D13 | `tref` removed | — | EXCLUDE |
| 07.D14 | `textLength` on tspan: Chrome yes, Firefox no | Spec intends it on all 3 | INCLUDE on text/tspan/textPath; impl note |

---

## Module 08 — Embedded / Linking / Interactivity (D1–D10)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 08.D1 | `zoomAndPan` attr-table initial `disable`; prose "magnify" | MDN/browsers: default magnify | Default = `magnify` (table value is editorial error) |
| 08.D2 | Only xlink:href/xlink:title listed; role/arcrole/show/actuate/type absent | Never implemented for SVG | EXCLUDE the 5 obsolete xlink attrs; KEEP href/title aliases |
| 08.D3 | `async`/`defer` on script not in attr tables | Browsers recognize them (HTML) | INCLUDE `async`/`defer` as boolean-presence attrs on `script`; note "from HTML" |
| 08.D4 | `referrerPolicy` (camel IDL) vs `referrerpolicy` (markup) | markup = lowercase | Grammar attr name = `referrerpolicy` |
| 08.D5 | `cursor` property absent from interact chapter | CSS UI property | Out of scope here; treat as external CSS leaf if needed |
| 08.D6 | `crossorigin` `[anonymous\|use-credentials]?` | empty string = anonymous (HTML) | Value = `"" \| anonymous \| use-credentials` |
| 08.D7 | `foreignObject` has no URL attrs | Confirmed | No href/crossorigin/preserveAspectRatio on foreignObject |
| 08.D8 | `_replace` target removed | Invalid in SVG 2 | EXCLUDE `_replace`; target = `_self\|_parent\|_top\|_blank\|<name>` |
| 08.D9 | SVG `pointer-events` set has no `auto` | Browsers accept `auto` (→ visiblePainted) | INCLUDE `auto` as a CSS-context alias arm (annotated) |
| 08.D10 | `image` href animatable; xlink:href "iff href animatable" | Consistent | No issue |

`view` has NO conditional-processing attrs (per spec); `viewTarget` EXCLUDED (removed).

---

## Module 09 — Painting & Markers (D1–D11)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 09.D1 | `<paint>` arms incl. context-fill/context-stroke | Good support in markers; variable in use-shadow | INCLUDE both arms; overlay flag |
| 09.D2 | `stroke-linejoin: arcs` | Firefox-only | INCLUDE `arcs` terminal; flag |
| 09.D3 | `stroke-linejoin: miter-clip` | Firefox; Chrome partial; not Safari | INCLUDE `miter-clip`; flag |
| 09.D4 | `color-interpolation` mixed-case (`sRGB`,`linearRGB`) | XML case-sensitive | Use canonical casing as terminals; case-lenience overlay |
| 09.D5 | refX/refY left/center/right keywords | Patchy (Chrome88+/FF89+/Safari15+) | INCLUDE all kws + `<length-percentage>\|<number>` |
| 09.D6 | `<dasharray>` `#*` allows empty | Empty ≡ none | `DasharrayType` requires ≥1 token; `none` is the keyword arm |
| 09.D7 | `orient` initial `0` | `0` ≡ `0deg` | `orient` = `auto\|auto-start-reverse\|AngleType\|NumberType` |
| 09.D8 | closed-subpath marker-start/end rendering | Rendering rule | Overlay note |
| 09.D9 | `<paint>` fallback changed SVG1.1→2 | SVG2 silent degradation | Optional fallback arm `[none\|currentColor\|<color>]` after url |
| 09.D10 | `paint-order` not in SVG1.1 | SVG2 | INCLUDE; old validators reject — overlay flag |
| 09.D11 | `stroke-miterlimit` ≥1 (1.1) vs ≥0 (2) | (0,1) valid but always falls back | `<number>`, non-negative overlay |

`<paint>` canonical structure (09 is the canonical home): `none \| currentColor \|
<color> \| url(#id) [none\|currentColor\|<color>]? \| context-fill \| context-stroke`.
**Leaf only the `<color>` and `<url>` parts** — the keyword arms are grammar terminals.

---

## Module 10 — Paint Servers (D1–D10)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 10.D1 | gradient x1/y1/x2/y2 typed `<length>` (table) but prose/IDL allow number+% | MDN: `<length-percentage> \| <number>` | Use `LengthPercentageType \| NumberType` for gradient coords |
| 10.D2 | stop-color incl. `<icccolor>` (SVG1.1) | Removed; not implemented | OMIT icccolor; `stop-color = currentColor \| <color>` |
| 10.D3 | stop-color: is `none` valid? | MDN: no | EXCLUDE `none` from stop-color |
| 10.D4 | `fr` new in SVG2 | Full support since 2020 | INCLUDE `fr` (no caveat) |
| 10.D5 | href + xlink:href both listed | href wins; xlink:href deprecated | INCLUDE both; href canonical, xlink:href deprecated alias |
| 10.D6 | fx/fy defaults computed (= cx/cy) | Run-time default | Mark optional; default is overlay |
| 10.D7 | pattern x/y/width/height typed `<length>` only (no %) | MDN confirms `<length>` | pattern geom = `LengthType` only (differs from gradient coords) |
| 10.D8 | patternContentUnits suppressed by viewBox | Context constraint | Overlay |
| 10.D9 | `stop` content excludes animateTransform + descriptive | Narrower than parents | `stop` content = `animate\|set\|script\|style` only |
| 10.D10 | gradientTransform/patternTransform are CSS-transform PAs | `<transform-list>` | Reuse `TransformList`; positional productions |

---

## Module 11 — Filters (D1–D16)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 11.D1 | feComposite operator: prose 7 (incl. `lighter`), IDL 6 | `lighter` supported (Porter-Duff plus) | INCLUDE `lighter`; note IDL gap |
| 11.D2 | feBlend mode: prose defers, IDL enumerates 16 | All 16 supported | Enumerate all 16 `<blend-mode>` terminals |
| 11.D3 | feGaussianBlur.edgeMode default `none` vs feConvolveMatrix `duplicate` | Distinct | **Separate positional productions** per element (defaults differ) |
| 11.D4 | feDisplacementMap impl ≠ spec (#113) | Known mismatch | Grammar follows spec; overlay note |
| 11.D5 | `kernelUnitLength` deprecated (3 elements) | May be ignored | INCLUDE for parse completeness; deprecated flag |
| 11.D6 | `seed` typed `<number>`, truncated to int | — | `NumberType`; int-trunc overlay |
| 11.D7 | `numOctaves` `<integer>` non-negative | — | `IntegerType` |
| 11.D8 | feImage/feTurbulence/feFlood have no `in` | Confirmed | Omit `in` from those 3 primitives' productions |
| 11.D9 | drop-shadow() param order ≠ feDropShadow | CSS fn vs element attrs | feDropShadow uses dx/dy/stdDeviation attrs; CSS fn separate |
| 11.D10 | grayscale/sepia/invert default(1) vs interp-initial(0) | Intentional | Syntactic default only; animation initial is overlay |
| 11.D11 | feBlend `no-composite` boolean-presence attr | Minimal impl | INCLUDE as optional boolean attr |
| 11.D12 | `filterRes` removed | — | EXCLUDE |
| 11.D13 | feImage `crossorigin` not animatable; IDL `crossOrigin` | — | INCLUDE `anonymous\|use-credentials` |
| 11.D14 | feMergeNode is NOT a filter primitive (category None) | No x/y/w/h/result | feMergeNode production has only `in` (+core) |
| 11.D15 | feSpecularLighting.kernelUnitLength → 2 IDL props | `<number-optional-number>` | One attr, two DOM props; `NumberOptionalNumberType` |
| 11.D16 | feDropShadow flood-color/flood-opacity are CSS props | PAs settable as attrs | Model as presentation-attr values |

`color-interpolation-filters` initial `linearRGB` (canonical home = filters 11).
`flood-color`/`flood-opacity` canonical home = filters; `lighting-color` canonical home
= filters.

---

## Module 12 — Masking & Clipping (D1–D10)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 12.1 | clipPath `use` must directly ref path/text/shape | Browsers vary (open issue #17) | Grammar allows `use` child; direct-ref rule is overlay |
| 12.2 | `<basic-shape>` lacks `shape()`/`polygon() round` (CSS Shapes L1) | Both supported (CSS Shapes L2) | INCLUDE all 6 basic-shape fns; mark L2 extensions |
| 12.3 | `margin-box` removed from mask-clip/mask-origin (kept in clip-path) | — | **Separate `<geometry-box>` productions**: clip-path keeps margin-box; mask-clip/origin omit it |
| 12.4 | `no-clip` valid for mask-clip only (not mask-origin) | — | Distinct productions |
| 12.5 | `clip` deprecated; comma vs space separation | Comma normative | `clip = rect(edge "," edge "," edge "," edge) \| auto`; comma-only |
| 12.6 | spec example `mask-mode: auto` | `auto` not valid | EXCLUDE `auto`; `mask-mode = alpha\|luminance\|match-source` |
| 12.7 | maskUnits default objectBoundingBox vs clipPathUnits userSpaceOnUse | Both correct | Encode distinct defaults (overlay annotation) |
| 12.8 | `-webkit-mask-composite` keyword set differs | Not interchangeable | Standard keywords only (`add\|subtract\|intersect\|exclude`) |
| 12.9 | `clip-rule` SVG-vs-CSS context | applies inside clipPath | Value `nonzero\|evenodd`; applies-to is overlay |
| 12.10 | `shape()` `<shape-command>` not defined here | CSS Shapes L2 | Model `<shape-command>` as a sub-grammar (or open leaf); flag |
| 12.note | mask content model lists removed SVG1.1 elements (`color-profile`,`cursor`,`font`,`font-face`,`altGlyphDef`) | Removed/unsupported | EXCLUDE those from mask content model (use reconciled categories) |

---

## Module 13 — Animation / SMIL (D1–D11)

| ID | SPEC | REAL | DECISION |
|---|---|---|---|
| 13.D1 | SVGwg-L2 removed `attributeType`; SVG1.1 had `CSS\|XML\|auto` | Browsers ignore it | EXCLUDE `attributeType` (auto-resolution always) |
| 13.D2 | `animateColor` deprecated/removed | MDN deprecated | EXCLUDE `animateColor` |
| 13.D3 | `discard` not in SVGwg-L2/SVG1.1 (SVG2 CR draft) | Chrome partial | INCLUDE stub (`href`, `begin`, core); low-confidence flag |
| 13.D4 | mpath href: SVG1.1 path-only, SVGwg-L2 path-or-shape | MDN behind spec | `href` typed IRI; target-type check overlay (path + shapes) |
| 13.D5 | `playbackorder`/`timelinebegin` on svg (SVGwg-L2) | Not in MDN, support unknown | EXCLUDE (not implemented) |
| 13.D6 | SVGwg-L2 adds `#xC` to whitespace `S` | — | INCLUDE `#xC` in `S` |
| 13.D7 | offset-value parenthesization irregular | Intent clear | `offset-value = [ S ("+"\|"-") S ] Clock-value` |
| 13.D8 | repeatCount "floating point" | — | `repeatCount = NumberType \| "indefinite"` |
| 13.D9 | SMIL browser support | Supported (Chrome/FF/Safari); discard/accessKey/wallclock weak | All forms in grammar; support gaps are overlay |
| 13.D10 | SMIL prose has `media-marker-value`; SVG omits | Not in SVG | EXCLUDE `media-marker-value` |
| 13.D11 | Event-symbol is open set | host-defined | `EventSymbol` scalarize leaf |

Animation `fill` (`freeze\|remove`) vs paint `fill` resolved **positionally**:
`AnimateFillAttr` vs `FillAttr`. `animateTransform` `type` (closed 5) drives specialized
per-type value productions (context-free).

---

## Cross-cutting roadblocks

1. **Gluon name-collision dedupe.** `norm(s)=lower, underscores stripped`. Any two rules
   whose normalized names collide are silently deduped by genproto. This means
   `LengthType`/`length_type`, `time_type`/`TimeType`, and (critically for SVG) every
   per-element attribute that shares a base name (`x`, `y`, `width`, `fill`, `type`,
   `offset`, …) MUST get a **distinct PascalCase production name** (e.g. `RectXAttr`,
   `FePointLightXAttr`). See GRAMMAR_PLAN §b. This is the #1 authoring hazard.

2. **Repeated-field omission (PIPELINE_PORTING §6.1).** The proto-css renderer drops
   every `{ }` repetition. For SVG that erases all children/attributes. The renderer
   port MUST emit repetitions. Authors should keep list/child productions as honest
   `{ }` repetitions and rely on the fixed renderer, NOT smuggle tuples through reps.

3. **EBNF comment bug.** Gluon empties a rule body if a `(* *)` comment appears inside
   its RHS. genproto strips comments before parsing. Authors may use comments freely
   (they are stripped), but must not rely on them structurally.

4. **Presentation-attribute §4.2 expansion.** Length/length-percentage/angle valued
   presentation attributes accept an extra unitless `<number>` arm. Encode this in the
   attribute productions (e.g. `StrokeWidthAttr = LengthPercentageType | NumberType`),
   NOT by mutating the shared leaf types (the CSS-property context must stay strict).

5. **`xml:` / `xlink:` colon in attribute names.** Attribute names with a `:` are
   terminals containing a colon (`"xlink:href"`, `"xml:lang"`, `"xml:space"`). They are
   modeled as quoted terminal literals in the attribute-name position; safe for the lexer.
