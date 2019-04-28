package model

// FaceApply face record.
type FaceApply struct {
	ID         int64  `json:"-"`
	Mid        int64  `json:"mid"`
	OldFace    string `json:"old_face"`
	NewFace    string `json:"new_face"`
	ApplyTime  int64  `json:"apply_time"`
	Status     string `json:"status"`
	Operator   string `json:"operator"`
	ModifyTime string `json:"modify_time"`
}
