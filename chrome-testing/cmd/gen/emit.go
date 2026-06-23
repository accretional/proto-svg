package main

import "strings"

// emit.go — the element×attribute APPLICABILITY classifier. nonVisualAttr reports
// whether an attribute has any rendering effect on a given element; the catalogue
// uses it to keep only the visually-meaningful attributes as controls/presets.

// page bundles an element's enumerated variants for emission.
type page struct {
	tag      string
	variants []Variant
}

// --- tag-aware non-visual grouping (QA round2) ------------------------------
//
// An attribute is "non-visual" for a given element when setting it produces no
// rendering difference on THAT element. Three sources of non-visual-ness:
//
//  1. Universally non-visual attrs (id, aria-*, on* events, conditional-
//     processing flags, …) — never affect rendering on any element.
//  2. Navigation-semantics attrs on <a> (download, rel, type, …).
//  3. ANY presentation attribute (the full SVG presentation-attribute set:
//     paint/stroke/marker/render hints AND text/font/clip/mask/display/transform/
//     pointer-events/cursor/opacity/visibility/overflow/…) that the grammar
//     over-approximates onto elements which do not paint geometry directly
//     (filter primitives, gradients, stops, defs, masks, animation elements, …).
//     They paint nothing, so these attrs have no effect — UNLESS the element has
//     its own element-specific paint attr (feFlood→flood-color, stop→stop-color,
//     lighting-color on lighting primitives, color-interpolation-filters on fe*).

// universalNonVisual is the set of attribute names with no visual effect on ANY
// element (semantics / accessibility / scripting / conditional-processing).
var universalNonVisual = map[string]bool{
	"id": true, "tabindex": true, "autofocus": true,
	"lang": true, "xml:lang": true, "xml:space": true,
	"class": true, "style": true, "role": true,
	// conditional-processing flags — selection, not rendering (no visual diff).
	"requiredExtensions": true, "systemLanguage": true,
}

// smilTimingNonVisual is the set of SMIL timing/control attributes (animation
// elements only). They modulate WHEN an animation runs, HOW LONG, and whether it
// repeats/freezes — never the appearance of any single frame, so a per-value
// specimen capture cannot tell them apart (every card already shows the host
// animating). They render valid values by construction (the clock grammar is
// unsigned; repeatCount is non-negative) but carry no distinct still/GIF, so they
// belong in the collapsed Non-visual section. `restart`/`fill`(freeze|remove)/
// `additive`/`accumulate`/`attributeType` are likewise control-only.
var smilTimingNonVisual = map[string]bool{
	"dur": true, "begin": true, "end": true, "min": true, "max": true,
	"repeatDur": true, "repeatCount": true, "restart": true,
	"additive": true, "accumulate": true, "attributeType": true,
	// `fill` on an animation element is the fill-MODE (freeze|remove): it controls
	// what happens AFTER the animation ends, not any single frame — non-visual.
	// (On a shape, `fill` is paint and stays visual; this is animation-gated.)
	"fill": true,
}

// aNonVisual is the set of <a>-specific navigation-semantics attrs (no visual
// rendering effect — they only control link behavior). href/target (and their
// xlink: aliases) merely set the link destination/framing; like download/ping/
// rel/hreflang/type/referrerpolicy they never change how the <a> renders.
var aNonVisual = map[string]bool{
	"href": true, "target": true,
	"download": true, "ping": true, "rel": true,
	"hreflang": true, "type": true, "referrerpolicy": true,
	"xlink:href": true, "xlink:title": true,
}

// presentationAttrs is the COMPLETE SVG presentation-attribute set — every name
// in lang/styling.ebnf's `PresentationAttribute` union. A presentation attr only
// has a visual effect on elements that paint geometry directly (paintingGraphics);
// on every other element (filter primitives, gradients, stops, pattern, masks,
// animation elements, …) the grammar over-approximates them and they have no
// rendering effect, so they belong in the collapsed non-visual section.
//
// Kept in sync with the union by name: painting (fill/stroke/marker/color/render
// hints/vector-effect/filter), filter (color-interpolation-filters/flood-*/
// lighting-color/image-rendering), masking (clip-*/clip/mask), text (anchor/
// baselines/direction/bidi/writing-mode/spacing/decoration/overflow/render/
// white-space/font-*/inline-size/glyph-orientation), embedded (pointer-events),
// pservers (stop-color/stop-opacity), styling (display/visibility/overflow/
// opacity/cursor/transform/transform-origin). The element-specific exceptions
// (stop-color, flood-*, lighting-color, color-interpolation-filters) are routed
// back into the main grid by keepPaintAttrForTag below.
var presentationAttrs = map[string]bool{
	// painting.ebnf
	"fill": true, "fill-rule": true, "fill-opacity": true,
	"stroke": true, "stroke-opacity": true, "stroke-width": true,
	"stroke-linecap": true, "stroke-linejoin": true, "stroke-miterlimit": true,
	"stroke-dasharray": true, "stroke-dashoffset": true,
	"paint-order": true,
	"marker": true, "marker-start": true, "marker-mid": true, "marker-end": true,
	"color": true, "color-interpolation": true, "color-rendering": true,
	"shape-rendering": true, "vector-effect": true,
	// filter.ebnf
	"filter": true, "color-interpolation-filters": true,
	"flood-color": true, "flood-opacity": true, "lighting-color": true,
	"image-rendering": true,
	// masking.ebnf
	"clip-path": true, "clip-rule": true, "clip": true, "mask": true,
	// text.ebnf
	"text-anchor": true, "dominant-baseline": true, "alignment-baseline": true,
	"baseline-shift": true, "direction": true, "unicode-bidi": true,
	"writing-mode": true, "letter-spacing": true, "word-spacing": true,
	"text-decoration": true, "text-overflow": true, "text-rendering": true,
	"white-space": true,
	"font": true, "font-family": true, "font-size": true, "font-size-adjust": true,
	"font-style": true, "font-variant": true, "font-weight": true, "font-stretch": true,
	"inline-size": true, "glyph-orientation-vertical": true, "kerning": true,
	// embedded.ebnf
	"pointer-events": true,
	// pservers.ebnf
	"stop-color": true, "stop-opacity": true,
	// styling.ebnf
	"display": true, "visibility": true, "overflow": true, "opacity": true,
	"cursor": true, "transform": true, "transform-origin": true,
}

// paintingGraphics is the set of elements that paint geometry directly (or are
// containers whose paint presentation attrs inherit to painted children). For
// these, the SHARED paint presentation attrs (fill/stroke/color/…) STAY in the
// main grid. This is the broadest paint-applicability set; narrower attribute
// categories (text, marker, render hints, …) refine it in appliesTo below.
var paintingGraphics = map[string]bool{
	"rect": true, "circle": true, "ellipse": true, "line": true,
	"polyline": true, "polygon": true, "path": true,
	"text": true, "tspan": true, "textPath": true,
	"use": true, "g": true, "svg": true, "a": true, "switch": true,
	"foreignObject": true, "image": true, "symbol": true,
}

// --- element×attribute APPLICABILITY (FIX 1) --------------------------------
//
// A presentation attribute only has a UNIQUE rendered effect on an element when
// the SVG "Applies to" rule says it applies to that element. Capturing a
// text/font prop on a <rect>, a marker-* on a <circle>, or a paint prop on an
// element that paints nothing yields identical cards (and pointless specimens).
//
// The attribute-CATEGORY sets below name each presentation attr by the kind of
// element it visually affects; the element-CATEGORY sets name the elements in
// each kind. appliesTo(attr, tag) is the membership test: it returns false (→
// collapsed / no specimen) when a presentation attr is applied to an element it
// does not visually affect.

// shapeTags are the basic shapes that paint geometry from vertices/outlines.
var shapeTags = map[string]bool{
	"rect": true, "circle": true, "ellipse": true, "line": true,
	"polyline": true, "polygon": true, "path": true,
}

// vertexShapeTags are the shapes that have vertices markers can be drawn at
// (marker-* "Applies to": path, line, polyline, polygon).
var vertexShapeTags = map[string]bool{
	"path": true, "line": true, "polyline": true, "polygon": true,
}

// textTags are the text content elements (font/text props "Applies to").
var textTags = map[string]bool{
	"text": true, "tspan": true, "textPath": true,
}

// paintContainerTags are containers that paint via their children, so inherited
// paint props (fill/stroke/color/…) set on them propagate down and DO show.
var paintContainerTags = map[string]bool{
	"use": true, "g": true, "svg": true, "a": true, "switch": true,
	"symbol": true, "marker": true, "pattern": true,
}

// shapeRenderingTags: shape-rendering "Applies to" shapes + use/g/svg.
var shapeRenderingTags = map[string]bool{
	"rect": true, "circle": true, "ellipse": true, "line": true,
	"polyline": true, "polygon": true, "path": true,
	"use": true, "g": true, "svg": true,
}

// fillRuleTags: fill-rule "Applies to" path, polygon, polyline, clipPath (+ text
// content). Narrower than the general paint set.
var fillRuleTags = map[string]bool{
	"path": true, "polygon": true, "polyline": true, "clipPath": true,
	"text": true, "tspan": true, "textPath": true,
}

// textProps are the font/text presentation attrs — ONLY text content elements.
var textProps = map[string]bool{
	"font": true, "font-family": true, "font-size": true, "font-size-adjust": true,
	"font-style": true, "font-variant": true, "font-weight": true, "font-stretch": true,
	"text-anchor": true, "dominant-baseline": true, "alignment-baseline": true,
	"baseline-shift": true, "direction": true, "unicode-bidi": true,
	"writing-mode": true, "letter-spacing": true, "word-spacing": true,
	"text-decoration": true, "text-overflow": true, "text-rendering": true,
	"white-space": true, "inline-size": true, "glyph-orientation-vertical": true,
	"kerning": true,
}

// markerProps are the marker-* attrs — ONLY vertex-bearing shapes.
var markerProps = map[string]bool{
	"marker": true, "marker-start": true, "marker-mid": true, "marker-end": true,
}

// generalPaintProps are the shared paint/stroke attrs — shapes + text content +
// paint containers (use/g/svg/a/switch/symbol/marker/pattern).
var generalPaintProps = map[string]bool{
	"fill": true, "fill-opacity": true,
	"stroke": true, "stroke-opacity": true, "stroke-width": true,
	"stroke-linecap": true, "stroke-linejoin": true, "stroke-miterlimit": true,
	"stroke-dasharray": true, "stroke-dashoffset": true,
	"paint-order": true, "color": true, "color-interpolation": true,
}

// broadGraphicsProps apply to most graphics; they STAY in the main grid for any
// painting graphic (clip/mask/filter/opacity/visibility/display/transform/…).
// They are NOT routed through the narrow per-category tests above.
var broadGraphicsProps = map[string]bool{
	"clip-path": true, "clip-rule": true, "clip": true, "mask": true,
	"filter": true, "opacity": true, "visibility": true, "display": true,
	"transform": true, "transform-origin": true, "cursor": true,
	"pointer-events": true, "overflow": true, "color-rendering": true,
	"vector-effect": true,
}

// elementSpecificPaintProps are paint attrs that ONLY ever apply to a small,
// fixed set of owning elements (filter primitives / stops). On every other
// element they are inert, so they must NOT fall through to the broad painting-
// graphic test below. keepPaintAttrForTag decides the owning element(s).
var elementSpecificPaintProps = map[string]bool{
	"flood-color": true, "flood-opacity": true, "lighting-color": true,
	"color-interpolation-filters": true, "stop-color": true, "stop-opacity": true,
}

// staticallyInertProps are attributes that, in a still SVG showcase, produce no
// distinct visual display: rendering/colour hints (no visible diff), interaction
// (pointer-events), clip-rule (only meaningful inside a clipPath, not on a stand-
// alone shape), vector-effect (non-scaling-stroke only manifests under a
// transform), and pathLength (only affects dash patterns). They are valid and
// stay in the grammar + the collapsed "Non-visual" gallery section, but they get
// no per-value specimen (they could never render a unique display).
//
// `result` and `kernelUnitLength` are filter-primitive WIRING attrs with no
// standalone display: `result` only NAMES the primitive's output (every value
// renders identically because nothing consumes the named output in a single-
// primitive specimen), and `kernelUnitLength` is a deprecated/no-op hint that
// browsers ignore. Both are valid grammar but never a unique still render.
var staticallyInertProps = map[string]bool{
	"pointer-events": true, "color-rendering": true, "color-interpolation": true,
	"clip-rule": true, "vector-effect": true, "pathLength": true,
	"result": true, "kernelUnitLength": true,
}

// appliesTo reports whether presentation attr has a UNIQUE visual effect on
// element tag per the SVG "Applies to" rules. Only consulted for names in
// presentationAttrs; non-presentation attrs are handled separately.
func appliesTo(attr, tag string) bool {
	switch {
	// Text/font props: ONLY text content elements.
	case textProps[attr]:
		return textTags[tag]
	// marker-*: ONLY vertex-bearing shapes.
	case markerProps[attr]:
		return vertexShapeTags[tag]
	// shape-rendering: shapes + use/g/svg.
	case attr == "shape-rendering":
		return shapeRenderingTags[tag]
	// image-rendering: image (+ feImage).
	case attr == "image-rendering":
		return tag == "image" || tag == "feImage"
	// fill-rule: path/polygon/polyline/clipPath + text content.
	case attr == "fill-rule":
		return fillRuleTags[tag]
	// Element-specific paint attrs (flood-*/lighting-color/stop-*/color-
	// interpolation-filters): apply ONLY on their owning element(s), nowhere else.
	case elementSpecificPaintProps[attr]:
		return keepPaintAttrForTag(attr, tag)
	// General paint/stroke props: shapes + text content + paint containers.
	case generalPaintProps[attr]:
		return shapeTags[tag] || textTags[tag] || paintContainerTags[tag]
	// Broad graphics props (clip/mask/filter/opacity/transform/…): any painting
	// graphic.
	case broadGraphicsProps[attr]:
		return paintingGraphics[tag]
	}
	// Any other presentation attr: fall back to the broad painting-graphic test.
	return paintingGraphics[tag]
}

// keepPaintAttrForTag handles element-specific paint attrs that DO have an
// effect on otherwise-non-painting elements. It returns true when attr must
// stay in the main grid for tag even though tag is not a painting graphic.
func keepPaintAttrForTag(attr, tag string) bool {
	switch attr {
	case "flood-color", "flood-opacity":
		return tag == "feFlood" || tag == "feDropShadow"
	case "lighting-color":
		return tag == "feDiffuseLighting" || tag == "feSpecularLighting"
	case "color-interpolation-filters":
		return isFilterPrimitive(tag)
	case "stop-color", "stop-opacity":
		return tag == "stop"
	}
	return false
}

// isFilterPrimitive reports whether tag is an fe* filter primitive element.
func isFilterPrimitive(tag string) bool {
	return strings.HasPrefix(tag, "fe")
}

// nonVisualAttr reports whether attr has no visual rendering effect on element
// tag (it only affects semantics/accessibility/scripting, or it is a
// presentation attr applied to an element the SVG "Applies to" rule says it does
// not affect). These flood galleries with identical cards, so emitPage groups
// them into the collapsed "Non-visual attributes" section AND specimen.go skips
// them. Tag-aware element×attribute applicability (FIX 1, QA round2/round4).
func nonVisualAttr(attr, tag string) bool {
	// 1. Universal non-visual attrs and prefixed families.
	if universalNonVisual[attr] {
		return true
	}
	if strings.HasPrefix(attr, "data-") ||
		strings.HasPrefix(attr, "aria-") ||
		strings.HasPrefix(attr, "on") {
		return true
	}
	// 1b. Attributes with no distinct STATIC display (hints / interaction / dash-
	//     only / transform-only) — valid but never a unique still render.
	if staticallyInertProps[attr] {
		return true
	}
	// 1b-ii. SMIL timing/control on animation elements — modulates timing, not
	//        appearance; no distinct per-value capture.
	if smilTimingNonVisual[attr] && isAnimationTag(tag) {
		return true
	}
	// 1c. Tag-specific loading / scripting / interaction attrs with no visual
	//     display (a MIME type, a media query, a CORS/decoding/async hint, the
	//     interactive zoomAndPan) — they never change the static render.
	switch tag {
	case "script":
		if attr == "type" || attr == "crossorigin" || attr == "async" || attr == "defer" {
			return true
		}
		// FIX 7: <script> renders nothing itself, so its reference/metadata attrs
		// (href/xlink:href point at the script SOURCE; xlink:title is metadata) have
		// no visual display — the specimens were identical blank cards. Mark them
		// non-visual so <script> gets 0 specimens (like <style>).
		if attr == "href" || attr == "xlink:href" || attr == "xlink:title" {
			return true
		}
	case "style":
		if attr == "type" || attr == "media" || attr == "title" {
			return true
		}
	case "image":
		if attr == "crossorigin" || attr == "decoding" {
			return true
		}
	case "feImage":
		if attr == "crossorigin" {
			return true
		}
	case "svg", "view":
		if attr == "zoomAndPan" {
			return true
		}
	}
	// 2. <a>-specific navigation-semantics attrs (href/target/download/ping/rel/
	//    hreflang/type/referrerpolicy) — pure link behavior, no rendering effect.
	if tag == "a" && aNonVisual[attr] {
		return true
	}
	// 3. Presentation attrs: non-visual unless the SVG "Applies to" rule says they
	//    visually affect THIS element (appliesTo). This collapses font/text props
	//    on a rect, marker-* on a circle, paint props on non-painting elements,
	//    etc. — anything that would otherwise render an identical duplicate card.
	if presentationAttrs[attr] {
		return !appliesTo(attr, tag)
	}
	return false
}

