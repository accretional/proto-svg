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
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
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

		// base render: the element at baseline (no varied attribute).
		baseMarkup, needsID := en.renderWithOneAttr(el, "", "", "")
		if needsID || blueprintSlotNeedsID(bp.blueprintFor(p.tag), p.tag) {
			baseMarkup = ensureSlotID(baseMarkup, p.tag)
		}
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
			c := inferControl(attr, p.tag, dedupStr(vals[attr]))
			ce.Attrs = append(ce.Attrs, c)
			// Default ONLY from the baseline render (which is tuned for visibility).
			// Attrs absent from the baseline are left unset so the element keeps its
			// natural initial value (e.g. display stays `inline`, not `none`).
			if d, ok := baseAttrs[attr]; ok && d != "" {
				ce.Defaults[attr] = d
			}
		}
		// presets: one per visual value-path.
		for _, attr := range order {
			for _, val := range dedupStr(vals[attr]) {
				ce.Presets = append(ce.Presets, catPreset{
					Name:   presetName(attr, val),
					Values: map[string]string{attr: val},
				})
			}
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

var kwRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
var numStartRe = regexp.MustCompile(`^[-+]?[.0-9]`)

func inferControl(attr, tag string, values []string) catControl {
	c := catControl{Key: attr, Label: labelFor(attr)}
	if paintAttrs[attr] {
		c.Control = "paint"
		return c
	}
	var kw, num []string
	for _, v := range values {
		if numStartRe.MatchString(v) {
			num = append(num, v)
		} else if kwRe.MatchString(v) {
			kw = append(kw, v)
		}
	}
	if len(num) > 0 {
		c.Control = "range"
		mn, mx, st := rangeFor(attr, tag, num)
		c.Min, c.Max, c.Step = &mn, &mx, &st
		return c
	}
	if len(kw) > 0 {
		c.Control = "select"
		c.Options = dedupStr(kw)
		return c
	}
	c.Control = "text"
	return c
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
