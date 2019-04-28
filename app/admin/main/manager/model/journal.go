package model

import (
	"encoding/json"
)

// LogRes .
type LogRes struct {
	Code int       `json:"code"`
	Data *LogChild `json:"data"`
}

// LogChild .
type LogChild struct {
	Result []json.RawMessage `json:"result"`
}
