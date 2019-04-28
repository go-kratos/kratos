package task

import (
	"time"
)

//AssignConfig 指派配置
type AssignConfig struct {
	ID          int64
	Pool        int8
	MinDuration int64
	MaxDuration int64
	MIDs        map[int64]struct{}
	TIDs        map[int16]struct{}
	UIDs        []int64
	AdminID     int64
	State       int8
	STime       time.Time
	ETime       time.Time
	Index       int
}
