package main

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// catalogue.go — emits chrome-testing/gallery/catalogue.json, the data contract
// the SVG Lab gallery renders from. For every element it produces, from the
// grammar's enumerated value-paths:
//
//   - a BASE svg (the blueprint scaffold with the element at its baseline),
//   - typed ATTRIBUTE CONTROLS (paint / range / select / number / text) derived
//     from each visual attribute's value-paths,
//   - DEFAULTS for those controls, and
//   - PRESETS: one per visual value-path (attr=value), which the viewer applies
//     and the shoot pipeline screenshots.
//
// This replaces the static per-element gallery HTML: the gallery is now one live
// app driven entirely by this catalogue.

const catalogueDir = "chrome-testing/gallery"

// animationTags are the SVG animation elements — their render is temporal (it
// animates), so the shoot pipeline captures a GIF frame sequence rather than a
// single PNG.
var animationTags = map[string]bool{
	"animate": true, "set": true, "animateMotion": true,
	"animateTransform": true, "mpath": true, "discard": true,
}

// liveSMILRe detects a live SMIL animation anywhere in a wrapped SVG (e.g. one
// the blueprint scaffold injects around a static element) — such an element is
// temporal even when the showcased element is itself static.
var liveSMILRe = regexp.MustCompile(`<(animate|animateTransform|animateMotion|set)[\s>/]`)

type catControl struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Control string   `json:"control"` // paint | range | select | number | text
	Min     *float64 `json:"min,omitempty"`
	Max     *float64 `json:"max,omitempty"`
	Step    *float64 `json:"step,omitempty"`
	Options []string `json:"options,omitempty"`
}

type catPreset struct {
	Name    string            `json:"name"`
	Values  map[string]string `json:"values"`
	Meaning string            `json:"meaning,omitempty"`
	// Prov is "grammar" when the preset is a raw enumerated value-path (leaf value
	// sampled by the overlay, companions by companionFor) or "curated" when the
	// curation layer selected a demonstrative value or expanded a showcase set —
	// so the walkthrough-vs-showcase split is auditable.
	Prov string `json:"prov,omitempty"`
	// Interactive marks a preset whose effect only appears under user input —
	// "hover" or "click". The shoot drives a real mouse Hover/Click and captures the
	// before/after, demonstrating that event attributes ARE visual in context.
	Interactive string `json:"interactive,omitempty"`
}

// eventRenderable are directly-rendered elements where an event handler that
// changes paint/opacity is visible. (Filter primitives, defs, gradients,
// animation elements render nothing of their own, so events are pointless there.)
var eventRenderable = map[string]bool{
	"rect": true, "circle": true, "ellipse": true, "line": true,
	"polyline": true, "polygon": true, "path": true,
	"text": true, "tspan": true, "a": true, "g": true, "switch": true,
	"image": true, "foreignObject": true, "use": true,
}

// interactivePresets adds showcase presets that demonstrate event ATTRIBUTES are
// visual once they act on other attributes: on hover the handler fades the
// element, on click it recolors it. (The grammar enumerates ~50 event attrs all
// sampled to the meaningless StringType "label"; these select the demonstrable
// ones with appearance-changing handlers — prov="curated".)
func interactivePresets(tag string) []catPreset {
	if !eventRenderable[tag] {
		return nil
	}
	// pointer-events:bounding-box makes the WHOLE bounding box a hit target, so the
	// hover/click fires even on a container whose centre is empty (g/a/switch). The
	// handlers touch fill+color+stroke so they recolor direct-fill shapes AND
	// currentColor-filled container children.
	return []catPreset{
		{
			Name: "onmouseover", Prov: "curated", Interactive: "hover",
			Meaning: "ON HOVER the handler fades the element (event attributes are visual in context)",
			Values: map[string]string{
				"onmouseover":    "this.style.opacity='0.3'",
				"onmouseout":     "this.style.opacity='1'",
				"pointer-events": "bounding-box",
			},
		},
		{
			Name: "onclick", Prov: "curated", Interactive: "click",
			Meaning: "ON CLICK the handler recolors the element to orange",
			Values: map[string]string{
				"onclick":        "this.style.fill='#f5a623';this.style.color='#f5a623';this.style.stroke='#f5a623'",
				"pointer-events": "bounding-box",
			},
		},
	}
}

type catElement struct {
	ID       string            `json:"id"`
	Tag      string            `json:"tag"`
	Name     string            `json:"name"`
	Cat      string            `json:"cat"`
	Desc     string            `json:"desc"`
	Temporal bool              `json:"temporal,omitempty"`
	Attrs    []catControl      `json:"attrs"`
	Defaults map[string]string `json:"defaults"`
	Base     string            `json:"base"`
	Presets  []catPreset       `json:"presets"`
}

type catalogue struct {
	Groups   []string     `json:"groups"`
	Elements []catElement `json:"elements"`
}

// runCataloguePass builds and writes catalogue.json from the enumerated pages.
func runCataloguePass(en *Enumerator, bp *blueprintProvider, els []element, pages []page) (int, int) {
	byTag := map[string]element{}
	for _, el := range els {
		byTag[el.tag] = el
	}

	var out catalogue
	out.Groups = groupOrder
	seenGroup := map[string]bool{}

	for _, p := range pages {
		el := byTag[p.tag]

		// base render: the element at baseline (no varied attribute). Route the
		// grammar element through the codec (renderer of record) BEFORE layering on
		// harness annotations (id="slot"/data-lab) and the scaffold — so what the
		// gallery shows is codec-emitted, while the harness chrome stays as-is.
		baseMarkup, needsID := en.renderWithOneAttr(el, "", "", "")
		baseMarkup = checkPath(p.tag, "", "", baseMarkup)
		if needsID || blueprintSlotNeedsID(bp.blueprintFor(p.tag), p.tag) {
			baseMarkup = ensureSlotID(baseMarkup, p.tag)
		}
		baseMarkup = markLabSlot(baseMarkup, p.tag)
		base := inject(bp.blueprintFor(p.tag), baseMarkup)
		baseAttrs := parseOpenTagAttrs(base, p.tag)

		// group the VISUAL value-paths by attribute, preserving first-seen order.
		var order []string
		vals := map[string][]string{}
		for _, v := range p.variants {
			if nonVisualAttr(v.Attr, p.tag) {
				continue
			}
			if _, ok := vals[v.Attr]; !ok {
				order = append(order, v.Attr)
			}
			vals[v.Attr] = append(vals[v.Attr], v.Value)
		}
		if len(order) == 0 {
			continue // nothing visual to show
		}

		ce := catElement{
			ID: p.tag, Tag: p.tag, Name: nameFor(p.tag), Cat: groupFor(p.tag),
			Desc: descFor(p.tag), Defaults: map[string]string{},
		}
		ce.Temporal = animationTags[p.tag] || liveSMILRe.MatchString(base)
		seenGroup[ce.Cat] = true

		for _, attr := range order {
			if !controlWorthy(attr) {
				continue // still a preset below, just not a cluttering control
			}
			c := inferControl(attr, p.tag, dedupStr(vals[attr]))
			ce.Attrs = append(ce.Attrs, c)
			// Default ONLY from the baseline render (which is tuned for visibility).
			// Attrs absent from the baseline are left unset so the element keeps its
			// natural initial value (e.g. display stays `inline`, not `none`).
			if d, ok := baseAttrs[attr]; ok && d != "" {
				ce.Defaults[attr] = d
			}
		}
		// presets: one per visual value-path. The OVERLAY samples leaf values
		// (refs→real defs, magnitudes→visible) so they're grammar-faithful; curate
		// adds showcase EXPANSIONS, captions, and display pruning. Each preset's
		// Values also gets companionFor's companions/overrides MERGED IN, so a varied
		// attr that only acts in context (k1→operator=arithmetic, feFunc amplitude→
		// type=gamma) carries that context into the gallery (which applies Values to
		// the base, not a per-preset specimen).
		usesDemoDefs := false
		for _, attr := range order {
			raw := dedupStr(vals[attr])
			curated := curateAttr(p.tag, attr, raw)
			prov := "curated"
			if curated == nil {
				prov = "grammar"
				for _, val := range raw {
					curated = append(curated, demoPreset{
						label:   val,
						values:  map[string]string{attr: val},
						meaning: meaningFor(p.tag, attr, val),
					})
				}
			}
			for _, dp := range curated {
				vals := withCompanions(p.tag, attr, dp.label, dp.values)
				for _, v := range vals {
					if strings.Contains(v, "#fx-") {
						usesDemoDefs = true
					}
				}
				ce.Presets = append(ce.Presets, catPreset{
					Name:    presetName(attr, dp.label),
					Values:  vals,
					Meaning: dp.meaning,
					Prov:    prov,
				})
			}
		}
		// interactivity: event-attribute presets demonstrated under real mouse input.
		ce.Presets = append(ce.Presets, interactivePresets(p.tag)...)
		if usesDemoDefs {
			base = injectDemoDefs(base)
		}
		ce.Base = base
		out.Elements = append(out.Elements, ce)
	}

	// keep only groups that actually have elements, in canonical order.
	var groups []string
	for _, g := range groupOrder {
		if seenGroup[g] {
			groups = append(groups, g)
		}
	}
	out.Groups = groups

	data, _ := json.MarshalIndent(out, "", " ")
	writeFile(filepath.Join(catalogueDir, "catalogue.json"), string(data)+"\n")

	temporal := 0
	for _, e := range out.Elements {
		if e.Temporal {
			temporal++
		}
	}
	return len(out.Elements), temporal
}

// ---- control inference ----

var paintAttrs = map[string]bool{
	"fill": true, "stroke": true, "stop-color": true, "flood-color": true,
	"lighting-color": true, "color": true, "solid-color": true,
}

// kwRe matches a single keyword OR a space-separated compound keyword
// ("xMidYMid meet"); numStartRe matches anything beginning like a number.
var kwRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9 _-]*$`)
var numStartRe = regexp.MustCompile(`^[-+]?[.0-9]`)

func inferControl(attr, tag string, values []string) catControl {
	c := catControl{Key: attr, Label: labelFor(attr)}
	if paintAttrs[attr] {
		c.Control = "paint"
		return c
	}
	var kw, singleNum, listNum []string
	for _, v := range values {
		switch {
		case numStartRe.MatchString(v):
			if strings.Contains(v, " ") {
				listNum = append(listNum, v)
			} else {
				singleNum = append(singleNum, v)
			}
		case kwRe.MatchString(v):
			kw = append(kw, v)
		}
	}
	num := append(append([]string{}, singleNum...), listNum...)
	switch {
	case forceRange[attr] && len(num) > 0:
		// number-optional-number filter attrs (stdDeviation, baseFrequency, …)
		// render as 1-or-2-number lists but read best as a slider on the primary
		// value; the "x y" forms stay available as presets.
		c.Control = "range"
		mn, mx, st := rangeFor(attr, tag, num)
		c.Min, c.Max, c.Step = &mn, &mx, &st
	case len(singleNum) > 0:
		c.Control = "range"
		mn, mx, st := rangeFor(attr, tag, singleNum)
		c.Min, c.Max, c.Step = &mn, &mx, &st
	case len(listNum) > 0:
		c.Control = "text" // multi-number lists (text x/rotate, points, …) — a slider can't hold them
	case len(kw) > 0:
		c.Control = "select"
		c.Options = dedupStr(kw)
	default:
		c.Control = "text"
	}
	return c
}

// forceRange are number-optional-number attributes whose value-paths are 1-or-2
// number lists but whose primary value reads best as a slider.
var forceRange = map[string]bool{
	"stdDeviation": true, "baseFrequency": true, "radius": true,
	"order": true, "kernelUnitLength": true,
}

// primaryPresentation is the small set of presentation attributes worth a live
// CONTROL (paint, stroke, opacity, dash, transform). Every other presentation
// attribute (cursor, color-rendering, clip, pointer-events, …) is dropped from
// the control panel — it still appears as a PRESET so its values are catalogued
// and screenshot, but it doesn't clutter the interactive controls.
var primaryPresentation = map[string]bool{
	"fill": true, "stroke": true, "stroke-width": true, "stroke-opacity": true,
	"fill-opacity": true, "opacity": true, "stroke-dasharray": true, "stroke-dashoffset": true,
	"stroke-linecap": true, "stroke-linejoin": true, "transform": true, "paint-order": true,
	"stop-color": true, "stop-opacity": true, "flood-color": true, "flood-opacity": true,
	"lighting-color": true, "color": true, "offset": true,
}

// controlWorthy reports whether an attribute should get an interactive control:
// element-specific attributes always do; shared presentation attributes only if
// they are in the primary set.
func controlWorthy(attr string) bool {
	if presentationAttrs[attr] {
		return primaryPresentation[attr]
	}
	return true
}

// rangeMeta is a curated slider range for common numeric attributes (the user
// sanctioned presetting how each attribute is displayed). Anything not listed is
// derived from the actual value-paths.
type rangeMeta struct{ mn, mx, st float64 }

var rangeTable = map[string]rangeMeta{
	"opacity": {0, 1, 0.05}, "fill-opacity": {0, 1, 0.05}, "stroke-opacity": {0, 1, 0.05},
	"flood-opacity": {0, 1, 0.05}, "stop-opacity": {0, 1, 0.05},
	"stroke-width": {0, 30, 1}, "stroke-dashoffset": {-40, 40, 1}, "stroke-miterlimit": {1, 20, 0.5},
	"font-size": {8, 120, 1}, "letter-spacing": {-4, 24, 1}, "word-spacing": {-4, 24, 1},
	"r": {0, 120, 1}, "rx": {0, 120, 1}, "ry": {0, 120, 1}, "fr": {0, 120, 1},
	"width": {0, 240, 1}, "height": {0, 240, 1},
	"x": {0, 240, 1}, "y": {0, 240, 1}, "cx": {0, 240, 1}, "cy": {0, 240, 1},
	"x1": {0, 240, 1}, "y1": {0, 240, 1}, "x2": {0, 240, 1}, "y2": {0, 240, 1},
	"dx": {-40, 40, 1}, "dy": {-40, 40, 1},
	"stdDeviation": {0, 24, 0.5}, "baseFrequency": {0, 0.4, 0.01}, "numOctaves": {1, 8, 1},
	"scale": {0, 80, 1}, "surfaceScale": {0, 20, 0.5}, "specularExponent": {1, 50, 1},
	"diffuseConstant": {0, 4, 0.1}, "specularConstant": {0, 4, 0.1},
	"azimuth": {0, 360, 5}, "elevation": {0, 180, 5}, "pathLength": {0, 400, 5},
	"offset": {0, 1, 0.05}, "amplitude": {0, 4, 0.1}, "exponent": {0, 4, 0.1},
	"slope": {0, 4, 0.1}, "intercept": {-1, 1, 0.05}, "k1": {0, 2, 0.1}, "k2": {0, 2, 0.1},
	"k3": {0, 2, 0.1}, "k4": {0, 2, 0.1}, "divisor": {1, 20, 1}, "bias": {-1, 1, 0.05},
}

func rangeFor(attr, tag string, num []string) (float64, float64, float64) {
	if r, ok := rangeTable[attr]; ok {
		return r.mn, r.mx, r.st
	}
	// derive from the numeric value-paths.
	mn, mx := 0.0, 0.0
	first := true
	for _, v := range num {
		f, ok := parseLeadingNumber(v)
		if !ok {
			continue
		}
		if first {
			mn, mx, first = f, f, false
		}
		if f < mn {
			mn = f
		}
		if f > mx {
			mx = f
		}
	}
	if first { // no parseable numbers
		return 0, 100, 1
	}
	if mn > 0 {
		mn = 0
	}
	if mx <= mn {
		mx = mn + 100
	}
	pad := (mx - mn) * 0.4
	mx += pad
	st := 1.0
	switch {
	case mx <= 1:
		st = 0.05
	case mx <= 10:
		st = 0.5
	}
	return mn, mx, st
}

// ---- helpers ----

var openTagRe = func(tag string) *regexp.Regexp {
	return regexp.MustCompile(`<` + regexp.QuoteMeta(tag) + `(\s[^>]*?)?/?>`)
}
var attrPairRe = regexp.MustCompile(`([\w:.-]+)="([^"]*)"`)

// parseOpenTagAttrs returns the attribute map of the first <tag ...> in svg.
func parseOpenTagAttrs(svg, tag string) map[string]string {
	m := openTagRe(tag).FindString(svg)
	out := map[string]string{}
	for _, p := range attrPairRe.FindAllStringSubmatch(m, -1) {
		out[p[1]] = p[2]
	}
	return out
}

func parseLeadingNumber(s string) (float64, bool) {
	re := regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+`)
	m := re.FindString(strings.TrimSpace(s))
	if m == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(m, 64)
	return f, err == nil
}

func dedupStr(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

func presetName(attr, val string) string {
	if len(val) > 16 {
		val = val[:15] + "…"
	}
	return attr + "=" + val
}

func labelFor(attr string) string { return attr }

// withCompanions merges companionFor's companion attributes and overrides into a
// preset's Values. The gallery applies Values to the base (not a per-preset
// specimen), so a varied attr that only acts in context must carry that context
// itself — k1→operator=arithmetic, feFunc amplitude→type=gamma, etc.
func withCompanions(tag, attr, value string, base map[string]string) map[string]string {
	comp, over := companionFor(tag, attr, value)
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}
	for _, c := range comp {
		if _, ok := out[c[0]]; !ok {
			out[c[0]] = c[1]
		}
	}
	for k, v := range over {
		out[k] = v
	}
	return out
}

// markLabSlot tags the showcased element's open tag with data-lab so the gallery
// (buildCode) mutates exactly THIS element. The base may contain other elements
// of the same tag — a <rect> inside the demo mask, a <path> inside the demo
// marker, the faint original <use> — which a bare first-match would hit instead.
func markLabSlot(markup, tag string) string {
	open := "<" + tag
	i := strings.Index(markup, open)
	if i < 0 {
		return markup
	}
	j := i + len(open)
	if j < len(markup) { // guard <line vs <linearGradient prefix overlap
		switch markup[j] {
		case ' ', '\t', '\n', '/', '>':
		default:
			return markup
		}
	}
	return markup[:j] + ` data-lab=""` + markup[j:]
}
