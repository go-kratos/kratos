package model

import (
	"go-common/library/time"
)

// CActivity creative activity
type CActivity struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name" form:"name" validate:"required"`
	Creator           string    `json:"creator"`
	SignedStart       time.Time `json:"signed_start" form:"signed_start" validate:"required"`
	SignedEnd         time.Time `json:"signed_end" form:"signed_end" validate:"required"`
	SignUp            int       `json:"sign_up" form:"sign_up" default:"0"` // 需要报名 0不需要,1需要
	SignUpStart       time.Time `json:"sign_up_start" form:"sign_up_start" validate:"required"`
	SignUpEnd         time.Time `json:"sign_up_end" form:"sign_up_end" validate:"required"`
	Object            int       `json:"object" form:"object" validate:"required"` // 1:uid, 2:avid
	UploadStart       time.Time `json:"upload_start" form:"upload_start" validate:"required"`
	UploadEnd         time.Time `json:"upload_end" form:"upload_end" validate:"required"`
	WinType           int       `json:"win_type" form:"win_type" validate:"required"`           // 1:达标型,2:排序型
	RequireItems      string    `json:"require_items" form:"require_items" validate:"required"` // 1:点赞,2:分享,3:播放,4:评论,5:弹幕, 多个用","分割
	RequireValue      int64     `json:"require_value" form:"require_value" validate:"required"`
	StatisticsStart   time.Time `json:"statistics_start" form:"statistics_start" validate:"required"`
	StatisticsEnd     time.Time `json:"statistics_end" form:"statistics_end" validate:"required"`
	BonusType         int       `json:"bonus_type" form:"bonus_type" validate:"required"`         // 1:平分,2:各得
	BonusMoney        []int64   `json:"bonus_money" form:"bonus_money,split" validate:"required"` // (多个","分割)
	BonusTime         time.Time `json:"bonus_time" form:"bonus_time" validate:"required"`
	ProgressFrequency int       `json:"progress_frequency" form:"progress_frequency" validate:"required"` // 进展更新频率 1:每天 2:每周
	UpdatePage        int       `json:"update_page" form:"update_page" default:"0"`                       // 更新活动页 0:否 1:是
	ProgressStart     time.Time `json:"progress_start" form:"progress_start" validate:"required"`
	ProgressEnd       time.Time `json:"progress_end" form:"progress_end" validate:"required"`
	ProgressSync      int       `json:"progress_sync" form:"progress_sync" default:"0"` // 进展同步 1共有,2已有,3共有/已有
	BonusQuery        int       `json:"bonus_query" form:"bonus_query" default:"0"`     // 开奖查询 0:否 1:是
	BonusQuerStart    time.Time `json:"bonus_query_start" form:"bonus_query_start" validate:"required"`
	BonusQueryEnd     time.Time `json:"bonus_query_end" form:"bonus_query_end" validate:"required"`
}

// ActivityBonus activity bonus
type ActivityBonus struct {
	BonusMoney int64
	Rank       int64
}

// UpActivity up activity
type UpActivity struct {
	MID         int64
	ActivityID  int64
	AIDs        []int64
	AIDNum      int64
	ItemVal     int64
	Rank        int
	Bonus       int64
	State       int
	Nickname    string
	SuccessTime time.Time
}

// AvItem av item
type AvItem struct {
	AvID  int64
	Value int64
}

// AvUpload av upload time
type AvUpload struct {
	ID         int64
	MID        int64
	AvID       int64
	UploadTime time.Time
}
