package service

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// TestASTBuildPOC proves option A's mechanical core: build a shipped-type SVG
// element as a proto AST (by reflection over the registered svg types) and let
// the gluon codec render it — no string concatenation. If this yields exact
// markup, the gallery gen can build ASTs and the codec becomes the renderer of
// record.
func TestASTBuildPOC(t *testing.T) {
	newMsg := func(name string) protoreflect.Message {
		mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(name))
		if err != nil {
			t.Fatalf("find %s: %v", name, err)
		}
		return mt.New()
	}
	setStr := func(m protoreflect.Message, field, val string) {
		fd := m.Descriptor().Fields().ByName(protoreflect.Name(field))
		if fd == nil {
			t.Fatalf("%s has no field %s", m.Descriptor().FullName(), field)
		}
		m.Set(fd, protoreflect.ValueOfString(val))
	}
	setMsg := func(m protoreflect.Message, field string, sub protoreflect.Message) {
		fd := m.Descriptor().Fields().ByName(protoreflect.Name(field))
		if fd == nil {
			t.Fatalf("%s has no field %s", m.Descriptor().FullName(), field)
		}
		m.Set(fd, protoreflect.ValueOfMessage(sub))
	}
	appendMsg := func(m protoreflect.Message, field string, sub protoreflect.Message) {
		fd := m.Descriptor().Fields().ByName(protoreflect.Name(field))
		m.Mutable(fd).List().Append(protoreflect.ValueOfMessage(sub))
	}

	// one <length> attribute: RectAttribute{ rect_Xattr: RectXattr{ value, " } }
	lengthAttr := func(union, arm, val string) protoreflect.Message {
		leaf := newMsg("svg.LengthPercentageType")
		setStr(leaf, "value", val)
		attr := newMsg(arm)
		setMsg(attr, "length_percentage_type", leaf)
		setMsg(attr, "quotation_mark_keyword", newMsg("svg.QuotationMarkKeyword"))
		u := newMsg("svg.RectAttribute")
		setMsg(u, union, attr)
		return u
	}

	rect := newMsg("svg.SvgrectElement")
	appendMsg(rect, "rect_attribute", lengthAttr("rect_xattr", "svg.RectXattr", "10"))
	appendMsg(rect, "rect_attribute", lengthAttr("rect_yattr", "svg.RectYattr", "20"))
	setMsg(rect, "greater_than_sign_keyword", newMsg("svg.GreaterThanSignKeyword"))
	setMsg(rect, "less_than_sign_solidus_rect_greater_than_sign_keyword",
		newMsg("svg.LessThanSignSolidusRectGreaterThanSignKeyword"))

	out, err := Render(rect.Interface())
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	const want = `<rect x="10" y="20"></rect>`
	if out != want {
		t.Fatalf("AST→codec render mismatch:\n got: %q\nwant: %q", out, want)
	}
}
