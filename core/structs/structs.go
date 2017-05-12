package structs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"image"
	"time"
)

type Result struct {
	ID           string    `json:"id"`
	Project      string    `json:"project"`
	Branch       string    `json:"branch"`
	Batch        string    `json:"batch"`
	Target       string    `json:"target"`
	Browser      string    `json:"browser"`
	MaskID       string    `json:"mask"`
	DiffScore    float64   `json:"diffscore"`
	ImageID      string    `json:"image"`
	BaseImageID  string    `json:"baseimage"`
	DiffImageID  string    `json:"diffimage"`
	DiffClusters Mask      `json:"diffclusters"`
	Timestamp    time.Time `json:"timestamp"`
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

type Mask []image.Rectangle

func (m *Mask) UnmarshalJSON(data []byte) error {

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

	newMask := Mask(make([]image.Rectangle, len(aux)))

	for i := 0; i < len(aux); i++ {
		newMask[i].Min = image.Point{X: aux[i].X, Y: aux[i].Y}
		newMask[i].Max = image.Point{X: aux[i].X + aux[i].Width, Y: aux[i].Y + aux[i].Height}
	}

	*m = newMask
	return nil
}

func (m Mask) MarshalJSON() ([]byte, error) {
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

func (m Mask) Value() (driver.Value, error) {
	b, err := m.MarshalJSON()

	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (m *Mask) Scan(value interface{}) error {
	// if value is nil, false
	if value == nil {
		*m = nil
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		// if this is a string type
		if v, ok := bv.(string); ok {
			return m.UnmarshalJSON([]byte(v))
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan Mask")
}
