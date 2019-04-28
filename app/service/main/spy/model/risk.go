package model

import "time"

// Tel Expire
const (
	TelExpireMonth = 3
	QcloudType     = 1
)

// TelRiskLevel def.
type TelRiskLevel struct {
	ID     int64     `json:"id"`
	Mid    int64     `json:"mid"`
	Level  int8      `json:"level"`
	Origin int8      `json:"origin"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
}
