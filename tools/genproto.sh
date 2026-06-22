#!/usr/bin/env bash
# genproto.sh — compile the SVG EBNF grammar into the protobuf schema and the
# grammar-derived lookup tables (prefix_map.go / separator_map.go). Idempotent:
# re-running regenerates the committed proto/ artifacts in place.
#
# This is the ONLY sanctioned way to compile/run Go in this repo.
set -euo pipefail

# Resolve the repo root (this script lives in <root>/tools/).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

# (a) Resolve module dependencies.
go mod tidy

# (b) Ensure the generated-Go output package directory exists.
mkdir -p proto/pb/svg

# (c) Compile the grammar with the default flags.
go run ./lang/cmd/genproto/ \
	-lang lang \
	-bundled proto/svg.proto \
	-fdset proto/svg.fdset \
	-prefix-map proto/pb/svg/prefix_map.go \
	-separator-map proto/pb/svg/separator_map.go \
	-package svg \
	-go-package 'github.com/accretional/proto-svg/proto/pb/svg;svgpb'
