#!/usr/bin/env bash
# LET_IT_RIP.sh — The full proto-svg pipeline, end to end:
#   setup -> build (EBNF -> proto -> gallery) -> test -> screenshot -> serve.
#
#   1. ./setup.sh                    prereqs + go mod tidy + chromerpc warmup
#   2. ./build.sh                    EBNF grammar -> proto -> per-element gallery
#   3. ./test.sh                     the sanctioned validation suite
#   4. chrome-testing/snap.sh        screenshot every generated gallery page
#                                    (directory mode: one PNG per element page)
#   5. python3 -m http.server        serve the gallery and print the URL
#
# CRITICAL: Run this before EVERY `git commit` and `git push`. No exceptions.
# It is the single command that proves the whole pipeline is green.
#
# Env toggles:
#   SKIP_SHOTS=1   skip the gallery screenshot step (steps 1-3, then serve)
#   SHOOT=1        also capture ONE PNG per attribute value of every element +
#                  animated GIFs for SMIL elements (chrome-testing/shoot.sh).
#                  Long-running (~thousands of captures); off by default.
#   SKIP_SERVE=1   stop after screenshots; do not start the HTTP server
#   SERVE_PORT=N   gallery server port (default 8899)
#
# Idempotent: kills stale servers on the port before binding, and cleans up the
# server it starts on exit (trap).

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

CT="$ROOT/chrome-testing"
GALLERY_DIR="$CT/gallery"
SERVE_PORT="${SERVE_PORT:-8899}"

echo ""
echo "############################################"
echo "#                LET IT RIP                #"
echo "#   setup · build · test · shoot · serve   #"
echo "############################################"
echo ""

# ── cleanup / server tracking ───────────────────────────────────────────────
SERVE_PID=""
cleanup() {
  if [[ -n "$SERVE_PID" ]]; then
    kill "$SERVE_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT INT TERM

# kill_port <port> — terminate any process already listening on the port so the
# script is re-runnable without "address already in use".
kill_port() {
  local port="$1" pids
  pids="$(lsof -ti tcp:"$port" 2>/dev/null || true)"
  if [[ -n "$pids" ]]; then
    echo "  killing stale process(es) on port $port: $pids"
    kill $pids 2>/dev/null || true
    sleep 1
  fi
}

# ── 1/5 Setup ───────────────────────────────────────────────────────────────
echo "============ Step 1/5: Setup ============"
"$ROOT/setup.sh"
echo ""

# ── 2/5 Build (EBNF -> proto -> gallery) ────────────────────────────────────
echo "============ Step 2/5: Build (EBNF -> proto -> gallery) ============"
"$ROOT/build.sh"
echo ""

# ── 3/5 Test (the sanctioned validation) ────────────────────────────────────
echo "============ Step 3/5: Test (validation) ============"
"$ROOT/test.sh"
echo ""

# ── 4/5 Per-preset gallery shoot + GIFs (opt-in; long-running) ──────────────
# The gallery is one live SPA; screenshots are captured by driving its viewer to
# each preset (= each attribute value) via chrome-testing/shoot.sh, not by
# snapping per-element pages.
echo "============ Step 4/5: Per-preset gallery shoot ============"
if [[ -n "${SHOOT:-}" ]]; then
  "$CT/shoot.sh"
else
  echo "  Skipped. Set SHOOT=1 to drive the gallery viewer through every preset"
  echo "  (one PNG per attribute value + animated GIFs for SMIL elements) into"
  echo "  screenshots/gallery/."
fi
echo ""

# ── 5/5 Serve the gallery ───────────────────────────────────────────────────
if [[ -n "${SKIP_SERVE:-}" ]]; then
  echo "############################################"
  echo "#  LET IT RIP complete (serve skipped)     #"
  echo "############################################"
  exit 0
fi

echo "============ Step 5/5: Serve gallery ============"
if [[ ! -f "$GALLERY_DIR/index.html" ]]; then
  echo "ERROR: $GALLERY_DIR/index.html not found — cannot serve." >&2
  exit 1
fi

kill_port "$SERVE_PORT"
python3 -m http.server "$SERVE_PORT" --directory "$GALLERY_DIR" &
SERVE_PID=$!

# Wait for the server to accept connections.
for _ in $(seq 1 20); do
  if bash -c ">/dev/tcp/localhost/$SERVE_PORT" 2>/dev/null; then
    break
  fi
  sleep 0.25
done

echo ""
echo "############################################"
echo "#  LET IT RIP complete — gallery is live   #"
echo "############################################"
echo ""
echo "  SVG Lab gallery:    http://localhost:$SERVE_PORT/index.html"
echo "  Serving directory:  $GALLERY_DIR"
echo "  Press Ctrl-C to stop the server."
echo ""

# Keep the server in the foreground until interrupted.
wait "$SERVE_PID"
