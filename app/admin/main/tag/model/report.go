package model

import "go-common/library/time"

// const const value.
const (

	// DataType tag
	DataType = int32(4)

	// ActionDel 删除
	ActionDel = int32(0)
	// ActionAdd 增加
	ActionAdd = int32(1)

	// AuditFirst 一审
	AuditFirst = int32(1)
	// AuditSecond 二审
	AuditSecond = int32(2)

	// AuditNotHanleFirst 一审未处理
	AuditNotHanleFirst = int32(0)
	// AuditHanledFirst 一审已处理
	AuditHanledFirst = int32(4)
	// AuditHanledSecond 二审已处理
	AuditHanledSecond = int32(1)
	// AuditNotHanleSecond 二审未处理
	AuditNotHanleSecond = int32(3)
	// AuditNotDealSecond 二审不处理
	AuditNotDealSecond = int32(2)

	RptUserIsFirst = int32(1) // 第一举报者
	RptUserIsMng   = int32(1) // 举报者是管理员

	HandleNull             = int32(-2)  // 空状态
	HandleWait             = int32(-1)  // 等待审核
	HandleIntegral         = int32(2)   // 扣节操
	HandleBlock            = int32(3)   // 封禁用户
	HandleDelFirst         = int32(5)   // 一审删除
	HandleDelSecond        = int32(0)   // 二审删除
	HandleAddFirst         = int32(6)   // 一审添加
	HandleAddSecond        = int32(1)   // 二审添加
	HandleIgnoreFirst      = int32(7)   // 一审忽略(待二审)
	HandleIgnoreSecond     = int32(4)   // 二审忽略
	HandleAddUserFirst     = int32(8)   // 一审添加(用户)
	HandleAddUserSecond    = int32(10)  // 二审添加(用户)
	HandleDelUserFirst     = int32(9)   // 一审删除(用户)
	HandleDelUserSecond    = int32(11)  // 二审删除(用户)
	HandleRestoreDelFirst  = int32(12)  // 一审恢复(删除)
	HandleRestoreDelSecond = int32(14)  // 二审恢复(删除)
	HandleRestoreAddFirst  = int32(13)  // 一审恢复(添加)
	HandleRestoreAddSecond = int32(15)  // 二审恢复(添加)
	HandleCommission       = int32(100) // 移交众裁

	// MoralNotDeducted 节操尚未扣除
	MoralNotDeducted = int32(0)
	// MoralHasDeducted 节操已经被扣除
	MoralHasDeducted = int32(1)

	// RptOriginReply 评论
	RptOriginReply = int32(1)
	// RptOriginDM 弹幕
	RptOriginDM = int32(2)
	// RptOriginLetter 私信
	RptOriginLetter = int32(3)
	// RptOriginTag 标签
	RptOriginTag = int32(4)
	// RptOriginProfile 个人资料
	RptOriginProfile = int32(5)
	// RptOriginVideoup 投稿
	RptOriginVideoup = int32(6)
)

// ReportInfo ReportInfo
type ReportInfo struct {
	// TypeID     int64        `json:"type_id"` // rid
	ID         int64        `json:"id"`
	Oid        int64        `json:"oid"`
	Type       int32        `json:"type"`
	Title      string       `json:"title"`
	Count      int32        `json:"count"`
	MissionID  int64        `json:"mission_id"`
	Tid        int64        `json:"tid"`
	TName      string       `json:"tname"`
	TagState   int32        `json:"tag_state"`
	Mid        int64        `json:"mid"`
	Action     int32        `json:"action"`
	Rid        int64        `json:"rid"`
	Reason     int32        `json:"reason"`
	IsDelMoral int32        `json:"is_del_moral"`
	Score      int32        `json:"score"`
	State      int32        `json:"state"`
	RptMid     int64        `json:"rpt_mid"`
	RptIsUp    int32        `json:"rpt_is_up"`
	MidIsUp    int32        `json:"mid_is_up"`
	Examine    int32        `json:"examine"`
	Log        []*ReportLog `json:"log"`
	CTime      time.Time    `json:"ctime"`
	MTime      time.Time    `json:"mtime"`
}

// ReportLog ReportLog.
type ReportLog struct {
	ID         int64     `json:"id"`
	RptID      int64     `json:"rpt_id"`
	UserName   string    `json:"username"`
	Points     int32     `json:"points"`
	Oid        int64     `json:"oid"`
	Type       int32     `json:"type"`
	Mid        int64     `json:"mid"`
	Tid        int64     `json:"tid"`
	Rid        int64     `json:"rid"`
	Reason     string    `json:"reason"`
	HandleType int32     `json:"handle_type"`
	Notice     int32     `json:"is_notice"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
	Tag        *Tag      `json:"tag,omitempty"`
}

// ReportUser ReportUser.
type ReportUser struct {
	ID    int64     `json:"id"`
	RptID int64     `json:"rpt_id"`
	Mid   int64     `json:"rpt_mid"`
	Attr  int32     `json:"attr"`
	CTime time.Time `json:"ctime"`
	MTime time.Time `json:"mtime"`
}

// Report Report.
type Report struct {
	ID      int64     `json:"id"`
	Oid     int64     `json:"oid"`
	Type    int32     `json:"type"`
	Mid     int64     `json:"mid"`
	Tid     int64     `json:"tid"`
	Prid    int64     `json:"prid"`
	Rid     int64     `json:"rid"`
	Action  int32     `json:"action"`
	Reason  int32     `json:"reason"`
	Count   int32     `json:"count"`
	Content string    `json:"content"`
	Moral   int32     `json:"moral"`
	Score   int32     `json:"score"`
	State   int32     `json:"state"`
	CTime   time.Time `json:"ctime"`
	MTime   time.Time `json:"mtime"`
}

// ReportDetail ReportDetail
type ReportDetail struct {
	ID         int64     `json:"id"`
	PRptID     int64     `json:"parent_id"`
	Oid        int64     `json:"oid"`
	Type       int32     `json:"type"`
	Title      string    `json:"title"`
	Tid        int64     `json:"tid"`
	TName      string    `json:"tname"`
	TagState   int32     `json:"tag_state"`
	Mid        int64     `json:"mid"`
	RptMid     int64     `json:"rpt_mid"`
	Action     int32     `json:"action"`
	Rid        int64     `json:"rid"`
	Reason     int32     `json:"reason"`
	IsDelMoral int32     `json:"is_del_moral"`
	Score      int32     `json:"score"`
	State      int32     `json:"state"`
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
}
