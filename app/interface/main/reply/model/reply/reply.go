package reply

import (
	"database/sql/driver"
	"strconv"
	"strings"
	"text/template"

	accmdl "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	xtime "go-common/library/time"
)

// subtype
const (
	Rule = "https://www.bilibili.com/blackboard/foldingreply.html"

	UserLevelFirst = 1
	UserLevelSnd   = 2

	SubTypeArchive      = int8(1)
	SubTypeTopic        = int8(2)
	SubTypeDrawyoo      = int8(3)
	SubTypeActivity     = int8(4)
	SubTypeLive         = int8(5)
	SubTypeForbiden     = int8(6) // reply forbiden info
	SubTypeNotice       = int8(7) //reply notice info
	SubTypeLiveAct      = int8(8)
	SubTypeActArc       = int8(9)
	SubTypeLiveNotice   = int8(10)
	SubTypeLivePicture  = int8(11) // 文画
	SubTypeArticle      = int8(12) // 文章
	SubTypeTicket       = int8(13) // 票务
	SubTypeMusic        = int8(14) // 音乐
	SubTypeCredit       = int8(15) // 风纪委
	SubTypePgcCmt       = int8(16) // pgc点评
	SubTypeDynamic      = int8(17) // 庐山动态
	SubTypePlaylist     = int8(18) // 播单
	SubTypeMusicList    = int8(19) // 音乐播单
	SubTypeComicSeason  = int8(20) // 漫画部评论
	SubTypeComicEpisode = int8(21) // 漫画话评论
	SubTypeHuoniao      = int8(22) // 火鸟
	SubTypeBBQ          = int8(23) // BBQ
	SubTypePGC          = int8(24) // PGC 全网落地页
	SubTypeGame         = int8(25) // 赛事库电竞
	SubTypeMedialist    = int8(26) // 播单(收藏夹)
	SubTypeEsportsInfo  = int8(27) // 电竞项目比赛数据页

	SubStateNormal = int8(0)
	SubStateForbid = int8(1)

	// Sub attr bit
	SubAttrAdminTop = uint32(0)
	SubAttrUpperTop = uint32(1)
	SubAttrMonitor  = uint32(2)
	SubAttrConfig   = uint32(3)
	SubAttrAudit    = uint32(4)
	SubAttrFrozen   = uint32(5)
	// 标识有被折叠的根评论
	SubAttrFolded = uint32(7)

	ReplyStateNormal    = int8(0)  // normal
	ReplyStateHidden    = int8(1)  // hidden by up
	ReplyStateFiltered  = int8(2)  // filtered
	ReplyStateAdminDel  = int8(3)  // delete by admin
	ReplyStateUserDel   = int8(4)  // delete by user
	ReplyStateMonitor   = int8(5)  // reply after audit
	ReplyStateGarbage   = int8(6)  // spam reply
	ReplyStateTop       = int8(7)  // top
	ReplyStateUpDel     = int8(8)  // delete by up
	ReplyStateBlacklist = int8(9)  // in a blacklist
	ReplyStateAssistDel = int8(10) // delete by assistant
	ReplyStateAudit     = int8(11) // 监管中
	ReplyStateFolded    = int8(12) // 被折叠

	// reply attr bit

	ReplyAttrAdminTop = uint32(0)
	ReplyAttrUpperTop = uint32(1)
	ReplyAttrGarbage  = uint32(2)
	ReplyAttrFilter   = uint32(3)
	// 标识有被折叠的子评论
	ReplyAttrFolded = uint32(7)

	SortByFloor = int8(0)
	SortByCount = int8(1)
	SortByLike  = int8(2)

	ReportStateNew         = int8(0) // 待一审
	ReportStateDelete      = int8(1) // 移除，已废弃
	ReportStateIgnore      = int8(2) // 忽略，已废弃
	ReportStateDeleteOne   = int8(3) // 一审移除
	ReportStateNewTwo      = int8(4) // 待二审
	ReportStateDeleteTwo   = int8(5) // 二审移除
	ReportStateIgnoreTwo   = int8(6) // 二审忽略
	ReportStateIgnoreOne   = int8(7) // 一审忽略
	ReportStateTransferred = int8(8) // 举报转移风纪委

	ReportAttrTransferred = uint(0) //标识举报State是否曾今从待一审\二审转换成待二审\一审

	ReportUserStateNew      = int8(0) // 新增
	ReportUserStateReported = int8(1) // 已反馈

	OpCancel = int8(0) // 取消赞踩
	OpAdd    = int8(1) // 添加赞踩

	ActionNormal = int8(0) // 未踩赞
	ActionLike   = int8(1) // 赞
	ActionHate   = int8(2) // 踩

	UpperOpHide = int8(1)
	UpperOpShow = int8(2)

	ReportReasonAd                 = int8(1)  // 广告
	ReportReasonPorn               = int8(2)  // 色情
	ReportReasonMeaningless        = int8(3)  // 刷屏
	ReportReasonProvoke            = int8(4)  // 引战
	ReportReasonSpoiler            = int8(5)  // 剧透
	ReportReasonPolitic            = int8(6)  // 政治
	ReportReasonAttack             = int8(7)  // 人身攻击
	ReportReasonUnrelated          = int8(8)  // 视频不相关
	ReportReasonProhibited         = int8(9)  // 违禁
	ReportReasonVulgar             = int8(10) // 低俗
	ReportReasonIllegalWebsite     = int8(11) // 非法网站
	ReportReasonGamblingFraud      = int8(12) // 赌博诈骗
	ReportReasonRumor              = int8(13) // 传播不实信息
	ReportReasonAbetting           = int8(14) // 怂恿教唆信息
	ReportReasonPrivacyInvasion    = int8(15) // 侵犯隐私
	ReportReasonUnlimitedSign      = int8(16) // 抢楼
	ReportReasonYouthInappropriate = int8(17) // 青少年不良信息
	ReportReasonOther              = int8(0)  // 其他

	PlatUnknow   = int8(0)
	PlatWeb      = int8(1)
	PlatAndroid  = int8(2)
	PlatIPhone   = int8(3)
	PlatWpM      = int8(4) // wp mobile
	PlatIPad     = int8(5)
	PlatPadHd    = int8(6) // ipad hd
	PlatWpPc     = int8(7) // win10
	PlatAndroidI = int8(8) // 国际版安卓

	AdminOperDelete         = int8(0)  // admin delete
	AdminOperDeleteByReport = int8(1)  // admin delete by report
	AdminOperIgnoreReport   = int8(2)  // admin ignore report
	AdminOperRecover        = int8(3)  // admin recover
	AdminOperEdit           = int8(4)  // admin edit reply content
	AdminOperPass           = int8(5)  // admin pass
	AdminOperSubState       = int8(6)  // admin change subject state
	AdminOperSubTop         = int8(7)  // top reply
	AdminOperSubMid         = int8(8)  // admin change subject mid
	AdminOperRptIgnore1     = int8(9)  // admin report ignore 1
	AdminOperRptIgnore2     = int8(10) // admin report ignore 2
	AdminOperRptDel1        = int8(11) // admin report del 1
	AdminOperRptDel2        = int8(12) // admin report del 2
	AdminOperRptRecover1    = int8(13) // admin report recover 1
	AdminOperRptRecover2    = int8(14) // admin report recover 2
	AdminOperActionSet      = int8(15) // admin action set

	AdminIsNotReport = int8(0)
	AdminIsReport    = int8(1)
	AdminIsNotNew    = int8(0)
	AdminIsNew       = int8(1)

	AuditTypeOne = int8(1) // 一审
	AuditTypeTwo = int8(2) // 二审

	BlacklistRelation = int16(-1)

	ReportReplyAdd = "reply_add"
	ReportReplyDel = "reply_del"

	ReportReplyLike       = "reply_like"
	ReportReplyHate       = "reply_hate"
	ReportReplyCancelLike = "reply_cancel_like"
	ReportReplyCancelHate = "reply_cancel_hate"

	ReportReplyTop   = "reply_top"
	ReportReplyUntop = "reply_untop"

	ReportReplyReport = "reply_report"
)

// ActionCount ActionCount
type ActionCount struct {
	Like int32 `json:"like"`
	Hate int32 `json:"hate"`
}

// Subject ReplySubject
type Subject struct {
	ID     int64      `json:"-"`
	Oid    int64      `json:"oid"`
	Type   int8       `json:"type"`
	Mid    int64      `json:"mid"`
	Count  int        `json:"count"`
	RCount int        `json:"rcount"`
	ACount int        `json:"acount"`
	State  int8       `json:"state"`
	Attr   uint32     `json:"attr"`
	Meta   string     `json:"meta"`
	CTime  xtime.Time `json:"ctime"`
	MTime  xtime.Time `json:"-"`
}

// AttrVal return val of subject'attr.
func (s *Subject) AttrVal(bit uint32) uint32 {
	return (s.Attr >> bit) & uint32(1)
}

// Folder ...
func (s *Subject) Folder() (f Folder) {
	f.HasFolded = s.HasFolded()
	f.Rule = Rule
	return
}

// HasFolded ...
func (s *Subject) HasFolded() bool {
	return s.AttrVal(SubAttrFolded) == AttrYes
}

// SubjectMeta SubjectMeta
type SubjectMeta struct {
	AdminTop int64 `json:"atop,omitempty"`
	UpperTop int64 `json:"utop,omitempty"`
}

// IsNormal IsNormal
func (s *Subject) IsNormal() bool {
	return s.State == SubStateNormal
}

// LegalSubjectType LegalSubjectType
func LegalSubjectType(tp int8) bool {
	return SubTypeArchive <= tp && tp <= SubTypeEsportsInfo
}

// CheckSubForbid CheckSubForbid
func CheckSubForbid(state int8) (err error) {
	if state == SubStateForbid {
		err = ecode.ReplyForbidReply
	}
	return
}

// CheckSubState CheckSubState
func CheckSubState(state int8) (err error) {
	if state < SubStateNormal || state > SubStateForbid {
		err = ecode.ReplyIllegalSubState
	}
	return
}

// Counts ReplyCounts
type Counts struct {
	SubjectState int8  `json:"sub_state"`
	Counts       int64 `json:"count"`
}

// Reply Reply
type Reply struct {
	RpID      int64      `json:"rpid"`
	Oid       int64      `json:"oid"`
	Type      int8       `json:"type"`
	Mid       int64      `json:"mid"`
	Root      int64      `json:"root"`
	Parent    int64      `json:"parent"`
	Dialog    int64      `json:"dialog"`
	Count     int        `json:"count"`
	RCount    int        `json:"rcount"`
	Floor     int        `json:"floor,omitempty"`
	State     int8       `json:"state"`
	FansGrade int8       `json:"fansgrade"`
	Attr      uint32     `json:"attr"`
	CTime     xtime.Time `json:"ctime"`
	MTime     xtime.Time `json:"-"`
	// string
	RpIDStr   string `json:"rpid_str,omitempty"`
	RootStr   string `json:"root_str,omitempty"`
	ParentStr string `json:"parent_str,omitempty"`
	DialogStr string `json:"dialog_str",omitempty`
	// action count, from ReplyAction count
	Like   int  `json:"like"`
	Hate   int  `json:"-"`
	Action int8 `json:"action"`
	// member info
	Member *Member `json:"member"`
	// other
	Content *Content `json:"content"`
	Replies []*Reply `json:"replies"`
	Assist  int      `json:"assist"`
	// 是否有折叠评论
	Folder Folder `json:"folder"`
}

type Folder struct {
	HasFolded bool   `json:"has_folded"`
	IsFolded  bool   `json:"is_folded"`
	Rule      string `json:"rule"`
}

// FillFolder ...
func (r *Reply) FillFolder() {
	if r.IsRoot() {
		r.Folder.HasFolded = r.HasFolded()
		r.Folder.IsFolded = r.IsFolded()
		r.Folder.Rule = Rule
	}
}

// HasFolded ...
func (r *Reply) HasFolded() bool {
	return r.AttrVal(ReplyAttrFolded) == AttrYes
}

// IsFolded ...
func (r *Reply) IsFolded() bool {
	return r.State == ReplyStateFolded
}

// UnmarkHasFolded ...
func (r *Reply) UnmarkHasFolded() {
	r.AttrSet(AttrNo, ReplyAttrFolded)
}

// AttrSet set attr of reply'attr
func (r *Reply) AttrSet(v uint32, bit uint32) {
	r.Attr = r.Attr&(^(1 << bit)) | (v << bit)
}

// AttrVal return val of reply'attr
func (r *Reply) AttrVal(bit uint32) uint32 {
	if r.Attr == 0 {
		return uint32(0)
	}
	return (r.Attr >> bit) & uint32(1)
}

// IsRoot IsRoot
func (r *Reply) IsRoot() bool {
	return r.Root == 0 && r.Parent == 0
}

// IsTop IsTop
func (r *Reply) IsTop() bool {
	return r.AttrVal(ReplyAttrAdminTop) == AttrYes || r.AttrVal(ReplyAttrUpperTop) == AttrYes
}

// IsNormal IsNormal
func (r *Reply) IsNormal() bool {
	return r.State == ReplyStateNormal || r.State == ReplyStateHidden || r.State == ReplyStateFiltered || r.State == ReplyStateGarbage || r.State == ReplyStateMonitor || r.State == ReplyStateTop || r.State == ReplyStateFolded
}

// IsDeleted IsDeleted
func (r *Reply) IsDeleted() bool {
	return r.State == ReplyStateUserDel || r.State == ReplyStateUpDel || r.State == ReplyStateAdminDel || r.State == ReplyStateAssistDel
}

// FillStr FillStr
func (r *Reply) FillStr(isEscape bool) {
	r.RpIDStr = strconv.FormatInt(r.RpID, 10)
	r.RootStr = strconv.FormatInt(r.Root, 10)
	r.ParentStr = strconv.FormatInt(r.Parent, 10)
	if r.Content != nil {
		if isEscape {
			r.Content.Message = template.HTMLEscapeString(r.Content.Message)
		}
		r.Content.IP = 0
		r.Content.Version = ""
		r.Content.Members = []*Info{}
	}
}

// Clone clone a reply content.
func (r *Reply) Clone() (res *Reply) {
	content := new(Content)
	if r.Content != nil {
		*content = *(r.Content)
	}
	res = new(Reply)
	*res = *r
	res.Content = content
	return
}

// WithinSortRange WithinSortRange
func WithinSortRange(sort int8) bool {
	return SortByFloor <= sort && sort <= SortByLike
}

// CheckSort WithinSortRange
func CheckSort(sort int8) bool {
	return SortByFloor <= sort && sort <= SortByLike
}

// CheckPlat CheckPlat
func CheckPlat(plat int8) bool {
	return PlatUnknow <= plat && plat <= PlatWpPc
}

// Content ReplyContent
type Content struct {
	RpID    int64  `json:"-"`
	Message string `json:"message"`

	Ats     Int64Bytes `json:"ats,omitempty"`
	Topics  Mstr       `json:"topics,omitempty"`
	IP      uint32     `json:"ipi,omitempty"`
	Plat    int8       `json:"plat"`
	Device  string     `json:"device"`
	Version string     `json:"version,omitempty"`
	CTime   xtime.Time `json:"-"`
	MTime   xtime.Time `json:"-"`
	// ats member info
	Members []*Info `json:"members"`
}

// FillAts FillAts
func (rc *Content) FillAts(cards map[int64]*accmdl.Card) {
	rc.Members = make([]*Info, 0, len(rc.Ats))
	for _, at := range rc.Ats {
		if card, ok := cards[at]; ok {
			i := &Info{}
			i.FromCard(card)
			rc.Members = append(rc.Members, i)
		}
	}
	rc.Ats = nil
}

// Action Action
type Action struct {
	ID     int64      `json:"-"`
	RpID   int64      `json:"rpid"`
	Action int8       `json:"action"`
	Mid    int64      `json:"mid"`
	CTime  xtime.Time `json:"-"`
}

// CheckAction CheckAction
func CheckAction(act int8) (err error) {
	if act != OpAdd && act != OpCancel {
		err = ecode.ReplyIllegalAction
	}
	return
}

// Report Report
type Report struct {
	ID      int64      `json:"id"`
	RpID    int64      `json:"rpid"`
	Oid     int64      `json:"oid"`
	Type    int8       `json:"type"`
	Mid     int64      `json:"mid"`
	Reason  int8       `json:"reason"`
	Content string     `json:"content"`
	Count   int        `json:"count"`
	Score   int        `json:"score"`
	State   int8       `json:"state"`
	CTime   xtime.Time `json:"ctime"`
	MTime   xtime.Time `json:"-"`
	Attr    int8       `json:"attr"`
}

// AttrVal AttrVal
func (rpt Report) AttrVal(bit uint) int8 {
	return (rpt.Attr >> bit) & int8(1)
}

// IsTransffered IsTransffered
func (rpt Report) IsTransffered() bool {
	return rpt.AttrVal(ReportAttrTransferred) == 1
}

// ReportUser report user.
type ReportUser struct {
	ID      int64      `json:"id"`
	Oid     int64      `json:"oid"`
	Type    int8       `json:"type"`
	RpID    int64      `json:"rpid"`
	Mid     int64      `json:"mid"`
	Reason  int8       `json:"reason"`
	Content string     `json:"content"`
	State   int8       `json:"state"`
	CTime   xtime.Time `json:"ctime"`
	MTime   xtime.Time `json:"mtime"`
}

// CheckReportReason CheckReportReason
func CheckReportReason(reason int8) (err error) {
	if !(ReportReasonOther <= reason && reason <= ReportReasonYouthInappropriate) {
		err = ecode.ReplyIllegalReport
	}
	return
}

// GetReportType GetReportType
func GetReportType(reason int8) int8 {
	if (reason >= ReportReasonMeaningless && reason <= ReportReasonSpoiler) || reason == ReportReasonUnrelated || reason == ReportReasonOther {
		return ReportStateNewTwo
	}
	return ReportStateNew
}

// Mstr Mstr
type Mstr []string

// Scan Scan
func (ms *Mstr) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case []byte:
		if len(sc) == 0 {
			return
		}
		res := strings.Split(string(sc), ",")
		for i := range res {
			res[i] = strings.Replace(res[i], "%2c", ",", -1)
		}
		*ms = res
	}
	return
}

// Value Value
func (ms Mstr) Value() (driver.Value, error) {
	return ms.Bytes(), nil
}

// Bytes Bytes
func (ms Mstr) Bytes() []byte {
	var res string
	for i := range ms {
		str := strings.Replace(ms[i], ",", "%2c", -1)
		res += str
		if i != (len(ms) - 1) {
			res += ","
		}
	}
	return []byte(res)
}

// Info Info
type Info struct {
	Mid         string `json:"mid"`
	Name        string `json:"uname"`
	Sex         string `json:"sex"`
	Sign        string `json:"sign"`
	Avatar      string `json:"avatar"`
	Rank        string `json:"rank"`
	DisplayRank string `json:"DisplayRank"`
	LevelInfo   struct {
		Cur     int `json:"current_level"`
		Min     int `json:"current_min"`
		NowExp  int `json:"current_exp"`
		NextExp int `json:"next_exp"`
	} `json:"level_info"`
	Pendant        accmdl.PendantInfo   `json:"pendant"`
	Nameplate      accmdl.NameplateInfo `json:"nameplate"`
	OfficialVerify struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"`
	Vip struct {
		Type          int    `json:"vipType"`
		DueDate       int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
	} `json:"vip"`
}

// FromCard FromCard
func (i *Info) FromCard(c *accmdl.Card) {
	i.Mid = strconv.FormatInt(c.Mid, 10)
	i.Name = c.Name
	i.Sex = c.Sex
	i.Sign = c.Sign
	i.Avatar = c.Face
	i.Rank = strconv.FormatInt(int64(c.Rank), 10)
	i.DisplayRank = "0"
	i.LevelInfo.Cur = int(c.Level)
	i.Pendant = c.Pendant
	i.Nameplate = c.Nameplate
	if c.Official.Role == 0 {
		i.OfficialVerify.Type = -1
		i.OfficialVerify.Desc = ""
	} else {
		if c.Official.Role <= 2 {
			i.OfficialVerify.Type = 0
			i.OfficialVerify.Desc = c.Official.Title
		} else {
			i.OfficialVerify.Type = 1
			i.OfficialVerify.Desc = c.Official.Title
		}
	}
	i.Vip.Type = int(c.Vip.Type)
	i.Vip.VipStatus = int(c.Vip.Status)
	i.Vip.DueDate = c.Vip.DueDate
}

// Business Business
type Business struct {
	Type   int32  `json:"type"`
	Alias  string `json:"alias"`
	Appkey string `json:"appkey"`
}

type DialogCursor struct {
	MinFloor int `json:"min_floor"`
	MaxFloor int `json:"max_floor"`
	Size     int `json:"size"`
}

type DialogMeta struct {
	MinFloor int `json:"min_floor"`
	MaxFloor int `json:"max_floor"`
}

// ShouldShowFolded ...
func ShouldShowFolded(mobi_app string, build, scene int64) bool {
	if mobi_app == "android" && build < 5365000 {
		return true
	}
	if mobi_app == "iphone" && build < 8310 {
		return true
	}
	if scene == 1 {
		return true
	}
	return false
}
