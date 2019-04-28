package up

//Switch for up switch controll.
type Switch struct {
	State int32  `json:"state"`
	Show  int    `json:"show"`
	Face  string `json:"face"`
}

// SpecialGroup UP主分组关联关系
type SpecialGroup struct {
	ID        int64  `json:"id"`
	MID       int64  `json:"mid"`
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}
