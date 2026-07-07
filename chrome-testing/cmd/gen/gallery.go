package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// gallery.go — the enumerate → inject → emit driver. It runs the full
// all-value-paths walk over every element, wraps each variant in the element's
// blueprint (assigning id="slot" when the blueprint references it), and writes
// the gallery pages, index, values.json and manifest.tsv.

const templateHTMLDir = "chrome-testing/html/template"

func runGalleryPass(byFQN map[string]*descriptorpb.DescriptorProto, kw map[string]string, r *Renderer) {
	en := newEnumerator(byFQN, kw, r)
	bp := newBlueprintProvider(templateHTMLDir)

	els := en.allElements()
	var pages []page
	total := 0
	var unrendered []string

	for _, el := range els {
		// Skip elements pruned from the shipped grammar (present only in the gen's
		// no-strip grammar, e.g. the non-rendering SVGUnknownElement): the codec
		// can't render them and they aren't real gallery elements.
		if !inShippedGrammar(el.tag) {
			unrendered = append(unrendered, "<"+el.tag+"> (not in shipped grammar — skipped)")
			continue
		}
		variants := en.enumerateElement(el)
		if len(variants) == 0 {
			unrendered = append(unrendered, "<"+el.tag+"> (no enumerable attributes)")
			continue
		}
		blueprint := bp.blueprintFor(el.tag)
		bpRefsSlot := strings.Contains(blueprint, "url(#"+slotID+")") ||
			strings.Contains(blueprint, `href="#`+slotID+`"`) ||
			strings.Contains(blueprint, `marker-end="url(#`+slotID) ||
			strings.Contains(blueprint, `clip-path="url(#`+slotID) ||
			strings.Contains(blueprint, `mask="url(#`+slotID)
		// The blueprint may also place {{ELEMENT}} directly inside a referenced
		// container (filter/gradient with id="slot" already on the wrapper). In
		// those scaffolds {{ELEMENT}} is the def itself and must carry id="slot".
		slotIsElement := blueprintSlotNeedsID(blueprint, el.tag)

		for i := range variants {
			markup := variants[i].Markup
			// Round-trip every walked path through the codec (the renderer of record):
			// Parse structures the markup into the shipped AST (Any seams and all),
			// Render re-emits it. Records any path the generators don't handle exactly.
			checkPath(el.tag, variants[i].Attr, variants[i].Value, markup)
			// id-shadow guard: when the VARIED attribute is `id` itself, the markup
			// already carries that id; injecting id="slot" too would produce an
			// invalid duplicate (`id="slot" id="circle1"`). Skip the injected id for
			// that one card (its url(#slot)/href reference simply won't resolve — an
			// acceptable edge for the id-variation card; all other cards are fine).
			injectSlot := (bpRefsSlot && slotIsElement) || variants[i].NeedsID
			if variants[i].Attr == "id" {
				injectSlot = false
			}
			if injectSlot {
				markup = ensureSlotID(markup, el.tag)
			}
			variants[i].WrappedSVG = inject(blueprint, markup)

			if ok, msg := checkWellFormed(variants[i].WrappedSVG); !ok {
				unrendered = append(unrendered,
					fmt.Sprintf("<%s> %s=%q: not well-formed (%s)", el.tag, variants[i].Attr, variants[i].Value, msg))
			}
		}
		pages = append(pages, page{tag: el.tag, variants: variants})
		total += len(variants)
	}

	// The gallery is now the SVG Lab app driven by a single catalogue.json (the
	// per-element static HTML + specimen technique is retired). Emit the catalogue:
	// per element, typed attribute controls + presets (one per visual value-path)
	// + a base SVG.
	catEls, catTemporal := runCataloguePass(en, bp, els, pages)

	fmt.Printf("\n=== catalogue ===\n")
	fmt.Printf("elements: %d / %d  variants: %d  (temporal: %d)  ->  %s/catalogue.json\n",
		catEls, len(pages), total, catTemporal, catalogueDir)

	// The codec is the renderer of record: report every walked path it did not
	// round-trip faithfully (both Parse and Render exercised per path).
	codecReport(total)
	if len(unrendered) > 0 {
		fmt.Printf("value-paths that did not render meaningfully (%d):\n", len(unrendered))
		max := len(unrendered)
		if max > 20 {
			max = 20
		}
		for _, u := range unrendered[:max] {
			fmt.Printf("  - %s\n", u)
		}
		if len(unrendered) > max {
			fmt.Printf("  ... and %d more\n", len(unrendered)-max)
		}
	}
}

// blueprintSlotNeedsID reports whether the {{ELEMENT}} placeholder in the
// blueprint sits where the referenced def lives (so the injected element must
// carry id="slot"). True for paint-server/filter/clip/mask scaffolds where the
// element IS the referenced def; false for self-rendering scaffolds and those
// where the element is nested inside a wrapper that already owns id="slot"
// (filter primitives inside <filter id="slot">, stops inside a gradient, etc.).
func blueprintSlotNeedsID(blueprint, tag string) bool {
	switch category(tag) {
	case catGradient, catPattern, catMarker, catClip, catMask, catFilter:
		return true
	}
	// symbol is catSelf but its blueprint references the injected element via
	// <use href="#slot">, so the <symbol> itself must carry id="slot". (defs is
	// NOT here: its scaffold references a CHILD def carrying id="slot" — bodyFor —
	// not the <defs> wrapper, so no id is injected onto <defs> itself.)
	if tag == "symbol" {
		return true
	}
	return false
}

// ensureSlotID inserts id="slot" into the element's open tag if it is not
// already present. It targets the FIRST occurrence of the element's open tag.
func ensureSlotID(markup, tag string) string {
	if strings.Contains(markup, `id="`+slotID+`"`) {
		return markup
	}
	open := "<" + tag
	idx := strings.Index(markup, open)
	if idx < 0 {
		return markup
	}
	insert := idx + len(open)
	return markup[:insert] + ` id="` + slotID + `"` + markup[insert:]
}
