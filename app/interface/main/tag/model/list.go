package model

import "go-common/library/time"

const (
	// ArcListIsDelBit .
	ArcListIsDelBit = uint(0)
	// ArcListIsAddBit .
	ArcListIsAddBit = uint(1)
)

// LimitArc .
type LimitArc struct {
	ID     int64     `json:"id"`
	Aid    int64     `json:"aid"`
	UpName string    `json:"author"`
	IsAdd  int8      `json:"is_add"`
	IsDel  int8      `json:"is_del"`
	CTime  time.Time `json:"-"`
	MTime  time.Time `json:"-"`
}
