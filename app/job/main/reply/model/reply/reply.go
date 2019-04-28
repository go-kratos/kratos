package reply

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"database/sql/driver"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"
)

// reply cnost
const (
	AttrYes        = uint32(1)
	AttrNo         = uint32(0)
	FolderKindSub  = "s"
	FolderKindRoot = "r"

	SortByFloor = int8(0)
	SortByCount = int8(1)
	SortByLike  = int8(2)

	//SubTypeVideo = define.TypeVideo
	//SubTypeTopic = define.TypeTopic

	SubTypeVideo         = int8(1)
	SubTypeTopic         = int8(2)
	SubTypeDrawyoo       = int8(3)
	SubTypeActivity      = int8(4)
	SubTypeLiveVideo     = int8(5)
	SubTypeForbiden      = int8(6) // reply forbiden info
	SubTypeNotice        = int8(7) // reply notice info
	SubTypeLiveAct       = int8(8)
	SubTypeActArc        = int8(9)
	SubTypeLiveNotice    = int8(10)
	SubTypeLivePicture   = int8(11) // 文画
	SubTypeArticle       = int8(12) // 文章
	SubTypeTicket        = int8(13) // 票务
	SubTypeAudio         = int8(14) // 音乐
	SubTypeCredit        = int8(15) // 风纪委
	SubTypeDynamic       = int8(17) // 动态
	SubTypePlaylist      = int8(18) // 播单
	SubTypeAudioPlaylist = int8(19) // 音乐播单

	SubStateNormal          = int8(0)
	SubStateForbid          = int8(1)
	SubStateReplyAfterAudit = int8(2)

	// Sub attr bit
	SubAttrAdminTop = uint32(0)
	SubAttrUpperTop = uint32(1)
	SubAttrFolded   = uint32(7)

	ReplyStateNormal    = int8(0)  // normal
	ReplyStateHidden    = int8(1)  // hidden by up
	ReplyStateFiltered  = int8(2)  // filtered
	ReplyStateAdminDel  = int8(3)  // delete by admin
	ReplyStateUserDel   = int8(4)  // delele by user
	ReplyStateMonitor   = int8(5)  // reply monitor
	ReplyStateGarbage   = int8(6)  // bigdata filter reply
	ReplyStateTop       = int8(7)  // top by admin
	ReplyStateUpDel     = int8(8)  // delete by up
	ReplyStateBlacklist = int8(9)  // in a blacklist
	ReplyStateAssistDel = int8(10) // delete by assistant
	ReplyStateAudit     = int8(11) // reply after audit
	ReplyStateFolded    = int8(12) // 被折叠

	// reply attr bit
	ReplyAttrAdminTop = uint32(0)
	ReplyAttrUpperTop = uint32(1)
	ReplyAttrGarbage  = uint32(2)
	ReplyAttrFolded   = uint32(7)

	ReportStateNew         = int8(0) // 待一审
	ReportStateDelete      = int8(1) // 移除
	ReportStateIgnore      = int8(2) // 忽略
	ReportStateDeleteOne   = int8(3) // 一审移除
	ReportStateNewTwo      = int8(4) // 待二审
	ReportStateDeleteTwo   = int8(5) // 二审移除
	ReportStateIgnoreTwo   = int8(6) // 二审忽略
	ReportStateIgnoreOne   = int8(7) // 一审忽略
	ReportStateTransferred = int8(8) // 举报转移风纪委

	ReportAttrTransferred = uint(0) //State状态是否曾今从待一审\二审转换成待二审\一审

	ReportUserStateNew      = int8(0) // 新增
	ReportUserStateReported = int8(1) // 已反馈

	OpCancel = int8(0) // 取消赞踩
	OpAdd    = int8(1) // 添加赞踩

	ActionNormal = int8(0) // 未踩赞
	ActionLike   = int8(1) // 赞
	ActionHate   = int8(2) // 踩

	UpperOpHide = int8(1)
	UpperOpShow = int8(2)

	ReportReasonOther           = int8(0)  // 其他
	ReportReasonAd              = int8(1)  // 广告
	ReportReasonPorn            = int8(2)  // 色情
	ReportReasonMeaningless     = int8(3)  // 刷屏
	ReportReasonProvoke         = int8(4)  // 引站
	ReportReasonSpoiler         = int8(5)  // 剧透
	ReportReasonPolitic         = int8(6)  // 政治
	ReportReasonAttack          = int8(7)  // 人身攻击
	ReportReasonUnrelated       = int8(8)  // 视频不相关
	ReportReasonProhibited      = int8(9)  // 违禁
	ReportReasonVulgar          = int8(10) // 低俗
	ReportReasonIllegalWebsite  = int8(11) // 非法网站
	ReportReasonGamblingFraud   = int8(12) // 赌博诈骗
	ReportReasonRumor           = int8(13) // 传播不实信息
	ReportReasonAbetting        = int8(14) // 怂恿教唆信息
	ReportReasonPrivacyInvasion = int8(15) // 侵犯隐私
	ReportReasonUnlimitedSign   = int8(16) // 抢楼

	ForbidReasonSpoiler         = int8(10) // 发布剧透信息
	ForbidReasonAd              = int8(6)  // 发布垃圾广告信息
	ForbidReasonUnlimitedSign   = int8(2)  // 抢楼
	ForbidReasonMeaningless     = int8(1)  // 刷屏
	ForbidReasonProvoke         = int8(9)  // 发布引战言论
	ForbidReasonVulgar          = int8(14) // 发布低俗信息
	ForbidReasonGamblingFraud   = int8(4)  // 发布赌博诈骗信息
	ForbidReasonPorn            = int8(13) // 发布色情信息
	ForbidReasonRumor           = int8(18) // 发布传播不实信息
	ForbidReasonIllegalWebsite  = int8(17) // 发布非法网站信息
	ForbidReasonAbetting        = int8(19) // 发布怂恿教唆信息
	ForbidReasonProhibited      = int8(5)  // 发布违禁信息
	ForbidReasonPrivacyInvasion = int8(8)  // 涉及侵犯他人隐私
	ForbidReasonAttack          = int8(7)  // 发布人身攻击言论

	PlatUnknow  = int8(0)
	PlatWeb     = int8(1)
	PlatAndroid = int8(2)
	PlatIPhone  = int8(3)
	PlatWpM     = int8(4) // wp mobile
	PlatIPad    = int8(5)
	PlatWpPc    = int8(6) // wp win10

	AdminOperDelete                 = int8(0)  // admin delete
	AdminOperDeleteByReport         = int8(1)  // admin delete by report
	AdminOperIgnoreReport           = int8(2)  // admin ignore report
	AdminOperRecover                = int8(3)  // admin recover
	AdminOperEdit                   = int8(4)  // admin edit reply content
	AdminOperPass                   = int8(5)  // admin pass
	AdminOperSubState               = int8(6)  // admin change subject state
	AdminOperSubTop                 = int8(7)  // top reply
	AdminOperSubMid                 = int8(8)  // admin change subject mid
	AdminOperRptIgnore1             = int8(9)  // admin report ignore 1
	AdminOperRptIgnore2             = int8(10) // admin report ignore 2
	AdminOperRptDel1                = int8(11) // admin report del 1
	AdminOperRptDel2                = int8(12) // admin report del 2
	AdminOperRptRecover1            = int8(13) // admin report recover 1
	AdminOperRptRecover2            = int8(14) // admin report recover 2
	AdminOperActionSet              = int8(15) // admin action set
	AdminOperDeleteUp               = int8(16) // admin delete by up
	AdminOperDeleteUser             = int8(17) // admin delete by user
	AdminOperDeleteAssist           = int8(18) // admin delete by assist
	AdminOperRptTransfer1           = int8(20) // admin transfer report to 1
	AdminOperRptTransfer2           = int8(21) // admin transfer report to 2
	AdminOperRptTransferArbitration = int8(22) // admin transfer report to arbitration
	AdminOperRptStateSet            = int8(23) // admin set report state

	AdminIsNotReport = int8(0)
	AdminIsReport    = int8(1)
	AdminIsNotNew    = int8(0)
	AdminIsNew       = int8(1)

	AuditTypeOne = int8(1) // 一审
	AuditTypeTwo = int8(2) // 二审
)

var (
	// NotifyComRules 社区规则
	NotifyComRules = fmt.Sprintf(`评论区是公众场所，而非私人场所，具体规范烦请参阅#{《社区规则》}{"%s"}，良好的社区氛围需要大家一起维护！`, "http://www.bilibili.com/blackboard/blackroom.html")
	// NotifyComRulesReport NotifyComRulesReport
	NotifyComRulesReport = "感谢您对bilibili社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~"
	// NotifyComUnrelated NotifyComUnrelated
	NotifyComUnrelated = "bilibili倡导发送与视频相关的评论，希望大家尊重作品，尊重UP主。良好的社区氛围需要大家一起维护！"
	// NotifyComProvoke NotifyComProvoke
	NotifyComProvoke = "bilibili倡导平等友善的交流。良好的社区氛围需要大家一起维护！"
	// NofityComProhibited NofityComProhibited
	NofityComProhibited = fmt.Sprintf(`请自觉遵守国家相关法律法规及#{《社区规则》}{"%s"}，bilibili良好的社区氛围需要大家一起维护！`, "http://www.bilibili.com/blackboard/blackroom.html")
	// ReportReason 举报理由类型
	ReportReason = map[int8]string{
		ReportReasonAd:              "内容涉及垃圾广告",
		ReportReasonPorn:            "内容涉及色情",
		ReportReasonMeaningless:     "刷屏",
		ReportReasonProvoke:         "内容涉及引战",
		ReportReasonSpoiler:         "内容涉及视频剧透",
		ReportReasonPolitic:         "内容涉及政治相关",
		ReportReasonAttack:          "内容涉及人身攻击",
		ReportReasonUnrelated:       "视频不相关",
		ReportReasonProhibited:      "内容涉及违禁相关",
		ReportReasonVulgar:          "内容涉及低俗信息",
		ReportReasonIllegalWebsite:  "内容涉及非法网站信息",
		ReportReasonGamblingFraud:   "内容涉及赌博诈骗信息",
		ReportReasonRumor:           "内容涉及传播不实信息",
		ReportReasonAbetting:        "内容不适宜",
		ReportReasonPrivacyInvasion: "内容涉及侵犯他人隐私",
		ReportReasonUnlimitedSign:   "抢楼",
	}
	// ForbidReason 封禁理由类型
	ForbidReason = map[int8]string{
		ForbidReasonSpoiler:         "发布剧透信息",
		ForbidReasonAd:              "发布垃圾广告信息",
		ForbidReasonUnlimitedSign:   "抢楼",
		ForbidReasonMeaningless:     "刷屏",
		ForbidReasonProvoke:         "发布引战言论",
		ForbidReasonVulgar:          "发布低俗信息",
		ForbidReasonGamblingFraud:   "发布赌博诈骗信息",
		ForbidReasonPorn:            "发布色情信息",
		ForbidReasonRumor:           "发布传播不实信息",
		ForbidReasonIllegalWebsite:  "发布非法网站信息",
		ForbidReasonAbetting:        "发布不适宜内容",
		ForbidReasonProhibited:      "发布违禁信息",
		ForbidReasonPrivacyInvasion: "涉及侵犯他人隐私",
		ForbidReasonAttack:          "发布人身攻击言论",
	}
)

// Subject is subject of reply
type Subject struct {
	ID     int64      `json:"-"`
	Oid    int64      `json:"oid"`
	Type   int8       `json:"type"`
	Mid    int64      `json:"mid"`
	Count  int        `json:"count"`
	RCount int        `json:"rcount"`
	ACount int        `json:"acount"`
	MCount int        `json:"mcount"`
	State  int8       `json:"state"`
	Attr   uint32     `json:"attr"`
	Meta   string     `json:"meta"`
	CTime  xtime.Time `json:"ctime"`
	MTime  xtime.Time `json:"-"`
}

// HasFolded ...
func (s *Subject) HasFolded() bool {
	return s.AttrVal(SubAttrFolded) == AttrYes
}

// UnmarkHasFolded ...
func (s *Subject) UnmarkHasFolded() {
	s.AttrSet(AttrNo, SubAttrFolded)
}

// SubjectMeta SubjectMeta
type SubjectMeta struct {
	AdminTop int64 `json:"atop,omitempty"`
	UpperTop int64 `json:"utop,omitempty"`
}

// TopCount TopCount
func (s *Subject) TopCount() int {
	return int(s.AttrVal(SubAttrUpperTop) + s.AttrVal(SubAttrAdminTop))
}

// AttrVal return val of subject'attr
func (s *Subject) AttrVal(bit uint32) uint32 {
	if s.Attr == 0 {
		return uint32(0)
	}
	return (s.Attr >> bit) & uint32(1)
}

// TopSet TopSet
func (s *Subject) TopSet(top int64, typ uint32, act uint32) (err error) {
	var meta SubjectMeta
	if s.Meta != "" {
		err = json.Unmarshal([]byte(s.Meta), &meta)
		if err != nil {
			return
		}
	}
	if act == 1 {
		if typ == 0 {
			if meta.AdminTop == top {
				return fmt.Errorf("subject already have the same top")
			}
			meta.AdminTop = top
		} else {
			if meta.UpperTop == top {
				return fmt.Errorf("subject already have the same top")
			}
			meta.UpperTop = top
		}
	} else {
		if typ == 0 {
			meta.AdminTop = 0
		} else {
			meta.UpperTop = 0
		}
	}
	var content []byte
	content, err = json.Marshal(meta)
	if err != nil {
		return err
	}
	s.Meta = string(content)
	return
}

// AttrSet set val of subject'attr
func (s *Subject) AttrSet(v uint32, bit uint32) {
	s.Attr = s.Attr&(^(1 << bit)) | (v << bit)
}

// IsNormal check if reply subject normal
func (s *Subject) IsNormal() bool {
	return s.State == SubStateNormal
}

// IsAudit check reply subject is in audit
func (s *Subject) IsAudit() bool {
	return s.State == SubStateReplyAfterAudit
}

// CheckSubState check
func CheckSubState(state int8) (err error) {
	if state < SubStateNormal || state > SubStateReplyAfterAudit {
		err = ecode.ReplyIllegalSubState
	}
	return
}

// CheckSubForbid check if subject forbiden reply
func CheckSubForbid(state int8) (err error) {
	if state == SubStateForbid {
		err = ecode.ReplyForbidReply
	}
	return
}

// Reply define reply object
type Reply struct {
	RpID   int64      `json:"rpid"`
	Oid    int64      `json:"oid"`
	Type   int8       `json:"type"`
	Mid    int64      `json:"mid"`
	Root   int64      `json:"root"`
	Parent int64      `json:"parent"`
	Dialog int64      `json:"dialog"`
	Count  int        `json:"count"`
	RCount int        `json:"rcount"`
	Floor  int        `json:"floor"`
	State  int8       `json:"state"`
	Attr   uint32     `json:"attr"`
	CTime  xtime.Time `json:"ctime"`
	MTime  xtime.Time `json:"-"`
	// string
	RpIDStr   string `json:"rpid_str,omitempty"`
	RootStr   string `json:"root_str,omitempty"`
	ParentStr string `json:"parent_str,omitempty"`
	// action count, from ReplyAction count
	Like   int  `json:"like"`
	Hate   int  `json:"hate"`
	Action int8 `json:"action"`
	// member info
	Member *accmdl.Info `json:"member"`
	// other
	Content *Content `json:"content"`
	Replies []*Reply `json:"replies"`
}

// IsFolded ...
func (r *Reply) IsFolded() bool {
	return r.State == ReplyStateFolded
}

// HasFolded ...
func (r *Reply) HasFolded() bool {
	return r.AttrVal(ReplyAttrFolded) == AttrYes
}

// UnmarkHasFolded ...
func (r *Reply) UnmarkHasFolded() {
	r.AttrSet(AttrNo, ReplyAttrFolded)
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

// IsNormal check if reply normal
func (r *Reply) IsNormal() bool {
	return r.State == ReplyStateNormal || r.State == ReplyStateHidden || r.State == ReplyStateFiltered || r.State == ReplyStateMonitor || r.State == ReplyStateGarbage || r.State == ReplyStateTop || r.State == ReplyStateFolded
}

// IsDeleted deleted.
func (r *Reply) IsDeleted() bool {
	return r.State == ReplyStateUserDel || r.State == ReplyStateUpDel || r.State == ReplyStateAdminDel || r.State == ReplyStateAssistDel
}

// IsRoot IsRoot
func (r *Reply) IsRoot() bool {
	return r.Root == 0 && r.Parent == 0
}

// IsTop top.
func (r *Reply) IsTop() bool {
	if r.Attr != 0 && (r.AttrVal(ReplyAttrAdminTop) == 1 || r.AttrVal(ReplyAttrUpperTop) == 1) {
		return true
	}
	return false
}

// IsAdminTop admin top.
func (r *Reply) IsAdminTop() bool {
	return r.AttrVal(ReplyAttrAdminTop) == 1
}

// IsUpTop up top.
func (r *Reply) IsUpTop() bool {
	return r.AttrVal(ReplyAttrUpperTop) == 1
}

// FillStr fill reply string info
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
	}
}

// CheckSort check sort type
func CheckSort(sort int8) bool {
	return SortByFloor <= sort && sort <= SortByLike
}

// CheckPlat check plat type
func CheckPlat(plat int8) bool {
	return PlatUnknow <= plat && plat <= PlatWpPc
}

// Content define reply content
type Content struct {
	RpID    int64      `json:"-"`
	Message string     `json:"message"`
	Ats     Int64Bytes `json:"ats,omitempty"`
	Topics  Mstr       `json:"topics,omitempty"`
	IP      uint32     `json:"ipi,omitempty"`
	Plat    int8       `json:"plat"`
	Device  string     `json:"device"`
	Version string     `json:"version,omitempty"`
	CTime   xtime.Time `json:"-"`
	MTime   xtime.Time `json:"-"`
	// ats member info
	// Members []*accmdl.Info `json:"members"`
}

// // FillAts fill user ifo of @
// func (rc *ReplyContent) FillAts(mis map[int64]*accmdl.Info) {
// 	rc.Members = make([]*accmdl.Info, 0, len(rc.Ats))
// 	for _, at := range rc.Ats {
// 		if mi, ok := mis[at]; ok {
// 			rc.Members = append(rc.Members, mi)
// 		}
// 	}
// 	rc.Ats = nil
// }

// Action reply action info
type Action struct {
	ID     int64      `json:"-"`
	RpID   int64      `json:"rpid"`
	Action int8       `json:"action"`
	Mid    int64      `json:"mid"`
	CTime  xtime.Time `json:"-"`
}

// CheckAction check action operation
func CheckAction(act int8) (err error) {
	if act != OpAdd && act != OpCancel {
		err = ecode.ReplyIllegalAction
	}
	return
}

// Report define reply report
type Report struct {
	ID      int64      `json:"id,omitempty"`
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

// IsTransferred AttrVal
func (rpt Report) IsTransferred() bool {
	return rpt.AttrVal(ReportAttrTransferred) == 1

}

// AttrSet set attr of ReplyReport'attr
func (rpt *Report) AttrSet(v int8, bit uint) {
	rpt.Attr = rpt.Attr&(^(1 << bit)) | (v << bit)
}

// SetTransferred SetTransferred
func (rpt *Report) SetTransferred() {
	rpt.AttrSet(1, 0)
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

// CheckReportReason  check if reprot reason illegal
func CheckReportReason(reason int8) (err error) {
	if !(ReportReasonOther <= reason && reason <= ReportReasonUnrelated) {
		err = ecode.ReplyIllegalReport
	}
	return
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

// Business Business
type Business struct {
	Type  int32  `json:"type"`
	Alias string `json:"alias"`
}

// RpItem fix dialog model
type RpItem struct {
	ID     int64
	Parent int64
	Floor  int
	Next   *RpItem
}

type RpItems []*RpItem

func (c RpItems) Len() int {
	return len(c)
}

func (c RpItems) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c RpItems) Less(i, j int) bool {
	return c[i].ID < c[j].ID
}
