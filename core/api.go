package core

import (
	"image"
	"time"

	"github.com/pkg/errors"
	"github.com/theopticians/optician-api/core/store"
	"github.com/theopticians/optician-api/core/store/bolt"
	"github.com/theopticians/optician-api/core/structs"
)

//var db store.Store = sql.NewSqlStore("postgres", "postgresql://root@localhost:26257/optician?sslmode=disable")
var db store.Store = bolt.NewBoltStore("./optician.db")

func Batchs() ([]structs.BatchInfo, error) {
	return db.GetBatchs()
}

// TESTS
func Results() ([]structs.Result, error) {
	return db.GetResults()
}

func ResultsByBatchs(batch string) ([]structs.Result, error) {
	return db.GetResultsByBatch(batch)
}

func AddCase(c structs.Case) (structs.Result, error) {

	testImage := c.Image
	projectID := c.ProjectID
	branch := c.Branch
	target := c.Target
	browser := c.Browser
	batch := c.Batch

	if batchIsOld(batch) {
		return structs.Result{}, errors.New("The batch " + batch + " is too old, start a new one")
	}

	if batchHasTest(batch, projectID, branch, target, browser) {
		return structs.Result{}, errors.New("The batch " + batch + " already has this test")
	}

	if batchHasDifferentBranch(batch, branch) {
		return structs.Result{}, errors.New("The same batch was used for a different branch. Only one branch can be tested in a batch")
	}

	randID := RandStringBytes(14)

	imgID, err := db.StoreImage(testImage)

	baseImgID, err := db.GetBaseImageID(projectID, branch, target, browser)
	if err != nil {
		if err == store.NotFoundError {
			// IF no base image found, set this as base image
			baseImgID = imgID
			db.SetBaseImageID(baseImgID, projectID, branch, target, browser)
		} else {
			return structs.Result{}, errors.Wrap(err, "error getting base image ID")
		}
	}

	maskID, err := db.GetBaseMaskID(projectID, branch, target, browser)
	if err != nil {
		if err == store.NotFoundError {
			maskID = "nomask"
		} else {
			return structs.Result{}, errors.Wrap(err, "error getting base mask id")
		}
	}

	results := structs.Result{
		ID:          randID,
		Project:     projectID,
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

	err = db.StoreResult(results)

	return results, err
}

func GetTest(id string) (structs.Result, error) {
	return db.GetResult(id)
}

func AcceptTest(testID string) error {
	test, err := db.GetResult(testID)

	if err != nil {
		return err
	}

	lastTest, err := db.GetLastResult(test.Project, test.Branch, test.Target, test.Browser)

	if err != nil {
		return err
	}

	if testID != lastTest.ID {
		return errors.New("Cannot accept an old test. Last test is " + lastTest.ID)
	}

	return db.SetBaseImageID(test.ImageID, test.Project, test.Branch, test.Target, test.Browser)
}

// IMAGES

func GetImage(id string) image.Image {
	img, err := db.GetImage(id)
	if err != nil {
		panic(err)
	}

	return img
}

// MASKS

func GetMask(id string) ([]image.Rectangle, error) {
	return db.GetMask(id)
}

func MaskTest(testID string, mask []image.Rectangle) (structs.Result, error) {
	test, err := GetTest(testID)
	if err != nil {
		return structs.Result{}, err
	}

	lastTest, err := db.GetLastResult(test.Project, test.Branch, test.Target, test.Browser)

	if err != nil {
		return structs.Result{}, err
	}

	if testID != lastTest.ID {
		return structs.Result{}, errors.New("Cannot add masks based on an old test")
	}

	maskID, err := db.StoreMask(mask)
	if err != nil {
		return structs.Result{}, err
	}

	err = db.SetBaseMaskID(maskID, test.Project, test.Branch, test.Target, test.Browser)
	if err != nil {
		return structs.Result{}, err
	}

	test.MaskID = maskID

	err = RunTest(&test)
	if err != nil {
		return structs.Result{}, err
	}

	err = db.StoreResult(test)

	if err != nil {
		return structs.Result{}, err
	}

	return test, nil
}
