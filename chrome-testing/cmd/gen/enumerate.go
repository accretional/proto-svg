package main

import (
	"sort"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// enumerate.go — the ALL-VALUE-PATHS walk.
//
// USER REQUIREMENT (hard): "walk through the grammar for ALL POSSIBLE PATHS and
// show them as SVGs". For EACH element message E we enumerate, per attribute,
// EVERY value-path the grammar admits, and render E with that single attribute
// set to that single value plus a minimal visible baseline. Nothing about the
// element is predetermined — the open tag, the attribute name, and the value are
// all read out of the compiled grammar; only the NEIGHBORS come from the
// blueprint (blueprint.go).
//
// Per attribute arm A (a sequence ` name="`, VALUE, `"`):
//   - closed enum (oneof of keyword terminals)  → ONE variant per keyword.
//   - leaf in reps (free terminal)              → one variant per rep sample.
//   - structured value (PaintType, ViewBox, …)  → a few canonical instances.
// The overlay (overlay.go) is consulted so each emitted value is valid by
// construction (refs → slot, alpha ∈ [0,1], non-negative magnitudes, …).

// Variant is one enumerated (attribute, value, element-markup) tuple.
type Variant struct {
	Attr       string // attribute name, e.g. "x", "fill", "stdDeviation"
	Value      string // the value text inside the quotes, e.g. "10", "url(#slot)"
	Markup     string // the element rendered with exactly this one attribute
	WrappedSVG string // the blueprint-injected full <svg> (set in the emit pass)
	NeedsID    bool   // the variant references #slot, so the element must carry id="slot"
}

// Enumerator walks element messages and produces value-path variants. It reuses
// the engine's Renderer for the structured/baseline rendering and the kw/reps
// maps for token resolution.
type Enumerator struct {
	byFQN map[string]*descriptorpb.DescriptorProto
	kw    map[string]string
	r     *Renderer

	// maxLeaf caps rep samples per leaf attribute (every sample is a value-path,
	// but the gallery would be unreadable with infinite numeric variety).
	maxLeaf int
}

func newEnumerator(byFQN map[string]*descriptorpb.DescriptorProto, kw map[string]string, r *Renderer) *Enumerator {
	return &Enumerator{byFQN: byFQN, kw: kw, r: r, maxLeaf: 6}
}

// element describes one element message: its open tag, the FQN, the attribute
// union field's target message, and the content (children) field's target.
type element struct {
	tag        string
	fqn        string
	attrUnion  string // FQN of the repeated attribute-union message
	contentMsg string // FQN of the repeated content message (may be "")
}

// allElements returns every element message (leading terminal is "<tag", not a
// close tag), sorted by tag.
func (e *Enumerator) allElements() []element {
	var els []element
	for fqn, m := range e.byFQN {
		fields := m.GetField()
		if len(fields) == 0 {
			continue
		}
		lit := e.kw[fields[0].GetTypeName()]
		if !strings.HasPrefix(lit, "<") || strings.HasPrefix(lit, "</") {
			continue
		}
		attrUnion, content := e.attrAndContent(m)
		els = append(els, element{
			tag:        strings.TrimPrefix(lit, "<"),
			fqn:        fqn,
			attrUnion:  attrUnion,
			contentMsg: content,
		})
	}
	sort.Slice(els, func(i, j int) bool { return els[i].tag < els[j].tag })
	return els
}

// attrAndContent finds an element's attribute-union field (the FIRST repeated
// field, which holds the { Attribute } list) and its content field (the SECOND
// repeated field, the { Content } children list).
func (e *Enumerator) attrAndContent(m *descriptorpb.DescriptorProto) (attrUnion, content string) {
	reps := 0
	for _, f := range m.GetField() {
		if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			reps++
			switch reps {
			case 1:
				attrUnion = f.GetTypeName()
			case 2:
				content = f.GetTypeName()
			}
		}
	}
	return
}

// enumerateElement produces every value-path variant for one element. It walks
// the attribute union's oneof arms; for each concrete attribute production it
// enumerates that attribute's value-paths. Shared attribute groups
// (CoreAttribute, PresentationAttribute, …) are themselves oneofs and are
// flattened so their members are enumerated too — but to keep the gallery
// focused and finite we only descend ONE level into shared groups and cap the
// presentation/aria/event groups (they are huge and mostly value-identical
// strings), keeping every ELEMENT-SPECIFIC attribute exhaustive.
func (e *Enumerator) enumerateElement(el element) []Variant {
	um := e.byFQN[el.attrUnion]
	if um == nil {
		return nil
	}
	var out []Variant
	for _, arm := range um.GetField() {
		armMsg := e.byFQN[arm.GetTypeName()]
		if armMsg == nil {
			continue
		}
		if len(armMsg.GetOneofDecl()) > 0 {
			// Shared attribute group (CoreAttribute / PresentationAttribute / …):
			// descend one level into its members.
			out = append(out, e.enumerateGroup(el, arm.GetTypeName())...)
			continue
		}
		// A concrete attribute production sequence.
		out = append(out, e.enumerateAttr(el, arm.GetTypeName())...)
	}
	return out
}

// sharedGroupCap caps how many members of a large shared group (presentation,
// aria, events) we enumerate, and how many value-paths per member. Element-
// specific attributes are never capped.
var sharedGroupMemberCap = map[string]int{
	"AriaAttribute":                  3,
	"GlobalEventAttribute":           2,
	"DocumentElementEventAttribute":  1,
	"GraphicalEventAttribute":        1,
	"DocumentEventAttribute":         1,
	"AnimationEventAttribute":        1,
	"ConditionalProcessingAttribute": 2,
	// PresentationAttribute holds ~52 members; the painting attrs fill positions
	// 1-24, so a low cap cut off opacity/transform/transform-origin/display/
	// visibility/clip-path/mask for container elements (g/defs/switch/a/symbol) —
	// QA round2. Raise the cap to 64 so every presentation attr is enumerated and
	// the visually-important ones reach the main grid. FIX 1's tag-aware grouping
	// collapses the no-effect ones, so enumerating them all is safe.
	"PresentationAttribute": 64,
	"CoreAttribute":         6,
}

func (e *Enumerator) enumerateGroup(el element, groupFQN string) []Variant {
	gm := e.byFQN[groupFQN]
	if gm == nil {
		return nil
	}
	cap := sharedGroupMemberCap[simpleName(groupFQN)]
	if cap == 0 {
		cap = 6
	}
	var out []Variant
	members := 0
	for _, arm := range gm.GetField() {
		if members >= cap {
			break
		}
		armMsg := e.byFQN[arm.GetTypeName()]
		if armMsg == nil {
			continue
		}
		if len(armMsg.GetOneofDecl()) > 0 {
			// nested group (rare) — descend once more, but lightly
			out = append(out, e.enumerateGroup(el, arm.GetTypeName())...)
			members++
			continue
		}
		vs := e.enumerateAttr(el, arm.GetTypeName())
		// For big shared groups, keep just the first value-path per member so the
		// gallery stays element-focused.
		if simpleName(groupFQN) != "CoreAttribute" && len(vs) > 1 {
			vs = vs[:1]
		}
		out = append(out, vs...)
		members++
	}
	return out
}

// enumerateAttr enumerates every value-path of ONE attribute production. The
// production is a sequence: leading ` name="` keyword, a value field (or a few),
// trailing `"` keyword. We classify the value field and enumerate accordingly.
func (e *Enumerator) enumerateAttr(el element, attrFQN string) []Variant {
	attrMsg := e.byFQN[attrFQN]
	if attrMsg == nil {
		return nil
	}
	name, prefix, suffix, valFields := e.splitAttr(attrMsg)
	if name == "" || len(valFields) == 0 {
		return nil
	}
	var values []valuePath
	if len(valFields) == 1 {
		values = e.enumerateValue(el.tag, name, valFields[0].GetTypeName())
	} else {
		// A multi-field value (e.g. a structured sequence inlined into the attr).
		// Render one canonical instance via the engine renderer.
		v := e.r.render(attrFQN, 0, map[string]int{})
		v = stripQuotes(v, prefix, suffix)
		if v != "" {
			values = []valuePath{{value: v}}
		}
	}
	var out []Variant
	for _, vp := range values {
		markup, needsID := e.renderWithOneAttr(el, prefix, vp.value, suffix)
		out = append(out, Variant{
			Attr:    name,
			Value:   vp.value,
			Markup:  markup,
			NeedsID: needsID || vp.needsID,
		})
	}
	return out
}

// valuePath is one enumerated value for an attribute.
type valuePath struct {
	value   string
	needsID bool
}

// splitAttr decomposes an attribute production sequence into its name, the
// opening ` name="` prefix literal, the trailing `"` suffix literal, and the
// value field(s) in between.
func (e *Enumerator) splitAttr(m *descriptorpb.DescriptorProto) (name, prefix, suffix string, valFields []*descriptorpb.FieldDescriptorProto) {
	fields := m.GetField()
	if len(fields) == 0 {
		return
	}
	// prefix = leading keyword literal (e.g. ` x="`).
	prefix = e.kw[fields[0].GetTypeName()]
	name = attrNameFromPrefix(prefix)
	// suffix = trailing keyword literal (the closing quote).
	if last := fields[len(fields)-1]; e.kw[last.GetTypeName()] != "" {
		suffix = e.kw[last.GetTypeName()]
	}
	// value fields = everything between the first and last keyword fields.
	for i := 1; i < len(fields)-1; i++ {
		valFields = append(valFields, fields[i])
	}
	// If there was no trailing keyword (no suffix), the value runs to the end.
	if suffix == "" && len(fields) > 1 {
		valFields = fields[1:]
	}
	return
}

// attrNameFromPrefix extracts the attribute name from a ` name="` keyword
// literal, e.g. ` stdDeviation="` → "stdDeviation".
func attrNameFromPrefix(prefix string) string {
	s := strings.TrimSpace(prefix)
	s = strings.TrimSuffix(s, "=\"")
	s = strings.TrimSuffix(s, "=")
	return strings.TrimSpace(s)
}

// enumerateValue produces every value-path for a value TYPE in the context of
// attribute attrName. This is the classifier the user requirements describe:
//   - keyword-only oneof → one variant per keyword
//   - reps leaf          → one variant per rep sample (overlay may override)
//   - structured value   → a few canonical instances
func (e *Enumerator) enumerateValue(tag, attrName, valFQN string) []valuePath {
	// Overlay may fully determine the value (refs, alpha, monotone lists).
	if ov, ok := overlaySample(tag, attrName, simpleName(valFQN)); ok {
		return []valuePath{{value: ov, needsID: strings.Contains(ov, "#"+slotID)}}
	}

	// Leaf type with rep samples.
	if samples, ok := reps[simpleName(valFQN)]; ok {
		var out []valuePath
		n := len(samples)
		if n > e.maxLeaf {
			n = e.maxLeaf
		}
		for _, s := range samples[:n] {
			out = append(out, valuePath{value: s})
		}
		return out
	}

	m := e.byFQN[valFQN]
	if m == nil {
		// Maybe the value field is itself a keyword (boolean-presence etc.).
		if lit, ok := e.kw[valFQN]; ok && lit != "" {
			return []valuePath{{value: lit}}
		}
		return nil
	}

	// Closed enum: a oneof whose arms are ALL keyword terminals → one per keyword.
	if len(m.GetOneofDecl()) > 0 {
		if kws := e.keywordArms(m); kws != nil {
			var out []valuePath
			for _, kw := range kws {
				out = append(out, valuePath{value: kw})
			}
			return out
		}
		// Mixed oneof (keywords + leaves + structured): enumerate each arm's
		// value-paths and concatenate (every arm is a distinct value-path).
		return e.enumerateMixedOneof(tag, attrName, m)
	}

	// Structured value (sequence): canonical instances.
	return e.structuredInstances(attrName, valFQN, m)
}

// keywordArms returns the keyword literals if EVERY arm of the oneof is a bare
// keyword terminal; otherwise nil (mixed oneof).
func (e *Enumerator) keywordArms(m *descriptorpb.DescriptorProto) []string {
	var out []string
	for _, f := range m.GetField() {
		lit, ok := e.kw[f.GetTypeName()]
		if !ok || lit == "" {
			return nil
		}
		out = append(out, lit)
	}
	return out
}

// enumerateMixedOneof handles a oneof whose arms mix keywords, leaves, and
// structured values (e.g. PaintType, RxAttr's auto|length). Each arm yields its
// own value-paths.
func (e *Enumerator) enumerateMixedOneof(tag, attrName string, m *descriptorpb.DescriptorProto) []valuePath {
	var out []valuePath
	for _, f := range m.GetField() {
		tn := f.GetTypeName()
		if lit, ok := e.kw[tn]; ok && lit != "" {
			out = append(out, valuePath{value: lit})
			continue
		}
		// reps leaf or nested structured/oneof → recurse, but cap leaves to keep
		// the union readable.
		sub := e.enumerateValue(tag, attrName, tn)
		if _, isLeaf := reps[simpleName(tn)]; isLeaf && len(sub) > 2 {
			sub = sub[:2]
		}
		out = append(out, sub...)
	}
	return dedupePaths(out)
}

// structuredInstances returns a few canonical concrete instances of a structured
// value type. For the well-known shared sub-grammars we hand-pick instances that
// exercise distinct shapes; otherwise we fall back to rendering the sequence via
// the engine renderer a couple of times (rotation yields variety).
func (e *Enumerator) structuredInstances(attrName, valFQN string, m *descriptorpb.DescriptorProto) []valuePath {
	switch simpleName(valFQN) {
	case "PaintType", "PaintRef":
		return paths("none", "#e94560", "url(#"+slotID+")", "currentColor")
	case "ViewBox":
		return paths("0 0 100 100", "0 0 50 50", "-10 -10 120 120")
	case "PreserveAspectRatio":
		return paths("none", "xMidYMid meet", "xMinYMin slice")
	case "TransformList", "Transform":
		return paths("translate(20 10)", "rotate(45)", "scale(1.5)", "skewX(20)")
	case "SvgPath", "PathData":
		return paths("M10 10 L90 90", "M10 50 Q50 10 90 50", "M20 20 H80 V80 H20 Z")
	case "Points", "PointsList":
		return paths("20,20 80,20 50,80", "10,80 50,10 90,80 50,50")
	case "Align":
		return paths("none", "xMidYMid", "xMinYMin")
	}
	// Generic structured value: render the sequence a couple of times.
	var out []valuePath
	for i := 0; i < 2; i++ {
		v := e.r.render(valFQN, 0, map[string]int{})
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, valuePath{value: v})
		}
	}
	return dedupePaths(out)
}

// renderWithOneAttr produces the element markup carrying EXACTLY one attribute
// (prefix+value+suffix) plus a minimal visible baseline so the varied attribute
// is observable. The baseline is injected via baselineAttrs (size/fill for a
// shape, etc.). Returns the markup and whether the element references #slot.
func (e *Enumerator) renderWithOneAttr(el element, prefix, value, suffix string) (string, bool) {
	open := "<" + el.tag
	baseline, needsID := baselineFor(el.tag, prefix)
	needsRef := strings.Contains(value, "#"+slotID)
	if needsRef {
		// The element references itself by id — give it id="slot".
		baseline = ` id="` + slotID + `"` + baseline
	}
	attr := prefix + value + suffix
	closeTag := "</" + el.tag + ">"
	// Self-closing-friendly: emit open + baseline + this attr + ">" + close.
	return open + baseline + attr + ">" + bodyFor(el.tag) + closeTag, needsID || needsRef
}

func paths(vs ...string) []valuePath {
	out := make([]valuePath, 0, len(vs))
	for _, v := range vs {
		out = append(out, valuePath{value: v, needsID: strings.Contains(v, "#"+slotID)})
	}
	return out
}

func dedupePaths(in []valuePath) []valuePath {
	seen := map[string]bool{}
	var out []valuePath
	for _, p := range in {
		if p.value == "" || seen[p.value] {
			continue
		}
		seen[p.value] = true
		out = append(out, p)
	}
	return out
}

// stripQuotes removes a leading prefix and trailing suffix keyword from a
// rendered attribute string, leaving the bare value text.
func stripQuotes(s, prefix, suffix string) string {
	s = strings.TrimPrefix(s, prefix)
	s = strings.TrimSuffix(s, suffix)
	return strings.TrimSpace(s)
}
