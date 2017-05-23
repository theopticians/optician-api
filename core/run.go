package core

import (
	"image"

	"github.com/pkg/errors"
	"github.com/theopticians/optician-api/core/imgdiff"
	"github.com/theopticians/optician-api/core/structs"
)

func RunTest(r *structs.Result) error {

	baseImg, err := db.GetImage(r.BaseImageID)
	if err != nil {
		return errors.Wrap(err, "error getting base image")
	}

	testImg, err := db.GetImage(r.ImageID)
	if err != nil {
		return errors.Wrap(err, "error getting test image")
	}

	var mask []image.Rectangle
	if r.MaskID == "nomask" {
		mask = []image.Rectangle{}
	} else {
		mask, err = db.GetMask(r.MaskID)
		if err != nil {
			return errors.Wrap(err, "error getting mask")
		}
	}

	diffImg, diffScore := imgdiff.ComputeDiffImage(baseImg, testImg, mask)

	diffImageID, err := db.StoreImage(diffImg)

	if err != nil {
		return errors.Wrap(err, "error getting storing diff image")
	}

	r.DiffClusters = imgdiff.PerformClustering(diffImg)
	r.DiffImageID = diffImageID
	r.DiffScore = diffScore

	return nil
}
