package core

import "time"

type Test struct {
	TestID      string    `json:"testid"`
	ProjectID   string    `json:"projectid"`
	Branch      string    `json:"branch"`
	Batch       string    `json:"batch"`
	Target      string    `json:"target"`
	Browser     string    `json:"browser"`
	MaskID      string    `json:"mask"`
	DiffScore   float64   `json:"diffscore"`
	ImageID     string    `json:"image"`
	BaseImageID string    `json:"baseimage"`
	DiffImageID string    `json:"diffimage"`
	Timestamp   time.Time `json:"timestamp"`
}
