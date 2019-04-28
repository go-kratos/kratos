package model

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// All const variable used in job
const (
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
	AttrSubMaskReady     = uint(7) // 蒙版生产完成

	MaskPlatWeb int8 = 0   // web端
	MaskPlatMbl int8 = 1   // 移动端
	MaskPlatAll int8 = 100 //全部端
)

// Subject dm_subject.
type Subject struct {
	ID        int64 `json:"id"`
	Type      int32 `json:"type"`
	Oid       int64 `json:"oid"`
	Pid       int64 `json:"pid"`
	Mid       int64 `json:"mid"`
	State     int32 `json:"state"`
	Attr      int32 `json:"attr"`
	ACount    int64 `json:"acount"`
	Count     int64 `json:"count"`
	MCount    int64 `json:"mcount"`
	MoveCnt   int64 `json:"move_count"`
	Maxlimit  int64 `json:"maxlimit"`
	Childpool int32 `json:"childpool"`
	Ctime     stime `json:"ctime"`
	Mtime     stime `json:"mtime"`
}

// ConvertStime .
func ConvertStime(t time.Time) stime {
	return stime(t.Unix())
}

type stime int64

// Scan scan time.
func (st *stime) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*st = stime(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*st = stime(i)
	}
	return
}

// Value get time value.
func (st stime) Value() (driver.Value, error) {
	return time.Unix(int64(st), 0), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (st *stime) UnmarshalJSON(data []byte) error {
	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		*st = stime(timestamp)
		return nil
	}
	t, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*st = stime(t.Unix())
	return err
}

// AttrVal return val of subject'attr
func (s *Subject) AttrVal(bit uint) int32 {
	return (s.Attr >> bit) & int32(1)
}

// AttrSet set val of subject'attr
func (s *Subject) AttrSet(v int32, bit uint) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}
