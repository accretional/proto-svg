#!/usr/bin/env bash
# shoot.sh — serve the specimen gallery, start headless chromerpc, and screenshot
# every per-value specimen (static → one PNG; temporal → a frame folder → GIF).
# Idempotent: builds chromerpc/automate from GitHub into a cache, kills its
# servers on exit, and removes the per-run sequence chunks.
#
#   ONLY="rect,animate" ./chrome-testing/shoot.sh   # limit to some tags
#   RESUME=1 ./chrome-testing/shoot.sh              # skip specimens already shot
#   REBUILD_CHROMERPC=1 ./chrome-testing/shoot.sh   # force-rebuild chromerpc
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"   # go-run + textproto output paths below are repo-root-relative
CT="$ROOT/chrome-testing"
SCREENS="$CT/screenshots/gallery"
REL_SCREENS="chrome-testing/screenshots/gallery"
CATALOGUE="$CT/gallery/catalogue.json"

CHROMERPC_GIT="${CHROMERPC_GIT:-https://github.com/accretional/chromerpc}"
CACHE="/tmp/chromerpc-testing"
BIN="$CACHE/bin"
CHROMERPC_BIN="$BIN/chromerpc"
AUTOMATE_BIN="$BIN/automate"

ONLY="${ONLY:-}"
RESUME="${RESUME:-}"
RESTART_EVERY="${RESTART_EVERY:-40}"   # relaunch chromerpc every N chunks (Chrome leaks)
SEQDIR="$CT/generated/_seq"
SEQ="$SEQDIR/shots.textproto"

HTTP_PID=""
RPC_PID=""
RPC_PORT=""
cleanup() {
  [[ -n "$HTTP_PID" ]] && kill "$HTTP_PID" 2>/dev/null || true
  [[ -n "$RPC_PID" ]] && kill "$RPC_PID" 2>/dev/null || true
  rm -rf "$SEQDIR" 2>/dev/null || true
}
trap cleanup EXIT

free_port() { python3 -c "import socket;s=socket.socket();s.bind(('',0));print(s.getsockname()[1]);s.close()"; }

# start (or restart) the headless chromerpc server on a fresh port
start_rpc() {
  [[ -n "$RPC_PID" ]] && kill "$RPC_PID" 2>/dev/null || true
  RPC_PORT="$(free_port)"
  "$CHROMERPC_BIN" --headless --addr ":$RPC_PORT" &>"$CACHE/chromerpc.log" &
  RPC_PID=$!
  for i in $(seq 1 40); do
    if bash -c ">/dev/tcp/localhost/$RPC_PORT" 2>/dev/null; then return 0; fi
    sleep 0.5
  done
  echo "ERROR: chromerpc not ready"; cat "$CACHE/chromerpc.log"; return 1
}

# ── build chromerpc + automate from GitHub (reproducible) ─────────────────────
# Cloned + cached in /tmp/chromerpc-testing; rebuild with REBUILD_CHROMERPC=1.
if [[ ! -x "$CHROMERPC_BIN" || ! -x "$AUTOMATE_BIN" || -n "${REBUILD_CHROMERPC:-}" ]]; then
  echo "==> Fetching + building chromerpc from $CHROMERPC_GIT ..."
  mkdir -p "$BIN"
  if [[ -d "$CACHE/src/.git" ]]; then
    git -C "$CACHE/src" pull --quiet || true
  else
    git clone --quiet "$CHROMERPC_GIT" "$CACHE/src"
  fi
  ( cd "$CACHE/src" \
      && go build -o "$CHROMERPC_BIN" ./cmd/chromerpc \
      && go build -o "$AUTOMATE_BIN"  ./cmd/automate )
  echo "==> chromerpc ready: $CHROMERPC_BIN"
else
  echo "==> Using cached chromerpc: $CHROMERPC_BIN"
fi

# ── serve chrome-testing/ so /gallery/index.html + catalogue.json are reachable ─
SERVE_PORT="$(free_port)"
echo "==> Serving $CT on :$SERVE_PORT"
python3 -m http.server "$SERVE_PORT" --directory "$CT" &>/dev/null &
HTTP_PID=$!
BASE="http://localhost:$SERVE_PORT/gallery"

# ── build the automation sequence (chunked to keep gRPC responses small) ──────
mkdir -p "$SCREENS" "$SEQDIR"
rm -f "$SEQDIR"/shots-*.textproto

# Clean stale per-tag screenshot dirs first so values that no longer exist don't
# linger. ONLY → just those tags; otherwise every tag in catalogue.json. RESUME
# preserves existing output (shoot skips already-shot presets), so don't clean.
if [[ -z "$RESUME" ]]; then
  if [[ -n "$ONLY" ]]; then
    IFS=',' read -ra _clean_tags <<< "$ONLY"
    for _t in "${_clean_tags[@]}"; do _t="${_t// /}"; [[ -n "$_t" ]] && rm -rf "$SCREENS/$_t"; done
  else
    python3 -c "import json;[print(e['tag']) for e in json.load(open('$CATALOGUE'))['elements']]" \
      | while IFS= read -r _t; do [[ -n "$_t" ]] && rm -rf "$SCREENS/$_t"; done
  fi
fi

# Emit ROOT-relative output_path values in the textprotos (chromerpc + automate
# are launched below from $ROOT, so relative paths resolve to the same dir).
RESUME="$RESUME" go run ./chrome-testing/cmd/shoot/ \
  -catalogue "$CATALOGUE" \
  -base "$BASE" -outdir "$REL_SCREENS" -seq "$SEQ" \
  ${ONLY:+-only "$ONLY"}

# ── start headless chromerpc ──────────────────────────────────────────────────
echo "==> Starting chromerpc (headless)"
start_rpc || exit 1
echo "==> chromerpc ready on :$RPC_PORT"

# ── run each chunk; recycle the browser periodically and on failure ───────────
echo "==> Capturing screenshots (this can take a while) ..."
shopt -s nullglob
chunks=( "$SEQDIR"/shots-*.textproto )
total=${#chunks[@]}
[[ $total -eq 0 ]] && { echo "Nothing to capture (all present?)."; exit 0; }
i=0
for c in "${chunks[@]}"; do
  i=$((i+1))
  printf "\r    chunk %d/%d   " "$i" "$total"
  if (( i % RESTART_EVERY == 0 )); then start_rpc || exit 1; fi
  if ! "$AUTOMATE_BIN" -addr "localhost:$RPC_PORT" -input "$c" -timeout 180s >/dev/null 2>&1; then
    echo " (chunk $i failed — restarting chromerpc and retrying)"
    start_rpc || exit 1
    "$AUTOMATE_BIN" -addr "localhost:$RPC_PORT" -input "$c" -timeout 180s >/dev/null 2>&1 \
      || echo " (chunk $i failed again — skipping)"
  fi
done
echo ""

# ── encode temporal frame folders → GIFs ──────────────────────────────────────
echo "==> Encoding temporal frame sequences into GIFs..."
go run ./chrome-testing/cmd/gifenc/ -dir "$REL_SCREENS"

echo "==> Screenshots under $SCREENS"
find "$SCREENS" -name '*.png' | wc -l | xargs echo "Total PNGs:"
find "$SCREENS" -name '*.gif' | wc -l | xargs echo "Total GIFs:"
