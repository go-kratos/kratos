package model

import (
	xtime "go-common/library/time"
)

// ArchiveStatStream is table archive_stat
type ArchiveStatStream struct {
	ID        int64      `json:"id"`
	Business  int        `json:"business"`
	StatType  int        `json:"stat_type"`
	TypeID    int        `json:"typeid"`
	GroupID   int        `json:"group_id"`
	UID       int        `json:"uid"`
	StatTime  xtime.Time `json:"stat_time"`
	StatValue int64      `json:"stat_value"`
	Ctime     xtime.Time `json:"ctime"`
	Mtime     xtime.Time `json:"mtime"`
	State     int        `json:"state"`
}

const (
	// business字段枚举值

	// Recheck123 is 稿件123回查移入移出
	Recheck123 = 1

	// stat_type字段枚举值

	// FIRST_RECHECK_IN is 一查移入
	FIRST_RECHECK_IN = 1
	// FIRST_RECHECK_OUT is 一查移出
	FIRST_RECHECK_OUT = 2
	// SECOND_RECHECK_IN is 二查移入
	SECOND_RECHECK_IN = 3
	// SECOND_RECHECK_OUT is 二查移出
	SECOND_RECHECK_OUT = 4
	// THIRD_RECHECK_IN is 三查移入
	THIRD_RECHECK_IN = 5
	// THIRD_RECHECK_OUT is 三查移出
	THIRD_RECHECK_OUT = 6
)
