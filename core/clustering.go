package core

import (
	"errors"
	"image"
	"image/color"
)

var NoPixelFoundErr = errors.New("No pixel found")

func clusterDiffImage(img image.Image) []image.Rectangle {
	mask := image.NewAlpha(img.Bounds())

	var err error
	clusters := []image.Rectangle{}

	var pix image.Point
	pix, err = findUnmaskedPixel(img, mask)

	for err == nil {
		pixels := findAdjacentPixels(img, pix, mask)
		if len(pixels) == 0 {
			panic("No pixels returned by adjacent pixels")
		}
		cluster := pointsBounds(pixels)
		clusters = append(clusters, cluster)
		pix, err = findUnmaskedPixel(img, mask)
	}

	return mergeOverlappingClusters(clusters)
}

// If needed, makes a rect bigger to fit the point
func growRect(rect *image.Rectangle, point image.Point) {
	if !point.In(*rect) {

		if rect.Min.X > point.X {
			rect.Min.X = point.X
		}

		if rect.Max.X < point.X {
			rect.Max.X = point.X
		}

		if rect.Min.Y > point.Y {
			rect.Min.Y = point.Y
		}

		if rect.Max.Y < point.Y {
			rect.Max.Y = point.Y
		}

	}
}

func mergeClustersByCondition(c []image.Rectangle, condition func(image.Rectangle, image.Rectangle) bool) []image.Rectangle {
	clusters := c
	for i := 0; i < len(clusters); i++ {
		for j := i + 1; j < len(clusters); j++ {
			if condition(clusters[i], clusters[j]) {
				growRect(&clusters[j], clusters[i].Min)
				growRect(&clusters[j], clusters[i].Max)
				clusters = append(clusters[:i], clusters[i+1:]...)

				return mergeClustersByCondition(clusters, condition)
			}
		}
	}

	return clusters
}

func mergeOverlappingClusters(c []image.Rectangle) []image.Rectangle {
	return mergeClustersByCondition(c, func(r1, r2 image.Rectangle) bool {
		return r1.Overlaps(r2)
	})
}

// Finds the smallest rect that contains all points
func pointsBounds(points []image.Point) image.Rectangle {
	rect := image.Rectangle{points[0], points[0]}
	for i := 1; i < len(points); i++ {
		if !points[i].In(rect) {
			growRect(&rect, points[i])
		}
	}

	return rect
}

func getAlpha(color color.Color) int {
	_, _, _, a := color.RGBA()
	return int(a)
}

func findAdjacentPixels(img image.Image, start image.Point, mask *image.Alpha) []image.Point {
	if !start.In(img.Bounds()) || mask.AlphaAt(start.X, start.Y).A != 0 {
		return []image.Point{}
	}

	if getAlpha(img.At(start.X, start.Y)) == 0 {
		return []image.Point{}
	}

	mask.SetAlpha(start.X, start.Y, color.Alpha{255})

	pixels := []image.Point{start}

	pixels = append(pixels, findAdjacentPixels(img, image.Point{start.X - 1, start.Y}, mask)...)
	pixels = append(pixels, findAdjacentPixels(img, image.Point{start.X + 1, start.Y}, mask)...)
	pixels = append(pixels, findAdjacentPixels(img, image.Point{start.X, start.Y + 1}, mask)...)
	pixels = append(pixels, findAdjacentPixels(img, image.Point{start.X, start.Y - 1}, mask)...)

	return pixels
}

func findUnmaskedPixel(img image.Image, mask *image.Alpha) (image.Point, error) {
	if !img.Bounds().Eq(mask.Bounds()) {
		return image.Point{}, errors.New("Image and mask have different bounds")
	}

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			_, _, _, ia := img.At(x, y).RGBA()
			_, _, _, ma := mask.At(x, y).RGBA()
			if ia > 0 && ma == 0 {
				return image.Point{x, y}, nil
			}
		}
	}

	return image.Point{}, NoPixelFoundErr
}
