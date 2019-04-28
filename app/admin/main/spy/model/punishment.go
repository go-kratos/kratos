package model

import (
	"go-common/library/time"
)

// Punishment def.
type Punishment struct {
	ID     int64
	Mid    int64
	Type   int8   //惩罚类型 1.自动封禁
	Reason string //惩罚原因
	CTime  time.Time
}
