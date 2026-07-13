// Command server runs the SvgService gRPC server.
//
// The CSS and HTML service packages are linked so the codec descends into
// svg's cross-grammar seams (CSS in <style> and presentation-attribute
// values, HTML flow content inside <foreignObject>), and their proto
// descriptors land in the global registry — gRPC reflection therefore serves
// the css.* and html.* schemas alongside svg.*, even though svg.proto itself
// only imports any.proto.
package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	svgservicepb "github.com/accretional/proto-svg/proto/pb/svgservice"
	"github.com/accretional/proto-svg/service"

	// Link the CSS and HTML grammars (seam descent + reflectable schemas).
	_ "github.com/accretional/proto-css/service"
	_ "github.com/accretional/proto-html/service"
)

func main() {
	addr := flag.String("addr", ":50053", "listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen %s: %v", *addr, err)
	}

	s := grpc.NewServer(grpc.MaxRecvMsgSize(64*1024*1024), grpc.MaxSendMsgSize(64*1024*1024))
	svgservicepb.RegisterSvgServiceServer(s, service.NewServer())
	reflection.Register(s) // enables grpc reflection / discovery over the schema

	log.Printf("SvgService listening on %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
