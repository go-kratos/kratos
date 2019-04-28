package model

import (
	"time"
)

// Factor def.
type Factor struct {
	ID         int64     `json:"id"`
	NickName   string    `json:"nick_name"`   //风险因子名字
	ServiceID  int64     `json:"service_id"`  //服务ID
	EventID    int64     `json:"event_id"`    //事件ID
	GroupID    int64     `json:"group_id"`    //分组ID
	RiskLevel  int8      `json:"risk_level"`  //风险等级：1-9
	FactorVal  float32   `json:"factor_val"`  //因子值
	CTime      time.Time `json:"ctime"`       //创建时间
	MTime      time.Time `json:"mtime"`       //修改时间
	CategoryID int8      `json:"category_id"` //CategoryID
}

// Factors etc.
type Factors struct {
	ID        int64   `json:"id"`
	GroupID   int64   `json:"group_id"`   //分组ID
	GroupName string  `json:"group_name"` //分组名称
	NickName  string  `json:"nick_name"`  //风险因子名字
	FactorVal float32 `json:"factor_val"` //因子值
}
