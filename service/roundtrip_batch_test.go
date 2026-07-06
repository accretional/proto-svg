package service

import (
	"strings"
	"testing"

	svgpb "github.com/accretional/proto-svg/proto/pb/svg"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	_ "github.com/accretional/proto-css/service"
	_ "github.com/accretional/proto-html/service"
)

// rootForMarkup resolves the shipped element type for a markup fragment by exact
// open-tag prefix match (skipping ghost keyword entries), mirroring the gallery
// gen's tagRoot.
func tagOf(openLit string) string { // "<svg " / "<rect" -> "svg" / "rect"
	t := strings.TrimPrefix(openLit, "<")
	if i := strings.IndexByte(t, ' '); i >= 0 {
		t = t[:i]
	}
	return t
}

func rootForMarkup(markup string) string {
	if !strings.HasPrefix(markup, "<") {
		return ""
	}
	tag := tagOf(markup[strings.IndexByte(markup, '<'):])
	for fqn, pfx := range svgpb.MessagePrefix {
		if len(pfx) != 1 || !strings.HasPrefix(pfx[0], "<") || strings.HasPrefix(pfx[0], "</") {
			continue
		}
		if tagOf(pfx[0]) != tag {
			continue
		}
		name := strings.TrimPrefix(fqn, ".")
		if strings.Contains(name, "Keyword") {
			continue
		}
		if _, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(name)); err == nil {
			return name
		}
	}
	return ""
}

// TestRoundTripBatch round-trips one representative of each failure category
// through the codec, reporting parse errors and byte mismatches separately.
func TestRoundTripBatch(t *testing.T) {
	cases := []string{
		// space-separated scalar lists (the parseRepeated scalar bug)
		`<rect x="10" y="10" width="80" height="80" stroke-dasharray="4 2"></rect>`,
		`<feConvolveMatrix order="3" kernelMatrix="0 -1 0 -1 5 -1 0 -1 0" in="SourceGraphic"></feConvolveMatrix>`,
		`<polyline points="20,20 80,20 50,80"></polyline>`,
		`<text x="10 20 30" y="40">hi</text>`,
		// semicolon lists (parse should now succeed; render may canonicalize spacing)
		`<animate attributeName="x" dur="2s" keyTimes="0; 0.5; 1"></animate>`,
		`<animateMotion dur="2s" keyPoints="0; 0.5; 1" path="M0 0 L60 60"></animateMotion>`,
		// preserveAspectRatio (grammar space fix)
		`<view preserveAspectRatio="xMidYMid meet" viewBox="0 0 100 100"></view>`,
		`<image preserveAspectRatio="xMidYMid slice" href="#x"></image>`,
		// svg root: xmlns/version baked into the open tag (grammar SvgXmlnsAttr)
		`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" width="120" height="120" viewBox="0 0 100 100" font-family="serif"><rect x="10" y="10" width="80" height="80" fill="#e94560"></rect></svg>`,
		// foreignObject → html.FlowContent seam (the svg→html cycle)
		`<foreignObject x="10" y="5" width="90" height="90"><div xmlns="http://www.w3.org/1999/xhtml" style="width:100%;height:100%;background:#4d8bff">HTML in SVG</div></foreignObject>`,
	}
	for _, in := range cases {
		root := rootForMarkup(in)
		if root == "" {
			t.Errorf("no root type for %q", in)
			continue
		}
		msg, err := ParseAs(in, root)
		if err != nil {
			t.Errorf("PARSE-FAIL %s\n  in : %s\n  err: %v", root, in, err)
			continue
		}
		out, err := Render(msg)
		if err != nil {
			t.Errorf("RENDER-FAIL %s: %v", root, err)
			continue
		}
		if out != in {
			t.Errorf("MISMATCH %s\n  in : %s\n  out: %s", root, in, out)
			continue
		}
		t.Logf("OK  %s", in)
	}
}
