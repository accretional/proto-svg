// shoot reads chrome-testing/generated/specimens.json and emits chromerpc
// AutomationSequence textproto chunks that screenshot every per-value specimen:
//
//	screenshots/specimens/<tag>/NN-<slug>.png          — one shot per STATIC specimen
//	screenshots/specimens/<tag>/NN-<slug>/frame-NN.png — a frame sequence per TEMPORAL specimen
//
// Each specimen is a standalone HTML page served by shoot.sh at
// <BASE>/html/<file>; we just navigate to it, settle, and screenshot. Temporal
// specimens host a live SMIL animation, so we capture several frames spaced
// -framewait apart while the animation plays. Output is chunked so each
// RunAutomation response (which carries the screenshot bytes) stays small.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// specimen mirrors one entry of specimens.json:
//
//	{ "i":0, "label":"x=\"10\"", "value":"10", "file":"specimen/rect/00-10.html", "temporal":false }
type specimen struct {
	I        int    `json:"i"`
	Label    string `json:"label"`
	Value    string `json:"value"`
	File     string `json:"file"`
	Temporal bool   `json:"temporal"`
}

func main() {
	specimens := flag.String("specimens", "chrome-testing/generated/specimens.json", "specimens.json path")
	base := flag.String("base", "", "served gallery base URL (e.g. http://localhost:PORT)")
	outdir := flag.String("outdir", "", "screenshots output dir (ROOT-relative, e.g. chrome-testing/screenshots/specimens)")
	seq := flag.String("seq", "", "output textproto path prefix")
	only := flag.String("only", "", "comma-separated tags to limit to")
	perChunk := flag.Int("chunk", 20, "screenshots per chunk")
	frameWait := flag.Int("framewait", 200, "ms between frames of a temporal capture")
	frames := flag.Int("frames", 6, "temporal frame count")
	flag.Parse()
	if *base == "" || *outdir == "" || *seq == "" {
		fmt.Fprintln(os.Stderr, "usage: shoot -specimens FILE -base URL -outdir DIR -seq PREFIX")
		os.Exit(2)
	}
	// RESUME=1 (set by shoot.sh) skips specimens whose output already exists.
	resume := os.Getenv("RESUME") != ""

	raw, err := os.ReadFile(*specimens)
	if err != nil {
		panic(err)
	}
	byTag := map[string][]specimen{}
	if err := json.Unmarshal(raw, &byTag); err != nil {
		panic(err)
	}

	limit := map[string]bool{}
	for _, n := range strings.Split(*only, ",") {
		if n = strings.TrimSpace(n); n != "" {
			limit[n] = true
		}
	}

	tags := make([]string, 0, len(byTag))
	for t := range byTag {
		if len(limit) > 0 && !limit[t] {
			continue
		}
		tags = append(tags, t)
	}
	sort.Strings(tags)

	base = trimSlash(base)
	prefix := strings.TrimSuffix(*seq, filepath.Ext(*seq))
	if err := os.MkdirAll(filepath.Dir(prefix), 0o755); err != nil {
		panic(err)
	}
	w := &chunkWriter{prefix: prefix, perChunk: *perChunk}
	w.start()

	shots, temporal := 0, 0
	for _, tag := range tags {
		items := byTag[tag]
		if len(items) == 0 {
			continue
		}
		dir := filepath.Join(*outdir, tag)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
		for _, sp := range items {
			slug := strings.TrimSuffix(filepath.Base(sp.File), ".html") // "NN-<slug>"
			url := *base + "/html/" + sp.File
			if !sp.Temporal {
				out := filepath.Join(dir, slug+".png")
				if resume {
					if _, e := os.Stat(out); e == nil {
						continue
					}
				}
				w.static(url, out)
				shots++
				continue
			}
			// temporal: capture a live SMIL frame sequence into a per-value folder.
			framedir := filepath.Join(dir, slug)
			if resume {
				if _, e := os.Stat(framedir + ".gif"); e == nil {
					continue
				}
			}
			if err := os.MkdirAll(framedir, 0o755); err != nil {
				panic(err)
			}
			w.temporal(url, framedir, *frames, *frameWait)
			shots += *frames
			temporal++
		}
	}
	w.flush()
	fmt.Printf("Wrote %d chunks (%s-NNN.textproto) — %d screenshots (%d temporal frame-sets)\n",
		w.nChunks, prefix, shots, temporal)
}

func trimSlash(s *string) *string {
	t := strings.TrimRight(*s, "/")
	return &t
}

// chunkWriter emits steps into capped textproto chunks. Each chunk sets the
// viewport once, then navigates fresh to each specimen page and screenshots.
type chunkWriter struct {
	prefix         string
	perChunk       int
	b              strings.Builder
	shots, nChunks int
}

func (w *chunkWriter) start() {
	w.b.Reset()
	fmt.Fprintf(&w.b, "name: \"svg-shots-%03d\"\n", w.nChunks)
	fmt.Fprintf(&w.b, "steps { set_viewport { width: 480 height: 480 device_scale_factor: 1 } }\n")
}
func (w *chunkWriter) line(s string) { w.b.WriteString(s); w.b.WriteByte('\n') }
func (w *chunkWriter) rollover() {
	path := fmt.Sprintf("%s-%03d.textproto", w.prefix, w.nChunks)
	if err := os.WriteFile(path, []byte(w.b.String()), 0o644); err != nil {
		panic(err)
	}
	w.nChunks++
	w.shots = 0
	w.start()
}

// static navigates to the specimen page, settles, and captures one screenshot.
func (w *chunkWriter) static(url, out string) {
	if w.shots >= w.perChunk {
		w.rollover()
	}
	w.line(fmt.Sprintf("steps { navigate { url: %q } }", url))
	w.line("steps { wait { milliseconds: 350 } }")
	w.line(fmt.Sprintf("steps { screenshot { output_path: %q format: \"png\" } }", out))
	w.shots++
}

// temporal navigates to a specimen hosting a live SMIL animation and captures
// `frames` screenshots spaced `frameWait` ms apart so the animation progresses
// across frames. The whole sequence is kept in ONE chunk so the running
// animation's state is preserved (no reload mid-sequence).
func (w *chunkWriter) temporal(url, framedir string, frames, frameWait int) {
	// keep the full sequence within one chunk; roll over first if it won't fit.
	if w.shots > 0 && w.shots+frames > w.perChunk {
		w.rollover()
	}
	w.line(fmt.Sprintf("steps { navigate { url: %q } }", url))
	w.line("steps { wait { milliseconds: 250 } }") // short settle before frame 0
	for i := 0; i < frames; i++ {
		w.line(fmt.Sprintf("steps { wait { milliseconds: %d } }", frameWait))
		w.line(fmt.Sprintf("steps { screenshot { output_path: %q format: \"png\" } }",
			filepath.Join(framedir, fmt.Sprintf("frame-%02d.png", i))))
		w.shots++
	}
}

func (w *chunkWriter) flush() {
	if w.b.Len() > 0 {
		path := fmt.Sprintf("%s-%03d.textproto", w.prefix, w.nChunks)
		if err := os.WriteFile(path, []byte(w.b.String()), 0o644); err != nil {
			panic(err)
		}
		w.nChunks++
	}
}
