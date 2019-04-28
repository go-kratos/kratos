package model

import (
	"sort"
	"strconv"
	"time"

	xtime "go-common/library/time"
)

// RelationLog is
type RelationLog struct {
	Mid           int64  `json:"mid"`
	Fid           int64  `json:"fid"`
	MemberName    string `json:"member_name"`
	FollowingName string `json:"following_name"`

	Source int32      `json:"source"`
	MTime  xtime.Time `json:"mtime"`

	Attention  int32  `json:"attention"`
	Black      int32  `json:"black"`
	Whisper    int32  `json:"whisper"`
	AttrField  string `json:"attr_field"`
	AttrChange string `json:"attr_change"`
}

// FillAttrField is
func (l *RelationLog) FillAttrField() {
	if l.Attention > 0 {
		l.AttrField = "attention"
		return
	}
	if l.Black > 0 {
		l.AttrField = "black"
		return
	}
	if l.Whisper > 0 {
		l.AttrField = "whisper"
		return
	}
}

// RelationLogList is
type RelationLogList []*RelationLog

// Len is
func (rl RelationLogList) Len() int {
	return len(rl)
}

// Swap is
func (rl RelationLogList) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}

// Less is
func (rl RelationLogList) Less(i, j int) bool {
	return rl[i].MTime < rl[j].MTime
}

// OrderByMTime is
func (rl RelationLogList) OrderByMTime(desc bool) {
	sort.Sort(rl)
}

// ParseAction is
func ParseAction(act string) int32 {
	i, _ := strconv.ParseInt(act, 10, 32)
	return int32(i)
}

// ParseSource is
func ParseSource(src string) int32 {
	i, _ := strconv.ParseInt(src, 10, 64)
	return int32(i)
}

// ParseLogTime is
func ParseLogTime(ts string) (xt xtime.Time, err error) {
	var (
		t time.Time
	)
	if t, err = time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local); err != nil {
		return
	}
	xt.Scan(t)
	return
}
