package model

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// all const variable used in job
const (
	AttrNo  = int32(0)
	AttrYes = int32(1)

	// platform
	PlatUnknow  = int32(0)
	PlatWeb     = int32(1)
	PlatAndroid = int32(2)
	PlatIPhone  = int32(3)
	PlatWpM     = int32(4) // wp mobile
	PlatIPad    = int32(5)
	PlatPadHd   = int32(6) // ipad hd
	PlatWpPc    = int32(7) // win10

	// dm state
	StateNormal        = int32(0) // 普通状态
	StateDelete        = int32(1) // 删除状态
	StateHide          = int32(2) // 隐藏状态
	StateBlock         = int32(3) // 屏蔽状态
	StateFilter        = int32(4) // 过滤状态
	StateMonitorBefore = int32(5) // 先审后发
	StateMonitorAfter  = int32(6) // 先发后审

	// dm attribute
	AttrProtect = uint(0) // 保护弹幕

	// dm pool
	PoolNormal   = int32(0) // 普通弹幕池
	PoolSubtitle = int32(1) // 字幕弹幕池
	PoolSpecial  = int32(2) // 特殊弹幕池

	// dm mode
	ModeNormal  = int32(1) // 正常滚动弹幕
	ModeBottom  = int32(4) // 底部弹幕
	ModeTop     = int32(5) // 顶部弹幕
	ModeReverse = int32(6) // 逆向滚动弹幕
	ModeAdvance = int32(7) // 高级弹幕
	ModeCode    = int32(8) // 代码弹幕

	NotFound = int64(-1)
)

// AttrVal return val of index'attr
func (d *DM) AttrVal(bit uint) int32 {
	return (d.Attr >> bit) & int32(1)
}

// AttrSet set val of index'attr
func (d *DM) AttrSet(v int32, bit uint) {
	d.Attr = d.Attr&(^(1 << bit)) | (v << bit)
}

// NeedDisplay 判断该条弹幕是否需要展示
func (d *DM) NeedDisplay() bool {
	return d.State == StateNormal || d.State == StateMonitorAfter
}

// NeedStateNormal 判断是否更新状态
// pool0 变为 pool1 状态正常
// 变为保护弹幕 状态正常
func (d *DM) NeedStateNormal(old *DM) bool {
	if (d.Pool != old.Pool) && d.Pool == PoolSubtitle {
		return !(d.State == StateNormal)
	}
	if d.AttrVal(AttrProtect) == AttrYes && old.AttrVal(AttrProtect) == AttrNo {
		return !(d.State == StateNormal)
	}
	return false
}

// Trim dmid and it's progress time will be trimed.
type Trim struct {
	ID   int64 `json:"id"`
	Attr int32 `json:"-"`
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
