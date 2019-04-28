package model

import (
	"hash/crc32"
	"strconv"

	"go-common/library/time"
)

// All const variable used in dm2
const (
	AttrNo  = int32(0) // 属性位为0
	AttrYes = int32(1) // 属性位为1

	AttrProtect = uint(0) // 保护弹幕

	StateNormal        = int32(0)  // 普通状态
	StateDelete        = int32(1)  // 删除状态
	StateHide          = int32(2)  // 隐藏状态
	StateBlock         = int32(3)  // 屏蔽状态
	StateFilter        = int32(4)  // 过滤状态
	StateMonitorBefore = int32(5)  // 先审后发
	StateMonitorAfter  = int32(6)  // 先发后审
	StateSystemFilter  = int32(7)  // 敏感词过滤
	StateReportDelete  = int32(8)  // 举报删除
	StateAdminDelete   = int32(9)  // 弹幕管理删除
	StateUserDelete    = int32(10) // 用户删除
	StateScriptDelete  = int32(11) // 举报脚本删除

	PoolNormal   = int32(0) // 普通弹幕池
	PoolSubtitle = int32(1) // 字幕弹幕池
	PoolSpecial  = int32(2) // 特殊弹幕池

	NotFound = -1
)

// Hash 用户匿名弹幕uid hash
func Hash(uid int64, ip uint32) string {
	var s uint32
	if uid != 0 {
		s = crc32.ChecksumIEEE([]byte(strconv.FormatInt(uid, 10)))
		return strconv.FormatInt(int64(s), 16)
	}
	s = crc32.ChecksumIEEE([]byte(strconv.FormatInt(int64(ip), 10)))
	return "D" + strconv.FormatInt(int64(s), 16)
}

// DM dm_index and dm_content
type DM struct {
	ID         int64           `json:"id"`
	Type       int32           `json:"type"`
	Oid        int64           `json:"oid"`
	Mid        int64           `json:"mid"`
	Progress   int32           `json:"progress"`
	Pool       int32           `json:"pool"`
	Attr       int32           `json:"attr"`
	State      int32           `json:"state"`
	Ctime      time.Time       `json:"ctime"`
	Mtime      time.Time       `json:"mtime"`
	Content    *Content        `json:"content,omitempty"`
	ContentSpe *ContentSpecial `json:"content_special,omitempty"`
}

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

// Content dm_content
type Content struct {
	ID       int64     `json:"id"`
	FontSize int32     `json:"fontsize"`
	Color    int64     `json:"color"`
	Mode     int32     `json:"mode"`
	IP       int64     `json:"ip"`
	Plat     int32     `json:"plat"`
	Msg      string    `json:"msg"`
	Ctime    time.Time `json:"ctime"`
	Mtime    time.Time `json:"mtime"`
}

// ContentSpecial dm_content_special
type ContentSpecial struct {
	ID    int64     `json:"id"`
	Msg   string    `json:"msg"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}
