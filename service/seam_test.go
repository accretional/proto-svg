package service

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"

	// Link the CSS grammar so the codec can descend into SVG's css.* seams.
	_ "github.com/accretional/proto-css/service"
)

const svgOpen = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" width="120" height="120">`

// TestStyleSeam: an SVG <style> parses into a google.protobuf.Any embedding a
// structured css.CssStyleSheet, and round-trips.
func TestStyleSeam(t *testing.T) {
	in := svgOpen + `<style>circle{fill:red}</style><circle cx="60" cy="60" r="40"></circle></svg>`

	doc, err := Parse(in)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	sheet := findEmbedded(doc.ProtoReflect(), "css.CssStyleSheet")
	if sheet == nil {
		t.Fatal("no css.CssStyleSheet embedded — the SVG <style> seam produced no structured CSS AST")
	}
	if fieldLen(sheet, "alt1") == 0 {
		t.Fatal("embedded CSS AST has no rules — opaque, not structured")
	}

	out, err := Render(doc)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	t.Logf("round-trip: %s", out)
	for _, want := range []string{"<svg ", "<style>", "circle", "fill:red", "</style>", "<circle ", "</svg>"} {
		if !strings.Contains(out, want) {
			t.Errorf("rendered SVG missing %q:\n%s", want, out)
		}
	}
	// idempotency
	doc2, err := Parse(out)
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	out2, _ := Render(doc2)
	if out != out2 {
		t.Errorf("render not idempotent:\n %q\n %q", out, out2)
	}
}

// TestPresentationColorSeam: fill="red" parses into an Any embedding a
// structured css.ColorType.
func TestPresentationColorSeam(t *testing.T) {
	in := svgOpen + `<circle cx="60" cy="60" r="40" fill="red"></circle></svg>`
	doc, err := Parse(in)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if findEmbedded(doc.ProtoReflect(), "css.ColorType") == nil {
		t.Fatal("no css.ColorType embedded — the fill presentation-attribute seam did not fire")
	}
	out, err := Render(doc)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, `fill="red"`) {
		t.Errorf("rendered SVG missing fill=\"red\":\n%s", out)
	}
	if out != in {
		t.Errorf("round-trip mismatch:\n in: %q\nout: %q", in, out)
	}
}

// TestPlainSVG: the non-seam path round-trips exactly.
func TestPlainSVG(t *testing.T) {
	in := svgOpen + `<rect x="10" y="10" width="100" height="100"></rect></svg>`
	doc, err := Parse(in)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	out, err := Render(doc)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if out != in {
		t.Errorf("round-trip mismatch:\n in: %q\nout: %q", in, out)
	}
}

// findEmbedded walks a message tree (unpacking google.protobuf.Any seams) and
// returns the first message with the given fully-qualified name.
func findEmbedded(m protoreflect.Message, name string) protoreflect.Message {
	if string(m.Descriptor().FullName()) == name {
		return m
	}
	if m.Descriptor().FullName() == "google.protobuf.Any" {
		if sub := unpackAny(m); sub != nil {
			return findEmbedded(sub, name)
		}
		return nil
	}
	var found protoreflect.Message
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.Kind() != protoreflect.MessageKind && fd.Kind() != protoreflect.GroupKind {
			return true
		}
		if fd.IsList() {
			l := v.List()
			for i := 0; i < l.Len() && found == nil; i++ {
				found = findEmbedded(l.Get(i).Message(), name)
			}
		} else {
			found = findEmbedded(v.Message(), name)
		}
		return found == nil
	})
	return found
}

func unpackAny(m protoreflect.Message) protoreflect.Message {
	a, ok := m.Interface().(*anypb.Any)
	if !ok {
		b, err := proto.Marshal(m.Interface())
		if err != nil {
			return nil
		}
		a = &anypb.Any{}
		if proto.Unmarshal(b, a) != nil {
			return nil
		}
	}
	sub, err := a.UnmarshalNew()
	if err != nil {
		return nil
	}
	return sub.ProtoReflect()
}

func fieldLen(m protoreflect.Message, name string) int {
	fd := m.Descriptor().Fields().ByName(protoreflect.Name(name))
	if fd == nil || !fd.IsList() {
		return 0
	}
	return m.Get(fd).List().Len()
}
