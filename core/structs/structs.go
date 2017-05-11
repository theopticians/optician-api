package structs

import (
	"image"
	"time"
)

type Result struct {
	ID           string            `json:"id"`
	Project      string            `json:"project"`
	Branch       string            `json:"branch"`
	Batch        string            `json:"batch"`
	Target       string            `json:"target"`
	Browser      string            `json:"browser"`
	MaskID       string            `json:"mask"`
	DiffScore    float64           `json:"diffscore"`
	ImageID      string            `json:"image"`
	BaseImageID  string            `json:"baseimage"`
	DiffImageID  string            `json:"diffimage"`
	DiffClusters []image.Rectangle `json:"diffclusters"`
	Timestamp    time.Time         `json:"timestamp"`
}

type BatchInfo struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Failed    int       `json:"failed"`
	Project   string    `json:"project"`
}

type Case struct {
	ProjectID string `json:"projectid"`
	Branch    string `json:"branch"`
	Target    string `json:"target"`
	Browser   string `json:"browser"`
	Batch     string `json:"batch"`
	Image     image.Image
}
