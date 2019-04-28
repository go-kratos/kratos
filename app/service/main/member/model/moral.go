package model

const (
	// DMReasonType is.
	DMReasonType = 1
	// ReplyReasonType is.
	ReplyReasonType = 2
	// TagReasonType is.
	TagReasonType = 3
	// ElecReasonType is.
	ElecReasonType = 4
	// AccountReasonType is.
	AccountReasonType = 5
	// SysReasonType is.
	SysReasonType = 6

	// RevocableMoralStatus is.
	RevocableMoralStatus = 0
	// RevokedMoralStatus is.
	RevokedMoralStatus = 1
	// IrrevocableMoralStatus is.
	IrrevocableMoralStatus = 2

	// ReportRewardType is.
	ReportRewardType = 1
	// PunishmentType is.
	PunishmentType = 2
	// CancelRewardType is.
	CancelRewardType = 3
	// CancelPunishType is.
	CancelPunishType = 4
	// ManualRecoveryType is.
	ManualRecoveryType = 5
	// ManualChangeType is.
	ManualChangeType = 6
)

//ArgUpdateMorals argUpdateMorals.
type ArgUpdateMorals struct {
	Mids       []int64 `form:"mids,split" validate:"required"`
	Delta      int64   `form:"delta" validate:"required"`
	Origin     int64   `form:"origin" validate:"required"`
	Reason     string  `form:"reason" validate:"required"`
	ReasonType int64   `form:"reason_type"`
	Operator   string  `form:"operator" validate:"required"`
	Remark     string  `form:"remark" validate:"required"`
	Status     int64   `form:"status"`
	IsNotify   bool    `form:"is_notify"`
	IP         string  `form:"ip"`
}

//ArgUpdateMoral argUpdateMoral.
type ArgUpdateMoral struct {
	Mid        int64  `form:"mid" validate:"required"`
	Delta      int64  `form:"delta" validate:"required"`
	Origin     int64  `form:"origin" validate:"required"`
	Reason     string `form:"reason" validate:"required"`
	ReasonType int64  `form:"reason_type"`
	Operator   string `form:"operator" validate:"required"`
	Remark     string `form:"remark" validate:"required"`
	Status     int64  `form:"status"`
	IsNotify   bool   `form:"is_notify"`
	IP         string `form:"ip"`
}

//ArgUndo argUndo.
type ArgUndo struct {
	LogID    string `form:"log_id" validate:"required"`
	Remark   string `form:"remark" validate:"required"`
	Operator string `form:"operator" validate:"required"`
}

//ReasonType reasonType
type ReasonType struct {
	Name       string
	NotifyType string
}

//OriginType originType
type OriginType struct {
	Name       string
	NeedReason bool
}

//Notice notice
type Notice struct {
	Title      string
	Message    string
	NoticeType string
}

var (
	// ReasonTypes ...
	ReasonTypes = map[int64]*ReasonType{
		DMReasonType:      {"弹幕", "2_1_4"},
		ReplyReasonType:   {"评论", "2_1_3"},
		TagReasonType:     {"TAG", ""},
		ElecReasonType:    {"电波", ""},
		AccountReasonType: {"账号", ""},
		SysReasonType:     {"管理系统", ""},
	}
	// OriginTypes ...
	OriginTypes = map[int64]*OriginType{
		ReportRewardType:   {"举报奖励", true},
		PunishmentType:     {"违规惩罚", true},
		CancelRewardType:   {"撤销奖励", true},
		CancelPunishType:   {"撤销惩罚", true},
		ManualRecoveryType: {"自动恢复", true},
		ManualChangeType:   {"手动修改", false},
	}
	// Less6000Notice is.
	Less6000Notice = &Notice{Title: "你的节操值已低于60", Message: "抱歉，你的节操值已低于60，社交类功能将不能正常使用，更多加减明细请查看 #{节操记录}{\"https://account.bilibili.com/site/record?type=moral\"}", NoticeType: "2_1_5"}
	// Less3000Notice is.
	Less3000Notice = &Notice{Title: "你的节操值已低于30", Message: "抱歉，你的节操值已低于30，社交类功能将不能正常使用，更多加减明细请查看 #{节操记录}{\"https://account.bilibili.com/site/record?type=moral\"}", NoticeType: "2_1_6"}
	// Greater6000Notice is .
	Greater6000Notice = &Notice{Title: "你的节操值已恢复至60以上", Message: "恭喜，你的节操值已恢复至60以上，所有功能将回复正常使用，更多加减明细请查看 #{节操记录}{\"https://account.bilibili.com/site/record?type=moral\"}", NoticeType: "2_1_7"}
	// PunishmentNotice is.
	PunishmentNotice = &Notice{Title: "你被举报处理扣除了%s节操值", Message: "由于发布了违规内容，你被举报处理扣除了%s节操值，具体原因请看 #{节操记录}{\"https://account.bilibili.com/site/record?type=moral\"}"}
	// SysPunishmentNotice is.
	SysPunishmentNotice = &Notice{Title: "你被举报处理扣除了%s节操值", Message: "由于发布了违规内容，你被系统处理扣除了%s节操值，具体原因请看 #{节操记录}{\"https://account.bilibili.com/site/record?type=moral\"}"}
	// RewardNotice is.
	RewardNotice = &Notice{Title: "你举报的%s已被处理", Message: "您举报的%s已被管理员处理，获得了%s节操值奖励，具体详情请看 #{节操记录}{\"https://account.bilibili.com/site/record?type=moral\"}"}
)
