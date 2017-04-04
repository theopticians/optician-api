package core

import "image"

type Store interface {
	Close()
	GetTestList() []string
	GetResults(string) (Results, error)
	GetMasks(projectID, branch, target, browser string) ([]image.Rectangle, error)
	StoreMasks(masks []image.Rectangle, projectID, branch, target, browser string) error
	StoreResults(Results) error
	GetImage(string) (image.Image, error)
	StoreImage(image.Image) (string, error)
	GetBaseImageID(projectID, branch, target, browser string) (string, error)
	SetBaseImageID(baseImageID, projectID, branch, target, browser string) error
}
