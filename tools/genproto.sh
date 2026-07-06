#!/usr/bin/env bash
# genproto.sh — compile the SVG EBNF grammar into the protobuf schema, the
# grammar-derived lookup tables (prefix_map.go / separator_map.go), and the Go
# bindings (proto/pb/svg/svg.pb.go). Idempotent: re-running regenerates the
# committed proto/ artifacts in place.
#
# The CSS seam is wired here. Driven by grammar_deps.bzl, genproto's externalize
# pass makes svg.proto `import "css.proto"` and retype the <style> content field
# to css.CssStyleSheet; protoc's Mcss.proto mapping makes svg.pb.go import the
# proto-css csspb package instead of regenerating CSS. proto-css's compiled proto
# (css.proto + css.fdset) is read from a sibling checkout (override CSS_PROTO_DIR).
#
# This is the ONLY sanctioned way to compile/run Go in this repo.
set -euo pipefail

# Resolve the repo root (this script lives in <root>/tools/).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

CSS_PROTO_DIR="${CSS_PROTO_DIR:-../proto-css/proto}"
export PATH="$PATH:$(go env GOPATH)/bin"

# (a) Resolve module dependencies.
go mod tidy

# (b) Ensure the generated-Go output package directory exists.
mkdir -p proto/pb/svg

# (c) Compile the grammar (externalizes the css seam per grammar_deps.bzl).
go run ./lang/cmd/genproto/ \
	-lang lang \
	-deps grammar_deps.bzl \
	-bundled proto/svg.proto \
	-fdset proto/svg.fdset \
	-prefix-map proto/pb/svg/prefix_map.go \
	-separator-map proto/pb/svg/separator_map.go \
	-package svg \
	-go-package 'github.com/accretional/proto-svg/proto/pb/svg;svgpb'

# (d) protoc: svg.proto -> Go (css.proto -> csspb).
protoc -I proto -I "$CSS_PROTO_DIR" \
	--go_out=. --go_opt=module=github.com/accretional/proto-svg \
	--go_opt=Mcss.proto=github.com/accretional/proto-css/proto/pb/css \
	proto/svg.proto

# (e) Sanity-check the generated packages compile.
go build ./proto/... ./service/...
