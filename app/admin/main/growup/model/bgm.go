package model

import (
	"go-common/library/time"
)

// BGM bgm
type BGM struct {
	ID     int64
	MID    int64
	CID    int64
	AID    int64
	SID    int64
	JoinAt time.Time
}

// Business business
type Business struct {
	ID           int64 `json:"id"`
	MID          int64 `json:"mid"`
	AccountState int   `json:"account_state"`
	Type         int   `json:"archive_type"`
	CTime        int64 `json:"-"`
	MTime        int64 `json:"-"`
}
