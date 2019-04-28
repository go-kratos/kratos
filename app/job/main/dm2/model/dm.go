package model

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strconv"
)

var (
	bAmp = []byte(`&amp;`)
	bGt  = []byte(`&gt;`)
	bLt  = []byte(`&lt;`)
	bSp  = []byte(` `)

	// <d p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id">弹幕内容</d>
	_xmlFmt = `<d p="%.5f,%d,%d,%d,%d,%d,%s,%d">%s</d>`
	// <d p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id，用户id">弹幕内容</d>
	_rnameFmt = `<d p="%.5f,%d,%d,%d,%d,%d,%s,%d,%d">%s</d>`
)

// All const variable use in job
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
	StateTaskDel       = int32(12) //弹幕任务删除

	// 弹幕模式
	ModeRolling = int32(1)
	ModeBottom  = int32(4)
	ModeTop     = int32(5)
	ModeReverse = int32(6)
	ModeSpecial = int32(7)
	ModeCode    = int32(8)
	ModeBAS     = int32(9)

	PoolNormal   = int32(0) // 普通弹幕池
	PoolSubtitle = int32(1) // 字幕弹幕池
	PoolSpecial  = int32(2) // 特殊弹幕池

	MaskPriorityHgih = int32(1) // 弹幕蒙版优先级高
	MaskPriorityLow  = int32(0) // 弹幕蒙版优先级低

	NotFound = int64(-1)
)

// BinlogMsg binlog msg produced by canal
type BinlogMsg struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// AttrVal return val of index'attr
func (d *DM) AttrVal(bit uint) int32 {
	return (d.Attr >> bit) & int32(1)
}

// AttrSet set val of index'attr
func (d *DM) AttrSet(v int32, bit uint) {
	d.Attr = d.Attr&(^(1 << bit)) | (v << bit)
}

// DMSlice dm array
type DMSlice []*DM

func (d DMSlice) Len() int           { return len(d) }
func (d DMSlice) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DMSlice) Less(i, j int) bool { return d[i].ID < d[j].ID }

// ToXML convert dm struct to xml.
func (d *DM) ToXML(realname bool) (s string) {
	if d.Content == nil {
		return
	}
	msg := d.Content.Msg
	if d.ContentSpe != nil {
		msg = d.ContentSpe.Msg
	}
	if len(msg) == 0 {
		return
	}
	if realname {
		// <d e="用户id" p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id">弹幕内容</d>
		s = fmt.Sprintf(_rnameFmt, float64(d.Progress)/1000.0, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, hash(d.Mid, uint32(d.Content.IP)), d.ID, d.Mid, xmlReplace([]byte(msg)))
	} else {
		// <d p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id">弹幕内容</d>
		s = fmt.Sprintf(_xmlFmt, float64(d.Progress)/1000.0, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, hash(d.Mid, uint32(d.Content.IP)), d.ID, xmlReplace([]byte(msg)))
	}
	return
}

// xmlReplace replace special char in xml.
func xmlReplace(bi []byte) (bo []byte) {
	for _, b := range bi {
		if b == 0 {
			continue
		} else if b == '&' {
			bo = append(bo, bAmp...)
			continue
		} else if b == '>' {
			bo = append(bo, bGt...)
			continue
		} else if b == '<' {
			bo = append(bo, bLt...)
			continue
		} else if (b >= 0x01 && b <= 0x08) || (b >= 0x0b && b <= 0x0c) || (b >= 0x0e && b <= 0x1f) || (b == 0x7f) {
			bo = append(bo, bSp...)
		} else {
			bo = append(bo, b)
		}
	}
	return
}

// hash return hash string.
func hash(mid int64, ip uint32) string {
	var s uint32
	if mid != 0 {
		s = crc32.ChecksumIEEE([]byte(strconv.FormatInt(mid, 10)))
		return strconv.FormatInt(int64(s), 16)
	}
	s = crc32.ChecksumIEEE([]byte(strconv.FormatInt(int64(ip), 10)))
	return "D" + strconv.FormatInt(int64(s), 16)
}

// GetSpecialSeg .
func (d *DM) GetSpecialSeg() (msg string) {
	if d.Content == nil || d.Pool != PoolSpecial {
		return
	}
	msg = d.Content.Msg
	if d.ContentSpe != nil {
		msg = d.ContentSpe.Msg
	}
	return
}

// NeedDisplay 判断该条弹幕是否需要展示
func (d *DM) NeedDisplay() bool {
	return d.State == StateNormal || d.State == StateMonitorAfter
}

// NeedUpdateSpecial .
func (d *DM) NeedUpdateSpecial(old *DM) bool {
	if (d.Pool == PoolSpecial || old.Pool == PoolSpecial) && d.Pool != old.Pool {
		return true
	}
	if d.Pool == PoolSpecial && d.NeedDisplay() && !old.NeedDisplay() {
		return true
	}
	if d.Pool == PoolSpecial && old.NeedDisplay() && !d.NeedDisplay() {
		return true
	}
	return false
}
