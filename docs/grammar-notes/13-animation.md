# Animation (SMIL) grammar notes

## Source

- PRIMARY: `/docs/specs/cache/svgwg-animations.txt` — SVG Animations Level 2, W3C Editor's Draft 14 September 2025  
- SECONDARY: `/docs/specs/cache/smil-animation.txt` — SMIL Animation, W3C Recommendation 04-September-2001  
- CROSS-CHECK: `/docs/specs/cache/svg11-animate.txt` — SVG 1.1 Second Edition, Chapter 19 (16 August 2011)  
- POLICY: `/docs/CONTEXT_SENSITIVITY.md`

---

## Elements

### `animate`

**Category:** Animation element  
**Content model:** Any number of descriptive elements (`desc`, `title`, `metadata`) and `script`, in any order. (SVGwg-L2 §2.13; SVG11 §19.2.12 same minus `script`.)  
**DOM interface:** `SVGAnimateElement : SVGAnimationElement`

**Attribute groups on `animate`:**
- animation target element attributes: `href` (URL, no default — parent element if absent)
- animation attribute target attributes: `attributeName` (Name, required; `attributeType` is NOT supported in SVG Animations Level 2 — removed from SVGwg-L2 §7)
- animation timing attributes: `begin`, `dur`, `end`, `min`, `max`, `restart`, `repeatCount`, `repeatDur`, `fill`
- animation value attributes: `calcMode`, `values`, `keyTimes`, `keySplines`, `from`, `to`, `by`
- animation addition attributes: `additive`, `accumulate`
- animation event attributes: `onbegin`, `onend`, `onrepeat`
- conditional processing: `requiredExtensions`, `systemLanguage`
- core attributes: `id`, `tabindex`, `autofocus`, `lang`, `xml:space`, `class`, `style`
- global event attributes (full set: `oncancel` ... `onwheel` — see spec §2.13 for complete list)
- presentation attributes (inheritable set)

**Note on `attributeType`:** Present in SVG 1.1 with values `"CSS" | "XML" | "auto"` (default `"auto"`). **Removed in SVG Animations Level 2.** SVGwg-L2 always applies auto-resolution (CSS properties checked first, then XML attributes). The SVG 1.1 cross-check source preserves `attributeType`; the SVGwg-L2 primary source explicitly removed it in §7 "Changes since SVG 1.1 Second Edition." Grammar decision: `attributeType` is NOT included for SVGwg-L2 target; note as a removed attribute for SVG 1.1 compatibility layer if needed.

---

### `set`

**Category:** Animation element  
**Content model:** Any number of descriptive elements (`desc`, `title`, `metadata`) and `script`, in any order.  
**DOM interface:** `SVGSetElement : SVGAnimationElement`

**Attribute groups:**
- animation event attributes: `onbegin`, `onend`, `onrepeat`
- animation target element attributes: `href`
- animation attribute target attributes: `attributeName`
- animation timing attributes: `begin`, `dur`, `end`, `min`, `max`, `restart`, `repeatCount`, `repeatDur`, `fill`
- conditional processing, core, global event attributes (as animate)
- `to` (AnimationValue — the single value to set)

**Excluded:** `additive`, `accumulate`, `calcMode`, `values`, `keyTimes`, `keySplines`, `from`, `by` — not present on `set`. (Spec: "The set element is non-additive. The additive and accumulate attributes are not allowed, and will be ignored if specified.")

---

### `animateMotion`

**Category:** Animation element  
**Content model:** Any number of descriptive elements and `script`, and **at most one** `mpath` element, in any order.  
**DOM interface:** `SVGAnimateMotionElement : SVGAnimationElement`

**Attribute groups:**
- animation addition attributes: `additive`, `accumulate`
- animation event attributes: `onbegin`, `onend`, `onrepeat`
- animation target element attributes: `href`
- animation timing attributes: `begin`, `dur`, `end`, `min`, `max`, `restart`, `repeatCount`, `repeatDur`, `fill`
- animation value attributes: `calcMode` (default `"paced"` — differs from `animate`!), `values`, `keyTimes`, `keySplines`, `from`, `to`, `by`
- conditional processing, core, global event attributes
- `path` (svg-path data — full SVG path syntax)
- `keyPoints` (`<number> [; <number>]* ;?` — semicolon-separated list of [0,1] floats)
- `rotate` (`<number> | "auto" | "auto-reverse"`, default `0`)
- `origin` (`"default"`, only value; no effect in SVG)

**Note:** `attributeName` and `attributeType` are NOT used on `animateMotion`; the position attribute is implicitly targeted as defined by the host language (SVG: CTM supplemental translation).

**animateMotion coordinate pairs:** `from`, `to`, `by`, `values` specify x,y pairs. Format: `<number> (","? | S) <number>` per pair; `values` is semicolon-separated list of coordinate pairs. Each coordinate is a `<length>`.

**Priority override rule (not grammar, informative):** `mpath` overrides `path` overrides `values` overrides `from`/`by`/`to`.

---

### `mpath`

**Category:** (none)  
**Content model:** Any number of descriptive elements (`desc`, `title`, `metadata`) and `script`, in any order.  
**DOM interface:** `SVGMPathElement : SVGElement` (includes `SVGURIReference`)

**Attributes:**
- core attributes: `id`, `tabindex`, `autofocus`, `lang`, `xml:space`, `class`, `style`
- global event attributes
- `href` (URL — required; references a `path` element or shape element; resolution is overlay)

**Note (SVG 1.1 vs SVGwg-L2):** SVG 1.1 restricts `mpath href` to `path` elements only. SVGwg-L2 extends this to any shape element. Grammar: `href` is typed as `IRI`; target-type check is overlay.

---

### `animateTransform`

**Category:** Animation element  
**Content model:** Any number of descriptive elements (`desc`, `title`, `metadata`) and `script`, in any order.  
**DOM interface:** `SVGAnimateTransformElement : SVGAnimationElement`

**Attribute groups:**
- animation addition attributes: `additive`, `accumulate`
- animation event attributes: `onbegin`, `onend`, `onrepeat`
- animation target element attributes: `href`
- animation attribute target attributes: `attributeName`
- animation timing attributes: `begin`, `dur`, `end`, `min`, `max`, `restart`, `repeatCount`, `repeatDur`, `fill`
- animation value attributes: `calcMode`, `values`, `keyTimes`, `keySplines`, `from`, `to`, `by`
- conditional processing, core, global event attributes
- `type` (`"translate" | "scale" | "rotate" | "skewX" | "skewY"`, default `"translate"`)

**Note:** `from`/`to`/`by`/`values` interpretation depends on `type` (see animateTransform per-type value productions below). This dependency is closed-set and resolved structurally in the grammar.

---

### `discard`

**Note:** The `discard` element is referenced in some SVG contexts for progressive loading/long-running documents. It is **not defined in SVG Animations Level 2** (the primary spec). It appears in SVG 2 (draft) as an extension. MDN documents it. Chrome/Firefox support is incomplete (Chrome unflagged SMIL generally; `discard` support is sparse). Grammar decision: include a stub production but flag as low-confidence/not-in-primary-spec; record in discrepancies.

Proposed attributes (from SVG 2 draft / MDN cross-check, not in primary spec text):
- `href` (URL — target element to discard)
- `begin` (begin-value-list — when to discard)
- core attributes

---

## SMIL timing BNF (verbatim)

### From SMIL Animation spec (§3.6.7, normative)

```
S ::= (#x20 | #x9 | #xD | #xA)*
```

**Begin values:**
```
begin-value-list ::= begin-value (S ";" S begin-value-list )?
begin-value      ::= (offset-value | syncbase-value
                      | event-value
                      | repeat-value | accessKey-value
                      | wallclock-sync-value | "indefinite" )
```

**End values:**
```
end-value-list ::= end-value (S ";" S end-value-list )?
end-value      ::= (offset-value | syncbase-value
                      | event-value
                      | repeat-value | accessKey-value
                      | wallclock-sync-value | "indefinite" )
```

**Offset values (from §3.6.7):**
```
offset-value   ::= (( S "+" | "-" S )? ( Clock-value )
```

**ID-Reference values:**
```
Id-value   ::= IDREF
```

**Syncbase values (from §3.6.7):**
```
Syncbase-value   ::= ( Syncbase-element "." Time-symbol )
                     ( S ("+"|"-") S Clock-value )?
Syncbase-element ::= Id-value
Time-symbol      ::= "begin" | "end"
```

**Event values (from §3.6.7):**
```
Event-value        ::= ( Eventbase-element "." )? Event-symbol
                        ( S ("+"|"-") S Clock-value )?
Eventbase-element  ::= ID
```

**Repeat values (from §3.6.7):**
```
Repeat-value       ::= ( Eventbase-element "." )? "repeat(" iteration ")"
                        ( S ("+"|"-") S Clock-value )?
iteration          ::= DIGIT+
```

**AccessKey values (from §3.6.7):**
```
AccessKey-value ::= "accessKey(" character ")"
                    ( S ("+"|"-") S Clock-value )?
```

**Wallclock-sync values (from §3.6.7):**
```
wallclock-val  ::= "wallclock(" S (DateTime | WallTime)  S ")"
DateTime       ::= Date "T" WallTime
Date           ::= Years "-" Months "-" Days
WallTime       ::= (HHMM-Time | HHMMSS-Time)(TZD)?
HHMM-Time      ::= Hours24 ":" Minutes
HHMMSS-Time    ::= Hours24 ":" Minutes ":" Seconds ("." Fraction)?
Years          ::= 4DIGIT;
Months         ::= 2DIGIT; range from 01 to 12
Days           ::= 2DIGIT; range from 01 to 31
Hours24        ::= 2DIGIT; range from 00 to 23
4DIGIT         ::= DIGIT DIGIT DIGIT DIGIT
TZD            ::= "Z" | (("+" | "-") Hours24 ":" Minutes )
```

**Clock values (from SMIL §3.6.7 and SVGwg-L2 §2.9.1 — both agree verbatim):**
```
Clock-value         ::= ( Full-clock-value | Partial-clock-value
                          | Timecount-value )
Full-clock-value    ::= Hours ":" Minutes ":" Seconds ("." Fraction)?
Partial-clock-value ::= Minutes ":" Seconds ("." Fraction)?
Timecount-value     ::= Timecount ("." Fraction)? (Metric)?
Metric              ::= "h" | "min" | "s" | "ms"
Hours               ::= DIGIT+; any positive number
Minutes             ::= 2DIGIT; range from 00 to 59
Seconds             ::= 2DIGIT; range from 00 to 59
Fraction            ::= DIGIT+
Timecount           ::= DIGIT+
2DIGIT              ::= DIGIT DIGIT
DIGIT               ::= [0-9]
```

**Note (SVGwg-L2 §2.9 editorial comment):** "Align with whitespace used in CSS and SVG, adding #xC to S." — SVGwg-L2 suggests extending S to include U+000C (form feed). Grammar decision: include `#xC` in S for SVG.

**keySplines control-point syntax (from SMIL §3.2.3):**
```
control-pt-set ::= ( fpval comma-wsp fpval comma-wsp fpval comma-wsp fpval )
fpval          ::= Floating point number
comma-wsp      ::= S (spacechar|",") S
```

**SVGwg-L2 formulation (§2.10):**
```
<control-point> = <number> ,? <number> ,? <number> ,? <number>
```

---

## SMIL timing EBNF-ready

The following productions are expressed in proto-svg EBNF style, ready for grammar authoring. All are context-free.

```ebnf
(* Whitespace — SVGwg-L2 adds #xC *)
S = { #x20 | #x9 | #xD | #xA | #xC } ;

(* ---------- Clock values ---------- *)
DIGIT              = [0-9] ;
TWO_DIGIT          = DIGIT, DIGIT ;
Clock-value        = Full-clock-value | Partial-clock-value | Timecount-value ;
Full-clock-value   = Hours, ":", Minutes, ":", Seconds, [ ".", Fraction ] ;
Partial-clock-value = Minutes, ":", Seconds, [ ".", Fraction ] ;
Timecount-value    = Timecount, [ ".", Fraction ], [ Metric ] ;
Metric             = "h" | "min" | "s" | "ms" ;
Hours              = DIGIT, { DIGIT } ;       (* any positive integer *)
Minutes            = TWO_DIGIT ;              (* 00-59, range overlay *)
Seconds            = TWO_DIGIT ;              (* 00-59, range overlay *)
Fraction           = DIGIT, { DIGIT } ;
Timecount          = DIGIT, { DIGIT } ;

(* ---------- Offset value ---------- *)
offset-value = [ S, ( "+" | "-" ), S ], Clock-value ;

(* ---------- ID reference ---------- *)
Id-value = IDREF ;   (* XML IDREF — open scalarize leaf *)

(* ---------- Syncbase value ---------- *)
syncbase-value    = Id-value, ".", ( "begin" | "end" ),
                    [ S, ( "+" | "-" ), S, Clock-value ] ;

(* ---------- Event-symbol (open leaf — any DOM event name) ---------- *)
event-symbol = EventSymbol ;   (* scalarize: open set, defined by host language *)

(* ---------- Event value ---------- *)
event-value = [ Id-value, "." ], event-symbol,
              [ S, ( "+" | "-" ), S, Clock-value ] ;

(* ---------- Repeat value ---------- *)
repeat-value = [ Id-value, "." ], "repeat(", iteration, ")",
               [ S, ( "+" | "-" ), S, Clock-value ] ;
iteration    = DIGIT, { DIGIT } ;

(* ---------- AccessKey value ---------- *)
(* character = single Unicode code point, open scalarize leaf *)
accessKey-value = "accessKey(", character, ")",
                  [ S, ( "+" | "-" ), S, Clock-value ] ;

(* ---------- Wallclock value ---------- *)
wallclock-sync-value = "wallclock(", S, ( DateTime | WallTime ), S, ")" ;
DateTime             = Date, "T", WallTime ;
Date                 = Years, "-", Months, "-", Days ;
WallTime             = ( HHMM-Time | HHMMSS-Time ), [ TZD ] ;
HHMM-Time            = Hours24, ":", Minutes ;
HHMMSS-Time          = Hours24, ":", Minutes, ":", Seconds, [ ".", Fraction ] ;
Years                = DIGIT, DIGIT, DIGIT, DIGIT ;
Months               = TWO_DIGIT ;   (* 01-12, range overlay *)
Days                 = TWO_DIGIT ;   (* 01-31, range overlay *)
Hours24              = TWO_DIGIT ;   (* 00-23, range overlay *)
TZD                  = "Z" | ( ( "+" | "-" ), Hours24, ":", Minutes ) ;

(* ---------- Begin/end values ---------- *)
begin-value = offset-value | syncbase-value | event-value
            | repeat-value | accessKey-value
            | wallclock-sync-value | "indefinite" ;

begin-value-list = begin-value, { S, ";", S, begin-value } ;

end-value = offset-value | syncbase-value | event-value
          | repeat-value | accessKey-value
          | wallclock-sync-value | "indefinite" ;

end-value-list = end-value, { S, ";", S, end-value } ;

(* ---------- keySplines control-point set ---------- *)
comma-wsp       = S, ( "," | S ), S ;
control-pt-set  = NumberType, comma-wsp, NumberType, comma-wsp,
                  NumberType, comma-wsp, NumberType ;
keySplines-list = control-pt-set, { S, ";", S, control-pt-set }, [ S, ";" ] ;

(* ---------- keyTimes / keyPoints lists ---------- *)
keyTimes-list  = NumberType, { S, ";", S, NumberType }, [ S, ";" ] ;
keyPoints-list = NumberType, { S, ";", S, NumberType }, [ S, ";" ] ;
```

**Notes on EBNF-ready form:**
- `IDREF`, `NumberType`, `EventSymbol`, `character` are open scalarize leaves (context-free placeholders replaced at generation time).
- Range constraints on `Minutes` (00–59), `Seconds` (00–59), `Hours24` (00–23), `Months` (01–12), `Days` (01–31), `keyTimes` entries ([0,1]), `keyPoints` entries ([0,1]), `keySplines` entries ([0,1]) are overlay checks (Rule 2).
- `keyTimes` monotonicity (non-decreasing), first=0 and last=1 (for linear/spline), first=0 (for discrete) are overlay checks.
- `keyTimes`/`keyPoints`/`keySplines` cardinality matching against `values` count is overlay (Rule 3, item 8).
- Syncbase `Id-value` and event `Id-value` resolution (must reference existing elements) is overlay (Rule 3, item 9).

---

## animateTransform per-type value productions

Per CONTEXT_SENSITIVITY.md Rule 1: `type` is a closed set of 5, so the value shapes are specialized productions in the grammar (context-free).

```ebnf
(* animateTransform: per-type single-value productions *)

(* type="translate": <tx> [,<ty>] *)
AnimateTransformTranslateValue = NumberType, [ S, ","?, S, NumberType ] ;

(* type="scale": <sx> [,<sy>]  (if sy omitted, sy = sx) *)
AnimateTransformScaleValue = NumberType, [ S, ","?, S, NumberType ] ;

(* type="rotate": <rotate-angle> [<cx> <cy>] *)
(* space-or-comma separated; cx,cy are the center of rotation *)
AnimateTransformRotateValue = NumberType,
                              [ S, NumberType, S, NumberType ] ;

(* type="skewX": <skew-angle> *)
AnimateTransformSkewXValue = NumberType ;

(* type="skewY": <skew-angle> *)
AnimateTransformSkewYValue = NumberType ;

(* Union used in CFG over-approximation for from/to/by values on animateTransform *)
(* The correct arm is selected by overlay based on the type attribute *)
AnimateTransformSingleValue =
    AnimateTransformTranslateValue
  | AnimateTransformScaleValue
  | AnimateTransformRotateValue
  | AnimateTransformSkewXValue
  | AnimateTransformSkewYValue ;

(* values attribute = semicolon-separated list of single values *)
AnimateTransformValues = AnimateTransformSingleValue,
                         { S, ";", S, AnimateTransformSingleValue },
                         [ S, ";" ] ;
```

**Overlay note:** The grammar carries the full union `AnimateTransformSingleValue`. The overlay resolves `type` and selects the correct arm when validating or generating. In generation mode, `type` is chosen first, then the matching arm is generated (per CONTEXT_SENSITIVITY.md §Generation mode).

**Disambiguation of scale vs translate:** Both have the same structural shape (`number [, number]`). They are distinguished only by `type` attribute value — this is a pure overlay concern; the CFG arms are identical in structure (hence the union over-approximation is reasonable).

**Paced animation distances** (informative, affects generation sampling):
- translate: Euclidean distance in (tx, ty) space
- scale: Euclidean distance in (sx, sy) space (sy defaults to sx if omitted)
- rotate, skewX, skewY: distance on the angle component only (cx, cy not used for distance)

---

## AnimationValue union (CFG over-approximation) + overlay narrowing note

Per CONTEXT_SENSITIVITY.md Rule 2: `from`, `to`, `by`, `values` on `animate` and `set` carry an over-approximating union of all animatable value types. The overlay resolves `attributeName` to the actual type.

```ebnf
(* CFG over-approximation — the grammar accepts any of these arms.
   Overlay narrows to the single arm matching the resolved attributeName type. *)
AnimationValue =
    LengthType            (* <length>: e.g. stroke-width, cx, cy, r, x, y, width, height *)
  | NumberType            (* <number>: e.g. opacity, stroke-miterlimit *)
  | PercentageType        (* <percentage> *)
  | AngleType             (* <angle> *)
  | ColorType             (* <color> / <paint>: e.g. fill, stroke *)
  | PaintType             (* <paint>: url(#id) | <color> | none | ... *)
  | TransformListType     (* <transform-list>: for animateTransform *)
  | PathDataType          (* <path-data>: for d attribute *)
  | IntegerType           (* <integer> *)
  | StringType            (* arbitrary string: for non-interpolable attrs *)
  | IriType               (* <IRI>: for href-like attrs *)
  | KeywordType           (* enumeration keywords: visibility, display, etc. *)
  ;

(* values attribute = semicolon-separated list of AnimationValue *)
AnimationValueList = AnimationValue, { S, ";", S, AnimationValue }, [ S, ";" ] ;
```

**Overlay narrowing (type oracle):**
- Input: `attributeName` string (resolved against target element in context)
- Output: one arm of `AnimationValue`
- The oracle is derived from the attribute grammar itself: `attributeName="stroke-width"` → `LengthType | PercentageType`; `attributeName="opacity"` → `NumberType`; etc.
- The table is generated (not hand-maintained) from the attribute productions in the grammar.
- `attributeType` (removed in SVGwg-L2) no longer participates in resolution; the auto-resolution rule (CSS properties checked first) applies universally.

**`set` element:** `to` uses the same `AnimationValue` over-approximation. Only one value, no interpolation.

**`animateMotion` special case:** `from`, `to`, `by`, `values` use a coordinate-pair type (not the general `AnimationValue`):
```ebnf
MotionCoordinatePair = NumberType, ( ( S?, "," , S?) | S ), NumberType ;
MotionValues = MotionCoordinatePair, { S, ";", S, MotionCoordinatePair }, [ S, ";" ] ;
```

---

## Animation attributes (each: value, default, constraints)

### Target identification attributes

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `href` | `IRI` | (none) — parent element used | no | Common to all animation elements. Overlay: must resolve to a valid target. xlink:href deprecated, same semantics. |
| `attributeName` | `Name` (XML Name) | (none) — required | no | Closed oneof over all animatable attribute names (overlay-generated from attribute grammar). Resolution: CSS properties first, then XML attributes. |
| `attributeType` | `"CSS" \| "XML" \| "auto"` | `"auto"` | no | **SVG 1.1 only.** Removed in SVGwg-L2. Not in grammar for SVGwg-L2 target. |

### Timing attributes (common to all animation elements)

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `begin` | `begin-value-list` | `"0s"` (implicit) | no | Semicolon-separated list. All value forms context-free. ID resolution overlay. |
| `end` | `end-value-list` | (none) | no | Semicolon-separated list. Same value grammar as `begin`. |
| `dur` | `Clock-value \| "media" \| "indefinite"` | `"indefinite"` (if attr absent) | no | `"media"` is accepted syntax but ignored for SVG animation elements. Value > 0 is overlay. |
| `min` | `Clock-value \| "media"` | `"0"` (no constraint) | no | `"media"` ignored for SVG. Value >= 0 is overlay. |
| `max` | `Clock-value \| "media"` | (none — no constraint) | no | `"media"` ignored for SVG. Value > 0 is overlay. |
| `restart` | `"always" \| "whenNotActive" \| "never"` | `"always"` | no | Closed 3-keyword enum. |
| `repeatCount` | `NumberType \| "indefinite"` | (none) | no | NumberType must be > 0 (overlay). Fractional values allowed. |
| `repeatDur` | `Clock-value \| "indefinite"` | (none) | no | |
| `fill` | `"freeze" \| "remove"` | `"remove"` | no | **Animation fill** (distinct from presentation `fill` attribute on shape elements — see CONTEXT_SENSITIVITY.md §"What is not context-sensitive": positional distinction). |

### Value attributes (animate, animateMotion, animateTransform)

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `calcMode` | `"discrete" \| "linear" \| "paced" \| "spline"` | `"linear"` (animate/animateTransform); `"paced"` (animateMotion) | no | Different defaults per element — modeled as distinct element-level attribute definitions. |
| `values` | `AnimationValueList` (semicolon-sep) | (none) | no | If present, `from`/`to`/`by` are ignored. Overlay: each entry must match resolved `attributeName` type. |
| `keyTimes` | `keyTimes-list` (semicolon-sep `NumberType`s) | (none) | no | Cardinality, range [0,1], monotonicity, first=0 last=1 — all overlay. |
| `keySplines` | `keySplines-list` (semicolon-sep `control-pt-set`s) | (none) | no | Ignored unless `calcMode="spline"`. Count = count(keyTimes) - 1 is overlay. All values in [0,1] — overlay. |
| `from` | `AnimationValue` | (none) | no | Overlay: type must match resolved `attributeName`. |
| `to` | `AnimationValue` | (none) | no | Overlay: type must match resolved `attributeName`. |
| `by` | `AnimationValue` | (none) | no | Overlay: type must match resolved `attributeName`. Addition-compatible types only (overlay). |

### Addition attributes (animate, animateMotion, animateTransform)

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `additive` | `"replace" \| "sum"` | `"replace"` | no | Closed 2-keyword enum. |
| `accumulate` | `"none" \| "sum"` | `"none"` | no | Closed 2-keyword enum. Ignored if attr type doesn't support addition (overlay). |

### animateMotion-specific attributes

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `path` | `svg-path` (full SVG path data) | (none) | no | Reuses path.ebnf. If present, overrides `values`, `from`, `to`, `by`. |
| `keyPoints` | `keyPoints-list` (semicolon-sep `NumberType`s in [0,1]) | (none) | no | Must have same count as `keyTimes` (overlay). Range [0,1] is overlay. |
| `rotate` | `NumberType \| "auto" \| "auto-reverse"` | `"0"` | no | Mixed domain: keywords + open number type. |
| `origin` | `"default"` | `"default"` | no | Only one legal value in SVG. No effect. |

### animateTransform-specific attributes

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `type` | `"translate" \| "scale" \| "rotate" \| "skewX" \| "skewY"` | `"translate"` | no | Closed 5-keyword enum. Controls interpretation of `from`/`to`/`by`/`values`. |

### `set`-only attribute

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `to` | `AnimationValue` | (none) | no | The target value. Overlay: type must match resolved `attributeName`. |

### `mpath`-specific attribute

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `href` | `IRI` | (none) — required | no | References a `path` element or shape element (SVGwg-L2) / `path` element only (SVG 1.1). Overlay: must resolve to path/shape. |

### `svg` element extensions (SVGwg-L2 §4.1)

| Attribute | Value syntax | Default | Animatable | Notes |
|-----------|-------------|---------|------------|-------|
| `playbackorder` | `"forwardonly" \| "all"` | `"all"` | no | On `svg` element. Informs UA re: backward seek. |
| `timelinebegin` | `"loadend" \| "loadbegin"` | `"loadend"` | no | On `svg` element. Controls when the timeline starts. |

---

## Open datatypes used

The following open/scalarize leaf types are referenced by the animation grammar:

| Leaf | Description | Source in grammar |
|------|-------------|-------------------|
| `NumberType` | Floating-point number (signed or unsigned) | `repeatCount`, `keyTimes` entries, `keyPoints` entries, `keySplines` values, `rotate` (number case), transform values |
| `LengthType` | CSS/SVG length with optional unit | `AnimationValue` arm; `animateMotion` coordinate |
| `PercentageType` | Number with `%` suffix | `AnimationValue` arm |
| `AngleType` | Number with optional angle unit | `AnimationValue` arm |
| `ColorType` | CSS color value | `AnimationValue` arm |
| `PaintType` | SVG paint value (`<color>`, `url()`, `none`, etc.) | `AnimationValue` arm |
| `TransformListType` | SVG transform-list | `AnimationValue` arm |
| `PathDataType` | SVG path data (d attribute syntax) | `AnimationValue` arm; `path` attribute |
| `IntegerType` | Integer (base 10) | `AnimationValue` arm |
| `StringType` | Arbitrary string | `AnimationValue` arm (for non-interpolable attributes) |
| `IriType` | IRI / URL reference | `href` attributes; `AnimationValue` arm |
| `KeywordType` | Enumeration keywords (closed per-attribute, open as union) | `AnimationValue` arm |
| `IDREF` | XML IDREF (Id-value) | Syncbase/event timing |
| `EventSymbol` | DOM event name string | Event-value timing |
| `character` | Single Unicode character | AccessKey timing |

---

## Context-sensitive constraints (overlay only)

These are NOT encoded in the grammar. They live in the parser validation pass and the generator overlay.

1. **Animation value typing (the dependent type):** `from`, `to`, `by`, `values` must have the type named by the resolved `attributeName` together with the host element. The CFG carries `AnimationValue` union; overlay selects the matching arm.

2. **Animatable-ness:** `attributeName` must name an attribute that is animatable on the resolved target element. Subsumed by type resolution (item 1).

3. **Cross-attribute cardinality:**
   - `keyTimes` count must equal `values` count (or 2 for from/to/by animations).
   - `keySplines` count must equal `keyTimes` count minus 1.
   - `keyPoints` count must equal `keyTimes` count.
   - `calcMode="spline"` requires `keyTimes` and `keySplines`.

4. **`values` vs `from`/`to`/`by` mutual exclusion:** If `values` is present, `from`, `to`, `by` are ignored. Grammar accepts both; overlay enforces the priority rule.

5. **`keyTimes` range and monotonicity:**
   - Each entry in [0.0, 1.0].
   - Non-decreasing sequence.
   - For `linear`/`spline`: first = 0, last = 1.
   - For `discrete`: first = 0.
   - If `calcMode="paced"`: `keyTimes` and `keySplines` ignored.

6. **`keyPoints` range:** Each entry in [0.0, 1.0]. Monotonicity not required.

7. **`keySplines` range:** All four numbers in each control-pt-set must be in [0.0, 1.0].

8. **Syncbase and event-base ID resolution:** The `Id-value` in syncbase and event-value timing must reference an existing animation element in the document. Overlay (special case of IDREF resolution, item 9 in CONTEXT_SENSITIVITY.md).

9. **`href` target resolution:** The IRI must reference an element capable of being animated. Overlay.

10. **`repeatCount` > 0, `dur` > 0, `min` >= 0, `max` > 0:** Numeric range constraints. Overlay.

11. **`animateTransform` `type` → value interpretation:** Overlay selects the correct `AnimateTransformSingleValue` arm based on `type` value. (CFG carries union over-approximation.)

12. **`additive`/`accumulate` ignored when target type does not support addition:** Runtime/overlay concern.

13. **`min` ≤ `max` when both specified:** Overlay; if violated, both are ignored.

---

## Discrepancies, doc gaps, and roadblocks

### D1: `attributeType` removed in SVGwg-L2, present in SVG 1.1
- SVGwg-L2 §7 explicitly lists "Removed the attributeType attribute" as a change.
- SVG 1.1 examples use `attributeType="XML"` and `attributeType="CSS"` extensively.
- MDN documents `attributeType` as present with values `"CSS" | "XML" | "auto"`.
- **Decision:** Grammar does not include `attributeType` as a first-class attribute for the SVGwg-L2 target. If SVG 1.1 compatibility mode is needed, add it as an optional attribute with `"CSS" | "XML" | "auto"` enum. Record in CLAUDE.md.

### D2: `animateColor` — deprecated in SVG 1.1, absent in SVGwg-L2
- SVG 1.1 §19.2.15 defines `animateColor` but marks it deprecated.
- SVGwg-L2 §7 explicitly lists "Removed the animateColor element."
- SMIL Animation §4.4 defines it.
- MDN lists it as deprecated/non-standard.
- **Decision:** `animateColor` is excluded from the grammar. Note in CLAUDE.md for legacy fallback.

### D3: `discard` element — out of scope for primary spec
- `discard` is not defined in SVGwg-L2 or SVG 1.1. It appears in the SVG 2 CR draft.
- MDN documents it; Chrome support is partial.
- Original scope in the task prompt includes `discard`.
- **Decision:** Include a minimal stub production for `discard` in the grammar with attributes `href` (IRI) and `begin` (begin-value-list). Flag as "SVG 2 only / limited browser support." Verification needed.

### D4: `mpath` href target scope — SVG 1.1 vs SVGwg-L2
- SVG 1.1 §19.2.14: `mpath` references a `'path' element`.
- SVGwg-L2 §2.15: "SVG allows an 'animateMotion' element to contain a child 'mpath' element which references an SVG 'path' element **or shape element**."
- SVGwg-L2 §7 explicitly notes "Allowed the 'mpath' element to refer to a shape element."
- MDN: "refers to a `<path>` element" — MDN is behind the spec.
- **Decision:** Grammar types `mpath href` as `IRI`; overlay checks target type against `{path, circle, ellipse, rect, line, polyline, polygon}`. Note the SVG 1.1/MDN discrepancy.

### D5: `playbackorder` and `timelinebegin` on `svg` element
- These are defined in SVGwg-L2 §4.1 as extensions to the `svg` element.
- Not in SVG 1.1 or the SMIL spec.
- MDN does not document these.
- Browser support unknown.
- **Decision:** Include in grammar as SVGwg-L2-only extensions on the `svg` element production. Flag as low-confidence/limited browser support.

### D6: SVGwg-L2 S production — #xC addition
- SVGwg-L2 §2.9 editorial note: "Align with whitespace used in CSS and SVG, adding #xC to S."
- SMIL spec S definition does not include #xC.
- **Decision:** Include #xC in S for the SVGwg-L2 target grammar as indicated.

### D7: `offset-value` parenthesis irregularity in SVGwg-L2 text
- SVGwg-L2 §2.9 shows: `offset-value ::= ( S? "+" | "-" S? )? ( Clock-value )`
- SMIL §3.6.7 (normative) shows: `offset-value ::= (( S "+" | "-" S )? ( Clock-value )`
- Both are slightly irregular in parenthesization. The intent is: optional `+` or `-` preceded/followed by optional whitespace, then a Clock-value.
- **Decision:** EBNF-ready form uses: `offset-value = [ S, ( "+" | "-" ), S ], Clock-value ;` which captures the intent correctly.

### D8: `repeatCount` value type — "floating point" vs `NumberType`
- Both specs call `repeatCount` a "floating point" numeric value (> 0).
- The grammar uses `NumberType` (the open scalarize leaf for numbers).
- `"indefinite"` is the other allowed value.
- **Decision:** `repeatCount = NumberType | "indefinite"`. Range > 0 is overlay.

### D9: SMIL deprecation in browsers — browser support status
- Chrome deprecated SMIL in 2015, then un-deprecated (unflagged) it in 2016. As of 2024-2025, SMIL animations are supported in Chrome/Edge (Blink), Firefox (Gecko), and Safari (WebKit).
- `animateMotion` with `mpath` is broadly supported.
- `discard` has very limited support.
- `accessKey` timing: limited/inconsistent implementation.
- `wallclock-sync-value`: not implemented in any major browser (unresolvable in practice due to privacy concerns).
- **Grammar decision:** All forms are included in the grammar (they are syntactically valid); browser support gaps are runtime/implementation concerns, not grammar concerns.

### D10: `begin-value` in SMIL includes `media-marker-value` not in SVG
- SMIL §3.2.1 (prose) lists `media-marker-value` as a begin-value form.
- SVGwg-L2 §2.9 and SVG 1.1 §19.2.8 do NOT include `media-marker-value` in the SVG begin-value syntax.
- The SMIL §3.6.7 normative BNF also omits `media-marker-value` from the grammar (only the prose mentions it).
- **Decision:** `media-marker-value` is excluded from the SVG animation grammar. SVG does not define media markers.

### D11: `Event-symbol` is an open set
- The set of legal event names is defined by the host language. SVGwg-L2 §2.9 says "The list of event-symbols available for a given event-base element is the list of event attributes available for the given element as defined in the Scripting and Interactivity chapter, with the one difference that the leading 'on' is removed from the event name."
- This means `click`, `mouseenter`, `focus`, `beginEvent`, `endEvent`, `repeatEvent`, etc.
- The set is large and can include dynamically-created events.
- **Decision:** `EventSymbol` is a scalarize leaf (open domain). A representative sample list is provided in `reps.go`-style: `{click, mouseenter, mouseleave, mouseover, mouseout, mousedown, mouseup, focus, blur, keydown, keyup, beginEvent, endEvent, repeatEvent, load}`.

---

## 15-line summary

1. **Elements in scope:** `animate`, `set`, `animateMotion`, `mpath`, `animateTransform`, `discard` (stub). `animateColor` excluded (removed in SVGwg-L2).
2. **`animate` content model:** descriptive elements + `script`, zero or more.
3. **`set` excludes** `additive`, `accumulate`, `calcMode`, `values`, `keyTimes`, `keySplines`, `from`, `by`; carries only `to` plus timing/target attrs.
4. **`animateMotion` default `calcMode` is `"paced"`** (differs from other elements where it is `"linear"`).
5. **`mpath`** references a path or shape element (SVGwg-L2); only path elements (SVG 1.1). `href` required, target resolution is overlay.
6. **`animateTransform` type enum** `"translate" | "scale" | "rotate" | "skewX" | "skewY"` is fully enumerated; value shapes per type are specialized productions (context-free).
7. **`attributeType`** (`"CSS" | "XML" | "auto"`) is SVG 1.1 only; removed in SVGwg-L2.
8. **SMIL timing BNF** is fully context-free and belongs in the grammar: offset, syncbase, event, repeat, accessKey, wallclock, "indefinite" — all transcribed verbatim and rendered to EBNF-ready form.
9. **Clock-value grammar** is identical across all three sources; SVGwg-L2 adds `#xC` to the whitespace production `S`.
10. **`AnimationValue` union** over-approximates `from`/`to`/`by`/`values` types; overlay resolves `attributeName` to narrow it.
11. **Overlay constraints:** value typing, keyTimes/values/keySplines cardinality, calcMode=spline requires keyTimes/keySplines, values-vs-from/to/by priority, keyTimes/keyPoints range+monotonicity, syncbase/event id resolution, href target resolution, repeatCount/dur/min/max range checks.
12. **`wallclock-sync-value`** is syntactically valid in grammar; not implemented in any major browser.
13. **`discard`** is SVG 2 only, not in SVGwg-L2 or SVG 1.1; stub included with low-confidence flag.
14. **`playbackorder`** and `timelinebegin`** are SVGwg-L2 extensions to the `svg` element; enumerated as two-value closed sets.
15. **No `media-marker-value`** in SVG's timing grammar (SVG does not define media markers despite SMIL prose mentioning it).
