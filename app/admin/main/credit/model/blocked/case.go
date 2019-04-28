package blocked

import (
	"strconv"

	xtime "go-common/library/time"
)

// const case
const (
	// JudeCaseTypePrivate case type private
	JudeCaseTypePrivate = int8(0) // 小众众裁
	// JudeCaseTypePublic case type public
	JudeCaseTypePublic = int8(1) // 大众众裁
	// case status.
	CaseStatusGranting  = int8(1) // 发放中
	CaseStatusGrantStop = int8(2) // 停止发放
	CaseStatusDealing   = int8(3) // 结案中
	CaseStatusDealed    = int8(4) // 已裁决
	CaseStatusRestart   = int8(5) // 待重启
	CaseStatusUndealed  = int8(6) // 未裁决
	CaseStatusFreeze    = int8(7) // 冻结中
	CaseStatusQueueing  = int8(8) // 队列中
	// vote status.
	VoteTypeUndo    = int8(0) // 未投票
	VoteTypeViolate = int8(1) // 违规-封禁
	VoteTypeLegal   = int8(2) // 不违规
	VoteTypeGiveUp  = int8(3) // 放弃投票
	VoteTypeDelete  = int8(4) // 违规-删除
	// origin_type.
	OriginReply    = int8(1)  // 评论
	OriginDM       = int8(2)  // 弹幕
	OriginMsg      = int8(3)  // 私信
	OriginTag      = int8(4)  // 标签
	OriginMember   = int8(5)  // 个人资料
	OriginArchive  = int8(6)  // 投稿
	OriginMusic    = int8(7)  // 音频
	OriginArticle  = int8(8)  // 专栏
	OriginSpaceTop = int8(9)  // 空间头图
	OriginDsynamic = int8(10) // 动态
	OriginPhoto    = int8(11) // 相册
	OriginMinVideo = int8(12) // 小视频
	// punish status
	PunishhNon           = int8(0)
	PunishBlockedTree    = int8(1)
	PunishBlockedSeven   = int8(2)
	PunishBlockedEver    = int8(3)
	PunishBlockedCustom  = int8(4)
	PunishBlockedFifteen = int8(5)
	PunishJustDelete     = int8(6)

	// block time.
	BlockMoralNum = -2 // 扣除节操
	BlockCustom   = -1 // N天封禁
	BlockForever  = 0  // 永久封禁
	BlockThree    = 3  // 3天封禁
	BlockSeven    = 7  // 7天封禁
	BlockFifteen  = 15 // 15天封禁
)

// var case
var (
	StatusDesc = map[int8]string{
		CaseStatusGranting:  "发放中",
		CaseStatusGrantStop: "停止发放",
		CaseStatusDealing:   "结案中",
		CaseStatusDealed:    "已裁决",
		CaseStatusRestart:   "待重启",
		CaseStatusUndealed:  "未裁决",
		CaseStatusFreeze:    "冻结中",
		CaseStatusQueueing:  "队列中",
	}
	OriginTypeDesc = map[int8]string{
		OriginReply:    "评论",
		OriginDM:       "弹幕",
		OriginMsg:      "私信",
		OriginTag:      "标签",
		OriginMember:   "个人资料",
		OriginArchive:  "投稿",
		OriginMusic:    "音频",
		OriginArticle:  "专栏",
		OriginSpaceTop: "空间头图",
		OriginDsynamic: "动态",
		OriginPhoto:    "相册",
		OriginMinVideo: "小视频",
	}
	PunishDesc = map[int8]string{
		PunishhNon:           "不违规不封禁",
		PunishBlockedTree:    "封禁3天",
		PunishBlockedSeven:   "封禁7天",
		PunishBlockedEver:    "永久封禁",
		PunishBlockedCustom:  "自定义封禁",
		PunishBlockedFifteen: "封禁15天",
		PunishJustDelete:     "违规仅删除",
	}
	blockedDesc = map[int]string{
		BlockMoralNum: "扣除%d节操",
		BlockCustom:   "天封禁",
		BlockForever:  "永久封禁",
		BlockThree:    "3天封禁",
		BlockSeven:    "7天封禁",
		BlockFifteen:  "15天封禁",
	}
	CaseTypeDesc = map[int8]string{
		JudeCaseTypePrivate: "非公开众裁",
		JudeCaseTypePublic:  "公开众裁",
	}
)

// Case is blocked_case model.
type Case struct {
	ID             int64      `gorm:"column:id" json:"id"`
	MID            int64      `gorm:"column:mid" json:"uid"`
	OPID           int64      `gorm:"column:oper_id" json:"oper_id"`
	Status         int8       `gorm:"column:status" json:"status"`
	OriginType     int8       `gorm:"column:origin_type" json:"origin_type"`
	ReasonType     int8       `gorm:"column:reason_type" json:"reason_type"`
	PunishResult   int8       `gorm:"column:punish_result" json:"punish_result"`
	JudgeType      int        `gorm:"column:judge_type" json:"judge_type"`
	CaseType       int8       `gorm:"column:case_type" json:"case_type"`
	BlockedDays    int        `gorm:"column:blocked_days" json:"blocked_days"`
	PutTotal       int        `gorm:"column:put_total" json:"put_total"`
	VoteRule       int64      `gorm:"column:vote_rule" json:"vote_rule"`
	VoteBreak      int64      `gorm:"column:vote_break" json:"vote_break"`
	VoteDelete     int64      `gorm:"column:vote_delete" json:"vote_delete"`
	VoteTotal      int64      `gorm:"-" json:"vote_total"` // 总得票数
	StartTime      xtime.Time `gorm:"column:start_time" json:"start_time"`
	EndTime        xtime.Time `gorm:"column:end_time" json:"end_time"`
	CTime          xtime.Time `gorm:"column:ctime" json:"ctime"`
	MTime          xtime.Time `gorm:"column:mtime" json:"-"`
	OriginURL      string     `gorm:"column:origin_url" json:"origin_url"`
	OriginTitle    string     `gorm:"column:origin_title" json:"origin_title"`
	OriginContent  string     `gorm:"column:origin_content" json:"origin_content"`
	BusinessTime   xtime.Time `gorm:"column:business_time" json:"business_time"`
	Uname          string     `gorm:"-" json:"uname"`
	StatusDesc     string     `gorm:"-" json:"status_desc"`
	OriginTypeDesc string     `gorm:"-" json:"origin_type_desc"`
	ReasonTypeDesc string     `gorm:"-" json:"reason_type_desc"`
	CaseTypeDesc   string     `gorm:"-" json:"case_type_desc"`
	RelationID     string     `gorm:"column:relation_id" json:"relation_id"`
	RulePercent    string     `gorm:"-" json:"rule_percent"`    // 不违规得票率
	BlockedPercent string     `gorm:"-" json:"blocked_percent"` // 违规（封禁）得票率
	DeletePercent  string     `gorm:"-" json:"delete_percent"`  // 违规（仅删）得票率
	PunishDesc     string     `gorm:"-" json:"punish_desc"`
	OPName         string     `gorm:"-" json:"oname"` // 操作人
	Fans           int64      `gorm:"-" json:"fans"`  // 粉丝数
}

// CaseVote is blocked_case_vote model.
type CaseVote struct {
	ID       int64      `gorm:"column:id" json:"id"`
	CID      int64      `gorm:"column:cid" json:"cid"`
	UID      int64      `gorm:"column:mid" json:"uid"`
	VoteType int8       `gorm:"column:vote_type" json:"vote_type"`
	Expired  xtime.Time `gorm:"column:expired" json:"expired"`
	CTime    xtime.Time `gorm:"column:ctime" json:"-"`
	MTime    xtime.Time `gorm:"column:mtime" json:"-"`
}

// TableName case tablename
func (*Case) TableName() string {
	return "blocked_case"
}

// TableName CaseVote tablename
func (*CaseVote) TableName() string {
	return "blocked_case_vote"
}

// CaseList is case list.
type CaseList struct {
	Count int
	Order string
	Sort  string
	PN    int
	PS    int
	IDs   []int64
	List  []*Case
}

// RulePercent is  rule_percent.
func RulePercent(voteRule, voteBreak, voteDelete int64) string {
	return strconv.FormatFloat(float64(voteRule)/float64(voteRule+voteBreak+voteDelete), 'f', 2, 64)
}

// BreakPercent is blocked percent.
func BreakPercent(voteRule, voteBreak, voteDelete int64) string {
	return strconv.FormatFloat(float64(voteBreak)/float64(voteRule+voteBreak+voteDelete), 'f', 2, 64)
}

// DeletePercent is  delete_percent.
func DeletePercent(voteRule, voteBreak, voteDelete int64) string {
	return strconv.FormatFloat(float64(voteDelete)/float64(voteRule+voteBreak+voteDelete), 'f', 2, 64)
}

// VoteTotal is  vote_total.
func VoteTotal(voteRule, voteBreak, voteDelete int64) int64 {
	return voteRule + voteBreak + voteDelete
}
