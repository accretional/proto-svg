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
