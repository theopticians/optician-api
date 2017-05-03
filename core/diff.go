package core

import (
	"github.com/pkg/errors"
	"image"
	"image/color"
)

func RunTest(r *Result) error {

	baseImg, err := store.GetImage(r.BaseImageID)
	if err != nil {
		return errors.Wrap(err, "error getting base image")
	}

	testImg, err := store.GetImage(r.ImageID)
	if err != nil {
		return errors.Wrap(err, "error getting test image")
	}

	var mask []image.Rectangle
	if r.MaskID == "nomask" {
		mask = []image.Rectangle{}
	} else {
		mask, err = store.GetMask(r.MaskID)
		if err != nil {
			return errors.Wrap(err, "error getting mask")
		}
	}

	diffImg, diffScore := computeDiffImage(baseImg, testImg, mask)

	diffImageID, err := store.StoreImage(diffImg)

	if err != nil {
		return errors.Wrap(err, "error getting storing diff image")
	}

	r.DiffClusters = clusterDiffImage(diffImg)
	r.DiffImageID = diffImageID
	r.DiffScore = diffScore

	return nil
}

func computeDiffImage(img1, img2 image.Image, masks []image.Rectangle) (image.Image, float64) {
	diffImg, n, _ := compareImagesBin(img1, img2, masks)
	return diffImg, float64(n)
}

// CompareImagesBin compares a and b using binary comparison.
func compareImagesBin(a, b image.Image, masks []image.Rectangle) (image.Image, int, error) {
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
			if d > 0 && !pixelInMask(x, y, masks) {
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

func diffColor(c1, c2 color.Color) int64 {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	var diff int64
	diff += abs(int64(r1) - int64(r2))
	diff += abs(int64(g1) - int64(g2))
	diff += abs(int64(b1) - int64(b2))
	diff += abs(int64(a1) - int64(a2))
	return diff
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
