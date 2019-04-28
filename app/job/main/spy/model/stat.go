package model

import "time"

const (
	// IncreaseStat increase spy stat.
	IncreaseStat int8 = 1
	// ResetStat reset spy stat.
	ResetStat int8 = 2

	// WaiteCheck waite for check.
	WaiteCheck = 0
)

// type def.
const (
	ArchiveType  int8 = 1
	ActivityType int8 = 2
)

// Statistics def.
type Statistics struct {
	ID        int64     `json:"id"`
	TargetMid int64     `json:"target_mid"`
	TargetID  int64     `json:"target_id"`
	EventID   int64     `json:"event_id"`
	State     int8      `json:"state"`
	Type      int8      `json:"type"`
	Quantity  int64     `json:"quantity"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// SpyLog def.
type SpyLog struct {
	ID        int64     `json:"id"`
	TargetMid int64     `json:"target_mid"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}
