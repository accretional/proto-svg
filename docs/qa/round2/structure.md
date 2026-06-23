# Round 2 QA — Structure / Container Elements

**Batch:** svg, g, defs, symbol, use, switch, a, desc, title, metadata, unknown  
**Date:** 2026-06-22  
**Screenshots:** `chrome-testing/screenshots/review/{tag}.png` (re-shot fresh)

---

## Per-element verdicts

### `<svg>` — MOSTLY PERFECT with NO_EFFECT_IN_MAIN_GRID / WEAK_EFFECT

**Round-1 fixes confirmed:**
- viewBox, preserveAspectRatio, zoomAndPan, x, y, width, height all render correct geometry effects.
- transform cards (translate, rotate, scale, skewX) show visually distinct transformed rects.
- Presentation attributes (fill, stroke, fill-opacity, etc.) show inherited effects on the child rect.

**Remaining issues:**

| Card label | Category | Detail |
|---|---|---|
| `requiredExtensions=""` | NO_EFFECT_IN_MAIN_GRID | Chrome hides the `<svg>` when `requiredExtensions=""` (empty string interpreted as failing test — the rect disappears). The overlay sets `""` intending "no requirements" but Chrome disagrees. |
| `systemLanguage="en"` | PERFECT | Renders (browser matches "en"). |
| `width="auto"` | WEAK_EFFECT | Square collapses to ~0 width; card nearly blank. Auto on inline svg root needs a CSS context. |
| `height="auto"` | WEAK_EFFECT | Same collapse. |

**Assessment:** MOSTLY PERFECT. The `requiredExtensions` empty-value Chrome behavior is a browser quirk, not a generator bug. The `auto` dimension collapse is acceptable (correct SVG semantics).

---

### `<g>` — ISSUES: NO_EFFECT_IN_MAIN_GRID, missing `transform` / `opacity`

**Round-1 fixes confirmed:**
- Child rect (60×60 blue) is present in all main-grid cards except `requiredExtensions` and `systemLanguage`.

**Remaining issues:**

| Card label | Category | Detail |
|---|---|---|
| `requiredExtensions=""` | EMPTY/BLANK | Chrome hides `<g requiredExtensions="">` (same browser bug as svg). Card shows dark background only. |
| `systemLanguage="en"` | EMPTY/BLANK | Unexpectedly empty — Chrome appears to hide the `<g>` for `systemLanguage="en"` in this context. Needs investigation (may be a language mismatch in the snap headless browser). |
| `transform` | MISSING FROM GRID | `transform` is the most important `<g>`-specific attr but is absent — cut off by `PresentationAttribute` cap=24. `TransformAttr` is member #56 of the union, well past the 24-member cap. |
| `opacity` | MISSING FROM GRID | Same cap issue. `OpacityAttr` is member #53 of `PresentationAttribute` — past cap=24. |
| `display`, `visibility`, `clip-path`, `mask`, all font/text attrs, `pointer-events`, `stop-color` | MISSING FROM GRID | All fall after position 24 in `PresentationAttribute` order. Not shown for `<g>`, `<defs>`, `<switch>`, `<a>`, `<symbol>`, `<unknown>`. |

**Fix target:** `chrome-testing/cmd/gen/enumerate.go` — raise `sharedGroupMemberCap["PresentationAttribute"]` from 24 to at least 40, OR reorder `PresentationAttribute` in `lang/styling.ebnf` to put `OpacityAttr` and `TransformAttr` near the front (positions 1-2).

**Recommended fix (styling.ebnf reorder):** Move `OpacityAttr | TransformAttr | TransformOriginAttr` to be the FIRST three entries of `PresentationAttribute`, ahead of painting attrs. This ensures they always appear in the cap window on every container element.

---

### `<defs>` — CRITICAL: ALL MAIN-GRID CARDS BLANK

**Round-1 fixes attempted:** `blueprintSlotNeedsID` returns `true` for `defs`, so `id="slot"` is injected. The generated markup is:
```
<defs><defs id="slot" fill="none"><rect x="20" y="20" width="60" height="60" fill="#4d8bff"/></defs></defs>
<use href="#slot" x="10" y="10"/>
```
**Root cause:** `<use href="#slot">` references the inner `<defs id="slot">` element. But `<defs>` is a non-rendering container — browsers never paint it directly. So every `<use>` instantiation of a `<defs>` produces no visible output. All 26 main-grid cards are empty dark squares.

The blueprint (from `defaultScaffold`) for `defs` falls into `catSelf` (default), yielding `<svg>{{ELEMENT}}</svg>`. The `<defs>` body gets a rect child via `bodyFor("defs")`. But the full scaffold for defs should instead place the `<defs>` in a `<defs>` wrapper and reference a SHAPE inside it — not the `<defs>` itself.

**Also missing:** Same `transform`/`opacity` cap issue as `<g>`.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go`  
Change the `defs` category from `catSelf` (implicit) to a new explicit scaffold. In `defaultScaffold`, add a case for `"defs"` that wraps the defs child rect in a `<g>` instead of using `<use>`. Since `<defs>` attrs (fill, stroke, etc.) only matter as presentation inheritance for any children it might have in the real doc, the most honest demo is a simple `<svg><defs ...><rect id="slot" .../></defs><use href="#slot"/></svg>` where the `<rect>` inside `<defs>` carries the slot id (not the `<defs>` itself).

**Alternatively:** Change `blueprintSlotNeedsID` to return `false` for `defs` and instead give the `<defs>` a custom blueprint (template html) that embeds the element's child rect directly in the SVG root to show presentation attribute inheritance.

---

### `<symbol>` — PERFECT

**Round-1 fixes confirmed:** All cards show a blue rect. The `<symbol id="slot">` is properly defined and `<use href="#slot">` renders it.  
viewBox, preserveAspectRatio, refX, refY, x, y, width, height all show correct geometry/clipping effects.  
Presentation attrs (fill, stroke, opacity group, markers, colors) all show inherited effects.

**Remaining issues:**
- Same `transform`/`opacity` cap issue applies but is less severe since `symbol` already shows 88 variants covering the most useful attrs.
- `transform` is missing from the symbol gallery (same cap=24 on PresentationAttribute).

**Assessment:** PERFECT for the rendered variants. Cap issue is a systematic problem addressed under `<g>`.

---

### `<use>` — MOSTLY PERFECT

**Round-1 fixes confirmed:** The blueprint defines `<rect id="slot">` in defs. Every card shows the referenced pink/red dot (a circle is used as body — wait, the generated shape is actually a red rect from the blueprint `catContainerRef` scaffold). All x, y, width, height, href variants are visible.

**Remaining issues:**

| Card label | Category | Detail |
|---|---|---|
| `fill="none"` | WEAK_EFFECT | The `<use>` fill attr is overridden by the referenced shape's own fill — inherited fill on `<use>` is overridden by the `<rect fill="#e94560">` inside defs. Card shows same appearance as neighbors. |
| `fill-rule="nonzero"` | NO_EFFECT_IN_MAIN_GRID | fill-rule on `<use>` has no effect on this rectangular geometry. |
| `stroke-linecap`, `stroke-linejoin`, `stroke-miterlimit` | NO_EFFECT_IN_MAIN_GRID | No stroke applied to the referenced rect, so these attrs have no visible effect. |
| `marker`, `marker-start`, `marker-mid`, `marker-end` | NO_EFFECT_IN_MAIN_GRID | Markers apply to path/line endpoints; the rect has none. |
| `paint-order` | NO_EFFECT_IN_MAIN_GRID | No stroke applied; paint-order is invisible. |
| `xlink:href="#slot"` | PERFECT | Deprecated alias still resolves and renders. |

**Assessment:** MOSTLY PERFECT. The presentation attr NO_EFFECT cases on `<use>` are correct SVG behavior (inheritance overridden by child's own fill; no-stroke geometry). Not generator bugs.

---

### `<switch>` — ISSUES: EMPTY FIRST TWO CARDS + missing transform/opacity

**Round-1 fixes confirmed:** Most cards show a blue rect child.

**Remaining issues:**

| Card label | Category | Detail |
|---|---|---|
| `requiredExtensions=""` | EMPTY/BLANK | Same Chrome behavior as `<g>` — element hidden. |
| `systemLanguage="en"` | EMPTY/BLANK | Same issue — element hidden in Chrome headless. |
| `transform` / `opacity` | MISSING FROM GRID | Cap=24 on PresentationAttribute. |

**Assessment:** Same systematic issues as `<g>`. The `<switch>` structural behavior (showing first matching child) is not explicitly testable with this gallery approach but the presentation attr inheritance cards render correctly.

---

### `<a>` — MOSTLY PERFECT

**Round-1 fixes confirmed:** All cards show a blue rect. href, target, download, ping, rel, hreflang, type, referrerpolicy variants all render (they have no visual effect but show correctly).

**Remaining issues:**

| Card label | Category | Detail |
|---|---|---|
| `href="#slot"` | PERFECT | Renders rect. href is an IRI reference in SVG, works as expected. |
| `target="_self"`, `_parent"`, etc. | NO_EFFECT_IN_MAIN_GRID | Link target is a navigation hint — no visual difference. Acceptable. |
| `download`, `ping`, `rel`, `hreflang`, `type`, `referrerpolicy` | NO_EFFECT_IN_MAIN_GRID | All are navigation-semantics attrs with no visual rendering effect. Consider moving to the non-visual "meta" section. |
| `transform` / `opacity` | MISSING FROM GRID | Cap=24. |

**Fix target (optional):** `chrome-testing/cmd/gen/emit.go` — extend `nonVisualAttr` to include `download`, `ping`, `rel`, `hreflang`, `type`, `referrerpolicy` for `<a>`. Currently they appear in the main grid with identical-looking blue rects, which is confusing. They belong in the collapsed "Non-visual attributes" section.

---

### `<desc>` — CORRECT (pure non-visual element)

All 28 attrs are Core/GlobalEvent/DocumentElementEvent — correctly classified as non-visual and placed in the collapsed section. The main grid is empty, which is correct: `<desc>` has NO presentation attributes and no visual rendering.

**Assessment:** PERFECT. The empty main grid is correct behavior for a metadata-only element. No fix needed.

---

### `<title>` — CORRECT (pure non-visual element)

Identical structure to `<desc>`. All 28 attrs non-visual, main grid empty.

**Assessment:** PERFECT. Correct for a tooltip/accessibility-only element.

---

### `<metadata>` — CORRECT (pure non-visual element)

Identical structure to `<desc>` / `<title>`. All 28 attrs non-visual, main grid empty.

**Assessment:** PERFECT. Metadata is intentionally non-visual.

---

### `<unknown>` — CRITICAL: ALL MAIN-GRID CARDS BLANK

The `<unknown>` element is not a recognized SVG tag. Every card shows an empty dark square because:
1. `bodyFor("unknown")` returns `""` — no child shape is provided.
2. Even if children were provided, `<unknown>` is treated as an unknown element by browsers and its presentation attributes have no rendering effect on it (they cannot be inherited because there's nothing to inherit to).

Current markup for each card: `<svg><unknown fill="none"></unknown></svg>` — nothing to show.

**Two options:**
- **Option A (correct semantics):** Accept that `<unknown>` shows all-blank cards in the main grid. Move ALL its attrs to the non-visual section (same as desc/title/metadata). Add a note in the page description: "Unknown element — no visual rendering; attributes shown for grammar completeness."
- **Option B (illustrative):** Add a child rect inside `<unknown>` via `bodyFor("unknown")` to show that SVG browsers parse unknown elements but don't render their presentation attrs — producing blank-rect cards that at least prove the markup is parsed.

**Assessment:** Option A is the most honest. The main grid should be empty with a note. The `unknown` case is semantically correct as-is but needs a UI explanation so viewers don't think it's a generator bug.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go` — add `"unknown"` to `bodyFor` to return a rect, AND/OR change `nonVisualAttr` to flag all attrs on `<unknown>` as non-visual via an element-aware predicate.

---

## Systematic Issues Summary

### ISSUE 1 (CRITICAL — affects g, defs, switch, a, symbol, unknown): `PresentationAttribute` cap=24 cuts off `transform`, `opacity`, `display`, `visibility`, `clip-path`, `mask`, all font/text, `pointer-events`

**Root cause:** `sharedGroupMemberCap["PresentationAttribute"] = 24` in `enumerate.go`. The union has 52 members; painting attrs fill positions 1-24, so `OpacityAttr` (pos ~53) and `TransformAttr` (pos ~56) are never reached for `<g>`, `<defs>`, `<switch>`, `<a>`, `<symbol>`, `<unknown>`.

**`<svg>` is exempt** because it has its OWN `SvgTransformAttr` defined explicitly in `SvgAttribute` (not through `PresentationAttribute`), so it gets transform regardless of the cap.

**Fix (preferred):** Reorder `PresentationAttribute` in `lang/styling.ebnf` to put the most container-relevant attrs first:
```ebnf
PresentationAttribute =
  (* move these to TOP so they're always within any cap *)
    OpacityAttr | TransformAttr | TransformOriginAttr
  | DisplayAttr | VisibilityAttr | OverflowAttr
  (* painting — unchanged order after the above *)
  | FillAttr | FillRuleAttr | FillOpacityAttr
  | ...
```

**Fix (alternative):** Raise cap to 40 in `enumerate.go`:
```go
"PresentationAttribute": 40,
```

---

### ISSUE 2 (CRITICAL — defs): `<defs>` scaffold uses `<use href="#defs-id">` which renders nothing

`<use>` of a `<defs>` element is invisible. The `defs` blueprint needs a complete redesign. Suggested fix in `blueprint.go` `defaultScaffold`:

```go
case "defs":
    // Show the child rect directly; defs presentation attrs inherit to children
    // only, so we render with the child inline in the SVG.
    return svgOpen + `{{ELEMENT}}</svg>`
```
And keep `bodyFor("defs")` returning a rect. This way `<defs fill="none">` cards will show a rect (the child rect does NOT inherit fill="none" from defs — that's actually the correct visual demonstration that defs has no direct visual output but its presentation attrs ARE inherited by children if they lack their own).

Actually, for the clearest visual: change the `defs` blueprint to render the `<defs>` element with the child rect inside the SVG root (not via `<use>`), so every card shows the rect, and attrs like `opacity="0.5"` on the `<defs>` are visibly inherited.

---

### ISSUE 3 (MINOR — g, switch): `requiredExtensions=""` and `systemLanguage` hide elements in Chrome

Chrome's SVG engine hides `<g requiredExtensions="">` (treating empty string as a required-extension that isn't available) and may hide elements with `systemLanguage` mismatch. These cards show empty dark squares.

**Fix option 1:** Move `requiredExtensions` and `systemLanguage` cards to the non-visual section (extend `nonVisualAttr` to cover them) — they are conditional-processing attributes with no predictable visual effect.

**Fix option 2:** Change the overlay to produce `systemLanguage="en"` for the attribute card value (it already does this), and reconsider `requiredExtensions` value. Note: `requiredExtensions=""` is already the overlay override — Chrome simply doesn't honor the "empty = always true" semantic.

**Simplest fix:** `chrome-testing/cmd/gen/emit.go` — add `"requiredextensions"` and `"systemlanguage"` to `nonVisualAttr()` so they land in the collapsed section. They are conditional-processing flags, not painting attrs.

---

### ISSUE 4 (MINOR — a): Navigation attrs in main grid

`download`, `ping`, `rel`, `hreflang`, `type`, `referrerpolicy` for `<a>` are navigation semantics with no visual effect. They currently appear in the main grid with identical blue rect cards.

**Fix:** `chrome-testing/cmd/gen/emit.go` — extend `nonVisualAttr` to return true for these attr names (possibly with an element-aware check, or universally since they never affect SVG rendering).

---

### ISSUE 5 (MINOR — unknown): No page-level note explaining expected blankness

`<unknown>` correctly produces blank cards for presentation attrs (browser parses but doesn't render an unknown element), but without a page note, this looks like a generator failure.

**Fix:** Add element-specific `<p>` notes in the page template or in `emitPage` for elements where blank is correct. Or: classify all `<unknown>` attrs as non-visual via element-aware logic.

---

## Fix Priority

| Priority | Issue | File | Change |
|---|---|---|---|
| P0 | Defs all-blank (scaffold broken) | `blueprint.go` | Remove `defs` from `blueprintSlotNeedsID`; change `defaultScaffold("defs")` to drop the `<use>` approach |
| P0 | transform/opacity missing from g/switch/a/defs/symbol | `styling.ebnf` | Reorder `PresentationAttribute` to put `OpacityAttr`, `TransformAttr` first |
| P1 | requiredExtensions/systemLanguage empty cards | `emit.go` | Add to `nonVisualAttr()` |
| P1 | `<a>` nav attrs in main grid | `emit.go` | Add download/ping/rel/hreflang/type/referrerpolicy to `nonVisualAttr()` |
| P2 | `<unknown>` no explanatory note | `emit.go` or `blueprint.go` | Add page note or classify all attrs as non-visual |
