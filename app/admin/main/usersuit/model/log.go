package model

import (
	xtime "go-common/library/time"
)

// OperationLog operation log .
type OperationLog struct {
	ID     int64      `json:"id"`
	OID    int64      `json:"oid"`
	Action string     `json:"action"`
	CTime  xtime.Time `json:"ctime"`
	MTime  xtime.Time `json:"mtime"`
}
