package model

import (
	"go-common/library/time"
)

// VipPoint vip_point table
type VipPoint struct {
	ID           int64 `json:"id"`
	Mid          int64 `json:"mid"`
	PointBalance int64 `json:"point_balance"`
	Ver          int64 `json:"ver"`
}

//VipPointChangeHistory vip_point_change_history table
type VipPointChangeHistory struct {
	ID           int64     `json:"id"`
	Mid          int64     `json:"mid"`
	Point        int64     `json:"point"`
	OrderID      string    `json:"order_id"`
	ChangeType   int8      `json:"change_type"`
	ChangeTime   time.Time `json:"change_time"`
	RelationID   string    `json:"relation_id"`
	PointBalance int64     `json:"point_balance"`
	Remark       string    `json:"remark"`
	Operator     string    `json:"operator"`
}

//VipPointChangeHistoryMsg get databus json data
type VipPointChangeHistoryMsg struct {
	ID           int64  `json:"id"`
	Mid          int64  `json:"mid"`
	Point        int64  `json:"point"`
	OrderID      string `json:"order_id"`
	ChangeType   int8   `json:"change_type"`
	ChangeTime   string `json:"change_time"`
	RelationID   string `json:"relation_id"`
	PointBalance int64  `json:"point_balance"`
	Remark       string `json:"remark"`
	Operator     string `json:"operator"`
}

// PointInfo def.
type PointInfo struct {
	Mid          int64 `protobuf:"varint,1,opt,name=Mid,proto3" json:"mid"`
	PointBalance int64 `protobuf:"varint,2,opt,name=PointBalance,proto3" json:"pointBalance"`
	Ver          int64 `protobuf:"varint,3,opt,name=Ver,proto3" json:"ver"`
}

//PointHistory  point history
type PointHistory struct {
	ID           int64     `json:"id"`
	Mid          int64     `json:"mid"`
	Point        int64     `json:"point"`
	OrderID      string    `json:"orderID"`
	ChangeType   int       `json:"changeType"`
	ChangeTime   time.Time `json:"changeTime"`
	RelationID   string    `json:"relationID"`
	PointBalance int64     `json:"pointBalance"`
	Remark       string    `json:"remark"`
	Operator     string    `json:"operator"`
}
