package model

import (
	"strconv"

	"go-common/library/time"
	"go-common/library/xstr"
)

//CondIntNil cond int nil
const CondIntNil = -10516

// batch operation deleted code
const (
	StatusNormal  = iota // 正常弹幕
	StatusDelete         // 删除弹幕
	StatusProtect        // 保护弹幕

	DMIndexInactive = int8(0)
	DMIndexActive   = int8(1)

	PoolNormal   = int32(0) // 普通弹幕池
	PoolSubtitle = int32(1) // 字幕弹幕池
	PoolSpecial  = int32(2) // 特殊弹幕池

	AttrProtect = uint(0) // 保护弹幕

	StateNormal        = int32(0)  // 普通状态
	StateDelete        = int32(1)  // 删除状态
	StateHide          = int32(2)  // 隐藏状态
	StateBlock         = int32(3)  // 屏蔽状态
	StateFilter        = int32(4)  // 过滤状态
	StateMonitorBefore = int32(5)  // 先审后发
	StateMonitorAfter  = int32(6)  // 先发后审
	StateSensBlock     = int32(7)  // 敏感词过滤
	StateReportDelete  = int32(8)  // 举报删除
	StateAdminDelete   = int32(9)  // 后台管理删除
	StateUserDelete    = int32(10) // 用户删除
	StateRptAutoDelete = int32(11) // 举报脚本自动删除
	StateTaskDelete    = int32(12) // 弹幕任务删除
	StateAiDelete      = int32(13) // ai删除

	DMLogBizID = int(31) // dm日志平台business id

	// mask platform
	MaskPlatWeb int8 = 0
	MaskPlatMbl int8 = 1
	MaskPlatAll int8 = 100
)

// StateDesc get a state description
func StateDesc(state int32) (description string) {
	switch state {
	case StateNormal:
		description = "正常弹幕"
	case StateDelete:
		description = "删除状态"
	case StateHide:
		description = "隐藏状态"
	case StateBlock:
		description = "屏蔽状态"
	case StateFilter:
		description = "过滤状态"
	case StateMonitorBefore:
		description = "先审后发"
	case StateMonitorAfter:
		description = "先发后审"
	case StateSensBlock:
		description = "敏感词过滤"
	case StateReportDelete:
		description = "举报删除"
	case StateAdminDelete:
		description = "弹幕管理删除"
	case StateUserDelete:
		description = "用户删除"
	case StateRptAutoDelete:
		description = "举报脚本删除"
	case StateTaskDelete:
		description = "弹幕任务删除"
	default:
		description = "未知状态"
	}
	return
}

// DM dm info for new database
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

// Content dm content info
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

// ContentSpecial special dm data
type ContentSpecial struct {
	ID    int64     `json:"id"`
	Msg   string    `json:"msg"`
	Ctime time.Time `json:"ctime"`
	Mtime time.Time `json:"mtime"`
}

// DMVisible check dm is visible or not.
func DMVisible(state int32) bool {
	if state == StateNormal || state == StateHide || state == StateMonitorAfter {
		return true
	}
	return false
}

// SearchDMParams dm search params
type SearchDMParams struct {
	Type         int32  `form:"type"  validate:"required"`
	Oid          int64  `form:"oid"  validate:"required"`
	Keyword      string `form:"keyword"`
	Mid          int64  `form:"mid" default:"-10516"`
	IP           string `form:"ip"`
	State        string `form:"state"`
	Pool         string `form:"pool"`
	Attrs        string `form:"attrs"`
	ProgressFrom int64  `form:"progress_from" default:"-10516"`
	ProgressTo   int64  `form:"progress_to" default:"-10516"`
	CtimeFrom    int64  `form:"ctime_from" default:"-10516"`
	CtimeTo      int64  `form:"ctime_to" default:"-10516"`
	Page         int64  `form:"page" default:"1"`
	Size         int64  `form:"page_size" default:"100"`
	Sort         string `form:"sort"`
	Order        string `form:"order"`
}

// SearchDMData dm meta data from search
type SearchDMData struct {
	Order  string `json:"order"`
	Sort   string `json:"sort"`
	Result []*struct {
		ID int64 `json:"id"`
	} `json:"result"`
	Page *Page `json:"page"`
}

//SearchDMResult dm list
type SearchDMResult struct {
	Total     int64     `json:"total"`
	Count     int64     `json:"count"`
	MaxLimit  int64     `json:"max_limit"`
	Protected int64     `json:"protected"`
	Deleted   int64     `json:"deleted"`
	Page      int64     `json:"page"`
	Pagesize  int64     `json:"pagesize"`
	Result    []*DMItem `json:"result"`
}

//ListItem dm list item
type ListItem struct {
	ID       int64     `json:"id"`
	CID      int64     `json:"cid"`
	PoolID   int       `json:"pool_id"`
	Deleted  int       `json:"deleted"`
	UID      int64     `json:"uid"`
	Uname    string    `json:"uname"`
	IP       string    `json:"ip"`
	Playtime float64   `json:"playtime"`
	Model    int       `json:"model"`
	Msg      string    `json:"msg"`
	Fontsize int       `json:"fontsize"`
	Color    string    `json:"color"`
	Ctime    time.Time `json:"ctime"`
}

// DMItem dm list item from new db
type DMItem struct {
	IDStr    string    `json:"id_str"`
	ID       int64     `json:"id"`
	Type     int32     `json:"type"`
	Oid      int64     `json:"oid"`
	Mid      int64     `json:"mid"`
	Pool     int32     `json:"pool"`
	State    int32     `json:"state"`
	Attrs    string    `json:"attrs"`
	IP       int64     `json:"ip"`
	Progress int32     `json:"progress"`
	Mode     int32     `json:"mode"`
	Msg      string    `json:"msg"`
	Fontsize int32     `json:"fontsize"`
	Color    string    `json:"color"`
	Ctime    time.Time `json:"ctime"`
	Uname    string    `json:"uname"`
}

// DMSubject dm_inid info
type DMSubject struct {
	OID    int64     `json:"oid"`
	Type   int32     `json:"type"`
	AID    int64     `json:"aid"`
	MID    int64     `json:"uid"`
	ACount int64     `json:"count"`
	Limit  int64     `json:"limit"`
	TID    int64     `json:"tid"`
	TName  string    `json:"tname"`
	State  int32     `json:"state"`
	ETitle string    `json:"ep_title"`
	Title  string    `json:"title"`
	CTime  time.Time `json:"ctime"`
	MTime  time.Time `json:"mtime"`
}

//ArchiveResult archive list
type ArchiveResult struct {
	ArcLists []*DMSubject `json:"archives"`
	Page     *Page        `json:"page"`
}

//DMIndexInfo dm_inid index info
type DMIndexInfo struct {
	CID      int64  `json:"dm_inid"`
	AID      int64  `json:"aid"`
	MID      int64  `json:"mid"`
	UName    string `json:"u_name"`
	Duration int64  `json:"duration"`
	Limit    int64  `json:"limit"`
	Active   int64  `json:"dm_active"`
	ETitle   string `json:"ep_title"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	CTime    int64  `json:"ctime"`
	MTime    int64  `json:"mtime"`
}

// ArcVideo arc+video info by api
type ArcVideo struct {
	Archive *struct {
		AID       int64  `json:"aid"`
		MID       int64  `json:"mid"`
		TID       int64  `json:"tid"`
		Title     string `json:"title"`
		Cover     string `json:"cover"`
		RjReason  string `json:"reject_reason"`
		Tag       string `json:"tag"`
		Duration  int64  `json:"duration"`
		Copyright int64  `json:"copyright"`
		Desc      string `json:"desc"`
		MissionID int64  `json:"mission_id"`
		Attribute int64  `json:"attribute"`
		State     int64  `json:"state"`
		Source    string `json:"source"`
		NoReprint int64  `json:"no_reprint"`
		OrderID   int64  `json:"order_id"`
		DTime     int64  `json:"dtime"`
		PTime     int64  `json:"ptime"`
		CTime     int64  `json:"ctime"`
	} `json:"archive"`
	Videos []*struct {
		AID      int64  `json:"aid"`
		Title    string `json:"title"`
		Desc     string `json:"desc"`
		Filename string `json:"filename"`
		CID      int64  `json:"cid"`
		Index    int64  `json:"index"`
		Status   int64  `json:"status"`
		FailCode int64  `json:"fail_code"`
		XState   int64  `json:"xcode_state"`
		RjReason string `json:"reject_reason"`
		CTime    int64  `json:"ctime"`
	} `json:"videos"`
}

// ArchiveType archive type info
type ArchiveType struct {
	ID   int64  `json:"id"`
	PID  int64  `json:"pid"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

// ArchiveListReq archive list request
type ArchiveListReq struct {
	IDType string
	ID     int64
	Page   int64
	State  int64
	Attrs  []int64
	Pn     int64
	Ps     int64
	Sort   string
	Order  string
}

// UptSearchDMState update search dm state
type UptSearchDMState struct {
	ID    int64  `json:"id"`
	Oid   int64  `json:"oid"`
	Type  int32  `json:"type"`
	State int32  `json:"state"`
	Mtime string `json:"mtime"`
}

// UptSearchDMPool update search dm pool
type UptSearchDMPool struct {
	ID    int64  `json:"id"`
	Oid   int64  `json:"oid"`
	Type  int32  `json:"type"`
	Pool  int32  `json:"pool"`
	Mtime string `json:"mtime"`
}

// UptSearchDMAttr update search dm attr
type UptSearchDMAttr struct {
	ID         int64   `json:"id"`
	Oid        int64   `json:"oid"`
	Type       int32   `json:"type"`
	Attr       int32   `json:"attr"`
	Mtime      string  `json:"mtime"`
	AttrFormat []int64 `json:"attr_format"`
}

// MaskUp mask up info.
type MaskUp struct {
	ID      int64     `json:"id"`
	Mid     int64     `json:"mid"`
	Name    string    `json:"name"`
	State   int32     `json:"state"`
	Comment string    `json:"comment"`
	CTime   time.Time `json:"ctime"`
	MTime   time.Time `json:"mtime"`
}

// MaskUpRes maskUp and page info
type MaskUpRes struct {
	Result []*MaskUp `json:"result"`
	Page   *PageInfo `json:"page"`
}
