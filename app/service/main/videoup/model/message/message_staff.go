package message

const (
	//RouteBplusStaff pgc提交
	RouteBplusStaff = "bplus_staff"
)

type StaffBox struct {
	RID      int64        `json:"rid"`
	Type     int8         `json:"type"`
	ADDStaff []*StaffItem `json:"added_staffs"`
	DelStaff []*StaffItem `json:"removed_staffs"`
}

type StaffItem struct {
	Type int8  `json:"uid_type"`
	UID  int64 `json:"uid"`
}

// BplusCardMsg 粉丝动态databus消息
type BplusCardMsg struct {
	Outbox *StaffBox `json:"outbox"`
}
