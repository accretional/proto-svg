#!/usr/bin/env bash
# serve.sh — serve the SVG Lab gallery locally and print the URL.
#
# Standalone (no build): just hosts chrome-testing/gallery/ over HTTP so you can
# open the live element explorer in a browser. The gallery reads catalogue.json,
# which ./tools/gen.sh generates — run ./build.sh first if it is missing.
#
#   ./chrome-testing/serve.sh            # serve on :8899
#   SERVE_PORT=3000 ./chrome-testing/serve.sh
#
# Idempotent: kills any stale listener on the port before binding.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
GALLERY="$ROOT/chrome-testing/gallery"
PORT="${SERVE_PORT:-8899}"

if [[ ! -f "$GALLERY/index.html" ]]; then
  echo "ERROR: $GALLERY/index.html not found." >&2
  exit 1
fi
if [[ ! -f "$GALLERY/catalogue.json" ]]; then
  echo "WARN: catalogue.json not found — run ./build.sh to generate it." >&2
fi

# Free the port if something is already bound to it.
pids="$(lsof -ti tcp:"$PORT" 2>/dev/null || true)"
if [[ -n "$pids" ]]; then
  echo "  killing stale listener(s) on :$PORT: $pids"
  kill $pids 2>/dev/null || true
  sleep 1
fi

echo ""
echo "  SVG Lab gallery  →  http://localhost:$PORT/"
echo "  serving:            $GALLERY"
echo "  Ctrl-C to stop."
echo ""

exec python3 -m http.server "$PORT" --directory "$GALLERY"
