package main

import (
	"encoding/json"
	"github.com/theopticians/optician-api/core"
	"image"
)

type Result core.Result
type Case core.Case

func (u *Case) UnmarshalJSON(data []byte) error {
	type Alias Case
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

func (r *Result) MarshalJSON() ([]byte, error) {
	mask, err := core.GetMask(r.MaskID)
	if err != nil {
		return nil, err
	}

	type Alias Result
	return json.Marshal(&struct {
		Mask   []image.Rectangle `json:"mask"`
		MaskID string            `json:"-"`
		*Alias
	}{
		Mask:  mask,
		Alias: (*Alias)(r),
	})
}
