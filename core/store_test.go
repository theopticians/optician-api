package core

import (
	"image"
	"os"
	"testing"

	stores "github.com/theopticians/optician-api/core/store"
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

func GenericTestStore(t *testing.T, newStore func() stores.Store) {

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

		_, diffPixels, err := compareImagesBin(testImg1, img1Retrieved, []image.Rectangle{}, 0)
		if err != nil {
			t.Fatal("Error comparing images:", err)
		}

		if diffPixels > 0 {
			t.Fatal("Retrieved image is not equal to original")
		}
	})

	t.Run("base image id", func(t *testing.T) {
		const baseImageID = "abc"
		const project = "project"
		const branch = "branch"
		const target = "target"
		const browser = "browser"

		s := newStore()

		_, err := s.GetBaseImageID(project, branch, target, browser)
		if err == nil {
			t.Fatal("Expected error when getting unexistant base image ID")
		}

		err = s.SetBaseImageID(baseImageID, project, branch, target, browser)
		if err != nil {
			t.Fatal("Error setting base image:", err)
		}

		retrieved, err := s.GetBaseImageID(project, branch, target, browser)
		if err != nil {
			t.Fatal("Error getting base image:", err)
		}

		if retrieved != baseImageID {
			t.Fatal("Expected retrieved base image to be ", baseImageID, " got ", retrieved)
		}

	})

}
