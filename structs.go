package main

import (
	"encoding/json"

	"github.com/theopticians/optician-api/core"
	"github.com/theopticians/optician-api/core/structs"
)

type ApiResult structs.Result
type ApiCase structs.Case

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
		Mask   structs.Mask `json:"mask"`
		MaskID string       `json:"-"`
		*Alias
	}{
		Mask:  structs.Mask(mask),
		Alias: (*Alias)(r),
	})
}
