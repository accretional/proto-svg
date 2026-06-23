package main

import (
	"encoding/json"
	"fmt"
	"html"
	"sort"
	"strings"
)

// emit.go — write the generated galleries. Per element a dark-theme showcase
// page (one card per enumerated value-path: blueprint-wrapped SVG + a monospace
// attr="value" label), an index.html linking them, values.json
// (element→[{attr,value}]) and manifest.tsv.

// page bundles an element's enumerated variants for emission.
type page struct {
	tag      string
	variants []Variant
}

// galleryCSS is the showcase CSS from TEMPLATE_GUIDE (dark theme).
const galleryCSS = `body{margin:0;background:#1a1a2e;color:#e6e6e6;font:14px/1.4 ui-monospace,Menlo,monospace;padding:24px}
h1{color:#16c79a;font-size:18px;margin:0 0 4px}
p.desc{color:#9aa;margin:0 0 20px}
a{color:#4d8bff}
.grid{display:flex;flex-wrap:wrap;gap:16px}
.card{background:#0f1530;border:1px solid #26305a;border-radius:8px;padding:10px;width:160px;min-width:160px;min-height:200px;box-sizing:border-box;display:flex;flex-direction:column}
.card svg{display:block;background:#161c3a;border-radius:4px;width:140px;height:140px}
.card .label{margin-top:8px;color:#f5a623;font-size:12px;word-break:break-word;overflow-wrap:anywhere}
.card .attr{color:#16c79a}
details.meta{margin-top:28px}
details.meta>summary{color:#9aa;font-size:13px;cursor:pointer;margin-bottom:12px}
.card--meta{opacity:.5;border-style:dashed}`

// emitPage renders one element's gallery HTML. Visually-meaningful cards come
// first in the main grid; no-visual-effect metadata attributes (id, aria-*, on*
// events, etc.) are exhaustively kept but moved into a dimmed, collapsed
// <details> section after them so they don't flood the grid (QA round1).
func emitPage(p page) string {
	var b strings.Builder
	t := html.EscapeString("<" + p.tag + ">")

	var visual, meta []Variant
	for _, v := range p.variants {
		if nonVisualAttr(v.Attr, p.tag) {
			meta = append(meta, v)
		} else {
			visual = append(visual, v)
		}
	}

	fmt.Fprintf(&b, `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>%s</title>
<style>%s</style></head><body>
<h1>%s</h1>
<p class="desc">%d enumerated value-paths walked from the grammar. <a href="index.html">&larr; index</a></p>
<div class="grid">
`, t, galleryCSS, t, len(p.variants))
	for _, v := range visual {
		fmt.Fprintf(&b, `  <div class="card">%s<div class="label">%s</div></div>
`, v.WrappedSVG, cardLabel(v))
	}
	b.WriteString("</div>\n")

	if len(meta) > 0 {
		fmt.Fprintf(&b, `<details class="meta"><summary>Non-visual attributes (%d) — accessibility / scripting / metadata, no rendering effect</summary>
<div class="grid">
`, len(meta))
		for _, v := range meta {
			fmt.Fprintf(&b, `  <div class="card card--meta">%s<div class="label">%s</div></div>
`, v.WrappedSVG, cardLabel(v))
		}
		b.WriteString("</div></details>\n")
	}

	b.WriteString("</body></html>\n")
	return b.String()
}

// cardLabel renders the monospace attr="value" caption for a variant.
func cardLabel(v Variant) string {
	return fmt.Sprintf(`<span class="attr">%s</span>="%s"`,
		html.EscapeString(v.Attr), html.EscapeString(v.Value))
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

// aNonVisual is the set of <a>-specific navigation-semantics attrs (no visual
// rendering effect — they only control link behavior).
var aNonVisual = map[string]bool{
	"download": true, "ping": true, "rel": true,
	"hreflang": true, "type": true, "referrerpolicy": true,
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
	"inline-size": true, "glyph-orientation-vertical": true,
	// embedded.ebnf
	"pointer-events": true,
	// pservers.ebnf
	"stop-color": true, "stop-opacity": true,
	// styling.ebnf
	"display": true, "visibility": true, "overflow": true, "opacity": true,
	"cursor": true, "transform": true, "transform-origin": true,
}

// paintingGraphics is the set of elements that paint geometry directly (or are
// containers whose presentation attrs inherit to painted children). For these,
// presentationAttrs STAY in the main grid.
var paintingGraphics = map[string]bool{
	"rect": true, "circle": true, "ellipse": true, "line": true,
	"polyline": true, "polygon": true, "path": true,
	"text": true, "tspan": true, "textPath": true,
	"use": true, "g": true, "svg": true, "a": true, "switch": true,
	"foreignObject": true, "image": true, "symbol": true,
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
// tag (it only affects semantics/accessibility/scripting, or it is a paint
// presentation attr applied to an element that paints no geometry). These flood
// galleries with identical cards, so emitPage groups them into the collapsed
// "Non-visual attributes" section. Tag-aware (QA round2).
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
	// 2. <a>-specific navigation-semantics attrs.
	if tag == "a" && aNonVisual[attr] {
		return true
	}
	// 3. The FULL presentation-attribute set (paint/stroke/marker/render hints,
	//    text/font/clip/mask/display/transform/pointer-events/cursor/opacity/
	//    visibility/overflow/…) is non-visual on elements that do not paint
	//    geometry directly. Element-specific paint attrs (flood-color, stop-color,
	//    lighting-color, color-interpolation-filters) keep their main-grid card on
	//    the few elements where they actually apply (keepPaintAttrForTag).
	if presentationAttrs[attr] && !paintingGraphics[tag] {
		return !keepPaintAttrForTag(attr, tag)
	}
	return false
}

// emitIndex renders the index linking every element page.
func emitIndex(pages []page, totalVariants int) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>SVG grammar gallery</title>
<style>%s
ul{columns:4;list-style:none;padding:0}li{margin:2px 0}.n{color:#9aa}</style></head><body>
<h1>SVG grammar — all value-paths</h1>
<p class="desc">%d elements, %d enumerated value-path variants. Each page walks one element's grammar for every attribute and every value, rendered as an SVG.</p>
<ul>
`, galleryCSS, len(pages), totalVariants)
	for _, p := range pages {
		fmt.Fprintf(&b, `  <li><a href="%s.html">&lt;%s&gt;</a> <span class="n">%d</span></li>
`, p.tag, html.EscapeString(p.tag), len(p.variants))
	}
	b.WriteString("</ul></body></html>\n")
	return b.String()
}

// emitValuesJSON builds the element→[{attr,value}] map.
func emitValuesJSON(pages []page) string {
	type av struct {
		Attr  string `json:"attr"`
		Value string `json:"value"`
	}
	m := map[string][]av{}
	for _, p := range pages {
		for _, v := range p.variants {
			m[p.tag] = append(m[p.tag], av{Attr: v.Attr, Value: v.Value})
		}
	}
	data, _ := json.MarshalIndent(m, "", "  ")
	return string(data) + "\n"
}

// emitManifest builds the TSV: tag⇥attr⇥value⇥needsID.
func emitManifest(pages []page) string {
	var b strings.Builder
	b.WriteString("element\tattr\tvalue\tneeds_id\n")
	tags := make([]string, 0, len(pages))
	for _, p := range pages {
		tags = append(tags, p.tag)
	}
	sort.Strings(tags)
	byTag := map[string]page{}
	for _, p := range pages {
		byTag[p.tag] = p
	}
	for _, tag := range tags {
		for _, v := range byTag[tag].variants {
			needs := "0"
			if v.NeedsID {
				needs = "1"
			}
			fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n",
				tag, v.Attr, strings.ReplaceAll(v.Value, "\t", " "), needs)
		}
	}
	return b.String()
}
