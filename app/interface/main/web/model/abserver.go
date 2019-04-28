package model

import "encoding/json"

// AbServer .
type AbServer struct {
	Hit    json.RawMessage `json:"hit"`
	Expire int             `json:"expire"`
	Vars   json.RawMessage `json:"vars"`
}
