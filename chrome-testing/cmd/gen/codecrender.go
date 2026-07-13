package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	svgpb "github.com/accretional/proto-svg/proto/pb/svg"
	svgservicepb "github.com/accretional/proto-svg/proto/pb/svgservice"
	"github.com/accretional/proto-svg/service"

	// Link the embedded grammars so the codec structures and renders the seams:
	// css for paint colors (fill / stroke / stop-color / …), html for the
	// <foreignObject> flow-content seam (the svg→html cycle).
	_ "github.com/accretional/proto-css/service"
	_ "github.com/accretional/proto-html/service"
)

// codecrender.go — make the gluon codec the RENDERER OF RECORD for the gallery.
//
// The enumerator walks every grammar path and emits the element as MARKUP TEXT
// (its overlay/reps knowledge is string-valued). We do not re-implement that as a
// structural builder; instead each walked path's text is handed to codec.Parse —
// which structures it into the shipped AST, Any seams and all — and codec.Render,
// which re-emits it. Parse and Render are the two generators under test, so every
// gallery specimen exercises both. This gen owns the WALK; the codec owns the
// AST↔text boundary.

// tagRoot maps an element tag ("rect") to its shipped root message type name
// ("svg.SvgrectElement"), read out of the grammar's prefix table.
var tagRoot = buildTagRoot()

func buildTagRoot() map[string]string {
	m := map[string]string{}
	for fqn, pfx := range svgpb.MessagePrefix {
		if len(pfx) != 1 {
			continue
		}
		lit := pfx[0]
		if !strings.HasPrefix(lit, "<") || strings.HasPrefix(lit, "</") {
			continue
		}
		name := strings.TrimPrefix(fqn, ".")
		// The prefix table carries GHOST entries: dead keyword rules (e.g.
		// LessThanSignUseKeyword → "<use") share an element's open-tag prefix but
		// are not real messages. Resolve each tag to the real, registered ELEMENT
		// message, never a ghost — otherwise ParseAs("svg.LessThanSign…Keyword")
		// fails "not found".
		if strings.Contains(name, "Keyword") {
			continue
		}
		if _, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(name)); err != nil {
			continue
		}
		m[tagName(lit)] = name
	}
	return m
}

// svgService is the SvgService the gallery renders through — the same server
// service/cmd/server exposes over gRPC, driven in-process here so the gen
// exercises the exact Parse/Render surface clients see.
var svgService = service.NewServer()

// codecRender routes grammar-element markup through the SvgService (Parse with
// the element's root type → Render of the Any-packed node).
// On success it returns the codec's rendering — byte-identical to the input when
// the round-trip is faithful. On any failure it returns the ORIGINAL markup plus
// the error, so the caller can fall back and record the gap rather than abort.
func codecRender(tag, markup string) (string, error) {
	root := tagRoot[tag]
	if root == "" {
		return markup, fmt.Errorf("no shipped root type for <%s>", tag)
	}
	ctx := context.Background()
	parsed, err := svgService.Parse(ctx, &svgservicepb.ParseRequest{Svg: markup, Type: root})
	if err != nil {
		return markup, fmt.Errorf("parse: %w", err)
	}
	rendered, err := svgService.Render(ctx, &svgservicepb.RenderRequest{Node: parsed.GetNode()})
	if err != nil {
		return markup, fmt.Errorf("render: %w", err)
	}
	return rendered.GetSvg(), nil
}

// ── per-path codec report ─────────────────────────────────────────────────────
// Every enumerated path is round-tripped through the codec; mismatches and errors
// are collected so a gen run reports exactly which grammar paths the generators do
// not yet handle faithfully.

type codecFailure struct {
	tag, attr, value, kind, detail, in, out string
}

var codecFailures []codecFailure

func recordCodecFailure(tag, attr, value, kind, detail, in, out string) {
	codecFailures = append(codecFailures, codecFailure{tag, attr, value, kind, detail, in, out})
}

// inShippedGrammar reports whether tag resolves to a real element type in the
// shipped grammar. Elements present only in the gen's no-strip grammar but pruned
// as unreachable from the shipped one (e.g. SVGUnknownElement, a non-rendering
// fallback) are not gallery-able and are skipped.
func inShippedGrammar(tag string) bool { return tagRoot[tag] != "" }

// checkPath round-trips one walked path's markup through the codec and records a
// failure if it errors or is not byte-exact. Returns the codec rendering (or the
// input on failure). Elements absent from the shipped grammar are skipped (not
// recorded as failures). Self-closing child tags are normalized to explicit-close
// first, since SVG has no void elements and the grammar models explicit close.
func checkPath(tag, attr, value, markup string) string {
	if !inShippedGrammar(tag) {
		return markup
	}
	markup = expandSelfClosing(markup)
	out, err := codecRender(tag, markup)
	switch {
	case err != nil:
		recordCodecFailure(tag, attr, value, "error", err.Error(), markup, out)
	case out != markup:
		recordCodecFailure(tag, attr, value, "mismatch", "round-trip not byte-exact", markup, out)
	}
	return out
}

// expandSelfClosing rewrites every self-closing start tag `<name …/>` to the
// explicit-close form `<name …></name>`, quote-aware so a `/` or `>` inside an
// attribute value is not mistaken for the tag end. SVG has no void elements, so
// this preserves meaning; it aligns the gen's body markup with the grammar's
// explicit-close model (the codec's canonical form).
func expandSelfClosing(markup string) string {
	isNameStart := func(c byte) bool {
		return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
	}
	isNameChar := func(c byte) bool {
		return isNameStart(c) || (c >= '0' && c <= '9') || c == '-' || c == ':'
	}
	var b strings.Builder
	i, n := 0, len(markup)
	for i < n {
		c := markup[i]
		if c != '<' || i+1 >= n || !isNameStart(markup[i+1]) {
			b.WriteByte(c)
			i++
			continue
		}
		j := i + 1
		for j < n && isNameChar(markup[j]) {
			j++
		}
		name := markup[i+1 : j]
		k := j
		inQuote := false
		for k < n {
			switch markup[k] {
			case '"':
				inQuote = !inQuote
			case '>':
				if !inQuote {
					goto found
				}
			}
			k++
		}
	found:
		if k >= n { // malformed; emit the rest verbatim
			b.WriteString(markup[i:])
			break
		}
		if markup[k-1] == '/' { // self-closing
			open := strings.TrimRight(markup[i:k-1], " ")
			b.WriteString(open)
			b.WriteString("></")
			b.WriteString(name)
			b.WriteString(">")
		} else {
			b.WriteString(markup[i : k+1])
		}
		i = k + 1
	}
	return b.String()
}

// codecFailuresPath is where the FULL failure list is dumped (the console summary
// caps per tag). One line per failure: tag<TAB>kind<TAB>attr<TAB>value<TAB>detail.
const codecFailuresPath = "chrome-testing/generated/_codec_failures.tsv"

// dumpCodecFailures writes every failure to codecFailuresPath for offline analysis
// (grouped console output truncates to a few per tag).
func dumpCodecFailures() {
	var b strings.Builder
	b.WriteString("tag\tkind\tattr\tvalue\tdetail\tin\tout\n")
	for _, f := range codecFailures {
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", f.tag, f.kind, f.attr, f.value, f.detail, f.in, f.out)
	}
	writeFile(codecFailuresPath, b.String())
}

// codecReport prints a summary of every grammar path the codec did not round-trip
// faithfully, grouped by element, and returns the failure count.
func codecReport(total int) int {
	dumpCodecFailures()
	if len(codecFailures) == 0 {
		fmt.Printf("codec round-trip: all %d walked paths render faithfully ✓\n", total)
		return 0
	}
	byTag := map[string][]codecFailure{}
	var tags []string
	for _, f := range codecFailures {
		if _, ok := byTag[f.tag]; !ok {
			tags = append(tags, f.tag)
		}
		byTag[f.tag] = append(byTag[f.tag], f)
	}
	sort.Strings(tags)
	fmt.Printf("\ncodec round-trip: %d/%d walked paths FAILED:\n", len(codecFailures), total)
	for _, t := range tags {
		fs := byTag[t]
		fmt.Printf("  <%s> (%d):\n", t, len(fs))
		for i, f := range fs {
			if i >= 4 {
				fmt.Printf("      … and %d more\n", len(fs)-4)
				break
			}
			fmt.Printf("      [%s] %s=%q: %s\n", f.kind, f.attr, f.value, f.detail)
			if f.kind == "mismatch" {
				fmt.Printf("         in : %s\n         out: %s\n", f.in, f.out)
			}
		}
	}
	return len(codecFailures)
}
