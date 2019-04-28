package model

import (
	"go-common/library/time"
)

// CActivity creative activity
type CActivity struct {
	ID              int64        `json:"id"`
	Name            string       `json:"name"`
	SignedStart     time.Time    `json:"-"`
	SignedEnd       time.Time    `json:"-"`
	SignUpStart     time.Time    `json:"-"`
	SignUpEnd       time.Time    `json:"-"`
	SignUp          int          `json:"sign_up"`
	SignUpState     int          `json:"sign_up_state"` // 0可以报名,1已报名,2已获奖,3不能报名
	WinType         int          `json:"win_type"`
	ProgressStart   time.Time    `json:"-"`
	ProgressEnd     time.Time    `json:"-"`
	ProgressState   int          `json:"progress_state"` // 展示状态 0不展示 1展示
	ProgressSync    int          `json:"progress_sync"`
	UpdatePage      int          `json:"-"`
	BonusQuery      int          `json:"bonus_query"`
	BonusQueryStart time.Time    `json:"-"`
	BonusQueryEnd   time.Time    `json:"-"`
	Background      string       `json:"background"`
	WinDesc         string       `json:"win_desc"`
	UnwinDesc       string       `json:"unwin_desc"`
	Details         string       `json:"details"`
	Enrollment      int          `json:"enrollment"`
	WinNum          int          `json:"win_num"`
	Ranking         []*ActUpInfo `json:"ranking"`
}

// UpBonus up bonus
type UpBonus struct {
	MID        int64
	ActivityID int64
	Nickname   string
	Rank       int
	State      int
	SignUpTime time.Time
}
