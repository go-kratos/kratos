package task

import (
	"errors"

	"context"
	"go-common/app/admin/main/aegis/model/common"
	libtime "go-common/library/time"
)

//...
const (
	ActionConsumerOff = int8(0)
	ActionConsumerOn  = int8(1)

	TaskStateInit     = int8(0)
	TaskStateDispatch = int8(1)
	TaskStateDelay    = int8(2)
	TaskStateSubmit   = int8(3)
	TaskStateRscSb    = int8(4)
	TaskStateClosed   = int8(5)

	TaskConfigAssign      = int8(1)
	TaskConfigRangeWeight = int8(2)
	TaskConfigEqualWeight = int8(3)

	TaskRoleMember = int8(1) //组员
	TaskRoleLeader = int8(2) //组长
	TaskNoRole     = int8(0) //无身份
)

// ErrEmpty empty pool
var (
	ErrEmpty = errors.New("empty pool")
	ErrRole  = errors.New("不在用户组内")
)

// Task ..
type Task struct {
	ID         int64          `form:"id" json:"id" gorm:"AUTO_INCREMENT;primary_key;"`
	BusinessID int64          `form:"business_id" json:"business_id" gorm:"column:business_id"`
	FlowID     int64          `form:"flow_id" json:"flow_id" gorm:"column:flow_id"`
	RID        int64          `form:"rid" json:"rid" gorm:"column:rid"`
	AdminID    int64          `form:"admin_id" json:"admin_id" gorm:"column:admin_id"`
	UID        int64          `form:"uid" json:"uid" gorm:"column:uid"`
	MID        int64          `form:"mid" json:"mid" gorm:"column:mid"`
	State      int8           `form:"state" json:"state" gorm:"column:state"`
	Weight     int64          `form:"weight" json:"weight" gorm:"column:weight"`
	Utime      int64          `form:"utime" json:"utime" gorm:"column:utime"`
	Gtime      common.IntTime `form:"gtime" json:"gtime" gorm:"column:gtime"`
	Fans       int64          `form:"fans" json:"fans" gorm:"column:fans"`
	Group      string         `form:"group" json:"group" gorm:"column:group"`
	Reason     string         `form:"reason" json:"reason" grom:"column:reason"`
	Ctime      common.IntTime `form:"ctime" json:"ctime" gorm:"column:ctime"`
	Mtime      common.IntTime `form:"mtime" json:"mtime" gorm:"column:mtime"`
}

// TableName ...
func (t Task) TableName() string {
	return "task"
}

// TempOptions 中间参数
type TempOptions struct {
	BisLeader bool // 是否组长
	NoCache   bool // 不使用缓存
	Action    string
}

// NextOptions options for Next
type NextOptions struct {
	common.BaseOptions
	TempOptions
	SeizeCount    int64 `form:"seize_count" default:"10"`   // 抢占多少个
	DispatchCount int64 `form:"dispatch_count" default:"1"` // 领取多少个
}

// ListOptions options for List
type ListOptions struct {
	common.BaseOptions
	common.Pager
	TempOptions
	BisShow bool // 用于列表展示还是直接派发
	State   int8 `form:"state"`
}

// SubmitOptions options for Submit
type SubmitOptions struct {
	common.BaseOptions
	TempOptions
	TaskID   int64 `form:"task_id"`
	Utime    uint64
	OldUID   int64
	OldState int8
}

// DelayOptions options for Delay
type DelayOptions struct {
	common.BaseOptions
	TaskID int64  `form:"task_id"`
	Reason string `form:"reason"`
}

// ConfigOption .
type ConfigOption struct {
	common.BaseOptions
	ID          int64  `form:"id"`
	Btime       string `form:"btime"`
	Etime       string `form:"etime"`
	Description string `form:"description"`
	ConfType    int8   `form:"conf_type" validate:"required"`
	ConfJSON    string `form:"conf_json" validate:"required"`
}

// Config .
type Config struct {
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
func (c Config) TableName() string {
	return "task_config"
}

// EqualWeightConfig 等值权重
type EqualWeightConfig struct {
	Name   string `json:"name"` // taskid 或者 mid
	IDs    string `json:"ids"`
	Weight int64  `json:"weight"`
	Type   int8   `json:"type"` // 周期或者定值
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
	MIDs []int64 `json:"mids"`
	UIDs []int64 `json:"uids"`
}

// QueryParams 配置筛选参数
type QueryParams struct {
	common.Pager
	ConfType   int8   `form:"conf_type"`
	State      int8   `form:"state"`
	BusinessID int64  `form:"business_id"`
	FlowID     int64  `form:"flow_id"`
	Btime      string `form:"mtime_from"`
	Etime      string `form:"mtime_to"`

	ConfName   string `form:"conf_name"`   // 筛选配置具体类型，fans,group,waittime,mid,taskid
	IDFilter   string `form:"id_filter"`   // 筛选具体的ID
	TypeFilter string `form:"type_filter"` // 筛选动态或静态权重
}

// History 任务日志
type History struct {
	TaskID  int64 `json:"task_id"`
	AdminID int64
	UID     int64
	Reason  string
	Uname   string
	Action  int8
}

// UnDOStat undo stat
type UnDOStat struct {
	Assign int64 `json:"assign_count" gorm:"column:assign"`
	Delay  int64 `json:"delay_count" gorm:"column:delay"`
	Normal int64 `json:"normal_count" gorm:"column:normal"`
}

// Stat 列表页最上方
type Stat struct {
	Normal int64 `json:"normal_count"  gorm:"column:normal"`
	Assign int64 `json:"assign_count"  gorm:"column:assign"`

	DelayTotal    int64 `json:"delay_total"  gorm:"column:delayTotal"`
	DelayPersonal int64 `json:"delay_personal"  gorm:"column:delayPersonal"`

	ReviewTotal    int64 `json:"review_total"`
	ReviewPersonal int64 `json:"review_personal"`
}

// RangeFunc .
type RangeFunc func(context.Context, *ListOptions) (map[int64]*Task, int64, []int64, []int64, error)

// RemoveFunc .
type RemoveFunc func(context.Context, *common.BaseOptions, ...interface{}) error

// ListFuncDB .
type ListFuncDB func(context.Context, map[int64]*Task, []int64, ...interface{}) (map[int64]struct{}, error)
