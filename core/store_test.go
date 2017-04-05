package core

import (
	"image"
	"os"
	"testing"
)

var (
	testImg1        = readImage("./testimages/so_1.png")
	testImg2        = readImage("./testimages/so_2.png")
	testMask1       = image.Rectangle{image.Point{0, 0}, image.Point{300, 300}}
	testMaskInvalid = image.Rectangle{image.Point{10, 10}, image.Point{2, 30}}
)

func readImage(path string) image.Image {
	reader, _ := os.Open(path)
	defer reader.Close()
	im, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	return im
}

func testStore(t *testing.T, newStore func() Store) {

	t.Run("image storage", func(t *testing.T) {
		s := newStore()

		img1ID, err := s.StoreImage(testImg1)
		if err != nil {
			t.Fatal("Error storing image:", err)
		}

		img1Retrieved, err := s.GetImage(img1ID)
		if err != nil {
			t.Fatal("Error retrieving image:", err)
		}

		_, diffPixels, err := compareImagesBin(testImg1, img1Retrieved, []image.Rectangle{})
		if err != nil {
			t.Fatal("Error comparing images:", err)
		}

		if diffPixels > 0 {
			t.Fatal("Retrieved image is not equal to original")
		}
	})

}
