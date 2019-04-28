package model

import (
	"go-common/library/time"
)

const (
	// SubStateNormal .
	SubStateNormal = 0 // SubState default
	// SubTagLoad arc sub operation type
	SubTagLoad = int8(0)
	// SubTagAdd .
	SubTagAdd = int8(1)
	// SubTagDel .
	SubTagDel = int8(2)

	// SortOrderDESC sort order desc.
	SortOrderDESC = int(-1)
	// SortOrderASC order asc.
	SortOrderASC = int(1)

	// SubTagMaxNum sub tag total number .
	SubTagMaxNum = 400
)

// Attention .
type Attention struct {
	ID    int64     `json:"id"`
	Tid   int64     `json:"tag_id"`
	Mid   int64     `json:"mid"`
	State int8      `json:"state"`
	CTime time.Time `json:"ctime"`
	MTime time.Time `json:"-"`
}

// SubArcs .
type SubArcs struct {
	Tag  *Tag    `json:"tag"`
	Aids []int64 `json:"aids"`
}
