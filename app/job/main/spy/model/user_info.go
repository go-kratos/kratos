package model

import (
	"go-common/library/time"
)

const (
	// StateNormal UserInfo 正常状态
	StateNormal = 0
	// StateBlock UserInfo 封禁状态
	StateBlock = 1
)

// UserInfo def.
type UserInfo struct {
	ID         int64     `json:"id"`
	Mid        int64     `json:"mid"`
	Score      int8      `json:"score"`       //真实得分
	BaseScore  int8      `json:"base_score"`  //基础信息得分
	EventScore int8      `json:"event_score"` //事件得分
	State      int8      `json:"state"`       //状态 : 正常/封禁
	CTime      time.Time `json:"ctime"`
	MTime      time.Time `json:"mtime"`
}
