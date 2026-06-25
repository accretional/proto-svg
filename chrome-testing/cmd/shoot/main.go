// shoot reads chrome-testing/gallery/catalogue.json and emits chromerpc
// AutomationSequence textproto chunks that drive the SVG Lab gallery SPA and
// screenshot the VIEWER for every preset (one preset = one attribute value):
//
//	screenshots/gallery/<tag>/NN-<slug>.png          — one shot per STATIC preset
//	screenshots/gallery/<tag>/NN-<slug>/frame-NN.png — a frame sequence per preset
//	                                                    of a TEMPORAL (animation) element
//
// The gallery is a single SPA. Each chunk navigates to it once, then sets
// location.hash = '#/embed/<tag>/<idx>' per preset (which renders the chrome-free
// embed viewer with that preset applied) and screenshots. Output is chunked so
// each RunAutomation response stays small.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type catPreset struct {
	Name        string `json:"name"`
	Interactive string `json:"interactive,omitempty"` // "hover" | "click"
}
type catElement struct {
	Tag      string      `json:"tag"`
	Temporal bool        `json:"temporal"`
	Presets  []catPreset `json:"presets"`
}
type catalogue struct {
	Elements []catElement `json:"elements"`
}

func main() {
	catPath := flag.String("catalogue", "chrome-testing/gallery/catalogue.json", "catalogue.json path")
	base := flag.String("base", "", "served gallery base URL (e.g. http://localhost:PORT/gallery)")
	outdir := flag.String("outdir", "", "screenshots output dir (ROOT-relative, e.g. chrome-testing/screenshots/gallery)")
	seq := flag.String("seq", "", "output textproto path prefix")
	only := flag.String("only", "", "comma-separated tags to limit to")
	perChunk := flag.Int("chunk", 20, "screenshots per chunk")
	navWait := flag.Int("navwait", 320, "ms after a hash change before screenshot")
	frameWait := flag.Int("framewait", 200, "ms between frames of a temporal capture")
	frames := flag.Int("frames", 6, "temporal frame count")
	flag.Parse()
	if *base == "" || *outdir == "" || *seq == "" {
		fmt.Fprintln(os.Stderr, "usage: shoot -catalogue FILE -base URL -outdir DIR -seq PREFIX")
		os.Exit(2)
	}
	resume := os.Getenv("RESUME") != ""

	raw, err := os.ReadFile(*catPath)
	if err != nil {
		panic(err)
	}
	var cat catalogue
	if err := json.Unmarshal(raw, &cat); err != nil {
		panic(err)
	}

	limit := map[string]bool{}
	for _, n := range strings.Split(*only, ",") {
		if n = strings.TrimSpace(n); n != "" {
			limit[n] = true
		}
	}

	baseURL := strings.TrimRight(*base, "/")
	prefix := strings.TrimSuffix(*seq, filepath.Ext(*seq))
	if err := os.MkdirAll(filepath.Dir(prefix), 0o755); err != nil {
		panic(err)
	}
	w := &chunkWriter{prefix: prefix, perChunk: *perChunk, baseURL: baseURL, navWait: *navWait}
	w.start()

	shots, temporal := 0, 0
	for _, el := range cat.Elements {
		if len(limit) > 0 && !limit[el.Tag] {
			continue
		}
		dir := filepath.Join(*outdir, el.Tag)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
		for i, p := range el.Presets {
			slug := fmt.Sprintf("%02d-%s", i, slugify(p.Name))
			hash := fmt.Sprintf("location.hash='#/embed/%s/%d';void 0", el.Tag, i)
			// Interactive presets: drive a real Hover/Click and capture before/after.
			if p.Interactive != "" {
				framedir := filepath.Join(dir, slug)
				if resume {
					if _, e := os.Stat(framedir + ".gif"); e == nil {
						continue
					}
				}
				if err := os.MkdirAll(framedir, 0o755); err != nil {
					panic(err)
				}
				w.interactive(hash, framedir, p.Interactive)
				shots += 2
				temporal++
				continue
			}
			if !el.Temporal {
				out := filepath.Join(dir, slug+".png")
				if resume {
					if _, e := os.Stat(out); e == nil {
						continue
					}
				}
				w.static(hash, out)
				shots++
				continue
			}
			framedir := filepath.Join(dir, slug)
			if resume {
				if _, e := os.Stat(framedir + ".gif"); e == nil {
					continue
				}
			}
			if err := os.MkdirAll(framedir, 0o755); err != nil {
				panic(err)
			}
			w.temporal(hash, framedir, *frames, *frameWait)
			shots += *frames
			temporal++
		}
	}
	w.flush()
	fmt.Printf("Wrote %d chunks (%s-NNN.textproto) — %d screenshots (%d temporal frame-sets)\n",
		w.nChunks, prefix, shots, temporal)
}

// chunkWriter emits steps into capped textproto chunks. Each chunk sets the
// viewport, NAVIGATES to the gallery once (loading the SPA), then drives it by
// setting the embed hash per preset and screenshotting the viewer.
type chunkWriter struct {
	prefix         string
	perChunk       int
	baseURL        string
	navWait        int
	b              strings.Builder
	shots, nChunks int
}

func (w *chunkWriter) start() {
	w.b.Reset()
	fmt.Fprintf(&w.b, "name: \"svg-shots-%03d\"\n", w.nChunks)
	// embed viewer is 440px centred; a 560x560 viewport frames it with padding.
	fmt.Fprintf(&w.b, "steps { set_viewport { width: 560 height: 560 device_scale_factor: 2 } }\n")
	fmt.Fprintf(&w.b, "steps { navigate { url: %q } }\n", w.baseURL+"/index.html")
	fmt.Fprintf(&w.b, "steps { wait { milliseconds: 450 } }\n")
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

// static sets the embed hash, settles, and captures one viewport screenshot.
func (w *chunkWriter) static(hash, out string) {
	if w.shots >= w.perChunk {
		w.rollover()
	}
	w.line(fmt.Sprintf("steps { evaluate_script { expression: %q } }", hash))
	w.line(fmt.Sprintf("steps { wait { milliseconds: %d } }", w.navWait))
	w.line(fmt.Sprintf("steps { screenshot { output_path: %q format: \"png\" } }", out))
	w.shots++
}

// temporal sets the embed hash for an animated element and captures `frames`
// screenshots. It STEPS the SVG document clock deterministically with
// setCurrentTime() rather than relying on wall-clock waits — headless Chrome does
// not reliably advance the compositor between rapid screenshots, so a plain
// wait+screenshot loop captures the same frame repeatedly. Stepping the clock
// guarantees distinct frames spanning the (≈2s) animation. The whole sequence
// stays in one chunk so the SVG is not re-rendered between frames.
func (w *chunkWriter) temporal(hash, framedir string, frames, frameWait int) {
	if w.shots > 0 && w.shots+frames > w.perChunk {
		w.rollover()
	}
	w.line(fmt.Sprintf("steps { evaluate_script { expression: %q } }", hash))
	w.line("steps { wait { milliseconds: 350 } }") // let the embed SVG mount
	step := 0.33                                    // seconds of SVG time per frame (~1.65s span over a 2s dur)
	for i := 0; i < frames; i++ {
		t := float64(i) * step
		set := fmt.Sprintf("var s=document.querySelector('#viewer svg');"+
			"if(s&&s.setCurrentTime){if(s.pauseAnimations)s.pauseAnimations();s.setCurrentTime(%.2f);}void 0", t)
		w.line(fmt.Sprintf("steps { evaluate_script { expression: %q } }", set))
		w.line(fmt.Sprintf("steps { wait { milliseconds: %d } }", frameWait))
		w.line(fmt.Sprintf("steps { screenshot { output_path: %q format: \"png\" } }",
			filepath.Join(framedir, fmt.Sprintf("frame-%02d.png", i))))
		w.shots++
	}
}

// interactive captures an event-attribute preset under REAL user input: frame-00
// is the resting state, then a Hover or Click fires the inline handler (which
// fades / recolors the element), and frame-01 is the reacted state. gifenc turns
// the two frames into a looping before↔after GIF (+ strobe). The showcased element
// carries data-lab, so "#viewer [data-lab]" targets exactly it.
func (w *chunkWriter) interactive(hash, framedir, kind string) {
	if w.shots > 0 && w.shots+2 > w.perChunk {
		w.rollover()
	}
	w.line(fmt.Sprintf("steps { evaluate_script { expression: %q } }", hash))
	w.line("steps { wait { milliseconds: 350 } }")
	w.line(fmt.Sprintf("steps { screenshot { output_path: %q format: \"png\" } }",
		filepath.Join(framedir, "frame-00.png")))
	const sel = "#viewer [data-lab]"
	if kind == "click" {
		w.line(fmt.Sprintf("steps { click { selector: %q } }", sel))
	} else {
		w.line(fmt.Sprintf("steps { hover { selector: %q } }", sel))
	}
	w.line("steps { wait { milliseconds: 220 } }")
	w.line(fmt.Sprintf("steps { screenshot { output_path: %q format: \"png\" } }",
		filepath.Join(framedir, "frame-01.png")))
	w.shots += 2
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

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = slugRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 40 {
		s = strings.Trim(s[:40], "-")
	}
	if s == "" {
		s = "v"
	}
	return s
}
