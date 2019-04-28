package model

import (
	"go-common/library/time"
)

// All const variable used in dm subject
const (
	AttrNo  = int32(0) // no
	AttrYes = int32(1) // yes

	SubTypeVideo = int32(1) // 主题类型

	SubStateOpen   = int32(0) // 主题打开
	SubStateClosed = int32(1) // 主题关闭

	AttrSubGuest         = uint(0) // 允许游客弹幕
	AttrSubSpolier       = uint(1) // 允许剧透弹幕
	AttrSubMission       = uint(2) // 允许活动弹幕
	AttrSubAdvance       = uint(3) // 允许高级弹幕
	AttrSubMonitorBefore = uint(4) // 弹幕先审后发
	AttrSubMonitorAfter  = uint(5) // 弹幕先发后审
	AttrSubMaskOpen      = uint(6) // 开启蒙版
	AttrSubMblMaskReady  = uint(7) // 移动端蒙版生产完成
	AttrSubWebMaskReady  = uint(8) // web端蒙版生产完成
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

// SubjectLog subject log
type SubjectLog struct {
	UID     int64  `json:"uid"`
	Uname   string `json:"uname"`
	Oid     int64  `json:"oid"`
	Action  string `json:"action"`
	Comment string `json:"comment"`
	Ctime   string `json:"ctime"`
}

// SeasonInfo season info.
type SeasonInfo struct {
	Aid    int64  `json:"aid"`
	Cid    int64  `json:"cid"`
	Epid   int64  `json:"ep_id"`
	Ssid   int64  `json:"season_id"`
	State  int64  `json:"is_delete"`
	LTitle string `json:"long_title"`
	Title  string `json:"title"`
}

// SearchSubjectReq search subject request.
type SearchSubjectReq struct {
	Oids, Aids, Mids, Attrs []int64
	State                   int64
	Pn, Ps                  int64
	Sort, Order             string
}

// SearchSubjectResult result from search
type SearchSubjectResult struct {
	Page   *Page
	Result []*struct {
		Oid int64 `json:"oid"`
	} `json:"result"`
}

// SearchSubjectLog get subject logs
type SearchSubjectLog struct {
	Page   *Page
	Result []*struct {
		UID       int64  `json:"uid"`
		Uname     string `json:"uname"`
		Oid       int64  `json:"oid"`
		Action    string `json:"action"`
		ExtraData string `json:"extra_data"`
		Ctime     string `json:"ctime"`
	}
}

// AttrVal return val of subject'attr
func (s *Subject) AttrVal(bit uint) int32 {
	return (s.Attr >> bit) & int32(1)
}

// AttrSet set val of subject'attr
func (s *Subject) AttrSet(v int32, bit uint) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}

// IsMonitoring check if the subject is monitoring or not.
func (s *Subject) IsMonitoring() bool {
	return s.AttrVal(AttrSubMonitorBefore) == AttrYes ||
		s.AttrVal(AttrSubMonitorAfter) == AttrYes
}
