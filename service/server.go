package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	svgservicepb "github.com/accretional/proto-svg/proto/pb/svgservice"
)

// streamChunkBytes bounds each RenderStream chunk so a single gRPC message
// stays small (keep streamed payloads well under 4MB).
const streamChunkBytes = 64 * 1024

// Server implements svgservicepb.SvgServiceServer. It holds no per-element
// state — Render and Parse are reflection over the generated svg schema and
// tables, rooted at SvgDocument. Cross-grammar seams (CSS paint values,
// <foreignObject> HTML) descend when those grammars' service packages are
// linked into the binary.
type Server struct {
	svgservicepb.UnimplementedSvgServiceServer
}

// NewServer returns a ready SvgService server.
func NewServer() *Server { return &Server{} }

// renderRoot resolves a RenderRequest's root message: the Any-packed subtree
// in `node` when set (as produced by Parse with a `type`), else `document`.
func renderRoot(req *svgservicepb.RenderRequest) (proto.Message, error) {
	if n := req.GetNode(); n != nil {
		msg, err := n.UnmarshalNew()
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "node: %v", err)
		}
		return msg, nil
	}
	if doc := req.GetDocument(); doc != nil {
		return doc, nil
	}
	return nil, status.Error(codes.InvalidArgument, "document or node is required")
}

// Render serializes the request's SvgDocument (or Any-packed subtree) to SVG text.
func (s *Server) Render(_ context.Context, req *svgservicepb.RenderRequest) (*svgservicepb.RenderResponse, error) {
	root, err := renderRoot(req)
	if err != nil {
		return nil, err
	}
	svg, err := Render(root)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "render: %v", err)
	}
	return &svgservicepb.RenderResponse{Svg: svg}, nil
}

// RenderStream renders the document and streams the SVG text in chunks.
func (s *Server) RenderStream(req *svgservicepb.RenderRequest, stream svgservicepb.SvgService_RenderStreamServer) error {
	root, err := renderRoot(req)
	if err != nil {
		return err
	}
	svg, err := Render(root)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "render: %v", err)
	}
	for i := 0; i < len(svg); i += streamChunkBytes {
		end := i + streamChunkBytes
		if end > len(svg) {
			end = len(svg)
		}
		if err := stream.Send(&svgservicepb.RenderChunk{Svg: svg[i:end]}); err != nil {
			return err
		}
	}
	return nil
}

// Parse parses SVG text into an SvgDocument, or — when req.Type names a
// non-root grammar message — into that subtree, returned Any-packed in `node`.
func (s *Server) Parse(_ context.Context, req *svgservicepb.ParseRequest) (*svgservicepb.ParseResponse, error) {
	if typ := req.GetType(); typ != "" {
		msg, err := ParseAs(req.GetSvg(), typ)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "parse as %s: %v", typ, err)
		}
		node, err := anypb.New(msg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "pack %s: %v", typ, err)
		}
		return &svgservicepb.ParseResponse{Node: node}, nil
	}
	doc, err := Parse(req.GetSvg())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "parse: %v", err)
	}
	return &svgservicepb.ParseResponse{Document: doc}, nil
}
