package core

import "image"

type Store interface {
	Close()
	GetResults() ([]Result, error)
	GetResult(string) (Result, error)
	GetLastResult(projectID, branch, target, browser string) Result
	GetMask(string) ([]image.Rectangle, error)
	StoreMask(masks []image.Rectangle) (string, error)
	StoreResult(Result) error
	GetImage(string) (image.Image, error)
	StoreImage(image.Image) (string, error)
	GetBaseImageID(projectID, branch, target, browser string) (string, error)
	SetBaseImageID(baseImageID, projectID, branch, target, browser string) error
	GetBaseMaskID(projectID, branch, target, browser string) (string, error)
	SetBaseMaskID(baseImageID, projectID, branch, target, browser string) error
}
