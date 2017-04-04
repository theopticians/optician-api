package core

import (
	"image"
	"os"
	"testing"
)

const (
	img1Path = "./testimages/so_1.png"
	img2Path = "./testimages/so_2.png"
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
	img1 := readImage(img1Path)
	//img2 := readImage(img2Path)

	t.Run("image storage", func(t *testing.T) {

		s := newStore()

		img1ID, err := s.StoreImage(img1)
		if err != nil {
			t.Fatal("Error storing image:", err)
		}

		img1Retrieved, err := s.GetImage(img1ID)
		if err != nil {
			t.Fatal("Error retrieving image:", err)
		}

		_, n, err := compareImagesBin(img1, img1Retrieved, []image.Rectangle{})
		if err != nil {
			t.Fatal("Error comparing images:", err)
		}

		if n > 0 {
			t.Fatal("Retrieved image is not equal to original")
		}

	})
}
