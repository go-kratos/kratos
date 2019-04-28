package model

import (
	"go-common/library/time"
)

// UserEventHistory history of user event def.
type UserEventHistory struct {
	ID         int64
	Mid        int64
	EventID    int64
	Score      int8
	BaseScore  int8
	EventScore int8
	Remark     string
	Reason     string // 事件原因 == eventmsg.Effect
	FactorVal  float64
	CTime      time.Time
}
