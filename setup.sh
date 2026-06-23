#!/usr/bin/env bash
# setup.sh — Check prerequisites, warm the chromerpc build, and tidy Go modules.
#
# This is part of the proto-svg orchestration. ALL build/test/run goes through
# the four root scripts (setup.sh / build.sh / test.sh / LET_IT_RIP.sh). Never
# use a bare `go run` / `go test` as final validation — see CLAUDE.md.
#
# Idempotent: safe to re-run at any time. chromerpc setup is best-effort and
# never fatal (screenshots are a downstream concern).

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

echo "========================================="
echo "  setup.sh — proto-svg setup"
echo "========================================="

OK=0
FAIL=0
ok()   { echo "  [ok]   $1"; OK=$((OK + 1)); }
warn() { echo "  [warn] $1"; }
bad()  { echo "  [FAIL] $1" >&2; FAIL=$((FAIL + 1)); }

# ── Prerequisites ───────────────────────────────────────────────────────────

echo ""
echo "--- Checking prerequisites ---"

check_cmd() {
  local name="$1" hint="$2"
  if command -v "$name" &>/dev/null; then
    ok "$name found: $(command -v "$name")"
  else
    bad "$name not installed — $hint"
  fi
}

check_cmd go      "Install from https://go.dev/dl/"
check_cmd python3 "Install Python 3 from https://python.org"
check_cmd git     "Install git from https://git-scm.com"

# Google Chrome detection (same logic as chrome-testing/snap.sh).
CHROME=""
if command -v google-chrome &>/dev/null; then
  CHROME="$(command -v google-chrome)"
elif command -v chromium-browser &>/dev/null; then
  CHROME="$(command -v chromium-browser)"
elif command -v chromium &>/dev/null; then
  CHROME="$(command -v chromium)"
elif [[ -x "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" ]]; then
  CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
fi
if [[ -n "$CHROME" ]]; then
  ok "Google Chrome found: $CHROME"
else
  bad "Google Chrome not found — install from https://www.google.com/chrome/"
fi

if [[ $FAIL -ne 0 ]]; then
  echo ""
  echo "FATAL: missing prerequisites. Install them and re-run." >&2
  exit 1
fi

# ── Go modules ──────────────────────────────────────────────────────────────
# Module github.com/accretional/proto-svg, with replace directives to the
# sibling checkouts ../gluon and ../proto-merge. go mod tidy resolves them.

echo ""
echo "--- Tidying Go modules ---"
if [[ -f "$ROOT/go.mod" ]]; then
  if go mod tidy; then
    ok "go mod tidy complete"
  else
    bad "go mod tidy failed (check ../gluon and ../proto-merge sibling checkouts)"
  fi
else
  bad "no go.mod found at repo root"
fi

# ── chromerpc warmup (best-effort, NON-FATAL) ───────────────────────────────
# chrome-testing/snap.sh builds chromerpc+automate FROM GitHub and caches the
# binaries in /tmp/chromerpc-testing/bin. We pre-build them here using the same
# clone+build so the first real screenshot is fast. Failure here is fine — the
# grammar/proto/gallery build does not depend on it.

echo ""
echo "--- Warming chromerpc (screenshots; best-effort) ---"
CHROMERPC_GIT="${CHROMERPC_GIT:-https://github.com/accretional/chromerpc}"
CACHE="/tmp/chromerpc-testing"
CHROMERPC_BIN="$CACHE/bin/chromerpc"
AUTOMATE_BIN="$CACHE/bin/automate"

if [[ -x "$CHROMERPC_BIN" && -x "$AUTOMATE_BIN" && -z "${REBUILD_CHROMERPC:-}" ]]; then
  ok "chromerpc already built: $CHROMERPC_BIN"
else
  if (
    set -e
    mkdir -p "$CACHE/bin"
    if [[ -d "$CACHE/src/.git" ]]; then
      git -C "$CACHE/src" pull --quiet || true
    else
      git clone --quiet "$CHROMERPC_GIT" "$CACHE/src"
    fi
    cd "$CACHE/src"
    go build -o "$CHROMERPC_BIN" ./cmd/chromerpc
    go build -o "$AUTOMATE_BIN"  ./cmd/automate
  ); then
    ok "chromerpc + automate built (cached in $CACHE/bin)"
  else
    warn "chromerpc warmup failed — screenshots may be unavailable, but the build/test pipeline is unaffected"
  fi
fi

# ── EBNF grammar present ────────────────────────────────────────────────────
# The grammar is the source of truth: lang/*.ebnf (14 modules, svg.ebnf is the
# root file). genproto and gen concatenate these in order.

echo ""
echo "--- Checking EBNF grammar (source of truth) ---"
LANG_DIR="$ROOT/lang"
if [[ -d "$LANG_DIR" ]]; then
  EBNF_COUNT=$(find "$LANG_DIR" -maxdepth 1 -name '*.ebnf' | wc -l | tr -d ' ')
  if [[ "$EBNF_COUNT" -gt 0 ]]; then
    ok "$EBNF_COUNT EBNF grammar files found in lang/"
    if [[ ! -f "$LANG_DIR/svg.ebnf" ]]; then
      bad "lang/svg.ebnf (the root grammar file) is missing"
    fi
  else
    bad "no *.ebnf files in lang/"
  fi
else
  bad "no lang/ directory found"
fi

# ── Generators build ────────────────────────────────────────────────────────
# Verify the two Go entry points compile: the grammar compiler (genproto) and
# the gallery generator (gen). A quick compile check here catches breakage
# early; the real run happens through build.sh.

echo ""
echo "--- Checking Go generators build ---"
if go build ./lang/cmd/genproto/ ./chrome-testing/cmd/gen/; then
  ok "genproto + gen build successfully"
else
  bad "go build of ./lang/cmd/genproto/ and/or ./chrome-testing/cmd/gen/ failed"
fi

# ── Done ────────────────────────────────────────────────────────────────────

echo ""
echo "========================================="
if [[ $FAIL -eq 0 ]]; then
  echo "  setup.sh OK — $OK check(s) passed"
  echo "========================================="
  exit 0
else
  echo "  setup.sh FAILED — $FAIL check(s) failed" >&2
  echo "========================================="
  exit 1
fi
