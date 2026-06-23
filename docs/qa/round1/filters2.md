# Visual QA Round 1 — filters2.md
Batch: feBlend, feComposite, feMerge, feMergeNode, feMorphology, feOffset, feFlood, feDropShadow, feImage, feTile, feTurbulence, feConvolveMatrix, feDisplacementMap
Date: 2026-06-22

---

## feBlend

**Category: IDENTICAL_NO_EFFECT**

All 100 cards render the same solid blue (`#4d8bff`) filled rectangle. No blend effect is visible on any card.

**Root cause:** The blueprint wires only `in` (or no `in2`) on the tested attribute. The `feBlend` primitive requires both `in` and `in2` to produce a visible composite. The supporting `feFlood` result (`layer2`) is declared but never wired as `in2="layer2"`. Without `in2`, the browser blends `in` over a transparent black — producing the plain source colour.

**Problem cards (all groups):**
- `in=*` group (all 10 in-value variants): plain blue, no blend visible
- `in2=*` group: identical blue, in2 is set but in defaults to SourceGraphic with no second layer wired
- `mode=*` group (normal, multiply, screen, darken, lighten, overlay, color-dodge, color-burn, hard-light, soft-light, difference, exclusion, hue, saturation, color, luminosity): all cards look identical blue — no two-source blend shown
- All other attribute-only cards (x, y, width, height, result, id, tabindex, lang, etc.): plain blue rect, no visual change

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go`: For `feBlend`, the blueprint must wire `in2="layer2"` on the `<feBlend>` element (i.e., `<feBlend in="SourceGraphic" in2="layer2">`). The existing `feFlood result="layer2"` is correct but unused. The attribute under test should override only the specific attribute (e.g., when testing `mode`, keep `in="SourceGraphic" in2="layer2"`; when testing `in`, override `in` but still keep `in2="layer2"`; when testing `in2`, keep `in="SourceGraphic"` and override `in2`).

---

## feComposite

**Category: HAS_EMPTY_CARDS / IDENTICAL_NO_EFFECT**

All 113 cards appear nearly white (empty) or show only the card background — no visible composite shape in any card. The two-layer blueprint (two `feFlood` rects) is present but the `feComposite` primitive is not wired to use them; it takes `in` from the slot attribute under test (often `SourceGraphic` or `SourceAlpha` on the white source rect) and has no `in2`, so the composite result is either fully transparent or fully white.

**Problem cards:**
- All `in=*` cards (10 variants): empty/white
- All `in2=*` implied cards (operator tests): empty/white
- All `operator=*` cards (over, in, out, atop, xor, arithmetic): empty/white with no two-layer overlap visible
- All `k1`, `k2`, `k3`, `k4` cards: empty/white — arithmetic formula has no visible inputs
- All position/dimension/metadata cards: empty/white

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go`: For `feComposite`, the blueprint must wire `in="layer1" in2="layer2"` on the `<feComposite>` element. Currently the two `feFlood` layers are created but the element uses `in="SourceGraphic"` (the white rect), not the flood layers. When testing non-`in`/`in2` attributes (operator, k values), always hard-wire `in="layer1" in2="layer2"` in the base blueprint; the attribute under test overrides only its own attr.

---

## feMerge

**Category: IDENTICAL_NO_EFFECT**

All 66 cards show an identical solid orange (`#f5a623`) filled rectangle. The merge effect itself is invisible because a single solid fill is being shown — there is no multi-layer merge or transparency to demonstrate. However, the cards are not blank; the filter does produce output.

**Assessment:** This is the expected minimal rendering (the orange rect is the baseline `feMerge` output from a single `feMergeNode` piping `SourceGraphic` through). The attribute variations (x, y, width, height, result, lang, etc.) are all presentation attributes that do not alter the rendered appearance of a single-node merge. No cards are blank.

**Problem cards:**
- None are empty, but all are visually identical — the merge effect with a single input cannot be distinguished from a plain rect.

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feMerge): Use at least two `<feMergeNode>` children — one with the source, one with a coloured `feFlood` result — so the merge layering is visible. Example: `<feMergeNode in="layer1"/><feMergeNode in="SourceGraphic"/>`. This would show the orange over the orange flood, or a semi-transparent overlay, making the merge visible.

---

## feMergeNode

**Category: IDENTICAL_NO_EFFECT**

All 57 cards show a uniform pink/red solid rectangle. The `feMergeNode` is the child element; the parent `feMerge` is presumably a single-node merge. Like `feMerge`, no layering effect is apparent.

**Problem cards:**
- `in=*` cards (SourceGraphic, SourceAlpha, BackgroundImage, BackgroundAlpha, FillPaint, StrokePaint, blur1, result1): all same solid pink — `SourceAlpha` gives opaque alpha which is invisible on dark bg, `BackgroundImage`/`BackgroundAlpha` are unsupported in most contexts (blank or fallback)
- All other attribute cards: identical solid pink

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feMergeNode): Blueprint should include two `feMergeNode` children with different coloured flood layers so the merge stacking is visible. The attribute under test should be on the node being varied.

---

## feMorphology

**Category: IDENTICAL_NO_EFFECT**

All 77 cards render bold "SVG" text in orange — which is correct source content — but all look identical. The morphology effect (erode/dilate) is not visible because no `operator` or `radius` is varied in the non-operator cards, and the baseline uses default `operator` (erode with radius 0), which produces no change.

**Problem cards:**
- `in=*` cards: "SVG" text renders correctly when `in="SourceGraphic"` but appears identical across all in-value variants. Cards with `in="BackgroundImage"`, `in="BackgroundAlpha"`, `in="blur1"`, `in="result1"` reference non-existent filter intermediates and silently fall back.
- `operator=*` and `radius=*` cards (if present): would show erode/dilate — these would be visually useful if the baseline radius is set to a visible value.

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feMorphology): Set a non-zero default `radius` (e.g. `radius="3"`) and `operator="dilate"` in the baseline so the effect is visible across all attribute-sweep cards. The attribute under test overrides its own attr only.

---

## feOffset

**Category: MOSTLY_PERFECT**

All 86 cards render a visible orange rectangle. The offset effect (`dx`, `dy`) is not dramatically visible because the orange fill fills the entire filter viewport — the shift moves the rect but the card background (near-black) fills the void, making small offsets indistinguishable from no-offset. The `in=*` and pure metadata cards are indistinguishable from each other.

**Problem cards (minor):**
- Cards with `dx` or `dy` values are nominally correct but visual difference is subtle — the rect shifts but the card area fills uniformly with orange since the source is a full-size rect. A small rect or non-filling source would make the offset visible.
- Cards with `in="BackgroundImage"` / `in="BackgroundAlpha"` reference unsupported inputs — these fall back silently and look identical to `in="SourceGraphic"`.

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feOffset): Change the source from a full-canvas rect to a smaller centered rect (e.g., `x="25" y="25" width="50" height="50"`) so that the offset displacement is clearly visible against the dark background.

---

## feFlood

**Category: IDENTICAL_NO_EFFECT**

All 66 cards render the same solid pink/crimson rectangle. The `feFlood` fills the filter region with the flood colour, which is correct for the baseline, but:
- The `flood-color` attribute variants are not being swept (the label shown is position/metadata attrs)
- All cards look identical regardless of the attribute tested (x, y, width, height, lang, tabindex, etc.) because none of these affect the solid fill output

**Problem cards:**
- All 66 cards: identical solid pink/crimson — no variation in colour, opacity, or region is visible.

**Root cause:** `feFlood` produces output regardless of input, but the blueprint is only sweeping non-colour/non-opacity attrs. The `flood-color` and `flood-opacity` attrs (which are the meaningful ones for this primitive) should be the primary sweep targets, with visible value contrast.

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feFlood): Ensure `flood-color` is swept with clearly distinct values (e.g. red, green, blue, transparent). Sweep `flood-opacity` with values 0, 0.5, 1. The current blueprint sweeps position/metadata attrs which produce no perceptible change.
- `reps.go` or grammar `.ebnf` for feFlood: Confirm `flood-color` and `flood-opacity` are enumerated as primary attrs.

---

## feDropShadow

**Category: MOSTLY_PERFECT**

All 87 cards show a visible orange rectangle with what appears to be a drop shadow. The effect is present and cards render correctly. Variations in `dx`, `dy`, `stdDeviation`, position/size attrs are apparent from labels.

**Minor issues:**
- Cards with `in="BackgroundImage"` / `in="BackgroundAlpha"` silently fall back (no `in` attr on feDropShadow in SVG spec — this may be from a shared attr sweep).
- The shadow darkness and blur are subtle due to the dark card background — visible but hard to compare. Using a lighter card background for shadow demonstration would help.
- `stdDeviation` card variations are present. Cards with large `dx`/`dy` values (e.g. `dx="3.14"`, `dy="-1"`) show shifted shadows.

**Fix targets:**
- No critical fixes required. Optional improvement: `chrome-testing/cmd/gen/blueprint.go` (feDropShadow): set a white or light-grey source rect fill so drop shadows are more visible against the dark card background.

---

## feImage

**Category: HAS_EMPTY_CARDS (all)**

All 74 cards are entirely dark/empty — no image content renders in any card.

**Root cause:** The blueprint uses `href="#slot"` which is a self-referential filter ID fragment — a filter cannot reference itself as its own image source. The element has no actual image URL or embedded data URI, so the `feImage` loads nothing and produces transparent output. All cards that don't override `href` (i.e., the 68+ position/metadata cards) also have no image source.

**Problem cards:**
- `href="#slot"` — self-referential, produces nothing (all cards share this issue)
- `preserveAspectRatio=*` — no image source, all dark
- `crossorigin=*` — no image source, all dark
- `x=*`, `y=*`, `width=*`, `height=*`, `result=*`, `id=*`, `tabindex=*`, `lang=*`, `xml:lang=*`, `xml:space=*`, `fill=*`, `stroke=*`, etc. — all dark

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feImage): Replace the placeholder `href="#slot"` with either:
  a. An inline SVG data URI: `href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' width='100' height='100'><rect width='100' height='100' fill='%23f5a623'/><circle cx='50' cy='50' r='30' fill='%234d8bff'/></svg>"`, or
  b. A reference to a small PNG asset bundled with the test suite.
  The source element in the SVG should not have `filter="url(#slot)"` applied to the same element used as `href` — use a separate visible element as the source.

---

## feTile

**Category: HAS_EMPTY_CARDS / IDENTICAL_NO_EFFECT**

All 74 cards render as a solid white rectangle. The tiling effect is completely absent.

**Root cause:** The blueprint creates a `feFlood` patch (`x=0, y=0, width=25%, height=25%, result="patch"`) but `feTile` uses `in="SourceGraphic"` instead of `in="patch"`. `SourceGraphic` is the full white source rect, so tiling it produces a full white rect — no tile pattern is visible.

**Problem cards (all 74):** Identical solid white, no tile pattern.

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feTile): Wire `in="patch"` on `<feTile>` so it tiles the small `feFlood` patch. The blueprint should read:
  ```
  <feFlood flood-color="#f5a623" x="0" y="0" width="25%" height="25%" result="patch"/>
  <feTile in="patch"/>
  ```
  When testing `in=*` attribute variants, override only `in`; for all other attrs, lock `in="patch"`.

---

## feTurbulence

**Category: MOSTLY_PERFECT — with GRAMMAR_ISSUE**

Most cards show turbulence/noise texture. The screenshot is very dense (87 cards) and the noise is rendered. However:

**Issues identified:**
- `numOctaves="-1"` card: negative octave count is invalid per the SVG spec (must be integer ≥ 0). Chrome silently clamps or ignores it, likely producing the same as `numOctaves="0"` (a flat field). This is a **GRAMMAR_ISSUE** — negative integers should not be enumerated for `numOctaves`.
- `numOctaves="100"` card: extremely high octave count. Browser likely clamps to a maximum (often 8). Produces a render but the value is out of practical range.
- `baseFrequency="2"`, `"2 3"`, `"3 3"`, `"1 2"` cards: very high frequencies produce uniform noise (near-solid). The effect is technically correct but visually indistinguishable from a solid dark field. Lowering the baseline frequency range (0.01–0.5) would show more useful texture variation.
- Cards for metadata attrs (lang, tabindex, xml:space, fill, stroke, etc.) all show identical noise — correct, as these don't affect feTurbulence output.

**Fix targets:**
- `reps.go` or grammar `.ebnf` for feTurbulence `numOctaves`: Remove negative integer values from the enumeration. The SVG spec states `numOctaves` is a non-negative integer.
- `chrome-testing/cmd/gen/blueprint.go` (feTurbulence): Use `baseFrequency="0.05"` as the baseline (instead of defaulting to `"2"` or higher) so noise patterns are clearly visible on all cards.

---

## feConvolveMatrix

**Category: HAS_EMPTY_CARDS**

All 118 cards render as a near-black/very-dark rectangle. No convolution effect is visible.

**Root cause:** `feConvolveMatrix` requires a valid `kernelMatrix` that matches the `order` (number of rows × columns). Without a valid matching `kernelMatrix` + `order` pair, the browser either produces transparent/black output or applies an identity that nulls out the source. The baseline sets no `kernelMatrix` on cards that only sweep `in`, `order`, or `in` variants — those cards produce black.

**Problem cards:**
- All `in=*` cards (8 variants): black — no `kernelMatrix` specified, browser rejects the filter
- `order="0.5"` card: non-integer order is invalid (must be positive integer) — **GRAMMAR_ISSUE**
- `kernelMatrix="1 2 3 4"` with default `order="3"` (9 elements expected, only 4 provided): mismatch → black
- `kernelMatrix="0 1 0 0"` (4 values) with default `order`: mismatch → black
- `kernelMatrix="1 0.5"` (2 values): mismatch → black
- All metadata attr cards (lang, tabindex, etc.): black

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feConvolveMatrix): Always include a matching valid `kernelMatrix` + `order` pair in the baseline. Use `order="3" kernelMatrix="0 -1 0 -1 5 -1 0 -1 0"` (sharpening kernel) as the baseline. When sweeping `kernelMatrix` values, ensure the number of values matches the `order`. When sweeping `order`, update `kernelMatrix` to match.
- `reps.go` or grammar `.ebnf` for `order`: Remove fractional values (`0.5`) — `order` must be a positive integer per spec.
- `reps.go` or grammar `.ebnf` for `kernelMatrix`: Enumerated values should all be length-matched to the default order (9 values for order=3).

---

## feDisplacementMap

**Category: MOSTLY_PERFECT**

Most cards render a visible orange rectangle with a displacement/warping effect. The blueprint correctly includes a `feTurbulence` noise map (`result="noiseMap"`) as a supporting primitive. The displacement is subtle (low `scale`) but the effect is present.

**Issues:**
- `in2=*` cards: these cards set `in2` to various named inputs (SourceGraphic, SourceAlpha, etc.) but the noise map is named `noiseMap`. When `in2` is overridden to `"SourceGraphic"`, the displacement uses the source itself as the map — producing no visible displacement (uniform colour source has no luminance variation). Cards with `in2="SourceAlpha"` similarly produce no visible warp. These look identical to the un-displaced rect.
- `in="BackgroundImage"`, `in="BackgroundAlpha"` cards: unsupported inputs — fall back silently to SourceGraphic-like behaviour.
- `scale="0"` card (if present): no displacement — visually identical to baseline.
- All metadata cards (lang, tabindex, xml:space, fill, stroke, etc.): identical to baseline orange rect — correct and expected.

**Fix targets:**
- `chrome-testing/cmd/gen/blueprint.go` (feDisplacementMap): When sweeping `in2`, keep the noise map available but also add a note that `in2="SourceGraphic"` will show a no-op visually. Consider adding a descriptive comment or fallback in the gallery. Alternatively, when `in2` is being swept, set `scale` to a higher value (e.g., `scale="30"`) so even moderate displacement from the source is visible.
- Consider wiring `in2="noiseMap"` as the locked baseline for all non-`in2` attr cards so displacement is consistently visible.

---

## Summary Table

| Element | Category | Key Issue |
|---|---|---|
| feBlend | IDENTICAL_NO_EFFECT | `in2` not wired to `layer2` — no two-source blend |
| feComposite | HAS_EMPTY_CARDS | `in`/`in2` not wired to flood layers — all empty/white |
| feMerge | IDENTICAL_NO_EFFECT | Single node merge on solid fill — effect invisible |
| feMergeNode | IDENTICAL_NO_EFFECT | Same as feMerge — single layer, no stacking visible |
| feMorphology | IDENTICAL_NO_EFFECT | Default radius=0, no erode/dilate visible |
| feOffset | MOSTLY_PERFECT | Full-canvas source hides offset; subtle only |
| feFlood | IDENTICAL_NO_EFFECT | Correct output but all cards same colour; no colour/opacity sweep |
| feDropShadow | MOSTLY_PERFECT | Shadow present; dark background reduces contrast |
| feImage | HAS_EMPTY_CARDS (all 74) | Self-referential href, no image source — all blank |
| feTile | HAS_EMPTY_CARDS / IDENTICAL_NO_EFFECT | `in` not wired to the patch flood — tiles full white rect |
| feTurbulence | MOSTLY_PERFECT | `numOctaves="-1"` grammar issue; high baseFrequency too coarse |
| feConvolveMatrix | HAS_EMPTY_CARDS (all 118) | Missing/mismatched kernelMatrix+order — all black; `order="0.5"` grammar issue |
| feDisplacementMap | MOSTLY_PERFECT | `in2` sweep cards show no warp (uniform source used as map) |

---

## Top Recurring Issues

1. **`in`/`in2` not wired to supporting primitives** (feBlend, feComposite, feTile): The blueprint declares named intermediate results (flood layers, patches) but the tested primitive uses `SourceGraphic` or the wrong input, making supporting layers invisible. **Fix in `blueprint.go`**: lock `in`/`in2` to the named intermediates for all non-`in`/`in2` attribute sweeps.

2. **No valid baseline parameters for complex primitives** (feConvolveMatrix, feImage): These primitives require specific mandatory or co-dependent attrs (`kernelMatrix`+`order`, `href` with real content) to produce any output. Without them, all cards are blank/black. **Fix in `blueprint.go`**: always emit a complete valid baseline set of required attrs.

3. **Effect invisible on uniform sources** (feMorphology, feOffset, feMerge): When the source is a solid full-canvas rect, operations like morphology, offset, and merge produce no visually distinguishable result. **Fix in `blueprint.go`**: use smaller source shapes (small centred rects, text, circles) and add a second contrast layer so the effect has something to act on.

4. **Grammar issues — invalid values enumerated** (feTurbulence `numOctaves="-1"`, feConvolveMatrix `order="0.5"`): The grammar/reps generate values outside the spec's valid range. **Fix in `.ebnf` or `reps.go`**: clamp numOctaves to non-negative integers; clamp order to positive integers only.

5. **feImage self-referential href** (feImage `href="#slot"`): The filter ID is reused as the image href, creating a circular reference that renders nothing. **Fix in `blueprint.go`**: embed a data URI SVG or reference an external asset.
