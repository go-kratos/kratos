package model

// consts
const (
	ActUpdateExp   = "updateExp"
	ActUpdateLevel = "updateLevel"
	ActUpdateFace  = "updateFace"
	ActUpdateMoral = "updateMoral"
	ActUpdateUname = "updateUname"
)

// NotifyInfo notify info.
type NotifyInfo struct {
	Uname   string `json:"uname"`
	Mid     int64  `json:"mid"`
	Type    string `json:"type"`
	NewName string `json:"newName"`
	Action  string `json:"action"`
}
