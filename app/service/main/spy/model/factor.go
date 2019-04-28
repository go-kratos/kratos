package model

import (
	"go-common/library/time"
)

// Factor def.
type Factor struct {
	ID        int64
	NickName  string //风险因子名字
	ServiceID int64
	EventID   int64
	GroupID   int64
	RiskLevel int8    //风险等级：1-9
	FactorVal float64 //风险得分
	CTime     time.Time
	MTime     time.Time
}
