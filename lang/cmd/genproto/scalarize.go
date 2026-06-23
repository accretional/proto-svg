package main

import (
	"strings"

	"github.com/accretional/gluon/v2/compiler"
	pb "github.com/accretional/gluon/v2/pb"
)

// leafTypes are the atomic, open-ended SVG value leaf types — the scalars and
// tokens that bottom out in infinite character ranges (digits, units, ident
// characters, hex digits, string chars). The grammar fully expands them, but a
// caller constructing a value structurally has no business building a digit
// tree to spell "16px"; these collapse to a proto3 `string value = 1` field
// that the caller fills with the literal token. Everything composite above a
// leaf (PaintType, ViewBox, PreserveAspectRatio, transform lists, path data,
// point lists) stays a structured message graph walked by the renderer.
//
// Both the snake_case atom and the PascalCase wrapper normalize to the same
// message name and are scalarized together.
//
// This mirrors proto-css's scalarizeLeaves: a small, local, grammar-specific
// transform composed into the otherwise grammar-agnostic gluon pipeline.
//
// Note: PaintType / ViewBox / PreserveAspectRatio / UnitType / AlphaValue are
// deliberately ABSENT — they are structured value grammars (keyword arms plus
// ColorType/UrlType/NumberType sub-parts) and must be rendered, not scalarized.
var leafTypes = normSet(
	"number_type",
	"non_negative_number_type",
	"integer_type",
	"non_negative_integer_type",
	"positive_integer_type",
	"length_type",
	"non_negative_length_type",
	"percentage_type",
	"length_percentage_type",
	"non_negative_length_percentage_type",
	"miter_limit_type",
	"alpha_value",
	"coordinate_type",
	"angle_type",
	"time_type",
	// color_type is NOT a leaf — it is a structured oneof (HexColor |
	// FunctionalColor | NamedColor). Its two open arms ARE leaves; NamedColor is
	// the closed 148-keyword enum, walked as a oneof.
	"hex_color",
	"functional_color",
	"iri_type",
	"url_type",
	"string_type",
	"char_type",
	"custom_ident_type",
	"xml_name_type",
	"bcp47_type",
	// list value types are NO LONGER leaves — they are structured `repeated`
	// fields the renderer walks (list_of_numbers_type, number_optional_number_type,
	// dasharray_type). list_of_lengths_type was dead (zero grammar consumers) and
	// has been deleted from datatype.ebnf.
	"event_symbol_type",
	"character_data_type",
)

// normSet builds a set of names normalized to lowercase-without-underscores, so
// "time_type" and "TimeType" map to the same key.
func normSet(names ...string) map[string]bool {
	s := make(map[string]bool, len(names))
	for _, n := range names {
		s[norm(n)] = true
	}
	return s
}

func norm(s string) string { return strings.ToLower(strings.ReplaceAll(s, "_", "")) }

// scalarizeLeaves walks the AST and replaces the body of every leaf rule with a
// single scalar node, so it lowers to `message X { string value = 1; }`. A rule
// is a leaf when either:
//
//   - its normalized name is in leafTypes (the curated open-ended value types), or
//   - its body contains a character range (e.g. digit = "0"…"9"). A range is a
//     lexical primitive; collapsing it keeps svg.proto self-contained (ranges
//     otherwise lower to .unicode.UTF8 fields and pull in unicode/utf_8.proto)
//     and is the right shape for construction anyway — callers spell the token
//     directly.
//
// Rules are matched by normalized name, so both the snake_case atom and its
// PascalCase wrapper collapse to the same scalar message. The input is not
// mutated; a deep copy is returned.
func scalarizeLeaves(root *pb.ASTNode) *pb.ASTNode {
	if root == nil {
		return nil
	}
	if root.GetKind() == compiler.KindRule && (leafTypes[norm(root.GetValue())] || hasRange(root)) {
		return &pb.ASTNode{
			Kind:     compiler.KindRule,
			Value:    root.GetValue(),
			Children: []*pb.ASTNode{{Kind: compiler.KindScalar, Value: "value"}},
		}
	}
	kids := make([]*pb.ASTNode, 0, len(root.GetChildren()))
	for _, c := range root.GetChildren() {
		kids = append(kids, scalarizeLeaves(c))
	}
	return &pb.ASTNode{Kind: root.GetKind(), Value: root.GetValue(), Children: kids}
}

// hasRange reports whether the node's own subtree contains a character range.
// Nonterminal references are leaf nodes pointing at other rules, so this only
// sees ranges written inline in this rule's body.
func hasRange(node *pb.ASTNode) bool {
	if node == nil {
		return false
	}
	if node.GetKind() == compiler.KindRange {
		return true
	}
	for _, c := range node.GetChildren() {
		if hasRange(c) {
			return true
		}
	}
	return false
}
