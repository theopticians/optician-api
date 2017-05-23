package imgdiff

import (
	"image"
	_ "image/png"
	"os"
	"testing"
)

var (
	testImg1        = readImage("../testimages/so_1.png")
	testImg2        = readImage("../testimages/so_2.png")
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

func TestBinDiff(t *testing.T) {
	_, diffPixels, err := compareImagesBin(testImg1, testImg2, []image.Rectangle{}, 0)

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 33454 {
		t.Fatal("Expected number of pixel differences between testImg1 and testImg2 to be 33454, got ", diffPixels)
	}

	_, diffPixels, err = compareImagesBin(testImg1, testImg1, []image.Rectangle{}, 0)

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 0 {
		t.Fatal("Expected number of pixel differences between equal images to be 0, got ", diffPixels)
	}
}

func TestBinDiffMaskInvalid(t *testing.T) {
	_, _, err := compareImagesBin(testImg1, testImg2, []image.Rectangle{testMask1, testMaskInvalid}, 0)

	if err == nil {
		t.Fatal("Expected compareImagesBin to return error when passed invalid mask")
	}

}

func TestBinDiffMask(t *testing.T) {
	_, diffPixels, err := compareImagesBin(testImg1, testImg2, []image.Rectangle{testMask1}, 0)

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 33351 {
		t.Fatal("Expected number of pixel differences between testImg1 and testImg2 with testMask1 to be 33351, got ", diffPixels)
	}

	_, diffPixels, err = compareImagesBin(testImg1, testImg1, []image.Rectangle{testMask1}, 0)

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 0 {
		t.Fatal("Expected number of pixel differences between equal images to be 0, got ", diffPixels)
	}
}
