package oversea

import (
	"go-common/library/time"
)

const (
	StateOK      = 1
	StateDeleted = 0
)

type ArchiveRelation struct {
	ID      int64 `json:"id" gorm:"column:id"`
	GroupID int64 `json:"policy_id" gorm:"column:policy_id"`
	Aid     int64 `json:"aid" gorm:"column:aid"`
}

func (ArchiveRelation) TableName() string {
	return "archive_relation"
}

type PolicyGroup struct {
	ID       int64         `json:"id" gorm:"column:id"`
	IsGlobal int8          `json:"is_global" gorm:"column:is_global"`
	Aid      int64         `json:"aid" gorm:"column:aid"`
	Name     string        `json:"name" gorm:"column:name"`
	Type     int8          `json:"type" gorm:"column:type"`
	State    int8          `json:"state" gorm:"column:state"`
	UserName string        `json:"username" gorm:"-"`
	UID      int64         `json:"uid" gorm:"column:uid"`
	Remark   string        `json:"remark" gorm:"column:remark;size:140"`
	CTime    time.Time     `json:"ctime" gorm:"column:ctime"`
	MTime    time.Time     `json:"mtime" gorm:"column:mtime"`
	Items    []*PolicyItem `json:"items" gorm:"-"`
}

func (PolicyGroup) TableName() string {
	return "policy_group"
}

type PolicyGroupData struct {
	Items []*PolicyGroup `json:"items"`
	Pager *Pager         `json:"pager"`
}

type Pager struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

type PolicyItem struct {
	ID       int64  `json:"id" gorm:"column:id"`
	GroupID  int64  `json:"group_id" gorm:"column:group_id"`
	PlayAuth int8   `json:"play_auth" gorm:"column:play_auth"`
	DownAuth int8   `json:"down_auth" gorm:"column:down_auth"`
	AreaID   string `json:"area_id" gorm:"column:area_id"`
	ZoneID   string `json:"zone_id" gorm:"column:zone_id"`
	State    int8   `json:"state" gorm:"column:state"`
}

func (PolicyItem) TableName() string {
	return "policy_item"
}

type PolicyParams struct {
	ID       int64   `json:"id"`
	DownAuth int8    `json:"down_auth"`
	PlayAuth int8    `json:"play_auth"`
	AreaIds  []int64 `json:"area_id"`
}

type Zone struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Pid    int64  `json:"pid"`
	Ppid   int64  `json:"ppid"`
	ZoneID int64  `json:"zone_id"`
	Status int8   `json:"status"`
}

func (Zone) TableName() string {
	return "locality"
}
