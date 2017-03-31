package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
)

func imageToBase64(img image.Image) string {
	var buf []byte
	b := bytes.NewBuffer(buf)

	png.Encode(b, img)

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func base64ToImage(b64 string) image.Image {
	pngbytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer(pngbytes)

	img, _, err := image.Decode(b)
	if err != nil {
		panic(err)
	}

	return img
}
