package model

// all const variable used in mask
const (
	MaskPlatWeb int8 = 0
	MaskPlatMbl int8 = 1
)

// MaskLists mask lists
type MaskLists struct {
	Onoff int64    `json:"onoff"`
	Cid   int64    `json:"cid,omitempty"`
	Fps   int64    `json:"fps,omitempty"`
	Time  int64    `json:"time,omitempty"`
	Lists []string `json:"list,omitempty"`
}

// Mask mask info
type Mask struct {
	Cid     int64  `json:"cid,omitempty"`
	Plat    int8   `json:"plat,omitempty"`
	FPS     int32  `json:"fps,omitempty"`
	Time    int64  `json:"time,omitempty"`
	MaskURL string `json:"mask_url,omitempty"`
}
