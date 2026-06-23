package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// specimen.go — the per-value SPECIMEN emitter (proto-css-style per-value
// screenshot pipeline). For every MAIN-GRID value-path of every element it
// writes a standalone, label-free dark page containing exactly ONE centered SVG
// — the variant's blueprint-wrapped render, scaled to fill a ~360×360 frame — so
// a screenshot of the page is purely that one value's render. Alongside the
// files it writes a manifest (specimens.json) describing each specimen, flagging
// the ones that need multi-frame (GIF) capture because they animate.
//
// "MAIN-GRID" = the same partition emitPage uses: variants where
// nonVisualAttr(attr, tag) is FALSE. The no-visual-effect attributes (id, aria,
// on* events, presentation attrs on non-painting elements, …) are documented to
// render identically, so they are skipped here — there is nothing to shoot.

// specimenHTMLDir is the root of the per-value specimen HTML tree
// (<tag>/<NN>-<slug>.html). It mirrors the gallery dir layout but is a sibling
// "specimen/" subtree so the shoot pipeline can walk it independently.
const specimenHTMLDir = "chrome-testing/html/specimen"

// specimenJSONDir is where specimens.json (the shoot manifest) is written. It is
// the engine's --out directory (chrome-testing/generated), NOT the gallery HTML
// dir — specimens.json is build metadata consumed by the screenshot pipeline,
// kept next to sample-document.svg / sample-rect.svg.
const specimenJSONDir = "chrome-testing/generated"

// animationTags is the set of SVG animation elements. When the ELEMENT itself is
// one of these, its render is temporal (it animates) and needs GIF capture.
var animationTags = map[string]bool{
	"animate": true, "set": true, "animateMotion": true,
	"animateTransform": true, "mpath": true, "discard": true,
}

// liveSMILRe detects a live SMIL animation element anywhere in the wrapped SVG
// (e.g. an animation injected by the blueprint scaffold around a static
// element). Such a specimen is temporal even when the varied element is static.
var liveSMILRe = regexp.MustCompile(`<(animate|animateTransform|animateMotion|set)[\s>/]`)

// specimenEntry is one element's manifest record for one main-grid value-path.
type specimenEntry struct {
	I        int    `json:"i"`        // main-grid index (matches the NN in the filename)
	Label    string `json:"label"`    // the card label, e.g. `x="10"`
	Value    string `json:"value"`    // the raw value text
	File     string `json:"file"`     // path relative to the html dir, e.g. specimen/rect/00-10.html
	Temporal bool   `json:"temporal"` // true → needs multi-frame GIF capture
}

// runSpecimenPass writes the per-value specimen HTML files and specimens.json.
// It consumes the SAME pages the gallery pass built (variants already carry
// their blueprint-wrapped WrappedSVG), partitions each into main-grid value-
// paths (nonVisualAttr == false), and emits one standalone page per main-grid
// variant plus the manifest. It returns (files written, temporal count).
func runSpecimenPass(pages []page) (written, temporal int) {
	manifest := map[string][]specimenEntry{}

	// Clean the specimen tree so values removed by a grammar/generator change
	// (e.g. an attribute reclassified non-visual) don't linger as stale files.
	os.RemoveAll(specimenHTMLDir)

	for _, p := range pages {
		var entries []specimenEntry
		i := 0
		for _, v := range p.variants {
			// Skip the documented-identical no-effect attributes — only the
			// main-grid (visually-meaningful) value-paths get a specimen.
			if nonVisualAttr(v.Attr, p.tag) {
				continue
			}
			slug := slugify(v.Value)
			rel := filepath.Join("specimen", p.tag, fmt.Sprintf("%02d-%s.html", i, slug))
			temporalV := isTemporal(p.tag, v.WrappedSVG)

			writeFile(filepath.Join(specimenHTMLDir, p.tag, fmt.Sprintf("%02d-%s.html", i, slug)),
				specimenHTML(p.tag, v))

			entries = append(entries, specimenEntry{
				I:        i,
				Label:    specimenLabel(v),
				Value:    v.Value,
				File:     rel,
				Temporal: temporalV,
			})
			written++
			if temporalV {
				temporal++
			}
			i++
		}
		if len(entries) > 0 {
			manifest[p.tag] = entries
		}
	}

	data, _ := json.MarshalIndent(manifest, "", "  ")
	writeFile(filepath.Join(specimenJSONDir, "specimens.json"), string(data)+"\n")
	return written, temporal
}

// isTemporal reports whether a specimen needs multi-frame GIF capture: either
// the element itself is an animation element, or the wrapped SVG contains a live
// SMIL element (an animation the blueprint scaffold injected around it).
func isTemporal(tag, wrappedSVG string) bool {
	if animationTags[tag] {
		return true
	}
	return liveSMILRe.MatchString(wrappedSVG)
}

// specimenLabel renders the plain-text attr="value" caption used in the
// manifest (mirrors emit.go's cardLabel but without HTML markup).
func specimenLabel(v Variant) string {
	return fmt.Sprintf(`%s="%s"`, v.Attr, v.Value)
}

// specimenHTML builds a standalone, minimal dark page containing exactly ONE
// centered SVG: the variant's blueprint-wrapped render, scaled to fill a
// ~360×360 frame. No labels, no cards — a screenshot is purely the render.
func specimenHTML(tag string, v Variant) string {
	// The wrapped SVG has a viewBox, so overriding its width/height via CSS
	// scales the whole render to fill the frame while preserving aspect ratio.
	return `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>` +
		tag + ` specimen</title>
<style>body{margin:0;background:#1a1a2e;display:grid;place-items:center;height:100vh}
svg{width:360px;height:360px;display:block}</style></head><body>
` + v.WrappedSVG + `
</body></html>
`
}

// slugRe matches runs of characters that are not lowercase ASCII alphanumerics.
var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

// slugify converts an attribute value into a filesystem-safe slug: lowercase,
// every non-[a-z0-9] run collapsed to "-", trimmed of leading/trailing "-",
// truncated to 40 chars, and "v" when the result is empty.
func slugify(value string) string {
	s := strings.ToLower(value)
	s = slugRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 40 {
		s = s[:40]
		s = strings.Trim(s, "-")
	}
	if s == "" {
		s = "v"
	}
	return s
}
