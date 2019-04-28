package net

const (
	//TableFlowResource .
	TableFlowResource = "net_flow_resource"

	//FRStateDeleted 流程取消状态
	FRStateDeleted = -1
	//FRStateRunning 流程运行状态
	FRStateRunning = 0
)

//FlowResource 资源在审核网的运行现状
type FlowResource struct {
	ID     int64 `gorm:"primary_key" json:"id"`
	RID    int64 `gorm:"column:rid" json:"rid"`
	FlowID int64 `gorm:"column:flow_id" json:"flow_id"`
	State  int8  `gorm:"column:state" json:"state"`
	NetID  int64 `gorm:"column:net_id" json:"net_id"`
}

// TableName .
func (fr *FlowResource) TableName() string {
	return TableFlowResource
}
