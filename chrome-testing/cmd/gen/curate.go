package main

import (
	"strconv"
	"strings"
)

// curate.go — the SEMANTIC layer over enumerated presets. The grammar + overlay
// produce every legal attribute value; this layer turns each into a gallery
// preset that actually DEMONSTRATES what the attribute means and is visually
// distinct from its siblings. Three jobs:
//
//  1. Substitute no-op / default values (filter=none, visibility=visible, …)
//     with a value that shows the effect (a real blur, hidden, a clip).
//  2. Expand a few attributes whose single enumerated value hides their range
//     (animate attributeName → x/y/r/opacity/fill, each a distinct motion).
//  3. Attach a plain-language caption (meaning) to every preset.
//
// Anything not curated here falls back to the raw enumerated value with a
// generic caption from meaningFor.

// demoPreset is one curated gallery preset: the attribute(s) to set on the
// showcased element, a friendly label for the preset name, and a caption.
type demoPreset struct {
	label   string            // shown as "attr=label"; defaults to the single value
	values  map[string]string // attributes applied to the element (≥1)
	meaning string            // one-line plain-language explanation
}

func one(attr, val, meaning string) demoPreset {
	return demoPreset{label: val, values: map[string]string{attr: val}, meaning: meaning}
}

// dropAttr signals curateAttr that the attribute should produce NO presets (it
// has no visual effect on this element). Distinct from a nil return, which means
// "fall back to the raw enumerated values".
var dropAttr = []demoPreset{}

// openStrokeTags are elements with open stroked paths, where stroke-linecap is
// visible; corneredStrokeTags have angular joins, where stroke-linejoin /
// stroke-miterlimit are visible. Everywhere else those attrs are dropped.
var openStrokeTags = map[string]bool{"line": true, "polyline": true, "path": true}
var corneredStrokeTags = map[string]bool{"polyline": true, "polygon": true, "path": true, "rect": true}

// noFillArea elements have stroke but no fillable interior at default, so the
// fill family (fill-opacity / fill-rule / paint-order) can't be demonstrated.
var openStrokeOnly = map[string]bool{"line": true, "polyline": true, "path": true}

// curateAttr returns curated presets for (tag, attr) given the enumerated values,
// or nil to fall back to the raw values (each captioned by meaningFor). When it
// returns presets that reference a demo def (url(#fx-…)) the caller injects the
// shared demoDefs into the element's base.
func curateAttr(tag, attr string, vals []string) []demoPreset {
	// animation elements: attributeName/href/by/etc. need their own treatment.
	if isAnimationTag(tag) {
		if dp := curateAnimAttr(tag, attr, vals); dp != nil {
			return dp
		}
	}

	switch attr {
	// mask / clip-path / marker references are resolved to real demo defs by the
	// OVERLAY (overlaySample → #fx-mask / #fx-clip / #fx-marker) and surfaced by the
	// enumeration's demonstrative-arm pick. `filter` is the exception: its grammar
	// value (FilterValueList) samples a messy multi-item list with a dangling url()
	// and no-op items, so select a single clean, grammar-legal url() reference.
	case "filter":
		return []demoPreset{one("filter", "url(#fx-blur)", "runs the element through a blur filter")}
	case "visibility":
		return []demoPreset{one("visibility", "hidden", "hides the element (it still occupies its space)")}
	case "display":
		return []demoPreset{one("display", "none", "removes the element from the render — it vanishes")}
	case "overflow":
		return dropAttr // specimens never overflow their viewport, so nothing clips
	case "fill":
		if tag == "line" {
			return dropAttr // a line has no fill area
		}
		return nil
	case "fill-opacity", "fill-rule":
		if openStrokeOnly[tag] {
			return dropAttr // no filled interior to modulate on open/stroked shapes
		}
		return nil
	case "paint-order":
		if openStrokeOnly[tag] {
			return dropAttr // no fill, so stroke/fill paint order is invisible
		}
		return []demoPreset{one("paint-order", "stroke", "draws the stroke UNDER the fill (vs over it)")}
	case "font-style", "font-variant", "font-weight", "font-stretch":
		// drop the `normal` default (no visible change); keep the variants that do
		// render (italic / bold / small-caps / condensed …).
		var out []demoPreset
		for _, v := range vals {
			if v == "normal" {
				continue
			}
			out = append(out, one(attr, v, meaningFor(tag, attr, v)))
		}
		return out
	case "text-decoration":
		return []demoPreset{one("text-decoration", "underline", "underlines the text")}
	case "letter-spacing":
		return []demoPreset{one("letter-spacing", "6", "adds 6 units of space between letters")}
	case "dominant-baseline", "alignment-baseline":
		return []demoPreset{one(attr, "central", "aligns the glyphs by their CENTER baseline (vs alphabetic)")}
	case "transform-origin":
		return dropAttr // no effect without a transform on the same element
	case "stroke-dasharray":
		// `none` (the only enumerated value) is the solid default = a no-op. Show a
		// real dash pattern instead.
		return []demoPreset{one("stroke-dasharray", "7 5", "dashes the outline (7 on, 5 off)")}
	case "stroke-dashoffset":
		// only meaningful with a dash pattern, so set one alongside the offset.
		return []demoPreset{{
			label:   "10",
			values:  map[string]string{"stroke-dasharray": "8 6", "stroke-dashoffset": "10"},
			meaning: "shifts the start of the dash pattern along the outline",
		}}
	case "stroke-linecap":
		if !openStrokeTags[tag] {
			return dropAttr // closed/filled shapes have no open stroke ends
		}
		return []demoPreset{
			one("stroke-linecap", "round", "rounds the open ends of the stroke"),
			one("stroke-linecap", "square", "square caps that extend past the stroke ends"),
		}
	case "stroke-linejoin":
		if !corneredStrokeTags[tag] {
			return dropAttr
		}
		return []demoPreset{
			one("stroke-linejoin", "round", "rounds the corners where stroke segments meet"),
			one("stroke-linejoin", "bevel", "flattens (bevels) the stroke corners"),
		}
	case "stroke-miterlimit":
		if !corneredStrokeTags[tag] {
			return dropAttr
		}
		return []demoPreset{one("stroke-miterlimit", "1", "clips sharp miter joins to bevels below the limit")}
	case "d":
		// drop the empty/none path (renders nothing); caption the real geometries.
		var out []demoPreset
		for _, v := range vals {
			if v == "none" || strings.TrimSpace(v) == "" {
				continue
			}
			out = append(out, one("d", v, "path geometry: "+pathKind(v)))
		}
		return out
	case "width", "height":
		if tag == "use" {
			return dropAttr // only sizes referenced <svg>/<symbol>; a plain shape ref ignores it
		}
		// `auto` is the intrinsic-size default (no visible change); keep the numeric
		// presets that visibly resize, and add one if auto was the only value.
		var out []demoPreset
		for _, v := range vals {
			if v == "auto" {
				continue
			}
			out = append(out, one(attr, v, meaningFor(tag, attr, v)))
		}
		if len(out) == 0 {
			out = append(out, one(attr, "60", attr+" of 60 units"))
		}
		return out
	case "rx", "ry":
		if tag != "rect" {
			return nil // on ellipse etc. rx/ry are the radii — keep raw + captioned
		}
		var out []demoPreset
		for _, v := range vals {
			if v == "auto" {
				out = append(out, one(attr, "16", "rounds the corners ("+attr+" radius 16)"))
			} else {
				out = append(out, one(attr, v, meaningFor(tag, attr, v)))
			}
		}
		return out

	// ---- filter-primitive inputs & magnitudes ----
	case "in", "in2":
		// BackgroundImage/BackgroundAlpha are unsupported in browsers; FillPaint/
		// StrokePaint collapse to the source here — all render identically. Keep only
		// the inputs that visibly differ.
		var out []demoPreset
		for _, v := range vals {
			switch v {
			case "BackgroundImage", "BackgroundAlpha", "FillPaint", "StrokePaint":
				continue // unsupported in browsers — render identically to the source
			case "result1":
				continue // a second handle on the same upstream result as blur1
			}
			out = append(out, one(attr, v, meaningFor(tag, attr, v)))
		}
		return out
		// (dx/dy, scale, azimuth, light x/y, pointsAtX/Y are pinned to visible values
		// by the OVERLAY now; k1-k4→operator=arithmetic by companionFor.)
	case "x", "y":
		if tag == "pattern" {
			return dropAttr // an origin phase-shift is invisible on a seamless tile
		}
		if tag == "filter" {
			return dropAttr // large region offsets push the filter region off-canvas → blank
		}
		return nil

	// ---- filter-primitive no-op / out-of-range pruning ----
	case "numOctaves":
		var out []demoPreset
		for _, v := range vals {
			if v == "0" {
				continue // zero octaves → blank
			}
			out = append(out, one("numOctaves", v, v+" octaves of noise detail"))
		}
		return out
	case "stitchTiles":
		return dropAttr // only differs when the noise tiles across regions
	case "targetX", "targetY":
		return dropAttr // valid only in 0..order-1; the enumerated values blank out
	case "order":
		return dropAttr // pinned to match the 3×3 kernel — no real variation
	case "preserveAlpha":
		return dropAttr // no visible difference on opaque content
	case "edgeMode":
		var out []demoPreset
		for _, v := range vals {
			if v == "duplicate" || v == "wrap" {
				continue // only a 1px border differs — imperceptible
			}
			out = append(out, one("edgeMode", v, v+" handling at the filter edges"))
		}
		return out
	case "bias":
		return []demoPreset{one("bias", "0.2", "lifts every filtered channel by 0.2")}
	case "diffuseConstant":
		if tag == "feDiffuseLighting" {
			return []demoPreset{
				one("diffuseConstant", "0.4", "dim diffuse reflection"),
				one("diffuseConstant", "1", "normal diffuse reflection"),
				one("diffuseConstant", "1.6", "bright diffuse reflection"),
			}
		}
		return nil

	// (feFunc type / slope / amplitude / exponent / offset / tableValues are paired
	// with the param that makes them visible by companionFor — the grammar's own
	// value-path is kept; the blueprint supplies the companion type=gamma/table etc.)
	case "kernelMatrix":
		if tag == "feConvolveMatrix" {
			return []demoPreset{
				{label: "edge", values: map[string]string{"kernelMatrix": "0 -1 0 -1 4 -1 0 -1 0", "divisor": "1", "bias": "0"}, meaning: "edge-detection kernel"},
				{label: "sharpen", values: map[string]string{"kernelMatrix": "0 -1 0 -1 5 -1 0 -1 0", "divisor": "1", "bias": "0"}, meaning: "sharpen kernel"},
				{label: "emboss", values: map[string]string{"kernelMatrix": "-2 -1 0 -1 1 1 0 1 2", "divisor": "1", "bias": "0"}, meaning: "emboss kernel"},
				{label: "blur", values: map[string]string{"kernelMatrix": "1 1 1 1 1 1 1 1 1", "divisor": "9", "bias": "0"}, meaning: "box-blur kernel"},
			}
		}
		return nil

	// ---- reference-only attrs that can't show standalone ----
	case "preserveAspectRatio":
		switch tag {
		case "svg", "symbol", "image", "feImage":
			return dropAttr // content aspect matches the viewport — no visible effect
		}
		return nil // marker/pattern preserveAspectRatio is visible — keep
	case "refX", "refY":
		if tag == "symbol" {
			return dropAttr // only shifts relative to a <use> point — invisible standalone
		}
		return nil // marker refX/refY visibly move the arrowhead — keep
	case "patternContentUnits", "patternUnits":
		return dropAttr // objectBoundingBox collapses/blanks the tile; userSpaceOnUse = base
	}
	return nil
}

// curateAnimAttr curates animation-element attributes. attributeName expands to
// several animatable properties (each a visibly different motion); calcMode and
// the timing/value attrs get captions tuned to what they control. The host is the
// catAnimation rect (id="target", 40×40 at 10,10).
func curateAnimAttr(tag, attr string, vals []string) []demoPreset {
	a := strings.ToLower(attr)
	switch a {
	case "attributename":
		if tag == "animateTransform" || tag == "animateMotion" {
			return nil // these animate a fixed channel (transform / motion path)
		}
		mk := func(name, from, to, meaning string) demoPreset {
			return demoPreset{
				label:   name,
				values:  map[string]string{"attributeName": name, "from": from, "to": to},
				meaning: meaning,
			}
		}
		return []demoPreset{
			mk("x", "10", "55", "animates the x coordinate — slides right"),
			mk("y", "10", "55", "animates the y coordinate — slides down"),
			mk("width", "20", "78", "animates the width — grows wider"),
			mk("opacity", "1", "0.15", "animates opacity — fades out"),
			mk("fill", "#4d8bff", "#f5a623", "animates the fill — recolors blue→orange"),
		}
	case "href", "xlink:href":
		return []demoPreset{{
			label:   "#target",
			values:  map[string]string{attr: "#target"},
			meaning: "retargets the animation to the element with id=target",
		}}
	case "calcmode":
		if tag == "animateTransform" || tag == "animateMotion" {
			// these animate the transform / motion-path; calcMode applies to the base
			// animation directly. Injecting an x-based values list would break it, so
			// just caption the raw value.
			capm := map[string]string{
				"discrete": "jumps between values in steps — no tween",
				"linear":   "constant-speed interpolation",
				"paced":    "constant VELOCITY along the value path",
				"spline":   "eased via keySplines (slow → fast → slow)",
			}
			var out []demoPreset
			for _, v := range vals {
				out = append(out, one("calcMode", v, defStr(capm[v], "calcMode: "+v)))
			}
			return out
		}
		// animate/set: each calcMode preset is self-contained (sets attributeName + a
		// multi-stop values list) so the STROBE shows the velocity profile: discrete =
		// clustered stamps, linear/paced = evenly spaced, spline = bunched at the ends.
		mk := func(mode, vals, meaning string, extra map[string]string) demoPreset {
			m := map[string]string{"calcMode": mode, "attributeName": "x", "values": vals}
			for k, x := range extra {
				m[k] = x
			}
			return demoPreset{label: mode, values: m, meaning: meaning}
		}
		return []demoPreset{
			mk("discrete", "10; 30; 52; 78", "jumps between values in steps — no tween", nil),
			mk("linear", "10; 30; 52; 78", "constant-speed interpolation between values", nil),
			mk("paced", "10; 30; 52; 78", "constant VELOCITY along the value path", nil),
			mk("spline", "10; 78", "eased via keySplines (slow → fast → slow)",
				map[string]string{"keyTimes": "0; 1", "keySplines": "0.85 0 0.15 1"}),
		}
	case "values":
		if tag == "animateTransform" || tag == "animateMotion" {
			return captioned("values", vals, "animates through an explicit list of values")
		}
		return []demoPreset{{
			label:   "10; 45; 80",
			values:  map[string]string{"attributeName": "x", "values": "10; 45; 80"},
			meaning: "animates through an explicit list of values in order",
		}}
	case "keytimes":
		if tag == "animateTransform" || tag == "animateMotion" {
			return captioned("keyTimes", vals, "sets WHEN (0–1 of the duration) each value is reached")
		}
		return []demoPreset{{
			label:   "0; 0.8; 1",
			values:  map[string]string{"attributeName": "x", "values": "10; 44; 78", "keyTimes": "0; 0.8; 1", "calcMode": "linear"},
			meaning: "warps WHEN each value is reached — most of the time on the first leg",
		}}
	case "keysplines":
		if tag == "animateTransform" || tag == "animateMotion" {
			return captioned("keySplines", vals, "Bézier easing control points (with calcMode=spline)")
		}
		return []demoPreset{{
			label:   "ease",
			values:  map[string]string{"attributeName": "x", "values": "10; 78", "keyTimes": "0; 1", "calcMode": "spline", "keySplines": "0.85 0 0.15 1"},
			meaning: "Bézier easing control points for spline calcMode (slow → fast → slow)",
		}}
	case "from":
		if tag == "animateMotion" {
			return dropAttr // path-driven base ignores from/to/by
		}
		return captioned("from", vals, "the value the animation starts FROM")
	case "to":
		if tag == "animateMotion" {
			return dropAttr
		}
		return captioned("to", vals, "the value the animation ends AT")
	case "by":
		if tag == "animateMotion" {
			return dropAttr
		}
		return captioned("by", vals, "a relative delta added over the animation (instead of to)")
	case "type":
		if tag == "animateTransform" {
			cap := map[string]string{
				"rotate":    "rotates the element over time",
				"scale":     "scales the element over time",
				"translate": "translates (moves) the element over time",
				"skewX":     "skews the element horizontally over time",
				"skewY":     "skews the element vertically over time",
			}
			var out []demoPreset
			for _, v := range vals {
				out = append(out, demoPreset{label: v, values: map[string]string{"type": v, "from": transFrom(v), "to": transTo(v)}, meaning: defStr(cap[v], "transform type "+v)})
			}
			return out
		}
	case "path":
		if tag == "animateMotion" {
			return captioned("path", vals, "the motion path the element travels along")
		}
	case "rotate":
		if tag == "animateMotion" {
			return captioned("rotate", vals, "orients the element along the motion path ("+strings.Join(vals, ", ")+")")
		}
	}
	return nil
}

// transFrom/transTo give type-correct from/to for animateTransform `type` presets.
func transFrom(t string) string {
	switch t {
	case "scale":
		return "1"
	case "translate":
		return "0 0"
	case "skewX", "skewY":
		return "0"
	default: // rotate
		return "0 30 30"
	}
}
func transTo(t string) string {
	switch t {
	case "scale":
		return "1.8"
	case "translate":
		return "40 30"
	case "skewX", "skewY":
		return "35"
	default: // rotate
		return "360 30 30"
	}
}

func captioned(label string, vals []string, meaning string) []demoPreset {
	var out []demoPreset
	for _, v := range vals {
		out = append(out, one(label, v, meaning))
	}
	return out
}

func defStr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// meaningFor returns a plain-language caption for an UN-curated preset (the raw
// enumerated value). Keyed by attribute, with the value interpolated.
func meaningFor(tag, attr, val string) string {
	a := strings.ToLower(attr)
	switch a {
	case "fill":
		if isAnimationTag(tag) {
			return "fill mode after the animation ends (" + val + ")"
		}
		if val == "none" {
			return "no fill — the interior is transparent"
		}
		return "fills the interior with " + val
	case "stroke":
		if val == "none" {
			return "removes the outline stroke"
		}
		return "outlines the shape in " + val
	case "color":
		return "the currentColor descendants inherit (" + val + ")"
	case "filter":
		if strings.Contains(val, "fx-blur") {
			return "runs the element through a blur filter"
		}
		return "applies the filter " + val
	case "mask":
		if strings.Contains(val, "fx-mask") {
			return "fades the element out left→right via a gradient mask"
		}
		return "masks the element with " + val
	case "clip-path":
		if strings.Contains(val, "fx-clip") {
			return "clips the element to a star"
		}
		return "clips the element with " + val
	case "marker", "marker-start", "marker-mid", "marker-end":
		where := map[string]string{
			"marker": "every vertex", "marker-start": "the start vertex",
			"marker-mid": "each mid vertex", "marker-end": "the end vertex",
		}
		return "draws an arrowhead marker at " + where[a]
	case "href", "xlink:href":
		switch {
		case tag == "use":
			return "clones the element with id " + val
		case isGradientTag(tag):
			return "inherits gradient stops from " + val
		case tag == "image" || tag == "feImage":
			return "loads the referenced image source"
		case tag == "textPath":
			return "lays the text along path " + val
		case tag == "mpath":
			return "the path the motion follows (" + val + ")"
		case tag == "pattern" || tag == "filter":
			return "inherits content from " + val
		case isAnimationTag(tag):
			return "retargets the animation to " + val
		}
		return "references " + val
	case "rotate":
		if tag == "text" || tag == "tspan" {
			return "rotates each glyph by the listed angles (" + val + ")"
		}
		return "rotates by " + val
	case "dominant-baseline", "alignment-baseline":
		return "aligns the glyphs to the " + val + " baseline"
	case "baseline-shift":
		return "shifts the text baseline (" + val + " — e.g. super/subscript)"
	case "direction":
		return val + " text direction"
	case "writing-mode":
		return val + " writing mode"
	case "unicode-bidi":
		return "bidirectional handling: " + val
	case "z":
		return "light-source depth z = " + val
	case "seed":
		return "random seed " + val + " for the noise field"
	case "keypoints":
		return "fraction along the motion path at each keyTime (" + val + ")"
	case "origin":
		return "motion origin (" + val + ")"
	case "attributename":
		return "animates the " + val + " attribute"
	case "stop-color", "flood-color", "lighting-color", "solid-color":
		return "the " + strings.TrimSuffix(a, "-color") + " color (" + val + ")"
	case "fill-opacity":
		return "fill at " + pct(val) + " opacity"
	case "stroke-opacity":
		return "outline at " + pct(val) + " opacity"
	case "opacity", "flood-opacity", "stop-opacity", "solid-opacity":
		return "at " + pct(val) + " opacity"
	case "stroke-width":
		return "outline " + val + " units thick"
	case "stroke-dasharray":
		if val == "none" {
			return "solid (un-dashed) outline"
		}
		return "dashed outline, pattern " + val
	case "stroke-dashoffset":
		return "shifts the dash pattern start by " + val
	case "stroke-linecap":
		return val + " caps on the open ends of the stroke"
	case "stroke-linejoin":
		return val + " corners where stroke segments meet"
	case "stroke-miterlimit":
		return "miter join cutoff ratio " + val
	case "transform":
		return "applies the transform " + val
	case "transform-origin":
		return "moves the transform pivot to " + val
	case "x", "y":
		return "positions the " + a + "-coordinate at " + val
	case "cx", "cy":
		return "centers the " + a + " at " + val
	case "r":
		return "radius " + val
	case "rx", "ry":
		return a + " corner/axis radius " + val
	case "width", "height":
		return a + " of " + val
	case "d":
		return "path geometry: " + pathKind(val)
	case "path":
		return "follows the path: " + pathKind(val)
	case "points":
		return "vertex list " + val
	case "viewbox":
		return "maps the coordinate window to " + val
	case "preserveaspectratio":
		return "fits content using " + val
	case "offset":
		return "gradient stop at " + pct(val) + " along the vector"
	case "gradientunits", "patterncontentunits", "clippathunits", "maskunits", "patternunits", "primitiveunits", "filterunits", "maskcontentunits":
		return "interprets coordinates in " + val
	case "spreadmethod":
		return val + " spread beyond the gradient ends"
	case "font-size":
		return "text " + val + " units tall"
	case "font-weight":
		return val + " stroke weight text"
	case "font-style":
		return val + " text"
	case "font-family":
		return val + " typeface"
	case "text-anchor":
		return "anchors text to its " + val + " edge"
	case "letter-spacing":
		return "extra " + val + " between letters"
	case "textlength":
		return "stretches/squeezes the text to " + val + " units wide"
	case "lengthadjust":
		return "fits text to length via " + val
	case "type":
		return val + " type"
	case "in", "in2":
		return "takes " + val + " as input"
	case "mode":
		return val + " blend mode"
	case "operator":
		return val + " compositing operator"
	case "values":
		if tag == "feColorMatrix" {
			return "colour-transform parameter (" + val + ")"
		}
		return "animates through the values " + val
	case "stddeviation":
		return "blur radius " + val
	case "scale":
		return "displacement scale " + val
	case "basefrequency":
		return "noise frequency " + val
	case "numoctaves":
		return val + " octaves of noise detail"
	case "dx":
		return "offsets " + val + " units horizontally"
	case "dy":
		return "offsets " + val + " units vertically"
	case "x1", "y1", "x2", "y2":
		if isGradientTag(tag) {
			return "gradient-vector " + a + " endpoint at " + val
		}
		return a + " endpoint at " + val
	case "fx", "fy":
		return "radial gradient focal " + strings.TrimPrefix(a, "f") + " at " + val
	case "fr":
		return "radial gradient focal radius " + val
	case "refx", "refy":
		return "anchor/reference " + strings.TrimPrefix(a, "ref") + " at " + val
	case "markerwidth", "markerheight":
		return "marker viewport " + strings.TrimPrefix(a, "marker") + " of " + val
	case "markerunits":
		return "marker sized in " + val
	case "orient":
		return "orients the marker: " + val
	case "startoffset":
		return "starts the text " + val + " along the path"
	case "method":
		return val + " glyph placement along the path"
	case "side":
		return "renders text on the " + val + " side of the path"
	case "spacing":
		return val + " glyph spacing along the path"
	case "gradienttransform", "patterntransform":
		return "transforms the " + strings.TrimSuffix(a, "transform") + " by " + val
	case "k1", "k2", "k3", "k4":
		return "arithmetic-composite coefficient " + a + " = " + val
	case "order":
		return "convolution kernel " + val + "×" + val
	case "kernelmatrix":
		return "convolution kernel [" + val + "]"
	case "divisor":
		return "divides the kernel sum by " + val
	case "bias":
		return "adds " + val + " to each filtered channel"
	case "targetx", "targety":
		return "kernel target " + strings.TrimPrefix(a, "target") + " = " + val
	case "edgemode":
		return val + " handling at the filter edges"
	case "preservealpha":
		return "preserve-alpha = " + val
	case "slope":
		return "transfer-function slope " + val
	case "intercept":
		return "transfer-function intercept " + val
	case "amplitude":
		return "gamma transfer amplitude " + val
	case "exponent":
		return "gamma transfer exponent " + val
	case "tablevalues":
		return "remaps the channel through the table [" + val + "]"
	case "azimuth":
		return "light direction (azimuth) " + val + "°"
	case "elevation":
		return "light elevation " + val + "°"
	case "surfacescale":
		return "bump-map height scale " + val
	case "diffuseconstant":
		return "diffuse reflection constant " + val
	case "specularconstant":
		return "specular reflection constant " + val
	case "specularexponent":
		return "specular highlight sharpness " + val
	case "pointsatx", "pointsaty", "pointsatz":
		return "spotlight aim " + strings.TrimPrefix(a, "pointsat") + " = " + val
	case "limitingconeangle":
		return "spotlight cone half-angle " + val + "°"
	case "xchannelselector", "ychannelselector":
		return "displaces along " + strings.TrimSuffix(strings.ToUpper(a[:1])+a[1:], "channelselector") + " using the " + val + " channel"
	case "radius":
		return "morphology radius " + val
	case "mask-type":
		return "masks using the " + val + " channel"
	case "fill-rule", "clip-rule":
		return val + " rule for self-intersecting/overlapping subpaths"
	case "stitchtiles":
		return val + " noise across tile edges"
	}
	if val == "none" || val == "auto" || val == "normal" || val == "visible" {
		return "the default " + attr + " (" + val + ")"
	}
	return "sets " + attr + " to " + val
}

// ---- small helpers ----

func pct(val string) string {
	if f, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err == nil {
		if f <= 1 {
			return strconv.Itoa(int(f*100+0.5)) + "%"
		}
		return val + "%"
	}
	return val
}

func pathKind(d string) string {
	u := strings.ToUpper(d)
	switch {
	case strings.ContainsAny(u, "CS"):
		return "cubic Bézier curve"
	case strings.Contains(u, "Q") || strings.Contains(u, "T"):
		return "quadratic curve"
	case strings.Contains(u, "A"):
		return "elliptical arc"
	case strings.Contains(u, "Z"):
		return "closed polygon (lines)"
	case strings.ContainsAny(u, "HVL"):
		return "straight line segments"
	}
	return d
}

// demoDefs is a shared <defs> block of reusable demo resources that substituted
// preset values (filter=url(#fx-blur), mask=url(#fx-mask), …) reference. It is
// injected into an element's base only when a curated preset uses one.
func demoDefs() string {
	return `<defs>` +
		`<filter id="fx-blur" x="-40%" y="-40%" width="180%" height="180%"><feGaussianBlur stdDeviation="3.2"/></filter>` +
		`<linearGradient id="fx-maskgrad"><stop offset="0.15" stop-color="#fff"/><stop offset="1" stop-color="#000"/></linearGradient>` +
		`<mask id="fx-mask"><rect x="0" y="0" width="100" height="100" fill="url(#fx-maskgrad)"/></mask>` +
		`<marker id="fx-marker" markerWidth="7" markerHeight="7" refX="3.5" refY="3.5" orient="auto"><path d="M0 0 L7 3.5 L0 7 Z" fill="#f5a623"/></marker>` +
		`<clipPath id="fx-clip"><polygon points="50,6 61,38 95,38 67,58 78,92 50,72 22,92 33,58 5,38 39,38"/></clipPath>` +
		`</defs>`
}

// injectDemoDefs inserts demoDefs() immediately after the root <svg …> open tag
// of base, if not already present.
func injectDemoDefs(base string) string {
	if strings.Contains(base, `id="fx-blur"`) {
		return base
	}
	i := strings.Index(base, ">")
	if i < 0 {
		return base
	}
	return base[:i+1] + demoDefs() + base[i+1:]
}
