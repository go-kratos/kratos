package model

import (
	libtime "go-common/library/time"
)

//..
const (
	//初始状态
	TaskStateInit = int8(0)
	//已派发
	TaskStateDispatch = int8(1)
	//延迟
	TaskStateDelay = int8(2)
	//任务提交
	TaskStateSubmit = int8(3)
	//资源列表提交
	TaskStateRscSb = int8(4)
	//任务关闭
	TaskStateClosed = int8(5)

	// ActionCreate 生成任务
	ActionCreate = uint8(0)
	// ActionSeize 抢占任务
	ActionSeize = uint8(1)
	// ActionRelease 释放任务
	ActionRelease = uint8(2)
	// ActionDelay 延迟任务
	ActionDelay = uint8(3)
	// ActionSubmit 提交任务
	ActionSubmit = uint8(4)
	// ActionUnknow 其他变更
	ActionUnknow = uint8(5)

	LogBusinessTask     = int(232)
	LogTypeTaskDispatch = int(1)
	LogTypeTaskConsumer = int(2)
	LogTYpeTaskWeight   = int(3)

	// WeightTypeCycle 周期权重
	WeightTypeCycle = int8(0)
	// WeightTypeConst 定值权重
	WeightTypeConst = int8(1)
)

const (
	// ConfigStateOn .
	ConfigStateOn = int8(0)
	// ConfigStateOff .
	ConfigStateOff = int8(1)

	// ConsumerStateOn on
	ConsumerStateOn = int8(1)
	// ConsumerStateOff off
	ConsumerStateOff = int8(0)

	// ActionConsumerOff .
	ActionConsumerOff = int8(0)
	// ActionConsumerOn .
	ActionConsumerOn = int8(1)

	// TaskConfigAssign 指派
	TaskConfigAssign = int8(1)
	// TaskConfigRangeWeight 权重
	TaskConfigRangeWeight = int8(2)
	// TaskConfigEqualWeight 权重
	TaskConfigEqualWeight = int8(3)

	// TaskRoleMember 组员
	TaskRoleMember = int8(1)
	// TaskRoleLeader 组长
	TaskRoleLeader = int8(2)
)

// WeightItem 权重值
type WeightItem struct {
	ID     int64
	Weight int64
}

// Task ..
type Task struct {
	ID         int64   `form:"id" json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	BusinessID int64   `form:"business_id" json:"business_id" gorm:"column:business_id"`
	FlowID     int64   `form:"flow_id" json:"flow_id" gorm:"column:flow_id"`
	RID        int64   `form:"rid" json:"rid" gorm:"column:rid"`
	AdminID    int64   `form:"admin_id" json:"admin_id" gorm:"column:admin_id"`
	UID        int64   `form:"uid" json:"uid" gorm:"column:uid"`
	State      int8    `form:"state" json:"state" gorm:"column:state"`
	Weight     int64   `form:"weight" json:"weight" gorm:"column:weight"`
	Utime      int64   `form:"utime" json:"utime" gorm:"column:utime"`
	Gtime      IntTime `form:"gtime" json:"gtime" gorm:"column:gtime"`
	MID        int64   `form:"mid" json:"mid" gorm:"column:mid"`
	Fans       int64   `form:"fans" json:"fans" gorm:"column:fans"`
	Group      string  `form:"group" json:"group" gorm:"column:group"`
	Reason     string  `form:"reason" json:"reason" grom:"column:reason"`
	Ctime      IntTime `form:"ctime" json:"ctime" gorm:"column:ctime"`
	Mtime      IntTime `form:"mtime" json:"mtime" gorm:"column:mtime"`
}

// WeightLog task log
type WeightLog struct {
	UPtime      string        `json:"uptime"`
	Mid         int64         `json:"mid"`
	Fans        int64         `json:"fans"`
	FansWeight  int64         `json:"fans_weight"`
	Group       string        `json:"group"`
	GroupWeight int64         `json:"group_weight"`
	WaitTime    string        `json:"wait_time"`
	WaitWeight  int64         `json:"wait_weight"`
	EqualWeight int64         `json:"config_weight"`
	ConfigItems []*ConfigItem `json:"config_items"`
	Weight      int64         `json:"weight"`
}

// ConfigItem .
type ConfigItem struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Uname string `json:"uname"`
}

// EqualWeightConfig 等值权重
type EqualWeightConfig struct {
	Uname       string //  配置人
	Description string //  描述
	Name        string `json:"name"` // taskid 或者 mid
	IDs         string `json:"ids"`
	Weight      int64  `json:"weight"`
	Type        int8   `json:"type"` // 周期或者定值
}

// RangeWeightConfig 权重
type RangeWeightConfig struct {
	Name  string         `json:"name"`
	Range []*RangeConfig `json:"range"`
}

// RangeConfig 范围配置
type RangeConfig struct {
	Threshold int64 `json:"threshold"`
	Weight    int64 `json:"weight"`
}

// AssignConfig 指派
type AssignConfig struct {
	Admin int64   `json:"-"`
	Mids  []int64 `json:"mids"`
	Uids  []int64 `json:"uids"`
}

// TaskConfig .
type TaskConfig struct {
	ID          int64        `form:"id" json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	ConfJSON    string       `json:"conf_json" gorm:"column:conf_json"`
	ConfType    int8         `form:"conf_type" json:"conf_type" gorm:"column:conf_type"`
	BusinessID  int64        `form:"business_id" json:"business_id" gorm:"column:business_id"`
	FlowID      int64        `form:"flow_id" json:"flow_id" gorm:"column:flow_id"`
	Btime       libtime.Time `form:"btime" json:"btime" gorm:"column:btime"`
	Etime       libtime.Time `form:"etime" json:"etime" gorm:"column:etime"`
	State       int8         `form:"state" json:"state" gorm:"column:state"`
	UID         int64        `form:"uid" json:"uid" gorm:"column:uid"`
	Uname       string       `form:"uname" json:"uname" gorm:"column:uname"`
	Description string       `form:"description" json:"description" gorm:"column:description"`
	Ctime       libtime.Time `form:"ctime" json:"ctime" gorm:"column:ctime"`
	Mtime       libtime.Time `form:"mtime" json:"mtime" gorm:"column:mtime"`
}

// TableName for orm
func (TaskConfig) TableName() string {
	return "task_config"
}

// WeightOPT .
type WeightOPT struct {
	BusinessID   int64
	FlowID       int64
	TopListLen   int64
	BatchListLen int64
	RedisListLen int64
	DbListLen    int64
	AssignLen    int64
	Minute       int64
}
