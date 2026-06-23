# Visual QA — Round 1: pservers + masking batch

**Date:** 2026-06-22  
**Elements reviewed:** `linearGradient`, `radialGradient`, `stop`, `pattern`, `marker`, `clipPath`, `mask`  
**Method:** Full-page screenshot (`SNAP_FULLPAGE=1 SNAP_SCALE=1`) + HTML source inspection per card.

---

## `<linearGradient>` — MOSTLY_PERFECT

**80 cards. Gradient effect visible and url(#slot) resolves in every card.**

### Issues found

#### 1. Duplicate `id` attribute (GRAMMAR_ISSUE)
Cards for the `id` attribute group (`id="circle1"`, `id="grad-a"`, `id="myId"`, `id="node3"`, `id="r1"`) emit:
```html
<linearGradient id="slot" id="circle1">
```
Two `id` attributes on the same element is invalid XML/HTML. The first attribute wins in most parsers, so `url(#slot)` still resolves and the card renders, but the intent (showing the `id` attribute) is broken — the second `id` is silently dropped.

**Affected cards:** `id="circle1"`, `id="grad-a"`, `id="myId"`, `id="node3"`, `id="r1"` (same pattern on all elements that use `id` from `CoreAttribute`).

**Fix target:** `blueprint.go` → `defaultScaffold` for `catGradient`. Instead of injecting `id="slot"` as a hardcoded element attribute, supply the slot id through the referencing rect's `fill="url(#slot)"` and let the generator-emitted element carry whatever `id` it needs. Alternative: when the varying attribute IS `id`, the overlay in `overlay.go` should detect the conflict and override the element's `id` to equal `slotID`, skipping the blueprint's injected `id`. This is the "id-shadowing" fix — add a guard in the code that injects baseline attributes that checks `if varying == "id" { use varying value as slot id }`.

#### 2. Self-referencing `href`/`xlink:href` (GRAMMAR_ISSUE)
Cards `href="#slot"` and `xlink:href="#slot"` emit:
```html
<linearGradient id="slot" href="#slot">
```
A gradient referencing itself via `href` is a circular reference — browsers silently ignore the inheritance and render the gradient as if no `href` is set. Since the gradient still has its own stops, the card renders fine visually, but it demonstrates a semantically degenerate case.

**Fix target:** `overlay.go` → `overlaySample`. When `attrName == "href"` or `attrName == "xlink:href"` and `tag` is a gradient element, point the reference at a sibling gradient definition (e.g. a separately defined `<linearGradient id="base">`) rather than at itself. Currently `IriType` always returns `"#" + refTarget` which equals `"#slot"` — that is self-referential for gradient elements.

#### 3. Degenerate gradient endpoints (`x2="0"`, `y2="0"`) (IDENTICAL_NO_EFFECT)
When `x2="0"` the gradient vector becomes `(0%,0) → (0,0)` (zero-length horizontal extent in `objectBoundingBox` space), producing a solid color — visually indistinguishable from the `x2="1"` (default) card.

**Affected attr groups:** `{x2="0"}`, `{y2="0"}`.

**Fix target:** `reps.go` → `CoordinateType` / `LengthPercentageType` samples used for gradient axes: replace the `"0"` sample for `x2`/`y2` with `"50"` or `"50%"` so the endpoint is always non-degenerate. Alternatively add a guard in `overlay.go` `isNonNegativeAttr` that keeps gradient `x2`/`y2` > 0.

---

## `<radialGradient>` — MOSTLY_PERFECT

**All cards show a pink-to-teal radial blob effect. url(#slot) resolves in every card.**

### Issues found

Same duplicate-`id` issue as `linearGradient` (cards `id="circle1"`, `id="grad-a"`, `id="myId"`, `id="node3"`, `id="r1"`) — second `id` silently dropped.

Same self-referencing `href`/`xlink:href` issue.

No empty or invisible cards observed.

**Fix target:** Same as linearGradient: `blueprint.go` id-shadowing guard and `overlay.go` href redirect.

---

## `<stop>` — HAS_EMPTY_CARDS (ALL 56 cards are blank/black)

**Primary category: HAS_EMPTY_CARDS**

**All 56 cards are solid black.** The generated `<stop>` element has no `stop-color` attribute and defaults to black (transparent alpha = 0 in most SVG engines = shows the dark SVG background). The scaffold places the varied stop at the end of the gradient:

```html
<linearGradient id="slot">
  <stop offset="0" stop-color="#e94560"/>  <!-- baseline stop -->
  <stop offset="0"></stop>                  <!-- varied stop: NO stop-color -->
</linearGradient>
```

Because the varied stop has no `stop-color`, it defaults to black/transparent. Only the baseline stop at offset=0 contributes colour, but since the varied stop overrides offset 0 or follows at offset 1 with black, the visible gradient fades from red/pink to black — too dim on the dark card background to see the gradient effect clearly. For many attribute paths (e.g. `tabindex`, `lang`, `xml:lang`, etc.) the stop has no colour-related attributes at all, so the entire rect appears as a gradient from red-ish to black against a dark background — which reads as near-black.

### Root cause
`blueprint.go` `bodyFor("stop")` returns `""` (empty), which is correct since the stop is the generated element. However, the **baseline stop** at offset=0 in `catStop` scaffold provides a red anchor, but the **varied stop** at offset=1 needs its own `stop-color` in the baseline attributes so every non-colour-varying card still shows a visible 2-colour gradient.

**Fix target — `blueprint.go` `catStop` scaffold:**

Change the scaffold from:
```go
`<defs><linearGradient id="slot"><stop offset="0" stop-color="#e94560"/>{{ELEMENT}}<stop offset="1" stop-color="#16c79a"/></linearGradient></defs>`
```
to placing the varied element as the **middle stop** with the two coloured boundary stops intact:
```go
`<defs><linearGradient id="slot"><stop offset="0" stop-color="#e94560"/>{{ELEMENT}}<stop offset="1" stop-color="#16c79a"/></linearGradient></defs>`
```
AND ensure the baseline in `baselineFor("stop", ...)` adds `stop-color="#f5a623"` and `offset="0.5"` so the varied stop always has a visible colour independent of which attribute is being varied. Add to `baselineFor`:
```go
case "stop":
    return add([2]string{"offset", "0.5"}, [2]string{"stop-color", "#f5a623"}), false
```
This ensures even a `tabindex="0"` card shows a 3-stop gradient (red → orange → teal) making the gradient fully visible.

**Additionally:** The rect in the stop scaffold uses `width="100" height="100"` with no `x`/`y`, which means it covers the full SVG including the rounded-corner background. Fix: use `x="5" y="5" width="90" height="90"` consistent with other gradient scaffolds.

---

## `<pattern>` — MOSTLY_PERFECT

**All cards show orange circle tiling pattern. Effect is visible in every card.**

### Issues found

#### 1. `patternUnits="objectBoundingBox"` with absolute coordinates (IDENTICAL_NO_EFFECT)
Card `patternUnits="objectBoundingBox"` emits:
```html
<pattern id="slot" width="20" height="20" patternUnits="objectBoundingBox">
  <circle cx="10" cy="10" r="6" fill="#f5a623"/>
</pattern>
```
When `patternUnits="objectBoundingBox"`, `width="20"` means 2000% of the bounding box — a single enormous tile, and the circle at `cx="10"` is at 1000% of bounding box. This produces what appears to be a solid fill or a very large tile with the circle off-screen. The card visually looks different from others but may be showing an incorrect / degenerate pattern.

**Fix target:** `blueprint.go` `baselineFor("pattern", ...)`. When generating pattern baseline attributes, the baseline `width="20" height="20"` is only correct for `patternUnits="userSpaceOnUse"`. The scaffold or overlay should detect when `patternUnits="objectBoundingBox"` is being varied and switch to fractional width/height (e.g., `width="0.2" height="0.2"`) and fractional circle coordinates (`cx="0.5" cy="0.5" r="0.3"`).

Practical approach: add a template override in `chrome-testing/html/template/pattern.html` with two blueprints branching on the attribute, or add a special case in `baselineFor` that checks `varyingPrefix == "patternUnits"` and uses `objectBoundingBox`-compatible values.

#### 2. `patternContentUnits="objectBoundingBox"` mispositioned content (IDENTICAL_NO_EFFECT)
Same problem: the child circle uses absolute SVG coordinates while content units is `objectBoundingBox`, making the circle's `cx="10"`, `cy="10"`, `r="6"` refer to 10× the bounding box size — circle renders off-canvas.

**Fix target:** `blueprint.go` `bodyFor("pattern")` should produce fractional coordinates when `patternContentUnits="objectBoundingBox"` is active. Since `bodyFor` is unaware of which attribute is varying, the simplest fix is to add the same guard as above.

#### 3. Duplicate `id` attribute (same as gradient elements, GRAMMAR_ISSUE)
Cards for `id="circle1"` etc. emit `id="slot" id="circle1"` — second `id` dropped.

---

## `<marker>` — MOSTLY_PERFECT

**84 cards. Arrowhead-on-line effect is visible in every card. Effect clearly distinguishable.**

### Issues found

#### 1. `refX="left"`, `refX="center"`, `refX="right"` — keyword values cause misaligned arrowheads (GRAMMAR_ISSUE)
Cards `refX="left"`, `refX="center"`, `refX="right"`, `refY="_top"`, `refY="center"`, `refY="bottom"` show the arrowhead shifted noticeably from the line end. The grammar correctly enumerates these SVG 2 keyword values (`left`, `center`, `right`, `top`, `bottom` are valid for `refX`/`refY`). The visual shift is expected behaviour — these are semantically valid cards.

No actual grammar issue here; however `refY="_top"` is suspicious — `_top` is not a valid SVG keyword (valid are `top`, `center`, `bottom`).

**Fix target:** `lang/pservers.ebnf` or the appropriate marker EBNF — verify `_top` is not a typo for `top`. If it appears as `"_top"` in the grammar, fix it there. Check `lang/marker.ebnf` or wherever `refY` keywords are defined.

#### 2. Duplicate `id` attribute (GRAMMAR_ISSUE)
Same as other elements: `id="slot" id="circle1"` etc.

#### 3. `markerWidth="2"`, `markerHeight="2"` — tiny arrowhead (IDENTICAL_NO_EFFECT)
With `markerWidth="2"` the arrowhead is rendered 2 units wide (the default `markerUnits="strokeWidth"` scales by stroke-width=3, so the arrowhead is 6 user-units — still visible but very small). This is correct SVG behaviour; the card is distinguishable from `markerWidth="20"`.

No change needed; this is working as expected.

---

## `<clipPath>` — MOSTLY_PERFECT

**62 cards. Circle clip mask effect visible in all cards. url(#slot) resolves correctly.**

### Issues found

#### 1. Duplicate `id` attribute (GRAMMAR_ISSUE)
Cards `id="circle1"`, `id="grad-a"`, etc. emit `id="slot" id="..."`. Same root cause as other elements.

#### 2. `role="combobox complementarycontentinfo"` — concatenated tokens without space separator (GRAMMAR_ISSUE)
The card for `role` shows:
```html
role="combobox complementarycontentinfo"
```
The grammar `RoleValue = RoleToken , { RoleToken }` concatenates tokens without a space separator. `"complementarycontentinfo"` is not a valid role; it should be `"complementary contentinfo"` (two separate tokens). The grammar is missing a `" "` separator between tokens.

**Fix target:** `lang/svg.ebnf` line 82:
```
RoleValue = RoleToken , { " " , RoleToken } ;
```
Add a literal space `" "` before each additional token in the repetition.

#### 3. `fill="none"` on clipPath (GRAMMAR_ISSUE — incorrect attribute placement)
Card `fill="none"` shows a clipPath element with `fill="none"`:
```html
<clipPath id="slot" fill="none"><circle cx="50" cy="50" r="35"/></clipPath>
```
`fill` is a presentation attribute that is valid on `clipPath` per the SVG spec (as an inherited property for its content), but `fill="none"` on the clipPath container itself has no visible effect — the fill of the child circle (which has no explicit fill) would inherit `none`, making the clip region empty, and the clipped rect would disappear entirely. In Chrome the circle still clips (clipPath ignores fill for clipping geometry), so this card visually shows a normal clipped circle.

This is acceptable behaviour but worth noting: the clip effect is still visible so no empty card issue.

#### 4. `clipPathUnits="objectBoundingBox"` — circle coordinates are wrong units (IDENTICAL_NO_EFFECT)
Card `clipPathUnits="objectBoundingBox"` has `<circle cx="50" cy="50" r="35"/>` where units are in object bounding box fractions (0–1 scale). `cx="50"` means 5000% of the bounding box — the circle is entirely off-canvas and the clip region is empty, making the clipped rect invisible.

In the screenshot this card shows a fully clipped (invisible content) rect — confirming a blank clipping result.

**Fix target:** `blueprint.go` `bodyFor("clipPath")`. When `clipPathUnits="objectBoundingBox"` is the varying attribute, the child circle should use fractional coordinates: `<circle cx="0.5" cy="0.5" r="0.4"/>`. Add a parallel guard like `baselineFor` but for the child body, or add a template override.

Practical fix in `blueprint.go`: add a `case "clipPath":` branch in `baselineFor` that includes `clipPathUnits="userSpaceOnUse"` as the baseline so only the card where `clipPathUnits` is the varied attribute will produce the `objectBoundingBox` variant, and teach the blueprint how to size its children for `objectBoundingBox` mode.

---

## `<mask>` — MOSTLY_PERFECT

**75 cards. Masked square effect visible in nearly all cards.**

### Issues found

#### 1. `x="100%"`, `y="100%"` — mask region shifted entirely off-canvas (HAS_EMPTY_CARDS)
Cards `x="100%"` and `y="100%"` position the mask region at 100% offset, moving the white mask rectangle entirely off the element's bounding box. The masked rect becomes invisible.

**Affected cards:** `x="100%"`, `y="100%"`.

**Fix target:** `overlay.go` → add `"x"` and `"y"` for mask (and similar geometry elements) to `isNonNegativeAttr` AND clamp percentage values to a visible range. For `LengthPercentageType` on mask geometry attributes, the overlay should intercept and return `"20%"` rather than allowing `"100%"`.

Alternatively, remove `"100%"` from `reps["LengthPercentageType"]` and replace it with `"80%"` to keep the mask region on-canvas for all clients.

#### 2. `width="20"`, `height="20"` — tiny mask window (IDENTICAL_NO_EFFECT)
Cards `width="20"` and `height="20"` set the mask's bounding box to 20×20 user units (the mask content rect is 60×60, so it clips to a 20-unit window). The masked result shows a very small bright square — effect is still visible and distinguishable.

These cards are working correctly. No fix needed.

#### 3. `maskContentUnits="objectBoundingBox"` — mispositioned mask content (HAS_EMPTY_CARDS)
Card `maskContentUnits="objectBoundingBox"` has mask content:
```html
<rect x="20" y="20" width="60" height="60" fill="#fff"/>
```
When `maskContentUnits="objectBoundingBox"`, `x="20"` means 2000% of the element's bounding box — the white rect is off-canvas, the mask is entirely transparent, and the masked element is invisible.

**Fix target:** `blueprint.go` `bodyFor("mask")`. Add a conditional: when this attribute is varying, use fractional coordinates `<rect x="0.1" y="0.1" width="0.8" height="0.8" fill="#fff"/>`. As with `clipPath`, the simplest robust fix is a `chrome-testing/html/template/mask.html` template that supplies a properly-coordinated mask body for the `objectBoundingBox` variant.

#### 4. `role="region rowrowgroup"` — concatenated tokens, missing space, token collision (GRAMMAR_ISSUE)
Same root cause as `clipPath`: `RoleValue` grammar lacks a space separator. `"rowrowgroup"` is not a valid role — should be `"row rowgroup"`.

**Fix target:** `lang/svg.ebnf` line 82 — same fix as clipPath.

#### 5. `fill="none"` on mask container (GRAMMAR_ISSUE)
Card `fill="none"`:
```html
<mask id="slot" fill="none"><rect x="20" y="20" width="60" height="60" fill="#fff"/></mask>
```
The inner rect has an explicit `fill="#fff"` which overrides the inherited `none`, so the mask still works. The card is visually correct (bright square visible). Low severity.

#### 6. Duplicate `id` attribute (GRAMMAR_ISSUE)
Same pattern as all other elements.

---

## Cross-cutting issues summary

| Issue | Affected Elements | Fix Target |
|---|---|---|
| **Duplicate `id` attribute** (`id="slot" id="circle1"`) | All 7 elements | `blueprint.go` or `overlay.go` — id-shadowing guard when varying attr is `id` |
| **`RoleValue` missing space separator** (`"rowrowgroup"`, `"combobox complementarycontentinfo"`) | `clipPath`, `mask` (propagates to all elements with `AriaAttribute`) | `lang/svg.ebnf` line 82: `RoleValue = RoleToken , { " " , RoleToken } ;` |
| **`stop` scaffold has no baseline colour** (all 56 cards black) | `stop` | `blueprint.go` `baselineFor("stop", ...)` — add `stop-color="#f5a623" offset="0.5"` |
| **`objectBoundingBox` child coordinates wrong** (empty/invisible clip/mask region) | `clipPath`, `mask`, `pattern` | `blueprint.go` `bodyFor` + coordinate-space-aware child generation |
| **Gradient `href` self-reference** | `linearGradient`, `radialGradient` | `overlay.go` — emit a standalone `<linearGradient id="base">` def and point href there |
| **Mask/clipPath geometry at 100%** offsets masking out canvas | `mask` | `overlay.go` — clamp `LengthPercentageType` on geometry attrs to ≤ 80% |

---

## Per-element category summary

| Element | Category | Notes |
|---|---|---|
| `linearGradient` | MOSTLY_PERFECT | Duplicate `id`, self-ref `href`, degenerate `x2=0`/`y2=0` endpoints |
| `radialGradient` | MOSTLY_PERFECT | Duplicate `id`, self-ref `href` |
| `stop` | HAS_EMPTY_CARDS | **All 56 cards black** — no baseline `stop-color` on varied stop |
| `pattern` | MOSTLY_PERFECT | `objectBoundingBox` tiles invisible; duplicate `id` |
| `marker` | MOSTLY_PERFECT | `refY="_top"` typo candidate; duplicate `id` |
| `clipPath` | MOSTLY_PERFECT | `clipPathUnits=objectBoundingBox` clips to empty; `role` missing space; duplicate `id` |
| `mask` | MOSTLY_PERFECT | `x=100%`/`y=100%` hides mask; `maskContentUnits=objectBoundingBox` hides mask; `role` missing space; duplicate `id` |

---

## Top recurring issues

1. **`stop` baseline missing colour** — highest severity; makes the entire `stop` gallery useless. Fix in `blueprint.go` `baselineFor("stop")`.
2. **`RoleValue` missing space separator** — produces invalid role tokens on every element that enumerates `AriaAttribute`. Fix in `lang/svg.ebnf`.
3. **`objectBoundingBox` coordinate mismatch** — causes empty/invisible effect in `clipPath`, `mask`, and `pattern` when that specific attribute is varied. Fix in `blueprint.go` `bodyFor` with coordinate-space-aware child geometry.
4. **Duplicate `id` attribute** — propagates to every element in the gallery. Fix with an id-shadowing guard in `blueprint.go` or `overlay.go`.
5. **Self-referencing `href`** — semantic no-op for gradient inheritance. Fix in `overlay.go` by providing a separate reference target.
