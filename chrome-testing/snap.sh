#!/usr/bin/env bash
# snap.sh — screenshot HTML (or URLs) via chromerpc, using the LOCAL ../chromerpc.
#
# Usage:
#   ./chrome-testing/snap.sh <input.html> <output.png>     # single file
#   ./chrome-testing/snap.sh <input_dir/> <output_dir/>    # one PNG per .html
#   ./chrome-testing/snap.sh <https://url> <out.png>       # external URL
#   AUTOMATION=<file.textproto> ./chrome-testing/snap.sh <input.html> <out.png>
#       # drive a custom interaction sequence (click/type/hover/wait) instead of
#       # the default navigate+screenshot; {{URL}}/{{OUT}} are substituted.
#
# Builds chromerpc + automate from the local sibling checkout ../chromerpc
# (go.mod + cmd/) into /tmp/chromerpc-proto-svg/bin. Idempotent.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CHROMERPC_GIT="${CHROMERPC_GIT:-https://github.com/accretional/chromerpc}"

INPUT="${1:-}"
OUTPUT="${2:-}"
if [[ -z "$INPUT" || -z "$OUTPUT" ]]; then
  echo "usage: $0 <input.html|dir|url> <output.png|dir>" >&2
  exit 1
fi

# ── cleanup ──────────────────────────────────────────────────────────────────
WORK_DIR=""; CHROMERPC_PID=""; HTTP_PID=""; HTTP_PORT=""; CHROMERPC_ADDR=""
cleanup() {
  [[ -n "$CHROMERPC_PID" ]] && kill "$CHROMERPC_PID" 2>/dev/null || true
  [[ -n "$HTTP_PID" ]]      && kill "$HTTP_PID"      2>/dev/null || true
  [[ -n "$WORK_DIR" ]]      && rm -rf "$WORK_DIR"
}
trap cleanup EXIT

# ── chrome detection ─────────────────────────────────────────────────────────
find_chrome() {
  command -v google-chrome 2>/dev/null && return
  command -v chromium-browser 2>/dev/null && return
  command -v chromium 2>/dev/null && return
  local m="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
  [[ -x "$m" ]] && { echo "$m"; return; }
  echo ""
}
CHROME="$(find_chrome)"
[[ -z "$CHROME" ]] && { echo "ERROR: Chrome not found." >&2; exit 1; }

# ── build chromerpc + automate from GitHub (reproducible) ────────────────────
# Cloned + cached in /tmp/chromerpc-testing; rebuild with REBUILD_CHROMERPC=1.
CACHE="/tmp/chromerpc-testing"
CHROMERPC_BIN="$CACHE/bin/chromerpc"
AUTOMATE_BIN="$CACHE/bin/automate"
if [[ ! -x "$CHROMERPC_BIN" || ! -x "$AUTOMATE_BIN" || -n "${REBUILD_CHROMERPC:-}" ]]; then
  echo "Fetching + building chromerpc from $CHROMERPC_GIT ..."
  mkdir -p "$CACHE/bin"
  if [[ -d "$CACHE/src/.git" ]]; then
    git -C "$CACHE/src" pull --quiet || true
  else
    git clone --quiet "$CHROMERPC_GIT" "$CACHE/src"
  fi
  ( cd "$CACHE/src" \
      && go build -o "$CHROMERPC_BIN" ./cmd/chromerpc \
      && go build -o "$AUTOMATE_BIN"  ./cmd/automate )
  echo "chromerpc ready: $CHROMERPC_BIN"
else
  echo "Using cached chromerpc: $CHROMERPC_BIN"
fi

find_free_port() { python3 -c "import socket; s=socket.socket(); s.bind(('',0)); p=s.getsockname()[1]; s.close(); print(p)"; }

is_url=false; [[ "$INPUT" == http://* || "$INPUT" == https://* ]] && is_url=true
is_dir=false; [[ -d "$INPUT" ]] && is_dir=true

# ── local HTTP server for local inputs ───────────────────────────────────────
if [[ "$is_url" == false ]]; then
  if [[ "$is_dir" == true ]]; then SERVE_DIR="$(cd "$INPUT" && pwd)"; else SERVE_DIR="$(cd "$(dirname "$INPUT")" && pwd)"; fi
  HTTP_PORT="$(find_free_port)"
  python3 -m http.server "$HTTP_PORT" --directory "$SERVE_DIR" &>/dev/null &
  HTTP_PID=$!; disown "$HTTP_PID"
fi

# ── start chromerpc ──────────────────────────────────────────────────────────
CHROMERPC_PORT="$(find_free_port)"
CHROMERPC_ADDR="localhost:$CHROMERPC_PORT"
"$CHROMERPC_BIN" --headless --addr ":$CHROMERPC_PORT" &>/dev/null &
CHROMERPC_PID=$!; disown "$CHROMERPC_PID"
for i in $(seq 1 40); do
  bash -c ">/dev/tcp/localhost/$CHROMERPC_PORT" 2>/dev/null && break
  sleep 0.5
  [[ $i -eq 40 ]] && { echo "ERROR: chromerpc not ready" >&2; exit 1; }
done

WORK_DIR="$(mktemp -d)"

take_screenshot() {
  local url="$1" out_file="$2"
  mkdir -p "$(dirname "$out_file")"
  local abs_out; abs_out="$(cd "$(dirname "$out_file")" && pwd)/$(basename "$out_file")"
  local slug; slug="$(basename "$out_file" .png | tr ' /' '_-')"
  local textproto="$WORK_DIR/${slug}.textproto"

  if [[ -n "${AUTOMATION:-}" && -f "${AUTOMATION:-}" ]]; then
    sed -e "s#{{URL}}#${url}#g" -e "s#{{OUT}}#${abs_out}#g" "$AUTOMATION" > "$textproto"
  else
    cat > "$textproto" <<PROTO
name: "shot_${slug}"
steps: { label: "viewport" set_viewport: { width: 1280 height: 800 device_scale_factor: 2 } }
steps: { label: "navigate" navigate: { url: "$url" } }
steps: { label: "settle"   wait: { milliseconds: ${SNAP_WAIT_MS:-600} } }
steps: { label: "capture"  screenshot: { output_path: "$abs_out" format: "png" } }
PROTO
  fi
  echo "snap: $url -> $abs_out"
  "$AUTOMATE_BIN" -addr "$CHROMERPC_ADDR" -input "$textproto"
}

if [[ "$is_url" == true ]]; then
  take_screenshot "$INPUT" "$OUTPUT"
elif [[ "$is_dir" == true ]]; then
  mkdir -p "$OUTPUT"
  shopt -s nullglob
  html_files=("$INPUT"/*.html)
  [[ ${#html_files[@]} -eq 0 ]] && { echo "ERROR: no .html in $INPUT" >&2; exit 1; }
  for html in "${html_files[@]}"; do
    fn="$(basename "$html")"; slug="${fn%.html}"
    take_screenshot "http://localhost:$HTTP_PORT/$fn" "$OUTPUT/${slug}.png"
  done
else
  fn="$(basename "$INPUT")"
  take_screenshot "http://localhost:$HTTP_PORT/$fn" "$OUTPUT"
fi
echo "Done."
