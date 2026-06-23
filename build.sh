#!/usr/bin/env bash
# build.sh — The EBNF → proto → gallery build for proto-svg.
#
# Pipeline (the grammar is the source of truth):
#   1. ./setup.sh            — prereqs, go mod tidy, generator compile check
#   2. ./tools/genproto.sh   — compile lang/*.ebnf via gluon genproto into
#                              proto/svg.proto, proto/svg.fdset, and the
#                              grammar-derived lookup tables in proto/pb/svg/
#   3. ./tools/gen.sh        — walk the compiled grammar and emit every
#                              element's every value-path into the per-element
#                              gallery pages under chrome-testing/html/generated/
#
# This is the ONLY sanctioned way to build proto-svg. Never run the underlying
# `go run` commands directly as final validation — see CLAUDE.md.
#
# Idempotent: safe to re-run; regenerates committed proto/ and gallery artifacts
# in place.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

echo "########################################"
echo "#  build.sh — EBNF -> proto -> gallery  #"
echo "########################################"

# ── 1/3 Setup (idempotent) ──────────────────────────────────────────────────
echo ""
echo "============ Step 1/3: Setup ============"
"$ROOT/setup.sh"

# ── 2/3 Compile the grammar into the proto schema ───────────────────────────
echo ""
echo "============ Step 2/3: Compile grammar (genproto) ============"
"$ROOT/tools/genproto.sh"

# ── 3/3 Generate the all-value-paths galleries ──────────────────────────────
echo ""
echo "============ Step 3/3: Generate galleries (gen) ============"
"$ROOT/tools/gen.sh"

echo ""
echo "########################################"
echo "#         build.sh complete            #"
echo "########################################"
