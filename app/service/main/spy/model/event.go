package model

import (
	"go-common/library/time"
)

// Event def.
type Event struct {
	ID        int64
	Name      string // 事件标示
	NickName  string // 事件可读昵称
	ServiceID int64
	Status    int8 // 0:删除 1:未删除
	CTime     time.Time
	MTime     time.Time
}
