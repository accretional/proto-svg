// gifenc walks the screenshots tree and turns every temporal frame folder
// (a directory containing frame-00.png, frame-01.png, …) into a single animated
// GIF sibling (<folder>.gif), downscaled 2× to keep the file small. Pure stdlib
// (image/png, image/gif) — no external tools.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func main() {
	root := flag.String("dir", "chrome-testing/screenshots/specimens", "screenshots root")
	delay := flag.Int("delay", 28, "GIF frame delay in 1/100s")
	flag.Parse()

	var folders []string
	filepath.WalkDir(*root, func(p string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		if _, e := os.Stat(filepath.Join(p, "frame-00.png")); e == nil {
			folders = append(folders, p)
		}
		return nil
	})
	fmt.Printf("Encoding %d frame folders -> GIF\n", len(folders))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	var done, failed int
	var mu sync.Mutex
	for _, f := range folders {
		wg.Add(1)
		sem <- struct{}{}
		go func(folder string) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := encodeFolder(folder, *delay); err != nil {
				mu.Lock()
				failed++
				mu.Unlock()
				return
			}
			// Best-effort strobe still (onion-skin of the frames); never fatal.
			_ = strobeFolder(folder)
			mu.Lock()
			done++
			mu.Unlock()
		}(f)
	}
	wg.Wait()
	fmt.Printf("Wrote %d GIFs · %d failed\n", done, failed)
}

func encodeFolder(folder string, delay int) error {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	var frames []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "frame-") && strings.HasSuffix(e.Name(), ".png") {
			frames = append(frames, filepath.Join(folder, e.Name()))
		}
	}
	if len(frames) == 0 {
		return fmt.Errorf("no frames")
	}
	sort.Strings(frames)

	g := &gif.GIF{LoopCount: 0}
	for _, fp := range frames {
		img, err := loadPNG(fp)
		if err != nil {
			return err
		}
		small := downscale2x(img)
		pal := image.NewPaletted(small.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(pal, small.Bounds(), small, image.Point{})
		g.Image = append(g.Image, pal)
		g.Delay = append(g.Delay, delay)
		g.Disposal = append(g.Disposal, gif.DisposalNone)
	}
	out := folder + ".gif"
	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()
	return gif.EncodeAll(w, g)
}

// strobeFolder composites the frame sequence into a single onion-skin "strobe"
// still (<folder>.strobe.png): each frame is alpha-ramped (early → faint, late →
// solid) and LIGHTENED onto the accumulator, so the moving shape leaves a trail
// of stamps whose SPACING reveals the velocity profile at a glance — clustered
// for calcMode=discrete, evenly spaced for linear, bunched at the ends for
// spline. (Lighten suits the dark gallery canvas where shapes are brighter than
// the background.)
func strobeFolder(folder string) error {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	var frames []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "frame-") && strings.HasSuffix(e.Name(), ".png") {
			frames = append(frames, filepath.Join(folder, e.Name()))
		}
	}
	if len(frames) < 2 {
		return fmt.Errorf("need ≥2 frames")
	}
	sort.Strings(frames)

	imgs := make([]image.Image, 0, len(frames))
	for _, fp := range frames {
		img, err := loadPNG(fp)
		if err != nil {
			return err
		}
		imgs = append(imgs, img)
	}
	b := imgs[0].Bounds()
	// Interactive captures are exactly two frames (rest, reacted): a side-by-side
	// reads far clearer than an onion-skin, so emit that.
	if len(imgs) == 2 {
		w, h := b.Dx(), b.Dy()
		sxs := image.NewRGBA(image.Rect(0, 0, w*2, h))
		draw.Draw(sxs, image.Rect(0, 0, w, h), imgs[0], b.Min, draw.Src)
		draw.Draw(sxs, image.Rect(w, 0, w*2, h), imgs[1], b.Min, draw.Src)
		out := folder + ".strobe.png"
		f, err := os.Create(out)
		if err != nil {
			return err
		}
		defer f.Close()
		return png.Encode(f, downscale2x(sxs))
	}
	// Background = the top-left corner pixel of the first frame.
	bgr, bgg, bgb, _ := imgs[0].At(b.Min.X, b.Min.Y).RGBA()
	acc := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			acc.Set(x, y, color.RGBA{uint8(bgr >> 8), uint8(bgg >> 8), uint8(bgb >> 8), 255})
		}
	}
	n := len(imgs)
	for i, img := range imgs {
		alpha := 0.30 + 0.70*float64(i)/float64(n-1)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				fr, fg, fb, _ := img.At(x, y).RGBA()
				// blend this frame toward the background by alpha
				br := float64(bgr>>8) + (float64(fr>>8)-float64(bgr>>8))*alpha
				bgn := float64(bgg>>8) + (float64(fg>>8)-float64(bgg>>8))*alpha
				bb := float64(bgb>>8) + (float64(fb>>8)-float64(bgb>>8))*alpha
				cur := acc.RGBAAt(x, y)
				acc.SetRGBA(x, y, color.RGBA{
					maxu8(cur.R, uint8(br)),
					maxu8(cur.G, uint8(bgn)),
					maxu8(cur.B, uint8(bb)),
					255,
				})
			}
		}
	}
	out := folder + ".strobe.png"
	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()
	return png.Encode(w, downscale2x(acc))
}

func maxu8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func loadPNG(p string) (image.Image, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

// downscale2x averages each 2×2 block into one pixel.
func downscale2x(src image.Image) *image.RGBA {
	b := src.Bounds()
	w, h := b.Dx()/2, b.Dy()/2
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx, sy := b.Min.X+x*2, b.Min.Y+y*2
			var r, g, bl, a uint32
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					cr, cg, cb, ca := src.At(sx+dx, sy+dy).RGBA()
					r += cr
					g += cg
					bl += cb
					a += ca
				}
			}
			dst.Set(x, y, color.RGBA{uint8((r / 4) >> 8), uint8((g / 4) >> 8), uint8((bl / 4) >> 8), uint8((a / 4) >> 8)})
		}
	}
	return dst
}
