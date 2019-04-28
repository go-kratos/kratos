package model

import (
	"time"
)

// UserInfo def.
type UserInfo struct {
	ID          int64     `json:"id"`
	Mid         int64     `json:"mid"`
	State       int8      `json:"state"`        //状态  0.正常  1.封禁
	Score       int8      `json:"score"`        //真实分值
	BaseScore   int8      `json:"base_score"`   //基础信息得分
	EventScore  int8      `json:"event_score"`  //事件得分
	ReliveTimes int8      `json:"relive_times"` //累计重绘次数
	Mtime       time.Time `json:"mtime"`
}

// UserInfoDto dto.
type UserInfoDto struct {
	ID          int64  `json:"id"`
	Mid         int64  `json:"mid"`
	Name        string `json:"name"`         //昵称
	State       int8   `json:"state"`        //状态  0.正常  1.封禁
	Score       int8   `json:"score"`        //真实分值
	BaseScore   int8   `json:"base_score"`   //基础信息得分
	EventScore  int8   `json:"event_score"`  //事件得分
	ReliveTimes int8   `json:"relive_times"` //累计重绘次数
	Mtime       int64  `json:"mtime"`        //更新时间
}
