package main

import (
	"strings"
)

// overlay.go — the constraint overlay as a FORWARD GUIDED SAMPLER (per
// docs/CONTEXT_SENSITIVITY.md "Generation mode"). The CFG (grammar) is
// over-approximated; the overlay narrows the values it samples so that every
// emitted attribute value is VALID BY CONSTRUCTION rather than rejected after
// the fact. It is consulted by enumerate.go when it picks a concrete value for
// an attribute value-path.
//
// The overlay never invents structure the grammar does not have; it only
// chooses among the values the grammar already admits (and, for references,
// points them at an id the blueprint guarantees to define — "slot" — or a tiny
// emitted def). The classes implemented here are exactly the ones
// CONTEXT_SENSITIVITY lists as "provably not context-free":
//
//   1. Reference/IDREF resolution — url(#…)/href point at "slot" (the injected
//      element's id) so they never dangle.
//   2. Numeric ranges & monotonicity — opacity/alpha ∈ [0,1]; non-negative
//      where required; keyTimes/keyPoints ∈ [0,1] non-decreasing starting at 0.
//   5. Mutual exclusion — animation values XOR from/to/by.
//   6. Animation value typing — from/to/by/values typed to the chosen
//      attributeName; values count == keyTimes count.

// slotID is the id the blueprint guarantees the injected element carries when it
// is referenced (TEMPLATE_GUIDE convention). All generator-emitted references
// resolve to it by construction.
const slotID = "slot"

// overlaySample chooses a context-valid value string for a leaf/structured value
// in the context of an attribute named attrName on element tag. It returns the
// chosen value and true when the overlay deliberately steered the choice;
// otherwise ("", false) and the caller falls back to the raw grammar sample.
// valueKind is the simple message name of the value type (e.g. "UrlType",
// "NumberType", "IriType").
func overlaySample(tag, attrName, valueKind string) (string, bool) {
	an := strings.ToLower(attrName)

	// Reference target id: animation elements reference the animatable sibling
	// ("target", defined by the animation scaffold); everyone else references the
	// element's own slot id.
	refTarget := slotID
	if isAnimationTag(tag) {
		refTarget = "target"
	}

	// 1a. Gradient href self-reference: a gradient's href/xlink:href pointing at
	//     its own slot id is a circular no-op. Point it at a separately-defined
	//     template gradient ("refgrad") so the inheritance actually resolves.
	//     NOTE for the blueprint owner: the def <linearGradient id="refgrad"> (with
	//     stops) must be emitted in the scaffold for this reference to resolve.
	if isGradientTag(tag) && (an == "href" || an == "xlink:href") {
		return "#refgrad", true
	}

	// 1b. image + href/xlink:href: "#slot" is not a loadable image resource →
	//     broken-image icon. Return the inline data: URI so the image renders.
	if tag == "image" && (an == "href" || an == "xlink:href") {
		return imageDataURI, true
	}

	// 1c. feImage + href/xlink:href: "#slot" is a circular self-reference (the
	//     feImage element itself has id="slot") → blank output. Use the data: URI.
	if tag == "feImage" && (an == "href" || an == "xlink:href") {
		return imageDataURI, true
	}

	// 1. Reference resolution: any url()/IRI value, and any href-like attribute,
	//    points at an id the blueprint defines — never a dangling ref.
	switch valueKind {
	case "UrlType":
		return "url(#" + refTarget + ")", true
	case "IriType":
		// href/xlink:href forms reference the target fragment.
		return "#" + refTarget, true
	}
	if isRefAttr(an) {
		switch valueKind {
		case "UrlType":
			return "url(#" + refTarget + ")", true
		default:
			return "#" + refTarget, true
		}
	}

	// 1d. Per-attribute valid-value overrides (QA round-1 + round-2). These
	//     attributes get a fixed, valid, visible value BY CONSTRUCTION because the
	//     grammar's over-approximation otherwise yields values that hide the element
	//     (conditional processing), are out of range, collapse the geometry, or
	//     produce near-zero/black output. Keyed on the lowercased attribute name
	//     so they win regardless of the value kind the grammar would have sampled.
	switch an {
	// Animation fill-mode: on animation elements, `fill` is "freeze|remove" (the
	// animation fill mode), NOT a paint value. The grammar reuses the paint <fill>
	// production so "none"/colors leak in. Return "freeze" so the card is valid and
	// shows a distinct value from the "remove" enumeration arm.
	case "fill":
		if isAnimationTag(tag) {
			return "freeze", true
		}

	// text/tspan dx/dy: large percentage values (50%/75%) push the first glyph far
	// off the baseline, clipping it out of the viewport. Return a small absolute
	// value so per-glyph offsets stay within the card. x/y are intentionally not
	// overridden here (they are handled by the baseline positioning).
	case "dx", "dy":
		if tag == "text" || tag == "tspan" {
			return "4", true
		}

	// feTurbulence baseFrequency: defaults to 0 → uniform black output on every
	// non-baseFrequency card. "0.05" produces clearly visible noise at 100px scale.
	case "basefrequency":
		if tag == "feTurbulence" {
			return "0.05", true
		}
	// Conditional processing: make the element always satisfiable (it was hidden
	// when these got a StringType "label" value).
	case "systemlanguage":
		return "en", true // always matches the browser language
	case "requiredextensions":
		return "", true // empty list = no required extension → never hidden

	// Stroke miter limit must be >= 1 (was "0", invalid). 4 is the SVG default.
	case "stroke-miterlimit":
		return "4", true

	// Path data: "none" renders blank; give visible geometry.
	case "d":
		return "M10 50 Q50 10 90 50", true

	// Text length: small values (0/1/10) collapse the glyphs.
	case "textlength":
		return "80", true

	// Light / lighting domain: NumberType samples (0/1/-1/3.14/0.5/2) give
	// near-zero illumination → black. Use values that produce visible lighting.
	case "elevation":
		return "45", true // 45° always clearly lit
	case "z":
		// Only meaningful as a light-source depth on point/spot lights.
		if tag == "fePointLight" || tag == "feSpotLight" {
			return "50", true
		}
	case "limitingconeangle":
		return "30", true
	case "specularexponent":
		return "20", true
	case "surfacescale":
		return "5", true

	// feTurbulence: numOctaves is a non-negative int; seed a stable value.
	case "numoctaves":
		return "3", true
	case "seed":
		return "5", true

	// feConvolveMatrix: order is a positive int (was "0.5", invalid).
	case "order":
		return "3", true

	// SMIL clock values: the grammar yields malformed clocks with negative
	// components (e.g. "-1:100.3", "10:0:1.-1"). Override with simple valid
	// unsigned clock values so the timing is always well-formed.
	case "dur":
		return "2s", true
	case "begin":
		return "0s", true
	case "end":
		return "4s", true
	case "min":
		return "0s", true
	case "max":
		return "indefinite", true
	case "repeatdur":
		return "4s", true

	// Animation value typing. The host baseline fixes attributeName="x" (a length
	// attribute), so from/to/by/values must be numeric/length — not colors,
	// transforms, refs, or path-data (those produce no-effect cards). The
	// animateTransform host fixes attributeName="transform", so its values are
	// transform-shaped instead.
	case "attributename":
		// Restrict seeds to real animatable attribute names.
		if tag == "animateTransform" {
			return "transform", true
		}
		return "x", true
	case "from":
		return animFrom(tag), true
	case "to":
		return animTo(tag), true
	case "by":
		return animBy(tag), true
	case "values":
		// FIX 6: feColorMatrix `values` is a numeric matrix/scalar, NOT an animation
		// value list. With the hueRotate baseline it takes a single number; "120" is
		// a clearly-visible rotation. (The animation `values` list applies to
		// animate/set/animateMotion/animateTransform, handled by animValues.)
		if tag == "feColorMatrix" {
			return "120", true
		}
		return animValues(tag), true
	}

	// 3. Numeric ranges & monotonicity.
	switch {
	case isAlphaAttr(an): // opacity / *-opacity ∈ [0,1]
		return "0.5", true
	case an == "keytimes":
		return "0; 0.5; 1", true // non-decreasing, [0,1], first 0, last 1
	case an == "keypoints":
		return "0; 0.5; 1", true
	case an == "keysplines":
		return "0 0 1 1", true // one control-point set ∈ [0,1]
	case isNonNegativeAttr(an):
		switch valueKind {
		case "NumberType":
			return "2", true
		case "LengthPercentageType", "LengthType":
			return "20", true
		case "NumberOptionalNumberType":
			return "2", true
		}
	}

	return "", false
}

// isAnimationTag reports whether tag is an animation element (whose href targets
// an animatable sibling rather than a paint-server slot).
func isAnimationTag(tag string) bool {
	switch tag {
	case "animate", "animateTransform", "animateMotion", "set", "discard", "mpath":
		return true
	}
	return false
}

// isRefAttr reports whether the attribute names a reference (its value is an
// IRI/url to another element). Used so the overlay points it at the slot.
func isRefAttr(an string) bool {
	switch an {
	case "href", "xlink:href", "clip-path", "mask", "filter",
		"marker", "marker-start", "marker-mid", "marker-end",
		"fill", "stroke": // fill/stroke only when they carry a url() — handled by UrlType above
		return false // handled per-value-kind, not per-name, to keep keyword arms intact
	}
	return strings.HasSuffix(an, "href")
}

// isAlphaAttr reports whether the attribute's value is an <alpha-value> in [0,1].
func isAlphaAttr(an string) bool {
	return an == "opacity" || strings.HasSuffix(an, "-opacity")
}

// isNonNegativeAttr reports whether the attribute requires a non-negative value
// (a magnitude). These are the ones the overlay must keep ≥ 0 because the CFG
// over-approximates with a signed NumberType.
func isNonNegativeAttr(an string) bool {
	switch an {
	case "width", "height", "r", "rx", "ry", "stroke-width", "stroke-dashoffset",
		"stddeviation", "stdeviation", "radius", "markerwidth", "markerheight",
		"font-size", "numoctaves", "pathlength", "fr":
		return true
	}
	return false
}

// animationValueFor returns a value of the correct type for an animation
// from/to/by/values, given the resolved attributeName. This is the dependent-type
// oracle (CONTEXT_SENSITIVITY item 6): pick attributeName first, then a matching
// value. It returns a single typed value; the caller composes from/to or values.
func animationValueFor(attrName string) (typeName, fromVal, toVal string) {
	switch strings.ToLower(attrName) {
	case "opacity", "fill-opacity", "stroke-opacity":
		return "number", "0", "1"
	case "fill", "stroke":
		return "color", "#e94560", "#16c79a"
	case "x", "y", "cx", "cy", "r", "rx", "ry", "width", "height",
		"stroke-width", "x1", "y1", "x2", "y2":
		return "length", "10", "80"
	case "transform":
		return "transform", "0", "45"
	case "d":
		return "path", "M10 10 H 90", "M10 90 H 90"
	default:
		return "number", "0", "1"
	}
}

// isGradientTag reports whether tag is a gradient element (whose href inherits
// from another gradient definition rather than self-referencing).
func isGradientTag(tag string) bool {
	return tag == "linearGradient" || tag == "radialGradient"
}

// hostAttributeName returns the attributeName the animation host baseline fixes
// for tag. animate/set animate a length attribute ("x"); animateTransform
// animates the transform; animateMotion is path/coordinate driven.
func hostAttributeName(tag string) string {
	switch tag {
	case "animateTransform":
		return "transform"
	default:
		return "x"
	}
}

// animFrom/animTo/animBy/animValues return from/to/by/values values typed to the
// host's attributeName so the animation is type-correct and observable. For
// animateMotion (path/coordinate driven) they use on-canvas coordinate pairs.
// For animateTransform the values are typed by the baseline transform type
// ("rotate" is the baseline type the blueprint fixes for non-type-varying cards).
func animFrom(tag string) string {
	if tag == "animateMotion" {
		return "10,10"
	}
	if tag == "animateTransform" {
		// Rotate-format: angle cx cy — matches the blueprint baseline type="rotate".
		// translate → "tx ty"; scale → "sx"; skewX/skewY → "angle".
		// Since the overlay does not know the current type value being varied, we
		// return the rotate-format which is correct for the baseline type and is
		// accepted (though cosmetically suboptimal) for other types by browsers.
		return "0 25 25"
	}
	_, from, _ := animationValueFor(hostAttributeName(tag))
	return from
}

func animTo(tag string) string {
	if tag == "animateMotion" {
		return "80,50"
	}
	if tag == "animateTransform" {
		return "360 25 25"
	}
	_, _, to := animationValueFor(hostAttributeName(tag))
	return to
}

func animBy(tag string) string {
	if tag == "animateMotion" {
		return "40,40"
	}
	// "by" is a relative delta of the same type as from/to.
	switch hostAttributeName(tag) {
	case "transform":
		return "45"
	default:
		return "40"
	}
}

func animValues(tag string) string {
	if tag == "animateMotion" {
		return "10,10; 50,5; 80,50"
	}
	if hostAttributeName(tag) == "transform" {
		return "0 25 25; 360 25 25"
	}
	return "10; 45; 80"
}
