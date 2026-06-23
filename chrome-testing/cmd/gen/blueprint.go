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
	"use":            true,
	"switch":         true,
	"script":         true,
	"mpath":          true,
	"discard":        true,
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
		return svgOpen +
			`<defs>` +
			`<linearGradient id="bpGrad" x1="0" y1="0" x2="1" y2="0">` +
			`<stop offset="0" stop-color="#ff0000"/><stop offset="0.5" stop-color="#00ff00"/><stop offset="1" stop-color="#0000ff"/>` +
			`</linearGradient>` +
			slot +
			`</defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="url(#bpGrad)" filter="url(#slot)"/></svg>`
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
		return `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 100 100" color="#16c79a">` +
			`{{ELEMENT}}</svg>`
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
		return svgOpen +
			`<defs><circle id="slot" cx="50" cy="50" r="32" fill="#e94560"/></defs>{{ELEMENT}}</svg>`
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
func baselineFor(tag, varyingPrefix string) (string, bool) {
	varying := attrNameFromPrefix(varyingPrefix)
	add := func(pairs ...[2]string) string {
		var b strings.Builder
		for _, p := range pairs {
			if p[0] == varying {
				continue
			}
			b.WriteString(` ` + p[0] + `="` + p[1] + `"`)
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
		return add([2]string{"d", "M10 50 Q50 10 90 50 T90 90"}, [2]string{"fill", "none"}, [2]string{"stroke", "#16c79a"}, [2]string{"stroke-width", "2"}), false
	case "text", "tspan":
		return add([2]string{"x", "10"}, [2]string{"y", "55"}, [2]string{"fill", "#e6e6e6"}, [2]string{"font-size", "20"}), false
	case "textPath":
		// textPath ignores x/y; it follows a referenced path. The blueprint defines
		// <path id="slot"> in defs, so href="#slot" makes every card's text follow it.
		return add([2]string{"href", "#slot"}, [2]string{"fill", "#e6e6e6"}, [2]string{"font-size", "20"}), false
	case "image":
		// An inline data: URI so the image renders offline; omitted when href varies.
		return add([2]string{"x", "10"}, [2]string{"y", "10"}, [2]string{"width", "80"}, [2]string{"height", "80"}, [2]string{"href", imageDataURI}), false
	case "foreignObject":
		// Without an explicit content box the XHTML child is clipped to 0x0 and the
		// card is blank (QA round2). Supply a 90x90 box (minus the varied dimension).
		return add([2]string{"x", "5"}, [2]string{"y", "5"}, [2]string{"width", "90"}, [2]string{"height", "90"}), false
	case "use":
		// The blueprint's defs defines id="slot"; reference it (was "#ref" → dangling).
		return add([2]string{"href", "#slot"}, [2]string{"x", "10"}, [2]string{"y", "10"}), false
	case "g", "svg", "defs", "switch", "a", "symbol":
		return "", false
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
		return add([2]string{"x", "50"}, [2]string{"y", "50"}, [2]string{"z", "50"},
			[2]string{"pointsAtX", "50"}, [2]string{"pointsAtY", "50"}, [2]string{"pointsAtZ", "0"},
			[2]string{"specularExponent", "5"}, [2]string{"limitingConeAngle", "45"}), false
	case "discard":
		// begin="60s" so discard does not remove the host shape at t=0.
		return add([2]string{"begin", "60s"}), false
	case "animate", "set":
		// A working animation baseline so the varied attribute is observable on a
		// live animation; the overlay (animationValueFor) types from/to to the
		// chosen attributeName ("x" here). values XOR from/to → use from/to.
		return add([2]string{"attributeName", "x"}, [2]string{"from", "10"},
			[2]string{"to", "60"}, [2]string{"dur", "2s"},
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

// bodyFor returns the child content a freshly-generated element needs to be
// visible (e.g. a gradient/pattern needs stops/children; a marker/clipPath needs
// a shape). For self-rendering shapes this is empty.
func bodyFor(tag string) string {
	switch tag {
	case "linearGradient", "radialGradient":
		return `<stop offset="0" stop-color="#e94560"/><stop offset="1" stop-color="#16c79a"/>`
	case "pattern":
		return `<circle cx="10" cy="10" r="6" fill="#f5a623"/>`
	case "marker":
		return `<path d="M0 0 L10 5 L0 10 Z" fill="#e94560"/>`
	case "clipPath":
		return `<circle cx="50" cy="50" r="35"/>`
	case "mask":
		return `<rect x="20" y="20" width="60" height="60" fill="#fff"/>`
	case "filter":
		return `<feGaussianBlur stdDeviation="3"/>`
	case "g", "a", "svg", "defs", "symbol":
		return `<rect x="20" y="20" width="60" height="60" fill="#4d8bff"/>`
	case "switch":
		// FIX 8: the rendered child uses fill="currentColor" so presentation attrs
		// set on the <switch> parent (color/fill-opacity/opacity) actually propagate
		// and produce visibly distinct cards. The first child whose conditional-
		// processing attrs evaluate true renders; a bare rect always qualifies.
		return `<rect x="20" y="20" width="60" height="60" fill="currentColor"/>`
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
		// foreignObject needs an XHTML child or it renders nothing.
		return `<div xmlns="http://www.w3.org/1999/xhtml" style="background:#4d8bff;color:#fff;font:12px sans-serif;padding:6px">HTML in SVG</div>`
	case "script":
		// A JS body that recolors the blueprint's target rect (id="slot-target") to
		// teal. The rect is PRE-COLORED red in the scaffold (FIX 9) so non-JS-MIME
		// cards still show a colored square; an executing card additionally shifts it
		// to teal, demonstrating the script ran.
		return `(function(){var t=document.getElementById('slot-target');if(t){t.setAttribute('fill','#16c79a');}})();`
	case "style":
		// A CSS body that styles the blueprint's .slot shapes so cards differ.
		return `.slot{fill:#e94560;stroke:#16c79a;stroke-width:3}`
	case "text", "tspan", "textPath":
		return `Ag`
	case "desc", "title", "metadata":
		return `info`
	}
	return ""
}
