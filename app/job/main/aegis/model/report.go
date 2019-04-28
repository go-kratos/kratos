package model

import (
	"fmt"
	"strconv"
	"strings"
)

//hash fields
const (
	Dispatch = "ds"
	Delay    = "dy"
	Submit   = "st_%d_%d" // 参数1:提交状态(任务提交，资源提交，任务关闭) 参数2:提交前属于谁
	Release  = "rl"
	RscState = "rs_%d"
	UseTime  = "ut"
	SetKey   = "report_set"

	//type
	TypeMeta  = int8(0)
	TypeTotal = int8(1)
)

//RIR resource item report
type RIR struct {
	BizID  int64
	FlowID int64
	UID    int64
	RID    int64
}

//Report .
type Report struct {
	ID         int64  `gorm:"AUTO_INCREMENT;primary_key;"`
	BusinessID int64  `gorm:"column:business_id"`
	FlowID     int64  `gorm:"column:flow_id"`
	UID        int64  `gorm:"column:uid"`
	TYPE       int8   `gorm:"column:type"`
	Content    []byte `gorm:"column:content"`
}

//TableName .
func (r Report) TableName() string {
	return "task_report"
}

//PersonalHashKey .
func PersonalHashKey(bizid, flowid, uid int64) string {
	return fmt.Sprintf("report_hash_%d_%d_%d", bizid, flowid, uid)
}

//TotalHashKey .
func TotalHashKey(bizid, flowid int64) string {
	return fmt.Sprintf("total_inout_%d_%d_%d", bizid, flowid, 0)
}

//ParseKey .
func ParseKey(key string) (tp int8, bizid, flowid, uid int, err error) {
	arr := strings.Split(key, "_")
	if len(arr) != 5 {
		err = fmt.Errorf(key)
		return
	}
	prefix := arr[0] + "_" + arr[1]
	switch prefix {
	case "report_hash":
		tp = TypeMeta
	case "total_inout":
		tp = TypeTotal
	default:
		err = fmt.Errorf(key)
		return
	}

	if bizid, err = strconv.Atoi(arr[2]); err != nil || bizid == 0 {
		err = fmt.Errorf(key)
		return
	}
	if flowid, err = strconv.Atoi(arr[3]); err != nil || flowid == 0 {
		err = fmt.Errorf(key)
		return
	}
	if uid, err = strconv.Atoi(arr[4]); err != nil {
		err = fmt.Errorf(key)
		return
	}
	return
}
