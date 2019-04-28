package model

// PendantInfo is.
type PendantInfo struct {
	Pid    int    `json:"pid"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Expire int    `json:"expire"`
}

// NameplateInfo is.
type NameplateInfo struct {
	Nid        int    `json:"nid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageSmall string `json:"image_small"`
	Level      string `json:"level"`
	Condition  string `json:"condition"`
}
