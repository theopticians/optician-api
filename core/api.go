package core

import (
	"image"
	"time"

	"github.com/pkg/errors"
)

var store Store = NewBoltStore("optician.db")

// TESTS
func TestList() []string {
	return store.GetTestList()
}

func NewTest(testImage image.Image, projectID, branch, target, browser string) (Test, error) {
	randID := RandStringBytes(14)

	imgID, err := store.StoreImage(testImage)

	baseImgID, err := store.GetBaseImageID(projectID, branch, target, browser)
	if err != nil {
		if err == NotFoundError {
			// IF no base image found, set this as base image
			baseImgID = imgID
			store.SetBaseImageID(baseImgID, projectID, branch, target, browser)
		} else {
			return Test{}, errors.Wrap(err, "error getting base image ID")
		}
	}

	maskID, err := store.GetBaseMaskID(projectID, branch, target, browser)
	if err != nil {
		if err == NotFoundError {
			maskID = "nomask"
		} else {
			return Test{}, errors.Wrap(err, "error getting base mask id")
		}
	}

	results := Test{
		TestID:      randID,
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

	err = store.StoreTest(results)

	return results, err
}

func GetTest(id string) (Test, error) {
	return store.GetTest(id)
}

func AcceptTest(testID string) error {
	test, err := store.GetTest(testID)

	if err != nil {
		return err
	}

	if testID != GetLastTest(test.ProjectID, test.Branch, test.Target, test.Browser) {
		return errors.New("Cannot accept an old test")
	}

	return store.SetBaseImageID(test.ImageID, test.ProjectID, test.Branch, test.Target, test.Browser)
}

func GetLastTest(projectID, branch, target, browser string) string {
	return store.GetLastTest(projectID, branch, target, browser).TestID
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

func MaskTest(testID string, mask []image.Rectangle) (Test, error) {
	test, err := GetTest(testID)
	if err != nil {
		return Test{}, err
	}

	if testID != GetLastTest(test.ProjectID, test.Branch, test.Target, test.Browser) {
		return Test{}, errors.New("Cannot add masks based on an old test")
	}

	maskID, err := store.StoreMask(mask)
	if err != nil {
		return Test{}, err
	}

	err = store.SetBaseMaskID(maskID, test.ProjectID, test.Branch, test.Target, test.Browser)
	if err != nil {
		return Test{}, err
	}

	test.MaskID = maskID

	err = RunTest(&test)
	if err != nil {
		return Test{}, err
	}

	return test, nil
}
