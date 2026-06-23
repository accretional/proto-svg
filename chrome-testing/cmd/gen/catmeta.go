package main

import (
	"regexp"
	"strings"
)

// catmeta.go — presentational metadata for the catalogue: the rail groups, the
// human element names, and one-line descriptions. Kept separate from the data
// derivation in catalogue.go.

var groupOrder = []string{
	"Shapes", "Lines & paths", "Containers", "Text",
	"Paint & fills", "Clipping & masking", "Filters", "Content & motion", "Descriptive",
}

var groupByTag = map[string]string{
	"rect": "Shapes", "circle": "Shapes", "ellipse": "Shapes",
	"line": "Lines & paths", "polyline": "Lines & paths", "polygon": "Lines & paths", "path": "Lines & paths",
	"svg": "Containers", "g": "Containers", "defs": "Containers", "use": "Containers",
	"symbol": "Containers", "switch": "Containers", "a": "Containers", "foreignObject": "Containers", "view": "Containers",
	"text": "Text", "tspan": "Text", "textPath": "Text",
	"linearGradient": "Paint & fills", "radialGradient": "Paint & fills", "stop": "Paint & fills", "pattern": "Paint & fills",
	"clipPath": "Clipping & masking", "mask": "Clipping & masking", "marker": "Clipping & masking",
	"animate": "Content & motion", "animateMotion": "Content & motion", "animateTransform": "Content & motion",
	"set": "Content & motion", "mpath": "Content & motion", "discard": "Content & motion", "image": "Content & motion",
	"desc": "Descriptive", "title": "Descriptive", "metadata": "Descriptive", "script": "Descriptive", "style": "Descriptive",
}

func groupFor(tag string) string {
	if g, ok := groupByTag[tag]; ok {
		return g
	}
	if strings.HasPrefix(tag, "fe") {
		return "Filters"
	}
	if tag == "filter" {
		return "Filters"
	}
	return "Containers"
}

// nameFor renders a human title from a tag: fe-primitives drop the "fe" and read
// as "Gaussian blur"; camelCase splits into words; bare tags are capitalised.
var camelRe = regexp.MustCompile(`([a-z0-9])([A-Z])`)

var nameTable = map[string]string{
	"rect": "Rectangle", "g": "Group", "a": "Anchor", "svg": "SVG root",
	"defs": "Definitions", "use": "Use", "tspan": "Text span", "textPath": "Text path",
	"foreignObject": "Foreign object", "mpath": "Motion path", "desc": "Description",
	"linearGradient": "Linear gradient", "radialGradient": "Radial gradient",
	"clipPath": "Clip path", "animateMotion": "Animate motion", "animateTransform": "Animate transform",
}

func nameFor(tag string) string {
	if n, ok := nameTable[tag]; ok {
		return n
	}
	if strings.HasPrefix(tag, "fe") && len(tag) > 2 {
		rest := tag[2:] // feGaussianBlur -> GaussianBlur
		rest = camelRe.ReplaceAllString(rest, "$1 $2")
		return strings.ToUpper(rest[:1]) + strings.ToLower(rest[1:])
	}
	s := camelRe.ReplaceAllString(tag, "$1 $2")
	return strings.ToUpper(s[:1]) + s[1:]
}

var descTable = map[string]string{
	"rect":           "The fundamental box primitive — position, size and corner radius.",
	"circle":         "A perfect circle defined by a centre point and radius.",
	"ellipse":        "A circle with independent horizontal and vertical radii.",
	"line":           "A single straight segment between two points.",
	"polyline":       "A connected run of straight segments through a list of points.",
	"polygon":        "A closed shape built from a list of points.",
	"path":           "Arbitrary curves and lines from a `d` command string.",
	"text":           "Real, selectable type rendered as vectors.",
	"tspan":          "A styled run inside a text element.",
	"textPath":       "Text laid out along an arbitrary path.",
	"linearGradient": "A smooth colour blend along a vector, referenced as a fill.",
	"radialGradient": "A colour blend radiating from a focal point.",
	"stop":           "One colour stop within a gradient.",
	"pattern":        "A tiled motif that repeats across any fill.",
	"filter":         "A container for a chain of filter primitives.",
	"feGaussianBlur": "Softens an input by a standard deviation.",
	"feDropShadow":   "Offset, blurred drop shadow with its own colour.",
	"clipPath":       "Clips an element to an arbitrary shape.",
	"mask":           "Masks an element by luminance or alpha.",
	"marker":         "A symbol drawn at the vertices of a path or shape.",
	"animate":        "Declarative animation — tween an attribute over time.",
	"animateTransform": "Animate a transform (translate, scale, rotate…).",
	"animateMotion":  "Animate an element along a motion path.",
	"set":            "Set an attribute to a value for a span of time.",
	"image":          "An embedded raster or SVG image.",
	"use":            "Reuse a defined element by reference.",
	"g":              "Groups elements so they share attributes and transforms.",
	"svg":            "The root canvas (or a nested viewport).",
}

func descFor(tag string) string {
	if d, ok := descTable[tag]; ok {
		return d
	}
	if strings.HasPrefix(tag, "fe") {
		return "A filter primitive (" + tag + ") in the filter chain."
	}
	return "The <" + tag + "> SVG element."
}
