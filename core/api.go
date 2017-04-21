package core

import (
	"image"
	"time"

	"github.com/pkg/errors"
)

var store Store = NewBoltStore("optician.db")

// TESTS
func Results() ([]Result, error) {
	return store.GetResults()
}

func AddCase(c Case) (Result, error) {
	randID := RandStringBytes(14)

	testImage := c.Image
	projectID := c.ProjectID
	branch := c.Branch
	target := c.Target
	browser := c.Browser

	imgID, err := store.StoreImage(testImage)

	baseImgID, err := store.GetBaseImageID(projectID, branch, target, browser)
	if err != nil {
		if err == NotFoundError {
			// IF no base image found, set this as base image
			baseImgID = imgID
			store.SetBaseImageID(baseImgID, projectID, branch, target, browser)
		} else {
			return Result{}, errors.Wrap(err, "error getting base image ID")
		}
	}

	maskID, err := store.GetBaseMaskID(projectID, branch, target, browser)
	if err != nil {
		if err == NotFoundError {
			maskID = "nomask"
		} else {
			return Result{}, errors.Wrap(err, "error getting base mask id")
		}
	}

	results := Result{
		ID:          randID,
		ProjectID:   projectID,
		Branch:      branch,
		Target:      target,
		Browser:     browser,
		ImageID:     imgID,
		MaskID:      maskID,
		BaseImageID: baseImgID,
		Timestamp:   time.Now(),
	}

	err = RunTest(&results)

	err = store.StoreResult(results)

	return results, err
}

func GetTest(id string) (Result, error) {
	return store.GetResult(id)
}

func AcceptTest(testID string) error {
	test, err := store.GetResult(testID)

	if err != nil {
		return err
	}

	if testID != GetLastTest(test.ProjectID, test.Branch, test.Target, test.Browser) {
		return errors.New("Cannot accept an old test")
	}

	return store.SetBaseImageID(test.ImageID, test.ProjectID, test.Branch, test.Target, test.Browser)
}

func GetLastTest(projectID, branch, target, browser string) string {
	return store.GetLastResult(projectID, branch, target, browser).ID
}

// IMAGES

func GetImage(id string) image.Image {
	img, err := store.GetImage(id)
	if err != nil {
		panic(err)
	}

	return img
}

// MASKS

func GetMask(id string) ([]image.Rectangle, error) {
	return store.GetMask(id)
}

func MaskTest(testID string, mask []image.Rectangle) (Result, error) {
	test, err := GetTest(testID)
	if err != nil {
		return Result{}, err
	}

	if testID != GetLastTest(test.ProjectID, test.Branch, test.Target, test.Browser) {
		return Result{}, errors.New("Cannot add masks based on an old test")
	}

	maskID, err := store.StoreMask(mask)
	if err != nil {
		return Result{}, err
	}

	err = store.SetBaseMaskID(maskID, test.ProjectID, test.Branch, test.Target, test.Browser)
	if err != nil {
		return Result{}, err
	}

	test.MaskID = maskID

	err = RunTest(&test)
	if err != nil {
		return Result{}, err
	}

	return test, nil
}
