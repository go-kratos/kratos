package model

import "go-common/library/time"

// UpSpy up spy.
type UpSpy struct {
	MID            int64     `json:"mid"`
	Nickname       string    `json:"nickname"`
	SignedAt       time.Time `json:"signed_at"`
	Fans           int       `json:"fans"`
	CheatFans      int       `json:"cheat_fans"`
	PlayCount      int       `json:"play_count"`
	CheatPlayCount int       `json:"cheat_play_count"`
	AccountState   int       `json:"account_state"`
}

// ArchiveSpy archive spy.
type ArchiveSpy struct {
	ArchiveID      int64     `json:"archive_id"`
	MID            int64     `json:"mid"`
	Nickname       string    `json:"nickname"`
	UploadTime     time.Time `json:"pub_time"`
	TotalIncome    int       `json:"total_income"`
	CheatFavorite  int       `json:"cheat_favorite"`
	CheatPlayCount int       `json:"cheat_play_count"`
	CheatCoin      int       `json:"cheat_coin"`
	Deducted       int       `json:"deducted"`
}

// CheatFans cheat fans.
type CheatFans struct {
	MID       int64     `json:"mid"`
	Nickname  string    `json:"nickname"`
	RealFans  int       `json:"real_fans"`
	CheatFans int       `json:"cheat_fans"`
	SignedAt  time.Time `json:"signed_at"`
	DeductAt  time.Time `json:"deduct_at"`
}

// CheatCount cheat count.
type CheatCount struct {
	Quantity int    `json:"quantity"`
	EventID  string `json:"event_name"`
}
