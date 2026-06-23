# Animation Elements — Visual QA Round 1

Screenshotted 2026-06-22. All 6 pages snapped at SNAP_FULLPAGE=1 SNAP_SCALE=1.
Host shape: `<rect width="40" height="40" fill="#4d8bff">` (animate / set / animateTransform / animateMotion / discard)
or `<rect width="20" height="20" fill="#f5a623">` (mpath).

---

## `<animate>` — 195 cards

**Category: MOSTLY_PERFECT**

### Visual summary

The host rect is visible in every card. No blank cards. Cards are visually distinct in the
areas that matter (animated position shifts are mid-animation at screenshot time, so most
cards show the rect at a mid-point on the x-axis — that is expected and fine).

### Issues found

#### 1. `attributeName` values are non-SVG identifiers (GRAMMAR_ISSUE)
Cards: `attributeName="circle1"`, `attributeName="grad-a"`, `attributeName="myId"`,
`attributeName="node3"`, `attributeName="r1"`

The grammar emits arbitrary NCName-style tokens as the `attributeName` value. These names
(`circle1`, `grad-a`, etc.) do not correspond to any attribute of `<rect>`, so the browser
silently ignores the animation — the rect stays static at its initial position (x=0).
The rect IS visible, but there is no observable effect. This is technically correct
grammar coverage (attributeName accepts any string), but these cards cannot demonstrate
the attribute's purpose.

Fix target: **blueprint.go** `baselineFor("animate")` — when `varyingPrefix == "attributeName"`,
leave `attributeName` out of the baseline as it already does, but ensure the grammar seeds
for `attributeName` are restricted to real animatable SVG attribute names (e.g. `x`, `y`,
`fill`, `opacity`, `r`, `cx`) rather than arbitrary NCNames. Or add a note in overlay.go
`animationValueFor` to seed meaningful `attributeName` values.
Alternatively fix in **overlay.go** by having `animatedAttributeNames()` return a curated
list: `["x", "y", "r", "cx", "cy", "opacity", "fill", "stroke-width"]`.

#### 2. `dur` values with invalid clock-value formats (GRAMMAR_ISSUE)
Cards: `dur="1:-1:100.3"`, `dur="-1:100.3"`, `dur="1.-1h"`, and several others with
negative components such as `-1:100.3`.

SMIL clock-values must have non-negative components (hours, minutes, seconds are unsigned).
Negative clock parts like `-1:100.3` or `1.-1h` are malformed — browsers treat them as
invalid and the animation does not run (rect stays static). The rect is still visible, but
the card demonstrates a broken value rather than a valid grammar variant.

The same issue appears for `min`, `max`, `repeatDur`, `begin`, and `end` when they embed
time values with negative components.

Fix target: **overlay.go** clock-value generator — clamp numeric components to >= 0.
Specifically `clockValue`, `timecountValue` should not produce negative sub-values.
Also guard in the EBNF (if there is one for clock-value) so the grammar only walks
non-negative digit sequences for hours/minutes/seconds sub-parts.

#### 3. `from`/`to`/`by` values typed to wrong domain (GRAMMAR_ISSUE)
Cards: `from="0deg"`, `from="45deg"`, `from="#e94560"`, `from="red"`, `from="none"`,
`from="currentColor"`, `from="context-fill"`, `from="translate(20 10)"`, `from="rotate(45)"`,
`from="M10 10 L90 90"`, etc. (and matching `to`/`by` variants).

The baseline fixes `attributeName="x"`, which expects a numeric/length value. Supplying
`from="red"` (a color) or `from="rotate(45)"` (a transform) with `attributeName="x"` is
a type mismatch — the browser ignores the animation, so the rect is static. Host is
visible but animation has no effect; cards are IDENTICAL_NO_EFFECT for all type-mismatched
from/to/by groups.

Groups with no visible effect:
- All `from`/`to`/`by` = color values (`#e94560`, `red`, `none`, `currentColor`,
  `context-fill`, `context-stroke`)
- All `from`/`to`/`by` = degree values (`0deg`, `45deg`)
- All `from`/`to`/`by` = transform functions (`translate(...)`, `rotate(...)`,
  `scale(...)`, `skewX(...)`)
- All `from`/`to`/`by` = path data (`M10 10 L90 90`, `M10 50 Q50 10 90 50`, etc.)
- `from`/`to`/`by` = reference values (`url(#slot)`, `#target`, `blur1`, `result1`)
- `from`/`to`/`by` = text tokens (`label`, `Aa`)

Fix target: **overlay.go** `animationValueFor(attributeName, ...)` — when
`attributeName == "x"` (or any length attribute), restrict `from`/`to`/`by` seed values
to numerics and length units (`10`, `60`, `24px`, `50%`). For a richer showcase, the
baseline could be varied: different `attributeName` values paired with type-appropriate
from/to seeds in a lookup table inside `animationValueFor`.

#### 4. `dur="media"` / `min="media"` etc. — media duration (minor)
`dur="media"` is only valid for media elements; on SVG animations the browser falls back
to the element's implicit duration (which for a non-media SVG animation is `indefinite`
until `begin` fires). The rect remains visible but the animation timing is undefined.
Not a broken card but semantically incorrect for SVG.

Fix target: **overlay.go** — remove `"media"` from the clock-value alternatives pool for
SVG animation elements (it is defined in SMIL but not meaningful in SVG 2 context).

---

## `<set>` — 129 cards

**Category: MOSTLY_PERFECT**

### Visual summary

All cards show the host rect. No blank cards. Identical-position rects appear across many
cards (all the `to` type-mismatch cards look like the same static rect at x=0).

### Issues found

#### 1. Same `attributeName` non-SVG-attribute issue as `<animate>` (GRAMMAR_ISSUE)
Cards: `attributeName="circle1"`, `"grad-a"`, `"myId"`, `"node3"`, `"r1"`.
Same root cause as animate. Fix same as above.

#### 2. Same invalid `dur`/`min`/`max`/`repeatDur` negative clock components (GRAMMAR_ISSUE)
Cards: `dur="10:0:1.-1"`, `dur="100:3:10.0"` (valid but huge), `dur="1:-1.100"`,
`dur="1.-1min"`, etc.
Negative sub-values make the duration invalid; browser ignores timing.
Fix same as animate item 2.

#### 3. `to` values type-mismatched with `attributeName="x"` (GRAMMAR_ISSUE / IDENTICAL_NO_EFFECT)
Same groups as animate item 3 — all color, transform, path-data, and reference
`to` values produce a static rect identical to the zero-animation state.

Fix same as animate item 3.

#### 4. `fill="none"` on `<set>` (GRAMMAR_ISSUE — wrong attribute)
Card: `fill="none"`.

`fill` on a `<set>` means the animation fill behavior (`freeze`/`remove`). The value
`"none"` is not a valid animation-fill value; it is a paint value inherited from the
presentation attribute namespace. Browsers treat this as `remove` (the default) or ignore
it. The card is visually fine but represents a grammar error — `fill` for animation
elements should only enumerate `freeze` and `remove`.

Fix target: **overlay.go** or the EBNF for `<set>`'s `fill` attribute — restrict
enumeration to `["freeze", "remove"]` and exclude paint values (`none`, color strings,
etc.). The same issue exists on `<animate>`, `<animateMotion>`, `<animateTransform>`.

#### 5. Presentation attributes surfaced on `<set>` (GRAMMAR_ISSUE)
Cards: `fill-rule="nonzero"`, `fill-opacity="0.5"`, `stroke="none"`, `stroke-opacity`,
`stroke-width`, `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`,
`stroke-dasharray`, `stroke-dashoffset`, `paint-order`, `marker`, `marker-start`,
`marker-mid`, `marker-end`, `color`, `color-interpolation`, `color-rendering`,
`shape-rendering`, `vector-effect`, `filter`, `color-interpolation-filters`,
`flood-color`.

These are presentation/styling attributes. They are technically inheritable on any SVG
element (SVG 2 allows them broadly), so the browser accepts them without error, and the
rect is visible. However, they do not interact with the `<set>` animation mechanics in any
useful or observable way. The cards all look identical to the baseline. This inflates the
gallery with many visually indistinct cards.

Fix target: **overlay.go** `attributesFor("set")` — restrict the attribute set to
`<set>`-specific attributes (`attributeName`, `to`, `begin`, `end`, `dur`, `min`, `max`,
`restart`, `repeatCount`, `repeatDur`, `fill`, `id`, `href`, `tabindex`, global XML
attrs). Remove inherited presentation attributes from the enumeration for animation
elements.

---

## `<animateMotion>` — 123 cards

**Category: MOSTLY_PERFECT**

### Visual summary

The majority of cards show the host rect. One notable exception visually: cards with
`path="M10 10 L90 90"`, `path="M10 50 Q50 10 90 50"`, `path="M20 20 H80 V80 H20 Z"` —
the rect moves along those paths, so at screenshot time it may appear in various positions
within the card, but it is still visible. One card with `path` starting at `M10 10` places
the rect near the top-left corner of the card, which is partly clipped but still visible.

Cards for `dur="1.-1s"`, `dur="100.3ms"` (fast), `dur="100:3:10.0"` (huge duration),
`dur="1:-1.100"` (invalid negative) follow the same clock-value issues as animate.

The `values="1,-1;3.14,0.5;2,0"` card shows a slightly different rect position (the
animateMotion `values` are motion coordinates — small numbers near origin, so the rect
barely moves from (0,0)).

### Issues found

#### 1. Invalid clock-value negative components (GRAMMAR_ISSUE)
Same as `<animate>` item 2. Fix same.

#### 2. `from`/`to`/`by` for animateMotion use coordinate pairs but path baseline conflicts (GRAMMAR_ISSUE)
The baseline already sets `path="M0 0 L60 60"`. When `from="1,-1"` or `to="2,0"` is also
set, the browser is given both `path` and `from`/`to` — per SMIL, `path` takes precedence
and `from`/`to` are ignored. The card is visually indistinguishable from the `path`-only
baseline cards.

Fix target: **blueprint.go** `baselineFor("animateMotion")` — when the varying attribute
is `from`, `to`, `by`, or `values`, omit `path` from the baseline so from/to/by are
actually effective.

#### 3. `values="1,-1;3.14,0.5;2,0"` — tiny motion coordinates (minor)
These motion coordinates are near (0,0) in the 100x100 viewBox. The rect moves only a
few pixels. Marginally useful. Suggest larger coordinate pairs in overlay.go's motion
value generator, e.g. `"10,10;50,5;80,50"`.

Fix target: **overlay.go** motion values seed — use on-canvas coordinate pairs.

#### 4. Same presentation-attribute inflation as `<set>` (GRAMMAR_ISSUE)
`fill-rule`, `stroke`, `stroke-width`, `fill-opacity`, etc. appear as cards (bottom
portion of the gallery). Same diagnosis as set item 5 — these are visually indistinct
from the baseline.

Fix target: **overlay.go** `attributesFor("animateMotion")`.

#### 5. `dur="media"` (minor, same as animate item 4)
Fix same.

---

## `<animateTransform>` — 146 cards

**Category: MOSTLY_PERFECT**

### Visual summary

All cards show the host rect, usually rotated (the baseline is a `rotate` animation).
Cards are clearly visually distinct: the rect appears at various rotated angles at
screenshot time. This is the strongest-performing gallery in the batch.

### Issues found

#### 1. `attributeName` non-transform values (GRAMMAR_ISSUE)
Cards: `attributeName="circle1"`, `"grad-a"`, `"myId"`, `"node3"`, `"r1"`.

`animateTransform` requires `attributeName="transform"` to function. Setting it to
`"circle1"` etc. produces a no-op animation. The rect is visible at initial position (not
rotating). Cards look like static rotated-at-time-0 rects.

Fix target: **overlay.go** `attributesFor("animateTransform")` — for `attributeName`,
restrict seed values to `["transform"]` only. Or in **blueprint.go**
`baselineFor("animateTransform")`, always inject `attributeName="transform"` and
never allow the grammar to vary it to arbitrary strings (make `attributeName` a
`<animateTransform>`-specific constant in the overlay).

#### 2. `type` varied but baseline `from`/`to` are rotate-format values (GRAMMAR_ISSUE / IDENTICAL_NO_EFFECT)
Cards: `type="translate"`, `type="scale"`, `type="skewX"`, `type="skewY"`.

The baseline uses `from="0 25 25" to="360 25 25"` (rotate-format with cx/cy). When
`type="translate"`, those values are treated as translation distances, causing the rect
to animate from (0,25) to (360,25) — the rect flies off the right edge and is not visible
at mid-animation in the screenshot. When `type="scale"`, `from="0 25 25"` means
scale(0) which collapses the rect to zero size at time=0 (blank-looking card at that frame).

Affected cards:
- `type="translate"` — rect may be off-canvas at screenshot time (partial blank)
- `type="scale"` — rect may be scaled to near-zero (effectively blank at start)
- `type="skewX"` / `type="skewY"` — rect skewed off canvas

Fix target: **overlay.go** or **blueprint.go** — when varying `type`, pair it with
type-appropriate from/to seeds:
- `translate`: `from="0 0" to="50 0"` (stays on canvas)
- `scale`: `from="1" to="2"` (grows, stays visible)
- `skewX`: `from="0" to="30"` (mild skew)
- `skewY`: `from="0" to="30"`

This requires `overlay.go` to detect `type` as the varying attribute and inject matching
`from`/`to` values accordingly.

#### 3. `from="3.140.5 2"` — malformed number (GRAMMAR_ISSUE)
Card: `from="3.140.5 2"` and matching `to`/`by` variants.

`3.140.5` is not a valid SVG number (two decimal points). This is a grammar bug — the
walk produced an invalid number literal. The browser will likely reject the entire
`from` value and the animation will not run.

Fix target: **overlay.go** or the underlying number-generation grammar — the number
tokenizer must not allow two decimal points in a single numeric token. Check EBNF for
`number` or `coordinate` production.

#### 4. Invalid negative clock components (same as animate item 2)
Cards: `dur="-1:100.3"`, `dur="0:1:-1.100"`, `dur="-1.100min"`, etc.
Fix same as animate.

#### 5. Same presentation-attribute inflation (same as set item 5)
Fix same.

---

## `<mpath>` — 30 cards

**Category: HAS_EMPTY_CARDS**

### Visual summary

All 30 cards show the host rect (a small 20x20 blue-orange square). However, none of the
cards show the rect animating along the path because the `<mpath>` element in every card
does not have a working `href`. The `defs` block defines `<path id="slot-path" ...>` but
the `<mpath>` children use `href="#target"` (which does not exist), `id="circle1"`, etc.
— none of which resolve to the defined path.

The rect is therefore stationary in all cards. Since the rect is small (20x20) and placed
at `x=0 y=0`, it is barely visible in the top-left corner of each card's SVG viewport.

**Blank/barely-visible cards (all 30):** Every card. The rect is technically rendered but
sits at the top-left corner and does not animate because `<mpath>` cannot resolve its
path reference.

### Issues found

#### 1. `href="#target"` points to non-existent element (GRAMMAR_ISSUE / HAS_EMPTY_CARDS)

The blueprint (blueprint.go `catMpath` case) defines `<path id="slot-path" ...>` in
`<defs>`, but the first (working) `<mpath>` card uses `href="#target"` which does not
exist. The blueprint path id is `"slot-path"`, not `"target"`. This means the
`href="#target"` card has a broken reference and the animateMotion does not follow any
path.

For the non-href cards (`id="circle1"`, `tabindex="0"`, `lang="en"`, etc.), the `<mpath>`
has no `href` at all, so the motion path is also undefined.

Fix target: **blueprint.go** `catMpath` scaffold — change the `<path>` id from
`"slot-path"` to `"target"` so that `href="#target"` resolves correctly in the first
card. All other cards (which vary attributes other than `href`) will then also benefit
because the containing `<animateMotion>` can still follow `#target`.

Currently:
```
`<defs><path id="slot-path" d="M10 50 Q50 10 90 50"/></defs>`
```
Should be:
```
`<defs><path id="target" d="M10 50 Q50 10 90 50"/></defs>`
```

Also fix the `<animateMotion>` wrapper in the same scaffold: the animateMotion should
include `href="#target"` on mpath, not leave it up to the varied attribute alone.
A better scaffold:
```
<defs><path id="target" d="M10 50 Q50 10 90 50"/></defs>
<rect width="20" height="20" fill="#f5a623">
  <animateMotion dur="2s" repeatCount="indefinite"><mpath href="#target">{{ELEMENT-ATTRS-HERE}}</mpath></animateMotion>
</rect>
```

Actually the correct fix is simply: in the default blueprint, have a working `<mpath
href="#target">` with the `{{ELEMENT}}` injected around the mpath's attributes, OR rename
the defs path to `id="target"`.

Simplest fix: in `blueprint.go` line 134, change `id="slot-path"` to `id="target"`:
```go
`<defs><path id="target" d="M10 50 Q50 10 90 50"/></defs>` +
`<rect x="0" y="0" width="20" height="20" fill="#f5a623"><animateMotion dur="2s">{{ELEMENT}}</animateMotion></rect></svg>`
```

#### 2. Rect starts at (0,0) and is clipped at top-left (UI_ISSUE)
Even when the motion path would resolve, the rect starts at (0,0) before the first tick.
With `x=0 y=0` and `width=20 height=20`, the rect is partially clipped by the card
border at (0,0). Change the path start point to begin mid-canvas, e.g.
`d="M30 50 Q50 10 70 50"` so the rect is visually centered on motion start.

Fix target: **blueprint.go** `catMpath` scaffold — adjust the path data and possibly
add `x="30" y="30"` to the rect.

---

## `<discard>` — 28 cards

**Category: HAS_EMPTY_CARDS**

### Visual summary

The gallery shows 28 cards. In the screenshot, approximately 22-24 cards show the rect.
However the remaining 4-6 cards appear to show a blank SVG area where the rect should be
but is not visible.

The key issue: `<discard>` is designed to remove its parent element from the DOM after a
trigger. The baseline blueprint for `<discard>` uses the same `catAnimation` scaffold as
`<animate>`:

```
<rect id="target" x="10" y="10" width="40" height="40" fill="#4d8bff"><discard>{{ELEMENT}}</discard></rect>
```

This means `<discard>` is immediately active with no `begin` (default begin=0), so the
rect discards itself instantly on page load. Cards where the `begin` attribute fires
immediately (or has an immediately-satisfied condition) will show a blank SVG canvas.

**Blank cards (discard fires immediately):**
- `begin="+−1.100min;circle1.begin+3:10:0.1;grad-a.click+−1:100.3"` — the first term
  `+−1.100min` is a negative offset from now, which fires immediately.
  → Rect discarded, blank canvas.
- Any `begin` value that resolves to t<=0 will fire and discard the rect.

Cards with well-formed positive `begin` values (e.g. event-based begins like
`circle1.begin+3:10:0.1`) that won't fire at static screenshot time correctly show the
rect.

### Issues found

#### 1. `<discard>` fires immediately on several cards — rect discarded before screenshot (HAS_EMPTY_CARDS)
The default blueprint gives `<discard>` no `begin`, meaning it fires at t=0. For most
cards, the varied attribute is not timing-related (`id`, `tabindex`, `lang`, etc.),
so `begin` defaults to 0 and the rect is discarded instantly on page load.

However in the screenshot most of those cards DO still show the rect — this is because
the `snap.sh` tool uses a brief `settle` wait, and `<discard>` with `begin="0"` is
synchronous. Looking more carefully at the screenshot, the rect IS present in most
cards (Chrome may delay discard until next frame). The blank-looking cards appear to
be those with `begin` values that resolve to a negative offset (which fires at t=0).

Fix target: **blueprint.go** `catAnimation` default scaffold for `<discard>` specifically
— add a large `begin` on the `<discard>` element so it never fires during the screenshot
settle period:
```go
case catAnimation:
    if tag == "discard" {
        return svgOpen +
            `<rect id="target" x="10" y="10" width="40" height="40" fill="#4d8bff">` +
            `{{ELEMENT}}</rect></svg>`
        // Note: baseline for discard in baselineFor() should inject begin="60s" so it
        // never fires during the 1-2s screenshot window.
    }
```
And in **blueprint.go** `baselineFor("discard")` (currently missing — falls through to `""`):
add:
```go
case "discard":
    return add([2]string{"begin", "60s"}), false
```

This ensures the rect is always visible at screenshot time for non-begin-varying cards.
For cards varying `begin`, the grammar-generated begin values should be audited to exclude
immediately-firing negative offsets (same as clock-value fix above).

#### 2. `begin` values with complex syntax that are hard to reason about (GRAMMAR_ISSUE)
Card: `begin="accessKey(specimen)+0:1.-1;wallclock(100-3-10T0:1Z);indefinite"` — the
`accessKey` begin never fires without keyboard interaction, `wallclock` requires an
absolute time in the past/future, and `indefinite` means it waits for a script trigger.
None fire during screenshot. Rect remains visible (correct for QA) but this card does
not demonstrate `<discard>` working.

This is acceptable for grammar coverage but worth noting.

---

## Cross-Cutting Issues (Recurring Across All 6 Elements)

### A. Invalid clock-value grammar: negative sub-components
**Affects:** `animate`, `set`, `animateMotion`, `animateTransform` (and `discard` begin values).
**Attributes:** `dur`, `min`, `max`, `repeatDur`, `begin`, `end`.
**Examples:** `-1:100.3`, `1.-1h`, `0:1:-1.100`, `1:-1.100`.

SMIL clock-value BNF requires all numeric fields to be unsigned. Negative components
produce invalid values that browsers reject, leaving the animation non-functional.

**Fix target: overlay.go** — in the clock-value / timecount-value generator, ensure all
numeric parts are clamped to >= 0. If the grammar EBNF for clock-value exists as a
separate `.ebnf` file, add `[0-9]+` instead of allowing a leading `-` on sub-components.

### B. Presentation attributes on animation elements
**Affects:** `animate`, `set`, `animateMotion`, `animateTransform`.
**Cards:** `fill-rule`, `fill-opacity`, `stroke`, `stroke-opacity`, `stroke-width`,
`stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit`, `stroke-dasharray`,
`stroke-dashoffset`, `paint-order`, `marker`, `marker-start`, `marker-mid`, `marker-end`,
`color`, `color-interpolation`, `color-rendering`, `shape-rendering`, `vector-effect`,
`filter`, `color-interpolation-filters`, `flood-color`.

These inflate the gallery with ~20 identical-looking cards per element (all show the same
static-positioned rect with the presentation attribute having no visible effect on the
animation).

**Fix target: overlay.go** `attributesFor(tag)` for animation elements — strip
presentation / painting attributes that have no observable impact when placed on an
animation child element.

### C. `fill="none"` and animation fill attribute vs. paint attribute collision
**Affects:** `animate`, `set`, `animateMotion`, `animateTransform`.

The `fill` attribute on animation elements means animation fill behavior (`freeze`/`remove`).
The grammar additionally seeds `"none"` (a paint value) which is invalid for animation
fill. Cards for `fill="none"` are structurally fine but carry an invalid value.

**Fix target: overlay.go** — for the `fill` attribute on animation elements, restrict
enumeration to `["freeze", "remove"]`.

---

## Summary Table

| Element          | Cards | Category           | Top Issues |
|------------------|-------|--------------------|------------|
| animate          | 195   | MOSTLY_PERFECT     | Invalid clock-values; from/to type mismatches; attributeName non-SVG names; presentation attr inflation |
| set              | 129   | MOSTLY_PERFECT     | Same as animate; fill="none" grammar error; presentation attr inflation |
| animateMotion    | 123   | MOSTLY_PERFECT     | Invalid clock-values; path overrides from/to; tiny motion coords; presentation attr inflation |
| animateTransform | 146   | MOSTLY_PERFECT     | type-mismatched from/to (off-canvas); malformed number `3.140.5`; attributeName non-transform values |
| mpath            | 30    | HAS_EMPTY_CARDS    | All cards: href resolves to missing element (`#target` vs `#slot-path`); rect barely visible at (0,0) |
| discard          | 28    | HAS_EMPTY_CARDS    | No baseline begin="60s"; negative begin offsets fire immediately; rect discarded before screenshot |

## Priority Fix Order

1. **mpath / blueprint.go** — rename `id="slot-path"` to `id="target"` (one-line fix, makes all 30 cards live).
2. **discard / blueprint.go + baselineFor** — add `begin="60s"` baseline so rect survives screenshot.
3. **overlay.go clock-values** — clamp all clock sub-components >= 0 (fixes ~30-40 cards across 4 elements).
4. **overlay.go animationValueFor** — restrict from/to/by to values compatible with `attributeName="x"` (or pair attributeName with type-appropriate seeds).
5. **overlay.go attributesFor(animation elements)** — remove presentation attrs from animation element galleries.
6. **animateTransform type-paired from/to** — when type is being varied, pair with type-appropriate from/to (blueprint.go or overlay.go).
7. **overlay.go fill** — restrict `fill` on animation elements to `["freeze", "remove"]` only.
