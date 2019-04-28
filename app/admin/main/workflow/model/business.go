package model

import xtime "go-common/library/time"

const (
	// 业务bid 对应manager项目子业务
	// ArchiveComplain 稿件投诉
	ArchiveComplain = 1
	// ArchiveAppeal 稿件申诉
	ArchiveAppeal = 2
	// ReviewShortComplain 短点评投诉
	ReviewShortComplain = 3
	// ReviewLongComplain 长点评投诉
	ReviewLongComplain = 4
	// CreditAppeal 小黑屋申诉
	CreditAppeal = 5
	// ArchiveAudit 稿件审核
	ArchiveAudit = 6
	// ArchiveVT 任务质检
	ArchiveVT = 7
	// ChannelComplain 频道举报
	ChannelComplain = 9
	// CommentComplain 评论举报
	CommentComplain = 13
	// SubtitleComplain 字幕举报
	SubtitleComplain = 14
)

// Business will record any business properties
type Business struct {
	Bid      int32      `json:"-" gorm:"column:id"`
	Gid      int64      `json:"gid" gorm:"column:gid"`
	Cid      int64      `json:"cid" gorm:"column:cid"`
	Oid      int64      `json:"oid" gorm:"column:oid"`
	Business int8       `json:"business" gorm:"column:business"`
	TypeID   int32      `json:"typeid" gorm:"column:typeid"`
	Title    string     `json:"title" gorm:"column:title"`
	Content  string     `json:"content" gorm:"column:content"`
	Mid      int64      `json:"mid" gorm:"column:mid"`
	Extra    string     `json:"extra" gorm:"column:extra"`
	CTime    xtime.Time `json:"-" gorm:"column:ctime"`
	MTime    xtime.Time `json:"-" gorm:"column:mtime"`
}

// Meta is the model to store business metadata
type Meta struct {
	Business int8          `json:"business"`
	Name     string        `json:"name"`
	ItemType string        `json:"item_type"`
	Rounds   []*Round      `json:"rounds"`
	Attr     *BusinessAttr `json:"attr"`
}

// MetaSlice is used to support sort Metas
type MetaSlice []*Meta

func (ms MetaSlice) Len() int {
	return len(ms)
}

func (ms MetaSlice) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms MetaSlice) Less(i, j int) bool {
	return ms[i].Business < ms[j].Business
}

// Round is the model to describe how many business rounds are
type Round struct {
	ID   int8   `json:"id"`
	Name string `json:"name"`
}

// RoundSlice is used to support sort Rounds
type RoundSlice []*Round

func (rs RoundSlice) Len() int {
	return len(rs)
}

func (rs RoundSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs RoundSlice) Less(i, j int) bool {
	return rs[i].ID < rs[j].ID
}
