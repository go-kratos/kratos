package model

import (
	"go-common/library/time"
)

//MoralLog is.
type MoralLog struct {
	Mid        int64     `json:"mid"`
	IP         int64     `json:"ip"`
	Operater   string    `json:"operater"`
	Origin     int8      `json:"origin"`
	Reason     string    `json:"reason"`
	Remark     string    `json:"remark"`
	Status     int8      `json:"status"`
	FromMoral  int64     `json:"from_moral"`
	ToMoral    int64     `json:"to_moral"`
	ModifyTime time.Time `json:"modify_time"`
}
