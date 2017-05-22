package core

import (
	"image"
	"testing"
)

func BenchmarkGerardClustering(b *testing.B) {
	img, _, err := compareImagesBin(testImg1, testImg2, []image.Rectangle{}, 0)

	if err != nil {
		b.Fatal("Error comparing images:", err)
	}

	benchmarkclusterer(b, naiveClusterer, img)
}

func benchmarkclusterer(b *testing.B, c clusterer, img image.Image) {
	for n := 0; n < b.N; n++ {
		c(img)
	}
}
