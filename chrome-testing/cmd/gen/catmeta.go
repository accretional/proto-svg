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
	// shapes
	"rect":    "The fundamental box primitive — position, size and corner radius.",
	"circle":  "A perfect circle defined by a centre point and radius.",
	"ellipse": "A circle with independent horizontal and vertical radii.",
	// lines & paths
	"line":     "A single straight segment between two points.",
	"polyline": "A connected run of straight segments through a list of points.",
	"polygon":  "A closed shape built from a list of points.",
	"path":     "Arbitrary curves and lines from a `d` command string.",
	// containers
	"svg":           "The root canvas, or a nested viewport with its own coordinates.",
	"g":             "Groups elements so they share attributes and transforms.",
	"defs":          "Holds reusable definitions that render only when referenced.",
	"use":           "Stamps a copy of a defined element by reference.",
	"symbol":        "A reusable template instantiated with its own viewport via `use`.",
	"switch":        "Renders the first child whose conditional tests pass.",
	"a":             "A hyperlink wrapping SVG content.",
	"foreignObject": "Embeds HTML (or other XML) inside SVG.",
	"view":          "A named viewport recalled by a URL fragment.",
	// text
	"text":     "Real, selectable type rendered as vectors.",
	"tspan":    "A styled run inside a text element.",
	"textPath": "Text laid out along an arbitrary path.",
	// paint & fills
	"linearGradient": "A smooth colour blend along a vector, referenced as a fill.",
	"radialGradient": "A colour blend radiating from a focal point.",
	"stop":           "One colour stop within a gradient.",
	"pattern":        "A tiled motif that repeats across any fill or stroke.",
	// clipping & masking
	"clipPath": "Clips an element to an arbitrary shape — hard, 1-bit edges.",
	"mask":     "Masks an element by luminance or alpha — soft, graduated.",
	"marker":   "A symbol drawn at the vertices of a path or shape.",
	// filters
	"filter":              "A container for a chain of filter primitives.",
	"feBlend":             "Blends two inputs with a Photoshop-style blend mode.",
	"feColorMatrix":       "Transforms colours through a matrix — saturate, hue-rotate, recolour.",
	"feComponentTransfer": "Remaps each colour channel through a transfer function.",
	"feComposite":         "Combines two inputs with Porter-Duff or arithmetic compositing.",
	"feConvolveMatrix":    "Convolves pixels with a kernel — blur, sharpen, emboss, edges.",
	"feDiffuseLighting":   "Lights an alpha bump-map with matte (diffuse) shading.",
	"feDisplacementMap":   "Warps one input using another input's pixel values.",
	"feDistantLight":      "An infinitely-distant directional light for a lighting filter.",
	"feDropShadow":        "Offset, blurred drop shadow with its own colour.",
	"feFlood":             "Fills the filter region with a solid colour.",
	"feFuncR":             "The red-channel transfer function inside feComponentTransfer.",
	"feFuncG":             "The green-channel transfer function inside feComponentTransfer.",
	"feFuncB":             "The blue-channel transfer function inside feComponentTransfer.",
	"feFuncA":             "The alpha-channel transfer function inside feComponentTransfer.",
	"feGaussianBlur":      "Softens an input by a standard deviation.",
	"feImage":             "Loads an external image or element as a filter input.",
	"feMerge":             "Stacks several filter results into one layered output.",
	"feMergeNode":         "One input layer inside an feMerge.",
	"feMorphology":        "Fattens (dilate) or thins (erode) an input.",
	"feOffset":            "Shifts an input by dx/dy — the basis of shadows.",
	"fePointLight":        "A point light source at x/y/z for a lighting filter.",
	"feSpecularLighting":  "Lights an alpha bump-map with glossy (specular) highlights.",
	"feSpotLight":         "A cone-shaped spotlight for a lighting filter.",
	"feTile":              "Tiles a filter result across the filter region.",
	"feTurbulence":        "Generates Perlin noise — clouds, marble, organic texture.",
	// content & motion
	"animate":          "Declarative animation — tween an attribute over time.",
	"animateMotion":    "Animate an element along a motion path.",
	"animateTransform": "Animate a transform (translate, scale, rotate…).",
	"set":              "Set an attribute to a value for a span of time.",
	"mpath":            "References a path for animateMotion to follow.",
	"discard":          "Removes its target element at a set time.",
	"image":            "An embedded raster or SVG image.",
	// descriptive / metadata
	"desc":     "A long-form accessible description of its parent element.",
	"title":    "An accessible name / tooltip for its parent element.",
	"metadata": "Machine-readable metadata — RDF, licensing, authorship.",
	"script":   "Embedded or referenced ECMAScript.",
	"style":    "Embedded CSS that styles the document.",
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

// docFor returns the hand-authored long-form documentation for an element (the
// About panel), or nil if none is authored.
func docFor(tag string) *elementDoc {
	if d, ok := docTable[tag]; ok {
		return &d
	}
	return nil
}
