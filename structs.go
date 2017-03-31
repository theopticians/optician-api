package main

import (
	"encoding/json"
	"image"
)

type Test struct {
	ProjectID string `json:"projectid"`
	Branch    string `json:"branch"`
	Target    string `json:"target"`
	Browser   string `json:"browser"`
	Batch     string `json:"browser"`
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
