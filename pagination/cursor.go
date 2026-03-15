package pagination

import (
	"encoding/base64"
	"encoding/json"
)

type Cursor struct {
	Limit int    `query:"limit" form:"limit" json:"limit" default:"20" max:"100"`
	After string `query:"after" form:"after" json:"after"`
}

func (p *Cursor) Normalize() {
	Normalize(p)
}

func EncodeCursor(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func DecodeCursor(encoded string, dst any) error {
	b, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
