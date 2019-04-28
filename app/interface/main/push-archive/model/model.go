package model

import (
	"encoding/json"
	"time"

	"go-common/app/service/main/archive/api"
	relmdl "go-common/app/service/main/relation/model"
)

const (
	// PushTypeUnknown 用户未上报推送设置
	PushTypeUnknown = iota
	// PushTypeForbid 禁止推送稿件更新通知
	PushTypeForbid
	// PushTypeSpecial  推送特别关注的upper的更新
	PushTypeSpecial
	// PushTypeAttention 推送关注的upper的更新
	PushTypeAttention
)

const (
	// RelationAttention 关注
	RelationAttention = iota + 1
	// RelationSpecial 特别关注
	RelationSpecial
)

const (
	// StatisticsUnpush 命中分组但未推送
	StatisticsUnpush = iota
	// StatisticsPush 命中分组且推送
	StatisticsPush = 1
)

const (
	// GroupDataTypeDefault 默认
	GroupDataTypeDefault = "default"
	// GroupDataTypeHBase AI脚本提供的hbase数据
	GroupDataTypeHBase = "hbase"
	// GroupDataTypeAbtest ab实验数据
	GroupDataTypeAbtest = "ab_test"
	// GroupDataTypeAbComparison ab对照数据
	GroupDataTypeAbComparison = "ab_comparison"
)

const (
	// AttrBitIsPGC pgc稿件的属性位
	AttrBitIsPGC = 9
)

// Setting user push setting.
type Setting struct {
	Type int `json:"type"`
}

// Message canal databus message.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// Relation user relation.
type Relation struct {
	Mid       int64  `json:"mid,omitempty"`
	Fid       int64  `json:"fid,omitempty"`
	Attribute uint32 `json:"attribute"`
	Status    int    `json:"status"`
	MTime     string `json:"mtime"`
	CTime     string `json:"ctime"`
}

// Following judge that whether has following relation.
func (r *Relation) Following() bool {
	attr := relmdl.Following{Attribute: r.Attribute}
	return attr.Following()
}

// RelationTagUser user relatino tag.
type RelationTagUser struct {
	Mid   int64  `json:"mid,omitempty"`
	Fid   int64  `json:"fid,omitempty"`
	Tag   string `json:"tag"`
	MTime string `json:"mtime"`
	CTime string `json:"ctime"`
}

// HasTag judge that whether has specified tag.
func (r *RelationTagUser) HasTag(tag int64) bool {
	i := new(Ints)
	i.Scan([]byte(r.Tag))
	return i.Exist(tag)
}

// Archive model
type Archive struct {
	ID        int64  `json:"aid"`
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"typeid"`
	HumanRank int    `json:"humanrank"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Tag       string `json:"tag"`
	Attribute int32  `json:"attribute"`
	Copyright int8   `json:"copyright"`
	AreaLimit int8   `json:"arealimit"`
	State     int    `json:"state"`
	Author    string `json:"author"`
	Access    int    `json:"access"`
	Forward   int    `json:"forward"`
	PubTime   string `json:"pubtime"`
	Round     int8   `json:"round"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
}

// IsNormal judge that whether archive's state is normally.
func (a *Archive) IsNormal() bool {
	arc := api.Arc{State: int32(a.State)}
	return arc.IsNormal()
}

// PushStatistic 推送统计数据对象
type PushStatistic struct {
	Aid         int64     `json:"aid"`
	Group       string    `json:"group"`
	Type        int       `json:"type"`
	Mids        string    `json:"mids"`
	MidsCounter int       `json:"mids_counter"`
	CTime       time.Time `json:"ctime"`
}
