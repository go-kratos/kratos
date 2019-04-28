package model

import (
	xtime "go-common/library/time"
)

// Log Module Field defination
const (
	// report business = 12
	WLogModuleChallenge     = 1
	WLogModuleTag           = 2
	WLogModuleControl       = 3
	WLogModuleGroup         = 4
	WLogModuleReply         = 5 // modify business_state
	WLogModulePublicReferee = 6
	WLogModuleRoleShift     = 7  // 流转 (同一个执行方)
	WLogModuleDispose       = 8  // content dispose 操作内容对象
	WLogModuleAddMoral      = 20 // 扣节操
	WLogModuleBlock         = 21 // 封禁

	// report business = 11
	FeedBackTypeNotifyUserReceived = 2
	FeedBackTypeNotifyUserDisposed = 3
	FeedBackTypeReply              = 5
)

// LogSlice is a Log slice struct
type LogSlice []*WLog

// Log model is the universal model
// Will record any management actions
type WLog struct {
	Lid         int32       `json:"lid"`
	AdminID     int64       `json:"adminid"`
	Admin       string      `json:"admin"`
	Oid         int64       `json:"oid"`
	Business    int8        `json:"business"`
	Target      int64       `json:"target"`
	Module      int8        `json:"module"`
	Remark      string      `json:"remark"`
	Note        string      `json:"note"`
	CTime       xtime.Time  `json:"ctime"`
	MTime       xtime.Time  `json:"mtime"`
	Meta        interface{} `json:"meta"`
	ReportCTime string      `json:"report_ctime"`
	Mid         int64       `json:"mid"`
	TypeID      int64       `json:"type_id"`
	TimeConsume int64       `json:"time_consume"`
	OpType      string      `json:"op_type"`
	PreRid      string      `json:"pre_rid"`
	Param       interface{} `json:"param"`
	Mids        []int64     `json:"mids"` //对被举报人的批量操作
}
