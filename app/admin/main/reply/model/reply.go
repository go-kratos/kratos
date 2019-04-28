package model

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"

	rl "go-common/app/service/main/relation/model"
	xtime "go-common/library/time"
)

// State 评论状态
const (
	StateNormal    int32 = 0  // 正常
	StateHidden    int32 = 1  // up主隐藏
	StateFiltered  int32 = 2  // 敏感词过滤 @Deprecated to use attr 3
	StateDelAdmin  int32 = 3  // 管理员删除
	StateDelUser   int32 = 4  // 用户删除
	StateMonitor   int32 = 5  // 监控中
	StateGarbage   int32 = 6  // 大数据过虑 @Deprecated to use attr 2
	StateTopAdmin  int32 = 7  // 管理员置顶 @Deprecated to use attr 1
	StateDelUpper  int32 = 8  // up主删除
	StateBlacklist int32 = 9  // 黑名单屏蔽
	StateDelAssist int32 = 10 // 协管删除
	StatePending   int32 = 11 // 先审后发
	StateFolded    int32 = 12 // 被折叠
)

// Attr 评论属性
const (
	AttrTopAdmin uint32 = 0 // 管理员置顶
	AttrTopUpper uint32 = 1 // up主置顶
	AttrGarbage  uint32 = 2 // 大数据过滤
	AttrFiltered uint32 = 3 // 敏感词过滤

	// 有子评论被折叠
	AttrFolded uint32 = 7
)

// SortBy 排序方式
const (
	SortByFloor int32 = 0 // 按楼层排序
	SortByCount int32 = 1 // 按评论数排序
	SortByLike  int32 = 2 // 按点赞数排序
)

// 折叠评论类型
const (
	FolderKindSub  = "s"
	FolderKindRoot = "r"
)

// SearchParams reply params.
type SearchParams struct {
	Type        int32
	Oid         int64
	TypeIds     string
	Keyword     string
	KeywordHigh string
	UID         int64
	Uname       string
	AdminID     int64
	AdminName   string
	Begin       time.Time
	End         time.Time
	States      string
	IP          int64
	Attr        string
	Sort        string
	Order       string
}

// ReplySearchResponse ReplySearchResponse
type ReplySearchResponse struct {
	SearchResult
	Pager Pager `json:"pager"`
}

// SearchResult search result.
type SearchResult struct {
	Code      int            `json:"code"`
	Message   string         `json:"msg,omitempty"`
	Order     string         `json:"order"`
	Page      int64          `json:"page"`
	PageSize  int64          `json:"pagesize"`
	PageCount int64          `json:"pagecount"`
	Total     int64          `json:"total"`
	Result    []*SearchReply `json:"result"`
}

// SearchReply search reply.
type SearchReply struct {
	// 评论基础信息
	ID     int64  `json:"id"`
	Type   int8   `json:"type"`
	Oid    int64  `json:"oid"`
	OidStr string `json:"oid_str"`
	State  int8   `json:"state"`
	Floor  int64  `json:"floor"`
	Ctime  string `json:"ctime"`
	Mtime  string `json:"mtime"`
	Attr   Attr   `json:"attr"`
	Title  string `json:"title"`

	// 评论人的相关信息
	Mid     int64    `json:"mid"`
	Stat    *rl.Stat `json:"stat"`
	Replier string   `json:"replier"`
	IP      IP       `json:"ip"`
	Message string   `json:"message"`
	Typeid  int      `json:"typeid"`
	Root    int      `json:"root"`

	// 后台操作信息
	AdminID     int64  `json:"adminid"`
	AdminName   string `json:"admin_name"`
	Opremark    string `json:"opremark"`
	Opresult    string `json:"opresult"`
	OpCtime     string `json:"opctime"`
	RedirectURL string `json:"redirect_url"`
	// 搜索返回的额外数据
	DocID string `json:"doc_id"`
}
type IP int64

func (ip *IP) UnmarshalJSON(b []byte) error {
	if string(b) == "" {
		return nil
	}
	str := strings.Trim(string(b), `"`)
	ipI := net.ParseIP(str).To4()
	if ipI == nil {
		return nil
	}
	*ip = IP(int64(binary.BigEndian.Uint32(ipI)))
	return nil
}

type Attr []int32

func (a *Attr) UnmarshalJSON(b []byte) error {
	var s []int32
	err := json.Unmarshal(b, &s)
	if err != nil {
		attr, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
		var i int32 = 1
		for attr != 0 && i < 64 {
			if attr&1 == 1 {
				*a = append(*a, i)
			}
			attr = attr >> 1
			i++
		}
	} else {
		*a = s
	}
	return nil
}

// Reply info.b
type ReplyEx struct {
	Reply
	IsUp      bool  `json:"is_up"`
	RootFloor int32 `json:"root_floor"`
}

// Reply info.
type Reply struct {
	ID     int64      `json:"rpid"`
	Oid    int64      `json:"oid"`
	Type   int32      `json:"type"`
	Mid    int64      `json:"mid"`
	Root   int64      `json:"root"`
	Parent int64      `json:"parent"`
	Dialog int64      `json:"dialog"`
	Count  int32      `json:"count"`
	MCount int32      `json:"mcount"`
	RCount int32      `json:"rcount"`
	Floor  int32      `json:"floor"`
	State  int32      `json:"state"`
	Attr   uint32     `json:"attr"`
	CTime  xtime.Time `json:"ctime"`
	MTime  xtime.Time `json:"-"`
	// action info
	Like    int32         `json:"like"`
	Hate    int32         `json:"hate"`
	Action  int32         `json:"action"`
	Content *ReplyContent `json:"content"`
}

// HasFolded ...
func (r *Reply) HasFolded() bool {
	return r.AttrVal(AttrFolded) == AttrYes
}

// MarkHasFolded ...
func (r *Reply) MarkHasFolded() {
	r.AttrSet(AttrYes, AttrFolded)
}

// UnmarkHasFolded ...
func (r *Reply) UnmarkHasFolded() {
	r.AttrSet(AttrNo, AttrFolded)
}

// DenyFolded ...
func (r *Reply) DenyFolded() bool {
	return r.IsTop() || !r.AllowFoldState() || r.Type == SubTypeArticle
}

// AllowFoldState ...
func (r *Reply) AllowFoldState() bool {
	return r.State == StateNormal || r.State == StateHidden || r.State == StateFiltered || r.State == StateGarbage
}

// IsFolded .
func (r *Reply) IsFolded() bool {
	return r.State == StateFolded
}

// IsRoot root.
func (r *Reply) IsRoot() bool {
	return r.Root == 0
}

// IsTop top.
func (r *Reply) IsTop() bool {
	if r.Attr != 0 && (r.AttrVal(AttrTopAdmin) == 1 || r.AttrVal(AttrTopUpper) == 1) {
		return true
	}
	return false
}

// IsDeleted deleted.
func (r *Reply) IsDeleted() bool {
	return r.State == StateDelUser || r.State == StateDelUpper || r.State == StateDelAdmin
}

// AttrVal return val of reply'attr
func (r *Reply) AttrVal(bit uint32) uint32 {
	if r.Attr == 0 {
		return uint32(0)
	}
	return (r.Attr >> bit) & uint32(1)
}

// AttrSet set attr of reply'attr
func (r *Reply) AttrSet(v uint32, bit uint32) {
	r.Attr = r.Attr&(^(1 << bit)) | (v << bit)
}

// IsNormal IsNormal
func (r *Reply) IsNormal() bool {
	return r.State == StateNormal || r.State == StateHidden || r.State == StateFiltered || r.State == StateMonitor || r.State == StateGarbage || r.State == StateTopAdmin || r.State == StateFolded
}

func (r *Reply) IsPending() bool {
	return r.State == StatePending
}

// LegalSubjectType LegalSubjectType
func LegalSubjectType(tp int32) bool {
	return SubTypeArchive <= tp && tp <= SubTypeComicEpisode
}

// ReplyContent define reply content
type ReplyContent struct {
	ID      int64      `json:"-"`
	Message string     `json:"message"`
	Ats     Int64Bytes `json:"ats,omitempty"`
	IP      uint32     `json:"ipi,omitempty"`
	Plat    int8       `json:"plat"`
	Device  string     `json:"device"`
	Version string     `json:"version,omitempty"`
	CTime   xtime.Time `json:"-"`
	MTime   xtime.Time `json:"-"`
}

// LogSearchParam LogSearchParam
type LogSearchParam struct {
	Oid       int64
	Type      int32
	Mid       int64
	CtimeFrom string
	CtimeTo   string
	Action    string
	Pn        int64
	Ps        int64
	Other     int64
	Sort      string
	Order     string
	Group     string
	Appid     string
}

// ReplyTopLogResult ReplyTopLogResult
type ReplyTopLogResult struct {
	Logs  []*ReplyTopLog `json:"logs"`
	Page  Page           `json:"page"`
	Order string         `json:"order"`
	Sort  string         `json:"sort"`
}

// ReplyTopLog ReplyTopLog
type ReplyTopLog struct {
	AdminID     int64  `json:"adminid"`
	AdminName   string `json:"admin_name"`
	Oid         int64  `json:"oid"`
	Type        int32  `json:"type"`
	Title       string `json:"title"`
	RedirectURL string `json:"redirect_url"`
	Remark      string `json:"remark"`
	UserName    string `json:"username"`
	Mid         int64  `json:"mid"`
	CTime       string `json:"ctime"`
	RpID        int64  `json:"rpid"`
	Action      int64  `json:"action"`
}

// ExportedReply exported reply struct
type ExportedReply struct {
	ID      int64     `json:"rpid"`
	Oid     int64     `json:"oid"`
	Type    int32     `json:"type"`
	Mid     int64     `json:"mid"`
	Root    int64     `json:"root"`
	Parent  int64     `json:"parent"`
	Count   int32     `json:"count"`
	RCount  int32     `json:"rcount"`
	Like    int32     `json:"like"`
	Hate    int32     `json:"hate"`
	Floor   int32     `json:"floor"`
	State   int32     `json:"state"`
	Attr    int32     `json:"attr"`
	CTime   time.Time `json:"ctime"`
	Message string    `json:"message"`
}

// String convert ExportedReply to string
func (e *ExportedReply) String() (s []string) {
	s = append(s, strconv.FormatInt(e.ID, 10))
	s = append(s, strconv.FormatInt(e.Oid, 10))
	s = append(s, strconv.FormatInt(int64(e.Type), 10))
	s = append(s, strconv.FormatInt(e.Mid, 10))
	s = append(s, strconv.FormatInt(e.Root, 10))
	s = append(s, strconv.FormatInt(e.Parent, 10))
	s = append(s, strconv.FormatInt(int64(e.Count), 10))
	s = append(s, strconv.FormatInt(int64(e.RCount), 10))
	s = append(s, strconv.FormatInt(int64(e.Like), 10))
	s = append(s, strconv.FormatInt(int64(e.Hate), 10))
	s = append(s, strconv.FormatInt(int64(e.Floor), 10))
	s = append(s, strconv.FormatInt(int64(e.State), 10))
	s = append(s, strconv.FormatInt(int64(e.Attr), 10))
	s = append(s, e.CTime.String())
	s = append(s, e.Message)
	return
}
