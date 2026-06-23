# SVG Grammar: Context-Sensitivity Policy

How proto-svg splits SVG into a context-free core and a context-sensitive overlay.

**Governing principle.** We push all context-free parts of SVG into the self-describing EBNF grammar, sometimes deliberately over-approximating, and we enforce the context-sensitive parts as a minimal set of constraints layered over the CFG. The grammar is the source of truth. Gluon turns it into a typed reflection surface. The constraint overlay lives in the parser, renderer, and generator, never in the grammar.

The spec link sets that back this grammar live in [`specs/`](specs/): the MDN index, the W3C SVG 2 sections, and companion-spec pointers.

---

## The two layers

| Layer | What it is | Where it lives |
| --- | --- | --- |
| CFG skeleton | The self-describing EBNF grammar. Elements, attributes, content models, value syntaxes, sub-grammars (path data, transforms, timing), and enumerations. All of it, exhaustively. | `lang/*.ebnf`, compiled by gluon to a proto `FileDescriptorProto` |
| Constraint overlay | The minimal, necessary-only set of checks and choices that a CFG provably cannot express. It refines the CFG when validating and steers it when generating. | generator and renderer (`reps.go`-style), plus a parser validation pass |

This mirrors proto-css. The CSS grammar over-approximates (for example, it does not encode numeric upper bounds), and `reps.go` together with the renderer apply what the CFG cannot (range-bounded samples, leaf sampling, dedup). SVG follows the same pattern, with more relational constraints.

---

## Rule 1: If it can be context-free, it must be context-free

No shortcuts, no "any ident", no free-text escape hatches where a closed form exists. This is mandatory even when it is verbose:

- **Enumerate every enumeration in full.** Write `stroke-linecap = "butt" | "round" | "square"`, not `stroke-linecap = ident`. Every keyword set is listed.
- **List every element's full attribute set explicitly**, assembled from shared non-terminal groups (`CoreAttrs`, `PresentationAttrs`, `ConditionalProcessingAttrs`, `AriaAttrs`, `EventAttrs`, `XlinkAttrs`). No element accepts "any attribute".
- **List every element's content model explicitly:** which children are allowed, drawn from content categories such as shape elements, container elements, descriptive elements, gradient elements, filter primitives, and animation elements. For instance, `feMerge` admits only `feMergeNode`, so spell that out.
- **Write out every sub-grammar in full:** path data BNF, transform-list, `preserveAspectRatio` (align and meetOrSlice), `viewBox` (four numbers), `points` (a list of coordinate pairs), and SMIL timing values (`begin`, `end`). These are genuine context-free languages, so they belong in the grammar.
- **Specialize finite dependent types structurally.** When a "depends-on" relationship ranges over a small closed set, encode the cross-product as productions instead of pushing it to the overlay. For example, the value grammar of `animateTransform` depends on `type`, but `type` is one of exactly five, so write `AnimateTransformRotateValues`, `AnimateTransformTranslateValues`, and the rest directly.
- **Model references as closed oneofs where the target set is closed.** `attributeName` is a closed oneof over all animatable attribute names (generated from the attribute grammar), not a free string.

The test is this: can a single production (with reuse) describe the legal strings without consulting another node's runtime value or an external symbol table? If it can, it goes in the grammar, however long the listing.

### Rule 1a: The scalarize placeholder is for open domains only

The freeform leaf placeholder borrowed from proto-css is a scalarize mechanism: the grammar names an open leaf type and the renderer substitutes sample strings. The proto-css examples are the representative samples in `reps.go`, such as `<length>` becoming `{24px, 48px, ...}`, plus `<number>`, `<color>`, `<string>`, and `<ident>`. It exists for one reason only: the value domain is genuinely open-ended or impractically large to enumerate (real numbers, lengths, arbitrary strings, URLs, the full color space).

It is not a shortcut for closed value sets. If an attribute's legal values are an enumerable set of terminal values, meaning the special keywords that attribute accepts, those terminals must appear in the grammar as an explicit alternation. Never write them as a freeform `x` or leaf and then enforce the allowed set at runtime.

Correct:

- `spreadMethod = "pad" | "reflect" | "repeat" ;`
- `stroke-linecap = "butt" | "round" | "square" ;`
- `gradientUnits = "userSpaceOnUse" | "objectBoundingBox" ;`

Wrong:

- `spreadMethod = ident ;` and then checking it is one of pad, reflect, or repeat in the renderer.

**Split the mixed domains.** When a value is some keywords or an open type, enumerate the keywords as terminals and use a scalarize leaf only for the open part:

- `Width = "auto" | LengthType | PercentageType ;` where `auto` is a grammar terminal and only `LengthType` and `PercentageType` scalarize.
- `MotionRotate = "auto" | "auto-reverse" | NumberType ;`

**Boundary.** A value set goes in the grammar as terminals when it is a closed enumeration defined as part of that attribute's value syntax. The scalarize leaf is reserved for generic, reusable, open data types such as `<number>`, `<length>`, `<percentage>`, `<angle>`, `<color>`, `<string>`, and `<IRI>`. The discriminator: a closed list specific to one attribute becomes terminals, and an open generic datatype becomes a leaf.

This bounds Rule 2. Over-approximation may widen ranges within an open type (accept any number, clamp later), but it may never replace a closed keyword enumeration with a freeform leaf. Enumerations are always exact in the grammar.

---

## Rule 2: Over-approximate in the CFG, tighten in the overlay

When a constraint is not context-free, the grammar still carries a superset so the structure is present and walkable, and the overlay narrows it. The deliberate over-approximations are:

| CFG accepts (the superset) | Overlay tightens it to |
| --- | --- |
| `AnimationValue = LengthType \| NumberType \| ColorType \| PaintType \| ...`, the union of every animatable value type | the single arm matching the resolved `attributeName` type |
| `{ Attr }`, meaning any attribute in any order at any multiplicity | each attribute at most once, with required attributes present |
| `href = IRI`, any reference token | a referent that exists and has the required type |

### What MOVED into the CFG (grammar refactor)

Several bounds that were previously over-approximated and clamped in the overlay
are now **structural** — a bounded-decimal interval and a sign restriction are
regular (hence context-free) languages, so they belong in the grammar (Rule 1):

| Was overlay (clamp) | Now structural (grammar) |
| --- | --- |
| `opacity`/`*-opacity` clamped to [0,1] | `AlphaValue` — a bounded-decimal leaf (integer part exactly `0` or `1`) |
| non-negative magnitudes (`width`, `height`, `r`, `rx`, `ry`, `stroke-width`, `markerWidth/Height`, radial `r`/`fr`, `pathLength`, …) | `NonNegativeNumberType` / `NonNegativeLengthType` / `NonNegativeLengthPercentageType` (no leading minus) |
| `stroke-miterlimit ≥ 1`; `numOctaves ≥ 0` | `MiterLimitType`; `NonNegativeIntegerType` |
| malformed SMIL clocks with negative components (`-1:100.-1`) | clock components are `NonNegativeIntegerType`; `repeatCount` is non-negative — `dur`/`min`/`max`/`repeatDur` are valid by construction (no pin) |
| `named_color = letter+` (accepted any word) | `NamedColor` — the 148 CSS Color L4 names enumerated in full; `ColorType` is a structured `HexColor \| FunctionalColor \| NamedColor` oneof |
| `feColorMatrix values` as a generic number list | a per-`type` union (`FeColorMatrixMatrixValues` = 20 numbers \| scalar), like animateTransform |
| list values (`stroke-dasharray`, `stdDeviation`, `tableValues`, `kernelMatrix`) as opaque string leaves | structured `repeated` fields (the list structure is context-free) |
| `<position>` (`CssPosition`) as `StringType` | a structured edge-keyword + length-percentage grammar |

What STAYS overlay (genuinely non-context-free): list **monotonicity**
(`keyTimes` non-decreasing — a relation between unbounded list elements);
**cross-list cardinality** (`values`↔`keyTimes`↔`keySplines` counts; `kernelMatrix`
length `== order²`); the type→arm **selection** for dependent value unions
(`animateTransform`, `feColorMatrix`); attribute-**presence** dependencies
(`feComposite` `k1-4` require `operator="arithmetic"`; `feFunc` params keyed by
`type`) — these gate *separate sibling attributes* on an enum, which the per-
attribute generator resolves by companion injection rather than by binding them
into one ordered production (that would forfeit the free `{ Attr }` model);
syncbase/event timing reference resolution (`begin`/`end`); and IDREF resolution.

Over-approximation is a feature. It keeps the proto reflection surface complete, so the renderer can always produce a value, and it confines the non-CFG logic to a thin predicate.

---

## Rule 3: The overlay is minimal and necessary-only

A constraint earns a place in the overlay only if it is provably not context-free. The complete, intended overlay for SVG is just these classes.

### General SVG (applies across modules)

1. **Reference and IDREF resolution.** This covers `url(#id)` paint references (`fill` and `stroke` pointing at a gradient or pattern), `clip-path`, `mask`, `marker-*`, and `href` or `xlink:href` (on `use`, `image`, `textPath`, `mpath`, `feImage`, and gradient or pattern templates). The CFG parses the `#id` or IRI token, and the overlay checks that the referent exists and has the required type.
2. **`id` uniqueness.** Each `id` must be unique within the document.
3. **Numeric ranges and monotonicity.** Some values are restricted to a range (`opacity` from 0 to 1, non-negative lengths, `keyTimes` and `keyPoints` from 0 to 1) or to an ordering (`keyTimes` non-decreasing, first entry 0). Lower-bound sign constraints (no leading minus) are encoded structurally in the CFG, as in proto-css; arbitrary ranges and monotonicity are overlay.
4. **Required attributes and "exactly one of".** A `{ }` repetition cannot make an XML attribute mandatory, so "this element requires X" and "exactly one of A or B" are overlay.
5. **Mutual exclusion.** For example, an animation may carry `values`, or it may carry `from`, `to`, and `by`, but not both.

### Animation module (SMIL)

6. **Animation value typing (the dependent type).** The values of `from`, `to`, `by`, and `values` must have the type named by `attributeName` together with the host element. The CFG carries the `AnimationValue` union, and the overlay resolves `attributeName` to a value type and selects the matching arm. The type oracle is derived from the attribute grammar itself: `attributeName="stroke-width"` resolves to the existing `StrokeWidthAttr` value production, so the table is generated, not hand-maintained.
7. **Animatable-ness.** `attributeName` must name an attribute that is actually animatable on the resolved target element. This is subsumed by the resolution step in item 6.
8. **Cross-attribute cardinality.** The length of `keyTimes` relates to the length of `values` (they can differ by one depending on `calcMode`); the length of `keySplines` equals the length of `values` minus one; and `keyTimes` is required when `calcMode="spline"`.
9. **Syncbase and event-base timing resolution.** The timing syntax in `begin="other.end-1s"` is context-free (see Rule 1), but resolving `other` and the sync semantics is overlay, a special case of item 1.

Anything not on this list must be expressible context-free, and therefore belongs in the grammar.

---

## What is not context-sensitive (resolved structurally in the CFG)

Some context that looks like it needs an overlay is actually positional, so it stays in the grammar.

- **Attribute-name collisions across elements.** Consider `fill="freeze|remove"` on animation elements versus `fill=<paint>` on shapes, or `type` on `animateTransform` (the transform kind) versus `type` on `<script>` or `<style>` (a media type) versus `type` on filter components. Each is a distinct field on a distinct element's production (`AnimationFillAttr` versus `FillAttr`), so the correct value set is selected by position in the tree, not by inspecting another value. No overlay is needed.
- **The boundary rule.** Context that is positional or structural goes in the CFG. Context that is value-dependent, meaning it needs another node's runtime value or an external symbol table, goes in the overlay.

---

## Generation mode: the overlay is a guided sampler, not a rejecter

proto-svg generates example SVG the way proto-css generates CSS, so the overlay runs in the forward direction. Instead of accepting only when constraints hold, it chooses values so that the constraints hold by construction:

- Pick `attributeName` first, look up its type, then generate matching `from`, `to`, and `values`.
- Generate `keyTimes` with the count that `calcMode` and `values` require.
- Reference only the `id`s the generator has already emitted, so IDREFs resolve by construction.
- Sample numbers from inside the legal range, so the range is satisfied by construction.

This is the same role `reps.go` plays in proto-css. The CFG enumerates structure, and the overlay supplies the non-CFG specifics. The result is marked in provenance: `pure` when no overlay sampling was needed, and `assisted` when it was.

---

## Animation module: recorded decisions

The factoring agreed for the animation module, applying the rules above:

| Concern | Decision | Layer |
| --- | --- | --- |
| Element and attribute structure (`animate`, `set`, `animateMotion`, `animateTransform`, `mpath`, `discard`) | Full per-element productions with explicit attribute sets and content models | CFG |
| SMIL timing syntax (`begin`, `end`, `dur`, `min`, `max`, `repeatCount`, `repeatDur`, `restart`) | Full context-free grammar of offset, syncbase, event, repeat, accessKey, wallclock, and `indefinite` | CFG |
| Enums (`calcMode`, `additive`, `accumulate`, `fill`, `attributeType`, `rotate`, `type`) | Enumerated in full | CFG |
| List values (`keyTimes`, `keyPoints`, `keySplines`) | List and list-of-tuple productions | CFG |
| `path` on `animateMotion` | Reuse `path.ebnf` (the path-data grammar) | CFG |
| `animateTransform` values keyed by `type` (a closed set of five) | Specialized productions per type (`...RotateValues`, `...TranslateValues`, and so on) | CFG |
| `attributeName` | Closed oneof over all animatable attribute names, generated from the attribute grammar | CFG |
| General typing of `from`, `to`, `by`, and `values` | `AnimationValue` union (an over-approximation) | CFG, narrowed by overlay |
| Resolve `attributeName` to a value type and narrow the union | Type oracle derived from the attribute grammar | Overlay |
| Count matching across `keyTimes`, `values`, and `keySplines`, and `calcMode="spline"` requiring `keyTimes` | Cardinality checks | Overlay |
| `values` versus `from`, `to`, and `by` | Mutual exclusion | Overlay |
| Range and monotonicity of `keyTimes` and `keyPoints` | Range and order checks | Overlay |
| Targets of `href` and `mpath`, and `begin` or `end` syncbase ids | Reference resolution | Overlay |

In short, about 90 percent of the animation module is context-free (structure, timing syntax, enums, lists, transform-value specialization, and the closed-set `attributeName`), and the overlay is a compact set of checks that is mostly derived from the rest of the grammar rather than written by hand.
