package model

import (
	"time"
)

// Tag is the universal tag model, contains any type of tags
// The Business field and the Round field will from any business definition
type Tag struct {
	Tid      int32      `gorm:"column:id" json:"tid"`
	Name     string     `gorm:"column:name" json:"name"`
	Business int8       `gorm:"column:business" json:"business"`
	Weight   int16      `gorm:"column:weight" json:"weight"`
	Round    int8       `gorm:"column:round" json:"round"`
	State    int8       `gorm:"column:state" json:"state"`
	Remark   string     `gorm:"column:remark" json:"remark"`
	Controls []*Control `gorm:"-" json:"controls"`
	CTime    time.Time  `gorm:"column:ctime" json:"ctime"`
	MTime    time.Time  `gorm:"column:mtime" json:"mtime"`
}

// TableName Tag tablename
func (*Tag) TableName() string {
	return "workflow_tag"
}

// TagsCache tag cache
type TagsCache struct {
	TagMap     map[int8]map[int32]*Tag
	TagSlice   map[int8][]*Tag
	TagMap3    map[int64]map[int64][]*Tag3
	TagMap3Tid map[int64]map[int64]*Tag3
}

// ResponseTag3 .
type ResponseTag3 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int32  `json:"ttl"`
	Data    *Data
}

// Data .
type Data struct {
	Data  []*Tag3 `json:"data"`
	Pager *Pager  `json:"page"`
}

// Tag3 .
type Tag3 struct {
	BID      int64       `json:"bid"`
	TagID    int64       `json:"tag_id"`
	RID      int64       `json:"rid"`
	Name     string      `json:"name"`
	Weight   int64       `json:"weight"`
	State    int64       `json:"state"`
	Desc     string      `json:"description"`
	Controls []*Control3 `json:"controls"`
	Ctime    int64       `json:"ctime"`
	Mtime    int64       `json:"mtime"`
}

// Pager .
type Pager struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}
