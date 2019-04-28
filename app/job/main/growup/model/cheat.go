package model

import "go-common/library/time"

// Spy spy from http.
type Spy struct {
	ID        int64  `json:"id"`
	LogDate   string `json:"log_date"`
	TargetMID int64  `json:"target_mid"`
	TargetID  int64  `json:"target_id"`
	EventID   int64  `json:"event_id"`
	State     int    `json:"state"`
	Type      int    `json:"type"`
	Quantity  int    `json:"quantity"`
	IsDeleted int    `json:"is_del"`
}

// Cheating cheat
type Cheating struct {
	MID            int64
	Nickname       string
	AvID           int64
	PlayCount      int64
	Fans           int
	CheatFans      int
	CheatPlayCount int
	CheatCoin      int
	CheatFavorite  int
	UploadTime     time.Time
	SignedAt       time.Time
	TotalIncome    int
	AccountState   int
	Deducted       int
	IsDeleted      int
}
