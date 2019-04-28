package model

import (
	"database/sql"

	"go-common/library/time"
)

const (
	// ArcTagStateNormal tag archive
	ArcTagStateNormal = 0

	// ArcTagLogAdd tag Log
	ArcTagLogAdd = int8(0)
	// ArcTagLogDel tag Log
	ArcTagLogDel = int8(1)
	// ArcTagLogOpen tag Log
	ArcTagLogOpen = int8(0)
	// ArcTagLogClose tag Log
	ArcTagLogClose = int8(1)

	// ArcTagAdd  operation type
	ArcTagAdd = int8(1)
	// ArcTagDel  operation type
	ArcTagDel = int8(2)
	// ArcTagLike  operation type
	ArcTagLike = int8(3)
	// ArcTagHate  operation type
	ArcTagHate = int8(4)
	// ArcTagRpt  operation type
	ArcTagRpt = int8(7)

	// LogAddReport log report
	LogAddReport = int8(8)
	// LogDelReport  archive log report
	LogDelReport = int8(9)

	// ArcTagOpRoleUp arc tag operation role
	ArcTagOpRoleUp = 0
	// ArcTagOpRoleUser  arc tag operation role
	ArcTagOpRoleUser = 1
	// ArcTagOpRoleAdmin  arc tag operation role
	ArcTagOpRoleAdmin = 2

	// ArcTagLockBit  arctaglock
	ArcTagLockBit = 0
	// ArcTagRptPass  arctaglock
	ArcTagRptPass = 1

	// LimitArcDelbit arc bit
	LimitArcDelbit = uint(0)
	// LimitArcAddbit  arc bit
	LimitArcAddbit = uint(1)

	// ActionStateOpen tag action
	ActionStateOpen = int8(0)
	// ActionStateClose  tag action
	ActionStateClose = int8(1)
	// ActionLikeIncr  tag action
	ActionLikeIncr = int8(1)
	// ActionLikedecr tag action
	ActionLikedecr = int8(-1)
	// ActionHateIncr  tag action
	ActionHateIncr = int8(1)
	// ActionHatedecr  tag action
	ActionHatedecr = int8(-1)
)

// ArcTag  archive_tag
type ArcTag struct {
	ID        int64     `json:"-"`
	Aid       int64     `json:"aid"`
	Mid       int64     `json:"mid"`
	Tid       int64     `json:"tag_id"`
	Likes     int64     `json:"likes"`
	Hates     int64     `json:"hates"`
	Attribute int8      `json:"attribute"`
	Role      int8      `json:"-"`
	State     int8      `json:"state"`
	CTime     time.Time `json:"ctime"`
	MTime     time.Time `json:"-"`
}

// IsLock .
func (at *ArcTag) IsLock() bool {
	return at.Attribute&0x1 == 1
}

// ArcTagLog archive tag log
type ArcTagLog struct {
	Lid    int64  `json:"id"`
	Aid    int64  `json:"aid"`
	Tid    int64  `json:"tag_id"`
	Tname  string `json:"tag_name"`
	Mid    int64  `json:"mid"`
	Face   string `json:"face"`
	UName  string `json:"uname"`
	Role   int8   `json:"role"`
	Action int8   `json:"action"`
	Remark string `json:"-"`
	Lstate int8   `json:"-"`
	// report info
	Rid      sql.NullInt64 `json:"-"`
	State    sql.NullInt64 `json:"-"` // report state
	IsDeal   int8          `json:"is_deal"`
	IsReport int8          `json:"is_report"`
	CTime    time.Time     `json:"ctime"`
	MTime    time.Time     `json:"-"`
}

// ArcTagAction archive_tag_action
type ArcTagAction struct {
	ID     int64
	Aid    int64
	Tid    int64
	Mid    int64
	Action int8
	State  int8
}
