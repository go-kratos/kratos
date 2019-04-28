package model

import (
	"go-common/library/time"
)

// UpInfo upinfo of video
type UpInfo struct {
	ID                   int64     `json:"id" gorm:"id"`
	MID                  int64     `json:"mid" gorm:"column:mid"`
	Nickname             string    `json:"nickname" gorm:"column:nickname"`
	AccountType          int       `json:"account_type" gorm:"column:account_type"`
	OriginalArchiveCount int       `json:"original_archive_count" gorm:"column:original_archive_count"`
	ArticleCount         int       `json:"article_count" gorm:"column:article_count"`
	Avs                  int       `json:"avs" gorm:"column:avs"`
	BgmPlayCount         int       `json:"bgm_play_count"`
	BgmApplyCount        int       `json:"bgm_apply_count"`
	TotalPlayCount       int       `json:"total_play_count" gorm:"column:total_play_count"`
	TotalViewCount       int       `json:"total_view_count" gorm:"column:total_view_count"`
	MainCategory         int       `json:"category_id" gorm:"column:category_id"`
	Fans                 int       `json:"fans" gorm:"column:fans"`
	BGMs                 int       `json:"bgms"`
	AccountState         int       `json:"account_state" gorm:"column:account_state"`
	SignType             int       `json:"sign_type,omitempty" gorm:"column:sign_type"`
	Reason               string    `json:"reason" gorm:"column:reason"`
	ApplyAt              time.Time `json:"apply_at" gorm:"column:apply_at"`
	SignedAt             time.Time `json:"signed_at" gorm:"column:signed_at"`
	RejectAt             time.Time `json:"reject_at" gorm:"column:reject_at"`
	ForbidAt             time.Time `json:"forbid_at" gorm:"column:forbid_at"`
	QuitAt               time.Time `json:"quit_at,omitempty" gorm:"column:quit_at"`
	DismissAt            time.Time `json:"dismiss_at" gorm:"column:dismiss_at"`
	ExpiredIn            time.Time `json:"expired_in,omitempty" gorm:"column:expired_in"`
	CTime                time.Time `json:"ctime" gorm:"column:ctime"`
	MTime                time.Time `json:"mtime" gorm:"column:mtime"`
	CreditScore          int       `json:"credit_score"`
	IsDeleted            int       `json:"-" gorm:"column:is_deleted"`
	SignedType           []int     `json:"signed_type"`
	OtherType            []int     `json:"other_type"`
}

// Blocked up in blacklist
type Blocked struct {
	ID                   int64     `json:"id" gorm:"id"`
	MID                  int64     `json:"mid" gorm:"column:mid"`
	Nickname             string    `json:"nickname" gorm:"column:nickname"`
	OriginalArchiveCount int       `json:"original_archive_count" gorm:"column:original_archive_count"`
	MainCategory         int       `json:"category_id" gorm:"column:category_id"`
	Fans                 int       `json:"fans" gorm:"column:fans"`
	ApplyAt              time.Time `json:"apply_at" gorm:"column:apply_at"`
	CTime                time.Time `json:"ctime" gorm:"column:ctime"`
	MTime                time.Time `json:"mtime" gorm:"column:mtime"`
	IsDeleted            int       `json:"-" gorm:"column:is_deleted"`
}

// SimpleUpInfo include mid and signedAt for up-allowance-data compute
type SimpleUpInfo struct {
	MID      int64     `json:"mid"`
	SignedAt time.Time `json:"signed_at"`
}

// CreditRecord credit deducted and recover record
type CreditRecord struct {
	ID        int64     `json:"id"`
	MID       int64     `json:"mid"`
	OperateAt time.Time `json:"operate_at"`
	Operator  string    `json:"operator"`
	Reason    int       `json:"reason"`
	Deducted  int       `json:"deducted"`
	Remaining int       `json:"remaining"`
	IsDeleted int       `json:"recovered"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"mtime"`
}
