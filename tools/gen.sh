#!/usr/bin/env bash
# gen.sh — run the SVG generation ENGINE: compile the SVG EBNF grammar in
# memory, walk the proto message graph, and emit valid SVG markup into
# chrome-testing/generated/ (sample-document.svg, sample-rect.svg).
#
# This is the ONLY sanctioned way to build/run the gen command in proto-svg.
set -euo pipefail

# Resolve the repo root (this script lives in <root>/tools/).
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

# Resolve module dependencies, then run the generator.
go mod tidy
go run ./chrome-testing/cmd/gen/ -lang lang -out chrome-testing/generated
