# Round 2 QA — Animation Elements

**Batch:** animate, set, animateMotion, animateTransform, mpath, discard  
**Date:** 2026-06-22  
**Screenshots:** `chrome-testing/screenshots/review/{tag}.png` (re-shot fresh)

---

## Per-element verdicts

### `<animate>` — MOSTLY PERFECT with NO_EFFECT_IN_MAIN_GRID

**Round-1 fixes confirmed:**
- Clock values are valid (`dur="2s"`, `begin="0s"`, `end="4s"`, `min="0s"`, `max="indefinite"`, `repeatDur="4s"`). No malformed clocks.
- Host rect (40×40 blue square) is visible in all 56 main-grid cards.
- from/to typed to `attributeName="x"` — numeric values animate x position correctly.

**Remaining issues:**

| Card label | Category | Detail |
|---|---|---|
| `fill="none"` | NO_EFFECT_IN_MAIN_GRID | `fill` on `<animate>` is animation fill-mode (freeze/remove); "none" is not a valid animation-fill value. The rect is visible but this is a **grammar issue** — the `fill` attribute on animation elements should only enumerate `freeze` \| `remove`, not the full paint `<fill>` type. Card looks indistinct from the `fill="freeze"` / `fill="remove"` pair. |
| `fill-rule="nonzero"` | NO_EFFECT_IN_MAIN_GRID | Presentation attribute on `<animate>`, not rendered on the animate element itself; no visible difference from baseline. Should move to non-visual / meta section. |
| `fill-opacity="0.5"` | NO_EFFECT_IN_MAIN_GRID | Same — presentation attr on the animate element has no rendering effect on static capture. |
| `stroke="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `stroke-opacity="0.5"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `stroke-width="20"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `stroke-linecap="butt"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `stroke-linejoin="miter"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `stroke-miterlimit="4"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `stroke-dasharray="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `stroke-dashoffset="20"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `paint-order="normal"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `marker="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `marker-start="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `marker-mid="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `marker-end="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `color="#e94560"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `color-interpolation="auto"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `color-rendering="auto"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `shape-rendering="auto"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `vector-effect="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `filter="none"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `color-interpolation-filters="auto"` | NO_EFFECT_IN_MAIN_GRID | Same. |
| `flood-color="currentColor"` | NO_EFFECT_IN_MAIN_GRID | Same. |

**Count:** 24 presentation attributes in the main grid that have no rendering effect on `<animate>`. These are inherited/CSS presentation properties that apply to graphic content elements, not to the animation element itself.

**Root cause:** The grammar includes all global SVG presentation attributes on `<animate>`, `<set>`, `<animateMotion>`, `<animateTransform>`. These are technically valid per the SVG spec (they are "all elements" global attributes) but produce no visible output on animation elements. The gallery generator's non-visual classifier does not yet recognize this category.

---

### `<set>` — MOSTLY PERFECT with NO_EFFECT_IN_MAIN_GRID (same pattern)

**Round-1 fixes confirmed:**
- Clock values valid. Host shape visible on all main-grid cards (17 main cards).
- `attributeName="x"`, `from="10"`, `to="80"` correctly typed.

**Note:** `<set>` correctly does NOT emit `calcMode`, `values`, `keyTimes`, `keySplines`, `by`, `accumulate` (these are `<animate>`-specific). 

**Remaining issues (same root cause as `<animate>`):**

Main-grid presentation attrs with no effect:
- `fill="none"` — again, invalid animation-fill value; grammar issue.
- `fill-rule="nonzero"`, `fill-opacity="0.5"`, `stroke="none"`, `stroke-opacity="0.5"`, `stroke-width="20"`, `stroke-linecap="butt"`, `stroke-linejoin="miter"`, `stroke-miterlimit="4"`, `stroke-dasharray="none"`, `stroke-dashoffset="20"`, `paint-order="normal"`, `marker="none"`, `marker-start="none"`, `marker-mid="none"`, `marker-end="none"`, `color="#e94560"`, `color-interpolation="auto"`, `color-rendering="auto"`, `shape-rendering="auto"`, `vector-effect="none"`, `filter="none"`, `color-interpolation-filters="auto"`, `flood-color="currentColor"` — all NO_EFFECT_IN_MAIN_GRID.

**Count:** ~23 presentation attrs in main grid with no rendering effect on `<set>`.

---

### `<animateMotion>` — MOSTLY PERFECT with NO_EFFECT_IN_MAIN_GRID (same pattern)

**Round-1 fixes confirmed:**
- Host rect visible (40×40 blue square) in all cards — the animation moves it.
- Clock values valid.
- Path-specific cards show distinct positions: `path="M10 10 L90 90"`, `path="M10 50 Q50 10 90 50"`, `path="M20 20 H80 V80 H20 Z"` — three visually different positions in the captured frame.
- `rotate="auto"` card shows the rect at a 45° tilt (diamond orientation) — clearly distinct. Good.
- `rotate="auto-reverse"` shows an arrowhead-like rotated shape — distinct. Good.
- `from="10,10"`, `to="80,50"`, `by="40,40"` — coordinate pairs correctly formatted.

**Remaining issues:**

Same presentation-attribute-on-animation-element problem in main grid:
- `fill="none"` — grammar issue (same as animate/set).
- `fill-rule="nonzero"`, `fill-opacity="0.5"`, `stroke="none"`, `stroke-opacity="0.5"`, `stroke-width="20"`, `stroke-linecap="butt"`, `stroke-linejoin="miter"`, `stroke-miterlimit="4"`, `stroke-dasharray="none"`, `stroke-dashoffset="20"`, `paint-order="normal"`, `marker="none"`, `marker-start="none"`, `marker-mid="none"`, `marker-end="none"`, `color="#e94560"`, `color-interpolation="auto"`, `color-rendering="auto"`, `shape-rendering="auto"`, `vector-effect="none"`, `filter="none"`, `color-interpolation-filters="auto"`, `flood-color="currentColor"` — all NO_EFFECT_IN_MAIN_GRID.

**WEAK_EFFECT note:**
- `href="#target"` card: rect animates along an `href` to an un-resolved `#target` (that id is absent in the animateMotion scaffold — mpath-only variant). The rect is visible but the motion may be a no-op. **Minor** — host shape shows, so not EMPTY.

---

### `<animateTransform>` — MOSTLY PERFECT with NO_EFFECT_IN_MAIN_GRID

**Round-1 fixes confirmed:**
- Host rect visible, rotated in mid-spin in all cards (diamond orientation prominent).
- Clock values valid.
- `attributeName="transform"`, `type="rotate"`, `from="0 25 25"`, `to="360 25 25"` — correct baseline.
- `type="translate"`, `type="scale"`, `type="rotate"`, `type="skewX"`, `type="skewY"` — all show distinct geometry. Particularly:
  - `type="skewX"` → parallelogram shear. GOOD.
  - `type="skewY"` → vertical shear square. GOOD.
  - `type="scale"` → large square (scaled out). GOOD.
  - `type="translate"` → rect displaced to different position. GOOD.
- `from="0"` / `to="45"` solo scalar cards — shape is in a rotated pose. GOOD.
- `values="0 25 25; 360 25 25"` — distinct pose from baseline. GOOD.

**Remaining issues:**

- `type="translate"` card: baseline `from="0 25 25"` / `to="360 25 25"` are rotation-format triplets applied to translate type — translate from/to should be `"x y"` pairs, not `"angle cx cy"`. This is a **GRAMMAR_ISSUE** / type-mismatch. The rect still appears translated (browsers are lenient), but the values are semantically wrong.
- `type="scale"` card: same — `from="0 25 25"` is meaningless for a scale transform. The overlay `animFrom`/`animTo` do not disambiguate by `type`.
- `type="skewX"` / `type="skewY"` cards: same — three-component from/to applied to single-angle types. Browsers tolerate this (ignore extra components) but it is a grammar integrity issue.

Same presentation-attribute-on-animation-element NO_EFFECT_IN_MAIN_GRID pattern:
- `fill="none"`, `fill-rule="nonzero"`, `fill-opacity="0.5"`, `stroke="none"`, `stroke-opacity="0.5"`, `stroke-width="20"`, `stroke-linecap="butt"`, `stroke-linejoin="miter"`, `stroke-miterlimit="4"`, `stroke-dasharray="none"`, `stroke-dashoffset="20"`, `paint-order="normal"`, `marker="none"`, `marker-start="none"`, `marker-mid="none"`, `marker-end="none"`, `color="#e94560"`, `color-interpolation="auto"`, `color-rendering="auto"`, `shape-rendering="auto"`, `vector-effect="none"`, `filter="none"`, `color-interpolation-filters="auto"`, `flood-color="currentColor"` — all NO_EFFECT_IN_MAIN_GRID.

---

### `<mpath>` — PERFECT

**Round-1 fixes confirmed:**
- Host shape (20×20 orange/gold rect) renders in all 3 main-grid cards. The host is at (30,30) in the upper-left of the SVG canvas.
- `href="#target"` — resolves to the `<path id="target" d="M30 50 Q50 10 70 50"/>` in defs. Rect visible.
- `xlink:href="#target"` — same reference, visible.
- `xlink:title="label"` — metadata; rect visible (motion runs but no href to path, rect stays at origin — acceptable for a non-motion-path attr).
- 27 non-visual attributes are correctly collapsed into the `<details>` meta section.
- No blank cards. No grammar issues observed.

**Status: PERFECT** — only 3 visual attrs (href, xlink:href, xlink:title), all render correctly.

---

### `<discard>` — PERFECT

**Round-1 fixes confirmed:**
- Host shape visible in both main-grid cards.
- `href="#target"` card: `begin="60s"` baseline defers discard; host rect (40×40 blue) visible at t=0.
- `begin="0s"` card: The `<discard begin="0s">` fires immediately and removes the parent rect from the DOM. The card therefore renders as **empty canvas** at the snapshot time. This is **CORRECT BEHAVIOR** for `<discard>` — the element's purpose is removal. The card is not EMPTY in a broken sense; it demonstrates the element's effect.
- 25 non-visual attributes correctly in meta section.

**Status: PERFECT** — the `begin="0s"` empty canvas is correct semantics for `<discard>`.

---

## Systematic issues

### Issue 1 — NO_EFFECT_IN_MAIN_GRID: Presentation attributes on animation elements (HIGH priority)

**Affects:** animate, set, animateMotion, animateTransform (~24 cards each)

**Root cause:** The SVG grammar assigns all global presentation attributes (`fill`, `stroke`, `fill-opacity`, `stroke-width`, etc.) to animation elements. These are technically legal per spec but have no rendering effect on animation elements themselves (they do not inherit to the animated target). The gallery classifier puts them in the main grid instead of the non-visual meta section.

**Fix target:** `chrome-testing/cmd/gen/gallery.go` or `overlay.go`

**Recommended change:** Add a `isAnimationElement(tag)` guard in the non-visual attribute classifier. For animation tags, treat the entire set of CSS presentation properties (`fill`, `stroke`, `*-opacity`, `stroke-*`, `paint-order`, `marker*`, `color*`, `shape-rendering`, `vector-effect`, `filter`, `flood-color`) as non-visual when the element is an animation element. These ~24 attribute paths per animation element should move to the `<details class="meta">` section.

Specifically in `gallery.go` (or wherever `isNonVisual` / `shouldShowInMain` is decided), add a condition:

```go
// Presentation attributes on animation elements have no rendering effect on the
// animation element itself — they apply to the animated target, not as a visible
// property of the animate node. Classify them as non-visual for animation tags.
if isAnimationTag(tag) && isPresentationAttr(attrName) {
    return true // non-visual
}
```

Where `isPresentationAttr` covers: fill, stroke, fill-opacity, stroke-opacity, fill-rule, stroke-width, stroke-linecap, stroke-linejoin, stroke-miterlimit, stroke-dasharray, stroke-dashoffset, paint-order, marker, marker-start, marker-mid, marker-end, color, color-interpolation, color-rendering, color-interpolation-filters, shape-rendering, vector-effect, filter, flood-color.

---

### Issue 2 — GRAMMAR_ISSUE: `fill` on animation elements enumerates invalid value "none" (MEDIUM priority)

**Affects:** animate, set, animateMotion, animateTransform — the `fill="none"` card appears in the main grid.

**Root cause:** On animation elements, `fill` is the **animation-fill-mode** attribute with valid values `freeze | remove` only. The grammar's `fill` attribute for animation elements should be a separate production from the graphic `<fill>` paint type. The grammar currently reuses the same `fill` production (which includes `none`, color values, `url(#...)`) for animation elements.

**Fix target:** The animation element EBNF grammar file (whichever `.ebnf` defines `<animate>`, `<set>`, etc.) — restrict `fill` on animation elements to `"freeze" | "remove"`.

Also see overlay.go: add an override for `fill` when `isAnimationTag(tag)` to return one of `"freeze"` | `"remove"` rather than allowing the full paint grammar.

---

### Issue 3 — GRAMMAR_ISSUE: `animateTransform` from/to values not typed by `type` attribute (LOW priority)

**Affects:** animateTransform — `type="translate"`, `type="scale"`, `type="skewX"`, `type="skewY"` cards show rotation-format values (`"0 25 25"` / `"360 25 25"`) which are semantically wrong for those transform types.

**Root cause:** `overlay.go` `animFrom`/`animTo` always return the rotation-triplet format (`"0 25 25"` / `"360 25 25"`) regardless of the `type` attribute value. The `type` attribute is varied per card but the baseline always fixes rotation-format values.

**Fix target:** `chrome-testing/cmd/gen/overlay.go`

**Recommended change:** Extend `animFrom`/`animTo` (or add a `animFromForType(tag, transformType)` helper) that branches on the value of `type`:
- `rotate` → `"0 25 25"` / `"360 25 25"` (angle cx cy)
- `translate` → `"0 0"` / `"50 30"` (tx ty)
- `scale` → `"1"` / `"2"` (sx or sx sy)
- `skewX` / `skewY` → `"0"` / `"30"` (angle)

Since `type` is the varied attribute, the baseline `baselineFor("animateTransform")` already fixes `type="rotate"` for all other cards. The specific `type="X"` cards use the overlay's `from`/`to` values which don't know the current type. The cleanest fix is: in `overlay.go` `overlaySample`, detect that we are generating a `type` variant and cross-reference the chosen type — **or** accept browsers are lenient and treat this as cosmetic only.

---

## Summary table

| Element | Status | Main issues |
|---|---|---|
| `<animate>` | MOSTLY PERFECT | 24 NO_EFFECT_IN_MAIN_GRID (presentation attrs); `fill="none"` grammar issue |
| `<set>` | MOSTLY PERFECT | Same ~23 NO_EFFECT_IN_MAIN_GRID; `fill="none"` grammar issue |
| `<animateMotion>` | MOSTLY PERFECT | Same presentation attr NO_EFFECT pattern; `fill="none"` grammar issue |
| `<animateTransform>` | MOSTLY PERFECT | Same presentation attr NO_EFFECT; `fill="none"`; from/to type-mismatch for non-rotate types |
| `<mpath>` | PERFECT | None |
| `<discard>` | PERFECT | None (`begin="0s"` blank is correct behavior) |

**Top priority fix:** Move presentation attributes on animation elements to the meta/non-visual section in `gallery.go` (Issue 1). This will clean up ~24 identical-looking cards per animation element from the main grid, reducing noise significantly.
