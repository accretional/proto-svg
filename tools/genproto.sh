#!/usr/bin/env bash
# genproto.sh — compile the SVG EBNF grammar into the protobuf schema, the
# grammar-derived lookup tables (prefix_map / separator_map / seam_map), and the
# Go bindings (proto/pb/svg/svg.pb.go). Idempotent: re-running regenerates the
# committed proto/ artifacts in place.
#
# Cross-grammar seams (CSS in <style> and presentation-attribute values) are
# google.protobuf.Any: genproto's externalize pass (driven by grammar_deps.bzl)
# rewrites each seam field to Any and records the embedded type in seam_map.go.
# svg.proto therefore imports only any.proto — never css.proto — so SVG is a
# standalone grammar. The gluon codec descends into the seams at runtime when the
# CSS grammar is linked.
#
# This is the ONLY sanctioned way to compile/run Go in this repo.
set -euo pipefail

# Resolve the repo root (this script lives in <root>/tools/).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

export PATH="$PATH:$(go env GOPATH)/bin"

# (a) Resolve module dependencies.
go mod tidy

# (b) Ensure the generated-Go output package directory exists.
mkdir -p proto/pb/svg

# (c) Compile the grammar (rewrites cross-grammar seams to Any per grammar_deps.bzl).
go run ./lang/cmd/genproto/ \
	-lang lang \
	-deps grammar_deps.bzl \
	-bundled proto/svg.proto \
	-fdset proto/svg.fdset \
	-prefix-map proto/pb/svg/prefix_map.go \
	-separator-map proto/pb/svg/separator_map.go \
	-seam-map proto/pb/svg/seam_map.go \
	-package svg \
	-go-package 'github.com/accretional/proto-svg/proto/pb/svg;svgpb' \
	-service-proto proto/svg_service.proto \
	-service-go-package 'github.com/accretional/proto-svg/proto/pb/svgservice;svgservicepb'

# (d) protoc: svg.proto -> Go. any.proto is a well-known type; no other imports.
protoc -I proto \
	--go_out=. --go_opt=module=github.com/accretional/proto-svg \
	proto/svg.proto

# (d') protoc: svg_service.proto -> Go + gRPC (the repo-owned SvgService surface).
protoc -I proto \
	--go_out=. --go_opt=module=github.com/accretional/proto-svg \
	--go-grpc_out=. --go-grpc_opt=module=github.com/accretional/proto-svg \
	proto/svg_service.proto

# (e) Sanity-check the generated packages compile.
go build ./proto/...
