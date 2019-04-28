package model

import (
	"time"
)

// var const text
var (
	ReportReason = map[int8]string{
		1:  "内容涉及违禁相关",
		2:  "内容涉及非法网站信息",
		3:  "内容涉及赌博诈骗信息",
		4:  "内容涉及人身攻击",
		5:  "内容涉及侵犯他人隐私",
		6:  "内容涉及垃圾广告",
		7:  "内容涉及引战",
		8:  "内容涉及视频剧透",
		9:  "恶意刷屏",
		10: "视频不相关",
		11: "其他",
		12: "青少年不良信息",
	}
	RptMsgTitle    = "举报处理结果通知"
	RptMsgTemplate = `您好，您在视频#{%s}{"http://www.bilibili.com/av%d"}中举报的弹幕『%s』已被删除，原因是『%s』，感谢您对bilibili社区秩序的维护，哔哩哔哩 (゜-゜)つロ 干杯~`
)

// const var
const (
	// up主操作
	StatUpperInit   = int8(0) // up主未处理
	StatUpperIgnore = int8(1) // up主已忽略
	StatUpperDelete = int8(2) // up主已删除

	// 管理员操作
	StatFirstInit        = int8(0) // 待一审
	StatFirstDelete      = int8(1) // 一审删除
	StatSecondInit       = int8(2) // 待二审
	StatSecondIgnore     = int8(3) // 二审忽略
	StatSecondDelete     = int8(4) // 二审删除
	StatFirstIgnore      = int8(5) // 一审忽略
	StatSecondAutoDelete = int8(6) // 二审脚本删除
	// 处理结果通知
	NoticeUnsend = int8(0) // 未通知用户
	NoticeSend   = int8(1) // 已通知用户

	// 举报原因
	ReportReasonProhibited  = int8(1)  // 违禁
	ReportReasonPorn        = int8(2)  // 色情
	ReportReasonFraud       = int8(3)  // 赌博诈骗
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

// Report dm report info
type Report struct {
	ID      int64     `json:"id"`      // 主键id
	Cid     int64     `json:"cid"`     // 视频id
	Did     int64     `json:"dmid"`    // 弹幕id
	UID     int64     `json:"uid"`     // 举报用户的id
	Reason  int8      `json:"reason"`  // 举报原因类型
	Content string    `json:"content"` // 举报内容：reason为其它时有值
	Count   int64     `json:"count"`   // 被举报次数
	State   int8      `json:"state"`   // 举报状态
	UpOP    int8      `json:"up_op"`   // up主操作
	Score   int32     `json:"score"`   // 举报分
	RpTime  time.Time `json:"rp_time"` // 举报时间
	Ctime   time.Time `json:"ctime"`   // 插入时间
	Mtime   time.Time `json:"mtime"`   // 更新时间
}

// User report user info
type User struct {
	ID      int64     `json:"id"`
	Did     int64     `json:"dmid"`
	UID     int64     `json:"uid"`
	Reason  int8      `json:"reason"`
	State   int8      `json:"state"`
	Content string    `json:"content"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// RptLog dm admin log
type RptLog struct {
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

// RptSearch report info from search
type RptSearch struct {
	ID      int64  `json:"id"`
	Cid     int64  `json:"cid"`          // 视频的id
	Did     int64  `json:"dmid"`         // 弹幕id
	Aid     int64  `json:"arc_aid"`      // 稿件的id
	Tid     int64  `json:"arc_typeid"`   // 稿件的分区id
	Owner   int64  `json:"dm_owner_uid"` // 弹幕发送者的uid
	Msg     string `json:"dm_msg"`       // 弹幕内容
	Count   int64  `json:"count"`        // 被举报次数
	Content string `json:"content"`      // 举报内容:只有类别其他才有值
	UpOP    int8   `json:"up_op"`        // up主操作状态
	State   int8   `json:"state"`        // 举报状态
	UID     int64  `json:"uid"`          // 举报用户id
	RpTime  string `json:"rp_time"`      // 举报时间
	Reason  int64  `json:"reason"`       // 举报原因类型
	Ctime   string `json:"ctime"`        // 插入时间
	Mtime   string `json:"mtime"`        // 更新时间
	Title   string `json:"arc_title"`    // 稿件标题
	Deleted int64  `json:"dm_deleted"`   // 弹幕状态
	UPUid   int64  `json:"arc_mid"`      // up主id
	Cover   string `json:"arc_cover"`    // 稿件的封面图
}

// RptSearchs report list
type RptSearchs struct {
	Page      int64        `json:"page"`
	PageSize  int64        `json:"pagesize"`
	PageCount int64        `json:"pagecount"`
	Total     int64        `json:"total"`
	Result    []*RptSearch `json:"result"`
}

// UptSearchReport update search report
type UptSearchReport struct {
	DMid  int64  `json:"dmid"`
	Upop  int8   `json:"up_op"`
	Ctime string `json:"ctime"`
	Mtime string `json:"mtime"`
}

// Page search page
type Page struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// SearchReportResult dm repost list from search
type SearchReportResult struct {
	Page   *Page        `json:"page"`
	Result []*RptSearch `json:"result"`
}

// SearchReportAidResult dm repost archive list from search
type SearchReportAidResult struct {
	Page   *Page `json:"page"`
	Result map[string][]struct {
		Key string `json:"key"`
	} `json:"result"`
}

// RptMsg dm report message
type RptMsg struct {
	Aid    int64
	UID    int64
	Did    int64
	Title  string
	Msg    string
	State  int8
	Reason int8
}

// Archives report archive list
type Archives struct {
	Result []*struct {
		Aid   int64  `json:"aid"`
		Title string `json:"title"`
	} `json:"result"`
}

// ReportAction send dm info and hidetime
type ReportAction struct {
	Cid      int64 `json:"cid"`       // 视频id
	Did      int64 `json:"dmid"`      // 弹幕id
	HideTime int64 `json:"hide_time"` // 弹幕隐藏截止时间
}
