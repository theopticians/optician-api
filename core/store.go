package core

import (
	"image"

	"github.com/theopticians/optician-api/core/structs"
)

type Store interface {
	Close()
	GetResults() ([]structs.Result, error)
	GetResultsByBatch(string) ([]structs.Result, error)
	GetBatchs() ([]structs.BatchInfo, error)
	GetResult(string) (structs.Result, error)
	GetLastResult(projectID, branch, target, browser string) (structs.Result, error)
	GetMask(string) ([]image.Rectangle, error)
	StoreMask(masks []image.Rectangle) (string, error)
	StoreResult(structs.Result) error
	GetImage(string) (image.Image, error)
	StoreImage(image.Image) (string, error)
	GetBaseImageID(projectID, branch, target, browser string) (string, error)
	SetBaseImageID(baseImageID, projectID, branch, target, browser string) error
	GetBaseMaskID(projectID, branch, target, browser string) (string, error)
	SetBaseMaskID(baseImageID, projectID, branch, target, browser string) error
}
