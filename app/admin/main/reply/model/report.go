package model

import xtime "go-common/library/time"

var (
	// DateFormat date time format
	DateFormat = "2006-01-02 15:04:05"
	// DateSimpleFormat date simple format
	DateSimpleFormat = "2006-01-02"
)

const (
	// ReportStateNew 待一审
	ReportStateNew = int32(0)
	// ReportStateDelete 移除
	ReportStateDelete = int32(1)
	// ReportStateIgnore 忽略
	ReportStateIgnore = int32(2)
	// ReportStateDelete1 一审移除
	ReportStateDelete1 = int32(3)
	// ReportStateNew2 待二审
	ReportStateNew2 = int32(4)
	// ReportStateDelete2 二审移除
	ReportStateDelete2 = int32(5)
	// ReportStateIgnore2 二审忽略
	ReportStateIgnore2 = int32(6)
	// ReportStateIgnore1 一审忽略
	ReportStateIgnore1 = int32(7)
	// ReportStateTransferred 举报转移风纪委
	ReportStateTransferred = int32(8)
	// ReportUserStateNew 新增
	ReportUserStateNew = int32(0)
	// ReportUserStateReported 已反馈
	ReportUserStateReported = int32(1)
	// ReportAttrTransferred 是否从待一审\二审转换成待二审\一审
	ReportAttrTransferred = uint32(0)
	// AuditTypeFirst 一审
	AuditTypeFirst = int32(1)
	// AuditTypeSecond 二审
	AuditTypeSecond = int32(2)
	// ReportActionReplyPass action reply_pass
	ReportActionReplyPass = "reply_pass"
	// ReportActionReplyDel action reply_del
	ReportActionReplyDel = "reply_del"
	// ReportActionReplyEdit action reply_edit
	ReportActionReplyEdit = "reply_edit"
	// ReportActionReplyRecover action reply_recover
	ReportActionReplyRecover = "reply_recover"
	// ReportActionReplyTop action 置顶
	ReportActionReplyTop = "top"
	// ReportActionReplyMonitor action 监控
	ReportActionReplyMonitor = "monitor"
	// ReportActionReplyGarbage action reply_garbage
	ReportActionReplyGarbage = "reply_garbage"
	// ReportActionReportIgnore1 action report_ignore_1
	ReportActionReportIgnore1 = "report_ignore_1"
	// ReportActionReportIgnore2 action report_ignore_2
	ReportActionReportIgnore2 = "report_ignore_2"
	// ReportActionReportDel1 action report_del_1
	ReportActionReportDel1 = "report_del_1"
	// ReportActionReportDel2 action report_del_2
	ReportActionReportDel2 = "report_del_2"
	// ReportActionReport1To2 action report_1to2
	ReportActionReport1To2 = "report_1to2"
	// ReportActionReport2To1 action report_2to1
	ReportActionReport2To1 = "report_2to1"
	// ReportActionReportArbitration action 众裁
	ReportActionReportArbitration = "report_arbitration"
)

// Report report info.
type Report struct {
	ID         int64      `json:"id"`
	RpID       int64      `json:"rpid"`
	Oid        int64      `json:"oid"`
	Type       int32      `json:"type"`
	Mid        int64      `json:"mid"`
	Reason     int32      `json:"reason"`
	Content    string     `json:"content"`
	Count      int32      `json:"count"`
	Score      int        `json:"score"`
	State      int32      `json:"state"`
	CTime      xtime.Time `json:"ctime"`
	MTime      xtime.Time `json:"mtime"`
	Attr       uint32     `json:"attr"`
	ReplyCtime xtime.Time `json:"-"`
}

// AttrVal return attr val.
func (r *Report) AttrVal(bit uint32) uint32 {
	return (r.Attr >> bit) & uint32(1)
}

// AttrSet set attr of ReplyReport'attr
func (r *Report) AttrSet(v uint32, bit uint32) {
	r.Attr = r.Attr&(^(1 << bit)) | (v << bit)
}

// SearchReportParams search params.
type SearchReportParams struct {
	Type      int32
	Oid       int64
	UID       int64
	StartTime string
	EndTime   string
	Reason    string
	Typeids   string
	Keyword   string
	Nickname  string
	States    string
	Order     string
	Sort      string
}

// ReportUser report user.
type ReportUser struct {
	ID      int64      `json:"id"`
	Oid     int64      `json:"oid"`
	Type    int8       `json:"type"`
	RpID    int64      `json:"rpid"`
	Mid     int64      `json:"mid"`
	Reason  int32      `json:"reason"`
	Content string     `json:"content"`
	State   int32      `json:"state"`
	CTime   xtime.Time `json:"ctime"`
	MTime   xtime.Time `json:"mtime"`
}

// SearchReport search report.
type SearchReport struct {
	ID          int64  `json:"id"`
	Oid         int64  `json:"oid"`
	OidStr      string `json:"oid_str"`
	Type        int8   `json:"type"`
	RpID        int64  `json:"rpid"`
	Mid         int64  `json:"mid"`
	Reason      int8   `json:"reason"`
	Content     string `json:"content"`
	State       int8   `json:"state"`
	CTime       string `json:"ctime"`
	MTime       string `json:"mtime"`
	Parent      int64  `json:"parent"`
	Like        int64  `json:"like"`
	ReplyState  int64  `json:"reply_state"`
	Opremark    string `json:"opremark"`
	Count       int64  `json:"count"`
	Message     string `json:"message"`
	Title       string `json:"title"`
	Opresult    string `json:"opresult"`
	ReplyMid    int64  `json:"reply_mid"`
	Floor       int64  `json:"floor"`
	Root        int64  `json:"root"`
	ReportMid   int64  `json:"report_mid"`
	ArcMid      int64  `json:"arc_mid"`
	Reporter    string `json:"reporter"`
	Replier     string `json:"replier"`
	IsUp        int64  `json:"is_up"`
	AdminID     int64  `json:"adminid"`
	AdminName   string `json:"admin_name"`
	Opctime     string `json:"opctime"`
	DocID       string `json:"doc_id"`
	Score       int64  `json:"score"`
	Attr        []int8 `json:"attr"`
	RedirectURL string `json:"redirect_url"`
}

// SearchReportResult search result.
type SearchReportResult struct {
	Code      int             `json:"code,omitempty"`
	Page      int64           `json:"page"`
	PageSize  int64           `json:"pagesize"`
	PageCount int             `json:"pagecount"`
	Total     int64           `json:"total"`
	Order     string          `json:"order"`
	Result    []*SearchReport `json:"result"`
	Message   string          `json:"msg,omitempty"`
}

// ReportSearchResponse search result.
type ReportSearchResponse struct {
	SearchReportResult
	Pager Pager `json:"pager"`
}
