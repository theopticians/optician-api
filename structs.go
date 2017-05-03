package main

import (
	"encoding/json"
	"image"

	"github.com/theopticians/optician-api/core"
)

type ApiResult core.Result
type ApiCase core.Case
type ApiMask []image.Rectangle

func (u *ApiCase) UnmarshalJSON(data []byte) error {
	type Alias ApiCase
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

func (r *ApiResult) MarshalJSON() ([]byte, error) {
	mask, err := core.GetMask(r.MaskID)
	if err != nil {
		return nil, err
	}

	type Alias ApiResult
	return json.Marshal(&struct {
		Mask   ApiMask `json:"mask"`
		MaskID string  `json:"-"`
		*Alias
	}{
		Mask:  ApiMask(mask),
		Alias: (*Alias)(r),
	})
}

func (m *ApiMask) UnmarshalJSON(data []byte) error {

	aux := []struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	}{}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	newMask := ApiMask(make([]image.Rectangle, len(aux)))

	for i := 0; i < len(aux); i++ {
		newMask[i].Min = image.Point{X: aux[i].X, Y: aux[i].Y}
		newMask[i].Max = image.Point{X: aux[i].X + aux[i].Width, Y: aux[i].Y + aux[i].Height}
	}

	*m = newMask
	return nil
}

func (m ApiMask) MarshalJSON() ([]byte, error) {
	aux := make([]struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	}, len(m))

	for i := 0; i < len(aux); i++ {
		aux[i].X = m[i].Min.X
		aux[i].Y = m[i].Min.Y
		aux[i].Width = m[i].Dx()
		aux[i].Height = m[i].Dy()
	}

	return json.Marshal(&aux)
}
