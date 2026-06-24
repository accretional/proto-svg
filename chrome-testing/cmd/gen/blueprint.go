package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// blueprint.go — the scaffold that supplies ONLY an element's NEIGHBORS, with a
// {{ELEMENT}} slot the generator fills with grammar-generated element markup
// (TEMPLATE_GUIDE convention). Everything about the element itself is generated;
// the blueprint never predetermines the element's own attributes.
//
// For element E:
//   - if chrome-testing/html/template/<tag>.html exists, extract the
//     <script type="application/svg-blueprint" id="blueprint"> body and use it.
//   - else use a built-in DEFAULT scaffold derived from the content model:
//       shapes/text/image/containers  → self-render: <svg>{{ELEMENT}}</svg>
//       gradient/pattern              → {{ELEMENT}} in <defs>, a rect fill=url(#slot)
//       marker                        → {{ELEMENT}} in <defs>, a path marker-end=url(#slot)
//       clipPath/mask                 → {{ELEMENT}} in <defs>, a rect clip-path/mask=url(#slot)
//       filter primitive              → {{ELEMENT}} inside <filter id="slot">, rect filter=url(#slot)
//       filter                        → {{ELEMENT}} is the filter (id="slot"), rect filter=url(#slot)
//       stop                          → {{ELEMENT}} inside a gradient a rect uses
//       gradient stop/feFunc/light    → nested in the proper parent
//       animation elements            → {{ELEMENT}} as a child of an animated shape

// blueprintCache memoizes resolved blueprints per tag.
type blueprintProvider struct {
	templateDir string
	cache       map[string]string
}

func newBlueprintProvider(templateDir string) *blueprintProvider {
	return &blueprintProvider{templateDir: templateDir, cache: map[string]string{}}
}

var blueprintRe = regexp.MustCompile(`(?s)<script[^>]*type="application/svg-blueprint"[^>]*>(.*?)</script>`)

// builtinScaffoldWins lists tags whose corrected built-in scaffold (in
// defaultScaffold) must take precedence over any on-disk template HTML. These
// scaffolds were fixed in QA round2 (defs renders a child shape; lighting/
// transfer/convolve get a textured input) and the older template files would
// otherwise shadow the fix.
var builtinScaffoldWins = map[string]bool{
	"defs": true,
	// feBlend/feComposite need the corrected catFilterPrimitive scaffold: a
	// multi-color SourceGraphic blended against a CONTRASTING full-region in2
	// flood so the blend MODE cards visibly differ. The on-disk templates blend a
	// flat source against a flat flood (modes indistinguishable), so override them.
	"feBlend":             true,
	"feComposite":         true,
	"feDiffuseLighting":   true,
	"feSpecularLighting":  true,
	"feDistantLight":      true,
	"fePointLight":        true,
	"feSpotLight":         true,
	"feComponentTransfer": true,
	"feFuncR":             true,
	"feFuncG":             true,
	"feFuncB":             true,
	"feFuncA":             true,
	"feConvolveMatrix":    true,
	"feDropShadow":        true,
	// QA round3: the on-disk templates for these tags shadow the corrected
	// built-in scaffolds. Force the built-ins so the round3 fixes take effect:
	//   linearGradient → refgrad sibling def + non-zero x1 baseline (FIX 2)
	//   stop           → anchor moved to 0.3 + color on root (FIX 3)
	//   clipPath/mask  → bodyOverride fractional OBB children (FIX 1)
	//   pattern        → bodyOverride fractional OBB content (FIX 1)
	//   feColorMatrix  → hueRotate baseline (FIX 6, via baselineFor+overlay)
	//   feTile         → non-uniform patch tiling (FIX 5, via catFilterPrimitive)
	//   use            → larger referenced target circle (FIX 7)
	//   switch         → fill=currentColor child + root color (FIX 8)
	//   script         → pre-colored target rect (FIX 9)
	//   mpath/discard  → larger brightly-colored host shape (FIX 10)
	"linearGradient": true,
	"stop":           true,
	"clipPath":       true,
	"mask":           true,
	"pattern":        true,
	"feColorMatrix":  true,
	"feTile":         true,
	// FIX 5: the on-disk feTurbulence template uses a tight filter region, so the
	// enumerated primitive-subregion x/y/width/height values clip the noise to a
	// zero/empty area off the visible rect → blank cards. Force the corrected
	// built-in, which uses a wide filter region and a full-host rect so the
	// subregion always intersects the visible area.
	"feTurbulence": true,
	"use":            true,
	"switch":         true,
	"script":         true,
	"mpath":          true,
	"discard":        true,
	// FIX 6: the on-disk set template hosts the animation on a 40-wide rect at x=0,
	// so set attributeName="x" to="80" slides it mostly off the 100-wide viewBox
	// (a narrow sliver). Force the built-in, which uses a small centered host so the
	// whole from/to range keeps the shape fully visible.
	"set": true,
	// The on-disk animate/animateTransform/animateMotion blueprints host the
	// animation on a bare <rect> with NO id="target", so the href="#target" preset
	// (and the overlay's target resolution) dangles and the animation freezes.
	// Force the built-in catAnimation scaffold, whose host carries id="target"
	// (and a starting x=10), so every preset animates and href resolves.
	"animate":          true,
	"animateTransform": true,
	"animateMotion":    true,
	// FIX 2 (wrapper inheritance): the on-disk templates for these wrapper/
	// container tags use a colorless root, so the child's fill="currentColor"
	// resolves to black (invisible on the dark card) and paint set on the wrapper
	// does not produce a visible base. Force the corrected built-ins, which seed a
	// neutral root color (svgOpenColor) and a currentColor-painting child so
	// fill(currentColor)/color/opacity/transform varied on the wrapper give
	// distinct, visible cards.
	"g":             true,
	"a":             true,
	"svg":           true,
	"symbol":        true,
	"foreignObject": true,
	// FIX 3 (text color): the on-disk text/tspan/textPath templates hardcode the
	// glyph fill (#e6e6e6 / #e94560), so the CSS `color` property never reaches the
	// text and `color`/`fill`(currentColor) cards look identical. Force the
	// corrected built-ins, which seed a neutral root color and paint the text with
	// fill="currentColor" (baselineFor) so varying color visibly recolors it.
	"text":     true,
	"tspan":    true,
	"textPath": true,
}

// blueprintFor returns the scaffold for tag (with a {{ELEMENT}} placeholder).
func (b *blueprintProvider) blueprintFor(tag string) string {
	if s, ok := b.cache[tag]; ok {
		return s
	}
	var s string
	if builtinScaffoldWins[tag] {
		// Use the corrected built-in scaffold; ignore any stale template file.
		s = defaultScaffold(tag)
	} else {
		s = b.fromTemplate(tag)
		if s == "" {
			s = defaultScaffold(tag)
		}
	}
	b.cache[tag] = s
	return s
}

// fromTemplate reads chrome-testing/html/template/<tag>.html and extracts the
// blueprint script body, or "" if absent.
func (b *blueprintProvider) fromTemplate(tag string) string {
	path := filepath.Join(b.templateDir, tag+".html")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	m := blueprintRe.FindSubmatch(data)
	if m == nil {
		return ""
	}
	body := strings.TrimSpace(string(m[1]))
	if !strings.Contains(body, "{{ELEMENT}}") {
		return ""
	}
	return body
}

// inject substitutes the variant's element markup into the blueprint's
// {{ELEMENT}} slot.
func inject(blueprint, elementMarkup string) string {
	return strings.ReplaceAll(blueprint, "{{ELEMENT}}", elementMarkup)
}

// svgOpen is the standard root used by the default scaffolds.
const svgOpen = `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 100 100">`

// svgOpenColor is the root used by WRAPPER/CONTAINER scaffolds (a/g/svg/symbol/
// use/foreignObject). It seeds a neutral inherited color so a child painted with
// fill="currentColor" is visible by default; when the wrapper's fill (via
// currentColor)/color/opacity/transform is varied, the inheriting child visibly
// changes and the cards become distinct (FIX 2).
const svgOpenColor = `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 100 100" color="#16c79a">`

// imageDataURI is a self-contained 80x80 SVG (dark square + teal circle) used as
// the baseline href for <image> and <feImage> so they render offline instead of
// dangling (QA round1). It is an inline base64 data: URI — no network/fetch.
const imageDataURI = `data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI4MCIgaGVpZ2h0PSI4MCI+PHJlY3Qgd2lkdGg9IjgwIiBoZWlnaHQ9IjgwIiBmaWxsPSIjMjIyIi8+PGNpcmNsZSBjeD0iNDAiIGN5PSI0MCIgcj0iMzAiIGZpbGw9IiMxNmM3OWEiLz48L3N2Zz4=`

// defaultScaffold returns a built-in blueprint derived from tag's content model.
func defaultScaffold(tag string) string {
	// Tag-specific scaffolds (QA round2 baseline/visibility fixes) take priority
	// over the generic category scaffolds below.
	switch tag {
	case "defs":
		// <defs> is a non-rendering container, so a <use> of the <defs> itself
		// paints nothing (QA round2: all cards blank). Instead reference a CHILD
		// shape: the generated <defs> (carrying the varied attribute) sits beside
		// a sibling <defs> holding a renderable <rect id="defskid">, and <use>
		// instantiates THAT rect so every card shows the red square.
		return svgOpen +
			`{{ELEMENT}}` +
			`<defs><rect id="defskid" x="20" y="20" width="60" height="60" fill="#e94560"/></defs>` +
			`<use href="#defskid"/></svg>`
	case "feComponentTransfer", "feFuncR", "feFuncG", "feFuncB", "feFuncA":
		// A flat source makes every transfer curve a single flat shade. Feed a
		// horizontal black→white gradient so each pixel samples a different input
		// value and the transfer curve's shape (slope/gamma/table/discrete) is
		// visible across the rect (QA round2).
		// feComponentTransfer IS the filter child; feFunc* must be nested inside a
		// feComponentTransfer.
		var slot string
		if tag == "feComponentTransfer" {
			slot = `<filter id="slot">{{ELEMENT}}</filter>`
		} else {
			slot = `<filter id="slot"><feComponentTransfer>{{ELEMENT}}</feComponentTransfer></filter>`
		}
		// FIX 4: a COLORED source gradient (red→green→blue) instead of a black→white
		// greyscale ramp. With an achromatic source R=G=B, so manipulating a single
		// channel (feFuncR/G/B/A) is nearly imperceptible. A colored ramp makes each
		// channel's transfer curve show as a visible hue shift across the rect.
		//
		// COLLISION FIX: the RGB ramp is fully OPAQUE, so feFuncA (the ALPHA transfer)
		// sees a constant alpha=1 → every feFuncA type/param card is uniform and
		// identical. Add a VERTICAL alpha ramp (a white→transparent mask gradient,
		// vertical axis, orthogonal to the horizontal RGB axis) so the source alpha
		// varies top→bottom. feFuncA's transfer then reshapes a real alpha ramp
		// (distinct per type), while feFuncR/G/B keep reading the horizontal hue ramp.
		// The mask must be applied to the source BEFORE the filter consumes it: the
		// inner <rect> carries the alpha mask and is rendered first; the wrapping <g>
		// carries the filter, so the filter's SourceGraphic already has the vertical
		// alpha ramp (the filter→mask order on a single element would mask AFTER the
		// filter and feFuncA would still see constant alpha).
		return svgOpen +
			`<defs>` +
			`<linearGradient id="bpGrad" x1="0" y1="0" x2="1" y2="0">` +
			`<stop offset="0" stop-color="#ff0000"/><stop offset="0.5" stop-color="#00ff00"/><stop offset="1" stop-color="#0000ff"/>` +
			`</linearGradient>` +
			`<linearGradient id="bpAlpha" x1="0" y1="0" x2="0" y2="1">` +
			`<stop offset="0" stop-color="#fff" stop-opacity="1"/><stop offset="1" stop-color="#fff" stop-opacity="0.1"/>` +
			`</linearGradient>` +
			`<mask id="bpAlphaMask"><rect x="5" y="5" width="90" height="90" fill="url(#bpAlpha)"/></mask>` +
			slot +
			`</defs>` +
			`<g filter="url(#slot)"><rect x="5" y="5" width="90" height="90" fill="url(#bpGrad)" mask="url(#bpAlphaMask)"/></g></svg>`
	case "feDiffuseLighting", "feSpecularLighting":
		// Lighting computes a surface normal from the ALPHA gradient of its input.
		// A flat opaque rect has no alpha relief → uniform, flat shading. Prepend a
		// blurred SourceAlpha bump map and feed it via in="bumpMap" (set in
		// baselineFor) so surfaceScale/azimuth/elevation produce directional 3D
		// shading. A bright source rect keeps the lit result visible (QA round2).
		return svgOpen +
			`<defs><filter id="slot">` +
			`<feGaussianBlur in="SourceAlpha" stdDeviation="6" result="bumpMap"/>` +
			`{{ELEMENT}}</filter></defs>` +
			`<rect x="20" y="20" width="60" height="60" fill="#ffffff" filter="url(#slot)"/></svg>`
	case "feDistantLight", "fePointLight", "feSpotLight":
		// The generated element is a LIGHT SOURCE nested in a lighting primitive.
		// Prepend a blurred SourceAlpha bump map and feed it (in="bumpMap") so the
		// light's azimuth/elevation/position produces visible directional relief
		// instead of flat shading (QA round2). Bright source rect keeps it visible.
		return svgOpen +
			`<defs><filter id="slot">` +
			`<feGaussianBlur in="SourceAlpha" stdDeviation="6" result="bumpMap"/>` +
			`<feDiffuseLighting in="bumpMap" surfaceScale="5" diffuseConstant="1" lighting-color="#ffffff">{{ELEMENT}}</feDiffuseLighting>` +
			`</filter></defs>` +
			`<rect x="20" y="20" width="60" height="60" fill="#ffffff" filter="url(#slot)"/></svg>`
	case "feConvolveMatrix":
		// A convolution kernel on a flat-color rect produces the same flat color
		// (no edges to operate on). Wire the textured feTurbulence result
		// ("noiseMap", supplied below) via in="noiseMap" (baselineFor) so the
		// kernel effect is visible (QA round2).
		return svgOpen +
			`<defs><filter id="slot">` +
			`<feTurbulence type="turbulence" baseFrequency="0.08" numOctaves="2" seed="1" result="noiseMap"/>` +
			`{{ELEMENT}}</filter></defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="#f5a623" filter="url(#slot)"/></svg>`
	case "feDropShadow":
		// The default shadow is black and the card background is dark, so the
		// shadow is invisible. Use a bright source on a roomy filter region so the
		// offset/blurred shadow reads against the source (QA round2, best-effort).
		return svgOpen +
			`<defs><filter id="slot" x="-20%" y="-20%" width="150%" height="150%">{{ELEMENT}}</filter></defs>` +
			`<rect x="20" y="20" width="55" height="55" fill="#f5a623" filter="url(#slot)"/></svg>`
	case "set":
		// FIX 6: <set attributeName="x" to="80"> slides the host rect to x=80. With a
		// 40-wide host on a 100-wide viewBox the rect ran mostly off-canvas (a narrow
		// 20px sliver). Use a SMALL host (18×18, vertically centered) so the whole
		// from/to range (x: 10→80) keeps the shape fully visible (right edge ≤ 98).
		// The element is a child of the animated rect (id="target" so href resolves).
		return svgOpen +
			`<rect id="target" x="10" y="41" width="18" height="18" fill="#4d8bff">{{ELEMENT}}</rect></svg>`
	case "discard":
		// FIX 10: <discard> removes its target after `begin`; the static snapshot is
		// taken before the (60s, via baselineFor) discard fires, so the host should
		// be present and clearly visible. Use a large, bright centered circle (r=20)
		// instead of a tiny corner square. discard is a child of the target it
		// removes; host it directly inside a rendered shape.
		return svgOpen +
			`<circle id="target" cx="50" cy="50" r="20" fill="#4d8bff">{{ELEMENT}}</circle></svg>`
	case "script":
		// FIX 9: pre-color the target rect so cards whose <script type="…"> is a
		// non-JS MIME type (which the browser refuses to execute) are not a blank
		// grey square. The <script> is illustrative — it would recolor the rect IF
		// it ran — but the base color does not depend on JS execution. The script
		// body (bodyFor) targets id="slot-target".
		return svgOpen +
			`<rect id="slot-target" x="20" y="20" width="60" height="60" fill="#e94560"/>` +
			`{{ELEMENT}}</svg>`
	case "switch":
		// FIX 8: <switch> is a container; presentation attrs set on it (color,
		// fill-opacity, opacity, …) only show if its rendered child INHERITS them.
		// The child rect uses fill="currentColor" (see bodyFor), so a root
		// color="#16c79a" gives a visible base while attrs varied on <switch>
		// (e.g. color="#e94560", fill-opacity, opacity) produce distinct cards
		// (mirrors the <g> approach).
		return svgOpenColor + `{{ELEMENT}}</svg>`
	case "text":
		// FIX 3: text content paints with fill="currentColor" (see baselineFor) so
		// the CSS `color` property reaches the glyphs. The root seeds a neutral
		// color so the base card is visible; varying `color` (and `fill` via
		// currentColor) then visibly recolors the text.
		return svgOpenColor + `{{ELEMENT}}</svg>`
	case "tspan":
		// FIX 3: the wrapping <text> carries NO fill, so the tspan's own
		// fill="currentColor" (baselineFor) governs the glyph color and a varied
		// `color`/`fill` recolors it. The parent run has visible neutral text around
		// the {{ELEMENT}} slot so the tspan reads as a styled SUB-SPAN in context.
		return svgOpenColor +
			`<text x="8" y="55" fill="#7b8a82" font-size="16">Hi {{ELEMENT}}!</text></svg>`
	case "textPath":
		// FIX 3: <path id="slot"> for href="#slot" to follow; the wrapping <text>
		// carries no fill so the textPath's own fill="currentColor" (baselineFor)
		// governs the glyph color. Root seeds a neutral color for the base card.
		return svgOpenColor +
			`<defs><path id="slot" d="M 5 70 Q 50 10 95 70"/></defs>` +
			`<text font-size="20">{{ELEMENT}}</text></svg>`
	case "g", "a", "svg", "foreignObject":
		// FIX 2: self-rendering wrapper/container elements. Their rendered child
		// paints with fill="currentColor" (see bodyFor); the neutral root color makes
		// the base card visible, and varying the wrapper's inheritable paint (fill via
		// currentColor, color, opacity) — or its transform — now visibly changes the
		// child, so cards that were identical become distinct.
		return svgOpenColor + `{{ELEMENT}}</svg>`
	case "symbol":
		// FIX 2: <symbol> is NOT rendered directly — it must be instanced. Place the
		// generated <symbol> (carrying the varied attr, id="slot" via the slot
		// machinery) in <defs> and instance it with <use>. The symbol's child paints
		// fill="currentColor" and the root seeds a neutral color, so color/opacity/
		// fill-opacity/transform varied on the <symbol> inherit through the instance
		// and yield distinct, visible cards (base is teal, not invisible black).
		return svgOpenColor +
			`<defs>{{ELEMENT}}</defs>` +
			`<use href="#slot" x="10" y="10" width="80" height="80"/></svg>`
	case "feColorMatrix":
		// FIX 4: feColorMatrix must operate on a CLEAN, fully-saturated multi-color
		// source so the type=matrix channel-swap (and saturate/hueRotate) is visibly
		// distinct. The generic filter-primitive scaffold chains it off a low-
		// saturation turbulence result, where a matrix vs hueRotate difference is
		// muted. Operate directly on SourceGraphic (four bright quadrants) instead.
		// The 20-number matrix companion (companionFor) recolors those quadrants.
		return svgOpen +
			`<defs><filter id="slot">{{ELEMENT}}</filter></defs>` +
			`<g filter="url(#slot)">` +
			`<rect x="5" y="5" width="45" height="45" fill="#e94560"/>` +
			`<rect x="50" y="5" width="45" height="45" fill="#16c79a"/>` +
			`<rect x="5" y="50" width="45" height="45" fill="#4d8bff"/>` +
			`<rect x="50" y="50" width="45" height="45" fill="#ffd166"/>` +
			`</g></svg>`
	case "feTurbulence":
		// FIX 5: feTurbulence's x/y/width/height define the PRIMITIVE SUBREGION that
		// clips the generated noise. With the default tight filter region, percentage
		// subregion values (50%/75%) and the small length values resolve to a zero/
		// empty area off the visible rect → blank cards. Give the filter a WIDE region
		// (-20%..140%) and fill the whole 0..100 host so the subregion always lands
		// inside the visible area and the clipped noise patch shows on every card.
		return svgOpen +
			`<defs><filter id="slot" x="-20%" y="-20%" width="140%" height="140%">{{ELEMENT}}</filter></defs>` +
			`<rect x="0" y="0" width="100" height="100" fill="#4d8bff" filter="url(#slot)"/></svg>`
	case "feTile":
		// FIX 5: tiling a UNIFORM source produces a uniform output (no pattern), so
		// the round2 white-flood feTile showed nothing. Two things are needed:
		//  (1) A NON-UNIFORM 25×25 tile cell to tile — built from filter primitives
		//      only (a teal field feFlood + an orange corner feFlood composited
		//      over it → result="patch"); feImage data:/fragment patches do not
		//      paint reliably headless / under the gallery's duplicate ids.
		//  (2) The HOST graphic must itself be non-uniform. The `in` attribute is
		//      enumerated, so the FIRST card on the page is `in="SourceGraphic"`;
		//      because every card shares id="slot", that first filter (tiling a
		//      flat host) otherwise poisoned ALL cards to a uniform fill. Filling
		//      the host with a repeating <pattern> makes even `in="SourceGraphic"`
		//      tile a visible dot grid, so no card is blank.
		// feTile (in="patch", via baselineFor) repeats the cell across the filter
		// region; varying its x/y/width/height shifts/scales the tiling region.
		return svgOpen +
			`<defs>` +
			`<pattern id="tilehost" width="25" height="25" patternUnits="userSpaceOnUse">` +
			`<rect width="25" height="25" fill="#16c79a"/><circle cx="12" cy="12" r="6" fill="#f5a623"/></pattern>` +
			`<filter id="slot" x="0" y="0" width="100%" height="100%">` +
			`<feFlood flood-color="#16c79a" x="0" y="0" width="25" height="25" result="base"/>` +
			`<feFlood flood-color="#f5a623" x="0" y="0" width="13" height="13" result="dot"/>` +
			`<feComposite in="dot" in2="base" operator="over" x="0" y="0" width="25" height="25" result="patch"/>` +
			`{{ELEMENT}}</filter></defs>` +
			`<rect x="0" y="0" width="100" height="100" fill="url(#tilehost)" filter="url(#slot)"/></svg>`
	}
	switch category(tag) {
	case catGradient:
		// FIX 2(a): provide a concrete <linearGradient id="refgrad"> sibling so a
		// generated href="#refgrad" / xlink:href="#refgrad" actually resolves
		// instead of dangling (overlay routes gradient href to #refgrad).
		return svgOpen +
			`<defs>` +
			`<linearGradient id="refgrad"><stop offset="0" stop-color="#e94560"/><stop offset="1" stop-color="#16c79a"/></linearGradient>` +
			`{{ELEMENT}}` +
			`</defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="url(#slot)"/></svg>`
	case catPattern:
		return svgOpen +
			`<defs>{{ELEMENT}}</defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="url(#slot)"/></svg>`
	case catMarker:
		return svgOpen +
			`<defs>{{ELEMENT}}</defs>` +
			`<path d="M15 50 H85" fill="none" stroke="#16c79a" stroke-width="3" marker-end="url(#slot)" marker-start="url(#slot)"/></svg>`
	case catClip:
		return svgOpen +
			`<defs>{{ELEMENT}}</defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="#e94560" clip-path="url(#slot)"/></svg>`
	case catMask:
		return svgOpen +
			`<defs>{{ELEMENT}}</defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="#16c79a" mask="url(#slot)"/></svg>`
	case catFilter:
		return svgOpen +
			`<defs>{{ELEMENT}}</defs>` +
			`<rect x="20" y="20" width="60" height="60" fill="#e94560" filter="url(#slot)"/></svg>`
	case catFilterPrimitive:
		// Supply the supporting inputs the multi-input primitives reference by
		// result name (layer1/layer2/layerA/layerB/patch/noiseMap) so feBlend/
		// feComposite/feTile/feMerge/feDisplacementMap have a visible second input.
		// Single-input primitives ignore the unused results.
		//
		// feBlend mode cards only differ if the two blended layers CONTRAST and
		// overlap. So: SourceGraphic is a MULTI-COLOR graphic (two overlapping
		// colored rects in the host below) and `layer2` is a full-region CONTRASTING
		// solid flood that covers the whole source, so each blend MODE (normal/
		// multiply/screen/darken/lighten/…) produces a visibly different result.
		return svgOpen +
			`<defs><filter id="slot" x="-50%" y="-50%" width="200%" height="200%">` +
			`<feFlood flood-color="#f5a623" x="0" y="0" width="60%" height="60%" result="layer1"/>` +
			`<feFlood flood-color="#16c79a" x="0" y="0" width="100%" height="100%" result="layer2"/>` +
			`<feFlood flood-color="#e94560" x="0" y="0" width="60%" height="60%" result="layerA"/>` +
			`<feFlood flood-color="#4d8bff" x="40%" y="40%" width="60%" height="60%" result="layerB"/>` +
			`<feFlood flood-color="#f5a623" x="0" y="0" width="25%" height="25%" result="patch"/>` +
			`<feTurbulence type="turbulence" baseFrequency="0.05" numOctaves="2" seed="1" result="noiseMap"/>` +
			`{{ELEMENT}}</filter></defs>` +
			`<g filter="url(#slot)">` +
			`<rect x="20" y="20" width="50" height="50" fill="#e94560"/>` +
			`<rect x="40" y="40" width="45" height="45" fill="#f5a623"/>` +
			`</g></svg>`
	case catTransferFn:
		return svgOpen +
			`<defs><filter id="slot"><feComponentTransfer>{{ELEMENT}}</feComponentTransfer></filter></defs>` +
			`<rect x="10" y="10" width="80" height="80" fill="#4d8bff" filter="url(#slot)"/></svg>`
	case catLight:
		// The generated element is a LIGHT SOURCE (feDistantLight/fePointLight/
		// feSpotLight); it must be nested inside a lighting primitive. The baseline
		// gives it visible angles/positions so the lit surface is not black.
		return svgOpen +
			`<defs><filter id="slot"><feDiffuseLighting surfaceScale="5" diffuseConstant="1" lighting-color="#ffffff">{{ELEMENT}}</feDiffuseLighting></filter></defs>` +
			`<rect x="10" y="10" width="80" height="80" fill="#e94560" filter="url(#slot)"/></svg>`
	case catStop:
		// Two CONTRASTING fixed anchor stops bracket the varied stop so any
		// offset/stop-color change produces a visible gradient band that shifts and
		// recolors (QA round2). FIX 3: the leading anchor sits at offset="0.3" (was
		// 0.1) so a varied stop at offset="0" is clearly separated from it by a wide
		// red band; the trailing anchor is a strong blue for high contrast. The root
		// <svg> carries color="#f5a623" so a varied stop-color="currentColor"
		// resolves to a visible orange instead of the near-white CSS text color.
		return `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 100 100" color="#f5a623">` +
			`<defs><linearGradient id="slot"><stop offset="0.3" stop-color="#e94560"/>{{ELEMENT}}<stop offset="1" stop-color="#4d8bff"/></linearGradient></defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="url(#slot)"/></svg>`
	case catMergeNode:
		// Named colored layers so a feMergeNode baseline in="layerA" resolves and the
		// merge stacking is visible.
		return svgOpen +
			`<defs><filter id="slot">` +
			`<feFlood flood-color="#f5a623" x="0" y="0" width="70%" height="70%" result="layerA"/>` +
			`<feFlood flood-color="#4d8bff" x="30%" y="30%" width="70%" height="70%" result="layerB"/>` +
			`<feMerge>{{ELEMENT}}</feMerge></filter></defs>` +
			`<rect x="10" y="10" width="80" height="80" filter="url(#slot)"/></svg>`
	case catAnimation:
		return svgOpen +
			`<rect id="target" x="10" y="10" width="40" height="40" fill="#4d8bff">{{ELEMENT}}</rect></svg>`
	case catMpath:
		// The mpath's href resolves (via overlay) to "#target"; the motion-path def
		// id must match. Host shape starts mid-canvas so it is not clipped at (0,0).
		// FIX 10: a larger, bright host (40×40 #4d8bff, centered) so the element is
		// clearly visible in the static snapshot instead of a tiny corner square.
		return svgOpen +
			`<defs><path id="target" d="M30 50 Q50 10 70 50"/></defs>` +
			`<rect x="30" y="30" width="40" height="40" fill="#4d8bff"><animateMotion dur="2s" repeatCount="indefinite">{{ELEMENT}}</animateMotion></rect></svg>`
	case catDescriptive:
		return svgOpen +
			`<rect x="10" y="10" width="80" height="80" fill="#16c79a">{{ELEMENT}}</rect></svg>`
	case catContainerRef:
		// <use> references a defined shape; the baseline href is "#slot", so the
		// referenced def carries id="slot" (was "#ref" → dangling, QA round1).
		// FIX 7: reference a LARGE centered circle (was a tiny 40×40 top-left rect)
		// so the <use> instance renders prominently and x/y offsets shift it
		// visibly across cards.
		// FIX 2: the referenced circle paints with fill="currentColor" so paint
		// inherited through the <use> (fill via currentColor, color, opacity) reaches
		// the instance and varying it on the <use> yields distinct cards; the neutral
		// root color keeps the base instance visible.
		return svgOpenColor +
			`<defs><circle id="slot" cx="30" cy="40" r="20" fill="currentColor"/></defs>` +
			`<use href="#slot" opacity="0.3"/>` + // the faint ORIGINAL, so the <use> reads as a stamped COPY of it
			`{{ELEMENT}}</svg>`
	default: // shapes, text, image, g, svg, etc. — self-render
		return svgOpen + `{{ELEMENT}}</svg>`
	}
}

// content-model categories that drive scaffold selection.
type contentCategory int

const (
	catSelf contentCategory = iota
	catGradient
	catPattern
	catMarker
	catClip
	catMask
	catFilter
	catFilterPrimitive
	catTransferFn
	catLight
	catStop
	catMergeNode
	catAnimation
	catMpath
	catDescriptive
	catContainerRef
)

func category(tag string) contentCategory {
	switch tag {
	case "linearGradient", "radialGradient":
		return catGradient
	case "pattern":
		return catPattern
	case "marker":
		return catMarker
	case "clipPath":
		return catClip
	case "mask":
		return catMask
	case "filter":
		return catFilter
	case "stop":
		return catStop
	case "feMergeNode":
		return catMergeNode
	case "feFuncR", "feFuncG", "feFuncB", "feFuncA":
		return catTransferFn
	case "feDistantLight", "fePointLight", "feSpotLight":
		return catLight
	case "animate", "animateTransform", "set", "discard", "animateMotion":
		return catAnimation
	case "mpath":
		return catMpath
	case "desc", "title", "metadata":
		return catDescriptive
	case "use":
		return catContainerRef
	}
	if strings.HasPrefix(tag, "fe") {
		return catFilterPrimitive
	}
	return catSelf
}

// baselineFor returns the minimal extra attributes that make a freshly-generated
// element VISIBLE, so the one varied attribute can be observed. It deliberately
// omits the attribute the variant is varying (passed as varyingPrefix) to avoid
// a duplicate. The second result reports whether the element needs an id (unused
// for baselines but kept for symmetry with the slot mechanism).
//
// varyingValue is the concrete value chosen for the varied attribute. It lets
// the baseline carry a DEPENDENT COMPANION attribute that makes the varied attr
// effective: several filter-primitive attrs are inert unless a sibling attr is
// set (feComposite k1-k4 need operator="arithmetic"; feFunc amplitude/exponent/
// offset need type="gamma"; tableValues needs type="table"; …). companionFor
// supplies those companions (and may override a baseline pair, e.g. flipping
// feComposite's baseline operator="over" to "arithmetic").
func baselineFor(tag, varyingPrefix, varyingValue string) (string, bool) {
	varying := attrNameFromPrefix(varyingPrefix)
	companions, overrides := companionFor(tag, varying, varyingValue)
	add := func(pairs ...[2]string) string {
		var b strings.Builder
		for _, p := range pairs {
			if p[0] == varying {
				continue
			}
			// A companion may override a baseline pair (e.g. operator over→arithmetic).
			val := p[1]
			if ov, ok := overrides[p[0]]; ok {
				val = ov
			}
			b.WriteString(` ` + p[0] + `="` + val + `"`)
		}
		// Append companion attrs that are not already part of the baseline pairs.
		present := map[string]bool{varying: true}
		for _, p := range pairs {
			present[p[0]] = true
		}
		for _, c := range companions {
			if present[c[0]] {
				continue
			}
			b.WriteString(` ` + c[0] + `="` + c[1] + `"`)
		}
		return b.String()
	}
	switch tag {
	case "rect":
		// stroke baseline so fill="none" still shows an outline (QA round1).
		return add([2]string{"x", "10"}, [2]string{"y", "10"}, [2]string{"width", "80"}, [2]string{"height", "80"}, [2]string{"fill", "#e94560"}, [2]string{"stroke", "#16c79a"}, [2]string{"stroke-width", "2"}), false
	case "circle":
		return add([2]string{"cx", "50"}, [2]string{"cy", "50"}, [2]string{"r", "40"}, [2]string{"fill", "#16c79a"}, [2]string{"stroke", "#16c79a"}, [2]string{"stroke-width", "2"}), false
	case "ellipse":
		return add([2]string{"cx", "50"}, [2]string{"cy", "50"}, [2]string{"rx", "40"}, [2]string{"ry", "25"}, [2]string{"fill", "#f5a623"}, [2]string{"stroke", "#16c79a"}, [2]string{"stroke-width", "2"}), false
	case "line":
		return add([2]string{"x1", "10"}, [2]string{"y1", "10"}, [2]string{"x2", "90"}, [2]string{"y2", "90"}, [2]string{"stroke", "#4d8bff"}, [2]string{"stroke-width", "4"}), false
	case "polyline":
		return add([2]string{"points", "10,80 40,20 70,60 90,10"}, [2]string{"fill", "none"}, [2]string{"stroke", "#e94560"}, [2]string{"stroke-width", "3"}), false
	case "polygon":
		return add([2]string{"points", "50,10 90,80 10,80"}, [2]string{"fill", "#16c79a"}, [2]string{"stroke", "#16c79a"}, [2]string{"stroke-width", "2"}), false
	case "path":
		// path baseline is stroke-only; the explicit stroke keeps fill="none" visible.
		// A zig-zag with sharp corners and open ends so stroke-linejoin / miterlimit
		// (corners) and stroke-linecap (ends) have geometry to act on; the `d` presets
		// still replace it with their own line / curve / closed shapes.
		return add([2]string{"d", "M6 64 L28 22 L50 64 L72 22 L94 64"}, [2]string{"fill", "none"}, [2]string{"stroke", "#16c79a"}, [2]string{"stroke-width", "2"}), false
	case "text":
		// FIX 3: fill="currentColor" so a varied CSS `color` (or `fill`) recolors the
		// glyphs; the scaffold root seeds a neutral color so the base card is visible.
		return add([2]string{"x", "10"}, [2]string{"y", "55"}, [2]string{"fill", "currentColor"}, [2]string{"font-size", "20"}), false
	case "tspan":
		// NO baseline x/y: the tspan FLOWS inside its parent text run (the scaffold
		// supplies surrounding "Hi …!" text) so it reads as a styled sub-span, not a
		// detached label. When x/y ARE the varied attribute they reposition it,
		// visibly demonstrating absolute tspan positioning.
		return add([2]string{"fill", "currentColor"}, [2]string{"font-size", "22"}), false
	case "textPath":
		// textPath ignores x/y; it follows a referenced path. The blueprint defines
		// <path id="slot"> in defs, so href="#slot" makes every card's text follow it.
		// FIX 3: fill="currentColor" so a varied `color`/`fill` recolors the glyphs.
		return add([2]string{"href", "#slot"}, [2]string{"fill", "currentColor"}, [2]string{"font-size", "20"}), false
	case "image":
		// An inline data: URI so the image renders offline; omitted when href varies.
		return add([2]string{"x", "10"}, [2]string{"y", "10"}, [2]string{"width", "80"}, [2]string{"height", "80"}, [2]string{"href", imageDataURI}), false
	case "foreignObject":
		// Without an explicit content box the XHTML child is clipped to 0x0 and the
		// card is blank (QA round2). Supply a 90x90 box (minus the varied dimension).
		return add([2]string{"x", "5"}, [2]string{"y", "5"}, [2]string{"width", "90"}, [2]string{"height", "90"}), false
	case "use":
		// The blueprint's defs defines id="slot"; reference it (was "#ref" → dangling).
		// Baseline stroke inherits to the referenced shape so the stroke-* family is
		// demonstrable (the referenced shape carries no stroke of its own).
		return add([2]string{"href", "#slot"}, [2]string{"x", "10"}, [2]string{"y", "10"},
			[2]string{"stroke", "#0b3b2e"}, [2]string{"stroke-width", "3"}), false
	case "g", "switch", "symbol":
		// A baseline stroke INHERITS to the container's currentColor-filled children,
		// so stroke / stroke-width / dasharray / opacity / paint-order have a visible
		// outline to modulate instead of collapsing to identical strokeless tiles.
		return add([2]string{"stroke", "#0b3b2e"}, [2]string{"stroke-width", "3"}), false
	case "svg", "defs":
		return "", false
	case "a":
		// A baseline stroke on the <a> wrapper INHERITS to the rendered pill (whose
		// child carries no stroke of its own), so the stroke-* family (width / color /
		// dasharray / opacity) and paint-order have a visible outline to modulate
		// instead of all collapsing to the identical fill-only pill.
		return add([2]string{"stroke", "#0b3b2e"}, [2]string{"stroke-width", "3"}), false
	case "linearGradient":
		// FIX 2(b): a non-zero x1 baseline so the x2="0" card still has a non-zero
		// gradient vector (x1=0 default + x2=0 → zero-length → solid last-stop fill).
		// Skipped automatically when x1 itself is the varied attribute.
		return add([2]string{"x1", "20%"}), false
	case "radialGradient":
		return "", false
	case "stop":
		// Give the varied stop a visible color + mid offset so every card (even
		// non-color attribute paths) shows a 3-stop gradient instead of all-black.
		return add([2]string{"offset", "0.5"}, [2]string{"stop-color", "#f5a623"}), false
	case "pattern":
		return add([2]string{"width", "20"}, [2]string{"height", "20"}, [2]string{"patternUnits", "userSpaceOnUse"}), false
	case "marker":
		return add([2]string{"markerWidth", "10"}, [2]string{"markerHeight", "10"}, [2]string{"refX", "5"}, [2]string{"refY", "5"}), false
	case "feTurbulence":
		// FIX 5: baseFrequency defaults to 0 → uniform black on every non-frequency
		// card (including the x/y/width/height subregion cards). Seed a visible
		// frequency + octaves so the noise is clearly textured and the varied
		// subregion clips a recognizable patch of it.
		return add([2]string{"type", "turbulence"}, [2]string{"baseFrequency", "0.06"}, [2]string{"numOctaves", "3"}, [2]string{"seed", "5"}), false
	case "feGaussianBlur":
		return add([2]string{"stdDeviation", "3"}), false
	case "feOffset":
		return add([2]string{"dx", "4"}, [2]string{"dy", "4"}), false
	case "feFlood":
		return add([2]string{"flood-color", "#e94560"}), false
	case "feColorMatrix":
		// FIX 6: a hueRotate baseline takes a single number ("90") that is clearly
		// visible (rotates hues 90°), so every non-type/values card shows a strong,
		// well-formed effect. The previous saturate+"0.3" paired badly with the
		// type="matrix" card (matrix needs 20 numbers; "0.3" → silent identity).
		// The values override (overlay.go) returns "120" — valid for hueRotate.
		return add([2]string{"type", "hueRotate"}, [2]string{"values", "90"}), false
	case "feFuncR", "feFuncG", "feFuncB", "feFuncA":
		// type defaults to identity (ignores slope/intercept); linear makes them show.
		return add([2]string{"type", "linear"}, [2]string{"slope", "1.5"}), false
	case "feBlend":
		// The blueprint supplies feFlood result="layer2"; wire it as the 2nd input.
		return add([2]string{"in", "SourceGraphic"}, [2]string{"in2", "layer2"}), false
	case "feComposite":
		// The blueprint supplies feFlood result="layer1"/"layer2"; composite them.
		return add([2]string{"in", "layer1"}, [2]string{"in2", "layer2"}, [2]string{"operator", "over"}), false
	case "feTile":
		// The blueprint supplies a small feFlood result="patch"; tile it.
		return add([2]string{"in", "patch"}), false
	case "feMergeNode":
		// Reference a colored blueprint layer so the merge stacking is visible.
		return add([2]string{"in", "layerA"}), false
	case "feMorphology":
		// A non-zero radius + dilate so the erode/dilate effect is observable.
		return add([2]string{"operator", "dilate"}, [2]string{"radius", "3"}), false
	case "feDisplacementMap":
		// The blueprint supplies feTurbulence result="noiseMap"; use it as the map.
		return add([2]string{"in", "SourceGraphic"}, [2]string{"in2", "noiseMap"}, [2]string{"scale", "20"}), false
	case "feImage":
		// "#slot" is self-referential (the filter itself) → blank; use a data: URI.
		return add([2]string{"href", imageDataURI}), false
	case "feConvolveMatrix":
		// kernelMatrix must match order (3x3 → 9 values); a sharpen kernel.
		// in="noiseMap" feeds the textured feTurbulence result (supplied by the
		// scaffold) so the kernel has edges to operate on and the effect shows
		// instead of a flat color (QA round2).
		return add([2]string{"in", "noiseMap"}, [2]string{"order", "3"}, [2]string{"kernelMatrix", "0 -1 0 -1 5 -1 0 -1 0"}), false
	case "feDiffuseLighting":
		// camelCase SVG attribute names (match the grammar); lighting-color is a
		// hyphenated CSS presentation property. in="bumpMap" feeds the blurred
		// SourceAlpha bump (supplied by the scaffold) so surfaceScale/azimuth/
		// elevation produce visible directional 3D shading (QA round2).
		return add([2]string{"in", "bumpMap"}, [2]string{"surfaceScale", "5"}, [2]string{"diffuseConstant", "1"}, [2]string{"lighting-color", "#ffffff"}), false
	case "feSpecularLighting":
		return add([2]string{"in", "bumpMap"}, [2]string{"surfaceScale", "5"}, [2]string{"specularConstant", "1"}, [2]string{"specularExponent", "20"}, [2]string{"lighting-color", "#ffffff"}), false
	case "feDistantLight":
		// elevation/azimuth small samples read as black; 45° lights the surface.
		return add([2]string{"azimuth", "45"}, [2]string{"elevation", "45"}), false
	case "fePointLight":
		return add([2]string{"x", "50"}, [2]string{"y", "50"}, [2]string{"z", "50"}), false
	case "feSpotLight":
		// pointsAtZ shifts the cone's target depth. With a high z and the target
		// directly under the light, S_z = pointsAtZ−z barely changes over the small
		// NumberType samples (0/1/−1/3.14/0.5/2) → identical cones. Use a LOW light
		// z (10) and an OFF-AXIS target (pointsAt at 20,20 vs light at 70,70) so the
		// cone is angled and a few units of pointsAtZ visibly tilt where it lands; a
		// high surfaceScale + tight cone make the bright lobe move card-to-card.
		return add([2]string{"x", "70"}, [2]string{"y", "70"}, [2]string{"z", "10"},
			[2]string{"pointsAtX", "20"}, [2]string{"pointsAtY", "20"}, [2]string{"pointsAtZ", "0"},
			[2]string{"specularExponent", "8"}, [2]string{"limitingConeAngle", "25"}), false
	case "discard":
		// begin="1s": the host shape is present in the first frame and DISCARDED in
		// the later frames (the clock-stepped capture spans ~0–1.65s), so the GIF
		// actually demonstrates the removal instead of a static shape.
		return add([2]string{"begin", "1s"}), false
	case "animate":
		// A working animation baseline so the varied attribute is observable on a
		// live animation; the overlay (animationValueFor) types from/to to the
		// chosen attributeName ("x" here). values XOR from/to → use from/to.
		return add([2]string{"attributeName", "x"}, [2]string{"from", "10"},
			[2]string{"to", "60"}, [2]string{"dur", "2s"},
			[2]string{"repeatCount", "indefinite"}), false
	case "set":
		// <set> applies a DISCRETE value for its active interval (no tween). begin=1s
		// so the clock-stepped capture (~0–1.65s) catches the host at its original x
		// in the early frames and SNAPPED to the set value in the later ones — the
		// before→after jump that defines <set>. (from is ignored by set.)
		return add([2]string{"attributeName", "x"}, [2]string{"to", "70"},
			[2]string{"begin", "1s"}, [2]string{"dur", "2s"},
			[2]string{"repeatCount", "indefinite"}), false
	case "animateTransform":
		return add([2]string{"attributeName", "transform"}, [2]string{"type", "rotate"},
			[2]string{"from", "0 25 25"}, [2]string{"to", "360 25 25"},
			[2]string{"dur", "2s"}, [2]string{"repeatCount", "indefinite"}), false
	case "animateMotion":
		return add([2]string{"dur", "2s"}, [2]string{"repeatCount", "indefinite"},
			[2]string{"path", "M0 0 L60 60"}), false
	}
	return "", false
}

// bodyOverride returns an alternate inner body for tag when the single varied
// attribute (attrName=attrValue) changes how the CHILD geometry must be
// expressed. The motivating case is objectBoundingBox units: when a
// clipPath/mask/pattern switches its child/content coordinate system to
// objectBoundingBox, the standard userSpaceOnUse children (pixel coords like
// cx="50") map far off-canvas and the card looks blank/identical to the
// userSpaceOnUse card. Returning FRACTIONAL-coordinate children keeps the
// element visible and distinct in OBB mode. Returns ("", false) when no
// override applies and the normal bodyFor(tag) should be used.
func bodyOverride(tag, attrName, attrValue string) (string, bool) {
	if attrValue != "objectBoundingBox" {
		return "", false
	}
	switch {
	case tag == "clipPath" && attrName == "clipPathUnits":
		// Fractional circle: cx/cy/r as fractions of the clipped element's bbox.
		return `<circle cx="0.5" cy="0.5" r="0.4"/>`, true
	case tag == "mask" && attrName == "maskContentUnits":
		// Fractional white rect so the mask content covers most of the bbox.
		return `<rect x="0.1" y="0.1" width="0.8" height="0.8" fill="#fff"/>`, true
	case tag == "pattern" && attrName == "patternContentUnits":
		// A small fractional tile so the pattern content is visible in OBB mode.
		return `<circle cx="0.25" cy="0.25" r="0.2" fill="#16c79a"/>`, true
	}
	return "", false
}

// companionFor supplies DEPENDENT-COMPANION attributes for a filter primitive
// whose varied attribute is INERT unless a sibling attribute is set. It returns
// (a) companion attrs to ADD to the element's baseline and (b) overrides that
// REPLACE a baseline pair's value (keyed by attribute name). Both are keyed on
// (tag, the varied attribute name, the varied value) so the element always
// carries the context that makes the varied attribute produce a visible effect.
//
//   - feComposite varying k1/k2/k3/k4 → operator must be "arithmetic" (the k's
//     only weight the output under arithmetic mode); override the baseline
//     operator="over". When varying `operator` itself we add nothing — each
//     operator value (over/in/out/atop/xor/arithmetic) already differs.
//   - feFunc* varying amplitude/exponent/offset → type="gamma" (the gamma curve
//     C' = amplitude·Cᵉˣᵖ + offset); pair in the OTHER gamma params so the curve
//     bends visibly. Varying tableValues → type="table". Varying slope/intercept
//     → type="linear" (already the baseline). Varying `type` itself → pair each
//     type value with the param that makes that type visible.
//   - feConvolveMatrix varying kernelMatrix → order="3" (the kernels are 9-value;
//     baselineFor already sets order="3", distinctValueSet supplies the kernels).
func companionFor(tag, varied, value string) (companions [][2]string, overrides map[string]string) {
	overrides = map[string]string{}
	switch tag {
	case "feColorMatrix":
		// FIX 4: the baseline is type="hueRotate" values="90" (a valid 1-number
		// rotation). When the enumerated `type` value is "matrix", that 1 number is
		// invalid (matrix needs 20) → near-identity, so override `values` with a
		// VALID 20-number channel-swap matrix that visibly transforms color. The
		// other type values (saturate/hueRotate) take a 1-number `values`, so keep
		// the baseline "90" for them. When `values` itself is the varied attribute,
		// baselineFor pins type="hueRotate" (a 1-number type) so the values cards
		// stay valid (overlay.go returns "120").
		if varied == "type" && value == "matrix" {
			// channel-swap: R↔B (and keep G, A) — clearly distinct from hueRotate.
			overrides["values"] = "0 0 1 0 0  0 1 0 0 0  1 0 0 0 0  0 0 0 1 0"
		}
		// The arithmetic coefficients k1-k4 are ignored in every mode but
		// arithmetic, so flip the baseline operator when one of them is varied.
		switch varied {
		case "k1", "k2", "k3", "k4":
			overrides["operator"] = "arithmetic"
			// Seed the OTHER coefficients so the varied k is not the lone non-zero
			// term (k1·i1·i2 + k2·i1 + k3·i2 + k4). A mid blend makes each k visible.
			seed := map[string]string{"k1": "0.5", "k2": "0.5", "k3": "0.5", "k4": "0"}
			delete(seed, varied)
			for _, k := range []string{"k1", "k2", "k3", "k4"} {
				if v, ok := seed[k]; ok {
					companions = append(companions, [2]string{k, v})
				}
			}
		}
	case "feFuncR", "feFuncG", "feFuncB", "feFuncA":
		switch varied {
		case "amplitude", "exponent", "offset":
			// gamma transfer: C' = amplitude·Cᵉˣᵖᵒⁿᵉⁿᵗ + offset.
			overrides["type"] = "gamma"
			seed := map[string]string{"amplitude": "1", "exponent": "3", "offset": "0"}
			delete(seed, varied)
			for _, k := range []string{"amplitude", "exponent", "offset"} {
				if v, ok := seed[k]; ok {
					companions = append(companions, [2]string{k, v})
				}
			}
		case "tableValues":
			overrides["type"] = "table"
		case "type":
			// Pair each enumerated type value with the param that makes it visible:
			// table/discrete need tableValues; gamma needs exponent; linear needs
			// slope (already in the baseline, kept as-is). identity needs nothing.
			switch value {
			case "table", "discrete":
				companions = append(companions, [2]string{"tableValues", "0 0.3 0.6 1"})
			case "gamma":
				companions = append(companions, [2]string{"amplitude", "1"}, [2]string{"exponent", "3"})
			case "linear":
				companions = append(companions, [2]string{"slope", "1.5"}, [2]string{"intercept", "0"})
			}
		}
	}
	return companions, overrides
}

// distinctValueSet returns a hand-picked list of VALUES for a (tag, attr) whose
// shared grammar leaf samples all collapse to one render. Returns nil when the
// normal enumeration applies. Consulted by enumerateValue before the generic
// leaf/overlay sampling.
//
// feConvolveMatrix kernelMatrix: the convolution is only applied when the kernel
// length equals order² (=9 for order=3, set in baselineFor); the shared
// ListOfNumbersType reps (1 2 3 4, 0 1 0 0, …) never have 9 elements, so every
// kernel is invalid → identity → identical cards. Supply five 9-value 3×3
// kernels with visibly different effects.
func distinctValueSet(tag, attrName string) []string {
	// FIX 1: ellipse rx/ry. The grammar's ( "auto" | LengthPercentageType ) admits
	// rx="auto"/ry="auto", but Chrome does not implement SVG2 `auto` ellipse radii
	// (it treats them as 0 → the ellipse is invisible). Restrict the showcase to
	// CONCRETE radii so every card paints a real ellipse shape. (rect rx/ry keep
	// "auto" — there it just means "no corner rounding", which renders fine.)
	if tag == "ellipse" && (attrName == "rx" || attrName == "ry") {
		return []string{"30", "45", "20%", "60"}
	}
	// FIX 2: foreignObject width/height. `auto` is invalid for foreignObject
	// dimensions (the HTML box collapses to 0×0 → blank card). Use concrete
	// LengthPercentageType values only so the embedded HTML box always has size.
	if tag == "foreignObject" && (attrName == "width" || attrName == "height") {
		return []string{"30", "60", "50%", "90"}
	}
	if tag == "feConvolveMatrix" && attrName == "kernelMatrix" {
		return []string{
			"0 0 0 0 1 0 0 0 0",     // identity — passthrough
			"0 -1 0 -1 4 -1 0 -1 0", // edge detect (Laplacian)
			"-2 -1 0 -1 1 1 0 1 2",  // emboss
			"1 1 1 1 1 1 1 1 1",     // box blur
			"0 -1 0 -1 5 -1 0 -1 0", // sharpen
		}
	}
	if tag == "feSpotLight" && attrName == "pointsAtZ" {
		// The shared NumberType samples (0/1/−1/3.14/0.5/2) are all tiny relative to
		// the light position, so the cone target barely moves → identical cones.
		// Sweep pointsAtZ across a wide depth range (in front of, on, and behind the
		// surface plane) so the cone axis tilt and the lit lobe shift card-to-card.
		return []string{"-40", "-10", "0", "10", "40", "80"}
	}
	return nil
}

// bodyFor returns the child content a freshly-generated element needs to be
// visible (e.g. a gradient/pattern needs stops/children; a marker/clipPath needs
// a shape). For self-rendering shapes this is empty.
func bodyFor(tag string) string {
	switch tag {
	case "linearGradient", "radialGradient":
		return `<stop offset="0" stop-color="#e94560"/><stop offset="1" stop-color="#16c79a"/>`
	case "pattern":
		// A STRUCTURED, asymmetric tile (background + diagonal + dot) so that
		// patternTransform (rotate/scale), patternUnits and viewBox visibly change the
		// tiling — a single symmetric dot looks identical under every transform.
		return `<rect width="20" height="20" fill="#11201c"/>` +
			`<path d="M0 20 L20 0" stroke="#4ee39a" stroke-width="2.5"/>` +
			`<circle cx="6" cy="6" r="3.5" fill="#f5a623"/>`
	case "marker":
		return `<path d="M0 0 L10 5 L0 10 Z" fill="#e94560"/>`
	case "clipPath":
		// A STAR clip so the clipped host reads unmistakably as "clipped to a shape"
		// (a circle clip is hard to tell from a plain circle).
		return `<polygon points="50,6 61,38 95,38 67,58 78,92 50,72 22,92 33,58 5,38 39,38"/>`
	case "mask":
		// Concentric luminance bands (black bg → white → grey → white) so the mask
		// REVEALS the host at varying opacity (a bullseye) — showing luminance
		// masking, not just a hard rectangular cut-out. mask-type=luminance vs alpha
		// then differ visibly.
		return `<rect width="100" height="100" fill="#000"/>` +
			`<circle cx="50" cy="50" r="44" fill="#fff"/>` +
			`<circle cx="50" cy="50" r="30" fill="#888"/>` +
			`<circle cx="50" cy="50" r="15" fill="#fff"/>`
	case "filter":
		return `<feGaussianBlur stdDeviation="3"/>`
	case "g":
		// MULTIPLE distinct children so grouping is self-evident: a transform /
		// opacity / color set on the <g> propagates to all three shapes at once
		// (they inherit paint via fill="currentColor"). One square cannot show what
		// a group does.
		return `<circle cx="32" cy="38" r="18" fill="currentColor"/>` +
			`<rect x="50" y="20" width="34" height="34" rx="6" fill="currentColor"/>` +
			`<polygon points="50,62 80,90 20,90" fill="currentColor"/>`
	case "a":
		// A link-like "button": a pill the wrapper recolours (currentColor) plus a
		// dark caption, so the <a> reads as a clickable hyperlink, not a bare square.
		return `<rect x="12" y="36" width="76" height="28" rx="14" fill="currentColor"/>` +
			`<text x="50" y="55" text-anchor="middle" font-size="13" font-family="monospace" fill="#0c100e" stroke="none">link</text>`
	case "svg":
		// A NESTED viewport: a bordered box establishing its own coordinate system,
		// with off-centre content the viewport crops — so viewBox / preserveAspectRatio
		// visibly reframe it instead of looking like a plain square.
		return `<rect x="3" y="3" width="94" height="94" fill="#0e1a17" stroke="currentColor" stroke-width="3"/>` +
			`<circle cx="74" cy="74" r="36" fill="currentColor"/>`
	case "symbol":
		// Off-centre icon content (a circle + square) inside the symbol so that, once
		// instantiated via <use>, its viewBox / preserveAspectRatio / refX / refY
		// visibly reframe and anchor it.
		return `<circle cx="34" cy="34" r="22" fill="currentColor"/>` +
			`<rect x="52" y="52" width="34" height="34" rx="4" fill="currentColor"/>`
	case "defs":
		// <defs> never renders directly; its <use href="#defskid"> instantiates a
		// fixed rect (the defs child can't inherit paint set on <defs>). Keep a
		// concrete fill so the card is visible.
		return `<rect x="20" y="20" width="60" height="60" fill="#4d8bff"/>`
	case "switch":
		// <switch> renders the FIRST child whose conditional-processing attrs pass; a
		// bare child always qualifies. Wrapping in a labelled <g> (rendered as one
		// child) makes the "selected child" reading clear instead of an anonymous
		// square. fill="currentColor" still lets the <switch> parent's paint propagate.
		return `<g><rect x="18" y="22" width="64" height="56" rx="8" fill="currentColor"/>` +
			`<text x="50" y="55" text-anchor="middle" font-size="10" font-family="monospace" fill="#0c100e">switch</text></g>`
	case "feMerge":
		// Two nodes stacking colored blueprint layers so the merge is visible.
		return `<feMergeNode in="layerA"/><feMergeNode in="layerB"/>`
	case "feComponentTransfer":
		return `<feFuncR type="linear" slope="1.5"/>`
	case "feDiffuseLighting", "feSpecularLighting":
		// A distant light at 45°/45° reliably lights the surface (point/spot lights
		// need a large z which the bare scaffold can't guarantee).
		return `<feDistantLight azimuth="45" elevation="45"/>`
	case "foreignObject":
		// foreignObject needs an XHTML child or it renders nothing. The div fills
		// the foreignObject box (width:100%;height:100%) so that varying the
		// foreignObject's width OR height visibly resizes the painted area — a
		// content-sized div would leave every height looking identical (the one
		// text line is shorter than any tested height).
		return `<div xmlns="http://www.w3.org/1999/xhtml" style="box-sizing:border-box;width:100%;height:100%;background:#4d8bff;color:#fff;font:12px sans-serif;padding:6px">HTML in SVG</div>`
	case "script":
		// A JS body that recolors the blueprint's target rect (id="slot-target") to
		// teal. The rect is PRE-COLORED red in the scaffold (FIX 9) so non-JS-MIME
		// cards still show a colored square; an executing card additionally shifts it
		// to teal, demonstrating the script ran.
		return `(function(){var t=document.getElementById('slot-target');if(t){t.setAttribute('fill','#16c79a');}})();`
	case "style":
		// A CSS body that styles the blueprint's .slot shapes so cards differ.
		return `.slot{fill:#e94560;stroke:#16c79a;stroke-width:3}`
	case "text":
		return `Text` // a real word so x/y/dx/dy/rotate/textLength visibly deform it
	case "textPath":
		return `Following` // long enough to trace and reveal the curved path
	case "tspan":
		return `Ag` // the scaffold wraps this in a surrounding parent text run
	case "desc", "title", "metadata":
		return `info`
	}
	return ""
}
