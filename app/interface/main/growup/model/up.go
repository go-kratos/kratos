package model

import (
	"time"

	xtime "go-common/library/time"
)

// UpInfo is users of growup/video/column who applied for.
type UpInfo struct {
	ID                   int64      `json:"id"`
	MID                  int64      `json:"mid"`
	Nickname             string     `json:"nickname"`
	AccountType          int        `json:"account_type"`
	OriginalArchiveCount int        `json:"original_archive_count"`
	MainCategory         int        `json:"category_id"`
	Bgms                 int        `json:"bgms"`
	Fans                 int        `json:"fans"`
	TotalPlayCount       int64      `json:"total_play_count"`
	AccountState         int        `json:"account_state"`
	SignType             int        `json:"sign_type,omitempty"`
	Reason               string     `json:"reason"`
	ApplyAt              xtime.Time `json:"apply_at"`
	SignedAt             xtime.Time `json:"signed_at"`
	RejectAt             xtime.Time `json:"reject_at"`
	ForbidAt             xtime.Time `json:"forbid_at"`
	QuitAt               xtime.Time `json:"quit_at"`
	DismissAt            xtime.Time `json:"dismiss_at"`
	ExpiredIn            xtime.Time `json:"expired_in"`
	IsDeleted            int        `json:"-"`
}

// UpStatus is user status of growup plan
type UpStatus struct {
	Status  []*BusinessStatus `json:"status"`
	Blocked bool              `json:"blocked"`
}

// BusinessStatus type: 1.视频 2.专栏 3.素材
type BusinessStatus struct {
	IsWhite      bool       `json:"in_white_list"`
	AccountState int        `json:"account_state"`
	AccountType  int        `json:"account_type"`
	Type         int        `json:"type"`
	Reason       string     `json:"reason"`
	ShowPanel    bool       `json:"show_panel"`
	ExpiredIn    xtime.Time `json:"expired_in"`
	QuitAt       time.Time  `json:"-"`
	CTime        time.Time  `json:"-"`
}

// CreditRecord credit record
type CreditRecord struct {
	ID        int64      `json:"id"`
	MID       int64      `json:"mid"`
	OperateAt xtime.Time `json:"operate_at"`
	Operator  string     `json:"operator"`
	Reason    int        `json:"reason"`
	Deducted  int        `json:"deducted"`
	Remaining int        `json:"remaining"`
	IsDeleted int        `json:"recovered"`
	CTime     xtime.Time `json:"ctime"`
	MTime     xtime.Time `json:"mtime"`
}
