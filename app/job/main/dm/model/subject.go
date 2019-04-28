package model

import (
	"go-common/library/time"
)

var (
	_defaultSeg = &Segment{Start: 0, End: DefaultVideoEnd, Cnt: 1, Num: 1, Duration: 0}
)

const (
	// segmentLength 分段长度，根据视频时长做分段，单位：毫秒
	segmentLength = int64(6 * 60 * 1000)
	// DefaultVideoEnd 当视频时长不存在或者为0时的默认视频结尾时间点
	DefaultVideoEnd = int64(10 * 60 * 60 * 1000)

	// SubTypeVideo 主题类型
	SubTypeVideo = int32(1)
	// SubStateOpen 主题打开
	SubStateOpen = int32(0)
	// SubStateClosed 主题关闭
	SubStateClosed = int32(1)

	// AttrSubGuest 允许游客弹幕
	AttrSubGuest = uint(0)
	// AttrSubSpolier 允许剧透弹幕
	AttrSubSpolier = uint(1)
	// AttrSubMission 允许活动弹幕
	AttrSubMission = uint(2)
	// AttrSubAdvance 允许高级弹幕
	AttrSubAdvance = uint(3)
	// AttrSubMonitorBefore 弹幕先审后发
	AttrSubMonitorBefore = uint(4)
	// AttrSubMonitorAfter 弹幕先发后审
	AttrSubMonitorAfter = uint(5)
)

// Subject dm_subject.
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
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}

// AttrVal return val of subject'attr.
func (s *Subject) AttrVal(bit uint) int32 {
	return (s.Attr >> bit) & int32(1)
}

// AttrSet set val of subject'attr.
func (s *Subject) AttrSet(v int32, bit uint) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}

// Segment dm segment struct
type Segment struct {
	Start    int64 `json:"ps"`       // 分段起始时间
	End      int64 `json:"pe"`       // 分段结束时间
	Cnt      int64 `json:"cnt"`      // 总分段数
	Num      int64 `json:"num"`      // 当前第几段
	Duration int64 `json:"duration"` // 视频总时长
}

// SegmentInfo get segment info by start time and video duration.
func SegmentInfo(ps, duration int64) (s *Segment) {
	var cnt, num, pe int64
	if duration == 0 {
		s = _defaultSeg
		return
	}
	cnt = duration / segmentLength
	if duration%segmentLength > 0 {
		cnt++
	}
	for i := int64(0); i < cnt; i++ {
		if ps >= i*segmentLength && ps < (i+1)*segmentLength {
			ps = i * segmentLength
			pe = (i + 1) * segmentLength
			num = i + 1
		}
	}
	if pe > duration {
		pe = duration
	}
	if ps > duration {
		ps = duration
		pe = duration
		num = cnt
	}
	s = &Segment{
		Start:    ps,
		End:      pe,
		Cnt:      cnt,
		Num:      num,
		Duration: duration,
	}
	return
}

// SegmentPoint 根据当前段数和视频总时长计算分段的起始时间点
func SegmentPoint(num, duration int64) (ps, pe int64) {
	if duration == 0 {
		ps = 0
		pe = DefaultVideoEnd
		return
	}
	pe = num * segmentLength
	ps = pe - segmentLength
	if pe > duration {
		pe = duration
	}
	if ps < 0 {
		ps = 0
	}
	return
}
