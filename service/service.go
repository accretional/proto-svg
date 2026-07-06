// Package service renders typed SVG messages to SVG text and parses SVG text
// back, carrying no SVG logic: it registers the grammar-derived tables
// (svgpb.MessagePrefix / FieldSeparator / SeamType) with gluon's generic codec
// and delegates to it. Importing this package links the SVG grammar. SVG's
// cross-grammar seams (CSS in <style> and presentation-attribute values) are
// google.protobuf.Any; the codec descends into them when the CSS grammar is
// also linked (import proto-css's service), and leaves them opaque otherwise.
package service

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/accretional/gluon/v2/codec"
	svgpb "github.com/accretional/proto-svg/proto/pb/svg"
)

func init() {
	codec.Register(&codec.Grammar{
		Package:      "svg",
		Prefix:       svgpb.MessagePrefix,
		Separator:    svgpb.FieldSeparator,
		Seam:         svgpb.SeamType,
		SmartSpacing: false, // SVG: markup; spacing baked into terminals.
	})
}

// Render serializes any typed SVG message back into SVG text.
func Render(msg proto.Message) (string, error) {
	return codec.Render(codec.Default, msg)
}

// Parse parses SVG text into a typed SvgDocument (the start symbol).
func Parse(input string) (*svgpb.SvgDocument, error) {
	msg, err := ParseAs(input, "svg.SvgDocument")
	if err != nil {
		return nil, err
	}
	doc, ok := msg.(*svgpb.SvgDocument)
	if !ok {
		return nil, fmt.Errorf("internal: parsed %T", msg)
	}
	return doc, nil
}

// ParseAs parses SVG text against an arbitrary svg.* message type.
func ParseAs(input, typeName string) (proto.Message, error) {
	return codec.Parse(codec.Default, input, typeName)
}
