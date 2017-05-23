package imgdiff

import (
	"image"
	"image/color"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"
)

func ComputeDiffImage(img1, img2 image.Image, masks []image.Rectangle) (image.Image, float64) {
	diffImg, n, _ := compareImagesBin(img1, img2, masks, 0.05)
	return diffImg, float64(n)
}

// CompareImagesBin compares a and b using binary comparison.
func compareImagesBin(a, b image.Image, masks []image.Rectangle, threshold float64) (image.Image, int, error) {
	ab, bb := a.Bounds(), b.Bounds()
	w, h := ab.Dx(), ab.Dy()
	if w != bb.Dx() || h != bb.Dy() {
		return nil, -1, errors.New("Different image sizes")
	}

	if err := checkMasks(masks); err != nil {
		return nil, -1, err
	}

	diff := image.NewNRGBA(image.Rect(0, 0, w, h))
	n := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d := diffColor(a.At(ab.Min.X+x, ab.Min.Y+y), b.At(bb.Min.X+x, bb.Min.Y+y))
			c := color.RGBA{0, 0, 0, 0}
			if d > threshold && !pixelInMask(x, y, masks) {
				c.R = 0xff
				c.A = 0xff
				//c.A = uint8(100 + d*0xff/0xffff)
				n++
			}
			diff.Set(x, y, c)
		}
	}
	return diff, n, nil
}

func goColorful(c color.Color) colorful.Color {
	r, g, b, _ := c.RGBA()
	return colorful.Color{R: float64(r) / float64(0xffff), G: float64(g) / float64(0xffff), B: float64(b) / float64(0xffff)}
}

func diffColor(c1, c2 color.Color) float64 {
	co1 := goColorful(c1)
	co2 := goColorful(c2)
	return co1.DistanceCIE76(co2)
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// pixelInMask checks if a pixel is inside a series of masks
func pixelInMask(x, y int, masks []image.Rectangle) bool {
	for _, m := range masks {
		if x >= m.Min.X && x <= m.Max.X {
			if y >= m.Min.Y && y <= m.Max.Y {
				return true
			}
		}
	}

	return false
}

func checkMasks(masks []image.Rectangle) error {
	for _, m := range masks {
		if m.Max.X < m.Min.X || m.Max.Y < m.Min.Y {
			return errors.New("Mask is invalid")
		}
	}

	return nil
}
