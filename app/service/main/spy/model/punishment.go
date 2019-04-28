package model

import (
	"go-common/library/time"
)

const (
	//PunishmentTypeBlock 封禁
	PunishmentTypeBlock = 1
)

// Punishment def.
type Punishment struct {
	ID     int64
	Mid    int64
	Type   int8   //封禁原因
	Reason string //惩罚原因
	CTime  time.Time
	MTime  time.Time
}
