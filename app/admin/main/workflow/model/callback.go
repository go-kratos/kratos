package model

import (
	xtime "go-common/library/time"
)

// consts for callback
const (
	GroupSetResult        = "group.SetResult"
	BatchGroupSetResult   = "group.BatchSetResult"
	ChallSetResult        = "chall.SetResult"
	BatchChallSetResult   = "chall.BatchSetResult"
	GroupSetState         = "group.SetState"
	GroupSetPublicReferee = "group.SetPublicReferee"

	CallbackDisable = 0
	CallbackEnable  = 1
)

// Slice is the slice model for callback
type CallbackSlice []*Callback

// Callback is the workflow callback model
type Callback struct {
	CbID        int32      `json:"cbid" gorm:"column:id"`
	URL         string     `json:"url" gorm:"column:url"`
	Business    int8       `json:"business" gorm:"column:business"`
	IsSobot     bool       `json:"is_sobot" gorm:"column:is_sobot"`
	State       int8       `json:"state" gorm:"column:state"`
	ExternalAPI string     `json:"external_api" gorm:"column:external_api"`
	SourceAPI   string     `json:"source_api" gorm:"column:source_api"`
	CTime       xtime.Time `json:"ctime" gorm:"column:ctime"`
	MTime       xtime.Time `json:"mtime" gorm:"column:mtime"`
}

// Actor for callback payload
type Actor struct {
	AdminID   int64  `json:"admin_id"`
	AdminName string `json:"admin_name"`
}

// Payload is the payload model for callback
type Payload struct {
	Bid       int           `json:"bid"`
	Verb      string        `json:"verb"`
	Actor     Actor         `json:"actor"`
	CTime     xtime.Time    `json:"ctime"`
	Object    interface{}   `json:"object"`    //处理请求参数
	Target    interface{}   `json:"target"`    //被修改的工单或工单详情
	Targets   []interface{} `json:"targets"`   //所有被修改的工单或工单详情
	Influence interface{}   `json:"influence"` //业务自定义 Deprecated
	Extra     interface{}   `json:"extra"`     //业务自定义
}

// TableName is used to identify table name for gorm
func (Callback) TableName() string {
	return "workflow_callback"
}
