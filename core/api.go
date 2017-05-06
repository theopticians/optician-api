package core

import (
	"image"
	"time"

	"github.com/pkg/errors"
)

var store Store = NewBoltStore("optician.db")

func Batchs() ([]BatchInfo, error) {
	return store.GetBatchs()
}

// TESTS
func Results() ([]Result, error) {
	return store.GetResults()
}

func ResultsByBatchs(batch string) ([]Result, error) {
	return store.GetResultsByBatch(batch)
}

func AddCase(c Case) (Result, error) {

	testImage := c.Image
	projectID := c.ProjectID
	branch := c.Branch
	target := c.Target
	browser := c.Browser
	batch := c.Batch

	if batchIsOld(batch) {
		return Result{}, errors.New("The batch " + batch + " is too old, start a new one")
	}

	if batchHasTest(batch, projectID, branch, target, browser) {
		return Result{}, errors.New("The batch " + batch + " already has this test")
	}

	if batchHasDifferentBranch(batch, branch) {
		return Result{}, errors.New("The same batch was used for a different branch. Only one branch can be tested in a batch")
	}

	randID := RandStringBytes(14)

	imgID, err := store.StoreImage(testImage)

	baseImgID, err := store.GetBaseImageID(projectID, branch, target, browser)
	if err != nil {
		if err == KeyNotFoundError {
			// IF no base image found, set this as base image
			baseImgID = imgID
			store.SetBaseImageID(baseImgID, projectID, branch, target, browser)
		} else {
			return Result{}, errors.Wrap(err, "error getting base image ID")
		}
	}

	maskID, err := store.GetBaseMaskID(projectID, branch, target, browser)
	if err != nil {
		if err == KeyNotFoundError {
			maskID = "nomask"
		} else {
			return Result{}, errors.Wrap(err, "error getting base mask id")
		}
	}

	results := Result{
		ID:          randID,
		ProjectID:   projectID,
		Branch:      branch,
		Batch:       c.Batch,
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

	lastTest, err := store.GetLastResult(test.ProjectID, test.Branch, test.Target, test.Browser)

	if err != nil {
		return err
	}

	if testID != lastTest.ID {
		return errors.New("Cannot accept an old test. Last test is " + lastTest.ID)
	}

	return store.SetBaseImageID(test.ImageID, test.ProjectID, test.Branch, test.Target, test.Browser)
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

	lastTest, err := store.GetLastResult(test.ProjectID, test.Branch, test.Target, test.Browser)

	if err != nil {
		return Result{}, err
	}

	if testID != lastTest.ID {
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

	err = store.StoreResult(test)

	if err != nil {
		return Result{}, err
	}

	return test, nil
}
