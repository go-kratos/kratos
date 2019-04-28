package model

// Statistics def.
type Statistics struct {
	Quantity  int64  `json:"quantity"`
	EventID   int64  `json:"event_id"`
	EventName string `json:"event_name"`
}
