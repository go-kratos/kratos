package model

import (
	"fmt"
	"time"
)

// const const var
const (
	// up主操作
	StatUpperInit   = int8(0) // up主未处理
	StatUpperIgnore = int8(1) // up主已忽略
	StatUpperDelete = int8(2) // up主已删除

	// 管理员操作
	StatFirstInit        = int8(0) // 待一审
	StatFirstDelete      = int8(1) // 一审删除
	StatSecondInit       = int8(2) // 待二审
	StatSecondDelete     = int8(3) // 二审删除
	StatSecondIgnore     = int8(4) // 二审忽略
	StatFirstIgnore      = int8(5) // 一审忽略
	StatSecondAutoDelete = int8(6) // 二审脚本删除
	StatJudgeInit        = int8(7) // 风纪委待审(二审)
	StatJudgeDelete      = int8(8) // 风纪委删除(二审)
	StatJudgeIgnore      = int8(9) // 风纪委忽略(二审)

	// 处理结果通知
	NoticeUnsend = int8(0) // 未通知用户
	NoticeSend   = int8(1) // 已通知用户

	// 举报通知状态
	NoticeReporter = int8(1)
	NoticePoster   = int8(2)
	NoticeAll      = int8(3)

	// 举报原因
	ReportReasonProhibited  = int8(1)  // 违禁
	ReportReasonPorn        = int8(2)  // 色情
	RptReasonFraud          = int8(3)  // 赌博诈骗
	ReportReasonAttack      = int8(4)  // 人身攻击
	ReportReasonPrivate     = int8(5)  // 隐私
	ReportReasonAd          = int8(6)  // 广告
	ReportReasonWar         = int8(7)  // 引战
	ReportReasonSpoiler     = int8(8)  // 剧透
	ReportReasonMeaningless = int8(9)  // 刷屏
	ReportReasonUnrelated   = int8(10) // 视频不相关
	ReportReasonOther       = int8(11) // 其他
	ReportReasonTeenagers   = int8(12) // 青少年不良信息
)

// var const map
var (
	RptTemplate = map[string]string{
		"del":    `您好，您在视频#{%s}{"http://www.bilibili.com/av%d"}中举报的弹幕『%s』已被删除%s，原因是『%s』，感谢您对bilibili社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~`,
		"ignore": `您好，非常感谢您的举报，您在视频#{%s}{"http://www.bilibili.com/av%d"}中举报的弹幕『%s』暂未认定为违规内容，具体弹幕规范烦请参阅 #{《弹幕礼仪》}{"http://www.bilibili.com/blackboard/help.html#d5"}，哔哩哔哩 (゜-゜)つロ 干杯~`,
	}
	PostTemplate = map[int8]string{
		1: `您好，根据用户举报，您在视频#{%s}{"http://www.bilibili.com/av%d"}中的弹幕『%s』已被删除%s，原因是『%s』，请自觉遵守国家相关法律法规，具体弹幕规范烦请参阅#{《弹幕礼仪》}{"http://www.bilibili.com/blackboard/help.html#d5"}，bilibili良好的社区氛围需要大家一起维护！`,
		2: `您好，根据用户举报，您在视频#{%s}{"http://www.bilibili.com/av%d"}中的弹幕『%s』已被删除%s，原因是『%s』，bilibili倡导平等友善的交流，具体弹幕规范烦请参阅#{《弹幕礼仪》}{"http://www.bilibili.com/blackboard/help.html#d5"}，良好的社区氛围需要大家一起维护！`,
		3: `您好，根据用户举报，您在视频#{%s}{"http://www.bilibili.com/av%d"}中的弹幕『%s』已被删除%s，原因是『%s』，弹幕是公众场所而非私人场所，具体弹幕规范烦请参阅#{《弹幕礼仪》}{"http://www.bilibili.com/blackboard/help.html#d5"}，良好的社区氛围需要大家一起维护！`,
		4: `您好，根据用户举报，您在视频#{%s}{"http://www.bilibili.com/av%d"}中的弹幕『%s』已被删除%s，原因是『%s』，bilibili倡导发送与视频相关、有用的弹幕，具体弹幕规范烦请参阅#{《弹幕礼仪》}{"http://www.bilibili.com/blackboard/help.html#d5"}，良好的社区氛围需要大家一起维护！`,
	}
	AdminRptReason = map[int8]string{
		1:  "内容涉及传播不实信息",
		2:  "内容涉及非法网站信息",
		3:  "内容涉及怂恿教唆信息",
		4:  "内容涉及低俗信息",
		5:  "内容涉及色情",
		6:  "内容涉及赌博诈骗信息",
		7:  "内容涉及人身攻击",
		8:  "内容涉及侵犯他人隐私",
		9:  "内容涉及垃圾广告",
		10: "内容涉及引战",
		11: "内容涉及视频剧透",
		12: "恶意刷屏",
		13: "视频不相关",
		14: "其他",
		15: "内容涉及违禁相关",
		16: "内容不适宜",
		17: "内容涉及青少年不良信息",
	}
	BlockReason = map[int8]string{
		4:  "发布赌博诈骗信息",
		5:  "发布违禁相关信息",
		6:  "发布垃圾广告信息",
		7:  "发布人身攻击言论",
		8:  "发布侵犯他人隐私信息",
		9:  "发布引战言论",
		10: "发布剧透信息",
		13: "发布色情信息",
		14: "发布低俗信息",
		17: "发布非法网站信息",
		18: "发布传播不实信息",
		19: "发布怂恿教唆信息",
		20: "恶意刷屏",
		24: "发布青少年不良内容",
	}
)

// ReportListParams .
type ReportListParams struct {
	States   []int64 `form:"state,split"`
	UpOps    []int64 `form:"upop,split"`
	Tids     []int64 `form:"tid,split"`
	Aid      int64   `form:"aid"`
	Cid      int64   `form:"cid"`
	UID      int64   `form:"uid"`
	RpUID    int64   `form:"rp_user"`
	RpTypes  []int64 `form:"rp_type,split"`
	Start    string  `form:"start"`
	End      string  `form:"end"`
	Sort     string  `form:"sort"`
	Order    string  `form:"order"`
	Keyword  string  `form:"keyword"`
	Page     int32   `form:"page" default:"1"`
	PageSize int32   `form:"page_size" default:"100" validate:"max=1000"`
}

// Report dm report struct.
type Report struct {
	DidStr   string        `json:"dmid_str"` // str id
	ID       int64         `json:"id"`
	Did      int64         `json:"dmid"`         // 弹幕id
	Cid      int64         `json:"cid"`          // 视频的id
	Aid      int64         `json:"arc_aid"`      // 稿件的id
	Tid      int64         `json:"arc_typeid"`   // 稿件的分区id
	UID      int64         `json:"dm_owner_uid"` // 弹幕发送者的uid
	Msg      string        `json:"dm_msg"`       // 弹幕内容
	Count    int64         `json:"count"`        // 被举报次数
	Content  string        `json:"content"`      // 举报内容:只有类别其他才有值
	UpOP     int8          `json:"up_op"`        // up主操作状态
	State    int8          `json:"state"`        // 举报状态
	RpUID    int64         `json:"uid"`          // 最后一个举报用户id
	RpTime   string        `json:"rp_time"`      // 举报时间
	RpType   int64         `json:"reason"`       // 举报类型
	Title    string        `json:"arc_title"`    // 稿件标题
	Deleted  int64         `json:"dm_deleted"`   // 弹幕状态
	UPUid    int64         `json:"arc_mid"`      // up主id
	PoolID   int64         `json:"pool_id"`      // 弹幕池
	Model    int64         `json:"model"`        // 弹幕model
	Score    int32         `json:"score"`        // 举报分
	SendTime string        `json:"dm_ctime"`     // 弹幕发送时间
	Ctime    string        `json:"ctime"`        // 插入时间
	Mtime    string        `json:"mtime"`        // 更新时间
	RptUsers []*ReportUser `json:"user"`         // 举报用户列表
}

// ReportMsg report message
type ReportMsg struct {
	Aid         int64
	Uids        string
	Did         int64
	Title       string
	Msg         string
	State       int8
	RptReason   int8
	BlockReason int8
	Block       int64
}

// ReportJudge report judge
type ReportJudge struct {
	AID        int64  `json:"aid"`
	MID        int64  `json:"mid"`
	Operator   string `json:"operator"`
	OperID     int64  `json:"oper_id"`
	OContent   string `json:"origin_content"`
	OTitle     string `json:"origin_title"`
	OType      int64  `json:"origin_type"`
	OURL       string `json:"origin_url"`
	ReasonType int64  `json:"reason_type"`
	OID        int64  `json:"oid"`
	RPID       int64  `json:"rp_id"`
	TagID      int64  `json:"tag_id"`
	Type       int64  `json:"type"`
	Page       int64  `json:"page"`
	BTime      int64  `json:"business_time"`
}

// SearchReportResult dm repost list from search
type SearchReportResult struct {
	Code  int64  `json:"code"`
	Order string `json:"order"`
	Sort  string `json:"sort"`
	Page  *struct {
		Num   int64 `json:"num"`
		Size  int64 `json:"size"`
		Total int64 `json:"total"`
	} `json:"page"`
	Result []*Report `json:"result"`
}

// UptSearchReport update search report
type UptSearchReport struct {
	DMid  int64  `json:"dmid"`
	State int8   `json:"state"`
	Ctime string `json:"ctime"`
	Mtime string `json:"mtime"`
}

// ReportList dm report list
type ReportList struct {
	Code      int64     `json:"code"`
	Order     string    `json:"order"`
	Page      int64     `json:"page"`
	PageSize  int64     `json:"pagesize"`
	PageCount int64     `json:"pagecount"`
	Total     int64     `json:"total"`
	Result    []*Report `json:"result"`
}

// ReduceMoral reduce moral
type ReduceMoral struct {
	UID        int64
	Moral      int64
	Origin     int8
	Reason     int8
	ReasonType int8
	Operator   string
	IsNotify   int8
	Remark     string
}

// BlockUser block user
type BlockUser struct {
	UID             int64
	BlockForever    int64
	BlockTimeLength int64
	BlockRemark     string
	Operator        string
	OriginType      int64
	Moral           int64
	ReasonType      int64
	OriginTitle     string
	OriginContent   string
	OriginURL       string
	IsNotify        int64
}

// ReportUser report user
type ReportUser struct {
	ID     int64     `json:"id"`
	Did    int64     `json:"dmid"`
	UID    int64     `json:"uid"`
	Reason int64     `json:"reason"`
	State  int8      `json:"state"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}

// ReportLog dm admin log
type ReportLog struct {
	ID      int64     `json:"id"`
	Did     int64     `json:"dmid"`
	AdminID int64     `json:"admin_id"`
	Reason  int8      `json:"reason"`
	Result  int8      `json:"result"`
	Remark  string    `json:"remark"`
	Elapsed int64     `json:"elapsed"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// PosterBlockMsg return report msg template by blockReason
func PosterBlockMsg(reason int8) (str string, err error) {
	var (
		tmplKey int8
	)
	switch reason {
	case 4, 5, 13, 14, 17, 18, 19, 20, 24:
		tmplKey = 1
	case 6, 9:
		tmplKey = 2
	case 7, 8, 10, 12:
		tmplKey = 3
	default:
		err = fmt.Errorf("BlockReason %d not exist", reason)
		return
	}
	str = PostTemplate[tmplKey]
	return
}

// PosterAdminRptMsg return report msg template by adminRptReason
func PosterAdminRptMsg(reason int8) (str string, err error) {
	var (
		tmplKey int8
	)
	switch reason {
	case 1, 2, 3, 4, 5, 6, 15, 17:
		tmplKey = 1
	case 7, 10:
		tmplKey = 2
	case 8, 9, 11, 12:
		tmplKey = 3
	case 13, 14, 16:
		tmplKey = 4
	default:
		err = fmt.Errorf("adminRptReason %d not exist", reason)
		return
	}
	str = PostTemplate[tmplKey]
	return
}

// RpReasonToJudgeReason 修改弹幕风纪委的理由
func RpReasonToJudgeReason(r int8) (j int8) {
	switch r {
	case ReportReasonProhibited:
		j = 5
	case ReportReasonPorn:
		j = 13
	case RptReasonFraud:
		j = 4
	case ReportReasonAttack:
		j = 7
	case ReportReasonPrivate:
		j = 8
	case ReportReasonAd:
		j = 6
	case ReportReasonWar:
		j = 9
	case ReportReasonSpoiler:
		j = 10
	case ReportReasonMeaningless:
		j = 20
	}
	return
}

// CheckStateBelong check state first or second check
func CheckStateBelong(state int8) string {
	if state == StatFirstInit || state == StatFirstDelete || state == StatFirstIgnore {
		return "弹幕举报一审"
	}
	return "弹幕举报二审"
}
