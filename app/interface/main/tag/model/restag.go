package model

import (
	"go-common/library/time"
)

const (
	// ResTagStateNormal .
	ResTagStateNormal = int32(0)
	// ResTagStateDel .
	ResTagStateDel = int32(1)
	// ResTagStateHide .
	ResTagStateHide = int32(2)
	// ResTagStateRegion .
	ResTagStateRegion = int32(3)

	// ResTagAdd .
	ResTagAdd = int8(0)
	// ResTagDel .
	ResTagDel = int8(1)
)

// ResTag .
type ResTag struct {
	ID    int64     `json:"id"`
	Oid   int64     `json:"oid"`
	Tid   int64     `json:"tag_id"`
	Type  int8      `json:"type"`
	Mid   int64     `json:"mid"`
	Likes int64     `json:"likes"`
	Hates int64     `json:"hates"`
	Attr  int8      `json:"attribute"`
	Role  int8      `json:"-"`
	State int8      `json:"state"`
	CTime time.Time `json:"ctime"`
	MTime time.Time `json:"-"`
}

// ResTagLog .
type ResTagLog struct {
	ID     int64     `json:"id"`
	Oid    int64     `json:"oid"`
	Tid    int64     `json:"tag_id"`
	Type   int8      `json:"type"`
	Tname  string    `json:"tag_name"`
	Mid    int64     `json:"mid"`
	Role   int8      `json:"role"`
	Reason string    `json:"reason"`
	Action int8      `json:"action"`
	State  int8      `json:"state"`
	CTime  time.Time `json:"ctime"`
	MTime  time.Time `json:"-"`
}
