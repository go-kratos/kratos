package model

import "time"

// config properties.
const (
	LimitBlockCount = "limitBlockCount"
	LessBlockScore  = "lessBlockScore"
	AutoBlock       = "autoBlock"
	AutoBlockOpen   = 1
)

const (
	//BlockReasonSize block reason size
	BlockReasonSize = 4
	//BlockLockKey cycle block
	BlockLockKey = "cycleblock"
	//VipEnableStatus enable status
	VipEnableStatus int32 = 1
	//VipNonType non vip
	VipNonType int32 = 0
	// ReliveCheckTimes relive check times
	ReliveCheckTimes = 1
	// DoubleCheckRemake double check remake
	DoubleCheckRemake = "导入二次验证,恢复行为得分"
	// RetryTimes retry times
	RetryTimes = 3
	// SpyInitScore spy init score
	SpyInitScore = 100
)

// Config def.
type Config struct {
	ID       int64
	Property string
	Name     string
	Val      string
	Ctime    time.Time
}
