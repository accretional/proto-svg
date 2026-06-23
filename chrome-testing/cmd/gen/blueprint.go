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

// blueprintFor returns the scaffold for tag (with a {{ELEMENT}} placeholder).
func (b *blueprintProvider) blueprintFor(tag string) string {
	if s, ok := b.cache[tag]; ok {
		return s
	}
	s := b.fromTemplate(tag)
	if s == "" {
		s = defaultScaffold(tag)
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

// defaultScaffold returns a built-in blueprint derived from tag's content model.
func defaultScaffold(tag string) string {
	switch category(tag) {
	case catGradient:
		return svgOpen +
			`{{ELEMENT}}` +
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
		return svgOpen +
			`<defs><filter id="slot" x="-50%" y="-50%" width="200%" height="200%">{{ELEMENT}}</filter></defs>` +
			`<rect x="20" y="20" width="60" height="60" fill="#e94560" filter="url(#slot)"/></svg>`
	case catTransferFn:
		return svgOpen +
			`<defs><filter id="slot"><feComponentTransfer>{{ELEMENT}}</feComponentTransfer></filter></defs>` +
			`<rect x="10" y="10" width="80" height="80" fill="#4d8bff" filter="url(#slot)"/></svg>`
	case catLight:
		return svgOpen +
			`<defs><filter id="slot"><feDiffuseLighting surface-scale="2" diffuse-constant="1" lighting-color="#ffffff">{{ELEMENT}}</feDiffuseLighting></filter></defs>` +
			`<rect x="10" y="10" width="80" height="80" fill="#888" filter="url(#slot)"/></svg>`
	case catStop:
		return svgOpen +
			`<defs><linearGradient id="slot"><stop offset="0" stop-color="#e94560"/>{{ELEMENT}}<stop offset="1" stop-color="#16c79a"/></linearGradient></defs>` +
			`<rect x="5" y="5" width="90" height="90" fill="url(#slot)"/></svg>`
	case catMergeNode:
		return svgOpen +
			`<defs><filter id="slot"><feFlood flood-color="#e94560"/><feMerge>{{ELEMENT}}</feMerge></filter></defs>` +
			`<rect x="10" y="10" width="80" height="80" filter="url(#slot)"/></svg>`
	case catAnimation:
		return svgOpen +
			`<rect id="target" x="10" y="10" width="40" height="40" fill="#4d8bff">{{ELEMENT}}</rect></svg>`
	case catMpath:
		return svgOpen +
			`<defs><path id="slot" d="M10 50 Q50 10 90 50"/></defs>` +
			`<rect x="0" y="0" width="20" height="20" fill="#f5a623"><animateMotion dur="2s">{{ELEMENT}}</animateMotion></rect></svg>`
	case catDescriptive:
		return svgOpen +
			`<rect x="10" y="10" width="80" height="80" fill="#16c79a">{{ELEMENT}}</rect></svg>`
	case catContainerRef:
		// use references a defined shape
		return svgOpen +
			`<defs><rect id="ref" x="0" y="0" width="40" height="40" fill="#e94560"/></defs>{{ELEMENT}}</svg>`
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
		return add([2]string{"x", "10"}, [2]string{"y", "10"}, [2]string{"width", "80"}, [2]string{"height", "80"}, [2]string{"fill", "#e94560"}), false
	case "circle":
		return add([2]string{"cx", "50"}, [2]string{"cy", "50"}, [2]string{"r", "40"}, [2]string{"fill", "#16c79a"}), false
	case "ellipse":
		return add([2]string{"cx", "50"}, [2]string{"cy", "50"}, [2]string{"rx", "40"}, [2]string{"ry", "25"}, [2]string{"fill", "#f5a623"}), false
	case "line":
		return add([2]string{"x1", "10"}, [2]string{"y1", "10"}, [2]string{"x2", "90"}, [2]string{"y2", "90"}, [2]string{"stroke", "#4d8bff"}, [2]string{"stroke-width", "4"}), false
	case "polyline":
		return add([2]string{"points", "10,80 40,20 70,60 90,10"}, [2]string{"fill", "none"}, [2]string{"stroke", "#e94560"}, [2]string{"stroke-width", "3"}), false
	case "polygon":
		return add([2]string{"points", "50,10 90,80 10,80"}, [2]string{"fill", "#16c79a"}), false
	case "path":
		return add([2]string{"d", "M10 50 Q50 10 90 50 T90 90"}, [2]string{"fill", "none"}, [2]string{"stroke", "#f5a623"}, [2]string{"stroke-width", "3"}), false
	case "text", "tspan", "textPath":
		return add([2]string{"x", "10"}, [2]string{"y", "55"}, [2]string{"fill", "#e6e6e6"}, [2]string{"font-size", "20"}), false
	case "image":
		return add([2]string{"x", "10"}, [2]string{"y", "10"}, [2]string{"width", "80"}, [2]string{"height", "80"}), false
	case "use":
		return add([2]string{"href", "#ref"}, [2]string{"x", "10"}, [2]string{"y", "10"}), false
	case "g", "svg", "defs", "switch", "a", "symbol":
		return "", false
	case "linearGradient":
		return "", false // stops supply the visible color; baseline kept minimal
	case "radialGradient":
		return "", false
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
	case "g", "a", "switch", "svg", "defs", "symbol":
		return `<rect x="20" y="20" width="60" height="60" fill="#4d8bff"/>`
	case "feMerge":
		return `<feMergeNode in="SourceGraphic"/>`
	case "feComponentTransfer":
		return `<feFuncR type="linear" slope="1.5"/>`
	case "feDiffuseLighting", "feSpecularLighting":
		return `<fePointLight x="50" y="50" z="40"/>`
	case "text", "tspan", "textPath":
		return `Ag`
	case "desc", "title", "metadata":
		return `info`
	}
	return ""
}
