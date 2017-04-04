package core

import "image"
import "github.com/pkg/errors"

var store Store = NewBoltStore("optician.db")

func TestList() []string {
	return store.GetTestList()
}

func TestImage(image image.Image, projectID, branch, target, browser string) (Results, error) {
	randID := RandStringBytes(14)

	imgID, err := store.StoreImage(image)

	if err != nil {
		return Results{}, errors.Wrap(err, "error getting storing image")
	}

	baseImgID, err := store.GetBaseImageID(projectID, branch, target, browser)
	if err != nil {
		if err == errNotFound {
			// IF no base image found, set this as base image
			baseImgID = imgID
			store.SetBaseImageID(baseImgID, projectID, branch, target, browser)
		} else {
			return Results{}, errors.Wrap(err, "error getting base image ID")
		}
	}

	baseImg, err := store.GetImage(baseImgID)
	if err != nil {
		return Results{}, errors.Wrap(err, "error getting base image")
	}

	masks, err := store.GetMasks(projectID, branch, target, browser)
	if err != nil {
		if err == errNotFound {
			return Results{}, err
		}
	}

	diffImg, diffScore := computeDiffImage(baseImg, image, masks)

	diffImgID, err := store.StoreImage(diffImg)

	if err != nil {
		return Results{}, errors.Wrap(err, "error getting storing diff image")
	}

	results := Results{
		TestID:      randID,
		ProjectID:   projectID,
		Branch:      branch,
		Target:      target,
		Browser:     browser,
		ImageID:     imgID,
		BaseImageID: baseImgID,
		DiffScore:   diffScore,
		DiffImageID: diffImgID,
	}

	err = store.StoreResults(results)

	return results, err
}

func GetResults(id string) (Results, error) {
	return store.GetResults(id)
}

func AcceptTest(testID string) error {
	r, err := store.GetResults(testID)

	if err != nil {
		return err
	}

	return store.SetBaseImageID(r.ImageID, r.ProjectID, r.Branch, r.Target, r.Browser)
}

func GetImage(id string) image.Image {
	img, err := store.GetImage(id)
	if err != nil {
		panic(err)
	}

	return img
}
