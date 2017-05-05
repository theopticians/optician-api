package core

import (
	"fmt"
	"image"
	"image/color"
	"testing"
)

func TestDiffImageClustering(t *testing.T) {
	diffImg := image.NewAlpha(image.Rect(0, 0, 4, 4))

	diffImg.SetAlpha(0, 1, color.Alpha{255})
	diffImg.SetAlpha(1, 1, color.Alpha{255})
	diffImg.SetAlpha(2, 1, color.Alpha{255})
	diffImg.SetAlpha(3, 1, color.Alpha{255})

	rect := naiveClusterer(diffImg)

	fmt.Printf("%v \n", rect)
}

func TestPointsBounds(t *testing.T) {
	points := []image.Point{
		image.Point{3, 3},
		image.Point{7, 9},
	}

	fmt.Println(pointsBounds(points))
}
