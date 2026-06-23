# proto-svg Grammar Coverage Audit

**Question:** Does the proto-svg EBNF grammar (`lang/*.ebnf`) express every *in-scope*
SVG element, attribute, and property?

**Method:** Mechanically extracted the grammar's element and attribute inventory,
built authoritative in-scope name lists from the SVG 2 indices
(`svg2-eltindex.txt`, `svg2-attindex.txt` G.1+G.2, `svg2-propidx.txt`,
`svg2-geometry.txt`) plus companion-spec attrs, cross-checked MDN
(`mdn_docs_elements.md`, `mdn_docs_attributes.md`), applied the scope decision in
`docs/DOC_GAPS.md`, and reconciled every authoritative name into
COVERED / EXCLUDED_BY_SCOPE / MISSING.

**Date:** 2026-06-23

---

## VERDICT (summary)

| Dimension | Authoritative in-scope | Covered | Excluded-by-scope | **MISSING (real gaps)** |
|---|---|---|---|---|
| **Elements** | 64 (SVG 2 index, minus 4 unimplemented) | **64 + `unknown`** | 4 (`audio`/`canvas`/`iframe`/`video`) | **0** |
| **Attributes (G.1 regular)** | 261 names | all in-scope covered | 15 | **0** |
| **Presentation attrs (G.2)** | 59 names | 58 | 1 (`glyph-orientation-horizontal`) | **0** |
| **Properties (propidx)** | 41 names | 40 | 0 | **1 (`line-height`) — ambiguous; see below** |

**One-line answer:** Effectively **YES** for elements and attributes — every in-scope
SVG element and every in-scope attribute/presentation-attribute is expressed in the
grammar. The **single** unresolved item is the property **`line-height`**, which is in
the SVG 2 property index and flagged "Include / presentation attr: yes" in the
grammar-notes, but has **no production in the grammar** and is **not** in the SVG 2 G.2
presentation-attribute list (so it has no authored XML-attribute form). It is either a
genuine gap or should be documented as CSS-only (like `isolation`/`mix-blend-mode`).

---

## STEP 1 — Grammar inventory (extracted)

- **Elements:** 64 real SVG element tags + `unknown` (the `SVGUnknownElement` container).
  (`tag` from `grep` is a false positive — it is the `"<tag"` placeholder in the EBNF
  header-comment template, not a real production.)
- **Attribute-name-bearing productions:** 326 distinct attribute names
  (the `' name="'` leading-terminal across all modules; `attr` excluded as the
  header-comment template token). Presentation attributes are shared via the
  `PresentationAttribute` union in `lang/styling.ebnf`.

---

## ELEMENTS — full reconciliation

### COVERED (65 productions = 64 SVG elements + `unknown`)

Structural/container: `svg g defs symbol use switch a unknown` ·
Descriptive: `desc title metadata` ·
Shapes: `rect circle ellipse line polyline polygon` · Path: `path` ·
Text: `text tspan textPath` ·
Embedded/linking/scripting/styling: `image foreignObject view script style` ·
Paint servers: `linearGradient radialGradient stop pattern` · Painting: `marker` ·
Masking: `clipPath mask` · Filter container: `filter` ·
Filter primitives (16): `feBlend feColorMatrix feComponentTransfer feComposite
feConvolveMatrix feDiffuseLighting feDisplacementMap feDropShadow feFlood feGaussianBlur
feImage feMerge feMorphology feOffset feSpecularLighting feTile feTurbulence` (17) ·
Transfer-fn children: `feFuncR feFuncG feFuncB feFuncA` ·
Light sources: `feDistantLight fePointLight feSpotLight` · Merge child: `feMergeNode` ·
Animation: `animate set animateMotion animateTransform mpath discard`.

Element diff (SVG 2 index vs grammar) is **clean both directions**: zero elements in
the grammar that are absent from the spec; the only spec elements absent from the
grammar are the four excluded ones below.

### EXCLUDED_BY_SCOPE (4)

| Element | Reason (DOC_GAPS Scope decision / 00.1) |
|---|---|
| `audio` | SVG 2 proposed HTML-media embedding; no browser shipped it — unimplemented |
| `canvas` | SVG 2 proposed HTML-canvas embedding; no browser shipped it — unimplemented |
| `iframe` | SVG 2 proposed embedding; no browser shipped it — unimplemented |
| `video` | SVG 2 proposed HTML-media embedding; no browser shipped it — unimplemented |

Also formally excluded but already absent from SVG 2 / MDN (SVG 1.1-removed,
not in the authoritative SVG 2 element index): `altGlyph*`, `glyph`, `glyphRef`,
`missing-glyph`, `font`, `font-face*`, `hkern`, `vkern`, `tref`, `cursor` (element),
`color-profile`, `animateColor`.

### MISSING elements

**None.** Element coverage = **64/64 in-scope** (plus `unknown`).

---

## ATTRIBUTES — full reconciliation

Authoritative attribute/property union = SVG 2 G.1 (261) + G.2 (59) + propidx (41) +
geometry properties `d rx ry cx cy r x y width height` (svg2-geometry §7) =
**324 distinct names**. Grammar covers **307** of them; the 17 not present are all
classified below.

### COVERED (307 / 324) — highlights confirming no oversight

Every attribute the task flagged to "watch for" is present:

| Attribute | Status | Attribute | Status |
|---|---|---|---|
| `systemLanguage` | COVERED | `transform-origin` | COVERED |
| `requiredExtensions` | COVERED | `paint-order` | COVERED |
| `crossorigin` (image/script/feImage) | COVERED | `clip` | COVERED |
| `decoding` (image) | COVERED | `cursor` | COVERED |
| `autofocus` | COVERED | `vector-effect` | COVERED |
| `xlink:href` / `xlink:title` | COVERED | `pointer-events` | COVERED |
| `xml:lang` / `xml:space` (deprecated aliases, kept) | COVERED | `white-space` | COVERED |
| full mask-* longhands: `mask-clip mask-composite mask-image mask-mode mask-origin mask-position mask-repeat mask-size mask-type` | COVERED | `mask` (shorthand) | COVERED |
| full font/text set: `font-family font-size font-size-adjust font-style font-variant font-weight font-stretch text-anchor text-decoration text-rendering letter-spacing word-spacing dominant-baseline alignment-baseline baseline-shift direction unicode-bidi writing-mode text-overflow inline-size glyph-orientation-vertical` | COVERED | geometry `d rx ry cx cy r x y width height x1 x2 y1 y2` | COVERED |
| filter `in in2 k1 k2 k3 k4` + all primitive attrs | COVERED | `data-*` (modeled as `data-name`) | COVERED |

Event handlers: the grammar's `GlobalEventAttribute` (57), `DocumentEventAttribute`
(svg-only `onunload onabort onerror onresize onscroll`), `DocumentElementEventAttribute`
(`oncopy oncut onpaste`), `GraphicalEventAttribute` (`onfocusin onfocusout`),
`AnimationEventAttribute` (`onbegin onend onrepeat`) cover **every** event handler that
the SVG 2 index lists with a non-empty element column, plus `onwheel`.

### EXCLUDED_BY_SCOPE (16 attribute/presentation names)

| Name | Spec source | Documented reason (DOC_GAPS) |
|---|---|---|
| `glyph-orientation-horizontal` | G.2 / propidx (removed) | 07.D9 — removed in SVG 2; EXCLUDE entirely |
| `playbackorder` | G.1 (on `svg`) | Scope table / 00.4 / 13.D5 — SVGwg-L2 addition, not implemented in any browser |
| `timelinebegin` | G.1 (on `svg`) | Scope table / 00.4 / 13.D5 — SVGwg-L2 addition, not implemented in any browser |
| `ondragexit` | G.1 (global event) | 03.D6 — non-standard alias of `ondragleave`, HTML removed it |
| `onmousewheel` | G.1 (global event) | 03.D7 — non-standard; `onwheel` included instead |
| `onshow` | G.1 (global event) | 03.D9 — `show` event removed from HTML; vestigial |
| `onafterprint` | G.1 — **empty element column** | Window event; no SVG element carries it — out of scope by construction |
| `onbeforeprint` | G.1 — empty element column | Window event; no SVG element carries it |
| `onhashchange` | G.1 — empty element column | Window event; no SVG element carries it |
| `onmessage` | G.1 — empty element column | Window event; no SVG element carries it |
| `onoffline` | G.1 — empty element column | Window event; no SVG element carries it |
| `ononline` | G.1 — empty element column | Window event; no SVG element carries it |
| `onpagehide` | G.1 — empty element column | Window event; no SVG element carries it |
| `onpageshow` | G.1 — empty element column | Window event; no SVG element carries it |
| `onpopstate` | G.1 — empty element column | Window event; no SVG element carries it |
| `onstorage` | G.1 — empty element column | Window event; no SVG element carries it |

The 10 window-event handlers above appear in the SVG 2 G.1 table but with **no element**
in their "Elements on which the attribute may be specified" column (verified directly in
`svg2-attindex.txt`): they are `Window`/`Document`-object events, not attributes of any
SVG element, so they are not authorable SVG markup.

Additional EXCLUDED names surfaced only by the **MDN** cross-check (not even in the SVG 2
index, all confirmed by their MDN status markers ⚠️/🧪/🔶):

| Name | MDN status | Reason (DOC_GAPS) |
|---|---|---|
| `attributeType` | ⚠️ | 00.5 / 13.D1 — removed in SVG Animations L2; browsers parse-but-ignore |
| `baseProfile` | ⚠️ | 03.D4 / scope — SVG 1.1 conformance signal, ignored in SVG 2 |
| `version` | ⚠️ | 03.D4 / scope — SVG 1.1 conformance signal, ignored in SVG 2 |
| `requiredFeatures` | ⚠️ | 03.D5 / scope — SVG 1.1 conditional processing, removed in SVG 2 |
| `xlink:arcrole` | ⚠️ | 08.D2 — obsolete XLink, never implemented for SVG |
| `xlink:show` | ⚠️ | 08.D2 — obsolete XLink, never implemented for SVG |
| `xlink:type` | ⚠️ | 08.D2 — obsolete XLink, never implemented for SVG |
| `fetchpriority` | 🧪🔶 | scope — non-standard, not cross-browser |
| `font-width` | 🧪🔶 | scope — non-standard, not cross-browser |

(`xlink:role`, `xlink:actuate` are in the same 08.D2 exclusion; they do not even appear
in the MDN SVG attribute list.) Also `defer`/`viewTarget`/`filterRes`/`kerning` are
excluded per the scope table — none appear in the authoritative attribute lists for any
in-scope element. (`defer` as a **`<script>`** attribute is a different attribute and IS
present in the grammar; only `defer` inside `preserveAspectRatio` is excluded.)

CSS-only properties required by SVG 2 but with **no presentation-attribute (XML) form**,
therefore not authorable markup and out of the attribute grammar by construction
(DOC_GAPS 02.D4): **`isolation`**, **`mix-blend-mode`**. Not in G.2, not in the MDN SVG
attribute reference. Correctly absent.

### MISSING attributes (real gaps)

**None.** Every in-scope attribute and presentation attribute is expressed.

---

## PROPERTIES — full reconciliation

The SVG 2 property index (`svg2-propidx.txt`, 41 property rows) reconciles as:

- **40 COVERED** — every property whose name is also a presentation attribute (G.2) is
  present in the grammar (the `PresentationAttribute` union plus the geometry/font/text
  modules). This includes `marker` (shorthand), `marker-start/mid/end`, the full
  fill/stroke families, `opacity`, `paint-order`, `pointer-events`, `display`,
  `visibility`, `overflow`, `white-space`, `writing-mode`, etc.

- **1 NOT EXPRESSED — `line-height`** (see below).

### MISSING / AMBIGUOUS property: `line-height`

| Aspect | Finding |
|---|---|
| In SVG 2 property index? | **Yes** — `svg2-propidx.txt` line 289: `normal \| <number> \| <length-percentage>`, applies to `text`. |
| In SVG 2 G.2 presentation-attribute list? | **No** — not listed; SVG 2 defines no `line-height=""` presentation attribute. |
| In MDN SVG attribute reference? | **No.** |
| Grammar-notes decision? | `00-indices.md` line 671: **"Include"**; `07-text.md` line 847–862 documents it with **"Presentation attr: yes"**. |
| In the grammar? | **No production** — no `LineHeight*` rule, name appears in no `lang/*.ebnf` file. |

This is the one genuine inconsistency. The grammar-notes intended to include
`line-height` and assert it is a presentation attribute, but (a) it has no production,
and (b) per SVG 2 G.2 it is **not** a presentation attribute (no XML-attribute form) — so
on the spec's terms it belongs to the same CSS-only category as `isolation` /
`mix-blend-mode`, which DOC_GAPS 02.D4 deliberately keeps out of the attribute grammar.

**Resolution required (one of):**
1. If treated like `isolation`/`mix-blend-mode` (CSS-only, no presentation attribute) —
   then it is **correctly absent** and the grammar-notes ("presentation attr: yes") and
   DOC_GAPS should be corrected to record the exclusion. This is the spec-accurate call:
   SVG 2 G.2 does not list `line-height`.
2. If the project chooses to model it as a presentation attribute anyway (browsers accept
   `line-height` via `style=""` but generally not as a bare SVG attribute) — then add a
   `LineHeightAttr` production and wire it into the text presentation-attribute union.

Recommended: option 1 (document as CSS-only, mirroring 02.D4), since no spec defines a
`line-height` presentation attribute and MDN does not list it as an SVG attribute.

### Text-wrapping properties named "grammar-present" in DOC_GAPS but absent

DOC_GAPS "Features INCLUDED despite caveats" (lines 100–102) claims `inline-size`,
`shape-inside`, `shape-subtract`, `shape-padding`, `shape-margin` are "grammar-present".
Verified:

| Property | In grammar? | In SVG 2 indices / MDN? | In scope per authoritative lists? |
|---|---|---|---|
| `inline-size` | **Yes** (covered) | SVG 2 text §; present | yes |
| `shape-inside` | **No** | not in SVG 2 propidx/attindex, not MDN SVG attr | **not in authoritative in-scope set** (CSS-spec, "not implemented in browsers" per 07-text.md) |
| `shape-subtract` | **No** | not in SVG 2 propidx/attindex, not MDN SVG attr | not in authoritative in-scope set |
| `shape-padding` | **No** | not in SVG 2 propidx/attindex, not MDN SVG attr | not in authoritative in-scope set |
| `shape-margin` | **No** | not in SVG 2 propidx/attindex, not MDN SVG attr | not in authoritative in-scope set |

These four are **not** counted as MISSING here because they are not in the authoritative
in-scope name lists (SVG 2 indices + MDN); 07-text.md itself marks them as deferred to CSS
and "not implemented in browsers." However, **the DOC_GAPS claim that they are
"grammar-present" is inaccurate** and should be corrected (they are not in the grammar).
Same `text-align`/`text-indent`/`hyphens`/`word-break`/`line-break` family — documented as
"deferred to CSS specifications" (open leaves), correctly not emitted.

---

## Final tallies

- **Elements:** 64/64 in-scope COVERED (+`unknown`); 4 EXCLUDED_BY_SCOPE; **0 MISSING**.
- **Attributes (G.1 + G.2 + geometry):** all in-scope COVERED; 25 EXCLUDED_BY_SCOPE
  (16 attribute-index names incl. 10 window events + 3 non-standard/deprecated MDN-only
  +; plus xlink/version/baseProfile/etc.); **0 MISSING**.
- **Properties (propidx):** 40/41 COVERED; **`line-height` not expressed** — ambiguous
  (spec-accurate to exclude as CSS-only, but DOC_GAPS/notes currently say "include").

**Is every in-scope SVG element / attribute / property expressed?**
**Elements: yes. Attributes/presentation attributes: yes. Properties: all except
`line-height`** — and `line-height` is most likely *correctly* absent (CSS-only, not an
SVG presentation attribute per SVG 2 G.2), but the project's own notes claim otherwise, so
it must be reconciled: either add a `LineHeightAttr` production or amend the notes/DOC_GAPS
to record it as CSS-only (the spec-accurate choice). No other gaps exist.
