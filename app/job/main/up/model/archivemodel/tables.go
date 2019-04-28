package archivemodel

//ArchiveCanal struct from cannal
type ArchiveCanal struct {
	ID          int64  `json:"id"`
	AID         int64  `json:"aid"`
	Mid         int64  `json:"mid"`
	TypeID      int16  `json:"typeid"`
	Videos      int    `json:"videos"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Content     string `json:"content"`
	Duration    int    `json:"duration"`
	Attribute   int32  `json:"attribute"`
	Copyright   int8   `json:"copyright"`
	Access      int    `json:"access"`
	State       int    `json:"state"`
	MissionID   int64  `json:"mission_id"`
	OrderID     int64  `json:"order_id"`
	RedirectURL string `json:"redirect_url"`
	Forward     int64  `json:"forward"`
	Dynamic     string `json:"dynamic"`
}

// ArchiveStaff state值
const (
	StaffStateNormal    = 1 // 正常
	StaffStateDismissed = 2 // 解除
)

//ArchiveStaff .
type ArchiveStaff struct {
	ID           int64  `json:"id"`
	Aid          int64  `json:"aid"`
	Mid          int64  `json:"mid"`
	StaffMid     int64  `json:"staff_mid"`
	StaffTitle   string `json:"staff_title"`
	StaffTitleId int64  `json:"staff_title_id"`
	State        int64  `json:"state"`
}
