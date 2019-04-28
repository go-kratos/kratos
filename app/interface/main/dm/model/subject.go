package model

import (
	"go-common/library/time"
)

// All const variable used in dm subject
const (
	SubTypeVideo = int32(1) // 主题类型

	SubStateOpen   = int32(0) // 主题打开
	SubStateClosed = int32(1) // 主题关闭

	AttrSubGuest   = uint(0) // 允许游客弹幕
	AttrSubSpolier = uint(1) // 允许剧透弹幕
	AttrSubMission = uint(2) // 允许活动弹幕
	AttrSubAdvance = uint(3) // 允许高级弹幕

	// StateProtect protected dm
	StateProtect = 2
)

// Subject dm_subject
type Subject struct {
	ID        int64     `json:"id"`
	Type      int32     `json:"type"`
	Oid       int64     `json:"oid"`
	Pid       int64     `json:"pid"`
	Mid       int64     `json:"mid"`
	State     int32     `json:"state"`
	Attr      int32     `json:"attr"`
	ACount    int64     `json:"acount"`
	Count     int64     `json:"count"`
	MCount    int64     `json:"mcount"`
	MoveCnt   int64     `json:"move_count"`
	Maxlimit  int64     `json:"maxlimit"`
	Childpool int32     `json:"childpool"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// AttrVal return val of subject'attr
func (s *Subject) AttrVal(bit uint) int32 {
	return (s.Attr >> bit) & int32(1)
}

// AttrSet set val of subject'attr
func (s *Subject) AttrSet(v int32, bit uint) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}
