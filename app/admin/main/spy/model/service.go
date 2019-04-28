package model

import (
	"time"
)

// Service def.
type Service struct {
	ID       int64
	Name     string //服务标识
	NickName string //服务可读昵称
	Status   int8
	CTime    time.Time
	MTime    time.Time
}
