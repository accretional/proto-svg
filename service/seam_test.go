package service

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	csspb "github.com/accretional/proto-css/proto/pb/css"
	svgpb "github.com/accretional/proto-svg/proto/pb/svg"
)

const svgOpen = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" width="120" height="120">`

// TestStyleSeam is the crux: parsing an SVG whose <style> carries CSS must yield
// an SvgDocument whose css_style_sheet is a structured css.CssStyleSheet subtree
// (not an opaque string), and rendering it back must reproduce the markup + CSS.
func TestStyleSeam(t *testing.T) {
	in := svgOpen + `<style>circle{fill:red}</style><circle cx="60" cy="60" r="40"></circle></svg>`

	doc, err := Parse(in)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	sheet := findCSSAST(doc.ProtoReflect())
	if sheet == nil {
		t.Fatal("no css.CssStyleSheet embedded — the SVG <style> seam produced no structured CSS AST")
	}
	if len(sheet.GetAlt1()) == 0 {
		t.Fatal("embedded CSS AST has no rules — opaque, not structured")
	}

	out, err := Render(doc)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	t.Logf("render round-trip: %s", out)
	for _, want := range []string{"<svg ", "<style>", "circle", "fill:red", "</style>", "<circle ", "</svg>"} {
		if !strings.Contains(out, want) {
			t.Errorf("rendered SVG missing %q:\n%s", want, out)
		}
	}

	// Idempotency: re-parse + re-render must be stable.
	doc2, err := Parse(out)
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	out2, err := Render(doc2)
	if err != nil {
		t.Fatalf("re-render: %v", err)
	}
	if out != out2 {
		t.Errorf("render not idempotent:\n out:  %q\n out2: %q", out, out2)
	}
}

// TestPlainSVG exercises the non-seam path (a plain shape) to confirm the SVG
// engine itself round-trips.
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

func findCSSAST(m protoreflect.Message) *csspb.CssStyleSheet {
	var found *csspb.CssStyleSheet
	var walk func(protoreflect.Message)
	walk = func(m protoreflect.Message) {
		if found != nil {
			return
		}
		if m.Descriptor().FullName() == "css.CssStyleSheet" {
			if s, ok := m.Interface().(*csspb.CssStyleSheet); ok {
				found = s
				return
			}
		}
		m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			if fd.Kind() != protoreflect.MessageKind && fd.Kind() != protoreflect.GroupKind {
				return found == nil
			}
			if fd.IsList() {
				l := v.List()
				for i := 0; i < l.Len(); i++ {
					walk(l.Get(i).Message())
				}
			} else {
				walk(v.Message())
			}
			return found == nil
		})
	}
	walk(m)
	return found
}

var _ = svgpb.SvgDocument{}
