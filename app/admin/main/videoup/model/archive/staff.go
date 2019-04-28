package archive

//StaffParam .
type StaffParam struct {
	MID   int64  `json:"mid"`
	Title string `json:"title"`
}

//StaffBatchParam  批量提交的staff参数
type StaffBatchParam struct {
	AID      int64         `json:"aid"`
	SyncAttr bool          `json:"sync_attr"`
	Staffs   []*StaffParam `json:"staffs"`
}

//Staff .
type Staff struct {
	ID           int64  `json:"id"`
	AID          int64  `json:"aid"`
	MID          int64  `json:"mid"`
	StaffMID     int64  `json:"staff_mid"`
	StaffTitle   string `json:"staff_title"`
	StaffName    string `json:"staff_name"`
	StaffTitleID int64  `json:"staff_title_id"`
	State        int8   `json:"state"`
}
