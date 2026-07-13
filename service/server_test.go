package service

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	svgservicepb "github.com/accretional/proto-svg/proto/pb/svgservice"
)

const serverDoc = svgOpen + `<style>circle{fill:red}</style>` +
	`<circle cx="60" cy="60" r="40" fill="red"></circle></svg>`

// TestServerParseRenderRoundTrip: the gRPC surface round-trips a document
// (including css seams). The css seam re-renders with its own smart spacing,
// so the invariant is the codec's fixed point — rendering is stable across a
// second parse — not byte-equality with the input.
func TestServerParseRenderRoundTrip(t *testing.T) {
	s := NewServer()
	parsed, err := s.Parse(context.Background(), &svgservicepb.ParseRequest{Svg: serverDoc})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if parsed.GetDocument() == nil {
		t.Fatal("Parse returned nil document")
	}
	rendered, err := s.Render(context.Background(), &svgservicepb.RenderRequest{Document: parsed.GetDocument()})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	for _, w := range []string{"<svg ", "<style>", "fill:red", "<circle ", "</svg>"} {
		if !strings.Contains(rendered.GetSvg(), w) {
			t.Errorf("rendered output missing %q: %s", w, rendered.GetSvg())
		}
	}
	reparsed, err := s.Parse(context.Background(), &svgservicepb.ParseRequest{Svg: rendered.GetSvg()})
	if err != nil {
		t.Fatalf("re-Parse: %v", err)
	}
	rerendered, err := s.Render(context.Background(), &svgservicepb.RenderRequest{Document: reparsed.GetDocument()})
	if err != nil {
		t.Fatalf("re-Render: %v", err)
	}
	if rerendered.GetSvg() != rendered.GetSvg() {
		t.Errorf("render not a fixed point:\n 1st: %s\n 2nd: %s", rendered.GetSvg(), rerendered.GetSvg())
	}
}

// TestServerParseAsSubtree: ParseRequest.type parses an arbitrary grammar
// subtree (here a lone <rect> element), returns it Any-packed in `node`, and
// RenderRequest.node renders it back byte-exactly. This is the path the
// gallery generator drives for every walked grammar path.
func TestServerParseAsSubtree(t *testing.T) {
	s := NewServer()
	const markup = `<rect x="1" y="1" width="10" height="10"></rect>`
	parsed, err := s.Parse(context.Background(), &svgservicepb.ParseRequest{Svg: markup, Type: "svg.SvgrectElement"})
	if err != nil {
		t.Fatalf("Parse(type): %v", err)
	}
	if parsed.GetDocument() != nil {
		t.Error("Parse(type) should not set document")
	}
	node := parsed.GetNode()
	if node == nil {
		t.Fatal("Parse(type) returned nil node")
	}
	if want := "type.googleapis.com/svg.SvgrectElement"; node.GetTypeUrl() != want {
		t.Errorf("node type_url = %q, want %q", node.GetTypeUrl(), want)
	}
	rendered, err := s.Render(context.Background(), &svgservicepb.RenderRequest{Node: node})
	if err != nil {
		t.Fatalf("Render(node): %v", err)
	}
	if rendered.GetSvg() != markup {
		t.Errorf("subtree round-trip:\n  in : %s\n  out: %s", markup, rendered.GetSvg())
	}

	if _, err := s.Parse(context.Background(), &svgservicepb.ParseRequest{Svg: markup, Type: "svg.NoSuchType"}); status.Code(err) != codes.InvalidArgument {
		t.Errorf("Parse(unknown type): got %v, want InvalidArgument", err)
	}
}

// renderStreamCollector implements SvgService_RenderStreamServer for tests;
// only Send is called by the server.
type renderStreamCollector struct {
	svgservicepb.SvgService_RenderStreamServer
	chunks []string
}

func (c *renderStreamCollector) Send(chunk *svgservicepb.RenderChunk) error {
	c.chunks = append(c.chunks, chunk.GetSvg())
	return nil
}

// TestServerRenderStream: chunk concatenation equals Render's svg.
func TestServerRenderStream(t *testing.T) {
	s := NewServer()
	parsed, err := s.Parse(context.Background(), &svgservicepb.ParseRequest{Svg: serverDoc})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	rendered, err := s.Render(context.Background(), &svgservicepb.RenderRequest{Document: parsed.GetDocument()})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	col := &renderStreamCollector{}
	if err := s.RenderStream(&svgservicepb.RenderRequest{Document: parsed.GetDocument()}, col); err != nil {
		t.Fatalf("RenderStream: %v", err)
	}
	if got := strings.Join(col.chunks, ""); got != rendered.GetSvg() {
		t.Errorf("stream mismatch:\n render: %s\n stream: %s", rendered.GetSvg(), got)
	}
}

// TestServerInvalidArgument: missing inputs surface as InvalidArgument.
func TestServerInvalidArgument(t *testing.T) {
	s := NewServer()
	if _, err := s.Render(context.Background(), &svgservicepb.RenderRequest{}); status.Code(err) != codes.InvalidArgument {
		t.Errorf("Render(nil document): got %v, want InvalidArgument", err)
	}
	if err := s.RenderStream(&svgservicepb.RenderRequest{}, &renderStreamCollector{}); status.Code(err) != codes.InvalidArgument {
		t.Errorf("RenderStream(nil document): got %v, want InvalidArgument", err)
	}
	if _, err := s.Parse(context.Background(), &svgservicepb.ParseRequest{Svg: "<not svg"}); status.Code(err) != codes.InvalidArgument {
		t.Errorf("Parse(garbage): got %v, want InvalidArgument", err)
	}
}
