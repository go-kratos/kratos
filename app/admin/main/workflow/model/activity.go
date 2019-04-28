package model

// Activities is the model for all challenge activities
type Activities struct {
	Logs   []*WLog  `json:"logs"`
	Events []*Event `json:"events"`
}
