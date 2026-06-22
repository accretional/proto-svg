package main

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// Renderer walks the compiled SVG grammar (a FileDescriptorProto message graph,
// indexed by fully-qualified message name) and emits valid SVG markup.
//
// The defining difference from the proto-css reference renderer is that
// REPEATED fields are RENDERED, not dropped. In SVG the EBNF `{ }` repetition
// IS the document structure: a `{ SvgAttribute }` is the element's attribute
// list and a `{ ContainerContent }` is its child-element list. Omitting them —
// as proto-css does — would erase every attribute and child and leave a
// structurally-empty document. So renderField emits a repeated field's message
// type childCap times (rotating oneof arms each iteration) and concatenates.
//
// Dispatch order for render(fqn):
//  1. depth > maxDepth        → ""              (depth guard)
//  2. reps[simpleName] exists → one rotating sample (open-ended leaf)
//  3. kw[fqn] exists          → kw[fqn]         (markup token / enum keyword)
//  4. byFQN[fqn] == nil       → ""              (unresolved)
//  5. seen[fqn] >= 2          → ""              (cycle guard)
//  6. message has a oneof     → pick ONE arm, rotating across calls
//  7. plain message           → concat renderField over fields in order
type Renderer struct {
	byFQN    map[string]*descriptorpb.DescriptorProto // ".svg.Foo" → message
	kw       map[string]string                        // FQN → markup/keyword literal
	optional map[string]bool                          // "<msgFQN>/<field>" → EBNF-optional

	maxDepth int // recursion cap (14)
	childCap int // repetitions emitted per repeated field (2)

	ctr      map[string]int // per-leaf rotation counter (reps variety)
	oneofCtr map[string]int // per-oneof rotation counter (arm variety)
}

func newRenderer(byFQN map[string]*descriptorpb.DescriptorProto, kw map[string]string, optional map[string]bool) *Renderer {
	return &Renderer{
		byFQN:    byFQN,
		kw:       kw,
		optional: optional,
		maxDepth: 14,
		childCap: 2,
		ctr:      map[string]int{},
		oneofCtr: map[string]int{},
	}
}

func simpleName(fqn string) string {
	if i := strings.LastIndex(fqn, "."); i >= 0 {
		return fqn[i+1:]
	}
	return fqn
}

// render walks the message at fqn and returns its concrete markup string.
func (r *Renderer) render(fqn string, depth int, seen map[string]int) string {
	if depth > r.maxDepth {
		return ""
	}

	// 1. Open-ended leaf type → ONE representative sample, rotated by a per-FQN
	//    counter so successive instances of the same leaf vary (distinct
	//    attribute values rather than duplicates).
	if samples, ok := reps[simpleName(fqn)]; ok && len(samples) > 0 {
		v := samples[r.ctr[fqn]%len(samples)]
		r.ctr[fqn]++
		return v
	}

	// 2. Keyword/terminal message → its literal (a markup token like `<rect`,
	//    `="`, `>`, `</rect>`, or an enum keyword like `evenodd`).
	if lit, ok := r.kw[fqn]; ok {
		return lit
	}

	m := r.byFQN[fqn]
	if m == nil {
		return "" // unresolved reference
	}

	// 3. Cycle guard: allow a message to appear at most twice on the active path.
	if seen[fqn] >= 2 {
		return ""
	}
	seen[fqn]++
	defer func() { seen[fqn]-- }()

	// 4. Alternation (proto oneof): pick ONE arm, rotating the chosen arm by a
	//    per-oneof counter so repeated instances of an attribute/content union
	//    yield DIFFERENT members instead of duplicates. Skip arms that just
	//    recurse into a cycle (would render empty) so rotation lands on a
	//    productive member.
	if len(m.GetOneofDecl()) > 0 {
		fields := m.GetField()
		if len(fields) == 0 {
			return ""
		}
		n := len(fields)
		start := r.oneofCtr[fqn]
		for i := 0; i < n; i++ {
			idx := (start + i) % n
			f := fields[idx]
			// Guard against an arm whose only target is a message already on the
			// active path (would render to "" via the cycle guard); skip it.
			if tn := f.GetTypeName(); tn != "" && r.byFQN[tn] != nil && seen[tn] >= 2 {
				continue
			}
			out := r.renderField(fqn, f, depth, seen)
			if out != "" {
				r.oneofCtr[fqn] = (idx + 1) % n
				return out
			}
		}
		// Nothing productive this round — advance the counter and give up.
		r.oneofCtr[fqn] = (start + 1) % n
		return ""
	}

	// 5. Plain sequence message: concatenate every field in order. This is where
	//    an element's open tag, attribute list, child list, and close tag are
	//    assembled into markup. Between two ADJACENT NUMERIC-LEAF fields (the
	//    four numbers of a viewBox, the x/y of a coordinate pair, the members of
	//    a number list) a separating space is inserted so they do not collapse
	//    into one invalid token (`01-13.14`). Markup terminals and structured
	//    punctuation are concatenated verbatim (their own literals carry the
	//    spacing, e.g. the leading space in ` x="`).
	var b strings.Builder
	prevNumeric, prevWord := false, false
	for _, f := range m.GetField() {
		part := r.renderField(fqn, f, depth, seen)
		if part == "" {
			continue
		}
		// Numeric-leaf separation is STRUCTURAL (by field type): two adjacent
		// numeric leaves — the four numbers of a viewBox, the x/y of a coordinate
		// pair — get a space even when the second is negative (`1` `-1`), which a
		// character rule could not distinguish.
		curNumeric := f.GetLabel() != descriptorpb.FieldDescriptorProto_LABEL_REPEATED &&
			numericLeaf[simpleName(f.GetTypeName())]
		// Keyword-word separation is by RENDERED FRAGMENT: two adjacent plain
		// word fragments (e.g. preserveAspectRatio's `none` `meet`) get a space.
		// Markup-bearing fragments (containing `<`, `>`, `"`) and structured
		// value strings (SMIL clock values, digit-spelled numbers) never match.
		curWord := isPlainWord(part)
		if (prevNumeric && curNumeric) || (prevWord && curWord) {
			b.WriteByte(' ')
		}
		b.WriteString(part)
		prevNumeric, prevWord = curNumeric, curWord
	}
	return b.String()
}

// isPlainWord reports whether s is a single bare keyword word (a letter followed
// by letters/hyphens, e.g. `none`, `meet`, `userSpaceOnUse`). Used to separate
// adjacent enum-keyword arms with a space.
func isPlainWord(s string) bool {
	if len(s) == 0 || !isLetter(s[0]) {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !isLetter(c) && c != '-' {
			return false
		}
	}
	return true
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// numericLeaf is the set of leaf *Type messages that render as bare numeric
// tokens (viewBox numbers, coordinate pairs, number lists).
var numericLeaf = map[string]bool{
	"NumberType":           true,
	"IntegerType":          true,
	"CoordinateType":       true,
	"LengthType":           true,
	"PercentageType":       true,
	"LengthPercentageType": true,
	"AngleType":            true,
}

// renderField renders one field of the message at parentFQN.
//
//   - REPEATED (the EBNF `{ }` repetition — attribute list, child list, point
//     list, …) → render the field's message type childCap times, rotating oneof
//     arms each iteration, concatenated. THIS IS THE FIX: never drop a repeated
//     field, or the document would be empty.
//   - TYPE_STRING (a bare scalar `string value` — only on un-scalarized leaves,
//     which are already handled by reps above) → "".
//   - TYPE_MESSAGE → recurse into the field's target.
func (r *Renderer) renderField(parentFQN string, f *descriptorpb.FieldDescriptorProto, depth int, seen map[string]int) string {
	if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		target := f.GetTypeName()
		if target == "" {
			return ""
		}
		var b strings.Builder
		for i := 0; i < r.childCap; i++ {
			b.WriteString(r.render(target, depth+1, seen))
		}
		return b.String()
	}
	if f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_STRING {
		return ""
	}
	target := f.GetTypeName()
	if target == "" {
		return ""
	}
	return r.render(target, depth+1, seen)
}

// RenderDocument renders the whole SVG document from the root production.
func (r *Renderer) RenderDocument() string {
	return r.render(".svg.SvgDocument", 0, map[string]int{})
}

// RenderElement renders a single element (by its compiled FQN) and, if the
// result is not already a complete `<svg …>…</svg>` document, wraps it in a
// minimal valid SVG root so it is independently well-formed and renderable.
func (r *Renderer) RenderElement(elementFQN string) string {
	body := r.render(elementFQN, 0, map[string]int{})
	if strings.HasPrefix(strings.TrimSpace(body), "<svg") {
		return body
	}
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">` + body + `</svg>`
}

// RenderElementWithAttrs renders one element with attrCount attributes and NO
// child elements, then wraps it in a minimal valid SVG root. It is used for the
// sample-rect.svg snapshot, which should show SEVERAL distinct attributes (so
// repeated-field rendering + oneof rotation are clearly visible) while staying a
// minimal, empty-bodied element. The element compiles to a sequence:
// open-tag, repeated attribute-list, ">", repeated content-list, close-tag — so
// we walk the fields and emit the first repetition (attributes) attrCount times
// with rotation, skip the content repetition, and emit the rest verbatim.
func (r *Renderer) RenderElementWithAttrs(elementFQN string, attrCount int) string {
	m := r.byFQN[elementFQN]
	if m == nil {
		return r.RenderElement(elementFQN)
	}
	seen := map[string]int{}
	var b strings.Builder
	repeatedSeen := 0
	for _, f := range m.GetField() {
		if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			repeatedSeen++
			if repeatedSeen == 1 {
				// attribute list — emit attrCount distinct attributes
				for i := 0; i < attrCount; i++ {
					b.WriteString(r.render(f.GetTypeName(), 1, seen))
				}
			}
			// content list (repeatedSeen >= 2): emit nothing → empty element
			continue
		}
		b.WriteString(r.renderField(elementFQN, f, 0, seen))
	}
	inner := b.String()
	if strings.HasPrefix(strings.TrimSpace(inner), "<svg") {
		return inner
	}
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">` + inner + `</svg>`
}
