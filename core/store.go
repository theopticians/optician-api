package core

import "image"

type Store interface {
	Close()
	GetTestList() []string
	GetTest(string) (Test, error)
	GetMask(string) ([]image.Rectangle, error)
	StoreMask(masks []image.Rectangle) (string, error)
	StoreTest(Test) error
	GetImage(string) (image.Image, error)
	StoreImage(image.Image) (string, error)
	GetBaseImageID(projectID, branch, target, browser string) (string, error)
	SetBaseImageID(baseImageID, projectID, branch, target, browser string) error
	GetBaseMaskID(projectID, branch, target, browser string) (string, error)
	SetBaseMaskID(baseImageID, projectID, branch, target, browser string) error
}
