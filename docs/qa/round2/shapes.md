# Visual QA — Round 2: Shape Elements

Reviewer: automated visual QA pass  
Date: 2026-06-22  
Batch: `rect`, `circle`, `ellipse`, `line`, `polyline`, `polygon`, `path`  
Method: fresh full-page screenshots (`SNAP_FULLPAGE=1 SNAP_SCALE=1`) + HTML source inspection per card.  
Screenshots: `chrome-testing/screenshots/review/<tag>.png`

---

## What was fixed since Round 1

| Round 1 issue | Fixed? |
|---------------|--------|
| `x="100%"` / `y="100%"` / `cx="100%"` / `cy="100%"` off-canvas | YES — now `75%`, shape stays on canvas |
| `fill="none"` blank on rect / circle / ellipse / polygon | YES — baselines now have `stroke="#16c79a" stroke-width="2"` |
| `requiredExtensions="label"` hidden | YES — now `""` (empty list = always satisfied) |
| `systemLanguage="label"` hidden | YES — now `"en"` (matches browser language) |
| `d="none"` blank on path | YES — now `"M10 50 Q50 10 90 50"` |
| `stroke-miterlimit="0"` invalid | YES — now `"4"` (SVG default, spec minimum 1) |
| id / tabindex / lang / xml:lang / xml:space / role / aria-* / on* events in main grid | YES — collapsed into `<details class="meta">` Non-visual attributes section |

---

## Summary table

| Element | Main-grid cards | Category | Remaining issues |
|---------|-----------------|----------|-----------------|
| `rect` | 43 | ISSUES_REMAIN | NO_EFFECT_IN_MAIN_GRID: 15 attrs; WEAK_EFFECT: 3 |
| `circle` | 38 | ISSUES_REMAIN | NO_EFFECT_IN_MAIN_GRID: 14 attrs; WEAK_EFFECT: 2; EMPTY: 1 |
| `ellipse` | 41 | ISSUES_REMAIN | NO_EFFECT_IN_MAIN_GRID: 15 attrs; WEAK_EFFECT: 2; EMPTY: 1 |
| `line` | 43 | ISSUES_REMAIN | NO_EFFECT_IN_MAIN_GRID: 8 attrs; WEAK_EFFECT: 3; EMPTY: 1 |
| `polyline` | 29 | NEARLY_PERFECT | NO_EFFECT_IN_MAIN_GRID: 11 attrs; WEAK_EFFECT: 2 |
| `polygon` | 29 | NEARLY_PERFECT | NO_EFFECT_IN_MAIN_GRID: 11 attrs; WEAK_EFFECT: 1; EMPTY: 1 |
| `path` | 28 | ISSUES_REMAIN | NO_EFFECT_IN_MAIN_GRID: 10 attrs; WEAK_EFFECT: 4; EMPTY: 1 |

---

## `<rect>` — ISSUES_REMAIN

43 main-grid cards, 32 in the non-visual details section.

### EMPTY / BLANK (0 cards)

None — all main-grid cards now render a visible shape.

### NO_EFFECT_IN_MAIN_GRID (15 cards)

These cards render a rect that is **visually indistinguishable from the baseline**:

| Card label | Why no effect |
|------------|--------------|
| `requiredExtensions=""` | Empty list = always satisfied; renders identical baseline rect |
| `systemLanguage="en"` | Matches browser; renders identical baseline rect |
| `pathLength="2"` | Only affects dash-pattern scaling; baseline has solid stroke, no visible change |
| `fill-rule="nonzero"` | `rect` is a convex shape; `nonzero` and `evenodd` are identical on convex fills |
| `stroke-linecap="butt"` | `rect` is a closed shape with no open endpoints; linecap has no effect |
| `stroke-linejoin="miter"` | `"miter"` is the SVG default; card looks identical to baseline |
| `stroke-miterlimit="4"` | `4` is the SVG default; no visible change |
| `paint-order="normal"` | `"normal"` is the default; no visible change |
| `marker="none"` | No marker drawn; identical to baseline (markers apply at path endpoints) |
| `marker-start="none"` | Same as above |
| `marker-mid="none"` | Same as above |
| `marker-end="none"` | Same as above |
| `color-interpolation="auto"` | Rendering hint; imperceptible difference on static rasterized fill |
| `color-rendering="auto"` | Rendering hint; imperceptible |
| `shape-rendering="auto"` | Rendering hint; imperceptible |

### WEAK_EFFECT (3 cards)

These have a technically valid effect but it is either invisible or imperceptible at the card scale:

| Card label | Issue |
|------------|-------|
| `color="#e94560"` | Sets `currentColor`, but baseline `fill` and `stroke` use explicit hex; no visible change |
| `vector-effect="none"` | `"none"` is the default; no visible change |
| `filter="none"` | No filter applied; identical to baseline |
| `color-interpolation-filters="auto"` | Only affects filter primitives; baseline has no filter; no effect |
| `flood-color="currentColor"` | `flood-color` only affects `feFlood` filter primitive; no filter on rect; no effect |

*(The last three could also be classified NO_EFFECT; they are listed here because they are valid presentation attributes that have an effect in other contexts.)*

### Fix targets

| Issue | File | Change |
|-------|------|--------|
| `requiredExtensions` / `systemLanguage` identical to baseline | `emit.go` `nonVisualAttr()` | Add `"requiredextensions"` and `"systemlanguage"` to `nonVisualAttr` — they are conditional-processing attributes with no rendering delta on the element itself; move to details section |
| `pathLength` no effect without dasharray | `overlay.go` `overlaySample` | Either add `pathLength` to a "conditional-demo" note, OR keep it in main grid but add a `stroke-dasharray="10 5"` to the baseline so `pathLength` visibly rescales the dash; alternatively move `pathLength` to non-visual |
| `fill-rule` / `stroke-linecap` / `stroke-linejoin` / `stroke-miterlimit` on rect | `emit.go` `nonVisualAttr(tag, attr)` | These are element-dependent: linecap has no effect on closed shapes, fill-rule has no effect on convex shapes. Needs a per-tag version of `nonVisualAttr` |
| `paint-order="normal"` / `marker*="none"` identical defaults | `emit.go` `nonVisualAttr()` | Move to non-visual section; OR in `overlay.go` steer to a non-default value (e.g. `paint-order="stroke fill"` would be visibly different with thick stroke) |
| `color` / `vector-effect` / `filter` / `color-interpolation-filters` / `flood-color` | `emit.go` `nonVisualAttr()` | Add these to the non-visual set for shape elements: `"color-interpolation-filters"`, `"flood-color"`, `"color-rendering"`, `"shape-rendering"`, `"color-interpolation"` when no filter is applied; `"vector-effect"` when `"none"` |

---

## `<circle>` — ISSUES_REMAIN

38 main-grid cards, 32 in the non-visual details section.

### EMPTY / BLANK (1 card)

| Card label | Why empty |
|------------|-----------|
| *(none confirmed blank)* | All main-grid circle cards render a visible circle |

*(Note: `fill="none"` now shows the circle outline due to the stroke added to the baseline — correctly fixed.)*

### NO_EFFECT_IN_MAIN_GRID (14 cards)

Same set as `rect` minus `rx`/`ry`-specific entries:

`requiredExtensions=""`, `systemLanguage="en"`, `pathLength="2"`, `fill-rule="nonzero"` (convex circle), `stroke-linecap="butt"` (closed shape), `stroke-linejoin="miter"` (default), `stroke-miterlimit="4"` (default), `paint-order="normal"` (default), `marker="none"`, `marker-start="none"`, `marker-mid="none"`, `marker-end="none"`, `color-interpolation="auto"`, `color-rendering="auto"`, `shape-rendering="auto"`

### WEAK_EFFECT (2+ cards)

`color="#e94560"` (explicit hex fill/stroke, currentColor unused), `vector-effect="none"` (default), `filter="none"` (no filter), `color-interpolation-filters="auto"` (no filter), `flood-color="currentColor"` (not a filter element)

### Fix targets

Same as `rect`.

---

## `<ellipse>` — ISSUES_REMAIN

41 main-grid cards, 32 in the non-visual details section.

### EMPTY / BLANK (1 card)

| Card label | Why empty |
|------------|-----------|
| `stroke="none"` | Baseline ellipse has fill (`#f5a623`) so `stroke="none"` still shows the filled ellipse — NOT actually blank |

*(Re-checking screenshot: `stroke="none"` on ellipse shows the orange fill correctly — not blank. ZERO empty cards confirmed.)*

### NO_EFFECT_IN_MAIN_GRID (15 cards)

Same set as `circle`: `requiredExtensions=""`, `systemLanguage="en"`, `pathLength="2"`, `fill-rule="nonzero"` (convex), `stroke-linecap="butt"` (closed shape), `stroke-linejoin="miter"` (default), `stroke-miterlimit="4"` (default), `paint-order="normal"`, `marker="none"`, `marker-start="none"`, `marker-mid="none"`, `marker-end="none"`, `color-interpolation="auto"`, `color-rendering="auto"`, `shape-rendering="auto"`.

Also: `fill-opacity="0.5"` — baseline has an opaque orange fill; 0.5 IS visually distinct (semi-transparent ellipse over dark background shows differently) — actually **VISIBLE EFFECT**, keep in main grid. Confirmed in screenshot: row 3 col 7 shows a clearly darker/translucent ellipse.

### WEAK_EFFECT (2 cards)

`color="#e94560"`, `vector-effect="none"`, `filter="none"`, `color-interpolation-filters="auto"`, `flood-color="currentColor"` — same reasoning as `rect`.

### Fix targets

Same as `rect`.

---

## `<line>` — ISSUES_REMAIN

43 main-grid cards, 32 in the non-visual details section.

### EMPTY / BLANK (1 card)

| Card label | Why empty |
|------------|-----------|
| `stroke="none"` | A `<line>` has no fill area; removing stroke makes the line completely invisible |

**Fix:** In `overlay.go`, add a guard so `stroke="none"` is not emitted for `line` (and other stroke-only elements). Or accept this as a valid/intentional demo of "no visible stroke". Best fix: annotate the card differently, or in `blueprint.go` add a faint background reference shape so the card is not fully blank.

### NO_EFFECT_IN_MAIN_GRID (8 cards)

| Card label | Why no effect on `line` |
|------------|------------------------|
| `fill="none"` | `line` has no fill area; `fill="none"` is the effective default |
| `fill-rule="nonzero"` | `line` has no fill area; fill-rule has no effect |
| `fill-opacity="0.5"` | `line` has no fill; opacity of zero-fill is still zero |
| `stroke-linejoin="miter"` | A single-segment line has no joins; linejoin has no effect |
| `requiredExtensions=""` | Renders identical baseline line |
| `systemLanguage="en"` | Renders identical baseline line |
| `pathLength="2"` | No dasharray in baseline; no visible rescaling effect |
| `paint-order="normal"` | Default value; no change |

### WEAK_EFFECT (3 cards)

| Card label | Issue |
|------------|-------|
| `stroke-linecap="butt"` | `"butt"` is the default; shows line identical to baseline. Should use `"round"` or `"square"` to demonstrate visible cap extension |
| `stroke-linejoin="miter"` | (listed above as NO_EFFECT for single-segment line) |
| `color="#e94560"` | No `currentColor` reference in fill/stroke |

### Fix targets

| Issue | File | Change |
|-------|------|--------|
| `stroke="none"` blank | `overlay.go` | Add `case tag == "line" && an == "stroke"` → skip (or add a background SVG element to the scaffold in `blueprint.go`) |
| `fill` / `fill-rule` / `fill-opacity` on stroke-only element | `emit.go` per-tag `nonVisualAttr` | These three have no effect on `line`; move to non-visual section for `line` |
| `stroke-linecap="butt"` default value | grammar `painting.ebnf` `StrokeLinecapAttr` enumeration order | Move `"round"` before `"butt"` so the enumerator picks a visually distinct value; or add `case an == "stroke-linecap": return "round", true` to `overlay.go` |
| `stroke-linejoin` on single-segment line | `emit.go` per-tag `nonVisualAttr` | Move to non-visual for `line` |

---

## `<polyline>` — NEARLY_PERFECT

29 main-grid cards, 32 in the non-visual details section.

### EMPTY / BLANK (0 cards)

None.

### NO_EFFECT_IN_MAIN_GRID (11 cards)

| Card label | Why no effect |
|------------|--------------|
| `requiredExtensions=""` | Renders identical baseline polyline |
| `systemLanguage="en"` | Renders identical baseline polyline |
| `pathLength="2"` | No dasharray in baseline; no visible effect |
| `fill-rule="nonzero"` | Baseline is `fill="none"`; fill-rule has no effect |
| `paint-order="normal"` | Default value |
| `marker="none"` | No markers drawn; identical to baseline |
| `marker-start="none"` | Same |
| `marker-mid="none"` | Same |
| `marker-end="none"` | Same |
| `color-interpolation="auto"` | Rendering hint; imperceptible |
| `color-rendering="auto"` | Rendering hint; imperceptible |
| `shape-rendering="auto"` | Rendering hint; imperceptible |

*(12 if including `shape-rendering`.)*

### WEAK_EFFECT (2 cards)

| Card label | Issue |
|------------|-------|
| `stroke-linecap="butt"` | `"butt"` is the default; polyline HAS open endpoints so linecap IS meaningful, but the default value shows no change vs baseline. Use `"round"` for a visible demo |
| `stroke-linejoin="miter"` | `"miter"` is the default; polyline has joins (interior vertices) so linejoin IS meaningful, but default shows no change. Use `"round"` or `"bevel"` |

### Fix targets

| Issue | File | Change |
|-------|------|--------|
| `stroke-linecap` default value | `overlay.go` | Add `case an == "stroke-linecap": return "round", true` — this shows a visually distinct round cap at the open ends of the polyline |
| `stroke-linejoin` default value | `overlay.go` | Add `case an == "stroke-linejoin": return "round", true` — shows rounded interior vertex joins |
| `fill-rule` / `pathLength` / `paint-order` / `marker*` no-effect | `emit.go` per-tag `nonVisualAttr` | Move to non-visual section for `polyline` when the element's baseline makes them invisible |

---

## `<polygon>` — NEARLY_PERFECT

29 main-grid cards, 32 in the non-visual details section.

### EMPTY / BLANK (1 card)

| Card label | Why empty |
|------------|-----------|
| *(none)* | `fill="none"` now shows triangle outline due to baseline stroke — correctly fixed |

Actually re-checking polygon screenshot: all visible cards. **Zero blank cards confirmed.**

### NO_EFFECT_IN_MAIN_GRID (11 cards)

Same set as `polyline`: `requiredExtensions=""`, `systemLanguage="en"`, `pathLength="2"`, `fill-rule="nonzero"` (convex triangle — nonzero and evenodd give identical result), `paint-order="normal"`, `marker="none"`, `marker-start="none"`, `marker-mid="none"`, `marker-end="none"`, `color-interpolation="auto"`, `color-rendering="auto"`, `shape-rendering="auto"`.

### WEAK_EFFECT (1 card)

| Card label | Issue |
|------------|-------|
| `stroke-linejoin="miter"` | `"miter"` is the default; triangle corners ARE join points so linejoin IS meaningful, but default shows no change. Use `"round"` or `"bevel"` for a visibly distinct corner shape |

### Fix targets

| Issue | File | Change |
|-------|------|--------|
| `stroke-linejoin` default | `overlay.go` | Add `case an == "stroke-linejoin": return "round", true` |
| `fill-rule="nonzero"` on convex polygon | `emit.go` per-tag or `overlay.go` | Either move to non-visual for convex shapes, or change the polygon `points` baseline to a self-intersecting star shape so `nonzero` vs `evenodd` is visually distinct (more impactful fix) |

---

## `<path>` — ISSUES_REMAIN

28 main-grid cards, 32 in the non-visual details section.

### EMPTY / BLANK (1 card)

| Card label | Why empty |
|------------|-----------|
| `stroke="none"` | Path baseline is `fill="none"` (stroke-only arc); `stroke="none"` removes the only visible paint → completely invisible |

**Fix:** Same approach as `line` — either annotate the card, add a background reference shape, or in `overlay.go` redirect `stroke="none"` for stroke-only elements to a non-"none" value.

### NO_EFFECT_IN_MAIN_GRID (10 cards)

| Card label | Why no effect |
|------------|--------------|
| `requiredExtensions=""` | Renders identical baseline arc |
| `systemLanguage="en"` | Renders identical baseline arc |
| `pathLength="2"` | No dasharray in baseline; no visible effect |
| `fill="none"` | Baseline is already `fill="none"`; this card is identical to baseline |
| `fill-rule="nonzero"` | Baseline is `fill="none"`; fill-rule has no effect on unfilled path |
| `fill-opacity="0.5"` | Baseline is `fill="none"`; 0.5 × 0 fill = same as none |
| `paint-order="normal"` | Default value |
| `marker="none"` | No markers drawn; identical to baseline |
| `marker-start="none"` | Same |
| `marker-mid="none"` | Same |
| `marker-end="none"` | Same |

### WEAK_EFFECT (4 cards)

| Card label | Issue |
|------------|-------|
| `stroke-linecap="butt"` | `"butt"` is the default; arc has open endpoints so linecap IS meaningful, but default shows no delta from baseline. Use `"round"` or `"square"` |
| `stroke-linejoin="miter"` | `"miter"` is the default; the S-curve (T command) has an inflection but no sharp joins. Use `"round"` to show effect at the join |
| `stroke-miterlimit="4"` | `4` is the SVG default; no visible change |
| `color="#e94560"` | No `currentColor` reference in stroke; no visible change |

### Fix targets

| Issue | File | Change |
|-------|------|--------|
| `stroke="none"` blank | `overlay.go` | Add steering for stroke-only element+attr combination, or add a fallback fill in path blueprint |
| `fill="none"` identical to baseline | `blueprint.go` `baselineFor("path")` | The path baseline already uses `fill="none"`; either add a non-none fill to the baseline so the `fill="none"` card is meaningfully distinct, or suppress this card for path |
| `fill-rule` / `fill-opacity` on fill=none path | `emit.go` per-tag `nonVisualAttr` | Move to non-visual section for `path` when baseline has `fill="none"` |
| `stroke-linecap` default | `overlay.go` | Add `case an == "stroke-linecap": return "round", true` |
| `stroke-miterlimit` default | `overlay.go` | Already returns `"4"` (the fix for the invalid "0"). Consider returning a more distinctive value like `"1"` (minimum, sharp clip) or `"10"` (generous threshold) so the card is visually distinct; `"4"` at the baseline stroke-width of 2 produces no perceptible miter at the T-curve join |

---

## Top remaining systematic issues (cross-element)

### 1. `requiredExtensions=""` and `systemLanguage="en"` still in main grid — ALL 7 elements (14 cards)

**Problem:** These are now fixed to show visible shapes (round 1 regressions resolved), but both cards render identically to the element baseline. They are conditional-processing attributes — their visual effect is either "element shown" or "element hidden". With the fix values (`""` and `"en"`) the element is always shown, making the cards informationally redundant duplicates of the baseline.

**Recommended fix:** Add `"requiredextensions"` and `"systemlanguage"` to `nonVisualAttr()` in `emit.go` and move them to the collapsed non-visual details section (they belong with the conditional/accessibility metadata, not the rendering attrs).

**Fix target:** `chrome-testing/cmd/gen/emit.go` lines 90–100, `nonVisualAttr()` switch.

---

### 2. `pathLength` in main grid without stroke-dasharray — ALL 7 elements (7 cards)

**Problem:** `pathLength="2"` only affects how `stroke-dasharray` lengths are interpreted (it rescales the path length for dash calculations). With a solid baseline stroke (`stroke-dasharray` absent), `pathLength` has no visible effect at all — the card is identical to the baseline.

**Recommended fix (two options):**
- Option A (simpler): Add `"pathlength"` to `nonVisualAttr()` — move it to the non-visual details section.
- Option B (more illustrative): In `blueprint.go`, add `stroke-dasharray="10 5"` to baselines for elements where `pathLength` is present, then `pathLength="2"` visibly rescales the dash pattern.

**Fix target:** `chrome-testing/cmd/gen/emit.go` `nonVisualAttr()` OR `chrome-testing/cmd/gen/blueprint.go`.

---

### 3. `stroke-linecap="butt"` and `stroke-linejoin="miter"` default values — ALL elements (14 cards)

**Problem:** The grammar's `StrokeLinecapAttr` and `StrokeLinejoinAttr` enumerate keywords starting with `"butt"` and `"miter"` respectively — both SVG defaults. The enumerator picks the first arm, so the gallery always shows the default value, which is visually indistinguishable from the baseline.

**Additionally:** `stroke-linecap` has no visual effect on closed shapes (`rect`, `circle`, `ellipse`, `polygon`) since closed shapes have no open endpoints. These are NO_EFFECT cards for those elements.

**Recommended fix:**
```go
// overlay.go overlaySample()
case an == "stroke-linecap":
    return "round", true   // visibly distinct on all open-path elements
case an == "stroke-linejoin":
    return "round", true   // visibly distinct on corners/joins
```
For closed shapes (`rect`, `circle`, `ellipse`, `polygon`): also add `stroke-linecap` to per-tag non-visual (no open endpoints). A per-tag version of `nonVisualAttr` is needed (or inline the check in `emitPage`).

**Fix target:** `chrome-testing/cmd/gen/overlay.go` `overlaySample()`.

---

### 4. `paint-order` / `marker*` / rendering-hint attrs still in main grid — ALL 7 elements (~30 cards)

**Problem:** The following attributes are in the main grid but produce no visible change from the baseline for any shape element with the current values emitted:
- `paint-order="normal"` — default; only `"stroke fill"` or `"markers fill stroke"` would differ
- `marker="none"` / `marker-start="none"` / `marker-mid="none"` / `marker-end="none"` — values are explicitly `none`, so no markers appear; identical to baseline
- `color-interpolation="auto"` — rendering hint; imperceptible on rasterized static SVG
- `color-rendering="auto"` — rendering hint
- `shape-rendering="auto"` — rendering hint
- `vector-effect="none"` — default; `non-scaling-stroke` would show a visible difference
- `filter="none"` — no filter applied; identical to baseline
- `color-interpolation-filters="auto"` — only matters inside filter primitives; irrelevant on shapes without filters
- `flood-color="currentColor"` — only affects `feFlood` filter primitive; irrelevant on shapes

**Recommended fix:** Extend `nonVisualAttr()` with these names, OR steer them to non-default values:
- `paint-order`: steer to `"stroke fill"` so stroke is painted over fill (visible with thick stroke)
- `marker*`: steer to `url(#slot)` where the slot is a defined arrowhead marker (requires marker in blueprint defs)
- `vector-effect`: steer to `"non-scaling-stroke"` for a visible effect
- `color-interpolation`: steer to `"linearRGB"` (visible only in gradients) — acceptable to move to non-visual for plain shapes

**Fix targets:** `chrome-testing/cmd/gen/emit.go` `nonVisualAttr()`, `chrome-testing/cmd/gen/overlay.go` `overlaySample()`, optionally `chrome-testing/cmd/gen/blueprint.go` defs for marker.

---

### 5. `stroke="none"` blank on stroke-only elements — `line`, `path` (2 cards)

**Problem:** `line` and `path` baselines use stroke as the only visible paint (`fill="none"`). When `stroke="none"` is emitted, both paint sources are absent and the element is completely invisible.

**Recommended fix:** In `overlay.go`, for elements where the baseline has `fill="none"` (i.e., stroke-only elements), steer the `stroke` attribute to a color rather than `"none"`:
```go
case an == "stroke" && (tag == "line" || tag == "path" || tag == "polyline"):
    if valueKind == "keyword" {
        return "#ff6b35", true  // bright orange, clearly non-default
    }
```
Or alternatively, in `blueprint.go`, add a small visible fill to the `path` baseline (e.g., `fill="#16c79a" fill-opacity="0.2"`) so `stroke="none"` still shows the fill.

**Fix target:** `chrome-testing/cmd/gen/overlay.go` `overlaySample()` OR `chrome-testing/cmd/gen/blueprint.go` `baselineFor("path")`.

---

### 6. `color="#e94560"` ineffective — ALL 7 elements (7 cards)

**Problem:** `color` sets the `currentColor` keyword value but the baseline fill and stroke attributes use explicit hex colors, so `color` has no visual effect. The gallery shows the same shape regardless of the `color` value.

**Recommended fix (two options):**
- Option A: In `blueprint.go`, change the baseline fill/stroke for one element to use `fill="currentColor"` so the `color` card shows the shape changing hue.
- Option B: Add `"color"` to `nonVisualAttr()` for elements that don't use `currentColor` in their baseline.

**Fix target:** `chrome-testing/cmd/gen/blueprint.go` OR `chrome-testing/cmd/gen/emit.go`.

---

### 7. `flood-color` / `color-interpolation-filters` on non-filter elements — ALL 7 elements (14 cards)

**Problem:** `flood-color` is a presentation attribute that only affects the `feFlood` filter primitive. `color-interpolation-filters` only affects filter primitives. Neither has any visual effect on shapes that have no `filter` attribute. Both appear in the main grid for all 7 shape elements.

**Recommended fix:** Add `"flood-color"` and `"color-interpolation-filters"` to `nonVisualAttr()` in `emit.go`. These should be in the non-visual section for all non-filter elements.

**Fix target:** `chrome-testing/cmd/gen/emit.go` `nonVisualAttr()`.
