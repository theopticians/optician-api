package main

import (
	"encoding/json"
	"github.com/theopticians/optician-api/core"
	"image"
)

type Test struct {
	ProjectID string `json:"projectid"`
	Branch    string `json:"branch"`
	Target    string `json:"target"`
	Browser   string `json:"browser"`
	Batch     string `json:"batch"`
	Image     image.Image
}

func (u *Test) UnmarshalJSON(data []byte) error {
	type Alias Test
	aux := struct {
		Image string `json:"image"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	u.Image = base64ToImage(aux.Image)

	return nil
}

type Results core.Test

func (r *Results) MarshalJSON() ([]byte, error) {
	mask, err := core.GetMask(r.MaskID)
	if err != nil {
		return nil, err
	}

	type Alias Results
	return json.Marshal(&struct {
		Mask   []image.Rectangle `json:"mask"`
		MaskID string            `json:"-"`
		*Alias
	}{
		Mask:  mask,
		Alias: (*Alias)(r),
	})
}
