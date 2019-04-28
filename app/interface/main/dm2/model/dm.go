package model

import (
	"fmt"
	"hash/crc32"
	"strconv"

	"go-common/library/time"
	"go-common/library/xstr"
)

var (
	bAmp = []byte(`&amp;`)
	bGt  = []byte(`&gt;`)
	bLt  = []byte(`&lt;`)
	bSp  = []byte(` `)
	// LimitPerMin 每个rank权限允许的发送速度
	LimitPerMin = map[int32]int64{
		0:     10,
		10000: 30,
		15000: 30,
		20000: 300,
		25000: 60,
		30000: 300,
		32000: 300,
	}
)

// All const variable used in dm2
const (
	// <d p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id">弹幕内容</d>
	_xmlFmt = `<d p="%.5f,%d,%d,%d,%d,%d,%s,%d">%s</d>`
	// <d p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id，用户id">弹幕内容</d>
	_rnameFmt = `<d p="%.5f,%d,%d,%d,%d,%d,%s,%d,%d">%s</d>`

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
	StateAiDelete      = int32(13) // ai删除

	PoolNormal   = int32(0) // 普通弹幕池
	PoolSubtitle = int32(1) // 字幕弹幕池
	PoolSpecial  = int32(2) // 特殊弹幕池

	PlatUnknow  = int32(0)
	PlatWeb     = int32(1)
	PlatAndroid = int32(2)
	PlatIPhone  = int32(3)
	PlatWpM     = int32(4) // wp mobile
	PlatIPad    = int32(5)
	PlatPadHd   = int32(6) // ipad hd
	PlatWpPc    = int32(7) // win10

	MaxLenDefMsg = 100
	MaxLen7Msg   = 300
	// 弹幕模式
	ModeRolling = int32(1)
	ModeBottom  = int32(4)
	ModeTop     = int32(5)
	ModeReverse = int32(6)
	ModeSpecial = int32(7)
	ModeCode    = int32(8)
	ModeBAS     = int32(9)

	SpamBlack    = 52001
	SpamOverflow = 52002
	SpamRestrict = 52005
)

// AdvanceCmt struct
type AdvanceCmt struct {
	ID        int64
	Owner     int64
	Oid       int64
	Type      string // request:申请中，buy:已购买，accept:已通过，deny:拒绝
	Mode      string // 枚举类型，sp:高级弹幕，advance:pool2特殊弹幕
	Timestamp int64
	Mid       int64
	Refund    int
}

// DailyLimiter daily limiter
type DailyLimiter struct {
	Date  string
	Count int64
}

// Limiter retrun
type Limiter struct {
	Allowance int64
	Timestamp int64
}

// AttrVal return val of index'attr
func (d *DM) AttrVal(bit uint) int32 {
	return (d.Attr >> bit) & int32(1)
}

// AttrSet set val of index'attr
func (d *DM) AttrSet(v int32, bit uint) {
	d.Attr = d.Attr&(^(1 << bit)) | (v << bit)
}

// DMSlice dm slice
type DMSlice []*DM

func (d DMSlice) Len() int           { return len(d) }
func (d DMSlice) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DMSlice) Less(i, j int) bool { return d[i].ID < d[j].ID }

// DMSlice2 sort dm slice by progress
type DMSlice2 []*DM

func (d DMSlice2) Len() int           { return len(d) }
func (d DMSlice2) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DMSlice2) Less(i, j int) bool { return d[i].Progress < d[j].Progress }

// JudgeDMList dm list of judge
type JudgeDMList struct {
	List  []*JDM  `json:"list"`
	Index []int64 `json:"index"`
}

// JDM judge dm
type JDM struct {
	ID       int64     `json:"id"`
	Msg      string    `json:"msg"`
	Mid      int64     `json:"mid"`
	Progress string    `json:"progress"`
	CTime    time.Time `json:"ctime"`
}

// FilterData filter-service data
type FilterData struct {
	Level  int64    `json:"level"`
	Limit  int64    `json:"limit"`
	Msg    string   `json:"msg"`
	TypeID []int64  `json:"typeid"`
	Hit    []string `json:"hit"`
}

// ToXMLSpecialSeg .
func (d *DM) ToXMLSpecialSeg(realname bool) (s string) {
	if d.Content == nil || d.Pool != PoolSpecial {
		return
	}
	msg := d.Content.Msg
	if d.ContentSpe != nil {
		msg = d.ContentSpe.Msg
	}
	if len(msg) == 0 {
		return
	}
	s = fmt.Sprintf(`<d id="%d">%s</d>`, d.ID, xmlReplace([]byte(msg)))
	return
}

// ToXMLSeg convert dm struct to xml.
func (d *DM) ToXMLSeg(realname bool) (s string) {
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
	if d.Pool == PoolSpecial {
		msg = ""
	}
	if realname {
		// <d p="弹幕ID,弹幕属性,播放时间,弹幕模式,字体大小,颜色,发送时间,弹幕池,用户hash id,用户mid">弹幕内容</d>
		s = fmt.Sprintf(_xmlSegRealnameFmt, d.ID, d.Attr, d.Progress, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, Hash(d.Mid, uint32(d.Content.IP)), d.Mid, xmlReplace([]byte(msg)))

	} else {
		// <d p="弹幕ID,弹幕属性,播放时间,弹幕模式,字体大小,颜色,发送时间,弹幕池,用户id">弹幕内容</d>
		s = fmt.Sprintf(_xmlSegFmt, d.ID, d.Attr, d.Progress, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, Hash(d.Mid, uint32(d.Content.IP)), xmlReplace([]byte(msg)))
	}
	return
}

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
		// <d p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id，用户id">弹幕内容</d>
		s = fmt.Sprintf(_rnameFmt, float64(d.Progress)/1000.0, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, Hash(d.Mid, uint32(d.Content.IP)), d.ID, d.Mid, xmlReplace([]byte(msg)))
	} else {
		// <d p="播放时间，弹幕模式，字体大小，颜色，发送时间，弹幕池，用户hash，弹幕id">弹幕内容</d>
		s = fmt.Sprintf(_xmlFmt, float64(d.Progress)/1000.0, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, Hash(d.Mid, uint32(d.Content.IP)), d.ID, xmlReplace([]byte(msg)))
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
			// 替换掉控制字符，保留空格、回车、制表符
			bo = append(bo, bSp...)
		} else {
			bo = append(bo, b)
		}
	}
	return
}

// Hash return mid hash string.
func Hash(mid int64, ip uint32) string {
	var s uint32
	if mid != 0 {
		s = crc32.ChecksumIEEE([]byte(strconv.FormatInt(mid, 10)))
		return strconv.FormatInt(int64(s), 16)
	}
	s = crc32.ChecksumIEEE([]byte(strconv.FormatInt(int64(ip), 10)))
	return "D" + strconv.FormatInt(int64(s), 16)
}

// IsDMVisible dm can visible
func IsDMVisible(state int32) bool {
	if state == StateNormal || state == StateMonitorAfter {
		return true
	}
	return false
}

// IsDMEditAble dm can edit
func IsDMEditAble(state int32) bool {
	if state == StateNormal || state == StateHide || state == StateMonitorAfter {
		return true
	}
	return false
}

// AttrNtoA convert uint to string format,eg:5-->101-->1,3.
func (d *DM) AttrNtoA() string {
	if d.Attr == 0 {
		return ""
	}
	var bits []int64
	for k, v := range strconv.FormatInt(int64(d.Attr), 2) {
		if v == 49 {
			bits = append(bits, int64(k+1))
		}
	}
	return xstr.JoinInts(bits)
}

// DMAttrNtoA convert uint to string format,eg:5-->101-->1,3.
func DMAttrNtoA(attr int32) string {
	if attr == 0 {
		return ""
	}
	var bits []int64
	for k, v := range strconv.FormatInt(int64(attr), 2) {
		if v == 49 {
			bits = append(bits, int64(k+1))
		}
	}
	return xstr.JoinInts(bits)
}
