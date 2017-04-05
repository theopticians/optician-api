package core

import (
	"image"
	"testing"
)

func TestBinDiff(t *testing.T) {
	_, diffPixels, err := compareImagesBin(testImg1, testImg2, []image.Rectangle{})

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 33454 {
		t.Fatal("Expected number of pixel differences between testImg1 and testImg2 to be 33454, got ", diffPixels)
	}

	_, diffPixels, err = compareImagesBin(testImg1, testImg1, []image.Rectangle{})

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 0 {
		t.Fatal("Expected number of pixel differences between equal images to be 0, got ", diffPixels)
	}
}

func TestBinDiffMaskInvalid(t *testing.T) {
	_, _, err := compareImagesBin(testImg1, testImg2, []image.Rectangle{testMask1, testMaskInvalid})

	if err == nil {
		t.Fatal("Expected compareImagesBin to return error when passed invalid mask")
	}

}

func TestBinDiffMask(t *testing.T) {
	_, diffPixels, err := compareImagesBin(testImg1, testImg2, []image.Rectangle{testMask1})

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 33351 {
		t.Fatal("Expected number of pixel differences between testImg1 and testImg2 with testMask1 to be 33351, got ", diffPixels)
	}

	_, diffPixels, err = compareImagesBin(testImg1, testImg1, []image.Rectangle{testMask1})

	if err != nil {
		t.Fatal("Error comparing images:", err)
	}

	if diffPixels != 0 {
		t.Fatal("Expected number of pixel differences between equal images to be 0, got ", diffPixels)
	}
}
