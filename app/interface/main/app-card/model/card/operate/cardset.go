package operate

import (
	"encoding/json"
)

type CardSet struct {
	ID        int64           `json:"id,omitempty"`
	Type      string          `json:"type,omitempty"`
	Value     int64           `json:"value,omitempty"`
	Title     string          `json:"title,omitempty"`
	LongTitle string          `json:"long_title,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
}
