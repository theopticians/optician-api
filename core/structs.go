package core

type Results struct {
	TestID      string  `json:"testid"`
	ProjectID   string  `json:"projectid"`
	Branch      string  `json:"branch"`
	Batch       string  `json:"batch"`
	Target      string  `json:"target"`
	Browser     string  `json:"browser"`
	DiffScore   float64 `json:"diffscore"`
	ImageID     string  `json:"image"`
	BaseImageID string  `json:"baseimage"`
	DiffImageID string  `json:"diffimage"`
}
