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
