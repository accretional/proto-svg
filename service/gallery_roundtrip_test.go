package service

import (
	"strings"
	"testing"

	// Link CSS so the fill/stroke color seams structure into css.ColorType.
	_ "github.com/accretional/proto-css/service"
)

// TestGalleryRoundTrip proves the gallery pipeline: a walked grammar path's
// MARKUP (what the enumerator emits today) is structured by codec.Parse into the
// shipped AST — including the cross-grammar fill→Any(css.ColorType) seam — and
// re-emitted by codec.Render byte-for-byte. Parse and Render are the two
// generators under test; every gallery path exercises both.
func TestGalleryRoundTrip(t *testing.T) {
	// A realistic gallery specimen: a shape with element-specific geometry plus a
	// presentation-attribute paint value (the fill color is a css seam).
	const markup = `<rect x="10" y="10" width="80" height="80" fill="#e94560" stroke="#16c79a" stroke-width="2"></rect>`

	ast, err := ParseAs(markup, "svg.SvgrectElement")
	if err != nil {
		t.Fatalf("codec.Parse (the parse generator) failed on a walked path: %v", err)
	}

	// The parser must have STRUCTURED the paint value into the embedded css grammar,
	// not left it a string — this is the cross-grammar seam the gallery must exercise.
	if findEmbedded(ast.ProtoReflect(), "css.ColorType") == nil {
		t.Fatal("fill color did not structure into Any(css.ColorType) — seam not walked")
	}

	out, err := Render(ast)
	if err != nil {
		t.Fatalf("codec.Render (the render generator) failed: %v", err)
	}
	if out != markup {
		t.Fatalf("round-trip not byte-exact:\n walked: %q\n codec : %q", markup, out)
	}
	if !strings.Contains(out, `fill="#e94560"`) {
		t.Fatalf("rendered markup lost the seam value: %q", out)
	}
}
