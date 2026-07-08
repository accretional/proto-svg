#!/usr/bin/env bash
# test.sh — The ONLY sanctioned pre-commit validation for proto-svg.
#
# Runs the full build, then asserts the pipeline's outputs are present, correct,
# and well-formed:
#   - ./build.sh                              (setup -> genproto -> gen)
#   - codec round-trip gate                   (chrome-testing/generated/
#                                              _codec_failures.tsv must be empty)
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

# ── Codec round-trip (every walked gallery path must be faithful) ───────────
echo ""
echo "============ Codec round-trip ============"
CODEC_TSV="$ROOT/chrome-testing/generated/_codec_failures.tsv"
if [[ -f "$CODEC_TSV" ]]; then
  CODEC_ROWS=$(( $(wc -l <"$CODEC_TSV") - 1 ))
  if [[ "$CODEC_ROWS" -le 0 ]]; then
    pass "codec round-trip: every walked gallery path renders faithfully"
  else
    fail "codec round-trip: $CODEC_ROWS failing path(s) — see $CODEC_TSV"
  fi
else
  fail "codec failures report missing: $CODEC_TSV (build did not run the walk)"
fi

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

# ── Gallery app + generated catalogue ───────────────────────────────────────
echo ""
echo "============ Gallery app + catalogue ============"
GALLERY_DIR="$ROOT/chrome-testing/gallery"
for f in index.html app.js; do
  if [[ -f "$GALLERY_DIR/$f" ]]; then
    pass "gallery/$f exists"
  else
    fail "gallery/$f not found"
  fi
done

CATALOGUE="$GALLERY_DIR/catalogue.json"
CAT_ELS=0
if [[ -f "$CATALOGUE" ]]; then
  CAT_ELS=$(python3 -c "import json,sys;print(len(json.load(open('$CATALOGUE')).get('elements',[])))" 2>/dev/null || echo 0)
fi
if [[ "$CAT_ELS" -gt 0 ]]; then
  pass "catalogue.json parses with $CAT_ELS element(s)"
else
  fail "catalogue.json missing/empty/invalid — run ./build.sh"
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

# ── Well-formedness of the catalogue base SVGs ──────────────────────────────
# Every element's base render must be a balanced <svg>…</svg> the gallery can
# mount. Checked in one pass over catalogue.json.
echo ""
echo "============ Catalogue base SVG well-formedness ============"
if [[ "$CAT_ELS" -gt 0 ]]; then
  if python3 - "$CATALOGUE" <<'PY'
import json, sys, re
els = json.load(open(sys.argv[1]))["elements"]
bad = []
for e in els:
    b = e.get("base", "")
    if b.count("<svg") < 1 or b.count("</svg>") < 1 or "{{ELEMENT}}" in b:
        bad.append(e["tag"])
    if not e.get("attrs") or not e.get("presets"):
        bad.append(e["tag"] + "(empty)")
if bad:
    print("    malformed/empty:", ", ".join(sorted(set(bad)))); sys.exit(1)
print("    %d catalogue base SVGs well-formed" % len(els))
PY
  then
    pass "all $CAT_ELS catalogue base SVGs well-formed with controls + presets"
  else
    fail "some catalogue entries are malformed/empty"
  fi
else
  echo "  [skip] no catalogue to check"
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
