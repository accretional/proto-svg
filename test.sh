#!/usr/bin/env bash
# test.sh — The ONLY sanctioned pre-commit validation for proto-svg.
#
# Runs the full build, then asserts the pipeline's outputs are present, correct,
# and well-formed:
#   - ./build.sh                              (setup -> genproto -> gen)
#   - go vet ./...                            (static analysis)
#   - go test -count=1 ./...                  (unit tests; skipped if none exist)
#   - proto/svg.proto + proto/svg.fdset exist and are non-empty
#   - chrome-testing/html/generated/index.html exists with >0 per-element pages
#   - every hand-authored template starts with <!DOCTYPE
#   - a sample of generated pages is well-formed
#
# Never use a bare `go test` as final validation — always run THIS script. See
# CLAUDE.md.
#
# Idempotent: safe to re-run. Exits non-zero on any failure.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

PASS=0
FAILURES=0
pass() { echo "  [PASS] $1"; PASS=$((PASS + 1)); }
fail() { echo "  [FAIL] $1" >&2; FAILURES=$((FAILURES + 1)); }

echo "########################################"
echo "#  test.sh — validate (pre-commit)      #"
echo "########################################"

# ── Full build (setup -> genproto -> gen) ───────────────────────────────────
echo ""
echo "============ Build ============"
"$ROOT/build.sh"

# ── go vet ──────────────────────────────────────────────────────────────────
echo ""
echo "============ go vet ============"
if [[ -f "$ROOT/go.mod" ]]; then
  if go vet ./...; then
    pass "go vet clean"
  else
    fail "go vet found issues"
  fi
else
  echo "  [skip] no go.mod"
fi

# ── go test (skip gracefully if no tests exist) ─────────────────────────────
echo ""
echo "============ go test ============"
if [[ -f "$ROOT/go.mod" ]]; then
  TEST_FILES=$(find "$ROOT" -name '*_test.go' -not -path '*/.*' | wc -l | tr -d ' ')
  if [[ "$TEST_FILES" -eq 0 ]]; then
    echo "  [skip] no *_test.go files found"
  elif go test -count=1 ./...; then
    pass "go test passed ($TEST_FILES test file(s))"
  else
    fail "go test failed"
  fi
else
  echo "  [skip] no go.mod"
fi

# ── Proto schema artifacts ──────────────────────────────────────────────────
echo ""
echo "============ Proto schema artifacts ============"
check_nonempty() {
  local path="$1" name="$2"
  if [[ -f "$path" ]]; then
    if [[ -s "$path" ]]; then
      pass "$name exists and is non-empty"
    else
      fail "$name exists but is empty"
    fi
  else
    fail "$name not found — run ./build.sh"
  fi
}
check_nonempty "$ROOT/proto/svg.proto"           "proto/svg.proto"
check_nonempty "$ROOT/proto/svg.fdset"           "proto/svg.fdset"
check_nonempty "$ROOT/proto/pb/svg/prefix_map.go"    "proto/pb/svg/prefix_map.go"
check_nonempty "$ROOT/proto/pb/svg/separator_map.go" "proto/pb/svg/separator_map.go"

# ── Generated gallery ───────────────────────────────────────────────────────
echo ""
echo "============ Generated gallery ============"
GEN_DIR="$ROOT/chrome-testing/html/generated"
if [[ -f "$GEN_DIR/index.html" ]]; then
  pass "generated index.html exists"
else
  fail "generated index.html not found — run ./build.sh"
fi

# Per-element pages = generated *.html minus index.html.
GEN_PAGES=0
if [[ -d "$GEN_DIR" ]]; then
  GEN_PAGES=$(find "$GEN_DIR" -maxdepth 1 -name '*.html' ! -name 'index.html' | wc -l | tr -d ' ')
fi
if [[ "$GEN_PAGES" -gt 0 ]]; then
  pass "$GEN_PAGES per-element gallery page(s) generated"
else
  fail "no per-element gallery pages in $GEN_DIR — run ./build.sh"
fi

# ── Hand-authored templates start with <!DOCTYPE ────────────────────────────
echo ""
echo "============ Hand-authored templates ============"
TPL_DIR="$ROOT/chrome-testing/html/template"
if [[ -d "$TPL_DIR" ]]; then
  TPL_COUNT=$(find "$TPL_DIR" -maxdepth 1 -name '*.html' | wc -l | tr -d ' ')
  if [[ "$TPL_COUNT" -gt 0 ]]; then
    pass "$TPL_COUNT hand-authored template(s) found"
  else
    fail "no hand-authored templates in $TPL_DIR"
  fi

  BAD_TPL=0
  while IFS= read -r -d '' tpl; do
    # head -c keeps this robust against leading BOM/whitespace.
    if ! head -c 200 "$tpl" | grep -qi '<!DOCTYPE'; then
      echo "    missing <!DOCTYPE: $(basename "$tpl")" >&2
      BAD_TPL=$((BAD_TPL + 1))
    fi
  done < <(find "$TPL_DIR" -maxdepth 1 -name '*.html' -print0)

  if [[ $BAD_TPL -eq 0 ]]; then
    pass "every hand-authored template starts with <!DOCTYPE"
  else
    fail "$BAD_TPL hand-authored template(s) do not start with <!DOCTYPE"
  fi
else
  fail "hand-authored template dir not found: $TPL_DIR"
fi

# ── Well-formedness of a sample of generated pages ──────────────────────────
# Lightweight, dependency-free balance check: each sampled page must contain a
# DOCTYPE, balanced <html>…</html>, and the grammar-injected SVG markup.
echo ""
echo "============ Generated page well-formedness (sample) ============"
wellformed() {
  local f="$1"
  head -c 200 "$f" | grep -qi '<!DOCTYPE'  || { echo "    no DOCTYPE: $(basename "$f")" >&2; return 1; }
  grep -qi '<html'  "$f"                    || { echo "    no <html>: $(basename "$f")"  >&2; return 1; }
  grep -qi '</html>' "$f"                   || { echo "    no </html>: $(basename "$f")" >&2; return 1; }
  grep -q  '<svg'   "$f"                     || { echo "    no <svg>: $(basename "$f")"   >&2; return 1; }
  return 0
}

if [[ "$GEN_PAGES" -gt 0 ]]; then
  # Sample up to 8 pages deterministically (sorted, evenly spaced). Built with a
  # while-read loop, not mapfile, so it works on bash 3.2 (macOS default).
  SAMPLE_COUNT=0
  BAD_GEN=0
  while IFS= read -r f; do
    [[ -z "$f" ]] && continue
    SAMPLE_COUNT=$((SAMPLE_COUNT + 1))
    wellformed "$f" || BAD_GEN=$((BAD_GEN + 1))
  done < <(find "$GEN_DIR" -maxdepth 1 -name '*.html' ! -name 'index.html' | sort | awk 'NR%8==1' | head -8)
  if [[ $BAD_GEN -eq 0 ]]; then
    pass "$SAMPLE_COUNT sampled generated page(s) well-formed"
  else
    fail "$BAD_GEN of $SAMPLE_COUNT sampled generated page(s) malformed"
  fi
else
  echo "  [skip] no generated pages to sample"
fi

# ── Summary ─────────────────────────────────────────────────────────────────
echo ""
echo "########################################"
if [[ $FAILURES -eq 0 ]]; then
  echo "#  test.sh PASSED — $PASS check(s) green"
  echo "########################################"
  exit 0
else
  echo "#  test.sh FAILED — $FAILURES failure(s)" >&2
  echo "########################################"
  exit 1
fi
