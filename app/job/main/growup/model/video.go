package model

import (
	"go-common/library/time"
)

// CategoryInfo query from up_category_info,used to update up_info_xxx
type CategoryInfo struct {
	MID        int64
	Nickname   string
	CategoryID int
}

// UpBaseInfo query from up_base_statistics, used to update up_info_video
type UpBaseInfo struct {
	MID                  int64
	Fans                 int
	Avs                  int
	OriginalArchiveCount int
	TotalPlayCount       int
}

// UpInfoVideo up_info_video
type UpInfoVideo struct {
	MID            int64
	Nickname       string
	AccountType    int
	AccountState   int
	SignedAt       time.Time
	Fans           int64
	TotalPlayCount int64
	CreditScore    int64
	IsDeleted      int
}
